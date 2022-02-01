package types

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
)

type updateCollectionList struct{}

type Executor interface {
	actor.Actor
	UpdateCollectionList(context actor.Context) error
	OnHistoricalNftTransferDataRequest(context actor.Context) error
	GetLogger() *log.Logger
	Initialize(context actor.Context) error
	Clean(context actor.Context) error
}

type BaseExecutor struct {
}

func ReceiveExecutor(state Executor, context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.GetLogger().Error("error initializing", log.Error(err))
			panic(err)
		}
		state.GetLogger().Info("actor started")
		go func(pid *actor.PID) {
			time.Sleep(time.Minute)
			context.Send(pid, &updateCollectionList{})
		}(context.Self())

	case *actor.Stopping:
		if err := state.Clean(context); err != nil {
			state.GetLogger().Error("error stopping", log.Error(err))
			panic(err)
		}
		state.GetLogger().Info("actor stopping")

	case *actor.Stopped:
		state.GetLogger().Info("actor stopped")

	case *actor.Restarting:
		if err := state.Clean(context); err != nil {
			state.GetLogger().Error("error restarting", log.Error(err))
			// Attention, no panic in restarting or infinite loop
		}
		state.GetLogger().Info("actor restarting")

	case *updateCollectionList:
		if err := state.UpdateCollectionList(context); err != nil {
			state.GetLogger().Info("error updating security list", log.Error(err))
		}
		go func(pid *actor.PID) {
			time.Sleep(time.Minute)
			context.Send(pid, &updateCollectionList{})
		}(context.Self())
	case *messages.HistoricalNftTransferDataRequest:
		if err := state.OnHistoricalNftTransferDataRequest(context); err != nil {
			state.GetLogger().Error("error processing HistoricalNftTransferDataRequest", log.Error(err))
		}
	}

	return
}

func (state *BaseExecutor) UpdateCollectionList(context actor.Context) error {
	return nil
}

func (state *BaseExecutor) OnHistoricalNftTransferDataRequest(context actor.Context) error {
	req := context.Message().(*messages.HistoricalNftTransferDataRequest)
	resp := &messages.HistoricalNftTransferDataResponse{
		RequestID:       req.RequestID,
		ResponseID:      uint64(time.Now().UnixNano()),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	}
	context.Respond(resp)
	return nil
}

func (state *BaseExecutor) GetLogger() *log.Logger {
	return nil
}

func (state *BaseExecutor) Initialize(context actor.Context) error {
	return nil
}

func (state *BaseExecutor) Clean(context actor.Context) error {
	return nil
}
