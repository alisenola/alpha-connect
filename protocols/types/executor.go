package types

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
)

type updateProtocolAssetList struct{}

type Executor interface {
	actor.Actor
	UpdateProtocolAssetList(context actor.Context) error
	OnProtocolAssetListRequest(context actor.Context) error
	OnHistoricalProtocolAssetTransferRequest(context actor.Context) error
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
			context.Send(pid, &updateProtocolAssetList{})
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

	case *updateProtocolAssetList:
		if err := state.UpdateProtocolAssetList(context); err != nil {
			state.GetLogger().Info("error updating security list", log.Error(err))
		}
		go func(pid *actor.PID) {
			time.Sleep(10 * time.Minute)
			context.Send(pid, &updateProtocolAssetList{})
		}(context.Self())

	case *messages.HistoricalProtocolAssetTransferRequest:
		if err := state.OnHistoricalProtocolAssetTransferRequest(context); err != nil {
			state.GetLogger().Error("error processing HistoricalProtocolAssetTransferRequest", log.Error(err))
			panic(err)
		}
	case *messages.ProtocolAssetListRequest:
		if err := state.OnProtocolAssetListRequest(context); err != nil {
			state.GetLogger().Error("error processing ProtocolAssetListRequest", log.Error(err))
			panic(err)
		}
	}
}

func (state *BaseExecutor) UpdateProtocolAssetList(context actor.Context) error {
	return nil
}

func (state *BaseExecutor) OnProtocolAssetListRequest(context actor.Context) error {
	req := context.Message().(*messages.ProtocolAssetListRequest)
	context.Respond(&messages.ProtocolAssetList{
		RequestID:       req.RequestID,
		ResponseID:      uint64(time.Now().UnixNano()),
		Success:         false,
		RejectionReason: messages.RejectionReason_UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnHistoricalProtocolAssetTransferRequest(context actor.Context) error {
	req := context.Message().(*messages.HistoricalProtocolAssetTransferRequest)
	context.Respond(&messages.HistoricalProtocolAssetTransferResponse{
		RequestID:       req.RequestID,
		ResponseID:      uint64(time.Now().UnixNano()),
		Success:         false,
		RejectionReason: messages.RejectionReason_UnsupportedRequest,
	})
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
