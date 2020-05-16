package coinbasepro

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/exchanges/interface"
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

type CoinbaseProFixExecutor struct {
	//fixClient 		*http.Client
	fixRateLimit *exchanges.RateLimit
	logger       *log.Logger
}

func NewCoinbaseProFixExecutor() actor.Actor {
	return &CoinbaseProFixExecutor{
		fixRateLimit: nil,
	}
}

func (state *CoinbaseProFixExecutor) Receive(context actor.Context) {
	_interface.ExchangeExecutorReceive(state, context)
}

func (state *CoinbaseProFixExecutor) GetLogger() *log.Logger {
	return state.logger
}

func (state *CoinbaseProFixExecutor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))
	state.fixRateLimit = exchanges.NewRateLimit(50, time.Second)
	return nil
}

func (state *CoinbaseProFixExecutor) Clean(context actor.Context) error {
	return nil
}

func (state *CoinbaseProFixExecutor) UpdateSecurityList(context actor.Context) error {
	return nil
}

func (state *CoinbaseProFixExecutor) OnMarketDataRequest(context actor.Context) error {
	return nil
}

func (state *CoinbaseProFixExecutor) OnSecurityListRequest(context actor.Context) error {
	return nil
}
