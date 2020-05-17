package exchanges

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
	orders     map[string]*Order
	securities map[uint64]*Security
}

func NewAccount(ID string, securities []*models.Security) *Account {
	accnt := &Account{
		ID:         ID,
		orders:     make(map[string]*Order),
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
		accnt.orders[o.ClientOrderID] = &Order{
			Order:          o,
			previousStatus: o.OrderStatus,
		}
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
			return fmt.Errorf("order with nil instrument")
		}
		if p.Instrument.SecurityID == nil {
			return fmt.Errorf("order with nil security ID")
		}
		sec, ok := accnt.securities[p.Instrument.SecurityID.Value]
		if !ok {
			return fmt.Errorf("security %d for order not found", p.Instrument.SecurityID.Value)
		}
		sec.RawQuantity = int64(p.Quantity * sec.LotPrecision)
	}

	return nil
}

func (accnt *Account) NewOrder(order *models.Order) (*messages.ExecutionReport, error) {
	if _, ok := accnt.orders[order.ClientOrderID]; ok {
		return nil, fmt.Errorf("client order ID already exists")
	}
	if order.OrderStatus != models.PendingNew {
		return nil, fmt.Errorf("new order must have a PendingNew order status")
	}
	if order.Instrument == nil {
		return nil, fmt.Errorf("order has no instrument")
	}
	if order.Instrument.SecurityID == nil {
		return nil, fmt.Errorf("order has no security ID, ref from symbol not supported")
	}
	sec, ok := accnt.securities[order.Instrument.SecurityID.Value]
	if !ok {
		return nil, fmt.Errorf("unknown security %d", order.Instrument.SecurityID.Value)
	}
	rawCumQty := int64(order.CumQuantity * sec.LotPrecision)
	if rawCumQty > 0 {
		return nil, fmt.Errorf("new order must have zero CumQuantity")
	}
	accnt.orders[order.ClientOrderID] = &Order{
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
		Side:            order.Side,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: types.TimestampNow(),
	}, nil
}

func (accnt *Account) ConfirmNewOrder(clientID string) (*messages.ExecutionReport, error) {
	order, ok := accnt.orders[clientID]
	if !ok {
		return nil, fmt.Errorf("unknown order %s", clientID)
	}
	order.OrderStatus = models.New
	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &types.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.New,
		OrderStatus:     order.OrderStatus,
		Instrument:      order.Instrument,
		Side:            order.Side,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: types.TimestampNow(),
	}, nil
}

func (accnt *Account) RejectNewOrder(clientID string, reason messages.OrderRejectionReason) (*messages.ExecutionReport, error) {
	order, ok := accnt.orders[clientID]
	if !ok {
		return nil, fmt.Errorf("unknown order %s", clientID)
	}
	order.OrderStatus = models.Rejected
	delete(accnt.orders, clientID)

	return &messages.ExecutionReport{
		OrderID:              order.OrderID,
		ClientOrderID:        &types.StringValue{Value: order.ClientOrderID},
		ExecutionID:          "", // TODO
		ExecutionType:        messages.Rejected,
		OrderStatus:          models.Rejected,
		Instrument:           order.Instrument,
		Side:                 order.Side,
		LeavesQuantity:       order.LeavesQuantity,
		CumQuantity:          order.CumQuantity,
		TransactionTime:      types.TimestampNow(),
		OrderRejectionReason: reason,
	}, nil
}

func (accnt *Account) CancelOrder(clientID string) (*messages.ExecutionReport, error) {
	order, ok := accnt.orders[clientID]
	if !ok {
		return nil, fmt.Errorf("unknown order %s", clientID)
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
		Side:            order.Side,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: types.TimestampNow(),
	}, nil
}

func (accnt *Account) ConfirmCancelOrder(clientID string) (*messages.ExecutionReport, error) {
	order, ok := accnt.orders[clientID]
	if !ok {
		return nil, fmt.Errorf("unknown order %s", clientID)
	}
	order.OrderStatus = models.Canceled

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &types.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.Canceled,
		OrderStatus:     models.Canceled,
		Instrument:      order.Instrument,
		Side:            order.Side,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: types.TimestampNow(),
	}, nil
}

func (accnt *Account) RejectCancelOrder(clientID string) (*messages.ExecutionReport, error) {
	order, ok := accnt.orders[clientID]
	if !ok {
		return nil, fmt.Errorf("unknown order %s", clientID)
	}
	order.OrderStatus = order.previousStatus

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &types.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.Canceled,
		OrderStatus:     order.OrderStatus,
		Instrument:      order.Instrument,
		Side:            order.Side,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: types.TimestampNow(),
	}, nil
}

func (accnt *Account) ConfirmFill(clientID string, tradeID string, quantity float64) (*messages.ExecutionReport, error) {
	order, ok := accnt.orders[clientID]
	if !ok {
		return nil, fmt.Errorf("unknown order %s", clientID)
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
		ExecutionType:   messages.Fill,
		OrderStatus:     order.OrderStatus,
		Instrument:      order.Instrument,
		Side:            order.Side,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: types.TimestampNow(),
		TradeID:         &types.StringValue{Value: tradeID},
	}, nil
}
