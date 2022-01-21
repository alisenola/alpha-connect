package ftx

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
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/ftx"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"math"
	"reflect"
	"sort"
	"sync"
	"time"
)

var oid = 2 * time.Minute
var oidLock = sync.RWMutex{}

type checkSockets struct{}
type updateOpenInterest struct{}

type InstrumentData struct {
	orderBook      *gorderbook.OrderBookL2
	seqNum         uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	lastAggTradeTs uint64
}

type Listener struct {
	ws                 *ftx.Websocket
	security           *models.Security
	dialerPool         *xchangerUtils.DialerPool
	instrumentData     *InstrumentData
	executor           *actor.PID
	logger             *log.Logger
	lastPingTime       time.Time
	socketTicker       *time.Ticker
	openInterestTicker *time.Ticker
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

	case *updateOpenInterest:
		if err := state.updateOpenInterest(context); err != nil {
			state.logger.Error("error updating open interest", log.Error(err))
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
	state.executor = actor.NewPID(context.ActorSystem().Address(), "executor/"+constants.FTX.Name+"_executor")

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

	if state.openInterestTicker != nil {
		state.openInterestTicker.Stop()
		state.openInterestTicker = nil
	}

	return nil
}

func (state *Listener) subscribeInstrument(context actor.Context) error {
	if state.ws != nil {
		_ = state.ws.Disconnect()
	}

	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))

	ws := ftx.NewWebsocket()
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
	_, ok := ws.Msg.Message.(ftx.WSSubscribeResponse)
	if !ok {
		if err, ok := ws.Msg.Message.(ftx.WSError); ok {
			return fmt.Errorf("got WSError trying to subscribe to ob: %d: %s", err.Msg, err.Code)
		} else {
			return fmt.Errorf("was expecting WSSubscribeResponse, got %s", reflect.TypeOf(ws.Msg.Message).String())
		}
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	obUpdate, ok := ws.Msg.Message.(ftx.WSOrderBookGroupedUpdate)
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
		return fmt.Errorf("crossed orderbook")
	}
	state.instrumentData.orderBook = ob
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())
	state.instrumentData.lastUpdateTime = ts

	if err := ws.Subscribe(state.security.Symbol, ftx.WSTradeChannel); err != nil {
		return fmt.Errorf("error subscribing to trade stream: %v", err)
	}

	state.ws = ws

	go func(ws *ftx.Websocket, pid *actor.PID) {
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

	if msg.Subscribe {
		for _, stat := range msg.Stats {
			if stat == models.OpenInterest && state.openInterestTicker != nil {
				if state.security.SecurityType == enum.SecurityType_CRYPTO_PERP || state.security.SecurityType == enum.SecurityType_CRYPTO_FUT {
					openInterestTicker := time.NewTicker(10 * time.Second)
					state.openInterestTicker = openInterestTicker
					go func(pid *actor.PID) {
						for {
							select {
							case _ = <-openInterestTicker.C:
								context.Send(pid, &updateOpenInterest{})
							case <-time.After(11 * time.Second):
								if state.openInterestTicker != openInterestTicker {
									return
								}
							}
						}
					}(context.Self())
				}
			}
		}
	}

	context.Respond(response)
	return nil
}

func (state *Listener) onWebsocketMessage(context actor.Context) error {
	msg := context.Message().(*xchanger.WebsocketMessage)
	switch msg.Message.(type) {

	case error:
		return fmt.Errorf("socket error: %v", msg)

	case ftx.WSOrderBookGroupedUpdate:
		obData := msg.Message.(ftx.WSOrderBookGroupedUpdate)
		instr := state.instrumentData
		nLevels := len(obData.Snapshot.Bids) + len(obData.Snapshot.Asks)

		ts := uint64(msg.ClientTime.UnixNano()) / 1000000
		obDelta := &models.OBL2Update{
			Levels:    make([]gmodels.OrderBookLevel, nLevels, nLevels),
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

	case ftx.WSTradeUpdate:
		tradeData := msg.Message.(ftx.WSTradeUpdate)
		ts := uint64(msg.ClientTime.UnixNano() / 1000000)

		sort.Slice(tradeData.Trades, func(i, j int) bool {
			return tradeData.Trades[i].Time.Before(tradeData.Trades[j].Time)
		})

		var aggTrade *models.AggregatedTrade
		for _, trade := range tradeData.Trades {
			if trade.Liquidation {
				// Liquidation limit order, so was on the bid if order.Side == "sell"
				context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
					Liquidation: &models.Liquidation{
						Bid:       trade.Side == "buy",
						Timestamp: utils.MilliToTimestamp(ts),
						OrderID:   trade.ID,
						Price:     trade.Price,
						Quantity:  trade.Size,
					},
					SeqNum: state.instrumentData.seqNum + 1,
				})
				state.instrumentData.seqNum += 1
			}
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

	if time.Now().Sub(state.lastPingTime) > 10*time.Second {
		// "Ping" by resubscribing to the topic
		if err := state.ws.SubscribeGroupedOrderBook(state.security.Symbol, state.security.MinPriceIncrement.Value); err != nil {
			return fmt.Errorf("error subscribing to OBL2 stream: %v", err)
		}
		if err := state.ws.Subscribe(state.security.Symbol, ftx.WSTradeChannel); err != nil {
			return fmt.Errorf("error subscribing to trade stream: %v", err)
		}
		state.lastPingTime = time.Now()
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

func (state *Listener) updateOpenInterest(context actor.Context) error {
	fmt.Println("UPDATE OI", oid)
	fut := context.RequestFuture(
		state.executor,
		&messages.MarketStatisticsRequest{
			RequestID: uint64(time.Now().UnixNano()),
			Instrument: &models.Instrument{
				SecurityID: &types.UInt64Value{Value: state.security.SecurityID},
				Exchange:   state.security.Exchange,
				Symbol:     &types.StringValue{Value: state.security.Symbol},
			},
			Statistics: []models.StatType{models.OpenInterest, models.FundingRate},
		}, 2*time.Second)

	context.AwaitFuture(fut, func(res interface{}, err error) {
		if err != nil {
			if err == actor.ErrTimeout {
				oidLock.Lock()
				oid = time.Duration(float64(oid) * 1.01)
				state.openInterestTicker.Reset(oid)
				oidLock.Unlock()
			}
			state.logger.Info("error fetching market statistics", log.Error(err))
			return
		}
		msg := res.(*messages.MarketStatisticsResponse)
		if !msg.Success {
			// We want to converge towards the right value,
			if msg.RejectionReason == messages.RateLimitExceeded || msg.RejectionReason == messages.HTTPError {
				oidLock.Lock()
				oid = time.Duration(float64(oid) * 1.01)
				state.openInterestTicker.Reset(oid)
				oidLock.Unlock()
			}
			state.logger.Info("error fetching market statistics", log.Error(errors.New(msg.RejectionReason.String())))
			return
		}

		fmt.Println(msg.Statistics)
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			Stats:  msg.Statistics,
			SeqNum: state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1

		// Reduce delay
		oidLock.Lock()
		oid = time.Duration(float64(oid) * 0.99)
		if oid < 15*time.Second {
			oid = 15 * time.Second
		}
		state.openInterestTicker.Reset(oid)
		oidLock.Unlock()
	})

	return nil
}
