package modeling

import (
	"math"
	"math/rand"
)

type MarketModel interface {
	GetPrice(ID uint64) float64
	GetSamplePrices(ID uint64, time uint64, sampleSize int) []float64
	GetSampleMatchBid(ID uint64, time uint64, sampleSize int) []float64
	GetSampleMatchAsk(ID uint64, time uint64, sampleSize int) []float64
	SetPriceModel(ID uint64, model PriceModel)
	SetBuyTradeModel(ID uint64, model BuyTradeModel)
	SetSellTradeModel(ID uint64, model SellTradeModel)
}

type MapModel struct {
	buyTradeModels  map[uint64]BuyTradeModel
	sellTradeModels map[uint64]SellTradeModel
	priceModels     map[uint64]PriceModel
}

func NewMapModel() *MapModel {
	return &MapModel{
		buyTradeModels:  make(map[uint64]BuyTradeModel),
		sellTradeModels: make(map[uint64]SellTradeModel),
		priceModels:     make(map[uint64]PriceModel),
	}
}

func (m *MapModel) GetPrice(ID uint64) float64 {
	return m.priceModels[ID].GetPrice(ID)
}

func (m *MapModel) GetSamplePrices(ID uint64, time uint64, sampleSize int) []float64 {
	return m.priceModels[ID].GetSamplePrices(ID, time, sampleSize)
}

func (m *MapModel) GetSampleMatchBid(securityID uint64, time uint64, sampleSize int) []float64 {
	return m.sellTradeModels[securityID].GetSampleMatchBid(securityID, time, sampleSize)
}

func (m *MapModel) GetSampleMatchAsk(securityID uint64, time uint64, sampleSize int) []float64 {
	return m.buyTradeModels[securityID].GetSampleMatchAsk(securityID, time, sampleSize)
}

func (m *MapModel) SetPriceModel(ID uint64, model PriceModel) {
	m.priceModels[ID] = model
}

func (m *MapModel) SetBuyTradeModel(securityID uint64, model BuyTradeModel) {
	m.buyTradeModels[securityID] = model
}

func (m *MapModel) SetSellTradeModel(securityID uint64, model SellTradeModel) {
	m.sellTradeModels[securityID] = model
}

type Model interface {
	UpdatePrice(feedID uint64, tick uint64, price float64)
	UpdateTrade(feedID uint64, tick uint64, price, size float64)
	Progress(tick uint64)
}

type PriceModel interface {
	Model
	Frequency() uint64
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

func (m *ConstantPriceModel) UpdatePrice(_ uint64, _ uint64, _ float64) {

}

func (m *ConstantPriceModel) UpdateTrade(_ uint64, _ uint64, _ float64, _ float64) {

}

func (m *ConstantPriceModel) Progress(_ uint64) {

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
