package exchanges

import (
	"fmt"
	"reflect"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/utils"
)

// The market data manager spawns an instrument listener and multiplex its messages
// to actors who subscribed

type DataManager struct {
	subscribers map[uint64]*actor.PID
	listener    *actor.PID
	security    *models.Security
	dialerPool  *utils.DialerPool
	logger      *log.Logger
}

func NewDataManagerProducer(security *models.Security, dialerPool *utils.DialerPool) actor.Producer {
	return func() actor.Actor {
		return NewDataManager(security, dialerPool)
	}
}

func NewDataManager(security *models.Security, dialerPool *utils.DialerPool) actor.Actor {
	return &DataManager{
		security:   security,
		dialerPool: dialerPool,
		logger:     nil,
	}
}

func (state *DataManager) Receive(context actor.Context) {
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

	case *messages.MarketDataResponse:
		if err := state.OnMarketDataResponse(context); err != nil {
			state.logger.Error("error processing OnMarketDataResponse", log.Error(err))
			panic(err)
		}

	case *messages.MarketDataIncrementalRefresh:
		if err := state.OnMarketDataIncrementalRefresh(context); err != nil {
			state.logger.Error("error processing OnMarketDataIncrementalRefresh", log.Error(err))
			panic(err)
		}

	case *messages.UnipoolV3DataRequest:
		if err := state.OnUnipoolV3DataRequest(context); err != nil {
			state.logger.Error("error processing OnUnipoolV3DataRequest", log.Error(err))
			panic(err)
		}

	case *messages.UnipoolV3DataResponse:
		if err := state.OnUnipoolV3DataResponse(context); err != nil {
			state.logger.Error("error processing OnUnipoolV3DataResponse", log.Error(err))
			panic(err)
		}

	case *messages.UnipoolV3DataIncrementalRefresh:
		if err := state.OnUnipoolV3DataIncrementalRefresh(context); err != nil {
			state.logger.Error("error processing OnUnipoolV3DataIncrementalRefresh", log.Error(err))
			panic(err)
		}

	case *actor.Terminated:
		if err := state.OnTerminated(context); err != nil {
			state.logger.Error("error processing OnTerminated", log.Error(err))
			panic(err)
		}
	}
}

func (state *DataManager) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	state.subscribers = make(map[uint64]*actor.PID)
	producer := NewInstrumentListenerProducer(state.security.SecurityID, state.security.Exchange.ID, state.dialerPool)
	if producer == nil {
		return fmt.Errorf("error getting instrument listener")
	}
	props := actor.PropsFromProducer(producer)
	state.listener = context.Spawn(props)

	return nil
}

func (state *DataManager) Clean(context actor.Context) error {
	return nil
}

func (state *DataManager) OnMarketDataRequest(context actor.Context) error {
	request := context.Message().(*messages.MarketDataRequest)

	if request.Subscribe {
		state.subscribers[request.RequestID] = request.Subscriber
		context.Watch(request.Subscriber)
	}

	context.Forward(state.listener)

	return nil
}

func (state *DataManager) OnUnipoolV3DataRequest(context actor.Context) error {
	request := context.Message().(*messages.UnipoolV3DataRequest)

	if request.Subscribe {
		state.subscribers[request.RequestID] = request.Subscriber
		context.Watch(request.Subscriber)
	}

	context.Forward(state.listener)

	return nil
}

func (state *DataManager) OnMarketDataResponse(context actor.Context) error {
	snapshot := context.Message().(*messages.MarketDataResponse)
	for k, v := range state.subscribers {
		forward := &messages.MarketDataResponse{
			RequestID:  k,
			ResponseID: uint64(time.Now().UnixNano()),
			SnapshotL2: snapshot.SnapshotL2,
			SnapshotL3: snapshot.SnapshotL3,
			Trades:     snapshot.Trades,
		}
		context.Send(v, forward)
	}
	return nil
}

func (state *DataManager) OnMarketDataIncrementalRefresh(context actor.Context) error {
	refresh := context.Message().(*messages.MarketDataIncrementalRefresh)
	for k, v := range state.subscribers {
		forward := &messages.MarketDataIncrementalRefresh{
			RequestID:   k,
			ResponseID:  uint64(time.Now().UnixNano()),
			UpdateL2:    refresh.UpdateL2,
			UpdateL3:    refresh.UpdateL3,
			Trades:      refresh.Trades,
			Liquidation: refresh.Liquidation,
			Stats:       refresh.Stats,
			SeqNum:      refresh.SeqNum,
		}
		context.Send(v, forward)
	}
	return nil
}

func (state *DataManager) OnUnipoolV3DataResponse(context actor.Context) error {
	update := context.Message().(*messages.UnipoolV3DataResponse)
	for k, v := range state.subscribers {
		forward := &messages.UnipoolV3DataResponse{
			RequestID:  k,
			ResponseID: uint64(time.Now().UnixNano()),
			Update:     update.Update,
		}
		context.Send(v, forward)
	}
	return nil
}

func (state *DataManager) OnUnipoolV3DataIncrementalRefresh(context actor.Context) error {
	refresh := context.Message().(*messages.UnipoolV3DataIncrementalRefresh)
	for k, v := range state.subscribers {
		forward := &messages.UnipoolV3DataIncrementalRefresh{
			RequestID:  k,
			ResponseID: uint64(time.Now().UnixNano()),
			Update:     refresh.Update,
			SeqNum:     refresh.SeqNum,
		}
		context.Send(v, forward)
	}
	return nil
}

func (state *DataManager) OnTerminated(context actor.Context) error {
	// Handle subscriber krash
	msg := context.Message().(*actor.Terminated)
	for k, v := range state.subscribers {
		if v.String() == msg.Who.String() {
			delete(state.subscribers, k)
		}
	}
	if len(state.subscribers) == 0 {
		// Sudoku
		context.Stop(context.Self())
	}

	return nil
}
