package account

import (
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/modeling"
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

type Portfolio struct {
	sampleSize        int
	accountPortfolios map[string]*Account
}

func NewPortfolio(sampleSize int) *Portfolio {
	p := &Portfolio{
		sampleSize:        sampleSize,
		accountPortfolios: make(map[string]*Account),
	}

	return p
}

func (p *Portfolio) GetAccount(ID string) *Account {
	return p.accountPortfolios[ID]
}

func (p *Portfolio) AddAccount(account *Account) {
	p.accountPortfolios[account.Name] = account
}

func (p *Portfolio) GetMargin(model modeling.Market) (float64, error) {
	value := 0.
	for _, exch := range p.accountPortfolios {
		v := exch.GetMargin(model)
		value += v
	}

	return value, nil
}

func (p *Portfolio) Value(model modeling.Market) (float64, error) {
	value := 0.
	for _, exch := range p.accountPortfolios {
		v, err := exch.GetNetMargin(model)
		if err != nil {
			return 0, fmt.Errorf("error while computing net margin: %v", err)
		}
		value += v
	}

	return value, nil
}

func (p *Portfolio) GetExposure(asset uint32) float64 {
	var exposure float64 = 0
	for _, exch := range p.accountPortfolios {
		exposure += exch.GetExposure(asset)
	}

	return exposure
}

/*
func (p *Portfolio) GetAsset(asset uint32) float64 {
	value := 0.
	for _, exch := range p.portfolios {
		value += exch.GetAsset(asset)
	}

	return value
}
*/

// Return the expected log return
func (p *Portfolio) ExpectedLogReturn(model modeling.MarketModel, time uint64) (float64, error) {
	N := p.sampleSize

	var value float64 = 0
	values := make([]float64, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	// We need to ensure consistent sample prices across the computation
	for _, exch := range p.accountPortfolios {
		exch.AddSampleValues(model, time, values)
		v, err := exch.GetNetMargin(model)
		if err != nil {
			return 0, fmt.Errorf("error while computing net margin: %v", err)
		}
		value += v
	}

	// TODO what to do if value is 0 ?
	var expectedLogReturn float64 = 0

	// Compute expected log return
	for i := 0; i < N; i++ {
		expectedLogReturn += math.Log(values[i] / value)
	}

	return expectedLogReturn / float64(N), nil
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

func (p *Portfolio) GetELR(model modeling.MarketModel, time uint64) (float64, error) {
	N := p.sampleSize
	var value float64 = 0
	values := make([]float64, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	for _, exch := range p.accountPortfolios {
		exch.AddSampleValues(model, time, values)
		v, err := exch.GetNetMargin(model)
		if err != nil {
			return 0, fmt.Errorf("error while computing net margin: %v", err)
		}
		value += v
	}
	elr := 0.
	for i := 0; i < N; i++ {
		elr += math.Log(values[i] / value)
	}
	return elr / float64(N), nil
}

func (p *Portfolio) GetELROnCancelBid(accountID string, securityID uint64, orderID string, model modeling.MarketModel, time uint64) (float64, error) {
	// Need to compute the expected log return
	N := p.sampleSize
	var value float64 = 0
	values := make([]float64, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	for _, exch := range p.accountPortfolios {
		exch.AddSampleValues(model, time, values)
		v, err := exch.GetNetMargin(model)
		if err != nil {
			return 0, fmt.Errorf("error while computing net margin: %v", err)
		}
		value += v
	}
	target := p.accountPortfolios[accountID]

	return target.GetELROnCancelBid(securityID, orderID, model, time, values, value), nil
}

func (p *Portfolio) GetELROnCancelAsk(accountID string, securityID uint64, orderID string, model modeling.MarketModel, time uint64) (float64, error) {
	// Need to compute the expected log return
	N := p.sampleSize
	var value float64 = 0
	values := make([]float64, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	for _, exch := range p.accountPortfolios {
		exch.AddSampleValues(model, time, values)
		v, err := exch.GetNetMargin(model)
		if err != nil {
			return 0, fmt.Errorf("error while computing net margin: %v", err)
		}
		value += v
	}
	target := p.accountPortfolios[accountID]

	return target.GetELROnCancelAsk(securityID, orderID, model, time, values, value), nil
}

func (p *Portfolio) GetELROnLimitBid(accountID string, securityID uint64, model modeling.MarketModel, time uint64, prices []float64, queues []float64, maxQuantity float64) (float64, *COrder, error) {
	// Need to compute the expected log return
	N := p.sampleSize
	var value float64 = 0
	values := make([]float64, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	for _, exch := range p.accountPortfolios {
		exch.AddSampleValues(model, time, values)
		v, err := exch.GetNetMargin(model)
		if err != nil {
			return 0, nil, fmt.Errorf("error while computing net margin: %v", err)
		}
		value += v
	}
	target := p.accountPortfolios[accountID]
	elr, o := target.GetELROnLimitBid(securityID, model, time, values, value, prices, queues, maxQuantity)

	return elr, o, nil
}

func (p *Portfolio) GetELROnLimitAsk(accountID string, securityID uint64, model modeling.MarketModel, time uint64, prices []float64, queues []float64, maxQuantity float64) (float64, *COrder, error) {
	// Need to compute the expected log return
	N := p.sampleSize

	var value float64 = 0
	values := make([]float64, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	for _, exch := range p.accountPortfolios {
		exch.AddSampleValues(model, time, values)
		v, err := exch.GetNetMargin(model)
		if err != nil {
			return 0, nil, fmt.Errorf("error while computing net margin: %v", err)
		}
		value += v
	}

	target := p.accountPortfolios[accountID]
	elr, o := target.GetELROnLimitAsk(securityID, model, time, values, value, prices, queues, maxQuantity)

	return elr, o, nil
}

// We make multiple assumptions here. At each level, we assume the expectedMatch will be what will be set as order quantity
// This assumption kind of make sense, as if, in case of bid, the expected match is big at one level, might as well place the order
// at the level bellow to get a higher spread.
// When at ask, a // TODO think more about that
func (p *Portfolio) GetELROnLimitBidChange(accountID string, securityID uint64, orderID string, model modeling.MarketModel, time uint64, prices []float64, queues []float64, maxQuantity float64) (float64, *COrder, error) {
	// Need to compute the expected log return
	N := p.sampleSize
	var value float64 = 0
	values := make([]float64, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	for _, exch := range p.accountPortfolios {
		exch.AddSampleValues(model, time, values)
		v, err := exch.GetNetMargin(model)
		if err != nil {
			return 0, nil, fmt.Errorf("error while computing net margin: %v", err)
		}
		value += v
	}
	target := p.accountPortfolios[accountID]

	elr, o := target.GetELROnLimitBidChange(securityID, orderID, model, time, values, value, prices, queues, maxQuantity)
	return elr, o, nil
}

func (p *Portfolio) GetELROnLimitAskChange(accountID string, securityID uint64, orderID string, model modeling.MarketModel, time uint64, prices []float64, queues []float64, maxQuantity float64) (float64, *COrder, error) {
	// Need to compute the expected log return
	N := p.sampleSize

	var value float64 = 0
	values := make([]float64, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	for _, exch := range p.accountPortfolios {
		exch.AddSampleValues(model, time, values)
		v, err := exch.GetNetMargin(model)
		if err != nil {
			return 0, nil, fmt.Errorf("error while computing net margin: %v", err)
		}
		value += v
	}

	target := p.accountPortfolios[accountID]

	elr, o := target.GetELROnLimitAskChange(securityID, orderID, model, time, values, value, prices, queues, maxQuantity)
	return elr, o, nil
}

func (p *Portfolio) GetELROnMarketBuy(accountID string, securityID uint64, model modeling.MarketModel, time uint64, price float64, quantity float64, maxQuantity float64) (float64, *COrder, error) {
	// Need to compute the expected log return
	N := p.sampleSize

	var value float64 = 0
	values := make([]float64, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	for _, exch := range p.accountPortfolios {
		exch.AddSampleValues(model, time, values)
		v, err := exch.GetNetMargin(model)
		if err != nil {
			return 0, nil, fmt.Errorf("error while computing net margin: %v", err)
		}
		value += v
	}

	target := p.accountPortfolios[accountID]

	elr, o := target.GetELROnMarketBuy(securityID, model, time, values, value, price, quantity, maxQuantity)
	return elr, o, nil
}

func (p *Portfolio) GetELROnMarketSell(accountID string, securityID uint64, model modeling.MarketModel, time uint64, price float64, quantity float64, maxQuantity float64) (float64, *COrder, error) {
	// Need to compute the expected log return
	N := p.sampleSize

	var value float64 = 0
	values := make([]float64, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	for _, exch := range p.accountPortfolios {
		exch.AddSampleValues(model, time, values)
		v, err := exch.GetNetMargin(model)
		if err != nil {
			return 0, nil, fmt.Errorf("error while computing net margin: %v", err)
		}
		value += v
	}

	target := p.accountPortfolios[accountID]
	elr, o := target.GetELROnMarketSell(securityID, model, time, values, value, price, quantity, maxQuantity)

	return elr, o, nil
}
