package account

import (
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/modeling"
	"math"
)

// TODO what about updating securities ?

func (accnt *Account) AddSampleValues(model modeling.MarketModel, time uint64, values []float64) {
	accnt.RLock()
	defer accnt.RUnlock()
	N := len(values)
	for k, v := range accnt.balances {
		samplePrices := model.GetSamplePairPrices(k, accnt.quoteCurrency.ID, time, N)
		for i := 0; i < N; i++ {
			values[i] += (float64(v) / accnt.MarginPrecision) * samplePrices[i]
		}
	}
	if accnt.MarginCurrency != nil {
		marginPrices := model.GetSamplePairPrices(accnt.MarginCurrency.ID, accnt.quoteCurrency.ID, time, N)
		for k, p := range accnt.positions {
			if p.rawSize == 0 {
				continue
			}
			cost := float64(p.cost) / p.marginPrecision
			size := float64(p.rawSize) / p.lotPrecision
			factor := size * p.multiplier
			var exp float64
			if p.inverse {
				exp = -1
			} else {
				exp = 1
			}
			samplePrices := model.GetSamplePrices(k, time, N)
			for i := 0; i < N; i++ {
				values[i] -= (cost - math.Pow(samplePrices[i], exp)*factor) * marginPrices[i]
			}
		}
		margin := float64(accnt.margin) / accnt.MarginPrecision
		for i := 0; i < N; i++ {
			values[i] += margin * marginPrices[i]
		}
	}

	for _, s := range accnt.securities {
		s.AddSampleValueChange(model, time, values)
	}
	// TODO handle neg values
}

func (accnt *Account) GetELROnCancelBid(securityID uint64, orderID string, model modeling.MarketModel, time uint64, values []float64, value float64) float64 {
	return accnt.securities[securityID].GetELROnCancelBid(orderID, model, time, values, value)
}

func (accnt *Account) GetELROnCancelAsk(securityID uint64, orderID string, model modeling.MarketModel, time uint64, values []float64, value float64) float64 {
	return accnt.securities[securityID].GetELROnCancelAsk(orderID, model, time, values, value)
}

func (accnt *Account) GetELROnLimitBid(securityID uint64, model modeling.MarketModel, time uint64, values []float64, value float64, prices []float64, queues []float64, maxQuote float64) (float64, *COrder) {
	return accnt.securities[securityID].GetELROnLimitBidChange("", model, time, values, value, prices, queues, maxQuote)
}

func (accnt *Account) GetELROnLimitAsk(securityID uint64, model modeling.MarketModel, time uint64, values []float64, value float64, prices []float64, queues []float64, maxBase float64) (float64, *COrder) {
	return accnt.securities[securityID].GetELROnLimitAskChange("", model, time, values, value, prices, queues, maxBase)
}

func (accnt *Account) GetELROnLimitBidChange(securityID uint64, orderID string, model modeling.MarketModel, time uint64, values []float64, value float64, prices []float64, queues []float64, maxQuote float64) (float64, *COrder) {
	return accnt.securities[securityID].GetELROnLimitBidChange(orderID, model, time, values, value, prices, queues, maxQuote)
}

func (accnt *Account) GetELROnLimitAskChange(securityID uint64, orderID string, model modeling.MarketModel, time uint64, values []float64, value float64, prices []float64, queues []float64, maxBase float64) (float64, *COrder) {
	return accnt.securities[securityID].GetELROnLimitAskChange(orderID, model, time, values, value, prices, queues, maxBase)
}

func (accnt *Account) GetELROnMarketBuy(securityID uint64, model modeling.MarketModel, time uint64, values []float64, value, price, quantity, maxQuantity float64) (float64, *COrder) {
	return accnt.securities[securityID].GetELROnMarketBuy(model, time, values, value, price, quantity, maxQuantity)
}

func (accnt *Account) GetELROnMarketSell(securityID uint64, model modeling.MarketModel, time uint64, values []float64, value, price, quantity, maxQuantity float64) (float64, *COrder) {
	return accnt.securities[securityID].GetELROnMarketSell(model, time, values, value, price, quantity, maxQuantity)
}

func (accnt *Account) GetExposure(asset uint32) float64 {
	accnt.RLock()
	defer accnt.RUnlock()
	exposure := float64(accnt.balances[asset]) / accnt.MarginPrecision
	for _, pos := range accnt.baseToPositions[asset] {
		if pos.rawSize == 0 {
			continue
		}
		exposure += pos.Size()
	}
	return exposure
}

func (accnt *Account) GetLeverage(model modeling.Market) float64 {
	accnt.RLock()
	defer accnt.RUnlock()
	var availableMargin int64 = 0
	availableMargin += accnt.margin
	availableMargin += accnt.balances[accnt.MarginCurrency.ID]

	usedMargin := 0.
	for k, p := range accnt.positions {
		if p.rawSize == 0 {
			continue
		}
		exitPrice, ok := model.GetPrice(k)
		if !ok {
			panic("no price for position")
		}
		if p.inverse {
			usedMargin += (math.Abs(float64(p.rawSize)) / p.lotPrecision) * (1. / exitPrice) * p.multiplier
		} else {
			usedMargin += (math.Abs(float64(p.rawSize)) / p.lotPrecision) * exitPrice * p.multiplier
		}
	}
	return usedMargin / (float64(availableMargin) / accnt.MarginPrecision)
}

func (accnt *Account) GetMarginInfo(model modeling.Market) (float64, float64, float64) {
	margin := accnt.GetMargin(model)
	accnt.RLock()
	defer accnt.RUnlock()

	var longMargin, shortMargin float64
	for k, p := range accnt.positions {
		if p.rawSize == 0 {
			continue
		}
		price, ok := model.GetPrice(k)
		if !ok {
			panic("no price for position")
		}
		if p.inverse {
			notional := (float64(p.rawSize) / p.lotPrecision) * (1. / price) * p.multiplier
			if notional < 0 {
				shortMargin -= notional
			} else {
				longMargin += notional
			}
		} else {
			notional := (float64(p.rawSize) / p.lotPrecision) * price * p.multiplier
			if notional < 0 {
				shortMargin -= notional
			} else {
				longMargin += notional
			}
		}
	}
	return margin, longMargin, shortMargin
}

func (accnt *Account) GetNetMargin(model modeling.Market) (float64, error) {
	netMargin := accnt.GetMargin(model)
	accnt.RLock()
	defer accnt.RUnlock()
	//fmt.Println("AV MARGIN", availableMargin)
	for k, p := range accnt.positions {
		// Entry price not defined if size = 0, division by 0 !
		if p.rawSize == 0 {
			continue
		}
		exitPrice, ok := model.GetPrice(k)
		if !ok {
			return 0., fmt.Errorf("no price for position %d", k)
		}
		cost := float64(p.cost) / accnt.MarginPrecision

		if p.inverse {
			unrealizedPnL := (1./exitPrice)*p.multiplier*(float64(p.rawSize)/p.lotPrecision) - cost
			netMargin += unrealizedPnL
		} else {
			unrealizedPnL := exitPrice*p.multiplier*(float64(p.rawSize)/p.lotPrecision) - cost
			netMargin += unrealizedPnL
		}
	}

	return math.Max(netMargin, 0.), nil
}

// TODO improve speed of that one, called often
func (accnt *Account) UpdateAskOrderQueue(securityID uint64, orderID string, queue float64) {
	accnt.securities[securityID].UpdateAskOrderQueue(orderID, queue)
}

func (accnt *Account) UpdateBidOrderQueue(securityID uint64, orderID string, queue float64) {
	accnt.securities[securityID].UpdateBidOrderQueue(orderID, queue)
}

func (accnt *Account) UpdateMarkPrice(securityID uint64, price float64) {
	accnt.positions[securityID].UpdateMarkPrice(price)
}
