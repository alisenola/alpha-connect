package types

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
)

type Listener interface {
	actor.Actor
	GetLogger() *log.Logger
	Initialize(context actor.Context) error
	Clean(context actor.Context) error
	OnTransferDataRequest(context actor.Context) error
}

type BaseListener struct{}

func (state *BaseListener) GetLogger() *log.Logger {
	panic("not implemented")
}

func (state *BaseListener) Initialize(context actor.Context) error {
	panic("not implemented")
}

func (state *BaseListener) Clean(context actor.Context) error {
	panic("not implemented")
}

func (state *BaseListener) OnTransferDataRequest(context actor.Context) error {

	return nil
}

func (state *BaseListener) Receive(context actor.Context) {
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

	case *messages.TransferDataRequest:
		if err := state.OnTransferDataRequest(context); err != nil {
			state.GetLogger().Error("error processing OnTransferDataRequest", log.Error(err))
			panic(err)
		}
	}
}
