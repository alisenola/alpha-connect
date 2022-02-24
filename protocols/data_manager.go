package protocols

import (
	"fmt"
	"reflect"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
)

// The market data manager spawns an instrument listener and multiplex its messages
// to actors who subscribed

type DataManager struct {
	subscribers map[uint64]*actor.PID
	listener    *actor.PID
	asset       *models.ProtocolAsset
	logger      *log.Logger
}

func NewDataManagerProducer(protocol *models.ProtocolAsset) actor.Producer {
	return func() actor.Actor {
		return NewDataManager(protocol)
	}
}

func NewDataManager(protocol *models.ProtocolAsset) actor.Actor {
	return &DataManager{
		asset:  protocol,
		logger: nil,
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

	case *messages.ProtocolAssetTransferRequest:
		if err := state.OnProtocolAssetTransferRequest(context); err != nil {
			state.logger.Error("error processing OnProtocolAssetTransferRequest", log.Error(err))
			panic(err)
		}

	case *messages.ProtocolAssetTransferResponse:
		if err := state.OnProtocolAssetTransferResponse(context); err != nil {
			state.logger.Error("error processing OnProtocolAssetTransferResponse", log.Error(err))
			panic(err)
		}

	case *messages.ProtocolAssetDataIncrementalRefresh:
		if err := state.OnProtocolAssetDataIncrementalRefresh(context); err != nil {
			state.logger.Error("error processing ProtocolAssetDataIncrementalRefresh", log.Error(err))
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
	producer := NewAssetListenerProfucer(state.asset)
	if producer == nil {
		return fmt.Errorf("error getting asset listener")
	}
	props := actor.PropsFromProducer(producer)
	state.listener = context.Spawn(props)

	return nil
}

func (state *DataManager) OnProtocolAssetTransferRequest(context actor.Context) error {
	request := context.Message().(*messages.ProtocolAssetTransferRequest)

	if request.Subscribe {
		state.subscribers[request.RequestID] = request.Subscriber
		context.Watch(request.Subscriber)
	}

	context.Forward(state.listener)

	return nil
}

func (state *DataManager) OnProtocolAssetTransferResponse(context actor.Context) error {
	update := context.Message().(*messages.ProtocolAssetTransferResponse)
	for k, v := range state.subscribers {
		forward := &messages.ProtocolAssetTransferResponse{
			RequestID:  k,
			ResponseID: uint64(time.Now().UnixNano()),
			Update:     update.Update,
		}
		context.Send(v, forward)
	}
	return nil
}

func (state *DataManager) OnProtocolAssetDataIncrementalRefresh(context actor.Context) error {
	refresh := context.Message().(*messages.ProtocolAssetDataIncrementalRefresh)
	for k, v := range state.subscribers {
		forward := &messages.ProtocolAssetDataIncrementalRefresh{
			RequestID:  k,
			ResponseID: uint64(time.Now().UnixNano()),
			Update:     refresh.Update,
			SeqNum:     refresh.SeqNum,
		}
		context.Send(v, forward)
	}
	return nil
}

func (state *DataManager) Clean(context actor.Context) error {
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
