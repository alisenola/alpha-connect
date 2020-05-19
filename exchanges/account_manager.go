package exchanges

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"reflect"
)

// The account manager spawns an account listener and multiplex its messages
// to actors who subscribed

type AccountManager struct {
	execSubscribers map[uint64]*actor.PID
	trdSubscribers  map[uint64]*actor.PID
	listener        *actor.PID
	account         *models.Account
	logger          *log.Logger
}

func NewAccountManagerProducer(account *models.Account) actor.Producer {
	return func() actor.Actor {
		return NewAccountManager(account)
	}
}

func NewAccountManager(account *models.Account) actor.Actor {
	return &AccountManager{
		account: account,
		logger:  nil,
	}
}

func (state *AccountManager) Receive(context actor.Context) {
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

	case *messages.OrderStatusRequest:
		if err := state.OnOrderStatusRequest(context); err != nil {
			state.logger.Error("error processing OnOrderStatusRequest", log.Error(err))
			panic(err)
		}

	case *messages.PositionsRequest:
		if err := state.OnPositionsRequest(context); err != nil {
			state.logger.Error("error processing OnPositionsRequest", log.Error(err))
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

	case *messages.ExecutionReport:
		if err := state.OnExecutionReport(context); err != nil {
			state.logger.Error("error processing OnExecutionReport", log.Error(err))
			panic(err)
		}

	case *actor.Terminated:
		if err := state.OnTerminated(context); err != nil {
			state.logger.Error("error processing OnTerminated", log.Error(err))
			panic(err)
		}
	}
}

func (state *AccountManager) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	state.trdSubscribers = make(map[uint64]*actor.PID)
	state.execSubscribers = make(map[uint64]*actor.PID)
	producer := NewAccountListenerProducer(state.account)
	if producer == nil {
		return fmt.Errorf("error getting account listener")
	}
	props := actor.PropsFromProducer(producer)
	state.listener = context.Spawn(props)

	return nil
}

func (state *AccountManager) Clean(context actor.Context) error {
	return nil
}

func (state *AccountManager) OnOrderStatusRequest(context actor.Context) error {
	request := context.Message().(*messages.OrderStatusRequest)

	if request.Subscribe {
		state.execSubscribers[request.RequestID] = request.Subscriber
		context.Watch(request.Subscriber)
	}

	context.Forward(state.listener)

	return nil
}

func (state *AccountManager) OnPositionsRequest(context actor.Context) error {
	request := context.Message().(*messages.PositionsRequest)

	if request.Subscribe {
		state.trdSubscribers[request.RequestID] = request.Subscriber
		context.Watch(request.Subscriber)
	}

	context.Forward(state.listener)

	return nil
}

func (state *AccountManager) OnNewOrderSingleRequest(context actor.Context) error {
	context.Forward(state.listener)
	return nil
}

func (state *AccountManager) OnOrderCancelRequest(context actor.Context) error {
	context.Forward(state.listener)
	return nil
}

func (state *AccountManager) OnExecutionReport(context actor.Context) error {
	report := context.Message().(*messages.ExecutionReport)
	for _, v := range state.execSubscribers {
		context.Send(v, report)
	}
	return nil
}

func (state *AccountManager) OnTerminated(context actor.Context) error {
	// Handle subscriber krash
	msg := context.Message().(*actor.Terminated)
	for k, v := range state.trdSubscribers {
		if v.Id == msg.Who.Id {
			delete(state.trdSubscribers, k)
		}
	}
	for k, v := range state.execSubscribers {
		if v.Id == msg.Who.Id {
			delete(state.execSubscribers, k)
		}
	}

	return nil
}
