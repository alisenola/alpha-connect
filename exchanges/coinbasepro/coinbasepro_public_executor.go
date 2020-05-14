package coinbasepro

import (
	"encoding/json"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/enum"
	"gitlab.com/alphaticks/alphac/jobs"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/coinbasepro"
	"net/http"
	"reflect"
	"time"
)

// Execute api calls
// Contains rate limit
// Spawn a query actor for each request
// and pipe its result back

// 429 rate limit
// 418 IP ban

// The role of a CoinbasePro Executor is to
// process api request
type CoinbaseProPublicExecutor struct {
	client           *http.Client
	rateLimit        *exchanges.RateLimit
	orderBookL2Cache *utils.TTLMap
	orderBookL3Cache *utils.TTLMap
	queryRunner      *actor.PID
	logger           *log.Logger
}

func NewCoinbaseProPublicExecutor() actor.Actor {
	return &CoinbaseProPublicExecutor{
		client:           nil,
		rateLimit:        nil,
		orderBookL2Cache: nil,
		orderBookL3Cache: nil,
		queryRunner:      nil,
		logger:           nil,
	}
}

func (state *CoinbaseProPublicExecutor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.GetLogger().Error("error initializing", log.Error(err))
			panic(err)
		}
		state.GetLogger().Info("actor started")

	case *actor.Stopping:
		if err := state.Clean(context); err != nil {
			state.GetLogger().Error("error stopping", log.Error(err))
			panic(err)
		}
		state.GetLogger().Info("actor stopping")

	case *actor.Stopped:
		state.GetLogger().Info("actor stopped")

	case *actor.Restarting:
		if err := state.Clean(context); err != nil {
			state.GetLogger().Error("error restarting", log.Error(err))
			// Attention, no panic in restarting or infinite loop
		}
		state.GetLogger().Info("actor restarting")

	case *messages.SecurityListRequest:
		if err := state.OnSecurityListRequest(context); err != nil {
			state.GetLogger().Error("error processing OnSecurityListRequest", log.Error(err))
			panic(err)
		}

	case *messages.MarketDataRequest:
		if err := state.OnMarketDataRequest(context); err != nil {
			state.GetLogger().Error("error processing OnMarketDataRequest", log.Error(err))
			panic(err)
		}
	}
}

func (state *CoinbaseProPublicExecutor) GetLogger() *log.Logger {
	return state.logger
}

func (state *CoinbaseProPublicExecutor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	state.client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1024,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
	state.rateLimit = exchanges.NewRateLimit(3, time.Second)

	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewAPIQuery(state.client)
	})
	state.queryRunner = context.Spawn(props)

	// 5 seconds cache for orderbooks
	state.orderBookL2Cache = utils.NewTTLMap(5)
	state.orderBookL3Cache = utils.NewTTLMap(5)
	return nil
}

func (state *CoinbaseProPublicExecutor) Clean(context actor.Context) error {
	return nil
}

func (state *CoinbaseProPublicExecutor) OnSecurityListRequest(context actor.Context) error {
	// Get http request and the expected response
	msg := context.Message().(*messages.SecurityListRequest)
	request, weight, err := coinbasepro.GetProducts()
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		time.Sleep(state.rateLimit.DurationBeforeNextRequest(weight))
	}

	state.rateLimit.Request(weight)

	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res interface{}, err error) {
		var securities []*models.Security
		if err != nil {
			context.Respond(&messages.SecurityList{
				RequestID:  msg.RequestID,
				ResponseID: uint64(time.Now().UnixNano()),
				Error:      err.Error(),
				Securities: nil})
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"http client error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				context.Respond(&messages.SecurityList{
					RequestID:  msg.RequestID,
					ResponseID: uint64(time.Now().UnixNano()),
					Error:      err.Error(),
					Securities: nil})
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"http server error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				context.Respond(&messages.SecurityList{
					RequestID:  msg.RequestID,
					ResponseID: uint64(time.Now().UnixNano()),
					Error:      err.Error(),
					Securities: nil})
			}
			return
		}
		var products []coinbasepro.Product
		err = json.Unmarshal(queryResponse.Response, &products)
		if err != nil {
			err = fmt.Errorf(
				"error unmarshaling response: %v",
				err)
			context.Respond(&messages.SecurityList{
				RequestID:  msg.RequestID,
				ResponseID: uint64(time.Now().UnixNano()),
				Error:      err.Error(),
				Securities: nil})
			return
		}

		for _, product := range products {
			baseCurrency, ok := constants.SYMBOL_TO_ASSET[product.BaseCurrency]
			if !ok {
				continue
			}
			quoteCurrency, ok := constants.SYMBOL_TO_ASSET[product.QuoteCurrency]
			if !ok {
				continue
			}
			security := models.Security{}
			security.Symbol = product.ID
			security.Underlying = &baseCurrency
			security.QuoteCurrency = &quoteCurrency
			security.Enabled = true
			security.Exchange = &constants.COINBASEPRO
			security.SecurityType = enum.SecurityType_CRYPTO_SPOT
			security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name)
			security.MinPriceIncrement = product.QuoteIncrement
			security.RoundLot = 1. / 100000000.
			securities = append(securities, &security)
		}
		context.Respond(&messages.SecurityList{
			RequestID:  msg.RequestID,
			ResponseID: uint64(time.Now().UnixNano()),
			Error:      "",
			Securities: securities})
	})

	return nil
}

func (state *CoinbaseProPublicExecutor) OnMarketDataRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketDataRequest)
	if msg.Subscribe {
		context.Respond(&messages.MarketDataRequestReject{
			RequestID: msg.RequestID,
			Reason:    "market data subscription not supported on executor"})
	}
	symbol := msg.Instrument.Symbol
	if msg.Aggregation == messages.L2 {
		var snapshot *models.OBL2Snapshot

		// Check if we don't have it already cached
		if ob, ok := state.orderBookL2Cache.Get(symbol); ok {
			context.Respond(&messages.MarketDataSnapshot{
				RequestID:  msg.RequestID,
				ResponseID: uint64(time.Now().UnixNano()),
				SnapshotL2: ob.(*models.OBL2Snapshot)})
			return nil
		}

		// Get http request and the expected response
		request, weight, err := coinbasepro.GetProductOrderBook(symbol, coinbasepro.L2ORDERBOOK)
		if err != nil {
			return err
		}

		if state.rateLimit.IsRateLimited() {
			time.Sleep(state.rateLimit.DurationBeforeNextRequest(weight))
		}

		state.rateLimit.Request(weight)
		future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

		context.AwaitFuture(future, func(res interface{}, err error) {
			if err != nil {
				context.Respond(&messages.MarketDataRequestReject{
					RequestID: msg.RequestID,
					Reason:    err.Error()})
				return
			}
			queryResponse := res.(*jobs.PerformQueryResponse)
			if queryResponse.StatusCode != 200 {
				if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
					err := fmt.Errorf(
						"http client error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
					context.Respond(&messages.MarketDataRequestReject{
						RequestID: msg.RequestID,
						Reason:    err.Error()})
				} else if queryResponse.StatusCode >= 500 {
					err := fmt.Errorf(
						"http server error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
					context.Respond(&messages.MarketDataRequestReject{
						RequestID: msg.RequestID,
						Reason:    err.Error()})
				}
				return
			}
			var obData coinbasepro.OrderBookL2
			err = json.Unmarshal(queryResponse.Response, &obData)
			if err != nil {
				err = fmt.Errorf("error decoding query response: %v", err)
				context.Respond(&messages.MarketDataRequestReject{
					RequestID: msg.RequestID,
					Reason:    err.Error()})
				return
			}

			bids, asks := obData.ToBidAsk()
			snapshot = &models.OBL2Snapshot{
				Bids:      bids,
				Asks:      asks,
				Timestamp: utils.MilliToTimestamp(0),
				SeqNum:    obData.Sequence,
			}
			state.orderBookL2Cache.Put(symbol, snapshot)
			context.Respond(&messages.MarketDataSnapshot{
				RequestID:  msg.RequestID,
				ResponseID: uint64(time.Now().UnixNano()),
				SnapshotL2: snapshot})
		})

		return nil
	} else {
		var snapshot *models.OBL3Snapshot

		// Check if we don't have it already cached

		if ob, ok := state.orderBookL3Cache.Get(symbol); ok {
			context.Respond(&messages.MarketDataSnapshot{
				RequestID:  msg.RequestID,
				ResponseID: uint64(time.Now().UnixNano()),
				SnapshotL3: ob.(*models.OBL3Snapshot)})
			return nil
		}

		// Get http request and the expected response
		request, weight, err := coinbasepro.GetProductOrderBook(symbol, coinbasepro.L3ORDERBOOK)
		if err != nil {
			return err
		}

		if state.rateLimit.IsRateLimited() {
			time.Sleep(state.rateLimit.DurationBeforeNextRequest(weight))
		}

		state.rateLimit.Request(weight)
		future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

		context.AwaitFuture(future, func(res interface{}, err error) {
			if err != nil {
				context.Respond(&messages.MarketDataRequestReject{
					RequestID: msg.RequestID,
					Reason:    err.Error()})
				return
			}
			queryResponse := res.(*jobs.PerformQueryResponse)
			if queryResponse.StatusCode != 200 {
				if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
					err := fmt.Errorf(
						"http client error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
					context.Respond(&messages.MarketDataRequestReject{
						RequestID: msg.RequestID,
						Reason:    err.Error()})
				} else if queryResponse.StatusCode >= 500 {
					err := fmt.Errorf(
						"http server error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
					context.Respond(&messages.MarketDataRequestReject{
						RequestID: msg.RequestID,
						Reason:    err.Error()})
				}
				return
			}

			var obData coinbasepro.OrderBookL3
			err = json.Unmarshal(queryResponse.Response, &obData)
			if err != nil {
				err = fmt.Errorf("error decoding query response: %v", err)
				context.Respond(&messages.MarketDataRequestReject{
					RequestID: msg.RequestID,
					Reason:    err.Error()})
				return
			}

			bids, asks := obData.ToBidAsk()
			snapshot = &models.OBL3Snapshot{
				Bids:      bids,
				Asks:      asks,
				Timestamp: nil,
				SeqNum:    obData.Sequence,
			}
			state.orderBookL3Cache.Put(symbol, snapshot)
			context.Respond(&messages.MarketDataSnapshot{
				RequestID:  msg.RequestID,
				ResponseID: uint64(time.Now().UnixNano()),
				SnapshotL3: snapshot})
		})

		return nil
	}
}
