package account

import (
	"container/list"
	"sync"
	"time"
)

type fill struct {
	price    float64
	taker    bool
	buy      bool
	time     int64
	bucketed map[int64]bool
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
					buyMoves, ok := sc.buyMakerMoves[securityID]
					if !ok {
						buyMoves = make(map[int64]float64)
						sc.buyMakerMoves[securityID] = buyMoves
					}
					buyMoves[bucket] = alpha*buyMoves[bucket] + (1-alpha)*move
				} else {
					sellMoves, ok := sc.sellMakerMoves[securityID]
					if !ok {
						sellMoves = make(map[int64]float64)
						sc.sellMakerMoves[securityID] = sellMoves
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
					buyMoves, ok := sc.buyTakerMoves[securityID]
					if !ok {
						buyMoves = make(map[int64]float64)
						sc.buyTakerMoves[securityID] = buyMoves
					}
					buyMoves[bucket] = alpha*buyMoves[bucket] + (1-alpha)*move
				} else {
					sellMoves, ok := sc.sellTakerMoves[securityID]
					if !ok {
						sellMoves = make(map[int64]float64)
						sc.sellTakerMoves[securityID] = sellMoves
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
	for _, v := range sc.buyMakerMoves {
		for i, mv := range v {
			buyMakerMoves[i] += mv
		}
	}
	for i := range sc.buyMakerMoves {
		buyMakerMoves[i] /= float64(len(sc.buyMakerMoves))
	}

	for _, v := range sc.buyTakerMoves {
		for i, mv := range v {
			buyTakerMoves[i] += mv
		}
	}
	for i := range sc.buyTakerMoves {
		buyTakerMoves[i] /= float64(len(sc.buyTakerMoves))
	}

	for _, v := range sc.sellMakerMoves {
		for i, mv := range v {
			sellMakerMoves[i] += mv
		}
	}
	for i := range sc.sellMakerMoves {
		sellMakerMoves[i] /= float64(len(sc.sellMakerMoves))
	}

	for _, v := range sc.sellTakerMoves {
		for i, mv := range v {
			sellTakerMoves[i] += mv
		}
	}
	for i := range sc.sellTakerMoves {
		sellTakerMoves[i] /= float64(len(sc.sellTakerMoves))
	}

	return buyMakerMoves, buyTakerMoves, sellMakerMoves, sellTakerMoves
}
