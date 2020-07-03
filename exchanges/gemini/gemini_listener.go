package gemini

import (
	"container/list"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/exchanges/gemini"
	"math"
	"reflect"
	"time"
)

type readSocket struct{}
type postAggTrade struct{}

type InstrumentData struct {
	orderBook      *gorderbook.OrderBookL2
	seqNum         uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	aggTrade       *models.AggregatedTrade
	lastAggTradeTs uint64
}

type Listener struct {
	ws             *gemini.Websocket
	wsChan         chan *gemini.WebsocketMessage
	security       *models.Security
	instrumentData *InstrumentData
	logger         *log.Logger
	stashedTrades  *list.List
}

func NewListenerProducer(security *models.Security) actor.Producer {
	return func() actor.Actor {
		return NewListener(security)
	}
}

func NewListener(security *models.Security) actor.Actor {
	return &Listener{
		ws:             nil,
		wsChan:         nil,
		security:       security,
		instrumentData: nil,
		logger:         nil,
	}
}

func (state *Listener) Receive(context actor.Context) {
	fmt.Println("LISTENER MESSAGE", reflect.TypeOf(context.Message()).String())
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

	state.wsChan = make(chan *gemini.WebsocketMessage, 10000)
	state.stashedTrades = list.New()

	state.instrumentData = &InstrumentData{
		orderBook:      nil,
		seqNum:         uint64(time.Now().UnixNano()),
		lastUpdateTime: 0,
		lastHBTime:     time.Now(),
		aggTrade:       nil,
		lastAggTradeTs: 0,
	}

	if err := state.subscribeInstrument(context); err != nil {
		return fmt.Errorf("error subscribing to order book: %v", err)
	}

	context.Send(context.Self(), &readSocket{})

	return nil
}

func (state *Listener) Clean(context actor.Context) error {
	if state.ws != nil {
		if err := state.ws.Disconnect(); err != nil {
			state.logger.Info("error disconnecting socket", log.Error(err))
		}
	}

	return nil
}

func (state *Listener) subscribeInstrument(context actor.Context) error {
	if state.ws != nil {
		_ = state.ws.Disconnect()
	}

	ws := gemini.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to gemini websocket: %v", err)
	}

	if err := ws.Subscribe(state.security.Symbol, gemini.WSL2Feed); err != nil {
		return fmt.Errorf("error subscribing to l2 feed")
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	update, ok := ws.Msg.Message.(gemini.WSL2Updates)
	if !ok {
		return fmt.Errorf("was expecting WSL2Updates, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	bids, asks := update.ToBidAsk()

	// TODO
	/*
		bestAsk := float64(asks[0]) / float64(msg.Snapshot.Instrument.TickPrecision)
		depth := int(((bestAsk * 1.1) - bestAsk) * float64(msg.Snapshot.Instrument.TickPrecision))

		if depth > 10000 {
			depth = 10000
		}
	*/
	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot))
	ob := gorderbook.NewOrderBookL2(
		tickPrecision,
		lotPrecision,
		10000)
	ob.Sync(bids, asks)

	ts := uint64(ws.Msg.Time.UnixNano() / 1000000)
	state.instrumentData.orderBook = ob
	state.instrumentData.lastUpdateTime = ts
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())

	state.ws = ws

	go func(ws *gemini.Websocket) {
		for ws.ReadMessage() {
			state.wsChan <- ws.Msg
		}
	}(ws)

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

		case gemini.WSL2Updates:
			updates := msg.Message.(gemini.WSL2Updates)

			instr := state.instrumentData

			bids, asks := updates.ToBidAsk()
			ts := uint64(msg.Time.UnixNano() / 1000000)

			obDelta := &models.OBL2Update{
				Levels:    nil,
				Timestamp: utils.MilliToTimestamp(ts),
				Trade:     false,
			}

			for _, bid := range bids {
				instr.orderBook.UpdateOrderBookLevel(bid)
				obDelta.Levels = append(obDelta.Levels, bid)
			}
			for _, ask := range asks {
				instr.orderBook.UpdateOrderBookLevel(ask)
				obDelta.Levels = append(obDelta.Levels, ask)
			}

			if state.instrumentData.orderBook.Crossed() {
				state.logger.Info("crossed order book")
				// Stop the socket, we will restart instrument at the end
				if err := state.ws.Disconnect(); err != nil {
					state.logger.Info("error disconnecting from socket", log.Error(err))
				}
				break
			}

			context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
				UpdateL2: obDelta,
				SeqNum:   state.instrumentData.seqNum + 1,
			})
			state.instrumentData.seqNum += 1
			instr.lastUpdateTime = uint64(msg.Time.UnixNano() / 1000000)

		case gemini.WSTrade:
			trade := msg.Message.(gemini.WSTrade)
			ts := uint64(msg.Time.UnixNano()) / 1000000

			aggID := trade.Timestamp * 10
			if trade.Side == "sell" {
				aggID += 1
			}
			if state.instrumentData.aggTrade == nil || state.instrumentData.aggTrade.AggregateID != aggID {
				// ensure increasing timestamp
				if ts <= state.instrumentData.lastAggTradeTs {
					ts = state.instrumentData.lastAggTradeTs + 1
				}
				aggTrade := &models.AggregatedTrade{
					Bid:         trade.Side == "sell",
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
			state.instrumentData.aggTrade.Trades = append(
				state.instrumentData.aggTrade.Trades,
				models.Trade{
					Price:    trade.Price,
					Quantity: trade.Quantity,
					ID:       trade.EventID,
				})
		}

		state.postHeartBeat(context)
		if err := state.checkSockets(context); err != nil {
			return fmt.Errorf("error checking sockets: %v", err)
		}
		context.Send(context.Self(), &readSocket{})
		return nil

	case <-time.After(1 * time.Second):
		state.postHeartBeat(context)
		if err := state.checkSockets(context); err != nil {
			return fmt.Errorf("error checking sockets: %v", err)
		}
		context.Send(context.Self(), &readSocket{})
		return nil
	}
}

func (state *Listener) checkSockets(context actor.Context) error {
	// Don't re-subscribe as ping, we have heartbeats
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
