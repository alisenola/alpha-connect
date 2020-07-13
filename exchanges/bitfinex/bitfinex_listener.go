package bitfinex

import (
	"container/list"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/exchanges/bitfinex"
	"math"
	"reflect"
	"time"
)

type readSocket struct{}
type postAggTrade struct{}

type OBL2Request struct {
	requester *actor.PID
	requestID int64
}

type InstrumentData struct {
	tickPrecision  uint64
	lotPrecision   uint64
	orderBook      *gorderbook.OrderBookL2
	obDelta        *models.OBL2Update
	nCrossed       int
	seqNum         uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	lastSequence   uint64
	aggTrade       *models.AggregatedTrade
	lastAggTradeTs uint64
}

// OBType: OBL3
// OBL3 timestamps: per connection
// Status: not ready problem with timestamps

type Listener struct {
	obWs           *bitfinex.Websocket
	tradeWs        *bitfinex.Websocket
	wsChan         chan *bitfinex.WebsocketMessage
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

// Limit of 30 subscription
func NewListener(security *models.Security) actor.Actor {
	return &Listener{
		obWs:           nil,
		tradeWs:        nil,
		wsChan:         nil,
		security:       security,
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

	state.wsChan = make(chan *bitfinex.WebsocketMessage, 10000)
	state.stashedTrades = list.New()

	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot))

	state.instrumentData = &InstrumentData{
		tickPrecision:  tickPrecision,
		lotPrecision:   lotPrecision,
		orderBook:      nil,
		seqNum:         uint64(time.Now().UnixNano()),
		lastUpdateTime: 0,
		lastHBTime:     time.Now(),
		aggTrade:       nil,
		lastAggTradeTs: 0,
		obDelta:        nil,
		nCrossed:       0,
	}

	if err := state.subscribeOrderBook(context); err != nil {
		return fmt.Errorf("error subscribing to order book: %v", err)
	}
	if err := state.subscribeTrades(context); err != nil {
		return fmt.Errorf("error subscribing to trades: %v", err)
	}
	context.Send(context.Self(), &readSocket{})

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

	ws := bitfinex.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to bitfinex websocket: %v", err)
	}

	err := ws.SubscribeDepth(
		state.security.Symbol,
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
		return fmt.Errorf("was expecting WSSubscribeDepthResponse, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	snapshot, ok := ws.Msg.Message.(bitfinex.WSSpotDepthSnapshot)
	if !ok {
		return fmt.Errorf("was expecting WSSpotDepthL3Snapshot, got %s", reflect.TypeOf(ws.Msg.Message).String())
	}

	state.instrumentData.lastSequence = snapshot.Sequence

	bids, asks := snapshot.ToBidAsk(state.instrumentData.tickPrecision, state.instrumentData.lotPrecision)
	bestAsk := float64(asks[0].Price) / float64(state.instrumentData.tickPrecision)
	//Allow a 10% price variation
	depth := int(((bestAsk * 1.1) - bestAsk) * float64(state.instrumentData.tickPrecision))
	if depth > 10000 {
		depth = 10000
	}

	ob := gorderbook.NewOrderBookL2(
		state.instrumentData.tickPrecision,
		state.instrumentData.lotPrecision,
		depth,
	)

	ts := uint64(ws.Msg.Time.UnixNano() / 1000000)

	ob.RawSync(bids, asks)

	state.instrumentData.orderBook = ob
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())
	state.instrumentData.lastUpdateTime = ts
	state.instrumentData.obDelta = nil
	state.instrumentData.nCrossed = 0

	state.obWs = ws

	go func(ws *bitfinex.Websocket) {
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

	ws := bitfinex.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to bitfinex websocket: %v", err)
	}
	err := ws.SubscribeTrades(state.security.Symbol)
	if err != nil {
		return err
	}

	state.tradeWs = ws

	go func(ws *bitfinex.Websocket) {
		for ws.ReadMessage() {
			state.wsChan <- ws.Msg
		}
	}(ws)

	return nil
}

func (state *Listener) OnMarketDataRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketDataRequest)
	snapshot := &models.OBL2Snapshot{
		Bids:      state.instrumentData.orderBook.GetBids(0),
		Asks:      state.instrumentData.orderBook.GetAsks(0),
		Timestamp: utils.MilliToTimestamp(state.instrumentData.lastUpdateTime),
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

func (state *Listener) readSocket(context actor.Context) error {
	select {
	case msg := <-state.wsChan:
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
			ts := uint64(msg.Time.UnixNano() / 1000000)
			if state.instrumentData.obDelta == nil {
				state.instrumentData.obDelta = &models.OBL2Update{
					Levels:    []gorderbook.OrderBookLevel{},
					Timestamp: utils.MilliToTimestamp(ts),
					Trade:     false,
				}
			}

			state.instrumentData.orderBook.UpdateOrderBookLevel(level)
			state.instrumentData.obDelta.Levels = append(state.instrumentData.obDelta.Levels, level)

			if state.instrumentData.orderBook.Crossed() {
				state.instrumentData.nCrossed += 1
				fmt.Println(state.instrumentData.nCrossed)
				if state.instrumentData.nCrossed > 5 {
					// Stop the socket, we will restart instrument at the end
					if err := state.obWs.Disconnect(); err != nil {
						state.logger.Info("error disconnecting from socket", log.Error(err))
					}
					break
				}
			} else {
				context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
					UpdateL2: state.instrumentData.obDelta,
					SeqNum:   state.instrumentData.seqNum + 1,
				})
				state.instrumentData.seqNum += 1
				state.instrumentData.lastUpdateTime = utils.TimestampToMilli(state.instrumentData.obDelta.Timestamp)
				state.instrumentData.obDelta = nil
			}

		case bitfinex.WSSpotTrade:
			tradeData := msg.Message.(bitfinex.WSSpotTrade)
			ts := uint64(msg.Time.UnixNano() / 1000000)

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
					models.Trade{
						Price:    tradeData.Price,
						Quantity: -tradeData.Amount,
						ID:       tradeData.ID,
					})
			} else {
				state.instrumentData.aggTrade.Trades = append(
					state.instrumentData.aggTrade.Trades,
					models.Trade{
						Price:    tradeData.Price,
						Quantity: tradeData.Amount,
						ID:       tradeData.ID,
					})
			}

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

func (state *Listener) checkSockets(context actor.Context) error {
	// No need to ping HB mechanism already
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
