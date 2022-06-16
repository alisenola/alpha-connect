package fbinance

import (
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/fbinance"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func StatusToModel(status fbinance.OrderStatus) *models.OrderStatus {
	switch status {
	case fbinance.NEW_ORDER:
		v := models.OrderStatus_New
		return &v
	case fbinance.CANCELED_ORDER:
		v := models.OrderStatus_Canceled
		return &v
	case fbinance.PARTIALLY_FILLED:
		v := models.OrderStatus_PartiallyFilled
		return &v
	case fbinance.FILLED:
		v := models.OrderStatus_Filled
		return &v
	case fbinance.EXPIRED:
		v := models.OrderStatus_Expired
		return &v
	default:
		return nil
	}
}

func WSOrderToModel(o *fbinance.WSExecution) *models.Order {
	ord := &models.Order{
		OrderID:       fmt.Sprintf("%d", o.OrderID),
		ClientOrderID: o.ClientOrderID,
		Instrument: &models.Instrument{
			Exchange: constants.FBINANCE,
			Symbol:   &wrapperspb.StringValue{Value: o.Symbol},
		},
		LeavesQuantity: o.OrigQuantity - o.CumQuantity,
		CumQuantity:    o.CumQuantity,
	}
	if o.ReduceOnly {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ExecutionInstruction_ReduceOnly)
	}

	switch o.OrderStatus {
	case fbinance.NEW_ORDER:
		ord.OrderStatus = models.OrderStatus_New
	case fbinance.CANCELED_ORDER:
		ord.OrderStatus = models.OrderStatus_Canceled
	case fbinance.PARTIALLY_FILLED:
		ord.OrderStatus = models.OrderStatus_PartiallyFilled
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
		ord.OrderType = models.OrderType_Limit
	case fbinance.MARKET_ORDER:
		ord.OrderType = models.OrderType_Market
	case fbinance.STOP_LOSS_ORDER:
		ord.OrderType = models.OrderType_Stop
	case fbinance.STOP_LOSS_LIMIT_ORDER:
		ord.OrderType = models.OrderType_StopLimit
	default:
		fmt.Println("UNKNOWN ORDER TYPE", o.OrderType)
	}

	switch o.Side {
	case fbinance.BUY_ORDER:
		ord.Side = models.Side_Buy
	case fbinance.SELL_ODER:
		ord.Side = models.Side_Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", o.Side)
	}

	switch o.TimeInForce {
	case fbinance.GOOD_TILL_CANCEL:
		ord.TimeInForce = models.TimeInForce_GoodTillCancel
	case fbinance.FILL_OR_KILL:
		ord.TimeInForce = models.TimeInForce_FillOrKill
	case fbinance.IMMEDIATE_OR_CANCEL:
		ord.TimeInForce = models.TimeInForce_ImmediateOrCancel
	case fbinance.GOOD_TILL_CROSSING:
		ord.TimeInForce = models.TimeInForce_GoodTillCancel
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ExecutionInstruction_ParticipateDoNotInitiate)
	default:
		fmt.Println("UNKNOWN TOF", o.TimeInForce)
	}

	ord.Price = &wrapperspb.DoubleValue{Value: o.OrigPrice}
	return ord
}

func orderToModel(o *fbinance.OrderData) *models.Order {
	ord := &models.Order{
		OrderID:       fmt.Sprintf("%d", o.OrderID),
		ClientOrderID: o.ClientOrderID,
		Instrument: &models.Instrument{
			Exchange: constants.FBINANCE,
			Symbol:   &wrapperspb.StringValue{Value: o.Symbol},
		},
		LeavesQuantity: o.OriginalQuantity - o.ExecutedQuantity, // TODO check
		CumQuantity:    o.CumQuantity,
	}
	if o.ReduceOnly {
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ExecutionInstruction_ReduceOnly)
	}

	switch o.Status {
	case fbinance.NEW_ORDER:
		ord.OrderStatus = models.OrderStatus_New
	case fbinance.CANCELED_ORDER:
		ord.OrderStatus = models.OrderStatus_Canceled
	case fbinance.PARTIALLY_FILLED:
		ord.OrderStatus = models.OrderStatus_PartiallyFilled
	case fbinance.FILLED:
		ord.OrderStatus = models.OrderStatus_Filled
	case fbinance.EXPIRED:
		ord.OrderStatus = models.OrderStatus_Expired
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
		ord.OrderType = models.OrderType_Limit
	case fbinance.MARKET_ORDER:
		ord.OrderType = models.OrderType_Market
	case fbinance.STOP_LOSS_ORDER:
		ord.OrderType = models.OrderType_Stop
	case fbinance.STOP_LOSS_LIMIT_ORDER:
		ord.OrderType = models.OrderType_StopLimit
	default:
		fmt.Println("UNKNOWN ORDER TYPE", o.Type)
	}

	switch o.Side {
	case fbinance.BUY_ORDER:
		ord.Side = models.Side_Buy
	case fbinance.SELL_ODER:
		ord.Side = models.Side_Sell
	default:
		fmt.Println("UNKNOWN ORDER SIDE", o.Side)
	}

	switch o.TimeInForce {
	case fbinance.GOOD_TILL_CANCEL:
		ord.TimeInForce = models.TimeInForce_GoodTillCancel
	case fbinance.FILL_OR_KILL:
		ord.TimeInForce = models.TimeInForce_FillOrKill
	case fbinance.IMMEDIATE_OR_CANCEL:
		ord.TimeInForce = models.TimeInForce_ImmediateOrCancel
	case fbinance.GOOD_TILL_CROSSING:
		ord.TimeInForce = models.TimeInForce_GoodTillCancel
		ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ExecutionInstruction_ParticipateDoNotInitiate)
	default:
		fmt.Println("UNKNOWN TOF", o.TimeInForce)
	}

	ord.Price = &wrapperspb.DoubleValue{Value: o.Price}
	return ord
}

func buildPostOrderRequest(symbol string, order *messages.NewOrder, tickPrecision, lotPrecision int) (fbinance.NewOrderRequest, *messages.RejectionReason) {
	var side fbinance.OrderSide
	var typ fbinance.OrderType
	if order.OrderSide == models.Side_Buy {
		side = fbinance.BUY_ORDER
	} else {
		side = fbinance.SELL_ODER
	}
	switch order.OrderType {
	case models.OrderType_Limit:
		typ = fbinance.LIMIT_ORDER
	case models.OrderType_Market:
		typ = fbinance.MARKET_ORDER
	case models.OrderType_Stop:
		typ = fbinance.STOP_LOSS_ORDER
	case models.OrderType_StopLimit:
		typ = fbinance.STOP_LOSS_LIMIT_ORDER
	default:
		rej := messages.RejectionReason_UnsupportedOrderType
		return nil, &rej
	}

	request := fbinance.NewNewOrderRequest(symbol, side, typ)

	request.SetQuantity(order.Quantity, lotPrecision)
	request.SetNewClientOrderID(order.ClientOrderID)

	if order.OrderType != models.OrderType_Market {
		switch order.TimeInForce {
		case models.TimeInForce_Session:
			request.SetTimeInForce(fbinance.GOOD_TILL_CANCEL)
		case models.TimeInForce_GoodTillCancel:
			request.SetTimeInForce(fbinance.GOOD_TILL_CANCEL)
		case models.TimeInForce_ImmediateOrCancel:
			request.SetTimeInForce(fbinance.IMMEDIATE_OR_CANCEL)
		case models.TimeInForce_FillOrKill:
			request.SetTimeInForce(fbinance.FILL_OR_KILL)
		default:
			rej := messages.RejectionReason_UnsupportedOrderTimeInForce
			return nil, &rej
		}
	}

	if order.Price != nil {
		request.SetPrice(order.Price.Value, tickPrecision)
	}

	for _, exec := range order.ExecutionInstructions {
		rej := messages.RejectionReason_UnsupportedOrderCharacteristic
		switch exec {
		case models.ExecutionInstruction_ReduceOnly:
			if err := request.SetReduceOnly(true); err != nil {
				return nil, &rej
			}
		case models.ExecutionInstruction_ParticipateDoNotInitiate:
			request.SetTimeInForce(fbinance.GOOD_TILL_CROSSING)
		default:
			return nil, &rej
		}
	}
	request.SetNewOrderResponseType(fbinance.RESULT_RESPONSE_TYPE)

	return request, nil
}
