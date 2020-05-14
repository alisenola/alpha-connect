package coinbasepro

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alphac/models/messages"
)

// Execute api calls
// Contains rate limit
// Spawn a query actor for each request
// and pipe its result back

// 429 rate limit
// 418 IP ban

// The role of a CoinbasePro Executor is to
// process api request
type Executor struct {
	privateExecutor *actor.PID
	publicExecutor  *actor.PID
	fixExecutor     *actor.PID
}

func NewExecutor() actor.Actor {
	return &Executor{
		privateExecutor: nil,
		publicExecutor:  nil,
		fixExecutor:     nil}
}

func (state *Executor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		pid := context.Spawn(actor.PropsFromProducer(NewCoinbaseProPublicExecutor))
		state.publicExecutor = pid

		pid = context.Spawn(actor.PropsFromProducer(NewCoinbaseProPrivateExecutor))
		state.privateExecutor = pid

		pid = context.Spawn(actor.PropsFromProducer(NewCoinbaseProFixExecutor))
		state.fixExecutor = pid

	case *messages.SecurityListRequest:
		context.Forward(state.publicExecutor)

	case *messages.MarketDataRequest:
		context.Forward(state.publicExecutor)
	}
}
