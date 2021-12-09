package bitstamp

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
	"gitlab.com/alphaticks/xchanger"
	"gitlab.com/alphaticks/xchanger/exchanges/bitstamp"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"math"
	"reflect"
	"time"
)

type checkSockets struct{}
type postAggTrade struct{}

type InstrumentData struct {
	tickPrecision  uint64
	lotPrecision   uint64
	orderBook      *gorderbook.OrderBookL2
	seqNum         uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	aggTrade       *models.AggregatedTrade
	lastAggTradeTs uint64
}

// OBType: OBL2

type ListenerL2 struct {
	obWs            *bitstamp.Websocket
	tradeWs         *bitstamp.Websocket
	security        *models.Security
	dialerPool      *xchangerUtils.DialerPool
	instrumentData  *InstrumentData
	executorManager *actor.PID
	logger          *log.Logger
	lastPingTime    time.Time
	stashedTrades   *list.List
	socketTicker    *time.Ticker
}

func NewListenerL2Producer(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Producer {
	return func() actor.Actor {
		return NewListenerL2(security, dialerPool)
	}
}

func NewListenerL2(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Actor {
	return &ListenerL2{
		obWs:            nil,
		tradeWs:         nil,
		security:        security,
		dialerPool:      dialerPool,
		instrumentData:  nil,
		executorManager: nil,
		logger:          nil,
		stashedTrades:   nil,
	}
}

func (state *ListenerL2) Receive(context actor.Context) {
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

func (state *ListenerL2) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()),
		log.String("exchange", state.security.Exchange.Name),
		log.String("symbol", state.security.Symbol))

	state.lastPingTime = time.Now()
	state.stashedTrades = list.New()

	if state.security.MinPriceIncrement == nil || state.security.RoundLot == nil {
		return fmt.Errorf("security is missing MinPriceIncrement or RoundLot")
	}
	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot.Value))

	state.instrumentData = &InstrumentData{
		tickPrecision:  tickPrecision,
		lotPrecision:   lotPrecision,
		orderBook:      nil,
		seqNum:         uint64(time.Now().UnixNano()),
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

func (state *ListenerL2) Clean(context actor.Context) error {
	if state.tradeWs != nil {
		if err := state.tradeWs.Disconnect(); err != nil {
			state.logger.Info("error disconnecting socket", log.Error(err))
		}
	}
	if state.obWs != nil {
		if err := state.obWs.Disconnect(); err != nil {
			state.logger.Info("error disconnecting socket", log.Error(err))
		}
	}
	if state.socketTicker != nil {
		state.socketTicker.Stop()
		state.socketTicker = nil
	}

	return nil
}

func (state *ListenerL2) subscribeOrderBook(context actor.Context) error {
	if state.obWs != nil {
		_ = state.obWs.Disconnect()
	}

	ws := bitstamp.NewWebsocket()
	if err := ws.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to bitstamp websocket: %v", err)
	}

	if err := ws.Subscribe(state.security.Symbol, bitstamp.WSOrderBookChannel); err != nil {
		return fmt.Errorf("error subscribing to depth stream for symbol %s", state.security.Symbol)
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	_, ok := ws.Msg.Message.(bitstamp.WSSubscribedMessage)
	if !ok {
		return fmt.Errorf("was expecting WSSubsribed message, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	depthData, ok := ws.Msg.Message.(bitstamp.OrderBookL2)
	if !ok {
		return fmt.Errorf("was expecting depth data, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	bids, asks := depthData.ToBidAsk()
	ob := gorderbook.NewOrderBookL2(
		state.instrumentData.tickPrecision,
		state.instrumentData.lotPrecision,
		1000)
	ob.Sync(bids, asks)

	ts := uint64(ws.Msg.ClientTime.UnixNano()) / 1000

	state.instrumentData.seqNum = uint64(time.Now().UnixNano())
	state.instrumentData.lastUpdateTime = ts
	state.instrumentData.orderBook = ob

	state.obWs = ws

	go func(ws *bitstamp.Websocket, pid *actor.PID) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
		}
	}(ws, context.Self())

	return nil
}

func (state *ListenerL2) subscribeTrades(context actor.Context) error {
	if state.tradeWs != nil {
		_ = state.tradeWs.Disconnect()
	}
	ws := bitstamp.NewWebsocket()
	if err := ws.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to bitstamp websocket: %v", err)
	}

	if err := ws.Subscribe(state.security.Symbol, bitstamp.WSLiveTradesChannel); err != nil {
		return fmt.Errorf("error subscribing to trade stream for symbol %s", state.security.Symbol)
	}

	state.tradeWs = ws

	go func(ws *bitstamp.Websocket, pid *actor.PID) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
		}
	}(ws, context.Self())

	return nil
}

func (state *ListenerL2) OnMarketDataRequest(context actor.Context) error {
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
			Timestamp:     utils.MicroToTimestamp(state.instrumentData.lastUpdateTime),
			TickPrecision: &types.UInt64Value{Value: state.instrumentData.orderBook.TickPrecision},
			LotPrecision:  &types.UInt64Value{Value: state.instrumentData.orderBook.LotPrecision},
		}
		response.SnapshotL2 = snapshot
	}

	context.Respond(response)
	return nil
}

func (state *ListenerL2) onWebsocketMessage(context actor.Context) error {
	msg := context.Message().(*xchanger.WebsocketMessage)
	switch msg.Message.(type) {

	case error:
		return fmt.Errorf("socket error: %v", msg)

	case bitstamp.OrderBookL2:
		obData := msg.Message.(bitstamp.OrderBookL2)

		newOb := gorderbook.NewOrderBookL2(
			state.instrumentData.tickPrecision,
			state.instrumentData.lotPrecision,
			1000)

		bids, asks := obData.ToBidAsk()

		newOb.Sync(bids, asks)

		deltas := state.instrumentData.orderBook.Diff(newOb)

		ts := uint64(msg.ClientTime.UnixNano() / 1000)
		obDelta := &models.OBL2Update{
			Levels:    deltas,
			Timestamp: utils.MicroToTimestamp(ts),
			Trade:     false,
		}

		state.instrumentData.orderBook = newOb
		state.instrumentData.lastUpdateTime = ts

		if state.instrumentData.orderBook.Crossed() {
			state.logger.Info("crossed orderbook", log.Error(errors.New("crossed")))
			return state.subscribeOrderBook(context)
		}

		// Send OBData
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			UpdateL2: obDelta,
			SeqNum:   state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1

	case bitstamp.WSTrade:
		tradeData := msg.Message.(bitstamp.WSTrade)
		tradeData.MicroTimestamp = uint64(msg.ClientTime.UnixNano()) / 1000
		ts := tradeData.MicroTimestamp / 1000

		var aggID uint64
		if tradeData.Type == 1 {
			aggID = tradeData.SellOrderID
		} else {
			aggID = tradeData.BuyOrderID
		}

		if state.instrumentData.aggTrade == nil || state.instrumentData.aggTrade.AggregateID != aggID {
			if state.instrumentData.lastAggTradeTs >= ts {
				ts = state.instrumentData.lastAggTradeTs + 1
			}
			aggTrade := &models.AggregatedTrade{
				Bid:         tradeData.Type == 1,
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
				Price:    tradeData.Price,
				Quantity: tradeData.Amount,
				ID:       tradeData.ID,
			})
	}

	return nil
}

func (state *ListenerL2) checkSockets(context actor.Context) error {

	if time.Now().Sub(state.lastPingTime) > 10*time.Second {
		// "Ping" by resubscribing to the topic
		_ = state.obWs.Subscribe(state.security.Symbol, bitstamp.WSOrderBookChannel)
		_ = state.tradeWs.Subscribe(state.security.Symbol, bitstamp.WSLiveTradesChannel)
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

func (state *ListenerL2) postAggTrade(context actor.Context) {
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
