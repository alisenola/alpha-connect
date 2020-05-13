package binance

import (
	"encoding/json"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alphac/enum"
	"gitlab.com/alphaticks/alphac/exchanges/interface"
	"gitlab.com/alphaticks/alphac/jobs"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/binance"
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

// The role of a Binance Executor is to
// process api request

// The global rate limit is per IP and the orderRateLimit is per
// account.

type Executor struct {
	client               *http.Client
	secondOrderRateLimit *exchanges.RateLimit
	dayOrderRateLimit    *exchanges.RateLimit
	globalRateLimit      *exchanges.RateLimit
	queryRunner          *actor.PID
	logger               *log.Logger
}

func NewExecutor() actor.Actor {
	return &Executor{
		client:               nil,
		secondOrderRateLimit: nil,
		dayOrderRateLimit:    nil,
		globalRateLimit:      nil,
		queryRunner:          nil,
	}
}

func (state *Executor) Receive(context actor.Context) {
	_interface.ExchangeExecutorReceive(state, context)
}

func (state *Executor) GetLogger() *log.Logger {
	return state.logger
}

func (state *Executor) Initialize(context actor.Context) error {
	state.client = &http.Client{Timeout: 10 * time.Second}
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewAPIQuery(state.client)
	})
	state.queryRunner = context.Spawn(props)

	request, weight, err := binance.GetExchangeInfo()
	// Launch an APIQuery actor with the given request and target

	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
	res, err := future.Result()
	if err != nil {
		return err
	}
	queryResponse := res.(*jobs.PerformQueryResponse)
	if queryResponse.StatusCode != 200 {
		return fmt.Errorf("error getting exchange info: status code %d", queryResponse.StatusCode)
	}

	var exchangeInfo binance.ExchangeInfo
	err = json.Unmarshal(queryResponse.Response, &exchangeInfo)
	if err != nil {
		return fmt.Errorf("error decoding query response: %v", err)
	}
	if exchangeInfo.Code != 0 {
		return fmt.Errorf("error getting exchange info: %s", exchangeInfo.Msg)
	}

	// Initialize rate limit
	for _, rateLimit := range exchangeInfo.RateLimits {
		if rateLimit.RateLimitType == "ORDERS" {
			if rateLimit.Interval == "SECOND" {
				state.secondOrderRateLimit = exchanges.NewRateLimit(rateLimit.Limit, time.Second)
			} else if rateLimit.Interval == "DAY" {
				state.dayOrderRateLimit = exchanges.NewRateLimit(rateLimit.Limit, time.Hour*24)
			}
		} else if rateLimit.RateLimitType == "REQUEST_WEIGHT" {
			state.globalRateLimit = exchanges.NewRateLimit(rateLimit.Limit, time.Minute)
		}
	}
	if state.secondOrderRateLimit == nil || state.dayOrderRateLimit == nil {
		return fmt.Errorf("unable to set second or day rate limit")
	}

	// Update rate limit with weight from the current exchange info fetch
	state.globalRateLimit.Request(weight)
	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) OnSecurityListRequest(context actor.Context) error {
	// Get http request and the expected response
	msg := context.Message().(*messages.SecurityListRequest)
	request, weight, err := binance.GetExchangeInfo()
	if err != nil {
		return err
	}

	if state.globalRateLimit.IsRateLimited() {
		time.Sleep(state.globalRateLimit.DurationBeforeNextRequest(weight))
		return nil
	}

	state.globalRateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			context.Respond(&messages.SecurityList{
				RequestID:  msg.RequestID,
				ResponseID: uint64(time.Now().UnixNano()),
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
		var exchangeInfo binance.ExchangeInfo
		err = json.Unmarshal(queryResponse.Response, &exchangeInfo)
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
		if exchangeInfo.Code != 0 {
			err = fmt.Errorf(
				"binance api error: %d %s",
				exchangeInfo.Code,
				exchangeInfo.Msg)
			context.Respond(&messages.SecurityList{
				RequestID:  msg.RequestID,
				ResponseID: uint64(time.Now().UnixNano()),
				Error:      err.Error(),
				Securities: nil})
			return
		}

		var securities []*models.Security
		for _, symbol := range exchangeInfo.Symbols {
			baseCurrency, ok := constants.SYMBOL_TO_ASSET[symbol.BaseAsset]
			if !ok {
				continue
			}
			quoteCurrency, ok := constants.SYMBOL_TO_ASSET[symbol.QuoteAsset]
			if !ok {
				continue
			}
			security := models.Security{}
			security.Underlying = &baseCurrency
			security.QuoteCurrency = &quoteCurrency
			security.Enabled = symbol.Status == "TRADING"
			security.Exchange = &constants.BINANCE
			security.Type = enum.SecurityType_CRYPTO_SPOT
			for _, filter := range symbol.Filters {
				if filter.FilterType == "PRICE_FILTER" {
					security.MinPriceIncrement = filter.TickSize
				} else if filter.FilterType == "LOT_SIZE" {
					security.RoundLot = filter.StepSize
				}
			}
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

func (state *Executor) OnMarketDataRequest(context actor.Context) error {
	var snapshot *models.OBL2Snapshot
	msg := context.Message().(*messages.MarketDataRequest)
	if msg.Subscribe {
		context.Respond(&messages.MarketDataRequestReject{
			RequestID: msg.RequestID,
			Reason:    "market data subscription not supported on executor"})
	}
	symbol := msg.Instrument.Symbol
	// Get http request and the expected response
	request, weight, err := binance.GetOrderBook(symbol, 1000)
	if err != nil {
		return err
	}

	if state.globalRateLimit.IsRateLimited() {
		time.Sleep(state.globalRateLimit.DurationBeforeNextRequest(weight))
		return nil
	}

	state.globalRateLimit.Request(weight)
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
		var obData binance.OrderBookData
		err = json.Unmarshal(queryResponse.Response, &obData)
		if err != nil {
			err = fmt.Errorf("error decoding query response: %v", err)
			context.Respond(&messages.MarketDataRequestReject{
				RequestID: msg.RequestID,
				Reason:    err.Error()})
			return
		}
		if obData.Code != 0 {
			err = fmt.Errorf("error getting orderbook: %d %s", obData.Code, obData.Msg)
			context.Respond(&messages.MarketDataRequestReject{
				RequestID: msg.RequestID,
				Reason:    err.Error()})
			return
		}

		bids, asks, err := obData.ToBidAsk()
		if err != nil {
			err = fmt.Errorf("error converting orderbook: %v", err)
			context.Respond(&messages.MarketDataRequestReject{
				RequestID: msg.RequestID,
				Reason:    err.Error()})
			return
		}
		snapshot = &models.OBL2Snapshot{
			Bids:      bids,
			Asks:      asks,
			Timestamp: &types.Timestamp{Seconds: 0, Nanos: 0},
			SeqNum:    obData.LastUpdateID,
		}
		context.Respond(&messages.MarketDataSnapshot{
			RequestID:  msg.RequestID,
			ResponseID: uint64(time.Now().UnixNano()),
			SnapshotL2: snapshot})
	})

	return nil
}
