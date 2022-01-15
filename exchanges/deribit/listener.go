package deribit

import (
	"errors"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/gorderbook"
	gmodels "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	"gitlab.com/alphaticks/xchanger"
	"gitlab.com/alphaticks/xchanger/exchanges/deribit"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"math"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type checkSockets struct{}

type InstrumentData struct {
	orderBook      *gorderbook.OrderBookL2
	seqNum         uint64
	lastUpdateTime uint64
	lastUpdateID   uint64
	lastHBTime     time.Time
	lastAggTradeTs uint64
	openInterest   float64
	funding8h      float64
}

type Listener struct {
	ws             *deribit.Websocket
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

	ws := deribit.NewWebsocket()
	err := ws.Connect(state.dialerPool.GetDialer())
	if err != nil {
		return err
	}
	var interval string
	if os.Getenv("DERIBIT_KEY") != "" {
		interval = deribit.Interval0ms
		if err, _ := ws.Auth(&xchangerModels.APICredentials{
			APIKey:    os.Getenv("DERIBIT_KEY"),
			APISecret: os.Getenv("DERIBIT_SECRET"),
		}); err != nil {
			return err
		}
	} else {
		interval = deribit.Interval100ms
	}

	if err, _ := ws.SubscribeOrderBook(state.security.Symbol, interval); err != nil {
		return err
	}

	if !ws.ReadMessage() {
		return ws.Err
	}
	_, ok := ws.Msg.Message.(deribit.Subscription)
	if !ok {
		return fmt.Errorf("error casting message to Subscription")
	}

	if !ws.ReadMessage() {
		return ws.Err
	}

	update, ok := ws.Msg.Message.(deribit.OrderBookUpdate)
	if !ok {
		return fmt.Errorf("error casting message to OrderBookUpdate")
	}

	var bids, asks []gmodels.OrderBookLevel
	bids = make([]gmodels.OrderBookLevel, len(update.Bids), len(update.Bids))
	for i, bid := range update.Bids {
		bids[i] = gmodels.OrderBookLevel{
			Price:    bid.Price,
			Quantity: bid.Quantity,
			Bid:      true,
		}
	}
	asks = make([]gmodels.OrderBookLevel, len(update.Asks), len(update.Asks))
	for i, ask := range update.Asks {
		asks[i] = gmodels.OrderBookLevel{
			Price:    ask.Price,
			Quantity: ask.Quantity,
			Bid:      false,
		}
	}
	if len(update.Asks) == 0 || len(update.Bids) == 0 {
		return fmt.Errorf("empty bid or ask")
	}
	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot.Value))
	bestAsk := update.Asks[0].Price
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
		return fmt.Errorf("crossed orderbook")
	}
	state.instrumentData.orderBook = ob
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())
	state.instrumentData.lastUpdateTime = ts
	state.instrumentData.lastUpdateID = update.ChangeID

	if err, _ := ws.SubscribeTrade(state.security.Symbol, interval); err != nil {
		return fmt.Errorf("error subscribing to trade stream: %v", err)
	}

	if err, _ := ws.SubscribeTicker(state.security.Symbol, interval); err != nil {
		return fmt.Errorf("error subscribing to ticker stream: %v", err)
	}

	state.ws = ws

	go func(ws *deribit.Websocket, pid *actor.PID) {
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

	case deribit.OrderBookUpdate:
		obData := msg.Message.(deribit.OrderBookUpdate)
		instr := state.instrumentData

		if obData.ChangeID <= instr.lastUpdateID {
			break
		}

		if obData.PreviousChangeID != instr.lastUpdateID {
			state.logger.Info("error processing ob update", log.Error(fmt.Errorf("out of order sequence")))
			return state.subscribeInstrument(context)
		}

		nLevels := len(obData.Bids) + len(obData.Asks)

		ts := uint64(msg.ClientTime.UnixNano()) / 1000000
		obDelta := &models.OBL2Update{
			Levels:    make([]gmodels.OrderBookLevel, nLevels, nLevels),
			Timestamp: utils.MilliToTimestamp(ts),
			Trade:     false,
		}

		lvlIdx := 0
		for _, bid := range obData.Bids {
			level := gmodels.OrderBookLevel{
				Price:    bid.Price,
				Quantity: bid.Quantity,
				Bid:      true,
			}
			obDelta.Levels[lvlIdx] = level
			lvlIdx += 1
			instr.orderBook.UpdateOrderBookLevel(level)
		}
		for _, ask := range obData.Asks {
			level := gmodels.OrderBookLevel{
				Price:    ask.Price,
				Quantity: ask.Quantity,
				Bid:      false,
			}
			obDelta.Levels[lvlIdx] = level
			lvlIdx += 1
			instr.orderBook.UpdateOrderBookLevel(level)
		}

		if state.instrumentData.orderBook.Crossed() {
			fmt.Println(state.instrumentData.orderBook)
			state.logger.Info("crossed orderbook", log.Error(errors.New("crossed")))
			return state.subscribeInstrument(context)
		}
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			UpdateL2: obDelta,
			SeqNum:   state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1
		state.instrumentData.lastUpdateID = obData.ChangeID
		state.instrumentData.lastUpdateTime = ts

	case deribit.TradeUpdate:
		tradeData := msg.Message.(deribit.TradeUpdate)
		ts := uint64(msg.ClientTime.UnixNano() / 1000000)

		sort.Slice(tradeData, func(i, j int) bool {
			return tradeData[i].Timestamp < tradeData[j].Timestamp
		})

		var aggTrade *models.AggregatedTrade
		var aggHelpR uint64 = 0

		for _, trade := range tradeData {
			aggHelp := trade.Timestamp
			if trade.Direction == "buy" {
				aggHelp += 1
			}
			splits := strings.Split(trade.TradeID, "-")
			tradeID, err := strconv.ParseInt(splits[len(splits)-1], 10, 64)
			if err != nil {
				return fmt.Errorf("error parsing trade ID: %s %v", trade.TradeID, err)
			}
			if trade.Liquidation != nil {
				if *trade.Liquidation == "MT" {
					context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
						Liquidation: &models.Liquidation{
							Bid:       true,
							Timestamp: utils.MilliToTimestamp(ts),
							OrderID:   uint64(tradeID),
							Price:     trade.Price,
							Quantity:  trade.Amount,
						},
						SeqNum: state.instrumentData.seqNum + 1,
					})
					state.instrumentData.seqNum += 1
					context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
						Liquidation: &models.Liquidation{
							Bid:       false,
							Timestamp: utils.MilliToTimestamp(ts),
							OrderID:   uint64(tradeID),
							Price:     trade.Price,
							Quantity:  trade.Amount,
						},
						SeqNum: state.instrumentData.seqNum + 1,
					})
					state.instrumentData.seqNum += 1
				} else if *trade.Liquidation == "T" {
					// Taker side of the trade was liquidated
					// So if direction is "buy", the liquidation
					// was on the bid
					context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
						Liquidation: &models.Liquidation{
							Bid:       trade.Direction == "buy",
							Timestamp: utils.MilliToTimestamp(ts),
							OrderID:   uint64(tradeID),
							Price:     trade.Price,
							Quantity:  trade.Amount,
						},
						SeqNum: state.instrumentData.seqNum + 1,
					})
					state.instrumentData.seqNum += 1
				} else if *trade.Liquidation == "M" {
					// Maker side of the trade was liquidated
					// So if direction is "sell", the liquidation
					// was on the bid
					context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
						Liquidation: &models.Liquidation{
							Bid:       trade.Direction == "sell",
							Timestamp: utils.MilliToTimestamp(ts),
							OrderID:   uint64(tradeID),
							Price:     trade.Price,
							Quantity:  trade.Amount,
						},
						SeqNum: state.instrumentData.seqNum + 1,
					})
					state.instrumentData.seqNum += 1
				}
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
					Bid:         trade.Direction == "sell",
					Timestamp:   utils.MilliToTimestamp(ts),
					AggregateID: uint64(tradeID),
					Trades:      nil,
				}
				aggHelpR = aggHelp
			}

			trade := models.Trade{
				Price:    trade.Price,
				Quantity: trade.Amount,
				ID:       uint64(tradeID),
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

	case deribit.TickerUpdate:
		tickerData := msg.Message.(deribit.TickerUpdate)
		ts := uint64(msg.ClientTime.UnixNano() / 1000000)
		refresh := &messages.MarketDataIncrementalRefresh{
			SeqNum: state.instrumentData.seqNum + 1,
		}
		update := false
		if state.instrumentData.openInterest != tickerData.OpenInterest {
			refresh.Stats = append(refresh.Stats, &models.Stat{
				Timestamp: utils.MilliToTimestamp(ts),
				StatType:  models.OpenInterest,
				Value:     tickerData.OpenInterest,
			})
			update = true
			state.instrumentData.openInterest = tickerData.OpenInterest
		}
		if state.security.SecurityType == enum.SecurityType_CRYPTO_PERP && state.instrumentData.funding8h != tickerData.Funding8h {
			refresh.Funding = &models.Funding{
				Timestamp: utils.MilliToTimestamp(ts),
				Rate:      tickerData.Funding8h,
			}
			state.instrumentData.funding8h = tickerData.Funding8h
		}
		if update {
			context.Send(context.Parent(), refresh)
			state.instrumentData.seqNum += 1
		}
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
