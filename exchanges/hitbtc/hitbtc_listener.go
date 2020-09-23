package hitbtc

import (
	"errors"
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

type checkSockets struct{}

type InstrumentData struct {
	orderBook          *gorderbook.OrderBookL2
	seqNum             uint64
	lastUpdateSequence uint64
	lastUpdateTime     uint64
	lastHBTime         time.Time
	lastAggTradeTs     uint64
}

type Listener struct {
	ws                *hitbtc.Websocket
	security          *models.Security
	instrumentData    *InstrumentData
	logger            *log.Logger
	lastPingTime      time.Time
	lastSubscribeTime time.Time
	socketTicker      *time.Ticker
}

func NewListenerProducer(security *models.Security) actor.Producer {
	return func() actor.Actor {
		return NewListener(security)
	}
}

func NewListener(security *models.Security) actor.Actor {
	return &Listener{
		ws:             nil,
		security:       security,
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

	case *hitbtc.WebsocketMessage:
		if err := state.onWebsocketMessage(context); err != nil {
			state.logger.Error("error processing websocket message", log.Error(err))
			panic(err)
		}

	case *checkSockets:
		if err := state.checkSockets(context); err != nil {
			state.logger.Error("error checking socket", log.Error(err))
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

	if state.security.MinPriceIncrement == nil || state.security.RoundLot == nil {
		return fmt.Errorf("security is missing MinPriceIncrement or RoundLot")
	}
	state.lastPingTime = time.Now()
	state.lastSubscribeTime = time.Now()

	state.instrumentData = &InstrumentData{
		orderBook:          nil,
		seqNum:             uint64(time.Now().UnixNano()),
		lastUpdateSequence: 0,
		lastUpdateTime:     0,
		lastHBTime:         time.Now(),
		lastAggTradeTs:     0,
	}

	if err := state.subscribeInstrument(context); err != nil {
		return fmt.Errorf("error subscribing to order book: %v", err)
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

func (state *Listener) Clean(context actor.Context) error {
	if state.ws != nil {
		if err := state.ws.Disconnect(); err != nil {
			state.logger.Info("error disconnecting socket", log.Error(err))
		}
	}
	if state.socketTicker != nil {
		state.socketTicker.Stop()
		state.socketTicker = nil
	}

	return nil
}

func (state *Listener) subscribeInstrument(context actor.Context) error {
	if state.ws != nil {
		_ = state.ws.Disconnect()
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

	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot.Value))
	ob := gorderbook.NewOrderBookL2(
		tickPrecision,
		lotPrecision,
		10000)
	ob.Sync(bids, asks)

	ts := uint64(ws.Msg.Time.UnixNano()) / 1000000
	state.instrumentData.orderBook = ob
	state.instrumentData.lastUpdateTime = ts
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())
	state.instrumentData.lastUpdateSequence = snapshot.Sequence

	if err := ws.SubscribeTrades(state.security.Symbol); err != nil {
		return fmt.Errorf("error subscribing to trades")
	}

	state.ws = ws

	go func(ws *hitbtc.Websocket, pid *actor.PID) {
		for ws.ReadMessage() {
			actor.EmptyRootContext.Send(pid, ws.Msg)
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
			Bids:      state.instrumentData.orderBook.GetBids(0),
			Asks:      state.instrumentData.orderBook.GetAsks(0),
			Timestamp: utils.MilliToTimestamp(state.instrumentData.lastUpdateTime),
		}
		response.SnapshotL2 = snapshot
	}

	context.Respond(response)
	return nil
}

func (state *Listener) onWebsocketMessage(context actor.Context) error {
	msg := context.Message().(*hitbtc.WebsocketMessage)
	switch msg.Message.(type) {

	case error:
		return fmt.Errorf("socket error: %v", msg)

	case hitbtc.WSUpdateOrderBook:
		ts := uint64(msg.Time.UnixNano() / 1000000)
		update := msg.Message.(hitbtc.WSUpdateOrderBook)

		instr := state.instrumentData

		if update.Sequence != instr.lastUpdateSequence+1 {
			state.logger.Info("error processing ob update for "+update.Symbol, log.Error(fmt.Errorf("out of order sequence")))
			return state.subscribeInstrument(context)
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
			state.logger.Info("crossed orderbook", log.Error(errors.New("crossed")))
			return state.subscribeInstrument(context)
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
		var aggHelpR uint64 = 0
		for _, trade := range trades.Data {
			aggHelp := uint64(trade.Timestamp.UnixNano()/1000000) * 10
			// do that so new agg trade if side changes
			if trade.Side == "sell" {
				aggHelp += 1
			}
			if aggTrade == nil || aggHelpR != aggHelp {

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
					AggregateID: trade.ID,
					Trades:      nil,
				}
				aggHelpR = aggHelp
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

	return nil
}

func (state *Listener) checkSockets(context actor.Context) error {

	if time.Now().Sub(state.lastPingTime) > 10*time.Second {
		_ = state.ws.Ping()
		state.lastPingTime = time.Now()
	}

	/*
		if time.Now().Sub(state.lastSubscribeTime) > 1*time.Minute {
			_ = state.obWs.SubscribeOrderBook(state.security.Symbol)
			_ = state.tradeWs.SubscribeTrades(state.security.Symbol)
			state.lastSubscribeTime = time.Now()
		}
	*/

	if state.ws.Err != nil || !state.ws.Connected {
		if state.ws.Err != nil {
			state.logger.Info("error on ob socket", log.Error(state.ws.Err))
		}
		if err := state.subscribeInstrument(context); err != nil {
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
