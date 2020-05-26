package account

import (
	"gitlab.com/alphaticks/alphac/modeling"
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
		size := float64(v.rawSize) / v.marginPrecision
		value += (cost - model.GetSecurityPrice(k)*size*v.multiplier) * marginPrice
	}

	margin := float64(accnt.margin) / accnt.marginPrecision
	value += margin * marginPrice

	return value
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
		size := float64(v.rawSize) / v.marginPrecision
		factor := size * v.multiplier
		samplePrices := model.GetSampleSecurityPrices(k, time, N)
		for i := 0; i < N; i++ {
			values[i] -= (cost - samplePrices[i]*factor) * marginPrices[i]
		}
	}

	margin := float64(accnt.margin) / accnt.marginPrecision
	for i := 0; i < N; i++ {
		values[i] += margin * marginPrices[i]
	}
	for _, s := range accnt.securities {
		s.AddSampleValueChange(model, time, values)
	}
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
