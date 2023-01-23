package modeling

import (
	"gitlab.com/alphaticks/tickstore-types/tickobjects"
	"math"
	"math/rand"
	"sync"
)

type Market interface {
	GetPrice(ID uint64) (float64, bool)
	GetPairPrice(base, quote uint32) (float64, bool)
}

type Model interface {
	Forward(selector int, tick, objectID uint64, object tickobjects.TickObject) error
	Backward()
	Ready() bool
	Frequency() uint64
	GetSelectors() []string
}

type MarketModel interface {
	Market
	Model
	GetSamplePairPrices(base uint32, quote uint32, time uint64, sampleSize int) []float64
	GetSamplePrices(ID uint64, time uint64, sampleSize int) []float64
	GetSampleMatchBid(ID uint64, time uint64, sampleSize int) []float64
	GetSampleMatchAsk(ID uint64, time uint64, sampleSize int) []float64
	SetPriceModel(ID uint64, model PriceModel)
	SetBuyTradeModel(ID uint64, model BuyTradeModel)
	SetSellTradeModel(ID uint64, model SellTradeModel)
}

type AllocationModel interface {
	Model
	GetAllocations(allocations, risks *sync.Map, margin, cost float64) map[uint64]float64
	GetAllocationDelta(ID uint64, allocations, risks *sync.Map, margin, cost float64) (float64, bool)
}

type LongShortModel interface {
	Model
	GetPenalty(fee float64) float64
	GetScore(ID uint64) (float64, bool)
	GetSelectors() []string
}

type LongModel interface {
	Model
	GetLongScore(ID uint64) float64
	//GetExitLongELR(ID uint64, fee, lambda float64) float64
}

type ShortModel interface {
	Model
	GetShortScore(ID uint64) float64
	//GetExitShortELR(ID uint64, fee, lambda float64) float64
}

type PriceModel interface {
	Model
	GetPrice(feedID uint64) (float64, bool)
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

type MarketMap struct {
	sync.RWMutex
	prices map[uint64]float64
}

func NewMarketMap() *MarketMap {
	return &MarketMap{
		prices: make(map[uint64]float64),
	}
}

func (m *MarketMap) SetPrice(ID uint64, p float64) {
	m.Lock()
	defer m.Unlock()
	m.prices[ID] = p
}

func (m *MarketMap) SetPairPrice(base, quote uint32, p float64) {
	m.Lock()
	defer m.Unlock()
	ID := uint64(base)<<32 | uint64(quote)
	m.prices[ID] = p
}

func (m *MarketMap) GetPrice(ID uint64) (float64, bool) {
	m.RLock()
	defer m.RUnlock()
	p, ok := m.prices[ID]
	return p, ok
}

func (m *MarketMap) GetPairPrice(base, quote uint32) (float64, bool) {
	m.RLock()
	defer m.RUnlock()
	ID := uint64(base)<<32 | uint64(quote)
	p, ok := m.prices[ID]
	return p, ok
}

type MarketAllocationModel struct {
	sync.RWMutex
	AllocationModel
	PriceModel
	*MarketMap
}

func NewMarketAllocationModel(al AllocationModel, price PriceModel) *MarketAllocationModel {
	return &MarketAllocationModel{
		AllocationModel: al,
		PriceModel:      price,
		MarketMap:       NewMarketMap(),
	}
}

func (m *MarketAllocationModel) GetPrice(ID uint64) (float64, bool) {
	if m.PriceModel != nil {
		p, ok := m.PriceModel.GetPrice(ID)
		if ok {
			return p, ok
		} else {
			return m.MarketMap.GetPrice(ID)
		}
	} else {
		return m.MarketMap.GetPrice(ID)
	}
}

func (m *MarketAllocationModel) GetPairPrice(base, quote uint32) (float64, bool) {
	ID := uint64(base)<<32 | uint64(quote)
	return m.GetPrice(ID)
}

type MarketLongShortModel struct {
	sync.RWMutex
	model LongShortModel
	*MarketMap
}

func NewMarketLongShortModel(model LongShortModel) *MarketLongShortModel {
	return &MarketLongShortModel{
		model: model,
	}
}

func (m *MarketLongShortModel) GetPenalty(fee float64) float64 {
	return m.model.GetPenalty(fee)
}

func (m *MarketLongShortModel) GetScore(ID uint64) (float64, bool) {
	return m.model.GetScore(ID)
}

type MapMarketModel struct {
	buyTradeModels  map[uint64]BuyTradeModel
	sellTradeModels map[uint64]SellTradeModel
	priceModels     map[uint64]PriceModel
	selectors       []string
}

func (m *MapMarketModel) Forward(selector int, tick, objectID uint64, object tickobjects.TickObject) error {
	return nil
}

func (m *MapMarketModel) Backward() {
	return
}

func (m *MapMarketModel) Ready() bool {
	return true
}

func (m *MapMarketModel) Frequency() uint64 {
	return 0
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
		return pm.GetPrice(ID)
	} else {
		return 0., false
	}
}

func (m *MapMarketModel) GetPairPrice(base uint32, quote uint32) (float64, bool) {
	ID := uint64(base)<<32 | uint64(quote)
	if pm, ok := m.priceModels[ID]; ok {
		return pm.GetPrice(ID)
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

type ConstantLongShortModel struct {
	penalty float64
	scores  map[uint64]float64
}

func NewConstantLongShortModel(penalty float64) *ConstantLongShortModel {
	return &ConstantLongShortModel{
		penalty: penalty,
		scores:  make(map[uint64]float64),
	}
}

func (m *ConstantLongShortModel) GetPenalty(_ float64) float64 {
	return m.penalty
}

func (m *ConstantLongShortModel) SetScore(ID uint64, score float64) {
	m.scores[ID] = score
}

func (m *ConstantLongShortModel) GetScore(ID uint64) (float64, bool) {
	s, ok := m.scores[ID]
	return s, ok
}

func (m *ConstantLongShortModel) GetSelectors() []string {
	return nil
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

func (m *ConstantPriceModel) Forward(selector int, tick uint64, objectID uint64, object tickobjects.TickObject) error {
	return nil
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
		m.samplePrices = make([]float64, sampleSize)
		for i := 0; i < sampleSize; i++ {
			m.samplePrices[i] = m.price
		}
	}
	return m.samplePrices
}

func (m *ConstantPriceModel) GetPrice(ID uint64) (float64, bool) {
	return m.price, true
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
		m.samplePrices = make([]float64, sampleSize)
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

func (m *ConstantTradeModel) Forward(selector int, tick uint64, objectID uint64, object tickobjects.TickObject) error {
	return nil
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
		m.sampleMatch = make([]float64, sampleSize)
		for i := 0; i < sampleSize; i++ {
			m.sampleMatch[i] = m.match
		}
	}
	return m.sampleMatch
}

func (m *ConstantTradeModel) GetSampleMatchBid(ID, time uint64, sampleSize int) []float64 {
	if m.sampleMatch == nil || len(m.sampleMatch) != sampleSize {
		m.sampleMatch = make([]float64, sampleSize)
		for i := 0; i < sampleSize; i++ {
			m.sampleMatch[i] = m.match
		}
	}
	return m.sampleMatch
}
