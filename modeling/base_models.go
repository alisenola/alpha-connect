package modeling

import (
	"math"
	"math/rand"
)

type MarketModel interface {
	GetSecurityPrice(securityID uint64) float64
	GetAssetPrice(assetID uint32) float64
	GetSampleAssetPrices(assetID uint32, time uint64, sampleSize int) []float64
	GetSampleSecurityPrices(securityID uint64, time uint64, sampleSize int) []float64
	GetSampleMatchBid(securityID uint64, time uint64, sampleSize int) []float64
	GetSampleMatchAsk(securityID uint64, time uint64, sampleSize int) []float64
	SetSecurityPriceModel(securityID uint64, model SecurityPriceModel)
	SetBuyTradeModel(securityID uint64, model BuyTradeModel)
	SetSellTradeModel(securityID uint64, model SellTradeModel)
	SetAssetPriceModel(assetID uint32, model AssetPriceModel)
}

type MapModel struct {
	buyTradeModels      map[uint64]BuyTradeModel
	sellTradeModels     map[uint64]SellTradeModel
	securityPriceModels map[uint64]SecurityPriceModel
	assetPriceModels    map[uint32]AssetPriceModel
}

func NewMapModel() *MapModel {
	return &MapModel{
		buyTradeModels:      make(map[uint64]BuyTradeModel),
		sellTradeModels:     make(map[uint64]SellTradeModel),
		securityPriceModels: make(map[uint64]SecurityPriceModel),
		assetPriceModels:    make(map[uint32]AssetPriceModel),
	}
}

func (m *MapModel) GetSecurityPrice(securityID uint64) float64 {
	return m.securityPriceModels[securityID].GetSecurityPrice(securityID)
}

func (m *MapModel) GetAssetPrice(assetID uint32) float64 {
	return m.assetPriceModels[assetID].GetAssetPrice(assetID)
}

func (m *MapModel) GetSampleAssetPrices(assetID uint32, time uint64, sampleSize int) []float64 {
	return m.assetPriceModels[assetID].GetSampleAssetPrices(assetID, time, sampleSize)
}

func (m *MapModel) GetSampleSecurityPrices(securityID uint64, time uint64, sampleSize int) []float64 {
	return m.securityPriceModels[securityID].GetSampleSecurityPrices(securityID, time, sampleSize)
}

func (m *MapModel) GetSampleMatchBid(securityID uint64, time uint64, sampleSize int) []float64 {
	return m.sellTradeModels[securityID].GetSampleMatchBid(securityID, time, sampleSize)
}

func (m *MapModel) GetSampleMatchAsk(securityID uint64, time uint64, sampleSize int) []float64 {
	return m.buyTradeModels[securityID].GetSampleMatchAsk(securityID, time, sampleSize)
}

func (m *MapModel) SetSecurityPriceModel(securityID uint64, model SecurityPriceModel) {
	m.securityPriceModels[securityID] = model
}

func (m *MapModel) SetAssetPriceModel(assetID uint32, model AssetPriceModel) {
	m.assetPriceModels[assetID] = model
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

type AssetPriceModel interface {
	Model
	GetAssetPrice(assetID uint32) float64
	GetSampleAssetPrices(assetID uint32, time uint64, sampleSize int) []float64
	Frequency() uint64
}

type SecurityPriceModel interface {
	Model
	GetSecurityPrice(securityID uint64) float64
	GetSampleSecurityPrices(securityID uint64, time uint64, sampleSize int) []float64
	Frequency() uint64
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

func (m *ConstantPriceModel) GetAssetPrice(assetID uint32) float64 {
	return m.price
}

func (m *ConstantPriceModel) GetSampleAssetPrices(time uint64, sampleSize int) []float64 {
	if m.samplePrices == nil || len(m.samplePrices) != sampleSize {
		m.samplePrices = make([]float64, sampleSize, sampleSize)
		for i := 0; i < sampleSize; i++ {
			m.samplePrices[i] = m.price
		}
	}
	return m.samplePrices
}

func (m *ConstantPriceModel) GetSecurityPrice(securityID uint64) float64 {
	return m.price
}

func (m *ConstantPriceModel) GetSampleSecurityPrices(time uint64, sampleSize int) []float64 {
	if m.samplePrices == nil || len(m.samplePrices) != sampleSize {
		m.samplePrices = make([]float64, sampleSize, sampleSize)
		for i := 0; i < sampleSize; i++ {
			m.samplePrices[i] = m.price
		}
	}
	return m.samplePrices
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

func (m *GBMPriceModel) GetSampleAssetPrices(assetID uint32, time uint64, sampleSize int) []float64 {
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

func (m *GBMPriceModel) GetSampleSecurityPrices(securityID uint64, time uint64, sampleSize int) []float64 {
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

func (m *GBMPriceModel) GetAssetPrice(assetID uint32) float64 {
	return m.price
}

func (m *GBMPriceModel) GetSecurityPrice(securityID uint64) float64 {
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
