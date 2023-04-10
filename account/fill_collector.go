package account

import (
	"container/list"
	"math"
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
	alpha          float64
	volumeTau      time.Duration
	cutoff         int64 // drop trades older than this
	takerFills     map[uint64]*list.List
	makerFills     map[uint64]*list.List
	buyMakerMoves  map[uint64]map[int64][2]float64
	buyTakerMoves  map[uint64]map[int64][2]float64
	sellMakerMoves map[uint64]map[int64][2]float64
	sellTakerMoves map[uint64]map[int64][2]float64
	volume         map[uint64][2]float64
}

func NewFillCollector(cutoff int64, alpha float64, volumeTau time.Duration, securityIDs []uint64) *FillCollector {
	fc := &FillCollector{
		alpha:          alpha,
		volumeTau:      volumeTau,
		cutoff:         cutoff,
		takerFills:     make(map[uint64]*list.List),
		makerFills:     make(map[uint64]*list.List),
		buyMakerMoves:  make(map[uint64]map[int64][2]float64),
		buyTakerMoves:  make(map[uint64]map[int64][2]float64),
		sellMakerMoves: make(map[uint64]map[int64][2]float64),
		sellTakerMoves: make(map[uint64]map[int64][2]float64),
		volume:         make(map[uint64][2]float64),
	}
	for _, sec := range securityIDs {
		fc.takerFills[sec] = list.New()
		fc.makerFills[sec] = list.New()
		fc.buyMakerMoves[sec] = make(map[int64][2]float64)
		fc.buyTakerMoves[sec] = make(map[int64][2]float64)
		fc.sellMakerMoves[sec] = make(map[int64][2]float64)
		fc.sellTakerMoves[sec] = make(map[int64][2]float64)
	}
	return fc
}

func (sc *FillCollector) AddFill(securityID uint64, price, quantity float64, buy, taker bool, ts time.Time) {
	sc.Lock()
	defer sc.Unlock()
	milli := ts.UnixMilli()
	if taker {
		m := sc.takerFills[securityID]
		m.PushFront(&fill{price: price, buy: buy, taker: taker, time: milli, bucketed: make(map[int64]bool)})
	} else {
		m := sc.makerFills[securityID]
		m.PushFront(&fill{price: price, buy: buy, taker: taker, time: milli, bucketed: make(map[int64]bool)})
	}
	// TODO use half life to compute alpha

	delta := float64(milli) - sc.volume[securityID][1]
	w := math.Exp(-delta / float64(sc.volumeTau.Milliseconds()))

	// Decay volume + add new quantity
	sc.volume[securityID] = [2]float64{w*sc.volume[securityID][0] + price*quantity, float64(milli)}
}

func (sc *FillCollector) Collect(securityID uint64, price float64) {
	sc.Lock()
	defer sc.Unlock()
	ts := time.Now().UnixMilli()
	m := sc.makerFills[securityID]
	for e := m.Front(); e != nil; e = e.Next() {
		if e.Value == nil {
			m.Remove(e)
			continue
		}
		f := e.Value.(*fill)
		delta := ts - f.time
		if delta >= sc.cutoff {
			m.Remove(e)
		} else {
			bucket := delta / 1000 // 1s buckets
			if !f.bucketed[bucket] {
				if f.buy {
					move := price/f.price - 1
					buyMoves, ok := sc.buyMakerMoves[securityID]
					if !ok {
						buyMoves = make(map[int64][2]float64)
						sc.buyMakerMoves[securityID] = buyMoves
					}
					buyMoves[bucket] = [2]float64{sc.alpha*buyMoves[bucket][0] + (1-sc.alpha)*move, float64(ts)}
				} else {
					move := f.price/price - 1
					sellMoves, ok := sc.sellMakerMoves[securityID]
					if !ok {
						sellMoves = make(map[int64][2]float64)
						sc.sellMakerMoves[securityID] = sellMoves
					}
					sellMoves[bucket] = [2]float64{sc.alpha*sellMoves[bucket][0] + (1-sc.alpha)*move, float64(ts)}
				}
				f.bucketed[bucket] = true
			}
		}
	}
	m = sc.takerFills[securityID]
	for e := m.Front(); e != nil; e = e.Next() {
		if e.Value == nil {
			m.Remove(e)
			continue
		}
		f := e.Value.(*fill)
		delta := ts - f.time
		if delta >= sc.cutoff {
			m.Remove(e)
		} else {
			bucket := delta / 1000 // 1s buckets
			if !f.bucketed[bucket] {
				if f.buy {
					move := price/f.price - 1
					buyMoves, ok := sc.buyTakerMoves[securityID]
					if !ok {
						buyMoves = make(map[int64][2]float64)
						sc.buyTakerMoves[securityID] = buyMoves
					}
					buyMoves[bucket] = [2]float64{sc.alpha*buyMoves[bucket][0] + (1-sc.alpha)*move, float64(ts)}
				} else {
					move := f.price/price - 1
					sellMoves, ok := sc.sellTakerMoves[securityID]
					if !ok {
						sellMoves = make(map[int64][2]float64)
						sc.sellTakerMoves[securityID] = sellMoves
					}
					sellMoves[bucket] = [2]float64{sc.alpha*sellMoves[bucket][0] + (1-sc.alpha)*move, float64(ts)}
				}
				f.bucketed[bucket] = true
			}
		}
	}
}

func (sc *FillCollector) GetVolume(securityID uint64) float64 {
	sc.RLock()
	defer sc.RUnlock()
	ts := float64(time.Now().UnixMilli()) + 10
	v := sc.volume[securityID]
	delta := ts - v[1]
	w := math.Exp(-delta / float64(sc.volumeTau.Milliseconds()))
	return w * v[0]
}

func (sc *FillCollector) GetVolumes() map[uint64]float64 {
	sc.RLock()
	defer sc.RUnlock()
	vols := make(map[uint64]float64)
	ts := float64(time.Now().UnixMilli()) + 10
	for k, v := range sc.volume {
		delta := ts - v[1]
		w := math.Exp(-delta / float64(sc.volumeTau.Milliseconds()))
		vols[k] = v[0] * w
	}
	return vols
}

func (sc *FillCollector) GetSecurityMoveAfterFill(securityID uint64) ([]float64, []float64, []float64, []float64) {
	// TODO
	sc.RLock()
	defer sc.RUnlock()
	buyMakerMoves := make([]float64, 10)
	buyTakerMoves := make([]float64, 10)
	sellMakerMoves := make([]float64, 10)
	sellTakerMoves := make([]float64, 10)
	for i := 0; i < 10; i++ {
		v := sc.buyMakerMoves[securityID]
		if mv, ok := v[int64(i)]; ok {
			buyMakerMoves[i] = mv[0]
		}
	}
	for i := 0; i < 10; i++ {
		v := sc.buyTakerMoves[securityID]
		if mv, ok := v[int64(i)]; ok {
			buyTakerMoves[i] = mv[0]
		}
	}
	for i := 0; i < 10; i++ {
		v := sc.sellMakerMoves[securityID]
		if mv, ok := v[int64(i)]; ok {
			sellMakerMoves[i] += mv[0]
		}
	}
	for i := 0; i < 10; i++ {
		v := sc.sellTakerMoves[securityID]
		if mv, ok := v[int64(i)]; ok {
			sellTakerMoves[i] += mv[0]
		}
	}
	return buyMakerMoves, buyTakerMoves, sellMakerMoves, sellTakerMoves
}

func (sc *FillCollector) GetMoveAfterFill() ([]float64, []float64, []float64, []float64) {
	// TODO
	sc.RLock()
	defer sc.RUnlock()
	buyMakerMoves := make([]float64, 10)
	buyTakerMoves := make([]float64, 10)
	sellMakerMoves := make([]float64, 10)
	sellTakerMoves := make([]float64, 10)
	ts := float64(time.Now().UnixMilli()) + 10
	for i := 0; i < 10; i++ {
		sum := 0.
		for _, v := range sc.buyMakerMoves {
			if mv, ok := v[int64(i)]; ok {
				w := 1 / (ts - mv[1])
				buyMakerMoves[i] += mv[0] * w
				sum += w
			}
		}
		if sum != 0. {
			buyMakerMoves[i] /= sum
		}
	}
	for i := 0; i < 10; i++ {
		sum := 0.
		for _, v := range sc.buyTakerMoves {
			if mv, ok := v[int64(i)]; ok {
				w := 1 / (ts - mv[1])
				buyTakerMoves[i] += mv[0] * w
				sum += w
			}
		}
		if sum != 0. {
			buyTakerMoves[i] /= sum
		}
	}
	for i := 0; i < 10; i++ {
		sum := 0.
		for _, v := range sc.sellMakerMoves {
			if mv, ok := v[int64(i)]; ok {
				w := 1 / (ts - mv[1])
				sellMakerMoves[i] += mv[0] * w
				sum += w
			}
		}
		if sum != 0 {
			sellMakerMoves[i] /= sum
		}
	}
	for i := 0; i < 10; i++ {
		sum := 0.
		for _, v := range sc.sellTakerMoves {
			if mv, ok := v[int64(i)]; ok {
				w := 1 / (ts - mv[1])
				sellTakerMoves[i] += mv[0] * w
				sum += w
			}
		}
		if sum != 0 {
			sellTakerMoves[i] /= sum
		}
	}
	return buyMakerMoves, buyTakerMoves, sellMakerMoves, sellTakerMoves
}

// 100
// 200
// you will have weight of 0.01 and 0.005, so one will weight for 0.66 and the other for 0.33
