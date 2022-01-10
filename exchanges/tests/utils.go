package tests

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/account"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/gorderbook"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
)

func StartExecutor(t *testing.T, exchange *xchangerModels.Exchange, acc *models.Account) (*actor.ActorSystem, *actor.PID, func()) {
	exch := []*xchangerModels.Exchange{
		exchange,
	}
	as := actor.NewActorSystem()
	var accnts []*account.Account
	if acc != nil {
		accnt, err := exchanges.NewAccount(acc)
		if err != nil {
			t.Fatal(err)
		}
		accnts = append(accnts, accnt)
	}

	cfg := &exchanges.ExecutorConfig{
		Exchanges: exch,
		Strict:    true,
		Accounts:  accnts,
	}
	executor, _ := as.Root.SpawnNamed(actor.PropsFromProducer(exchanges.NewExecutorProducer(cfg)), "executor")
	return as, executor, func() { _ = as.Root.PoisonFuture(executor).Wait() }
}

type GetStat struct {
	Error     error
	Trades    int
	AggTrades int
	OBUpdates int
}

type MDChecker struct {
	test          MDTest
	security      *models.Security
	orderbook     *gorderbook.OrderBookL2
	tickPrecision uint64
	lotPrecision  uint64
	seqNum        uint64
	synced        bool
	trades        int
	aggTrades     int
	aggTradeIDs   map[uint64]bool
	OBUpdates     int
	err           error
}

func NewMDCheckerProducer(security *models.Security, test MDTest) actor.Producer {
	return func() actor.Actor {
		return NewMDChecker(security, test)
	}
}

func NewMDChecker(security *models.Security, test MDTest) actor.Actor {
	return &MDChecker{
		test:        test,
		security:    security,
		orderbook:   nil,
		seqNum:      0,
		synced:      false,
		trades:      0,
		aggTrades:   0,
		aggTradeIDs: make(map[uint64]bool),
		OBUpdates:   0,
		err:         nil,
	}
}

func (state *MDChecker) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		fmt.Println("INITIALIZING")
		if err := state.Initialize(context); err != nil {
			state.err = err
		}
		fmt.Println("INITIALIZED")

	case *messages.MarketDataIncrementalRefresh:
		if state.err == nil {
			if err := state.OnMarketDataIncrementalRefresh(context); err != nil {
				state.err = err
			}
		}

	case *GetStat:
		context.Respond(&GetStat{
			Error:     state.err,
			Trades:    state.trades,
			AggTrades: state.aggTrades,
			OBUpdates: state.OBUpdates,
		})
	}
}

func (state *MDChecker) Initialize(context actor.Context) error {
	executor := context.ActorSystem().NewLocalPID("executor")
	res, err := context.RequestFuture(executor, &messages.MarketDataRequest{
		RequestID:  0,
		Subscribe:  true,
		Subscriber: context.Self(),
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: state.security.SecurityID},
			Exchange:   state.security.Exchange,
			Symbol:     &types.StringValue{Value: state.security.Symbol},
		},
		Aggregation: models.L2,
	}, 80*time.Second).Result()
	if err != nil {
		return err
	}
	response, ok := res.(*messages.MarketDataResponse)
	if !ok {
		return fmt.Errorf("was expecting market data snapshot, got %s", reflect.TypeOf(res).String())
	}

	var tickPrecision uint64
	if response.SnapshotL2.TickPrecision != nil {
		tickPrecision = response.SnapshotL2.TickPrecision.Value
	} else if state.security.MinPriceIncrement != nil {
		tickPrecision = uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))
	} else {
		return fmt.Errorf("unable to get tick precision")
	}

	var lotPrecision uint64
	if response.SnapshotL2.LotPrecision != nil {
		lotPrecision = response.SnapshotL2.LotPrecision.Value
	} else if state.security.RoundLot != nil {
		lotPrecision = uint64(math.Ceil(1. / state.security.RoundLot.Value))
	} else {
		return fmt.Errorf("unable to get lo precision")
	}

	state.tickPrecision = tickPrecision
	state.lotPrecision = lotPrecision

	for _, b := range response.SnapshotL2.Bids {
		if !state.test.IgnoreSizeResidue {
			rawQty := b.Quantity * float64(lotPrecision)
			if (math.Round(rawQty) - rawQty) > 0.01 {
				return fmt.Errorf("residue in qty: %f %f", rawQty, math.Round(rawQty))
			}
		}
		if !state.test.IgnorePriceResidue {
			rawPrice := b.Price * float64(tickPrecision)
			if (math.Round(rawPrice) - rawPrice) > 0.00001 {
				return fmt.Errorf("residue in price: %f %f", rawPrice, math.Round(rawPrice))
			}
		}
	}
	for _, a := range response.SnapshotL2.Asks {
		if !state.test.IgnoreSizeResidue {
			rawQty := a.Quantity * float64(lotPrecision)
			if (math.Round(rawQty) - rawQty) > 0.01 {
				return fmt.Errorf("residue in qty: %f %f", rawQty, math.Round(rawQty))
			}
		}
		if !state.test.IgnorePriceResidue {
			rawPrice := a.Price * float64(tickPrecision)
			if (math.Round(rawPrice) - rawPrice) > 0.00001 {
				return fmt.Errorf("residue in price: %f %f", rawPrice, math.Round(rawPrice))
			}
		}
	}
	state.OBUpdates += 1
	state.orderbook = gorderbook.NewOrderBookL2(uint64(tickPrecision), uint64(lotPrecision), 10000)
	state.orderbook.Sync(response.SnapshotL2.Bids, response.SnapshotL2.Asks)
	state.seqNum = response.SeqNum
	if state.orderbook.Crossed() {
		return fmt.Errorf("crossed OB on snapshot \n" + state.orderbook.String())
	}
	return nil
}

func (state *MDChecker) OnMarketDataIncrementalRefresh(context actor.Context) error {
	refresh := context.Message().(*messages.MarketDataIncrementalRefresh)

	if !state.synced && refresh.SeqNum <= state.seqNum {
		//fmt.Println("SKIPPING", refresh.SeqNum, state.securityInfo.seqNum)
		return nil
	}
	state.synced = true
	if state.seqNum+1 != refresh.SeqNum {
		//fmt.Println("OUT OF SYNC", state.securityInfo.seqNum, refresh.SeqNum)
		return fmt.Errorf("out of order sequence %d %d", state.seqNum, refresh.SeqNum)
	}

	if refresh.UpdateL2 != nil {
		for _, l := range refresh.UpdateL2.Levels {
			if !state.test.IgnoreSizeResidue {
				rawQty := l.Quantity * float64(state.lotPrecision)
				if (math.Round(rawQty) - rawQty) > 0.01 {
					return fmt.Errorf("residue in qty: %f %f", rawQty, math.Round(rawQty))
				}
			}
			if !state.test.IgnorePriceResidue {
				rawPrice := l.Price * float64(state.tickPrecision)
				if (math.Round(rawPrice) - rawPrice) > 0.00001 {
					return fmt.Errorf("residue in price: %f %f", rawPrice, math.Round(rawPrice))
				}
			}
			state.orderbook.UpdateOrderBookLevel(l)
		}

		if state.orderbook.Crossed() {
			fmt.Println("CROSSED")
			for _, l := range refresh.UpdateL2.Levels {
				fmt.Println(l)
			}
			return fmt.Errorf("crossed OB \n" + state.orderbook.String())
		}
		state.OBUpdates += 1
	}

	if refresh.Funding != nil {
		fmt.Println(refresh.Funding.Rate, refresh.Funding.Timestamp)
	}

	for _, aggT := range refresh.Trades {
		if _, ok := state.aggTradeIDs[aggT.AggregateID]; ok {
			return fmt.Errorf("duplicate aggregate ID")
		}
		state.aggTradeIDs[aggT.AggregateID] = true
		state.aggTrades += 1
		for _, t := range aggT.Trades {
			rawPrice := t.Price * float64(state.tickPrecision)
			rawQty := t.Quantity * float64(state.lotPrecision)
			if (math.Round(rawPrice) - rawPrice) > 0.00001 {
				return fmt.Errorf("residue in trade price: %f %f", rawPrice, math.Round(rawPrice))
			}
			if (math.Round(rawQty) - rawQty) > 0.0001 {
				return fmt.Errorf("residue in trade qty: %f %f", rawQty, math.Round(rawQty))
			}
			state.trades += 1
		}
	}
	state.seqNum = refresh.SeqNum
	return nil
}

func (state *MDChecker) OnMarketDataSnapshot(context actor.Context) error {
	return nil
}

type PoolV3Checker struct {
	test          MDTest
	security      *models.Security
	pool          *gorderbook.UnipoolV3
	tickPrecision uint64
	lotPrecision  uint64
	seqNum        uint64
	synced        bool
	trades        int
	aggTrades     int
	aggTradeIDs   map[uint64]bool
	OBUpdates     int
	err           error
}

func NewPoolV3CheckerProducer(security *models.Security, test MDTest) actor.Producer {
	return func() actor.Actor {
		return NewPoolV3Checker(security, test)
	}
}

func NewPoolV3Checker(security *models.Security, test MDTest) actor.Actor {
	return &PoolV3Checker{
		test:        test,
		security:    security,
		pool:        nil,
		seqNum:      0,
		synced:      false,
		trades:      0,
		aggTrades:   0,
		aggTradeIDs: make(map[uint64]bool),
		OBUpdates:   0,
		err:         nil,
	}
}

func (state *PoolV3Checker) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		fmt.Println("INITIALIZING")
		if err := state.Initialize(context); err != nil {
			state.err = err
		}
		fmt.Println("INITIALIZED")

	case *messages.UnipoolV3DataIncrementalRefresh:
		if state.err == nil {
			if err := state.OnUnipoolV3DataIncrementalRefresh(context); err != nil {
				state.err = err
			}
		}

	case *GetStat:
		context.Respond(&GetStat{
			Error:     state.err,
			Trades:    state.trades,
			AggTrades: state.aggTrades,
			OBUpdates: state.OBUpdates,
		})
	}
}

func (state *PoolV3Checker) Initialize(context actor.Context) error {
	executor := context.ActorSystem().NewLocalPID("executor")
	res, err := context.RequestFuture(executor, &messages.UnipoolV3DataRequest{
		RequestID:  0,
		Subscribe:  true,
		Subscriber: context.Self(),
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: state.security.SecurityID},
			Exchange:   state.security.Exchange,
			Symbol:     &types.StringValue{Value: state.security.Symbol},
		},
	}, 80*time.Second).Result()
	if err != nil {
		return err
	}
	response, ok := res.(*messages.UnipoolV3DataResponse)
	if !ok {
		return fmt.Errorf("was expecting market data snapshot, got %s", reflect.TypeOf(res).String())
	}

	state.OBUpdates += 1
	feeTier := int32(state.security.TakerFee.Value * 1e6)
	state.pool = gorderbook.NewUnipoolV3(feeTier)
	state.seqNum = response.SeqNum
	return nil
}

func (state *PoolV3Checker) OnUnipoolV3DataIncrementalRefresh(context actor.Context) error {
	refresh := context.Message().(*messages.UnipoolV3DataIncrementalRefresh)

	if !state.synced && refresh.SeqNum <= state.seqNum {
		//fmt.Println("SKIPPING", refresh.SeqNum, state.securityInfo.seqNum)
		return nil
	}
	state.synced = true
	if state.seqNum+1 != refresh.SeqNum {
		//fmt.Println("OUT OF SYNC", state.securityInfo.seqNum, refresh.SeqNum)
		return fmt.Errorf("out of order sequence %d %d", state.seqNum, refresh.SeqNum)
	}

	state.seqNum = refresh.SeqNum
	return nil
}

func (state *PoolV3Checker) OnMarketDataSnapshot(context actor.Context) error {
	return nil
}
