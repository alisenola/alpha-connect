package huobil

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
	"gitlab.com/alphaticks/xchanger/exchanges/huobil"
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
	lastUpdateID   uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	lastAggTradeTs uint64
}

type Listener struct {
	mdws               *huobil.MDWebsocket
	opws               *huobil.OPWebsocket
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
		security:   security,
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
	state.executor = actor.NewPID(context.ActorSystem().Address(), "executor/exchanges/"+constants.HUOBIL.Name+"_executor")
	state.lastPingTime = time.Now()

	state.instrumentData = &InstrumentData{
		orderBook:      nil,
		seqNum:         uint64(time.Now().UnixNano()),
		lastUpdateID:   0,
		lastUpdateTime: 0,
		lastHBTime:     time.Now(),
	}

	if err := state.subscribeMarketData(context); err != nil {
		return fmt.Errorf("error subscribing to order book: %v", err)
	}
	if err := state.subscribeOrderPush(context); err != nil {
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
	if state.mdws != nil {
		if err := state.mdws.Disconnect(); err != nil {
			state.logger.Info("error disconnecting socket", log.Error(err))
		}
	}
	if state.opws != nil {
		if err := state.opws.Disconnect(); err != nil {
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

func (state *Listener) subscribeMarketData(context actor.Context) error {
	if state.mdws != nil {
		_ = state.mdws.Disconnect()
	}

	mdws := huobil.NewMDWebsocket()
	if err := mdws.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to huobi websocket: %v", err)
	}

	if err := mdws.SubscribeMarketDepth(state.security.Symbol, huobil.WSOBLevel150, true); err != nil {
		return fmt.Errorf("error subscribing to orderbook for %s", state.security.Symbol)
	}

	var ob *gorderbook.OrderBookL2
	nTries := 0
	for nTries < 100 {
		if !mdws.ReadMessage() {
			return fmt.Errorf("error reading message: %v", mdws.Err)
		}

		switch mdws.Msg.Message.(type) {
		case huobil.WSMarketDepthTick:
			res := mdws.Msg.Message.(huobil.WSMarketDepthTick)
			if res.Event != "snapshot" {
				return fmt.Errorf("was expecting snapshot as first event, got %s", res.Event)
			}
			bids, asks := res.ToBidAsk()
			tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))
			lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot.Value))
			ob = gorderbook.NewOrderBookL2(
				tickPrecision,
				lotPrecision,
				10000)
			ob.Sync(bids, asks)
			if ob.Crossed() {
				return fmt.Errorf("cossed order book")
			}
			ts := uint64(mdws.Msg.ClientTime.UnixNano()) / 1000000
			state.instrumentData.orderBook = ob
			state.instrumentData.lastUpdateTime = ts
			state.instrumentData.lastUpdateID = res.Version
			state.instrumentData.seqNum = uint64(time.Now().UnixNano())
			nTries = 100

		case huobil.WSError:
			err := fmt.Errorf("error getting orderbook: %s", mdws.Msg.Message.(huobil.WSError).ErrMsg)
			return err
		}
		nTries += 1
	}

	if ob == nil {
		return fmt.Errorf("error getting orderbook")
	}

	if err := mdws.SubscribeMarketTradeDetail(state.security.Symbol); err != nil {
		return fmt.Errorf("error subscribing to trades for %s", state.security.Symbol)
	}

	state.mdws = mdws

	go func(ws *huobil.MDWebsocket, pid *actor.PID) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
		}
	}(mdws, context.Self())

	return nil
}

func (state *Listener) subscribeOrderPush(context actor.Context) error {
	if state.opws != nil {
		_ = state.opws.Disconnect()
	}

	opws := huobil.NewOPWebsocket()
	if err := opws.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to huobi websocket: %v", err)
	}

	if err := opws.SubscribeLiquidationOrders(state.security.Symbol); err != nil {
		return fmt.Errorf("error subscribing to liquidation orders for %s", state.security.Symbol)
	}

	if err := opws.SubscribeFundingRate(state.security.Symbol); err != nil {
		return fmt.Errorf("error subscribing to funding rate for %s", state.security.Symbol)
	}

	state.opws = opws

	go func(ws *huobil.OPWebsocket, pid *actor.PID) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
		}
	}(opws, context.Self())

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
	switch res := msg.Message.(type) {

	case error:
		return fmt.Errorf("OB socket error: %v", msg)

	case huobil.WSMarketDepthTick:
		ts := uint64(msg.ClientTime.UnixNano() / 1000000)

		instr := state.instrumentData
		if res.Version <= instr.lastUpdateID {
			break
		}

		// This, implies that instr.lastUpdateID > seqNum
		// We want update.SeqNum > instr.lastUpdateID
		if instr.lastUpdateID+1 != res.Version {
			state.logger.Info("error processing ob update for "+res.Symbol, log.Error(fmt.Errorf("out of order sequence")))
			return state.subscribeMarketData(context)
		}

		obDelta := &models.OBL2Update{
			Levels:    nil,
			Timestamp: utils.MilliToTimestamp(ts),
			Trade:     false,
		}

		for _, bid := range res.Bids {
			level := gmodels.OrderBookLevel{
				Price:    bid.Price,
				Quantity: bid.Quantity,
				Bid:      true,
			}
			instr.orderBook.UpdateOrderBookLevel(level)
			obDelta.Levels = append(obDelta.Levels, level)
		}
		for _, ask := range res.Asks {
			level := gmodels.OrderBookLevel{
				Price:    ask.Price,
				Quantity: ask.Quantity,
				Bid:      false,
			}
			instr.orderBook.UpdateOrderBookLevel(level)
			obDelta.Levels = append(obDelta.Levels, level)
		}

		if state.instrumentData.orderBook.Crossed() {
			state.logger.Info("crossed orderbook", log.Error(errors.New("crossed")))
			return state.subscribeMarketData(context)
		}

		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			UpdateL2: obDelta,
			SeqNum:   state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1

		instr.lastUpdateTime = ts
		instr.lastUpdateID = res.Version
		//state.postSnapshot(context)

	case huobil.WSMarketTradeDetailTick:
		ts := uint64(msg.ClientTime.UnixNano() / 1000000)
		if len(res.Data) == 0 {
			break
		}

		sort.Slice(res.Data, func(i, j int) bool {
			return res.Data[i].Timestamp < res.Data[j].Timestamp
		})

		var aggTrade *models.AggregatedTrade
		var aggHelpR uint64 = 0
		for _, trade := range res.Data {
			aggHelp := trade.Timestamp * 10
			// do that so new agg trade if side changes
			if trade.Direction == "sell" {
				aggHelp += 1
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
					AggregateID: trade.ID,
					Trades:      nil,
				}
				aggHelpR = aggHelp
			}

			trd := models.Trade{
				Price:    trade.Price,
				Quantity: trade.Amount,
				ID:       trade.ID,
			}

			aggTrade.Trades = append(aggTrade.Trades, trd)
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

	case []huobil.WSFundingRate:
		for _, r := range res {
			fmt.Println("FUNDING", r.FundingRate, r.SettlementTime, r.FundingTime)
			refresh := &messages.MarketDataIncrementalRefresh{
				SeqNum: state.instrumentData.seqNum + 1,
			}
			refresh.Stats = append(refresh.Stats, &models.Stat{
				Timestamp: utils.MilliToTimestamp(r.SettlementTime),
				StatType:  models.FundingRate,
				Value:     r.FundingRate,
			})
			context.Send(context.Parent(), refresh)
			state.instrumentData.seqNum += 1
		}

	case []huobil.WSLiquidationOrder:
		for _, r := range res {
			refresh := &messages.MarketDataIncrementalRefresh{
				// If the position was a short, the direction == 'sell', therefore
				// the liquidation is a buy
				Liquidation: &models.Liquidation{
					Bid:       r.Direction == "sell",
					Timestamp: utils.MilliToTimestamp(uint64(r.CreatedAt)),
					OrderID:   uint64(r.CreatedAt),
					Price:     r.Price,
					Quantity:  r.Amount,
				},
				SeqNum: state.instrumentData.seqNum + 1,
			}
			context.Send(context.Parent(), refresh)
			state.instrumentData.seqNum += 1
		}

	case huobil.MDWSSubscribeResponse, huobil.OPWSSubscribeResponse:
		// pass

	case huobil.MDWSPing:
		ping := msg.Message.(huobil.MDWSPing)
		if err := state.mdws.Pong(ping.Ping); err != nil {
			return fmt.Errorf("error sending pong to websocket")
		}

	case huobil.OPWSPing:
		ping := msg.Message.(huobil.OPWSPing)
		if err := state.opws.Pong(ping.Ts); err != nil {
			return fmt.Errorf("error sending pong to websocket")
		}

	case huobil.WSError:
		msg := msg.Message.(huobil.WSError)
		state.logger.Info("got WSError message",
			log.String("message", msg.ErrMsg),
			log.String("code", msg.ErrCode))

	default:
		state.logger.Info("received unknown message",
			log.String("message_type",
				reflect.TypeOf(msg.Message).String()))
	}

	return nil
}

func (state *Listener) checkSockets(context actor.Context) error {
	// If haven't sent anything for 2 seconds, send heartbeat
	if time.Now().Sub(state.instrumentData.lastHBTime) > 2*time.Second {
		// Send an empty refresh
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			SeqNum: state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1
		state.instrumentData.lastHBTime = time.Now()
	}

	if state.mdws.Err != nil || !state.mdws.Connected {
		if state.mdws.Err != nil {
			state.logger.Info("error on socket", log.Error(state.mdws.Err))
		}
		if err := state.subscribeMarketData(context); err != nil {
			return fmt.Errorf("error subscribing to instrument: %v", err)
		}
	}
	if state.opws.Err != nil || !state.opws.Connected {
		if state.opws.Err != nil {
			state.logger.Info("error on socket", log.Error(state.opws.Err))
		}
		if err := state.subscribeOrderPush(context); err != nil {
			return fmt.Errorf("error subscribing to instrument: %v", err)
		}
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
			Statistics: []models.StatType{models.OpenInterest},
		}, 2*time.Second)

	context.AwaitFuture(fut, func(res interface{}, err error) {
		if err != nil {
			if err == actor.ErrTimeout {
				oidLock.Lock()
				oid = time.Duration(float64(oid) * 1.01)
				if state.openInterestTicker != nil {
					state.openInterestTicker.Reset(oid)
				}
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
				if state.openInterestTicker != nil {
					state.openInterestTicker.Reset(oid)
				}
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
		if state.openInterestTicker != nil {
			state.openInterestTicker.Reset(oid)
		}
		oidLock.Unlock()
	})

	return nil
}
