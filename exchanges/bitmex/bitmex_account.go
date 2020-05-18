package bitmex

import (
	"errors"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/account"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/bitmex"
	"reflect"
	"sort"
	"time"
)

type AccountListener struct {
	account         *account.Account
	accountM        *models.Account
	seqNum          uint64
	bitmexExecutor  *actor.PID
	ws              *bitmex.Websocket
	executorManager *actor.PID
	logger          *log.Logger
}

func NewAccountListenerProducer(account *models.Account) actor.Producer {
	return func() actor.Actor {
		return NewAccountListener(account)
	}
}

func NewAccountListener(account *models.Account) actor.Actor {
	return &AccountListener{
		account:         nil,
		accountM:        account,
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

	case *messages.NewOrderSingle:
		if err := state.OnNewOrderSingle(context); err != nil {
			state.logger.Error("error processing OnNewOrderSingle", log.Error(err))
			panic(err)
		}

	case *messages.ExecutionReport:
		if err := state.OnExecutionReport(context); err != nil {
			state.logger.Error("error processing OnExecutionReport", log.Error(err))
			panic(err)
		}

	case *bitmex.WebsocketMessage:
		if err := state.onWebsocketMessage(context); err != nil {
			state.logger.Error("error processing onWebocketMessage", log.Error(err))
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
	state.bitmexExecutor = actor.NewLocalPID("executor/" + constants.BITMEX.Name + "_executor")

	if err := state.subscribeAccount(context); err != nil {
		return fmt.Errorf("error subscribing to order book: %v", err)
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
	if securityList.Error != "" {
		return fmt.Errorf("error getting securities: %s", securityList.Error)
	}

	// Instantiate account
	state.account = account.NewAccount(state.accountM.AccountID, securityList.Securities)

	// Then fetch positions
	res, err = context.RequestFuture(state.bitmexExecutor, &messages.PositionsRequest{
		Instrument: nil,
		Account:    state.accountM,
	}, 10*time.Second).Result()

	if err != nil {
		return fmt.Errorf("error getting orders from executor: %v", err)
	}

	positionList, ok := res.(*messages.PositionList)
	if !ok {
		return fmt.Errorf("was expecting PositionList, got %s", reflect.TypeOf(res).String())
	}

	if positionList.Error != "" {
		return errors.New(positionList.Error)
	}

	// Then fetch orders
	res, err = context.RequestFuture(state.bitmexExecutor, &messages.OrderStatusRequest{
		OrderID:       nil,
		ClientOrderID: nil,
		Instrument:    nil,
		Account:       state.accountM,
	}, 10*time.Second).Result()

	if err != nil {
		return fmt.Errorf("error getting orders from executor: %v", err)
	}

	orderList, ok := res.(*messages.OrderList)
	if !ok {
		return fmt.Errorf("was expecting OrderList, got %s", reflect.TypeOf(res).String())
	}

	if orderList.Error != "" {
		return errors.New(orderList.Error)
	}

	// Sync account
	if err := state.account.Sync(orderList.Orders, positionList.Positions); err != nil {
		return fmt.Errorf("error syncing account: %v", err)
	}
	state.seqNum = 0

	context.Send(context.Self(), &readSocket{})

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
	positions := state.account.GetPositions()
	context.Respond(&messages.PositionList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Error:      "",
		Positions:  positions,
	})
	return nil
}

func (state *AccountListener) OnOrderStatusRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderStatusRequest)
	// TODO filtering
	orders := state.account.GetOrders()
	context.Respond(&messages.OrderList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Error:      "",
		Orders:     orders,
	})
	return nil
}

func (state *AccountListener) OnNewOrderSingle(context actor.Context) error {
	msg := context.Message().(*messages.NewOrderSingle)
	msg.Account = state.accountM
	order := &models.Order{
		OrderID:        "",
		ClientOrderID:  msg.ClientOrderID,
		Instrument:     msg.Instrument,
		OrderStatus:    models.PendingNew,
		OrderType:      msg.OrderType,
		Side:           msg.OrderSide,
		TimeInForce:    msg.TimeInForce,
		LeavesQuantity: msg.Quantity,
		CumQuantity:    0,
	}
	report, err := state.account.NewOrder(order)
	if err != nil {
		return err
	}
	if report != nil {
		report.SeqNum = state.seqNum + 1
		state.seqNum += 1
		context.Send(context.Parent(), report)
		if report.ExecutionType == messages.PendingNew {
			context.Request(state.bitmexExecutor, msg)
		}
	}

	return nil
}

func (state *AccountListener) OnExecutionReport(context actor.Context) error {
	// We can get execution report from the executers
	report := context.Message().(*messages.ExecutionReport)
	fmt.Println("GOT REPORT BACK ")
	if report.ClientOrderID == nil {
		return fmt.Errorf("report has no ClientOrderID")
	}
	if report.ExecutionType == messages.Rejected {
		nReport, err := state.account.RejectNewOrder(report.ClientOrderID.Value, report.OrderRejectionReason)
		if err != nil {
			return fmt.Errorf("error rejecting order: %v", err)
		}
		if nReport != nil {
			nReport.SeqNum = state.seqNum + 1
			state.seqNum += 1
			context.Send(context.Parent(), nReport)
		}
	} else if report.ExecutionType == messages.New {
		nReport, err := state.account.ConfirmNewOrder(report.ClientOrderID.Value)
		if err != nil {
			return fmt.Errorf("error confirming new order: %v", err)
		}
		if nReport != nil {
			nReport.SeqNum = state.seqNum + 1
			state.seqNum += 1
			context.Send(context.Parent(), nReport)
		}
	}
	return nil
}

func (state *AccountListener) onWebsocketMessage(context actor.Context) error {
	msg := context.Message().(*bitmex.WebsocketMessage)
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
		switch data.ExecType {
		case "New":
			// New order
			fmt.Println("GOT ORDER EWWS")
			if data.ClOrdID == nil {
				return fmt.Errorf("got an order with nil ClOrdID")
			}
			report, err := state.account.ConfirmNewOrder(*data.ClOrdID)
			if err != nil {
				return fmt.Errorf("error confirming new order: %v", err)
			}
			if report != nil {
				report.SeqNum = state.seqNum + 1
				state.seqNum += 1
				context.Send(context.Parent(), report)
			}
		}
	}

	return nil
}

func (state *AccountListener) subscribeAccount(context actor.Context) error {
	if state.ws != nil {
		_ = state.ws.Disconnect()
	}

	ws := bitmex.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to bitmex websocket: %v", err)
	}

	if err := ws.Auth(state.accountM.Credentials); err != nil {
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
			actor.EmptyRootContext.Send(pid, ws.Msg)
		}
	}(ws, context.Self())
	state.ws = ws

	return nil
}
