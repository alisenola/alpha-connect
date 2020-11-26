package legacy

import (
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/models"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"math"
	"sync"
)

type Position struct {
	Cost    int64
	RawSize int64
	Cross   bool
}

type AccountPortfolio struct {
	exchange   *xchangerModels.Exchange
	sampleSize int
	balances   map[uint32]float64
	positions  map[uint64]*Position
	securities map[uint64]Security
	margin     float64
}

type Security interface {
	updateSampleValueChange(time uint64)
	AddSampleValueChange(time uint64, values []float64)
	GetELROnCancelBid(time uint64, values []float64, value float64) float64
	GetELROnCancelAsk(time uint64, values []float64, value float64) float64
	GetELROnBidChange(time uint64, values []float64, value float64, prices []float64, queues []float64, maxQuote float64) (float64, *Order)
	GetELROnAskChange(time uint64, values []float64, value float64, prices []float64, queues []float64, maxBase float64) (float64, *Order)
}

type SpotSecurity struct {
	sync.RWMutex
	*models.Security
	BuyTradeModels    map[uint64]BuyTradeModel
	SellTradeModels   map[uint64]SellTradeModel
	BasePriceModel    PriceModel
	QuotePriceModel   PriceModel
	SampleValueChange []float64
	SampleMatchBid    []float64
	SampleMatchAsk    []float64
	SampleTime        uint64
	BidOrder          *Order
	AskOrder          *Order
}

type MarginSecurity struct {
	sync.RWMutex
	*models.Security
	MarginPriceModel    PriceModel
	SecurityPriceModel  PriceModel
	BuyTradeModel       BuyTradeModel
	SellTradeModel      SellTradeModel
	SampleValueChange   []float64
	SampleMatchBid      []float64
	SampleMatchAsk      []float64
	SampleMarginPrice   []float64
	SampleSecurityPrice []float64
	SampleTime          uint64
	BidOrder            *Order
	AskOrder            *Order
}

// Security sample value will be either margin or base or quote
// you can then use margin currency to eval portfolio value

func (sec *MarginSecurity) updateSampleValueChange(time uint64) {
	sampleValueChange := sec.SampleValueChange
	N := len(sampleValueChange)
	sampleMatchBid := sec.SellTradeModel.GetSampleMatchBid(time, N)
	sampleMatchAsk := sec.BuyTradeModel.GetSampleMatchAsk(time, N)
	sampleMarginPrice := sec.MarginPriceModel.GetSamplePrices(time, N)
	sampleSecurityPrice := sec.SecurityPriceModel.GetSamplePrices(time, N)

	if sec.BidOrder != nil {
		var price float64
		if sec.Security.IsInverse {
			price = 1. / sec.BidOrder.Price
		} else {
			price = sec.BidOrder.Price
		}
		queue := sec.BidOrder.Queue
		quantity := sec.BidOrder.Quantity
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * sec.Multiplier.Value
			fee := math.Abs(contractMarginValue) * sec.MakerFee.Value

			// Compute the cost
			cost := contractMarginValue + fee - (contractChange * sampleSecurityPrice[i] * sec.Multiplier.Value)
			sampleValueChange[i] -= cost * sampleMarginPrice[i]
		}
	}
	if sec.AskOrder != nil {
		var price float64
		if sec.Security.IsInverse {
			price = 1. / sec.AskOrder.Price
		} else {
			price = sec.AskOrder.Price
		}
		queue := sec.AskOrder.Queue
		quantity := sec.AskOrder.Quantity
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * sec.Multiplier.Value
			fee := math.Abs(contractMarginValue) * sec.MakerFee.Value

			// Compute the cost
			cost := -contractMarginValue + fee + (contractChange * sampleSecurityPrice[i] * sec.Multiplier.Value)
			sampleValueChange[i] -= cost * sampleMarginPrice[i]
		}
	}

	sec.SampleSecurityPrice = sampleSecurityPrice
	sec.SampleMarginPrice = sampleMarginPrice
	sec.SampleMatchBid = sampleMatchBid
	sec.SampleMatchAsk = sampleMatchAsk
	sec.SampleTime = time
}

func (sec *MarginSecurity) AddSampleValueChange(time uint64, values []float64) {
	sec.Lock()
	// Update this instrument sample value change only if bid order or ask order set
	if sec.BidOrder != nil || sec.AskOrder != nil {
		if sec.SampleTime != time {
			sec.updateSampleValueChange(time)
		}
		N := len(values)
		sampleValueChange := sec.SampleValueChange
		for i := 0; i < N; i++ {
			values[i] += sampleValueChange[i]
		}
	}

	sec.Unlock()
}

func (sec *MarginSecurity) GetELROnCancelBid(time uint64, values []float64, value float64) float64 {
	N := len(values)
	sec.Lock()

	if sec.SampleTime != time {
		sec.updateSampleValueChange(time)
	}

	sampleMatchBid := sec.SampleMatchBid
	sampleSecurityPrice := sec.SampleSecurityPrice
	sampleMarginPrice := sec.SampleMarginPrice

	// Remove bid order from value
	if sec.BidOrder != nil {
		var price float64
		if sec.Security.IsInverse {
			price = 1. / sec.BidOrder.Price
		} else {
			price = sec.BidOrder.Price
		}
		queue := sec.BidOrder.Queue
		quantity := sec.BidOrder.Quantity
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * sec.Multiplier.Value
			fee := math.Abs(contractMarginValue) * sec.MakerFee.Value

			// Remove cost from values
			cost := contractMarginValue + fee - (contractChange * sampleSecurityPrice[i] * sec.Multiplier.Value)
			values[i] += cost * sampleMarginPrice[i]
		}
	}
	sec.Unlock()

	expectedLogReturn := 0.
	for i := 0; i < N; i++ {
		expectedLogReturn += math.Log(values[i] / value)
	}
	expectedLogReturn /= float64(N)

	return expectedLogReturn
}

func (sec *MarginSecurity) GetELROnCancelAsk(time uint64, values []float64, value float64) float64 {
	N := len(values)
	sec.Lock()

	if sec.SampleTime != time {
		sec.updateSampleValueChange(time)
	}

	sampleMatchAsk := sec.SampleMatchAsk
	sampleSecurityPrice := sec.SampleSecurityPrice
	sampleMarginPrice := sec.SampleMarginPrice

	// Remove ask order from value
	if sec.AskOrder != nil {
		var price float64
		if sec.Security.IsInverse {
			price = 1. / sec.AskOrder.Price
		} else {
			price = sec.AskOrder.Price
		}
		queue := sec.AskOrder.Queue
		quantity := sec.AskOrder.Quantity
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * sec.Multiplier.Value
			fee := math.Abs(contractMarginValue) * sec.MakerFee.Value

			// Remove cost from values
			cost := -contractMarginValue + fee + (contractChange * sampleSecurityPrice[i] * sec.Multiplier.Value)
			values[i] += cost * sampleMarginPrice[i]
		}
	}
	sec.Unlock()

	expectedLogReturn := 0.
	for i := 0; i < N; i++ {
		expectedLogReturn += math.Log(values[i] / value)
	}
	expectedLogReturn /= float64(N)

	return expectedLogReturn
}

func (sec *MarginSecurity) GetELROnBidChange(time uint64, values []float64, value float64, prices []float64, queues []float64, maxQuote float64) (float64, *Order) {
	N := len(values)
	sec.Lock()
	// We want to see which option is the best, update, do nothing, or cancel
	if sec.SampleTime != time {
		sec.updateSampleValueChange(time)
	}

	sampleMatchBid := sec.SampleMatchBid
	sampleSecurityPrice := sec.SampleSecurityPrice
	sampleMarginPrice := sec.SampleMarginPrice

	var maxOrder *Order
	maxExpectedLogReturn := -999.

	// If we have an a bid order, we remove it from values
	if sec.BidOrder != nil {
		var price float64
		if sec.Security.IsInverse {
			price = 1. / sec.BidOrder.Price
		} else {
			price = sec.BidOrder.Price
		}
		expectedLogReturn := 0.
		queue := sec.BidOrder.Queue
		quantity := sec.BidOrder.Quantity
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * sec.Multiplier.Value
			fee := math.Abs(contractMarginValue) * sec.MakerFee.Value

			// Remove cost from values
			cost := contractMarginValue + fee - (contractChange * sampleSecurityPrice[i] * sec.Multiplier.Value)
			values[i] += cost * sampleMarginPrice[i]
		}
		// Compute ELR with current bid
		expectedLogReturn /= float64(N)
		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = sec.BidOrder
		}
	}

	// Compute ELR without any order
	expectedLogReturn := 0.
	for i := 0; i < N; i++ {
		expectedLogReturn += math.Log(values[i] / value)
	}
	expectedLogReturn /= float64(N)
	if expectedLogReturn > maxExpectedLogReturn {
		maxExpectedLogReturn = expectedLogReturn
		maxOrder = nil
	}

	for l := 0; l < len(prices); l++ {
		var price float64
		if sec.IsInverse {
			price = 1. / prices[l]
		} else {
			price = prices[l]
		}
		queue := queues[l]
		expectedLogReturn := 0.

		// TODO available margin
		for i := 0; i < N; i++ {
			contractChange := math.Max(sampleMatchBid[i]-queue, 0.)
			contractMarginValue := contractChange * price * sec.Multiplier.Value
			fee := math.Abs(contractMarginValue) * sec.MakerFee.Value
			cost := contractMarginValue + fee - (contractChange * sampleSecurityPrice[i] * sec.Multiplier.Value)
			// Compute ELR with cost
			expectedLogReturn += math.Log((values[i] - cost*sampleMarginPrice[i]) / value)
		}
		expectedLogReturn /= float64(N)

		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = &Order{
				Price:    price,
				Quantity: 0., // Computed at the end
				Queue:    queue,
			}
		}
	}

	if maxOrder != nil && maxOrder != sec.BidOrder {
		// Compute recommended order quantity based on price and queue
		// Recommended order quantity is the expected match at the level
		expectedMatch := 0.
		for i := 0; i < N; i++ {
			expectedMatch += math.Max(sampleMatchBid[i]-maxOrder.Queue, 0.)
		}
		expectedMatch /= float64(N)
		maxOrder.Quantity = math.Min(expectedMatch, maxQuote/maxOrder.Price)
	}
	sec.Unlock()

	return maxExpectedLogReturn, maxOrder
}

func (sec *MarginSecurity) GetELROnAskChange(time uint64, values []float64, value float64, prices []float64, queues []float64, maxBase float64) (float64, *Order) {
	N := len(values)
	sec.Lock()
	// We want to see which option is the best, update, do nothing, or cancel
	if sec.SampleTime != time {
		sec.updateSampleValueChange(time)
	}

	sampleMatchAsk := sec.SampleMatchAsk
	sampleSecurityPrice := sec.SampleSecurityPrice
	sampleMarginPrice := sec.SampleMarginPrice

	var maxOrder *Order
	maxExpectedLogReturn := -999.

	// Remove ask order from value
	if sec.AskOrder != nil {
		var price float64
		if sec.Security.IsInverse {
			price = 1. / sec.AskOrder.Price
		} else {
			price = sec.AskOrder.Price
		}
		queue := sec.AskOrder.Queue
		quantity := sec.AskOrder.Quantity
		expectedLogReturn := 0.
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * sec.Multiplier.Value
			fee := math.Abs(contractMarginValue) * sec.MakerFee.Value

			// Remove the cost from values
			cost := -contractMarginValue + fee + (contractChange * sampleSecurityPrice[i] * sec.Multiplier.Value)
			values[i] += cost * sampleMarginPrice[i]
		}
		// Compute ELR with current order
		expectedLogReturn /= float64(N)
		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = sec.AskOrder
		}
	}

	// Compute ELR without any order
	expectedLogReturn := 0.
	for i := 0; i < N; i++ {
		expectedLogReturn += math.Log(values[i] / value)
	}
	expectedLogReturn /= float64(N)
	if expectedLogReturn > maxExpectedLogReturn {
		maxExpectedLogReturn = expectedLogReturn
		maxOrder = nil
	}

	for l := 0; l < len(prices); l++ {
		var price float64
		if sec.IsInverse {
			price = 1. / prices[l]
		} else {
			price = prices[l]
		}
		queue := queues[l]
		expectedLogReturn := 0.

		// We need to compute what would happen if we changed
		// our ask order on the instrument.
		// TODO available margin
		for i := 0; i < N; i++ {
			contractChange := math.Max(sampleMatchAsk[i]-queue, 0.)
			contractMarginValue := contractChange * price * sec.Multiplier.Value
			fee := math.Abs(contractMarginValue) * sec.MakerFee.Value

			cost := -contractMarginValue + fee + (contractChange * sampleSecurityPrice[i] * sec.Multiplier.Value)
			// Compute ELR with new cost
			expectedLogReturn += math.Log((values[i] - cost*sampleMarginPrice[i]) / value)
		}

		expectedLogReturn /= float64(N)

		//fmt.Println(expectedLogReturn, maxExpectedLogReturn)
		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = &Order{
				Price:    price,
				Quantity: 0., // Computed at the end
				Queue:    queue,
			}
		}
	}

	if maxOrder != nil && maxOrder != sec.AskOrder {
		// Compute recommended order quantity based on price and queue
		// Recommended order quantity is the expected match at the level
		//fmt.Println("RETS", baseline, maxExpectedLogReturn)
		expectedMatch := 0.
		for i := 0; i < N; i++ {
			expectedMatch += math.Max(sampleMatchAsk[i]-maxOrder.Queue, 0.)
		}
		expectedMatch /= float64(N)
		maxOrder.Quantity = math.Min(expectedMatch, maxBase)
	}
	sec.Unlock()

	return maxExpectedLogReturn, maxOrder
}

func NewAccountPortFolio(exchange *xchangerModels.Exchange, marginCurrency *xchangerModels.Asset, securities []*models.Security, assetPriceModels map[uint32]PriceModel, securityPriceModels map[uint64]PriceModel, buyTradeModels map[uint64]BuyTradeModel, sellTradeModels map[uint64]SellTradeModel, sampleSize int) (*AccountPortfolio, error) {

	p := &AccountPortfolio{
		exchange:   exchange,
		sampleSize: sampleSize,
		balances:   make(map[uint32]float64),
		positions:  make(map[uint64]*Position),
		securities: make(map[uint64]Security),
		margin:     0.,
	}

	for _, security := range securities {
		ID := security.SecurityID
		if _, ok := buyTradeModels[ID]; !ok {
			return nil, fmt.Errorf("no trade model for security %d", ID)
		}
		if _, ok := sellTradeModels[ID]; !ok {
			return nil, fmt.Errorf("no trade model for security %d", ID)
		}
		switch security.SecurityType {
		case enum.SecurityType_CRYPTO_PERP:
			if _, ok := securityPriceModels[ID]; !ok {
				return nil, fmt.Errorf("no price model for security %d", ID)
			}
			if _, ok := assetPriceModels[marginCurrency.ID]; !ok {
				return nil, fmt.Errorf("no price model for asset %s", security.QuoteCurrency.Name)
			}
			p.securities[ID] = &MarginSecurity{
				RWMutex:             sync.RWMutex{},
				Security:            security,
				MarginPriceModel:    assetPriceModels[marginCurrency.ID],
				SecurityPriceModel:  securityPriceModels[ID],
				BuyTradeModel:       buyTradeModels[ID],
				SellTradeModel:      sellTradeModels[ID],
				SampleValueChange:   make([]float64, p.sampleSize, p.sampleSize),
				SampleMatchBid:      make([]float64, p.sampleSize, p.sampleSize),
				SampleMatchAsk:      make([]float64, p.sampleSize, p.sampleSize),
				SampleMarginPrice:   make([]float64, p.sampleSize, p.sampleSize),
				SampleSecurityPrice: make([]float64, p.sampleSize, p.sampleSize),
				SampleTime:          0,
				BidOrder:            nil,
				AskOrder:            nil,
			}
		}
	}

	return p, nil
}

func (p *AccountPortfolio) Sync(orders []*models.Order, positions []*models.Position, balances []*models.Balance, margin float64) error {
	p.margin = margin

	for _, o := range orders {
		ord := &Order{
			Order:          o,
			previousStatus: o.OrderStatus,
		}
		accnt.ordersID[o.OrderID] = ord
		accnt.ordersClID[o.ClientOrderID] = ord
	}

	// Reset securities
	for _, s := range p.securities {
		s.Position.Cost = 0.
		s.Position.RawSize = 0
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
		rawSize := int64(math.Round(sec.LotPrecision * p.Quantity))
		sec.Position.RawSize = rawSize
		sec.Position.Cross = p.Cross
		sec.Position.Cost = int64(p.Cost * accnt.marginPrecision)
	}

	for _, b := range balances {
		accnt.balances[b.Asset.ID] = b.Quantity
	}

	return nil
}
