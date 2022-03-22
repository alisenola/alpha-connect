package ftx

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/ftx"
)

func WSOrderToModel(o ftx.WSOrder) *models.Order {
	ord := &models.Order{
		OrderID: fmt.Sprintf("%d", o.ID),
		Instrument: &models.Instrument{
			Exchange: &constants.FTX,
			Symbol:   &types.StringValue{Value: o.Market},
		},
		LeavesQuantity: o.Size - o.FilledSize, // TODO check
		CumQuantity:    o.FilledSize,
	}
	if o.ClientID != nil {
		ord.ClientOrderID = *o.ClientID
	}

	if o.ReduceOnly {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ReduceOnly)
	}
	if o.PostOnly {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ParticipateDoNotInitiate)
	}

	switch o.Status {
	case ftx.NEW_ORDER:
		ord.OrderStatus = models.PendingNew
	case ftx.OPEN_ORDER:
		if o.FilledSize > 0 {
			ord.OrderStatus = models.PartiallyFilled
		} else {
			ord.OrderStatus = models.New
		}
	case ftx.CLOSED_ORDER:
		if o.FilledSize == o.Size {
			ord.OrderStatus = models.Filled
		} else {
			ord.OrderStatus = models.Canceled
		}
	default:
		fmt.Println("unknown ORDER STATUS", o.Status)
	}

	switch o.Type {
	case ftx.LIMIT_ORDER:
		ord.OrderType = models.Limit
	case ftx.MARKET_ORDER:
		ord.OrderType = models.Market
	default:
		fmt.Println("UNKNOWN ORDER TYPE", o.Type)
	}

	switch o.Side {
	case ftx.BUY:
		ord.Side = models.Buy
	case ftx.SELL:
		ord.Side = models.Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", o.Side)
	}

	if o.IOC {
		ord.TimeInForce = models.ImmediateOrCancel
	} else {
		ord.TimeInForce = models.GoodTillCancel
	}
	ord.Price = &types.DoubleValue{Value: o.Price}

	return ord
}

func OrderToModel(o ftx.Order) *models.Order {
	ord := &models.Order{
		OrderID: fmt.Sprintf("%d", o.ID),
		Instrument: &models.Instrument{
			Exchange: &constants.FTX,
			Symbol:   &types.StringValue{Value: o.Market},
		},
		LeavesQuantity: o.Size - o.FilledSize, // TODO check
		CumQuantity:    o.FilledSize,
	}
	ord.ClientOrderID = o.ClientID

	if o.ReduceOnly {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ReduceOnly)
	}
	if o.PostOnly {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ParticipateDoNotInitiate)
	}

	switch o.Status {
	case ftx.NEW_ORDER:
		ord.OrderStatus = models.PendingNew
	case ftx.OPEN_ORDER:
		if o.FilledSize > 0 {
			ord.OrderStatus = models.PartiallyFilled
		} else {
			ord.OrderStatus = models.New
		}
	case ftx.CLOSED_ORDER:
		if o.FilledSize == o.Size {
			ord.OrderStatus = models.Filled
		} else {
			ord.OrderStatus = models.Canceled
		}
	default:
		fmt.Println("unknown ORDER STATUS", o.Status)
	}

	switch o.Type {
	case ftx.LIMIT_ORDER:
		ord.OrderType = models.Limit
	case ftx.MARKET_ORDER:
		ord.OrderType = models.Market
	default:
		fmt.Println("UNKNOWN ORDER TYPE", o.Type)
	}

	switch o.Side {
	case ftx.BUY:
		ord.Side = models.Buy
	case ftx.SELL:
		ord.Side = models.Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", o.Side)
	}

	if o.Ioc {
		ord.TimeInForce = models.ImmediateOrCancel
	} else {
		ord.TimeInForce = models.GoodTillCancel
	}
	ord.Price = &types.DoubleValue{Value: o.Price}

	return ord
}
