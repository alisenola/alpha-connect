package ftx

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/exchanges/ftx"
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
	obWs           *ftx.Websocket
	tradeWs        *ftx.Websocket
	wsChan         chan *ftx.WebsocketMessage
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
			state.logger.Error("error processing GetOrderBookL2Request", log.Error(err))
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

	state.wsChan = make(chan *ftx.WebsocketMessage, 10000)
	state.lastPingTime = time.Now()

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
	context.Send(context.Self(), &readSocket{})

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

	obWs := ftx.NewWebsocket()
	err := obWs.Connect()
	if err != nil {
		return err
	}

	if err := obWs.Subscribe(state.security.Symbol, ftx.WSOrderBookChannel); err != nil {
		return fmt.Errorf("error subscribing to OBL2 stream: %v", err)
	}

	if !obWs.ReadMessage() {
		return fmt.Errorf("error reading message: %v", obWs.Err)
	}
	_, ok := obWs.Msg.Message.(ftx.WSSubscribeResponse)
	if !ok {
		if err, ok := obWs.Msg.Message.(ftx.WSError); ok {
			return fmt.Errorf("got WSError trying to subscribe to ob: %d: %s", err.Msg, err.Code)
		} else {
			return fmt.Errorf("was expecting WSSubscribeResponse, got %s", reflect.TypeOf(obWs.Msg.Message).String())
		}
	}

	if !obWs.ReadMessage() {
		return fmt.Errorf("error reading message: %v", obWs.Err)
	}
	obUpdate, ok := obWs.Msg.Message.(ftx.WSOrderBookUpdate)
	if !ok {
		return fmt.Errorf("was expecting WSOrderBookUpdate, got %s", reflect.TypeOf(obWs.Msg.Message).String())
	}

	var bids, asks []gorderbook.OrderBookLevel
	bids = make([]gorderbook.OrderBookLevel, len(obUpdate.Snapshot.Bids), len(obUpdate.Snapshot.Bids))
	for i, bid := range obUpdate.Snapshot.Bids {
		bids[i] = gorderbook.OrderBookLevel{
			Price:    bid.Price,
			Quantity: bid.Quantity,
			Bid:      true,
		}
	}
	asks = make([]gorderbook.OrderBookLevel, len(obUpdate.Snapshot.Asks), len(obUpdate.Snapshot.Asks))
	for i, ask := range obUpdate.Snapshot.Asks {
		asks[i] = gorderbook.OrderBookLevel{
			Price:    ask.Price,
			Quantity: ask.Quantity,
			Bid:      false,
		}
	}
	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot))
	bestAsk := obUpdate.Snapshot.Asks[0].Price
	depth := int(((bestAsk * 1.1) - bestAsk) * float64(tickPrecision))

	if depth > 10000 {
		depth = 10000
	}

	ts := uint64(obWs.Msg.Time.UnixNano()) / 1000000

	ob := gorderbook.NewOrderBookL2(
		tickPrecision,
		lotPrecision,
		depth,
	)

	ob.Sync(bids, asks)
	state.instrumentData.orderBook = ob
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())
	state.instrumentData.lastUpdateTime = ts

	state.obWs = obWs

	go func(ws *ftx.Websocket) {
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

	ws := ftx.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to websocket: %v", err)
	}

	if err := ws.Subscribe(state.security.Symbol, ftx.WSTradeChannel); err != nil {
		return fmt.Errorf("error subscribing to trade stream: %v", err)
	}

	state.tradeWs = ws

	go func(ws *ftx.Websocket) {
		for ws.ReadMessage() {
			state.wsChan <- ws.Msg
		}
	}(ws)

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

		case ftx.WSOrderBookUpdate:
			obData := msg.Message.(ftx.WSOrderBookUpdate)
			instr := state.instrumentData
			nLevels := len(obData.Snapshot.Bids) + len(obData.Snapshot.Asks)

			ts := uint64(msg.Time.UnixNano()) / 1000000
			obDelta := &models.OBL2Update{
				Levels:    make([]gorderbook.OrderBookLevel, nLevels, nLevels),
				Timestamp: utils.MilliToTimestamp(ts),
				Trade:     false,
			}

			lvlIdx := 0
			for _, bid := range obData.Snapshot.Bids {
				level := gorderbook.OrderBookLevel{
					Price:    bid.Price,
					Quantity: bid.Quantity,
					Bid:      true,
				}
				obDelta.Levels[lvlIdx] = level
				lvlIdx += 1
				instr.orderBook.UpdateOrderBookLevel(level)
			}
			for _, ask := range obData.Snapshot.Asks {
				level := gorderbook.OrderBookLevel{
					Price:    ask.Price,
					Quantity: ask.Quantity,
					Bid:      false,
				}
				obDelta.Levels[lvlIdx] = level
				lvlIdx += 1
				instr.orderBook.UpdateOrderBookLevel(level)
			}

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
				SeqNum:   state.instrumentData.seqNum + 1,
			})
			state.instrumentData.seqNum += 1

		case ftx.WSTradeUpdate:
			tradeData := msg.Message.(ftx.WSTradeUpdate)
			ts := uint64(msg.Time.UnixNano() / 1000000)

			sort.Slice(tradeData.Trades, func(i, j int) bool {
				return tradeData.Trades[i].Time.Before(tradeData.Trades[j].Time)
			})

			var aggTrade *models.AggregatedTrade
			for _, trade := range tradeData.Trades {
				aggID := uint64(trade.Time.UnixNano())
				if trade.Side == "buy" {
					aggID += 1
				}
				if aggTrade == nil || aggTrade.AggregateID != aggID {
					if aggTrade != nil {
						// Send aggregate trade
						context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
							Trades: []*models.AggregatedTrade{aggTrade},
							SeqNum: state.instrumentData.seqNum + 1,
						})
						state.instrumentData.seqNum += 1
						state.instrumentData.lastAggTradeTs = ts
					}

					if ts <= state.instrumentData.lastAggTradeTs {
						ts = state.instrumentData.lastAggTradeTs + 1
					}
					aggTrade = &models.AggregatedTrade{
						Bid:         trade.Side == "sell",
						Timestamp:   utils.MilliToTimestamp(ts),
						AggregateID: aggID,
						Trades:      nil,
					}
				}
				trade := models.Trade{
					Price:    trade.Price,
					Quantity: trade.Size,
					ID:       trade.ID,
				}
				aggTrade.Trades = append(aggTrade.Trades, trade)
			}
			if aggTrade != nil {
				// Send aggregate trade
				context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
					Trades: []*models.AggregatedTrade{aggTrade},
					SeqNum: state.instrumentData.seqNum + 1,
				})
				state.instrumentData.seqNum += 1
				state.instrumentData.lastAggTradeTs = ts
			}
			state.instrumentData.lastAggTradeTs = ts
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
		_ = state.tradeWs.Ping()
		_ = state.obWs.Ping()
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
