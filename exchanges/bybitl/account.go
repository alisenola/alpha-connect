package bybitl

import (
	"encoding/json"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/account"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/bybitl"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"net"
	"net/http"
	"reflect"
	"sort"
	"time"
)

type checkSocket struct{}
type checkAccount struct{}

type AccountListener struct {
	account            *account.Account
	seqNum             uint64
	bybitlExecutor     *actor.PID
	ws                 *bybitl.Websocket
	logger             *log.Logger
	checkAccountTicker *time.Ticker
	checkSocketTicker  *time.Ticker
	lastPingTime       time.Time
	securities         map[uint64]*models.Security
	symbolToSec        map[string]*models.Security
	client             *http.Client
	txs                *mongo.Collection
	execs              *mongo.Collection
	//TODO Needed?
	//reconciler *actor.PID
}

func NewAccountListenerProducer(account *account.Account, txs, execs *mongo.Collection) actor.Producer {
	return func() actor.Actor {
		return NewAccountListener(account, txs, execs)
	}
}

func NewAccountListener(account *account.Account, txs, execs *mongo.Collection) actor.Actor {
	return &AccountListener{
		account: account,
		seqNum:  0,
		ws:      nil,
		logger:  nil,
		txs:     txs,
		execs:   execs,
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
		}
		state.logger.Info("actor restarting")
	case *messages.AccountDataRequest:
		if err := state.OnAccountDataRequest(context); err != nil {
			state.logger.Error("error processing OnAccountDataRequest", log.Error(err))
			panic(err)
		}
	case *messages.BalancesRequest:
		if err := state.OnBalancesRequest(context); err != nil {
			state.logger.Error("error processing OnBalancesRequest", log.Error(err))
			panic(err)
		}
	case *messages.PositionsRequest:
		if err := state.OnPositionsRequest(context); err != nil {
			state.logger.Error("error processing OnPositionListRequest", log.Error(err))
			panic(err)
		}
	case *messages.OrderStatusRequest:
		if err := state.OnOrderStatusRequest(context); err != nil {
			state.logger.Error("error processing OnOrderStatusRequset", log.Error(err))
			panic(err)
		}
	case *messages.NewOrderSingleRequest:
		if err := state.OnNewOrderSingleRequest(context); err != nil {
			state.logger.Error("error processing OnNewOrderSingleRequest", log.Error(err))
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
	case *messages.OrderReplaceRequest:
		if err := state.OnOrderReplaceRequest(context); err != nil {
			state.logger.Error("error processing OnOrderReplaceRequest", log.Error(err))
			panic(err)
		}
	case *xchanger.WebsocketMessage:
		if err := state.onWebsocketMessage(context); err != nil {
			state.logger.Error("error processing onWebsocketMessage", log.Error(err))
			panic(err)
		}
	case *checkSocket:
		if err := state.checkSocket(context); err != nil {
			state.logger.Error("error checking socket", log.Error(err))
			panic(err)
		}
	case *checkAccount:
		if err := state.checkAccount(context); err != nil {
			state.logger.Error("error checking account", log.Error(err))
			panic(err)
		}
	}
}

func (state *AccountListener) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()),
	)
	state.bybitlExecutor = actor.NewPID(context.ActorSystem().Address(), "executor/exchanges/"+constants.BYBITL.Name+"_executor")
	state.client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1024,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	//Fetch the securities
	fmt.Println("Fetching the securities")
	ex := actor.NewPID(context.ActorSystem().Address(), "executor")
	res, err := context.RequestFuture(ex, &messages.SecurityListRequest{}, 10*time.Second).Result()
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
	var filteredSecurities []*models.Security
	for _, s := range securityList.Securities {
		if s.Exchange.ID == state.account.Exchange.ID {
			filteredSecurities = append(filteredSecurities, s)
		}
	}
	state.symbolToSec = make(map[string]*models.Security)
	for _, sec := range filteredSecurities {
		state.symbolToSec[sec.Symbol] = sec
	}

	if err := state.subscribeAccount(context); err != nil {
		return fmt.Errorf("error subscribing to account: %v", err)
	}

	//Fetch the current balance
	fmt.Println("Fetching the balance")
	bal := context.RequestFuture(state.bybitlExecutor, &messages.BalancesRequest{
		Account: state.account.Account,
	}, 10*time.Second)
	pos := context.RequestFuture(state.bybitlExecutor, &messages.PositionsRequest{
		Instrument: nil,
		Account:    state.account.Account,
	}, 10*time.Second)

	res, err = bal.Result()
	if err != nil {
		return fmt.Errorf("error getting balances from executor: %v", err)
	}
	balances, ok := res.(*messages.BalanceList)
	if !ok {
		return fmt.Errorf("was expecting *messages.BalanceList, got %s", reflect.TypeOf(res).String())
	}
	if !balances.Success {
		return fmt.Errorf("error getting balances: %s", balances.RejectionReason.String())
	}

	//Fetch the positions
	fmt.Println("Fetching the positions")
	res, err = pos.Result()
	if err != nil {
		return fmt.Errorf("error getting positions from executor: %v", err)
	}
	positions, ok := res.(*messages.PositionList)
	if !ok {
		return fmt.Errorf("was expecting *messages.PositionsList, got %s", reflect.TypeOf(res).String())
	}
	if !positions.Success {
		return fmt.Errorf("error getting positions: %s", positions.RejectionReason.String())
	}

	//Fetch the orders
	fmt.Println("Fetching the orders")
	var ords []*models.Order
	securityMap := make(map[uint64]*models.Security)
	var futs []*actor.Future
	for _, fs := range filteredSecurities {
		f := context.RequestFuture(state.bybitlExecutor, &messages.OrderStatusRequest{
			Account: state.account.Account,
			Filter: &messages.OrderFilter{
				Instrument: &models.Instrument{
					SecurityID: &types.UInt64Value{Value: fs.SecurityID},
					Symbol:     &types.StringValue{Value: fs.Symbol},
					Exchange:   fs.Exchange,
				},
			},
		}, 10*time.Second)
		futs = append(futs, f)
		time.Sleep(100 * time.Millisecond)
		securityMap[fs.SecurityID] = fs
	}
	for _, f := range futs {
		res, err := f.Result()
		if err != nil {
			return fmt.Errorf("error getting orders from executor: %v", err)
		}
		orders, ok := res.(*messages.OrderList)
		if !ok {
			return fmt.Errorf("was expecting *messages.OrderList, got %s", reflect.TypeOf(res).String())
		}
		if !orders.Success {
			return fmt.Errorf("error getting orders: %s", orders.RejectionReason.String())
		}
		ords = append(ords, orders.Orders...)
	}
	state.securities = securityMap
	state.seqNum = 0

	//Sync account
	makerFee := 0.0001
	takerFee := 0.0006
	if err := state.account.Sync(filteredSecurities, ords, positions.Positions, balances.Balances, &makerFee, &takerFee); err != nil {
		return fmt.Errorf("error syncing account: %v", err)
	}

	checkAccountTicker := time.NewTicker(5 * time.Minute)
	state.checkAccountTicker = checkAccountTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-checkAccountTicker.C:
				context.Send(pid, &checkAccount{})
			case <-time.After(6 * time.Minute):
				// timer stopped, we leave
				return
			}
		}
	}(context.Self())

	checkSocketTicker := time.NewTicker(5 * time.Second)
	state.checkSocketTicker = checkSocketTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-checkSocketTicker.C:
				context.Send(pid, &checkSocket{})
			case <-time.After(10 * time.Second):
				// timer stopped, we leave
				return
			}
		}
	}(context.Self())

	return nil
}

func (state *AccountListener) OnAccountDataRequest(context actor.Context) error {
	req := context.Message().(*messages.AccountDataRequest)
	resp := &messages.AccountDataResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Securities: state.account.GetSecurities(),
		Orders:     state.account.GetOrders(nil),
		Positions:  state.account.GetPositions(),
		Balances:   state.account.GetBalances(),
		Success:    true,
		SeqNum:     state.seqNum,
	}
	makerFee := state.account.GetMakerFee()
	takerFee := state.account.GetTakerFee()

	if makerFee != nil {
		resp.MakerFee = &types.DoubleValue{Value: *makerFee}
	}
	if takerFee != nil {
		resp.TakerFee = &types.DoubleValue{Value: *takerFee}
	}
	context.Respond(resp)
	return nil
}

func (state *AccountListener) OnPositionsRequest(context actor.Context) error {
	req := context.Message().(*messages.PositionsRequest)
	// TODO FILTER
	context.Respond(&messages.PositionList{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Positions:  state.account.GetPositions(),
	})
	return nil
}

func (state *AccountListener) OnBalancesRequest(context actor.Context) error {
	req := context.Message().(*messages.BalancesRequest)
	// TODO FILTER
	balances := state.account.GetBalances()
	context.Respond(&messages.BalanceList{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Balances:   balances,
	})
	return nil
}

func (state *AccountListener) OnOrderStatusRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderStatusRequest)
	orders := state.account.GetOrders(req.Filter)
	context.Respond(&messages.OrderList{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Orders:     orders,
		Success:    true,
	})
	return nil
}

func (state *AccountListener) OnNewOrderSingleRequest(context actor.Context) error {
	req := context.Message().(*messages.NewOrderSingleRequest)
	req.Account = state.account.Account
	// Check order quantity
	order := &models.Order{
		OrderID:               "",
		ClientOrderID:         req.Order.ClientOrderID,
		Instrument:            req.Order.Instrument,
		OrderStatus:           models.PendingNew,
		OrderType:             req.Order.OrderType,
		Side:                  req.Order.OrderSide,
		TimeInForce:           req.Order.TimeInForce,
		LeavesQuantity:        req.Order.Quantity,
		Price:                 req.Order.Price,
		CumQuantity:           0,
		ExecutionInstructions: req.Order.ExecutionInstructions,
		Tag:                   req.Order.Tag,
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
				fut := context.RequestFuture(state.bybitlExecutor, req, 10*time.Second)
				context.AwaitFuture(fut, func(res interface{}, err error) {
					if err != nil {
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
			Success:         false,
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
				fut := context.RequestFuture(state.bybitlExecutor, req, 10*time.Second)
				context.AwaitFuture(fut, func(res interface{}, err error) {
					if err != nil {
						report, err := state.account.RejectCancelOrder(ID, messages.Other)
						if err != nil {
							panic(err)
						}
						if report != nil {
							report.SeqNum = state.seqNum + 1
							state.seqNum += 1
							context.Send(context.Parent(), report)
						}
						return
					}
					response := res.(*messages.OrderCancelResponse)

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
			RequestID:  req.RequestID,
			ResponseID: uint64(time.Now().UnixNano()),
			Success:    true,
		})
		return nil
	}
	var reports []*messages.ExecutionReport
	for _, o := range orders {
		if o.OrderStatus != models.New && o.OrderStatus != models.PartiallyFilled {
			continue
		}
		report, res := state.account.CancelOrder(o.OrderID)
		if res != nil {
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
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
	})

	for _, report := range reports {
		report.SeqNum = state.seqNum + 1
		state.seqNum += 1
		context.Send(context.Parent(), report)
	}
	fut := context.RequestFuture(state.bybitlExecutor, req, 10*time.Second)
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
				report, err := state.account.RejectCancelOrder(r.ClientOrderID.Value, response.RejectionReason)
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
	req := context.Message().(*messages.OrderReplaceRequest)
	response := &messages.OrderReplaceResponse{
		ResponseID: uint64(time.Now().UnixNano()),
		RequestID:  req.RequestID,
		Success:    false,
	}
	Id := ""
	if req.Update != nil {
		if req.Update.OrderID != nil {
			Id = req.Update.OrderID.Value
		} else if req.Update.OrigClientOrderID != nil {
			Id = req.Update.OrigClientOrderID.Value
		}
	}
	report, rej := state.account.ReplaceOrder(Id, req.Update.Price, req.Update.Quantity)
	if rej != nil {
		response.RejectionReason = *rej
		context.Respond(response)
		return nil
	}
	report.SeqNum = state.seqNum + 1
	state.seqNum += 1
	context.Send(context.Parent(), report)
	if report.OrderStatus == models.PendingReplace {
		fut := context.RequestFuture(state.bybitlExecutor, req, 10*time.Second)
		context.AwaitFuture(fut, func(res interface{}, err error) {
			if err != nil {
				report, err := state.account.RejectReplaceOrder(Id, messages.Other)
				if err != nil {
					panic(err)
				}
				response.RejectionReason = messages.Other
				context.Respond(response)
				if report != nil {
					report.SeqNum = state.seqNum + 1
					state.seqNum += 1
					context.Send(context.Parent(), report)
				}
				return
			}
			replaceResponse := res.(*messages.OrderReplaceResponse)
			context.Respond(replaceResponse)
			if replaceResponse.Success {
				report, err := state.account.ConfirmReplaceOrder(Id, replaceResponse.OrderID)
				if err != nil {
					panic(err)
				}
				if report != nil {
					report.SeqNum = state.seqNum + 1
					state.seqNum += 1
					context.Send(context.Parent(), report)
				}
			} else {
				report, err := state.account.RejectReplaceOrder(Id, messages.Other)
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
	return nil
}

func (state *AccountListener) subscribeAccount(context actor.Context) error {
	if state.ws != nil {
		_ = state.ws.Disconnect()
	}

	ws := bybitl.NewWebsocket()
	if err := ws.ConnectPrivate(&net.Dialer{}); err != nil {
		return fmt.Errorf("error connection to bybitl websocket: %v", err)
	}
	if err := ws.Authenticate(state.account.ApiCredentials); err != nil {
		return fmt.Errorf("error authenticating for bybitl websocket: %v", err)
	}
	// Subscribe to orders
	if err := ws.SubscribeOrders(); err != nil {
		return fmt.Errorf("error subscribing to orders: %v", err)
	}
	// Subscribe to executions
	if err := ws.SubscribeExecutions(); err != nil {
		return fmt.Errorf("error subscribing to executions: %v", err)
	}
	go func(pid *actor.PID, ws *bybitl.Websocket) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
		}
	}(context.Self(), ws)

	state.ws = ws
	return nil
}

func (state *AccountListener) onWebsocketMessage(context actor.Context) error {
	msg := context.Message().(*xchanger.WebsocketMessage)
	if state.ws == nil || msg.WSID != state.ws.ID {
		return nil
	}

	if msg.Message == nil {
		return fmt.Errorf("reveived nil message")
	}
	b, _ := json.Marshal(msg.Message)
	fmt.Println(string(b))
	switch s := msg.Message.(type) {
	case bybitl.WSOrders:
		for _, order := range s {
			switch order.OrderStatus {
			case bybitl.OrderNew:
				// New Order
				if !state.account.HasOrder(order.OrderLinkId) {
					o := wsOrderToModel(&order)
					o.OrderStatus = models.PendingNew
					_, rej := state.account.NewOrder(o)
					if rej != nil {
						return fmt.Errorf("error creating new order: %s", rej.String())
					}
				}
				report, err := state.account.ConfirmNewOrder(order.OrderLinkId, order.OrderId)
				if err != nil {
					return fmt.Errorf("error confirming new order: %v", err)
				}
				if report != nil {
					report.SeqNum = state.seqNum + 1
					state.seqNum += 1
					context.Send(context.Parent(), report)
				}
			case bybitl.OrderCancelled:
				report, err := state.account.ConfirmCancelOrder(order.OrderLinkId)
				if err != nil {
					return fmt.Errorf("error confirming cancel order: %v", err)
				}
				if report != nil {
					report.SeqNum = state.seqNum + 1
					state.seqNum += 1
					context.Send(context.Parent(), report)
				}
			case bybitl.OrderFilled:
				ord, err := state.account.GetOrder(order.OrderLinkId)
				if err != nil {
					return fmt.Errorf("error getting order: %v", err)
				}
				// If instantly filled, we will not receive an update with a bybitl.OrderNew status
				// therefore we won't go through the branch above that confirms the new order.
				// Therefore, we check if it's pendingNew, and we confirm it here. We let
				// the confirmFill to the bybitl.WSExecutions
				if ord.OrderStatus == models.PendingNew {
					report, err := state.account.ConfirmNewOrder(order.OrderLinkId, order.OrderId)
					if err != nil {
						return fmt.Errorf("error confirming filled order: %v", err)
					}
					if report != nil {
						report.SeqNum = state.seqNum + 1
						state.seqNum += 1
						context.Send(context.Parent(), report)
					}
				}
			}
		}
	case bybitl.WSExecutions:
		for _, exec := range s {
			switch exec.ExecType {
			case "Trade":
				report, err := state.account.ConfirmFill(exec.OrderId, exec.ExecId, exec.Price, exec.ExecQty, !exec.IsMaker)
				if err != nil {
					return fmt.Errorf("error confirming filled order: %v", err)
				}
				if report != nil {
					report.SeqNum = state.seqNum + 1
					state.seqNum += 1
					context.Send(context.Parent(), report)
				}
			}
		}
	case bybitl.WSResponse:
		if !s.Success {
			return fmt.Errorf("error in WSResponse: %s", s.ReturnMessage)
		}
		return nil
	default:
		return fmt.Errorf("received unknown event type: %s", reflect.TypeOf(s).String())
	}
	return nil
}

func (state *AccountListener) Clean(context actor.Context) error {
	if state.ws != nil {
		if err := state.ws.Disconnect(); err != nil {
			state.logger.Info("error disconnecting websocket", log.Error(err))
		}
	}

	if state.checkAccountTicker != nil {
		state.checkAccountTicker.Stop()
		state.checkAccountTicker = nil
	}

	if state.checkSocketTicker != nil {
		state.checkSocketTicker.Stop()
		state.checkSocketTicker = nil
	}

	return nil
}

func (state *AccountListener) checkSocket(context actor.Context) error {
	if time.Since(state.lastPingTime) > 5*time.Second {
		_ = state.ws.Ping()
		state.lastPingTime = time.Now()
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
	fmt.Println("CHECKING ACCOUNT")

	// Fetch balances
	resp, err := context.RequestFuture(state.bybitlExecutor, &messages.BalancesRequest{
		Account: state.account.Account,
	}, 10*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error getting balances from executor: %v", err)
	}

	accntBalances := state.account.GetBalances()
	execBalanceList, ok := resp.(*messages.BalanceList)
	if !ok {
		return fmt.Errorf("was expecting *messages.BalanceList, got %s", reflect.TypeOf(resp).String())
	}
	if !execBalanceList.Success {
		return fmt.Errorf("error getting balances: %s", execBalanceList.RejectionReason.String())
	}
	if len(execBalanceList.Balances) != len(accntBalances) {
		return fmt.Errorf("was expecting %d balance, got %d", len(execBalanceList.Balances), len(accntBalances))
	}

	// Fetch positions
	resp, err = context.RequestFuture(state.bybitlExecutor, &messages.PositionsRequest{
		Instrument: nil,
		Account:    state.account.Account,
	}, 10*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error getting balances from executor: %v", err)
	}

	positionList, ok := resp.(*messages.PositionList)
	if !ok {
		return fmt.Errorf("was expecting *messages.PositionList, got %s", reflect.TypeOf(resp).String())
	}
	if !positionList.Success {
		return fmt.Errorf("error getting balances: %s", execBalanceList.RejectionReason.String())
	}
	execBalances := execBalanceList.Balances
	sort.Slice(accntBalances, func(i, j int) bool {
		return accntBalances[i].Asset.ID < accntBalances[j].Asset.ID
	})
	sort.Slice(execBalances, func(i, j int) bool {
		return execBalances[i].Asset.ID < execBalances[j].Asset.ID
	})
	for i, b1 := range execBalances {
		b2 := accntBalances[i]
		//rawB1 := int(math.Round(b1.Quantity * state.account.MarginPrecision))
		//rawB2 := int(math.Round(b2.Quantity * state.account.MarginPrecision))
		diff := math.Abs(1 - b1.Quantity/b2.Quantity)
		if diff > 0.01 {
			return fmt.Errorf("different margin amount: %f %f", b1.Quantity, b2.Quantity)
		}
	}

	pos1 := state.account.GetPositions()

	var pos2 []*models.Position
	for _, p := range positionList.Positions {
		if p.Quantity != 0 {
			pos2 = append(pos2, p)
		}
	}

	if len(pos1) != len(pos2) {
		return fmt.Errorf("different number of positions: %d vs %d", len(pos1), len(pos2))
	}

	sort.Slice(pos1, func(i, j int) bool {
		return pos1[i].Instrument.SecurityID.Value < pos1[j].Instrument.SecurityID.Value
	})
	sort.Slice(pos2, func(i, j int) bool {
		return pos2[i].Instrument.SecurityID.Value < pos2[j].Instrument.SecurityID.Value
	})

	for i := range pos1 {
		lp := math.Ceil(1. / state.securities[pos1[i].Instrument.SecurityID.Value].RoundLot.Value)
		if int(math.Round(pos1[i].Quantity*lp)) != int(math.Round(pos2[i].Quantity*lp)) {
			return fmt.Errorf("positions have different quantities: %f vs %f", pos1[i].Quantity, pos2[i].Quantity)
		}
		rawCost1 := int(math.Round(pos1[i].Cost * state.account.MarginPrecision))
		rawCost2 := int(math.Round(pos2[i].Cost * state.account.MarginPrecision))
		if rawCost1 != rawCost2 {
			return fmt.Errorf("positions have different costs: %f vs %f", pos1[i].Cost, pos2[i].Cost)
		}
	}

	return nil
}
