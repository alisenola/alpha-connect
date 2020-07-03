package coinbasepro

import (
	"encoding/binary"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/coinbasepro"
	"math"
	"reflect"
	"time"
)

type readSocket struct{}

type OBL2Request struct {
	requester *actor.PID
	requestID int64
}

type InstrumentData struct {
	orderBook      *gorderbook.OrderBookL3
	seqNum         uint64
	lastSequence   uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	aggTrade       *models.AggregatedTrade
	lastAggTradeTs uint64
}

// OBType: OBL3
// OBL3 timestamps: Per server, full order book ws events can arrive unordered even though sequence are ordered
// so we use local time
// Trades can be inferred with delta updates
// Status: Not ready, problem with unordered ob events

type Listener struct {
	ws                  *coinbasepro.Websocket
	wsChan              chan *coinbasepro.WebsocketMessage
	security            *models.Security
	instrumentData      *InstrumentData
	coinbaseproExecutor *actor.PID
	logger              *log.Logger
}

func NewListenerProducer(security *models.Security) actor.Producer {
	return func() actor.Actor {
		return NewListener(security)
	}
}

func NewListener(security *models.Security) actor.Actor {
	return &Listener{
		ws:                  nil,
		wsChan:              nil,
		security:            security,
		instrumentData:      nil,
		coinbaseproExecutor: nil,
		logger:              nil,
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

	state.coinbaseproExecutor = actor.NewLocalPID("executor/" + constants.COINBASEPRO.Name + "_executor")
	state.wsChan = make(chan *coinbasepro.WebsocketMessage, 10000)

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

	ws := coinbasepro.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to the websocket: %v", err)
	}

	err := ws.SubscribeFullChannel([]string{state.security.Symbol})
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	fut := context.RequestFuture(
		state.coinbaseproExecutor,
		&messages.MarketDataRequest{
			RequestID: uint64(time.Now().UnixNano()),
			Subscribe: false,
			Instrument: &models.Instrument{
				SecurityID: &types.UInt64Value{Value: state.security.SecurityID},
				Exchange:   state.security.Exchange,
				Symbol:     &types.StringValue{Value: state.security.Symbol},
			},
			Aggregation: models.L3,
		},
		5*time.Second)

	res, err := fut.Result()
	if err != nil {
		return fmt.Errorf("error getting OBL3")
	}
	msg, ok := res.(*messages.MarketDataResponse)
	if !ok {
		return fmt.Errorf("was expecting MarketDataSnapshot, got %s", reflect.TypeOf(msg).String())
	}
	if !msg.Success {
		return fmt.Errorf("error fetching snapshot: %s", msg.RejectionReason.String())
	}
	if msg.SnapshotL3 == nil {
		return fmt.Errorf("market data snapshot has no OBL3")
	}

	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot))
	ob := gorderbook.NewOrderBookL3(
		tickPrecision,
		lotPrecision,
		10000)

	ob.Sync(msg.SnapshotL3.Bids, msg.SnapshotL3.Asks)

	state.instrumentData.orderBook = ob
	state.instrumentData.lastUpdateTime = uint64(time.Now().UnixNano() / 1000000)
	state.instrumentData.seqNum = uint64(time.Now().UnixNano())
	state.instrumentData.lastSequence = msg.SeqNum
	state.ws = ws

	// Fetch messages until sync
	synced := false
	for !synced {
		if !ws.ReadMessage() {
			return fmt.Errorf("error reading message: %v", ws.Err)
		}
		switch ws.Msg.Message.(type) {
		case error:
			return fmt.Errorf("socket error: %v", msg)

		case coinbasepro.WSOpenOrder:
			order := ws.Msg.Message.(coinbasepro.WSOpenOrder)
			if order.Sequence < state.instrumentData.lastSequence {
				continue
			} else if order.Sequence == state.instrumentData.lastSequence {
				synced = true
			} else {
				return fmt.Errorf("out of order sequence %d:%d", order.Sequence, state.instrumentData.lastSequence)
			}
		case coinbasepro.WSChangeOrder:
			order := ws.Msg.Message.(coinbasepro.WSChangeOrder)
			if order.Sequence < state.instrumentData.lastSequence {
				continue
			} else if order.Sequence == state.instrumentData.lastSequence {
				synced = true
			} else {
				return fmt.Errorf("out of order sequence %d:%d", order.Sequence, state.instrumentData.lastSequence)
			}

		case coinbasepro.WSMatchOrder:
			order := ws.Msg.Message.(coinbasepro.WSMatchOrder)
			if order.Sequence < state.instrumentData.lastSequence {
				continue
			} else if order.Sequence == state.instrumentData.lastSequence {
				synced = true
			} else {
				return fmt.Errorf("out of order sequence %d:%d", order.Sequence, state.instrumentData.lastSequence)
			}

		case coinbasepro.WSDoneOrder:
			order := ws.Msg.Message.(coinbasepro.WSDoneOrder)
			if order.Sequence < state.instrumentData.lastSequence {
				continue
			} else if order.Sequence == state.instrumentData.lastSequence {
				synced = true
			} else {
				return fmt.Errorf("out of order sequence %d:%d", order.Sequence, state.instrumentData.lastSequence)
			}

		case coinbasepro.WSReceivedOrder:
			order := ws.Msg.Message.(coinbasepro.WSReceivedOrder)
			if order.Sequence < state.instrumentData.lastSequence {
				continue
			} else if order.Sequence == state.instrumentData.lastSequence {
				synced = true
			} else {
				return fmt.Errorf("out of order sequence %d:%d", order.Sequence, state.instrumentData.lastSequence)
			}

		case coinbasepro.WSSubscriptions:
			break
		}
	}

	go func(ws *coinbasepro.Websocket) {
		for ws.ReadMessage() {
			state.wsChan <- ws.Msg
		}
	}(state.ws)

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
	// TODO L3 ?

	context.Respond(response)
	return nil
}

func (state *Listener) readSocket(context actor.Context) error {
	select {
	case msg := <-state.wsChan:
		switch msg.Message.(type) {

		case error:
			return fmt.Errorf("socket error: %v", msg)

		case coinbasepro.WSOpenOrder:
			order := msg.Message.(coinbasepro.WSOpenOrder)
			// Replace time with local one

			order.Time = msg.Time
			if err := state.onOpenOrder(order, context); err != nil {
				state.logger.Info("error processing OpenOrder", log.Error(err))
				// Stop the socket, we will restart instrument at the end
				if err := state.ws.Disconnect(); err != nil {
					state.logger.Info("error disconnecting from socket", log.Error(err))
				}
			}

		case coinbasepro.WSChangeOrder:
			order := msg.Message.(coinbasepro.WSChangeOrder)

			order.Time = msg.Time
			err := state.onChangeOrder(order, context)
			if err != nil {
				state.logger.Info("error processing ChangeOrder", log.Error(err))
				// Stop the socket, we will restart instrument at the end
				if err := state.ws.Disconnect(); err != nil {
					state.logger.Info("error disconnecting from socket", log.Error(err))
				}
			}

		case coinbasepro.WSMatchOrder:
			order := msg.Message.(coinbasepro.WSMatchOrder)
			order.Time = msg.Time
			err := state.onMatchOrder(order, context)
			if err != nil {
				state.logger.Info("error processing MatchOrder", log.Error(err))
				// Stop the socket, we will restart instrument at the end
				if err := state.ws.Disconnect(); err != nil {
					state.logger.Info("error disconnecting from socket", log.Error(err))
				}
			}

		case coinbasepro.WSDoneOrder:
			order := msg.Message.(coinbasepro.WSDoneOrder)
			order.Time = msg.Time
			err := state.onDoneOrder(order, context)
			if err != nil {
				state.logger.Info("error processing DoneOrder", log.Error(err))
				// Stop the socket, we will restart instrument at the end
				if err := state.ws.Disconnect(); err != nil {
					state.logger.Info("error disconnecting from socket", log.Error(err))
				}
			}

		case coinbasepro.WSReceivedOrder:
			order := msg.Message.(coinbasepro.WSReceivedOrder)
			order.Time = msg.Time
			err := state.onReceivedOrder(order, context)
			if err != nil {
				state.logger.Info("error processing ReceivedOrder", log.Error(err))
				// Stop the socket, we will restart instrument at the end
				if err := state.ws.Disconnect(); err != nil {
					state.logger.Info("error disconnecting from socket", log.Error(err))
				}
			}

		case coinbasepro.WSSubscriptions:
			break

		case coinbasepro.WSError:
			// TODO handle error, skip unsubscribe error
			fmt.Println("ERROR MESSAGE:", msg)
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
	// TODO ping or HB ?
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

func (state *Listener) onOpenOrder(order coinbasepro.WSOpenOrder, context actor.Context) error {
	state.postAggTrade(context)
	symbol := order.ProductID
	// Check sequence consistency
	if order.Sequence <= state.instrumentData.lastSequence {
		return nil
	}
	if order.Sequence != state.instrumentData.lastSequence+1 {
		return fmt.Errorf("got inconsistent sequence for %s: %d, %d",
			symbol,
			order.Sequence,
			state.instrumentData.lastSequence+1)
	}

	orderID, err := uuid.FromString(order.OrderID)
	if err != nil {
		return fmt.Errorf("error parsing order uuid: %v", err)
	}

	obOrder := gorderbook.Order{
		Price:    order.Price,
		Quantity: order.RemainingSize,
		Bid:      order.Side == "buy",
		ID:       binary.LittleEndian.Uint64(orderID.Bytes()[0:8]),
	}

	state.instrumentData.orderBook.AddOrder(obOrder)

	var quantity float64
	if order.Side == "buy" {
		quantity = state.instrumentData.orderBook.GetBid(order.Price)
	} else {
		quantity = state.instrumentData.orderBook.GetAsk(order.Price)
	}

	levelDelta := gorderbook.OrderBookLevel{
		Price:    order.Price,
		Quantity: quantity,
		Bid:      order.Side == "buy",
	}

	ts := uint64(order.Time.UnixNano()) / 1000000

	if state.instrumentData.orderBook.Crossed() {
		return fmt.Errorf("crossed order book")
	}

	// SEND DELTA //
	obDelta := &models.OBL2Update{
		Levels:    []gorderbook.OrderBookLevel{levelDelta},
		Timestamp: utils.MilliToTimestamp(ts),
		Trade:     false,
	}

	context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
		UpdateL2: obDelta,
		SeqNum:   state.instrumentData.seqNum + 1,
	})
	state.instrumentData.seqNum += 1
	state.instrumentData.lastUpdateTime = ts
	/////////////////

	state.instrumentData.lastSequence = order.Sequence
	return nil
}

func (state *Listener) onChangeOrder(order coinbasepro.WSChangeOrder, context actor.Context) error {
	state.postAggTrade(context)
	instr := state.instrumentData
	// check sequence consistency
	if order.Sequence <= instr.lastSequence {
		return nil
	}
	if order.Sequence != instr.lastSequence+1 {
		return fmt.Errorf("change inconsistent sequence: %d, %d",
			order.Sequence,
			state.instrumentData.lastSequence+1)
	}

	orderUUID, err := uuid.FromString(order.OrderID)
	if err != nil {
		return fmt.Errorf("error parsing order uuid: %v", err)
	}

	orderID := binary.LittleEndian.Uint64(orderUUID.Bytes()[0:8])

	if instr.orderBook.HasOrder(orderID) {
		obOrder := gorderbook.Order{
			Price:    order.Price,
			Quantity: order.NewSize,
			Bid:      order.Side == "buy",
			ID:       orderID,
		}
		// TODO can an order change price ? here we assume not
		//lastRawOrder := state.instruments[order.ProductID].orderBook.GetRawOrder(obOrder.ID)

		instr.orderBook.UpdateOrder(obOrder)

		newOrder := instr.orderBook.GetOrder(obOrder.ID)

		var quantity float64
		if order.Side == "buy" {
			quantity = instr.orderBook.GetBid(newOrder.Price)
		} else {
			quantity = instr.orderBook.GetAsk(newOrder.Price)
		}

		levelDelta := gorderbook.OrderBookLevel{
			Price:    newOrder.Price,
			Quantity: quantity,
			Bid:      order.Side == "buy",
		}

		ts := uint64(order.Time.UnixNano()) / 1000000

		if state.instrumentData.orderBook.Crossed() {
			return fmt.Errorf("crossed order book")
		}
		// SEND DELTA //
		obDelta := &models.OBL2Update{
			Levels:    []gorderbook.OrderBookLevel{levelDelta},
			Timestamp: utils.MilliToTimestamp(ts),
			Trade:     false,
		}
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			UpdateL2: obDelta,
			SeqNum:   state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1

		instr.lastUpdateTime = ts
		//////////////

		// Send snapshot on sync

	}

	instr.lastSequence = order.Sequence

	return nil
}

func (state *Listener) onMatchOrder(order coinbasepro.WSMatchOrder, context actor.Context) error {
	instr := state.instrumentData
	// check sequence consistency
	if order.Sequence <= instr.lastSequence {
		return nil
	}
	if order.Sequence != instr.lastSequence+1 {
		return fmt.Errorf("match inconsistent sequence: %d, %d",
			order.Sequence,
			instr.lastSequence+1)
	}

	// MATCH PROCESSING
	orderUUID, err := uuid.FromString(order.MakerOrderID)
	if err != nil {
		return fmt.Errorf("error parsing order uuid: %v", err)
	}
	orderID := binary.LittleEndian.Uint64(orderUUID.Bytes()[0:8])

	rawOrder := instr.orderBook.GetRawOrder(orderID)
	rawMatchSize := uint64(math.Round(order.Size * float64(instr.orderBook.LotPrecision)))
	rawOrder.Quantity -= rawMatchSize

	instr.orderBook.UpdateRawOrder(rawOrder)

	// We want the quantity at the price of the maker !! not
	// at the price of the taker
	price := float64(rawOrder.Price) / float64(instr.orderBook.TickPrecision)
	var quantity float64
	if order.Side == "buy" {
		quantity = instr.orderBook.GetBid(price)
	} else {
		quantity = instr.orderBook.GetAsk(price)
	}

	levelDelta := gorderbook.OrderBookLevel{
		Price:    order.Price,
		Quantity: quantity,
		Bid:      order.Side == "buy",
	}

	ts := uint64(order.Time.UnixNano()) / 1000000

	if state.instrumentData.orderBook.Crossed() {
		return fmt.Errorf("crossed order book")
	}
	// SEND OBDELTA //
	obDelta := &models.OBL2Update{
		Levels:    []gorderbook.OrderBookLevel{levelDelta},
		Timestamp: utils.MilliToTimestamp(ts),
		Trade:     true,
	}

	context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
		UpdateL2: obDelta,
		SeqNum:   state.instrumentData.seqNum + 1,
	})
	state.instrumentData.seqNum += 1
	instr.lastUpdateTime = ts
	///////////////

	// HANDLE TRADE ///
	takerOrderUUID, err := uuid.FromString(order.TakerOrderID)
	if err != nil {
		return fmt.Errorf("error parsing order uuid: %v", err)
	}

	takerOrderID := binary.LittleEndian.Uint64(takerOrderUUID.Bytes()[0:8])
	aggID := takerOrderID

	if state.instrumentData.aggTrade == nil || state.instrumentData.aggTrade.AggregateID != aggID {
		state.postAggTrade(context)
		// ensure increasing timestamp
		if ts <= state.instrumentData.lastAggTradeTs {
			ts = state.instrumentData.lastAggTradeTs + 1
		}
		state.instrumentData.aggTrade = &models.AggregatedTrade{
			Bid:         levelDelta.Bid,
			Timestamp:   utils.MilliToTimestamp(ts),
			AggregateID: aggID,
			Trades:      nil,
		}
		state.instrumentData.lastAggTradeTs = ts
	}

	state.instrumentData.aggTrade.Trades = append(
		state.instrumentData.aggTrade.Trades,
		models.Trade{
			Price:    levelDelta.Price,
			Quantity: order.Size,
			ID:       order.Sequence,
		})
	/////////////////

	instr.lastSequence = order.Sequence
	return nil
}

func (state *Listener) onDoneOrder(order coinbasepro.WSDoneOrder, context actor.Context) error {
	instr := state.instrumentData
	// check sequence consistency
	if order.Sequence <= instr.lastSequence {
		return nil
	}
	if order.Sequence != instr.lastSequence+1 {
		return fmt.Errorf("done inconsistent sequence: %d, %d",
			order.Sequence,
			instr.lastSequence+1)
	}

	// DONE ORDER PROCESSING
	if order.Reason == "canceled" {
		state.postAggTrade(context)
		orderUUID, err := uuid.FromString(order.OrderID)
		if err != nil {
			return fmt.Errorf("error parsing order uuid: %v", err)
		}
		orderID := binary.LittleEndian.Uint64(orderUUID.Bytes()[0:8])

		// perfectly normal if OB hasn't the order, if a match
		// ate the order, it will already have been deleted
		if instr.orderBook.HasOrder(orderID) {
			obOrder := instr.orderBook.GetOrder(orderID)
			instr.orderBook.DeleteOrder(orderID)

			var quantity float64
			if order.Side == "buy" {
				quantity = instr.orderBook.GetBid(obOrder.Price)
			} else {
				quantity = instr.orderBook.GetAsk(obOrder.Price)
			}

			levelDelta := gorderbook.OrderBookLevel{
				Price:    obOrder.Price,
				Quantity: quantity,
				Bid:      obOrder.Bid,
			}

			if state.instrumentData.orderBook.Crossed() {
				return fmt.Errorf("crossed order book")
			}
			// SEND DELTA //
			ts := uint64(order.Time.UnixNano()) / 1000000

			obDelta := &models.OBL2Update{
				Levels:    []gorderbook.OrderBookLevel{levelDelta},
				Timestamp: utils.MilliToTimestamp(ts),
				Trade:     false,
			}

			context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
				UpdateL2: obDelta,
				SeqNum:   state.instrumentData.seqNum + 1,
			})
			state.instrumentData.seqNum += 1
			instr.lastUpdateTime = ts
			///////////////

		}
	}

	instr.lastSequence = order.Sequence

	return nil
}

func (state *Listener) onReceivedOrder(order coinbasepro.WSReceivedOrder, context actor.Context) error {
	state.postAggTrade(context)
	instr := state.instrumentData

	// check sequence consistency
	if order.Sequence <= instr.lastSequence {
		return nil
	}
	if order.Sequence != instr.lastSequence+1 {
		return fmt.Errorf("received inconsistent sequence: %d, %d",
			order.Sequence,
			instr.lastSequence+1)
	}

	// No sync, no update time, no update id

	instr.lastSequence = order.Sequence

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

// Note, here we don't do the delay 20ms trick like in other exchanges
// to determine the end of an agg trade, because we know when it ends
// depending on the event received after
func (state *Listener) postAggTrade(context actor.Context) {
	if state.instrumentData.aggTrade != nil {
		context.Send(context.Parent(), &messages.MarketDataIncrementalRefresh{
			Trades: []*models.AggregatedTrade{state.instrumentData.aggTrade},
			SeqNum: state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1
		state.instrumentData.aggTrade = nil
	}
}
