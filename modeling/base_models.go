package modeling

import (
	"gitlab.com/tachikoma.ai/tickobjects"
	"math"
	"math/rand"
)

type LongShortModel interface {
	GetEnterLongScore(ID uint64) float64
	GetExitLongScore(ID uint64) float64
	GetEnterShortScore(ID uint64) float64
	GetExitShortScore(ID uint64) float64
	SetLongModel(ID uint64, model LongModel)
	SetShortModel(ID uint64, model ShortModel)
}

type MarketModel interface {
	GetPrice(ID uint64) float64
	GetPairPrice(base uint32, quote uint32) float64
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

func (m *MapMarketModel) GetPrice(ID uint64) float64 {
	return m.priceModels[ID].GetPrice(ID)
}

func (m *MapMarketModel) GetPairPrice(base uint32, quote uint32) float64 {
	ID := uint64(base)<<32 | uint64(quote)
	return m.priceModels[ID].GetPrice(ID)
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
	longModels  map[uint64]LongModel
	shortModels map[uint64]ShortModel
	selectors   []string
}

func NewMapLongShortModel() *MapLongShortModel {
	return &MapLongShortModel{
		longModels:  make(map[uint64]LongModel),
		shortModels: make(map[uint64]ShortModel),
	}
}

func (m *MapLongShortModel) SetLongModel(ID uint64, lm LongModel) {
	m.longModels[ID] = lm
}

func (m *MapLongShortModel) SetShortModel(ID uint64, sm ShortModel) {
	m.shortModels[ID] = sm
}

func (m *MapLongShortModel) GetEnterLongScore(ID uint64) float64 {
	return m.longModels[ID].GetEnterLongScore(ID)
}

func (m *MapLongShortModel) GetExitLongScore(ID uint64) float64 {
	return m.longModels[ID].GetExitLongScore(ID)
}

func (m *MapLongShortModel) GetEnterShortScore(ID uint64) float64 {
	return m.shortModels[ID].GetEnterShortScore(ID)
}

func (m *MapLongShortModel) GetExitShortScore(ID uint64) float64 {
	return m.shortModels[ID].GetExitShortScore(ID)
}

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
	GetEnterLongScore(feedID uint64) float64
	GetExitLongScore(feedID uint64) float64
}

type ShortModel interface {
	Model
	GetEnterShortScore(feedID uint64) float64
	GetExitShortScore(feedID uint64) float64
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

func (m *ConstantPriceModel) Forward(_ uint64, _ int) {

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

func (m *ConstantTradeModel) Progress(_ uint64) {

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
