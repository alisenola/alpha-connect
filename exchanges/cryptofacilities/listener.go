package cryptofacilities

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/gorderbook"
	gmodels "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	"gitlab.com/alphaticks/xchanger"
	"gitlab.com/alphaticks/xchanger/exchanges/cryptofacilities"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"math"
	"reflect"
	"time"
)

type checkSockets struct{}
type postAggTrade struct{}

type InstrumentData struct {
	orderBook      *gorderbook.OrderBookL2
	seqNum         uint64
	lastUpdateID   uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	aggTrade       *models.AggregatedTrade
	lastAggTradeTs uint64
}

// OBType: OBL3
// No ID for the deltas..

type Listener struct {
	obWs           *cryptofacilities.Websocket
	tradeWs        *cryptofacilities.Websocket
	security       *models.Security
	dialerPool     *xchangerUtils.DialerPool
	instrumentData *InstrumentData
	logger         *log.Logger
	lastPingTime   time.Time
	stashedTrades  *list.List
	socketTicker   *time.Ticker
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

	state.instrumentData = &InstrumentData{
		orderBook:      nil,
		seqNum:         uint64(time.Now().UnixNano()),
		lastUpdateID:   0,
		lastUpdateTime: 0,
		lastHBTime:     time.Now(),
		aggTrade:       nil,
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
	if state.tradeWs != nil {
		if err := state.tradeWs.Disconnect(); err != nil {
			state.logger.Warn("error disconnecting socket", log.Error(err))
		}
	}
	if state.obWs != nil {
		if err := state.obWs.Disconnect(); err != nil {
			state.logger.Warn("error disconnecting socket", log.Error(err))
		}
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

	ws := cryptofacilities.NewWebsocket()
	if err := ws.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to websocket: %v", err)
	}

	if err := ws.Subscribe([]string{state.security.Symbol}, cryptofacilities.WSBookFeed); err != nil {
		return fmt.Errorf("error subscribing to OBL2 stream: %v", err)
	}
	if err := ws.Subscribe(nil, cryptofacilities.WSHeartBeatFeed); err != nil {
		return fmt.Errorf("error subscribing to WSHeartBeatFeed stream: %v", err)
	}

	trials := 0
	synced := false
	for !synced && trials < 4 {
		if !ws.ReadMessage() {
			return fmt.Errorf("error reading message: %v", ws.Err)
		}
		obData, ok := ws.Msg.Message.(cryptofacilities.WSBookSnapshot)
		if !ok {
			trials += 1
			continue
		}

		var bids, asks []gmodels.OrderBookLevel
		bids = make([]gmodels.OrderBookLevel, len(obData.Bids))
		for i, bid := range obData.Bids {
			bids[i] = gmodels.OrderBookLevel{
				Price:    bid.Price,
				Quantity: bid.Qty,
				Bid:      true,
			}
		}
		asks = make([]gmodels.OrderBookLevel, len(obData.Asks))
		for i, ask := range obData.Asks {
			asks[i] = gmodels.OrderBookLevel{
				Price:    ask.Price,
				Quantity: ask.Qty,
				Bid:      false,
			}
		}

		ts := uint64(ws.Msg.ClientTime.UnixNano()) / 1000000

		tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))
		lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot.Value))
		ob := gorderbook.NewOrderBookL2(
			tickPrecision,
			lotPrecision,
			10000)

		ob.Sync(bids, asks)
		if ob.Crossed() {
			return fmt.Errorf("crossed orderbook")
		}
		state.instrumentData.orderBook = ob
		state.instrumentData.lastUpdateID = obData.Seq
		state.instrumentData.lastUpdateTime = ts
		synced = true
	}
	if !synced {
		return fmt.Errorf("error getting WSBookSnapshot")
	}

	state.instrumentData.seqNum = uint64(time.Now().UnixNano())

	state.obWs = ws

	go func(ws *cryptofacilities.Websocket, pid *actor.PID) {
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

	ws := cryptofacilities.NewWebsocket()
	if err := ws.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to websocket: %v", err)
	}

	if err := ws.Subscribe([]string{state.security.Symbol}, cryptofacilities.WSTradeFeed); err != nil {
		return fmt.Errorf("error subscribing to trade stream: %v", err)
	}
	if err := ws.Subscribe(nil, cryptofacilities.WSHeartBeatFeed); err != nil {
		return fmt.Errorf("error subscribing to WSHeartBeatFeed stream: %v", err)
	}

	state.tradeWs = ws
	go func(ws *cryptofacilities.Websocket, pid *actor.PID) {
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
	switch msg.Message.(type) {

	case error:
		return fmt.Errorf("socket error: %v", msg)

	case cryptofacilities.WSAlert:
		err := msg.Message.(cryptofacilities.WSAlert)
		return fmt.Errorf("socket error: %v", err)

	case cryptofacilities.WSBookDelta:
		if state.obWs == nil || msg.WSID != state.obWs.ID {
			return nil
		}
		obData := msg.Message.(cryptofacilities.WSBookDelta)
		instr := state.instrumentData

		ts := uint64(msg.ClientTime.UnixNano()) / 1000000

		if obData.Seq != instr.lastUpdateID+1 {
			state.logger.Info("got inconsistent sequence")
			return state.subscribeOrderBook(context)
		}
		instr.lastUpdateID = obData.Seq

		level := gmodels.OrderBookLevel{
			Price:    obData.Price,
			Quantity: obData.Qty,
			Bid:      obData.Side == "buy",
		}

		instr.orderBook.UpdateOrderBookLevel(level)

		obDelta := &models.OBL2Update{
			Levels:    []gmodels.OrderBookLevel{level},
			Timestamp: utils.MilliToTimestamp(ts),
			Trade:     false,
		}

		instr.lastUpdateTime = ts

		if state.instrumentData.orderBook.Crossed() {
			state.logger.Warn("crossed orderbook", log.Error(errors.New("crossed")))
			return state.subscribeOrderBook(context)
		}

		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			UpdateL2: obDelta,
			SeqNum:   instr.seqNum + 1,
		})
		instr.seqNum += 1

	case cryptofacilities.WSTrade:
		trade := msg.Message.(cryptofacilities.WSTrade)

		aggID := trade.Time * 10
		// Add one to aggregatedID if it's a sell so that
		// buy and sell happening at the same time won't have the same ID
		if trade.Side == "sell" {
			aggID += 1
		}
		ts := uint64(msg.ClientTime.UnixNano()) / 1000000
		if state.instrumentData.aggTrade == nil || state.instrumentData.aggTrade.AggregateID != aggID {
			// ensure increasing timestamp
			if ts <= state.instrumentData.lastAggTradeTs {
				ts = state.instrumentData.lastAggTradeTs + 1
			}
			aggTrade := &models.AggregatedTrade{
				Bid:         trade.Side == "sell",
				Timestamp:   utils.MilliToTimestamp(ts),
				AggregateID: aggID,
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
				Price:    trade.Price,
				Quantity: trade.Qty,
				ID:       trade.Seq,
			})
	}

	return nil
}

func (state *Listener) checkSockets(context actor.Context) error {

	/*
		if time.Now().Sub(state.lastPingTime) > 10*time.Second {
			// "Ping" by resubscribing to the topic
			_ = state.tradeWs.Subscribe(nil, cryptofacilities.WSHeartBeatFeed)
			_ = state.obWs.Subscribe(nil, cryptofacilities.WSHeartBeatFeed)
			state.lastPingTime = time.Now()
		}

	*/
	if state.obWs.Err != nil || !state.obWs.Connected {
		if state.obWs.Err != nil {
			state.logger.Warn("error on socket", log.Error(state.obWs.Err))
		}
		if err := state.subscribeOrderBook(context); err != nil {
			return fmt.Errorf("error subscribing to instrument: %v", err)
		}
	}

	if state.tradeWs.Err != nil || !state.tradeWs.Connected {
		if state.tradeWs.Err != nil {
			state.logger.Warn("error on socket", log.Error(state.tradeWs.Err))
		}
		if err := state.subscribeTrades(context); err != nil {
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
				Trades: []*models.AggregatedTrade{trd},
				SeqNum: state.instrumentData.seqNum + 1,
			})
			state.instrumentData.seqNum += 1
			// At this point, the state.instrumentData.aggTrade can be our trade, or it can be a new one
			if state.instrumentData.aggTrade == trd {
				state.instrumentData.aggTrade = nil
			}
			state.stashedTrades.Remove(el)
		} else {
			break
		}
	}
}
