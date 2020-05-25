package account

import (
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alphac/modeling"
	"gitlab.com/alphaticks/alphac/models"
	"math"
	"sync"
)

type COrder struct {
	Price    float64
	Quantity float64
	Queue    float64
}

type MarginSecurity struct {
	sync.RWMutex
	*models.Security
	tickPrecision   float64
	lotPrecision    float64
	marginPrecision float64
	openBidOrders   map[string]*COrder
	openAskOrders   map[string]*COrder

	marginPriceModel    modeling.PriceModel
	securityPriceModel  modeling.PriceModel
	buyTradeModel       modeling.BuyTradeModel
	sellTradeModel      modeling.SellTradeModel
	sampleValueChange   []float64
	sampleMatchBid      []float64
	sampleMatchAsk      []float64
	sampleMarginPrice   []float64
	sampleSecurityPrice []float64
	sampleTime          uint64
}

func NewMarginSecurity(sec *models.Security, marginPrecision float64, marginPriceModel, securityPriceModel modeling.PriceModel, buyTradeModel modeling.BuyTradeModel, sellTradeModel modeling.SellTradeModel, sampleSize int) *MarginSecurity {
	return &MarginSecurity{
		RWMutex:             sync.RWMutex{},
		Security:            sec,
		tickPrecision:       math.Ceil(1. / sec.MinPriceIncrement),
		lotPrecision:        math.Ceil(1. / sec.RoundLot),
		marginPrecision:     marginPrecision,
		openBidOrders:       make(map[string]*COrder),
		openAskOrders:       make(map[string]*COrder),
		marginPriceModel:    marginPriceModel,
		securityPriceModel:  securityPriceModel,
		buyTradeModel:       buyTradeModel,
		sellTradeModel:      sellTradeModel,
		sampleValueChange:   make([]float64, sampleSize, sampleSize),
		sampleMatchBid:      make([]float64, sampleSize, sampleSize),
		sampleMatchAsk:      make([]float64, sampleSize, sampleSize),
		sampleMarginPrice:   make([]float64, sampleSize, sampleSize),
		sampleSecurityPrice: make([]float64, sampleSize, sampleSize),
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
	sec.Unlock()
}

func (sec *MarginSecurity) AddAskOrder(ID string, price, quantity, queue float64) {
	sec.Lock()
	sec.openAskOrders[ID] = &COrder{
		Price:    price,
		Quantity: quantity,
		Queue:    queue,
	}
	sec.Unlock()
}

func (sec *MarginSecurity) RemoveBidOrder(ID string) {
	sec.Lock()
	delete(sec.openBidOrders, ID)
	sec.Unlock()
}

func (sec *MarginSecurity) RemoveAskOrder(ID string) {
	sec.Lock()
	delete(sec.openAskOrders, ID)
	sec.Unlock()
}

func (sec *MarginSecurity) UpdateBidOrderQuantity(ID string, qty float64) {
	sec.Lock()
	// TODO update sample value change
	sec.openBidOrders[ID].Quantity = qty
	sec.Unlock()
}

func (sec *MarginSecurity) UpdateAskOrderQuantity(ID string, qty float64) {
	sec.Lock()
	// TODO update sample value change
	sec.openAskOrders[ID].Quantity = qty
	sec.Unlock()
}

func (sec *MarginSecurity) GetLotPrecision() float64 {
	return sec.lotPrecision
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

func (sec *MarginSecurity) updateSampleValueChange(time uint64) {
	sampleValueChange := sec.sampleValueChange
	N := len(sampleValueChange)
	sampleMatchBid := sec.sellTradeModel.GetSampleMatchBid(time, N)
	sampleMatchAsk := sec.buyTradeModel.GetSampleMatchAsk(time, N)
	sampleMarginPrice := sec.marginPriceModel.GetSamplePrices(time, N)
	sampleSecurityPrice := sec.securityPriceModel.GetSamplePrices(time, N)

	for _, o := range sec.openBidOrders {
		var price float64
		if sec.Security.IsInverse {
			price = 1. / o.Price
		} else {
			price = o.Price
		}
		queue := o.Queue
		quantity := o.Quantity
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * sec.Multiplier.Value
			fee := math.Abs(contractMarginValue) * sec.MakerFee.Value

			// Compute the cost
			cost := contractMarginValue + fee - (contractChange * sampleSecurityPrice[i] * sec.Multiplier.Value)
			sampleValueChange[i] -= cost * sampleMarginPrice[i]
		}
	}
	for _, o := range sec.openAskOrders {
		var price float64
		if sec.Security.IsInverse {
			price = 1. / o.Price
		} else {
			price = o.Price
		}
		queue := o.Queue
		quantity := o.Quantity
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * sec.Multiplier.Value
			fee := math.Abs(contractMarginValue) * sec.MakerFee.Value

			// Compute the cost
			cost := -contractMarginValue + fee + (contractChange * sampleSecurityPrice[i] * sec.Multiplier.Value)
			sampleValueChange[i] -= cost * sampleMarginPrice[i]
		}
	}

	sec.sampleSecurityPrice = sampleSecurityPrice
	sec.sampleMarginPrice = sampleMarginPrice
	sec.sampleMatchBid = sampleMatchBid
	sec.sampleMatchAsk = sampleMatchAsk
	sec.sampleTime = time
}

func (sec *MarginSecurity) AddSampleValueChange(time uint64, values []float64) {
	sec.Lock()
	// Update this instrument sample value change only if bid order or ask order set
	if len(sec.openBidOrders) > 0 || len(sec.openAskOrders) > 0 {
		if sec.sampleTime != time {
			sec.updateSampleValueChange(time)
		}
		N := len(values)
		sampleValueChange := sec.sampleValueChange
		for i := 0; i < N; i++ {
			values[i] += sampleValueChange[i]
		}
	}

	sec.Unlock()
}

func (sec *MarginSecurity) GetELROnCancelBid(ID string, time uint64, values []float64, value float64) float64 {
	N := len(values)
	sec.Lock()

	if sec.sampleTime != time {
		sec.updateSampleValueChange(time)
	}

	sampleMatchBid := sec.sampleMatchBid
	sampleSecurityPrice := sec.sampleSecurityPrice
	sampleMarginPrice := sec.sampleMarginPrice

	// Remove bid order from value
	if o, ok := sec.openBidOrders[ID]; ok {
		var price float64
		if sec.Security.IsInverse {
			price = 1. / o.Price
		} else {
			price = o.Price
		}
		queue := o.Queue
		quantity := o.Quantity
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

func (sec *MarginSecurity) GetELROnCancelAsk(ID string, time uint64, values []float64, value float64) float64 {
	N := len(values)
	sec.Lock()

	if sec.sampleTime != time {
		sec.updateSampleValueChange(time)
	}

	sampleMatchAsk := sec.sampleMatchAsk
	sampleSecurityPrice := sec.sampleSecurityPrice
	sampleMarginPrice := sec.sampleMarginPrice

	// Remove ask order from value
	if o, ok := sec.openAskOrders[ID]; ok {
		var price float64
		if sec.Security.IsInverse {
			price = 1. / o.Price
		} else {
			price = o.Price
		}
		queue := o.Queue
		quantity := o.Quantity
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

func (sec *MarginSecurity) GetELROnBidChange(ID string, time uint64, values []float64, value float64, prices []float64, queues []float64, maxQuote float64) (float64, *COrder) {
	N := len(values)
	sec.Lock()
	// We want to see which option is the best, update, do nothing, or cancel
	if sec.sampleTime != time {
		sec.updateSampleValueChange(time)
	}

	sampleMatchBid := sec.sampleMatchBid
	sampleSecurityPrice := sec.sampleSecurityPrice
	sampleMarginPrice := sec.sampleMarginPrice

	var maxOrder, currentOrder *COrder
	maxExpectedLogReturn := -999.

	// If we have an a bid order, we remove it from values
	if o, ok := sec.openBidOrders[ID]; ok {
		currentOrder = o
		var price float64
		if sec.Security.IsInverse {
			price = 1. / o.Price
		} else {
			price = o.Price
		}
		expectedLogReturn := 0.
		queue := o.Queue
		quantity := o.Quantity
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
			maxOrder = &COrder{
				Price:    price,
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
		maxOrder.Quantity = math.Min(expectedMatch, maxQuote/maxOrder.Price)
	}
	sec.Unlock()

	return maxExpectedLogReturn, maxOrder
}

func (sec *MarginSecurity) GetELROnAskChange(ID string, time uint64, values []float64, value float64, prices []float64, queues []float64, maxBase float64) (float64, *COrder) {
	N := len(values)
	sec.Lock()
	// We want to see which option is the best, update, do nothing, or cancel
	if sec.sampleTime != time {
		sec.updateSampleValueChange(time)
	}

	sampleMatchAsk := sec.sampleMatchAsk
	sampleSecurityPrice := sec.sampleSecurityPrice
	sampleMarginPrice := sec.sampleMarginPrice

	var maxOrder, currentOrder *COrder
	maxExpectedLogReturn := -999.

	// Remove ask order from value
	if o, ok := sec.openAskOrders[ID]; ok {
		currentOrder = o
		var price float64
		if sec.Security.IsInverse {
			price = 1. / o.Price
		} else {
			price = o.Price
		}
		queue := o.Queue
		quantity := o.Quantity
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
			maxOrder = &COrder{
				Price:    price,
				Quantity: 0., // Computed at the end
				Queue:    queue,
			}
		}
	}

	if maxOrder != nil && maxOrder != currentOrder {
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
