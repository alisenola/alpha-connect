package utils

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/gorderbook"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
	"math/rand"
	"sync"
	"time"
)

type SpawnStopContext interface {
	actor.SpawnerContext
	Stop(pid *actor.PID)
}

type MarketDataContext struct {
	sync.RWMutex
	MarketDataContextSettings
	ctx            SpawnStopContext
	receiver       *actor.PID
	OBL2           *gorderbook.OrderBookL2
	Stats          map[models.StatType]float64
	VWAP           float64
	lastVolumeTime time.Time
	Volume         float64
	vwapVolume     float64
	vwapPrice      float64
	TradeCond      *sync.Cond
	BookCond       *sync.Cond
}

func (s *MarketDataContext) Close() {
	if s.receiver == nil {
		return
	}
	s.ctx.Stop(s.receiver)
	s.receiver = nil
}

type MarketDataReceiver struct {
	seqNum       uint64
	securityID   uint64
	executor     *actor.PID
	ctx          *MarketDataContext
	lastMidPrice float64
	ticker       *time.Ticker
}

type checkTimeout struct{}

type MarketDataContextSettings struct {
	VWAPCoefficient float64
	VolumeTau       time.Duration
}

func NewMarketDataContext(parent SpawnStopContext, executor *actor.PID, securityID uint64, settings MarketDataContextSettings) *MarketDataContext {
	ctx := &MarketDataContext{
		MarketDataContextSettings: settings,
		ctx:                       parent,
		Stats:                     make(map[models.StatType]float64),
		TradeCond:                 sync.NewCond(&sync.Mutex{}),
		BookCond:                  sync.NewCond(&sync.Mutex{}),
	}
	ctx.receiver = parent.Spawn(actor.PropsFromProducer(NewMarketDataReceiverProducer(executor, securityID, ctx)))
	return ctx
}

func NewMarketDataReceiverProducer(executor *actor.PID, securityID uint64, ctx *MarketDataContext) actor.Producer {
	return func() actor.Actor {
		return NewMarketDataReceiver(executor, securityID, ctx)
	}
}

func NewMarketDataReceiver(executor *actor.PID, securityID uint64, ctx *MarketDataContext) actor.Actor {
	return &MarketDataReceiver{
		seqNum:     0,
		executor:   executor,
		securityID: securityID,
		ctx:        ctx,
	}
}

func (state *MarketDataReceiver) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			panic(err)
		}

	case *actor.Stopping:
		state.ticker = nil

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
	}
}

func (state *MarketDataReceiver) Initialize(context actor.Context) error {
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
	if !mdres.Success {
		return fmt.Errorf("failed to subscribe to market data: %s", mdres.RejectionReason.String())
	}
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

	ticker := time.NewTicker(10 * time.Second)
	state.ticker = ticker
	go func(pid *actor.PID) {
		for {
			select {
			case <-ticker.C:
				context.Send(pid, &checkTimeout{})
			case <-time.After(20 * time.Second):
				if state.ticker != ticker {
					// Only stop if socket ticker has changed
					return
				}
			}
		}
	}(context.Self())
	return nil
}

func (state *MarketDataReceiver) OnMarketDataIncrementalRefresh(context actor.Context) error {
	refresh := context.Message().(*messages.MarketDataIncrementalRefresh)
	if state.seqNum >= refresh.SeqNum {
		return nil
	}
	if state.seqNum+1 != refresh.SeqNum {
		return fmt.Errorf("out of order sequence")
	}
	var crossed bool
	state.ctx.Lock()
	if refresh.UpdateL2 != nil {
		for _, l := range refresh.UpdateL2.Levels {
			state.ctx.OBL2.UpdateOrderBookLevel(l)
		}
		crossed = state.ctx.OBL2.Crossed()
	}
	for _, s := range refresh.Stats {
		state.ctx.Stats[s.StatType] = s.Value
	}
	// if trade event, put trade flag to 1
	volume := 0.
	if len(refresh.Trades) > 0 {
		alpha := state.ctx.VWAPCoefficient
		for _, trds := range refresh.Trades {
			for _, trd := range trds.Trades {
				state.ctx.vwapVolume = alpha*state.ctx.vwapVolume + (1-alpha)*math.Abs(trd.Quantity)
				state.ctx.vwapPrice = alpha*state.ctx.vwapPrice + (1-alpha)*math.Abs(trd.Price*trd.Quantity)
				state.ctx.VWAP = state.ctx.vwapPrice / state.ctx.vwapVolume
				// compute alpha from half life
				/*
						def halfLifeToLambda(halfLife, sampFreq):
						    return 0.5 ** (1 / (halfLife / sampFreq))
					w = np.exp(-(delta) / tau)
				*/
				volume += math.Abs(trd.Quantity * trd.Price)
			}
		}
		state.ctx.TradeCond.Broadcast()
	}
	if refresh.UpdateL2 != nil {
		state.ctx.BookCond.Broadcast()
	}

	delta := float64(time.Since(state.ctx.lastVolumeTime))
	w := math.Exp(-delta / float64(state.ctx.VolumeTau))
	// Decay volume + add new quantity
	state.ctx.Volume = w * state.ctx.Volume
	state.ctx.Volume += volume
	state.ctx.lastVolumeTime = time.Now()
	state.ctx.Unlock()

	if crossed {
		return fmt.Errorf("crossed ob")
	}
	state.seqNum = refresh.SeqNum
	return nil
}

func (state *MarketDataReceiver) onCheckTimeout(context actor.Context) error {
	return nil
}
