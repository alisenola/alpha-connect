package fbinance

import (
	goContext "context"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/account"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/xchanger/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type AccountReconcile struct {
	account          *models.Account
	fbinanceExecutor *actor.PID
	executorManager  *actor.PID
	logger           *log.Logger
	securities       map[uint64]*models.Security
	txs              *mongo.Collection
	positions        map[string]*account.Position
	lastDepositTs    uint64
	lastWithdrawalTs uint64
	lastFundingTs    uint64
	lastTradeID      map[uint64]uint64
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

func (state *AccountReconcile) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.logger.Error("error initializing", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor started")

	case *actor.Stopping:
		if err := state.Clean(context); err != nil {
			state.logger.Error("error stopping", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor stopping")

	case *actor.Stopped:
		state.logger.Info("actor stopped")

	case *actor.Restarting:
		if err := state.Clean(context); err != nil {
			state.logger.Error("error restarting", log.Error(err))
			// Attention, no panic in restarting or infinite loop
		}
		state.logger.Info("actor restarting")

	case *messages.AccountMovementRequest:
		if err := state.OnAccountMovementRequest(context); err != nil {
			state.logger.Error("error processing OnAccountMovementRequest", log.Error(err))
			panic(err)
		}

	case *messages.TradeCaptureReportRequest:
		if err := state.OnTradeCaptureReportRequest(context); err != nil {
			state.logger.Error("error processing OnTradeCaptureReportRequest", log.Error(err))
			panic(err)
		}
	}
}

func (state *AccountReconcile) Initialize(context actor.Context) error {
	// When initialize is done, the account must be aware of all the settings / assets / portfolio
	// so as to be able to answer to FIX messages

	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))
	state.fbinanceExecutor = actor.NewPID(context.ActorSystem().Address(), "executor/"+constants.FBINANCE.Name+"_executor")

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
	state.lastTradeID = make(map[uint64]uint64)
	for _, sec := range state.securities {
		tp := math.Ceil(1. / sec.MinPriceIncrement.Value)
		lp := math.Ceil(1. / sec.RoundLot.Value)
		state.positions[fmt.Sprintf("%d", sec.SecurityID)] = account.NewPosition(
			sec.IsInverse, tp, lp, 1e8, sec.Multiplier.Value, 0, 0)

		state.lastTradeID[sec.SecurityID] = 0
	}
	// First, calculate current positions from historical
	cur, err := state.txs.Find(goContext.Background(), bson.D{
		{"account", state.account.Name},
	}, options.Find().SetSort(bson.D{
		{"time", 1},
		{"_id", 1},
	}))
	if err != nil {
		return fmt.Errorf("error reconcile trades: %v", err)
	}

	balances := make(map[uint32]float64)
	var realized2, realized1 float64
	for cur.Next(goContext.Background()) {
		var tx extypes.Transaction
		if err := cur.Decode(&tx); err != nil {
			return fmt.Errorf("error decoding transaction: %v", err)
		}
		switch tx.Type {
		case "TRADE":
			secID, _ := strconv.ParseUint(tx.Fill.SecurityID, 10, 64)
			sec := state.securities[secID]
			if sec.SecurityType == "CRPERP" {
				if tx.Fill.Quantity < 0 {
					_, realized := state.positions[tx.Fill.SecurityID].Sell(tx.Fill.Price, -tx.Fill.Quantity, false)
					realized1 -= float64(realized) / 1e8
				} else {
					_, realized := state.positions[tx.Fill.SecurityID].Buy(tx.Fill.Price, tx.Fill.Quantity, false)
					realized1 -= float64(realized) / 1e8
				}
			}
			tradeID, _ := strconv.ParseUint(strings.Split(tx.ID, "-")[0], 10, 64)
			state.lastTradeID[secID] = tradeID
		}
		for _, m := range tx.Movements {
			if messages.AccountMovementType(m.Reason) == messages.RealizedPnl {
				realized2 += m.Quantity
			}
			balances[m.AssetID] += m.Quantity
		}
	}
	fmt.Println("pnl", realized1, realized2)

	for k, b := range balances {
		a, _ := constants.GetAssetByID(k)
		fmt.Println(a.Symbol, b)
	}

	for k, pos := range state.positions {
		if ppos := pos.GetPosition(); ppos != nil {
			kid, _ := strconv.ParseUint(k, 10, 64)
			fmt.Println(state.securities[kid].Symbol, ppos.Quantity)
		}
	}

	sres := state.txs.FindOne(goContext.Background(), bson.D{
		{"account", state.account.Name},
		{"type", "FUNDING"},
	}, options.FindOne().SetSort(bson.D{{"_id", -1}}))
	if sres.Err() != nil {
		if sres.Err() != mongo.ErrNoDocuments {
			return fmt.Errorf("error getting last funding: %v", err)
		} else {
			state.lastFundingTs = uint64(time.Now().Add(-365*24*time.Hour).UnixNano() / 1000000)
		}
	} else {
		var tx extypes.Transaction
		if err := sres.Decode(&tx); err != nil {
			return fmt.Errorf("error decoding transaction: %v", err)
		}
		state.lastFundingTs = uint64(tx.Time.UnixNano() / 1000000)
	}

	sres = state.txs.FindOne(goContext.Background(), bson.D{
		{"account", state.account.Name},
		{"type", "DEPOSIT"},
	}, options.FindOne().SetSort(bson.D{{"_id", -1}}))
	if sres.Err() != nil {
		if sres.Err() != mongo.ErrNoDocuments {
			return fmt.Errorf("error getting last deposit: %v", err)
		} else {
			state.lastDepositTs = uint64(time.Now().Add(-365*24*time.Hour).UnixNano() / 1000000)
		}
	} else {
		var tx extypes.Transaction
		if err := sres.Decode(&tx); err != nil {
			return fmt.Errorf("error decoding transaction: %v", err)
		}
		state.lastDepositTs = uint64(tx.Time.UnixNano() / 1000000)
	}

	sres = state.txs.FindOne(goContext.Background(), bson.D{
		{"account", state.account.Name},
		{"type", "WITHDRAWAL"},
	}, options.FindOne().SetSort(bson.D{{"_id", -1}}))
	if sres.Err() != nil {
		if sres.Err() != mongo.ErrNoDocuments {
			return fmt.Errorf("error getting last withdrawal: %v", err)
		} else {
			state.lastWithdrawalTs = uint64(time.Now().Add(-365*24*time.Hour).UnixNano() / 1000000)
		}
	} else {
		var tx extypes.Transaction
		if err := sres.Decode(&tx); err != nil {
			return fmt.Errorf("error decoding transaction: %v", err)
		}
		state.lastWithdrawalTs = uint64(tx.Time.UnixNano() / 1000000)
	}

	/*
		if err := state.reconcileTrades(context); err != nil {
			return fmt.Errorf("error reconcile trade: %v", err)
		}

	*/

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
		{"type", msg.Type.String()},
	})
	return nil
}

func (state *AccountReconcile) OnTradeCaptureReportRequest(context actor.Context) error {
	//msg := context.Message().(*messages.TradeCaptureReportRequest)
	state.txs.Find(goContext.Background(), bson.D{})
	return nil
}

func (state *AccountReconcile) reconcileTrades(context actor.Context) error {
	for _, sec := range state.securities {
		instrument := &models.Instrument{
			SecurityID: &types.UInt64Value{Value: sec.SecurityID},
			Symbol:     &types.StringValue{Value: sec.Symbol},
		}
		done := false
		for !done {
			res, err := context.RequestFuture(state.fbinanceExecutor, &messages.TradeCaptureReportRequest{
				RequestID: 0,
				Filter: &messages.TradeCaptureReportFilter{
					FromID:     &types.StringValue{Value: fmt.Sprintf("%d", state.lastTradeID[sec.SecurityID]+1)},
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
				ts, _ := types.TimestampFromProto(trd.TransactionTime)
				secID := fmt.Sprintf("%d", trd.Instrument.SecurityID.Value)
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
						SecurityID: secID,
						Price:      trd.Price,
						Quantity:   trd.Quantity,
					},
				}
				// Realized PnL
				if realized != 0 {
					tx.Movements = append(tx.Movements, extypes.Movement{
						Reason:   int32(messages.RealizedPnl),
						AssetID:  constants.TETHER.ID,
						Quantity: -float64(realized) / 1e8,
					})
				}
				// Commission
				if trd.Commission != 0 {
					tx.Movements = append(tx.Movements, extypes.Movement{
						Reason:   int32(messages.Commission),
						AssetID:  constants.TETHER.ID,
						Quantity: -trd.Commission,
					})
				}

				tradeIDInt, _ := strconv.ParseUint(strings.Split(trd.TradeID, "-")[0], 10, 64)
				if _, err := state.txs.InsertOne(goContext.Background(), tx); err != nil {
					// TODO
					return fmt.Errorf("error inserting: %v", err)
				} else {
					state.lastTradeID[sec.SecurityID] = tradeIDInt
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
	// Get last account movement
	done := false
	for !done {
		res, err := context.RequestFuture(state.fbinanceExecutor, &messages.AccountMovementRequest{
			RequestID: 0,
			Type:      messages.FundingFee,
			Filter: &messages.AccountMovementFilter{
				From: utils.MilliToTimestamp(state.lastFundingTs + 1),
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
			ts, _ := types.TimestampFromProto(m.Time)
			tx := extypes.Transaction{
				Type:    "FUNDING",
				SubType: m.Subtype,
				Time:    ts,
				ID:      m.MovementID,
				Account: state.account.Name,
				Fill:    nil,
				Movements: []extypes.Movement{{
					Reason:   int32(messages.FundingFee),
					AssetID:  m.Asset.ID,
					Quantity: m.Change,
				}},
			}
			if _, err := state.txs.InsertOne(goContext.Background(), tx); err != nil {
				return fmt.Errorf("error inserting: %v", err)
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
		res, err := context.RequestFuture(state.fbinanceExecutor, &messages.AccountMovementRequest{
			RequestID: 0,
			Type:      messages.Deposit,
			Filter: &messages.AccountMovementFilter{
				From: utils.MilliToTimestamp(state.lastDepositTs + 1),
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
			ts, _ := types.TimestampFromProto(m.Time)
			tx := extypes.Transaction{
				Type:    "DEPOSIT",
				SubType: m.Subtype,
				Time:    ts,
				ID:      m.MovementID,
				Account: state.account.Name,
				Fill:    nil,
				Movements: []extypes.Movement{{
					Reason:   int32(messages.Deposit),
					AssetID:  m.Asset.ID,
					Quantity: m.Change,
				}},
			}
			if _, err := state.txs.InsertOne(goContext.Background(), tx); err != nil {
				return fmt.Errorf("error inserting: %v", err)
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
		res, err := context.RequestFuture(state.fbinanceExecutor, &messages.AccountMovementRequest{
			RequestID: 0,
			Type:      messages.Withdrawal,
			Filter: &messages.AccountMovementFilter{
				From: utils.MilliToTimestamp(state.lastWithdrawalTs + 1),
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
			ts, _ := types.TimestampFromProto(m.Time)
			tx := extypes.Transaction{
				Type:    "WITHDRAWAL",
				SubType: m.Subtype,
				Time:    ts,
				ID:      m.MovementID,
				Account: state.account.Name,
				Fill:    nil,
				Movements: []extypes.Movement{{
					Reason:   int32(messages.Withdrawal),
					AssetID:  m.Asset.ID,
					Quantity: m.Change,
				}},
			}
			if _, err := state.txs.InsertOne(goContext.Background(), tx); err != nil {
				return fmt.Errorf("error inserting: %v", err)
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
