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
					sc.takerMoves[bucket] += price / f.price
				} else {
					sc.makerMoves[bucket] += price / f.price
				}
				// TODO
			}
		}
	}
}

func (sc *FillCollector) GetMoveAfterFill(securityID uint64) float64 {
	// TODO
	fmt.Println(sc.makerFills[securityID])
	fmt.Println(sc.takerFills[securityID])
	return 0
}
