package account

import (
	"fmt"
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
	tickPrecision  float64
	lotPrecision   float64
	marginCurrency *xchangerModels.Asset
	openBidOrders  map[string]*COrder
	openAskOrders  map[string]*COrder
	size           float64

	sampleValueChange   []float64
	sampleMatchBid      []float64
	sampleMatchAsk      []float64
	sampleMarginPrice   []float64
	sampleSecurityPrice []float64
	sampleTime          uint64
}

func NewMarginSecurity(sec *models.Security, marginCurrency *xchangerModels.Asset) *MarginSecurity {
	return &MarginSecurity{
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
	makerFee := sec.MakerFee.Value
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
		cost := contractMarginValue + fee - (contractChange * math.Pow(sampleSecurityPrice[i], exp) * sec.Multiplier.Value)
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
	makerFee := sec.MakerFee.Value
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
	o := sec.openBidOrders[ID]
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
	makerFee := sec.MakerFee.Value
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
	sec.Unlock()
}

func (sec *MarginSecurity) RemoveAskOrder(ID string) {
	sec.Lock()
	o := sec.openAskOrders[ID]
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
	makerFee := sec.MakerFee.Value
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
	sec.Unlock()
}

func (sec *MarginSecurity) UpdateBidOrderQuantity(ID string, qty float64) {
	sec.Lock()
	// TODO update sample value change ?
	sec.openBidOrders[ID].Quantity = qty
	sec.Unlock()
}

func (sec *MarginSecurity) UpdateAskOrderQuantity(ID string, qty float64) {
	sec.Lock()
	// TODO update sample value change
	sec.openAskOrders[ID].Quantity = qty
	sec.Unlock()
}

func (sec *MarginSecurity) UpdateBidOrderQueue(ID string, queue float64) {
	sec.Lock()
	// TODO update sample value change
	sec.openBidOrders[ID].Queue = queue
	sec.Unlock()
}

func (sec *MarginSecurity) UpdateAskOrderQueue(ID string, queue float64) {
	sec.Lock()
	// TODO update sample value change
	sec.openAskOrders[ID].Queue = queue
	sec.Unlock()
}

func (sec *MarginSecurity) UpdatePositionSize(size float64) {
	sec.size = size
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
	makerFee := sec.MakerFee.Value

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
			fee := math.Abs(contractMarginValue) * sec.MakerFee.Value

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
			fee := math.Abs(contractMarginValue) * sec.MakerFee.Value

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
	makerFee := sec.MakerFee.Value

	var maxOrder, currentOrder *COrder
	maxExpectedLogReturn := -999.

	// If we have the bid order, we remove it from values
	if o, ok := sec.openBidOrders[ID]; ok {
		currentOrder = o
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
		expectedLogReturn := 0.
		for i := 0; i < N; i++ {
			// Compute ELR before removing order
			expectedLogReturn += math.Log(values[i] / value)
			contractChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * mul
			fee := math.Abs(contractMarginValue) * makerFee

			// Remove cost from values
			cost := contractMarginValue + fee - (contractChange * math.Pow(sampleSecurityPrice[i], exp) * sec.Multiplier.Value)

			values[i] += cost * sampleMarginPrice[i]
		}
		// Compute ELR with current bid
		expectedLogReturn /= float64(N)
		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = o
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
		expectedLogReturn := 0.

		for i := 0; i < N; i++ {
			availableContracts := correctedAvailableMargin / (math.Abs(sec.Multiplier.Value) * price)
			contractChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), availableContracts)
			contractMarginValue := contractChange * price * mul
			fee := math.Abs(contractMarginValue) * makerFee
			cost := contractMarginValue + fee - (contractChange * math.Pow(sampleSecurityPrice[i], exp) * sec.Multiplier.Value)
			// Compute ELR with cost
			expectedLogReturn += math.Log((values[i] - cost*sampleMarginPrice[i]) / value)
		}

		expectedLogReturn /= float64(N)
		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = &COrder{
				Price:    prices[l],
				Quantity: 0., // Computed at the end
				Queue:    queue,
			}
		}
	}

	if maxOrder != nil && maxOrder != currentOrder {
		// Compute recommended order quantity based on price and queue
		// Recommended order quantity is the expected match at the level
		expectedMatch := 0.
		for i := 0; i < N; i++ {
			expectedMatch += math.Max(sampleMatchBid[i]-maxOrder.Queue, 0.)
		}
		expectedMatch /= float64(N)

		correctedAvailableMargin := availableMargin
		var price float64
		if sec.IsInverse {
			price = 1. / maxOrder.Price
		} else {
			price = maxOrder.Price
		}
		if sec.size < 0 {
			// If we are short, we can go long the short size + the available margin size
			correctedAvailableMargin += math.Abs(price * mul * sec.size)
		}
		maxOrder.Quantity = math.Min(expectedMatch, correctedAvailableMargin/(price*math.Abs(mul)))
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
	makerFee := sec.MakerFee.Value

	var maxOrder, currentOrder *COrder
	maxExpectedLogReturn := -999.

	// Remove ask order from value
	if o, ok := sec.openAskOrders[ID]; ok {
		currentOrder = o
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
		expectedLogReturn := 0.
		for i := 0; i < N; i++ {
			// Compute ELR before removing order
			expectedLogReturn += math.Log(values[i] / value)
			contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * mul
			fee := math.Abs(contractMarginValue) * makerFee

			// Remove the cost from values
			cost := -contractMarginValue + fee + (contractChange * math.Pow(sampleSecurityPrice[i], exp) * mul)
			values[i] += cost * sampleMarginPrice[i]
		}
		// Compute ELR with current order
		expectedLogReturn /= float64(N)
		if expectedLogReturn > maxExpectedLogReturn {
			fmt.Println("current order elr", expectedLogReturn)
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = o
		}
	}

	// Compute ELR without any order
	expectedLogReturn := 0.
	for i := 0; i < N; i++ {
		expectedLogReturn += math.Log(values[i] / value)
	}
	expectedLogReturn /= float64(N)
	if expectedLogReturn > maxExpectedLogReturn {
		fmt.Println("naked elr", expectedLogReturn)
		maxExpectedLogReturn = expectedLogReturn
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
		expectedLogReturn := 0.

		// We need to compute what would happen if we changed
		// our ask order on the instrument.
		for i := 0; i < N; i++ {
			availableContracts := correctedAvailableMargin / (math.Abs(mul) * price)
			contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), availableContracts)
			contractMarginValue := contractChange * price * mul
			fee := math.Abs(contractMarginValue) * makerFee
			cost := -contractMarginValue + fee + (contractChange * math.Pow(sampleSecurityPrice[i], exp) * mul)
			// Compute ELR with new cost
			expectedLogReturn += math.Log((values[i] - cost*sampleMarginPrice[i]) / value)
		}

		expectedLogReturn /= float64(N)

		//fmt.Println(expectedLogReturn, maxExpectedLogReturn)
		if expectedLogReturn > maxExpectedLogReturn {
			fmt.Println("level elr", prices[l], expectedLogReturn)
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = &COrder{
				Price:    prices[l],
				Quantity: 0., // Computed at the end
				Queue:    queue,
			}
		}
	}

	if maxOrder != nil && maxOrder != currentOrder {
		// Compute recommended order quantity based on price and queue
		// Recommended order quantity is the expected match at the level
		//fmt.Println("RETS", baseline, maxExpectedLogReturn)
		correctedAvailableMargin := availableMargin
		expectedMatch := 0.
		for i := 0; i < N; i++ {
			expectedMatch += math.Max(sampleMatchAsk[i]-maxOrder.Queue, 0.)
		}
		expectedMatch /= float64(N)
		var price float64
		if sec.IsInverse {
			price = 1. / maxOrder.Price
		} else {
			price = maxOrder.Price
		}
		if sec.size > 0 {
			// If we are long, we can go short the long size + the available margin size
			correctedAvailableMargin += math.Abs(price * mul * sec.size)
		}
		maxOrder.Quantity = math.Min(expectedMatch, correctedAvailableMargin/(price*math.Abs(mul)))
	}
	sec.Unlock()

	return maxExpectedLogReturn, maxOrder
}
