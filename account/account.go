package account

import (
	"errors"
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/modeling"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
	"sync"
	"time"
)

func IsPending(status models.OrderStatus) bool {
	return status == models.OrderStatus_PendingFilled ||
		status == models.OrderStatus_PendingCancel ||
		status == models.OrderStatus_PendingNew ||
		status == models.OrderStatus_PendingReplace
}

func IsClosed(status models.OrderStatus) bool {
	return status == models.OrderStatus_Filled ||
		status == models.OrderStatus_Canceled ||
		status == models.OrderStatus_Expired
}

func IsOpen(status models.OrderStatus) bool {
	return status == models.OrderStatus_New ||
		status == models.OrderStatus_Created ||
		status == models.OrderStatus_PartiallyFilled
}

type Order struct {
	*models.Order
	lastEventTime          time.Time
	previousStatus         models.OrderStatus
	pendingAmendPrice      *wrapperspb.DoubleValue
	pendingAmendQty        *wrapperspb.DoubleValue
	unknownOrderErrorCount int
}

var ErrNotPendingReplace = errors.New("not pending replace")

type Security interface {
	GetSecurity() *models.Security
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
	symbolToSec     map[string]Security
	positions       map[uint64]*Position
	balances        map[uint32]int64
	assets          map[uint32]*xchangerModels.Asset
	margin          int64
	MarginCurrency  *xchangerModels.Asset
	MarginPrecision float64
	quoteCurrency   *xchangerModels.Asset
	takerFee        *float64
	makerFee        *float64
	expirationLimit time.Duration
}

func NewAccount(account *models.Account) (*Account, error) {
	quoteCurrency := constants.DOLLAR
	accnt := &Account{
		Account:         account,
		ordersID:        make(map[string]*Order),
		ordersClID:      make(map[string]*Order),
		securities:      make(map[uint64]Security),
		positions:       make(map[uint64]*Position),
		balances:        make(map[uint32]int64),
		assets:          make(map[uint32]*xchangerModels.Asset),
		margin:          0,
		quoteCurrency:   quoteCurrency,
		expirationLimit: 1 * time.Minute,
	}
	switch account.Exchange.ID {
	case constants.FBINANCE.ID:
		accnt.MarginCurrency = constants.TETHER
		if constants.TETHER == nil {
			return nil, fmt.Errorf("not loaded")
		}
		accnt.MarginPrecision = 100000000
	case constants.BITMEX.ID:
		accnt.MarginCurrency = constants.BITCOIN
		if constants.BITCOIN == nil {
			return nil, fmt.Errorf("not loaded")
		}
		accnt.MarginPrecision = 1. / 0.00000001
	case constants.FTX.ID:
		accnt.MarginCurrency = constants.DOLLAR
		if constants.DOLLAR == nil {
			return nil, fmt.Errorf("not loaded")
		}
		accnt.MarginPrecision = 100000000
	case constants.DYDX.ID:
		accnt.MarginCurrency = constants.USDC
		if constants.USDC == nil {
			return nil, fmt.Errorf("not loaded")
		}
		accnt.MarginPrecision = 1e6
	case constants.BYBITL.ID:
		if constants.TETHER == nil {
			return nil, fmt.Errorf("not loaded")
		}
		accnt.MarginCurrency = constants.TETHER
		accnt.MarginPrecision = 1e8
	}
	if accnt.MarginCurrency != nil {
		accnt.assets[accnt.MarginCurrency.ID] = accnt.MarginCurrency
	}

	return accnt, nil
}

func (accnt *Account) Sync(securities []*models.Security, orders []*models.Order, positions []*models.Position, balances []*models.Balance, makerFee, takerFee *float64) error {
	accnt.Lock()
	defer accnt.Unlock()
	accnt.ordersID = make(map[string]*Order)
	accnt.ordersClID = make(map[string]*Order)
	accnt.securities = make(map[uint64]Security)
	accnt.symbolToSec = make(map[string]Security)
	accnt.positions = make(map[uint64]*Position)
	accnt.balances = make(map[uint32]int64)
	accnt.assets = make(map[uint32]*xchangerModels.Asset)
	accnt.takerFee = takerFee
	accnt.makerFee = makerFee

	var err error
	for _, s := range securities {
		switch s.SecurityType {
		case enum.SecurityType_CRYPTO_PERP, enum.SecurityType_CRYPTO_FUT:
			accnt.securities[s.SecurityID], err = NewMarginSecurity(s, accnt.MarginCurrency, accnt.quoteCurrency, makerFee, takerFee)
			if err != nil {
				return err
			}
			var posMakerFee float64
			if makerFee != nil {
				posMakerFee = *makerFee
			} else if s.MakerFee != nil {
				posMakerFee = s.MakerFee.Value
			} else {
				return fmt.Errorf("unknown maker fee")
			}
			var posTakerFee float64
			if takerFee != nil {
				posTakerFee = *takerFee
			} else if s.TakerFee != nil {
				posTakerFee = s.TakerFee.Value
			} else {
				return fmt.Errorf("unknown taker fee")
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

	for _, s := range accnt.securities {
		s.Clear()
		accnt.symbolToSec[s.GetSecurity().Symbol] = s
	}

	for _, o := range orders {
		ord := &Order{
			Order:                  o,
			previousStatus:         o.OrderStatus,
			unknownOrderErrorCount: 0,
		}
		accnt.ordersID[o.OrderID] = ord
		accnt.ordersClID[o.ClientOrderID] = ord
		// Add orders to security
		if (o.OrderStatus == models.OrderStatus_New || o.OrderStatus == models.OrderStatus_PartiallyFilled) && o.OrderType == models.OrderType_Limit {
			if o.Side == models.Side_Buy {
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
		if p.Account != accnt.Name {
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

	accnt.margin = 0
	for _, b := range balances {
		accnt.assets[b.Asset.ID] = b.Asset
		accnt.balances[b.Asset.ID] = int64(math.Round(b.Quantity * accnt.MarginPrecision))
	}
	if accnt.MarginCurrency != nil {
		accnt.assets[accnt.MarginCurrency.ID] = accnt.MarginCurrency
	}

	return nil
}

func (accnt *Account) getSec(instrument *models.Instrument) (Security, *messages.RejectionReason) {
	var sec Security
	if instrument.SecurityID != nil {
		var ok bool
		sec, ok = accnt.securities[instrument.SecurityID.Value]
		if !ok {
			res := messages.RejectionReason_UnknownSecurityID
			return nil, &res
		}
	} else if instrument.Symbol != nil {
		var ok bool
		sec, ok = accnt.symbolToSec[instrument.Symbol.Value]
		if !ok {
			res := messages.RejectionReason_UnknownSymbol
			return nil, &res
		}
	} else {
		res := messages.RejectionReason_InvalidRequest
		return nil, &res
	}

	return sec, nil
}

func (accnt *Account) NewOrder(order *models.Order) (*messages.ExecutionReport, *messages.RejectionReason) {
	accnt.Lock()
	defer accnt.Unlock()
	if _, ok := accnt.ordersClID[order.ClientOrderID]; ok {
		res := messages.RejectionReason_DuplicateOrder
		return nil, &res
	}
	if order.OrderStatus != models.OrderStatus_PendingNew {
		res := messages.RejectionReason_Other
		return nil, &res
	}
	if order.Instrument == nil {
		res := messages.RejectionReason_UnknownSymbol
		return nil, &res
	}

	sec, rej := accnt.getSec(order.Instrument)
	if rej != nil {
		return nil, rej
	}
	if order.Instrument.SecurityID == nil {
		order.Instrument.SecurityID = &wrapperspb.UInt64Value{Value: sec.GetSecurity().SecurityID}
	}
	if order.Instrument.Symbol == nil {
		order.Instrument.Symbol = &wrapperspb.StringValue{Value: sec.GetSecurity().Symbol}
	}
	if order.CreationTime == nil {
		order.CreationTime = timestamppb.New(time.Now())
	}
	order.LastEventTime = order.CreationTime
	lotPrecision := sec.GetLotPrecision()
	rawLeavesQuantity := lotPrecision * order.LeavesQuantity
	if math.Abs(rawLeavesQuantity-math.Round(rawLeavesQuantity)) > 0.00001 {
		res := messages.RejectionReason_IncorrectQuantity
		return nil, &res
	}
	rawCumQty := int64(math.Round(order.CumQuantity * lotPrecision))
	if rawCumQty > 0 {
		res := messages.RejectionReason_IncorrectQuantity
		return nil, &res
	}
	if order.OrderType == models.OrderType_Limit && order.Price == nil {
		res := messages.RejectionReason_InvalidOrder
		return nil, &res
	}
	accnt.ordersClID[order.ClientOrderID] = &Order{
		Order:                  order,
		lastEventTime:          time.Now(),
		previousStatus:         order.OrderStatus,
		unknownOrderErrorCount: 0,
	}
	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &wrapperspb.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.ExecutionType_PendingNew,
		OrderStatus:     order.OrderStatus,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: timestamppb.Now(),
	}, nil
}

func (accnt *Account) ConfirmNewOrder(clientID string, ID string) (*messages.ExecutionReport, error) {
	accnt.Lock()
	defer accnt.Unlock()
	order, ok := accnt.ordersClID[clientID]
	if !ok {
		return nil, fmt.Errorf("unknown order %s", clientID)
	}
	if order.OrderStatus != models.OrderStatus_PendingNew {
		// Order already confirmed, nop
		return nil, nil
	}
	order.OrderID = ID
	order.OrderStatus = models.OrderStatus_New
	order.lastEventTime = time.Now()
	order.LastEventTime = timestamppb.New(order.lastEventTime)
	accnt.ordersID[ID] = order

	sec, rej := accnt.getSec(order.Instrument)
	if rej != nil {
		return nil, fmt.Errorf("unknown instrument: %s", rej.String())
	}

	if order.OrderType == models.OrderType_Limit {
		if order.Side == models.Side_Buy {
			accnt.securities[sec.GetSecurity().SecurityID].AddBidOrder(order.ClientOrderID, order.Price.Value, order.LeavesQuantity, 0)
		} else {
			accnt.securities[sec.GetSecurity().SecurityID].AddAskOrder(order.ClientOrderID, order.Price.Value, order.LeavesQuantity, 0)
		}
	}

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &wrapperspb.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.ExecutionType_New,
		OrderStatus:     order.OrderStatus,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: timestamppb.Now(),
	}, nil
}

func (accnt *Account) RejectNewOrder(clientID string, reason messages.RejectionReason) (*messages.ExecutionReport, error) {
	accnt.Lock()
	defer accnt.Unlock()
	order, ok := accnt.ordersClID[clientID]
	if !ok {
		return nil, fmt.Errorf("unknown order %s", clientID)
	}
	order.OrderStatus = models.OrderStatus_Rejected
	delete(accnt.ordersClID, clientID)

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &wrapperspb.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.ExecutionType_Rejected,
		OrderStatus:     models.OrderStatus_Rejected,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: timestamppb.Now(),
		RejectionReason: reason,
	}, nil
}

func (accnt *Account) ReplaceOrder(ID string, price *wrapperspb.DoubleValue, quantity *wrapperspb.DoubleValue) (*messages.ExecutionReport, *messages.RejectionReason) {
	accnt.Lock()
	defer accnt.Unlock()
	var order *Order
	order = accnt.ordersClID[ID]
	if order == nil {
		order = accnt.ordersID[ID]
	}
	if order == nil {
		res := messages.RejectionReason_UnknownOrder
		return nil, &res
	}
	if order.OrderStatus != models.OrderStatus_New && order.OrderStatus != models.OrderStatus_PartiallyFilled {
		res := messages.RejectionReason_NonReplaceableOrder
		return nil, &res
	}

	order.pendingAmendPrice = price
	order.pendingAmendQty = quantity
	order.previousStatus = order.OrderStatus
	order.OrderStatus = models.OrderStatus_PendingReplace
	order.lastEventTime = time.Now()
	order.LastEventTime = timestamppb.New(order.lastEventTime)

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &wrapperspb.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.ExecutionType_PendingReplace,
		OrderStatus:     models.OrderStatus_PendingReplace,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: timestamppb.Now(),
	}, nil
}

func (accnt *Account) ConfirmReplaceOrder(ID, newID string) (*messages.ExecutionReport, error) {
	accnt.Lock()
	defer accnt.Unlock()
	var order *Order
	order = accnt.ordersClID[ID]
	if order == nil {
		order = accnt.ordersID[ID]
	}
	if order == nil {
		return nil, fmt.Errorf("unknown order %s", ID)
	}
	if order.OrderStatus != models.OrderStatus_PendingReplace {
		return nil, ErrNotPendingReplace
	}
	order.OrderStatus = order.previousStatus
	order.lastEventTime = time.Now()
	order.LastEventTime = timestamppb.New(order.lastEventTime)

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
	if newID != "" {
		order.OrderID = newID
	}

	if order.OrderType == models.OrderType_Limit {
		if order.Side == models.Side_Buy {
			accnt.securities[order.Instrument.SecurityID.Value].UpdateBidOrder(order.ClientOrderID, order.Price.Value, order.LeavesQuantity)
		} else {
			accnt.securities[order.Instrument.SecurityID.Value].UpdateAskOrder(order.ClientOrderID, order.Price.Value, order.LeavesQuantity)
		}
	}

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &wrapperspb.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.ExecutionType_Canceled,
		OrderStatus:     models.OrderStatus_Canceled,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: timestamppb.Now(),
	}, nil
}

func (accnt *Account) RejectReplaceOrder(ID string, reason messages.RejectionReason) (*messages.ExecutionReport, error) {
	accnt.Lock()
	defer accnt.Unlock()
	var order *Order
	order = accnt.ordersClID[ID]
	if order == nil {
		order = accnt.ordersID[ID]
	}
	if order == nil {
		return nil, fmt.Errorf("unknown order %s", ID)
	}

	if order.OrderStatus == models.OrderStatus_PendingReplace {
		order.OrderStatus = order.previousStatus
		order.lastEventTime = time.Now()
		order.LastEventTime = timestamppb.New(order.lastEventTime)
	}

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &wrapperspb.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.ExecutionType_Rejected,
		OrderStatus:     order.OrderStatus,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: timestamppb.Now(),
		RejectionReason: reason,
	}, nil
}

func (accnt *Account) CancelOrder(ID string) (*messages.ExecutionReport, *messages.RejectionReason) {
	accnt.Lock()
	defer accnt.Unlock()
	var order *Order
	order = accnt.ordersClID[ID]
	if order == nil {
		order = accnt.ordersID[ID]
	}
	if order == nil {
		res := messages.RejectionReason_UnknownOrder
		return nil, &res
	}
	if order.OrderStatus == models.OrderStatus_PendingCancel {
		res := messages.RejectionReason_CancelAlreadyPending
		return nil, &res
	}
	if order.OrderStatus != models.OrderStatus_New && order.OrderStatus != models.OrderStatus_PartiallyFilled && order.OrderStatus != models.OrderStatus_PendingNew {
		res := messages.RejectionReason_NonCancelableOrder
		return nil, &res
	}

	// Save current order status in case cancel gets rejected
	order.previousStatus = order.OrderStatus
	order.OrderStatus = models.OrderStatus_PendingCancel
	order.lastEventTime = time.Now()
	order.LastEventTime = timestamppb.New(order.lastEventTime)

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &wrapperspb.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.ExecutionType_PendingCancel,
		OrderStatus:     models.OrderStatus_PendingCancel,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: timestamppb.Now(),
	}, nil
}

func (accnt *Account) ConfirmCancelOrder(ID string) (*messages.ExecutionReport, error) {
	accnt.Lock()
	defer accnt.Unlock()
	var order *Order
	order = accnt.ordersClID[ID]
	if order == nil {
		order = accnt.ordersID[ID]
	}
	if order == nil {
		return nil, fmt.Errorf("unknown order %s", ID)
	}
	if order.OrderStatus == models.OrderStatus_Canceled {
		return nil, nil
	}

	order.OrderStatus = models.OrderStatus_Canceled
	order.lastEventTime = time.Now()
	order.LastEventTime = timestamppb.New(order.lastEventTime)

	//order.LeavesQuantity = 0.
	if order.OrderType == models.OrderType_Limit {
		if order.Side == models.Side_Buy {
			accnt.securities[order.Instrument.SecurityID.Value].RemoveBidOrder(order.ClientOrderID)
		} else {
			accnt.securities[order.Instrument.SecurityID.Value].RemoveAskOrder(order.ClientOrderID)
		}
	}

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &wrapperspb.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.ExecutionType_Canceled,
		OrderStatus:     models.OrderStatus_Canceled,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: timestamppb.Now(),
	}, nil
}

func (accnt *Account) ConfirmExpiredOrder(ID string) (*messages.ExecutionReport, error) {
	accnt.Lock()
	defer accnt.Unlock()
	var order *Order
	order = accnt.ordersClID[ID]
	if order == nil {
		order = accnt.ordersID[ID]
	}
	if order == nil {
		return nil, fmt.Errorf("unknown order %s", ID)
	}
	if order.OrderStatus == models.OrderStatus_Expired {
		return nil, nil
	}

	order.OrderStatus = models.OrderStatus_Expired
	order.lastEventTime = time.Now()
	order.LastEventTime = timestamppb.New(order.lastEventTime)

	order.LeavesQuantity = 0.
	if order.OrderType == models.OrderType_Limit {
		if order.Side == models.Side_Buy {
			accnt.securities[order.Instrument.SecurityID.Value].RemoveBidOrder(order.ClientOrderID)
		} else {
			accnt.securities[order.Instrument.SecurityID.Value].RemoveAskOrder(order.ClientOrderID)
		}
	}

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &wrapperspb.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.ExecutionType_Expired,
		OrderStatus:     models.OrderStatus_Expired,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: timestamppb.Now(),
	}, nil
}

func (accnt *Account) RejectCancelOrder(ID string, reason messages.RejectionReason) (*messages.ExecutionReport, error) {
	accnt.Lock()
	defer accnt.Unlock()
	var order *Order
	order = accnt.ordersClID[ID]
	if order == nil {
		order = accnt.ordersID[ID]
	}
	if order == nil {
		return nil, fmt.Errorf("unknown order %s", ID)
	}
	if order.OrderStatus == models.OrderStatus_PendingCancel {
		if reason == messages.RejectionReason_UnknownOrder {
			order.unknownOrderErrorCount += 1
			if order.unknownOrderErrorCount > 3 {
				return nil, fmt.Errorf("unknown order %s, missed a fill", ID)
			}
		}
		order.OrderStatus = order.previousStatus
		order.lastEventTime = time.Now()
		order.LastEventTime = timestamppb.New(order.lastEventTime)
	}

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &wrapperspb.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.ExecutionType_Rejected,
		OrderStatus:     order.OrderStatus,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: timestamppb.Now(),
		RejectionReason: reason,
	}, nil
}

func (accnt *Account) PendingFilled(ID string) (*messages.ExecutionReport, error) {
	accnt.Lock()
	defer accnt.Unlock()
	var order *Order
	order = accnt.ordersClID[ID]
	if order == nil {
		order = accnt.ordersID[ID]
	}
	if order == nil {
		return nil, fmt.Errorf("unknown order %s", ID)
	}
	if order.OrderStatus == models.OrderStatus_Filled {
		return nil, nil
	}

	// Save current order status in case filled gets rejected
	order.previousStatus = order.OrderStatus
	order.OrderStatus = models.OrderStatus_PendingFilled
	order.lastEventTime = time.Now()
	order.LastEventTime = timestamppb.New(order.lastEventTime)

	return &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &wrapperspb.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.ExecutionType_PendingFilled,
		OrderStatus:     models.OrderStatus_PendingFilled,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: timestamppb.Now(),
	}, nil
}

func (accnt *Account) ConfirmFill(ID string, tradeID string, price, quantity float64, taker bool) (*messages.ExecutionReport, error) {
	accnt.Lock()
	defer accnt.Unlock()
	var order *Order
	order = accnt.ordersClID[ID]
	if order == nil {
		order = accnt.ordersID[ID]
	}
	if order == nil {
		return nil, fmt.Errorf("unknown order %s", ID)
	}
	sec := accnt.securities[order.Instrument.SecurityID.Value]
	lotPrecision := sec.GetLotPrecision()

	fmt.Println("FILL", price, quantity, order.LeavesQuantity, lotPrecision)
	rawFillQuantity := int64(math.Round(quantity * lotPrecision))
	rawLeavesQuantity := int64(math.Round(order.LeavesQuantity * lotPrecision))
	if rawFillQuantity > rawLeavesQuantity {
		return nil, fmt.Errorf("fill bigger than order leaves quantity %d %d", rawFillQuantity, rawLeavesQuantity)
	}
	rawCumQuantity := int64(math.Round(order.CumQuantity * lotPrecision))
	leavesQuantity := float64(rawLeavesQuantity-rawFillQuantity) / lotPrecision
	order.LeavesQuantity = leavesQuantity
	order.CumQuantity = float64(rawCumQuantity+rawFillQuantity) / lotPrecision

	if IsPending(order.OrderStatus) {
		// Pending, only set it to filled
		if rawFillQuantity == rawLeavesQuantity {
			order.OrderStatus = models.OrderStatus_Filled
		}
	} else if IsOpen(order.OrderStatus) {
		if rawFillQuantity == rawLeavesQuantity {
			order.OrderStatus = models.OrderStatus_Filled
		} else if order.OrderStatus == models.OrderStatus_New {
			// Only set it to partially filled if its in new status
			// otherwise, will overwrite pending state
			order.OrderStatus = models.OrderStatus_PartiallyFilled
		}
	}

	order.lastEventTime = time.Now()
	order.LastEventTime = timestamppb.New(order.lastEventTime)

	if order.OrderType == models.OrderType_Limit {
		if order.Side == models.Side_Buy {
			if order.OrderStatus == models.OrderStatus_Filled {
				accnt.securities[order.Instrument.SecurityID.Value].RemoveBidOrder(order.ClientOrderID)
			} else {
				accnt.securities[order.Instrument.SecurityID.Value].UpdateBidOrderQuantity(order.ClientOrderID, order.LeavesQuantity)
			}
		} else {
			if order.OrderStatus == models.OrderStatus_Filled {
				accnt.securities[order.Instrument.SecurityID.Value].RemoveAskOrder(order.ClientOrderID)
			} else {
				accnt.securities[order.Instrument.SecurityID.Value].UpdateAskOrderQuantity(order.ClientOrderID, order.LeavesQuantity)
			}
		}
	}

	er := &messages.ExecutionReport{
		OrderID:         order.OrderID,
		ClientOrderID:   &wrapperspb.StringValue{Value: order.ClientOrderID},
		ExecutionID:     "", // TODO
		ExecutionType:   messages.ExecutionType_Trade,
		OrderStatus:     order.OrderStatus,
		Instrument:      order.Instrument,
		LeavesQuantity:  order.LeavesQuantity,
		CumQuantity:     order.CumQuantity,
		TransactionTime: timestamppb.Now(),
		TradeID:         &wrapperspb.StringValue{Value: tradeID},
		FillPrice:       &wrapperspb.DoubleValue{Value: price},
		FillQuantity:    &wrapperspb.DoubleValue{Value: quantity},
	}

	switch sec := sec.(type) {
	case *SpotSecurity:
		if order.Side == models.Side_Buy {
			accnt.balances[sec.Underlying.ID] += int64(math.Round(quantity * accnt.MarginPrecision))
			accnt.balances[sec.QuoteCurrency.ID] -= int64(math.Round(quantity * price * accnt.MarginPrecision))
		} else {
			accnt.balances[sec.Underlying.ID] -= int64(math.Round(quantity))
			accnt.balances[sec.QuoteCurrency.ID] += int64(math.Round(quantity * price * accnt.MarginPrecision))
		}
	case *MarginSecurity:
		pos := accnt.positions[sec.SecurityID]
		// Margin
		//fmt.Println("MARGIN POS", float64(pos.cost) / accnt.MarginPrecision, float64(pos.rawSize) / accnt.MarginPrecision)
		var fee, cost int64
		if order.Side == models.Side_Buy {
			fee, cost = pos.Buy(price, quantity, taker)
		} else {
			fee, cost = pos.Sell(price, quantity, taker)
		}
		//fmt.Println("MARGIN POS", float64(pos.cost) / accnt.MarginPrecision, float64(pos.rawSize) / accnt.MarginPrecision)

		accnt.margin -= fee
		accnt.margin -= cost
		er.FeeAmount = &wrapperspb.DoubleValue{Value: float64(fee) / accnt.MarginPrecision}
		er.FeeCurrency = accnt.MarginCurrency
		er.FeeBasis = messages.FeeBasis_Percentage
		er.FeeType = messages.FeeType_ExchangeFees
		// TODO mutex on position ?
		sec.UpdatePositionSize(float64(pos.rawSize) / pos.lotPrecision)
	}

	return er, nil
}

func (accnt *Account) ConfirmFunding(security uint64, markPrice, fundingFee float64) (*messages.AccountUpdate, error) {
	accnt.RLock()
	defer accnt.RUnlock()
	if pos, ok := accnt.positions[security]; ok {
		fee := pos.Funding(markPrice, fundingFee)
		accnt.margin -= fee
		return nil, nil
	} else {
		return nil, fmt.Errorf("no open position for the funding")
	}
}

func (accnt *Account) UpdateBalance(asset *xchangerModels.Asset, balance float64, reason messages.AccountMovementType) (*messages.AccountUpdate, error) {
	accnt.Lock()
	defer accnt.Unlock()
	if _, ok := accnt.assets[asset.ID]; !ok {
		accnt.assets[asset.ID] = asset
	}
	// Reset margin, as we have a fresh balance for it
	if accnt.MarginCurrency != nil && asset.ID == accnt.MarginCurrency.ID {
		accnt.margin = 0.
	}
	accnt.balances[asset.ID] = int64(math.Round(balance * accnt.MarginPrecision))
	return &messages.AccountUpdate{
		Type:    reason,
		Asset:   accnt.assets[asset.ID],
		Balance: balance,
	}, nil
}

func (accnt *Account) Settle() {
	// Only for margin accounts
	if accnt.MarginCurrency != nil {
		accnt.Lock()
		defer accnt.Unlock()
		accnt.balances[accnt.MarginCurrency.ID] += accnt.margin
		accnt.margin = 0
	}
}

func (accnt *Account) HasOrder(ID string) bool {
	accnt.RLock()
	defer accnt.RUnlock()
	var order *Order
	order = accnt.ordersClID[ID]
	if order == nil {
		order = accnt.ordersID[ID]
	}
	return order != nil
}

func (accnt *Account) GetOrder(ID string) *models.Order {
	accnt.RLock()
	defer accnt.RUnlock()
	var order *Order
	order = accnt.ordersClID[ID]
	if order == nil {
		order = accnt.ordersID[ID]
	}
	if order == nil {
		return nil
	}
	return order.Order
}

func (accnt *Account) GetOrders(filter *messages.OrderFilter) []*models.Order {
	accnt.RLock()
	defer accnt.RUnlock()
	var orders []*models.Order
	for _, o := range accnt.ordersClID {
		if filter != nil {
			if filter.Instrument != nil {
				if filter.Instrument.SecurityID != nil && o.Instrument.SecurityID.Value != filter.Instrument.SecurityID.Value {
					continue
				}
				if filter.Instrument.Symbol != nil && o.Instrument.Symbol.Value != filter.Instrument.Symbol.Value {
					continue
				}
			}
			if filter.Side != nil && o.Side != filter.Side.Value {
				continue
			}
			if filter.OrderID != nil && o.OrderID != filter.OrderID.Value {
				continue
			}
			if filter.ClientOrderID != nil && o.ClientOrderID != filter.ClientOrderID.Value {
				continue
			}
			if filter.OrderStatus != nil && o.OrderStatus != filter.OrderStatus.Value {
				continue
			}
			if filter.Open != nil && IsOpen(o.OrderStatus) != filter.Open.Value {
				continue
			}
		}

		orders = append(orders, o.Order)
	}

	return orders
}

func (accnt *Account) GetMakerFee() *float64 {
	return accnt.makerFee
}

func (accnt *Account) GetTakerFee() *float64 {
	return accnt.takerFee
}

func (accnt *Account) GetSecurities() []*models.Security {
	var secs []*models.Security
	for _, sec := range accnt.securities {
		secs = append(secs, sec.GetSecurity())
	}
	return secs
}

func (accnt *Account) GetPositions() []*models.Position {
	accnt.RLock()
	defer accnt.RUnlock()
	var positions []*models.Position
	for k, pos := range accnt.positions {
		pos := pos.GetPosition()
		if pos != nil {
			pos.Instrument = accnt.securities[k].GetInstrument()
			pos.Account = accnt.Name
			positions = append(positions, pos)
		}
	}

	return positions
}

func (accnt *Account) GetAllPositions() []*models.Position {
	accnt.RLock()
	defer accnt.RUnlock()
	var positions []*models.Position
	for k, pos := range accnt.positions {
		p := &models.Position{
			Account:    accnt.Name,
			Instrument: accnt.securities[k].GetInstrument(),
			Cross:      false,
		}
		if pp := pos.GetPosition(); pp != nil {
			p.Quantity = pp.Quantity
			p.Cost = pp.Cost
		}
		positions = append(positions, p)
	}

	return positions
}

func (accnt *Account) GetPositionMap() map[uint64]*models.Position {
	accnt.RLock()
	defer accnt.RUnlock()
	var positions = make(map[uint64]*models.Position)
	for k, pos := range accnt.positions {
		pos := pos.GetPosition()
		if pos != nil {
			pos.Instrument = accnt.securities[k].GetInstrument()
			pos.Account = accnt.Name
			positions[k] = pos
		}
	}

	return positions
}

func (accnt *Account) GetBalances() []*models.Balance {
	accnt.Settle()
	accnt.RLock()
	defer accnt.RUnlock()
	var balances []*models.Balance
	for k, b := range accnt.balances {
		balances = append(balances, &models.Balance{
			Account:  accnt.Name,
			Asset:    accnt.assets[k],
			Quantity: float64(b) / accnt.MarginPrecision,
		})
	}

	return balances
}

func (accnt *Account) GetBalance(asset uint32) float64 {
	accnt.RLock()
	defer accnt.RUnlock()
	if b, ok := accnt.balances[asset]; ok {
		return float64(b) / accnt.MarginPrecision
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

func (accnt *Account) GetMargin(model modeling.Market) float64 {
	accnt.RLock()
	defer accnt.RUnlock()
	if accnt.MarginCurrency == nil {
		return 0.
	}
	var availableMargin int64 = 0
	availableMargin += accnt.margin
	availableMargin += accnt.balances[accnt.MarginCurrency.ID]

	if model != nil {
		for k, b := range accnt.balances {
			if k == accnt.MarginCurrency.ID {
				continue
			} else {
				pp, ok := model.GetPairPrice(k, accnt.MarginCurrency.ID)
				if ok {
					availableMargin += int64(math.Round(float64(b) * pp))
				}
			}
		}
	}
	return float64(availableMargin) / accnt.MarginPrecision
}

func (accnt *Account) CheckExpiration() error {
	accnt.RLock()
	defer accnt.RUnlock()
	for _, o := range accnt.ordersClID {
		if IsPending(o.OrderStatus) && (time.Since(o.lastEventTime) > accnt.expirationLimit) {
			return fmt.Errorf("order %s in unknown state %s", o.OrderID, o.OrderStatus.String())
		}
	}
	return nil
}

func (accnt *Account) CleanOrders() {
	accnt.Lock()
	defer accnt.Unlock()
	for k, o := range accnt.ordersClID {
		if IsClosed(o.OrderStatus) && (time.Since(o.lastEventTime) > time.Minute) {
			delete(accnt.ordersClID, k)
			delete(accnt.ordersID, o.OrderID)
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
			if b1-b2 != 0 {
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
