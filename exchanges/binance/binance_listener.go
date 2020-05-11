package binance

import (
	"container/list"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/quickfixgo/enum"
	mdir "github.com/quickfixgo/fix50/marketdataincrementalrefresh"
	"github.com/shopspring/decimal"
	exchangeModels "gitlab.com/alphaticks/alphac/messages/exchanges"
	"gitlab.com/alphaticks/alphac/messages/executor"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/binance"
	"reflect"
	"time"
)

// OBType: OBL2
// OBL2 Timestamps: ordered & consistent with sequence ID
// Trades: Impossible to infer from deltas
// Status: ready
type readSocket struct{}
type postAggTrade struct{}

type InstrumentData struct {
	orderBook        *gorderbook.OrderBookL2
	lastUpdateID     uint64
	lastUpdateTime   uint64
	lastHBTime       time.Time
	aggTrade         *mdir.NoMDEntriesRepeatingGroup
	aggTradeID       uint64
	lastAggTradeTime time.Time
}

type Listener struct {
	obWs            *binance.Websocket
	tradeWs         *binance.Websocket
	wsChan          chan *binance.WebsocketMessage
	instrument      *exchanges.Instrument
	instrumentData  *InstrumentData
	executorManager *actor.PID
	logger          *log.Logger
	stashedTrades   *list.List
}

func NewListenerProducer(instr *exchanges.Instrument) actor.Producer {
	return func() actor.Actor {
		return NewListener(instr)
	}
}

func NewListener(instr *exchanges.Instrument) actor.Actor {
	return &Listener{
		obWs:            nil,
		tradeWs:         nil,
		wsChan:          nil,
		instrument:      instr,
		instrumentData:  nil,
		executorManager: nil,
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

	case *executor.GetOrderBookL2Request:
		if err := state.GetOrderBookL2Request(context); err != nil {
			state.logger.Error("error processing GetOrderBookL2Request", log.Error(err))
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
		log.String("instrument", state.instrument.DefaultFormat()))

	state.executorManager = actor.NewLocalPID("exchange_executor_manager")
	state.wsChan = make(chan *binance.WebsocketMessage, 10000)
	state.stashedTrades = list.New()

	context.Send(context.Self(), &readSocket{})

	state.instrumentData = &InstrumentData{
		orderBook:        nil,
		lastUpdateID:     0,
		lastUpdateTime:   0,
		lastHBTime:       time.Now(),
		aggTrade:         nil,
		aggTradeID:       0,
		lastAggTradeTime: nil,
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
	err := obWs.Connect(
		state.instrument.Format(binance.WSSymbolFormat),
		[]string{binance.WSDepthStream100ms})
	if err != nil {
		return err
	}

	state.obWs = obWs

	time.Sleep(5 * time.Second)
	fut := context.RequestFuture(
		state.executorManager,
		&executor.GetOrderBookL2Request{Instrument: state.instrument, RequestID: 0},
		20*time.Second)

	res, err := fut.Result()
	if err != nil {
		return fmt.Errorf("error getting OBL2")
	}
	msg := res.(*executor.GetOrderBookL2Response)
	if msg.Error != nil {
		return fmt.Errorf("error while receiving OBL2 for %v", msg.Error)
	}

	bestAsk := float64(msg.Snapshot.Asks[0].Price) / float64(msg.Snapshot.Instrument.TickPrecision)
	depth := int(((bestAsk * 1.1) - bestAsk) * float64(msg.Snapshot.Instrument.TickPrecision))

	if depth > 10000 {
		depth = 10000
	}

	ob := gorderbook.NewOrderBookL2(
		state.instrument.TickPrecision,
		state.instrument.LotPrecision,
		depth,
	)

	ob.RawSync(msg.Snapshot.Bids, msg.Snapshot.Asks)
	state.instrumentData.orderBook = ob
	state.instrumentData.lastUpdateID = msg.Snapshot.ID
	state.instrumentData.lastUpdateTime = utils.TimestampToMilli(msg.Snapshot.Timestamp)

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
	err := tradeWs.Connect(
		state.instrument.Format(binance.WSSymbolFormat),
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

func (state *Listener) GetOrderBookL2Request(context actor.Context) error {
	msg := context.Message().(*executor.GetOrderBookL2Request)
	if msg.Instrument.DefaultFormat() != state.instrument.DefaultFormat() {
		err := fmt.Errorf("order book not tracked by listener")
		context.Respond(&executor.GetOrderBookL2Response{
			RequestID: msg.RequestID,
			Error:     err,
			Snapshot:  nil})
		return nil
	}

	snapshot := &exchangeModels.OBL2Snapshot{
		Instrument: msg.Instrument,
		Bids:       state.instrumentData.orderBook.GetAbsoluteRawBids(0),
		Asks:       state.instrumentData.orderBook.GetAbsoluteRawAsks(0),
		Timestamp:  utils.MilliToTimestamp(state.instrumentData.lastUpdateTime),
		ID:         state.instrumentData.lastUpdateID,
	}
	context.Respond(&executor.GetOrderBookL2Response{
		RequestID: msg.RequestID,
		Error:     nil,
		Snapshot:  snapshot})

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

			var entry mdir.NoMDEntries
			var tradeTime time.Time

			if state.instrumentData.aggTradeID != aggregateID {
				if state.instrumentData.lastAggTradeTime.Equal(msg.Time) {
					tradeTime = msg.Time.Add(time.Millisecond)
				} else {
					tradeTime = msg.Time
				}

				// Create a new one
				aggTrade := mdir.NewNoMDEntriesRepeatingGroup()
				entry = aggTrade.Add()
				state.instrumentData.aggTradeID = aggregateID
				state.instrumentData.aggTrade = &aggTrade
				state.stashedTrades.PushBack(&aggTrade)
				go func(pid *actor.PID) {
					time.Sleep(21 * time.Millisecond)
					context.Send(pid, &postAggTrade{})
				}(context.Self())
			} else {
				entry = state.instrumentData.aggTrade.Add()
				// Force same trade time for aggregate trade
				tradeTime = state.instrumentData.lastAggTradeTime
			}

			entry.SetMDUpdateAction(enum.MDUpdateAction_NEW)
			entry.SetMDEntryType(enum.MDEntryType_TRADE)
			entry.SetMDEntryID(fmt.Sprintf("%d", tradeData.TradeID))
			entry.SetSymbol(state.instrument.DefaultFormat())
			rawPrice := int64(tradeData.Price * float64(state.instrument.TickPrecision))
			entry.SetMDEntryPx(
				decimal.New(rawPrice, int32(state.instrument.TickPrecision)),
				int32(state.instrument.TickPrecision))
			rawSize := int64(tradeData.Quantity * float64(state.instrument.LotPrecision))
			entry.SetMDEntrySize(
				decimal.New(rawSize, int32(state.instrument.LotPrecision)),
				int32(state.instrument.LotPrecision))
			entry.SetMDEntryDate(tradeTime.Format(utils.FIX_UTCDateOnly_LAYOUT))
			entry.SetMDEntryTime(tradeTime.Format(utils.FIX_UTCTimeOnly_LAYOUT))
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

	bids, asks, err := depthData.ToRawBidAsk(state.instrument.TickPrecision, state.instrument.LotPrecision)
	if err != nil {
		return fmt.Errorf("error converting depth data: %s ", err.Error())
	}

	rg := mdir.NewNoMDEntriesRepeatingGroup()
	for _, bid := range bids {
		entry := rg.Add()
		entry.SetSymbol(state.instrument.DefaultFormat())
		entry.SetMDEntryType(enum.MDEntryType_BID)
		entry.SetMDEntryPx(decimal.New(int64(bid.Price), int32(state.instrument.TickPrecision)), int32(state.instrument.TickPrecision))
		entry.SetMDEntrySize(decimal.New(int64(bid.Quantity), int32(state.instrument.LotPrecision)), int32(state.instrument.LotPrecision))
		state.instrumentData.orderBook.UpdateRawOrderBookLevel(bid)
	}

	for _, ask := range asks {
		entry := rg.Add()
		entry.SetSymbol(state.instrument.DefaultFormat())
		entry.SetMDEntryType(enum.MDEntryType_OFFER)
		entry.SetMDEntryPx(decimal.New(int64(ask.Price), int32(state.instrument.TickPrecision)), int32(state.instrument.TickPrecision))
		entry.SetMDEntrySize(decimal.New(int64(ask.Quantity), int32(state.instrument.LotPrecision)), int32(state.instrument.LotPrecision))
		state.instrumentData.orderBook.UpdateRawOrderBookLevel(ask)
	}

	if state.instrumentData.orderBook.Crossed() {
		return fmt.Errorf("crossed order book")
	}

	state.instrumentData.lastUpdateID = depthData.FinalUpdateID
	state.instrumentData.lastUpdateTime = depthData.EventTime

	refresh := mdir.New()
	refresh.SetMDBookType(enum.MDBookType_PRICE_DEPTH)
	refresh.SetNoMDEntries(rg)

	context.Send(context.Parent(), &refresh)

	state.instrumentData.lastHBTime = time.Now()

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

func (state *Listener) postHeartBeat(context actor.Context) {
	// If haven't sent anything for 2 seconds, send heartbeat
	if time.Now().Sub(state.instrumentData.lastHBTime) > 2*time.Second {
		//topic := fmt.Sprintf("%s/HEARTBEAT", state.instrument.DefaultFormat())
		// TODO HB ?
		context.Send(context.Parent(), nil)
		state.instrumentData.lastHBTime = time.Now()
	}
}

func (state *Listener) postAggTrade(context actor.Context) {
	for el := state.stashedTrades.Front(); el != nil; el = state.stashedTrades.Front() {
		trd := el.Value.(*mdir.NoMDEntriesRepeatingGroup)
		entry := trd.Get(0)
		entryDate, _ := entry.GetMDEntryDate()
		entryTime, _ := entry.GetMDEntryTime()
		ts, _ := time.Parse(utils.FIX_UTCDateTime_LAYOUT, entryDate+"-"+entryTime)
		if time.Now().Sub(ts) > 20*time.Millisecond {
			refresh := mdir.New()
			refresh.SetNoMDEntries(*trd)
			context.Send(context.Parent(), &refresh)
			// At this point, the state.instrumentData.aggTrade can be our trade, or it can be a new one
			if state.instrumentData.aggTrade == trd {
				state.instrumentData.aggTrade = nil
				state.instrumentData.aggTradeID = 0
			}
			state.stashedTrades.Remove(el)
		} else {
			break
		}
	}
}
