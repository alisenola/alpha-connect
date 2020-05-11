package _interface

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	fix50mdr "github.com/quickfixgo/fix50/marketdatarequest"
)

type ExchangeListener interface {
	actor.Actor
	GetLogger() *log.Logger
	Initialize(context actor.Context) error
	Clean(context actor.Context) error
	OnFIX50MarketDataRequest(context actor.Context) error
}

func ExchangeListenerReceive(state ExchangeListener, context actor.Context) {
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

	case *fix50mdr.MarketDataRequest:
		if err := state.OnFIX50MarketDataRequest(context); err != nil {
			state.GetLogger().Error("error processing GetOrderBookL2Request", log.Error(err))
			panic(err)
		}
	}
}
