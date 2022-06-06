package bybitl

import (
	goContext "context"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/account"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	"gitlab.com/alphaticks/xchanger/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"reflect"
	"strconv"
	"time"
)

type AccountReconcile struct {
	extypes.BaseReconcile
	account          *models.Account
	executor         *actor.PID
	logger           *log.Logger
	securities       map[uint64]*registry.Security
	symbToSecs       map[string]*registry.Security
	txs              *mongo.Collection
	registry         registry.PublicRegistryClient
	positions        map[uint64]*account.Position
	lastDepositTs    uint64
	lastWithdrawalTs uint64
	lastFundingTs    uint64
	lastTradeTs      map[uint64]uint64
}

func NewAccountReconcileProducer(account *models.Account, registry registry.PublicRegistryClient, txs *mongo.Collection) actor.Producer {
	return func() actor.Actor {
		return NewAccountReconcile(account, registry, txs)
	}
}

func NewAccountReconcile(account *models.Account, registry registry.PublicRegistryClient, txs *mongo.Collection) actor.Actor {
	return &AccountReconcile{
		account:  account,
		txs:      txs,
		registry: registry,
	}
}

func (state *AccountReconcile) GetLogger() *log.Logger {
	return state.logger
}

func (state *AccountReconcile) GetTransactions() *mongo.Collection {
	return state.txs
}

func (state *AccountReconcile) Receive(context actor.Context) {
	extypes.ReconcileReceive(state, context)
}

func (state *AccountReconcile) Initialize(context actor.Context) error {
	// When initialize is done, the account must be aware of all the settings / assets / portfolio
	// so as to be able to answer to FIX messages
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))
	state.executor = actor.NewPID(context.ActorSystem().Address(), "executor/exchanges/"+constants.BYBITL.Name+"_executor")
	// Request securities
	res, err := state.registry.Securities(goContext.Background(), &registry.SecuritiesRequest{
		Filter: &registry.SecurityFilter{
			ExchangeId: []uint32{constants.BYBITL.ID},
		},
	})
	if err != nil {
		return fmt.Errorf("error fetching historical securities: %v", err)
	}

	state.symbToSecs = make(map[string]*registry.Security)
	securityMap := make(map[uint64]*registry.Security)
	for _, sec := range res.Securities {
		securityMap[sec.SecurityId] = sec
		state.symbToSecs[sec.Symbol] = sec
	}
	state.securities = securityMap
	state.lastTradeTs = make(map[uint64]uint64)

	// Start reconciliation
	state.positions = make(map[uint64]*account.Position)
	for _, sec := range state.securities {
		if sec.SecurityType == "CRPERP" {
			state.positions[sec.SecurityId] = account.NewPosition(
				sec.IsInverse, 1e8, 1e8, 1e8, 1, 0, 0)
		}
	}
	// First, calculate current positions from historical
	cur, err := state.txs.Find(goContext.Background(), bson.D{
		{Key: "account", Value: state.account.Name},
	}, options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}))
	if err != nil {
		return fmt.Errorf("error reconcile trades: %v", err)
	}
	balances := make(map[uint32]float64)
	for cur.Next(goContext.Background()) {
		var tx extypes.Transaction
		if err := cur.Decode(&tx); err != nil {
			return fmt.Errorf("error decoding transaction: %v", err)
		}
		switch tx.Type {
		case "TRADE":
			if tx.Fill.SecurityID != "" {
				secID, _ := strconv.ParseUint(tx.Fill.SecurityID, 10, 64)
				sec, ok := state.securities[secID]
				if ok && sec.SecurityType == "CRPERP" {
					if tx.Fill.Quantity < 0 {
						state.positions[secID].Sell(tx.Fill.Price, -tx.Fill.Quantity, false)
					} else {
						state.positions[secID].Buy(tx.Fill.Price, tx.Fill.Quantity, false)
					}
					state.lastTradeTs[sec.SecurityId] = uint64(tx.Time.UnixNano() / 1000000)
				}
			}
		}
		for _, m := range tx.Movements {
			balances[m.AssetID] += m.Quantity
		}
	}

	for k, b := range balances {
		a, _ := constants.GetAssetByID(k)
		if a != nil {
			fmt.Println(a.Symbol, b)
		}
	}

	for k, pos := range state.positions {
		if ppos := pos.GetPosition(); ppos != nil {
			fmt.Println(state.securities[k].Symbol, ppos.Quantity)
		}
	}

	sres := state.txs.FindOne(goContext.Background(), bson.D{
		{Key: "account", Value: state.account.Name},
		{Key: "type", Value: "FUNDING"},
	}, options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}}))
	if sres.Err() != nil {
		if sres.Err() != mongo.ErrNoDocuments {
			return fmt.Errorf("error getting last funding: %v", err)
		}
	} else {
		var tx extypes.Transaction
		if err := sres.Decode(&tx); err != nil {
			return fmt.Errorf("error decoding transaction: %v", err)
		}
		state.lastFundingTs = uint64(tx.Time.UnixNano() / 1000000)
	}

	sres = state.txs.FindOne(goContext.Background(), bson.D{
		{Key: "account", Value: state.account.Name},
		{Key: "type", Value: "DEPOSIT"},
	}, options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}}))
	if sres.Err() != nil {
		if sres.Err() != mongo.ErrNoDocuments {
			return fmt.Errorf("error getting last deposit: %v", err)
		}
	} else {
		var tx extypes.Transaction
		if err := sres.Decode(&tx); err != nil {
			return fmt.Errorf("error decoding transaction: %v", err)
		}
		state.lastDepositTs = uint64(tx.Time.UnixNano() / 1000000)
	}

	sres = state.txs.FindOne(goContext.Background(), bson.D{
		{Key: "account", Value: state.account.Name},
		{Key: "type", Value: "WITHDRAWAL"},
	}, options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}}))
	if sres.Err() != nil {
		if sres.Err() != mongo.ErrNoDocuments {
			return fmt.Errorf("error getting last withdrawal: %v", err)
		}
	} else {
		var tx extypes.Transaction
		if err := sres.Decode(&tx); err != nil {
			return fmt.Errorf("error decoding transaction: %v", err)
		}
		state.lastWithdrawalTs = uint64(tx.Time.UnixNano() / 1000000)
	}

	if err := state.reconcileTrades(context); err != nil {
		return fmt.Errorf("error reconcile trade: %v", err)
	}
	if err := state.reconcileMovements(context); err != nil {
		return fmt.Errorf("error reconcile movements: %v", err)
	}

	return nil
}

// TODO
func (state *AccountReconcile) Clean(context actor.Context) error {

	return nil
}

func (state *AccountReconcile) OnAccountMovementRequest(context actor.Context) error {
	return nil
}

func (state *AccountReconcile) reconcileTrades(context actor.Context) error {
	for _, sec := range state.securities {
		fmt.Println("REQ", sec.Symbol)
		instrument := &models.Instrument{
			SecurityID: &wrapperspb.UInt64Value{Value: sec.SecurityId},
			Symbol:     &wrapperspb.StringValue{Value: sec.Symbol},
		}
		done := false
		for !done {
			res, err := context.RequestFuture(state.executor, &messages.TradeCaptureReportRequest{
				RequestID: 0,
				Filter: &messages.TradeCaptureReportFilter{
					From:       utils.MilliToTimestamp(state.lastTradeTs[sec.SecurityId]),
					Instrument: instrument,
				},
				Account: state.account,
			}, 20*time.Second).Result()
			if err != nil {
				fmt.Println("error getting trade capture report", err)
				time.Sleep(1 * time.Second)
				continue
			}
			trds := res.(*messages.TradeCaptureReport)
			if !trds.Success {
				fmt.Println("error getting trade capture report", trds.RejectionReason.String())
				time.Sleep(1 * time.Second)
				continue
			}
			progress := false
			for _, trd := range trds.Trades {
				fmt.Println(trd)
				ts := trd.TransactionTime.AsTime()
				secID := trd.Instrument.SecurityID.Value
				var realized int64
				if trd.Quantity < 0 {
					_, realized = state.positions[secID].Sell(trd.Price, -trd.Quantity, false)
				} else {
					_, realized = state.positions[secID].Buy(trd.Price, trd.Quantity, false)
				}
				tx := &extypes.Transaction{
					Type:    "TRADE",
					Time:    ts,
					ID:      trd.TradeID,
					Account: state.account.Name,
					Fill: &extypes.Fill{
						SecurityID: fmt.Sprintf("%d", secID),
						Price:      trd.Price,
						Quantity:   trd.Quantity,
					},
				}
				// Realized PnL
				if realized != 0 {
					tx.Movements = append(tx.Movements, extypes.Movement{
						Reason:   int32(messages.AccountMovementType_RealizedPnl),
						AssetID:  constants.TETHER.ID,
						Quantity: -float64(realized) / 1e8,
					})
				}
				// Commission
				if trd.Commission != 0 {
					tx.Movements = append(tx.Movements, extypes.Movement{
						Reason:   int32(messages.AccountMovementType_Commission),
						AssetID:  constants.TETHER.ID,
						Quantity: -trd.Commission,
					})
				}

				if _, err := state.txs.InsertOne(goContext.Background(), tx); err != nil {
					// TODO
					//return fmt.Errorf("error inserting: %v", err)
				} else {
					state.lastTradeTs[sec.SecurityId] = uint64(ts.UnixMilli())
					progress = true
				}
			}
			if len(trds.Trades) == 0 || !progress {
				done = true
			}
		}
	}

	return nil
}

func (state *AccountReconcile) reconcileMovements(context actor.Context) error {

	return nil
}
