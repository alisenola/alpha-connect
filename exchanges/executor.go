package exchanges

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"reflect"
	"time"
)

// The executor routes all the request to the underlying exchange executor & listeners
// He is the main part of the whole software..
//
type Executor struct {
	exchanges     []*xchangerModels.Exchange
	executors     map[uint32]*actor.PID       // A map from exchange ID to executor
	securities    map[uint64]*models.Security // A map from security ID to security
	instruments   map[uint64]*actor.PID       // A map from security ID to instrument listener
	slSubscribers map[uint64]*actor.PID       // Map from request ID to security list subscribers
	logger        *log.Logger
}

func NewExecutorProducer(exchanges []*xchangerModels.Exchange) actor.Producer {
	return func() actor.Actor {
		return NewExecutor(exchanges)
	}
}

func NewExecutor(exchanges []*xchangerModels.Exchange) actor.Actor {
	return &Executor{
		exchanges: exchanges,
		logger:    nil,
	}
}

func (state *Executor) Receive(context actor.Context) {
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

	case *messages.MarketDataRequest:
		if err := state.OnMarketDataRequest(context); err != nil {
			state.logger.Error("error processing OnMarketDataRequest", log.Error(err))
			panic(err)
		}

	case *messages.SecurityListRequest:
		if err := state.OnSecurityListRequest(context); err != nil {
			state.logger.Error("error processing OnSecurityListRequest", log.Error(err))
			panic(err)
		}

	case *messages.SecurityList:
		if err := state.OnSecurityList(context); err != nil {
			state.logger.Error("error processing OnSecurityList", log.Error(err))
			panic(err)
		}

	case *actor.Terminated:
		if err := state.OnTerminated(context); err != nil {
			state.logger.Error("error processing OnTerminated", log.Error(err))
			panic(err)
		}
	}
}

func (state *Executor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	state.instruments = make(map[uint64]*actor.PID)
	state.slSubscribers = make(map[uint64]*actor.PID)
	// Spawn all exchange executors
	state.executors = make(map[uint32]*actor.PID)
	for _, exch := range state.exchanges {
		producer := NewExchangeExecutorProducer(exch)
		if producer == nil {
			return fmt.Errorf("unknown exchange %s", exch.Name)
		}
		props := actor.PropsFromProducer(producer)
		state.executors[exch.ID] = context.Spawn(props)
	}

	// Request securities for each one of them
	var futures []*actor.Future
	request := messages.SecurityListRequest{}
	for _, pid := range state.executors {
		fut := context.RequestFuture(pid, request, 20*time.Second)
		futures = append(futures, fut)
	}

	state.securities = make(map[uint64]*models.Security)
	for _, fut := range futures {
		res, err := fut.Result()
		if err != nil {
			return fmt.Errorf("error fetching securities: %v", err)
		}
		response, ok := res.(*messages.SecurityList)
		if !ok {
			return fmt.Errorf("was expecting GetSecuritiesResponse, got %s", reflect.TypeOf(res).String())
		}
		for _, s := range response.Securities {
			state.securities[s.SecurityID] = s
		}
	}
	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) OnMarketDataRequest(context actor.Context) error {
	request := context.Message().(*messages.MarketDataRequest)
	security, ok := state.securities[request.Instrument.SecurityID]
	if !ok {
		context.Respond(&messages.MarketDataRequestReject{
			RequestID: request.RequestID,
			Reason:    fmt.Sprintf("unknown security"),
		})
		return nil
	}
	if pid, ok := state.instruments[request.Instrument.SecurityID]; ok {
		context.Forward(pid)
	} else {
		props := actor.PropsFromProducer(NewMarketDataManagerProducer(security))
		pid := context.Spawn(props)
		state.instruments[request.Instrument.SecurityID] = pid
		context.Forward(pid)
	}

	return nil
}

func (state *Executor) OnSecurityListRequest(context actor.Context) error {
	request := context.Message().(*messages.SecurityListRequest)
	var securities []*models.Security
	for _, v := range state.securities {
		securities = append(securities, v)
	}
	response := &messages.SecurityList{
		RequestID:  request.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Securities: securities,
	}
	if request.Subscribe {
		context.Watch(context.Sender())
		state.slSubscribers[request.RequestID] = context.Sender()
	}
	context.Respond(response)

	return nil
}

func (state *Executor) OnSecurityList(context actor.Context) error {
	securityList := context.Message().(*messages.SecurityList)
	// Do nothing
	if len(securityList.Securities) == 0 {
		return nil
	}
	exchangeID := securityList.Securities[0].Exchange.ID
	// It has to come from one exchange only, so delete all the known securities from that exchange
	for k, v := range state.securities {
		if v.Exchange.ID == exchangeID {
			delete(state.securities, k)
		}
	}
	// re-add them
	for _, s := range securityList.Securities {
		state.securities[s.SecurityID] = s
	}

	var securities []*models.Security
	for _, v := range state.securities {
		securities = append(securities, v)
	}
	for k, v := range state.slSubscribers {
		securityList := &messages.SecurityList{
			RequestID:  k,
			ResponseID: uint64(time.Now().UnixNano()),
			Securities: securities,
		}
		context.Send(v, securityList)
	}

	return nil
}

func (state *Executor) OnTerminated(context actor.Context) error {
	// Handle subscriber krash
	msg := context.Message().(*actor.Terminated)
	for k, v := range state.slSubscribers {
		if v.Id == msg.Who.Id {
			delete(state.slSubscribers, k)
		}
	}

	return nil
}
