package utils

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/gorderbook"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math/rand"
	"sync"
	"time"
)

type MidPrice struct {
	Price float64
	Time  time.Time
}

type MarketDataContext struct {
	sync.RWMutex
	ctx           *actor.RootContext
	receiver      *actor.PID
	OBL2          *gorderbook.OrderBookL2
	LastMidPrices []MidPrice
}

func (s *MarketDataContext) Close() {
	if s.receiver == nil {
		return
	}
	s.ctx.Stop(s.receiver)
	s.receiver = nil
}

type MarketData struct {
	seqNum       uint64
	securityID   uint64
	executor     *actor.PID
	ctx          *MarketDataContext
	lastMidPrice float64
}

type checkTimeout struct{}

func NewMarketDataProducer(executor *actor.PID, securityID uint64, ctx *MarketDataContext) actor.Producer {
	return func() actor.Actor {
		return NewMarketData(executor, securityID, ctx)
	}
}

func NewMarketData(executor *actor.PID, securityID uint64, ctx *MarketDataContext) actor.Actor {
	return &MarketData{
		seqNum:     0,
		executor:   executor,
		securityID: securityID,
		ctx:        ctx,
	}
}

func (state *MarketData) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			panic(err)
		}

	case *actor.Stopping:

	case *actor.Stopped:

	case *actor.Restarting:

	case *messages.MarketDataIncrementalRefresh:
		if err := state.OnMarketDataIncrementalRefresh(context); err != nil {
			panic(err)
		}

	case *checkTimeout:
		if err := state.onCheckTimeout(context); err != nil {
			panic(err)
		}
		go func(pid *actor.PID) {
			time.Sleep(10 * time.Second)
			context.Send(pid, &checkTimeout{})
		}(context.Self())
	}
}

func (state *MarketData) Initialize(context actor.Context) error {
	state.ctx.receiver = context.Self()
	tmp, err := context.RequestFuture(state.executor, &messages.MarketDataRequest{
		RequestID:  rand.Uint64(),
		Subscribe:  true,
		Subscriber: context.Self(),
		Instrument: &models.Instrument{
			SecurityID: &wrapperspb.UInt64Value{Value: state.securityID},
		},
		Aggregation: models.OrderBookAggregation_L2,
	}, 10*time.Second).Result()
	if err != nil {
		panic(err)
	}
	mdres := tmp.(*messages.MarketDataResponse)

	state.ctx.Lock()
	defer state.ctx.Unlock()

	if mdres.SnapshotL2 != nil {
		var tickPrecision uint64
		if mdres.SnapshotL2.TickPrecision != nil {
			tickPrecision = mdres.SnapshotL2.TickPrecision.Value
		} else {
			panic("unable to get tick precision")
		}

		var lotPrecision uint64
		if mdres.SnapshotL2.LotPrecision != nil {
			lotPrecision = mdres.SnapshotL2.LotPrecision.Value
		} else {
			panic(fmt.Errorf("unable to get lot precision"))
		}

		ob := gorderbook.NewOrderBookL2(
			tickPrecision,
			lotPrecision,
			10000)
		ob.Sync(mdres.SnapshotL2.Bids, mdres.SnapshotL2.Asks)
		state.ctx.OBL2 = ob
		state.seqNum = mdres.SeqNum
	}
	return nil
}

func (state *MarketData) OnMarketDataIncrementalRefresh(context actor.Context) error {
	refresh := context.Message().(*messages.MarketDataIncrementalRefresh)
	if state.seqNum >= refresh.SeqNum {
		return nil
	}
	if state.seqNum+1 != refresh.SeqNum {
		return fmt.Errorf("out of order sequence")
	}
	state.ctx.Lock()
	if refresh.UpdateL2 != nil {
		for _, l := range refresh.UpdateL2.Levels {
			state.ctx.OBL2.UpdateOrderBookLevel(l)
		}
	}
	mp := (state.ctx.OBL2.BestAsk().Price + state.ctx.OBL2.BestBid().Price) / 2.
	if mp != state.lastMidPrice {
		state.ctx.LastMidPrices = append(state.ctx.LastMidPrices, MidPrice{
			Price: mp,
			Time:  time.Now(),
		})
		for len(state.ctx.LastMidPrices) > 0 && time.Now().Sub(state.ctx.LastMidPrices[0].Time) > 2*time.Second {
			state.ctx.LastMidPrices = state.ctx.LastMidPrices[1:]
		}
		state.lastMidPrice = mp
	}
	state.ctx.Unlock()
	state.seqNum = refresh.SeqNum
	return nil
}

func (state *MarketData) onCheckTimeout(context actor.Context) error {
	return nil
}
