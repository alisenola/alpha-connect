package types

import (
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	models2 "gitlab.com/alphaticks/xchanger/models"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
)

type Executor interface {
	actor.Actor
	GetLogger() *log.Logger
	Initialize(context actor.Context) error
	Clean(context actor.Context) error
	OnBlockNumberRequest(context actor.Context) error
	OnEVMContractCallRequest(context actor.Context) error
	OnEVMLogsQueryRequest(context actor.Context) error
	OnEVMLogsSubscribeRequest(context actor.Context) error
}

type BaseExecutor struct {
	Registry registry.PublicRegistryClient
	Chains   []*models2.Chain
}

func ReceiveExecutor(state Executor, context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.GetLogger().Error("error initializing", log.Error(err))
			panic(err)
		}
		state.GetLogger().Info("actor started")

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

	case *messages.BlockNumberRequest:
		if err := state.OnBlockNumberRequest(context); err != nil {
			state.GetLogger().Error("error processing OnBlockNumberRequest", log.Error(err))
			panic(err)
		}

	case *messages.EVMContractCallRequest:
		if err := state.OnEVMContractCallRequest(context); err != nil {
			state.GetLogger().Error("error processing OnEVMContractCallRequest", log.Error(err))
			panic(err)
		}

	case *messages.EVMLogsQueryRequest:
		if err := state.OnEVMLogsQueryRequest(context); err != nil {
			state.GetLogger().Error("error processing OnEVMLogsQueryRequest", log.Error(err))
			panic(err)
		}

	case *messages.EVMLogsSubscribeRequest:
		if err := state.OnEVMLogsSubscribeRequest(context); err != nil {
			state.GetLogger().Error("error processing OnEVMLogsSubscribeRequest", log.Error(err))
			panic(err)
		}
	}
}

func (state *BaseExecutor) OnBlockNumberRequest(context actor.Context) error {
	req := context.Message().(*messages.BlockNumberRequest)
	context.Respond(&messages.BlockNumberResponse{
		RequestID:       req.RequestID,
		ResponseID:      uint64(time.Now().UnixNano()),
		Success:         false,
		RejectionReason: messages.RejectionReason_UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnEVMContractCallRequest(context actor.Context) error {
	req := context.Message().(*messages.EVMContractCallRequest)
	context.Respond(&messages.EVMContractCallResponse{
		RequestID:       req.RequestID,
		ResponseID:      uint64(time.Now().UnixNano()),
		Success:         false,
		RejectionReason: messages.RejectionReason_UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnEVMLogsQueryRequest(context actor.Context) error {
	req := context.Message().(*messages.EVMLogsQueryRequest)
	context.Respond(&messages.EVMLogsQueryResponse{
		RequestID:       req.RequestID,
		ResponseID:      uint64(time.Now().UnixNano()),
		Success:         false,
		RejectionReason: messages.RejectionReason_UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnEVMLogsSubscribeRequest(context actor.Context) error {
	req := context.Message().(*messages.EVMLogsSubscribeRequest)
	context.Respond(&messages.EVMLogsSubscribeResponse{
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
