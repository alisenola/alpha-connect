package hitbtc

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/exchanges/hitbtc"
	"math"
	"reflect"
	"sort"
	"time"
)

type readSocket struct{}

type InstrumentData struct {
	orderBook          *gorderbook.OrderBookL2
	seqNum             uint64
	lastUpdateSequence uint64
	lastUpdateTime     uint64
	lastHBTime         time.Time
	lastAggTradeTs     uint64
}

type Listener struct {
	obWs              *hitbtc.Websocket
	tradeWs           *hitbtc.Websocket
	wsChan            chan *hitbtc.WebsocketMessage
	security          *models.Security
	instrumentData    *InstrumentData
	logger            *log.Logger
	lastPingTime      time.Time
	lastSubscribeTime time.Time
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

	state.wsChan = make(chan *hitbtc.WebsocketMessage, 10000)
	state.lastPingTime = time.Now()
	state.lastSubscribeTime = time.Now()

	context.Send(context.Self(), &readSocket{})

	state.instrumentData = &InstrumentData{
		orderBook:          nil,
		seqNum:             0,
		lastUpdateSequence: 0,
		lastUpdateTime:     0,
		lastHBTime:         time.Now(),
		lastAggTradeTs:     0,
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

	ws := hitbtc.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to hitbtc websocket: %v", err)
	}

	if err := ws.SubscribeOrderBook(state.security.Symbol); err != nil {
		return fmt.Errorf("error subscribing to orderbook")
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}

	snapshot, ok := ws.Msg.Message.(hitbtc.WSSnapshotOrderBook)
	if !ok {
		return fmt.Errorf("was expecting WSSnapshotOrderBook, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	bids, asks := snapshot.ToBidAsk()

	//estAsk := float64(asks[0].Price) / float64(state.instruments[obData.Symbol].instrument.TickPrecision)
	// Allow a 10% price variation
	//depth := int(((bestAsk * 1.1) - bestAsk) * float64(state.instruments[obData.Symbol].instrument.TickPrecision))

	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot))
	ob := gorderbook.NewOrderBookL2(
		tickPrecision,
		lotPrecision,
		10000)
	ob.Sync(bids, asks)

	ts := uint64(ws.Msg.Time.UnixNano()) / 1000000
	state.instrumentData.orderBook = ob
	state.instrumentData.lastUpdateTime = ts
	state.instrumentData.seqNum = 0
	state.instrumentData.lastUpdateSequence = snapshot.Sequence

	state.obWs = ws

	go func(ws *hitbtc.Websocket) {
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

	ws := hitbtc.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to hitbtc websocket: %v", err)
	}

	if err := ws.SubscribeTrades(state.security.Symbol); err != nil {
		return fmt.Errorf("error subscribing to orderbook")
	}

	state.tradeWs = ws

	go func(ws *hitbtc.Websocket) {
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

		case hitbtc.WSUpdateOrderBook:
			ts := uint64(msg.Time.UnixNano() / 1000000)
			update := msg.Message.(hitbtc.WSUpdateOrderBook)

			instr := state.instrumentData

			if update.Sequence != instr.lastUpdateSequence+1 {
				state.logger.Info("error processing ob update for "+update.Symbol, log.Error(fmt.Errorf("out of order sequence")))
				// Stop the socket, we will restart instrument at the end
				if err := state.obWs.Disconnect(); err != nil {
					state.logger.Info("error disconnecting from socket", log.Error(err))
				}
				break
			}

			obDelta := &models.OBL2Update{
				Levels:    nil,
				Timestamp: utils.MilliToTimestamp(ts),
				Trade:     false,
			}

			for _, bid := range update.Bid {
				level := gorderbook.OrderBookLevel{
					Price:    bid.Price,
					Quantity: bid.Size,
					Bid:      true,
				}
				instr.orderBook.UpdateOrderBookLevel(level)
				obDelta.Levels = append(obDelta.Levels, level)
			}
			for _, ask := range update.Ask {
				level := gorderbook.OrderBookLevel{
					Price:    ask.Price,
					Quantity: ask.Size,
					Bid:      false,
				}
				instr.orderBook.UpdateOrderBookLevel(level)
				obDelta.Levels = append(obDelta.Levels, level)
			}

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
			instr.lastUpdateTime = ts
			instr.lastUpdateSequence = update.Sequence

			//state.postSnapshot(context)

		case hitbtc.WSUpdateTrades:
			ts := uint64(msg.Time.UnixNano() / 1000000)
			trades := msg.Message.(hitbtc.WSUpdateTrades)

			sort.Slice(trades.Data, func(i, j int) bool {
				return trades.Data[i].ID < trades.Data[j].ID
			})
			var aggTrade *models.AggregatedTrade
			for _, trade := range trades.Data {
				aggID := uint64(trade.Timestamp.UnixNano()/1000000) * 10
				// do that so new agg trade if side changes
				if trade.Side == "sell" {
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

				trd := models.Trade{
					Price:    trade.Price,
					Quantity: trade.Quantity,
					ID:       trade.ID,
				}

				aggTrade.Trades = append(aggTrade.Trades, trd)
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

		case hitbtc.WSSnapshotTrades:
			break

		case hitbtc.WSSnapshotOrderBook:
			break

		case hitbtc.WSPong:
			break

		case hitbtc.WSError:
			msg := msg.Message.(hitbtc.WSError)
			state.logger.Info("got WSError message",
				log.String("message", msg.Message),
				log.Uint64("code", msg.Code),
				log.String("description", msg.Description))

		default:
			state.logger.Info("received unknown message",
				log.String("message_type",
					reflect.TypeOf(msg.Message).String()))
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
	if time.Now().Sub(state.lastPingTime) > 10*time.Second {
		_ = state.obWs.Ping()
		_ = state.tradeWs.Ping()
		state.lastPingTime = time.Now()
	}

	if time.Now().Sub(state.lastSubscribeTime) > 1*time.Minute {
		_ = state.obWs.SubscribeOrderBook(state.security.Symbol)
		_ = state.tradeWs.SubscribeTrades(state.security.Symbol)
		state.lastSubscribeTime = time.Now()
	}

	if state.obWs.Err != nil || !state.obWs.Connected {
		if state.obWs.Err != nil {
			state.logger.Info("error on ob socket", log.Error(state.obWs.Err))
		}
		if err := state.subscribeOrderBook(context); err != nil {
			return fmt.Errorf("error subscribing to instrument: %v", err)
		}
	}

	if state.tradeWs.Err != nil || !state.tradeWs.Connected {
		if state.tradeWs.Err != nil {
			state.logger.Info("error on trade socket", log.Error(state.tradeWs.Err))
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
