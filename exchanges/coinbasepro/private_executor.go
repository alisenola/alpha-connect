package coinbasepro

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"net/http"
	"reflect"
	"time"
)

// Execute private api calls
// Contains rate limit
// Spawn a query actor for each request
// and pipe its result back

// 429 rate limit
// 418 IP ban

type PrivateExecutor struct {
	httpClient    *http.Client
	httpRateLimit *exchanges.RateLimit
	logger        *log.Logger
}

func NewPrivateExecutor() actor.Actor {
	return &PrivateExecutor{
		httpClient:    nil,
		httpRateLimit: nil,
	}
}

func (state *PrivateExecutor) Receive(context actor.Context) {
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
	}
}

func (state *PrivateExecutor) GetLogger() *log.Logger {
	return state.logger
}

func (state *PrivateExecutor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	state.httpClient = &http.Client{Timeout: 10 * time.Second}
	state.httpRateLimit = exchanges.NewRateLimit(5, time.Second)
	return nil
}

func (state *PrivateExecutor) Clean(context actor.Context) error {
	return nil
}
