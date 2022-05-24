package ftx

import (
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/ftx"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func WSOrderToModel(o ftx.WSOrder) *models.Order {
	ord := &models.Order{
		OrderID: fmt.Sprintf("%d", o.ID),
		Instrument: &models.Instrument{
			Exchange: constants.FTX,
			Symbol:   &wrapperspb.StringValue{Value: o.Market},
		},
		LeavesQuantity: o.Size - o.FilledSize, // TODO check
		CumQuantity:    o.FilledSize,
	}
	if o.ClientID != nil {
		ord.ClientOrderID = *o.ClientID
	}

	if o.ReduceOnly {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ExecutionInstruction_ReduceOnly)
	}
	if o.PostOnly {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ExecutionInstruction_ParticipateDoNotInitiate)
	}

	switch o.Status {
	case ftx.NEW_ORDER:
		ord.OrderStatus = models.OrderStatus_PendingNew
	case ftx.OPEN_ORDER:
		if o.FilledSize > 0 {
			ord.OrderStatus = models.OrderStatus_PartiallyFilled
		} else {
			ord.OrderStatus = models.OrderStatus_New
		}
	case ftx.CLOSED_ORDER:
		if o.FilledSize == o.Size {
			ord.OrderStatus = models.OrderStatus_Filled
		} else {
			ord.OrderStatus = models.OrderStatus_Canceled
		}
	default:
		fmt.Println("unknown ORDER STATUS", o.Status)
	}

	switch o.Type {
	case ftx.LIMIT_ORDER:
		ord.OrderType = models.OrderType_Limit
	case ftx.MARKET_ORDER:
		ord.OrderType = models.OrderType_Market
	default:
		fmt.Println("UNKNOWN ORDER TYPE", o.Type)
	}

	switch o.Side {
	case ftx.BUY:
		ord.Side = models.Side_Buy
	case ftx.SELL:
		ord.Side = models.Side_Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", o.Side)
	}

	if o.IOC {
		ord.TimeInForce = models.TimeInForce_ImmediateOrCancel
	} else {
		ord.TimeInForce = models.TimeInForce_GoodTillCancel
	}
	ord.Price = &wrapperspb.DoubleValue{Value: o.Price}

	return ord
}

func OrderToModel(o ftx.Order) *models.Order {
	ord := &models.Order{
		OrderID: fmt.Sprintf("%d", o.ID),
		Instrument: &models.Instrument{
			Exchange: constants.FTX,
			Symbol:   &wrapperspb.StringValue{Value: o.Market},
		},
		LeavesQuantity: o.Size - o.FilledSize, // TODO check
		CumQuantity:    o.FilledSize,
	}
	ord.ClientOrderID = o.ClientID

	if o.ReduceOnly {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ExecutionInstruction_ReduceOnly)
	}
	if o.PostOnly {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ExecutionInstruction_ParticipateDoNotInitiate)
	}

	switch o.Status {
	case ftx.NEW_ORDER:
		ord.OrderStatus = models.OrderStatus_PendingNew
	case ftx.OPEN_ORDER:
		if o.FilledSize > 0 {
			ord.OrderStatus = models.OrderStatus_PartiallyFilled
		} else {
			ord.OrderStatus = models.OrderStatus_New
		}
	case ftx.CLOSED_ORDER:
		if o.FilledSize == o.Size {
			ord.OrderStatus = models.OrderStatus_Filled
		} else {
			ord.OrderStatus = models.OrderStatus_Canceled
		}
	default:
		fmt.Println("unknown ORDER STATUS", o.Status)
	}

	switch o.Type {
	case ftx.LIMIT_ORDER:
		ord.OrderType = models.OrderType_Limit
	case ftx.MARKET_ORDER:
		ord.OrderType = models.OrderType_Market
	default:
		fmt.Println("UNKNOWN ORDER TYPE", o.Type)
	}

	switch o.Side {
	case ftx.BUY:
		ord.Side = models.Side_Buy
	case ftx.SELL:
		ord.Side = models.Side_Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", o.Side)
	}

	if o.Ioc {
		ord.TimeInForce = models.TimeInForce_ImmediateOrCancel
	} else {
		ord.TimeInForce = models.TimeInForce_GoodTillCancel
	}
	ord.Price = &wrapperspb.DoubleValue{Value: o.Price}

	return ord
}

func buildPlaceOrderRequest(symbol string, order *messages.NewOrder, tickPrecision, lotPrecision int) (ftx.NewOrderRequest, *messages.RejectionReason) {
	request := ftx.NewOrderRequest{
		Market: symbol,
	}

	if order.OrderSide == models.Side_Buy {
		request.Side = ftx.BUY
	} else {
		request.Side = ftx.SELL
	}
	switch order.OrderType {
	case models.OrderType_Limit:
		request.Type = ftx.LIMIT_ORDER
	case models.OrderType_Market:
		request.Type = ftx.MARKET_ORDER
	default:
		rej := messages.RejectionReason_UnsupportedOrderType
		return request, &rej
	}

	request.Size = order.Quantity
	request.ClientID = order.ClientOrderID

	if order.OrderType != models.OrderType_Market {
		switch order.TimeInForce {
		case models.TimeInForce_ImmediateOrCancel:
			request.Ioc = true
		case models.TimeInForce_GoodTillCancel:
			request.Ioc = false
		default:
			rej := messages.RejectionReason_UnsupportedOrderTimeInForce
			return request, &rej
		}
	}

	if order.Price != nil {
		request.Price = &order.Price.Value
	}

	for _, exec := range order.ExecutionInstructions {
		switch exec {
		case models.ExecutionInstruction_ReduceOnly:
			request.ReduceOnly = true
		case models.ExecutionInstruction_ParticipateDoNotInitiate:
			request.PostOnly = true
		default:
			rej := messages.RejectionReason_UnsupportedOrderCharacteristic
			return request, &rej
		}
	}

	return request, nil
}
