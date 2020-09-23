package account

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alphac/modeling"
	"gitlab.com/alphaticks/alphac/models"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"math"
	"sync"
)

type SpotSecurity struct {
	sync.RWMutex
	*models.Security
	lotPrecision  float64
	openBidOrders map[string]*COrder
	openAskOrders map[string]*COrder

	modelCurrency     *xchangerModels.Asset
	sampleValueChange []float64
	sampleBasePrice   []float64
	sampleQuotePrice  []float64
	sampleMatchBid    []float64
	sampleMatchAsk    []float64
	sampleTime        uint64
	makerFee          float64
	takerFee          float64
}

func NewSpotSecurity(sec *models.Security, modelCurrency *xchangerModels.Asset, makerFee, takerFee *float64) (*SpotSecurity, error) {
	if sec.RoundLot == nil {
		return nil, fmt.Errorf("security is missing RoundLot")
	}
	spotSec := &SpotSecurity{
		RWMutex:           sync.RWMutex{},
		Security:          sec,
		lotPrecision:      math.Ceil(1. / sec.RoundLot.Value),
		openBidOrders:     make(map[string]*COrder),
		openAskOrders:     make(map[string]*COrder),
		modelCurrency:     modelCurrency,
		sampleValueChange: nil,
		sampleBasePrice:   nil,
		sampleQuotePrice:  nil,
		sampleMatchBid:    nil,
		sampleMatchAsk:    nil,
		sampleTime:        0,
	}

	if makerFee != nil {
		spotSec.makerFee = *makerFee
	} else {
		spotSec.makerFee = sec.MakerFee.Value
	}
	if takerFee != nil {
		spotSec.takerFee = *takerFee
	} else {
		spotSec.takerFee = sec.TakerFee.Value
	}
	return spotSec, nil
}

func (sec *SpotSecurity) AddBidOrder(ID string, price, quantity, queue float64) {
	sec.Lock()
	sec.openBidOrders[ID] = &COrder{
		Price:    price,
		Quantity: quantity,
		Queue:    queue,
	}
	sampleMatchBid := sec.sampleMatchBid
	sampleBasePrice := sec.sampleBasePrice
	sampleQuotePrice := sec.sampleQuotePrice
	sampleValueChange := sec.sampleValueChange
	N := len(sampleValueChange)
	for i := 0; i < N; i++ {
		baseChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
		// We get our reduced-by-fees base
		sampleValueChange[i] += (baseChange - (baseChange * sec.makerFee)) * sampleBasePrice[i]
		// We lose our quote
		sampleValueChange[i] -= baseChange * price * sampleQuotePrice[i]
	}
	sec.Unlock()
}

func (sec *SpotSecurity) AddAskOrder(ID string, price, quantity, queue float64) {
	sec.Lock()
	sec.openAskOrders[ID] = &COrder{
		Price:    price,
		Quantity: quantity,
		Queue:    queue,
	}
	sampleMatchAsk := sec.sampleMatchAsk
	sampleBasePrice := sec.sampleBasePrice
	sampleQuotePrice := sec.sampleQuotePrice
	sampleValueChange := sec.sampleValueChange
	N := len(sampleValueChange)
	for i := 0; i < N; i++ {
		baseChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
		// We lose our base
		sampleValueChange[i] -= baseChange * sampleBasePrice[i]
		// We get our reduced-by-fees quote
		sampleValueChange[i] += (baseChange - (baseChange * sec.makerFee)) * price * sampleQuotePrice[i]
	}
	sec.Unlock()
}

func (sec *SpotSecurity) RemoveBidOrder(ID string) {
	sec.Lock()
	o := sec.openBidOrders[ID]
	queue := o.Queue
	quantity := o.Quantity
	price := o.Price
	sampleMatchBid := sec.sampleMatchBid
	sampleBasePrice := sec.sampleBasePrice
	sampleQuotePrice := sec.sampleQuotePrice
	sampleValueChange := sec.sampleValueChange
	N := len(sampleValueChange)
	for i := 0; i < N; i++ {
		baseChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
		// We don't get our reduced-by-fees base
		sampleValueChange[i] -= (baseChange - (baseChange * sec.makerFee)) * sampleBasePrice[i]
		// We keep our quote
		sampleValueChange[i] += baseChange * price * sampleQuotePrice[i]
	}
	delete(sec.openBidOrders, ID)
	sec.Unlock()
}

func (sec *SpotSecurity) RemoveAskOrder(ID string) {
	sec.Lock()
	o := sec.openAskOrders[ID]
	queue := o.Queue
	quantity := o.Quantity
	price := o.Price
	sampleMatchAsk := sec.sampleMatchAsk
	sampleBasePrice := sec.sampleBasePrice
	sampleQuotePrice := sec.sampleQuotePrice
	sampleValueChange := sec.sampleValueChange
	N := len(sampleValueChange)
	for i := 0; i < N; i++ {
		baseChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
		// We keep our base
		sampleValueChange[i] += baseChange * sampleBasePrice[i]
		// We don't get our reduced-by-fees quote
		sampleValueChange[i] -= (baseChange - (baseChange * sec.makerFee)) * price * sampleQuotePrice[i]
	}
	delete(sec.openAskOrders, ID)
	sec.Unlock()
}

func (sec *SpotSecurity) UpdateFees(makerFee, takerFee float64) {
	sec.Lock()
	sec.makerFee = makerFee
	sec.takerFee = takerFee
	sec.Unlock()
}

func (sec *SpotSecurity) UpdateBidOrder(ID string, price, qty float64) {
	sec.Lock()
	sec.openBidOrders[ID].Quantity = qty
	sec.openBidOrders[ID].Price = price
	sec.Unlock()
}

func (sec *SpotSecurity) UpdateAskOrder(ID string, price, qty float64) {
	sec.Lock()
	sec.openAskOrders[ID].Quantity = qty
	sec.openAskOrders[ID].Price = price
	sec.Unlock()
}

func (sec *SpotSecurity) UpdateBidOrderPrice(ID string, price float64) {
	sec.Lock()
	o := sec.openBidOrders[ID]
	queue := o.Queue
	quantity := o.Quantity
	oldPrice := o.Price
	sampleMatchBid := sec.sampleMatchBid
	sampleQuotePrice := sec.sampleQuotePrice
	sampleValueChange := sec.sampleValueChange
	N := len(sampleValueChange)
	for i := 0; i < N; i++ {
		baseChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
		sampleValueChange[i] += baseChange * oldPrice * sampleQuotePrice[i]
		sampleValueChange[i] -= baseChange * price * sampleQuotePrice[i]
	}
	o.Price = price
	sec.Unlock()
}

func (sec *SpotSecurity) UpdateAskOrderPrice(ID string, price float64) {
	sec.Lock()
	o := sec.openAskOrders[ID]
	queue := o.Queue
	quantity := o.Quantity
	oldPrice := o.Price
	sampleMatchAsk := sec.sampleMatchAsk
	sampleQuotePrice := sec.sampleQuotePrice
	sampleValueChange := sec.sampleValueChange
	N := len(sampleValueChange)
	for i := 0; i < N; i++ {
		baseChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
		sampleValueChange[i] -= (baseChange - (baseChange * sec.makerFee)) * oldPrice * sampleQuotePrice[i]
		sampleValueChange[i] += (baseChange - (baseChange * sec.makerFee)) * price * sampleQuotePrice[i]
	}
	o.Price = price
	sec.Unlock()
}

func (sec *SpotSecurity) UpdateBidOrderQuantity(ID string, quantity float64) {
	sec.Lock()
	o := sec.openBidOrders[ID]
	queue := o.Queue
	oldQuantity := o.Quantity
	price := o.Price
	sampleMatchBid := sec.sampleMatchBid
	sampleBasePrice := sec.sampleBasePrice
	sampleQuotePrice := sec.sampleQuotePrice
	sampleValueChange := sec.sampleValueChange
	N := len(sampleValueChange)
	for i := 0; i < N; i++ {
		baseChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), oldQuantity)
		// We don't get our reduced-by-fees base
		sampleValueChange[i] -= (baseChange - (baseChange * sec.makerFee)) * sampleBasePrice[i]
		// We keep our quote
		sampleValueChange[i] += baseChange * price * sampleQuotePrice[i]

		baseChange = math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
		sampleValueChange[i] += (baseChange - (baseChange * sec.makerFee)) * sampleBasePrice[i]
		sampleValueChange[i] -= baseChange * price * sampleQuotePrice[i]
	}
	o.Quantity = quantity
	sec.Unlock()
}

func (sec *SpotSecurity) UpdateAskOrderQuantity(ID string, quantity float64) {
	sec.Lock()
	o := sec.openAskOrders[ID]
	queue := o.Queue
	oldQuantity := o.Quantity
	price := o.Price
	sampleMatchAsk := sec.sampleMatchAsk
	sampleBasePrice := sec.sampleBasePrice
	sampleQuotePrice := sec.sampleQuotePrice
	sampleValueChange := sec.sampleValueChange
	N := len(sampleValueChange)
	for i := 0; i < N; i++ {
		baseChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), oldQuantity)
		// We keep our base
		sampleValueChange[i] += baseChange * sampleBasePrice[i]
		// We don't get our reduced-by-fees quote
		sampleValueChange[i] -= (baseChange - (baseChange * sec.makerFee)) * price * sampleQuotePrice[i]

		baseChange = math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
		sampleValueChange[i] -= baseChange * sampleBasePrice[i]
		sampleValueChange[i] += (baseChange - (baseChange * sec.makerFee)) * price * sampleQuotePrice[i]
	}
	o.Quantity = quantity
	sec.Unlock()
}

func (sec *SpotSecurity) UpdateBidOrderQueue(ID string, queue float64) {
	sec.Lock()
	o := sec.openBidOrders[ID]
	oldQueue := o.Queue
	quantity := o.Quantity
	price := o.Price
	sampleMatchBid := sec.sampleMatchBid
	sampleBasePrice := sec.sampleBasePrice
	sampleQuotePrice := sec.sampleQuotePrice
	sampleValueChange := sec.sampleValueChange
	N := len(sampleValueChange)
	for i := 0; i < N; i++ {
		baseChange := math.Min(math.Max(sampleMatchBid[i]-oldQueue, 0.), quantity)
		// We don't get our reduced-by-fees base
		sampleValueChange[i] -= (baseChange - (baseChange * sec.makerFee)) * sampleBasePrice[i]
		// We keep our quote
		sampleValueChange[i] += baseChange * price * sampleQuotePrice[i]

		baseChange = math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
		sampleValueChange[i] += (baseChange - (baseChange * sec.makerFee)) * sampleBasePrice[i]
		sampleValueChange[i] -= baseChange * price * sampleQuotePrice[i]
	}
	o.Queue = queue
	sec.Unlock()
}

func (sec *SpotSecurity) UpdateAskOrderQueue(ID string, queue float64) {
	sec.Lock()
	o := sec.openAskOrders[ID]
	oldQueue := o.Queue
	quantity := o.Quantity
	price := o.Price
	sampleMatchAsk := sec.sampleMatchAsk
	sampleBasePrice := sec.sampleBasePrice
	sampleQuotePrice := sec.sampleQuotePrice
	sampleValueChange := sec.sampleValueChange
	N := len(sampleValueChange)
	for i := 0; i < N; i++ {
		baseChange := math.Min(math.Max(sampleMatchAsk[i]-oldQueue, 0.), quantity)
		// We keep our base
		sampleValueChange[i] += baseChange * sampleBasePrice[i]
		// We don't get our reduced-by-fees quote
		sampleValueChange[i] -= (baseChange - (baseChange * sec.makerFee)) * price * sampleQuotePrice[i]

		baseChange = math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
		sampleValueChange[i] -= baseChange * sampleBasePrice[i]
		sampleValueChange[i] += (baseChange - (baseChange * sec.makerFee)) * price * sampleQuotePrice[i]
	}
	o.Queue = queue
	sec.Unlock()
}

func (sec *SpotSecurity) UpdatePositionSize(_ float64) {}

func (sec *SpotSecurity) GetPosition() *models.Position {
	return nil
}

func (sec *SpotSecurity) GetLotPrecision() float64 {
	return sec.lotPrecision
}

func (sec *SpotSecurity) Clear() {
	for k, _ := range sec.openBidOrders {
		sec.RemoveBidOrder(k)
	}
	for k, _ := range sec.openAskOrders {
		sec.RemoveAskOrder(k)
	}
}

func (sec *SpotSecurity) GetInstrument() *models.Instrument {
	return &models.Instrument{
		SecurityID: &types.UInt64Value{Value: sec.SecurityID},
		Exchange:   sec.Exchange,
		Symbol:     &types.StringValue{Value: sec.Symbol},
	}
}

func (sec *SpotSecurity) updateSampleValueChange(model modeling.MarketModel, time uint64, sampleSize int) {
	// TODO refresh when sample size changes
	N := sampleSize
	sampleValueChange := make([]float64, sampleSize, sampleSize)
	sampleMatchBid := model.GetSampleMatchBid(sec.SecurityID, time, N)
	sampleMatchAsk := model.GetSampleMatchAsk(sec.SecurityID, time, N)
	sampleBasePrice := model.GetSamplePairPrices(sec.Underlying.ID, sec.modelCurrency.ID, time, N)
	sampleQuotePrice := model.GetSamplePairPrices(sec.QuoteCurrency.ID, sec.modelCurrency.ID, time, N)
	for i := 0; i < N; i++ {
		sampleValueChange[i] = 0.
	}
	for _, o := range sec.openBidOrders {
		queue := o.Queue
		quantity := o.Quantity
		price := o.Price
		for i := 0; i < N; i++ {
			baseChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			// We get our reduced-by-fees base
			sampleValueChange[i] += (baseChange - (baseChange * sec.makerFee)) * sampleBasePrice[i]
			// We lose our quote
			sampleValueChange[i] -= baseChange * price * sampleQuotePrice[i]
		}
	}
	for _, o := range sec.openAskOrders {
		queue := o.Queue
		quantity := o.Quantity
		price := o.Price
		for i := 0; i < N; i++ {
			baseChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			// We lose our base
			sampleValueChange[i] -= baseChange * sampleBasePrice[i]
			// We get our reduced-by-fees quote
			sampleValueChange[i] += (baseChange - (baseChange * sec.makerFee)) * price * sampleQuotePrice[i]
		}
	}
	sec.sampleValueChange = sampleValueChange
	sec.sampleBasePrice = sampleBasePrice
	sec.sampleQuotePrice = sampleQuotePrice
	sec.sampleMatchBid = sampleMatchBid
	sec.sampleMatchAsk = sampleMatchAsk
	sec.sampleTime = time
}

func (sec *SpotSecurity) AddSampleValueChange(model modeling.MarketModel, time uint64, values []float64) {
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

func (sec *SpotSecurity) GetELROnCancelBid(ID string, model modeling.MarketModel, time uint64, values []float64, value float64) float64 {
	N := len(values)
	sec.Lock()

	if sec.sampleTime != time {
		sec.updateSampleValueChange(model, time, N)
	}

	sampleMatchBid := sec.sampleMatchBid
	sampleBasePrice := sec.sampleBasePrice
	sampleQuotePrice := sec.sampleQuotePrice

	// Remove bid order from value
	if o, ok := sec.openBidOrders[ID]; ok {
		queue := o.Queue
		quantity := o.Quantity
		price := o.Price
		for i := 0; i < N; i++ {
			baseChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			// We don't get our reduced-by-fees base
			values[i] -= (baseChange - (baseChange * sec.makerFee)) * sampleBasePrice[i]
			// We keep our quote
			values[i] += baseChange * price * sampleQuotePrice[i]
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

func (sec *SpotSecurity) GetELROnCancelAsk(ID string, model modeling.MarketModel, time uint64, values []float64, value float64) float64 {
	N := len(values)
	sec.Lock()

	if sec.sampleTime != time {
		sec.updateSampleValueChange(model, time, N)
	}

	sampleMatchAsk := sec.sampleMatchAsk
	sampleBasePrice := sec.sampleBasePrice
	sampleQuotePrice := sec.sampleQuotePrice

	// Remove ask order from value
	if o, ok := sec.openAskOrders[ID]; ok {
		queue := o.Queue
		quantity := o.Quantity
		price := o.Price
		for i := 0; i < N; i++ {
			baseChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			// We keep our base
			values[i] += baseChange * sampleBasePrice[i]
			// We don't get our reduced-by-fees quote
			values[i] -= (baseChange - (baseChange * sec.makerFee)) * price * sampleQuotePrice[i]
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

func (sec *SpotSecurity) GetELROnLimitBidChange(ID string, model modeling.MarketModel, time uint64, values []float64, value float64, prices []float64, queues []float64, maxQuote float64) (float64, *COrder) {
	N := len(values)
	sec.Lock()
	// We want to see which option is the best, update, do nothing, or cancel
	if sec.sampleTime != time {
		sec.updateSampleValueChange(model, time, N)
	}

	sampleMatchBid := sec.sampleMatchBid
	sampleBasePrice := sec.sampleBasePrice
	sampleQuotePrice := sec.sampleQuotePrice

	var currentOrder *COrder
	var maxOrder *COrder
	maxExpectedLogReturn := -999.
	// If we have an a bid order, we remove it from values
	if o, ok := sec.openBidOrders[ID]; ok {
		currentOrder = o
		queue := o.Queue
		quantity := o.Quantity
		price := o.Price
		expectedLogReturn := 0.
		for i := 0; i < N; i++ {
			expectedLogReturn += math.Log(values[i] / value)
			baseChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			// We keep our base
			values[i] += baseChange * sampleBasePrice[i]
			// We don't get our reduced-by-fees quote
			values[i] -= (baseChange - (baseChange * sec.makerFee)) * price * sampleQuotePrice[i]
		}
		expectedLogReturn /= float64(N)
		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = o
		}
	}

	/*
		// Compute expected log return without any order
		expectedLogReturn := 0.
		for i := 0; i < N; i++ {
			expectedLogReturn += math.Log(values[i] / value)
		}
		expectedLogReturn /= float64(N)
		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = nil
		}
	*/

	for l := 0; l < len(prices); l++ {
		price := prices[l]
		queue := queues[l]
		expectedLogReturn := 0.

		for i := 0; i < N; i++ {
			expectedMatchQuote := math.Min(math.Max(sampleMatchBid[i]-queue, 0)*price, maxQuote)
			expectedMatchBase := expectedMatchQuote / price
			valueChange := (expectedMatchBase*(1-sec.makerFee))*sampleBasePrice[i] - expectedMatchQuote*sampleQuotePrice[i]
			expectedLogReturn += math.Log((values[i] + valueChange) / value)
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

func (sec *SpotSecurity) GetELROnLimitAskChange(ID string, model modeling.MarketModel, time uint64, values []float64, value float64, prices []float64, queues []float64, maxBase float64) (float64, *COrder) {
	N := len(values)
	sec.Lock()
	// We want to see which option is the best, update, do nothing, or cancel
	if sec.sampleTime != time {
		sec.updateSampleValueChange(model, time, N)
	}

	sampleMatchAsk := sec.sampleMatchAsk
	sampleBasePrice := sec.sampleBasePrice
	sampleQuotePrice := sec.sampleQuotePrice

	var maxOrder, currentOrder *COrder
	maxExpectedLogReturn := -999.

	// Remove ask order from value
	if o, ok := sec.openAskOrders[ID]; ok {
		currentOrder = o
		queue := o.Queue
		quantity := o.Quantity
		price := o.Price
		expectedLogReturn := 0.
		for i := 0; i < N; i++ {
			expectedLogReturn += math.Log(values[i] / value)
			baseChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			// We keep our base
			values[i] += baseChange * sampleBasePrice[i]
			// We don't get our reduced-by-fees quote
			values[i] -= (baseChange - (baseChange * sec.makerFee)) * price * sampleQuotePrice[i]
		}
		expectedLogReturn /= float64(N)
		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = o
		}
	}

	/*
		// Compute expected log return without any order
		expectedLogReturn := 0.
		for i := 0; i < N; i++ {
			expectedLogReturn += math.Log(values[i] / value)
		}
		expectedLogReturn /= float64(N)
		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = nil
		}
	*/

	for l := 0; l < len(prices); l++ {
		price := prices[l]
		queue := queues[l]
		expectedLogReturn := 0.

		// We need to compute what would happen if we changed
		// our ask order on the instrument. A change in the ask order
		// means a change in the sampleAsset
		for i := 0; i < N; i++ {
			// Now for the new order
			expectedMatch := math.Min(math.Max(sampleMatchAsk[i]-queue, 0), maxBase)
			// We are expected to lose expectedMatch base
			baseValueChange := -expectedMatch * sampleBasePrice[i]
			// We are expected to gain expectedMatch * price quote
			quoteValueChange := (expectedMatch - (sec.makerFee * expectedMatch)) * price * sampleQuotePrice[i]

			expectedLogReturn += math.Log((values[i] + baseValueChange + quoteValueChange) / value)
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

func (sec *SpotSecurity) GetELROnMarketBuy(model modeling.MarketModel, time uint64, values []float64, value float64, price, quantity float64, maxQuantity float64) (float64, *COrder) {
	return 0, nil
}

func (sec *SpotSecurity) GetELROnMarketSell(model modeling.MarketModel, time uint64, values []float64, value float64, price float64, quantity float64, maxQuantity float64) (float64, *COrder) {
	return 0, nil
}
