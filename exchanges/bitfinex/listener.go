package bitfinex

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
	gmodels "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	"gitlab.com/alphaticks/xchanger"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/bitfinex"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
	"reflect"
	"strings"
	"sync"
	"time"
)

var fbd = 5 * time.Minute
var fbdLock = sync.RWMutex{}

type checkSockets struct{}
type postAggTrade struct{}
type updateBook struct{}

type InstrumentData struct {
	tickPrecision   uint64
	lotPrecision    uint64
	orderBook       *gorderbook.OrderBookL2
	fullBook        *gorderbook.OrderBookL2
	seqNum          uint64
	lastUpdateTime  uint64
	lastFundingTime uint64
	lastHBTime      time.Time
	lastSequence    uint64
	aggTrade        *models.AggregatedTrade
	lastAggTradeTs  uint64
}

// OBType: OBL3
// OBL3 timestamps: per connection
// Status: not ready problem with timestamps

type Listener struct {
	ws             *bitfinex.Websocket
	dialerPool     *xchangerUtils.DialerPool
	security       *models.Security
	instrumentData *InstrumentData
	logger         *log.Logger
	stashedTrades  *list.List
	socketTicker   *time.Ticker
	fullBookTicker *time.Ticker
	executor       *actor.PID
}

func NewListenerProducer(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Producer {
	return func() actor.Actor {
		return NewListener(security, dialerPool)
	}
}

// Limit of 30 subscription
func NewListener(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Actor {
	return &Listener{
		ws:             nil,
		security:       security,
		dialerPool:     dialerPool,
		instrumentData: nil,
		logger:         nil,
		stashedTrades:  nil,
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

	case *messages.MarketDataResponse:
		if err := state.onMarketDataResponse(context); err != nil {
			state.logger.Error("error processing market data response", log.Error(err))
		}

	case *checkSockets:
		if err := state.checkSockets(context); err != nil {
			state.logger.Error("error checking socket", log.Error(err))
			panic(err)
		}

	case *updateBook:
		if err := state.updateFullBook(context); err != nil {
			state.logger.Error("error updating full book", log.Error(err))
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

	state.executor = actor.NewPID(context.ActorSystem().Address(), "executor/exchanges/"+constants.BITFINEX.Name+"_executor")
	state.stashedTrades = list.New()

	if state.security.RoundLot == nil {
		return fmt.Errorf("security is missing RoundLot")
	}
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot.Value))

	state.instrumentData = &InstrumentData{
		lotPrecision:    lotPrecision,
		orderBook:       nil,
		seqNum:          uint64(time.Now().UnixNano()),
		lastUpdateTime:  0,
		lastFundingTime: 0,
		lastHBTime:      time.Now(),
		aggTrade:        nil,
		lastAggTradeTs:  0,
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
					return
				}
			}
		}
	}(context.Self())

	fullBookTicker := time.NewTicker(fbd)
	state.fullBookTicker = fullBookTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-fullBookTicker.C:
				context.Send(pid, &updateBook{})
			case <-time.After(10 * time.Second):
				if state.fullBookTicker != fullBookTicker {
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

	if state.fullBookTicker != nil {
		state.fullBookTicker.Stop()
		state.fullBookTicker = nil
	}

	return nil
}

func (state *Listener) subscribeInstrument(context actor.Context) error {
	if state.ws != nil {
		_ = state.ws.Disconnect()
	}

	ws := bitfinex.NewWebsocket()

	if err := ws.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to bitfinex websocket: %v", err)
	}

	symbol := "t" + strings.ToUpper(state.security.Symbol)
	err := ws.SubscribeDepth(
		symbol,
		bitfinex.WSDepthPrecisionP0,
		bitfinex.WSDepthFrequency0,
		bitfinex.WSDepthLength100)
	if err != nil {
		return err
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	_, ok := ws.Msg.Message.(bitfinex.WSSubscribeDepthResponse)
	if !ok {
		if errMsg, ok := ws.Msg.Message.(bitfinex.WSErrorMessage); ok {
			return fmt.Errorf("error subscribing to depth: %s", errMsg.Msg)
		} else {
			return fmt.Errorf("was expecting WSSubscribeDepthResponse, got %s", reflect.TypeOf(ws.Msg.Message).String())
		}
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	snapshot, ok := ws.Msg.Message.(bitfinex.WSSpotDepthSnapshot)
	if !ok {
		return fmt.Errorf("was expecting WSSpotDepthSnapshot, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	state.instrumentData.lastSequence = snapshot.Sequence

	bids, asks := snapshot.ToBidAsk()

	avgPrice := 0.
	n := 0.
	for i := 0; i < 25 && i < len(bids); i++ {
		avgPrice += bids[i].Price
		n += 1.
	}
	avgPrice /= n

	lg := math.Log10(avgPrice)
	digits := int(math.Ceil(lg))
	maxTickPrecisionF := math.Pow10(5 - digits)
	tickPrecision := uint64(maxTickPrecisionF)
	if tickPrecision == 0 {
		tickPrecision = 1
	}
	state.instrumentData.tickPrecision = tickPrecision

	bestAsk := asks[0].Price
	//Allow a 10% price variation
	depth := int(((bestAsk * 1.1) - bestAsk) * float64(tickPrecision))
	if depth > 10000 {
		depth = 10000
	}

	ob := gorderbook.NewOrderBookL2(
		tickPrecision,
		state.instrumentData.lotPrecision,
		depth,
	)

	ts := uint64(ws.Msg.ClientTime.UnixNano() / 1000000)

	ob.Sync(bids, asks)
	if ob.Crossed() {
		return fmt.Errorf("crossed orderbook")
	}
	state.instrumentData.orderBook = ob
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())
	state.instrumentData.lastUpdateTime = ts

	// Try to fetch full book
	fullBook := gorderbook.NewOrderBookL2(
		tickPrecision,
		state.instrumentData.lotPrecision,
		depth,
	)
	state.instrumentData.fullBook = fullBook

	res, err := context.RequestFuture(
		state.executor,
		&messages.MarketDataRequest{
			RequestID: uint64(time.Now().UnixNano()),
			Subscribe: false,
			Instrument: &models.Instrument{
				SecurityID: &wrapperspb.UInt64Value{Value: state.security.SecurityID},
				Exchange:   state.security.Exchange,
				Symbol:     &wrapperspb.StringValue{Value: state.security.Symbol},
			},
			Aggregation: models.OrderBookAggregation_L2,
		}, 10*time.Second).Result()
	if err == nil {
		msg := res.(*messages.MarketDataResponse)
		if msg.Success && msg.SnapshotL2 != nil {
			// What to do ? ...
			// Remove all worstBid worstAsk
			bids := state.instrumentData.orderBook.GetBids(-1)
			worstBid := bids[len(bids)-1].Price
			asks := state.instrumentData.orderBook.GetAsks(-1)
			worstAsk := asks[len(asks)-1].Price

			bids = nil
			asks = nil

			for _, b := range msg.SnapshotL2.Bids {
				if b.Price < worstBid {
					// Add to book
					bids = append(bids, b)
				}
			}
			for _, a := range msg.SnapshotL2.Asks {
				if a.Price > worstAsk {
					// Add to book
					asks = append(asks, a)
				}
			}

			fullBook.Sync(bids, asks)
		}
	}

	if err := ws.SubscribeTrades(symbol); err != nil {
		return err
	}

	if state.security.SecurityType == enum.SecurityType_CRYPTO_PERP {
		fmt.Println("SUBSCRIBING PERP")
		if err := ws.SubscribeStatus(symbol); err != nil {
			return err
		}
		/*
			// TODO
			if err := ws.SubscribeLiquidation(); err != nil {
				return err
			}

		*/
	}

	state.ws = ws

	go func(ws *bitfinex.Websocket, pid *actor.PID) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
		}
	}(ws, context.Self())

	return nil
}

func (state *Listener) OnMarketDataRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketDataRequest)
	bids := state.instrumentData.fullBook.GetBids(0)
	asks := state.instrumentData.fullBook.GetAsks(0)
	bids = append(bids, state.instrumentData.orderBook.GetBids(0)...)
	asks = append(asks, state.instrumentData.orderBook.GetAsks(0)...)
	snapshot := &models.OBL2Snapshot{
		Bids:          bids,
		Asks:          asks,
		Timestamp:     utils.MilliToTimestamp(state.instrumentData.lastUpdateTime),
		TickPrecision: &wrapperspb.UInt64Value{Value: state.instrumentData.orderBook.TickPrecision},
		LotPrecision:  &wrapperspb.UInt64Value{Value: state.instrumentData.orderBook.LotPrecision},
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

	case bitfinex.WSErrorMessage:
		err := msg.Message.(bitfinex.WSErrorMessage)
		return fmt.Errorf("socket error: %v", err)

	case bitfinex.WSHeartBeat:
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			SeqNum: state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1

	case bitfinex.WSSpotDepthData:
		obData := msg.Message.(bitfinex.WSSpotDepthData)

		state.instrumentData.lastSequence = obData.Sequence

		level := obData.Depth.ToOrderBookLevel()
		// only allow 5 digits on price

		ts := uint64(msg.ClientTime.UnixNano() / 1000000)
		obDelta := &models.OBL2Update{
			Levels:    []*gmodels.OrderBookLevel{level},
			Timestamp: utils.MilliToTimestamp(ts),
			Trade:     false,
		}

		state.instrumentData.orderBook.UpdateOrderBookLevel(level)

		if state.instrumentData.orderBook.Crossed() {
			state.logger.Info("crossed orderbook", log.Error(errors.New("crossed")))
			return state.subscribeInstrument(context)
		}

		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			UpdateL2: obDelta,
			SeqNum:   state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1
		state.instrumentData.lastUpdateTime = utils.TimestampToMilli(obDelta.Timestamp)

	case bitfinex.WSSpotTrade:
		tradeData := msg.Message.(bitfinex.WSSpotTrade)
		ts := uint64(msg.ClientTime.UnixNano() / 1000000)

		state.instrumentData.lastSequence = tradeData.Sequence

		aggID := tradeData.Timestamp * 10
		if tradeData.Amount < 0 {
			aggID += 1
		}

		if state.instrumentData.aggTrade == nil || state.instrumentData.aggTrade.AggregateID != aggID {
			// Create new agg trade
			if state.instrumentData.lastAggTradeTs >= ts {
				ts = state.instrumentData.lastAggTradeTs + 1
			}
			aggTrade := &models.AggregatedTrade{
				Bid:         tradeData.Amount < 0,
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
				time.Sleep(21 * time.Millisecond)
				context.Send(pid, &postAggTrade{})
			}(context.Self())
		}

		// Unlike the depths, negative amounts mean bid change and floor
		if tradeData.Amount < 0 {
			state.instrumentData.aggTrade.Trades = append(
				state.instrumentData.aggTrade.Trades,
				&models.Trade{
					Price:    tradeData.Price,
					Quantity: -tradeData.Amount,
					ID:       tradeData.ID,
				})
		} else {
			state.instrumentData.aggTrade.Trades = append(
				state.instrumentData.aggTrade.Trades,
				&models.Trade{
					Price:    tradeData.Price,
					Quantity: tradeData.Amount,
					ID:       tradeData.ID,
				})
		}

	case bitfinex.WSStatus:
		status := msg.Message.(bitfinex.WSStatus)
		ts := uint64(msg.ClientTime.UnixNano() / 1000000)
		refresh := &messages.MarketDataIncrementalRefresh{
			Stats: []*models.Stat{{
				Timestamp: utils.MilliToTimestamp(ts),
				StatType:  models.StatType_OpenInterest,
				Value:     status.OpenInterest,
			}},
			SeqNum: state.instrumentData.seqNum + 1,
		}
		if state.instrumentData.lastFundingTime < status.FundingEventMs {
			refresh.Stats = append(refresh.Stats, &models.Stat{
				Timestamp: utils.MilliToTimestamp(status.FundingEventMs),
				StatType:  models.StatType_FundingRate,
				Value:     status.CurrentFunding,
			})
			state.instrumentData.lastFundingTime = status.FundingEventMs
		}
		context.Send(context.Parent(), refresh)
		state.instrumentData.seqNum += 1
		state.instrumentData.lastSequence = status.Sequence

	case bitfinex.WSLiquidation:
		liq := msg.Message.(bitfinex.WSLiquidation)
		ts := uint64(msg.ClientTime.UnixNano() / 1000000)
		fmt.Println("LIQ", liq)
		refresh := &messages.MarketDataIncrementalRefresh{
			Liquidation: &models.Liquidation{
				Bid:       liq.Amount > 0,
				Timestamp: utils.MilliToTimestamp(ts),
				OrderID:   liq.PositionID,
				Price:     liq.LiquidationPrice,
				Quantity:  math.Abs(liq.Amount),
			},
			SeqNum: state.instrumentData.seqNum + 1,
		}
		context.Send(context.Parent(), refresh)
		state.instrumentData.seqNum += 1

	case bitfinex.WSSpotTradeSnapshot:
		tradeData := msg.Message.(bitfinex.WSSpotTradeSnapshot)
		state.instrumentData.lastSequence = tradeData.Sequence

	case bitfinex.WSSpotTradeU:
		tradeData := msg.Message.(bitfinex.WSSpotTradeU)
		state.instrumentData.lastSequence = tradeData.Sequence

	case bitfinex.WSSubscribeSpotTradesResponse:
		break

	case bitfinex.WSSubscribeDepthResponse:
		break
	}

	return nil
}

func (state *Listener) onMarketDataResponse(context actor.Context) error {
	msg := context.Message().(*messages.MarketDataResponse)
	if !msg.Success {
		// We want to converge towards the right value,
		if msg.RejectionReason == messages.RejectionReason_RateLimitExceeded {
			fbdLock.Lock()
			fbd = time.Duration(float64(fbd) * 1.01)
			fbdLock.Unlock()
			if state.fullBookTicker != nil {
				state.fullBookTicker.Reset(fbd)
			}
			fmt.Println("INCREASE DELAY", fbd)
		}
		state.logger.Info("error fetching snapshot", log.Error(errors.New(msg.RejectionReason.String())))
		return nil
	}
	if msg.SnapshotL2 == nil {
		state.logger.Error("error fetching snapshot: no OBL2")
		return nil
	}

	// Reduce delay
	fbdLock.Lock()
	fbd = time.Duration(float64(fbd) * 0.999)
	fbdLock.Unlock()
	if state.fullBookTicker != nil {
		state.fullBookTicker.Reset(fbd)
	}
	fmt.Println("DECREASE DELAY", fbd)

	// What to do ? ...
	// Remove all worstBid worstAsk
	bids := state.instrumentData.orderBook.GetBids(-1)
	worstBid := bids[len(bids)-1].Price
	asks := state.instrumentData.orderBook.GetAsks(-1)
	worstAsk := asks[len(asks)-1].Price

	bids = nil
	asks = nil

	for _, b := range msg.SnapshotL2.Bids {
		if b.Price < worstBid {
			// Add to book
			bids = append(bids, b)
		}
	}
	for _, a := range msg.SnapshotL2.Asks {
		if a.Price > worstAsk {
			// Add to book
			asks = append(asks, a)
		}
	}

	fullBook := gorderbook.NewOrderBookL2(
		state.instrumentData.fullBook.TickPrecision,
		state.instrumentData.fullBook.LotPrecision,
		state.instrumentData.fullBook.Depth)

	fullBook.Sync(bids, asks)

	levels := state.instrumentData.fullBook.Diff(fullBook)

	state.instrumentData.fullBook = fullBook

	context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
		UpdateL2: &models.OBL2Update{
			Levels:    levels,
			Timestamp: msg.SnapshotL2.Timestamp,
			Trade:     false,
		},
		SeqNum: state.instrumentData.seqNum + 1,
	})
	state.instrumentData.seqNum += 1
	state.instrumentData.lastUpdateTime = utils.TimestampToMilli(msg.SnapshotL2.Timestamp)

	return nil
}

func (state *Listener) checkSockets(context actor.Context) error {
	// No need to ping HB mechanism already
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

func (state *Listener) updateFullBook(context actor.Context) error {
	context.Request(
		state.executor,
		&messages.MarketDataRequest{
			RequestID: uint64(time.Now().UnixNano()),
			Subscribe: false,
			Instrument: &models.Instrument{
				SecurityID: &wrapperspb.UInt64Value{Value: state.security.SecurityID},
				Exchange:   state.security.Exchange,
				Symbol:     &wrapperspb.StringValue{Value: state.security.Symbol},
			},
			Aggregation: models.OrderBookAggregation_L2,
		})

	return nil
}
