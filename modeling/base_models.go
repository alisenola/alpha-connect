package modeling

import (
	"gitlab.com/alphaticks/tickobjects"
	"math"
	"math/rand"
	"sync"
)

type Market interface {
	GetPrice(ID uint64) (float64, bool)
	GetPairPrice(base, quote uint32) (float64, bool)
}

type LongShortModel interface {
	Market
	GetLongScore(ID uint64, fee, lambda float64) float64
	GetShortScore(ID uint64, fee, lambda float64) float64
	SetLongModel(ID uint64, model LongModel)
	SetShortModel(ID uint64, model ShortModel)
	SetPrice(ID uint64, price float64)
	SetPairPrice(base, quote uint32, price float64)
}

type MarketModel interface {
	Market
	GetSamplePairPrices(base uint32, quote uint32, time uint64, sampleSize int) []float64
	GetSamplePrices(ID uint64, time uint64, sampleSize int) []float64
	GetSampleMatchBid(ID uint64, time uint64, sampleSize int) []float64
	GetSampleMatchAsk(ID uint64, time uint64, sampleSize int) []float64
	SetPriceModel(ID uint64, model PriceModel)
	SetBuyTradeModel(ID uint64, model BuyTradeModel)
	SetSellTradeModel(ID uint64, model SellTradeModel)
}

type MapMarketModel struct {
	buyTradeModels  map[uint64]BuyTradeModel
	sellTradeModels map[uint64]SellTradeModel
	priceModels     map[uint64]PriceModel
	selectors       []string
}

func NewMapMarketModel() *MapMarketModel {
	return &MapMarketModel{
		buyTradeModels:  make(map[uint64]BuyTradeModel),
		sellTradeModels: make(map[uint64]SellTradeModel),
		priceModels:     make(map[uint64]PriceModel),
	}
}

func (m *MapMarketModel) GetPrice(ID uint64) (float64, bool) {
	if pm, ok := m.priceModels[ID]; ok {
		return pm.GetPrice(ID), true
	} else {
		return 0., false
	}
}

func (m *MapMarketModel) GetPairPrice(base uint32, quote uint32) (float64, bool) {
	ID := uint64(base)<<32 | uint64(quote)
	if pm, ok := m.priceModels[ID]; ok {
		return pm.GetPrice(ID), true
	} else {
		return 0., false
	}
}

func (m *MapMarketModel) GetSamplePairPrices(base uint32, quote uint32, time uint64, sampleSize int) []float64 {
	ID := uint64(base)<<32 | uint64(quote)
	return m.priceModels[ID].GetSamplePrices(ID, time, sampleSize)
}

func (m *MapMarketModel) GetSamplePrices(ID uint64, time uint64, sampleSize int) []float64 {
	return m.priceModels[ID].GetSamplePrices(ID, time, sampleSize)
}

func (m *MapMarketModel) GetSampleMatchBid(securityID uint64, time uint64, sampleSize int) []float64 {
	return m.sellTradeModels[securityID].GetSampleMatchBid(securityID, time, sampleSize)
}

func (m *MapMarketModel) GetSampleMatchAsk(securityID uint64, time uint64, sampleSize int) []float64 {
	return m.buyTradeModels[securityID].GetSampleMatchAsk(securityID, time, sampleSize)
}

func (m *MapMarketModel) SetPriceModel(ID uint64, model PriceModel) {
	m.priceModels[ID] = model
}

func (m *MapMarketModel) SetBuyTradeModel(securityID uint64, model BuyTradeModel) {
	m.buyTradeModels[securityID] = model
}

func (m *MapMarketModel) SetSellTradeModel(securityID uint64, model SellTradeModel) {
	m.sellTradeModels[securityID] = model
}

func (m *MapMarketModel) PushSelectors(selectors []string) {
	m.selectors = append(m.selectors, selectors...)
}

func (m *MapMarketModel) GetSelectors() []string {
	return m.selectors
}

type MapLongShortModel struct {
	sync.RWMutex
	longModels  map[uint64]LongModel
	shortModels map[uint64]ShortModel
	selectors   []string
	prices      map[uint64]float64
}

func NewMapLongShortModel() *MapLongShortModel {
	return &MapLongShortModel{
		longModels:  make(map[uint64]LongModel),
		shortModels: make(map[uint64]ShortModel),
		prices:      make(map[uint64]float64),
	}
}

func (m *MapLongShortModel) SetLongModel(ID uint64, lm LongModel) {
	m.longModels[ID] = lm
}

func (m *MapLongShortModel) SetShortModel(ID uint64, sm ShortModel) {
	m.shortModels[ID] = sm
}

func (m *MapLongShortModel) SetPrice(ID uint64, p float64) {
	m.Lock()
	defer m.Unlock()
	m.prices[ID] = p
}

func (m *MapLongShortModel) SetPairPrice(base, quote uint32, p float64) {
	m.Lock()
	defer m.Unlock()
	ID := uint64(base)<<32 | uint64(quote)
	m.prices[ID] = p
}

func (m *MapLongShortModel) GetPrice(ID uint64) (float64, bool) {
	m.RLock()
	defer m.RUnlock()
	p, ok := m.prices[ID]
	return p, ok
}

func (m *MapLongShortModel) GetPairPrice(base, quote uint32) (float64, bool) {
	m.RLock()
	defer m.RUnlock()
	ID := uint64(base)<<32 | uint64(quote)
	p, ok := m.prices[ID]
	return p, ok
}

func (m *MapLongShortModel) GetLongScore(ID uint64, fee, lambda float64) float64 {
	return m.longModels[ID].GetLongScore(ID, fee, lambda)
}

func (m *MapLongShortModel) GetShortScore(ID uint64, fee, lambda float64) float64 {
	return m.shortModels[ID].GetShortScore(ID, fee, lambda)
}

/*
func (m *MapLongShortModel) GetEnterShortELR(ID uint64, fee, lambda float64) float64 {
	return m.shortModels[ID].GetEnterShortELR(ID, fee, lambda)
}

func (m *MapLongShortModel) GetExitShortELR(ID uint64, fee, lambda float64) float64 {
	return m.shortModels[ID].GetExitShortELR(ID, fee, lambda)
}
*/

/*
func (m *MapModel) Progress(time uint64) {
	for _, m := range m.priceModels {
		m.Progress(time)
	}
	for _, m := range m.sellTradeModels {
		m.Progress(time)
	}
	for _, m := range m.buyTradeModels {
		m.Progress(time)
	}
}
*/

type Model interface {
	Forward(selector int, tick, objectID uint64, object tickobjects.TickObject)
	Backward()
	Ready() bool
	Frequency() uint64
	GetSelectors() []string
}

type LongModel interface {
	Model
	GetLongScore(ID uint64, fee, lambda float64) float64
	//GetExitLongELR(ID uint64, fee, lambda float64) float64
}

type ShortModel interface {
	Model
	GetShortScore(ID uint64, fee, lambda float64) float64
	//GetExitShortELR(ID uint64, fee, lambda float64) float64
}

type PriceModel interface {
	Model
	GetPrice(feedID uint64) float64
	GetSamplePrices(feedID uint64, time uint64, sampleSize int) []float64
}

type SellTradeModel interface {
	Model
	GetSampleMatchBid(securityID uint64, time uint64, sampleSize int) []float64
}

type BuyTradeModel interface {
	Model
	GetSampleMatchAsk(securityID uint64, time uint64, sampleSize int) []float64
}

type ConstantPriceModel struct {
	price        float64
	samplePrices []float64
}

func NewConstantPriceModel(price float64) *ConstantPriceModel {
	return &ConstantPriceModel{
		price:        price,
		samplePrices: nil,
	}
}

func (m *ConstantPriceModel) SetSelectors(_ []tickobjects.TickObject) {

}

func (m *ConstantPriceModel) GetSelectors() []string {
	return nil
}

func (m *ConstantPriceModel) Forward(selector int, tick uint64, objectID uint64, object tickobjects.TickObject) {

}

func (m *ConstantPriceModel) Backward() {

}

func (m *ConstantPriceModel) Ready() bool {
	return true
}

func (m *ConstantPriceModel) Frequency() uint64 {
	return 0
}

func (m *ConstantPriceModel) GetSamplePrices(ID uint64, time uint64, sampleSize int) []float64 {
	if m.samplePrices == nil || len(m.samplePrices) != sampleSize {
		m.samplePrices = make([]float64, sampleSize, sampleSize)
		for i := 0; i < sampleSize; i++ {
			m.samplePrices[i] = m.price
		}
	}
	return m.samplePrices
}

func (m *ConstantPriceModel) GetPrice(ID uint64) float64 {
	return m.price
}

type GBMPriceModel struct {
	time         uint64
	price        float64
	freq         uint64
	samplePrices []float64
	sampleTime   uint64
}

func NewGBMPriceModel(price float64, freq uint64) *GBMPriceModel {
	return &GBMPriceModel{
		time:         0,
		price:        price,
		freq:         freq,
		samplePrices: nil,
		sampleTime:   0,
	}
}

func (m *GBMPriceModel) UpdatePrice(_ uint64, _ uint64, _ float64)            {}
func (m *GBMPriceModel) UpdateTrade(_ uint64, _ uint64, _ float64, _ float64) {}

func (m *GBMPriceModel) Progress(time uint64) {
	for m.time < time {
		m.price *= math.Exp(rand.NormFloat64())
		m.time += m.freq
	}
}

func (m *GBMPriceModel) Ready() bool {
	return true
}

func (m *GBMPriceModel) Frequency() uint64 {
	return m.freq
}

func (m *GBMPriceModel) GetSamplePrices(ID uint64, time uint64, sampleSize int) []float64 {
	if m.samplePrices == nil || len(m.samplePrices) != sampleSize || m.sampleTime != time {
		intervalLength := int((time - m.time) / m.freq)
		m.samplePrices = make([]float64, sampleSize, sampleSize)
		for i := 0; i < sampleSize; i++ {
			m.samplePrices[i] = m.price
			for j := 0; j < intervalLength; j++ {
				m.samplePrices[i] *= (rand.NormFloat64() / 10) + 1
			}
		}
		m.sampleTime = time
	}
	return m.samplePrices
}

func (m *GBMPriceModel) GetPrice(ID uint64) float64 {
	return m.price
}

type ConstantTradeModel struct {
	match       float64
	sampleMatch []float64
}

func NewConstantTradeModel(match float64) *ConstantTradeModel {
	return &ConstantTradeModel{
		match:       match,
		sampleMatch: nil,
	}
}

func (m *ConstantTradeModel) UpdatePrice(_ uint64, _ uint64, _ float64) {

}

func (m *ConstantTradeModel) UpdateTrade(_ uint64, _ uint64, _ float64, _ float64) {

}

func (m *ConstantTradeModel) Forward(selector int, tick uint64, objectID uint64, object tickobjects.TickObject) {

}

func (m *ConstantTradeModel) Backward() {

}

func (m *ConstantTradeModel) Frequency() uint64 {
	return 0
}

func (m *ConstantTradeModel) GetSelectors() []string {
	return nil
}

func (m *ConstantTradeModel) Ready() bool {
	return true
}

func (m *ConstantTradeModel) GetSampleMatchAsk(ID, time uint64, sampleSize int) []float64 {
	if m.sampleMatch == nil || len(m.sampleMatch) != sampleSize {
		m.sampleMatch = make([]float64, sampleSize, sampleSize)
		for i := 0; i < sampleSize; i++ {
			m.sampleMatch[i] = m.match
		}
	}
	return m.sampleMatch
}

func (m *ConstantTradeModel) GetSampleMatchBid(ID, time uint64, sampleSize int) []float64 {
	if m.sampleMatch == nil || len(m.sampleMatch) != sampleSize {
		m.sampleMatch = make([]float64, sampleSize, sampleSize)
		for i := 0; i < sampleSize; i++ {
			m.sampleMatch[i] = m.match
		}
	}
	return m.sampleMatch
}
