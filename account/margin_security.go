package account

import (
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alphac/modeling"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"math"
	"sync"
)

type COrder struct {
	Price    float64
	Quantity float64
	Queue    float64
	ID       string
}

type MarginSecurity struct {
	sync.RWMutex
	*models.Security
	tickPrecision       float64
	lotPrecision        float64
	marginCurrency      *xchangerModels.Asset
	openBidOrders       map[string]*COrder
	openAskOrders       map[string]*COrder
	size                float64
	makerFee            float64
	takerFee            float64
	sampleValueChange   []float64
	sampleMatchBid      []float64
	sampleMatchAsk      []float64
	sampleMarginPrice   []float64
	sampleSecurityPrice []float64
	sampleTime          uint64
}

func NewMarginSecurity(sec *models.Security, marginCurrency *xchangerModels.Asset, makerFee, takerFee *float64) *MarginSecurity {
	m := &MarginSecurity{
		RWMutex:             sync.RWMutex{},
		Security:            sec,
		tickPrecision:       math.Ceil(1. / sec.MinPriceIncrement),
		lotPrecision:        math.Ceil(1. / sec.RoundLot),
		marginCurrency:      marginCurrency,
		openBidOrders:       make(map[string]*COrder),
		openAskOrders:       make(map[string]*COrder),
		size:                0.,
		sampleValueChange:   nil,
		sampleMatchBid:      nil,
		sampleMatchAsk:      nil,
		sampleMarginPrice:   nil,
		sampleSecurityPrice: nil,
		sampleTime:          0,
	}

	if makerFee != nil {
		m.makerFee = *makerFee
	} else {
		m.makerFee = sec.MakerFee.Value
	}
	if takerFee != nil {
		m.takerFee = *takerFee
	} else {
		m.takerFee = sec.TakerFee.Value
	}

	return m
}

func (sec *MarginSecurity) AddBidOrder(ID string, price, quantity, queue float64) {
	sec.Lock()
	sec.openBidOrders[ID] = &COrder{
		Price:    price,
		Quantity: quantity,
		Queue:    queue,
	}

	exp := 1.
	if sec.Security.IsInverse {
		price = 1. / price
		exp = -1
	}
	mul := sec.Multiplier.Value
	makerFee := sec.makerFee
	N := len(sec.sampleValueChange)
	sampleValueChange := sec.sampleValueChange
	sampleMatchBid := sec.sampleMatchBid
	sampleSecurityPrice := sec.sampleSecurityPrice
	sampleMarginPrice := sec.sampleMarginPrice
	for i := 0; i < N; i++ {
		contractChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
		contractMarginValue := contractChange * price * mul
		fee := math.Abs(contractMarginValue) * makerFee

		// Compute the cost
		cost := contractMarginValue + fee - (contractChange * math.Pow(sampleSecurityPrice[i], exp) * mul)
		sampleValueChange[i] -= cost * sampleMarginPrice[i]
	}
	sec.Unlock()
}

func (sec *MarginSecurity) AddAskOrder(ID string, price, quantity, queue float64) {
	sec.Lock()
	sec.openAskOrders[ID] = &COrder{
		Price:    price,
		Quantity: quantity,
		Queue:    queue,
	}
	exp := 1.
	if sec.Security.IsInverse {
		price = 1. / price
		exp = -1
	}
	mul := sec.Multiplier.Value
	makerFee := sec.makerFee
	sampleValueChange := sec.sampleValueChange
	sampleMatchAsk := sec.sampleMatchAsk
	sampleSecurityPrice := sec.sampleSecurityPrice
	sampleMarginPrice := sec.sampleMarginPrice
	N := len(sampleValueChange)
	for i := 0; i < N; i++ {
		contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
		contractMarginValue := contractChange * price * mul
		fee := math.Abs(contractMarginValue) * makerFee

		// Compute the cost
		cost := -contractMarginValue + fee + (contractChange * math.Pow(sampleSecurityPrice[i], exp) * sec.Multiplier.Value)
		sampleValueChange[i] -= cost * sampleMarginPrice[i]
	}
	sec.Unlock()
}

func (sec *MarginSecurity) RemoveBidOrder(ID string) {
	sec.Lock()
	if o, ok := sec.openBidOrders[ID]; ok {

		var price, exp float64
		if sec.Security.IsInverse {
			price = 1. / o.Price
			exp = -1
		} else {
			price = o.Price
			exp = 1
		}
		queue := o.Queue
		quantity := o.Quantity
		mul := sec.Multiplier.Value
		makerFee := sec.makerFee
		sampleValueChange := sec.sampleValueChange
		sampleMatchBid := sec.sampleMatchBid
		sampleSecurityPrice := sec.sampleSecurityPrice
		sampleMarginPrice := sec.sampleMarginPrice
		N := len(sampleValueChange)
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * mul
			fee := math.Abs(contractMarginValue) * makerFee

			// Remove cost from values
			cost := contractMarginValue + fee - (contractChange * math.Pow(sampleSecurityPrice[i], exp) * sec.Multiplier.Value)
			sampleValueChange[i] += cost * sampleMarginPrice[i]
		}
		delete(sec.openBidOrders, ID)
	}
	sec.Unlock()
}

func (sec *MarginSecurity) RemoveAskOrder(ID string) {
	sec.Lock()
	if o, ok := sec.openAskOrders[ID]; ok {
		var price, exp float64
		if sec.Security.IsInverse {
			price = 1. / o.Price
			exp = -1
		} else {
			price = o.Price
			exp = 1
		}
		queue := o.Queue
		quantity := o.Quantity
		mul := sec.Multiplier.Value
		makerFee := sec.makerFee
		sampleValueChange := sec.sampleValueChange
		sampleMatchAsk := sec.sampleMatchAsk
		sampleSecurityPrice := sec.sampleSecurityPrice
		sampleMarginPrice := sec.sampleMarginPrice
		N := len(sampleValueChange)
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * mul
			fee := math.Abs(contractMarginValue) * makerFee

			// Remove cost from values
			cost := -contractMarginValue + fee + (contractChange * math.Pow(sampleSecurityPrice[i], exp) * sec.Multiplier.Value)
			sampleValueChange[i] += cost * sampleMarginPrice[i]
		}
		delete(sec.openAskOrders, ID)
	}
	sec.Unlock()
}

func (sec *MarginSecurity) UpdateFees(makerFee, takerFee float64) {
	sec.Lock()
	sec.makerFee = makerFee
	sec.takerFee = takerFee
	sec.Unlock()
}

func (sec *MarginSecurity) UpdateBidOrder(ID string, price float64, qty float64) {
	sec.Lock()
	// TODO update sample value change ?
	if o, ok := sec.openBidOrders[ID]; ok {
		o.Price = price
		o.Quantity = qty
	}
	sec.Unlock()
}

func (sec *MarginSecurity) UpdateAskOrder(ID string, price float64, qty float64) {
	sec.Lock()
	// TODO update sample value change
	if o, ok := sec.openAskOrders[ID]; ok {
		o.Price = price
		o.Quantity = qty
	}
	sec.Unlock()
}

func (sec *MarginSecurity) UpdateBidOrderPrice(ID string, price float64) {
	sec.Lock()
	// TODO update sample value change
	if o, ok := sec.openBidOrders[ID]; ok {
		o.Price = price
	}
	sec.Unlock()
}

func (sec *MarginSecurity) UpdateAskOrderPrice(ID string, price float64) {
	sec.Lock()
	// TODO update sample value change
	if o, ok := sec.openAskOrders[ID]; ok {
		o.Price = price
	}
	sec.Unlock()
}

func (sec *MarginSecurity) UpdateBidOrderQuantity(ID string, qty float64) {
	sec.Lock()
	// TODO update sample value change
	if o, ok := sec.openBidOrders[ID]; ok {
		o.Quantity = qty
	}
	sec.Unlock()
}

func (sec *MarginSecurity) UpdateAskOrderQuantity(ID string, qty float64) {
	sec.Lock()
	// TODO update sample value change
	if o, ok := sec.openAskOrders[ID]; ok {
		o.Quantity = qty
	}
	sec.Unlock()
}

func (sec *MarginSecurity) UpdateBidOrderQueue(ID string, queue float64) {
	sec.Lock()
	if o, ok := sec.openBidOrders[ID]; ok {
		var price, exp float64
		if sec.Security.IsInverse {
			price = 1. / o.Price
			exp = -1
		} else {
			price = o.Price
			exp = 1
		}
		oldQueue := o.Queue
		quantity := o.Quantity
		mul := sec.Multiplier.Value
		makerFee := sec.makerFee
		sampleValueChange := sec.sampleValueChange
		sampleMatchBid := sec.sampleMatchBid
		sampleSecurityPrice := sec.sampleSecurityPrice
		sampleMarginPrice := sec.sampleMarginPrice
		N := len(sampleValueChange)
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchBid[i]-oldQueue, 0.), quantity)
			contractMarginValue := contractChange * price * mul
			fee := math.Abs(contractMarginValue) * makerFee
			// Remove cost from values
			cost := contractMarginValue + fee - (contractChange * math.Pow(sampleSecurityPrice[i], exp) * sec.Multiplier.Value)
			sampleValueChange[i] += cost * sampleMarginPrice[i]

			contractChange = math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			contractMarginValue = contractChange * price * mul
			fee = math.Abs(contractMarginValue) * makerFee

			// Add cost to values
			cost = contractMarginValue + fee - (contractChange * math.Pow(sampleSecurityPrice[i], exp) * sec.Multiplier.Value)

			sampleValueChange[i] -= cost * sampleMarginPrice[i]
		}
		o.Queue = queue
	}
	sec.Unlock()
}

func (sec *MarginSecurity) UpdateAskOrderQueue(ID string, queue float64) {
	sec.Lock()
	if o, ok := sec.openAskOrders[ID]; ok {
		var price, exp float64
		if sec.Security.IsInverse {
			price = 1. / o.Price
			exp = -1
		} else {
			price = o.Price
			exp = 1
		}
		oldQueue := o.Queue
		quantity := o.Quantity
		mul := sec.Multiplier.Value
		makerFee := sec.makerFee
		sampleValueChange := sec.sampleValueChange
		sampleMatchAsk := sec.sampleMatchAsk
		sampleSecurityPrice := sec.sampleSecurityPrice
		sampleMarginPrice := sec.sampleMarginPrice
		N := len(sampleValueChange)
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchAsk[i]-oldQueue, 0.), quantity)
			contractMarginValue := contractChange * price * mul
			fee := math.Abs(contractMarginValue) * makerFee

			// Remove cost from values
			cost := -contractMarginValue + fee + (contractChange * math.Pow(sampleSecurityPrice[i], exp) * sec.Multiplier.Value)
			sampleValueChange[i] += cost * sampleMarginPrice[i]

			contractChange = math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			contractMarginValue = contractChange * price * mul
			fee = math.Abs(contractMarginValue) * makerFee

			// Remove cost from values
			cost = -contractMarginValue + fee + (contractChange * math.Pow(sampleSecurityPrice[i], exp) * sec.Multiplier.Value)
			sampleValueChange[i] -= cost * sampleMarginPrice[i]
		}
		o.Queue = queue
	}
	sec.Unlock()
}

func (sec *MarginSecurity) UpdatePositionSize(size float64) {
	sec.Lock()
	sec.size = size
	sec.Unlock()
}

func (sec *MarginSecurity) GetLotPrecision() float64 {
	return sec.lotPrecision
}

func (sec *MarginSecurity) Clear() {
	for k, _ := range sec.openBidOrders {
		sec.RemoveBidOrder(k)
	}
	for k, _ := range sec.openAskOrders {
		sec.RemoveAskOrder(k)
	}
	sec.size = 0.
}

func (sec *MarginSecurity) GetInstrument() *models.Instrument {
	return &models.Instrument{
		SecurityID: &types.UInt64Value{Value: sec.SecurityID},
		Exchange:   sec.Exchange,
		Symbol:     &types.StringValue{Value: sec.Symbol},
	}
}

// Security sample value will be either margin or base or quote
// you can then use margin currency to eval portfolio value

func (sec *MarginSecurity) updateSampleValueChange(model modeling.MarketModel, time uint64, sampleSize int) {
	N := sampleSize
	// TODO handle change in sample size, need to recompute too
	sampleValueChange := make([]float64, sampleSize, sampleSize)
	sampleMatchBid := model.GetSampleMatchBid(sec.SecurityID, time, N)
	sampleMatchAsk := model.GetSampleMatchAsk(sec.SecurityID, time, N)
	sampleMarginPrice := model.GetSamplePairPrices(sec.marginCurrency.ID, constants.DOLLAR.ID, time, N)
	sampleSecurityPrice := model.GetSamplePrices(sec.SecurityID, time, N)
	securityPrice := model.GetPrice(sec.SecurityID)

	mul := sec.Multiplier.Value
	makerFee := sec.makerFee

	if sec.IsInverse {
		for i := 0; i < sampleSize; i++ {
			sampleMatchBid[i] *= securityPrice
			sampleMatchAsk[i] *= securityPrice
		}
	}

	for _, o := range sec.openBidOrders {
		var price, exp float64
		if sec.Security.IsInverse {
			price = 1. / o.Price
			exp = -1
		} else {
			price = o.Price
			exp = 1
		}
		queue := o.Queue
		quantity := o.Quantity

		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * mul
			fee := math.Abs(contractMarginValue) * makerFee
			// Compute the cost
			cost := contractMarginValue + fee - (contractChange * math.Pow(sampleSecurityPrice[i], exp) * sec.Multiplier.Value)
			sampleValueChange[i] -= cost * sampleMarginPrice[i]
		}
	}
	for _, o := range sec.openAskOrders {
		var price, exp float64
		if sec.Security.IsInverse {
			price = 1. / o.Price
			exp = -1
		} else {
			price = o.Price
			exp = 1
		}
		queue := o.Queue
		quantity := o.Quantity
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * mul
			fee := math.Abs(contractMarginValue) * makerFee

			// Compute the cost
			cost := -contractMarginValue + fee + (contractChange * math.Pow(sampleSecurityPrice[i], exp) * sec.Multiplier.Value)
			sampleValueChange[i] -= cost * sampleMarginPrice[i]
		}
	}

	sec.sampleValueChange = sampleValueChange
	sec.sampleSecurityPrice = sampleSecurityPrice
	sec.sampleMarginPrice = sampleMarginPrice
	sec.sampleMatchBid = sampleMatchBid
	sec.sampleMatchAsk = sampleMatchAsk
	sec.sampleTime = time
}

func (sec *MarginSecurity) AddSampleValueChange(model modeling.MarketModel, time uint64, values []float64) {
	sec.Lock()
	// Update this instrument sample value change only if bid order or ask order set
	if len(sec.openBidOrders) > 0 || len(sec.openAskOrders) > 0 {
		N := len(values)
		if sec.sampleTime != time {
			sec.updateSampleValueChange(model, time, N)
		}
		sampleValueChange := sec.sampleValueChange
		for i := 0; i < N; i++ {
			values[i] += sampleValueChange[i]
		}
	}

	sec.Unlock()
}

func (sec *MarginSecurity) GetELROnCancelBid(ID string, model modeling.MarketModel, time uint64, values []float64, value float64) float64 {
	N := len(values)
	sec.Lock()

	if sec.sampleTime != time {
		sec.updateSampleValueChange(model, time, N)
	}

	sampleMatchBid := sec.sampleMatchBid
	sampleSecurityPrice := sec.sampleSecurityPrice
	sampleMarginPrice := sec.sampleMarginPrice

	// Remove bid order from value
	if o, ok := sec.openBidOrders[ID]; ok {
		var price, exp float64
		if sec.Security.IsInverse {
			price = 1. / o.Price
			exp = -1
		} else {
			price = o.Price
			exp = 1
		}
		queue := o.Queue
		quantity := o.Quantity
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * sec.Multiplier.Value
			fee := math.Abs(contractMarginValue) * sec.makerFee

			// Remove cost from values
			cost := contractMarginValue + fee - (contractChange * math.Pow(sampleSecurityPrice[i], exp) * sec.Multiplier.Value)
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

func (sec *MarginSecurity) GetELROnCancelAsk(ID string, model modeling.MarketModel, time uint64, values []float64, value float64) float64 {
	N := len(values)
	sec.Lock()

	if sec.sampleTime != time {
		sec.updateSampleValueChange(model, time, N)
	}

	sampleMatchAsk := sec.sampleMatchAsk
	sampleSecurityPrice := sec.sampleSecurityPrice
	sampleMarginPrice := sec.sampleMarginPrice

	// Remove ask order from value
	if o, ok := sec.openAskOrders[ID]; ok {
		var price, exp float64
		if sec.Security.IsInverse {
			price = 1. / o.Price
			exp = -1
		} else {
			price = o.Price
			exp = 1
		}
		queue := o.Queue
		quantity := o.Quantity
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * sec.Multiplier.Value
			fee := math.Abs(contractMarginValue) * sec.makerFee

			// Remove cost from values
			cost := -contractMarginValue + fee + (contractChange * math.Pow(sampleSecurityPrice[i], exp) * sec.Multiplier.Value)
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

func (sec *MarginSecurity) GetELROnLimitBidChange(ID string, model modeling.MarketModel, time uint64, values []float64, value float64, prices []float64, queues []float64, availableMargin float64) (float64, *COrder) {
	N := len(values)
	sec.Lock()
	// We want to see which option is the best, update, do nothing, or cancel
	if sec.sampleTime != time {
		sec.updateSampleValueChange(model, time, N)
	}

	sampleMatchBid := sec.sampleMatchBid
	sampleSecurityPrice := sec.sampleSecurityPrice
	sampleMarginPrice := sec.sampleMarginPrice
	mul := sec.Multiplier.Value
	makerFee := sec.makerFee

	var maxOrder *COrder = nil
	maxExpectedLogReturn := -999.

	// If we have the bid order, we remove it from values
	if o, ok := sec.openBidOrders[ID]; ok {
		var price, exp float64
		if sec.Security.IsInverse {
			price = 1. / o.Price
			exp = -1
		} else {
			price = o.Price
			exp = 1
		}
		queue := o.Queue
		quantity := o.Quantity
		for i := 0; i < N; i++ {
			// Compute ELR before removing order
			contractChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * mul
			fee := math.Abs(contractMarginValue) * makerFee

			// Remove cost from values
			cost := contractMarginValue + fee - (contractChange * math.Pow(sampleSecurityPrice[i], exp) * mul)
			values[i] += cost * sampleMarginPrice[i]
		}
	}

	// Compute ELR without any order
	noOrderELR := 0.
	for i := 0; i < N; i++ {
		noOrderELR += math.Log(values[i] / value)
	}
	noOrderELR /= float64(N)
	if noOrderELR > maxExpectedLogReturn {
		maxExpectedLogReturn = noOrderELR
		maxOrder = nil
	}

	for l := 0; l < len(prices); l++ {
		correctedAvailableMargin := availableMargin
		var price, exp float64
		if sec.IsInverse {
			price = 1. / prices[l]
			exp = -1.
		} else {
			price = prices[l]
			exp = 1.
		}
		if sec.size < 0 {
			// If we are short, we can go long the short size + the available margin size
			correctedAvailableMargin += math.Abs(price * mul * sec.size)
		}
		queue := queues[l]

		for q := 1; q < 5; q++ {
			expectedLogReturn := 0.
			availableContracts := math.Round((correctedAvailableMargin / float64(q)) / (math.Abs(mul) * price))
			for i := 0; i < N; i++ {
				contractChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), availableContracts)
				contractMarginValue := contractChange * price * mul
				fee := math.Abs(contractMarginValue) * makerFee
				cost := contractMarginValue + fee - (contractChange * math.Pow(sampleSecurityPrice[i], exp) * mul)
				// Compute ELR with cost
				expectedLogReturn += math.Log((values[i] - cost*sampleMarginPrice[i]) / value)
			}
			if expectedLogReturn > maxExpectedLogReturn {
				//fmt.Println("level ELR",  prices[l], expectedLogReturn)
				maxExpectedLogReturn = expectedLogReturn
				maxOrder = &COrder{
					Price:    prices[l],
					Quantity: availableContracts,
					Queue:    queue,
				}
			}
		}
	}

	sec.Unlock()
	return maxExpectedLogReturn, maxOrder
}

func (sec *MarginSecurity) GetELROnLimitAskChange(ID string, model modeling.MarketModel, time uint64, values []float64, value float64, prices []float64, queues []float64, availableMargin float64) (float64, *COrder) {
	N := len(values)
	sec.Lock()
	// We want to see which option is the best, update, do nothing, or cancel
	if sec.sampleTime != time {
		sec.updateSampleValueChange(model, time, N)
	}

	sampleMatchAsk := sec.sampleMatchAsk
	sampleSecurityPrice := sec.sampleSecurityPrice
	sampleMarginPrice := sec.sampleMarginPrice
	mul := sec.Multiplier.Value
	makerFee := sec.makerFee

	var maxOrder *COrder = nil
	maxExpectedLogReturn := -999.

	// Remove ask order from value
	if o, ok := sec.openAskOrders[ID]; ok {
		var price, exp float64
		if sec.Security.IsInverse {
			price = 1. / o.Price
			exp = -1.
		} else {
			price = o.Price
			exp = 1.
		}
		queue := o.Queue
		quantity := o.Quantity
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * mul
			fee := math.Abs(contractMarginValue) * makerFee

			// Remove the cost from values
			cost := -contractMarginValue + fee + (contractChange * math.Pow(sampleSecurityPrice[i], exp) * mul)
			values[i] += cost * sampleMarginPrice[i]
		}
	}

	// Compute ELR without any order
	noOrderELR := 0.
	for i := 0; i < N; i++ {
		noOrderELR += math.Log(values[i] / value)
	}
	noOrderELR /= float64(N)
	if noOrderELR > maxExpectedLogReturn {
		maxExpectedLogReturn = noOrderELR
		maxOrder = nil
	}

	for l := 0; l < len(prices); l++ {
		correctedAvailableMargin := availableMargin
		var price, exp float64
		if sec.IsInverse {
			price = 1. / prices[l]
			exp = -1
		} else {
			price = prices[l]
			exp = 1
		}
		if sec.size > 0 {
			// If we are long, we can go short the long size + the available margin size
			correctedAvailableMargin += math.Abs(price * mul * sec.size)
		}
		queue := queues[l]

		for q := 1; q < 5; q++ {
			expectedLogReturn := 0.
			availableContracts := (correctedAvailableMargin / float64(q)) / (math.Abs(mul) * price)
			// We need to compute what would happen if we changed
			// our ask order on the instrument.
			for i := 0; i < N; i++ {
				contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), availableContracts)
				contractMarginValue := contractChange * price * mul
				fee := math.Abs(contractMarginValue) * makerFee
				cost := -contractMarginValue + fee + (contractChange * math.Pow(sampleSecurityPrice[i], exp) * mul)
				// Compute ELR with new cost
				expectedLogReturn += math.Log((values[i] - cost*sampleMarginPrice[i]) / value)
			}

			expectedLogReturn /= float64(N)
			if expectedLogReturn > maxExpectedLogReturn {
				maxExpectedLogReturn = expectedLogReturn
				maxOrder = &COrder{
					Price:    prices[l],
					Quantity: (correctedAvailableMargin / float64(q)) / (price * math.Abs(mul)),
					Queue:    queue,
				}
			}
		}
	}

	sec.Unlock()
	return maxExpectedLogReturn, maxOrder
}

func (sec *MarginSecurity) GetELROnMarketBuy(model modeling.MarketModel, time uint64, values []float64, value float64, price, quantity float64, maxQuantity float64) (float64, *COrder) {
	N := len(values)
	sec.Lock()
	if sec.sampleTime != time {
		sec.updateSampleValueChange(model, time, N)
	}

	sampleSecurityPrice := sec.sampleSecurityPrice
	sampleMarginPrice := sec.sampleMarginPrice
	mul := sec.Multiplier.Value
	takerFee := sec.takerFee

	var maxOrder *COrder = nil
	correctedAvailableMargin := maxQuantity
	var exp, corrPrice float64
	if sec.IsInverse {
		corrPrice = 1. / price
		exp = -1.
	} else {
		corrPrice = price
		exp = 1.
	}
	if sec.size < 0 {
		// If we are short, we can go long the short size + the available margin size
		correctedAvailableMargin += math.Abs(corrPrice * mul * sec.size)
	}
	availableToMatch := math.Min(correctedAvailableMargin/(math.Abs(mul)*corrPrice), quantity)
	// Compute ELR without any order
	maxExpectedLogReturn := 0.
	for i := 0; i < N; i++ {
		maxExpectedLogReturn += math.Log(values[i] / value)
	}
	maxExpectedLogReturn /= float64(N)

	for q := 1; q < 5; q++ {
		expectedLogReturn := 0.
		contractChange := availableToMatch / float64(q)

		for i := 0; i < N; i++ {
			contractMarginValue := contractChange * corrPrice * mul
			fee := math.Abs(contractMarginValue) * takerFee
			cost := contractMarginValue + fee - (contractChange * math.Pow(sampleSecurityPrice[i], exp) * mul)
			// Compute ELR with cost
			expectedLogReturn += math.Log((values[i] - cost*sampleMarginPrice[i]) / value)
		}

		expectedLogReturn /= float64(N)
		if expectedLogReturn > maxExpectedLogReturn {
			//fmt.Println("level ELR",  prices[l], expectedLogReturn)
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = &COrder{
				Price:    price,
				Quantity: (correctedAvailableMargin / float64(q)) / (corrPrice * math.Abs(mul)),
				Queue:    0,
			}
		}
	}
	sec.Unlock()
	return maxExpectedLogReturn, maxOrder
}

func (sec *MarginSecurity) GetELROnMarketSell(model modeling.MarketModel, time uint64, values []float64, value float64, price float64, quantity float64, maxQuantity float64) (float64, *COrder) {
	N := len(values)
	sec.Lock()
	// We want to see which option is the best, update, do nothing, or cancel
	if sec.sampleTime != time {
		sec.updateSampleValueChange(model, time, N)
	}

	sampleSecurityPrice := sec.sampleSecurityPrice
	sampleMarginPrice := sec.sampleMarginPrice
	mul := sec.Multiplier.Value
	takerFee := sec.takerFee

	var maxOrder *COrder = nil
	correctedAvailableMargin := maxQuantity
	var exp, corrPrice float64
	if sec.IsInverse {
		corrPrice = 1. / price
		exp = -1.
	} else {
		corrPrice = price
		exp = 1.
	}
	if sec.size > 0 {
		// If we are long, we can go short the long size + the available margin size
		correctedAvailableMargin += math.Abs(corrPrice * mul * sec.size)
	}
	availableToMatch := math.Min(correctedAvailableMargin/(math.Abs(mul)*corrPrice), quantity)

	// Compute ELR without any order
	maxExpectedLogReturn := 0.
	for i := 0; i < N; i++ {
		maxExpectedLogReturn += math.Log(values[i] / value)
	}
	maxExpectedLogReturn /= float64(N)

	for q := 1; q < 5; q++ {
		expectedLogReturn := 0.
		contractChange := availableToMatch / float64(q)

		for i := 0; i < N; i++ {
			contractMarginValue := contractChange * corrPrice * mul
			fee := math.Abs(contractMarginValue) * takerFee
			cost := -contractMarginValue + fee + (contractChange * math.Pow(sampleSecurityPrice[i], exp) * mul)
			// Compute ELR with new cost
			expectedLogReturn += math.Log((values[i] - cost*sampleMarginPrice[i]) / value)
		}

		expectedLogReturn /= float64(N)
		if expectedLogReturn > maxExpectedLogReturn {
			//fmt.Println("level ELR",  prices[l], expectedLogReturn)
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = &COrder{
				Price:    price,
				Quantity: (correctedAvailableMargin / float64(q)) / (corrPrice * math.Abs(mul)),
				Queue:    0,
			}
		}
	}
	sec.Unlock()
	return maxExpectedLogReturn, maxOrder
}
