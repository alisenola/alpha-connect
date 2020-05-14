package bitstamp

import (
	"container/list"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/exchanges/bitstamp"
	"math"
	"reflect"
	"time"
)

type readSocket struct{}
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

type Listener struct {
	obWs            *bitstamp.Websocket
	tradeWs         *bitstamp.Websocket
	wsChan          chan *bitstamp.WebsocketMessage
	security        *models.Security
	instrumentData  *InstrumentData
	executorManager *actor.PID
	logger          *log.Logger
	lastPingTime    time.Time
	stashedTrades   *list.List
}

func NewListenerProducer(security *models.Security) actor.Producer {
	return func() actor.Actor {
		return NewListener(security)
	}
}

func NewListener(security *models.Security) actor.Actor {
	return &Listener{
		obWs:            nil,
		tradeWs:         nil,
		wsChan:          nil,
		security:        security,
		instrumentData:  nil,
		executorManager: nil,
		logger:          nil,
		stashedTrades:   nil,
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

	case *readSocket:
		if err := state.readSocket(context); err != nil {
			state.logger.Error("error processing readSocket", log.Error(err))
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

	state.wsChan = make(chan *bitstamp.WebsocketMessage, 10000)
	state.lastPingTime = time.Now()
	state.stashedTrades = list.New()

	context.Send(context.Self(), &readSocket{})
	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot))

	state.instrumentData = &InstrumentData{
		tickPrecision:  tickPrecision,
		lotPrecision:   lotPrecision,
		orderBook:      nil,
		seqNum:         0,
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

	return nil
}

func (state *Listener) Clean(context actor.Context) error {
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

	return nil
}

func (state *Listener) subscribeOrderBook(context actor.Context) error {
	if state.obWs != nil {
		_ = state.obWs.Disconnect()
	}

	ws := bitstamp.NewWebsocket()
	if err := ws.Connect(); err != nil {
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

	ts := uint64(ws.Msg.Time.UnixNano()) / 1000

	state.instrumentData.seqNum = 0
	state.instrumentData.lastUpdateTime = ts
	state.instrumentData.orderBook = ob

	state.obWs = ws

	go func(ws *bitstamp.Websocket) {
		for ws.ReadMessage() {
			state.wsChan <- ws.Msg
		}
	}(state.obWs)

	return nil
}

func (state *Listener) subscribeTrades(context actor.Context) error {
	if state.tradeWs != nil {
		_ = state.tradeWs.Disconnect()
	}
	ws := bitstamp.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to bitstamp websocket: %v", err)
	}

	if err := ws.Subscribe(state.security.Symbol, bitstamp.WSLiveTradesChannel); err != nil {
		return fmt.Errorf("error subscribing to trade stream for symbol %s", state.security.Symbol)
	}

	state.tradeWs = ws

	go func(ws *bitstamp.Websocket) {
		for ws.ReadMessage() {
			state.wsChan <- ws.Msg
		}
	}(state.tradeWs)

	return nil
}

func (state *Listener) OnMarketDataRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketDataRequest)

	snapshot := &models.OBL2Snapshot{
		Bids:      state.instrumentData.orderBook.GetBids(0),
		Asks:      state.instrumentData.orderBook.GetAsks(0),
		Timestamp: utils.MicroToTimestamp(state.instrumentData.lastUpdateTime),
		SeqNum:    state.instrumentData.seqNum,
	}
	context.Respond(&messages.MarketDataSnapshot{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		SnapshotL2: snapshot,
	})
	return nil
}

func (state *Listener) readSocket(context actor.Context) error {
	select {
	case msg := <-state.wsChan:
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

			ts := uint64(msg.Time.UnixNano() / 1000)
			obDelta := &models.OBL2Update{
				Levels:    deltas,
				Timestamp: utils.MicroToTimestamp(ts),
				SeqNum:    state.instrumentData.seqNum + 1,
				Trade:     false,
			}
			state.instrumentData.seqNum += 1

			state.instrumentData.orderBook = newOb
			state.instrumentData.lastUpdateTime = ts

			if state.instrumentData.orderBook.Crossed() {
				state.logger.Info("crossed order book")
				// Stop the socket, we will restart instrument at the end
				if err := state.obWs.Disconnect(); err != nil {
					state.logger.Info("error disconnecting from socket", log.Error(err))
				}
				break
			}

			// Send OBData
			context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
				UpdateL2: obDelta,
			})

		case bitstamp.WSTrade:
			tradeData := msg.Message.(bitstamp.WSTrade)
			tradeData.MicroTimestamp = uint64(msg.Time.UnixNano()) / 1000
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

		if err := state.checkSockets(context); err != nil {
			return fmt.Errorf("error checking sockets: %v", err)
		}
		state.postHeartBeat(context)
		context.Send(context.Self(), &readSocket{})
		return nil

	case <-time.After(1 * time.Second):
		if err := state.checkSockets(context); err != nil {
			return fmt.Errorf("error checking sockets: %v", err)
		}
		state.postHeartBeat(context)
		context.Send(context.Self(), &readSocket{})
		return nil
	}
}

func (state *Listener) checkSockets(context actor.Context) error {

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

	return nil
}

func (state *Listener) postHeartBeat(context actor.Context) {
	// If haven't sent anything for 2 seconds, send heartbeat
	if time.Now().Sub(state.instrumentData.lastHBTime) > 2*time.Second {
		// Send an empty refresh
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{})
		state.instrumentData.lastHBTime = time.Now()
	}
}

func (state *Listener) postAggTrade(context actor.Context) {
	nowMilli := uint64(time.Now().UnixNano() / 1000000)

	for el := state.stashedTrades.Front(); el != nil; el = state.stashedTrades.Front() {
		trd := el.Value.(*models.AggregatedTrade)
		if trd != nil && nowMilli-utils.TimestampToMilli(trd.Timestamp) > 20 {
			context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
				Trades: []*models.AggregatedTrade{trd},
			})
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
