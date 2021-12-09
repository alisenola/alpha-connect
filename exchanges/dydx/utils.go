package dydx

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/dydx"
	"time"
)

func statusToDydx(status models.OrderStatus) dydx.OrderStatus {
	switch status {
	case models.New:
		return dydx.ORDER_OPEN
	case models.PendingNew, models.PendingCancel, models.PendingReplace:
		return dydx.ORDER_PENDING
	case models.Filled:
		return dydx.ORDER_FILLED
	case models.Canceled:
		return dydx.ORDER_CANCELED
	default:
		return ""
	}
}

func sideToDydx(side models.Side) dydx.OrderSide {
	switch side {
	case models.Buy:
		return dydx.BUY
	case models.Sell:
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
			Exchange: &constants.DYDX,
			Symbol:   &types.StringValue{Value: o.Market},
		},
		LeavesQuantity: o.RemainingSize, // TODO check
		CumQuantity:    o.Size - o.RemainingSize,
	}

	switch o.Status {
	case dydx.ORDER_OPEN:
		if o.Size != o.RemainingSize {
			ord.OrderStatus = models.PartiallyFilled
		} else {
			ord.OrderStatus = models.New
		}
	case dydx.ORDER_FILLED:
		ord.OrderStatus = models.Filled
	case dydx.ORDER_CANCELED:
		ord.OrderStatus = models.Canceled
	case dydx.ORDER_PENDING:
		ord.OrderStatus = models.PendingNew
	default:
		fmt.Println("UNKNOWN ORDER STATUS", o.Status)
	}

	switch o.Type {
	case dydx.LIMIT:
		ord.OrderType = models.Limit
	case dydx.MARKET:
		ord.OrderType = models.Market
	case dydx.STOP_LIMIT:
		ord.OrderType = models.StopLimit
	case dydx.TAKE_PROFIT, dydx.TRAILING_STOP:
		ord.OrderType = models.TrailingStopLimit
	default:
		fmt.Println("UNKNOWN ORDER TYPE", o.Type)
	}

	switch o.Side {
	case dydx.BUY:
		ord.Side = models.Buy
	case dydx.SELL:
		ord.Side = models.Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", o.Side)
	}

	switch o.TimeInForce {
	case dydx.GOOD_TIL_TIME:
		ord.TimeInForce = models.GoodTillCancel
	case dydx.FILL_OR_KILL:
		ord.TimeInForce = models.FillOrKill
	case dydx.IMMEDIATE_OR_CANCEL:
		ord.TimeInForce = models.ImmediateOrCancel
	default:
		fmt.Println("UNKNOWN TOF", o.TimeInForce)
	}

	ord.Price = &types.DoubleValue{Value: o.Price}
	return ord
}

func buildCreateOrderParams(symbol string, o *messages.NewOrder) *dydx.CreateOrderParams {
	var side dydx.OrderSide
	var typ dydx.OrderType
	var tif dydx.OrderTIF
	var postOnly = false
	var price float64
	if o.OrderSide == models.Buy {
		side = dydx.BUY
	} else {
		side = dydx.SELL
	}
	switch o.OrderType {
	case models.Limit:
		typ = dydx.LIMIT
	case models.Market:
		typ = dydx.MARKET
	case models.StopLimit:
		typ = dydx.STOP_LIMIT
	case models.TrailingStopLimit:
		typ = dydx.TRAILING_STOP
	}
	switch o.TimeInForce {
	case models.GoodTillCancel:
		tif = dydx.GOOD_TIL_TIME
	case models.FillOrKill:
		tif = dydx.FILL_OR_KILL
	case models.ImmediateOrCancel:
		tif = dydx.IMMEDIATE_OR_CANCEL
	}
	for _, exec := range o.ExecutionInstructions {
		if exec == models.ParticipateDoNotInitiate {
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
