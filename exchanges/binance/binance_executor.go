package binance

import (
	"encoding/json"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	//fix50nso "github.com/quickfixgo/quickfix/fix50/newordersingle"
	"gitlab.com/alphaticks/alphac/exchanges/interface"
	"gitlab.com/alphaticks/alphac/jobs"
	models "gitlab.com/alphaticks/alphac/messages/exchanges"
	"gitlab.com/alphaticks/alphac/messages/executor"
	"gitlab.com/alphaticks/xchanger/asset"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/binance"
	"math"
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
		logger:               nil,
	}
}

func (state *Executor) Receive(context actor.Context) {
	_interface.ExchangeExecutorReceive(state, context)
}

func (state *Executor) GetLogger() *log.Logger {
	return state.logger
}

func (state *Executor) Initialize(context actor.Context) error {
	state.client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1024,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
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

func (state *Executor) GetInstrumentsRequest(context actor.Context) error {
	// Get http request and the expected response
	msg := context.Message().(*executor.GetInstrumentsRequest)
	request, weight, err := binance.GetExchangeInfo()
	if err != nil {
		return err
	}

	if state.globalRateLimit.IsRateLimited() {
		time.Sleep(state.globalRateLimit.DurationBeforeNextRequest(weight))
		return nil
	}

	// Launch an APIQuery actor with the given request and target
	state.globalRateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			context.Respond(&executor.GetInstrumentsResponse{
				RequestID:   msg.RequestID,
				Error:       err,
				Instruments: nil})
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"http client error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				context.Respond(&executor.GetInstrumentsResponse{
					RequestID:   msg.RequestID,
					Error:       err,
					Instruments: nil})
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"http server error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				context.Respond(&executor.GetInstrumentsResponse{
					RequestID:   msg.RequestID,
					Error:       err,
					Instruments: nil})
			}
			return
		}
		var exchangeInfo binance.ExchangeInfo
		err = json.Unmarshal(queryResponse.Response, &exchangeInfo)
		if err != nil {
			err = fmt.Errorf(
				"error unmarshaling response: %v",
				err)
			context.Respond(&executor.GetInstrumentsResponse{
				RequestID:   msg.RequestID,
				Error:       err,
				Instruments: nil})
			return
		}
		if exchangeInfo.Code != 0 {
			err = fmt.Errorf(
				"binance api error: %d %s",
				exchangeInfo.Code,
				exchangeInfo.Msg)
			context.Respond(&executor.GetInstrumentsResponse{
				RequestID:   msg.RequestID,
				Error:       err,
				Instruments: nil})
			return
		}

		var instruments []*exchanges.Instrument
		for _, symbol := range exchangeInfo.Symbols {
			if symbol.Status == "TRADING" {
				baseCurrency, ok := constants.SYMBOL_TO_ASSET[symbol.BaseAsset]
				if !ok {
					continue
				}
				quoteCurrency, ok := constants.SYMBOL_TO_ASSET[symbol.QuoteAsset]
				if !ok {
					continue
				}
				pair := asset.NewPair(&baseCurrency, &quoteCurrency)
				instrument := exchanges.Instrument{}
				instrument.Exchange = constants.BINANCE
				instrument.Pair = pair
				instrument.Type = exchanges.SPOT
				for _, filter := range symbol.Filters {
					if filter.FilterType == "PRICE_FILTER" {
						instrument.TickPrecision = uint64(math.Round(1.0 / filter.TickSize))
					} else if filter.FilterType == "LOT_SIZE" {
						instrument.LotPrecision = uint64(math.Round(1.0 / filter.StepSize))
					}
				}
				instruments = append(instruments, &instrument)
			}
		}
		context.Respond(&executor.GetInstrumentsResponse{
			RequestID:   msg.RequestID,
			Error:       nil,
			Instruments: instruments})
	})

	return nil
}

func (state *Executor) GetOrderBookL2Request(context actor.Context) error {
	var snapshot *models.OBL2Snapshot
	msg := context.Message().(*executor.GetOrderBookL2Request)

	symbol := msg.Instrument.Format(binance.SymbolFormat)
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
			context.Respond(&executor.GetOrderBookL2Response{
				RequestID: msg.RequestID,
				Error:     err,
				Snapshot:  nil})
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"http client error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				context.Respond(&executor.GetOrderBookL2Response{
					RequestID: msg.RequestID,
					Error:     err,
					Snapshot:  nil})
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"http server error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				context.Respond(&executor.GetOrderBookL2Response{
					RequestID: msg.RequestID,
					Error:     err,
					Snapshot:  nil})
			}
			return
		}
		var obData binance.OrderBookData
		err = json.Unmarshal(queryResponse.Response, &obData)
		if err != nil {
			err = fmt.Errorf("error decoding query response: %v", err)
			context.Respond(&executor.GetOrderBookL2Response{
				RequestID: msg.RequestID,
				Error:     err,
				Snapshot:  nil})
			return
		}
		if obData.Code != 0 {
			err = fmt.Errorf("error getting orderbook: %d %s", obData.Code, obData.Msg)
			context.Respond(&executor.GetOrderBookL2Response{
				RequestID: msg.RequestID,
				Error:     err,
				Snapshot:  nil})
			return
		}

		bids, asks, err := obData.ToRawBidAsk(msg.Instrument.TickPrecision, msg.Instrument.LotPrecision)
		if err != nil {
			err = fmt.Errorf("error converting orderbook: %v", err)
			context.Respond(&executor.GetOrderBookL2Response{
				RequestID: msg.RequestID,
				Error:     err,
				Snapshot:  nil})
			return
		}
		snapshot = &models.OBL2Snapshot{
			Instrument: msg.Instrument,
			Bids:       bids,
			Asks:       asks,
			Timestamp:  &types.Timestamp{Seconds: 0, Nanos: 0},
			ID:         obData.LastUpdateID,
		}
		context.Respond(&executor.GetOrderBookL2Response{
			RequestID: msg.RequestID,
			Error:     nil,
			Snapshot:  snapshot})
	})

	return nil
}

func (state *Executor) OnFIX50NewOrderSingle(context actor.Context) error {
	//msg := context.Message().(*fix50nso.NewOrderSingle)
	return nil
}

func (state *Executor) GetOrderBookL3Request(context actor.Context) error {
	return nil
}

func (state *Executor) GetOpenOrdersRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OpenOrdersRequest(context actor.Context) error {
	return nil
}

func (state *Executor) CloseOrdersRequest(context actor.Context) error {
	return nil
}

func (state *Executor) CloseAllOrdersRequest(context actor.Context) error {
	return nil
}
