package bitz

import (
	"errors"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger"
	"gitlab.com/alphaticks/xchanger/exchanges/bitz"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"math"
	"reflect"
	"sort"
	"time"
)

type checkSockets struct{}

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
	security       *models.Security
	dialerPool     *xchangerUtils.DialerPool
	instrumentData *InstrumentData
	logger         *log.Logger
	lastPingTime   time.Time
	socketTicker   *time.Ticker
}

func NewListenerProducer(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Producer {
	return func() actor.Actor {
		return NewListener(security, dialerPool)
	}
}

func NewListener(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Actor {
	return &Listener{
		obWs:           nil,
		tradeWs:        nil,
		security:       security,
		dialerPool:     dialerPool,
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

	case *xchanger.WebsocketMessage:
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

	socketTicker := time.NewTicker(5 * time.Second)
	state.socketTicker = socketTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-socketTicker.C:
				context.Send(pid, &checkSockets{})
			case <-time.After(10 * time.Second):
				if state.socketTicker != socketTicker {
					// Only stop if socket ticker has changed
					return
				}
			}
		}
	}(context.Self())

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
	if state.socketTicker != nil {
		state.socketTicker.Stop()
		state.socketTicker = nil
	}

	return nil
}

func (state *Listener) subscribeOrderBook(context actor.Context) error {
	if state.obWs != nil {
		_ = state.obWs.Disconnect()
	}

	ws := bitz.NewWebsocket()
	if err := ws.Connect(state.dialerPool.GetDialer()); err != nil {
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
	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot.Value))
	ob := gorderbook.NewOrderBookL2(
		tickPrecision,
		lotPrecision,
		10000)
	ob.Sync(bids, asks)
	if ob.Crossed() {
		return fmt.Errorf("crossed orderbook")
	}
	ts := uint64(ws.Msg.ClientTime.UnixNano() / 1000000)
	state.instrumentData.orderBook = ob
	state.instrumentData.lastUpdateTime = ts
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())
	state.obWs = ws

	go func(ws *bitz.Websocket, pid *actor.PID) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
		}
	}(ws, context.Self())

	return nil
}

func (state *Listener) subscribeTrades(context actor.Context) error {
	if state.tradeWs != nil {
		_ = state.tradeWs.Disconnect()
	}

	ws := bitz.NewWebsocket()
	if err := ws.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to bitz websocket: %v", err)
	}

	if err := ws.Subscribe(state.security.Symbol, []string{bitz.WSOrderTopic}); err != nil {
		return fmt.Errorf("error subscribing to trade stream for symbol %s", state.security.Symbol)
	}

	state.tradeWs = ws

	go func(ws *bitz.Websocket, pid *actor.PID) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
		}
	}(ws, context.Self())

	return nil
}

func (state *Listener) OnMarketDataRequest(context actor.Context) error {
	fmt.Println("MD REQUEST")
	msg := context.Message().(*messages.MarketDataRequest)

	snapshot := &models.OBL2Snapshot{
		Bids:          state.instrumentData.orderBook.GetBids(0),
		Asks:          state.instrumentData.orderBook.GetAsks(0),
		Timestamp:     utils.MilliToTimestamp(state.instrumentData.lastUpdateTime),
		TickPrecision: &types.UInt64Value{Value: state.instrumentData.orderBook.TickPrecision},
		LotPrecision:  &types.UInt64Value{Value: state.instrumentData.orderBook.LotPrecision},
	}
	context.Respond(&messages.MarketDataResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		SnapshotL2: snapshot,
		SeqNum:     state.instrumentData.seqNum,
		Success:    true,
	})
	return nil
}

func (state *Listener) onWebsocketMessage(context actor.Context) error {
	msg := context.Message().(*xchanger.WebsocketMessage)
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

		ts := uint64(msg.ClientTime.UnixNano() / 1000000)

		obDelta := &models.OBL2Update{
			Levels:    levels,
			Timestamp: utils.MilliToTimestamp(ts),
			Trade:     false,
		}

		instr.lastUpdateTime = uint64(msg.ClientTime.UnixNano() / 1000000)

		if state.instrumentData.orderBook.Crossed() {
			state.logger.Info("crossed orderbook", log.Error(errors.New("crossed")))
			return state.subscribeOrderBook(context)
		}

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

		ts := uint64(msg.ClientTime.UnixNano() / 1000000)

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
	return nil
}

func (state *Listener) checkSockets(context actor.Context) error {
	if time.Since(state.lastPingTime) > 5*time.Second {
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

	// If haven't sent anything for 2 seconds, send heartbeat
	if time.Since(state.instrumentData.lastHBTime) > 2*time.Second {
		// Send an empty refresh
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			SeqNum: state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1
		state.instrumentData.lastHBTime = time.Now()
	}

	return nil
}
