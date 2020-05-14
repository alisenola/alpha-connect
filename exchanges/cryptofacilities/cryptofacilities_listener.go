package cryptofacilities

import (
	"container/list"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/exchanges/cryptofacilities"
	"math"
	"reflect"
	"time"
)

type readSocket struct{}
type postAggTrade struct{}

type OBL2Request struct {
	requester *actor.PID
	requestID int64
}

type InstrumentData struct {
	orderBook      *gorderbook.OrderBookL2
	seqNum         uint64
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
	wsChan         chan *cryptofacilities.WebsocketMessage
	security       *models.Security
	instrumentData *InstrumentData
	logger         *log.Logger
	lastPingTime   time.Time
	stashedTrades  *list.List
}

func NewListenerProducer(security *models.Security) actor.Producer {
	return func() actor.Actor {
		return NewListener(security)
	}
}

func NewListener(security *models.Security) actor.Actor {
	return &Listener{
		obWs:           nil,
		tradeWs:        nil,
		wsChan:         nil,
		security:       security,
		instrumentData: nil,
		logger:         nil,
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

	state.wsChan = make(chan *cryptofacilities.WebsocketMessage, 10000)
	state.lastPingTime = time.Now()
	state.stashedTrades = list.New()

	context.Send(context.Self(), &readSocket{})

	state.instrumentData = &InstrumentData{
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

	ws := cryptofacilities.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to websocket: %v", err)
	}

	if err := ws.Subscribe([]string{state.security.Symbol}, cryptofacilities.WSBookFeed); err != nil {
		return fmt.Errorf("error subscribing to OBL2 stream: %v", err)
	}
	if err := ws.Subscribe(nil, cryptofacilities.WSHeartBeatFeed); err != nil {
		return fmt.Errorf("error subscribing to WSHeartBeatFeed stream: %v", err)
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	_, ok := ws.Msg.Message.(cryptofacilities.WSInfo)
	if !ok {
		return fmt.Errorf("was expecting WSInfo, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	_, ok = ws.Msg.Message.(cryptofacilities.WSSubscribeResponse)
	if !ok {
		return fmt.Errorf("was expecting WSSubscribeResponse, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}
	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	_, ok = ws.Msg.Message.(cryptofacilities.WSSubscribeResponse)
	if !ok {
		return fmt.Errorf("was expecting WSSubscribeResponse, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	obData, ok := ws.Msg.Message.(cryptofacilities.WSBookSnapshot)
	if !ok {
		return fmt.Errorf("was expecting WSBookSnapshot, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	var bids, asks []gorderbook.OrderBookLevel
	bids = make([]gorderbook.OrderBookLevel, len(obData.Bids), len(obData.Bids))
	for i, bid := range obData.Bids {
		bids[i] = gorderbook.OrderBookLevel{
			Price:    bid.Price,
			Quantity: bid.Qty,
			Bid:      true,
		}
	}
	asks = make([]gorderbook.OrderBookLevel, len(obData.Asks), len(obData.Asks))
	for i, ask := range obData.Asks {
		asks[i] = gorderbook.OrderBookLevel{
			Price:    ask.Price,
			Quantity: ask.Qty,
			Bid:      false,
		}
	}

	ts := uint64(ws.Msg.Time.UnixNano()) / 1000000

	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot))
	ob := gorderbook.NewOrderBookL2(
		tickPrecision,
		lotPrecision,
		10000)

	ob.Sync(bids, asks)
	state.instrumentData.orderBook = ob
	state.instrumentData.seqNum = obData.Seq
	state.instrumentData.lastUpdateTime = ts

	state.obWs = ws

	go func(ws *cryptofacilities.Websocket) {
		for ws.ReadMessage() {
			state.wsChan <- ws.Msg
		}
	}(ws)

	return nil
}

func (state *Listener) subscribeTrades(context actor.Context) error {
	if state.tradeWs != nil {
		_ = state.tradeWs.Disconnect()
	}

	ws := cryptofacilities.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to websocket: %v", err)
	}

	if err := ws.Subscribe([]string{state.security.Symbol}, cryptofacilities.WSTradeFeed); err != nil {
		return fmt.Errorf("error subscribing to trade stream: %v", err)
	}
	if err := ws.Subscribe(nil, cryptofacilities.WSHeartBeatFeed); err != nil {
		return fmt.Errorf("error subscribing to WSHeartBeatFeed stream: %v", err)
	}

	state.tradeWs = ws

	go func(ws *cryptofacilities.Websocket) {
		for ws.ReadMessage() {
			state.wsChan <- ws.Msg
		}
	}(ws)

	return nil
}

func (state *Listener) OnMarketDataRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketDataRequest)

	response := &messages.MarketDataSnapshot{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
	}

	if msg.Aggregation == messages.L2 {
		snapshot := &models.OBL2Snapshot{
			Bids:      state.instrumentData.orderBook.GetBids(0),
			Asks:      state.instrumentData.orderBook.GetAsks(0),
			Timestamp: utils.MilliToTimestamp(state.instrumentData.lastUpdateTime),
			SeqNum:    state.instrumentData.seqNum,
		}
		response.SnapshotL2 = snapshot
	}
	context.Respond(response)
	return nil
}

func (state *Listener) readSocket(context actor.Context) error {
	select {

	case msg := <-state.wsChan:
		switch msg.Message.(type) {

		case error:
			return fmt.Errorf("socket error: %v", msg)

		case cryptofacilities.WSError:
			err := msg.Message.(cryptofacilities.WSError)
			return fmt.Errorf("socket error: %v", err)

		case cryptofacilities.WSBookDelta:
			obData := msg.Message.(cryptofacilities.WSBookDelta)
			instr := state.instrumentData

			ts := uint64(msg.Time.UnixNano()) / 1000000

			if obData.Seq != instr.seqNum+1 {
				state.logger.Info("got inconsistent sequence")
				// Stop the socket, we will restart instrument at the end
				if err := state.obWs.Disconnect(); err != nil {
					state.logger.Info("error disconnecting from socket", log.Error(err))
				}
			}

			level := gorderbook.OrderBookLevel{
				Price:    obData.Price,
				Quantity: obData.Qty,
				Bid:      obData.Side == "buy",
			}

			instr.orderBook.UpdateOrderBookLevel(level)

			obDelta := &models.OBL2Update{
				Levels:    []gorderbook.OrderBookLevel{level},
				Timestamp: utils.MilliToTimestamp(ts),
				SeqNum:    instr.seqNum + 1,
				Trade:     false,
			}
			instr.seqNum += 1

			instr.lastUpdateTime = ts

			if state.instrumentData.orderBook.Crossed() {
				state.logger.Info("crossed order book")
				// Stop the socket, we will restart instrument at the end
				if err := state.obWs.Disconnect(); err != nil {
					state.logger.Info("error disconnecting from socket", log.Error(err))
				}
				break
			}

			context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
				UpdateL2: obDelta,
			})

		case cryptofacilities.WSTrade:
			trade := msg.Message.(cryptofacilities.WSTrade)

			aggID := trade.Time * 10
			// Add one to aggregatedID if it's a sell so that
			// buy and sell happening at the same time won't have the same ID
			if trade.Side == "sell" {
				aggID += 1
			}
			ts := uint64(msg.Time.UnixNano()) / 1000000
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
