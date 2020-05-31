package account

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alphac/enum"
	"gitlab.com/alphaticks/alphac/modeling"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"math"
)

type Order struct {
	*models.Order
	previousStatus models.OrderStatus
}

type Security interface {
	AddBidOrder(ID string, price, quantity, queue float64)
	AddAskOrder(ID string, price, quantity, queue float64)
	RemoveBidOrder(ID string)
	RemoveAskOrder(ID string)
	UpdateBidOrderQuantity(ID string, qty float64)
	UpdateAskOrderQuantity(ID string, qty float64)
	UpdateBidOrderQueue(ID string, queue float64)
	UpdateAskOrderQueue(ID string, queue float64)
	UpdatePositionSize(size float64)
	GetLotPrecision() float64
	GetInstrument() *models.Instrument
	Clear()
	AddSampleValueChange(model modeling.Model, time uint64, values []float64)
	GetELROnCancelBid(ID string, model modeling.Model, time uint64, values []float64, value float64) float64
	GetELROnCancelAsk(ID string, model modeling.Model, time uint64, values []float64, value float64) float64
	GetELROnLimitBidChange(ID string, model modeling.Model, time uint64, values []float64, value float64, prices []float64, queues []float64, maxQuote float64) (float64, *COrder)
	GetELROnLimitAskChange(ID string, model modeling.Model, time uint64, values []float64, value float64, prices []float64, queues []float64, maxBase float64) (float64, *COrder)
}

type Account struct {
	ID              string
	ordersID        map[string]*Order
	ordersClID      map[string]*Order
	securities      map[uint64]Security
	positions       map[uint64]*Position
	balances        map[uint32]float64
	assets          map[uint32]*xchangerModels.Asset
	margin          int64
	marginCurrency  *xchangerModels.Asset
	marginPrecision float64
}

func NewAccount(ID string, marginCurrency *xchangerModels.Asset, marginPrecision float64) *Account {
	accnt := &Account{
		ID:              ID,
		ordersID:        make(map[string]*Order),
		ordersClID:      make(map[string]*Order),
		securities:      make(map[uint64]Security),
		positions:       make(map[uint64]*Position),
		balances:        make(map[uint32]float64),
		assets:          make(map[uint32]*xchangerModels.Asset),
		margin:          0,
		marginCurrency:  marginCurrency,
		marginPrecision: marginPrecision,
	}

	return accnt
}

func (accnt *Account) Sync(securities []*models.Security, orders []*models.Order, positions []*models.Position, balances []*models.Balance, margin float64) error {

	for _, s := range securities {
		switch s.SecurityType {
		case enum.SecurityType_CRYPTO_PERP, enum.SecurityType_CRYPTO_FUT:
			accnt.securities[s.SecurityID] = NewMarginSecurity(s, accnt.marginCurrency)
			accnt.positions[s.SecurityID] = &Position{
				inverse:         s.IsInverse,
				tickPrecision:   math.Ceil(1. / s.MinPriceIncrement),
				lotPrecision:    math.Ceil(1. / s.RoundLot),
				marginPrecision: accnt.marginPrecision,
				multiplier:      s.Multiplier.Value,
				makerFee:        s.MakerFee.Value,
				takerFee:        s.TakerFee.Value,
			}
		case enum.SecurityType_CRYPTO_SPOT:
			accnt.securities[s.SecurityID] = NewSpotSecurity(s)
			accnt.assets[s.Underlying.ID] = s.Underlying
			accnt.assets[s.QuoteCurrency.ID] = s.QuoteCurrency
		}
	}

	accnt.margin = int64(math.Round(margin * accnt.marginPrecision))
	for _, s := range accnt.securities {
		s.Clear()
	}

	for _, o := range orders {
		ord := &Order{
			Order:          o,
			previousStatus: o.OrderStatus,
		}
		accnt.ordersID[o.OrderID] = ord
		accnt.ordersClID[o.ClientOrderID] = ord

		// Add orders to security
		if (o.OrderStatus == models.New || o.OrderStatus == models.PartiallyFilled) && o.OrderType == models.Limit {
			if o.Side == models.Buy {
				accnt.securities[o.Instrument.SecurityID.Value].AddBidOrder(o.ClientOrderID, o.Price.Value, o.LeavesQuantity, 0)
			} else {
				accnt.securities[o.Instrument.SecurityID.Value].AddAskOrder(o.ClientOrderID, o.Price.Value, o.LeavesQuantity, 0)
			}
		}
	}

	// Reset positions
	for _, pos := range accnt.positions {
		pos.Sync(0, 0)
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
		pos, ok := accnt.positions[p.Instrument.SecurityID.Value]
		if !ok {
			return fmt.Errorf("position %d for order not found", p.Instrument.SecurityID.Value)
		}
		pos.Sync(p.Cost, p.Quantity)

		accnt.securities[p.Instrument.SecurityID.Value].UpdatePositionSize(p.Quantity)
		// TODO cross
		//sec.Position.Cross = p.Cross
	}

	for _, b := range balances {
		accnt.balances[b.Asset.ID] = b.Quantity
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
	lotPrecision := sec.GetLotPrecision()
	rawLeavesQuantity := lotPrecision * order.LeavesQuantity
	if math.Abs(rawLeavesQuantity-math.Round(rawLeavesQuantity)) > 0.00001 {
		res := messages.IncorrectQuantity
		return nil, &res
	}
	rawCumQty := int64(order.CumQuantity * lotPrecision)
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

	if order.OrderType == models.Limit {
		if order.Side == models.Buy {
			accnt.securities[order.Instrument.SecurityID.Value].AddBidOrder(order.ClientOrderID, order.Price.Value, order.LeavesQuantity, 0)
		} else {
			accnt.securities[order.Instrument.SecurityID.Value].AddAskOrder(order.ClientOrderID, order.Price.Value, order.LeavesQuantity, 0)
		}
	}

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
	if order.OrderType == models.Limit {
		if order.Side == models.Buy {
			accnt.securities[order.Instrument.SecurityID.Value].RemoveBidOrder(order.ClientOrderID)
		} else {
			accnt.securities[order.Instrument.SecurityID.Value].RemoveAskOrder(order.ClientOrderID)
		}
	}

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

func (accnt *Account) ConfirmFill(ID string, tradeID string, price, quantity float64, taker bool) (*messages.ExecutionReport, error) {
	var order *Order
	order, _ = accnt.ordersClID[ID]
	if order == nil {
		order, _ = accnt.ordersID[ID]
	}
	if order == nil {
		return nil, fmt.Errorf("unknown order %s", ID)
	}
	sec := accnt.securities[order.Instrument.SecurityID.Value]
	lotPrecision := sec.GetLotPrecision()

	rawFillQuantity := int64(quantity * lotPrecision)
	rawLeavesQuantity := int64(order.LeavesQuantity * lotPrecision)
	if rawFillQuantity > rawLeavesQuantity {
		return nil, fmt.Errorf("fill bigger than order leaves quantity")
	}
	rawCumQuantity := int64(order.CumQuantity * lotPrecision)
	leavesQuantity := float64(rawLeavesQuantity-rawFillQuantity) / lotPrecision
	order.LeavesQuantity = leavesQuantity
	order.CumQuantity = float64(rawCumQuantity+rawFillQuantity) / lotPrecision
	if rawFillQuantity == rawLeavesQuantity {
		order.OrderStatus = models.Filled
	} else {
		order.OrderStatus = models.PartiallyFilled
	}

	if order.OrderType == models.Limit {
		if order.Side == models.Buy {
			accnt.securities[order.Instrument.SecurityID.Value].UpdateBidOrderQuantity(order.ClientOrderID, order.LeavesQuantity)
		} else {
			accnt.securities[order.Instrument.SecurityID.Value].UpdateAskOrderQuantity(order.ClientOrderID, order.LeavesQuantity)
		}
	}

	er := &messages.ExecutionReport{
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
	}

	switch sec.(type) {
	case *SpotSecurity:
		spotSec := sec.(*SpotSecurity)
		if order.Side == models.Buy {
			accnt.balances[spotSec.Underlying.ID] += quantity
			accnt.balances[spotSec.QuoteCurrency.ID] -= quantity * price
		} else {
			accnt.balances[spotSec.Underlying.ID] -= quantity
			accnt.balances[spotSec.QuoteCurrency.ID] += quantity * price
		}
	case *MarginSecurity:
		marginSec := sec.(*MarginSecurity)
		pos := accnt.positions[marginSec.SecurityID]
		// Margin
		var fee, cost int64
		if order.Side == models.Buy {
			fee, cost = pos.Buy(price, quantity, taker)
		} else {
			fee, cost = pos.Sell(price, quantity, taker)
		}
		accnt.margin -= fee
		accnt.margin -= cost
		er.FeeAmount = &types.DoubleValue{Value: float64(fee) / accnt.marginPrecision}
		er.FeeCurrency = accnt.marginCurrency
		er.FeeBasis = messages.Percentage
		er.FeeType = messages.ExchangeFees
		// TODO mutex on position ?
		marginSec.UpdatePositionSize(float64(pos.rawSize) / pos.lotPrecision)
	}

	return er, nil
}

func (accnt *Account) Settle() {
	accnt.balances[accnt.marginCurrency.ID] += float64(accnt.margin) / accnt.marginPrecision
	// TODO check less than zero
	accnt.margin = 0
}

func (accnt *Account) GetPositions() []*models.Position {
	var positions []*models.Position
	for k, pos := range accnt.positions {
		pos := pos.GetPosition()
		if pos != nil {
			pos.Instrument = accnt.securities[k].GetInstrument()
			pos.AccountID = accnt.ID
			positions = append(positions, pos)
		}
	}

	return positions
}

func (accnt *Account) GetBalances() []*models.Balance {
	var balances []*models.Balance
	for k, b := range accnt.balances {
		balances = append(balances, &models.Balance{
			AccountID: accnt.ID,
			Asset:     accnt.assets[k],
			Quantity:  b,
		})
	}

	return balances
}

func (accnt *Account) GetPositionSize(securityID uint64) float64 {
	if pos, ok := accnt.positions[securityID]; ok {
		// TODO mutex ?
		return float64(pos.rawSize) / pos.lotPrecision
	} else {
		return 0.
	}
}

func (accnt *Account) GetMargin() float64 {
	return float64(accnt.margin) / accnt.marginPrecision
}
