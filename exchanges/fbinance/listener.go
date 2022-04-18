package fbinance

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
	"gitlab.com/alphaticks/xchanger/exchanges/fbinance"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"math"
	"math/rand"
	"reflect"
	"strings"
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
	lastUpdateID   uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	lastAggTradeTs uint64
}

type Listener struct {
	ws                 *fbinance.Websocket
	security           *models.Security
	dialerPool         *xchangerUtils.DialerPool
	instrumentData     *InstrumentData
	executor           *actor.PID
	logger             *log.Logger
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
	state.executor = actor.NewPID(context.ActorSystem().Address(), "executor/exchanges/"+constants.FBINANCE.Name+"_executor")

	state.instrumentData = &InstrumentData{
		orderBook:      nil,
		seqNum:         uint64(time.Now().UnixNano()),
		lastUpdateID:   0,
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
			state.logger.Warn("error disconnecting socket", log.Error(err))
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

	ws := fbinance.NewWebsocket()
	symbol := strings.ToLower(state.security.Symbol)
	streams := []string{
		fbinance.WSDepthStreamRealTime,
		fbinance.WSAggregatedTradeStream,
		fbinance.WSForceOrderStream,
	}
	if state.security.SecurityType == enum.SecurityType_CRYPTO_PERP {
		streams = append(streams, fbinance.WSMarkPriceStream1000ms)
	}
	err := ws.Connect(
		symbol,
		streams,
		state.dialerPool.GetDialer())
	if err != nil {
		return err
	}

	state.ws = ws

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

	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot.Value))
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
	if ob.Crossed() {
		return fmt.Errorf("crossed orderbook")
	}
	state.instrumentData.lastUpdateID = msg.SeqNum
	state.instrumentData.lastUpdateTime = utils.TimestampToMilli(msg.SnapshotL2.Timestamp)

	synced := false
	for !synced {
		if !ws.ReadMessage() {
			return fmt.Errorf("error reading message: %v", ws.Err)
		}
		depthData, ok := ws.Msg.Message.(fbinance.WSDepthData)
		if !ok {
			// Trade message
			context.Send(context.Self(), ws.Msg)
			continue
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

		state.instrumentData.lastUpdateID = depthData.FinalUpdateID - 1
		state.instrumentData.lastUpdateTime = uint64(ws.Msg.ClientTime.UnixNano() / 1000000)

		synced = true
	}

	state.instrumentData.orderBook = ob
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())

	go func(ws *fbinance.Websocket, pid *actor.PID) {
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
			if stat == models.OpenInterest && state.openInterestTicker == nil {
				openInterestTicker := time.NewTicker(10 * time.Second)
				state.openInterestTicker = openInterestTicker
				go func(pid *actor.PID) {
					for {
						select {
						case <-openInterestTicker.C:
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

	context.Respond(response)
	return nil
}

func (state *Listener) onWebsocketMessage(context actor.Context) error {
	msg := context.Message().(*xchanger.WebsocketMessage)
	switch msg.Message.(type) {
	case error:
		return fmt.Errorf("socket error: %v", msg)

	case fbinance.WSDepthData:
		depthData := msg.Message.(fbinance.WSDepthData)

		// change event time
		depthData.EventTime = uint64(msg.ClientTime.UnixNano()) / 1000000
		err := state.onDepthData(context, depthData)
		if err != nil {
			state.logger.Info("error processing depth data for "+depthData.Symbol,
				log.Error(err))
			return state.subscribeInstrument(context)
		}

	case fbinance.WSAggregatedTradeData:
		tradeData := msg.Message.(fbinance.WSAggregatedTradeData)
		ts := uint64(msg.ClientTime.UnixNano() / 1000000)
		if ts <= state.instrumentData.lastAggTradeTs {
			ts = state.instrumentData.lastAggTradeTs + 1
		}
		trade := &models.AggregatedTrade{
			Bid:         tradeData.MarketSell,
			Timestamp:   utils.MilliToTimestamp(ts),
			AggregateID: uint64(tradeData.AggregateTradeID),
			Trades: []models.Trade{{
				Price:    tradeData.Price,
				Quantity: tradeData.Quantity,
				ID:       uint64(tradeData.AggregateTradeID),
			}},
		}
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			Trades: []*models.AggregatedTrade{trade},
			SeqNum: state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1
		state.instrumentData.lastAggTradeTs = ts

	case fbinance.WSForceOrderData:
		orderData := msg.Message.(fbinance.WSForceOrderData)
		ts := uint64(msg.ClientTime.UnixNano() / 1000000)
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			Liquidation: &models.Liquidation{
				Bid:       orderData.Order.Side == fbinance.BUY_ORDER,
				Timestamp: utils.MilliToTimestamp(ts),
				OrderID:   rand.Uint64(),
				Price:     orderData.Order.OrigPrice,
				Quantity:  orderData.Order.LastFilledQuantity,
			},
			SeqNum: state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1

	case fbinance.WSMarkPriceData:
		mpData := msg.Message.(fbinance.WSMarkPriceData)
		refresh := &messages.MarketDataIncrementalRefresh{
			SeqNum: state.instrumentData.seqNum + 1,
		}
		refresh.Stats = append(refresh.Stats, &models.Stat{
			Timestamp: utils.MilliToTimestamp(mpData.NextFundingTime),
			StatType:  models.FundingRate,
			Value:     mpData.FundingRate,
		})
		context.Send(context.Parent(), refresh)
		state.instrumentData.seqNum += 1
	}

	return nil
}

func (state *Listener) onDepthData(context actor.Context, depthData fbinance.WSDepthData) error {

	symbol := depthData.Symbol
	instr := state.instrumentData

	// Skip depth that are younger than OB
	if depthData.FinalUpdateID <= instr.lastUpdateID {
		return nil
	}

	// Check depth continuity
	if instr.lastUpdateID+1 != depthData.PreviousUpdateID {
		return fmt.Errorf("got wrong sequence ID for %s: %d, %d",
			symbol, instr.lastUpdateID, depthData.PreviousUpdateID)
	}

	bids, asks, err := depthData.ToBidAsk()
	if err != nil {
		return fmt.Errorf("error converting depth data: %s ", err.Error())
	}

	obDelta := &models.OBL2Update{
		Levels:    []gmodels.OrderBookLevel{},
		Timestamp: utils.MilliToTimestamp(depthData.EventTime),
		Trade:     false,
	}

	for _, bid := range bids {
		obDelta.Levels = append(
			obDelta.Levels,
			bid,
		)
		instr.orderBook.UpdateOrderBookLevel(bid)
	}

	for _, ask := range asks {
		obDelta.Levels = append(
			obDelta.Levels,
			ask,
		)
		instr.orderBook.UpdateOrderBookLevel(ask)
	}

	if state.instrumentData.orderBook.Crossed() {
		state.logger.Warn("crossed orderbook", log.Error(errors.New("crossed")))
		return state.subscribeInstrument(context)
	}

	state.instrumentData.lastUpdateID = depthData.FinalUpdateID - 1
	state.instrumentData.lastUpdateTime = depthData.EventTime
	context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
		UpdateL2: obDelta,
		SeqNum:   state.instrumentData.seqNum + 1,
	})
	state.instrumentData.seqNum += 1
	return nil
}

func (state *Listener) checkSockets(context actor.Context) error {
	// TODO ping or hb ?
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
			state.logger.Warn("error fetching market statistics", log.Error(err))
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
			state.logger.Warn("error fetching market statistics", log.Error(errors.New(msg.RejectionReason.String())))
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
