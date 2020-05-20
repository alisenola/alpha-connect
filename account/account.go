package account

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"math"
)

type Security struct {
	*models.Security
	RawQuantity   int64
	TickPrecision float64
	LotPrecision  float64
}

type Order struct {
	*models.Order
	previousStatus models.OrderStatus
}

type Account struct {
	ID         string
	ordersID   map[string]*Order
	ordersClID map[string]*Order
	securities map[uint64]*Security
}

func NewAccount(ID string, securities []*models.Security) *Account {
	accnt := &Account{
		ID:         ID,
		ordersID:   make(map[string]*Order),
		ordersClID: make(map[string]*Order),
		securities: make(map[uint64]*Security),
	}
	for _, s := range securities {
		accnt.securities[s.SecurityID] = &Security{
			Security:      s,
			RawQuantity:   0,
			TickPrecision: math.Ceil(1. / s.MinPriceIncrement),
			LotPrecision:  math.Ceil(1. / s.RoundLot),
		}
	}

	return accnt
}

func (accnt *Account) Sync(orders []*models.Order, positions []*models.Position) error {

	for _, o := range orders {
		ord := &Order{
			Order:          o,
			previousStatus: o.OrderStatus,
		}
		accnt.ordersID[o.OrderID] = ord
		accnt.ordersClID[o.ClientOrderID] = ord
	}

	// Reset securities
	for _, s := range accnt.securities {
		s.RawQuantity = 0
	}
	for _, p := range positions {
		if p.AccountID != accnt.ID {
			return fmt.Errorf("got position for wrong account ID")
		}
		if p.Instrument == nil {
			return fmt.Errorf("position with nil instrument")
		}
		if p.Instrument.SecurityID == nil {
			return fmt.Errorf("position with nil security ID")
		}
		sec, ok := accnt.securities[p.Instrument.SecurityID.Value]
		if !ok {
			return fmt.Errorf("security %d for order not found", p.Instrument.SecurityID.Value)
		}
		sec.RawQuantity = int64(p.Quantity * sec.LotPrecision)
	}

	return nil
}

func (accnt *Account) NewOrder(order *models.Order) (*messages.ExecutionReport, *messages.RejectionReason) {
	if _, ok := accnt.ordersClID[order.ClientOrderID]; ok {
		res := messages.DuplicateOrder
		return nil, &res
	}
	if order.OrderStatus != models.PendingNew {
		res := messages.Other
		return nil, &res
	}
	if order.Instrument == nil || order.Instrument.SecurityID == nil {
		res := messages.UnknownSymbol
		return nil, &res
	}
	sec, ok := accnt.securities[order.Instrument.SecurityID.Value]
	if !ok {
		res := messages.UnknownSymbol
		return nil, &res
	}
	rawLeavesQuantity := sec.LotPrecision * order.LeavesQuantity
	if math.Abs(rawLeavesQuantity-math.Round(rawLeavesQuantity)) > 0.00001 {
		res := messages.IncorrectQuantity
		return nil, &res
	}
	rawCumQty := int64(order.CumQuantity * sec.LotPrecision)
	if rawCumQty > 0 {
		res := messages.IncorrectQuantity
		return nil, &res
	}
	accnt.ordersClID[order.ClientOrderID] = &Order{
		Order:          order,
		previousStatus: order.OrderStatus,
	}
	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &types.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.PendingNew,
		OrderStatus:     order.OrderStatus,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: types.TimestampNow(),
	}, nil
}

func (accnt *Account) ConfirmNewOrder(clientID string, ID string) (*messages.ExecutionReport, error) {
	order, ok := accnt.ordersClID[clientID]
	if !ok {
		return nil, fmt.Errorf("unknown order %s", clientID)
	}
	if order.OrderStatus != models.PendingNew {
		// Order already confirmed, nop
		return nil, nil
	}
	order.OrderID = ID
	order.OrderStatus = models.New
	accnt.ordersID[ID] = order
	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &types.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.New,
		OrderStatus:     order.OrderStatus,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: types.TimestampNow(),
	}, nil
}

func (accnt *Account) RejectNewOrder(clientID string, reason messages.RejectionReason) (*messages.ExecutionReport, error) {
	order, ok := accnt.ordersClID[clientID]
	if !ok {
		return nil, fmt.Errorf("unknown order %s", clientID)
	}
	order.OrderStatus = models.Rejected
	delete(accnt.ordersClID, clientID)

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &types.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.Rejected,
		OrderStatus:     models.Rejected,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: types.TimestampNow(),
		RejectionReason: reason,
	}, nil
}

func (accnt *Account) CancelOrder(ID string) (*messages.ExecutionReport, *messages.RejectionReason) {
	var order *Order
	order, _ = accnt.ordersClID[ID]
	if order == nil {
		order, _ = accnt.ordersID[ID]
	}
	if order == nil {
		res := messages.UnknownOrder
		return nil, &res
	}
	if order.OrderStatus == models.PendingCancel {
		res := messages.CancelAlreadyPending
		return nil, &res
	}

	// Save current order status in case cancel gets rejected
	order.previousStatus = order.OrderStatus
	order.OrderStatus = models.PendingCancel

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &types.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.PendingCancel,
		OrderStatus:     models.PendingCancel,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: types.TimestampNow(),
	}, nil
}

func (accnt *Account) GetOrders(filter *messages.OrderFilter) []*models.Order {
	var orders []*models.Order
	for _, o := range accnt.ordersClID {
		if filter != nil && filter.Instrument != nil && o.Instrument.SecurityID.Value != filter.Instrument.SecurityID.Value {
			continue
		}
		if filter != nil && filter.Side != nil && o.Side != filter.Side.Value {
			continue
		}
		if filter != nil && filter.OrderID != nil && o.OrderID != filter.OrderID.Value {
			continue
		}
		if filter != nil && filter.ClientOrderID != nil && o.ClientOrderID != filter.ClientOrderID.Value {
			continue
		}
		if filter != nil && filter.OrderStatus != nil && o.OrderStatus != filter.OrderStatus.Value {
			continue
		}
		orders = append(orders, o.Order)
	}

	return orders
}

func (accnt *Account) GetOpenOrders(instrument *models.Instrument, side *models.Side) []*models.Order {
	var orders []*models.Order
	for _, o := range accnt.ordersClID {
		if (instrument != nil && o.Instrument.SecurityID.Value != instrument.SecurityID.Value) ||
			(side != nil && o.Side != (*side)) {
			// pass
		} else {
			if o.OrderStatus == models.New || o.OrderStatus == models.PartiallyFilled {
				orders = append(orders, o.Order)
			}
		}
	}

	return orders
}

func (accnt *Account) ConfirmCancelOrder(ID string) (*messages.ExecutionReport, error) {
	var order *Order
	order, _ = accnt.ordersClID[ID]
	if order == nil {
		order, _ = accnt.ordersID[ID]
	}
	if order == nil {
		return nil, fmt.Errorf("unknown order %s", ID)
	}
	if order.OrderStatus == models.Canceled {
		return nil, nil
	}
	if order.OrderStatus != models.PendingCancel {
		return nil, fmt.Errorf("error not pending cancel")
	}

	order.OrderStatus = models.Canceled
	order.LeavesQuantity = 0.

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &types.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.Canceled,
		OrderStatus:     models.Canceled,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: types.TimestampNow(),
	}, nil
}

func (accnt *Account) RejectCancelOrder(ID string, reason messages.RejectionReason) (*messages.ExecutionReport, error) {
	var order *Order
	order, _ = accnt.ordersClID[ID]
	if order == nil {
		order, _ = accnt.ordersID[ID]
	}
	if order == nil {
		return nil, fmt.Errorf("unknown order %s", ID)
	}
	order.OrderStatus = order.previousStatus

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &types.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.Rejected,
		OrderStatus:     order.OrderStatus,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: types.TimestampNow(),
		RejectionReason: reason,
	}, nil
}

func (accnt *Account) ConfirmFill(ID string, tradeID string, price, quantity float64) (*messages.ExecutionReport, error) {
	var order *Order
	order, _ = accnt.ordersClID[ID]
	if order == nil {
		order, _ = accnt.ordersID[ID]
	}
	if order == nil {
		return nil, fmt.Errorf("unknown order %s", ID)
	}
	sec := accnt.securities[order.Instrument.SecurityID.Value]
	rawFillQuantity := int64(quantity * sec.LotPrecision)
	rawLeavesQuantity := int64(order.LeavesQuantity * sec.LotPrecision)
	if rawFillQuantity > rawLeavesQuantity {
		return nil, fmt.Errorf("fill bigger than order leaves quantity")
	}
	rawCumQuantity := int64(order.CumQuantity * sec.LotPrecision)
	order.LeavesQuantity = float64(rawLeavesQuantity-rawFillQuantity) / sec.LotPrecision
	order.CumQuantity = float64(rawCumQuantity+rawFillQuantity) / sec.LotPrecision
	if rawFillQuantity == rawLeavesQuantity {
		order.OrderStatus = models.Filled
	} else {
		order.OrderStatus = models.PartiallyFilled
	}

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &types.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.Trade,
		OrderStatus:     order.OrderStatus,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: types.TimestampNow(),
		TradeID:         &types.StringValue{Value: tradeID},
		FillPrice:       &types.DoubleValue{Value: price},
		FillQuantity:    &types.DoubleValue{Value: quantity},
	}, nil
}

func (accnt *Account) GetPositions() []*models.Position {
	var positions []*models.Position
	for _, s := range accnt.securities {
		if s.RawQuantity > 0 {
			positions = append(positions,
				&models.Position{
					AccountID: accnt.ID,
					Instrument: &models.Instrument{
						SecurityID: &types.UInt64Value{Value: s.SecurityID},
						Exchange:   s.Exchange,
						Symbol:     &types.StringValue{Value: s.Symbol},
					},
					Quantity: float64(s.RawQuantity) / s.LotPrecision,
				})
		}
	}

	return positions
}
