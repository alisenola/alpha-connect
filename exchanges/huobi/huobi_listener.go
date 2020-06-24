package huobi

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/exchanges/huobi"
	"math"
	"reflect"
	"sort"
	"time"
)

type readSocket struct{}

type InstrumentData struct {
	orderBook      *gorderbook.OrderBookL2
	seqNum         uint64
	lastUpdateID   uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	lastAggTradeTs uint64
}

type Listener struct {
	obWs            *huobi.Websocket
	tradeWs         *huobi.Websocket
	obWsChan        chan *huobi.WebsocketMessage
	tradeWsChan     chan *huobi.WebsocketMessage
	security        *models.Security
	instrumentData  *InstrumentData
	executorManager *actor.PID
	mediator        *actor.PID
	logger          *log.Logger
	lastPingTime    time.Time
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
		obWsChan:        nil,
		tradeWsChan:     nil,
		security:        security,
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

	state.mediator = actor.NewLocalPID("data_broker")
	state.executorManager = actor.NewLocalPID("exchange_executor_manager")
	state.obWsChan = make(chan *huobi.WebsocketMessage, 10000)
	state.tradeWsChan = make(chan *huobi.WebsocketMessage, 10000)
	state.lastPingTime = time.Now()

	context.Send(context.Self(), &readSocket{})

	state.instrumentData = &InstrumentData{
		orderBook:      nil,
		seqNum:         0,
		lastUpdateID:   0,
		lastUpdateTime: 0,
		lastHBTime:     time.Now(),
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

	ws := huobi.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to huobi websocket: %v", err)
	}

	if err := ws.SubscribeMarketByPrice(state.security.Symbol, huobi.WSOBLevel150); err != nil {
		return fmt.Errorf("error subscribing to orderbook for %s", state.security.Symbol)
	}

	time.Sleep(1 * time.Second)

	if err := ws.RequestMarketByPrice(state.security.Symbol, huobi.WSOBLevel150); err != nil {
		return fmt.Errorf("error requesting orderbook snapshot for %s", state.security.Symbol)
	}

	var ob *gorderbook.OrderBookL2
	nTries := 0
	for nTries < 100 {
		if !ws.ReadMessage() {
			return fmt.Errorf("error reading message: %v", ws.Err)
		}

		switch ws.Msg.Message.(type) {
		case huobi.WSMarketByPriceResponse:
			res := ws.Msg.Message.(huobi.WSMarketByPriceResponse)
			bids, asks := res.ToBidAsk()
			tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement))
			lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot))
			ob = gorderbook.NewOrderBookL2(
				tickPrecision,
				lotPrecision,
				10000)
			ob.Sync(bids, asks)
			ts := uint64(ws.Msg.ClientTime.UnixNano()) / 1000000
			state.instrumentData.orderBook = ob
			state.instrumentData.lastUpdateTime = ts
			state.instrumentData.lastUpdateID = res.SeqNum
			state.instrumentData.seqNum = 0
			nTries = 100

		case huobi.WSMarketByPriceTick:
			// buffer the update
			state.obWsChan <- ws.Msg

		case huobi.WSError:
			err := fmt.Errorf("error getting orderbook: %s", ws.Msg.Message.(huobi.WSError).ErrMsg)
			return err
		}
		nTries += 1
	}

	if ob == nil {
		return fmt.Errorf("error getting orderbook")
	}

	state.obWs = ws

	go func(ws *huobi.Websocket) {
		for ws.ReadMessage() {
			state.obWsChan <- ws.Msg
		}
	}(ws)

	return nil
}

func (state *Listener) subscribeTrades(context actor.Context) error {
	if state.tradeWs != nil {
		_ = state.tradeWs.Disconnect()
	}

	ws := huobi.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to websocket: %v", err)
	}

	if err := ws.SubscribeMarketTradeDetail(state.security.Symbol); err != nil {
		return fmt.Errorf("error subscribing to trades for %s", state.security.Symbol)
	}

	state.tradeWs = ws

	go func(ws *huobi.Websocket) {
		for ws.ReadMessage() {
			state.tradeWsChan <- ws.Msg
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
	case msg := <-state.obWsChan:
		switch msg.Message.(type) {

		case error:
			return fmt.Errorf("OB socket error: %v", msg)

		case huobi.WSMarketByPriceTick:
			ts := uint64(msg.ClientTime.UnixNano() / 1000000)
			update := msg.Message.(huobi.WSMarketByPriceTick)

			instr := state.instrumentData
			if update.SeqNum <= instr.lastUpdateID {
				break
			}

			// This, implies that instr.lastUpdateID > seqNum
			// We want update.SeqNum > instr.lastUpdateID
			if instr.lastUpdateID < update.PrevSeqNum {
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

			for _, bid := range update.Bids {
				level := gorderbook.OrderBookLevel{
					Price:    bid.Price,
					Quantity: bid.Quantity,
					Bid:      true,
				}
				instr.orderBook.UpdateOrderBookLevel(level)
				obDelta.Levels = append(obDelta.Levels, level)
			}
			for _, ask := range update.Asks {
				level := gorderbook.OrderBookLevel{
					Price:    ask.Price,
					Quantity: ask.Quantity,
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
			instr.lastUpdateID = update.SeqNum

			//state.postSnapshot(context)

		case huobi.WSSubscribeResponse:
			// pass

		case huobi.WSPing:
			msg := msg.Message.(huobi.WSPing)
			if err := state.tradeWs.Pong(msg.Ping); err != nil {
				return fmt.Errorf("error sending pong to websocket")
			}

		case huobi.WSError:
			msg := msg.Message.(huobi.WSError)
			state.logger.Info("got WSError message",
				log.String("message", msg.ErrMsg),
				log.String("code", msg.ErrCode))

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

	case msg := <-state.tradeWsChan:
		switch msg.Message.(type) {

		case error:
			return fmt.Errorf("trade socket error: %v", msg)

		case huobi.WSMarketTradeDetailTick:
			ts := uint64(msg.ClientTime.UnixNano() / 1000000)
			trades := msg.Message.(huobi.WSMarketTradeDetailTick)
			if len(trades.Data) == 0 {
				break
			}

			sort.Slice(trades.Data, func(i, j int) bool {
				return trades.Data[i].TradeId < trades.Data[j].TradeId
			})

			var aggTrade *models.AggregatedTrade
			for _, trade := range trades.Data {
				aggID := trade.Ts * 10
				// do that so new agg trade if side changes
				if trade.Direction == "sell" {
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
						Bid:         trade.Direction == "sell",
						Timestamp:   utils.MilliToTimestamp(ts),
						AggregateID: aggID,
						Trades:      nil,
					}
				}

				trd := models.Trade{
					Price:    trade.Price,
					Quantity: trade.Amount,
					ID:       trade.TradeId,
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

		case huobi.WSSubscribeResponse:
			// pass

		case huobi.WSPing:
			msg := msg.Message.(huobi.WSPing)
			if err := state.tradeWs.Pong(msg.Ping); err != nil {
				return fmt.Errorf("error sending pong to websocket")
			}

		case huobi.WSError:
			msg := msg.Message.(huobi.WSError)
			state.logger.Info("got WSError message",
				log.String("message", msg.ErrMsg),
				log.String("code", msg.ErrCode))

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
		_ = state.obWs.SubscribeMarketByPrice(state.security.Symbol, huobi.WSOBLevel150)
		_ = state.tradeWs.SubscribeMarketTradeDetail(state.security.Symbol)
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
