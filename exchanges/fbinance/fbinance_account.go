package fbinance

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/account"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/fbinance"
	"gitlab.com/alphaticks/xchanger/utils"
	"math"
	"net/http"
	"reflect"
	"sort"
	"time"
)

type checkSocket struct{}
type checkAccount struct{}

type AccountListener struct {
	account             *account.Account
	seqNum              uint64
	fbinanceExecutor    *actor.PID
	ws                  *fbinance.AuthWebsocket
	executorManager     *actor.PID
	logger              *log.Logger
	checkAccountPending bool
	checkSocketPending  bool
	lastPingTime        time.Time
	securities          []*models.Security
	client              *http.Client
}

func NewAccountListenerProducer(account *account.Account) actor.Producer {
	return func() actor.Actor {
		return NewAccountListener(account)
	}
}

func NewAccountListener(account *account.Account) actor.Actor {
	return &AccountListener{
		account:             account,
		seqNum:              0,
		ws:                  nil,
		executorManager:     nil,
		logger:              nil,
		checkSocketPending:  false,
		checkAccountPending: false,
	}
}

func (state *AccountListener) Receive(context actor.Context) {
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

	case *messages.PositionsRequest:
		if err := state.OnPositionsRequest(context); err != nil {
			state.logger.Error("error processing OnPositionListRequest", log.Error(err))
			panic(err)
		}

	case *messages.BalancesRequest:
		if err := state.OnBalancesRequest(context); err != nil {
			state.logger.Error("error processing OnBalancesRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderStatusRequest:
		if err := state.OnOrderStatusRequest(context); err != nil {
			state.logger.Error("error processing OnOrderStatusRequset", log.Error(err))
			panic(err)
		}

	case *messages.NewOrderSingleRequest:
		if err := state.OnNewOrderSingle(context); err != nil {
			state.logger.Error("error processing OnNewOrderSingle", log.Error(err))
			panic(err)
		}

	case *messages.NewOrderBulkRequest:
		if err := state.OnNewOrderBulkRequest(context); err != nil {
			state.logger.Error("error processing NewOrderBulkRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderReplaceRequest:
		if err := state.OnOrderReplaceRequest(context); err != nil {
			state.logger.Error("error processing OrderReplaceRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderBulkReplaceRequest:
		if err := state.OnBulkOrderReplaceRequest(context); err != nil {
			state.logger.Error("error processing BulkOrderReplaceRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderCancelRequest:
		if err := state.OnOrderCancelRequest(context); err != nil {
			state.logger.Error("error processing OnOrderCancelRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderMassCancelRequest:
		if err := state.OnOrderMassCancelRequest(context); err != nil {
			state.logger.Error("error processing OnOrderMassCancelRequest", log.Error(err))
			panic(err)
		}

	case *fbinance.AuthWebsocketMessage:
		if err := state.onAuthWebsocketMessage(context); err != nil {
			state.logger.Error("error processing AuthWebsocketMessage", log.Error(err))
			panic(err)
		}

	case *checkSocket:
		state.checkSocketPending = false
		if err := state.checkSocket(context); err != nil {
			state.logger.Error("error checking socket", log.Error(err))
			panic(err)
		}
		if !state.checkSocketPending {
			state.checkSocketPending = true
			go func(pid *actor.PID) {
				time.Sleep(5 * time.Second)
				context.Send(pid, &checkSocket{})
			}(context.Self())
		}

	case *checkAccount:
		state.checkAccountPending = false
		if err := state.checkAccount(context); err != nil {
			state.logger.Error("error checking socket", log.Error(err))
			panic(err)
		}
		if !state.checkAccountPending {
			state.checkAccountPending = true
			go func(pid *actor.PID) {
				time.Sleep(10 * time.Minute)
				context.Send(pid, &checkAccount{})
			}(context.Self())
		}
	}
}

func (state *AccountListener) Initialize(context actor.Context) error {
	// When initialize is done, the account must be aware of all the settings / assets / portofilio
	// so as to be able to answer to FIX messages

	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))
	state.fbinanceExecutor = actor.NewLocalPID("executor/" + constants.FBINANCE.Name + "_executor")
	state.client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1024,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	if err := state.subscribeAccount(context); err != nil {
		return fmt.Errorf("error subscribing to account: %v", err)
	}
	// Request securities
	executor := actor.NewLocalPID("executor")
	res, err := actor.EmptyRootContext.RequestFuture(executor, &messages.SecurityListRequest{}, 10*time.Second).Result()
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
	state.securities = filteredSecurities

	// Then fetch balances
	res, err = context.RequestFuture(state.fbinanceExecutor, &messages.BalancesRequest{
		Account: state.account.Account,
	}, 10*time.Second).Result()

	if err != nil {
		return fmt.Errorf("error getting balances from executor: %v", err)
	}

	balanceList, ok := res.(*messages.BalanceList)
	if !ok {
		return fmt.Errorf("was expecting BalanceList, got %s", reflect.TypeOf(res).String())
	}

	if !balanceList.Success {
		return fmt.Errorf("error getting balances: %s", balanceList.RejectionReason.String())
	}

	// Then fetch positions
	res, err = context.RequestFuture(state.fbinanceExecutor, &messages.PositionsRequest{
		Instrument: nil,
		Account:    state.account.Account,
	}, 10*time.Second).Result()

	if err != nil {
		return fmt.Errorf("error getting positions from executor: %v", err)
	}

	positionList, ok := res.(*messages.PositionList)
	if !ok {
		return fmt.Errorf("was expecting PositionList, got %s", reflect.TypeOf(res).String())
	}

	if !positionList.Success {
		return fmt.Errorf("error getting positions: %s", positionList.RejectionReason.String())
	}

	// Then fetch orders
	res, err = context.RequestFuture(state.fbinanceExecutor, &messages.OrderStatusRequest{
		Account: state.account.Account,
	}, 10*time.Second).Result()

	if err != nil {
		return fmt.Errorf("error getting orders from executor: %v", err)
	}

	orderList, ok := res.(*messages.OrderList)
	if !ok {
		return fmt.Errorf("was expecting OrderList, got %s", reflect.TypeOf(res).String())
	}

	if !orderList.Success {
		return fmt.Errorf("error fetching orders: %s", orderList.RejectionReason.String())
	}

	tetherMargin := 0.
	var filteredBalances []*models.Balance
	for _, b := range balanceList.Balances {
		if b.Asset.ID == constants.TETHER.ID {
			tetherMargin = b.Quantity
		} else {
			filteredBalances = append(filteredBalances, b)
		}
	}
	// Sync account
	if err := state.account.Sync(filteredSecurities, orderList.Orders, positionList.Positions, filteredBalances, tetherMargin, nil, nil); err != nil {
		return fmt.Errorf("error syncing account: %v", err)
	}

	state.seqNum = 0

	if !state.checkSocketPending {
		go func(pid *actor.PID) {
			time.Sleep(5 * time.Second)
			context.Send(pid, &checkSocket{})
		}(context.Self())
		state.checkSocketPending = true
	}
	if !state.checkAccountPending {
		go func(pid *actor.PID) {
			time.Sleep(10 * time.Minute)
			context.Send(pid, &checkAccount{})
		}(context.Self())
		state.checkAccountPending = true
	}
	return nil
}

// TODO
func (state *AccountListener) Clean(context actor.Context) error {
	if state.ws != nil {
		if err := state.ws.Disconnect(); err != nil {
			state.logger.Info("error disconnecting socket", log.Error(err))
		}
	}

	return nil
}

func (state *AccountListener) OnPositionsRequest(context actor.Context) error {
	msg := context.Message().(*messages.PositionsRequest)
	// TODO FILTER
	positions := state.account.GetPositions()
	context.Respond(&messages.PositionList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Positions:  positions,
	})
	return nil
}

func (state *AccountListener) OnBalancesRequest(context actor.Context) error {
	msg := context.Message().(*messages.BalancesRequest)
	// TODO FILTER
	state.account.Settle()
	balances := state.account.GetBalances()
	context.Respond(&messages.BalanceList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Balances:   balances,
	})
	return nil
}

func (state *AccountListener) OnOrderStatusRequest(context actor.Context) error {
	fmt.Println("GOT ORDERE STATUS REQUEST")
	req := context.Message().(*messages.OrderStatusRequest)
	orders := state.account.GetOrders(req.Filter)
	context.Respond(&messages.OrderList{
		RequestID: req.RequestID,
		Success:   true,
		Orders:    orders,
	})
	return nil
}

func (state *AccountListener) OnNewOrderSingle(context actor.Context) error {
	req := context.Message().(*messages.NewOrderSingleRequest)
	req.Account = state.account.Account
	// Check order quantity
	order := &models.Order{
		OrderID:        "",
		ClientOrderID:  req.Order.ClientOrderID,
		Instrument:     req.Order.Instrument,
		OrderStatus:    models.PendingNew,
		OrderType:      req.Order.OrderType,
		Side:           req.Order.OrderSide,
		TimeInForce:    req.Order.TimeInForce,
		LeavesQuantity: req.Order.Quantity,
		Price:          req.Order.Price,
		CumQuantity:    0,
	}
	report, res := state.account.NewOrder(order)
	if res != nil {
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID:       req.RequestID,
			Success:         false,
			RejectionReason: *res,
		})
	} else {
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID: req.RequestID,
			Success:   true,
		})
		if report != nil {
			report.SeqNum = state.seqNum + 1
			state.seqNum += 1
			context.Send(context.Parent(), report)
			if report.ExecutionType == messages.PendingNew {
				fmt.Println("PENDING NEW SENDING")
				fut := context.RequestFuture(state.fbinanceExecutor, req, 10*time.Second)
				context.AwaitFuture(fut, func(res interface{}, err error) {
					if err != nil {
						fmt.Println("REJECTING")
						report, err := state.account.RejectNewOrder(order.ClientOrderID, messages.Other)
						if err != nil {
							panic(err)
						}
						context.Respond(&messages.NewOrderSingleResponse{
							RequestID:       req.RequestID,
							Success:         false,
							RejectionReason: messages.Other,
						})
						if report != nil {
							report.SeqNum = state.seqNum + 1
							state.seqNum += 1
							context.Send(context.Parent(), report)
						}
						return
					}
					response := res.(*messages.NewOrderSingleResponse)
					context.Respond(response)

					if response.Success {
						nReport, _ := state.account.ConfirmNewOrder(order.ClientOrderID, response.OrderID)
						if nReport != nil {
							nReport.SeqNum = state.seqNum + 1
							state.seqNum += 1
							context.Send(context.Parent(), nReport)
						}
					} else {
						fmt.Println("REJECTING", response.RejectionReason.String())
						nReport, _ := state.account.RejectNewOrder(order.ClientOrderID, response.RejectionReason)
						if nReport != nil {
							nReport.SeqNum = state.seqNum + 1
							state.seqNum += 1
							context.Send(context.Parent(), nReport)
						}
					}
				})
			}
		}
	}

	return nil
}

func (state *AccountListener) OnNewOrderBulkRequest(context actor.Context) error {
	req := context.Message().(*messages.NewOrderBulkRequest)
	req.Account = state.account.Account
	reports := make([]*messages.ExecutionReport, 0, len(req.Orders))
	response := &messages.NewOrderBulkResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	symbol := ""
	if req.Orders[0].Instrument != nil && req.Orders[0].Instrument.Symbol != nil {
		symbol = req.Orders[0].Instrument.Symbol.Value
	}
	for _, order := range req.Orders {
		if order.Instrument == nil || order.Instrument.Symbol == nil {
			response.RejectionReason = messages.UnknownSymbol
			context.Respond(response)
			return nil
		} else if symbol != order.Instrument.Symbol.Value {
			response.RejectionReason = messages.DifferentSymbols
			context.Respond(response)
			return nil
		}
	}
	for _, reqOrder := range req.Orders {
		order := &models.Order{
			OrderID:        "",
			ClientOrderID:  reqOrder.ClientOrderID,
			Instrument:     reqOrder.Instrument,
			OrderStatus:    models.PendingNew,
			OrderType:      reqOrder.OrderType,
			Side:           reqOrder.OrderSide,
			TimeInForce:    reqOrder.TimeInForce,
			LeavesQuantity: reqOrder.Quantity,
			Price:          reqOrder.Price,
			CumQuantity:    0,
		}
		report, res := state.account.NewOrder(order)
		if res != nil {
			// Cancel all new order up until now
			for _, r := range reports {
				_, err := state.account.RejectNewOrder(r.ClientOrderID.Value, messages.Other)
				if err != nil {
					return err
				}
			}
			response.RejectionReason = *res
			context.Respond(response)
			return nil
		}
		if report != nil {
			reports = append(reports, report)
		}
	}

	response.Success = true
	context.Respond(response)

	for _, report := range reports {
		report.SeqNum = state.seqNum + 1
		state.seqNum += 1
		context.Send(context.Parent(), report)
	}
	fut := context.RequestFuture(state.fbinanceExecutor, req, 10*time.Second)
	context.AwaitFuture(fut, func(res interface{}, err error) {
		if err != nil {
			for _, r := range reports {
				report, err := state.account.RejectNewOrder(r.ClientOrderID.Value, messages.Other)
				if err != nil {
					panic(err)
				}

				if report != nil {
					report.SeqNum = state.seqNum + 1
					state.seqNum += 1
					context.Send(context.Parent(), report)
				}
			}
			context.Respond(&messages.NewOrderBulkResponse{
				RequestID:       req.RequestID,
				Success:         false,
				RejectionReason: messages.Other,
			})
			return
		}
		response := res.(*messages.NewOrderBulkResponse)
		context.Respond(response)
		if response.Success {
			for i, r := range reports {
				report, err := state.account.ConfirmNewOrder(r.ClientOrderID.Value, response.OrderIDs[i])
				if err != nil {
					panic(err)
				}

				if report != nil {
					report.SeqNum = state.seqNum + 1
					state.seqNum += 1
					context.Send(context.Parent(), report)
				}
			}
		} else {
			for _, r := range reports {
				report, err := state.account.RejectNewOrder(r.ClientOrderID.Value, messages.Other)
				if err != nil {
					panic(err)
				}

				if report != nil {
					report.SeqNum = state.seqNum + 1
					state.seqNum += 1
					context.Send(context.Parent(), report)
				}
			}
		}
	})
	return nil
}

func (state *AccountListener) OnOrderReplaceRequest(context actor.Context) error {
	// TODO
	return nil
}

func (state *AccountListener) OnBulkOrderReplaceRequest(context actor.Context) error {
	//TODO
	return nil
}

func (state *AccountListener) OnOrderCancelRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderCancelRequest)
	var ID string
	if req.ClientOrderID != nil {
		ID = req.ClientOrderID.Value
	} else if req.OrderID != nil {
		ID = req.OrderID.Value
	}
	report, res := state.account.CancelOrder(ID)
	if res != nil {
		context.Respond(&messages.OrderCancelResponse{
			RequestID:       req.RequestID,
			RejectionReason: *res,
		})
	} else {
		context.Respond(&messages.OrderCancelResponse{
			RequestID: req.RequestID,
			Success:   true,
		})
		if report != nil {
			report.SeqNum = state.seqNum + 1
			state.seqNum += 1
			context.Send(context.Parent(), report)
			if report.ExecutionType == messages.PendingCancel {
				fut := context.RequestFuture(state.fbinanceExecutor, req, 10*time.Second)
				context.AwaitFuture(fut, func(res interface{}, err error) {
					if err != nil {
						report, err := state.account.RejectCancelOrder(ID, messages.Other)
						if err != nil {
							panic(err)
						}
						context.Respond(&messages.OrderCancelResponse{
							RequestID:       req.RequestID,
							Success:         false,
							RejectionReason: messages.Other,
						})
						if report != nil {
							report.SeqNum = state.seqNum + 1
							state.seqNum += 1
							context.Send(context.Parent(), report)
						}
						return
					}
					response := res.(*messages.OrderCancelResponse)
					context.Respond(response)

					if !response.Success {
						report, err := state.account.RejectCancelOrder(ID, response.RejectionReason)
						if err != nil {
							panic(err)
						}
						if report != nil {
							report.SeqNum = state.seqNum + 1
							state.seqNum += 1
							context.Send(context.Parent(), report)
						}
					}
				})
			}
		}
	}

	return nil
}

func (state *AccountListener) OnOrderMassCancelRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderMassCancelRequest)
	orders := state.account.GetOrders(req.Filter)
	if len(orders) == 0 {
		context.Respond(&messages.OrderMassCancelResponse{
			RequestID: req.RequestID,
			Success:   true,
		})
		return nil
	}
	var reports []*messages.ExecutionReport
	for _, o := range orders {
		if o.OrderStatus != models.New && o.OrderStatus != models.PartiallyFilled {
			continue
		}
		report, res := state.account.CancelOrder(o.ClientOrderID)
		if res != nil {
			// Reject all cancel order up until now
			for _, r := range reports {
				_, err := state.account.RejectCancelOrder(r.ClientOrderID.Value, messages.Other)
				if err != nil {
					return err
				}
			}

			context.Respond(&messages.OrderMassCancelResponse{
				RequestID:       req.RequestID,
				Success:         false,
				RejectionReason: *res,
			})

			return nil
		} else if report != nil {
			reports = append(reports, report)
		}
	}

	context.Respond(&messages.OrderMassCancelResponse{
		RequestID: req.RequestID,
		Success:   true,
	})

	for _, report := range reports {
		report.SeqNum = state.seqNum + 1
		state.seqNum += 1
		context.Send(context.Parent(), report)
	}
	fmt.Println("REQUESTING AMSS CANCEL FUTURE")
	fut := context.RequestFuture(state.fbinanceExecutor, req, 10*time.Second)
	context.AwaitFuture(fut, func(res interface{}, err error) {
		if err != nil {
			for _, r := range reports {
				report, err := state.account.RejectCancelOrder(r.ClientOrderID.Value, messages.Other)
				if err != nil {
					panic(err)
				}
				if report != nil {
					report.SeqNum = state.seqNum + 1
					state.seqNum += 1
					context.Send(context.Parent(), report)
				}
			}
			context.Respond(&messages.OrderMassCancelResponse{
				RequestID:       req.RequestID,
				Success:         false,
				RejectionReason: messages.Other,
			})

			return
		}
		response := res.(*messages.OrderMassCancelResponse)
		context.Respond(response)
		if !response.Success {
			for _, r := range reports {
				report, err := state.account.RejectCancelOrder(r.ClientOrderID.Value, messages.Other)
				if err != nil {
					panic(err)
				}
				if report != nil {
					report.SeqNum = state.seqNum + 1
					state.seqNum += 1
					context.Send(context.Parent(), report)
				}
			}
		}
	})

	return nil
}

func (state *AccountListener) onAuthWebsocketMessage(context actor.Context) error {
	state.lastPingTime = time.Now()
	msg := context.Message().(*fbinance.AuthWebsocketMessage)
	if msg.Message == nil {
		return fmt.Errorf("received nil message")
	}
	switch msg.Message.Event {
	case fbinance.ORDER_TRADE_UPDATE:
		fmt.Println("EXECUTION", msg.Message.Execution)
		if msg.Message.Execution == nil {
			return fmt.Errorf("received ORDER_TRADE_UPDATE with no execution data")
		}
		exec := msg.Message.Execution
		switch exec.ExecutionType {
		case fbinance.ET_NEW:
			// New order
			orderID := fmt.Sprintf("%d", exec.OrderID)
			report, err := state.account.ConfirmNewOrder(exec.ClientOrderID, orderID)
			if err != nil {
				return fmt.Errorf("error confirming new order: %v", err)
			}
			if report != nil {
				report.SeqNum = state.seqNum + 1
				state.seqNum += 1
				context.Send(context.Parent(), report)
			}

		case fbinance.ET_CANCELED:
			report, err := state.account.ConfirmCancelOrder(exec.ClientOrderID)
			if err != nil {
				return fmt.Errorf("error confirming cancel order: %v", err)
			}
			if report != nil {
				report.SeqNum = state.seqNum + 1
				state.seqNum += 1
				context.Send(context.Parent(), report)
			}
		}
	case fbinance.ACCOUNT_UPDATE:
		fmt.Println("ACCOUNT UPDATE", msg.Message.Account)

	case fbinance.MARGIN_CALL:
		// TODO
	default:
		return fmt.Errorf("received unknown event type: %s", msg.Message.Event)
	}

	return nil
}

func (state *AccountListener) subscribeAccount(context actor.Context) error {
	if state.ws != nil {
		_ = state.ws.Disconnect()
	}

	req, _, err := fbinance.GetListenKey(state.account.Credentials)
	if err != nil {
		return fmt.Errorf("error getting listen key request: %v", err)
	}

	listenKey := fbinance.ListenKey{}
	if err := utils.PerformRequest(state.client, req, &listenKey); err != nil {
		return fmt.Errorf("error getting listen key: %v", err)
	}
	if listenKey.Code != 0 {
		return fmt.Errorf(listenKey.Message)
	}

	ws := fbinance.NewAuthWebsocket(listenKey.ListenKey)
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to fbinance websocket: %v", err)
	}

	go func(ws *fbinance.AuthWebsocket, pid *actor.PID) {
		for ws.ReadMessage() {
			actor.EmptyRootContext.Send(pid, ws.Msg)
		}
	}(ws, context.Self())
	state.ws = ws

	return nil
}

func (state *AccountListener) checkSocket(context actor.Context) error {

	if time.Now().Sub(state.lastPingTime) > 5*time.Second {
		//_ = state.ws.Ping()
	}

	if state.ws.Err != nil || !state.ws.Connected {
		if state.ws.Err != nil {
			state.logger.Info("error on socket", log.Error(state.ws.Err))
		}
		if err := state.subscribeAccount(context); err != nil {
			return fmt.Errorf("error subscribing to account: %v", err)
		}
	}

	return nil
}

func (state *AccountListener) checkAccount(context actor.Context) error {
	fmt.Println("CHECKING ACCOUNT !")
	// Fetch balances
	res, err := context.RequestFuture(state.fbinanceExecutor, &messages.BalancesRequest{
		Account: state.account.Account,
	}, 10*time.Second).Result()

	if err != nil {
		return fmt.Errorf("error getting balances from executor: %v", err)
	}

	balanceList, ok := res.(*messages.BalanceList)
	if !ok {
		return fmt.Errorf("was expecting BalanceList, got %s", reflect.TypeOf(res).String())
	}

	if !balanceList.Success {
		return fmt.Errorf("error getting balances: %s", balanceList.RejectionReason.String())
	}

	if len(balanceList.Balances) != 1 {
		return fmt.Errorf("was expecting 1 balance, got %d", len(balanceList.Balances))
	}

	// Fetch positions
	res, err = context.RequestFuture(state.fbinanceExecutor, &messages.PositionsRequest{
		Instrument: nil,
		Account:    state.account.Account,
	}, 10*time.Second).Result()

	if err != nil {
		return fmt.Errorf("error getting positions from executor: %v", err)
	}

	positionList, ok := res.(*messages.PositionList)
	if !ok {
		return fmt.Errorf("was expecting PositionList, got %s", reflect.TypeOf(res).String())
	}

	if !positionList.Success {
		return fmt.Errorf("error getting positions: %s", positionList.RejectionReason.String())
	}

	rawMargin1 := int(math.Round(state.account.GetMargin() * state.account.MarginPrecision))
	rawMargin2 := int(math.Round(balanceList.Balances[0].Quantity * state.account.MarginPrecision))
	if rawMargin1 != rawMargin2 {
		return fmt.Errorf("different margin amount: %f %f", state.account.GetMargin(), balanceList.Balances[0].Quantity)
	}

	pos1 := state.account.GetPositions()

	var pos2 []*models.Position
	for _, p := range positionList.Positions {
		if p.Quantity > 0 {
			pos2 = append(pos2, p)
		}
	}
	if len(pos1) != len(pos2) {
		return fmt.Errorf("different number of positions")
	}

	// sort
	sort.Slice(pos1, func(i, j int) bool {
		return pos1[i].Instrument.SecurityID.Value < pos1[j].Instrument.SecurityID.Value
	})
	sort.Slice(pos2, func(i, j int) bool {
		return pos2[i].Instrument.SecurityID.Value < pos2[j].Instrument.SecurityID.Value
	})

	for i := range pos1 {
		if int(pos1[i].Quantity) != int(pos2[i].Quantity) {
			return fmt.Errorf("position have different quantity: %f %f", pos1[i].Quantity, pos2[i].Quantity)
		}
		rawCost1 := int(pos1[i].Cost * state.account.MarginPrecision)
		rawCost2 := int(pos2[i].Cost * state.account.MarginPrecision)
		if rawCost1 != rawCost2 {
			return fmt.Errorf("position have different cost: %f %f", pos1[i].Cost, pos2[i].Cost)
		}
	}
	return nil
}