package bybitl

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/bybitl"
	"strconv"
)

func statusToBybitl(status models.OrderStatus) bybitl.OrderStatus {
	switch status {
	case models.New:
		return bybitl.OrderNew
	case models.PendingCancel:
		return bybitl.OrderPendingCancel
	case models.PendingNew, models.PendingReplace:
		return bybitl.OrderCreated
	case models.Filled:
		return bybitl.OrderFilled
	case models.Canceled:
		return bybitl.OrderCancelled
	case models.Rejected:
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
			Exchange: &constants.BYBITL,
			Symbol:   &types.StringValue{Value: order.Symbol},
		},
		Price:          &types.DoubleValue{Value: order.Price},
		LeavesQuantity: order.Qty - order.CumExecQty,
		CumQuantity:    order.CumExecQty,
		CreationTime:   &types.Timestamp{Seconds: order.CreatedTime.Unix()},
	}

	if order.CloseOnTrigger {
		o.ExecutionInstructions = append(o.ExecutionInstructions, models.CloseOnTrigger)
	}
	if order.ReduceOnly {
		o.ExecutionInstructions = append(o.ExecutionInstructions, models.ReduceOnly)
	}

	switch order.OrderStatus {
	case bybitl.OrderCreated:
		o.OrderStatus = models.Created
	case bybitl.OrderRejected:
		o.OrderStatus = models.Rejected
	case bybitl.OrderNew:
		o.OrderStatus = models.New
	case bybitl.OrderPartiallyFilled:
		o.OrderStatus = models.PartiallyFilled
	case bybitl.OrderFilled:
		o.OrderStatus = models.Filled
	case bybitl.OrderCancelled:
		o.OrderStatus = models.Canceled
	case bybitl.OrderPendingCancel:
		o.OrderStatus = models.PendingCancel
	default:
		fmt.Println("UNKNOWN ORDER STATUS", order.OrderStatus)
	}

	switch order.OrderType {
	case bybitl.LimitOrder:
		o.OrderType = models.Limit
	case bybitl.MarketOrder:
		o.OrderType = models.Market
	default:
		fmt.Println("UNKNOWN ORDER TYPE", order.OrderType)
	}

	switch order.Side {
	case bybitl.Buy:
		o.Side = models.Buy
	case bybitl.Sell:
		o.Side = models.Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", order.Side)
	}

	switch order.TimeInForce {
	case bybitl.GoodTillCancel:
		o.TimeInForce = models.GoodTillCancel
	case bybitl.ImmediateOrCancel:
		o.TimeInForce = models.ImmediateOrCancel
	case bybitl.FillOrKill:
		o.TimeInForce = models.FillOrKill
	case bybitl.PostOnly:
		o.TimeInForce = models.PostOnly
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
			Exchange: &constants.BYBITL,
			Symbol:   &types.StringValue{Value: order.Symbol},
		},
		LeavesQuantity: order.LeavesQty,
		CumQuantity:    order.CumExecQty,
		Price:          &types.DoubleValue{Value: order.Price},
	}

	if order.ReduceOnly {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ReduceOnly)
	}
	if order.CloseOnTrigger {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.CloseOnTrigger)
	}

	switch order.OrderStatus {
	case bybitl.OrderCreated:
		ord.OrderStatus = models.Created
	case bybitl.OrderRejected:
		ord.OrderStatus = models.Rejected
	case bybitl.OrderNew:
		ord.OrderStatus = models.New
	case bybitl.OrderPartiallyFilled:
		ord.OrderStatus = models.PartiallyFilled
	case bybitl.OrderFilled:
		ord.OrderStatus = models.Filled
	case bybitl.OrderCancelled:
		ord.OrderStatus = models.Canceled
	case bybitl.OrderPendingCancel:
		ord.OrderStatus = models.PendingCancel
	default:
		fmt.Println("UNKNOWN ORDER STATUS", order.OrderStatus)
	}

	switch order.OrderType {
	case bybitl.LimitOrder:
		ord.OrderType = models.Limit
	case bybitl.MarketOrder:
		ord.OrderType = models.Market
	default:
		fmt.Println("UNKNOWN ORDER TYPE", order.OrderType)
	}

	switch order.Side {
	case bybitl.Buy:
		ord.Side = models.Buy
	case bybitl.Sell:
		ord.Side = models.Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", order.Side)
	}

	switch order.TimeInForce {
	case bybitl.GoodTillCancel:
		ord.TimeInForce = models.GoodTillCancel
	case bybitl.ImmediateOrCancel:
		ord.TimeInForce = models.ImmediateOrCancel
	case bybitl.FillOrKill:
		ord.TimeInForce = models.FillOrKill
	case bybitl.PostOnly:
		ord.TimeInForce = models.PostOnly
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

	if order.OrderSide == models.Buy {
		side = bybitl.Buy
	} else {
		side = bybitl.Sell
	}
	switch order.OrderType {
	case models.Limit:
		typ = bybitl.LimitOrder
	case models.Market:
		typ = bybitl.MarketOrder
	default:
		rej := messages.UnsupportedOrderType
		return nil, &rej
	}

	switch order.TimeInForce {
	case models.Session:
		tif = bybitl.GoodTillCancel
	case models.GoodTillCancel:
		tif = bybitl.GoodTillCancel
	case models.ImmediateOrCancel:
		tif = bybitl.ImmediateOrCancel
	case models.FillOrKill:
		tif = bybitl.FillOrKill
	}

	for _, exec := range order.ExecutionInstructions {
		rej := messages.UnsupportedOrderCharacteristic
		switch exec {
		case models.ReduceOnly:
			redo = true
		case models.ParticipateDoNotInitiate:
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
