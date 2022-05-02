package ftxus

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
	gmodels "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	"gitlab.com/alphaticks/xchanger"
	"gitlab.com/alphaticks/xchanger/exchanges/ftxus"
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
	ws             *ftxus.Websocket
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
		ws:             nil,
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
			state.logger.Error("error processing GetOrderBookL2Request", log.Error(err))
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

	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))

	ws := ftxus.NewWebsocket()
	err := ws.Connect(state.dialerPool.GetDialer())
	if err != nil {
		return err
	}

	if err := ws.SubscribeGroupedOrderBook(state.security.Symbol, state.security.MinPriceIncrement.Value); err != nil {
		return fmt.Errorf("error subscribing to OBL2 stream: %v", err)
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	_, ok := ws.Msg.Message.(ftxus.WSSubscribeResponse)
	if !ok {
		if err, ok := ws.Msg.Message.(ftxus.WSError); ok {
			return fmt.Errorf("got WSError trying to subscribe to ob: %s: %d", err.Msg, err.Code)
		} else {
			return fmt.Errorf("was expecting WSSubscribeResponse, got %s", reflect.TypeOf(ws.Msg.Message).String())
		}
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	obUpdate, ok := ws.Msg.Message.(ftxus.WSOrderBookGroupedUpdate)
	if !ok {
		return fmt.Errorf("was expecting WSOrderBookGroupedUpdate, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	var bids, asks []gmodels.OrderBookLevel
	for _, bid := range obUpdate.Snapshot.Bids {
		if bid.Quantity == 0 {
			continue
		}
		bids = append(bids, gmodels.OrderBookLevel{
			Price:    bid.Price,
			Quantity: bid.Quantity,
			Bid:      true,
		})
	}
	for _, ask := range obUpdate.Snapshot.Asks {
		if ask.Quantity == 0 {
			continue
		}
		asks = append(asks, gmodels.OrderBookLevel{
			Price:    ask.Price,
			Quantity: ask.Quantity,
			Bid:      false,
		})
	}
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot.Value))
	bestAsk := obUpdate.Snapshot.Asks[0].Price
	depth := int(((bestAsk * 1.1) - bestAsk) * float64(tickPrecision))

	if depth > 10000 {
		depth = 10000
	}

	ts := uint64(ws.Msg.ClientTime.UnixNano()) / 1000000

	ob := gorderbook.NewOrderBookL2(
		tickPrecision,
		lotPrecision,
		depth,
	)

	ob.Sync(bids, asks)
	if ob.Crossed() {
		return fmt.Errorf("crossed order book")
	}
	state.instrumentData.orderBook = ob
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())
	state.instrumentData.lastUpdateTime = ts

	if err := ws.Subscribe(state.security.Symbol, ftxus.WSTradeChannel); err != nil {
		return fmt.Errorf("error subscribing to trade stream: %v", err)
	}

	state.ws = ws

	go func(ws *ftxus.Websocket, pid *actor.PID) {
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
	msg := context.Message().(*xchanger.WebsocketMessage)
	switch msg.Message.(type) {

	case error:
		return fmt.Errorf("socket error: %v", msg)

	case ftxus.WSOrderBookGroupedUpdate:
		obData := msg.Message.(ftxus.WSOrderBookGroupedUpdate)
		instr := state.instrumentData
		nLevels := len(obData.Snapshot.Bids) + len(obData.Snapshot.Asks)

		ts := uint64(msg.ClientTime.UnixNano()) / 1000000
		obDelta := &models.OBL2Update{
			Levels:    make([]gmodels.OrderBookLevel, nLevels),
			Timestamp: utils.MilliToTimestamp(ts),
			Trade:     false,
		}

		lvlIdx := 0
		for _, bid := range obData.Snapshot.Bids {
			level := gmodels.OrderBookLevel{
				Price:    bid.Price,
				Quantity: bid.Quantity,
				Bid:      true,
			}
			obDelta.Levels[lvlIdx] = level
			lvlIdx += 1
			instr.orderBook.UpdateOrderBookLevel(level)
		}
		for _, ask := range obData.Snapshot.Asks {
			level := gmodels.OrderBookLevel{
				Price:    ask.Price,
				Quantity: ask.Quantity,
				Bid:      false,
			}
			obDelta.Levels[lvlIdx] = level
			lvlIdx += 1
			instr.orderBook.UpdateOrderBookLevel(level)
		}

		instr.lastUpdateTime = ts

		if state.instrumentData.orderBook.Crossed() {
			state.logger.Info("crossed orderbook", log.Error(errors.New("crossed")))
			return state.subscribeInstrument(context)
		}
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			UpdateL2: obDelta,
			SeqNum:   state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1

	case ftxus.WSTradeUpdate:
		tradeData := msg.Message.(ftxus.WSTradeUpdate)
		ts := uint64(msg.ClientTime.UnixNano() / 1000000)

		sort.Slice(tradeData.Trades, func(i, j int) bool {
			return tradeData.Trades[i].Time.Before(tradeData.Trades[j].Time)
		})

		var aggTrade *models.AggregatedTrade
		for _, trade := range tradeData.Trades {
			aggID := uint64(trade.Time.UnixNano())
			if trade.Side == "buy" {
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
			trade := models.Trade{
				Price:    trade.Price,
				Quantity: trade.Size,
				ID:       trade.ID,
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
		state.instrumentData.lastAggTradeTs = ts
	}

	return nil
}

func (state *Listener) checkSockets(context actor.Context) error {

	if state.ws.Err != nil || !state.ws.Connected {
		if state.ws.Err != nil {
			state.logger.Info("error on socket", log.Error(state.ws.Err))
		}
		if err := state.subscribeInstrument(context); err != nil {
			return fmt.Errorf("error subscribing to instrument: %v", err)
		}
	}

	if time.Since(state.lastPingTime) > 10*time.Second {
		// "Ping" by resubscribing to the topic
		if err := state.ws.SubscribeGroupedOrderBook(state.security.Symbol, state.security.MinPriceIncrement.Value); err != nil {
			return fmt.Errorf("error subscribing to OBL2 stream: %v", err)
		}
		if err := state.ws.Subscribe(state.security.Symbol, ftxus.WSTradeChannel); err != nil {
			return fmt.Errorf("error subscribing to trade stream: %v", err)
		}
		state.lastPingTime = time.Now()
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
