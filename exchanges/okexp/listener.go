package okexp

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/okex"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
	"reflect"
	"sync"
	"time"
)

var liqd = 1 * time.Minute
var liqdLock = sync.RWMutex{}

type checkSockets struct{}
type postAggTrade struct{}
type updateLiquidations struct{}

type InstrumentData struct {
	orderBook         *gorderbook.OrderBookL2
	seqNum            uint64
	lastUpdateTime    uint64
	lastHBTime        time.Time
	aggTrade          *models.AggregatedTrade
	lastAggTradeTs    uint64
	lastLiquidationTs uint64
}

type Listener struct {
	obWs               *okex.Websocket
	tradeWs            *okex.Websocket
	security           *models.Security
	securityID         uint64
	dialerPool         *xchangerUtils.DialerPool
	instrumentData     *InstrumentData
	executor           *actor.PID
	logger             *log.Logger
	lastPingTime       time.Time
	stashedTrades      *list.List
	socketTicker       *time.Ticker
	liquidationsTicker *time.Ticker
}

func NewListenerProducer(securityID uint64, dialerPool *xchangerUtils.DialerPool) actor.Producer {
	return func() actor.Actor {
		return NewListener(securityID, dialerPool)
	}
}

func NewListener(securityID uint64, dialerPool *xchangerUtils.DialerPool) actor.Actor {
	return &Listener{
		securityID: securityID,
		dialerPool: dialerPool,
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

	case *updateLiquidations:
		if err := state.updateLiquidations(context); err != nil {
			state.logger.Error("error updating liquidations", log.Error(err))
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
		log.String("security-id", fmt.Sprintf("%d", state.securityID)))
	state.executor = actor.NewPID(context.ActorSystem().Address(), "executor/exchanges/"+constants.OKEXP.Name+"_executor")

	res, err := context.RequestFuture(state.executor, &messages.SecurityDefinitionRequest{
		RequestID:  0,
		Instrument: &models.Instrument{SecurityID: wrapperspb.UInt64(state.securityID)},
	}, 5*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error fetching security definition: %v", err)
	}
	def := res.(*messages.SecurityDefinitionResponse)
	if !def.Success {
		return fmt.Errorf("error fetching security definition: %s", def.RejectionReason.String())
	}
	state.security = def.Security
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()),
		log.String("security-id", fmt.Sprintf("%d", state.securityID)),
		log.String("exchange", state.security.Exchange.Name),
		log.String("symbol", state.security.Symbol))
	if state.security.MinPriceIncrement == nil || state.security.RoundLot == nil {
		return fmt.Errorf("security is missing MinPriceIncrement or RoundLot")
	}

	state.lastPingTime = time.Now()
	state.stashedTrades = list.New()
	state.instrumentData = &InstrumentData{
		orderBook:      nil,
		seqNum:         uint64(time.Now().UnixNano()),
		lastUpdateTime: 0,
		lastHBTime:     time.Now(),
		aggTrade:       nil,
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

	liquidationsTicker := time.NewTicker(10 * time.Second)
	state.liquidationsTicker = liquidationsTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-liquidationsTicker.C:
				context.Send(pid, &updateLiquidations{})
			case <-time.After(11 * time.Second):
				if state.liquidationsTicker != liquidationsTicker {
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

	ws := okex.NewWebsocket()
	if err := ws.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to okcoin websocket: %v", err)
	}

	if err := ws.Subscribe(state.security.Symbol, okex.WSBookL2Channel); err != nil {
		return fmt.Errorf("error subscribing to depth stream for symbol")
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	_, ok := ws.Msg.Message.(okex.WSSubscribe)
	if !ok {
		return fmt.Errorf("was expecting depth data, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	depthData, ok := ws.Msg.Message.(okex.WSDepthSnapshot)
	if !ok {
		return fmt.Errorf("was expecting depth data, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	bids, asks := depthData.ToBidAsk()
	//bestAsk := float64(asks[0].Price) / float64(state.instruments[obData.Symbol].instrument.TickPrecision)
	// Allow a 10% price variation
	//depth := int(((bestAsk * 1.1) - bestAsk) * float64(state.instruments[obData.Symbol].instrument.TickPrecision))

	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot.Value))
	ob := gorderbook.NewOrderBookL2(tickPrecision, lotPrecision, 10000)
	ob.Sync(bids, asks)
	if ob.Crossed() {
		return fmt.Errorf("cossed order book")
	}
	ts := uint64(ws.Msg.ClientTime.UnixNano() / 1000000)
	state.instrumentData.orderBook = ob
	state.instrumentData.lastUpdateTime = ts
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())

	state.obWs = ws

	go func(ws *okex.Websocket, pid *actor.PID) {
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

	ws := okex.NewWebsocket()
	if err := ws.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to okex websocket: %v", err)
	}

	if err := ws.Subscribe(state.security.Symbol, okex.WSTradesChannel); err != nil {
		return fmt.Errorf("error subscribing to trade stream")
	}
	if state.security.SecurityType == enum.SecurityType_CRYPTO_PERP {
		if err := ws.Subscribe(state.security.Symbol, okex.WSFundingRateChannel); err != nil {
			return fmt.Errorf("error subscribing to funding rate stream")
		}
		if err := ws.Subscribe(state.security.Symbol, okex.WSOpenInterestChannel); err != nil {
			return fmt.Errorf("error subscribing to open interest stream")
		}
	}

	state.tradeWs = ws

	go func(ws *okex.Websocket, pid *actor.PID) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
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

	if msg.Aggregation == models.OrderBookAggregation_L2 {
		snapshot := &models.OBL2Snapshot{
			Bids:          state.instrumentData.orderBook.GetBids(0),
			Asks:          state.instrumentData.orderBook.GetAsks(0),
			Timestamp:     utils.MilliToTimestamp(state.instrumentData.lastUpdateTime),
			TickPrecision: &wrapperspb.UInt64Value{Value: state.instrumentData.orderBook.TickPrecision},
			LotPrecision:  &wrapperspb.UInt64Value{Value: state.instrumentData.orderBook.LotPrecision},
		}
		response.SnapshotL2 = snapshot
	}

	context.Respond(response)
	return nil
}

func (state *Listener) onWebsocketMessage(context actor.Context) error {
	msg := context.Message().(*xchanger.WebsocketMessage)

	switch msg.Message.(type) {

	case error:
		return fmt.Errorf("socket error: %v", msg)

	case okex.WSDepthUpdate:
		if state.obWs == nil || msg.WSID != state.obWs.ID {
			return nil
		}

		obData := msg.Message.(okex.WSDepthUpdate)

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
			SeqNum:   state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1

	case okex.WSFundingRateUpdate:
		fmt.Println("FUNDING RATE UPDATE")
		fundData := msg.Message.(okex.WSFundingRateUpdate)
		refresh := &messages.MarketDataIncrementalRefresh{
			SeqNum: state.instrumentData.seqNum + 1,
		}
		refresh.Stats = append(refresh.Stats, &models.Stat{
			Timestamp: utils.MilliToTimestamp(fundData.FundingTime),
			StatType:  models.StatType_FundingRate,
			Value:     fundData.FundingRate,
		})
		context.Send(context.Parent(), refresh)

		state.instrumentData.seqNum += 1

	case okex.WSOpenInterestUpdate:
		fundData := msg.Message.(okex.WSOpenInterestUpdate)
		refresh := &messages.MarketDataIncrementalRefresh{
			SeqNum: state.instrumentData.seqNum + 1,
		}
		refresh.Stats = append(refresh.Stats, &models.Stat{
			Timestamp: utils.MilliToTimestamp(fundData.Ts),
			StatType:  models.StatType_OpenInterest,
			Value:     fundData.OpenInterest * state.security.Multiplier.Value,
		})
		context.Send(context.Parent(), refresh)

		state.instrumentData.seqNum += 1

	case []okex.WSTrade:
		trades := msg.Message.([]okex.WSTrade)
		// order trade by timestamps
		// if trades are on same side, aggregate. We aggregate because otherwise
		// we will have two trades on the same side with the same timestamp as
		// we use same current timestamp for timestamp of all the trades in the array
		// the timestamp of the first trade is going to be the aggregateID and each
		// trade in the aggregate will have as ID aggregateID + index of trade in aggregate

		if len(trades) == 0 {
			break
		}

		ts := uint64(msg.ClientTime.UnixNano() / 1000000)
		for _, t := range trades {
			aggID := t.Ts * 10
			// do that so new agg trade if side changes
			if t.Side == "sell" {
				aggID += 1
			}

			if state.instrumentData.aggTrade == nil || state.instrumentData.aggTrade.AggregateID != aggID {
				if state.instrumentData.lastAggTradeTs >= ts {
					ts = state.instrumentData.lastAggTradeTs + 1
				}
				aggTrade := &models.AggregatedTrade{
					Bid:         t.Side == "sell",
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
					time.Sleep(20 * time.Millisecond)
					context.Send(pid, &postAggTrade{})
				}(context.Self())
			}

			state.instrumentData.aggTrade.Trades = append(
				state.instrumentData.aggTrade.Trades,
				&models.Trade{
					Price:    t.Price,
					Quantity: t.Size,
					ID:       uint64(t.TradeID),
				})
		}
	}

	return nil
}

func (state *Listener) checkSockets(context actor.Context) error {
	if time.Since(state.lastPingTime) > 5*time.Second {
		// "Ping" by resubscribing to the topic
		_ = state.obWs.Ping()
		_ = state.tradeWs.Ping()
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

func (state *Listener) postAggTrade(context actor.Context) {
	nowMilli := uint64(time.Now().UnixNano() / 1000000)

	for el := state.stashedTrades.Front(); el != nil; el = state.stashedTrades.Front() {
		trd := el.Value.(*models.AggregatedTrade)
		if trd != nil && nowMilli-utils.TimestampToMilli(trd.Timestamp) > 20 {
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

func (state *Listener) updateLiquidations(context actor.Context) error {
	fut := context.RequestFuture(
		state.executor,
		&messages.HistoricalLiquidationsRequest{
			RequestID: uint64(time.Now().UnixNano()),
			Instrument: &models.Instrument{
				SecurityID: &wrapperspb.UInt64Value{Value: state.security.SecurityID},
				Exchange:   state.security.Exchange,
				Symbol:     &wrapperspb.StringValue{Value: state.security.Symbol},
			},
			From: utils.MilliToTimestamp(state.instrumentData.lastLiquidationTs + 1),
		}, 2*time.Second)

	context.ReenterAfter(fut, func(res interface{}, err error) {
		if err != nil {
			if err == actor.ErrTimeout {
				liqdLock.Lock()
				liqd = time.Duration(float64(liqd) * 1.01)
				liqdLock.Unlock()
				if state.liquidationsTicker != nil {
					state.liquidationsTicker.Reset(liqd)
				}
			}
			state.logger.Info("error fetching liquidations", log.Error(err))
			return
		}
		msg := res.(*messages.HistoricalLiquidationsResponse)
		if !msg.Success {
			// We want to converge towards the right value,
			if msg.RejectionReason == messages.RejectionReason_RateLimitExceeded || msg.RejectionReason == messages.RejectionReason_HTTPError {
				liqdLock.Lock()
				liqd = time.Duration(float64(liqd) * 1.01)
				liqdLock.Unlock()
				if state.liquidationsTicker != nil {
					state.liquidationsTicker.Reset(liqd)
				}
			}
			state.logger.Info("error fetching liquidations", log.Error(errors.New(msg.RejectionReason.String())))
			return
		}

		for _, liq := range msg.Liquidations {
			context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
				Liquidation: liq,
				SeqNum:      state.instrumentData.seqNum + 1,
			})
			state.instrumentData.seqNum += 1
			state.instrumentData.lastLiquidationTs = utils.TimestampToMilli(liq.Timestamp)
		}

		// Reduce delay
		liqdLock.Lock()
		liqd = time.Duration(float64(liqd) * 0.99)
		if liqd < 15*time.Second {
			liqd = 15 * time.Second
		}
		liqdLock.Unlock()
		if state.liquidationsTicker != nil {
			state.liquidationsTicker.Reset(liqd)
		}
	})

	return nil
}
