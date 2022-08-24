package coinbasepro

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
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
		log.String("type", reflect.TypeOf(state).String()))
	state.fixRateLimit = exchanges.NewRateLimit(50, time.Second)
	return nil
}

func (state *FixExecutor) Clean(context actor.Context) error {
	return nil
}

func (state *FixExecutor) OnHistoricalLiquidationsRequest(context actor.Context) error {
	return nil
}

func (state *FixExecutor) UpdateSecurityList(context actor.Context) error {
	return nil
}

func (state *FixExecutor) OnMarketStatisticsRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketStatisticsResponse)
	context.Respond(&messages.MarketStatisticsResponse{
		RequestID:       msg.RequestID,
		Success:         false,
		RejectionReason: messages.RejectionReason_UnsupportedRequest,
	})
	return nil
}

func (state *FixExecutor) OnMarketDataRequest(context actor.Context) error {
	return nil
}

func (state *FixExecutor) OnSecurityListRequest(context actor.Context) error {
	return nil
}

func (state *FixExecutor) OnOrderStatusRequest(context actor.Context) error {
	return nil
}

func (state *FixExecutor) OnPositionsRequest(context actor.Context) error {
	return nil
}

func (state *FixExecutor) OnBalancesRequest(context actor.Context) error {
	return nil
}

func (state *FixExecutor) OnNewOrderSingleRequest(context actor.Context) error {
	return nil
}

func (state *FixExecutor) OnNewOrderBulkRequest(context actor.Context) error {
	return nil
}

func (state *FixExecutor) OnOrderReplaceRequest(context actor.Context) error {
	return nil
}

func (state *FixExecutor) OnOrderBulkReplaceRequest(context actor.Context) error {
	return nil
}

func (state *FixExecutor) OnOrderCancelRequest(context actor.Context) error {
	return nil
}

func (state *FixExecutor) OnOrderMassCancelRequest(context actor.Context) error {
	return nil
}
