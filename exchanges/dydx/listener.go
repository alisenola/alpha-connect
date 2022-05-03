package dydx

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/gorderbook"
	gmodels "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	"gitlab.com/alphaticks/xchanger"
	"gitlab.com/alphaticks/xchanger/exchanges/dydx"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"math"
	"reflect"
	"time"
)

type checkSockets struct{}

type InstrumentData struct {
	orderBook      *gorderbook.OrderBookL2
	seqNum         uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	levelOffset    map[uint64]uint64
	lastAggTradeTs uint64
	levelDeltas    []gmodels.OrderBookLevel
	nCrossed       int
}

// OBType: OBL3
// No ID for the deltas..

type Listener struct {
	ws              *dydx.Websocket
	security        *models.Security
	dialerPool      *xchangerUtils.DialerPool
	instrumentData  *InstrumentData
	executorManager *actor.PID
	logger          *log.Logger
	lastPingTime    time.Time
	socketTicker    *time.Ticker
	lastMessageID   int
}

func NewListenerProducer(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Producer {
	return func() actor.Actor {
		return NewListener(security, dialerPool)
	}
}

func NewListener(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Actor {
	return &Listener{
		ws:              nil,
		security:        security,
		dialerPool:      dialerPool,
		instrumentData:  nil,
		executorManager: nil,
		logger:          nil,
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
	state.executorManager = actor.NewPID(context.ActorSystem().Address(), "exchange_executor_manager")

	state.instrumentData = &InstrumentData{
		orderBook:      nil,
		seqNum:         uint64(time.Now().UnixNano()),
		lastUpdateTime: 0,
		lastHBTime:     time.Now(),
		levelOffset:    make(map[uint64]uint64),
		lastAggTradeTs: 0,
	}

	if err := state.subscribeInstrument(context); err != nil {
		return fmt.Errorf("error subscribing to order book: %v", err)
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

	ws := dydx.NewWebsocket()
	if err := ws.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to websocket: %v", err)
	}
	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}

	if err := ws.SubscribeOrderBook(state.security.Symbol, false); err != nil {
		return fmt.Errorf("error subscribing to OBL2 stream: %v", err)
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	resp, ok := ws.Msg.Message.(dydx.WSOrderBookSubscribed)
	if !ok {
		return fmt.Errorf("was expecting WSOrderBookSubscribed, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot.Value))

	var bids, asks []gmodels.OrderBookLevel
	for _, l := range resp.Contents.Bids {
		if uint64(math.Round(l.Size*float64(lotPrecision))) == 0 {
			continue
		}
		bids = append(bids, gmodels.OrderBookLevel{
			Price:    l.Price,
			Quantity: l.Size,
			Bid:      true,
		})
		k := uint64(math.Round(l.Price * float64(tickPrecision)))
		state.instrumentData.levelOffset[k] = resp.Contents.Offset
	}
	for _, l := range resp.Contents.Asks {
		if uint64(math.Round(l.Size*float64(lotPrecision))) == 0 {
			continue
		}
		asks = append(asks, gmodels.OrderBookLevel{
			Price:    l.Price,
			Quantity: l.Size,
			Bid:      false,
		})
		k := uint64(math.Round(l.Price * float64(tickPrecision)))
		state.instrumentData.levelOffset[k] = resp.Contents.Offset
	}
	ts := uint64(ws.Msg.ClientTime.UnixNano()) / 1000000
	// TODO depth

	ob := gorderbook.NewOrderBookL2(
		tickPrecision,
		lotPrecision,
		10000)

	ob.Sync(bids, asks)
	if ob.Crossed() {
		return fmt.Errorf("crossed orderbook")
	}
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())
	state.instrumentData.orderBook = ob
	state.instrumentData.lastUpdateTime = ts

	if err := ws.SubscribeTrades(state.security.Symbol); err != nil {
		return fmt.Errorf("error subscribing to trade stream: %v", err)
	}

	state.ws = ws
	state.lastMessageID = resp.MessageID

	go func(ws *dydx.Websocket, pid *actor.PID) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
		}
	}(ws, context.Self())

	return nil
}

func (state *Listener) OnMarketDataRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketDataRequest)
	if state.instrumentData.orderBook.Crossed() {
		response := &messages.MarketDataResponse{
			RequestID:       msg.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			RejectionReason: messages.Other,
			Success:         false,
		}
		context.Respond(response)
		return nil
	}
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
	if state.ws == nil || msg.WSID != state.ws.ID {
		return nil
	}
	switch res := msg.Message.(type) {

	case error:
		return fmt.Errorf("socket error: %v", msg)

	case dydx.WSError:
		return fmt.Errorf("socket error: %v", res.Message)

	case dydx.WSOrderBookData:
		if res.MessageID != state.lastMessageID+1 {
			return state.subscribeInstrument(context)
		}
		state.lastMessageID = res.MessageID
		var bids, asks []gmodels.OrderBookLevel

		for _, l := range res.Contents.Bids {
			k := uint64(math.Round(l.Price * float64(state.instrumentData.orderBook.TickPrecision)))
			if state.instrumentData.levelOffset[k] < res.Contents.Offset {
				bids = append(bids, gmodels.OrderBookLevel{
					Price:    l.Price,
					Quantity: l.Size,
					Bid:      true,
				})
				state.instrumentData.levelOffset[k] = res.Contents.Offset
			}
		}
		for _, l := range res.Contents.Asks {
			k := uint64(math.Round(l.Price * float64(state.instrumentData.orderBook.TickPrecision)))
			if state.instrumentData.levelOffset[k] < res.Contents.Offset {
				asks = append(asks, gmodels.OrderBookLevel{
					Price:    l.Price,
					Quantity: l.Size,
					Bid:      false,
				})
				state.instrumentData.levelOffset[k] = res.Contents.Offset
			}
		}

		for _, bid := range bids {
			state.instrumentData.orderBook.UpdateOrderBookLevel(bid)
			state.instrumentData.levelDeltas = append(state.instrumentData.levelDeltas, bid)
		}
		for _, ask := range asks {
			state.instrumentData.orderBook.UpdateOrderBookLevel(ask)
			state.instrumentData.levelDeltas = append(state.instrumentData.levelDeltas, ask)
		}

		if !state.instrumentData.orderBook.Crossed() {
			// Send the deltas
			ts := uint64(msg.ClientTime.UnixNano()) / 1000000
			state.postDelta(context, ts)
		} else {
			state.instrumentData.nCrossed += 1
		}
		if state.instrumentData.nCrossed > 20 {
			fmt.Println("CROSSED")
			fmt.Println(state.instrumentData.orderBook)
			return state.subscribeInstrument(context)
		}

	case dydx.WSTradesSubscribed:
		// Ignore
		if res.MessageID != state.lastMessageID+1 {
			return state.subscribeInstrument(context)
		}
		state.lastMessageID = res.MessageID

	case dydx.WSTradesData:
		if res.MessageID != state.lastMessageID+1 {
			return state.subscribeInstrument(context)
		}
		state.lastMessageID = res.MessageID

		var aggTrade *models.AggregatedTrade
		ts := uint64(msg.ClientTime.UnixNano()) / 1000000
		for _, trade := range res.Contents.Trades {
			aggID := (uint64(trade.CreatedAt.UnixNano()) / 1000) * 10
			// Add one to aggregatedID if it's a sell so that
			// buy and sell happening at the same time won't have the same ID
			if trade.Side == "SELL" {
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
					Bid:         trade.Side == "SELL",
					Timestamp:   utils.MilliToTimestamp(ts),
					AggregateID: aggID,
					Trades:      nil,
				}
			}
			trade := models.Trade{
				Price:    trade.Price,
				Quantity: trade.Size,
				ID:       aggID,
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

	case dydx.WSPong:

	default:
		return fmt.Errorf("received unknown message: %s", reflect.TypeOf(msg.Message).String())
	}

	return nil
}

func (state *Listener) checkSockets(context actor.Context) error {

	if time.Since(state.lastPingTime) > 10*time.Second {
		_ = state.ws.Ping()

		state.lastPingTime = time.Now()
	}

	if state.ws.Err != nil || !state.ws.Connected {
		if state.ws.Err != nil {
			state.logger.Info("error on socket", log.Error(state.ws.Err))
		}
		if err := state.subscribeInstrument(context); err != nil {
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

func (state *Listener) postDelta(context actor.Context, ts uint64) {
	// Send the deltas

	if len(state.instrumentData.levelDeltas) > 1 {
		// Aggregate
		bids := make(map[uint64]gmodels.OrderBookLevel)
		asks := make(map[uint64]gmodels.OrderBookLevel)
		for _, l := range state.instrumentData.levelDeltas {
			k := uint64(math.Round(l.Price * float64(state.instrumentData.orderBook.TickPrecision)))
			if l.Bid {
				bids[k] = l
			} else {
				asks[k] = l
			}
		}
		state.instrumentData.levelDeltas = nil
		for _, l := range bids {
			state.instrumentData.levelDeltas = append(state.instrumentData.levelDeltas, l)
		}
		for _, l := range asks {
			state.instrumentData.levelDeltas = append(state.instrumentData.levelDeltas, l)
		}
	}

	obDelta := &models.OBL2Update{
		Levels:    state.instrumentData.levelDeltas,
		Timestamp: utils.MilliToTimestamp(ts),
		Trade:     false,
	}
	context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
		UpdateL2: obDelta,
		SeqNum:   state.instrumentData.seqNum + 1,
	})
	state.instrumentData.seqNum += 1
	state.instrumentData.lastUpdateTime = ts
	state.instrumentData.levelDeltas = nil
	state.instrumentData.nCrossed = 0
}
