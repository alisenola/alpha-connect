package fbinance

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/fbinance"
	"strconv"
)

func wsOrderToModel(o *fbinance.WSExecution) *models.Order {
	ord := &models.Order{
		OrderID:       fmt.Sprintf("%d", o.OrderID),
		ClientOrderID: o.ClientOrderID,
		Instrument: &models.Instrument{
			Exchange: &constants.FBINANCE,
			Symbol:   &types.StringValue{Value: o.Symbol},
		},
		LeavesQuantity: o.OrigQuantity - o.CumQuantity,
		CumQuantity:    o.CumQuantity,
	}
	if o.ReduceOnly {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ReduceOnly)
	}

	switch o.OrderStatus {
	case fbinance.NEW_ORDER:
		ord.OrderStatus = models.New
	case fbinance.CANCELED_ORDER:
		ord.OrderStatus = models.Canceled
	case fbinance.PARTIALLY_FILLED:
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
	case fbinance.LIMIT_ORDER:
		ord.OrderType = models.Limit
	case fbinance.MARKET_ORDER:
		ord.OrderType = models.Market
	case fbinance.STOP_LOSS_ORDER:
		ord.OrderType = models.Stop
	case fbinance.STOP_LOSS_LIMIT_ORDER:
		ord.OrderType = models.StopLimit
	default:
		fmt.Println("UNKNOWN ORDER TYPE", o.OrderType)
	}

	switch o.Side {
	case fbinance.BUY_ORDER:
		ord.Side = models.Buy
	case fbinance.SELL_ODER:
		ord.Side = models.Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", o.Side)
	}

	switch o.TimeInForce {
	case fbinance.GOOD_TILL_CANCEL:
		ord.TimeInForce = models.GoodTillCancel
	case fbinance.FILL_OR_KILL:
		ord.TimeInForce = models.FillOrKill
	case fbinance.IMMEDIATE_OR_CANCEL:
		ord.TimeInForce = models.ImmediateOrCancel
	case fbinance.GOOD_TILL_CROSSING:
		ord.TimeInForce = models.GoodTillCancel
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ParticipateDoNotInitiate)
	default:
		fmt.Println("UNKNOWN TOF", o.TimeInForce)
	}

	ord.Price = &types.DoubleValue{Value: o.OrigPrice}
	return ord
}

func orderToModel(o *fbinance.OrderData) *models.Order {
	ord := &models.Order{
		OrderID:       fmt.Sprintf("%d", o.OrderID),
		ClientOrderID: o.ClientOrderID,
		Instrument: &models.Instrument{
			Exchange: &constants.FBINANCE,
			Symbol:   &types.StringValue{Value: o.Symbol},
		},
		LeavesQuantity: o.OriginalQuantity - o.ExecutedQuantity, // TODO check
		CumQuantity:    o.CumQuantity,
	}
	if o.ReduceOnly {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ReduceOnly)
	}

	switch o.Status {
	case fbinance.NEW_ORDER:
		ord.OrderStatus = models.New
	case fbinance.CANCELED_ORDER:
		ord.OrderStatus = models.Canceled
	case fbinance.PARTIALLY_FILLED:
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
	case fbinance.LIMIT_ORDER:
		ord.OrderType = models.Limit
	case fbinance.MARKET_ORDER:
		ord.OrderType = models.Market
	case fbinance.STOP_LOSS_ORDER:
		ord.OrderType = models.Stop
	case fbinance.STOP_LOSS_LIMIT_ORDER:
		ord.OrderType = models.StopLimit
	default:
		fmt.Println("UNKNOWN ORDER TYPE", o.Type)
	}

	switch o.Side {
	case fbinance.BUY_ORDER:
		ord.Side = models.Buy
	case fbinance.SELL_ODER:
		ord.Side = models.Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", o.Side)
	}

	switch o.TimeInForce {
	case fbinance.GOOD_TILL_CANCEL:
		ord.TimeInForce = models.GoodTillCancel
	case fbinance.FILL_OR_KILL:
		ord.TimeInForce = models.FillOrKill
	case fbinance.IMMEDIATE_OR_CANCEL:
		ord.TimeInForce = models.ImmediateOrCancel
	case fbinance.GOOD_TILL_CROSSING:
		ord.TimeInForce = models.GoodTillCancel
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ParticipateDoNotInitiate)
	default:
		fmt.Println("UNKNOWN TOF", o.TimeInForce)
	}

	ord.Price = &types.DoubleValue{Value: o.Price}
	return ord
}

func buildPostOrderRequest(symbol string, order *messages.NewOrder, tickPrecision, lotPrecision int) (fbinance.NewOrderRequest, *messages.RejectionReason) {
	var side fbinance.OrderSide
	var typ fbinance.OrderType
	if order.OrderSide == models.Buy {
		side = fbinance.BUY_ORDER
	} else {
		side = fbinance.SELL_ODER
	}
	switch order.OrderType {
	case models.Limit:
		typ = fbinance.LIMIT_ORDER
	case models.Market:
		typ = fbinance.MARKET_ORDER
	case models.Stop:
		typ = fbinance.STOP_LOSS_ORDER
	case models.StopLimit:
		typ = fbinance.STOP_LOSS_LIMIT_ORDER
	default:
		rej := messages.UnsupportedOrderType
		return nil, &rej
	}

	request := fbinance.NewNewOrderRequest(symbol, side, typ)

	fmt.Println("SET QTY", order.Quantity, lotPrecision, strconv.FormatFloat(order.Quantity, 'f', int(lotPrecision), 64))
	request.SetQuantity(order.Quantity, lotPrecision)
	request.SetNewClientOrderID(order.ClientOrderID)

	if order.OrderType != models.Market {
		switch order.TimeInForce {
		case models.Session:
			request.SetTimeInForce(fbinance.GOOD_TILL_CANCEL)
		case models.GoodTillCancel:
			request.SetTimeInForce(fbinance.GOOD_TILL_CANCEL)
		case models.ImmediateOrCancel:
			request.SetTimeInForce(fbinance.IMMEDIATE_OR_CANCEL)
		case models.FillOrKill:
			request.SetTimeInForce(fbinance.FILL_OR_KILL)
		default:
			rej := messages.UnsupportedOrderTimeInForce
			return nil, &rej
		}
	}

	if order.Price != nil {
		request.SetPrice(order.Price.Value, tickPrecision)
	}

	for _, exec := range order.ExecutionInstructions {
		rej := messages.UnsupportedOrderCharacteristic
		switch exec {
		case models.ReduceOnly:
			if err := request.SetReduceOnly(true); err != nil {
				return nil, &rej
			}
		case models.ParticipateDoNotInitiate:
			request.SetTimeInForce(fbinance.GOOD_TILL_CROSSING)
		default:
			return nil, &rej
		}
	}
	request.SetNewOrderResponseType(fbinance.ACK_RESPONSE_TYPE)

	return request, nil
}
