package bitmex

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/account"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/bitmex"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
	"net"
	"reflect"
	"sort"
	"time"
)

type checkSocket struct{}
type checkAccount struct{}

type AccountListener struct {
	account         *account.Account
	seqNum          uint64
	bitmexExecutor  *actor.PID
	ws              *bitmex.Websocket
	executorManager *actor.PID
	logger          *log.Logger
	lastPingTime    time.Time
	securities      []*models.Security
	socketTicker    *time.Ticker
	accountTicker   *time.Ticker
}

func NewAccountListenerProducer(account *account.Account) actor.Producer {
	return func() actor.Actor {
		return NewAccountListener(account)
	}
}

func NewAccountListener(account *account.Account) actor.Actor {
	return &AccountListener{
		account:         account,
		seqNum:          0,
		ws:              nil,
		executorManager: nil,
		logger:          nil,
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

	case *messages.AccountDataRequest:
		if err := state.OnAccountDataRequest(context); err != nil {
			state.logger.Error("error processing OnAccountDataRequest", log.Error(err))
			panic(err)
		}

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
		if err := state.NewOrderSingleRequest(context); err != nil {
			state.logger.Error("error processing NewOrderSingleRequest", log.Error(err))
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

	case *xchanger.WebsocketMessage:
		if err := state.onWebsocketMessage(context); err != nil {
			state.logger.Error("error processing onWebocketMessage", log.Error(err))
			panic(err)
		}

	case *checkSocket:
		fmt.Println("CHECK SOCKET")
		if err := state.checkSocket(context); err != nil {
			state.logger.Error("error checking socket", log.Error(err))
			panic(err)
		}

	case *checkAccount:
		fmt.Println("CHEKC ACCOUNT")
		if err := state.checkAccount(context); err != nil {
			state.logger.Error("error checking socket", log.Error(err))
			panic(err)
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
	state.bitmexExecutor = actor.NewPID(context.ActorSystem().Address(), "executor/exchanges/"+constants.BITMEX.Name+"_executor")

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
	state.securities = filteredSecurities

	if err := state.Sync(context); err != nil {
		return fmt.Errorf("error syncing account: %v", err)
	}

	socketTicker := time.NewTicker(5 * time.Second)
	state.socketTicker = socketTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-socketTicker.C:
				fmt.Println("SENDING CHECK SOCKET")
				context.Send(pid, &checkSocket{})
			case <-time.After(10 * time.Second):
				if state.socketTicker != socketTicker {
					return
				}
			}
		}
	}(context.Self())

	accountTicker := time.NewTicker(10 * time.Minute)
	state.accountTicker = accountTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-accountTicker.C:
				fmt.Println("SENDING CHECK ACCOUUNT")
				context.Send(pid, &checkAccount{})
			case <-time.After(11 * time.Minute):
				// timer stopped, we leave
				if state.accountTicker != accountTicker {
					return
				}
			}
		}
	}(context.Self())

	return nil
}

func (state *AccountListener) Sync(context actor.Context) error {
	if err := state.subscribeAccount(context); err != nil {
		return fmt.Errorf("error subscribing to account: %v", err)
	}

	// Then fetch balances
	res, err := context.RequestFuture(state.bitmexExecutor, &messages.BalancesRequest{
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

	// Then fetch positions
	res, err = context.RequestFuture(state.bitmexExecutor, &messages.PositionsRequest{
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
	res, err = context.RequestFuture(state.bitmexExecutor, &messages.OrderStatusRequest{
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

	// Sync account
	if err := state.account.Sync(state.securities, orderList.Orders, positionList.Positions, balanceList.Balances, nil, nil); err != nil {
		return fmt.Errorf("error syncing account: %v", err)
	}

	state.seqNum = 0

	return nil
}

// TODO
func (state *AccountListener) Clean(context actor.Context) error {
	if state.ws != nil {
		if err := state.ws.Disconnect(); err != nil {
			state.logger.Info("error disconnecting socket", log.Error(err))
		}
	}
	if state.socketTicker != nil {
		state.socketTicker.Stop()
		state.socketTicker = nil
	}
	if state.accountTicker != nil {
		state.accountTicker.Stop()
		state.accountTicker = nil
	}

	return nil
}

func (state *AccountListener) OnAccountDataRequest(context actor.Context) error {
	msg := context.Message().(*messages.AccountDataRequest)
	res := &messages.AccountDataResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: state.account.GetSecurities(),
		Orders:     state.account.GetOrders(nil),
		Positions:  state.account.GetPositions(),
		Balances:   state.account.GetBalances(),
		SeqNum:     state.seqNum,
	}

	makerFee := state.account.GetMakerFee()
	takerFee := state.account.GetTakerFee()
	if makerFee != nil {
		res.MakerFee = &wrapperspb.DoubleValue{Value: *makerFee}
	}
	if takerFee != nil {
		res.TakerFee = &wrapperspb.DoubleValue{Value: *takerFee}
	}

	context.Respond(res)

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
	req := context.Message().(*messages.OrderStatusRequest)
	orders := state.account.GetOrders(req.Filter)
	context.Respond(&messages.OrderList{
		RequestID: req.RequestID,
		Success:   true,
		Orders:    orders,
	})
	return nil
}

func (state *AccountListener) NewOrderSingleRequest(context actor.Context) error {
	req := context.Message().(*messages.NewOrderSingleRequest)
	if req.Expire != nil && req.Expire.AsTime().Before(time.Now()) {
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID:       req.RequestID,
			Success:         false,
			RejectionReason: messages.RejectionReason_RequestExpired,
		})
		return nil
	}
	req.Account = state.account.Account
	// Check order quantity
	order := &models.Order{
		OrderID:               "",
		ClientOrderID:         req.Order.ClientOrderID,
		Instrument:            req.Order.Instrument,
		OrderStatus:           models.OrderStatus_PendingNew,
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
			if report.ExecutionType == messages.ExecutionType_PendingNew {
				fut := context.RequestFuture(state.bitmexExecutor, req, 10*time.Second)
				context.ReenterAfter(fut, func(res interface{}, err error) {
					if err != nil {
						report, err := state.account.RejectNewOrder(order.ClientOrderID, messages.RejectionReason_Other)
						if err != nil {
							panic(err)
						}
						context.Respond(&messages.NewOrderSingleResponse{
							RequestID:       req.RequestID,
							Success:         false,
							RejectionReason: messages.RejectionReason_Other,
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
						nReport, _ := state.account.ConfirmNewOrder(order.ClientOrderID, response.OrderID, nil)
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
			response.RejectionReason = messages.RejectionReason_UnknownSymbol
			context.Respond(response)
			return nil
		} else if symbol != order.Instrument.Symbol.Value {
			response.RejectionReason = messages.RejectionReason_DifferentSymbols
			context.Respond(response)
			return nil
		}
	}
	for _, reqOrder := range req.Orders {
		order := &models.Order{
			OrderID:               "",
			ClientOrderID:         reqOrder.ClientOrderID,
			Instrument:            reqOrder.Instrument,
			OrderStatus:           models.OrderStatus_PendingNew,
			OrderType:             reqOrder.OrderType,
			Side:                  reqOrder.OrderSide,
			TimeInForce:           reqOrder.TimeInForce,
			LeavesQuantity:        reqOrder.Quantity,
			Price:                 reqOrder.Price,
			CumQuantity:           0,
			ExecutionInstructions: reqOrder.ExecutionInstructions,
			Tag:                   reqOrder.Tag,
		}
		report, res := state.account.NewOrder(order)
		if res != nil {
			// Cancel all new order up until now
			for _, r := range reports {
				_, err := state.account.RejectNewOrder(r.ClientOrderID.Value, messages.RejectionReason_Other)
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
	fut := context.RequestFuture(state.bitmexExecutor, req, 10*time.Second)
	context.ReenterAfter(fut, func(res interface{}, err error) {
		if err != nil {
			for _, r := range reports {
				report, err := state.account.RejectNewOrder(r.ClientOrderID.Value, messages.RejectionReason_Other)
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
				RejectionReason: messages.RejectionReason_Other,
			})
			return
		}
		response := res.(*messages.NewOrderBulkResponse)
		context.Respond(response)
		if response.Success {
			for i, r := range reports {
				report, err := state.account.ConfirmNewOrder(r.ClientOrderID.Value, response.OrderIDs[i], nil)
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
				report, err := state.account.RejectNewOrder(r.ClientOrderID.Value, messages.RejectionReason_Other)
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
	var ID string
	if req.Update.OrigClientOrderID != nil {
		ID = req.Update.OrigClientOrderID.Value
	} else if req.Update.OrderID != nil {
		ID = req.Update.OrderID.Value
	}
	report, res := state.account.ReplaceOrder(ID, req.Update.Price, req.Update.Quantity)
	if res != nil {
		context.Respond(&messages.OrderReplaceResponse{
			RequestID:       req.RequestID,
			RejectionReason: *res,
		})
	} else {
		context.Respond(&messages.OrderReplaceResponse{
			RequestID: req.RequestID,
			Success:   true,
		})
		if report != nil {
			report.SeqNum = state.seqNum + 1
			state.seqNum += 1
			context.Send(context.Parent(), report)
			if report.ExecutionType == messages.ExecutionType_PendingReplace {
				fut := context.RequestFuture(state.bitmexExecutor, req, 10*time.Second)
				context.ReenterAfter(fut, func(res interface{}, err error) {
					if err != nil {
						report, err := state.account.RejectReplaceOrder(ID, messages.RejectionReason_Other)
						if err != nil {
							panic(err)
						}
						context.Respond(&messages.OrderReplaceResponse{
							RequestID:       req.RequestID,
							Success:         false,
							RejectionReason: messages.RejectionReason_Other,
						})
						if report != nil {
							report.SeqNum = state.seqNum + 1
							state.seqNum += 1
							context.Send(context.Parent(), report)
						}
						return
					}
					response := res.(*messages.OrderReplaceResponse)
					context.Respond(response)

					if response.Success {
						// TODO for now let the WS do the job
						/*
							report, err := state.account.ConfirmReplaceOrder(ID)
							if err != nil {
								panic(err)
							}
							if report != nil {
								report.SeqNum = state.seqNum + 1
								state.seqNum += 1
								context.Send(context.Parent(), report)
							}
						*/
					} else {
						report, err := state.account.RejectReplaceOrder(ID, response.RejectionReason)
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

func (state *AccountListener) OnBulkOrderReplaceRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderBulkReplaceRequest)
	var reports []*messages.ExecutionReport
	for _, u := range req.Updates {
		var ID string
		if u.OrigClientOrderID != nil {
			ID = u.OrigClientOrderID.Value
		} else if u.OrderID != nil {
			ID = u.OrderID.Value
		}
		report, res := state.account.ReplaceOrder(ID, u.Price, u.Quantity)
		if res != nil {
			// Reject all cancel order up until now
			for _, r := range reports {
				_, err := state.account.RejectReplaceOrder(r.ClientOrderID.Value, messages.RejectionReason_Other)
				if err != nil {
					return err
				}
			}

			context.Respond(&messages.OrderBulkReplaceResponse{
				RequestID:       req.RequestID,
				Success:         false,
				RejectionReason: *res,
			})

			return nil
		} else if report != nil {
			reports = append(reports, report)
		}
	}

	context.Respond(&messages.OrderBulkReplaceResponse{
		RequestID: req.RequestID,
		Success:   true,
	})

	for _, report := range reports {
		report.SeqNum = state.seqNum + 1
		state.seqNum += 1
		context.Send(context.Parent(), report)
	}
	fut := context.RequestFuture(state.bitmexExecutor, req, 10*time.Second)
	context.ReenterAfter(fut, func(res interface{}, err error) {
		if err != nil {
			for _, r := range reports {
				report, err := state.account.RejectReplaceOrder(r.ClientOrderID.Value, messages.RejectionReason_Other)
				if err != nil {
					panic(err)
				}
				if report != nil {
					report.SeqNum = state.seqNum + 1
					state.seqNum += 1
					context.Send(context.Parent(), report)
				}
			}
			context.Respond(&messages.OrderBulkReplaceResponse{
				RequestID:       req.RequestID,
				Success:         false,
				RejectionReason: messages.RejectionReason_Other,
			})

			return
		}
		response := res.(*messages.OrderBulkReplaceResponse)
		context.Respond(response)

		if response.Success {
			// TODO for now let the WS do the job
			/*
				for _, r := range reports {
					report, err := state.account.ConfirmReplaceOrder(r.ClientOrderID.Value)
					if err != nil {
						panic(err)
					}
					if report != nil {
						report.SeqNum = state.seqNum + 1
						state.seqNum += 1
						context.Send(context.Parent(), report)
					}
				}
			*/
		} else {
			for _, r := range reports {
				report, err := state.account.RejectReplaceOrder(r.ClientOrderID.Value, messages.RejectionReason_Other)
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
			if report.ExecutionType == messages.ExecutionType_PendingCancel {
				fut := context.RequestFuture(state.bitmexExecutor, req, 10*time.Second)
				context.ReenterAfter(fut, func(res interface{}, err error) {
					if err != nil {
						report, err := state.account.RejectCancelOrder(ID, messages.RejectionReason_Other)
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
						if response.RejectionReason == messages.RejectionReason_UnknownOrder {
							panic(fmt.Errorf("unknown order rejection"))
						}
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
		if o.OrderStatus != models.OrderStatus_New && o.OrderStatus != models.OrderStatus_PartiallyFilled {
			continue
		}
		report, res := state.account.CancelOrder(o.ClientOrderID)
		if res != nil {
			// Reject all cancel order up until now
			for _, r := range reports {
				_, err := state.account.RejectCancelOrder(r.ClientOrderID.Value, messages.RejectionReason_Other)
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
	fut := context.RequestFuture(state.bitmexExecutor, req, 10*time.Second)
	context.ReenterAfter(fut, func(res interface{}, err error) {
		if err != nil {
			for _, r := range reports {
				report, err := state.account.RejectCancelOrder(r.ClientOrderID.Value, messages.RejectionReason_Other)
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
				RejectionReason: messages.RejectionReason_Other,
			})

			return
		}
		response := res.(*messages.OrderMassCancelResponse)
		context.Respond(response)

		/*
			if response.Success {
				for _, r := range reports {
					report, err := state.account.ConfirmCancelOrder(r.ClientOrderID.Value)
					if err != nil {
						panic(err)
					}
					if report != nil {
						report.SeqNum = state.seqNum + 1
						state.seqNum += 1
						context.Send(context.Parent(), report)
					}
				}
		*/
		if !response.Success {
			for _, r := range reports {
				report, err := state.account.RejectCancelOrder(r.ClientOrderID.Value, messages.RejectionReason_Other)
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

func (state *AccountListener) onWebsocketMessage(context actor.Context) error {
	state.lastPingTime = time.Now()
	msg := context.Message().(*xchanger.WebsocketMessage)
	switch msg.Message.(type) {
	case error:
		return fmt.Errorf("socket error: %v", msg)

	case bitmex.WSExecutionData:
		execData := msg.Message.(bitmex.WSExecutionData)
		if err := state.onWSExecutionData(context, execData); err != nil {
			return err
		}
	}

	return nil
}

func (state *AccountListener) onWSExecutionData(context actor.Context, executionData bitmex.WSExecutionData) error {
	// Sort data by event time
	sort.Slice(executionData.Data, func(i, j int) bool {
		return executionData.Data[i].TransactTime.Before(executionData.Data[j].TransactTime)
	})
	for _, data := range executionData.Data {
		b, _ := json.Marshal(data)
		fmt.Println(string(b))
		switch data.ExecType {
		case "New":
			// New order
			if data.ClOrdID == nil {
				return fmt.Errorf("got an order with nil ClOrdID")
			}
			report, err := state.account.ConfirmNewOrder(*data.ClOrdID, data.OrderID, nil)
			if err != nil {
				return fmt.Errorf("error confirming new order: %v", err)
			}
			if report != nil {
				report.SeqNum = state.seqNum + 1
				state.seqNum += 1
				context.Send(context.Parent(), report)
			}

		case "Canceled":
			report, err := state.account.ConfirmCancelOrder(*data.ClOrdID)
			if err != nil {
				return fmt.Errorf("error confirming cancel order: %v", err)
			}
			if report != nil {
				report.SeqNum = state.seqNum + 1
				state.seqNum += 1
				context.Send(context.Parent(), report)
			}

		case "Rejected":
			report, err := state.account.RejectNewOrder(*data.ClOrdID, messages.RejectionReason_Other)
			if err != nil {
				return fmt.Errorf("error rejecting new order: %v", err)
			}
			if report != nil {
				report.SeqNum = state.seqNum + 1
				state.seqNum += 1
				context.Send(context.Parent(), report)
			}

		case "Replaced":
			report, err := state.account.ConfirmReplaceOrder(*data.ClOrdID, "")
			if err != nil {
				if err == account.ErrNotPendingReplace {
					return fmt.Errorf("error not pending replace: %v", err)
				}
				return fmt.Errorf("error confirming replace order: %v", err)
			}
			if report != nil {
				report.SeqNum = state.seqNum + 1
				state.seqNum += 1
				context.Send(context.Parent(), report)
			}

		case "Trade":
			report, err := state.account.ConfirmFill(*data.ClOrdID, *data.TrdMatchID, *data.LastPx, float64(*data.LastQty), *data.ExecComm > 0)
			if err != nil {
				return fmt.Errorf("error confirming fill: %v", err)
			}
			if report != nil {
				report.SeqNum = state.seqNum + 1
				state.seqNum += 1
				context.Send(context.Parent(), report)
			}

		case "Funding":
			/*
				{
				"execID":"d7ba8934-366f-0ddb-31b8-f745db13e01d",
				"orderID":"00000000-0000-0000-0000-000000000000",
				"clOrdID":"",
				"account":1502932,
				"symbol":"XBTUSD","
				side":"",
				"lastQty":169,
				"lastPx":9268.18,
				"underlyingLastPx":null,
				"lastMkt":"XBME",
				"lastLiquidityInd":"",
				"simpleOrderQty":null,
				"orderQty":169,
				"price":9268.18,
				"displayQty":null,
				"stopPx":null,
				"pegOffsetValue":null,
				"pegPriceType":"",
				"currency":"USD",
				"settlCurrency":"XBt",
				"execType":"Funding",
				"ordType":"Limit",
				"timeInForce":"AtTheClose",
				"execInst":"",
				"contingencyType":"",
				"exDestination":"XBME",
				"ordStatus":"Filled",
				"triggered":"",
				"workingIndicator":false,
				"ordRejReason":"",
				"simpleLeavesQty":null,
				"leavesQty":0,
				"simpleCumQty":null,
				"cumQty":169,
				"avgPx":9268.18,
				"commission":-0.0001,
				"tradePublishIndicator":"",
				"multiLegReportingType":"SingleSecurity",
				"text":"Funding",
				"trdMatchID":"5febca48-cb97-a9b3-1532-60e9418eb81d",
				"execCost":1823510,
				"execComm":-182,
				"homeNotional":-0.0182351,
				"foreignNotional":169,
				"transactTime":"2020-07-14T20:00:00Z",
				"timestamp":"2020-07-14T20:00:00.002Z"}

			*/
		default:
			return fmt.Errorf("got unknown exec type: %s", data.ExecType)
		}
	}

	return nil
}

func (state *AccountListener) subscribeAccount(context actor.Context) error {
	if state.ws != nil && state.ws.Err == nil && state.ws.Connected {
		// Skip if socket ok
		return nil
	}

	ws := bitmex.NewWebsocket()
	// TODO Dialer
	if err := ws.Connect(&net.Dialer{}); err != nil {
		return fmt.Errorf("error connecting to bitmex websocket: %v", err)
	}

	if err := ws.Auth(state.account.ApiCredentials); err != nil {
		return fmt.Errorf("error sending auth request: %v", err)
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	receivedMessage, ok := ws.Msg.Message.(bitmex.WSResponse)
	if !ok {
		errorMessage, ok := ws.Msg.Message.(bitmex.WSErrorResponse)
		if ok {
			return fmt.Errorf("error auth: %s", errorMessage.Error)
		}
		return fmt.Errorf("error casting message to WSResponse")
	}

	if !receivedMessage.Success {
		return fmt.Errorf("auth unsuccessful")
	}

	if err := ws.Subscribe(bitmex.WSExecutionStreamName); err != nil {
		return fmt.Errorf("error sending subscription request: %v", err)
	}
	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	subResponse, ok := ws.Msg.Message.(bitmex.WSSubscribeResponse)
	if !ok {
		errorMessage, ok := ws.Msg.Message.(bitmex.WSErrorResponse)
		if ok {
			return fmt.Errorf("error auth: %s", errorMessage.Error)
		}
		return fmt.Errorf("error casting message to WSSubscribeResponse")
	}
	if !subResponse.Success {
		return fmt.Errorf("subscription unsucessful")
	}

	go func(ws *bitmex.Websocket, pid *actor.PID) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
		}
	}(ws, context.Self())
	state.ws = ws

	return nil
}

func (state *AccountListener) checkSocket(context actor.Context) error {

	if time.Since(state.lastPingTime) > 5*time.Second {
		_ = state.ws.Ping()
	}

	if state.ws.Err != nil || !state.ws.Connected {
		if state.ws.Err != nil {
			state.logger.Info("error on socket", log.Error(state.ws.Err))
		}
		if err := state.Sync(context); err != nil {
			return fmt.Errorf("error syncing account: %v", err)
		}
	}

	return nil
}

func (state *AccountListener) checkAccount(context actor.Context) error {
	// Fetch balances
	res, err := context.RequestFuture(state.bitmexExecutor, &messages.BalancesRequest{
		Account: state.account.Account,
	}, 10*time.Second).Result()

	if err != nil {
		state.logger.Info("error getting balances from executor", log.Error(err))
		return nil
	}

	balanceList, ok := res.(*messages.BalanceList)
	if !ok {
		err := fmt.Errorf("was expecting BalanceList, got %s", reflect.TypeOf(res).String())
		state.logger.Info("error getting balances from executor", log.Error(err))
		return nil
	}

	if !balanceList.Success {
		state.logger.Info("error getting balances from executor", log.Error(errors.New(balanceList.RejectionReason.String())))
		return nil
	}

	if len(balanceList.Balances) != 1 {
		state.logger.Info("error getting balances from executor", log.Error(errors.New("was expecting one balance")))
		return nil
	}

	// Fetch positions
	res, err = context.RequestFuture(state.bitmexExecutor, &messages.PositionsRequest{
		Instrument: nil,
		Account:    state.account.Account,
	}, 10*time.Second).Result()

	if err != nil {
		state.logger.Info("error getting position from executor", log.Error(err))
		return nil
	}

	positionList, ok := res.(*messages.PositionList)
	if !ok {
		err := fmt.Errorf("was expecting PositionList, got %s", reflect.TypeOf(res).String())
		state.logger.Info("error getting position from executor", log.Error(err))
		return nil
	}

	if !positionList.Success {
		err := fmt.Errorf("%s", positionList.RejectionReason.String())
		state.logger.Info("error getting position from executor", log.Error(err))
		return nil
	}

	rawMargin1 := int(math.Round(state.account.GetMargin(nil) * state.account.MarginPrecision))
	rawMargin2 := int(math.Round(balanceList.Balances[0].Quantity * state.account.MarginPrecision))
	if rawMargin1 != rawMargin2 {
		err := fmt.Errorf("got different margin: %f %f", state.account.GetMargin(nil), balanceList.Balances[0].Quantity)
		state.logger.Info("re-syncing", log.Error(err))
		return state.Sync(context)
	}

	pos1 := state.account.GetPositions()

	var pos2 []*models.Position
	for _, p := range positionList.Positions {
		if math.Abs(p.Quantity) > 0 {
			pos2 = append(pos2, p)
		}
	}
	if len(pos1) != len(pos2) {
		// Re-sync
		err := fmt.Errorf("got different position number")
		state.logger.Info("re-syncing", log.Error(err))
		return state.Sync(context)
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
			// Re-sync
			err := fmt.Errorf("got different position quantity: %f %f", pos1[i].Quantity, pos2[i].Quantity)
			state.logger.Info("re-syncing", log.Error(err))
			return state.Sync(context)
		}
		rawCost1 := int(pos1[i].Cost * state.account.MarginPrecision)
		rawCost2 := int(pos2[i].Cost * state.account.MarginPrecision)
		if rawCost1 != rawCost2 {
			// Re-sync
			err := fmt.Errorf("got different position cost: %f %f", pos1[i].Cost, pos2[i].Cost)
			state.logger.Info("re-syncing", log.Error(err))
			return state.Sync(context)
		}
	}
	return nil
}
