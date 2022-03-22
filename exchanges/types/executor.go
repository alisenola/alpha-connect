package types

import (
	"math/rand"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
)

type updateSecurityList struct{}

type Executor interface {
	actor.Actor
	OnSecurityListRequest(context actor.Context) error
	OnHistoricalOpenInterestsRequest(context actor.Context) error
	OnHistoricalFundingRatesRequest(context actor.Context) error
	OnHistoricalLiquidationsRequest(context actor.Context) error
	OnMarketStatisticsRequest(context actor.Context) error
	OnMarketDataRequest(context actor.Context) error
	OnAccountMovementRequest(context actor.Context) error
	OnTradeCaptureReportRequest(context actor.Context) error
	OnOrderStatusRequest(context actor.Context) error
	OnPositionsRequest(context actor.Context) error
	OnBalancesRequest(context actor.Context) error
	OnNewOrderSingleRequest(context actor.Context) error
	OnNewOrderBulkRequest(context actor.Context) error
	OnOrderReplaceRequest(context actor.Context) error
	OnOrderBulkReplaceRequest(context actor.Context) error
	OnOrderCancelRequest(context actor.Context) error
	OnOrderMassCancelRequest(context actor.Context) error
	OnHistoricalUnipoolV3DataRequest(context actor.Context) error
	UpdateSecurityList(context actor.Context) error
	GetLogger() *log.Logger
	Initialize(context actor.Context) error
	Clean(context actor.Context) error
}

type BaseExecutor struct {
}

func ReceiveExecutor(state Executor, context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.GetLogger().Error("error initializing", log.Error(err))
			panic(err)
		}
		state.GetLogger().Info("actor started")
		go func(pid *actor.PID) {
			time.Sleep(time.Minute)
			context.Send(pid, &updateSecurityList{})
		}(context.Self())

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

	case *messages.SecurityListRequest:
		if err := state.OnSecurityListRequest(context); err != nil {
			state.GetLogger().Error("error processing SecurityListRequest", log.Error(err))
			panic(err)
		}

	case *messages.HistoricalLiquidationsRequest:
		if err := state.OnHistoricalLiquidationsRequest(context); err != nil {
			state.GetLogger().Error("error processing OnHistoricalLiquidationRequest", log.Error(err))
			panic(err)
		}

	case *messages.HistoricalOpenInterestsRequest:
		if err := state.OnHistoricalOpenInterestsRequest(context); err != nil {
			state.GetLogger().Error("error processing OnHistoricalOpenInterestsRequest", log.Error(err))
			panic(err)
		}

	case *messages.HistoricalFundingRatesRequest:
		if err := state.OnHistoricalFundingRatesRequest(context); err != nil {
			state.GetLogger().Error("error processing OnHistoricalFundingRateRequest", log.Error(err))
			panic(err)
		}

	case *messages.MarketStatisticsRequest:
		if err := state.OnMarketStatisticsRequest(context); err != nil {
			state.GetLogger().Error("error processing OnMarketStatisticsRequest", log.Error(err))
			panic(err)
		}

	case *messages.MarketDataRequest:
		if err := state.OnMarketDataRequest(context); err != nil {
			state.GetLogger().Error("error processing MarketDataRequest", log.Error(err))
			panic(err)
		}

	case *messages.PositionsRequest:
		if err := state.OnPositionsRequest(context); err != nil {
			state.GetLogger().Error("error processing OnPositionListRequest", log.Error(err))
			panic(err)
		}

	case *messages.BalancesRequest:
		if err := state.OnBalancesRequest(context); err != nil {
			state.GetLogger().Error("error processing OnBalancesRequest", log.Error(err))
			panic(err)
		}

	case *messages.AccountMovementRequest:
		if err := state.OnAccountMovementRequest(context); err != nil {
			state.GetLogger().Error("error processing OnAccountMovementRequest", log.Error(err))
			panic(err)
		}

	case *messages.TradeCaptureReportRequest:
		if err := state.OnTradeCaptureReportRequest(context); err != nil {
			state.GetLogger().Error("error processing OnTradeCaptureReportRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderStatusRequest:
		if err := state.OnOrderStatusRequest(context); err != nil {
			state.GetLogger().Error("error processing OnOrderStatusRequest", log.Error(err))
			panic(err)
		}

	case *messages.NewOrderSingleRequest:
		if err := state.OnNewOrderSingleRequest(context); err != nil {
			state.GetLogger().Error("error processing OnNewSingleOrderRequest", log.Error(err))
			panic(err)
		}

	case *messages.NewOrderBulkRequest:
		if err := state.OnNewOrderBulkRequest(context); err != nil {
			state.GetLogger().Error("error processing OnNewOrderBulkRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderReplaceRequest:
		if err := state.OnOrderReplaceRequest(context); err != nil {
			state.GetLogger().Error("error processing OnOrderReplaceRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderBulkReplaceRequest:
		if err := state.OnOrderBulkReplaceRequest(context); err != nil {
			state.GetLogger().Error("erro procesing OnOrderBulkReplaceRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderCancelRequest:
		if err := state.OnOrderCancelRequest(context); err != nil {
			state.GetLogger().Error("error processing OnOrderCancelRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderMassCancelRequest:
		if err := state.OnOrderMassCancelRequest(context); err != nil {
			state.GetLogger().Error("error processing OnOrderMassCancelRequest", log.Error(err))
			panic(err)
		}

	case *messages.HistoricalUnipoolV3DataRequest:
		if err := state.OnHistoricalUnipoolV3DataRequest(context); err != nil {
			state.GetLogger().Error("error processing HistoricalUnipoolV3DataRequest", log.Error(err))
			panic(err)
		}

	case *updateSecurityList:
		if err := state.UpdateSecurityList(context); err != nil {
			state.GetLogger().Info("error updating security list", log.Error(err))
		}
		go func(pid *actor.PID) {
			time.Sleep(time.Minute)
			context.Send(pid, &updateSecurityList{})
		}(context.Self())
	}
}

func (state *BaseExecutor) OnSecurityListRequest(context actor.Context) error {
	req := context.Message().(*messages.SecurityListRequest)
	context.Respond(&messages.SecurityList{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnHistoricalOpenInterestsRequest(context actor.Context) error {
	req := context.Message().(*messages.HistoricalOpenInterestsRequest)
	context.Respond(&messages.HistoricalOpenInterestsResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnHistoricalFundingRatesRequest(context actor.Context) error {
	req := context.Message().(*messages.HistoricalFundingRatesRequest)
	context.Respond(&messages.HistoricalFundingRatesResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnHistoricalLiquidationsRequest(context actor.Context) error {
	req := context.Message().(*messages.HistoricalLiquidationsRequest)
	context.Respond(&messages.HistoricalLiquidationsResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnMarketStatisticsRequest(context actor.Context) error {
	req := context.Message().(*messages.MarketStatisticsRequest)
	context.Respond(&messages.MarketStatisticsResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnMarketDataRequest(context actor.Context) error {
	req := context.Message().(*messages.MarketDataRequest)
	context.Respond(&messages.MarketDataResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnAccountMovementRequest(context actor.Context) error {
	req := context.Message().(*messages.AccountMovementRequest)
	context.Respond(&messages.AccountMovementResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnTradeCaptureReportRequest(context actor.Context) error {
	req := context.Message().(*messages.TradeCaptureReportRequest)
	context.Respond(&messages.TradeCaptureReport{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnOrderStatusRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderStatusRequest)
	context.Respond(&messages.OrderList{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnPositionsRequest(context actor.Context) error {
	req := context.Message().(*messages.PositionsRequest)
	context.Respond(&messages.PositionList{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnBalancesRequest(context actor.Context) error {
	req := context.Message().(*messages.BalancesRequest)
	context.Respond(&messages.BalanceList{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnNewOrderSingleRequest(context actor.Context) error {
	req := context.Message().(*messages.NewOrderSingleRequest)
	context.Respond(&messages.NewOrderSingleResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnNewOrderBulkRequest(context actor.Context) error {
	req := context.Message().(*messages.NewOrderBulkRequest)
	context.Respond(&messages.NewOrderBulkResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnOrderReplaceRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderReplaceRequest)
	context.Respond(&messages.OrderReplaceResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnOrderBulkReplaceRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderBulkReplaceRequest)
	context.Respond(&messages.OrderBulkReplaceResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnOrderCancelRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderCancelRequest)
	context.Respond(&messages.OrderCancelResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnOrderMassCancelRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderMassCancelRequest)
	context.Respond(&messages.OrderMassCancelResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) OnHistoricalUnipoolV3DataRequest(context actor.Context) error {
	req := context.Message().(*messages.HistoricalUnipoolV3DataRequest)
	context.Respond(&messages.HistoricalUnipoolV3DataResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *BaseExecutor) UpdateSecurityList(context actor.Context) error {
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
