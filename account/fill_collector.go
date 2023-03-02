package account

import (
	"container/list"
	"sync"
	"time"
)

type fill struct {
	price float64
	taker bool
	time  int64
}

type FillCollector struct {
	sync.RWMutex
	cutoff     int64 // drop trades older than this
	takerFills map[uint64]*list.List
	makerFills map[uint64]*list.List
	makerMoves map[int64]float64
	takerMoves map[int64]float64
}

func NewFillCollector(cutoff int64) *FillCollector {
	return &FillCollector{
		cutoff:     cutoff,
		takerFills: make(map[uint64]*list.List),
		makerFills: make(map[uint64]*list.List),
		makerMoves: make(map[int64]float64),
		takerMoves: make(map[int64]float64),
	}
}

func (sc *FillCollector) AddFill(securityID uint64, price float64, taker bool, time int64) {
	sc.Lock()
	defer sc.Unlock()
	if taker {
		m, ok := sc.takerFills[securityID]
		if !ok {
			m = list.New()
			sc.takerFills[securityID] = m
		}

		m.PushFront(&fill{price: price, taker: taker, time: time})
	} else {
		m, ok := sc.makerFills[securityID]
		if !ok {
			m = list.New()
			sc.makerFills[securityID] = m
		}

		m.PushFront(&fill{price: price, taker: taker, time: time})
	}
}

func (sc *FillCollector) Collect(securityID uint64, price float64) {
	sc.Lock()
	defer sc.Unlock()
	ts := time.Now().UnixMilli()
	m, ok := sc.makerFills[securityID]
	if ok {
		for e := m.Front(); e != nil; e = e.Next() {
			f := e.Value.(*fill)
			delta := ts - f.time
			if delta > sc.cutoff {
				m.Remove(e)
			} else {
				bucket := delta / 1000 // 1s buckets
				move := price/f.price - 1
				sc.makerMoves[bucket] = 0.99*sc.makerMoves[bucket] + 0.01*move
			}
		}
	}
	m, ok = sc.takerFills[securityID]
	if ok {
		for e := m.Front(); e != nil; e = e.Next() {
			f := e.Value.(*fill)
			delta := ts - f.time
			if delta > sc.cutoff {
				m.Remove(e)
			} else {
				bucket := delta / 1000 // 1s buckets
				move := price/f.price - 1
				sc.takerMoves[bucket] = 0.99*sc.takerMoves[bucket] + 0.01*move
			}
		}
	}
}

func (sc *FillCollector) GetMoveAfterFill() ([]float64, []float64) {
	// TODO
	sc.RLock()
	defer sc.RUnlock()
	makerMoves := make([]float64, 0, 10)
	takerMoves := make([]float64, 0, 10)
	for k, v := range sc.makerMoves {
		makerMoves[k] = v
	}
	for k, v := range sc.takerMoves {
		takerMoves[k] = v
	}
	return makerMoves, takerMoves
}
