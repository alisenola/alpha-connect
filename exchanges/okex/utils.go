package okex

import (
	"fmt"
	"strconv"
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
	var tdMode okex.TdMode
	if order.TimeInForce == models.TimeInForce_GoodTillCancel {
		tdMode = okex.CASH
	} else if order.TimeInForce == models.TimeInForce_ImmediateOrCancel {
		tdMode = okex.ISOLATED
	}

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

	request := okex.NewPlaceOrderRequest(strings.ToLower(string(side)), strings.ToUpper(symbol), string(tdMode), string(typ), strconv.FormatFloat(size, 'f', -1, 64))
	if order.ClientOrderID != "" {
		request.SetClOrdId(order.ClientOrderID)
	}
	if order.Tag != "" {
		request.SetTag(order.Tag)
	}
	if order.Price != nil {
		request.SetPx(strconv.FormatFloat(order.Price.Value, 'f', -1, 64))
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

func WSOrderToModel(o *okex.OrderData) *models.Order {
	sz, _ := strconv.ParseFloat(o.Sz, 64)
	fillSz, _ := strconv.ParseFloat(o.FillSz, 64)
	px, _ := strconv.ParseFloat(o.Px, 64)

	ord := &models.Order{
		OrderID: o.OrdId,
		Instrument: &models.Instrument{
			Exchange: constants.OKEX,
			Symbol:   &wrapperspb.StringValue{Value: o.InstId},
		},
		LeavesQuantity: sz,
		CumQuantity:    fillSz,
	}

	if o.ReduceOnly == "true" {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ExecutionInstruction_ReduceOnly)
	}

	if sz == 0 {
		ord.OrderStatus = models.OrderStatus_Filled
	} else if fillSz > 0 {
		ord.OrderStatus = models.OrderStatus_PartiallyFilled
	} else if o.CancelSourceReason != "" {
		ord.OrderStatus = models.OrderStatus_Canceled
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
	switch o.TdMode {
	case string(okex.CASH):
		ord.TimeInForce = models.TimeInForce_GoodTillCancel
	case string(okex.ISOLATED):
		ord.TimeInForce = models.TimeInForce_ImmediateOrCancel
	case string(okex.CROSS):
		ord.TimeInForce = models.TimeInForce_GoodTillCancel
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ExecutionInstruction_ParticipateDoNotInitiate)
	default:
		fmt.Println("UNKNOWN TOF", o.TdMode)
	}

	ord.Price = &wrapperspb.DoubleValue{Value: px}
	return ord
}

func SymbolToAsset(symbol string) *xmodels.Asset {
	if sym, ok := okex.OKEX_SYMBOL_TO_GLOBAL_SYMBOL[symbol]; ok {
		symbol = sym
	}
	asset, ok := constants.GetAssetBySymbol(symbol)
	if !ok {
		return nil
	}
	return asset
}
