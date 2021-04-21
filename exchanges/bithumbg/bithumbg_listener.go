package bithumbg

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
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/bithumbg"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"math"
	"reflect"
	"time"
)

const TRADE_DELAY_MS uint64 = 100

type checkSockets struct{}

type InstrumentData struct {
	orderBook      *gorderbook.OrderBookL2
	seqNum         uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	aggTrade       *models.AggregatedTrade
	lastAggTradeTs uint64
	lastSequence   uint64
}

type Listener struct {
	ws             *bithumbg.Websocket
	security       *models.Security
	dialerPool     *xchangerUtils.DialerPool
	instrumentData *InstrumentData
	executor       *actor.PID
	logger         *log.Logger
	lastPingTime   time.Time
	stashedTrades  *list.List
	socketTicker   *time.Ticker
}

func NewListenerProducer(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Producer {
	return func() actor.Actor {
		return NewListener(security, dialerPool)
	}
}

func NewListener(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Actor {
	return &Listener{
		ws:             nil,
		security:       security,
		dialerPool:     dialerPool,
		instrumentData: nil,
		executor:       nil,
		logger:         nil,
		stashedTrades:  nil,
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

	case *bithumbg.WebsocketMessage:
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
	state.executor = actor.NewPID(context.ActorSystem().Address(), "executor/"+constants.BITHUMBG.Name+"_executor")
	state.lastPingTime = time.Now()
	state.stashedTrades = list.New()

	state.instrumentData = &InstrumentData{
		orderBook:      nil,
		seqNum:         uint64(time.Now().UnixNano()),
		lastUpdateTime: 0,
		lastHBTime:     time.Now(),
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

	ws := bithumbg.NewWebsocket()
	if err := ws.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to bithumb websocket: %v", err)
	}

	if err := ws.SubscribeOrderBook(state.security.Symbol); err != nil {
		return fmt.Errorf("error subscribing to orderbook for %s", state.security.Symbol)
	}

	time.Sleep(5 * time.Second)
	fut := context.RequestFuture(
		state.executor,
		&messages.MarketDataRequest{
			RequestID: uint64(time.Now().UnixNano()),
			Subscribe: false,
			Instrument: &models.Instrument{
				SecurityID: &types.UInt64Value{Value: state.security.SecurityID},
				Exchange:   state.security.Exchange,
				Symbol:     &types.StringValue{Value: state.security.Symbol},
			},
			Aggregation: models.L2,
		},
		5*time.Second)

	res, err := fut.Result()
	if err != nil {
		return fmt.Errorf("error getting OBL2")
	}
	msg, ok := res.(*messages.MarketDataResponse)
	if !ok {
		return fmt.Errorf("was expecting MarketDataSnapshot, got %s", reflect.TypeOf(msg).String())
	}
	if !msg.Success {
		return fmt.Errorf("error fetching snapshot: %s", msg.RejectionReason.String())
	}
	if msg.SnapshotL2 == nil {
		return fmt.Errorf("market data snapshot has no OBL3")
	}

	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot.Value))
	ob := gorderbook.NewOrderBookL2(
		tickPrecision,
		lotPrecision,
		10000)

	ob.Sync(msg.SnapshotL2.Bids, msg.SnapshotL2.Asks)
	state.instrumentData.orderBook = ob
	state.instrumentData.lastUpdateTime = uint64(time.Now().UnixNano() / 1000000)
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())
	state.instrumentData.lastSequence = msg.SeqNum
	state.ws = ws

	// Fetch messages until sync
	synced := false
	for !synced {
		if !ws.ReadMessage() {
			return fmt.Errorf("error reading message: %v", ws.Err)
		}
		switch res := ws.Msg.Message.(type) {
		case error:
			return fmt.Errorf("socket error: %v", msg)

		case bithumbg.WSOrderBook:
			if uint64(res.Version) < state.instrumentData.lastSequence {
				continue
			} else if uint64(res.Version) == state.instrumentData.lastSequence {
				synced = true
			} else {
				return fmt.Errorf("out of order sequence %d:%d", res.Version, state.instrumentData.lastSequence)
			}
		}
	}

	if err := ws.SubscribeTrades(state.security.Symbol); err != nil {
		return fmt.Errorf("error subscribing to trades for %s", state.security.Symbol)
	}

	state.ws = ws

	go func(ws *bithumbg.Websocket, pid *actor.PID) {
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
	msg := context.Message().(*bithumbg.WebsocketMessage)
	switch res := msg.Message.(type) {

	case error:
		return fmt.Errorf("OB socket error: %v", msg)

	case bithumbg.WSOrderBook:
		ts := uint64(msg.Time.UnixNano() / 1000000)

		instr := state.instrumentData

		if uint64(res.Version) <= state.instrumentData.lastSequence {
			return nil
		}

		if uint64(res.Version) != state.instrumentData.lastSequence+1 {
			return fmt.Errorf("order book inconsistent sequence: %d, %d",
				res.Version,
				state.instrumentData.lastSequence+1)
		}

		obDelta := &models.OBL2Update{
			Levels:    nil,
			Timestamp: utils.MilliToTimestamp(ts),
			Trade:     false,
		}

		for _, bid := range res.Bids {
			level := gorderbook.OrderBookLevel{
				Price:    bid.Price,
				Quantity: bid.Quantity,
				Bid:      true,
			}
			instr.orderBook.UpdateOrderBookLevel(level)
			obDelta.Levels = append(obDelta.Levels, level)
		}
		for _, ask := range res.Asks {
			level := gorderbook.OrderBookLevel{
				Price:    ask.Price,
				Quantity: ask.Quantity,
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
		state.instrumentData.lastSequence += 1

		instr.lastUpdateTime = ts

	case []bithumbg.WSTrade:
		break

	case bithumbg.WSTrade:
		ts := uint64(msg.Time.UnixNano() / 1000000)
		aggTrade := &models.AggregatedTrade{
			Bid:         res.Side == "sell",
			Timestamp:   utils.MilliToTimestamp(ts),
			AggregateID: uint64(res.Version),
			Trades: []models.Trade{{
				Price:    res.Price,
				Quantity: res.Volume,
				ID:       uint64(res.Version),
			}},
		}

		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			Trades: []*models.AggregatedTrade{aggTrade},
			SeqNum: state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1

	case bithumbg.WSResponse:
		// pass

	default:
		state.logger.Info("received unknown message",
			log.String("message_type",
				reflect.TypeOf(msg.Message).String()))
	}

	return nil
}

func (state *Listener) checkSockets(context actor.Context) error {
	if time.Now().Sub(state.lastPingTime) > 10*time.Second {
		// "Ping" by resubscribing to the topic
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
