package coinbasepro

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	_interface "gitlab.com/alphaticks/alphac/exchanges/interface"
	"gitlab.com/alphaticks/alphac/models/messages"
	"reflect"
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
	logger          *log.Logger
}

func NewExecutor() actor.Actor {
	return &Executor{
		privateExecutor: nil,
		publicExecutor:  nil,
		fixExecutor:     nil,
		logger:          nil,
	}
}

func (state *Executor) Receive(context actor.Context) {
	if _, ok := context.Message().(*messages.SecurityList); ok {
		// Forward security list updates to parent
		context.Forward(context.Parent())
	} else {
		_interface.ExchangeExecutorReceive(state, context)
	}
}

func (state *Executor) GetLogger() *log.Logger {
	return state.logger
}

func (state *Executor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	state.publicExecutor = context.Spawn(actor.PropsFromProducer(NewCoinbaseProPublicExecutor))
	state.privateExecutor = context.Spawn(actor.PropsFromProducer(NewCoinbaseProPrivateExecutor))
	state.fixExecutor = context.Spawn(actor.PropsFromProducer(NewCoinbaseProFixExecutor))

	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	context.Forward(state.publicExecutor)
	return nil
}

func (state *Executor) OnSecurityListRequest(context actor.Context) error {
	context.Forward(state.publicExecutor)
	return nil
}

func (state *Executor) OnMarketDataRequest(context actor.Context) error {
	context.Forward(state.publicExecutor)
	return nil
}
