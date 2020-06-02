package kraken

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/exchanges/kraken"
	"math"
	"reflect"
	"sort"
	"time"
)

type readSocket struct{}

type OBL2Request struct {
	requester *actor.PID
	requestID int64
}

type InstrumentData struct {
	orderBook      *gorderbook.OrderBookL2
	seqNum         uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	lastAggTradeTs uint64
}

// OBType: OBL2
// OBL2 Timestamps: Per server, no ID so timestamps are the ID no risk of disorder
// Trades: Impossible to infer from delta
// Status: Not ready, different ID occurs for different listeners due to different obdelta received

// Note: on the kraken website, it can display wrong colors for the trades, but the
// websocket trades are correct

type Listener struct {
	obWs           *kraken.Websocket
	tradeWs        *kraken.Websocket
	wsChan         chan *kraken.WebsocketMessage
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

	state.wsChan = make(chan *kraken.WebsocketMessage, 10000)
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

	ws := kraken.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to kraken websocket: %v", err)
	}

	if err := ws.SubscribeDepth([]string{state.security.Symbol}, kraken.WSDepth1000); err != nil {
		return fmt.Errorf("error subscribing to depth stream")
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	s, ok := ws.Msg.Message.(kraken.WSSystemStatus)
	if !ok {
		return fmt.Errorf("was expecting WSOrderBookL2, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}
	if s.Status != "online" {
		return fmt.Errorf("was expecting online status, got %s", s.Status)
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	su, ok := ws.Msg.Message.(kraken.WSSubscriptionStatus)
	if !ok {
		return fmt.Errorf("was expecting WSOrderBookL2, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	if su.Status != "subscribed" {
		return fmt.Errorf("was expecting subscribed status, got %s", su.Status)
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	obData, ok := ws.Msg.Message.(kraken.WSOrderBookL2)
	if !ok {
		return fmt.Errorf("was expecting WSOrderBookL2, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	bids, asks := obData.ToBidAsk()

	//bestAsk := float64(asks[0].Price) / float64(state.instruments[obData.Symbol].instrument.TickPrecision)
	// Allow a 10% price variation
	//depth := int(((bestAsk * 1.1) - bestAsk) * float64(state.instruments[obData.Symbol].instrument.TickPrecision))

	var maxID uint64 = 0
	for _, ask := range obData.Asks {
		if ask.Time > maxID {
			maxID = ask.Time
		}
	}
	for _, bid := range obData.Bids {
		if bid.Time > maxID {
			maxID = bid.Time
		}
	}
	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot))
	ob := gorderbook.NewOrderBookL2(tickPrecision, lotPrecision, 10000)

	ob.Sync(bids, asks)

	ts := uint64(ws.Msg.Time.UnixNano() / 1000000)
	state.instrumentData.orderBook = ob
	state.instrumentData.seqNum = 0
	state.instrumentData.lastUpdateTime = ts

	state.obWs = ws

	go func(ws *kraken.Websocket) {
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

	ws := kraken.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to kraken websocket: %v", err)
	}

	if err := ws.SubscribeTrade([]string{state.security.Symbol}); err != nil {
		return fmt.Errorf("error subscribing to trade stream")
	}

	state.tradeWs = ws

	go func(ws *kraken.Websocket) {
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
		SeqNum:     state.instrumentData.seqNum,
	}
	if msg.Aggregation == models.L2 {
		snapshot := &models.OBL2Snapshot{
			Bids:      state.instrumentData.orderBook.GetBids(0),
			Asks:      state.instrumentData.orderBook.GetAsks(0),
			Timestamp: utils.MilliToTimestamp(state.instrumentData.lastUpdateTime),
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

		case kraken.WSOrderBookL2Update:
			obData := msg.Message.(kraken.WSOrderBookL2Update)

			// The idea of exploding bid and ask won't do, multiple listeners
			// is buggy

			nLevels := len(obData.Bids) + len(obData.Asks)

			if nLevels == 0 {
				break
			}
			instr := state.instrumentData

			ts := uint64(msg.Time.UnixNano() / 1000000)
			obDelta := &models.OBL2Update{
				Levels:    make([]gorderbook.OrderBookLevel, nLevels, nLevels),
				Timestamp: utils.MilliToTimestamp(ts),
				Trade:     false,
			}

			lvlIdx := 0
			for _, bid := range obData.Bids {
				level := gorderbook.OrderBookLevel{
					Price:    bid.Price,
					Quantity: bid.Amount,
					Bid:      true,
				}
				obDelta.Levels[lvlIdx] = level
				lvlIdx += 1
				instr.orderBook.UpdateOrderBookLevel(level)
			}

			for _, ask := range obData.Asks {
				level := gorderbook.OrderBookLevel{
					Price:    ask.Price,
					Quantity: ask.Amount,
					Bid:      false,
				}
				obDelta.Levels[lvlIdx] = level
				lvlIdx += 1
				instr.orderBook.UpdateOrderBookLevel(level)
			}

			if state.instrumentData.orderBook.Crossed() {
				state.logger.Info("crossed order book")
				// Stop the socket, we will restart instrument at the end
				if err := state.obWs.Disconnect(); err != nil {
					state.logger.Info("error disconnecting from socket", log.Error(err))
				}
				break
			}

			instr.lastUpdateTime = ts

			// Send OBData
			context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
				UpdateL2: obDelta,
				SeqNum:   state.instrumentData.seqNum + 1,
			})
			state.instrumentData.seqNum += 1

		case kraken.WSTradeUpdate:
			tradeUpdate := msg.Message.(kraken.WSTradeUpdate)
			// order trade by timestamps
			// if trades are on same side, aggregate. We aggregate because otherwise
			// we will have two trades on the same side with the same timestamp as
			// we use same current timestamp for timestamp of all the trades in the array
			// the timestamp of the first trade is going to be the aggregateID and each
			// trade in the aggregate will have as ID aggregateID + index of trade in aggregate

			if len(tradeUpdate.Trades) == 0 {
				break
			}

			ts := uint64(msg.Time.UnixNano() / 1000000)

			sort.Slice(tradeUpdate.Trades, func(i, j int) bool {
				return tradeUpdate.Trades[i].Time < tradeUpdate.Trades[j].Time
			})
			tradeID := tradeUpdate.Trades[0].Time

			var buyTrades []kraken.WSTrade
			var sellTrades []kraken.WSTrade

			for _, t := range tradeUpdate.Trades {
				if t.Side == "s" {
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
					AggregateID: tradeID,
					Trades:      nil,
				}
				for _, trade := range buyTrades {
					trd := models.Trade{
						Price:    trade.Price,
						Quantity: trade.Volume,
						ID:       tradeID,
					}
					tradeID += 1
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
					AggregateID: tradeID,
					Trades:      nil,
				}
				for _, trade := range sellTrades {
					trd := models.Trade{
						Price:    trade.Price,
						Quantity: trade.Volume,
						ID:       tradeID,
					}
					tradeID += 1
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

	if time.Now().Sub(state.lastPingTime) > 5*time.Second {
		// "Ping" by resubscribing to the topic
		_ = state.obWs.Ping()
		_ = state.tradeWs.SubscribeTrade([]string{state.security.Symbol})
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
