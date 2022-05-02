package bitstamp

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
	gmodels "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	"gitlab.com/alphaticks/xchanger"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/bitstamp"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"math"
	"reflect"
	"time"
)

type InstrumentDataL3 struct {
	tickPrecision  uint64
	lotPrecision   uint64
	orderBook      *gorderbook.OrderBookL3
	seqNum         uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	aggTrade       *models.AggregatedTrade
	lastAggTradeTs uint64
	levelDeltas    []gmodels.OrderBookLevel
}

// OBType: OBL3

type ListenerL3 struct {
	ws               *bitstamp.Websocket
	security         *models.Security
	dialerPool       *xchangerUtils.DialerPool
	instrumentData   *InstrumentDataL3
	bitstampExecutor *actor.PID
	logger           *log.Logger
	lastPingTime     time.Time
	stashedTrades    *list.List
	socketTicker     *time.Ticker
}

func NewListenerL3Producer(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Producer {
	return func() actor.Actor {
		return NewListenerL3(security, dialerPool)
	}
}

func NewListenerL3(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Actor {
	return &ListenerL3{
		ws:               nil,
		security:         security,
		dialerPool:       dialerPool,
		instrumentData:   nil,
		bitstampExecutor: nil,
		logger:           nil,
		stashedTrades:    nil,
	}
}

func (state *ListenerL3) Receive(context actor.Context) {
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

	case *postAggTrade:
		state.postAggTrade(context)
	}
}

func (state *ListenerL3) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()),
		log.String("exchange", state.security.Exchange.Name),
		log.String("symbol", state.security.Symbol))

	state.lastPingTime = time.Now()
	state.stashedTrades = list.New()
	state.bitstampExecutor = actor.NewPID(context.ActorSystem().Address(), "executor/exchanges/"+constants.BITSTAMP.Name+"_executor")

	if state.security.MinPriceIncrement == nil || state.security.RoundLot == nil {
		return fmt.Errorf("security is missing MinPriceIncrement or RoundLot")
	}
	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot.Value))

	state.instrumentData = &InstrumentDataL3{
		tickPrecision:  tickPrecision,
		lotPrecision:   lotPrecision,
		orderBook:      nil,
		seqNum:         uint64(time.Now().UnixNano()),
		lastUpdateTime: 0,
		lastHBTime:     time.Now(),
		aggTrade:       nil,
		lastAggTradeTs: 0,
		levelDeltas:    nil,
	}

	if err := state.subscribeInstrument(context); err != nil {
		return fmt.Errorf("error subscribing to instrument: %v", err)
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

func (state *ListenerL3) Clean(context actor.Context) error {
	if state.ws != nil {
		if err := state.ws.Disconnect(); err != nil {
			state.logger.Warn("error disconnecting socket", log.Error(err))
		}
	}
	if state.socketTicker != nil {
		state.socketTicker.Stop()
		state.socketTicker = nil
	}

	return nil
}

func (state *ListenerL3) subscribeInstrument(context actor.Context) error {
	if state.ws != nil {
		_ = state.ws.Disconnect()
	}

	ws := bitstamp.NewWebsocket()
	if err := ws.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to bitstamp websocket: %v", err)
	}

	if err := ws.Subscribe(state.security.Symbol, bitstamp.WSLiveOrdersChannel); err != nil {
		return fmt.Errorf("error subscribing to depth stream for symbol %s", state.security.Symbol)
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	_, ok := ws.Msg.Message.(bitstamp.WSSubscribedMessage)
	if !ok {
		return fmt.Errorf("was expecting WSSubscribed message, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	time.Sleep(5 * time.Second)
	fut := context.RequestFuture(
		state.bitstampExecutor,
		&messages.MarketDataRequest{
			RequestID: uint64(time.Now().UnixNano()),
			Subscribe: false,
			Instrument: &models.Instrument{
				SecurityID: &types.UInt64Value{Value: state.security.SecurityID},
				Exchange:   state.security.Exchange,
				Symbol:     &types.StringValue{Value: state.security.Symbol},
			},
			Aggregation: models.L3,
		},
		5*time.Second)

	res, err := fut.Result()
	if err != nil {
		return fmt.Errorf("error getting OBL3")
	}
	msg, ok := res.(*messages.MarketDataResponse)
	if !ok {
		return fmt.Errorf("was expecting MarketDataSnapshot, got %s", reflect.TypeOf(msg).String())
	}
	if !msg.Success {
		return fmt.Errorf("error fetching snapshot: %s", msg.RejectionReason.String())
	}
	if msg.SnapshotL3 == nil {
		return fmt.Errorf("market data snapshot has no OBL3")
	}

	/*
		for i := range msg.SnapshotL3.Bids {
			msg.SnapshotL3.Bids[i] = gorderbook.Order{
				Price:    msg.SnapshotL3.Bids[i].Price,
				Quantity: math.Round(msg.SnapshotL3.Bids[i].Quantity/state.security.RoundLot) * state.security.RoundLot,
				Bid:      msg.SnapshotL3.Bids[i].Bid,
				ID:       msg.SnapshotL3.Bids[i].ID,
			}
			fmt.Println(math.Round(msg.SnapshotL3.Bids[i].Quantity/state.security.RoundLot) - msg.SnapshotL3.Bids[i].Quantity / state.security.RoundLot)
		}
		for i := range msg.SnapshotL3.Asks {
			msg.SnapshotL3.Asks[i] = gorderbook.Order{
				Price:    msg.SnapshotL3.Asks[i].Price,
				Quantity: math.Round(msg.SnapshotL3.Asks[i].Quantity/state.security.RoundLot) * state.security.RoundLot,
				Bid:      msg.SnapshotL3.Asks[i].Bid,
				ID:       msg.SnapshotL3.Asks[i].ID,
			}
		}
	*/

	ob := gorderbook.NewOrderBookL3(
		state.instrumentData.tickPrecision,
		state.instrumentData.lotPrecision,
		10000)

	ob.Sync(msg.SnapshotL3.Bids, msg.SnapshotL3.Asks)
	if ob.Crossed() {
		return fmt.Errorf("crossed order book")
	}

	ts := uint64(ws.Msg.ClientTime.UnixNano()) / 1000000

	state.instrumentData.seqNum = uint64(time.Now().UnixNano())
	state.instrumentData.lastUpdateTime = ts
	state.instrumentData.orderBook = ob
	state.instrumentData.levelDeltas = nil

	if err := ws.Subscribe(state.security.Symbol, bitstamp.WSLiveTradesChannel); err != nil {
		return fmt.Errorf("error subscribing to trade stream for symbol %s", state.security.Symbol)
	}

	state.ws = ws

	go func(ws *bitstamp.Websocket, pid *actor.PID) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
		}
	}(ws, context.Self())

	return nil
}

func (state *ListenerL3) OnMarketDataRequest(context actor.Context) error {
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

func (state *ListenerL3) onWebsocketMessage(context actor.Context) error {
	msg := context.Message().(*xchanger.WebsocketMessage)
	switch msg.Message.(type) {

	case error:
		return fmt.Errorf("socket error: %v", msg)

	case bitstamp.WSCreatedOrder:
		o := msg.Message.(bitstamp.WSCreatedOrder)
		order := gmodels.Order{
			Price:    o.Price,
			Quantity: math.Round(o.Amount*float64(state.instrumentData.lotPrecision)) / float64(state.instrumentData.lotPrecision),
			Bid:      o.OrderType == 0,
			ID:       o.ID,
		}
		// If crossed before applying create, bug
		if state.instrumentData.orderBook.Crossed() {
			state.logger.Warn("crossed orderbook", log.Error(errors.New("crossed")))
			return state.subscribeInstrument(context)
		}
		if !state.instrumentData.orderBook.HasOrder(o.ID) {
			state.instrumentData.orderBook.AddOrder(order)
			var quantity float64
			if o.OrderType == 0 {
				quantity = state.instrumentData.orderBook.GetBid(order.Price)
			} else {
				quantity = state.instrumentData.orderBook.GetAsk(order.Price)
			}
			levelDelta := gmodels.OrderBookLevel{
				Price:    o.Price,
				Quantity: quantity,
				Bid:      o.OrderType == 0,
			}
			state.instrumentData.levelDeltas = append(state.instrumentData.levelDeltas, levelDelta)
		}
		if !state.instrumentData.orderBook.Crossed() {
			// Send the deltas
			ts := uint64(msg.ClientTime.UnixNano()) / 1000000
			state.postDelta(context, ts)
		}

	case bitstamp.WSChangedOrder:
		o := msg.Message.(bitstamp.WSChangedOrder)
		order := gmodels.Order{
			Price:    o.Price,
			Quantity: math.Round(o.Amount*float64(state.instrumentData.lotPrecision)) / float64(state.instrumentData.lotPrecision),
			Bid:      o.OrderType == 0,
			ID:       o.ID,
		}

		if state.instrumentData.orderBook.HasOrder(o.ID) {
			oldO := state.instrumentData.orderBook.GetOrder(o.ID)
			oldP := uint64(math.Round(oldO.Price * float64(state.instrumentData.tickPrecision)))
			newP := uint64(math.Round(o.Price * float64(state.instrumentData.tickPrecision)))
			state.instrumentData.orderBook.DeleteOrder(o.ID)
			if oldP != newP {
				var quantity float64
				if oldO.Bid {
					quantity = state.instrumentData.orderBook.GetBid(oldO.Price)
				} else {
					quantity = state.instrumentData.orderBook.GetAsk(oldO.Price)
				}

				state.instrumentData.levelDeltas = append(state.instrumentData.levelDeltas, gmodels.OrderBookLevel{
					Price:    oldO.Price,
					Quantity: quantity,
					Bid:      oldO.Bid,
				})
			}
		}
		state.instrumentData.orderBook.AddOrder(order)
		var quantity float64
		if order.Bid {
			quantity = state.instrumentData.orderBook.GetBid(o.Price)
		} else {
			quantity = state.instrumentData.orderBook.GetAsk(o.Price)
		}
		state.instrumentData.levelDeltas = append(state.instrumentData.levelDeltas, gmodels.OrderBookLevel{
			Price:    o.Price,
			Quantity: quantity,
			Bid:      order.Bid,
		})

		if !state.instrumentData.orderBook.Crossed() {
			// Send the deltas
			ts := uint64(msg.ClientTime.UnixNano()) / 1000000
			state.postDelta(context, ts)
		}

	case bitstamp.WSDeletedOrder:
		o := msg.Message.(bitstamp.WSDeletedOrder)
		if state.instrumentData.orderBook.HasOrder(o.ID) {
			oldO := state.instrumentData.orderBook.GetOrder(o.ID)
			state.instrumentData.orderBook.DeleteOrder(o.ID)
			var quantity float64
			if o.OrderType == 0 {
				quantity = state.instrumentData.orderBook.GetBid(oldO.Price)
			} else {
				quantity = state.instrumentData.orderBook.GetAsk(oldO.Price)
			}

			state.instrumentData.levelDeltas = append(state.instrumentData.levelDeltas, gmodels.OrderBookLevel{
				Price:    oldO.Price,
				Quantity: quantity,
				Bid:      o.OrderType == 0,
			})

			if !state.instrumentData.orderBook.Crossed() {
				ts := uint64(msg.ClientTime.UnixNano()) / 1000000
				state.postDelta(context, ts)
			}
		}

	case bitstamp.WSTrade:
		tradeData := msg.Message.(bitstamp.WSTrade)
		tradeData.MicroTimestamp = uint64(msg.ClientTime.UnixNano()) / 1000
		ts := tradeData.MicroTimestamp / 1000

		var aggID uint64
		if tradeData.Type == 1 {
			aggID = tradeData.SellOrderID
		} else {
			aggID = tradeData.BuyOrderID
		}

		if state.instrumentData.aggTrade == nil || state.instrumentData.aggTrade.AggregateID != aggID {
			if state.instrumentData.lastAggTradeTs >= ts {
				ts = state.instrumentData.lastAggTradeTs + 1
			}
			aggTrade := &models.AggregatedTrade{
				Bid:         tradeData.Type == 1,
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
			models.Trade{
				Price:    tradeData.Price,
				Quantity: tradeData.Amount,
				ID:       tradeData.ID,
			})
	}
	return nil
}

func (state *ListenerL3) checkSockets(context actor.Context) error {

	if time.Since(state.lastPingTime) > 10*time.Second {
		// "Ping" by resubscribing to the topic
		_ = state.ws.Subscribe(state.security.Symbol, bitstamp.WSLiveOrdersChannel)
		state.lastPingTime = time.Now()
	}

	if state.ws.Err != nil || !state.ws.Connected {
		if state.ws.Err != nil {
			state.logger.Warn("error on socket", log.Error(state.ws.Err))
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

func (state *ListenerL3) postAggTrade(context actor.Context) {
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

func (state *ListenerL3) postDelta(context actor.Context, ts uint64) {
	// Send the deltas

	if len(state.instrumentData.levelDeltas) > 1 {
		// Aggregate
		bids := make(map[uint64]gmodels.OrderBookLevel)
		asks := make(map[uint64]gmodels.OrderBookLevel)
		for _, l := range state.instrumentData.levelDeltas {
			k := uint64(math.Round(l.Price * float64(state.instrumentData.tickPrecision)))
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
}
