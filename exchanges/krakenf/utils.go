package krakenf

import (
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/models"
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
