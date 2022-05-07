package dydx

import (
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/dydx"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"
)

func statusToDydx(status models.OrderStatus) dydx.OrderStatus {
	switch status {
	case models.OrderStatus_New:
		return dydx.ORDER_OPEN
	case models.OrderStatus_PendingNew, models.OrderStatus_PendingCancel, models.OrderStatus_PendingReplace:
		return dydx.ORDER_PENDING
	case models.OrderStatus_Filled:
		return dydx.ORDER_FILLED
	case models.OrderStatus_Canceled:
		return dydx.ORDER_CANCELED
	default:
		return ""
	}
}

func sideToDydx(side models.Side) dydx.OrderSide {
	switch side {
	case models.Side_Buy:
		return dydx.BUY
	case models.Side_Sell:
		return dydx.SELL
	default:
		return ""
	}
}

func orderToModel(o dydx.Order) *models.Order {
	ord := &models.Order{
		OrderID:       o.ID,
		ClientOrderID: o.ClientID,
		Instrument: &models.Instrument{
			Exchange: constants.DYDX,
			Symbol:   &wrapperspb.StringValue{Value: o.Market},
		},
		LeavesQuantity: o.RemainingSize, // TODO check
		CumQuantity:    o.Size - o.RemainingSize,
	}

	switch o.Status {
	case dydx.ORDER_OPEN:
		if o.Size != o.RemainingSize {
			ord.OrderStatus = models.OrderStatus_PartiallyFilled
		} else {
			ord.OrderStatus = models.OrderStatus_New
		}
	case dydx.ORDER_FILLED:
		ord.OrderStatus = models.OrderStatus_Filled
	case dydx.ORDER_CANCELED:
		ord.OrderStatus = models.OrderStatus_Canceled
	case dydx.ORDER_PENDING:
		ord.OrderStatus = models.OrderStatus_PendingNew
	default:
		fmt.Println("UNKNOWN ORDER STATUS", o.Status)
	}

	switch o.Type {
	case dydx.LIMIT:
		ord.OrderType = models.OrderType_Limit
	case dydx.MARKET:
		ord.OrderType = models.OrderType_Market
	case dydx.STOP_LIMIT:
		ord.OrderType = models.OrderType_StopLimit
	case dydx.TAKE_PROFIT, dydx.TRAILING_STOP:
		ord.OrderType = models.OrderType_TrailingStopLimit
	default:
		fmt.Println("UNKNOWN ORDER TYPE", o.Type)
	}

	switch o.Side {
	case dydx.BUY:
		ord.Side = models.Side_Buy
	case dydx.SELL:
		ord.Side = models.Side_Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", o.Side)
	}

	switch o.TimeInForce {
	case dydx.GOOD_TIL_TIME:
		ord.TimeInForce = models.TimeInForce_GoodTillCancel
	case dydx.FILL_OR_KILL:
		ord.TimeInForce = models.TimeInForce_FillOrKill
	case dydx.IMMEDIATE_OR_CANCEL:
		ord.TimeInForce = models.TimeInForce_ImmediateOrCancel
	default:
		fmt.Println("UNKNOWN TOF", o.TimeInForce)
	}

	ord.Price = &wrapperspb.DoubleValue{Value: o.Price}
	return ord
}

func buildCreateOrderParams(symbol string, o *messages.NewOrder) *dydx.CreateOrderParams {
	var side dydx.OrderSide
	var typ dydx.OrderType
	var tif dydx.OrderTIF
	var postOnly = false
	var price float64
	if o.OrderSide == models.Side_Buy {
		side = dydx.BUY
	} else {
		side = dydx.SELL
	}
	switch o.OrderType {
	case models.OrderType_Limit:
		typ = dydx.LIMIT
	case models.OrderType_Market:
		typ = dydx.MARKET
	case models.OrderType_StopLimit:
		typ = dydx.STOP_LIMIT
	case models.OrderType_TrailingStopLimit:
		typ = dydx.TRAILING_STOP
	}
	if o.OrderType != models.OrderType_Market {
		switch o.TimeInForce {
		case models.TimeInForce_GoodTillCancel:
			tif = dydx.GOOD_TIL_TIME
		case models.TimeInForce_FillOrKill:
			tif = dydx.FILL_OR_KILL
		case models.TimeInForce_ImmediateOrCancel:
			tif = dydx.IMMEDIATE_OR_CANCEL
		}
	} else {
		tif = dydx.IMMEDIATE_OR_CANCEL
	}

	for _, exec := range o.ExecutionInstructions {
		if exec == models.ExecutionInstruction_ParticipateDoNotInitiate {
			postOnly = true
		}
	}
	if o.Price != nil {
		price = o.Price.Value
	}
	params := dydx.NewCreateOrderParams(
		symbol, side, typ, tif, postOnly, 0.1, price, o.Quantity, o.ClientOrderID, time.Now().Add(time.Hour))
	return params
}
