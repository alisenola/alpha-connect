package legacy

import (
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"math"
	"sync"
)

type BitmexPortfolio struct {
	sampleSize       int
	contracts        map[uint64]*Contract
	wallet           float64
	marginPriceModel PriceModel
}

type Contract struct {
	sync.RWMutex
	buyTradeModel  BuyTradeModel
	sellTradeModel SellTradeModel
	priceModel     PriceModel
	size           int
	margin         float64
	sampleSize     []float64
	sampleMargin   []float64
	sampleMatchBid []float64
	sampleMatchAsk []float64
	sampleTime     uint64
	multiplier     float64
	inverse        bool
	takerFee       float64
	makerFee       float64
	bidOrder       *Order
	askOrder       *Order
}

func (c *Contract) AddSampleValue(time uint64, values, marginPrices []float64) {
	c.Lock()
	if c.sampleTime != time {
		c.updateSampleMargin(time)
	}
	// TODO optimize, create sampleMarginChange for ex ?
	N := len(values)
	exitPrices := c.priceModel.GetSamplePrices(time, N)
	for i := 0; i < N; i++ {
		values[i] += c.sampleMargin[i] * marginPrices[i]
		if c.inverse {
			values[i] += (1. / exitPrices[i]) * c.multiplier * c.sampleSize[i] * marginPrices[i]
		} else {
			values[i] += exitPrices[i] * c.multiplier * c.sampleSize[i] * marginPrices[i]
		}
	}

	c.Unlock()
}

func (c *Contract) updateSampleMargin(time uint64) {
	N := len(c.sampleSize)
	for i := 0; i < N; i++ {
		c.sampleSize[i] = float64(c.size)
		c.sampleMargin[i] = c.margin
	}
	sampleMatchAsk := c.buyTradeModel.GetSampleMatchAsk(time, N)
	sampleMatchBid := c.sellTradeModel.GetSampleMatchBid(time, N)
	// Recompute sample margin and sample size
	if c.bidOrder != nil {
		var price float64
		if c.inverse {
			price = 1. / c.bidOrder.Price
		} else {
			price = c.bidOrder.Price
		}
		queue := c.bidOrder.Queue
		quantity := c.bidOrder.Quantity
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * c.multiplier

			// You post a bid order, you lose the contractMarginValue on match, and you lose the fees
			if c.inverse {
				c.sampleMargin[i] += contractMarginValue
				c.sampleSize[i] -= contractChange
			} else {
				c.sampleMargin[i] -= contractMarginValue
				c.sampleSize[i] += contractChange
			}
			c.sampleMargin[i] -= c.makerFee * contractMarginValue
		}
	}
	// Recompute sample asset
	if c.askOrder != nil {
		var price float64
		if c.inverse {
			price = 1. / c.askOrder.Price
		} else {
			price = c.askOrder.Price
		}
		queue := c.askOrder.Queue
		quantity := c.askOrder.Quantity
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)

			contractMarginValue := contractChange * price * c.multiplier
			// You post an ask order, you get the contractMarginValue, you lose the fees
			if c.inverse {
				c.sampleMargin[i] -= contractMarginValue
				c.sampleSize[i] += contractChange
			} else {
				c.sampleMargin[i] += contractMarginValue
				c.sampleSize[i] -= contractChange
			}
			c.sampleMargin[i] -= c.makerFee * contractMarginValue
		}
	}

	c.sampleMatchBid = sampleMatchBid
	c.sampleMatchAsk = sampleMatchAsk
	c.sampleTime = time
}

func NewBitmexPortfolio(securities []*models.Security, marginPriceModel PriceModel, priceModels map[uint64]PriceModel, buyTradeModels map[uint64]BuyTradeModel, sellTradeModels map[uint64]SellTradeModel, sampleSize int) (*BitmexPortfolio, error) {
	p := &BitmexPortfolio{
		sampleSize:       sampleSize,
		contracts:        make(map[uint64]*Contract),
		wallet:           0.,
		marginPriceModel: marginPriceModel,
	}

	for _, security := range securities {
		ID := security.SecurityID
		if _, ok := priceModels[ID]; !ok {
			return nil, fmt.Errorf("no price model for security %d", ID)
		}
		if _, ok := buyTradeModels[ID]; !ok {
			return nil, fmt.Errorf("no buy trade model for security %d", ID)
		}
		if _, ok := sellTradeModels[ID]; !ok {
			return nil, fmt.Errorf("no sell trade model for security %d", ID)
		}
		if security.GetMakerFee() == nil {
			return nil, fmt.Errorf("no maker fee for security %d", ID)
		}
		makerFee := security.GetMakerFee().Value
		if security.GetTakerFee() == nil {
			return nil, fmt.Errorf("no maker fee for security %d", ID)
		}
		takerFee := security.GetTakerFee().Value
		if security.GetMultiplier() == nil {
			return nil, fmt.Errorf("no multiplier for security %d", ID)
		}
		multiplier := security.GetMultiplier().Value
		switch security.SecurityType {
		case "CRPERP":
			c := &Contract{
				priceModel:     priceModels[ID],
				buyTradeModel:  buyTradeModels[ID],
				sellTradeModel: sellTradeModels[ID],
				size:           0,
				margin:         0.,
				sampleSize:     make([]float64, sampleSize, sampleSize),
				sampleMargin:   make([]float64, sampleSize, sampleSize),
				sampleMatchBid: make([]float64, sampleSize, sampleSize),
				sampleMatchAsk: make([]float64, sampleSize, sampleSize),
				sampleTime:     0,
				multiplier:     multiplier,
				inverse:        security.IsInverse,
				takerFee:       takerFee,
				makerFee:       makerFee,
				bidOrder:       nil,
				askOrder:       nil,
			}
			for i := 0; i < sampleSize; i++ {
				c.sampleSize[i] = 0.
				c.sampleMargin[i] = 0.
			}
			p.contracts[ID] = c
		default:
			return nil, fmt.Errorf("unknown contract type %s", security.SecurityType)
		}
	}

	return p, nil
}

func (p *BitmexPortfolio) Leverage() float64 {
	val := 0.
	for _, c := range p.contracts {
		c.RLock()
		// Entry price not defined if size = 0, division by 0 !
		if c.size == 0 {
			c.RUnlock()
			continue
		}
		exitPrice := c.priceModel.GetPrice()
		if c.inverse {
			val += math.Abs((1. / exitPrice) * c.multiplier * float64(c.size))
		} else {
			val += math.Abs(exitPrice * c.multiplier * float64(c.size))
		}
		c.RUnlock()
	}
	return val / p.wallet
}

func (p *BitmexPortfolio) AvailableMargin(leverage float64) float64 {
	CMV := 0.
	for _, c := range p.contracts {
		c.RLock()
		// Entry price not defined if size = 0, division by 0 !
		if c.size == 0 {
			c.RUnlock()
			continue
		}
		exitPrice := c.priceModel.GetPrice()
		if c.inverse {
			CMV += math.Abs((1. / exitPrice) * c.multiplier * float64(c.size))
		} else {
			CMV += math.Abs(exitPrice * c.multiplier * float64(c.size))
		}
		c.RUnlock()
	}

	return math.Max((leverage*p.wallet)-CMV, 0.)
}

// TODO add ExpectedLeverage

// The current portfolio value
func (p *BitmexPortfolio) Value() float64 {
	var val = p.wallet
	for _, c := range p.contracts {
		c.RLock()
		// Entry price not defined if size = 0, division by 0 !
		if c.size == 0 {
			c.RUnlock()
			continue
		}
		val += c.margin
		exitPrice := c.priceModel.GetPrice()
		if c.inverse {
			val += (1. / exitPrice) * c.multiplier * float64(c.size)
		} else {
			val += exitPrice * c.multiplier * float64(c.size)
		}
		c.RUnlock()
	}

	val *= p.marginPriceModel.GetPrice()

	return val
}

func (p *BitmexPortfolio) GetAsset(asset uint32) float64 {
	if asset == constants.BITCOIN.ID {
		val := p.wallet
		for _, c := range p.contracts {
			c.RLock()
			// Entry price not defined if size = 0, division by 0 !
			if c.size == 0 {
				c.RUnlock()
				continue
			}
			val += c.margin
			exitPrice := c.priceModel.GetPrice()
			if c.inverse {
				val += (1. / exitPrice) * c.multiplier * float64(c.size)
			} else {
				val += exitPrice * c.multiplier * float64(c.size)
			}
			c.RUnlock()
		}
		return val
	} else {
		return 0.
	}
}

func (p *BitmexPortfolio) GetContract(instrID uint64) int {
	c, ok := p.contracts[instrID]
	if !ok {
		return 0
	} else {
		c.RLock()
		defer c.RUnlock()
		if c.inverse {
			return -c.size
		} else {
			return c.size
		}
	}
}

func (p *BitmexPortfolio) GetPrice(instrID uint64) float64 {
	return p.contracts[instrID].priceModel.GetPrice()
}

// TODO set or update margin ?
func (p *BitmexPortfolio) SetWallet(wallet float64) {
	p.wallet = wallet
}

func (p *BitmexPortfolio) SetContract(instrID uint64, size int) {
	c := p.contracts[instrID]
	c.Lock()
	N := p.sampleSize
	if c.inverse {
		size = -size
	}
	for i := 0; i < N; i++ {
		c.sampleSize[i] -= float64(c.size)
		c.sampleSize[i] += float64(size)
	}

	c.size = size

	c.Unlock()
}

func (p *BitmexPortfolio) BuyContract(instrID uint64, size int, price float64) {
	c := p.contracts[instrID]

	if c.inverse {
		p.sellContract(instrID, size, 1./price)
	} else {
		p.buyContract(instrID, size, price)
	}
}

func (p *BitmexPortfolio) buyContract(instrID uint64, size int, price float64) {
	c := p.contracts[instrID]

	c.Lock()
	N := p.sampleSize

	if c.size < 0 {
		// We are closing our positions from c.size to c.size + size
		closedSize := size
		// we don't close 'size' if we go over 0 and re-open longs
		if -c.size < closedSize {
			closedSize = -c.size
		}
		contractMarginValue := float64(c.size) * price * c.multiplier
		PnL := c.margin + contractMarginValue
		realizedPnL := PnL * (float64(closedSize) / float64(-c.size))
		// Remove realized PnL from margin
		c.margin -= realizedPnL
		for i := 0; i < N; i++ {
			c.sampleMargin[i] -= realizedPnL
		}
		// Add it to wallet
		p.wallet += realizedPnL
	}

	deltaMarginValue := float64(size) * price * c.multiplier

	p.wallet -= c.makerFee * deltaMarginValue
	c.margin -= deltaMarginValue
	c.size += size

	// Update sample
	for i := 0; i < N; i++ {
		c.sampleMargin[i] -= deltaMarginValue
		c.sampleSize[i] += float64(size)
	}

	c.Unlock()
}

func (p *BitmexPortfolio) SellContract(instrID uint64, size int, price float64) {
	c := p.contracts[instrID]

	if c.inverse {
		p.buyContract(instrID, size, 1./price)
	} else {
		p.sellContract(instrID, size, price)
	}
}

func (p *BitmexPortfolio) sellContract(instrID uint64, size int, price float64) {
	c := p.contracts[instrID]

	c.Lock()

	N := p.sampleSize

	if c.size > 0 {
		// We are closing our positions from c.size to c.size - size
		closedSize := size
		// we don't close 'size' if we go under 0 and re-open shorts
		if c.size < closedSize {
			closedSize = c.size
		}
		contractMarginValue := float64(c.size) * price * c.multiplier
		PnL := c.margin + contractMarginValue
		realizedPnL := PnL * (float64(closedSize) / float64(c.size))

		// Remove realized PnL from margin
		c.margin -= realizedPnL
		for i := 0; i < N; i++ {
			c.sampleMargin[i] -= realizedPnL
		}
		// Add it to wallet
		p.wallet += realizedPnL
	}

	deltaMarginValue := float64(size) * price * c.multiplier

	p.wallet -= c.makerFee * deltaMarginValue
	c.margin += deltaMarginValue
	c.size -= size

	// Update sample
	for i := 0; i < N; i++ {
		c.sampleMargin[i] += deltaMarginValue
		c.sampleSize[i] -= float64(size)
	}

	c.Unlock()
}

func (p *BitmexPortfolio) SetMakerFee(instrID uint64, fee float64) {
	c := p.contracts[instrID]
	c.Lock()
	// We update sample margin by updating orders with old maker fee
	if c.bidOrder != nil {
		var price float64
		if c.inverse {
			price = 1. / c.bidOrder.Price
		} else {
			price = c.bidOrder.Price
		}
		queue := c.bidOrder.Queue
		quantity := c.bidOrder.Quantity
		sampleMatch := c.sampleMatchBid
		for i := 0; i < p.sampleSize; i++ {
			contractChange := math.Min(math.Max(sampleMatch[i]-queue, 0.), quantity)

			contractMarginValue := contractChange * price * c.multiplier
			// We update our sample margin
			c.sampleMargin[i] += contractMarginValue * c.makerFee
			c.sampleMargin[i] -= contractMarginValue * fee
		}
	}

	if c.askOrder != nil {
		var price float64
		if c.inverse {
			price = 1. / c.askOrder.Price
		} else {
			price = c.askOrder.Price
		}
		queue := c.askOrder.Queue
		quantity := c.askOrder.Quantity
		sampleMatch := c.sampleMatchAsk
		for i := 0; i < p.sampleSize; i++ {
			contractChange := math.Min(math.Max(sampleMatch[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * c.multiplier

			// We update our sample margin
			c.sampleMargin[i] += contractMarginValue * c.makerFee
			c.sampleMargin[i] -= contractMarginValue * fee
		}
	}
	c.makerFee = fee
	c.Unlock()
}

func (p *BitmexPortfolio) SetBidOrder(instrID uint64, newO Order) {

	// Update contract sample size

	c := p.contracts[instrID]
	c.Lock()

	sampleMatch := c.sampleMatchBid

	if c.bidOrder != nil {
		var price float64
		if c.inverse {
			price = 1. / c.bidOrder.Price
		} else {
			price = c.bidOrder.Price
		}
		queue := c.bidOrder.Queue
		quantity := c.bidOrder.Quantity
		for i := 0; i < p.sampleSize; i++ {
			contractChange := math.Min(math.Max(sampleMatch[i]-queue, 0.), quantity)

			contractMarginValue := contractChange * price * c.multiplier
			// You cancel a bid order, you get back the contractMarginValue you would have lost on match, and the fees
			if c.inverse {
				c.sampleMargin[i] -= contractMarginValue
				c.sampleSize[i] += contractChange
			} else {
				c.sampleMargin[i] += contractMarginValue
				c.sampleSize[i] -= contractChange
			}
			c.sampleMargin[i] += c.makerFee * contractMarginValue
		}
	}
	var price float64
	if c.inverse {
		price = 1. / newO.Price
	} else {
		price = newO.Price
	}
	for i := 0; i < p.sampleSize; i++ {
		contractChange := math.Min(math.Max(sampleMatch[i]-newO.Queue, 0.), newO.Quantity)

		contractMarginValue := contractChange * price * c.multiplier

		// You post a bid order, you lose the contractMarginValue on match, and you lose the fees
		if c.inverse {
			c.sampleMargin[i] += contractMarginValue
			c.sampleSize[i] -= contractChange
		} else {
			c.sampleMargin[i] -= contractMarginValue
			c.sampleSize[i] += contractChange
		}
		c.sampleMargin[i] -= c.makerFee * contractMarginValue
	}
	c.bidOrder = &newO
	c.Unlock()
}

func (p *BitmexPortfolio) SetAskOrder(instrID uint64, newO Order) {

	// Update contract sample size

	c := p.contracts[instrID]
	c.Lock()

	sampleMatchAsk := c.sampleMatchAsk
	if c.askOrder != nil {
		var price float64
		if c.inverse {
			price = 1. / c.askOrder.Price
		} else {
			price = c.askOrder.Price
		}
		queue := c.askOrder.Queue
		quantity := c.askOrder.Quantity
		for i := 0; i < p.sampleSize; i++ {
			contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)

			contractMarginValue := contractChange * price * c.multiplier
			// Update sample margin
			// You cancel ask order, you lose the contractMarginValue, you get back the fees
			if c.inverse {
				c.sampleMargin[i] += contractMarginValue
				c.sampleSize[i] -= contractChange
			} else {
				c.sampleMargin[i] -= contractMarginValue
				c.sampleSize[i] += contractChange
			}
			c.sampleMargin[i] += c.makerFee * contractMarginValue
		}
	}
	var price float64
	if c.inverse {
		price = 1. / newO.Price
	} else {
		price = newO.Price
	}
	for i := 0; i < p.sampleSize; i++ {
		contractChange := math.Min(math.Max(sampleMatchAsk[i]-newO.Queue, 0.), newO.Quantity)

		contractMarginValue := contractChange * price * c.multiplier
		// You post an ask order, you get the contractMarginValue, you lose the fees
		if c.inverse {
			c.sampleMargin[i] -= contractMarginValue
			c.sampleSize[i] += contractChange
		} else {
			c.sampleMargin[i] += contractMarginValue
			c.sampleSize[i] -= contractChange
		}
		c.sampleMargin[i] -= c.makerFee * contractMarginValue
	}
	c.askOrder = &newO
	c.Unlock()
}

func (p *BitmexPortfolio) CancelBidOrder(instrID uint64) {
	c := p.contracts[instrID]
	c.Lock()
	if c.bidOrder != nil {
		sampleMatchBid := c.sampleMatchBid
		var price float64
		if c.inverse {
			price = 1. / c.bidOrder.Price
		} else {
			price = c.bidOrder.Price
		}
		queue := c.bidOrder.Queue
		quantity := c.bidOrder.Quantity
		for i := 0; i < p.sampleSize; i++ {
			contractChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)

			contractMarginValue := contractChange * price * c.multiplier
			// You cancel a bid order, you get back the contractMarginValue you would have lost on match, and the fees
			if c.inverse {
				c.sampleMargin[i] -= contractMarginValue
				c.sampleSize[i] += contractChange
			} else {
				c.sampleMargin[i] += contractMarginValue
				c.sampleSize[i] -= contractChange
			}
			c.sampleMargin[i] += c.makerFee * contractMarginValue
		}
		c.bidOrder = nil
	}
	c.Unlock()
}

func (p *BitmexPortfolio) CancelAskOrder(instrID uint64) {
	c := p.contracts[instrID]
	c.Lock()
	// Update contract sample size
	if c.askOrder != nil {
		sampleMatchAsk := c.sampleMatchAsk
		var price float64
		if c.inverse {
			price = 1. / c.askOrder.Price
		} else {
			price = c.askOrder.Price
		}
		queue := c.askOrder.Queue
		quantity := c.askOrder.Quantity
		for i := 0; i < p.sampleSize; i++ {
			contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)

			contractMarginValue := contractChange * price * c.multiplier
			// You cancel ask order, you lose the contractMarginValue, you get back the fees
			if c.inverse {
				c.sampleMargin[i] += contractMarginValue
				c.sampleSize[i] -= contractChange
			} else {
				c.sampleMargin[i] -= contractMarginValue
				c.sampleSize[i] += contractChange
			}
			c.sampleMargin[i] += c.makerFee * contractMarginValue
		}
		c.askOrder = nil
	}
	c.Unlock()
}

func (p *BitmexPortfolio) AddSampleValues(time uint64, values []float64) {
	N := p.sampleSize
	marginPrices := p.marginPriceModel.GetSamplePrices(time, N)
	for i := 0; i < N; i++ {
		values[i] += p.wallet * marginPrices[i]
	}
	// TODO prevent negative PnL ?
	for _, c := range p.contracts {
		c.AddSampleValue(time, values, marginPrices)
	}
}

func (p *BitmexPortfolio) GetELROnCancelBid(instrID uint64, time uint64, values []float64, value float64) float64 {
	c := p.contracts[instrID]
	c.Lock()

	if c.sampleTime != time {
		c.updateSampleMargin(time)
	}

	N := p.sampleSize
	expectedLogReturn := 0.
	if c.bidOrder != nil {
		var price float64
		if c.inverse {
			price = 1. / c.bidOrder.Price
		} else {
			price = c.bidOrder.Price
		}
		queue := c.bidOrder.Queue
		quantity := c.bidOrder.Quantity
		// Need to know, for each sample, the value change if we remove it
		// so need to compute expected match for each sample, then compute the change in portfolio value
		// if this order is not executed
		marginPrices := p.marginPriceModel.GetSamplePrices(time, N)
		exitPrices := c.priceModel.GetSamplePrices(time, N)

		expectedLogReturn := 0.
		sampleMatchBid := c.sampleMatchBid
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * c.multiplier
			// You cancel a bid order, you get back the contractMarginValue you would have lost on match, and the fees
			// but you lose the contracts you would have won on match
			marginChange := c.makerFee * contractMarginValue
			if c.inverse {
				marginChange -= contractMarginValue
				marginChange += (1. / exitPrices[i]) * c.multiplier * contractChange
			} else {
				marginChange += contractMarginValue
				marginChange -= exitPrices[i] * c.multiplier * contractChange
			}
			expectedLogReturn += math.Log((values[i] + marginChange*marginPrices[i]) / value)
		}
	} else {
		for i := 0; i < N; i++ {
			expectedLogReturn += math.Log(values[i] / value)
		}
	}
	c.Unlock()

	return expectedLogReturn / float64(N)
}

func (p *BitmexPortfolio) GetELROnCancelAsk(instrID uint64, time uint64, values []float64, value float64) float64 {
	c := p.contracts[instrID]
	c.Lock()

	if c.sampleTime != time {
		c.updateSampleMargin(time)
	}

	N := p.sampleSize
	expectedLogReturn := 0.
	if c.askOrder != nil {
		var price float64
		if c.inverse {
			price = 1. / c.askOrder.Price
		} else {
			price = c.askOrder.Price
		}
		queue := c.askOrder.Queue
		quantity := c.askOrder.Quantity
		// Need to know, for each sample, the value change if we remove it
		// so need to compute expected match for each sample, then compute the change in portfolio value
		// if this order is not executed
		marginPrices := p.marginPriceModel.GetSamplePrices(time, N)
		exitPrices := c.priceModel.GetSamplePrices(time, N)

		expectedLogReturn := 0.
		sampleMatchAsk := c.sampleMatchAsk
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * c.multiplier
			// You cancel an ask order, you lose the contractMarginValue you would have gained on match
			// but you keep the contracts you would have lost on match, and the fees
			marginChange := c.makerFee * contractMarginValue
			if c.inverse {
				marginChange += contractMarginValue
				marginChange -= (1. / exitPrices[i]) * c.multiplier * contractChange
			} else {
				marginChange -= contractMarginValue
				marginChange += exitPrices[i] * c.multiplier * contractChange
			}
			expectedLogReturn += math.Log((values[i] + marginChange*marginPrices[i]) / value)
		}
	} else {
		for i := 0; i < N; i++ {
			expectedLogReturn += math.Log(values[i] / value)
		}
	}

	c.Unlock()
	return expectedLogReturn / float64(N)
}

func (p *BitmexPortfolio) GetELROnBidChange(instrID uint64, time uint64, values []float64, value float64, prices []float64, queues []float64, availableMargin float64) (float64, *Order) {
	N := p.sampleSize
	c := p.contracts[instrID]
	c.Lock()

	if c.sampleTime != time {
		c.updateSampleMargin(time)
	}

	sampleMatchBid := c.sampleMatchBid
	marginPrices := p.marginPriceModel.GetSamplePrices(time, p.sampleSize)
	exitPrices := c.priceModel.GetSamplePrices(time, p.sampleSize)

	var maxOrder *Order
	maxExpectedLogReturn := -999.

	// If we have a bid order, we remove it from values
	if c.bidOrder != nil {
		var price float64
		if c.inverse {
			price = 1. / c.bidOrder.Price
		} else {
			price = c.bidOrder.Price
		}
		queue := c.bidOrder.Queue
		quantity := c.bidOrder.Quantity

		expectedLogReturn := 0.
		for i := 0; i < N; i++ {
			expectedLogReturn += math.Log(values[i] / value)

			contractChange := math.Min(math.Max(sampleMatchBid[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * c.multiplier
			// You cancel a bid order, you get back the contractMarginValue you would have lost on match, and the fees
			// but you lose the contracts and their possible future margin value you would have won on match
			marginChange := c.makerFee * contractMarginValue
			if c.inverse {
				marginChange -= contractMarginValue
				marginChange += (1. / exitPrices[i]) * c.multiplier * contractChange
			} else {
				marginChange += contractMarginValue
				marginChange -= exitPrices[i] * c.multiplier * contractChange
			}
			values[i] -= marginChange * marginPrices[i]
		}
		expectedLogReturn /= float64(N)
		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = c.bidOrder
		}
	}

	// Compute expected log return without any bid order
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
		correctedAvailableMargin := availableMargin
		if c.inverse {
			price = 1. / prices[l]
			if c.size > 0 {
				// If we are short, we can go long the short size + the available margin size
				correctedAvailableMargin += math.Abs(price * c.multiplier * float64(c.size))
			}
		} else {
			price = prices[l]
			if c.size < 0 {
				// If we are short, we can go long the short size + the available margin size
				correctedAvailableMargin += math.Abs(price * c.multiplier * float64(c.size))
			}
		}
		queue := queues[l]
		expectedLogReturn := 0.

		for i := 0; i < N; i++ {
			contractMarginValue := math.Max(sampleMatchBid[i]-queue, 0.) * price * c.multiplier
			contractMarginValue = math.Min(contractMarginValue, correctedAvailableMargin)
			contractChange := contractMarginValue / (price * c.multiplier)

			// You create a bid order, you lose the fees and the contractMarginValue in your margin
			// but you get the contract and their possible future margin value depending on its exit price
			marginChange := -c.makerFee * contractMarginValue
			if c.inverse {
				marginChange += contractMarginValue
				marginChange -= (1. / exitPrices[i]) * c.multiplier * contractChange
			} else {
				marginChange -= contractMarginValue
				marginChange += exitPrices[i] * c.multiplier * contractChange
			}
			expectedLogReturn += math.Log((values[i] + marginChange*marginPrices[i]) / value)
		}
		expectedLogReturn /= float64(N)
		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = &Order{
				Price:    prices[l],
				Quantity: 0., // Computed at the end
				Queue:    queue,
			}
		}
	}

	if maxOrder != nil && maxOrder != c.bidOrder {
		// Compute recommended order quantity based on price and queue
		// Recommended order quantity is the expected match at the level
		expectedMatch := 0.
		for i := 0; i < N; i++ {
			expectedMatch += math.Max(sampleMatchBid[i]-maxOrder.Queue, 0.)
		}
		expectedMatch /= float64(N)

		correctedAvailableMargin := availableMargin
		var price float64
		if c.inverse {
			price = 1. / maxOrder.Price
			if c.size > 0 {
				// If we are short, we can go long the short size + the available margin size
				correctedAvailableMargin += math.Abs(price * c.multiplier * float64(c.size))
			}
		} else {
			price = maxOrder.Price
			if c.size < 0 {
				// If we are short, we can go long the short size + the available margin size
				correctedAvailableMargin += math.Abs(price * c.multiplier * float64(c.size))
			}
		}

		maxOrder.Quantity = math.Min(expectedMatch, correctedAvailableMargin/(price*c.multiplier))
	}
	c.Unlock()

	return maxExpectedLogReturn, maxOrder
}

func (p *BitmexPortfolio) GetELROnAskChange(instrID uint64, time uint64, values []float64, value float64, prices []float64, queues []float64, availableMargin float64) (float64, *Order) {
	N := p.sampleSize
	c := p.contracts[instrID]
	c.Lock()

	if c.sampleTime != time {
		c.updateSampleMargin(time)
	}

	sampleMatchAsk := c.buyTradeModel.GetSampleMatchAsk(time, N)
	marginPrices := p.marginPriceModel.GetSamplePrices(time, p.sampleSize)
	exitPrices := c.priceModel.GetSamplePrices(time, p.sampleSize)

	var maxOrder *Order
	maxExpectedLogReturn := -999.

	// If we have a ask order, we remove it from values
	if c.askOrder != nil {
		var price float64
		if c.inverse {
			price = 1. / c.askOrder.Price
		} else {
			price = c.askOrder.Price
		}
		queue := c.askOrder.Queue
		quantity := c.askOrder.Quantity
		for i := 0; i < N; i++ {
			contractChange := math.Min(math.Max(sampleMatchAsk[i]-queue, 0.), quantity)
			contractMarginValue := contractChange * price * c.multiplier
			// You cancel an ask order, you lose the contractMarginValue you would have gained on match
			// but you keep the contracts you would have lost on match, and the fees
			marginChange := c.makerFee * contractMarginValue
			if c.inverse {
				marginChange += contractMarginValue
				marginChange -= (1. / exitPrices[i]) * c.multiplier * contractChange
			} else {
				marginChange -= contractMarginValue
				marginChange += exitPrices[i] * c.multiplier * contractChange
			}
			values[i] -= marginChange * marginPrices[i]
		}
	}

	// Compute expected log return without any ask order
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
		var price float64
		if c.inverse {
			price = 1. / prices[l]
			if c.size < 0 {
				// If we are long, we can go short the long size + the available margin size
				correctedAvailableMargin += math.Abs(price * c.multiplier * float64(c.size))
			}
		} else {
			price = prices[l]
			if c.size > 0 {
				// If we are long, we can go short the long size + the available margin size
				correctedAvailableMargin += math.Abs(price * c.multiplier * float64(c.size))
			}
		}

		queue := queues[l]
		expectedLogReturn := 0.
		for i := 0; i < N; i++ {
			contractMarginValue := math.Max(sampleMatchAsk[i]-queue, 0.) * price * c.multiplier
			contractMarginValue = math.Min(contractMarginValue, correctedAvailableMargin)
			contractChange := contractMarginValue / (price * c.multiplier)
			// You create an ask order, you lose the fees and the contract and their possible future margin value depending on its exit price
			// but you get the margin value of the contract at the order price
			marginChange := -c.makerFee * contractMarginValue
			if c.inverse {
				marginChange -= contractMarginValue
				marginChange += (1. / exitPrices[i]) * c.multiplier * contractChange
			} else {
				marginChange += contractMarginValue
				marginChange -= exitPrices[i] * c.multiplier * contractChange
			}
			expectedLogReturn += math.Log((values[i] + marginChange*marginPrices[i]) / value)
		}
		expectedLogReturn /= float64(N)

		if expectedLogReturn > maxExpectedLogReturn {
			maxExpectedLogReturn = expectedLogReturn
			maxOrder = &Order{
				Price:    prices[l],
				Quantity: 0., // Computed at the end
				Queue:    queue,
			}
		}
	}

	if maxOrder != nil && maxOrder != c.askOrder {
		// Compute recommended order quantity based on price and queue
		// Recommended order quantity is the expected match at the level
		expectedMatch := 0.
		for i := 0; i < N; i++ {
			expectedMatch += math.Max(sampleMatchAsk[i]-maxOrder.Queue, 0.)
		}
		expectedMatch /= float64(N)

		correctedAvailableMargin := availableMargin
		var price float64
		if c.inverse {
			price = 1. / maxOrder.Price
			if c.size < 0 {
				// If we are short, we can go long the short size + the available margin size
				correctedAvailableMargin += math.Abs(price * c.multiplier * float64(c.size))
			}
		} else {
			price = maxOrder.Price
			if c.size > 0 {
				// If we are short, we can go long the short size + the available margin size
				correctedAvailableMargin += math.Abs(price * c.multiplier * float64(c.size))
			}
		}

		maxOrder.Quantity = math.Min(expectedMatch, correctedAvailableMargin/(price*c.multiplier))
	}
	c.Unlock()
	return maxExpectedLogReturn, maxOrder
}
