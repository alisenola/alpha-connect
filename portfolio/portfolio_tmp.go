package portfolio

/*
import (
	"fmt"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"math"
	"sync"
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

type Order struct {
	Price    float64
	Quantity float64
	Queue    float64
	Base     uint32
	Quote    uint32
}

type ExchangePortfolio struct {
	sampleSize       int
	prices           *sync.Map
	assetsLock       sync.RWMutex // one exchange update, others read to get value
	assets           map[uint32]uint64
	lotPrecisions    map[uint32]float64
	sampleTime       uint64
	instruments      map[uint64]*exchanges.Instrument
	sampleMatchBid   map[uint64][]uint64
	sampleMatchAsk   map[uint64][]uint64
	bidOrder         map[uint64]Order // One bid order on all instruments
	askOrder         map[uint64]Order // One ask order on all instruments
	sampleAssetsLock sync.RWMutex     // one exchange update, others read to get exp log ret
	sampleAssets     map[uint32][]uint64
	takerFee         float64
	makerFee         float64
}

func NewExchangePortfolio(instrs []*exchanges.Instrument, sampleSize int, prices *sync.Map) *ExchangePortfolio {

	p := &ExchangePortfolio{
		sampleSize:     sampleSize,
		prices:         prices,
		assets:         make(map[uint32]uint64),
		sampleTime:     0,
		sampleMatchBid: make(map[uint64][]uint64),
		sampleMatchAsk: make(map[uint64][]uint64),
		bidOrder:       make(map[uint64]Order),
		askOrder:       make(map[uint64]Order),
		sampleAssets:   make(map[uint32][]uint64),
		takerFee:       0.,
		makerFee:       0.,
	}

	p.lotPrecisions = make(map[uint32]float64)
	for _, instr := range instrs {
		lp, ok := p.lotPrecisions[instr.Pair.Base.ID]
		if !ok || float64(instr.LotPrecision) > lp {
			p.lotPrecisions[instr.Pair.Base.ID] = float64(instr.LotPrecision)
			p.lotPrecisions[instr.Pair.Quote.ID] = float64(instr.LotPrecision * instr.TickPrecision)
		}
		p.assets[instr.Pair.Base.ID] = 0
		p.assets[instr.Pair.Quote.ID] = 0
	}

	return p
}

// The current portfolio value
func (p *ExchangePortfolio) Value() float64 {
	var val float64 = 0
	// Lock assets
	p.assetsLock.RLock()
	defer p.assetsLock.RUnlock()

	for k, v := range p.assets {
		price, ok := p.prices.Load(k)
		if !ok {
			price = 0.
		}
		val += (float64(v) / p.lotPrecisions[k]) * price.(float64)
	}
	return val
}

func (p *ExchangePortfolio) UpdateInstruments(instrs []*exchanges.Instrument) {
	p.lotPrecisions = make(map[uint32]float64)
	for _, instr := range instrs {
		lp, ok := p.lotPrecisions[instr.Pair.Base.ID]
		if !ok || float64(instr.LotPrecision) > lp {
			p.lotPrecisions[instr.Pair.Base.ID] = float64(instr.LotPrecision)
		}
	}
	// TODO update sample assets
}

func (p *ExchangePortfolio) GetAsset(asset uint32) float64 {
	p.assetsLock.RLock()
	defer p.assetsLock.RUnlock()
	qty, ok := p.assets[asset]
	if !ok {
		qty = 0
	}
	return float64(qty) / p.lotPrecisions[asset]
}

func (p *ExchangePortfolio) SetAsset(asset uint32, quantityF float64) {
	p.assetsLock.Lock()
	lastValue, ok := p.assets[asset]
	if !ok {
		lastValue = 0
	}
	quantity := uint64(quantityF * p.lotPrecisions[asset])
	p.assets[asset] = quantity
	p.assetsLock.Unlock()

	N := p.sampleSize

	// Update sample assets
	p.sampleAssetsLock.Lock()
	if sampleAsset, ok := p.sampleAssets[asset]; ok {
		for i := 0; i < N; i++ {
			sampleAsset[i] -= lastValue
			sampleAsset[i] += quantity
		}
	} else {
		sampleAsset := make([]uint64, N, N)
		for i := 0; i < N; i++ {
			sampleAsset[i] = quantity
		}
		p.sampleAssets[asset] = sampleAsset
	}
	p.sampleAssetsLock.Unlock()
}

func (p *ExchangePortfolio) SetSampleMatch(time uint64, base, quote uint32, sampleBidTmp, sampleAskTmp []float64) {
	instrumentID := uint64(base)<<32 | uint64(quote)

	lpb := p.lotPrecisions[base]
	lpq := p.lotPrecisions[quote]
	N := len(sampleBidTmp)
	sampleBid := make([]uint64, N, N)
	sampleAsk := make([]uint64, N, N)
	for i := 0; i < N; i++ {
		sampleBid[i] = uint64(lpb * sampleBidTmp[i])
		sampleAsk[i] = uint64(lpb * sampleAskTmp[i])
	}
	p.sampleAssetsLock.Lock()
	// Recompute sample asset
	if o, ok := p.bidOrder[instrumentID]; ok {
		queue := uint64(o.Queue * lpb)
		qty := uint64(o.Quantity * lpb)
		// Remove old bid sim
		if oldSampleMatch, ok := p.sampleMatchBid[instrumentID]; ok {
			for i := 0; i < p.sampleSize; i++ {
				var baseChange uint64
				if oldSampleMatch[i] > queue {
					diff := oldSampleMatch[i] - queue
					if diff > qty {
						baseChange = qty
					} else {
						baseChange = diff
					}
				} else {
					baseChange = 0
				}
				// We don't get our reduced-by-fees base
				p.sampleAssets[base][i] -= baseChange - uint64(float64(baseChange) * p.makerFee)
				// We keep our quote
				p.sampleAssets[quote][i] += uint64(float64(baseChange) * o.Price * (lpq / lpb))
			}
		}
		for i := 0; i < p.sampleSize; i++ {
			// We update with new sample match
			var baseChange uint64
			if sampleBid[i] > queue {
				diff := sampleBid[i] - queue
				if diff > qty {
					baseChange = qty
				} else {
					baseChange = diff
				}
			} else {
				baseChange = 0
			}
			// We get our reduced-by-fees base
			p.sampleAssets[base][i] += baseChange - uint64(float64(baseChange) * p.makerFee)
			// We lose our quote
			p.sampleAssets[quote][i] -= uint64(float64(baseChange) * o.Price * (lpq / lpb))
		}
	}
	// Recompute sample asset
	if o, ok := p.askOrder[instrumentID]; ok {
		queue := uint64(o.Queue * lpb)
		qty := uint64(o.Quantity * lpb)
		if oldSampleMatch, ok := p.sampleMatchAsk[instrumentID]; ok {
			for i := 0; i < p.sampleSize; i++ {
				var baseChange uint64
				if oldSampleMatch[i] > queue {
					diff := oldSampleMatch[i] - queue
					if diff > qty {
						baseChange = qty
					} else {
						baseChange = diff
					}
				} else {
					baseChange = 0
				}
				// We keep our base
				p.sampleAssets[base][i] += baseChange
				// We don't get our reduced-by-fees quote
				fee := uint64(float64(baseChange) * p.makerFee)
				p.sampleAssets[quote][i] -= uint64(float64(baseChange - fee) * o.Price * (lpq / lpb))
			}
		}
		for i := 0; i < p.sampleSize; i++ {
			// We update with new sample match
			var baseChange uint64
			if sampleAsk[i] > queue {
				diff := sampleAsk[i] - queue
				if diff > qty {
					baseChange = qty
				} else {
					baseChange = diff
				}
			} else {
				baseChange = 0
			}
			// We lose our base
			p.sampleAssets[base][i] -= baseChange
			// We get our reduced-by-fees quote
			fee := uint64(float64(baseChange) * p.makerFee)
			p.sampleAssets[quote][i] += uint64(float64(baseChange - fee) * o.Price * (lpq / lpb))
		}
	}
	p.sampleAssetsLock.Unlock()

	p.sampleMatchBid[instrumentID] = sampleBid
	p.sampleMatchAsk[instrumentID] = sampleAsk
	p.sampleTime = time
}

func (p *ExchangePortfolio) SetMakerFee(fee float64) {
	// We update sampleAssets by updating orders with old maker fee
	for instrID, o := range p.bidOrder {
		base := uint32(instrID >> 32)
		lpb := p.lotPrecisions[base]
		queue := uint64(o.Queue * lpb)
		qty := uint64(o.Quantity * lpb)
		sampleMatch := p.sampleMatchBid[instrID]
		sampleBase := p.sampleAssets[base]
		for i := 0; i < p.sampleSize; i++ {
			var baseChange uint64
			if sampleMatch[i] > queue {
				diff := sampleMatch[i] - queue
				if diff > qty {
					baseChange = qty
				} else {
					baseChange = diff
				}
			} else {
				baseChange = 0
			}
			// We update our reduced-by-fees base
			sampleBase[i] -= baseChange - uint64(float64(baseChange) * p.makerFee)
			sampleBase[i] += baseChange - uint64(float64(baseChange) * fee)
		}
	}

	for instrID, o := range p.askOrder {
		base := uint32(instrID >> 32)
		quote := uint32(instrID & math.MaxUint32)
		lpb := p.lotPrecisions[base]
		lpq := p.lotPrecisions[quote]
		queue := uint64(o.Queue * lpq)
		qty := uint64(o.Quantity * lpq)
		sampleMatch := p.sampleMatchAsk[instrID]
		sampleQuote := p.sampleAssets[quote]
		for i := 0; i < p.sampleSize; i++ {
			var baseChange uint64
			if sampleMatch[i] > queue {
				diff := sampleMatch[i] - queue
				if diff > qty {
					baseChange = qty
				} else {
					baseChange = diff
				}
			} else {
				baseChange = 0
			}
			// We update our reduced-by-fees quote
			sampleQuote[i] -= uint64(float64(baseChange - uint64(float64(baseChange) * p.makerFee)) * o.Price * (lpq / lpb))
			sampleQuote[i] += uint64(float64(baseChange - uint64(float64(baseChange) * fee)) * o.Price * (lpq / lpb))
		}
	}

	p.makerFee = fee
}

func (p *ExchangePortfolio) SetBidOrder(base, quote uint32, newO Order) {
	instrumentID := uint64(base)<<32 | uint64(quote)
	lpb := p.lotPrecisions[base]
	lpq := p.lotPrecisions[quote]

	if sample, ok := p.sampleMatchBid[instrumentID]; ok {
		p.sampleAssetsLock.Lock()
		sampleBase, ok := p.sampleAssets[base]
		if !ok {
			sampleBase = make([]uint64, p.sampleSize, p.sampleSize)
			for i := 0; i < p.sampleSize; i++ {
				sampleBase[i] = 0
			}
			p.sampleAssets[base] = sampleBase
		}
		sampleQuote := p.sampleAssets[quote]
		if o, ok := p.bidOrder[instrumentID]; ok {
			queue := uint64(o.Queue * lpb)
			qty := uint64(o.Quantity * lpb)
			for i := 0; i < p.sampleSize; i++ {
				var baseChange uint64
				if sample[i] > queue {
					diff := sample[i] - queue
					if diff > qty {
						baseChange = qty
					} else {
						baseChange = diff
					}
				} else {
					baseChange = 0
				}
				// We don't get our reduced-by-fees base
				sampleBase[i] -= baseChange - uint64(float64(baseChange) * p.makerFee)
				// We keep our quote
				sampleQuote[i] += uint64(float64(baseChange) * o.Price * (lpq / lpb))
			}
		}

		queue := uint64(newO.Queue * lpb)
		qty := uint64(newO.Quantity * lpb)
		for i := 0; i < p.sampleSize; i++ {
			var baseChange uint64
			if sample[i] > queue {
				diff := sample[i] - queue
				if diff > qty {
					baseChange = qty
				} else {
					baseChange = diff
				}
			} else {
				baseChange = 0
			}
			// We get our reduced-by-fees base
			sampleBase[i] += baseChange -  uint64(float64(baseChange) * p.makerFee)
			// We lose our quote
			sampleQuote[i] -= uint64(float64(baseChange) * newO.Price * (lpq / lpb))
		}
		p.sampleAssetsLock.Unlock()
	}

	p.bidOrder[instrumentID] = newO
}

func (p *ExchangePortfolio) SetAskOrder(base, quote uint32, newO Order) {
	instrumentID := uint64(base)<<32 | uint64(quote)
	lpb := p.lotPrecisions[base]
	lpq := p.lotPrecisions[quote]
	if sample, ok := p.sampleMatchAsk[instrumentID]; ok {
		p.sampleAssetsLock.Lock()
		sampleBase, ok := p.sampleAssets[base]
		if !ok {
			sampleBase = make([]uint64, p.sampleSize, p.sampleSize)
			for i := 0; i < p.sampleSize; i++ {
				sampleBase[i] = 0
			}
			p.sampleAssets[base] = sampleBase
		}
		sampleQuote := p.sampleAssets[quote]
		if o, ok := p.askOrder[instrumentID]; ok {
			queue := uint64(o.Queue * lpb)
			qty := uint64(o.Quantity * lpb)
			for i := 0; i < p.sampleSize; i++ {
				var baseChange uint64
				if sample[i] > queue {
					diff := sample[i] - queue
					if diff > qty {
						baseChange = qty
					} else {
						baseChange = diff
					}
				} else {
					baseChange = 0
				}
				// We keep our base
				sampleBase[i] += baseChange
				// We don't get our reduced-by-fees quote
				sampleQuote[i] -= uint64(float64(baseChange - uint64(float64(baseChange) * p.makerFee)) * o.Price * (lpq / lpb))
			}
		}
		queue := uint64(newO.Queue * lpb)
		qty := uint64(newO.Quantity * lpb)
		for i := 0; i < p.sampleSize; i++ {
			var baseChange uint64
			if sample[i] > queue {
				diff := sample[i] - queue
				if diff > qty {
					baseChange = qty
				} else {
					baseChange = diff
				}
			} else {
				baseChange = 0
			}
			// We lose our base
			sampleBase[i] -= baseChange
			// We get our reduced-by-fees quote
			sampleQuote[i] += uint64(float64(baseChange - uint64(float64(baseChange) * p.makerFee)) * newO.Price * (lpq / lpb))
		}
		p.sampleAssetsLock.Unlock()
	}

	p.askOrder[instrumentID] = newO
}

func (p *ExchangePortfolio) CancelBidOrder(base, quote uint32) {
	// look up order by order ID
	instrumentID := uint64(base)<<32 | uint64(quote)

	p.sampleAssetsLock.Lock()
	if o, ok := p.bidOrder[instrumentID]; ok {
		lpb := p.lotPrecisions[base]
		lpq := p.lotPrecisions[quote]

		sample := p.sampleMatchBid[instrumentID]
		sampleBase := p.sampleAssets[base]
		sampleQuote := p.sampleAssets[quote]
		queue := uint64(o.Queue * lpb)
		qty := uint64(o.Quantity * lpb)
		for i := 0; i < p.sampleSize; i++ {
			var baseChange uint64
			if sample[i] > queue {
				diff := sample[i] - queue
				if diff > qty {
					baseChange = qty
				} else {
					baseChange = diff
				}
			} else {
				baseChange = 0
			}
			// We don't get our reduced-by-fees base
			sampleBase[i] -= baseChange - uint64(float64(baseChange) * p.makerFee)
			// We keep our quote
			sampleQuote[i] += uint64(float64(baseChange) * o.Price * (lpq / lpb))
		}
		delete(p.bidOrder, instrumentID)
	}
	p.sampleAssetsLock.Unlock()
}

func (p *ExchangePortfolio) CancelAskOrder(base, quote uint32) {
	// look up order by order ID
	instrumentID := uint64(base)<<32 | uint64(quote)

	p.sampleAssetsLock.Lock()
	if o, ok := p.askOrder[instrumentID]; ok {
		sample := p.sampleMatchAsk[instrumentID]
		lpb := p.lotPrecisions[base]
		lpq := p.lotPrecisions[quote]
		queue := uint64(o.Queue * lpb)
		qty := uint64(o.Quantity * lpb)
		sampleBase := p.sampleAssets[base]
		sampleQuote := p.sampleAssets[quote]

		for i := 0; i < p.sampleSize; i++ {
			var baseChange uint64
			if sample[i] > queue {
				diff := sample[i] - queue
				if diff > qty {
					baseChange = qty
				} else {
					baseChange = diff
				}
			} else {
				baseChange = 0
			}
			// We keep our base
			sampleBase[i] += baseChange
			// We don't get our reduced-by-fees quote
			sampleQuote[i] -= uint64(float64(baseChange - uint64(float64(baseChange) * p.makerFee)) * o.Price * (lpq / lpb))
		}
		delete(p.askOrder, instrumentID)
	}
	p.sampleAssetsLock.Unlock()
}

type Portfolio struct {
	SampleSize          int
	Prices              *sync.Map
	samplePrices        map[uint32][]float64 // Sample of possible future prices
	samplePricesLock    sync.RWMutex
	sampleTime          uint64
	exchangesPortfolios map[uint64]*ExchangePortfolio
}

func NewPortfolio(instruments []*exchanges.Instrument, sampleSize int) *Portfolio {
	p := &Portfolio{
		SampleSize:          sampleSize,
		Prices:              &sync.Map{},
		samplePrices:        make(map[uint32][]float64),
		exchangesPortfolios: make(map[uint64]*ExchangePortfolio),
	}
	exchangeInstruments := make(map[uint64][]*exchanges.Instrument)
	for _, instr := range instruments {
		exchangeID := constants.EXCHANGE_NAME_TO_EXCHANGE_ID[instr.Exchange]
		exchangeInstruments[exchangeID] = append(exchangeInstruments[exchangeID], instr)
	}
	for k, v := range exchangeInstruments {
		p.exchangesPortfolios[k] = NewExchangePortfolio(v, sampleSize, p.Prices)
	}

	return p
}

func (p *Portfolio) UpdateInstruments(instruments []*exchanges.Instrument) {
	exchangeInstruments := make(map[uint64][]*exchanges.Instrument)
	// Need to recompute asset lotPrecision
	for _, instr := range instruments {
		exchangeID := constants.EXCHANGE_NAME_TO_EXCHANGE_ID[instr.Exchange]
		exchangeInstruments[exchangeID] = append(exchangeInstruments[exchangeID], instr)
	}
	for k, v := range exchangeInstruments {
		p.exchangesPortfolios[k].UpdateInstruments(v)
	}
}

func (p *Portfolio) Value() float64 {
	value := 0.
	for _, exch := range p.exchangesPortfolios {
		value += exch.Value()
	}

	return value
}

func (p *Portfolio) SetTakerFee(exchange uint64, takerFee float64) {
	p.exchangesPortfolios[exchange].takerFee = takerFee
}

func (p *Portfolio) SetMakerFee(exchange uint64, makerFee float64) {
	p.exchangesPortfolios[exchange].SetMakerFee(makerFee)
}

func (p *Portfolio) SetAsset(exchange uint64, asset uint32, quantity float64) {
	p.exchangesPortfolios[exchange].SetAsset(asset, quantity)
}

func (p *Portfolio) GetAsset(exchange uint64, asset uint32) float64 {
	return p.exchangesPortfolios[exchange].GetAsset(asset)
}

func (p *Portfolio) SetSamplePrices(time uint64, asset uint32, prices []float64) {
	p.samplePricesLock.Lock()
	p.sampleTime = time
	p.samplePrices[asset] = prices
	p.samplePricesLock.Unlock()
}

func (p *Portfolio) SetSampleMatch(time uint64, exchange uint64, base uint32, quote uint32, sampleBids, sampleAsks []float64) {
	p.exchangesPortfolios[exchange].SetSampleMatch(time, base, quote, sampleBids, sampleAsks)
}

func (p *Portfolio) SetBidOrder(exchange uint64, order Order) {
	p.exchangesPortfolios[exchange].SetBidOrder(order.Base, order.Quote, order)
}

func (p *Portfolio) SetAskOrder(exchange uint64, order Order) {
	p.exchangesPortfolios[exchange].SetAskOrder(order.Base, order.Quote, order)
}

func (p *Portfolio) CancelBidOrder(exchange uint64, base uint32, quote uint32) {
	p.exchangesPortfolios[exchange].CancelBidOrder(base, quote)
}

func (p *Portfolio) CancelAskOrder(exchange uint64, base uint32, quote uint32) {
	p.exchangesPortfolios[exchange].CancelAskOrder(base, quote)
}

// Return the expected log return
func (p *Portfolio) ExpectedLogReturn() float64 {
	N := p.SampleSize

	var value float64 = 0
	values := make([]float64, N, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	// We need to ensure consistent sample prices across the computation
	p.samplePricesLock.Lock()
	for _, exch := range p.exchangesPortfolios {
		// Lock sample assets
		exch.sampleAssetsLock.RLock()
		for k, samples := range exch.sampleAssets {
			lp := exch.lotPrecisions[k]
			if samplePrices, ok := p.samplePrices[k]; ok {
				for i := 0; i < N; i++ {
					values[i] += (float64(samples[i]) / lp) * samplePrices[i]
				}
			} else {
				priceTmp, _ := p.Prices.Load(k)
				price := priceTmp.(float64)
				for i := 0; i < N; i++ {
					values[i] += (float64(samples[i]) / lp) * price
					fmt.Println(samples[i])
				}
			}
		}
		exch.sampleAssetsLock.RUnlock()

		value += exch.Value()
	}
	p.samplePricesLock.Unlock()

	var expectedLogReturn float64 = 0

	// Compute expected log return
	for i := 0; i < N; i++ {
		expectedLogReturn += math.Log(values[i] / value)
	}

	return expectedLogReturn / float64(N)
}

// Returns the order with the highest expected log return on cancel
func (p *Portfolio) ExpectedLogReturnOnCancel(exchange uint64) (float64, bool, *Order) {
	N := p.SampleSize

	var value float64 = 0
	values := make([]float64, N, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	// We need to ensure consistent sample prices across the computation
	p.samplePricesLock.Lock()
	for _, exch := range p.exchangesPortfolios {

		exch.sampleAssetsLock.RLock()
		for k, samples := range exch.sampleAssets {
			prices := p.samplePrices[k]
			lp := exch.lotPrecisions[k]
			for i := 0; i < N; i++ {
				values[i] += (float64(samples[i]) / lp) * prices[i]
			}
		}
		exch.sampleAssetsLock.RUnlock()

		value += exch.Value()
	}

	// Compute expected log return if you don't cancel anything
	maxExpectedLogReturn := 0.
	for i := 0; i < N; i++ {
		maxExpectedLogReturn += math.Log(values[i] / value)
	}
	maxExpectedLogReturn /= float64(N)

	exch := p.exchangesPortfolios[exchange]
	var maxOrder *Order

	for instrID, o := range exch.bidOrder {
		// Need to know, for each sample, the value change if we remove it
		// so need to compute expected match for each sample, then compute the change in portfolio value
		// if this order is not executed
		sampleMatch := exch.sampleMatchBid[instrID]
		base := uint32(instrID >> 32)
		quote := uint32(instrID & math.MaxUint32)

		basePrices := p.samplePrices[base]
		quotePrices := p.samplePrices[quote]
		lpb := exch.lotPrecisions[base]
		expectedLogReturn := 0.
		for i := 0; i < N; i++ {
			sm := float64(sampleMatch[i]) / lpb
			expectedMatch := math.Max(sm-o.Queue, 0)
			var valueChange float64
			// At this sample, we were expected to have gained expectedMatch base and lost expectedMatch * o.price quote
			// so we remove expectedMatch * basePrice
			baseValueChange := -(expectedMatch - (exch.makerFee * expectedMatch)) * basePrices[i]
			// but we got to keep our quote, so we add expectedMatch * o.price * quotePrice
			quoteValueChange := expectedMatch * o.Price * quotePrices[i]
			valueChange = baseValueChange + quoteValueChange
			expectedLogReturn += math.Log((values[i] + valueChange) / value)
		}
		expectedLogReturn /= float64(N)

		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = &o
		}
	}

	for instrID, o := range exch.askOrder {
		// Need to know, for each sample, the value change if we remove it
		// so need to compute expected match for each sample, then compute the change in portfolio value
		// if this order is not executed
		sampleMatch := exch.sampleMatchBid[instrID]
		base := uint32(instrID >> 32)
		quote := uint32(instrID & math.MaxUint32)

		basePrices := p.samplePrices[base]
		quotePrices := p.samplePrices[quote]
		lpb := exch.lotPrecisions[base]

		expectedLogReturn := 0.
		for i := 0; i < N; i++ {
			sm := float64(sampleMatch[i]) / lpb
			expectedMatch := math.Max(sm-o.Queue, 0)
			var valueChange float64
			// At this sample, we were expected to have lost expectedMatch base and gained expectedMatch * o.price quote
			// so we add expectedMatch * basePrice
			baseValueChange := expectedMatch * basePrices[i]
			// but we lose our quote, so we remove expectedMatch * o.price * quotePrice
			quoteValueChange := -(expectedMatch - (exch.makerFee * expectedMatch)) * o.Price * quotePrices[i]

			valueChange = baseValueChange + quoteValueChange
			expectedLogReturn += math.Log((values[i] + valueChange) / value)
		}
		expectedLogReturn /= float64(N)

		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = &o
		}
	}
	p.samplePricesLock.Unlock()

	return maxExpectedLogReturn, true, maxOrder
}

// We make multiple assumptions here. At each level, we assume the expectedMatch will be what will be set as order quantity
// This assumption kind of make sense, as if, in case of bid, the expected match is big at one level, might as well place the order
// at the level bellow to get a higher spread.
// When at ask, a // TODO think more about that
func (p *Portfolio) ExpectedLogReturnOnLimitBid(exchange uint64, base, quote uint32, prices []float64, queues []float64, maxQuote float64) (float64, *Order) {
	// Need to compute the expected log return
	exch := p.exchangesPortfolios[exchange]
	instrID := uint64(base)<<32 | uint64(quote)

	N := p.SampleSize

	var value float64 = 0
	values := make([]float64, N, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	// We need to ensure consistent sample prices across the computation
	p.samplePricesLock.Lock()
	for _, exch := range p.exchangesPortfolios {

		exch.sampleAssetsLock.RLock()
		for k, samples := range exch.sampleAssets {
			prices := p.samplePrices[k]
			lp := exch.lotPrecisions[k]
			for i := 0; i < N; i++ {
				values[i] += (float64(samples[i]) / lp) * prices[i]
			}
		}
		exch.sampleAssetsLock.RUnlock()

		value += exch.Value()
	}
	basePrices := p.samplePrices[base]
	quotePrices := p.samplePrices[quote]

	p.samplePricesLock.Unlock()
	// If we have a bid order already, it has been taken into account in the sampleAsset
	// so we need to remove it
	sampleMatch := exch.sampleMatchBid[instrID]
	lpb := exch.lotPrecisions[base]
	if o, ok := exch.bidOrder[instrID]; ok {
		for j := 0; j < N; j++ {
			sm := float64(sampleMatch[j]) / lpb
			expectedMatch := math.Max(sm-o.Queue, 0)
			// At this sample, we were expected to have gained expectedMatch base and lost expectedMatch * o.price quote
			// so we remove expectedMatch * basePrice
			baseValueChange := -(expectedMatch - (exch.makerFee * expectedMatch)) * basePrices[j]
			// but we got to keep our quote, so we add expectedMatch * o.price * quotePrice
			quoteValueChange := expectedMatch * o.Price * quotePrices[j]
			values[j] = values[j] + baseValueChange + quoteValueChange
		}
	}

	// Compute expected log return without any order
	maxExpectedLogReturn := 0.
	for i := 0; i < N; i++ {
		maxExpectedLogReturn += math.Log(values[i] / value)
	}
	maxExpectedLogReturn /= float64(N)
	var maxOrder *Order

	for i := 0; i < len(prices); i++ {
		price := prices[i]
		queue := queues[i]
		expectedLogReturn := 0.
		// We need to compute what would happen if we changed
		// our bid order on the instrument. A change in the bid order
		// means a change in the sampleAsset
		for j := 0; j < N; j++ {
			sm := float64(sampleMatch[j]) / lpb
			expectedMatchQuote := math.Max(sm-queue, 0) * price
			// We limit this by quote amount available in our portfolio
			expectedMatchQuote = math.Min(expectedMatchQuote, maxQuote)

			expectedMatchBase := expectedMatchQuote / price

			// We are expected to gain expectedMatch base
			baseValueChange := (expectedMatchBase - (exch.makerFee * expectedMatchBase)) * basePrices[j]

			// We are expected to lose expectedMatch * price quote
			quoteValueChange := -expectedMatchQuote * quotePrices[j]

			expectedLogReturn += math.Log((values[j] + baseValueChange + quoteValueChange) / value)
		}

		expectedLogReturn /= float64(N)
		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = &Order{
				Price:    price,
				Quantity: 0., // Computed at the end
				Queue:    queue,
				Base:     base,
				Quote:    quote,
			}
		}
	}

	if maxOrder != nil {
		// Compute recommended order quantity based on price and queue
		// Recommended order quantity is the expected match at the level
		expectedMatch := 0.
		for i := 0; i < N; i++ {
			expectedMatch += math.Max(float64(sampleMatch[i]) / lpb - maxOrder.Queue, 0.)
		}
		maxOrder.Quantity = expectedMatch / float64(N)

		//basePrice, _ := p.Prices.Load(base)
		//fmt.Println(basePrice, maxOrder.Price, maxOrder.Queue, maxExpectedLogReturn)
	}

	return maxExpectedLogReturn, maxOrder
}

func (p *Portfolio) ExpectedLogReturnOnLimitAsk(exchange uint64, base, quote uint32, prices []float64, queues []float64, maxBase float64) (float64, *Order) {
	// Need to compute the expected log return
	exch := p.exchangesPortfolios[exchange]
	instrID := uint64(base)<<32 | uint64(quote)

	N := p.SampleSize

	var value float64 = 0
	values := make([]float64, N, N)
	for i := 0; i < N; i++ {
		values[i] = 0.
	}

	p.samplePricesLock.Lock()
	for _, exch := range p.exchangesPortfolios {

		exch.sampleAssetsLock.RLock()
		for k, samples := range exch.sampleAssets {
			prices := p.samplePrices[k]
			lp := exch.lotPrecisions[k]
			for i := 0; i < N; i++ {
				values[i] += (float64(samples[i]) / lp) * prices[i]
			}
		}
		exch.sampleAssetsLock.RUnlock()

		value += exch.Value()
	}

	basePrices := p.samplePrices[base]
	quotePrices := p.samplePrices[quote]
	p.samplePricesLock.Unlock()

	// TODO maybe don't do that and let the user call cancel bid order
	// If we have a ask order already, it has been taken into account in the sampleAsset
	// so we need to remove it
	sampleMatch := exch.sampleMatchAsk[instrID]
	lpb := exch.lotPrecisions[base]
	if o, ok := exch.askOrder[instrID]; ok {
		for j := 0; j < N; j++ {
			sm := float64(sampleMatch[j]) / lpb
			expectedMatch := math.Max(sm-o.Queue, 0)
			// At this sample, we were expected to have lost expectedMatch base and gained expectedMatch * o.price quote
			// so we add expectedMatch * basePrice
			baseValueChange := expectedMatch * basePrices[j]
			// but we lose our quote, so we remove expectedMatch * o.price * quotePrice
			quoteValueChange := -(expectedMatch - (exch.makerFee * expectedMatch)) * o.Price * quotePrices[j]
			values[j] = values[j] + baseValueChange + quoteValueChange
		}
	}

	// Compute expected log return if you don't do anything
	maxExpectedLogReturn := 0.
	for i := 0; i < N; i++ {
		maxExpectedLogReturn += math.Log(values[i] / value)
	}
	maxExpectedLogReturn /= float64(N)
	//baseline := maxExpectedLogReturn
	var maxOrder *Order

	for i := 0; i < len(prices); i++ {
		price := prices[i]
		queue := queues[i]
		expectedLogReturn := 0.

		// We need to compute what would happen if we changed
		// our ask order on the instrument. A change in the ask order
		// means a change in the sampleAsset
		for j := 0; j < N; j++ {
			// Now for the new order
			sm := float64(sampleMatch[j]) / lpb
			expectedMatch := math.Max(sm-queue, 0)

			// We limit this by amount available in our portfolio
			expectedMatch = math.Min(expectedMatch, maxBase)
			// We are expected to lose expectedMatch base
			baseValueChange := -expectedMatch * basePrices[j]

			// We are expected to gain expectedMatch * price quote
			quoteValueChange := (expectedMatch - (exch.makerFee * expectedMatch)) * price * quotePrices[j]

			expectedLogReturn += math.Log((values[j] + baseValueChange + quoteValueChange) / value)
		}

		expectedLogReturn /= float64(N)


		//fmt.Println(expectedLogReturn, maxExpectedLogReturn)
		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = &Order{
				Price:    price,
				Quantity: 0., // Computed at the end
				Queue:    queue,
				Base:     base,
				Quote:    quote,
			}
		}
	}
	if maxOrder != nil {
		// Compute recommended order quantity based on price and queue
		// Recommended order quantity is the expected match at the level
		//fmt.Println("RETS", baseline, maxExpectedLogReturn)
		expectedMatch := 0.
		for i := 0; i < N; i++ {
			expectedMatch += math.Max(float64(sampleMatch[i]) / lpb - maxOrder.Queue, 0.)
		}
		maxOrder.Quantity = expectedMatch / float64(N)
	}

	return maxExpectedLogReturn, maxOrder
}
*/
