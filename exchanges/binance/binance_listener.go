package binance

import (
	"container/list"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/binance"
	"math"
	"reflect"
	"strings"
	"time"
)

// OBType: OBL2
// OBL2 Timestamps: ordered & consistent with sequence ID
// Trades: Impossible to infer from deltas
// Status: ready
type readSocket struct{}
type postAggTrade struct{}

type InstrumentData struct {
	orderBook      *gorderbook.OrderBookL2
	seqNum         uint64
	lastUpdateID   uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	aggTrade       *models.AggregatedTrade
	lastAggTradeTs uint64
}

type Listener struct {
	obWs            *binance.Websocket
	tradeWs         *binance.Websocket
	wsChan          chan *binance.WebsocketMessage
	security        *models.Security
	instrumentData  *InstrumentData
	binanceExecutor *actor.PID
	logger          *log.Logger
	stashedTrades   *list.List
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
		wsChan:          nil,
		security:        security,
		instrumentData:  nil,
		binanceExecutor: nil,
		logger:          nil,
		stashedTrades:   nil,
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

	state.binanceExecutor = actor.NewLocalPID("executor/" + constants.BINANCE.Name + "_executor")
	state.wsChan = make(chan *binance.WebsocketMessage, 10000)
	state.stashedTrades = list.New()

	context.Send(context.Self(), &readSocket{})

	state.instrumentData = &InstrumentData{
		orderBook:      nil,
		seqNum:         uint64(time.Now().UnixNano()),
		lastUpdateID:   0,
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

	obWs := binance.NewWebsocket()
	symbol := strings.ToLower(state.security.Symbol)
	err := obWs.Connect(
		symbol,
		[]string{binance.WSDepthStream100ms})
	if err != nil {
		return err
	}

	state.obWs = obWs

	time.Sleep(5 * time.Second)
	fut := context.RequestFuture(
		state.binanceExecutor,
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
		20*time.Second)

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
		return fmt.Errorf("market data snapshot has no OBL2")
	}

	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot))
	bestAsk := msg.SnapshotL2.Asks[0].Price
	depth := int(((bestAsk * 1.1) - bestAsk) * float64(tickPrecision))

	if depth > 10000 {
		depth = 10000
	}

	ob := gorderbook.NewOrderBookL2(
		tickPrecision,
		lotPrecision,
		depth,
	)

	ob.Sync(msg.SnapshotL2.Bids, msg.SnapshotL2.Asks)
	state.instrumentData.lastUpdateID = msg.SeqNum
	state.instrumentData.lastUpdateTime = utils.TimestampToMilli(msg.SnapshotL2.Timestamp)

	synced := false
	for !synced {
		if !obWs.ReadMessage() {
			return fmt.Errorf("error reading message: %v", obWs.Err)
		}
		depthData, ok := obWs.Msg.Message.(binance.WSDepthData)
		if !ok {
			return fmt.Errorf("was expecting depth data, got %s", reflect.TypeOf(obWs.Msg.Message).String())
		}

		if depthData.FinalUpdateID <= state.instrumentData.lastUpdateID {
			continue
		}

		bids, asks, err := depthData.ToBidAsk()
		if err != nil {
			return fmt.Errorf("error converting depth data: %s ", err.Error())
		}
		for _, bid := range bids {
			ob.UpdateOrderBookLevel(bid)
		}
		for _, ask := range asks {
			ob.UpdateOrderBookLevel(ask)
		}

		state.instrumentData.lastUpdateID = depthData.FinalUpdateID
		state.instrumentData.lastUpdateTime = uint64(obWs.Msg.Time.UnixNano() / 1000000)

		synced = true
	}

	state.instrumentData.orderBook = ob
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())

	go func(ws *binance.Websocket) {
		for ws.ReadMessage() {
			state.wsChan <- ws.Msg
		}
	}(state.obWs)

	return nil
}

func (state *Listener) subscribeTrades(context actor.Context) error {
	if state.tradeWs != nil {
		_ = state.tradeWs.Disconnect()
	}
	tradeWs := binance.NewWebsocket()
	symbol := strings.ToLower(state.security.Symbol)
	err := tradeWs.Connect(
		symbol,
		[]string{binance.WSTradeStream})
	if err != nil {
		return err
	}
	state.tradeWs = tradeWs

	go func(ws *binance.Websocket) {
		for ws.ReadMessage() {
			state.wsChan <- ws.Msg
		}
	}(state.tradeWs)

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

		case binance.WSDepthData:
			depthData := msg.Message.(binance.WSDepthData)

			// change event time
			depthData.EventTime = uint64(msg.Time.UnixNano()) / 1000000
			err := state.onDepthData(context, depthData)
			if err != nil {
				state.logger.Info("error processing depth data for "+depthData.Symbol,
					log.Error(err))
				// Stop the socket, we will restart instrument at the end
				if err := state.obWs.Disconnect(); err != nil {
					state.logger.Info("error disconnecting from socket", log.Error(err))
				}
			}

		case binance.WSTradeData:
			tradeData := msg.Message.(binance.WSTradeData)
			var aggregateID uint64
			if tradeData.MarketSell {
				aggregateID = uint64(tradeData.SellerOrderID)
			} else {
				aggregateID = uint64(tradeData.BuyerOrderID)
			}

			ts := uint64(msg.Time.UnixNano() / 1000000)

			if state.instrumentData.aggTrade == nil || state.instrumentData.aggTrade.AggregateID != aggregateID {
				if state.instrumentData.lastAggTradeTs >= ts {
					ts = state.instrumentData.lastAggTradeTs + 1
				}
				aggTrade := &models.AggregatedTrade{
					Bid:         tradeData.MarketSell,
					Timestamp:   utils.MilliToTimestamp(ts),
					AggregateID: aggregateID,
					Trades:      nil,
				}
				state.instrumentData.aggTrade = aggTrade
				state.instrumentData.lastAggTradeTs = ts

				// Stash the aggTrade
				state.stashedTrades.PushBack(aggTrade)
				// start the timer on trade creation, it will publish the trade in 20 ms
				go func(pid *actor.PID) {
					time.Sleep(21 * time.Millisecond)
					context.Send(pid, &postAggTrade{})
				}(context.Self())
			}

			state.instrumentData.aggTrade.Trades = append(
				state.instrumentData.aggTrade.Trades,
				models.Trade{
					Price:    tradeData.Price,
					Quantity: tradeData.Quantity,
					ID:       uint64(tradeData.TradeID),
				})
		}

		if err := state.checkSockets(context); err != nil {
			return fmt.Errorf("error checking sockets: %v", err)
		}
		state.postHeartBeat(context)
		context.Send(context.Self(), &readSocket{})
		return nil

	case <-time.After(1 * time.Second):
		if err := state.checkSockets(context); err != nil {
			return fmt.Errorf("error checking sockets: %v", err)
		}
		state.postHeartBeat(context)
		context.Send(context.Self(), &readSocket{})
		return nil
	}
}

func (state *Listener) onDepthData(context actor.Context, depthData binance.WSDepthData) error {

	symbol := depthData.Symbol

	// Skip depth that are younger than OB
	if depthData.FinalUpdateID <= state.instrumentData.lastUpdateID {
		return nil
	}

	// Check depth continuity
	if state.instrumentData.lastUpdateID+1 != depthData.FirstUpdateID {
		return fmt.Errorf("got wrong sequence ID for %s: %d, %d",
			symbol, state.instrumentData.lastUpdateID, depthData.FirstUpdateID)
	}

	bids, asks, err := depthData.ToBidAsk()
	if err != nil {
		return fmt.Errorf("error converting depth data: %s ", err.Error())
	}

	obDelta := &models.OBL2Update{
		Levels:    []gorderbook.OrderBookLevel{},
		Timestamp: utils.MilliToTimestamp(depthData.EventTime),
		Trade:     false,
	}

	for _, bid := range bids {
		obDelta.Levels = append(
			obDelta.Levels,
			bid,
		)
		state.instrumentData.orderBook.UpdateOrderBookLevel(bid)
	}

	for _, ask := range asks {
		obDelta.Levels = append(
			obDelta.Levels,
			ask,
		)
		state.instrumentData.orderBook.UpdateOrderBookLevel(ask)
	}

	if state.instrumentData.orderBook.Crossed() {
		fmt.Println("CROSSED")
		return fmt.Errorf("crossed order book")
	}

	state.instrumentData.lastUpdateID = depthData.FinalUpdateID
	state.instrumentData.lastUpdateTime = depthData.EventTime
	context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
		UpdateL2: obDelta,
		SeqNum:   state.instrumentData.seqNum + 1,
	})
	state.instrumentData.seqNum += 1

	return nil
}

func (state *Listener) checkSockets(context actor.Context) error {
	// TODO ping or HB ?
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
