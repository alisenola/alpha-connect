package legacy

/*
import (
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/bitfinex"
	"github.com/AsynkronIT/protoactor-go/actor"
	"fmt"
	"time"
	exchangeModels "gitlab.com/alphaticks/alphac/models/exchanges"
	"gitlab.com/alphaticks/xchanger/orderbook"
	"math"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/xchanger/constants"
)

type readSocket struct {}

type InstrumentData struct {
	Instrument 			exchanges.Instrument
	OrderBook 			*orderbook.OrderBookL2
	LastUpdateTime 		int64
	LastDepthTime		int64
	Depths				[]*bitfinex.WSSpotDepthData
	Trades 				[]*bitfinex.WSSpotTrade
	StashedOBRequests 	[]*actor.PID
}

type BitfinexListener struct {
	ws 			*bitfinex.BitfinexWebsocket
	instruments	map[string]*InstrumentData
	mediator 	*actor.PID
	executor 	*actor.PID
}

// the listener will publish all the collected data
// to a mediator.
// If someone wants to follow a topic ({exchange}/{symbol}/trades)

func NewBitfinexListener() actor.Actor {
	return &BitfinexListener{
		nil,
		nil,
		nil,
		nil,
	}
}

func (state *BitfinexListener) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		fmt.Println("Bitfinex listener started, initialize actor here")
		err := state.Initialize(context)
		if err != nil {
			panic(err)
		}
	case *actor.Stopping:
		fmt.Println("Bitfinex listener Stopping, actor is about shut down")

	case *actor.Stopped:
		fmt.Println("Bitfinex listener Stopped, actor and its children are stopped")

	case *actor.Restarting:
		fmt.Println("Bitfinex listener Restarting, actor is about restart")


	case *exchangeModels.SetupInstrumentsRequest:
		fmt.Println("setup instrument request received")
		// Setup the instruments listened by the socket
		err := state.SetupInstruments(context)
		if err != nil {
			panic(err)
		}

	case *exchangeModels.GetOrderBookL2Request:
		err := state.GetOrderBookL2Request(context)
		if err != nil {
			panic(err)
		}

	case *readSocket:
		err := state.readSocket(context)
		if err != nil {
			panic(err)
		}
	}
}

func (state *BitfinexListener) Initialize(context actor.Context) error {
	state.instruments = make(map[string]*InstrumentData)
	state.ws = bitfinex.NewBitfinexWebsocket()
	state.mediator = actor.NewLocalPID("data_broker")

	executorManager := actor.NewLocalPID("exchange_executor_manager")
	future := context.RequestFuture(
		executorManager,
		&exchangeModels.GetExchangeExecutorRequest{constants.BITFINEX},
		5 * time.Second)

	res, err := future.Result()
	if err != nil {
		return fmt.Errorf("error getting exchange executer: %v", err)
	}
	exchangeExecutorResponse := res.(*exchangeModels.GetExchangeExecutorResponse)

	state.executor = &exchangeExecutorResponse.PID

	return nil
}

func (state *BitfinexListener) SetupInstruments(context actor.Context) error {
	if !state.ws.Connected {
		err := state.ws.Connect()
		if err != nil {
			return fmt.Errorf("error connecting to the websocket: %v", err)
		}
		context.Send(context.Self(), &readSocket{})
	}
	instruments := context.Message().(*exchangeModels.SetupInstrumentsRequest).Instruments
	// We clear instruments and trades
	state.instruments = make(map[string]*InstrumentData)

	for _, instrument := range instruments {
		var symbol string
		if instrument.Type == exchanges.SPOT {
			symbol = bitfinex.InstrumentFormat(instrument, bitfinex.WSSpotSymbolFormat)
		} else {
			symbol = bitfinex.InstrumentFormat(instrument, bitfinex.WSFundingSymbolFormat)
		}
		err := state.ws.SubscribeDepth(
			symbol,
			bitfinex.WSDepthPrecisionR0,
			bitfinex.WSDepthFrequency0,
			bitfinex.WSDepthLength100)
		if err != nil {
			return err
		}

		err = state.ws.SubscribeTrades(symbol)
		if err != nil {
			return err
		}

		state.instruments[bitfinex.InstrumentFormat(instrument, bitfinex.WSSpotSymbolFormat)] = &InstrumentData{
			instrument,
			nil,
			0,
			0,
			[]*bitfinex.WSSpotDepthData{},
			[]*bitfinex.WSSpotTrade{},
			[]*actor.PID{},
		}
	}

	context.Respond(&exchangeModels.SetupInstrumentsResponse{})
	return nil
}

func (state *BitfinexListener) GetOrderBookL2Request(context actor.Context) error {
	msg := context.Message().(*exchangeModels.GetOrderBookL2Request)
	instrumentData, ok := state.instruments[bitfinex.InstrumentFormat(msg.Instrument, bitfinex.WSSpotSymbolFormat)]
	if !ok {
		return nil
	}

	if instrumentData.OrderBook == nil || !instrumentData.OrderBook.Synced {
		instrumentData.StashedOBRequests = append(instrumentData.StashedOBRequests, context.Sender())
		return nil
	}

	snapshot := exchangeModels.OBL2Snapshot{
		Instrument: &instrumentData.Instrument,
		Bids: instrumentData.OrderBook.GetAbsoluteRawBids(0),
		Asks: instrumentData.OrderBook.GetAbsoluteRawAsks(0),
		Timestamp: exchangeModels.MilliToTimestamp(instrumentData.LastUpdateTime),
		ID: instrumentData.LastUpdateTime,
	}
	context.Respond(&exchangeModels.GetOrderBookL2Response{nil, snapshot})
	return nil
}

func (state *BitfinexListener) readSocket(context actor.Context) error {
	msg, err := state.ws.ReadMessage()
	if err != nil {
		return err
	}
	switch msg.(type) {

	case bitfinex.WSErrorMessage:
		err := msg.(bitfinex.WSErrorMessage)
		return fmt.Errorf("socket error: %v", err)

	case bitfinex.WSSpotDepthSnapshot:
		obData := msg.(bitfinex.WSSpotDepthSnapshot)

		if _, ok := state.instruments[obData.Symbol]; !ok {
			return fmt.Errorf("received event for a non-subscribed symbol: %s", obData.Symbol)
		}

		bids, asks := obData.ToRawBidAsk(
			int(state.instruments[obData.Symbol].Instrument.TickPrecision),
			int(state.instruments[obData.Symbol].Instrument.LotPrecision))

		bestAsk := float64(asks[0].Price) / float64(state.instruments[obData.Symbol].Instrument.TickPrecision)

		// Allow a 10% price variation
		depth := int(((bestAsk * 1.1) - bestAsk) * float64(state.instruments[obData.Symbol].Instrument.TickPrecision))

		ob := orderbook.NewOrderBookL2(
			int(state.instruments[obData.Symbol].Instrument.TickPrecision),
			int(state.instruments[obData.Symbol].Instrument.LotPrecision),
			depth,
		)
		ob.RawSync(bids, asks)
		state.instruments[obData.Symbol].OrderBook = ob
		//fmt.Println(depth, obData.Symbol)

		state.instruments[obData.Symbol].LastUpdateTime = obData.Timestamp
		for _, pid := range state.instruments[obData.Symbol].StashedOBRequests {
			snapshot := exchangeModels.OBL2Snapshot{
				&state.instruments[obData.Symbol].Instrument,
				bids,
				asks,
				exchangeModels.MilliToTimestamp(obData.Timestamp),
				obData.Timestamp,
			}
			context.Send(pid, &exchangeModels.GetOrderBookL2Response{nil, snapshot})
		}

	case bitfinex.WSSpotDepthData:
		obData := msg.(bitfinex.WSSpotDepthData)
		if _, ok := state.instruments[obData.Symbol]; !ok {
			return fmt.Errorf("received event for a non-subscribed symbol: %s", obData.Symbol)
		}

		if state.instruments[obData.Symbol].OrderBook == nil || !state.instruments[obData.Symbol].OrderBook.Synced {
			break
		}
		// We agglomerate OBData by TS
		if state.instruments[obData.Symbol].LastDepthTime == 0 || state.instruments[obData.Symbol].LastDepthTime == obData.Timestamp {
			state.instruments[obData.Symbol].Depths = append(state.instruments[obData.Symbol].Depths, &obData)
			state.instruments[obData.Symbol].LastDepthTime = obData.Timestamp
		} else {
			symbol := obData.Symbol
			instrument := state.instruments[symbol].Instrument

			var bids []orderbook.RawOrderBookLevel
			var asks []orderbook.RawOrderBookLevel

			bidQuantity := make(map[int]int)
			askQuantity := make(map[int]int)

			for _, depth := range state.instruments[obData.Symbol].Depths {
				item := depth.Depth.ToRawOrderBookLevel(
					int(instrument.TickPrecision),
					int(instrument.LotPrecision))
				if depth.Depth.Amount > 0 {
					if _, ok := bidQuantity[int(item.Price)]; ok {
						fmt.Println("DOUBLE BID")
					}
					bidQuantity[int(item.Price)] = int(item.Quantity)
					bids = append(bids, item)
				} else {
					if _, ok := askQuantity[int(item.Price)]; ok {
						fmt.Println("DOUBLE ASK")
					}
					askQuantity[int(item.Price)] = int(item.Quantity)
					asks = append(asks, item)
				}
			}

			var deltasSlice []*exchangeModels.OBDelta

			var trades []*bitfinex.WSSpotTrade
			for i := 0; i < len(state.instruments[symbol].Trades); {
				if state.instruments[symbol].Trades[i].Timestamp <= state.instruments[symbol].LastDepthTime {
					trades = append(trades, state.instruments[symbol].Trades[i])
					state.instruments[symbol].Trades = append(
						state.instruments[symbol].Trades[:i],
						state.instruments[symbol].Trades[i+1:]...)
				} else {
					i += 1
				}
			}

			if len(trades) > 0 {
				ts := (state.instruments[symbol].LastUpdateTime + trades[0].Timestamp) / 2
				if ts <= state.instruments[symbol].LastUpdateTime {
					fmt.Println("INCREASE LAST UPDATE TIME")
					ts = state.instruments[symbol].LastUpdateTime + 1
				}
				firstObDelta := exchangeModels.OBDelta{
					TickPrecision: uint32(instrument.TickPrecision),
					LotPrecision: uint32(instrument.LotPrecision),
					LevelDeltas: []*exchangeModels.OBLevelDelta{},
					Timestamp: exchangeModels.MilliToTimestamp(ts),
					FirstID: ts,
					ID: ts,
					Trade: false,
				}
				state.instruments[symbol].LastUpdateTime = ts

				buyTradeAmounts := make(map[int]int)
				sellTradeAmounts := make(map[int]int)
				for _, trade := range trades {
					var rawTradePrice int
					var rawTradeQuantity int
					if trade.Amount > 0 {
						rawTradePrice = int(math.Ceil(trade.Price * float64(instrument.TickPrecision)))
						rawTradeQuantity = int(math.Round(trade.Amount * float64(instrument.LotPrecision)))
						if _, ok := buyTradeAmounts[rawTradePrice]; !ok {
							buyTradeAmounts[rawTradePrice] = 0
						}
						buyTradeAmounts[rawTradePrice] += rawTradeQuantity
					} else {
						rawTradePrice = int(math.Floor(trade.Price * float64(instrument.TickPrecision)))
						rawTradeQuantity = int(math.Round(-trade.Amount * float64(instrument.LotPrecision)))
						if _, ok := sellTradeAmounts[rawTradePrice]; !ok {
							sellTradeAmounts[rawTradePrice] = 0
						}
						sellTradeAmounts[rawTradePrice] += rawTradeQuantity
					}
				}

				// Check if there is more buy trades than ask
				for key, amount := range buyTradeAmounts {
					lastRawAsk := state.instruments[symbol].OrderBook.GetRawAsk(key)
					rawAsk, ok := askQuantity[key]
					if !ok {
						// There wasn't any change in the ask
						rawAsk = lastRawAsk
					}
					if lastRawAsk < amount {
						// There was more traded than quantity at the ask, add a limit
						// so that quantity - amount = rawAsk
						firstObDelta.LevelDeltas = append(
							firstObDelta.LevelDeltas,
							&exchangeModels.OBLevelDelta{
								RawPrice: int64(key),
								RawQuantity: int64(rawAsk + amount),
								RawChange: int64(rawAsk + amount - lastRawAsk),
								Bid: false,
							},
						)
					}
				}

				// Check if there is more sell trades than bid
				for key, amount := range sellTradeAmounts {
					lastRawBid := state.instruments[symbol].OrderBook.GetRawBid(key)
					rawBid, ok := bidQuantity[key]
					if !ok {
						// There wasn't any change in the bid
						rawBid = lastRawBid
					}
					if lastRawBid < amount {
						// There was more traded than quantity at the bid, add a limit
						// so that quantity - amount = rawBid
						firstObDelta.LevelDeltas = append(
							firstObDelta.LevelDeltas,
							&exchangeModels.OBLevelDelta{
								RawPrice: int64(key),
								RawQuantity: int64(rawBid + amount),
								RawChange: int64(rawBid + amount - lastRawBid),
								Bid: true,
							},
						)
					}
				}

				// Process the deltas
				for _, delta := range firstObDelta.LevelDeltas {
					state.instruments[symbol].OrderBook.UpdateRawOrderBookLevel(
						orderbook.RawOrderBookLevel{
							delta.RawPrice,
							delta.RawQuantity,
							delta.Bid})
				}

				if len(firstObDelta.LevelDeltas) > 0 {
					deltasSlice = append(deltasSlice, &firstObDelta)
				}

				for i := 0; i < len(trades); {
					tradeID := trades[i].Timestamp
					tradeTime := trades[i].Timestamp
					// We ensure that the trade time is no lower than last event time (last delta processed)
					if tradeTime <= state.instruments[symbol].LastUpdateTime {
						fmt.Println("INCREASE LAST UPDATE TIME")
						tradeTime = state.instruments[symbol].LastUpdateTime + 1
					}
					obDelta := exchangeModels.OBDelta{
						TickPrecision: uint32(instrument.TickPrecision),
						LotPrecision: uint32(instrument.LotPrecision),
						LevelDeltas: []*exchangeModels.OBLevelDelta{},
						Timestamp: exchangeModels.MilliToTimestamp(tradeTime),
						FirstID: tradeTime,
						ID: tradeTime,
						Trade: true,
					}
					state.instruments[symbol].LastUpdateTime = tradeTime

					// Aggregate trades with the same timestamp
					for ; i < len(trades) && tradeID == trades[i].Timestamp; {
						trade := trades[i]
						var rawTradePrice int
						var rawTradeQuantity int
						var lastRawQuantity int
						if trade.Amount > 0 {
							rawTradePrice = int(math.Ceil(trade.Price * float64(instrument.TickPrecision)))
							rawTradeQuantity = int(math.Round(trade.Amount * float64(instrument.LotPrecision)))
							lastRawQuantity = state.instruments[symbol].OrderBook.GetRawAsk(rawTradePrice)

						} else {
							rawTradePrice = int(math.Floor(trade.Price * float64(instrument.TickPrecision)))
							rawTradeQuantity = int(math.Round(-trade.Amount * float64(instrument.LotPrecision)))
							lastRawQuantity = state.instruments[symbol].OrderBook.GetRawBid(rawTradePrice)
						}
						// Aggregate OBItemDelta on the same price level
						lenItemDeltas := len(obDelta.LevelDeltas)
						if lenItemDeltas > 0 && obDelta.LevelDeltas[lenItemDeltas-1].RawPrice == int64(rawTradePrice) {
							obDelta.LevelDeltas[lenItemDeltas-1].RawQuantity -= int64(rawTradeQuantity)
							obDelta.LevelDeltas[lenItemDeltas-1].RawChange -= int64(rawTradeQuantity)
						} else {
							obDelta.LevelDeltas = append(
								obDelta.LevelDeltas,
								&exchangeModels.OBLevelDelta{
									RawPrice:    int64(rawTradePrice),
									RawQuantity: int64(lastRawQuantity - rawTradeQuantity),
									RawChange: int64(-rawTradeQuantity),
									Bid: trade.Amount < 0,
								})
						}

						i += 1
					}

					for _, delta := range obDelta.LevelDeltas {
						state.instruments[symbol].OrderBook.UpdateRawOrderBookLevel(
							orderbook.RawOrderBookLevel{
								delta.RawPrice,
								delta.RawQuantity,
								delta.Bid})
					}
					deltasSlice = append(deltasSlice, &obDelta)
				}
			}

			ts := state.instruments[symbol].LastDepthTime
			if ts <= state.instruments[symbol].LastUpdateTime {
				fmt.Println("INCREASE LAST UPDATE TIME")
				ts = state.instruments[symbol].LastUpdateTime + 1
			}
			lastObDelta := exchangeModels.OBDelta{
				TickPrecision: uint32(instrument.TickPrecision),
				LotPrecision: uint32(instrument.LotPrecision),
				LevelDeltas: []*exchangeModels.OBLevelDelta{},
				Timestamp: exchangeModels.MilliToTimestamp(ts),
				FirstID: ts,
				ID: ts,
				Trade: false,
			}
			state.instruments[obData.Symbol].LastUpdateTime = ts

			// Go for each bid/ask and check which ones have changed
			for _, bid := range bids {
				rawBid := state.instruments[symbol].OrderBook.GetRawBid(int(bid.Price))
				if rawBid != int(bid.Quantity) {
					// Add delta
					lastObDelta.LevelDeltas = append(
						lastObDelta.LevelDeltas,
						&exchangeModels.OBLevelDelta{
							RawPrice: bid.Price,
							RawQuantity: bid.Quantity,
							RawChange: bid.Quantity - int64(rawBid),
							Bid: true,
						},
					)
				}
				state.instruments[symbol].OrderBook.UpdateRawOrderBookLevel(bid)
			}

			for _, ask := range asks {
				rawAsk := state.instruments[symbol].OrderBook.GetRawAsk(int(ask.Price))
				if rawAsk != int(ask.Quantity) {
					// Add delta
					lastObDelta.LevelDeltas = append(
						lastObDelta.LevelDeltas,
						&exchangeModels.OBLevelDelta{
							RawPrice: ask.Price,
							RawQuantity: ask.Quantity,
							RawChange: ask.Quantity - int64(rawAsk),
							Bid: false,
						},
					)
				}
				state.instruments[symbol].OrderBook.UpdateRawOrderBookLevel(ask)
			}

			if len(lastObDelta.LevelDeltas) > 0 {
				deltasSlice = append(deltasSlice, &lastObDelta)
			}

			// Send the deltas to the mediator
			deltaTopic := fmt.Sprintf("%s/OBDELTA", instrument.DefaultFormat())
			for _, delta := range deltasSlice {
				context.Send(state.mediator, &models.PubSubMessage{deltaTopic, delta.ID,*delta})
			}


			state.instruments[symbol].Depths = []*bitfinex.WSSpotDepthData{&obData}
			state.instruments[symbol].LastDepthTime = obData.Timestamp

			rawBids := state.instruments[obData.Symbol].OrderBook.GetAbsoluteRawBids(25)
			rawAsks := state.instruments[obData.Symbol].OrderBook.GetAbsoluteRawAsks(25)

			// TODO out of range error
			snapshotL1Topic := fmt.Sprintf("%s/ORDERBOOKL1", instrument.DefaultFormat())
			snapshotL1 := exchangeModels.OBL2Snapshot{
				Instrument: &state.instruments[obData.Symbol].Instrument,
				Bids: rawBids[len(rawBids)-1:],
				Asks: rawAsks[0:1],
				Timestamp: exchangeModels.MilliToTimestamp(state.instruments[symbol].LastUpdateTime),
				ID: state.instruments[symbol].LastUpdateTime,
			}
			context.Send(state.mediator, &models.PubSubMessage{snapshotL1Topic, snapshotL1.ID,snapshotL1})

			snapshotL2Topic := fmt.Sprintf("%s/ORDERBOOKL2_25", instrument.DefaultFormat())
			snapshotL2_25 := exchangeModels.OBL2Snapshot{
				Instrument: &state.instruments[obData.Symbol].Instrument,
				Bids: rawBids,
				Asks: rawAsks,
				Timestamp: exchangeModels.MilliToTimestamp(state.instruments[symbol].LastUpdateTime),
				ID: state.instruments[symbol].LastUpdateTime,
			}
			context.Send(state.mediator, &models.PubSubMessage{snapshotL2Topic, snapshotL2_25.ID,snapshotL2_25})
		}

	case bitfinex.WSSpotTradeSnapshot:
		//
	case bitfinex.WSSpotTrade:
		trade := msg.(bitfinex.WSSpotTrade)
		if state.instruments[trade.Symbol].OrderBook !=  nil && state.instruments[trade.Symbol].OrderBook.Synced {
			// we can receive a depth update after the trade but which doesn't relate to the trade
			// even though the depth update timestamp is after the trade, so we increase the trade TS
			trade.Timestamp += 31
			state.instruments[trade.Symbol].Trades = append(
				state.instruments[trade.Symbol].Trades,
				&trade)
		}
	}

	context.Send(context.Self(), &readSocket{})
	return nil
}*/
