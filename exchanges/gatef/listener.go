package gatef

import (
	"container/list"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
	"reflect"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/gorderbook"
	gmodels "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	"gitlab.com/alphaticks/xchanger"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/gatef"
	xutils "gitlab.com/alphaticks/xchanger/utils"
)

type checkSockets struct{}

type postAggTrade struct{}

type InstrumentData struct {
	orderBook      *gorderbook.OrderBookL2
	seqNum         uint64
	lastUpdateTime uint64
	lastUpdateID   uint64
	lastHBTime     time.Time
	aggTrade       *models.AggregatedTrade
	lastAggTradeTs uint64
}

type Listener struct {
	sink           *xutils.WebsocketSink
	wsPool         *xutils.WebsocketPool
	security       *models.Security
	securityID     uint64
	instrumentData *InstrumentData
	logger         *log.Logger
	lastPingTime   time.Time
	socketTicker   *time.Ticker
	executor       *actor.PID
	stashedTrades  *list.List
}

func NewListenerProducer(securityID uint64, wsPool *xutils.WebsocketPool) actor.Producer {
	return func() actor.Actor {
		return NewListener(securityID, wsPool)
	}
}

func NewListener(securityID uint64, wsPool *xutils.WebsocketPool) actor.Actor {
	return &Listener{
		securityID: securityID,
		wsPool:     wsPool,
	}
}

func (state *Listener) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.logger.Error("error initializing", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor started")

	case *actor.Stopping:
		if err := state.Clean(context); err != nil {
			state.logger.Error("error stopping", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor stopping")

	case *actor.Stopped:
		state.logger.Info("actor stopped")

	case *actor.Restarting:
		if err := state.Clean(context); err != nil {
			state.logger.Error("error restarting", log.Error(err))
			// Attention, no panic in restarting or infinite loop
		}
		state.logger.Info("actor restarting")

	case *messages.MarketDataRequest:
		if err := state.OnMarketDataRequest(context); err != nil {
			state.logger.Error("error processing OnMarketDataRequest", log.Error(err))
			panic(err)
		}

	case *xchanger.WebsocketMessage:
		if err := state.onWebsocketMessage(context); err != nil {
			state.logger.Error("error processing websocket message", log.Error(err))
			panic(err)
		}

	case *checkSockets:
		if err := state.checkSockets(context); err != nil {
			state.logger.Error("error checking socket", log.Error(err))
			panic(err)
		}
	case *postAggTrade:
		state.postAggTrade(context)
	}
}

func (state *Listener) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()),
		log.String("security-id", fmt.Sprintf("%d", state.securityID)))
	state.executor = actor.NewPID(context.ActorSystem().Address(), "executor/exchanges/"+constants.GATEF.Name+"_executor")

	res, err := context.RequestFuture(state.executor, &messages.SecurityDefinitionRequest{
		RequestID:  0,
		Instrument: &models.Instrument{SecurityID: wrapperspb.UInt64(state.securityID)},
	}, 5*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error fetching security definition: %v", err)
	}
	def := res.(*messages.SecurityDefinitionResponse)
	if !def.Success {
		return fmt.Errorf("error fetching security definition: %s", def.RejectionReason.String())
	}
	state.security = def.Security
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()),
		log.String("security-id", fmt.Sprintf("%d", state.securityID)),
		log.String("exchange", state.security.Exchange.Name),
		log.String("symbol", state.security.Symbol))
	if state.security.MinPriceIncrement == nil || state.security.RoundLot == nil {
		return fmt.Errorf("security is missing MinPriceIncrement or RoundLot")
	}
	state.lastPingTime = time.Now()
	state.stashedTrades = list.New()

	state.instrumentData = &InstrumentData{
		orderBook:      nil,
		seqNum:         uint64(time.Now().UnixNano()),
		lastUpdateTime: 0,
		lastHBTime:     time.Now(),
		lastAggTradeTs: 0,
	}

	if err := state.subscribeInstrument(context); err != nil {
		return fmt.Errorf("error subscribing to order book: %v", err)
	}

	socketTicker := time.NewTicker(5 * time.Second)
	state.socketTicker = socketTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-socketTicker.C:
				context.Send(pid, &checkSockets{})
			case <-time.After(10 * time.Second):
				if state.socketTicker != socketTicker {
					// Only stop if socket ticker has changed
					return
				}
			}
		}
	}(context.Self())

	return nil
}

func (state *Listener) Clean(context actor.Context) error {
	if state.sink != nil {
		state.wsPool.Unsubscribe(state.security.Symbol, state.sink)
		state.sink = nil
	}
	if state.socketTicker != nil {
		state.socketTicker.Stop()
		state.socketTicker = nil
	}

	return nil
}

func (state *Listener) subscribeInstrument(context actor.Context) error {
	if state.sink != nil {
		state.wsPool.Unsubscribe(state.security.Symbol, state.sink)
		state.sink = nil
	}

	sink, err := state.wsPool.Subscribe(state.security.Symbol)
	if err != nil {
		return fmt.Errorf("error subscribing to websocket: %v", err)
	}
	state.sink = sink

	sync := func() (*gorderbook.OrderBookL2, error) {
		fut := context.RequestFuture(
			state.executor,
			&messages.MarketDataRequest{
				RequestID: uint64(time.Now().UnixNano()),
				Subscribe: false,
				Instrument: &models.Instrument{
					SecurityID: &wrapperspb.UInt64Value{Value: state.security.SecurityID},
					Symbol:     &wrapperspb.StringValue{Value: state.security.Symbol},
					Exchange:   state.security.Exchange,
				},
				Aggregation: models.OrderBookAggregation_L2,
			},
			20*time.Second)

		res, err := fut.Result()
		if err != nil {
			return nil, fmt.Errorf("error getting market data: %v", err)
		}

		mdres, ok := res.(*messages.MarketDataResponse)
		if !ok {
			return nil, fmt.Errorf("was expecting MarketDataResponse, got %s", reflect.TypeOf(res).String())
		}
		if !mdres.Success {
			if mdres.RejectionReason == messages.RejectionReason_IPRateLimitExceeded {
				return nil, nil
			} else {
				return nil, fmt.Errorf("error fetching the snapshot %s", mdres.RejectionReason.String())
			}
		}
		if mdres.SnapshotL2 == nil {
			return nil, fmt.Errorf("market data snapshot has no OBL2")
		}
		tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))
		lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot.Value))
		depth := 10000
		if len(mdres.SnapshotL2.Asks) > 0 {
			bestAsk := mdres.SnapshotL2.Asks[0].Price
			depth = int(((bestAsk * 1.1) - bestAsk) * float64(tickPrecision))
			if depth > 10000 {
				depth = 10000
			}
		}

		ob := gorderbook.NewOrderBookL2(
			tickPrecision,
			lotPrecision,
			depth,
		)
		ob.Sync(mdres.SnapshotL2.Bids, mdres.SnapshotL2.Asks)
		if ob.Crossed() {
			return nil, fmt.Errorf("crossed orderbook")
		}
		state.instrumentData.lastUpdateID = mdres.SeqNum
		state.instrumentData.lastUpdateTime = utils.TimestampToMilli(mdres.SnapshotL2.Timestamp)
		return ob, nil
	}

	ob, err := sync()
	if err != nil {
		return fmt.Errorf("error syncing order book: %v", err)
	}
	//fmt.Println("SEQ", state.instrumentData.lastUpdateID)
	//fmt.Println(ob)
	synced := false
	for !synced {
		msg, err := sink.ReadMessage()
		if err != nil {
			return fmt.Errorf("error reading message: %v", err)
		}
		obUpdate, ok := msg.Message.(gatef.WSFuturesOrderBookUpdate)
		if !ok {
			context.Send(context.Self(), msg)
		}
		for state.instrumentData.lastUpdateID+1 < obUpdate.FirstUpdateId {
			// We missed some updates, snapshot is too old
			//fmt.Println("syncing again.....", state.instrumentData.lastUpdateID+1, obUpdate.FirstUpdateId)
			time.Sleep(5 * time.Second)
			ob, err = sync()
			if err != nil {
				return fmt.Errorf("error syncing order book: %v", err)
			}
			//fmt.Println("SEQ", state.instrumentData.lastUpdateID+1, obUpdate.FirstUpdateId)
		}
		if obUpdate.LastUpdateId <= state.instrumentData.lastUpdateID {
			// We already have this update, snapshot is too new
			continue
		}

		bids, asks := obUpdate.ToBidAsk()
		for _, bid := range bids {
			ob.UpdateOrderBookLevel(bid)
		}
		for _, ask := range asks {
			ob.UpdateOrderBookLevel(ask)
		}
		//fmt.Println("SEQ", state.instrumentData.lastUpdateID+1, obUpdate.FirstUpdateId, obUpdate.LastUpdateId)
		//fmt.Println(ob)
		state.instrumentData.lastUpdateID = obUpdate.LastUpdateId
		state.instrumentData.lastUpdateTime = uint64(msg.ClientTime.UnixNano() / 1000000)
		synced = true
	}
	//fmt.Println("SYNCED")
	state.instrumentData.orderBook = ob
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())

	go func(sink *xutils.WebsocketSink, pid *actor.PID) {
		for {
			msg, err := sink.ReadMessage()
			if err == nil {
				context.Send(pid, msg)
			} else {
				return
			}
		}
	}(sink, context.Self())

	return nil
}

func (state *Listener) OnMarketDataRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketDataRequest)
	response := &messages.MarketDataResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		SeqNum:     state.instrumentData.seqNum,
		Success:    true,
	}

	if msg.Aggregation == models.OrderBookAggregation_L2 {
		snapshot := &models.OBL2Snapshot{
			Bids:          state.instrumentData.orderBook.GetBids(0),
			Asks:          state.instrumentData.orderBook.GetAsks(0),
			Timestamp:     utils.MilliToTimestamp(state.instrumentData.lastUpdateTime),
			TickPrecision: &wrapperspb.UInt64Value{Value: state.instrumentData.orderBook.TickPrecision},
			LotPrecision:  &wrapperspb.UInt64Value{Value: state.instrumentData.orderBook.LotPrecision},
		}
		response.SnapshotL2 = snapshot
	}

	context.Respond(response)
	return nil
}

func (state *Listener) onWebsocketMessage(context actor.Context) error {
	msg := context.Message().(*xchanger.WebsocketMessage)
	switch res := msg.Message.(type) {

	case error:
		return fmt.Errorf("socket error: %v", msg)

	case gatef.WSFuturesTrade:

		var aggregatefID = uint64(res.CreateTimeMs) * 10
		if res.Size < 0 {
			aggregatefID += 1
		}

		ts := uint64(msg.ClientTime.UnixNano() / 1000000)
		if state.instrumentData.aggTrade == nil || state.instrumentData.aggTrade.AggregateID != aggregatefID {
			if state.instrumentData.lastAggTradeTs >= ts {
				ts = state.instrumentData.lastAggTradeTs + 1
			}
			aggTrade := &models.AggregatedTrade{
				Bid:         res.Size < 0,
				Timestamp:   utils.MilliToTimestamp(ts),
				AggregateID: aggregatefID,
				Trades:      nil,
			}
			state.instrumentData.aggTrade = aggTrade
			state.instrumentData.lastAggTradeTs = ts

			// Stash the aggTrade
			state.stashedTrades.PushBack(aggTrade)
			// start the timer on trade creation, it will publish the trade in 20 ms
			go func(pid *actor.PID) {
				time.Sleep(20 * time.Millisecond)
				context.Send(pid, &postAggTrade{})
			}(context.Self())
		}

		state.instrumentData.aggTrade.Trades = append(
			state.instrumentData.aggTrade.Trades,
			&models.Trade{
				Price:    res.Price,
				Quantity: math.Abs(res.Size),
				ID:       res.Id,
			})

	case gatef.WSFuturesOrderBookUpdate:
		if state.sink == nil {
			return nil
		}
		symbol := res.Contract
		// Check depth continuity
		//fmt.Println("SEQ", state.instrumentData.lastUpdateID+1, res.FirstUpdateId)
		if res.LastUpdateId <= state.instrumentData.lastUpdateID {
			// We already have this update
			return nil
		}

		// Check depth continuity
		if state.instrumentData.lastUpdateID+1 != res.FirstUpdateId {
			return fmt.Errorf("got wrong sequence ID for %s: %d, %d",
				symbol, state.instrumentData.lastUpdateID, res.FirstUpdateId)
		}

		bids, asks := res.ToBidAsk()

		obDelta := &models.OBL2Update{
			Levels:    []*gmodels.OrderBookLevel{},
			Timestamp: utils.MilliToTimestamp(res.Time),
			Trade:     false,
		}

		for _, bid := range bids {
			obDelta.Levels = append(
				obDelta.Levels,
				bid,
			)
			state.instrumentData.orderBook.UpdateOrderBookLevel(bid)
		}

		for _, ask := range asks {
			obDelta.Levels = append(
				obDelta.Levels,
				ask,
			)
			state.instrumentData.orderBook.UpdateOrderBookLevel(ask)
		}

		//fmt.Println(state.instrumentData.orderBook)
		if state.instrumentData.orderBook.Crossed() {
			state.logger.Info("crossed orderbook", log.Error(errors.New("crossed")))
			return state.subscribeInstrument(context)
		}

		state.instrumentData.lastUpdateID = res.LastUpdateId
		state.instrumentData.lastUpdateTime = res.Time
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			UpdateL2: obDelta,
			SeqNum:   state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1

		return nil
	}

	return nil
}

func (state *Listener) checkSockets(context actor.Context) error {
	if state.sink.Error() != nil {
		state.logger.Info("error on ob socket", log.Error(state.sink.Error()))
		if err := state.subscribeInstrument(context); err != nil {
			return fmt.Errorf("error subscribing to instrument: %v", err)
		}
	}

	// If haven't sent anything for 2 seconds, send heartbeat
	if time.Since(state.instrumentData.lastHBTime) > 2*time.Second {
		// Send an empty refresh
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			SeqNum: state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1
		state.instrumentData.lastHBTime = time.Now()
	}

	return nil
}

func (state *Listener) postAggTrade(context actor.Context) {
	nowMilli := uint64(time.Now().UnixNano() / 1000000)

	for el := state.stashedTrades.Front(); el != nil; el = state.stashedTrades.Front() {
		trd := el.Value.(*models.AggregatedTrade)
		if trd != nil && nowMilli-utils.TimestampToMilli(trd.Timestamp) > 20 {
			context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
				SeqNum: state.instrumentData.seqNum + 1,
				Trades: []*models.AggregatedTrade{trd},
			})

			state.instrumentData.seqNum += 1

			//Check if the trd send is the last trade added to state
			if state.instrumentData.aggTrade == trd {
				state.instrumentData.aggTrade = nil
			}
			state.stashedTrades.Remove(el)
		} else {
			break
		}
	}
}
