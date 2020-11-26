package account

import (
	"errors"
	"fmt"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/modeling"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"math"
	"sync"
	"time"
)

type Order struct {
	*models.Order
	lastEventTime     time.Time
	previousStatus    models.OrderStatus
	pendingAmendPrice *types.DoubleValue
	pendingAmendQty   *types.DoubleValue
}

var ErrNotPendingReplace = errors.New("not pending replace")

type Security interface {
	AddBidOrder(ID string, price, quantity, queue float64)
	AddAskOrder(ID string, price, quantity, queue float64)
	RemoveBidOrder(ID string)
	RemoveAskOrder(ID string)
	UpdateBidOrder(ID string, price, qty float64)
	UpdateAskOrder(ID string, price, qty float64)
	UpdateBidOrderPrice(ID string, price float64)
	UpdateAskOrderPrice(ID string, price float64)
	UpdateBidOrderQuantity(ID string, qty float64)
	UpdateAskOrderQuantity(ID string, qty float64)
	UpdateBidOrderQueue(ID string, queue float64)
	UpdateAskOrderQueue(ID string, queue float64)
	UpdatePositionSize(size float64)
	GetLotPrecision() float64
	GetInstrument() *models.Instrument
	Clear()
	AddSampleValueChange(model modeling.MarketModel, time uint64, values []float64)
	GetELROnCancelBid(ID string, model modeling.MarketModel, time uint64, values []float64, value float64) float64
	GetELROnCancelAsk(ID string, model modeling.MarketModel, time uint64, values []float64, value float64) float64
	GetELROnLimitBidChange(ID string, model modeling.MarketModel, time uint64, values []float64, value float64, prices []float64, queues []float64, maxQuote float64) (float64, *COrder)
	GetELROnLimitAskChange(ID string, model modeling.MarketModel, time uint64, values []float64, value float64, prices []float64, queues []float64, maxBase float64) (float64, *COrder)
	GetELROnMarketBuy(model modeling.MarketModel, time uint64, values []float64, value float64, price float64, quantity float64, maxQuantity float64) (float64, *COrder)
	GetELROnMarketSell(model modeling.MarketModel, time uint64, values []float64, value float64, price float64, quantity float64, maxQuantity float64) (float64, *COrder)
}

type Account struct {
	sync.RWMutex
	*models.Account
	ordersID        map[string]*Order
	ordersClID      map[string]*Order
	securities      map[uint64]Security
	positions       map[uint64]*Position
	balances        map[uint32]float64
	assets          map[uint32]*xchangerModels.Asset
	margin          int64
	MarginCurrency  *xchangerModels.Asset
	MarginPrecision float64
	quoteCurrency   *xchangerModels.Asset
}

func NewAccount(account *models.Account, marginCurrency *xchangerModels.Asset, marginPrecision float64) (*Account, error) {
	quoteCurrency := &constants.DOLLAR
	accnt := &Account{
		Account:         account,
		ordersID:        make(map[string]*Order),
		ordersClID:      make(map[string]*Order),
		securities:      make(map[uint64]Security),
		positions:       make(map[uint64]*Position),
		balances:        make(map[uint32]float64),
		assets:          make(map[uint32]*xchangerModels.Asset),
		margin:          0,
		MarginCurrency:  marginCurrency,
		MarginPrecision: marginPrecision,
		quoteCurrency:   quoteCurrency,
	}

	return accnt, nil
}

func (accnt *Account) Sync(securities []*models.Security, orders []*models.Order, positions []*models.Position, balances []*models.Balance, margin float64, makerFee, takerFee *float64) error {
	accnt.Lock()
	defer accnt.Unlock()
	var err error
	for _, s := range securities {
		switch s.SecurityType {
		case enum.SecurityType_CRYPTO_PERP, enum.SecurityType_CRYPTO_FUT:
			accnt.securities[s.SecurityID], err = NewMarginSecurity(s, accnt.MarginCurrency, accnt.quoteCurrency, makerFee, takerFee)
			if err != nil {
				return err
			}
			posMakerFee := s.MakerFee.Value
			if makerFee != nil {
				posMakerFee = *makerFee
			}
			posTakerFee := s.TakerFee.Value
			if takerFee != nil {
				posTakerFee = *takerFee
			}
			if s.MinPriceIncrement == nil || s.RoundLot == nil {
				return fmt.Errorf("security is missing MinPriceIncrement or RoundLot")
			}
			accnt.positions[s.SecurityID] = &Position{
				inverse:         s.IsInverse,
				tickPrecision:   math.Ceil(1. / s.MinPriceIncrement.Value),
				lotPrecision:    math.Ceil(1. / s.RoundLot.Value),
				marginPrecision: accnt.MarginPrecision,
				multiplier:      s.Multiplier.Value,
				makerFee:        posMakerFee,
				takerFee:        posTakerFee,
			}
		case enum.SecurityType_CRYPTO_SPOT:
			accnt.securities[s.SecurityID], err = NewSpotSecurity(s, accnt.quoteCurrency, makerFee, takerFee)
			if err != nil {
				return err
			}
			accnt.assets[s.Underlying.ID] = s.Underlying
			accnt.assets[s.QuoteCurrency.ID] = s.QuoteCurrency
		}
	}

	accnt.margin = int64(math.Round(margin * accnt.MarginPrecision))
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
		if p.AccountID != accnt.AccountID {
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
	accnt.Lock()
	defer accnt.Unlock()
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
	rawCumQty := int64(math.Round(order.CumQuantity * lotPrecision))
	if rawCumQty > 0 {
		res := messages.IncorrectQuantity
		return nil, &res
	}
	if order.OrderType == models.Limit && order.Price == nil {
		res := messages.InvalidOrder
		return nil, &res
	}
	accnt.ordersClID[order.ClientOrderID] = &Order{
		Order:          order,
		lastEventTime:  time.Now(),
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
	accnt.Lock()
	defer accnt.Unlock()
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
	order.lastEventTime = time.Now()
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
	accnt.Lock()
	defer accnt.Unlock()
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

func (accnt *Account) ReplaceOrder(ID string, price *types.DoubleValue, quantity *types.DoubleValue) (*messages.ExecutionReport, *messages.RejectionReason) {
	accnt.Lock()
	defer accnt.Unlock()
	var order *Order
	order, _ = accnt.ordersClID[ID]
	if order == nil {
		order, _ = accnt.ordersID[ID]
	}
	if order == nil {
		res := messages.UnknownOrder
		return nil, &res
	}
	if order.OrderStatus != models.New && order.OrderStatus != models.PartiallyFilled {
		res := messages.NonReplaceableOrder
		return nil, &res
	}

	order.pendingAmendPrice = price
	order.pendingAmendQty = quantity
	order.previousStatus = order.OrderStatus
	order.OrderStatus = models.PendingReplace
	order.lastEventTime = time.Now()

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &types.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.PendingReplace,
		OrderStatus:     models.PendingReplace,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: types.TimestampNow(),
	}, nil
}

func (accnt *Account) ConfirmReplaceOrder(ID string) (*messages.ExecutionReport, error) {
	accnt.Lock()
	defer accnt.Unlock()
	var order *Order
	order, _ = accnt.ordersClID[ID]
	if order == nil {
		order, _ = accnt.ordersID[ID]
	}
	if order == nil {
		return nil, fmt.Errorf("unknown order %s", ID)
	}
	if order.OrderStatus != models.PendingReplace {
		return nil, ErrNotPendingReplace
	}
	order.OrderStatus = order.previousStatus
	order.lastEventTime = time.Now()

	if order.pendingAmendQty != nil {
		order.LeavesQuantity = order.pendingAmendQty.Value
		order.pendingAmendQty = nil
	}
	if order.pendingAmendPrice != nil {
		if order.Price == nil {
			return nil, fmt.Errorf("replacing price on an order without price")
		}
		order.Price.Value = order.pendingAmendPrice.Value
		order.pendingAmendPrice = nil
	}

	if order.OrderType == models.Limit {
		if order.Side == models.Buy {
			accnt.securities[order.Instrument.SecurityID.Value].UpdateBidOrder(order.ClientOrderID, order.Price.Value, order.LeavesQuantity)
		} else {
			accnt.securities[order.Instrument.SecurityID.Value].UpdateAskOrder(order.ClientOrderID, order.Price.Value, order.LeavesQuantity)
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

func (accnt *Account) RejectReplaceOrder(ID string, reason messages.RejectionReason) (*messages.ExecutionReport, error) {
	accnt.Lock()
	defer accnt.Unlock()
	var order *Order
	order, _ = accnt.ordersClID[ID]
	if order == nil {
		order, _ = accnt.ordersID[ID]
	}
	if order == nil {
		return nil, fmt.Errorf("unknown order %s", ID)
	}

	if order.OrderStatus == models.PendingReplace {
		order.OrderStatus = order.previousStatus
		order.lastEventTime = time.Now()
	}

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

func (accnt *Account) CancelOrder(ID string) (*messages.ExecutionReport, *messages.RejectionReason) {
	accnt.Lock()
	defer accnt.Unlock()
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
	if order.OrderStatus != models.New && order.OrderStatus != models.PartiallyFilled && order.OrderStatus != models.PendingNew {
		res := messages.NonCancelableOrder
		return nil, &res
	}

	// Save current order status in case cancel gets rejected
	order.previousStatus = order.OrderStatus
	order.OrderStatus = models.PendingCancel
	order.lastEventTime = time.Now()

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

func (accnt *Account) ConfirmCancelOrder(ID string) (*messages.ExecutionReport, error) {
	accnt.Lock()
	defer accnt.Unlock()
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

	order.OrderStatus = models.Canceled
	order.lastEventTime = time.Now()
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
	accnt.Lock()
	defer accnt.Unlock()
	var order *Order
	order, _ = accnt.ordersClID[ID]
	if order == nil {
		order, _ = accnt.ordersID[ID]
	}
	if order == nil {
		return nil, fmt.Errorf("unknown order %s", ID)
	}
	if order.OrderStatus == models.PendingCancel {
		order.OrderStatus = order.previousStatus
		order.lastEventTime = time.Now()
	}

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
	accnt.Lock()
	defer accnt.Unlock()
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

	//fmt.Println("FILL", price, quantity, lotPrecision)
	rawFillQuantity := int64(math.Round(quantity * lotPrecision))
	rawLeavesQuantity := int64(math.Round(order.LeavesQuantity * lotPrecision))
	if rawFillQuantity > rawLeavesQuantity {
		return nil, fmt.Errorf("fill bigger than order leaves quantity %d %d", rawFillQuantity, rawLeavesQuantity)
	}
	rawCumQuantity := int64(math.Round(order.CumQuantity * lotPrecision))
	leavesQuantity := float64(rawLeavesQuantity-rawFillQuantity) / lotPrecision
	order.LeavesQuantity = leavesQuantity
	order.CumQuantity = float64(rawCumQuantity+rawFillQuantity) / lotPrecision
	if rawFillQuantity == rawLeavesQuantity {
		order.OrderStatus = models.Filled
	} else {
		order.OrderStatus = models.PartiallyFilled
	}
	order.lastEventTime = time.Now()

	if order.OrderType == models.Limit {
		if order.Side == models.Buy {
			if order.OrderStatus == models.Filled {
				accnt.securities[order.Instrument.SecurityID.Value].RemoveBidOrder(order.ClientOrderID)
			} else {
				accnt.securities[order.Instrument.SecurityID.Value].UpdateBidOrderQuantity(order.ClientOrderID, order.LeavesQuantity)
			}
		} else {
			if order.OrderStatus == models.Filled {
				accnt.securities[order.Instrument.SecurityID.Value].RemoveAskOrder(order.ClientOrderID)
			} else {
				accnt.securities[order.Instrument.SecurityID.Value].UpdateAskOrderQuantity(order.ClientOrderID, order.LeavesQuantity)
			}
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
		//fmt.Println("MARGIN POS", float64(pos.cost) / accnt.MarginPrecision, float64(pos.rawSize) / accnt.MarginPrecision)
		var fee, cost int64
		if order.Side == models.Buy {
			fee, cost = pos.Buy(price, quantity, taker)
		} else {
			fee, cost = pos.Sell(price, quantity, taker)
		}
		//fmt.Println("MARGIN POS", float64(pos.cost) / accnt.MarginPrecision, float64(pos.rawSize) / accnt.MarginPrecision)

		accnt.margin -= fee
		accnt.margin -= cost
		er.FeeAmount = &types.DoubleValue{Value: float64(fee) / accnt.MarginPrecision}
		er.FeeCurrency = accnt.MarginCurrency
		er.FeeBasis = messages.Percentage
		er.FeeType = messages.ExchangeFees
		// TODO mutex on position ?
		marginSec.UpdatePositionSize(float64(pos.rawSize) / pos.lotPrecision)
	}

	return er, nil
}

func (accnt *Account) Settle() {
	accnt.balances[accnt.MarginCurrency.ID] += float64(accnt.margin) / accnt.MarginPrecision
	// TODO check less than zero
	accnt.margin = 0
}

func (accnt *Account) GetOrders(filter *messages.OrderFilter) []*models.Order {
	accnt.RLock()
	defer accnt.RUnlock()
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

func (accnt *Account) GetPositions() []*models.Position {
	accnt.RLock()
	defer accnt.RUnlock()
	var positions []*models.Position
	for k, pos := range accnt.positions {
		pos := pos.GetPosition()
		if pos != nil {
			pos.Instrument = accnt.securities[k].GetInstrument()
			pos.AccountID = accnt.AccountID
			positions = append(positions, pos)
		}
	}

	return positions
}

func (accnt *Account) GetBalances() []*models.Balance {
	accnt.RLock()
	defer accnt.RUnlock()
	var balances []*models.Balance
	for k, b := range accnt.balances {
		balances = append(balances, &models.Balance{
			AccountID: accnt.AccountID,
			Asset:     accnt.assets[k],
			Quantity:  b,
		})
	}

	return balances
}

func (accnt *Account) GetBalance(asset uint32) float64 {
	accnt.RLock()
	defer accnt.RUnlock()
	if b, ok := accnt.balances[asset]; ok {
		return b
	} else {
		return 0.
	}
}

func (accnt *Account) GetPositionSize(securityID uint64) float64 {
	accnt.RLock()
	defer accnt.RUnlock()
	if pos, ok := accnt.positions[securityID]; ok {
		// TODO mutex ?
		return float64(pos.rawSize) / pos.lotPrecision
	} else {
		return 0.
	}
}

func (accnt *Account) GetMargin() float64 {
	accnt.RLock()
	defer accnt.RUnlock()
	return float64(accnt.margin) / accnt.MarginPrecision
}

func (accnt *Account) CleanOrders() {
	accnt.RLock()
	defer accnt.RUnlock()
	for k, o := range accnt.ordersClID {
		if (o.OrderStatus == models.Filled || o.OrderStatus == models.Canceled) && (time.Now().Sub(o.lastEventTime) > time.Minute) {
			delete(accnt.ordersClID, k)
		}
	}
	for k, o := range accnt.ordersID {
		if (o.OrderStatus == models.Filled || o.OrderStatus == models.Canceled) && (time.Now().Sub(o.lastEventTime) > time.Minute) {
			delete(accnt.ordersClID, k)
		}
	}
}

func (accnt *Account) Compare(other *Account) bool {
	// Compare positions, balances and margin
	accnt.Lock()
	other.Lock()
	defer accnt.Unlock()
	defer other.Unlock()

	if accnt.margin != other.margin {
		fmt.Println("DIFF MARGIN")
		fmt.Println(accnt.margin, other.margin)
		return false
	}
	for k, p1 := range accnt.positions {
		if p2, ok := other.positions[k]; ok {
			if p1.rawSize != p2.rawSize {
				fmt.Println("DIFF POS")
				fmt.Println(p1)
				fmt.Println(p2)
				return false
			}
			if p1.cost != p2.cost {
				fmt.Println("DIFF POS")
				fmt.Println(p1)
				fmt.Println(p2)
				return false
			}
		} else {
			return false
		}
	}
	for k, b1 := range accnt.balances {
		if b2, ok := other.balances[k]; ok {
			// TODO
			if math.Abs(b1-b2) > 0.0000001 {
				fmt.Println("DIFF BALANCE")
				fmt.Println(b1, b2)
				return false
			}
		} else {
			return false
		}
	}

	return true
}
