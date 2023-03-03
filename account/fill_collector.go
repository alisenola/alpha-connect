package account

import (
	"container/list"
	"sync"
	"time"
)

type fill struct {
	securityID uint64
	price      float64
	taker      bool
	buy        bool
	time       int64
	bucketed   map[int64]bool
}

type FillCollector struct {
	sync.RWMutex
	cutoff         int64 // drop trades older than this
	takerFills     map[uint64]*list.List
	makerFills     map[uint64]*list.List
	buyMakerMoves  map[uint64]map[int64]float64
	buyTakerMoves  map[uint64]map[int64]float64
	sellMakerMoves map[uint64]map[int64]float64
	sellTakerMoves map[uint64]map[int64]float64
}

func NewFillCollector(cutoff int64, securityIDs []uint64) *FillCollector {
	fc := &FillCollector{
		cutoff:         cutoff,
		takerFills:     make(map[uint64]*list.List),
		makerFills:     make(map[uint64]*list.List),
		buyMakerMoves:  make(map[uint64]map[int64]float64),
		buyTakerMoves:  make(map[uint64]map[int64]float64),
		sellMakerMoves: make(map[uint64]map[int64]float64),
		sellTakerMoves: make(map[uint64]map[int64]float64),
	}
	for _, sec := range securityIDs {
		fc.takerFills[sec] = list.New()
		fc.makerFills[sec] = list.New()
		fc.buyMakerMoves[sec] = make(map[int64]float64)
		fc.buyTakerMoves[sec] = make(map[int64]float64)
		fc.sellMakerMoves[sec] = make(map[int64]float64)
		fc.sellTakerMoves[sec] = make(map[int64]float64)
	}
	return fc
}

func (sc *FillCollector) AddFill(securityID uint64, price float64, buy, taker bool, time int64) {
	if taker {
		m := sc.takerFills[securityID]
		m.PushFront(&fill{price: price, buy: buy, taker: taker, time: time, bucketed: make(map[int64]bool)})
	} else {
		m := sc.makerFills[securityID]
		m.PushFront(&fill{price: price, buy: buy, taker: taker, time: time, bucketed: make(map[int64]bool)})
	}
}

func (sc *FillCollector) Collect(securityID uint64, price float64) {
	sc.Lock()
	defer sc.Unlock()
	ts := time.Now().UnixMilli()
	alpha := 0.995
	m := sc.makerFills[securityID]
	for e := m.Front(); e != nil; e = e.Next() {
		f := e.Value.(*fill)
		delta := ts - f.time
		if delta >= sc.cutoff {
			m.Remove(e)
		} else {
			bucket := delta / 1000 // 1s buckets
			if !f.bucketed[bucket] {
				var move float64
				if f.buy {
					// price down, move is negative
					move = price/f.price - 1
				} else {
					move = f.price/price - 1
				}
				if f.buy {
					buyMoves, ok := sc.buyMakerMoves[f.securityID]
					if !ok {
						buyMoves = make(map[int64]float64)
						sc.buyMakerMoves[f.securityID] = buyMoves
					}
					buyMoves[bucket] = alpha*buyMoves[bucket] + (1-alpha)*move
				} else {
					sellMoves, ok := sc.sellMakerMoves[f.securityID]
					if !ok {
						sellMoves = make(map[int64]float64)
						sc.sellMakerMoves[f.securityID] = sellMoves
					}
					sellMoves[bucket] = alpha*sellMoves[bucket] + (1-alpha)*move
				}
				f.bucketed[bucket] = true
			}
		}
	}
	m = sc.takerFills[securityID]
	for e := m.Front(); e != nil; e = e.Next() {
		f := e.Value.(*fill)
		delta := ts - f.time
		if delta >= sc.cutoff {
			m.Remove(e)
		} else {
			bucket := delta / 1000 // 1s buckets
			if !f.bucketed[bucket] {
				var move float64
				if f.buy {
					move = price/f.price - 1
				} else {
					move = f.price/price - 1
				}
				if f.buy {
					buyMoves, ok := sc.buyTakerMoves[f.securityID]
					if !ok {
						buyMoves = make(map[int64]float64)
						sc.buyTakerMoves[f.securityID] = buyMoves
					}
					buyMoves[bucket] = alpha*buyMoves[bucket] + (1-alpha)*move
				} else {
					sellMoves, ok := sc.sellTakerMoves[f.securityID]
					if !ok {
						sellMoves = make(map[int64]float64)
						sc.sellTakerMoves[f.securityID] = sellMoves
					}
					sellMoves[bucket] = alpha*sellMoves[bucket] + (1-alpha)*move
				}
				f.bucketed[bucket] = true
			}
		}
	}
}

func (sc *FillCollector) GetMoveAfterFill() ([]float64, []float64, []float64, []float64) {
	// TODO
	sc.RLock()
	defer sc.RUnlock()
	buyMakerMoves := make([]float64, 10)
	buyTakerMoves := make([]float64, 10)
	sellMakerMoves := make([]float64, 10)
	sellTakerMoves := make([]float64, 10)
	for k, v := range sc.buyMakerMoves {
		sum := 0.0
		for _, mv := range v {
			sum += mv
		}
		buyMakerMoves[k] = sum / float64(len(v))
	}
	for k, v := range sc.buyTakerMoves {
		sum := 0.0
		for _, mv := range v {
			sum += mv
		}
		buyTakerMoves[k] = sum / float64(len(v))
	}
	for k, v := range sc.sellMakerMoves {
		sum := 0.0
		for _, mv := range v {
			sum += mv
		}
		sellMakerMoves[k] = sum / float64(len(v))
	}
	for k, v := range sc.sellTakerMoves {
		sum := 0.0
		for _, mv := range v {
			sum += mv
		}
		sellTakerMoves[k] = sum / float64(len(v))
	}
	return buyMakerMoves, buyTakerMoves, sellMakerMoves, sellTakerMoves
}
