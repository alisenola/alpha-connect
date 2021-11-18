package types

import (
	goContext "context"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type updateSecurityList struct{}

type ExchangeExecutor interface {
	actor.Actor
	OnSecurityListRequest(context actor.Context) error
	OnHistoricalOpenInterestsRequest(context actor.Context) error
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
	UpdateSecurityList(context actor.Context) error
	GetLogger() *log.Logger
	Initialize(context actor.Context) error
	Clean(context actor.Context) error
}

type AccountReconcile interface {
	actor.Actor
	OnAccountMovementRequest(context actor.Context) error
	OnTradeCaptureReportRequest(context actor.Context) error
	GetLogger() *log.Logger
	Initialize(context actor.Context) error
	Clean(context actor.Context) error
}

type ExchangeExecutorBase struct {
}

func (e *ExchangeExecutorBase) OnSecurityListRequest(context actor.Context) error {
	req := context.Message().(*messages.SecurityListRequest)
	context.Respond(&messages.SecurityList{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (e *ExchangeExecutorBase) OnHistoricalOpenInterestsRequest(context actor.Context) error {
	req := context.Message().(*messages.HistoricalOpenInterestsRequest)
	context.Respond(&messages.HistoricalOpenInterestsResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (e *ExchangeExecutorBase) OnHistoricalLiquidationsRequest(context actor.Context) error {
	req := context.Message().(*messages.HistoricalLiquidationsRequest)
	context.Respond(&messages.HistoricalLiquidationsResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (e *ExchangeExecutorBase) OnMarketStatisticsRequest(context actor.Context) error {
	req := context.Message().(*messages.MarketStatisticsRequest)
	context.Respond(&messages.MarketStatisticsResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (e *ExchangeExecutorBase) OnMarketDataRequest(context actor.Context) error {
	req := context.Message().(*messages.MarketDataRequest)
	context.Respond(&messages.MarketDataResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (e *ExchangeExecutorBase) OnAccountMovementRequest(context actor.Context) error {
	req := context.Message().(*messages.AccountMovementRequest)
	context.Respond(&messages.AccountMovementResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (e *ExchangeExecutorBase) OnTradeCaptureReportRequest(context actor.Context) error {
	req := context.Message().(*messages.TradeCaptureReportRequest)
	context.Respond(&messages.TradeCaptureReport{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (e *ExchangeExecutorBase) OnOrderStatusRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderStatusRequest)
	context.Respond(&messages.OrderList{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (e *ExchangeExecutorBase) OnPositionsRequest(context actor.Context) error {
	req := context.Message().(*messages.PositionsRequest)
	context.Respond(&messages.PositionList{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (e *ExchangeExecutorBase) OnBalancesRequest(context actor.Context) error {
	req := context.Message().(*messages.BalancesRequest)
	context.Respond(&messages.BalanceList{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (e *ExchangeExecutorBase) OnNewOrderSingleRequest(context actor.Context) error {
	req := context.Message().(*messages.NewOrderSingleRequest)
	context.Respond(&messages.NewOrderSingleResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (e *ExchangeExecutorBase) OnNewOrderBulkRequest(context actor.Context) error {
	req := context.Message().(*messages.NewOrderBulkRequest)
	context.Respond(&messages.NewOrderBulkResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (e *ExchangeExecutorBase) OnOrderReplaceRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderReplaceRequest)
	context.Respond(&messages.OrderReplaceResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (e *ExchangeExecutorBase) OnOrderBulkReplaceRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderBulkReplaceRequest)
	context.Respond(&messages.OrderBulkReplaceResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (e *ExchangeExecutorBase) OnOrderCancelRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderCancelRequest)
	context.Respond(&messages.OrderCancelResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (e *ExchangeExecutorBase) OnOrderMassCancelRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderMassCancelRequest)
	context.Respond(&messages.OrderMassCancelResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (e *ExchangeExecutorBase) UpdateSecurityList(context actor.Context) error {
	panic("not implemented")
}

func (e *ExchangeExecutorBase) GetLogger() *log.Logger {
	panic("not implemented")
}

func (e *ExchangeExecutorBase) Initialize(context actor.Context) error {
	panic("not implemented")
}

func (e *ExchangeExecutorBase) Clean(context actor.Context) error {
	panic("not implemented")
}

func ExchangeExecutorReceive(state ExchangeExecutor, context actor.Context) {
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

type AccountReconcileBase struct {
}

func (state *AccountReconcileBase) GetLogger() *log.Logger {
	panic("not implemented")
}

func (state *AccountReconcileBase) GetTransactions() *mongo.Collection {
	panic("not implemented")
}

func (state *AccountReconcileBase) OnTradeCaptureReportRequest(context actor.Context) error {
	msg := context.Message().(*messages.TradeCaptureReportRequest)
	res := &messages.TradeCaptureReport{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
	}
	defer context.Respond(res)
	filter := bson.D{
		{"account", msg.Account.Name},
		{"type", "TRADE"},
	}
	if msg.Filter.From != nil {
		from, err := types.TimestampFromProto(msg.Filter.From)
		if err != nil {
			res.Success = false
			res.RejectionReason = messages.InvalidRequest
			return nil
		}
		filter = append(filter, bson.E{"time", bson.D{{"$gt", from}}})
	}
	if msg.Filter.To != nil {
		to, err := types.TimestampFromProto(msg.Filter.To)
		if err != nil {
			res.Success = false
			res.RejectionReason = messages.InvalidRequest
			return nil
		}
		filter = append(filter, bson.E{"time", bson.D{{"$lt", to}}})
	}
	cur, err := state.GetTransactions().Find(goContext.Background(), filter)
	if err != nil {
		res.Success = false
		res.RejectionReason = messages.Other
		state.GetLogger().Error("error fetching trades", log.Error(err))
		return nil
	}
	defer cur.Close(goContext.Background())
	for cur.Next(goContext.Background()) {
		var tx Transaction
		if err := cur.Decode(tx); err != nil {
			state.GetLogger().Error("error decoding tx", log.Error(err))
			res.Success = false
			res.RejectionReason = messages.Other
			return nil
		}
		side := models.Buy
		if tx.Fill.Quantity < 0 {
			side = models.Sell
		}
		secID, _ := strconv.ParseUint(tx.Fill.SecurityID, 10, 64)
		ts, _ := types.TimestampProto(tx.Time)
		var commission float64
		var commissionAsset *xchangerModels.Asset
		for _, m := range tx.Movements {
			if m.Reason == int32(messages.Commission) {
				commission = m.Quantity
				commissionAsset, _ = constants.GetAssetByID(m.AssetID)
			}
		}
		res.Trades = append(res.Trades, &models.TradeCapture{
			Side:            side,
			Type:            models.Regular,
			Price:           tx.Fill.Price,
			Quantity:        math.Abs(tx.Fill.Quantity),
			Commission:      commission,
			CommissionAsset: commissionAsset,
			TradeID:         tx.ID,
			Instrument:      &models.Instrument{SecurityID: &types.UInt64Value{Value: secID}},
			Trade_LinkID:    nil,
			OrderID:         nil,
			ClientOrderID:   nil,
			TransactionTime: ts,
		})
	}
	return nil
}

func AccountReconcileReceive(state AccountReconcile, context actor.Context) {
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
	}
}
