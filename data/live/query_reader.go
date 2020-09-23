package live

import (
	"container/list"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"github.com/melaurent/gotickfile/v2"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/tachikoma.ai/tickobjects"
	"gitlab.com/tachikoma.ai/tickobjects/market"
	tickstore_grpc "gitlab.com/tachikoma.ai/tickstore-grpc"
	storeMessages "gitlab.com/tachikoma.ai/tickstore/actors/messages"
	"gitlab.com/tachikoma.ai/tickstore/db"
	"gitlab.com/tachikoma.ai/tickstore/parsing"
	"math"
	"reflect"
	"time"
	"unsafe"
)

type checkTimeout struct{}

type delayedReadQueryEventBatchRequest struct {
	message *storeMessages.ReadQueryEventBatchRequest
	sender  *actor.PID
}

type QueryReader struct {
	executor      *actor.PID
	index         *utils.TagIndex
	functor       parsing.Functor
	measurement   string
	deltaType     reflect.Type
	objects       map[uint64]tickobjects.TickFunctor
	logger        *log.Logger
	lastRead      time.Time
	feeds         map[uint64]*Feed
	subscriptions map[uint64]*Feed
	stashedMD     *list.List
}

func NewQueryReaderProducer(feeds map[uint64]*Feed, functor parsing.Functor) actor.Producer {
	return func() actor.Actor {
		return NewQueryReader(feeds, functor)
	}
}

func NewQueryReader(feeds map[uint64]*Feed, functor parsing.Functor) actor.Actor {
	return &QueryReader{
		feeds:       feeds,
		functor:     functor,
		measurement: "",
		logger:      nil,
		objects:     make(map[uint64]tickobjects.TickFunctor),
	}
}

// TODO timeout if no message, commit sudoku
func (state *QueryReader) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.logger.Error("error initializing actor", log.Error(err))
			context.Stop(context.Self())
		}
		state.logger.Info("actor started")

	case *actor.Stopping:
		state.logger.Info("actor stopping")
		state.Clean(context)

	case *actor.Stopped:
		state.logger.Info("actor stopped")

	case *actor.Restarting:
		state.logger.Info("actor restarting")
		state.Clean(context)

	case *storeMessages.ReadQueryEventBatchRequest:
		msg := context.Message().(*storeMessages.ReadQueryEventBatchRequest)
		dmsg := &delayedReadQueryEventBatchRequest{
			message: msg,
			sender:  context.Sender(),
		}
		go func(pid *actor.PID) {
			time.Sleep(time.Duration(msg.Timeout))
			context.Send(pid, dmsg)
		}(context.Self())

	case *delayedReadQueryEventBatchRequest:
		if err := state.ReadQueryEventBatchRequest(context); err != nil {
			state.logger.Error("error processing ReadTickEventBatchRequest", log.Error(err))
			context.Stop(context.Self())
		}

	case *storeMessages.CloseQueryReader:
		if err := state.CloseQueryReader(context); err != nil {
			state.logger.Error("error closing reader", log.Error(err))
			context.Stop(context.Self())
		}

	case *messages.MarketDataIncrementalRefresh:
		if err := state.OnMarketDataIncrementalRefresh(context); err != nil {
			state.logger.Error("error processing OnMarketDataIncrementalRefresh", log.Error(err))
			context.Stop(context.Self())
		}

	case *checkTimeout:
		if err := state.onCheckTimeout(context); err != nil {
			state.logger.Error("error checking timeout", log.Error(err))
			context.Stop(context.Self())
		}
		go func(pid *actor.PID) {
			time.Sleep(10 * time.Second)
			context.Send(pid, &checkTimeout{})
		}(context.Self())
	}
}

func (state *QueryReader) Clean(context actor.Context) {
	// TODO unsub
}

func (state *QueryReader) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	state.lastRead = time.Now()
	state.executor = actor.NewLocalPID("executor")
	state.subscriptions = make(map[uint64]*Feed)
	state.stashedMD = list.New()

	tmpFunctor := state.functor
	for tmpFunctor.Measurement == nil {
		tmpFunctor = *tmpFunctor.Application.Value
	}
	state.measurement = *tmpFunctor.Measurement

	if state.functor.Application != nil {
		_, deltaType, ok := tickobjects.GetTickObject(state.functor.Application.Apply)
		if !ok {
			return fmt.Errorf("unknwon functor %s", state.functor.Application.Apply)
		}
		state.deltaType = deltaType
	} else {
		return fmt.Errorf("live feed for raw measurement not supported")
	}

	// Sub to feeds
	for _, f := range state.feeds {
		state.subscriptions[f.requestID] = f
		if err := state.subscribeMarketData(context, f.security.securityID); err != nil {
			return fmt.Errorf("error subscribing to feed: %v", err)
		}
	}
	return nil
}

func (state *QueryReader) subscribeMarketData(context actor.Context, securityID uint64) error {
	feed := state.feeds[securityID]
	res, err := context.RequestFuture(state.executor, &messages.MarketDataRequest{
		RequestID:  feed.requestID,
		Subscribe:  true,
		Subscriber: context.Self(),
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: securityID},
		},
	}, 20*time.Second).Result()

	if err != nil {
		return fmt.Errorf("error getting snapshot: %v", err)
	}
	response, ok := res.(*messages.MarketDataResponse)
	if !ok {
		return fmt.Errorf("was expecting *MarketDataSnapshot, got %s", reflect.TypeOf(res).String())
	}
	state.stashedMD.PushBack(response)

	feed.seqNum = response.SeqNum
	feed.lastEventTime = time.Now()

	return nil
}

func (state *QueryReader) OnMarketDataIncrementalRefresh(context actor.Context) error {
	refresh := context.Message().(*messages.MarketDataIncrementalRefresh)
	feed := state.subscriptions[refresh.RequestID]
	if refresh.SeqNum <= feed.seqNum {
		//fmt.Println("SKIPPING", refresh.SeqNum, state.securityInfo.seqNum)
		return nil
	}
	if feed.seqNum+1 != refresh.SeqNum {
		//fmt.Println("OUT OF SYNC", secInfo.seqNum, refresh.SeqNum)
		return state.subscribeMarketData(context, feed.security.securityID)
	}

	state.stashedMD.PushBack(refresh)

	feed.seqNum = refresh.SeqNum
	feed.lastEventTime = time.Now()

	// Now use the sampler to update models

	return nil
}

func (state *QueryReader) ReadQueryEventBatchRequest(context actor.Context) error {
	msg := context.Message().(*delayedReadQueryEventBatchRequest)

	batch := &tickstore_grpc.QueryEventBatch{
		Events: nil,
	}

	request := msg.message
	events := make([]*tickstore_grpc.QueryEvent, 0, request.BatchSize)
	// Process stashed events
	for uint64(len(events)) < request.BatchSize && state.stashedMD.Front() != nil {
		el := state.stashedMD.Front()
		switch el.Value.(type) {
		case *messages.MarketDataResponse:
			msg := el.Value.(*messages.MarketDataResponse)
			if msg.SnapshotL2 == nil {
				return fmt.Errorf("response has no snapshotL2")
			}
			ts := utils.TimestampToMilli(msg.SnapshotL2.Timestamp)
			feed := state.subscriptions[msg.RequestID]
			if state.measurement == "orderbook" {

				var tickPrecision uint64
				if msg.SnapshotL2.TickPrecision != nil {
					tickPrecision = msg.SnapshotL2.TickPrecision.Value
				} else if feed.security.minPriceIncrement != nil {
					tickPrecision = uint64(math.Ceil(1. / feed.security.minPriceIncrement.Value))
				} else {
					return fmt.Errorf("unable to get tick precision")
				}
				feed.security.tickPrecision = tickPrecision

				var lotPrecision uint64
				if msg.SnapshotL2.LotPrecision != nil {
					lotPrecision = msg.SnapshotL2.LotPrecision.Value
				} else if feed.security.roundLot != nil {
					lotPrecision = uint64(math.Ceil(1. / feed.security.roundLot.Value))
				} else {
					return fmt.Errorf("unable to get lo precision")
				}
				feed.security.lotPrecision = lotPrecision

				var event *tickstore_grpc.QueryEvent

				ob := gorderbook.NewOrderBookL2(
					tickPrecision,
					lotPrecision,
					10000)
				ob.Sync(msg.SnapshotL2.Bids, msg.SnapshotL2.Asks)

				snapshot := market.NewRawOrderBook(ob).ToSnapshot()

				if functor, ok := state.objects[feed.groupID]; ok {
					deltas, err := functor.ApplySnapshot(snapshot, feed.security.securityID, ts)
					if err != nil {
						return fmt.Errorf("error applying snapshot: %v", err)
					}

					buff := make([]byte, int(state.deltaType.Size())*deltas.Len)
					ptr := uintptr(deltas.Pointer)
					for i := 0; i < len(buff); i++ {
						buff[i] = *(*byte)(unsafe.Pointer(ptr + uintptr(i)))
					}
					event = &tickstore_grpc.QueryEvent{
						Tick:    ts,
						GroupId: feed.groupID,
						Event: &tickstore_grpc.QueryEvent_Deltas{
							Deltas: buff,
						},
					}

				} else {
					functor, _, err := db.ConstructFunctor(state.functor)
					if err != nil {
						return fmt.Errorf("error creating functor: %v", err)
					}
					state.objects[feed.groupID] = functor

					_, err = functor.ApplySnapshot(snapshot, feed.security.securityID, ts)
					if err != nil {
						return fmt.Errorf("error applying snapshot: %v", err)
					}

					event = &tickstore_grpc.QueryEvent{
						Tick:    ts,
						GroupId: feed.groupID,
						Event: &tickstore_grpc.QueryEvent_Snapshot{
							Snapshot: &tickstore_grpc.QuerySnapshot{
								Tags:     feed.tags,
								Snapshot: functor.ToSnapshot(),
							},
						},
					}
				}

				events = append(events, event)
				feed.lastDeltaTime = ts
			}

		case *messages.MarketDataIncrementalRefresh:
			msg := el.Value.(*messages.MarketDataIncrementalRefresh)
			feed := state.subscriptions[msg.RequestID]
			if msg.UpdateL2 != nil && len(msg.UpdateL2.Levels) > 0 && state.measurement == "orderbook" {
				var event *tickstore_grpc.QueryEvent
				ts := utils.TimestampToMilli(msg.UpdateL2.Timestamp)

				var err error
				slice := make([]market.RawOrderBookDelta, len(msg.UpdateL2.Levels))
				// TODO check, not sure about that, do tests
				for i, l := range msg.UpdateL2.Levels {
					rawPrice := uint64(math.Round(l.Price * float64(feed.security.tickPrecision)))
					rawQty := uint64(math.Round(l.Quantity * float64(feed.security.lotPrecision)))
					slice[i], err = market.NewRawOrderBookDelta(rawPrice, rawQty, l.Bid, false)
					if err != nil {
						return fmt.Errorf("error building delta: %v", err)
					}
				}
				feed.lastDeltaTime = ts
				deltas := gotickfile.TickDeltas{
					Pointer: unsafe.Pointer(&slice[0]),
					Len:     len(slice),
				}

				if functor, ok := state.objects[feed.groupID]; ok {
					deltas, err := functor.ApplyDeltas(deltas, feed.security.securityID, ts)
					if err != nil {
						return fmt.Errorf("error applying deltas: %v", err)
					}

					buff := make([]byte, int(state.deltaType.Size())*deltas.Len)
					ptr := uintptr(deltas.Pointer)
					for i := 0; i < len(buff); i++ {
						buff[i] = *(*byte)(unsafe.Pointer(ptr + uintptr(i)))
					}
					event = &tickstore_grpc.QueryEvent{
						Tick:    ts,
						GroupId: feed.groupID,
						Event: &tickstore_grpc.QueryEvent_Deltas{
							Deltas: buff,
						},
					}

				} else {
					functor, _, err := db.ConstructFunctor(state.functor)
					if err != nil {
						return fmt.Errorf("error creating functor: %v", err)
					}
					state.objects[feed.groupID] = functor

					_, err = functor.ApplyDeltas(deltas, feed.security.securityID, ts)
					if err != nil {
						return fmt.Errorf("error applying snapshot: %v", err)
					}

					event = &tickstore_grpc.QueryEvent{
						Tick:    ts,
						GroupId: feed.groupID,
						Event: &tickstore_grpc.QueryEvent_Snapshot{
							Snapshot: &tickstore_grpc.QuerySnapshot{
								Tags:     feed.tags,
								Snapshot: functor.ToSnapshot(),
							},
						},
					}
				}

				events = append(events, event)
			}

			if state.measurement == "trade" {
				for _, _ = range msg.Trades {

				}
			}
		}
		state.stashedMD.Remove(el)
	}

	batch.Events = events
	errMsg := ""

	context.Send(msg.sender, &storeMessages.ReadQueryEventBatchResponse{
		RequestID: request.RequestID,
		Error:     errMsg,
		Batch:     batch,
	})

	// If we have an error, we stop
	if errMsg != "" {
		context.Stop(context.Self())
	}

	state.lastRead = time.Now()

	return nil
}

func (state *QueryReader) onCheckTimeout(context actor.Context) error {
	if time.Now().Sub(state.lastRead) > 180*time.Second {
		return state.CloseQueryReader(context)
	}
	return nil
}

func (state *QueryReader) CloseQueryReader(context actor.Context) error {
	// commit sudoku
	context.Stop(context.Self())

	return nil
}
