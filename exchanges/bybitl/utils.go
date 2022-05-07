package bybitl

import (
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/bybitl"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strconv"
)

func statusToBybitl(status models.OrderStatus) bybitl.OrderStatus {
	switch status {
	case models.OrderStatus_New:
		return bybitl.OrderNew
	case models.OrderStatus_PendingCancel:
		return bybitl.OrderPendingCancel
	case models.OrderStatus_PendingNew, models.OrderStatus_PendingReplace:
		return bybitl.OrderCreated
	case models.OrderStatus_Filled:
		return bybitl.OrderFilled
	case models.OrderStatus_Canceled:
		return bybitl.OrderCancelled
	case models.OrderStatus_Rejected:
		return bybitl.OrderRejected
	default:
		return ""
	}
}

func orderToModel(order *bybitl.ActiveOrder) *models.Order {
	o := &models.Order{
		OrderID:       order.OrderId,
		ClientOrderID: order.OrderLinkId,
		Instrument: &models.Instrument{
			Exchange: constants.BYBITL,
			Symbol:   &wrapperspb.StringValue{Value: order.Symbol},
		},
		Price:          &wrapperspb.DoubleValue{Value: order.Price},
		LeavesQuantity: order.Qty - order.CumExecQty,
		CumQuantity:    order.CumExecQty,
		CreationTime:   timestamppb.New(order.CreatedTime),
	}

	if order.CloseOnTrigger {
		o.ExecutionInstructions = append(o.ExecutionInstructions, models.ExecutionInstruction_CloseOnTrigger)
	}
	if order.ReduceOnly {
		o.ExecutionInstructions = append(o.ExecutionInstructions, models.ExecutionInstruction_ReduceOnly)
	}

	switch order.OrderStatus {
	case bybitl.OrderCreated:
		o.OrderStatus = models.OrderStatus_Created
	case bybitl.OrderRejected:
		o.OrderStatus = models.OrderStatus_Rejected
	case bybitl.OrderNew:
		o.OrderStatus = models.OrderStatus_New
	case bybitl.OrderPartiallyFilled:
		o.OrderStatus = models.OrderStatus_PartiallyFilled
	case bybitl.OrderFilled:
		o.OrderStatus = models.OrderStatus_Filled
	case bybitl.OrderCancelled:
		o.OrderStatus = models.OrderStatus_Canceled
	case bybitl.OrderPendingCancel:
		o.OrderStatus = models.OrderStatus_PendingCancel
	default:
		fmt.Println("UNKNOWN ORDER STATUS", order.OrderStatus)
	}

	switch order.OrderType {
	case bybitl.LimitOrder:
		o.OrderType = models.OrderType_Limit
	case bybitl.MarketOrder:
		o.OrderType = models.OrderType_Market
	default:
		fmt.Println("UNKNOWN ORDER TYPE", order.OrderType)
	}

	switch order.Side {
	case bybitl.Buy:
		o.Side = models.Side_Buy
	case bybitl.Sell:
		o.Side = models.Side_Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", order.Side)
	}

	switch order.TimeInForce {
	case bybitl.GoodTillCancel:
		o.TimeInForce = models.TimeInForce_GoodTillCancel
	case bybitl.ImmediateOrCancel:
		o.TimeInForce = models.TimeInForce_ImmediateOrCancel
	case bybitl.FillOrKill:
		o.TimeInForce = models.TimeInForce_FillOrKill
	case bybitl.PostOnly:
		o.TimeInForce = models.TimeInForce_PostOnly
	default:
		fmt.Println("UNKNOWN ORDER TIME IN FORCE", order.TimeInForce)
	}

	return o
}

func wsOrderToModel(order *bybitl.WSOrder) *models.Order {
	ord := &models.Order{
		OrderID:       order.OrderId,
		ClientOrderID: order.OrderLinkId,
		Instrument: &models.Instrument{
			Exchange: constants.BYBITL,
			Symbol:   &wrapperspb.StringValue{Value: order.Symbol},
		},
		LeavesQuantity: order.LeavesQty,
		CumQuantity:    order.CumExecQty,
		Price:          &wrapperspb.DoubleValue{Value: order.Price},
	}

	if order.ReduceOnly {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ExecutionInstruction_ReduceOnly)
	}
	if order.CloseOnTrigger {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ExecutionInstruction_CloseOnTrigger)
	}

	switch order.OrderStatus {
	case bybitl.OrderCreated:
		ord.OrderStatus = models.OrderStatus_Created
	case bybitl.OrderRejected:
		ord.OrderStatus = models.OrderStatus_Rejected
	case bybitl.OrderNew:
		ord.OrderStatus = models.OrderStatus_New
	case bybitl.OrderPartiallyFilled:
		ord.OrderStatus = models.OrderStatus_PartiallyFilled
	case bybitl.OrderFilled:
		ord.OrderStatus = models.OrderStatus_Filled
	case bybitl.OrderCancelled:
		ord.OrderStatus = models.OrderStatus_Canceled
	case bybitl.OrderPendingCancel:
		ord.OrderStatus = models.OrderStatus_PendingCancel
	default:
		fmt.Println("UNKNOWN ORDER STATUS", order.OrderStatus)
	}

	switch order.OrderType {
	case bybitl.LimitOrder:
		ord.OrderType = models.OrderType_Limit
	case bybitl.MarketOrder:
		ord.OrderType = models.OrderType_Market
	default:
		fmt.Println("UNKNOWN ORDER TYPE", order.OrderType)
	}

	switch order.Side {
	case bybitl.Buy:
		ord.Side = models.Side_Buy
	case bybitl.Sell:
		ord.Side = models.Side_Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", order.Side)
	}

	switch order.TimeInForce {
	case bybitl.GoodTillCancel:
		ord.TimeInForce = models.TimeInForce_GoodTillCancel
	case bybitl.ImmediateOrCancel:
		ord.TimeInForce = models.TimeInForce_ImmediateOrCancel
	case bybitl.FillOrKill:
		ord.TimeInForce = models.TimeInForce_FillOrKill
	case bybitl.PostOnly:
		ord.TimeInForce = models.TimeInForce_PostOnly
	default:
		fmt.Println("UNKNOWN ORDER TIME IN FORCE", order.TimeInForce)
	}

	return ord
}

func buildPostOrderRequest(symbol string, order *messages.NewOrder, tickPrecision, lotPrecision int) (bybitl.PostActiveOrderParams, *messages.RejectionReason) {
	var side bybitl.Side
	var typ bybitl.OrderType
	var tif bybitl.TimeInForce
	var redo = false

	if order.OrderSide == models.Side_Buy {
		side = bybitl.Buy
	} else {
		side = bybitl.Sell
	}
	switch order.OrderType {
	case models.OrderType_Limit:
		typ = bybitl.LimitOrder
	case models.OrderType_Market:
		typ = bybitl.MarketOrder
	default:
		rej := messages.RejectionReason_UnsupportedOrderType
		return nil, &rej
	}

	switch order.TimeInForce {
	case models.TimeInForce_Session:
		tif = bybitl.GoodTillCancel
	case models.TimeInForce_GoodTillCancel:
		tif = bybitl.GoodTillCancel
	case models.TimeInForce_ImmediateOrCancel:
		tif = bybitl.ImmediateOrCancel
	case models.TimeInForce_FillOrKill:
		tif = bybitl.FillOrKill
	}

	for _, exec := range order.ExecutionInstructions {
		rej := messages.RejectionReason_UnsupportedOrderCharacteristic
		switch exec {
		case models.ExecutionInstruction_ReduceOnly:
			redo = true
		case models.ExecutionInstruction_ParticipateDoNotInitiate:
			tif = bybitl.PostOnly
		default:
			return nil, &rej
		}
	}

	request := bybitl.NewPostActiveOrderParams(symbol, side, typ, tif, redo, false)

	fmt.Println("SET QTY", order.Quantity, lotPrecision, strconv.FormatFloat(order.Quantity, 'f', int(lotPrecision), 64))
	request.SetQuantity(order.Quantity, lotPrecision)
	request.SetOrderLinkId(order.ClientOrderID)
	// For now, only one way is supported
	if side == bybitl.Buy {
		if redo {
			// Buy, reduce only, we are working with a short
			request.SetPositionIdx(2)
		} else {
			// Buy, we are working with a long
			request.SetPositionIdx(1)
		}
	} else {
		if redo {
			// Sell, reduce only, we are working with a long
			request.SetPositionIdx(1)
		} else {
			// Sell, we are working with a short
			request.SetPositionIdx(2)
		}
	}

	if order.Price != nil {
		request.SetPrice(order.Price.Value, tickPrecision)
	}

	return request, nil
}
