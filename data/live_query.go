package data

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/types"
	"github.com/melaurent/gotickfile/v2"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/tickobjects/market"
	"gitlab.com/alphaticks/tickstore-types/tickobjects"
	"gitlab.com/alphaticks/tickstore/parsing"
	"math"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

func ConstructFunctor(f parsing.Functor) (tickobjects.TickFunctor, reflect.Type, error) {
	if f.Measurement != nil {
		switch *f.Measurement {
		case "orderbook":
			return &market.RawOrderBook{}, reflect.TypeOf(market.RawOrderBookDelta{}), nil
		case "trade":
			return &market.RawTrade{}, reflect.TypeOf(market.RawTradeDelta{}), nil
		default:
			return nil, nil, fmt.Errorf("unknown measurement %s", *f.Measurement)
		}
	} else {
		objectType, deltaType, ok := tickobjects.GetTickObject(f.Application.Apply)
		if !ok {
			return nil, deltaType, fmt.Errorf("unknown functor %s", f.Application.Apply)
		}
		functor, ok := reflect.New(objectType).Interface().(tickobjects.TickFunctor)
		if !ok {
			return nil, deltaType, fmt.Errorf("%s is not a functor", f.Application.Apply)
		}
		childFunctor, _, err := ConstructFunctor(*f.Application.Value)
		if err != nil {
			return nil, nil, err
		} else {
			if err := functor.Initialize(childFunctor, f.Application.Args); err != nil {
				return nil, nil, fmt.Errorf("error initializing functor: %v", err)
			}
		}
		return functor, deltaType, nil
	}
}

type LiveQuery struct {
	sync.RWMutex
	pid           *actor.PID
	ch            chan interface{}
	subscriptions map[uint64]*Feed
	measurement   string
	objects       map[uint64]tickobjects.TickFunctor
	tick          uint64
	groupID       uint64
	functor       parsing.Functor
	deltas        *gotickfile.TickDeltas
	nextDeadline  *time.Time
	err           error
}

func NewLiveQuery(as *actor.ActorSystem, executor *actor.PID, sel parsing.Selector, feeds map[uint64]*Feed) (*LiveQuery, error) {
	// ConstructFunctor
	// Spawn listener
	// Let them push events in the chan
	lq := &LiveQuery{
		ch:            make(chan interface{}, 10000),
		subscriptions: make(map[uint64]*Feed),
		objects:       make(map[uint64]tickobjects.TickFunctor),
		tick:          0,
		functor:       sel.TickSelector.Functor,
		deltas:        nil,
		err:           nil,
	}

	tmpFunctor := sel.TickSelector.Functor
	for tmpFunctor.Measurement == nil {
		tmpFunctor = *tmpFunctor.Application.Value
	}
	lq.measurement = *tmpFunctor.Measurement

	for _, f := range feeds {
		aggregation := models.L2
		receiver := utils.NewMDReceiver(as, executor, &models.Instrument{
			SecurityID: &types.UInt64Value{Value: f.security.securityID},
		}, aggregation, f.requestID, lq.ch)
		f.receiver = receiver
		lq.subscriptions[f.requestID] = f
	}

	return lq, nil
}

func (lq *LiveQuery) SetNextDeadline(time time.Time) {
	lq.Lock()
	defer lq.Unlock()
	lq.nextDeadline = &time
}

func (lq *LiveQuery) Next() bool {
	lq.Lock()
	defer lq.Unlock()
	res := lq.next()
	lq.nextDeadline = nil
	return res
}

func (lq *LiveQuery) next() bool {
	var el interface{}
	if lq.nextDeadline != nil {
		select {
		case el = <-lq.ch:
			break
		case <-time.After(lq.nextDeadline.Sub(time.Now())):
			return false
		}
	} else {
		el = <-lq.ch
	}

	switch msg := el.(type) {
	case *messages.MarketDataResponse:
		feed := lq.subscriptions[msg.RequestID]

		if msg.SnapshotL2 == nil {
			lq.err = fmt.Errorf("response has no snapshotL2")
			return false
		}

		var tickPrecision uint64
		if msg.SnapshotL2.TickPrecision != nil {
			tickPrecision = msg.SnapshotL2.TickPrecision.Value
		} else if feed.security.minPriceIncrement != nil {
			tickPrecision = uint64(math.Ceil(1. / feed.security.minPriceIncrement.Value))
		} else {
			lq.err = fmt.Errorf("unable to get tick precision")
			return false
		}
		feed.security.tickPrecision = tickPrecision

		var lotPrecision uint64
		if msg.SnapshotL2.LotPrecision != nil {
			lotPrecision = msg.SnapshotL2.LotPrecision.Value
		} else if feed.security.roundLot != nil {
			lotPrecision = uint64(math.Ceil(1. / feed.security.roundLot.Value))
		} else {
			lq.err = fmt.Errorf("unable to get lot precision")
			return false
		}
		feed.security.lotPrecision = lotPrecision

		if lq.measurement == "orderbook" {

			ts := utils.TimestampToMilli(msg.SnapshotL2.Timestamp)

			ob := gorderbook.NewOrderBookL2(
				tickPrecision,
				lotPrecision,
				10000)
			ob.Sync(msg.SnapshotL2.Bids, msg.SnapshotL2.Asks)

			snapshot := market.NewRawOrderBook(ob).ToSnapshot()

			if functor, ok := lq.objects[feed.groupID]; ok {
				deltas, err := functor.ApplySnapshot(snapshot, feed.security.securityID, ts)
				if err != nil {
					lq.err = fmt.Errorf("error applying snapshot: %v", err)
					return false
				}
				lq.tick = ts
				lq.groupID = feed.groupID
				lq.deltas = deltas
				return true
			} else {
				functor, _, err := ConstructFunctor(lq.functor)
				if err != nil {
					lq.err = fmt.Errorf("error creating functor: %v", err)
					return false
				}
				lq.objects[feed.groupID] = functor

				deltas, err := functor.ApplySnapshot(snapshot, feed.security.securityID, ts)
				if err != nil {
					lq.err = fmt.Errorf("error applying snapshot: %v", err)
					return false
				}
				lq.tick = ts
				lq.groupID = feed.groupID
				lq.deltas = deltas

				return true
			}
		} else {
			return false
		}

	case *messages.MarketDataIncrementalRefresh:
		feed := lq.subscriptions[msg.RequestID]

		if lq.measurement == "orderbook" && msg.UpdateL2 != nil && len(msg.UpdateL2.Levels) > 0 {
			ts := utils.TimestampToMilli(msg.UpdateL2.Timestamp)
			var err error
			slice := make([]market.RawOrderBookDelta, len(msg.UpdateL2.Levels))
			// TODO check, not sure about that, do tests
			for i, l := range msg.UpdateL2.Levels {
				rawPrice := uint64(math.Round(l.Price * float64(feed.security.tickPrecision)))
				rawQty := uint64(math.Round(l.Quantity * float64(feed.security.lotPrecision)))
				slice[i], err = market.NewRawOrderBookDelta(rawPrice, rawQty, l.Bid, false)
				if err != nil {
					lq.err = fmt.Errorf("error building delta: %v", err)
					return false
				}
			}
			deltas := gotickfile.TickDeltas{
				Pointer: unsafe.Pointer(&slice[0]),
				Len:     len(slice),
			}

			if functor, ok := lq.objects[feed.groupID]; ok {
				resdeltas, err := functor.ApplyDeltas(deltas, feed.security.securityID, ts)
				if err != nil {
					lq.err = fmt.Errorf("error applying deltas: %v", err)
					return false
				}

				lq.tick = ts
				lq.groupID = feed.groupID
				lq.deltas = resdeltas

				return true
			} else {
				functor, _, err := ConstructFunctor(lq.functor)
				if err != nil {
					lq.err = fmt.Errorf("error creating functor: %v", err)
					return false
				}
				lq.objects[feed.groupID] = functor

				resdeltas, err := functor.ApplyDeltas(deltas, feed.security.securityID, ts)
				if err != nil {
					lq.err = fmt.Errorf("error applying snapshot: %v", err)
					return false
				}

				lq.tick = ts
				lq.groupID = feed.groupID
				lq.deltas = resdeltas

				return true
			}
		} else if lq.measurement == "trade" && len(msg.Trades) > 0 {
			var tradeDeltas []market.RawTradeDelta
			var ts uint64
			for _, aggTrade := range msg.Trades {
				ts = utils.TimestampToMilli(aggTrade.Timestamp)
				for _, trade := range aggTrade.Trades {
					rawPrice := uint64(math.Round(float64(feed.security.tickPrecision) * trade.Price))
					rawQuantity := uint64(math.Round(float64(feed.security.lotPrecision) * trade.Quantity))
					// Create delta
					dlt, err := market.NewRawTradeDelta(
						rawPrice,
						rawQuantity,
						trade.ID,
						aggTrade.AggregateID,
						aggTrade.Bid)
					if err != nil {
						lq.err = fmt.Errorf("error creating raw trade delta: %v", err)
					}
					tradeDeltas = append(tradeDeltas, dlt)
				}
			}

			var ptr unsafe.Pointer
			if len(tradeDeltas) > 0 {
				ptr = unsafe.Pointer(&tradeDeltas[0])
			}
			deltas := gotickfile.TickDeltas{
				Pointer: ptr,
				Len:     len(tradeDeltas),
			}

			if functor, ok := lq.objects[feed.groupID]; ok {
				resdeltas, err := functor.ApplyDeltas(deltas, feed.security.securityID, ts)
				if err != nil {
					lq.err = fmt.Errorf("error applying deltas: %v", err)
					return false
				}

				lq.tick = ts
				lq.groupID = feed.groupID
				lq.deltas = resdeltas

				return true
			} else {
				functor, _, _ := ConstructFunctor(lq.functor)
				aggTrade := market.NewRawTrade(feed.security.tickPrecision, feed.security.lotPrecision)
				_, err := functor.ApplySnapshot(aggTrade.ToSnapshot(), feed.security.securityID, ts)
				if err != nil {
					lq.err = fmt.Errorf("error applying snap: %v", err)
				}
				resdeltas, err := functor.ApplyDeltas(deltas, feed.security.securityID, ts)
				if err != nil {
					lq.err = fmt.Errorf("error applying deltas: %v", err)
				}
				lq.objects[feed.groupID] = functor

				lq.tick = ts
				lq.groupID = feed.groupID
				lq.deltas = resdeltas

				return true
			}
		} else {
			return false
		}
	default:
		return false
	}
}

func (lq *LiveQuery) Progress(end uint64) bool {
	lq.Lock()
	defer lq.Unlock()
	for end < lq.tick {
		lq.next()
	}
	return lq.tick >= end
}

func (lq *LiveQuery) Read() (uint64, tickobjects.TickObject, uint64) {
	lq.RLock()
	defer lq.RUnlock()
	return lq.tick, lq.objects[lq.groupID], lq.groupID
}

func (lq *LiveQuery) Tags() map[string]string {
	// TODO
	lq.RLock()
	defer lq.RUnlock()
	return make(map[string]string)
}

func (lq *LiveQuery) Close() error {
	lq.Lock()
	defer lq.Unlock()
	for _, f := range lq.subscriptions {
		if f.receiver != nil {
			f.receiver.Close()
			f.receiver = nil
		}
	}
	close(lq.ch)
	lq.err = fmt.Errorf("query closed")
	return nil
}

func (lq *LiveQuery) Err() error {
	lq.RLock()
	defer lq.RUnlock()
	return lq.err
}
