package binance

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/binance"
)

func wsOrderToModel(o *binance.WSExecutionReport) *models.Order {
	ord := &models.Order{
		OrderID:       fmt.Sprintf("%d", o.OrderID),
		ClientOrderID: o.ClientOrderID,
		Instrument: &models.Instrument{
			Exchange: &constants.BINANCE,
			Symbol:   &types.StringValue{Value: o.Symbol},
		},
		LeavesQuantity: o.OrigQuantity - o.CumQuantity,
		CumQuantity:    o.CumQuantity,
	}

	switch o.OrderStatus {
	case binance.NEW_ORDER:
		ord.OrderStatus = models.New
	case binance.CANCELED_ORDER:
		ord.OrderStatus = models.Canceled
	case binance.PARTIALLY_FILLED:
		ord.OrderStatus = models.PartiallyFilled
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
		ord.OrderType = models.Limit
	case binance.MARKET:
		ord.OrderType = models.Market
	case binance.STOP_LOSS:
		ord.OrderType = models.Stop
	case binance.STOP_LOSS_LIMIT:
		ord.OrderType = models.StopLimit
	default:
		fmt.Println("UNKNOWN ORDER TYPE", o.OrderType)
	}

	switch o.Side {
	case binance.BUY_SIDE:
		ord.Side = models.Buy
	case binance.SELL_SIDE:
		ord.Side = models.Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", o.Side)
	}

	switch o.TimeInForce {
	case binance.GOOD_TILL_CANCEL:
		ord.TimeInForce = models.GoodTillCancel
	case binance.FILL_OR_KILL:
		ord.TimeInForce = models.FillOrKill
	case binance.IMMEDIATE_OR_CANCEL:
		ord.TimeInForce = models.ImmediateOrCancel
	case binance.GOOD_TILL_CROSSING:
		ord.TimeInForce = models.GoodTillCancel
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ParticipateDoNotInitiate)
	default:
		fmt.Println("UNKNOWN TOF", o.TimeInForce)
	}

	ord.Price = &types.DoubleValue{Value: o.OrigPrice}
	return ord
}

func orderToModel(o *binance.OrderData) *models.Order {
	ord := &models.Order{
		OrderID:       fmt.Sprintf("%d", o.OrderID),
		ClientOrderID: o.ClientOrderID,
		Instrument: &models.Instrument{
			Exchange: &constants.BINANCE,
			Symbol:   &types.StringValue{Value: o.Symbol},
		},
		LeavesQuantity: o.OrigQty - o.ExecutedQty, // TODO check
		CumQuantity:    o.OrigQty,
	}

	switch o.Status {
	case binance.NEW_ORDER:
		ord.OrderStatus = models.New
	case binance.CANCELED_ORDER:
		ord.OrderStatus = models.Canceled
	case binance.PARTIALLY_FILLED:
		ord.OrderStatus = models.PartiallyFilled
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
		ord.OrderType = models.Limit
	case binance.MARKET:
		ord.OrderType = models.Market
	case binance.STOP_LOSS:
		ord.OrderType = models.Stop
	case binance.STOP_LOSS_LIMIT:
		ord.OrderType = models.StopLimit
	default:
		fmt.Println("UNKNOWN ORDER TYPE", o.Type)
	}

	switch o.Side {
	case binance.BUY_SIDE:
		ord.Side = models.Buy
	case binance.SELL_SIDE:
		ord.Side = models.Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", o.Side)
	}

	switch o.TimeInForce {
	case binance.GOOD_TILL_CANCEL:
		ord.TimeInForce = models.GoodTillCancel
	case binance.FILL_OR_KILL:
		ord.TimeInForce = models.FillOrKill
	case binance.IMMEDIATE_OR_CANCEL:
		ord.TimeInForce = models.ImmediateOrCancel
	case binance.GOOD_TILL_CROSSING:
		ord.TimeInForce = models.GoodTillCancel
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ParticipateDoNotInitiate)
	default:
		fmt.Println("UNKNOWN TOF", o.TimeInForce)
	}

	ord.Price = &types.DoubleValue{Value: o.Price}
	return ord
}

func buildPostOrderRequest(symbol string, order *messages.NewOrder, tickPrecision, lotPrecision int) (binance.NewOrderRequest, *messages.RejectionReason) {
	var side binance.OrderSide
	var typ binance.OrderType
	if order.OrderSide == models.Buy {
		side = binance.BUY_SIDE
	} else {
		side = binance.SELL_SIDE
	}
	switch order.OrderType {
	case models.Limit:
		typ = binance.LIMIT
	case models.Market:
		typ = binance.MARKET
	case models.Stop:
		typ = binance.STOP_LOSS
	case models.StopLimit:
		typ = binance.STOP_LOSS_LIMIT
	default:
		rej := messages.UnsupportedOrderType
		return nil, &rej
	}
	tif := binance.GOOD_TILL_CANCEL
	for _, exec := range order.ExecutionInstructions {
		rej := messages.UnsupportedOrderCharacteristic
		switch exec {
		case models.ParticipateDoNotInitiate:
			tif = binance.GOOD_TILL_CROSSING
		default:
			return nil, &rej
		}
	}

	request := binance.NewNewOrderRequest(symbol, side, typ, tif)

	request.SetQuantity(order.Quantity, lotPrecision)
	request.SetNewClientOrderID(order.ClientOrderID)

	if order.OrderType != models.Market {
		switch order.TimeInForce {
		case models.Session:
			request.SetTimeInForce(binance.GOOD_TILL_CANCEL)
		case models.GoodTillCancel:
			request.SetTimeInForce(binance.GOOD_TILL_CANCEL)
		case models.ImmediateOrCancel:
			request.SetTimeInForce(binance.IMMEDIATE_OR_CANCEL)
		case models.FillOrKill:
			request.SetTimeInForce(binance.FILL_OR_KILL)
		default:
			rej := messages.UnsupportedOrderTimeInForce
			return nil, &rej
		}
	}

	if order.Price != nil {
		request.SetPrice(order.Price.Value, tickPrecision)
	}

	request.SetNewOrderResponseType(binance.ACK_RESPONSE_TYPE)

	return request, nil
}
