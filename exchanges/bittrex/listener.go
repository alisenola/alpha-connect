package bittrex

/*
import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	exchangeModels "gitlab.com/alphaticks/alpha-connect/models/messages/exchanges"
	"gitlab.com/alphaticks/alpha-connect/models/messages/executor"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/bittrex"
	"reflect"
	"time"
)

type readSocket struct{}

type InstrumentData struct {
	orderBook      *gorderbook.OrderBookL2
	lastUpdateID   uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	lastAggTradeTs uint64
}

type Listener struct {
	ws              *bittrex.Websocket
	wsChan          chan *bittrex.WebsocketMessage
	instrument      *exchanges.Instrument
	instrumentData  *InstrumentData
	executorManager *actor.PID
	mediator        *actor.PID
	logger          *log.Logger
	lastPingTime    time.Time
}

func NewListenerProducer(instr *exchanges.Instrument) actor.Producer {
	return func() actor.Actor {
		return NewListener(instr)
	}
}

func NewListener(instr *exchanges.Instrument) actor.Actor {
	return &Listener{
		ws:              nil,
		wsChan:          nil,
		instrument:      instr,
		instrumentData:  nil,
		executorManager: nil,
		mediator:        nil,
		logger:          nil,
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

	case *executor.GetOrderBookL2Request:
		if err := state.GetOrderBookL2Request(context); err != nil {
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
		log.String("instrument", state.instrument.DefaultFormat()))

	state.mediator = actor.NewLocalPID("data_broker")
	state.executorManager = actor.NewLocalPID("exchange_executor_manager")
	state.wsChan = make(chan *bittrex.WebsocketMessage, 10000)
	state.lastPingTime = time.Now()

	state.instrumentData = &InstrumentData{
		orderBook:      nil,
		lastUpdateID:   0,
		lastUpdateTime: 0,
		lastHBTime:     time.Now(),
		lastAggTradeTs: 0,
	}

	if err := state.subscribeInstrument(context); err != nil {
		return fmt.Errorf("error subscribing to instrument: %v", err)
	}
	context.Send(context.Self(), &readSocket{})

	return nil
}

func (state *Listener) Clean(context actor.Context) error {
	if state.ws != nil {
		if err := state.ws.Disconnect(); err != nil {
			state.logger.Info("error disconnecting socket", log.Error(err))
		}
	}

	return nil
}

func (state *Listener) subscribeInstrument(context actor.Context) error {
	if state.ws != nil {
		_ = state.ws.Disconnect()
	}

	symbol := state.instrument.Format(bittrex.SymbolFormat)

	ws := bittrex.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to bittrex websocket: %v", err)
	}

	if err := ws.SubscribeToExchangeDeltas(symbol); err != nil {
		return fmt.Errorf("error subscribing to exchange deltas for symbol %s", symbol)
	}

	if err := ws.QueryExchangeState(symbol); err != nil {
		return fmt.Errorf("error querying exchange state for symbol %s", symbol)
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	exchangeState, ok := ws.Msg.Message.(bittrex.WSExchangeState)

	if !ok {
		if !ws.ReadMessage() {
			return fmt.Errorf("error reading message: %v", ws.Err)
		}
		exchangeState, ok = ws.Msg.Message.(bittrex.WSExchangeState)
		if !ok {
			if !ws.ReadMessage() {
				return fmt.Errorf("error reading message: %v", ws.Err)
			}
			exchangeState, ok = ws.Msg.Message.(bittrex.WSExchangeState)
			if !ok {
				return fmt.Errorf("was expecting depth WSExchangeState, got %s", reflect.TypeOf(ws.Msg.Message).String())
			}
		}
	}

	bids, asks := exchangeState.ToBidAsk()
	bestAsk := asks[0].Price
	// Allow a 10% price variation
	depth := int(((bestAsk * 1.1) - bestAsk) * float64(state.instrument.TickPrecision))
	if depth > 10000 {
		depth = 10000
	}
	ob := gorderbook.NewOrderBookL2(
		state.instrument.TickPrecision,
		state.instrument.LotPrecision,
		1)
	ob.Sync(bids, asks)

	ts := ws.Msg.Time.UnixNano() / 1000000
	state.instrumentData.orderBook = ob
	state.instrumentData.lastUpdateTime = uint64(ts)
	state.instrumentData.lastUpdateID = exchangeState.Nonce
	state.logger.Info("subscribed")
	state.ws = ws

	go func(ws *bittrex.Websocket) {
		for ws.ReadMessage() {
			state.wsChan <- ws.Msg
		}
	}(state.ws)

	return nil
}

func (state *Listener) GetOrderBookL2Request(context actor.Context) error {
	msg := context.Message().(*executor.GetOrderBookL2Request)
	if msg.Instrument.DefaultFormat() != state.instrument.DefaultFormat() {
		err := fmt.Errorf("order book not tracked by listener")
		context.Respond(&executor.GetOrderBookL2Response{
			RequestID: msg.RequestID,
			Error:     err,
			Snapshot:  nil})
		return nil
	}

	snapshot := &exchangeModels.OBL2Snapshot{
		Instrument: msg.Instrument,
		Bids:       state.instrumentData.orderBook.GetAbsoluteRawBids(0),
		Asks:       state.instrumentData.orderBook.GetAbsoluteRawAsks(0),
		Timestamp:  models.MilliToTimestamp(state.instrumentData.lastUpdateTime),
		ID:         state.instrumentData.lastUpdateID,
	}
	context.Respond(&executor.GetOrderBookL2Response{
		RequestID: msg.RequestID,
		Error:     nil,
		Snapshot:  snapshot})

	return nil
}

func (state *Listener) readSocket(context actor.Context) error {
	select {
	case msg := <-state.wsChan:
		fmt.Println("got ws mess")
		switch msg.Message.(type) {

		case error:
			return fmt.Errorf("socket error: %v", msg)

		case bittrex.WSExchangeStateUpdate:
			exchangeState := msg.Message.(bittrex.WSExchangeStateUpdate)
			ts := uint64(msg.Time.UnixNano() / 1000000)

			// Skip
			if state.instrumentData.lastUpdateID >= exchangeState.Nonce {
				break
			}

			if exchangeState.Nonce != state.instrumentData.lastUpdateID+1 {
				// TODO avoid logging for now
				state.logger.Info("error processing ob update for "+state.instrument.DefaultFormat(),
					log.Error(fmt.Errorf("out of order sequence %d : %d", state.instrumentData.lastUpdateID, exchangeState.Nonce)))
				// Stop the socket, we will restart instrument at the end
				fmt.Println("disconnecting")
				if err := state.ws.Disconnect(); err != nil {
					state.logger.Info("error disconnecting from socket", log.Error(err))
				}
				fmt.Println("disconnected")
				break
			}

			obDelta := &exchangeModels.OBDelta{
				Levels:    nil,
				Timestamp: models.MilliToTimestamp(ts),
				FirstID:   state.instrumentData.lastUpdateID + 1,
				ID:        exchangeState.Nonce,
				Trade:     false,
			}
			for _, bid := range exchangeState.Buys {
				level := gmodels.OrderBookLevel{
					Price:    bid.Rate,
					Quantity: bid.Quantity,
					Bid:      true,
				}
				state.instrumentData.orderBook.UpdateOrderBookLevel(level)
				obDelta.Levels = append(obDelta.Levels, level)
			}
			for _, ask := range exchangeState.Sells {
				level := gmodels.OrderBookLevel{
					Price:    ask.Rate,
					Quantity: ask.Quantity,
					Bid:      false,
				}
				state.instrumentData.orderBook.UpdateOrderBookLevel(level)
				obDelta.Levels = append(obDelta.Levels, level)
			}

			if state.instrumentData.orderBook.Crossed() {
				state.logger.Info("crossed order book")
				// Stop the socket, we will restart instrument at the end
				if err := state.ws.Disconnect(); err != nil {
					state.logger.Info("error disconnecting from socket", log.Error(err))
				}
				break
			}

			// TODO trade not ok
			deltaTopic := fmt.Sprintf("%s/OBDELTA", state.instrument.DefaultFormat())
			context.Send(state.mediator, &messages.PubSubMessage{
				ID:      uint64(time.Now().UnixNano()),
				Topic:   deltaTopic,
				Message: obDelta})

			state.instrumentData.lastUpdateTime = uint64(msg.Time.UnixNano() / 1000000)
			state.instrumentData.lastUpdateID = exchangeState.Nonce

			state.instrumentData.lastHBTime = time.Now()

			var aggTrade *exchangeModels.AggregatedTrade
			tradeTopic := fmt.Sprintf("%s/TRADE", state.instrument.DefaultFormat())

			for _, trade := range exchangeState.Fills {

				aggID := (uint64(time.Time(trade.TimeStamp).UnixNano()) / 1000000) * 10
				// Prevent buy and sell with same ID to be aggregated
				if trade.OrderType == "SELL" {
					aggID += 1
				}
				if aggTrade == nil || aggTrade.AggregateID != aggID {

					if aggTrade != nil {
						// Send aggregate trade
						context.Send(state.mediator, &messages.PubSubMessage{
							ID:      uint64(time.Now().UnixNano()),
							Topic:   tradeTopic,
							Message: aggTrade})
						state.instrumentData.lastAggTradeTs = ts
					}

					if ts <= state.instrumentData.lastAggTradeTs {
						ts = state.instrumentData.lastAggTradeTs + 1
					}
					aggTrade = &exchangeModels.AggregatedTrade{
						Bid:         trade.OrderType == "SELL",
						Timestamp:   models.MilliToTimestamp(ts),
						AggregateID: aggID,
						Trades:      nil,
					}
				}
			}

			if aggTrade != nil {
				// Send aggregate trade
				context.Send(state.mediator, &messages.PubSubMessage{
					ID:      uint64(time.Now().UnixNano()),
					Topic:   tradeTopic,
					Message: aggTrade})
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
		fmt.Println("timeout")
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
		fmt.Println("PING")
		// "Ping" by resubscribing to the topic
		symbol := state.instrument.Format(bittrex.SymbolFormat)
		err := state.ws.SubscribeToExchangeDeltas(symbol)
		fmt.Println(err)

		state.lastPingTime = time.Now()
	}

	if state.ws.Err != nil || !state.ws.Connected {
		if state.ws.Err != nil {
			state.logger.Info("error on socket", log.Error(state.ws.Err))
		}

		// Try 10 times because bittrex sucks
		tries := 10
		var err error = nil
		for i := 0; i < tries; i++ {
			if err = state.subscribeInstrument(context); err == nil {
				break
			} else {
				state.logger.Info(fmt.Sprintf("error on socket: try %d out of %d", i, tries), log.Error(err))
			}
		}
		if err != nil {
			return fmt.Errorf("error subscribing to instrument: %v", err)
		}
	}

	return nil
}

func (state *Listener) postHeartBeat(context actor.Context) {
	// If haven't sent anything for 2 seconds, send heartbeat
	if time.Now().Sub(state.instrumentData.lastHBTime) > 2*time.Second {
		topic := fmt.Sprintf("%s/HEARTBEAT", state.instrument.DefaultFormat())
		context.Send(state.mediator, &messages.PubSubMessage{
			ID:      uint64(time.Now().UnixNano()),
			Topic:   topic,
			Message: models.MilliToTimestamp(uint64(time.Now().UnixNano() / 1000000)),
		})
		state.instrumentData.lastHBTime = time.Now()
	}
}

*/
