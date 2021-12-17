package bitmex

import (
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/exchanges/bitmex"
)

func buildPostOrderRequest(order *messages.NewOrder) (bitmex.PostOrderRequest, *messages.RejectionReason) {
	request := bitmex.NewPostOrderRequest(order.Instrument.Symbol.Value)

	request.SetOrderQty(order.Quantity)
	request.SetClOrdID(order.ClientOrderID)

	switch order.OrderSide {
	case models.Buy:
		request.SetSide(bitmex.BUY_ORDER_SIDE)
	case models.Sell:
		request.SetSide(bitmex.SELL_ORDER_SIDE)
	default:
		request.SetSide(bitmex.BUY_ORDER_SIDE)
	}

	switch order.OrderType {
	case models.Limit:
		request.SetOrderType(bitmex.LIMIT)
	case models.Market:
		request.SetOrderType(bitmex.MARKET)
	case models.Stop:
		request.SetOrderType(bitmex.STOP)
	case models.StopLimit:
		request.SetOrderType(bitmex.STOP_LIMIT)
	case models.LimitIfTouched:
		request.SetOrderType(bitmex.LIMIT_IF_TOUCHED)
	case models.MarketIfTouched:
		request.SetOrderType(bitmex.MARKET_IF_TOUCHED)
	default:
		rej := messages.UnsupportedOrderType
		return nil, &rej
	}

	switch order.TimeInForce {
	case models.Session:
		request.SetTimeInForce(bitmex.TIF_DAY)
	case models.GoodTillCancel:
		request.SetTimeInForce(bitmex.GOOD_TILL_CANCEL)
	case models.ImmediateOrCancel:
		request.SetTimeInForce(bitmex.IMMEDIATE_OR_CANCEL)
	case models.FillOrKill:
		request.SetTimeInForce(bitmex.FILL_OR_KILL)
	default:
		rej := messages.UnsupportedOrderTimeInForce
		return nil, &rej
	}

	if order.Price != nil {
		request.SetPrice(order.Price.Value)
	}

	// TODO handle multiple exec inst
	if len(order.ExecutionInstructions) > 0 {
		switch order.ExecutionInstructions[0] {
		case models.ParticipateDoNotInitiate:
			request.SetExecInst(bitmex.EI_PARTICIPATE_DO_NOT_INITIATE)
		}
	}

	return request, nil
}
