package types

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
)

type updateAssetList struct{}

type Executor interface {
	actor.Actor
	UpdateAssetList(context actor.Context) error
	OnAssetListRequest(context actor.Context) error
	OnHistoricalAssetTransferRequest(context actor.Context) error
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
			context.Send(pid, &updateAssetList{})
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

	case *updateAssetList:
		if err := state.UpdateAssetList(context); err != nil {
			state.GetLogger().Info("error updating security list", log.Error(err))
		}
		go func(pid *actor.PID) {
			time.Sleep(10 * time.Minute)
			context.Send(pid, &updateAssetList{})
		}(context.Self())

	case *messages.HistoricalAssetTransferRequest:
		if err := state.OnHistoricalAssetTransferRequest(context); err != nil {
			state.GetLogger().Error("error processing HistoricalNftTransferDataRequest", log.Error(err))
			panic(err)
		}
	case *messages.AssetListRequest:
		if err := state.OnAssetListRequest(context); err != nil {
			state.GetLogger().Error("error getting security list", log.Error(err))
			panic(err)
		}
	}
}

func (state *BaseExecutor) UpdateAssetList(context actor.Context) error {
	return nil
}

func (state *BaseExecutor) OnAssetListRequest(context actor.Context) error {
	req := context.Message().(*messages.AssetListRequest)
	context.Respond(&messages.AssetListResponse{
		RequestID:       req.RequestID,
		ResponseID:      uint64(time.Now().UnixNano()),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnHistoricalAssetTransferRequest(context actor.Context) error {
	req := context.Message().(*messages.HistoricalAssetTransferRequest)
	context.Respond(&messages.HistoricalAssetTransferResponse{
		RequestID:       req.RequestID,
		ResponseID:      uint64(time.Now().UnixNano()),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
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
