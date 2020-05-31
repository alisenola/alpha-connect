package account

import (
	"gitlab.com/alphaticks/alphac/modeling"
	"gitlab.com/alphaticks/alphac/models"
	"math"
)

// TODO what about updating securities ?

func (accnt *Account) Value(model modeling.Model) float64 {
	value := 0.
	for k, v := range accnt.balances {
		value += v * model.GetAssetPrice(k)
	}
	marginPrice := model.GetAssetPrice(accnt.marginCurrency.ID)
	for k, v := range accnt.positions {
		cost := float64(v.cost) / v.marginPrecision
		size := float64(v.rawSize) / v.lotPrecision
		if v.inverse {
			value += (cost - (1./model.GetSecurityPrice(k))*size*v.multiplier) * marginPrice
		} else {
			value += (cost - model.GetSecurityPrice(k)*size*v.multiplier) * marginPrice
		}
	}

	margin := float64(accnt.margin) / accnt.marginPrecision
	value += margin * marginPrice

	return math.Max(value, 0.)
}

func (accnt *Account) AddSampleValues(model modeling.Model, time uint64, values []float64) {
	N := len(values)
	for k, v := range accnt.balances {
		samplePrices := model.GetSampleAssetPrices(k, time, N)
		for i := 0; i < N; i++ {
			values[i] += v * samplePrices[i]
		}
	}
	marginPrices := model.GetSampleAssetPrices(accnt.marginCurrency.ID, time, N)
	for k, v := range accnt.positions {
		cost := float64(v.cost) / v.marginPrecision
		size := float64(v.rawSize) / v.lotPrecision
		factor := size * v.multiplier
		var exp float64
		if v.inverse {
			exp = -1
		} else {
			exp = 1
		}
		samplePrices := model.GetSampleSecurityPrices(k, time, N)
		for i := 0; i < N; i++ {
			values[i] -= (cost - math.Pow(samplePrices[i], exp)*factor) * marginPrices[i]
		}
	}

	margin := float64(accnt.margin) / accnt.marginPrecision
	for i := 0; i < N; i++ {
		values[i] += margin * marginPrices[i]
	}
	for _, s := range accnt.securities {
		s.AddSampleValueChange(model, time, values)
	}
	// TODO handle neg values
}

func (accnt Account) GetELROnCancelBid(orderID string, model modeling.Model, time uint64, values []float64, value float64) float64 {
	if order, ok := accnt.ordersClID[orderID]; ok {
		return accnt.securities[order.Instrument.SecurityID.Value].GetELROnCancelBid(orderID, model, time, values, value)
	} else {
		// TODO compute ELR using given values and val
		return 0.
	}
}

func (accnt Account) GetELROnCancelAsk(orderID string, model modeling.Model, time uint64, values []float64, value float64) float64 {
	if order, ok := accnt.ordersClID[orderID]; ok {
		return accnt.securities[order.Instrument.SecurityID.Value].GetELROnCancelAsk(orderID, model, time, values, value)
	} else {
		// TODO
		return 0.
	}
}

func (accnt Account) GetELROnLimitBid(securityID uint64, model modeling.Model, time uint64, values []float64, value float64, prices []float64, queues []float64, maxQuote float64) (float64, *COrder) {
	return accnt.securities[securityID].GetELROnLimitBidChange("", model, time, values, value, prices, queues, maxQuote)
}

func (accnt Account) GetELROnLimitAsk(securityID uint64, model modeling.Model, time uint64, values []float64, value float64, prices []float64, queues []float64, maxBase float64) (float64, *COrder) {
	return accnt.securities[securityID].GetELROnLimitAskChange("", model, time, values, value, prices, queues, maxBase)
}

func (accnt Account) GetELROnLimitBidChange(orderID string, model modeling.Model, time uint64, values []float64, value float64, prices []float64, queues []float64, maxQuote float64) (float64, *COrder) {
	if order, ok := accnt.ordersClID[orderID]; ok {
		return accnt.securities[order.Instrument.SecurityID.Value].GetELROnLimitBidChange(orderID, model, time, values, value, prices, queues, maxQuote)
	} else {
		// TODO ?
		return 0., nil
	}
}

func (accnt Account) GetELROnLimitAskChange(orderID string, model modeling.Model, time uint64, values []float64, value float64, prices []float64, queues []float64, maxBase float64) (float64, *COrder) {
	if order, ok := accnt.ordersClID[orderID]; ok {
		return accnt.securities[order.Instrument.SecurityID.Value].GetELROnLimitAskChange(orderID, model, time, values, value, prices, queues, maxBase)
	} else {
		// TODO ?
		return 0., nil
	}
}

func (accnt *Account) GetLeverage(model modeling.Model) float64 {
	availableMargin := accnt.balances[accnt.marginCurrency.ID] + float64(accnt.margin)/accnt.marginPrecision
	usedMargin := 0.
	for k, p := range accnt.positions {
		exitPrice := model.GetSecurityPrice(k)
		if p.inverse {
			usedMargin += (float64(p.rawSize) / p.lotPrecision) * (1. / exitPrice) * math.Abs(p.multiplier)
		} else {
			usedMargin += (float64(p.rawSize) / p.lotPrecision) * exitPrice * math.Abs(p.multiplier)
		}
	}
	return usedMargin / availableMargin
}

func (accnt *Account) GetAvailableMargin(model modeling.Model, leverage float64) float64 {
	availableMargin := accnt.balances[accnt.marginCurrency.ID] + float64(accnt.margin)/accnt.marginPrecision
	for k, p := range accnt.positions {
		// Entry price not defined if size = 0, division by 0 !
		if p.rawSize == 0 {
			continue
		}
		exitPrice := model.GetSecurityPrice(k)
		cost := float64(p.cost) / accnt.marginPrecision

		if p.inverse {
			unrealizedPnL := (1./exitPrice)*p.multiplier*(float64(p.rawSize)/p.lotPrecision) - cost
			// Cannot use unrealized profit in margin
			// TODO ? unrealizedPnL = math.Min(unrealizedPnL, 0)
			// Remove leveraged entry value and add PnL
			availableMargin = availableMargin - (math.Abs(cost) / leverage) + unrealizedPnL
		} else {
			unrealizedPnL := exitPrice*p.multiplier*(float64(p.rawSize)/p.lotPrecision) - cost
			// Cannot use unrealized profit in margin
			// TODO ? unrealizedPnL = math.Min(unrealizedPnL, 0)
			// Remove leveraged entry value and add PnL
			availableMargin = availableMargin - (math.Abs(cost) / leverage) + unrealizedPnL
		}
	}

	return math.Max(availableMargin, 0.)
}

// TODO improve speed of that one, called often
func (accnt *Account) UpdateQueue(orderID string, queue float64) {
	if o, ok := accnt.ordersClID[orderID]; ok {
		if o.Side == models.Buy {
			accnt.securities[o.Instrument.SecurityID.Value].UpdateBidOrderQueue(orderID, queue)
		} else {
			accnt.securities[o.Instrument.SecurityID.Value].UpdateAskOrderQueue(orderID, queue)
		}
	}
}
