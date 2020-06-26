package legacy

import (
	"fmt"
	"gitlab.com/alphaticks/alphac/modeling"
	"gitlab.com/alphaticks/alphac/models"
	"math"
	"sync"
)

type Instrument struct {
	sync.RWMutex
	MakerFee          float64
	BasePriceModel    PriceModel
	QuotePriceModel   PriceModel
	BuyTradeModel     BuyTradeModel
	SellTradeModel    SellTradeModel
	SampleValueChange []float64
	SampleMatchBid    []float64
	SampleMatchAsk    []float64
	SampleBasePrice   []float64
	SampleQuotePrice  []float64
	SampleTime        uint64
	BidOrder          *Order
	AskOrder          *Order
}

func (instr *Instrument) updateSampleValueChange(time uint64) {
	sampleValueChange := instr.SampleValueChange
	N := len(sampleValueChange)
	sampleMatchBid := instr.SellTradeModel.GetSampleMatchBid(time, N)
	sampleMatchAsk := instr.BuyTradeModel.GetSampleMatchAsk(time, N)
	sampleBasePrice := instr.BasePriceModel.GetSamplePrices(time, N)
	sampleQuotePrice := instr.QuotePriceModel.GetSamplePrices(time, N)
	for i := 0; i < N; i++ {
		sampleValueChange[i] = 0.
	}
	if instr.BidOrder != nil {
		queue := instr.BidOrder.Queue
		quantity := instr.BidOrder.Quantity
		price := instr.BidOrder.Price
		for i := 0; i < N; i++ {
			baseChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			// We get our reduced-by-fees base
			sampleValueChange[i] += (baseChange - (baseChange * instr.MakerFee)) * sampleBasePrice[i]
			// We lose our quote
			sampleValueChange[i] -= baseChange * price * sampleQuotePrice[i]
		}
	}
	if instr.AskOrder != nil {
		queue := instr.AskOrder.Queue
		quantity := instr.AskOrder.Quantity
		price := instr.AskOrder.Price
		for i := 0; i < N; i++ {
			baseChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			// We lose our base
			sampleValueChange[i] -= baseChange * sampleBasePrice[i]
			// We get our reduced-by-fees quote
			sampleValueChange[i] += (baseChange - (baseChange * instr.MakerFee)) * price * sampleQuotePrice[i]
		}
	}
	instr.SampleBasePrice = sampleBasePrice
	instr.SampleQuotePrice = sampleQuotePrice
	instr.SampleMatchBid = sampleMatchBid
	instr.SampleMatchAsk = sampleMatchAsk
	instr.SampleTime = time
}

func (instr *Instrument) AddSampleValueChange(time uint64, values []float64) {
	instr.Lock()
	// Update this instrument sample value change only if bid order or ask order set
	if instr.BidOrder != nil || instr.AskOrder != nil {
		if instr.SampleTime != time {
			instr.updateSampleValueChange(time)
		}
		N := len(values)
		sampleValueChange := instr.SampleValueChange
		for i := 0; i < N; i++ {
			values[i] += sampleValueChange[i]
		}
	}

	instr.Unlock()
}

func (instr *Instrument) GetELROnCancelBid(time uint64, values []float64, value float64) float64 {
	N := len(values)
	instr.Lock()

	if instr.SampleTime != time {
		instr.updateSampleValueChange(time)
	}

	sampleMatchBid := instr.SampleMatchBid
	sampleBasePrice := instr.SampleBasePrice
	sampleQuotePrice := instr.SampleQuotePrice

	// Remove bid order from value
	if instr.BidOrder != nil {
		queue := instr.BidOrder.Queue
		quantity := instr.BidOrder.Quantity
		price := instr.BidOrder.Price
		for i := 0; i < N; i++ {
			baseChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			// We don't get our reduced-by-fees base
			values[i] -= (baseChange - (baseChange * instr.MakerFee)) * sampleBasePrice[i]
			// We keep our quote
			values[i] += baseChange * price * sampleQuotePrice[i]
		}
	}
	instr.Unlock()

	expectedLogReturn := 0.
	for i := 0; i < N; i++ {
		expectedLogReturn += math.Log(values[i] / value)
	}
	expectedLogReturn /= float64(N)

	return expectedLogReturn
}

func (instr *Instrument) GetELROnCancelAsk(time uint64, values []float64, value float64) float64 {
	N := len(values)
	instr.Lock()

	if instr.SampleTime != time {
		instr.updateSampleValueChange(time)
	}

	sampleMatchAsk := instr.SampleMatchAsk
	sampleBasePrice := instr.SampleBasePrice
	sampleQuotePrice := instr.SampleQuotePrice

	// Remove ask order from value
	if instr.AskOrder != nil {
		queue := instr.AskOrder.Queue
		quantity := instr.AskOrder.Quantity
		price := instr.AskOrder.Price
		for i := 0; i < N; i++ {
			baseChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			// We keep our base
			values[i] += baseChange * sampleBasePrice[i]
			// We don't get our reduced-by-fees quote
			values[i] -= (baseChange - (baseChange * instr.MakerFee)) * price * sampleQuotePrice[i]
		}
	}
	instr.Unlock()

	expectedLogReturn := 0.
	for i := 0; i < N; i++ {
		expectedLogReturn += math.Log(values[i] / value)
	}
	expectedLogReturn /= float64(N)

	return expectedLogReturn
}

func (instr *Instrument) GetELROnBidChange(time uint64, values []float64, value float64, prices []float64, queues []float64, maxQuote float64) (float64, *Order) {
	N := len(values)
	instr.Lock()
	// We want to see which option is the best, update, do nothing, or cancel
	if instr.SampleTime != time {
		instr.updateSampleValueChange(time)
	}

	sampleMatchBid := instr.SampleMatchBid
	sampleBasePrice := instr.SampleBasePrice
	sampleQuotePrice := instr.SampleQuotePrice

	var maxOrder *Order
	maxExpectedLogReturn := -999.

	// If we have an a bid order, we remove it from values
	if instr.BidOrder != nil {
		queue := instr.BidOrder.Queue
		quantity := instr.BidOrder.Quantity
		price := instr.BidOrder.Price
		expectedLogReturn := 0.
		for i := 0; i < N; i++ {
			expectedLogReturn += math.Log(values[i] / value)
			baseChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			// We keep our base
			values[i] += baseChange * sampleBasePrice[i]
			// We don't get our reduced-by-fees quote
			values[i] -= (baseChange - (baseChange * instr.MakerFee)) * price * sampleQuotePrice[i]
		}
		expectedLogReturn /= float64(N)
		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = instr.BidOrder
		}
	}

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

	for l := 0; l < len(prices); l++ {
		price := prices[l]
		queue := queues[l]
		expectedLogReturn := 0.

		for i := 0; i < N; i++ {
			expectedMatchQuote := math.Min(math.Max(sampleMatchBid[i]-queue, 0)*price, maxQuote)
			expectedMatchBase := expectedMatchQuote / price
			valueChange := (expectedMatchBase*(1-instr.MakerFee))*sampleBasePrice[i] - expectedMatchQuote*sampleQuotePrice[i]
			expectedLogReturn += math.Log((values[i] + valueChange) / value)
		}
		expectedLogReturn /= float64(N)

		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = &Order{
				Price:    price,
				Quantity: 0., // Computed at the end
				Queue:    queue,
			}
		}
	}

	if maxOrder != nil && maxOrder != instr.BidOrder {
		// Compute recommended order quantity based on price and queue
		// Recommended order quantity is the expected match at the level
		expectedMatch := 0.
		for i := 0; i < N; i++ {
			expectedMatch += math.Max(sampleMatchBid[i]-maxOrder.Queue, 0.)
		}
		expectedMatch /= float64(N)
		maxOrder.Quantity = math.Min(expectedMatch, maxQuote/maxOrder.Price)
	}
	instr.Unlock()

	return maxExpectedLogReturn, maxOrder
}

func (instr *Instrument) GetELROnAskChange(time uint64, values []float64, value float64, prices []float64, queues []float64, maxBase float64) (float64, *Order) {
	N := len(values)
	instr.Lock()
	// We want to see which option is the best, update, do nothing, or cancel
	if instr.SampleTime != time {
		instr.updateSampleValueChange(time)
	}

	sampleMatchAsk := instr.SampleMatchAsk
	sampleBasePrice := instr.SampleBasePrice
	sampleQuotePrice := instr.SampleQuotePrice

	var maxOrder *Order
	maxExpectedLogReturn := -999.

	// Remove ask order from value
	if instr.AskOrder != nil {
		queue := instr.AskOrder.Queue
		quantity := instr.AskOrder.Quantity
		price := instr.AskOrder.Price
		expectedLogReturn := 0.
		for i := 0; i < N; i++ {
			expectedLogReturn += math.Log(values[i] / value)
			baseChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			// We keep our base
			values[i] += baseChange * sampleBasePrice[i]
			// We don't get our reduced-by-fees quote
			values[i] -= (baseChange - (baseChange * instr.MakerFee)) * price * sampleQuotePrice[i]
		}
		expectedLogReturn /= float64(N)
		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = instr.AskOrder
		}
	}

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
			quoteValueChange := (expectedMatch - (instr.MakerFee * expectedMatch)) * price * sampleQuotePrice[i]

			expectedLogReturn += math.Log((values[i] + baseValueChange + quoteValueChange) / value)
		}

		expectedLogReturn /= float64(N)

		//fmt.Println(expectedLogReturn, maxExpectedLogReturn)
		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = &Order{
				Price:    price,
				Quantity: 0., // Computed at the end
				Queue:    queue,
			}
		}
	}

	if maxOrder != nil && maxOrder != instr.AskOrder {
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
	instr.Unlock()

	return maxExpectedLogReturn, maxOrder
}

type SpotPortfolio struct {
	exchange    uint64
	sampleSize  int
	assets      map[uint32]float64
	instruments map[uint64]*Instrument
	priceModels map[uint32]PriceModel
}

func NewSpotPortfolio(exchange uint64, securities []*models.Security, priceModels map[uint32]PriceModel, buyTradeModels map[uint64]BuyTradeModel, sellTradeModels map[uint64]SellTradeModel, sampleSize int) (*SpotPortfolio, error) {

	p := &SpotPortfolio{
		exchange:    exchange,
		sampleSize:  sampleSize,
		assets:      make(map[uint32]float64),
		instruments: make(map[uint64]*Instrument),
		priceModels: priceModels,
	}

	for _, security := range securities {
		ID := security.SecurityID
		p.assets[security.Underlying.ID] = 0.
		p.assets[security.QuoteCurrency.ID] = 0.
		p.instruments[ID] = &Instrument{
			RWMutex:           sync.RWMutex{},
			BuyTradeModel:     buyTradeModels[ID],
			SellTradeModel:    sellTradeModels[ID],
			BasePriceModel:    priceModels[security.Underlying.ID],
			QuotePriceModel:   priceModels[security.QuoteCurrency.ID],
			MakerFee:          0,
			SampleValueChange: make([]float64, p.sampleSize, p.sampleSize),
			SampleMatchBid:    make([]float64, p.sampleSize, p.sampleSize),
			SampleMatchAsk:    make([]float64, p.sampleSize, p.sampleSize),
			SampleBasePrice:   make([]float64, p.sampleSize, p.sampleSize),
			SampleQuotePrice:  make([]float64, p.sampleSize, p.sampleSize),
			SampleTime:        0,
		}
		if _, ok := priceModels[security.Underlying.ID]; !ok {
			return nil, fmt.Errorf("no price model for asset %s", security.Underlying.Name)
		}
		if _, ok := priceModels[security.QuoteCurrency.ID]; !ok {
			return nil, fmt.Errorf("no price model for asset %s", security.QuoteCurrency.Name)
		}
		if _, ok := buyTradeModels[ID]; !ok {
			return nil, fmt.Errorf("no trade model for security %d", ID)
		}
		if _, ok := sellTradeModels[ID]; !ok {
			return nil, fmt.Errorf("no trade model for security %d", ID)
		}
	}

	return p, nil
}

// The current portfolio value
func (p *SpotPortfolio) Value() float64 {
	var val float64 = 0

	for k, v := range p.assets {
		price := p.priceModels[k].GetPrice()
		val += v * price
	}
	return val
}

func (p *SpotPortfolio) GetAsset(asset uint32) float64 {
	v, ok := p.assets[asset]
	if !ok {
		return 0
	} else {
		return v
	}
}

func (p *SpotPortfolio) GetPrice(asset uint32) float64 {
	return p.priceModels[asset].GetPrice()
}

func (p *SpotPortfolio) GetExpectedPrice(time uint64, asset uint32) float64 {
	samples := p.priceModels[asset].GetSamplePrices(time, p.sampleSize)

	price := 0.
	for i := 0; i < p.sampleSize; i++ {
		price += samples[i]
	}

	return price / float64(p.sampleSize)
}

func (p *SpotPortfolio) SetAsset(asset uint32, quantity float64) {
	p.assets[asset] = quantity
}

func (p *SpotPortfolio) SetMakerFee(fee float64) {
	// We update sampleAssets by updating orders with old maker fee
	for _, instr := range p.instruments {
		instr.Lock()
		if instr.BidOrder != nil {
			sampleMatch := instr.SampleMatchBid
			sampleBasePrice := instr.SampleBasePrice
			queue := instr.BidOrder.Queue
			quantity := instr.BidOrder.Quantity
			for i := 0; i < p.sampleSize; i++ {
				baseChange := math.Min(math.Max(sampleMatch[i]-queue, 0.), quantity)
				// We update the value change caused by a match with the new fee
				instr.SampleValueChange[i] -= (baseChange - (baseChange * instr.MakerFee)) * sampleBasePrice[i]
				instr.SampleValueChange[i] += (baseChange - (baseChange * fee)) * sampleBasePrice[i]
			}
		}
		if instr.AskOrder != nil {
			sampleMatch := instr.SampleMatchAsk
			sampleQuotePrice := instr.SampleQuotePrice
			queue := instr.AskOrder.Queue
			quantity := instr.AskOrder.Quantity
			price := instr.AskOrder.Price
			for i := 0; i < p.sampleSize; i++ {
				baseChange := math.Min(math.Max(sampleMatch[i]-queue, 0.), quantity)
				// We update the value change caused by a match with the new fee
				instr.SampleValueChange[i] -= (baseChange - (baseChange * instr.MakerFee)) * price * sampleQuotePrice[i]
				instr.SampleValueChange[i] += (baseChange - (baseChange * fee)) * price * sampleQuotePrice[i]
			}
		}
		instr.MakerFee = fee
		instr.Unlock()
	}
}

func (p *SpotPortfolio) SetBidOrder(instrID uint64, newO Order) {
	// Update bid order
	instr := p.instruments[instrID]
	instr.Lock()
	sampleBasePrice := instr.SampleBasePrice
	sampleQuotePrice := instr.SampleQuotePrice
	sampleValueChange := instr.SampleValueChange
	sampleMatchBid := instr.SampleMatchBid
	if instr.BidOrder != nil {
		queue := instr.BidOrder.Queue
		quantity := instr.BidOrder.Quantity
		price := instr.BidOrder.Price
		// Bid order change, so update sample value changes
		for i := 0; i < p.sampleSize; i++ {
			baseChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			// We don't get our reduced-by-fees base
			sampleValueChange[i] -= (baseChange - (baseChange * instr.MakerFee)) * sampleBasePrice[i]
			// We keep our quote
			sampleValueChange[i] += baseChange * price * sampleQuotePrice[i]
		}
	}

	queue := newO.Queue
	quantity := newO.Quantity
	price := newO.Price
	for i := 0; i < p.sampleSize; i++ {
		baseChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
		// We get our reduced-by-fees base
		sampleValueChange[i] += (baseChange - (baseChange * instr.MakerFee)) * sampleBasePrice[i]
		// We lose our quote
		sampleValueChange[i] -= baseChange * price * sampleQuotePrice[i]
	}
	instr.BidOrder = &newO
	instr.Unlock()
}

func (p *SpotPortfolio) SetAskOrder(instrID uint64, newO Order) {
	// Update bid order
	instr := p.instruments[instrID]
	instr.Lock()
	sampleBasePrice := instr.SampleBasePrice
	sampleQuotePrice := instr.SampleQuotePrice
	sampleValueChange := instr.SampleValueChange
	sampleMatchAsk := instr.SampleMatchAsk
	if instr.AskOrder != nil {
		queue := instr.AskOrder.Queue
		quantity := instr.AskOrder.Quantity
		price := instr.AskOrder.Price
		for i := 0; i < p.sampleSize; i++ {
			baseChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			// We keep our base
			sampleValueChange[i] += baseChange * sampleBasePrice[i]
			// We don't get our reduced-by-fees quote
			sampleValueChange[i] -= (baseChange - (baseChange * instr.MakerFee)) * price * sampleQuotePrice[i]
		}
	}
	queue := newO.Queue
	quantity := newO.Quantity
	price := newO.Price
	for i := 0; i < p.sampleSize; i++ {
		baseChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
		// We lose our base
		sampleValueChange[i] -= baseChange * sampleBasePrice[i]
		// We get our reduced-by-fees quote
		sampleValueChange[i] += (baseChange - (baseChange * instr.MakerFee)) * price * sampleQuotePrice[i]
	}
	instr.AskOrder = &newO
	instr.Unlock()
}

func (p *SpotPortfolio) CancelBidOrder(instrID uint64) {
	instr := p.instruments[instrID]
	instr.Lock()
	if instr.BidOrder != nil {
		sampleBasePrice := instr.SampleBasePrice
		sampleQuotePrice := instr.SampleQuotePrice
		sampleValueChange := instr.SampleValueChange
		sampleMatchBid := instr.SampleMatchBid
		queue := instr.BidOrder.Queue
		quantity := instr.BidOrder.Quantity
		price := instr.BidOrder.Price
		for i := 0; i < p.sampleSize; i++ {
			baseChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			// We don't get our reduced-by-fees base
			sampleValueChange[i] -= (baseChange - (baseChange * instr.MakerFee)) * sampleBasePrice[i]
			// We keep our quote
			sampleValueChange[i] += baseChange * price * sampleQuotePrice[i]
		}
		instr.BidOrder = nil
	}
	instr.Unlock()
}

func (p *SpotPortfolio) CancelAskOrder(instrID uint64) {
	instr := p.instruments[instrID]
	instr.Lock()
	if instr.AskOrder != nil {
		sampleBasePrice := instr.SampleBasePrice
		sampleQuotePrice := instr.SampleQuotePrice
		sampleValueChange := instr.SampleValueChange
		sampleMatchAsk := instr.SampleMatchAsk
		queue := instr.AskOrder.Queue
		quantity := instr.AskOrder.Quantity
		price := instr.AskOrder.Price
		for i := 0; i < p.sampleSize; i++ {
			baseChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			// We keep our base
			sampleValueChange[i] += baseChange * sampleBasePrice[i]
			// We don't get our reduced-by-fees quote
			sampleValueChange[i] -= (baseChange - (baseChange * instr.MakerFee)) * price * sampleQuotePrice[i]
		}
		instr.AskOrder = nil
	}
	instr.Unlock()
}

func (p *SpotPortfolio) UpdateSampleAssetValue(time uint64, values []float64) {

}

// What triggers a change in sample values ?
// - A change in sample asset
// - A change in sample prices
func (p *SpotPortfolio) AddSampleValues(model modeling.Model, time uint64, values []float64) {
	N := p.sampleSize
	// TODO add state var to prevent re-computing fixed assets value
	// TODO asset lock
	for k, v := range p.assets {
		samplePrices := model.GetSampleAssetPrices(k, time, p.sampleSize)
		for i := 0; i < N; i++ {
			values[i] += v * samplePrices[i]
		}
	}
	// Add value changes due to orders
	for _, instr := range p.instruments {
		instr.AddSampleValueChange(time, values)
	}
}

func (p *SpotPortfolio) GetELROnCancelBid(instrID uint64, time uint64, values []float64, value float64) float64 {
	return p.instruments[instrID].GetELROnCancelBid(time, values, value)
}

func (p *SpotPortfolio) GetELROnCancelAsk(instrID uint64, time uint64, values []float64, value float64) float64 {
	return p.instruments[instrID].GetELROnCancelAsk(time, values, value)
}

func (p *SpotPortfolio) GetELROnBidChange(instrID uint64, time uint64, values []float64, value float64, prices []float64, queues []float64, maxQuote float64) (float64, *Order) {
	return p.instruments[instrID].GetELROnBidChange(time, values, value, prices, queues, maxQuote)
}

func (p *SpotPortfolio) GetELROnAskChange(instrID uint64, time uint64, values []float64, value float64, prices []float64, queues []float64, maxBase float64) (float64, *Order) {
	return p.instruments[instrID].GetELROnAskChange(time, values, value, prices, queues, maxBase)
}
