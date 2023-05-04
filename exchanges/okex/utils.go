package okex

import (
	"fmt"
	"strings"

	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/exchanges/okex"
)

func buildPostOrderRequest(symbol string, order *messages.NewOrder) (okex.PlaceOrderRequest, *messages.RejectionReason) {
	var typ okex.OrderType
	size := order.Quantity
	tdMode := "cash"
	var side okex.OrderSide
	if order.OrderSide == models.Side_Buy {
		side = okex.BUY_ORDER
	} else {
		side = okex.SELL_ODER
	}
	switch order.OrderType {
	case models.OrderType_Limit:
		typ = okex.LimitOrderType
	case models.OrderType_Market:
		typ = okex.MarketOrderType
	default:
		rej := messages.RejectionReason_UnsupportedOrderType
		return nil, &rej
	}

	request := okex.NewPlaceOrderRequest(string(side), strings.ToLower(symbol), tdMode, string(typ), fmt.Sprintf("%.6f", size))
	if order.ClientOrderID != "" {
		request.SetClOrdId(order.ClientOrderID)
	}
	if order.Tag != "" {
		request.SetTag(order.Tag)
	}

	return request, nil
}
