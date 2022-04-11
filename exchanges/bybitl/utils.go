package bybitl

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/bybitl"
)

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
	case bybitl.Limit:
		o.OrderType = models.Limit
	case bybitl.Market:
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
