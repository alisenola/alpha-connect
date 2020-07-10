package exchanges

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/account"
	"gitlab.com/alphaticks/alphac/models/messages"
	"reflect"
)

// The account manager spawns an account listener and multiplex its messages
// to actors who subscribed

type AccountManager struct {
	execSubscribers map[uint64]*actor.PID
	trdSubscribers  map[uint64]*actor.PID
	blcSubscribers  map[uint64]*actor.PID
	listener        *actor.PID
	account         *account.Account
	logger          *log.Logger
	paperTrading    bool
}

func NewAccountManagerProducer(account *account.Account, paperTrading bool) actor.Producer {
	return func() actor.Actor {
		return NewAccountManager(account, paperTrading)
	}
}

func NewAccountManager(account *account.Account, paperTrading bool) actor.Actor {
	return &AccountManager{
		account:      account,
		logger:       nil,
		paperTrading: paperTrading,
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

	case *messages.BalancesRequest:
		if err := state.OnBalancesRequest(context); err != nil {
			state.logger.Error("error processing OnBalancesRequest", log.Error(err))
			panic(err)
		}

	case *messages.NewOrderSingleRequest,
		*messages.NewOrderBulkRequest,
		*messages.OrderCancelRequest,
		*messages.OrderMassCancelRequest,
		*messages.OrderReplaceRequest,
		*messages.OrderBulkReplaceRequest:

		context.Forward(state.listener)

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
	state.blcSubscribers = make(map[uint64]*actor.PID)
	var producer actor.Producer
	if state.paperTrading {
		producer = NewPaperAccountListenerProducer(state.account)
		if producer == nil {
			return fmt.Errorf("error getting account listener")
		}
	} else {
		producer = NewAccountListenerProducer(state.account)
		if producer == nil {
			return fmt.Errorf("error getting account listener")
		}
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

	fmt.Println("FORWARDING ORDER STATUS REQUEST")
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

func (state *AccountManager) OnBalancesRequest(context actor.Context) error {
	request := context.Message().(*messages.BalancesRequest)

	if request.Subscribe {
		state.blcSubscribers[request.RequestID] = request.Subscriber
		context.Watch(request.Subscriber)
	}

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
	for k, v := range state.blcSubscribers {
		if v.Id == msg.Who.Id {
			delete(state.blcSubscribers, k)
		}
	}
	return nil
}
