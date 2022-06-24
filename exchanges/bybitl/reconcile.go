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
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type AccountReconcile struct {
	extypes.BaseReconcile
	account          *models.Account
	dbAccount        *extypes.Account
	executor         *actor.PID
	logger           *log.Logger
	securities       map[uint64]*registry.Security
	symbToSecs       map[string]*registry.Security
	db               *gorm.DB
	registry         registry.PublicRegistryClient
	positions        map[uint64]*account.Position
	lastDepositTs    uint64
	lastWithdrawalTs uint64
	lastFundingTs    uint64
	lastTradeTs      map[uint64]uint64
}

func NewAccountReconcileProducer(account *models.Account, registry registry.PublicRegistryClient, db *gorm.DB) actor.Producer {
	return func() actor.Actor {
		return NewAccountReconcile(account, registry, db)
	}
}

func NewAccountReconcile(account *models.Account, registry registry.PublicRegistryClient, db *gorm.DB) actor.Actor {
	return &AccountReconcile{
		account:  account,
		db:       db,
		registry: registry,
	}
}

func (state *AccountReconcile) GetLogger() *log.Logger {
	return state.logger
}

func (state *AccountReconcile) GetTransactions() *gorm.DB {
	return state.db
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

	state.dbAccount = &extypes.Account{
		Name:       state.account.Name,
		ExchangeID: state.account.Exchange.ID,
	}
	// Check if account exists
	tx := state.db.Where("name=?", state.account.Name).FirstOrCreate(state.dbAccount)
	if tx.Error != nil {
		return fmt.Errorf("error creating account: %v", err)
	}
	var transactions []extypes.Transaction
	state.db.Debug().Model(&extypes.Transaction{}).Joins("Fill").Where(`"transactions"."account_id"=?`, state.dbAccount.ID).Order("time asc, execution_id asc").Find(&transactions)
	for _, tr := range transactions {
		switch tr.Type {
		case "TRADE":
			if tr.Fill != nil {
				secID := uint64(tr.Fill.SecurityID)
				sec, ok := state.securities[secID]
				if ok && sec.SecurityType == "CRPERP" {
					if tr.Fill.Quantity < 0 {
						state.positions[secID].Sell(tr.Fill.Price, -tr.Fill.Quantity, false)
					} else {
						state.positions[secID].Buy(tr.Fill.Price, tr.Fill.Quantity, false)
					}
					state.lastTradeTs[sec.SecurityId] = uint64(tr.Time.UnixNano() / 1000000)
				}
			}
		}
	}

	for k, pos := range state.positions {
		if ppos := pos.GetPosition(); ppos != nil {
			fmt.Println(state.securities[k].Symbol, ppos.Quantity)
		}
	}

	// Find funding transaction
	var cnt int64
	tx = state.db.Model(&extypes.Transaction{}).Where("type=?", "FUNDING").Count(&cnt)
	if tx.Error != nil {
		return fmt.Errorf("error getting funding transaction count: %v", err)
	}
	if cnt > 0 {
		var tr extypes.Transaction
		tx = state.db.Model(&extypes.Transaction{}).Where("type=?", "FUNDING").Order("time desc").First(&tr)
		if tx.Error != nil {
			return fmt.Errorf("error finding last funding transaction: %v", tx.Error)
		}
		state.lastFundingTs = uint64(tr.Time.UnixNano() / 1000000)
	}

	tx = state.db.Model(&extypes.Transaction{}).Where("type=?", "DEPOSIT").Count(&cnt)
	if tx.Error != nil {
		return fmt.Errorf("error getting deposit transaction count: %v", err)
	}
	if cnt > 0 {
		var tr extypes.Transaction
		tx = state.db.Model(&extypes.Transaction{}).Where("type=?", "DEPOSIT").Order("time desc").First(&tr)
		if tx.Error != nil {
			return fmt.Errorf("error finding last deposit transaction: %v", tx.Error)
		}
		state.lastDepositTs = uint64(tr.Time.UnixNano() / 1000000)
	}

	tx = state.db.Model(&extypes.Transaction{}).Where("type=?", "WITHDRAWAL").Count(&cnt)
	if tx.Error != nil {
		return fmt.Errorf("error getting withdrawal transaction count: %v", err)
	}
	if cnt > 0 {
		var tr extypes.Transaction
		tx = state.db.Model(&extypes.Transaction{}).Where("type=?", "WITHDRAWAL").Order("time desc").First(&tr)
		if tx.Error != nil {
			return fmt.Errorf("error finding last withdrawal transaction: %v", tx.Error)
		}
		state.lastWithdrawalTs = uint64(tr.Time.UnixNano() / 1000000)
	}

	if err := state.reconcileTrades(context); err != nil {
		return fmt.Errorf("error reconcile trade: %v", err)
	}

	// Fetch positions & compare

	// Start reconciliation
	state.positions = make(map[uint64]*account.Position)
	for _, sec := range state.securities {
		if sec.SecurityType == "CRPERP" {
			state.positions[sec.SecurityId] = account.NewPosition(
				sec.IsInverse, 1e8, 1e8, 1e8, 1, 0, 0)
		}
	}
	state.db.Debug().Model(&extypes.Transaction{}).Joins("Fill").Where(`"transactions"."account_id"=?`, state.dbAccount.ID).Order("time asc, execution_id asc").Find(&transactions)
	for _, tr := range transactions {
		switch tr.Type {
		case "TRADE":
			if tr.Fill != nil {
				secID := uint64(tr.Fill.SecurityID)
				sec, ok := state.securities[secID]
				if ok && sec.SecurityType == "CRPERP" {
					if tr.Fill.Quantity < 0 {
						state.positions[secID].Sell(tr.Fill.Price, -tr.Fill.Quantity, false)
					} else {
						state.positions[secID].Buy(tr.Fill.Price, tr.Fill.Quantity, false)
					}
					state.lastTradeTs[sec.SecurityId] = uint64(tr.Time.UnixNano() / 1000000)
				}
			}
		}
	}

	for k, pos := range state.positions {
		if ppos := pos.GetPosition(); ppos != nil {
			fmt.Println(state.securities[k].Symbol, ppos.Quantity)
		}
	}

	// Fetch positions
	resp, err := context.RequestFuture(state.executor, &messages.PositionsRequest{
		Instrument: nil,
		Account:    state.account,
	}, 10*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error getting positions from executor: %v", err)
	}

	positionList, ok := resp.(*messages.PositionList)
	if !ok {
		return fmt.Errorf("was expecting *messages.PositionList, got %s", reflect.TypeOf(resp).String())
	}
	if !positionList.Success {
		return fmt.Errorf("error getting balances: %s", positionList.RejectionReason.String())
	}

	execPositions := make(map[uint64]*models.Position)
	for _, pos := range positionList.Positions {
		execPositions[pos.Instrument.SecurityID.Value] = pos
	}

	for k, p1 := range state.positions {
		if p2, ok := execPositions[k]; ok {
			if p1.GetPosition() != nil {
				fmt.Println("position", p1.GetPosition().Quantity, p2.Quantity)
			}
		} else {
			fmt.Println("p1 not in exec")
		}
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
		fmt.Println("REQ", sec.Symbol, state.lastTradeTs[sec.SecurityId])
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
				// check if trade exists
				var cnt int64
				tx := state.db.Model(&extypes.Transaction{}).Where("execution_id=?", trd.TradeID).Count(&cnt)
				if tx.Error != nil {
					return fmt.Errorf("error getting trade transaction count: %v", err)
				}
				if cnt > 0 {
					fmt.Println("SKIP")
					continue
				}
				ts := trd.TransactionTime.AsTime()
				secID := trd.Instrument.SecurityID.Value
				var realized int64
				if trd.Quantity < 0 {
					_, realized = state.positions[secID].Sell(trd.Price, -trd.Quantity, false)
				} else {
					_, realized = state.positions[secID].Buy(trd.Price, trd.Quantity, false)
				}
				tr := &extypes.Transaction{
					Type:        "TRADE",
					Time:        ts,
					ExecutionID: trd.TradeID,
					AccountID:   state.dbAccount.ID,
					Fill: &extypes.Fill{
						AccountID:  state.dbAccount.ID,
						SecurityID: int64(secID),
						Price:      trd.Price,
						Quantity:   trd.Quantity,
					},
				}
				// Realized PnL
				if realized != 0 {
					tr.Movements = append(tr.Movements, extypes.Movement{
						AccountID: state.dbAccount.ID,
						Reason:    int32(messages.AccountMovementType_RealizedPnl),
						AssetID:   constants.TETHER.ID,
						Quantity:  -float64(realized) / 1e8,
					})
				}
				// Commission
				if trd.Commission != 0 {
					tr.Movements = append(tr.Movements, extypes.Movement{
						AccountID: state.dbAccount.ID,
						Reason:    int32(messages.AccountMovementType_Commission),
						AssetID:   constants.TETHER.ID,
						Quantity:  -trd.Commission,
					})
				}

				if tx := state.db.Create(tr); tx.Error != nil {
					return fmt.Errorf("error inserting: %v", err)
				}
				state.lastTradeTs[sec.SecurityId] = uint64(ts.UnixMilli())
				progress = true
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
		res, err := context.RequestFuture(state.executor, &messages.AccountMovementRequest{
			RequestID: 0,
			Type:      messages.AccountMovementType_FundingFee,
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
			ts := m.Time.AsTime()
			tr := &extypes.Transaction{
				Type:        "FUNDING",
				SubType:     m.Subtype,
				Time:        ts,
				ExecutionID: m.MovementID,
				AccountID:   state.dbAccount.ID,
				Fill:        nil,
				Movements: []extypes.Movement{{
					Reason:    int32(messages.AccountMovementType_FundingFee),
					AssetID:   m.Asset.ID,
					Quantity:  m.Change,
					AccountID: state.dbAccount.ID,
				}},
			}
			if tx := state.db.Create(tr); tx.Error != nil {
				return fmt.Errorf("error inserting: %v", err)
			}
			state.lastFundingTs = uint64(ts.UnixNano() / 1000000)
			progress = true
		}
		if len(mvts.Movements) == 0 || !progress {
			done = true
		}
	}
	return nil
}
