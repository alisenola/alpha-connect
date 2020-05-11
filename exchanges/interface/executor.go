package _interface

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	fix50sdr "github.com/quickfixgo/fix50/securitydefinitionrequest"
	fix50nso "github.com/quickfixgo/quickfix/fix50/newordersingle"
)

type ExchangeExecutor interface {
	actor.Actor
	GetLogger() *log.Logger
	Initialize(context actor.Context) error
	Clean(context actor.Context) error
	OnFIX50NewOrderSingle(context actor.Context) error
	OnFIX50SecurityDefinitionRequest(context actor.Context) error
}

func ExchangeExecutorReceive(state ExchangeExecutor, context actor.Context) {
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

	case *fix50nso.NewOrderSingle:
		if err := state.OnFIX50NewOrderSingle(context); err != nil {
			state.GetLogger().Error("error processing FIX50NewOrderSingle", log.Error(err))
			panic(err)
		}

	case *fix50sdr.SecurityDefinitionRequest:
		if err := state.OnFIX50SecurityDefinitionRequest(context); err != nil {
			state.GetLogger().Error("error processing FIX50SecurityDefinitionRequest", log.Error(err))
			panic(err)
		}
	}
}
