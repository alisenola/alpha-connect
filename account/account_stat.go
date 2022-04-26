package account

import (
	"gitlab.com/alphaticks/alpha-connect/modeling"
	"math"
)

// TODO what about updating securities ?

func (accnt *Account) Value(model modeling.Market) float64 {
	accnt.RLock()
	defer accnt.RUnlock()
	value := 0.
	for k, v := range accnt.balances {
		pp, ok := model.GetPairPrice(k, accnt.quoteCurrency.ID)
		if ok {
			value += (float64(v) / accnt.MarginPrecision) * pp
		}
	}

	if accnt.MarginCurrency != nil {
		marginPrice, ok := model.GetPairPrice(accnt.MarginCurrency.ID, accnt.quoteCurrency.ID)
		if !ok {
			panic("no price for margin currency")
		}
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
			pp, ok := model.GetPrice(k)
			if !ok {
				panic("no price for position")
			}
			value -= (cost - math.Pow(pp, exp)*factor) * marginPrice
		}

		margin := float64(accnt.margin) / accnt.MarginPrecision
		value += margin * marginPrice
	}

	return math.Max(value, 0.)
}

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

func (accnt *Account) GetLeverage(model modeling.Market) float64 {
	accnt.RLock()
	defer accnt.RUnlock()
	var availableMargin uint64 = 0
	if accnt.margin < 0 {
		availableMargin -= uint64(-accnt.margin)
	} else {
		availableMargin += uint64(accnt.margin)
	}
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
			usedMargin += (float64(p.rawSize) / p.lotPrecision) * (1. / exitPrice) * math.Abs(p.multiplier)
		} else {
			usedMargin += (float64(p.rawSize) / p.lotPrecision) * exitPrice * math.Abs(p.multiplier)
		}
	}
	return usedMargin / (float64(availableMargin) / accnt.MarginPrecision)
}

func (accnt *Account) GetAvailableMargin(model modeling.Market, leverage float64) float64 {
	availableMargin := accnt.GetMargin(model)
	accnt.RLock()
	defer accnt.RUnlock()
	availableMargin *= leverage
	//fmt.Println("AV MARGIN", availableMargin)
	for k, p := range accnt.positions {
		// Entry price not defined if size = 0, division by 0 !
		if p.rawSize == 0 {
			continue
		}
		exitPrice, ok := model.GetPrice(k)
		if !ok {
			panic("no price for position")
		}
		cost := float64(p.cost) / accnt.MarginPrecision

		if p.inverse {
			unrealizedPnL := (1./exitPrice)*p.multiplier*(float64(p.rawSize)/p.lotPrecision) - cost
			// Cannot use unrealized profit in margin
			// TODO ? unrealizedPnL = math.Min(unrealizedPnL, 0)
			// Remove leveraged entry value and add PnL
			availableMargin = availableMargin - math.Abs(cost) + unrealizedPnL
		} else {
			unrealizedPnL := exitPrice*p.multiplier*(float64(p.rawSize)/p.lotPrecision) - cost
			// Cannot use unrealized profit in margin
			// TODO ? unrealizedPnL = math.Min(unrealizedPnL, 0)
			// Remove leveraged entry value and add PnL
			//fmt.Println("REMOVE 1", availableMargin)
			availableMargin = availableMargin - math.Abs(cost) + unrealizedPnL
		}
	}

	// TODO remove open order from av margin
	/*
		for k, o := range accnt.ordersID {
			if (o.OrderType == models.Limit) && (o.OrderStatus == models.PendingNew || o.OrderStatus == models.PendingReplace || o.OrderStatus == models.PartiallyFilled || o.OrderStatus == models.New) {

			}
		}
	*/

	return math.Max(availableMargin, 0.)
}

// TODO improve speed of that one, called often
func (accnt *Account) UpdateAskOrderQueue(securityID uint64, orderID string, queue float64) {
	accnt.securities[securityID].UpdateAskOrderQueue(orderID, queue)
}

func (accnt *Account) UpdateBidOrderQueue(securityID uint64, orderID string, queue float64) {
	accnt.securities[securityID].UpdateBidOrderQueue(orderID, queue)
}
