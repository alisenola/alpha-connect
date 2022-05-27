package ftx

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
	"gitlab.com/alphaticks/xchanger/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type AccountReconcile struct {
	extypes.BaseReconcile
	account          *models.Account
	ftxExecutor      *actor.PID
	logger           *log.Logger
	securities       map[uint64]*models.Security
	txs              *mongo.Collection
	positions        map[string]*account.Position
	lastDepositTs    uint64
	lastWithdrawalTs uint64
	lastFundingTs    uint64
	lastTradeTs      uint64
}

func NewAccountReconcileProducer(account *models.Account, txs *mongo.Collection) actor.Producer {
	return func() actor.Actor {
		return NewAccountReconcile(account, txs)
	}
}

func NewAccountReconcile(account *models.Account, txs *mongo.Collection) actor.Actor {
	return &AccountReconcile{
		account: account,
		txs:     txs,
	}
}

func (state *AccountReconcile) GetLogger() *log.Logger {
	return state.logger
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
	state.ftxExecutor = actor.NewPID(context.ActorSystem().Address(), "executor/exchanges/"+constants.FTX.Name+"_executor")
	// Request securities
	executor := actor.NewPID(context.ActorSystem().Address(), "executor")
	res, err := context.RequestFuture(executor, &messages.SecurityListRequest{}, 10*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error getting securities: %v", err)
	}
	securityList, ok := res.(*messages.SecurityList)
	if !ok {
		return fmt.Errorf("was expecting *messages.SecurityList, got %s", reflect.TypeOf(res).String())
	}
	if !securityList.Success {
		return fmt.Errorf("error getting securities: %s", securityList.RejectionReason.String())
	}
	// TODO filtering should be done by the executor, when specifying exchange in the request
	var filteredSecurities []*models.Security
	for _, s := range securityList.Securities {
		if s.Exchange.ID == state.account.Exchange.ID {
			filteredSecurities = append(filteredSecurities, s)
		}
	}

	securityMap := make(map[uint64]*models.Security)
	for _, sec := range filteredSecurities {
		securityMap[sec.SecurityID] = sec
	}
	state.securities = securityMap

	// Start reconciliation
	state.positions = make(map[string]*account.Position)
	for _, sec := range state.securities {
		if sec.SecurityType == "CRPERP" {
			tp := math.Ceil(1. / sec.MinPriceIncrement.Value)
			lp := math.Ceil(1. / sec.RoundLot.Value)
			state.positions[fmt.Sprintf("%d", sec.SecurityID)] = account.NewPosition(
				sec.IsInverse, tp, lp, 1e8, sec.Multiplier.Value, 0, 0)
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
						state.positions[tx.Fill.SecurityID].Sell(tx.Fill.Price, -tx.Fill.Quantity, false)
					} else {
						state.positions[tx.Fill.SecurityID].Buy(tx.Fill.Price, tx.Fill.Quantity, false)
					}
				}
				state.lastTradeTs = uint64(tx.Time.UnixNano() / 1000000)
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
			kid, _ := strconv.ParseUint(k, 10, 64)
			fmt.Println(state.securities[kid].Symbol, ppos.Quantity)
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
	msg := context.Message().(*messages.AccountMovementRequest)
	state.txs.Find(goContext.Background(), bson.D{
		{Key: "type", Value: msg.Type.String()},
	})
	return nil
}

func (state *AccountReconcile) OnTradeCaptureReportRequest(context actor.Context) error {
	//msg := context.Message().(*messages.TradeCaptureReportRequest)
	state.txs.Find(goContext.Background(), bson.D{})
	return nil
}

func (state *AccountReconcile) reconcileTrades(context actor.Context) error {
	done := false
	for !done {
		res, err := context.RequestFuture(state.ftxExecutor, &messages.TradeCaptureReportRequest{
			RequestID: 0,
			Filter: &messages.TradeCaptureReportFilter{
				From: utils.MilliToTimestamp(state.lastTradeTs),
				To:   timestamppb.Now(),
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
			// Check if not already inserted
			if state.txs.FindOne(goContext.Background(), bson.D{{Key: "id", Value: trd.TradeID}}).Err() != mongo.ErrNoDocuments {
				fmt.Println("skip trd", trd)
				continue
			}
			ts := trd.TransactionTime.AsTime()
			var secID string
			if trd.Instrument.SecurityID != nil {
				secID = fmt.Sprintf("%d", trd.Instrument.SecurityID.Value)
			}
			tx := &extypes.Transaction{
				Type:    "TRADE",
				Time:    ts,
				ID:      trd.TradeID,
				Account: state.account.Name,
				Fill: &extypes.Fill{
					SecurityID: secID,
					Price:      trd.Price,
					Quantity:   trd.Quantity,
				},
			}
			if trd.Instrument.SecurityID != nil {
				sec := state.securities[trd.Instrument.SecurityID.Value]
				switch sec.SecurityType {
				case "CRPERP", "CRFUT":
					var realized int64
					if trd.Quantity < 0 {
						_, realized = state.positions[secID].Sell(trd.Price, -trd.Quantity, false)
					} else {
						_, realized = state.positions[secID].Buy(trd.Price, trd.Quantity, false)
					}
					// Realized PnL
					if realized != 0 {
						tx.Movements = append(tx.Movements, extypes.Movement{
							Reason:   int32(messages.AccountMovementType_RealizedPnl),
							AssetID:  constants.DOLLAR.ID,
							Quantity: -float64(realized) / 1e8,
						})
					}
				case "CRSPOT":
					// SPOT trade
					tx.Movements = append(tx.Movements, extypes.Movement{
						Reason:   int32(messages.AccountMovementType_Exchange),
						AssetID:  sec.Underlying.ID,
						Quantity: trd.Quantity,
					})
					tx.Movements = append(tx.Movements, extypes.Movement{
						Reason:   int32(messages.AccountMovementType_Exchange),
						AssetID:  sec.QuoteCurrency.ID,
						Quantity: -trd.Quantity * trd.Price,
					})
				default:
					return fmt.Errorf("unsupported type: %s", sec.SecurityType)
				}
			} else {
				splits := strings.Split(trd.Instrument.Symbol.Value, "-")
				base, ok := constants.GetAssetBySymbol(splits[0])
				if !ok {
					fmt.Println(trd)
					return fmt.Errorf("unknown symbol: %s", splits[0])
				}
				quote, ok := constants.GetAssetBySymbol(splits[1])
				if !ok {
					return fmt.Errorf("unknown symbol: %s", splits[1])
				}
				tx.Movements = append(tx.Movements, extypes.Movement{
					Reason:   int32(messages.AccountMovementType_Exchange),
					AssetID:  base.ID,
					Quantity: trd.Quantity,
				})
				tx.Movements = append(tx.Movements, extypes.Movement{
					Reason:   int32(messages.AccountMovementType_Exchange),
					AssetID:  quote.ID,
					Quantity: -trd.Quantity * trd.Price,
				})
			}

			// Commission
			if trd.Commission != 0 {
				tx.Movements = append(tx.Movements, extypes.Movement{
					Reason:   int32(messages.AccountMovementType_Commission),
					AssetID:  trd.CommissionAsset.ID,
					Quantity: -trd.Commission,
				})
			}

			if _, err := state.txs.InsertOne(goContext.Background(), tx); err != nil {
				return fmt.Errorf("error writing transaction: %v", err)
			} else {
				state.lastTradeTs = uint64(ts.UnixNano() / 1000000)
				progress = true
			}
		}
		if len(trds.Trades) == 0 || !progress {
			done = true
		}
	}

	return nil
}

func (state *AccountReconcile) reconcileMovements(context actor.Context) error {
	// Get last account movement
	done := false
	for !done {
		res, err := context.RequestFuture(state.ftxExecutor, &messages.AccountMovementRequest{
			RequestID: 0,
			Type:      messages.AccountMovementType_FundingFee,
			Filter: &messages.AccountMovementFilter{
				From: utils.MilliToTimestamp(state.lastFundingTs),
				To:   utils.MilliToTimestamp(uint64(time.Now().UnixNano() / 1000000)),
			},
			Account: state.account,
		}, 20*time.Second).Result()
		if err != nil {
			fmt.Println("error getting movement", err)
			time.Sleep(1 * time.Second)
			continue
		}
		mvts := res.(*messages.AccountMovementResponse)
		if !mvts.Success {
			fmt.Println("error getting account movements", mvts.RejectionReason.String())
			time.Sleep(1 * time.Second)
			continue
		}
		progress := false
		for _, m := range mvts.Movements {
			ts := m.Time.AsTime()
			tx := extypes.Transaction{
				Type:    "FUNDING",
				SubType: m.Subtype,
				Time:    ts,
				ID:      m.MovementID,
				Account: state.account.Name,
				Fill:    nil,
				Movements: []extypes.Movement{{
					Reason:   int32(messages.AccountMovementType_FundingFee),
					AssetID:  m.Asset.ID,
					Quantity: m.Change,
				}},
			}
			if _, err := state.txs.InsertOne(goContext.Background(), tx); err != nil {
				// TODO
				if wexc, ok := err.(mongo.WriteException); ok && wexc.WriteErrors[0].Code == 11000 {
					continue
				} else {
					return fmt.Errorf("error writing transaction: %v", err)
				}
			} else {
				state.lastFundingTs = uint64(ts.UnixNano() / 1000000)
				progress = true
			}
		}
		if len(mvts.Movements) == 0 || !progress {
			done = true
		}
	}

	done = false
	for !done {
		fmt.Println("FETCHING DEPOSIT FTX", state.lastDepositTs)
		res, err := context.RequestFuture(state.ftxExecutor, &messages.AccountMovementRequest{
			RequestID: 0,
			Type:      messages.AccountMovementType_Deposit,
			Filter: &messages.AccountMovementFilter{
				From: utils.MilliToTimestamp(state.lastDepositTs),
				To:   utils.MilliToTimestamp(uint64(time.Now().UnixNano() / 1000000)),
			},
			Account: state.account,
		}, 20*time.Second).Result()
		if err != nil {
			fmt.Println("error getting movement", err)
			time.Sleep(1 * time.Second)
			continue
		}
		mvts := res.(*messages.AccountMovementResponse)
		if !mvts.Success {
			fmt.Println("error getting account movements", mvts.RejectionReason.String())
			time.Sleep(1 * time.Second)
			continue
		}
		progress := false
		for _, m := range mvts.Movements {
			ts := m.Time.AsTime()
			tx := extypes.Transaction{
				Type:    "DEPOSIT",
				SubType: m.Subtype,
				Time:    ts,
				ID:      m.MovementID,
				Account: state.account.Name,
				Fill:    nil,
				Movements: []extypes.Movement{{
					Reason:   int32(messages.AccountMovementType_Deposit),
					AssetID:  m.Asset.ID,
					Quantity: m.Change,
				}},
			}
			if _, err := state.txs.InsertOne(goContext.Background(), tx); err != nil {
				// TODO
				if wexc, ok := err.(mongo.WriteException); ok && wexc.WriteErrors[0].Code == 11000 {
					continue
				} else {
					return fmt.Errorf("error writing transaction: %v", err)
				}
			} else {
				progress = true
				state.lastDepositTs = uint64(ts.UnixNano() / 1000000)
			}
		}
		if len(mvts.Movements) == 0 || !progress {
			done = true
		}
	}

	done = false
	for !done {
		res, err := context.RequestFuture(state.ftxExecutor, &messages.AccountMovementRequest{
			RequestID: 0,
			Type:      messages.AccountMovementType_Withdrawal,
			Filter: &messages.AccountMovementFilter{
				From: utils.MilliToTimestamp(state.lastWithdrawalTs),
				To:   utils.MilliToTimestamp(uint64(time.Now().UnixNano() / 1000000)),
			},
			Account: state.account,
		}, 20*time.Second).Result()
		if err != nil {
			fmt.Println("error getting movement", err)
			time.Sleep(1 * time.Second)
			continue
		}
		mvts := res.(*messages.AccountMovementResponse)
		if !mvts.Success {
			fmt.Println("error getting account movements", mvts.RejectionReason.String())
			time.Sleep(1 * time.Second)
			continue
		}
		progress := false
		for _, m := range mvts.Movements {
			ts := m.Time.AsTime()
			tx := extypes.Transaction{
				Type:    "WITHDRAWAL",
				SubType: m.Subtype,
				Time:    ts,
				ID:      m.MovementID,
				Account: state.account.Name,
				Fill:    nil,
				Movements: []extypes.Movement{{
					Reason:   int32(messages.AccountMovementType_Withdrawal),
					AssetID:  m.Asset.ID,
					Quantity: m.Change,
				}},
			}
			if _, err := state.txs.InsertOne(goContext.Background(), tx); err != nil {
				// TODO
				if wexc, ok := err.(mongo.WriteException); ok && wexc.WriteErrors[0].Code == 11000 {
					continue
				} else {
					return fmt.Errorf("error writing transaction: %v", err)
				}
			} else {
				progress = true
				state.lastWithdrawalTs = uint64(ts.UnixNano() / 1000000)
			}
		}
		if len(mvts.Movements) == 0 || !progress {
			done = true
		}
	}

	return nil
}
