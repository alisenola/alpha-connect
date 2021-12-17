package coinbasepro

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"reflect"
	"time"
)

// Execute private api calls
// Contains rate limit
// Spawn a query actor for each request
// and pipe its result back

// 429 rate limit
// 418 IP ban

type FixExecutor struct {
	extypes.BaseExecutor
	//fixClient 		*http.Client
	fixRateLimit *exchanges.RateLimit
	logger       *log.Logger
}

func NewFixExecutor() actor.Actor {
	return &FixExecutor{
		fixRateLimit: nil,
	}
}

func (state *FixExecutor) Receive(context actor.Context) {
	extypes.ReceiveExecutor(state, context)
}

func (state *FixExecutor) GetLogger() *log.Logger {
	return state.logger
}

func (state *FixExecutor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))
	state.fixRateLimit = exchanges.NewRateLimit(50, time.Second)
	return nil
}
