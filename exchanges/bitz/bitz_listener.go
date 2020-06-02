package bitz

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/exchanges/bitz"
	"math"
	"reflect"
	"sort"
	"time"
)

type readSocket struct{}

type InstrumentData struct {
	orderBook      *gorderbook.OrderBookL2
	seqNum         uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	lastAggTradeTs uint64
}

type Listener struct {
	obWs           *bitz.Websocket
	tradeWs        *bitz.Websocket
	wsChan         chan *bitz.WebsocketMessage
	security       *models.Security
	instrumentData *InstrumentData
	logger         *log.Logger
	lastPingTime   time.Time
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

	state.wsChan = make(chan *bitz.WebsocketMessage, 10000)
	state.lastPingTime = time.Now()

	context.Send(context.Self(), &readSocket{})

	state.instrumentData = &InstrumentData{
		orderBook:      nil,
		seqNum:         0,
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

	ws := bitz.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to bitz websocket: %v", err)
	}

	if err := ws.Subscribe(state.security.Symbol, []string{bitz.WSDepthTopic}); err != nil {
		return fmt.Errorf("error subscribing to depth stream for symbol %s", state.security.Symbol)
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	depthData, ok := ws.Msg.Message.(bitz.WSDepth)
	if !ok {
		return fmt.Errorf("was expecting depth data, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	bids, asks := depthData.ToBidAsk()

	//bestAsk := float64(asks[0].Price) / float64(state.instruments[obData.Symbol].instrument.TickPrecision)
	// Allow a 10% price variation
	//depth := int(((bestAsk * 1.1) - bestAsk) * float64(state.instruments[obData.Symbol].instrument.TickPrecision))
	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot))
	ob := gorderbook.NewOrderBookL2(
		tickPrecision,
		lotPrecision,
		10000)
	ob.Sync(bids, asks)

	ts := uint64(ws.Msg.Time.UnixNano() / 1000000)
	state.instrumentData.orderBook = ob
	state.instrumentData.lastUpdateTime = ts
	state.instrumentData.seqNum = 0

	state.obWs = ws

	go func(ws *bitz.Websocket) {
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

	ws := bitz.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to bitz websocket: %v", err)
	}

	if err := ws.Subscribe(state.security.Symbol, []string{bitz.WSOrderTopic}); err != nil {
		return fmt.Errorf("error subscribing to trade stream for symbol %s", state.security.Symbol)
	}

	state.tradeWs = ws

	go func(ws *bitz.Websocket) {
		for ws.ReadMessage() {
			state.wsChan <- ws.Msg
		}
	}(ws)

	return nil
}

func (state *Listener) OnMarketDataRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketDataRequest)

	snapshot := &models.OBL2Snapshot{
		Bids:      state.instrumentData.orderBook.GetBids(0),
		Asks:      state.instrumentData.orderBook.GetAsks(0),
		Timestamp: utils.MilliToTimestamp(state.instrumentData.lastUpdateTime),
	}
	context.Respond(&messages.MarketDataSnapshot{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		SnapshotL2: snapshot,
		SeqNum:     state.instrumentData.seqNum,
	})
	return nil
}

func (state *Listener) readSocket(context actor.Context) error {
	select {
	case msg := <-state.wsChan:
		switch msg.Message.(type) {

		case error:
			return fmt.Errorf("socket error: %v", msg)

		case bitz.WSDepth:
			obData := msg.Message.(bitz.WSDepth)

			instr := state.instrumentData

			bids, asks := obData.ToBidAsk()
			levels := append(bids, asks...)

			for _, l := range levels {
				instr.orderBook.UpdateOrderBookLevel(l)
			}

			ts := uint64(msg.Time.UnixNano() / 1000000)

			obDelta := &models.OBL2Update{
				Levels:    levels,
				Timestamp: utils.MilliToTimestamp(ts),
				Trade:     false,
			}

			instr.lastUpdateTime = uint64(msg.Time.UnixNano() / 1000000)

			if state.instrumentData.orderBook.Crossed() {
				state.logger.Info("crossed order book")
				// Stop the socket, we will restart instrument at the end
				if err := state.obWs.Disconnect(); err != nil {
					state.logger.Info("error disconnecting from socket", log.Error(err))
				}
				break
			}

			//fmt.Println(state.instrumentData.orderBook)

			// Send OBData
			context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
				UpdateL2: obDelta,
				SeqNum:   instr.seqNum + 1,
			})
			instr.seqNum += 1

		case []bitz.WSOrder:
			trades := msg.Message.([]bitz.WSOrder)

			if len(trades) == 0 {
				break
			}

			ts := uint64(msg.Time.UnixNano() / 1000000)

			sort.Slice(trades, func(i, j int) bool {
				return trades[i].ID < trades[j].ID
			})

			var buyTrades []bitz.WSOrder
			var sellTrades []bitz.WSOrder

			for _, t := range trades {
				if t.Direction == "sell" {
					sellTrades = append(sellTrades, t)
				} else {
					buyTrades = append(buyTrades, t)
				}
			}

			if len(buyTrades) > 0 {
				if ts <= state.instrumentData.lastAggTradeTs {
					ts = state.instrumentData.lastAggTradeTs + 1
				}
				aggBuyTrade := &models.AggregatedTrade{
					Bid:         false,
					Timestamp:   utils.MilliToTimestamp(ts),
					AggregateID: buyTrades[0].ID,
					Trades:      nil,
				}
				for _, trade := range buyTrades {
					trd := models.Trade{
						Price:    trade.Price,
						Quantity: trade.Quantity,
						ID:       trade.ID,
					}
					aggBuyTrade.Trades = append(aggBuyTrade.Trades, trd)
				}
				context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
					Trades: []*models.AggregatedTrade{aggBuyTrade},
					SeqNum: state.instrumentData.seqNum + 1,
				})
				state.instrumentData.seqNum += 1
				state.instrumentData.lastAggTradeTs = ts
			}

			if len(sellTrades) > 0 {
				if ts <= state.instrumentData.lastAggTradeTs {
					ts = state.instrumentData.lastAggTradeTs + 1
				}
				aggSellTrade := &models.AggregatedTrade{
					Bid:         true,
					Timestamp:   utils.MilliToTimestamp(ts),
					AggregateID: sellTrades[0].ID,
					Trades:      nil,
				}
				for _, trade := range sellTrades {
					trd := models.Trade{
						Price:    trade.Price,
						Quantity: trade.Quantity,
						ID:       trade.ID,
					}
					aggSellTrade.Trades = append(aggSellTrade.Trades, trd)
				}
				context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
					Trades: []*models.AggregatedTrade{aggSellTrade},
					SeqNum: state.instrumentData.seqNum + 1,
				})
				state.instrumentData.seqNum += 1
				state.instrumentData.lastAggTradeTs = ts
			}
		}

		state.postHeartBeat(context)
		if err := state.checkSockets(context); err != nil {
			return fmt.Errorf("error checking sockets: %v", err)
		}
		context.Send(context.Self(), &readSocket{})
		return nil

	case <-time.After(1 * time.Second):
		state.postHeartBeat(context)
		if err := state.checkSockets(context); err != nil {
			return fmt.Errorf("error checking sockets: %v", err)
		}
		context.Send(context.Self(), &readSocket{})
		return nil
	}
}

func (state *Listener) checkSockets(context actor.Context) error {

	if time.Now().Sub(state.lastPingTime) > 5*time.Second {
		// "Ping" by resubscribing to the topic
		_ = state.obWs.Subscribe(state.security.Symbol, []string{bitz.WSDepthTopic})
		_ = state.tradeWs.Subscribe(state.security.Symbol, []string{bitz.WSOrderTopic})
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
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			SeqNum: state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1
		state.instrumentData.lastHBTime = time.Now()
	}
}
