package huobif

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
	"gitlab.com/alphaticks/xchanger/exchanges/huobif"
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
	lastUpdateID   uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	lastAggTradeTs uint64
}

type Listener struct {
	ws              *huobif.Websocket
	security        *models.Security
	dialerPool      *xchangerUtils.DialerPool
	instrumentData  *InstrumentData
	executorManager *actor.PID
	logger          *log.Logger
	lastPingTime    time.Time
	socketTicker    *time.Ticker
}

func NewListenerProducer(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Producer {
	return func() actor.Actor {
		return NewListener(security, dialerPool)
	}
}

func NewListener(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Actor {
	return &Listener{
		ws:              nil,
		security:        security,
		dialerPool:      dialerPool,
		instrumentData:  nil,
		executorManager: nil,
		logger:          nil,
		socketTicker:    nil,
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
	state.executorManager = actor.NewPID(context.ActorSystem().Address(), "exchange_executor_manager")
	state.lastPingTime = time.Now()

	state.instrumentData = &InstrumentData{
		orderBook:      nil,
		seqNum:         uint64(time.Now().UnixNano()),
		lastUpdateID:   0,
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

	ws := huobif.NewWebsocket()
	if err := ws.Connect(state.dialerPool.GetDialer()); err != nil {
		return fmt.Errorf("error connecting to huobi websocket: %v", err)
	}

	// If contract deliver in more than 7 days and less than 14 days, it's NW
	/*
		maturityDate := time.Unix(state.security.MaturityDate.Seconds, int64(state.security.MaturityDate.Nanos))


		var symbol string
		nextWeek := time.Now().AddDate(0, 0, 7)
		nextNextWeek := time.Now().AddDate(0, 0, 14)
		nextQuarter := time.Now().AddDate(0, 3, 0)
		nextNextQuarter := time.Now().AddDate(0, 6, 0)
		if maturityDate.Before(nextWeek) {
			symbol = fmt.Sprintf("%s_CW", state.security.Underlying.Symbol)
		} else if maturityDate.Before(nextNextWeek) {
			symbol = fmt.Sprintf("%s_NW", state.security.Underlying.Symbol)
		} else if maturityDate.Before(nextQuarter) {
			symbol = fmt.Sprintf("%s_CQ", state.security.Underlying.Symbol)
		} else if maturityDate.Before(nextNextQuarter) {
			symbol = fmt.Sprintf("%s_NQ", state.security.Underlying.Symbol)
		} else {
			return fmt.Errorf("unknown delivery date %v", maturityDate)
		}
	*/

	if err := ws.SubscribeMarketDepth(state.security.Symbol, huobif.WSOBLevel150, true); err != nil {
		return fmt.Errorf("error subscribing to orderbook for %s", state.security.Symbol)
	}

	var ob *gorderbook.OrderBookL2
	nTries := 0
	for nTries < 100 {
		if !ws.ReadMessage() {
			return fmt.Errorf("error reading message: %v", ws.Err)
		}

		switch ws.Msg.Message.(type) {
		case huobif.WSMarketDepthTick:
			res := ws.Msg.Message.(huobif.WSMarketDepthTick)
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
				return fmt.Errorf("crossed order book")
			}
			ts := uint64(ws.Msg.ClientTime.UnixNano()) / 1000000
			state.instrumentData.orderBook = ob
			state.instrumentData.lastUpdateTime = ts
			state.instrumentData.lastUpdateID = res.Version
			state.instrumentData.seqNum = uint64(time.Now().UnixNano())
			nTries = 100

		case huobif.WSError:
			err := fmt.Errorf("error getting orderbook: %s", ws.Msg.Message.(huobif.WSError).ErrMsg)
			return err
		}
		nTries += 1
	}

	if ob == nil {
		return fmt.Errorf("error getting orderbook")
	}

	if err := ws.SubscribeMarketTradeDetail(state.security.Symbol); err != nil {
		return fmt.Errorf("error subscribing to trades for %s", state.security.Symbol)
	}

	state.ws = ws

	go func(ws *huobif.Websocket, pid *actor.PID) {
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
		return fmt.Errorf("OB socket error: %v", msg)

	case huobif.WSMarketDepthTick:
		ts := uint64(msg.ClientTime.UnixNano() / 1000000)
		update := msg.Message.(huobif.WSMarketDepthTick)

		instr := state.instrumentData
		if update.Version <= instr.lastUpdateID {
			break
		}

		// This, implies that instr.lastUpdateID > seqNum
		// We want update.SeqNum > instr.lastUpdateID
		if instr.lastUpdateID+1 != update.Version {
			state.logger.Info("error processing ob update for "+update.Symbol, log.Error(fmt.Errorf("out of order sequence")))
			return state.subscribeInstrument(context)
		}

		obDelta := &models.OBL2Update{
			Levels:    nil,
			Timestamp: utils.MilliToTimestamp(ts),
			Trade:     false,
		}

		for _, bid := range update.Bids {
			level := gmodels.OrderBookLevel{
				Price:    bid.Price,
				Quantity: bid.Quantity,
				Bid:      true,
			}
			instr.orderBook.UpdateOrderBookLevel(level)
			obDelta.Levels = append(obDelta.Levels, level)
		}
		for _, ask := range update.Asks {
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
			return state.subscribeInstrument(context)
		}

		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			UpdateL2: obDelta,
			SeqNum:   state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1

		instr.lastUpdateTime = ts
		instr.lastUpdateID = update.Version
		//state.postSnapshot(context)

	case huobif.WSMarketTradeDetailTick:
		ts := uint64(msg.ClientTime.UnixNano() / 1000000)
		trades := msg.Message.(huobif.WSMarketTradeDetailTick)
		if len(trades.Data) == 0 {
			break
		}

		sort.Slice(trades.Data, func(i, j int) bool {
			return trades.Data[i].Timestamp < trades.Data[j].Timestamp
		})

		var aggTrade *models.AggregatedTrade
		var aggHelpR uint64 = 0
		for _, trade := range trades.Data {
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

	case huobif.WSSubscribeResponse:
		// pass

	case huobif.WSPing:
		msg := msg.Message.(huobif.WSPing)
		if err := state.ws.Pong(msg.Ping); err != nil {
			return fmt.Errorf("error sending pong to websocket")
		}

	case huobif.WSError:
		msg := msg.Message.(huobif.WSError)
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
	if time.Since(state.instrumentData.lastHBTime) > 2*time.Second {
		// Send an empty refresh
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			SeqNum: state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1
		state.instrumentData.lastHBTime = time.Now()
	}

	if state.ws.Err != nil || !state.ws.Connected {
		if state.ws.Err != nil {
			state.logger.Info("error on socket", log.Error(state.ws.Err))
		}
		if err := state.subscribeInstrument(context); err != nil {
			return fmt.Errorf("error subscribing to instrument: %v", err)
		}
	}

	return nil
}
