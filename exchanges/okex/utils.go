package okex

import (
	"fmt"
	"strings"

	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/okex"
	xmodels "gitlab.com/alphaticks/xchanger/models"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

func WSPositionToModel(p *okex.WSPosition) *models.Position {
	pos := &models.Position{
		Instrument: &models.Instrument{
			Exchange: constants.OKEX,
			Symbol:   &wrapperspb.StringValue{Value: p.InstId},
		},
		Quantity:  p.Pos,
		Cost:      p.Pos * p.AvgPx,
		Cross:     false,
		MarkPrice: &wrapperspb.DoubleValue{Value: p.MarkPx},
	}

	return pos
}

func WSOrderToModel(o *okex.WSOrder) *models.Order {

	ord := &models.Order{
		OrderID: o.OrdId,
		Instrument: &models.Instrument{
			Exchange: constants.OKEX,
			Symbol:   &wrapperspb.StringValue{Value: o.InstId},
		},
		LeavesQuantity: o.Sz,
		CumQuantity:    o.FillSz,
	}

	if o.ReduceOnly {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ExecutionInstruction_ReduceOnly)
	}

	if o.Sz == 0. {
		ord.OrderStatus = models.OrderStatus_Filled
	} else if o.FillSz > 0. {
		ord.OrderStatus = models.OrderStatus_PartiallyFilled
	} else {
		ord.OrderStatus = models.OrderStatus_New
	}

	switch o.OrdType {
	case okex.LimitOrderType:
		ord.OrderType = models.OrderType_Limit
	case okex.MarketOrderType:
		ord.OrderType = models.OrderType_Market
	default:
		fmt.Println("UNKNOWN ORDER TYPE", o.OrdType)
		ord.OrderType = models.OrderType_Limit
	}

	switch o.Side {
	case string(okex.BUY_ORDER):
		ord.Side = models.Side_Buy
	case string(okex.SELL_ODER):
		ord.Side = models.Side_Sell
	}

	ord.Price = &wrapperspb.DoubleValue{Value: o.Px}
	return ord
}

func SymbolToAsset(symbol string) *xmodels.Asset {
	asset, ok := constants.GetAssetBySymbol(symbol)
	if !ok {
		return nil
	}
	return asset
}
