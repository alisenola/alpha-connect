package krakenf

import (
	"fmt"

	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/krakenf"
	xmodels "gitlab.com/alphaticks/xchanger/models"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func WSOrderToModel(o *krakenf.Order) *models.Order {

	ord := &models.Order{
		OrderID: o.OrderId,
		Instrument: &models.Instrument{
			Exchange: constants.FBINANCE,
			Symbol:   &wrapperspb.StringValue{Value: o.Instrument},
		},
		LeavesQuantity: o.Qty,
		CumQuantity:    o.Filled,
	}

	if o.ReduceOnly {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ExecutionInstruction_ReduceOnly)
	}

	if o.Qty == 0. {
		ord.OrderStatus = models.OrderStatus_Filled
	} else if o.Filled > 0. {
		ord.OrderStatus = models.OrderStatus_PartiallyFilled
	} else {
		ord.OrderStatus = models.OrderStatus_New
	}

	/*
		LimitOrderType        = "lmt"
		PostOnlyOrderType     = "post"
		IOCOrderType          = "ioc"
		MarketOrderType       = "mkt"
		StopOrderType         = "stp"
		TakeProfitOrderType   = "take_profit"
		TrailingStopOrderType = "trailing_stop"
	*/

	switch o.Type {
	case krakenf.WSOrderTypeLimit:
		ord.OrderType = models.OrderType_Limit
	case krakenf.WSOrderTypeStop:
		ord.OrderType = models.OrderType_Stop
	default:
		fmt.Println("UNKNOWN ORDER TYPE", o.Type)
		ord.OrderType = models.OrderType_Limit
	}

	switch o.Direction {
	case krakenf.BuyDirection:
		ord.Side = models.Side_Buy
	case krakenf.SellDirection:
		ord.Side = models.Side_Sell
	}

	ord.Price = &wrapperspb.DoubleValue{Value: o.LimitPrice}
	return ord
}

func WSPositionToModel(p *krakenf.Position) *models.Position {
	pos := &models.Position{
		Instrument: &models.Instrument{
			Exchange: constants.KRAKENF,
			Symbol:   &wrapperspb.StringValue{Value: p.Instrument},
		},
		Quantity:  p.Balance,
		Cost:      p.Balance * p.EntryPrice,
		Cross:     false,
		MarkPrice: &wrapperspb.DoubleValue{Value: p.MarkPrice},
	}

	return pos
}

func SymbolToAsset(symbol string) *xmodels.Asset {
	if sym, ok := krakenf.KRAKENF_SYMBOL_TO_GLOBAL_SYMBOL[symbol]; ok {
		symbol = sym
	}
	asset, ok := constants.GetAssetBySymbol(symbol)
	if !ok {
		return nil
	}
	return asset
}

func StatusToModel(status string) *models.OrderStatus {
	switch status {
	case "NEW":
		v := models.OrderStatus_New
		return &v
	case "CANCELED":
		v := models.OrderStatus_Canceled
		return &v
	case "PARTIALLY_FILLED":
		v := models.OrderStatus_PartiallyFilled
		return &v
	case "FILLED":
		v := models.OrderStatus_Filled
		return &v
	case "EXPIRED":
		v := models.OrderStatus_Expired
		return &v
	default:
		return nil
	}
}

func buildPostOrderRequest(symbol string, order *messages.NewOrder, tickPrecision, lotPrecision int) (krakenf.SendOrderRequest, *messages.RejectionReason) {
	var side string
	var typ krakenf.OrderType
	size := order.Size
	if order.OrderSide == models.Side_Buy {
		side = "BUY"
	} else {
		side = "SELL"
	}
	switch order.OrderType {
	case models.OrderType_Limit:
		typ = krakenf.LimitOrderType
	case models.OrderType_Market:
		typ = krakenf.MarketOrderType
	case models.OrderType_Stop:
		typ = krakenf.StopOrderType
	case models.OrderType_IOC:
		typ = krakenf.IOCOrderType
	case models.OrderType_PostOnly:
		typ = krakenf.PostOnlyOrderType
	case models.OrderType_TakeProfit:
		typ = krakenf.TakeProfitOrderType
	case models.OrderType_TrailingStopLimit:
		typ = krakenf.TrailingStopOrderType
	default:
		rej := messages.RejectionReason_UnsupportedOrderType
		return nil, &rej
	}

	request := krakenf.NewSendOrderRequest(symbol, string(typ), size, side)

	return request, nil
}
