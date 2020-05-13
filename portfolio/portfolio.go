package portfolio

import (
	"math"
)

// The portfolio data structure contains all the open orders and positions
// it is used to compute the expected log returns if you perform an action
// like cancelling an order or creating a new one.

// The actors using this DS must use it to compute all the possible
// expected log returns using available actions, and then choose the best action
// and update the portfolio accordingly

// Each time the OB changes, we must also update the queue behind an order

// Each time an order is filled or partially filled, we must update the portfolio too

// NB: This is one of the most critical data structure and it needs to be as efficient as
// possible. Try to prevent any heap allocation call here ?

type Order struct {
	Price    float64
	Quantity float64
	Queue    float64
}

type ExchangePortfolio interface {
	Value() float64
	GetAsset(asset uint32) float64
	AddSampleValues(time uint64, values []float64)
	GetELROnCancelBid(instrID uint64, time uint64, values []float64, value float64) float64
	GetELROnCancelAsk(instrID uint64, time uint64, values []float64, value float64) float64
	GetELROnBidChange(instrID uint64, time uint64, values []float64, value float64, prices []float64, queues []float64, maxQuantity float64) (float64, *Order)
	GetELROnAskChange(instrID uint64, time uint64, values []float64, value float64, prices []float64, queues []float64, maxQuantity float64) (float64, *Order)
}

type PriceModel interface {
	Update(feedID uint64, tick uint64, price float64)
	Progress(tick uint64)
	GetPrice() float64
	GetSamplePrices(time uint64, sampleSize int) []float64
	Frequency() uint64
}

type TradeModel interface {
	Update(feedID uint64, tick uint64, size float64)
	Progress(tick uint64)
}

type SellTradeModel interface {
	TradeModel
	GetSampleMatchBid(time uint64, sampleSize int) []float64
}

type BuyTradeModel interface {
	TradeModel
	GetSampleMatchAsk(time uint64, sampleSize int) []float64
}

type Portfolio struct {
	SampleSize int
	portfolios map[uint64]ExchangePortfolio
}

func NewPortfolio(portfolios map[uint64]ExchangePortfolio, sampleSize int) *Portfolio {
	p := &Portfolio{
		SampleSize: sampleSize,
		portfolios: portfolios,
	}

	return p
}

func (p *Portfolio) Value() float64 {
	value := 0.
	for _, exch := range p.portfolios {
		value += exch.Value()
	}

	return value
}

func (p *Portfolio) GetAsset(asset uint32) float64 {
	value := 0.
	for _, exch := range p.portfolios {
		value += exch.GetAsset(asset)
	}

	return value
}

// Return the expected log return
func (p *Portfolio) ExpectedLogReturn(time uint64) float64 {
	N := p.SampleSize

	var value float64 = 0
	values := make([]float64, N, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	// We need to ensure consistent sample prices across the computation
	for _, exch := range p.portfolios {
		exch.AddSampleValues(time, values)
		value += exch.Value()
	}

	// TODO what to do if value is 0 ?
	var expectedLogReturn float64 = 0

	// Compute expected log return
	for i := 0; i < N; i++ {
		expectedLogReturn += math.Log(values[i] / value)
	}

	return expectedLogReturn / float64(N)
}

/*
// Returns the order with the highest expected log return on cancel
func (p *Portfolio) ExpectedLogReturnOnCancel(time uint64, exchange uint64, instrID uint64) float64 {
	N := p.SampleSize

	var value float64 = 0
	values := make([]float64, N, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	for _, exch := range p.portfolios {
		exch.AddSampleValues(time, values)
		value += exch.Value()
	}

	// Compute expected log return if you don't cancel anything
	maxExpectedLogReturn := 0.
	for i := 0; i < N; i++ {
		maxExpectedLogReturn += math.Log(values[i] / value)
	}
	maxExpectedLogReturn /= float64(N)

	var maxOrder *Order

	target := p.portfolios[exchange]
	elr := target.GetELROnCancelBid(time, values, value, instrID)
	if elr > maxExpectedLogReturn {
		maxExpectedLogReturn = elr
	}

	elr := target.GetELROnCancelAsk(time, values, value, instrID)
	if elr > maxExpectedLogReturn {
		maxExpectedLogReturn = elr
	}

	return maxExpectedLogReturn, true
}
*/
func (p *Portfolio) ELROnCancelBid(time uint64, exchange uint64, instrumentID uint64) float64 {
	// Need to compute the expected log return
	N := p.SampleSize
	var value float64 = 0
	values := make([]float64, N, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	for _, exch := range p.portfolios {
		exch.AddSampleValues(time, values)
		value += exch.Value()
	}
	target := p.portfolios[exchange]

	return target.GetELROnCancelBid(instrumentID, time, values, value)
}

func (p *Portfolio) ELROnCancelAsk(time uint64, exchange uint64, instrumentID uint64) float64 {
	// Need to compute the expected log return
	N := p.SampleSize
	var value float64 = 0
	values := make([]float64, N, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	for _, exch := range p.portfolios {
		exch.AddSampleValues(time, values)
		value += exch.Value()
	}
	target := p.portfolios[exchange]

	return target.GetELROnCancelAsk(instrumentID, time, values, value)
}

// We make multiple assumptions here. At each level, we assume the expectedMatch will be what will be set as order quantity
// This assumption kind of make sense, as if, in case of bid, the expected match is big at one level, might as well place the order
// at the level bellow to get a higher spread.
// When at ask, a // TODO think more about that
func (p *Portfolio) ELROnLimitBid(time uint64, exchange uint64, instrumentID uint64, prices []float64, queues []float64, maxQuantity float64) (float64, *Order) {
	// Need to compute the expected log return
	N := p.SampleSize
	var value float64 = 0
	values := make([]float64, N, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	for _, exch := range p.portfolios {
		exch.AddSampleValues(time, values)
		value += exch.Value()
	}
	target := p.portfolios[exchange]

	return target.GetELROnBidChange(instrumentID, time, values, value, prices, queues, maxQuantity)
}

func (p *Portfolio) ELROnLimitAsk(time uint64, exchange uint64, instrumentID uint64, prices []float64, queues []float64, maxQuantity float64) (float64, *Order) {
	// Need to compute the expected log return
	N := p.SampleSize

	var value float64 = 0
	values := make([]float64, N, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	for _, exch := range p.portfolios {
		exch.AddSampleValues(time, values)
		value += exch.Value()
	}

	target := p.portfolios[exchange]

	return target.GetELROnAskChange(instrumentID, time, values, value, prices, queues, maxQuantity)
}
