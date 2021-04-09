package upbit

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
	"gitlab.com/alphaticks/xchanger/exchanges/upbit"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"math"
	"reflect"
	"time"
)

const TRADE_DELAY_MS uint64 = 100

type checkSockets struct{}
type postAggTrade struct{}

type InstrumentData struct {
	orderBook      *gorderbook.OrderBookL2
	seqNum         uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	aggTrade       *models.AggregatedTrade
	lastAggTradeTs uint64
}

type Listener struct {
	obWs              *upbit.Websocket
	tradeWs           *upbit.Websocket
	security          *models.Security
	dialerPool        *xchangerUtils.DialerPool
	instrumentData    *InstrumentData
	executorManager   *actor.PID
	logger            *log.Logger
	lastOBPingTime    time.Time
	lastTradePingTime time.Time
	stashedTrades     *list.List
	socketTicker      *time.Ticker
}

func NewListenerProducer(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Producer {
	return func() actor.Actor {
		return NewListener(security, dialerPool)
	}
}

func NewListener(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Actor {
	return &Listener{
		obWs:            nil,
		tradeWs:         nil,
		security:        security,
		dialerPool:      dialerPool,
		instrumentData:  nil,
		executorManager: nil,
		logger:          nil,
		stashedTrades:   nil,
		socketTicker:    nil,
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

	case *upbit.WebsocketMessage:
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

func (state *Listener) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()),
		log.String("exchange", state.security.Exchange.Name),
		log.String("symbol", state.security.Symbol))

	state.executorManager = actor.NewPID(context.ActorSystem().Address(), "exchange_executor_manager")
	state.stashedTrades = list.New()

	state.instrumentData = &InstrumentData{
		orderBook:      nil,
		seqNum:         uint64(time.Now().UnixNano()),
		lastUpdateTime: 0,
		lastHBTime:     time.Now(),
	}

	if err := state.subscribeOrderbook(context); err != nil {
		return fmt.Errorf("error subscribing to order book: %v", err)
	}
	if err := state.subscribeTrade(context); err != nil {
		return fmt.Errorf("error subscribing to trade: %v", err)
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
	if state.obWs != nil {
		if err := state.obWs.Disconnect(); err != nil {
			state.logger.Info("error disconnecting socket", log.Error(err))
		}
	}
	if state.tradeWs != nil {
		if err := state.tradeWs.Disconnect(); err != nil {
			state.logger.Info("error disconnecting socket", log.Error(err))
		}
	}

	if state.socketTicker != nil {
		state.socketTicker.Stop()
		state.socketTicker = nil
	}

	return nil
}

func (state *Listener) subscribeOrderbook(context actor.Context) error {
	if state.obWs != nil {
		_ = state.obWs.Disconnect()
	}
	state.lastOBPingTime = time.Now()

	obWs := upbit.NewWebsocket()
	if err := obWs.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to upbit websocket: %v", err)
	}

	if err := obWs.SubscribeOrderBook([]string{state.security.Symbol}); err != nil {
		return fmt.Errorf("error subscribing to orderbook for %s", state.security.Symbol)
	}

	var ob *gorderbook.OrderBookL2
	nTries := 0
	for nTries < 100 {
		if !obWs.ReadMessage() {
			return fmt.Errorf("error reading message: %v", obWs.Err)
		}

		switch res := obWs.Msg.Message.(type) {
		case upbit.WSOrderBook:
			if res.StreamType != "SNAPSHOT" {
				return fmt.Errorf("was expecting snapshot as first event, got %s", res.StreamType)
			}
			bids, asks := res.ToBidAsk()
			minPriceIncrement := math.MaxFloat64
			minLotIncrement := math.MaxFloat64

			for _, b := range bids {
				prec := 10.
				for math.Abs((math.Round(b.Quantity*prec)/prec)-b.Quantity) > 0. {
					prec *= 10
				}
				if 1/prec < minLotIncrement {
					minLotIncrement = 1 / prec
				}
			}

			for _, a := range asks {
				prec := 10.
				for math.Abs((math.Round(a.Quantity*prec)/prec)-a.Quantity) > 0. {
					prec *= 10
				}
				if 1/prec < minLotIncrement {
					minLotIncrement = 1 / prec
				}
			}

			for i := 1; i < len(bids); i++ {
				diff := math.Abs(bids[i-1].Price - bids[i].Price)
				if diff < minPriceIncrement {
					minPriceIncrement = diff
				}
			}
			for i := 1; i < len(asks); i++ {
				diff := math.Abs(asks[i].Price - asks[i-1].Price)
				if diff < minPriceIncrement {
					minPriceIncrement = diff
				}
			}
			tickPrecision := uint64(math.Ceil(1. / minPriceIncrement))

			lotPrecision := uint64(math.Ceil(1. / minLotIncrement))
			ob = gorderbook.NewOrderBookL2(
				tickPrecision,
				lotPrecision,
				100)
			ob.Sync(bids, asks)
			ts := uint64(obWs.Msg.Time.UnixNano()) / 1000000
			state.instrumentData.orderBook = ob
			state.instrumentData.lastUpdateTime = ts
			state.instrumentData.seqNum = uint64(time.Now().UnixNano())
			nTries = 100

		case upbit.WSStatus:
			if res.Status != "" {

			}
		}
		nTries += 1
	}

	if ob == nil {
		return fmt.Errorf("error getting orderbook")
	}

	state.obWs = obWs

	go func(ws *upbit.Websocket, pid *actor.PID) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
		}
	}(obWs, context.Self())

	return nil
}

func (state *Listener) subscribeTrade(context actor.Context) error {
	if state.tradeWs != nil {
		_ = state.tradeWs.Disconnect()
	}
	state.lastTradePingTime = time.Now()

	tradeWs := upbit.NewWebsocket()
	if err := tradeWs.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to upbit websocket: %v", err)
	}
	if err := tradeWs.SubscribeTrades([]string{state.security.Symbol}); err != nil {
		return fmt.Errorf("error subscribing to trades for %s", state.security.Symbol)
	}

	state.tradeWs = tradeWs

	go func(ws *upbit.Websocket, pid *actor.PID) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
		}
	}(tradeWs, context.Self())

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
			Bids:          state.instrumentData.orderBook.GetBids(0),
			Asks:          state.instrumentData.orderBook.GetAsks(0),
			Timestamp:     utils.MilliToTimestamp(state.instrumentData.lastUpdateTime),
			TickPrecision: &types.UInt64Value{Value: state.instrumentData.orderBook.TickPrecision},
			LotPrecision:  &types.UInt64Value{Value: state.instrumentData.orderBook.LotPrecision},
		}
		response.SnapshotL2 = snapshot
	}

	context.Respond(response)
	return nil
}

func (state *Listener) onWebsocketMessage(context actor.Context) error {
	msg := context.Message().(*upbit.WebsocketMessage)
	switch res := msg.Message.(type) {

	case error:
		return fmt.Errorf("OB socket error: %v", msg)

	case upbit.WSOrderBook:
		ts := uint64(msg.Time.UnixNano() / 1000000)

		instr := state.instrumentData

		obDelta := &models.OBL2Update{
			Levels:    nil,
			Timestamp: utils.MilliToTimestamp(ts),
			Trade:     false,
		}

		bids, asks := res.ToBidAsk()

		ob := gorderbook.NewOrderBookL2(
			instr.orderBook.TickPrecision,
			instr.orderBook.LotPrecision,
			100)
		ob.Sync(bids, asks)

		obDelta.Levels = instr.orderBook.Diff(ob)
		instr.orderBook = ob

		if state.instrumentData.orderBook.Crossed() {
			state.logger.Info("crossed orderbook", log.Error(errors.New("crossed")))
			return state.subscribeOrderbook(context)
		}
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			UpdateL2: obDelta,
			SeqNum:   state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1

		instr.lastUpdateTime = ts

	case upbit.WSTrade:
		//ts := uint64(msg.Time.UnixNano() / 1000000)
		ts := uint64(msg.Time.UnixNano() / 1000000)

		aggID := res.Timestamp * 10
		if res.AskBid == "BID" {
			aggID += 1
		}
		if state.instrumentData.aggTrade == nil || state.instrumentData.aggTrade.AggregateID != aggID {
			// Create new agg trade
			if state.instrumentData.lastAggTradeTs >= ts {
				ts = state.instrumentData.lastAggTradeTs + 1
			}
			aggTrade := &models.AggregatedTrade{
				Bid:         res.AskBid == "BID",
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
				time.Sleep(time.Duration(TRADE_DELAY_MS) * time.Millisecond)
				context.Send(pid, &postAggTrade{})
			}(context.Self())
		}

		state.instrumentData.aggTrade.Trades = append(
			state.instrumentData.aggTrade.Trades,
			models.Trade{
				Price:    res.TradePrice,
				Quantity: res.TradeVolume,
				ID:       res.SequentialID,
			})

	case upbit.WSStatus:
		// pass

	default:
		state.logger.Info("received unknown message",
			log.String("message_type",
				reflect.TypeOf(msg.Message).String()))
	}

	return nil
}

func (state *Listener) checkSockets(context actor.Context) error {
	if time.Now().Sub(state.lastOBPingTime) > 10*time.Second {
		// "Ping" by resubscribing to the topic
		_ = state.obWs.Ping()
		state.lastOBPingTime = time.Now()
	}

	if time.Now().Sub(state.lastTradePingTime) > 10*time.Second {
		// "Ping" by resubscribing to the topic
		_ = state.tradeWs.Ping()
		state.lastTradePingTime = time.Now()
	}

	if state.obWs.Err != nil || !state.obWs.Connected {
		if state.obWs.Err != nil {
			state.logger.Info("error on socket", log.Error(state.obWs.Err))
		}
		if err := state.subscribeOrderbook(context); err != nil {
			return fmt.Errorf("error subscribing to instrument: %v", err)
		}
	}

	if state.tradeWs.Err != nil || !state.tradeWs.Connected {
		if state.tradeWs.Err != nil {
			state.logger.Info("error on socket", log.Error(state.tradeWs.Err))
		}
		if err := state.subscribeTrade(context); err != nil {
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

func (state *Listener) postAggTrade(context actor.Context) {
	nowMilli := uint64(time.Now().UnixNano() / 1000000)

	for el := state.stashedTrades.Front(); el != nil; el = state.stashedTrades.Front() {
		trd := el.Value.(*models.AggregatedTrade)
		if trd != nil && nowMilli-utils.TimestampToMilli(trd.Timestamp) > TRADE_DELAY_MS {
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
