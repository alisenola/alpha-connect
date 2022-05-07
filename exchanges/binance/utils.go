package binance

import (
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/binance"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func WSOrderToModel(o *binance.WSExecutionReport) *models.Order {
	ord := &models.Order{
		OrderID:       fmt.Sprintf("%d", o.OrderID),
		ClientOrderID: o.ClientOrderID,
		Instrument: &models.Instrument{
			Exchange: constants.BINANCE,
			Symbol:   &wrapperspb.StringValue{Value: o.Symbol},
		},
		LeavesQuantity: o.OrigQuantity - o.CumQuantity,
		CumQuantity:    o.CumQuantity,
	}

	switch o.OrderStatus {
	case binance.NEW_ORDER:
		ord.OrderStatus = models.OrderStatus_New
	case binance.CANCELED_ORDER:
		ord.OrderStatus = models.OrderStatus_Canceled
	case binance.PARTIALLY_FILLED:
		ord.OrderStatus = models.OrderStatus_PartiallyFilled
	default:
		fmt.Println("UNKNOWN ORDER STATUS", o.OrderStatus)
	}

	/*
		const LIMIT = OrderType("LIMIT")
		const MARKET = OrderType("MARKET")
		const STOP_LOSS = OrderType("STOP_LOSS")
		const STOP_MARKET = OrderType("STOP_MARKET")
		const STOP_LOSS_LIMIT = OrderType("STOP_LOSS_LIMIT")
		const TAKE_PROFIT = OrderType("TAKE_PROFIT")
		const TAKE_PROFIT_LIMIT = OrderType("TAKE_PROFIT_LIMIT")
		const TAKE_PROFIT_MARKET = OrderType("TAKE_PROFIT_MARKET")
		const LIMIT_MAKER = OrderType("LIMIT_MAKER")
		const TRAILING_STOP_MARKET = OrderType("TRAILING_STOP_MARKET")

	*/

	switch o.OrderType {
	case binance.LIMIT:
		ord.OrderType = models.OrderType_Limit
	case binance.MARKET:
		ord.OrderType = models.OrderType_Market
	case binance.STOP_LOSS:
		ord.OrderType = models.OrderType_Stop
	case binance.STOP_LOSS_LIMIT:
		ord.OrderType = models.OrderType_StopLimit
	default:
		fmt.Println("UNKNOWN ORDER TYPE", o.OrderType)
	}

	switch o.Side {
	case binance.BUY_SIDE:
		ord.Side = models.Side_Buy
	case binance.SELL_SIDE:
		ord.Side = models.Side_Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", o.Side)
	}

	switch o.TimeInForce {
	case binance.GOOD_TILL_CANCEL:
		ord.TimeInForce = models.TimeInForce_GoodTillCancel
	case binance.FILL_OR_KILL:
		ord.TimeInForce = models.TimeInForce_FillOrKill
	case binance.IMMEDIATE_OR_CANCEL:
		ord.TimeInForce = models.TimeInForce_ImmediateOrCancel
	case binance.GOOD_TILL_CROSSING:
		ord.TimeInForce = models.TimeInForce_GoodTillCancel
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ExecutionInstruction_ParticipateDoNotInitiate)
	default:
		fmt.Println("UNKNOWN TOF", o.TimeInForce)
	}

	ord.Price = &wrapperspb.DoubleValue{Value: o.OrigPrice}
	return ord
}

func OrderToModel(o *binance.OrderData) *models.Order {
	ord := &models.Order{
		OrderID:       fmt.Sprintf("%d", o.OrderId),
		ClientOrderID: o.ClientOrderId,
		Instrument: &models.Instrument{
			Exchange: constants.BINANCE,
			Symbol:   &wrapperspb.StringValue{Value: o.Symbol},
		},
		LeavesQuantity: o.OrigQty - o.ExecutedQty, // TODO check
		CumQuantity:    o.OrigQty,
	}

	switch o.Status {
	case binance.NEW_ORDER:
		ord.OrderStatus = models.OrderStatus_New
	case binance.CANCELED_ORDER:
		ord.OrderStatus = models.OrderStatus_Canceled
	case binance.PARTIALLY_FILLED:
		ord.OrderStatus = models.OrderStatus_PartiallyFilled
	default:
		fmt.Println("UNKNOWN ORDER STATUS", o.Status)
	}

	/*
		const LIMIT = OrderType("LIMIT")
		const MARKET = OrderType("MARKET")
		const STOP_LOSS = OrderType("STOP_LOSS")
		const STOP_MARKET = OrderType("STOP_MARKET")
		const STOP_LOSS_LIMIT = OrderType("STOP_LOSS_LIMIT")
		const TAKE_PROFIT = OrderType("TAKE_PROFIT")
		const TAKE_PROFIT_LIMIT = OrderType("TAKE_PROFIT_LIMIT")
		const TAKE_PROFIT_MARKET = OrderType("TAKE_PROFIT_MARKET")
		const LIMIT_MAKER = OrderType("LIMIT_MAKER")
		const TRAILING_STOP_MARKET = OrderType("TRAILING_STOP_MARKET")

	*/

	switch o.Type {
	case binance.LIMIT:
		ord.OrderType = models.OrderType_Limit
	case binance.MARKET:
		ord.OrderType = models.OrderType_Market
	case binance.STOP_LOSS:
		ord.OrderType = models.OrderType_Stop
	case binance.STOP_LOSS_LIMIT:
		ord.OrderType = models.OrderType_StopLimit
	default:
		fmt.Println("UNKNOWN ORDER TYPE", o.Type)
	}

	switch o.Side {
	case binance.BUY_SIDE:
		ord.Side = models.Side_Buy
	case binance.SELL_SIDE:
		ord.Side = models.Side_Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", o.Side)
	}

	switch o.TimeInForce {
	case binance.GOOD_TILL_CANCEL:
		ord.TimeInForce = models.TimeInForce_GoodTillCancel
	case binance.FILL_OR_KILL:
		ord.TimeInForce = models.TimeInForce_FillOrKill
	case binance.IMMEDIATE_OR_CANCEL:
		ord.TimeInForce = models.TimeInForce_ImmediateOrCancel
	case binance.GOOD_TILL_CROSSING:
		ord.TimeInForce = models.TimeInForce_GoodTillCancel
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ExecutionInstruction_ParticipateDoNotInitiate)
	default:
		fmt.Println("UNKNOWN TOF", o.TimeInForce)
	}

	ord.Price = &wrapperspb.DoubleValue{Value: o.Price}
	return ord
}

func buildPostOrderRequest(symbol string, order *messages.NewOrder, tickPrecision, lotPrecision int) (binance.NewOrderRequest, *messages.RejectionReason) {
	var side binance.OrderSide
	var typ binance.OrderType
	if order.OrderSide == models.Side_Buy {
		side = binance.BUY_SIDE
	} else {
		side = binance.SELL_SIDE
	}
	switch order.OrderType {
	case models.OrderType_Limit:
		typ = binance.LIMIT
	case models.OrderType_Market:
		typ = binance.MARKET
	case models.OrderType_Stop:
		typ = binance.STOP_LOSS
	case models.OrderType_StopLimit:
		typ = binance.STOP_LOSS_LIMIT
	default:
		rej := messages.RejectionReason_UnsupportedOrderType
		return nil, &rej
	}
	tif := binance.GOOD_TILL_CANCEL
	for _, exec := range order.ExecutionInstructions {
		rej := messages.RejectionReason_UnsupportedOrderCharacteristic
		switch exec {
		case models.ExecutionInstruction_ParticipateDoNotInitiate:
			tif = binance.GOOD_TILL_CROSSING
		default:
			return nil, &rej
		}
	}

	request := binance.NewNewOrderRequest(symbol, side, typ, tif)

	request.SetQuantity(order.Quantity, lotPrecision)
	request.SetNewClientOrderID(order.ClientOrderID)

	if order.OrderType != models.OrderType_Market {
		switch order.TimeInForce {
		case models.TimeInForce_Session:
			request.SetTimeInForce(binance.GOOD_TILL_CANCEL)
		case models.TimeInForce_GoodTillCancel:
			request.SetTimeInForce(binance.GOOD_TILL_CANCEL)
		case models.TimeInForce_ImmediateOrCancel:
			request.SetTimeInForce(binance.IMMEDIATE_OR_CANCEL)
		case models.TimeInForce_FillOrKill:
			request.SetTimeInForce(binance.FILL_OR_KILL)
		default:
			rej := messages.RejectionReason_UnsupportedOrderTimeInForce
			return nil, &rej
		}
	}

	if order.Price != nil {
		request.SetPrice(order.Price.Value, tickPrecision)
	}

	request.SetNewOrderResponseType(binance.ACK_RESPONSE_TYPE)

	return request, nil
}
