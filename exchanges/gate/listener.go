package gate

import (
	"container/list"
	"errors"
	"fmt"
	gmodels "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	"math"
	"reflect"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/gate"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
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
	obWs           *gate.Websocket
	tradeWs        *gate.Websocket
	security       *models.Security
	dialerPool     *xchangerUtils.DialerPool
	instrumentData *InstrumentData
	logger         *log.Logger
	lastPingTime   time.Time
	socketTicker   *time.Ticker
	gateExecutor   *actor.PID
	stashedTrades  *list.List
}

func NewListenerProducer(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Producer {
	return func() actor.Actor {
		return NewListener(security, dialerPool)
	}
}

func NewListener(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Actor {
	return &Listener{
		obWs:           nil,
		tradeWs:        nil,
		security:       security,
		dialerPool:     dialerPool,
		instrumentData: nil,
		logger:         nil,
		socketTicker:   nil,
		stashedTrades:  nil,
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
		log.String("exchange", state.security.Exchange.Name),
		log.String("symbol", state.security.Symbol))

	if state.security.MinPriceIncrement == nil || state.security.RoundLot == nil {
		return fmt.Errorf("security is missing MinPriceIncrement or RoundLot")
	}
	state.lastPingTime = time.Now()
	state.stashedTrades = list.New()
	state.gateExecutor = actor.NewPID(context.ActorSystem().Address(), "executor/"+constants.GATE.Name+"_executor")

	state.instrumentData = &InstrumentData{
		orderBook:      nil,
		seqNum:         uint64(time.Now().UnixNano()),
		lastUpdateTime: 0,
		lastHBTime:     time.Now(),
		lastAggTradeTs: 0,
	}

	if err := state.subscribeOrderBook(context); err != nil {
		return fmt.Errorf("error subscribing to order book: %v", err)
	}
	if err := state.subscribeTrades(context); err != nil {
		return fmt.Errorf("error subscribing to trades: %v", err)
	}

	socketTicker := time.NewTicker(5 * time.Second)
	state.socketTicker = socketTicker
	go func(pid *actor.PID) {
		for {
			select {
			case _ = <-socketTicker.C:
				context.Send(pid, &checkSockets{})
			case <-time.After(10 * time.Second):
				// timer stopped, we leave
				return
			}
		}
	}(context.Self())

	return nil
}

func (state *Listener) Clean(context actor.Context) error {
	if state.tradeWs != nil {
		if err := state.tradeWs.Disconnect(); err != nil {
			state.logger.Info("error disconnecting socket", log.Error(err))
		}
		state.tradeWs = nil
	}
	if state.obWs != nil {
		if err := state.obWs.Disconnect(); err != nil {
			state.logger.Info("error disconnecting socket", log.Error(err))
		}
		state.obWs = nil
	}
	if state.socketTicker != nil {
		state.socketTicker.Stop()
		state.socketTicker = nil
	}

	return nil
}

func (state *Listener) subscribeOrderBook(context actor.Context) error {
	if state.obWs != nil {
		_ = state.obWs.Disconnect()
	}

	ws := gate.NewWebsocket()
	if err := ws.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to gate websocket: %v", err)
	}

	if err := ws.SubscribeOrderBookUpdate(state.security.Symbol, gate.FREQ_100MS); err != nil {
		return fmt.Errorf("error subscribing to the order book: %v", err)
	}
	state.obWs = ws

	time.Sleep(35 * time.Second)
	fut := context.RequestFuture(
		state.gateExecutor,
		&messages.MarketDataRequest{
			RequestID: uint64(time.Now().UnixNano()),
			Subscribe: false,
			Instrument: &models.Instrument{
				SecurityID: &types.UInt64Value{Value: state.security.SecurityID},
				Symbol:     &types.StringValue{Value: state.security.Symbol},
				Exchange:   state.security.Exchange,
			},
			Aggregation: models.L2,
		},
		20*time.Second)

	res, err := fut.Result()

	if err != nil {
		return fmt.Errorf("error getting OBL2")
	}

	msg, ok := res.(*messages.MarketDataResponse)
	if !ok {
		return fmt.Errorf("was expecting MarketDataResponse, got %s", reflect.TypeOf(msg).String())
	}
	if !msg.Success {
		return fmt.Errorf("error fetching the snapshot %s", msg.RejectionReason.String())
	}
	if msg.SnapshotL2 == nil {
		return fmt.Errorf("market data snapshot has no OBL2")
	}
	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot.Value))
	depth := 10000
	if len(msg.SnapshotL2.Asks) > 0 {
		bestAsk := msg.SnapshotL2.Asks[0].Price
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

	ob.Sync(msg.SnapshotL2.Bids, msg.SnapshotL2.Asks)
	if ob.Crossed() {
		return fmt.Errorf("crossed orderbook")
	}
	state.instrumentData.lastUpdateID = msg.SeqNum
	state.instrumentData.lastUpdateTime = utils.TimestampToMilli(msg.SnapshotL2.Timestamp)
	ws.ReadMessage()
	if !ws.ReadMessage() {
		return fmt.Errorf("error reading first message: %v", ws.Err)
	}
	firstMsg, ok := ws.Msg.Message.(gate.WSSpotOrderBookUpdate)
	if !ok {
		return fmt.Errorf("incorrect message type for first message: %v", reflect.TypeOf(ws.Msg.Message).String())
	}
	hb := false
	sync := false
	for !sync {
		if !hb {
			if err := ws.Ping(); err != nil {
				return fmt.Errorf("error ping request: %v", err)
			}
			hb = true
		}
		if !ws.ReadMessage() {
			return fmt.Errorf("error reading message: %v", ws.Err)
		}
		obUpdate, ok := ws.Msg.Message.(gate.WSSpotOrderBookUpdate)
		if !ok {
			if _, ok := ws.Msg.Message.(gate.WSPong); ok {
				hb = false
			}

			context.Send(context.Self(), ws.Msg)
		}
		if obUpdate.LastUpdateId < msg.SeqNum+1 {
			continue
		}
		if msg.SeqNum+1 < firstMsg.FirstUpdateId {
			return fmt.Errorf("error in order book ID %d %d", msg.SeqNum, firstMsg.FirstUpdateId)
		}

		bids, asks := obUpdate.ToBidAsk()
		for _, bid := range bids {
			ob.UpdateOrderBookLevel(bid)
		}
		for _, ask := range asks {
			ob.UpdateOrderBookLevel(ask)
		}
		state.instrumentData.lastUpdateID = obUpdate.LastUpdateId
		state.instrumentData.lastUpdateTime = uint64(ws.Msg.ClientTime.UnixNano() / 1000)
		sync = true
	}
	state.instrumentData.orderBook = ob
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())
	state.obWs = ws

	go func(ws *gate.Websocket, pid *actor.PID) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
		}
	}(ws, context.Self())

	return nil
}

func (state *Listener) subscribeTrades(context actor.Context) error {
	if state.tradeWs != nil {
		_ = state.tradeWs.Disconnect()
	}

	ws := gate.NewWebsocket()
	if err := ws.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to gate websocket: %v", err)
	}

	if err := ws.Subscribe([]string{state.security.Symbol}, gate.WSSpotTradesChannel); err != nil {
		return fmt.Errorf("error subscribing to the trades: %v", err)
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}

	_, ok := ws.Msg.Message.(gate.WSSubscription)
	if !ok {
		return fmt.Errorf("expected trade data, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	state.tradeWs = ws

	go func(ws *gate.Websocket, pid *actor.PID) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
		}
	}(ws, context.Self())

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

	if msg.Aggregation == models.L2 {
		snapshot := &models.OBL2Snapshot{
			Bids:          state.instrumentData.orderBook.GetBids(0),
			Asks:          state.instrumentData.orderBook.GetAsks(0),
			Timestamp:     utils.MilliToTimestamp(state.instrumentData.lastUpdateTime),
			TickPrecision: &types.UInt64Value{Value: state.instrumentData.orderBook.TickPrecision},
			LotPrecision:  &types.UInt64Value{Value: state.instrumentData.orderBook.LotPrecision},
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

	case gate.WSSpotTrade:

		var aggregateID = uint64(res.CreateTimeMs) * 10
		if res.Side == "sell" {
			aggregateID += 1
		}

		ts := uint64(msg.ClientTime.UnixMilli())
		if state.instrumentData.aggTrade == nil || state.instrumentData.aggTrade.AggregateID != aggregateID {
			if state.instrumentData.lastAggTradeTs >= ts {
				ts = state.instrumentData.lastAggTradeTs + 1
			}
			aggTrade := &models.AggregatedTrade{
				Bid:         res.Side == "sell",
				Timestamp:   utils.MilliToTimestamp(ts),
				AggregateID: aggregateID,
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
			models.Trade{
				Price:    res.Price,
				Quantity: res.Amount,
				ID:       res.Id,
			})

	case gate.WSSpotOrderBookUpdate:
		if state.obWs == nil || msg.WSID != state.obWs.ID {
			return nil
		}
		symbol := res.CurrencyPair
		// Check depth continuity
		if state.instrumentData.lastUpdateID+1 != res.FirstUpdateId {
			return fmt.Errorf("got wrong sequence ID for %s: %d, %d",
				symbol, state.instrumentData.lastUpdateID, res.FirstUpdateId)
		}

		bids, asks := res.ToBidAsk()

		obDelta := &models.OBL2Update{
			Levels:    []gmodels.OrderBookLevel{},
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

		if state.instrumentData.orderBook.Crossed() {
			state.logger.Info("crossed orderbook", log.Error(errors.New("crossed")))
			return state.subscribeOrderBook(context)
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
	if time.Now().Sub(state.lastPingTime) > 5*time.Second {
		_ = state.obWs.Ping()
		_ = state.tradeWs.Ping()
		state.lastPingTime = time.Now()
	}

	if state.obWs.Err != nil || !state.obWs.Connected {
		if state.obWs.Err != nil {
			state.logger.Info("error on socket", log.Error(state.obWs.Err))
		}
		if err := state.subscribeOrderBook(context); err != nil {
			return fmt.Errorf("error subscribing to instrument: %v", err)
		}
	}

	if state.tradeWs.Err != nil || !state.tradeWs.Connected {
		if state.tradeWs.Err != nil {
			state.logger.Info("error on socket", log.Error(state.tradeWs.Err))
		}
		if err := state.subscribeTrades(context); err != nil {
			return fmt.Errorf("error subscribing to instrument: %v", err)
		}
	}

	// If haven't sent anything for 2 seconds, send heartbeat
	if time.Now().Sub(state.instrumentData.lastHBTime) > 2*time.Second {
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
	nowMilli := uint64(time.Now().UnixMilli())

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
