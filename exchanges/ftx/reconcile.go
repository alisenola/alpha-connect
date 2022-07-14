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
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	"gitlab.com/alphaticks/xchanger/constants"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

type AccountReconcile struct {
	extypes.BaseReconcile
	account          *models.Account
	dbAccount        *extypes.Account
	ftxExecutor      *actor.PID
	logger           *log.Logger
	securities       map[uint64]*registry.Security
	symbToSecs       map[string]*registry.Security
	registry         registry.PublicRegistryClient
	positions        map[uint64]*account.Position
	db               *gorm.DB
	lastDepositTs    uint64
	lastWithdrawalTs uint64
	lastFundingTs    uint64
	lastTradeTs      uint64
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
	res, err := state.registry.Securities(goContext.Background(), &registry.SecuritiesRequest{
		Filter: &registry.SecurityFilter{
			ExchangeId: []uint32{constants.FTX.ID},
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

	var transactions []*extypes.Transaction
	state.db.Debug().Joins("Fill").Where(`"transactions"."account_id"=?`, state.dbAccount.ID).Order("time asc").Find(&transactions)
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
					state.lastTradeTs = uint64(tr.Time.UnixNano() / 1000000)
				}
			}
		}
	}

	// Find funding transaction
	var cnt int64
	tx = state.db.Model(&extypes.Transaction{}).Where("account_id=?", state.dbAccount.ID).Where("type=?", "FUNDING").Count(&cnt)
	if tx.Error != nil {
		return fmt.Errorf("error getting funding transaction count: %v", err)
	}
	if cnt > 0 {
		var tr extypes.Transaction
		tx = state.db.Model(&extypes.Transaction{}).Where("account_id=?", state.dbAccount.ID).Where("type=?", "FUNDING").Order("time desc").First(&tr)
		if tx.Error != nil {
			return fmt.Errorf("error finding last funding transaction: %v", tx.Error)
		}
		state.lastFundingTs = uint64(tr.Time.UnixNano() / 1000000)
	}

	tx = state.db.Model(&extypes.Transaction{}).Where("account_id=?", state.dbAccount.ID).Where("type=?", "DEPOSIT").Count(&cnt)
	if tx.Error != nil {
		return fmt.Errorf("error getting deposit transaction count: %v", err)
	}
	if cnt > 0 {
		var tr extypes.Transaction
		tx = state.db.Model(&extypes.Transaction{}).Where("account_id=?", state.dbAccount.ID).Where("type=?", "DEPOSIT").Order("time desc").First(&tr)
		if tx.Error != nil {
			return fmt.Errorf("error finding last deposit transaction: %v", tx.Error)
		}
		state.lastDepositTs = uint64(tr.Time.UnixNano() / 1000000)
	}

	tx = state.db.Model(&extypes.Transaction{}).Where("account_id=?", state.dbAccount.ID).Where("type=?", "WITHDRAWAL").Count(&cnt)
	if tx.Error != nil {
		return fmt.Errorf("error getting withdrawal transaction count: %v", err)
	}
	if cnt > 0 {
		var tr extypes.Transaction
		tx = state.db.Model(&extypes.Transaction{}).Where("account_id=?", state.dbAccount.ID).Where("type=?", "WITHDRAWAL").Order("time desc").First(&tr)
		if tx.Error != nil {
			return fmt.Errorf("error finding last withdrawal transaction: %v", tx.Error)
		}
		state.lastWithdrawalTs = uint64(tr.Time.UnixNano() / 1000000)
	}

	var movements []*extypes.Movement
	state.db.Debug().Joins("Transaction").Where(`"movements"."account_id"=?`, state.dbAccount.ID).Order("time asc").Find(&movements)
	fpnl := 0.
	balances := make(map[uint32]float64)
	for _, m := range movements {
		balances[m.AssetID] += m.Quantity
		if m.Reason == 5 {
			fpnl += m.Quantity
		}
	}
	fmt.Println(balances)

	if err := state.reconcileTrades(context); err != nil {
		return fmt.Errorf("error reconcile trade: %v", err)
	}
	if err := state.reconcileMovements(context); err != nil {
		return fmt.Errorf("error reconcile trade: %v", err)
	}

	return nil
}

// TODO
func (state *AccountReconcile) Clean(context actor.Context) error {

	return nil
}

func (state *AccountReconcile) OnAccountMovementRequest(context actor.Context) error {
	//msg := context.Message().(*messages.AccountMovementRequest)
	return nil
}

func (state *AccountReconcile) OnTradeCaptureReportRequest(context actor.Context) error {
	//msg := context.Message().(*messages.TradeCaptureReportRequest)
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
			// check if trade exists
			var cnt int64
			tx := state.db.Model(&extypes.Transaction{}).Where("execution_id=?", trd.TradeID).Count(&cnt)
			if tx.Error != nil {
				return fmt.Errorf("error getting trade transaction count: %v", err)
			}
			if cnt > 0 {
				continue
			}
			ts := trd.TransactionTime.AsTime()
			secID := trd.Instrument.SecurityID.Value
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
			securityID := trd.Instrument.SecurityID
			if securityID == nil {
				sec, ok := state.symbToSecs[trd.Instrument.Symbol.Value]
				if ok {
					securityID = &wrapperspb.UInt64Value{Value: sec.SecurityId}
				}
			}
			if securityID != nil {
				sec := state.securities[securityID.Value]
				switch sec.SecurityType {
				case "CRPERP", "CRFUT":
					var realized int64
					if trd.Quantity < 0 {
						_, realized = state.positions[sec.SecurityId].Sell(trd.Price, -trd.Quantity, false)
					} else {
						_, realized = state.positions[sec.SecurityId].Buy(trd.Price, trd.Quantity, false)
					}
					// Realized PnL
					if realized != 0 {
						tr.Movements = append(tr.Movements, extypes.Movement{
							Reason:    int32(messages.AccountMovementType_RealizedPnl),
							AssetID:   constants.DOLLAR.ID,
							Quantity:  -float64(realized) / 1e8,
							AccountID: state.dbAccount.ID,
						})
					}
				case "CRSPOT":
					// SPOT trade
					b, _ := constants.GetAssetBySymbol(sec.BaseCurrency)
					tr.Movements = append(tr.Movements, extypes.Movement{
						Reason:    int32(messages.AccountMovementType_Exchange),
						AssetID:   b.ID,
						Quantity:  trd.Quantity,
						AccountID: state.dbAccount.ID,
					})
					q, _ := constants.GetAssetBySymbol(sec.QuoteCurrency)
					tr.Movements = append(tr.Movements, extypes.Movement{
						Reason:    int32(messages.AccountMovementType_Exchange),
						AssetID:   q.ID,
						Quantity:  -trd.Quantity * trd.Price,
						AccountID: state.dbAccount.ID,
					})
				default:
					return fmt.Errorf("unsupported type: %s", sec.SecurityType)
				}
			} else {
				// OTC Conversion
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
				tr.Movements = append(tr.Movements, extypes.Movement{
					Reason:    int32(messages.AccountMovementType_Exchange),
					AssetID:   base.ID,
					Quantity:  trd.Quantity,
					AccountID: state.dbAccount.ID,
				})
				tr.Movements = append(tr.Movements, extypes.Movement{
					Reason:    int32(messages.AccountMovementType_Exchange),
					AssetID:   quote.ID,
					Quantity:  -trd.Quantity * trd.Price,
					AccountID: state.dbAccount.ID,
				})
			}

			// Commission
			if trd.Commission != 0 {
				tr.Movements = append(tr.Movements, extypes.Movement{
					Reason:    int32(messages.AccountMovementType_Commission),
					AssetID:   trd.CommissionAsset.ID,
					Quantity:  -trd.Commission,
					AccountID: state.dbAccount.ID,
				})
			}

			if tx := state.db.Create(tr); tx.Error != nil {
				return fmt.Errorf("error inserting: %v", err)
			}
			state.lastTradeTs = uint64(ts.UnixMilli())
			progress = true
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
			var cnt int64
			tx := state.db.Model(&extypes.Transaction{}).Where("execution_id=?", m.MovementID).Count(&cnt)
			if tx.Error != nil {
				return fmt.Errorf("error getting transaction count: %v", err)
			}
			if cnt > 0 {
				continue
			}
			ts := m.Time.AsTime()
			tr := &extypes.Transaction{
				Type:        "FUNDING",
				SubType:     m.Subtype,
				Time:        ts,
				ExecutionID: m.MovementID,
				AccountID:   state.dbAccount.ID,
				Movements: []extypes.Movement{{
					Reason:    int32(messages.AccountMovementType_FundingFee),
					AssetID:   m.Asset.ID,
					Quantity:  m.Change,
					AccountID: state.dbAccount.ID,
				}},
			}
			if tx := state.db.Debug().Create(tr); tx.Error != nil {
				return fmt.Errorf("error inserting: %v", err)
			}
			progress = true
			state.lastFundingTs = uint64(ts.UnixNano() / 1000000)
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
			var cnt int64
			tx := state.db.Model(&extypes.Transaction{}).Where("execution_id=?", m.MovementID).Count(&cnt)
			if tx.Error != nil {
				return fmt.Errorf("error getting transaction count: %v", err)
			}
			if cnt > 0 {
				continue
			}
			ts := m.Time.AsTime()
			tr := &extypes.Transaction{
				Type:        "DEPOSIT",
				SubType:     m.Subtype,
				Time:        ts,
				ExecutionID: m.MovementID,
				AccountID:   state.dbAccount.ID,
				Fill:        nil,
				Movements: []extypes.Movement{{
					Reason:    int32(messages.AccountMovementType_Deposit),
					AssetID:   m.Asset.ID,
					Quantity:  m.Change,
					AccountID: state.dbAccount.ID,
				}},
			}
			if tx := state.db.Create(tr); tx.Error != nil {
				return fmt.Errorf("error inserting: %v", err)
			}
			progress = true
			state.lastDepositTs = uint64(ts.UnixNano() / 1000000)
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
			var cnt int64
			tx := state.db.Model(&extypes.Transaction{}).Where("execution_id=?", m.MovementID).Count(&cnt)
			if tx.Error != nil {
				return fmt.Errorf("error getting transaction count: %v", err)
			}
			if cnt > 0 {
				continue
			}
			ts := m.Time.AsTime()
			tr := &extypes.Transaction{
				Type:        "WITHDRAWAL",
				SubType:     m.Subtype,
				Time:        ts,
				ExecutionID: m.MovementID,
				AccountID:   state.dbAccount.ID,
				Fill:        nil,
				Movements: []extypes.Movement{{
					Reason:    int32(messages.AccountMovementType_Withdrawal),
					AssetID:   m.Asset.ID,
					Quantity:  m.Change,
					AccountID: state.dbAccount.ID,
				}},
			}
			if tx := state.db.Create(tr); tx.Error != nil {
				return fmt.Errorf("error inserting: %v", err)
			}

			progress = true
			state.lastWithdrawalTs = uint64(ts.UnixNano() / 1000000)
		}
		if len(mvts.Movements) == 0 || !progress {
			done = true
		}
	}

	return nil
}
