package account

import (
	"container/list"
	"fmt"
	"time"
)

type fill struct {
	price float64
	taker bool
	time  int64
}

type FillCollector struct {
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
	m, ok := sc.takerFills[securityID]
	if !ok {
		m = list.New()
		sc.takerFills[securityID] = m
	}

	m.PushBack(&fill{price: price, taker: taker, time: time})
}

func (sc *FillCollector) Collect(securityID uint64, price float64) {
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
				if f.taker {
					sc.takerMoves[bucket] = 0.9*sc.takerMoves[bucket] + 0.1*(price/f.price)
				} else {
					sc.makerMoves[bucket] = 0.9*sc.makerMoves[bucket] + 0.1*(price/f.price)
				}
				// TODO
			}
		}
	}
}

func (sc *FillCollector) GetMoveAfterFill() float64 {
	// TODO
	for k, v := range sc.makerMoves {
		fmt.Println(k, v)
	}
	for k, v := range sc.takerMoves {
		fmt.Println(k, v)
	}
	return 0
}
