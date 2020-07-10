package bitfinex

import (
	"encoding/json"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/enum"
	"gitlab.com/alphaticks/alphac/exchanges/interface"
	"gitlab.com/alphaticks/alphac/jobs"
	models "gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/bitfinex"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
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
type Executor struct {
	client           *http.Client
	securities       []*models.Security
	obRateLimit      *exchanges.RateLimit
	symbolsRateLimit *exchanges.RateLimit
	queryRunner      *actor.PID
	logger           *log.Logger
}

func NewExecutor() actor.Actor {
	return &Executor{
		client:      nil,
		obRateLimit: nil,
		queryRunner: nil,
		logger:      nil,
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
	state.obRateLimit = exchanges.NewRateLimit(30, time.Minute)
	state.symbolsRateLimit = exchanges.NewRateLimit(10, time.Minute)

	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewAPIQuery(state.client)
	})
	state.queryRunner = context.Spawn(props)

	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	if state.symbolsRateLimit.IsRateLimited() {
		return fmt.Errorf("rate limited")
	}
	request, weight, err := bitfinex.GetSymbolsDetails()
	state.symbolsRateLimit.Request(weight)
	if err != nil {
		return err
	}
	resp, err := state.client.Do(request)
	if err != nil {
		return err
	}
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			err := fmt.Errorf(
				"http client error: %d %s",
				resp.StatusCode,
				string(response))
			return err
		} else if resp.StatusCode >= 500 {
			err := fmt.Errorf(
				"http server error: %d %s",
				resp.StatusCode,
				string(response))
			return err
		} else {
			err := fmt.Errorf("%d %s",
				resp.StatusCode,
				string(response))
			return err
		}
	}
	var symbolDetails []bitfinex.SymbolDetail
	err = json.Unmarshal(response, &symbolDetails)
	if err != nil {
		err = fmt.Errorf("error decoding query response: %v", err)
		return err
	}

	var securities []*models.Security
	for _, symbol := range symbolDetails {
		if len(symbol.Pair) != 6 {
			//state.logger.Warning(fmt.Sprintf("unknown symbol type: %s", symbol.Pair))
			continue
		}
		symbolStr := strings.ToUpper(symbol.Pair[:3])
		if sym, ok := bitfinex.BITFINEX_SYMBOL_TO_GLOBAL_SYMBOL[symbolStr]; ok {
			symbolStr = sym
		}
		baseCurrency, ok := constants.SYMBOL_TO_ASSET[symbolStr]
		if !ok {
			continue
		}
		symbolStr = strings.ToUpper(symbol.Pair[3:])
		if sym, ok := bitfinex.BITFINEX_SYMBOL_TO_GLOBAL_SYMBOL[symbolStr]; ok {
			symbolStr = sym
		}
		quoteCurrency, ok := constants.SYMBOL_TO_ASSET[symbolStr]
		if !ok {
			continue
		}

		pair := xchangerModels.Pair{
			Base:  &baseCurrency,
			Quote: &quoteCurrency,
		}
		security := models.Security{}
		security.Symbol = symbol.Pair
		security.Underlying = &baseCurrency
		security.QuoteCurrency = &quoteCurrency
		security.Exchange = &constants.BITFINEX
		security.SecurityType = enum.SecurityType_CRYPTO_SPOT
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name)
		security.RoundLot = 1. / 100000000.
		tickPrecision, ok := bitfinex.TickPrecisions[pair.String()]
		if !ok {
			//state.logger.Warning(fmt.Sprintf("TickPrecisions not defined for %s", instrument.DefaultFormat()))
			continue
		}
		security.MinPriceIncrement = 1. / float64(tickPrecision)
		securities = append(securities, &security)
	}

	state.securities = securities

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: state.securities})

	return nil
}

func (state *Executor) OnSecurityListRequest(context actor.Context) error {
	msg := context.Message().(*messages.SecurityListRequest)
	context.Respond(&messages.SecurityList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: state.securities})

	return nil
}

func (state *Executor) OnMarketDataRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketDataRequest)
	response := &messages.MarketDataResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	if msg.Subscribe {
		response.RejectionReason = messages.UnsupportedSubscription
		context.Respond(response)
		return nil
	}
	if msg.Instrument == nil {
		response.RejectionReason = messages.MissingInstrument
		context.Respond(response)
		return nil
	}
	symbol := msg.Instrument.Symbol.Value

	if state.obRateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}
	// Get http request and the expected response
	request, weight, err := bitfinex.GetOrderBook(symbol, 100, 100)
	if err != nil {
		return err
	}
	state.obRateLimit.Request(weight)

	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Info("http client error", log.Error(err))
			response.RejectionReason = messages.HTTPError
			context.Respond(response)
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)

		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"http client error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
				return
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"http server error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
				return
			}
			return
		}

		var obData bitfinex.OrderBookData
		err = json.Unmarshal(queryResponse.Response, &obData)
		if err != nil {
			err = fmt.Errorf("error decoding query response: %v", err)
			state.logger.Info("http client error", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		var bids []gorderbook.OrderBookLevel
		var asks []gorderbook.OrderBookLevel
		// TS is float in seconds, * 1000 + rounding to get millisecond
		var ts uint64 = 0
		for _, bid := range obData.Bids {
			if uint64(bid.Timestamp*1000) > ts {
				ts = uint64(bid.Timestamp * 1000)
			}
			bids = append(bids, gorderbook.OrderBookLevel{
				Price:    bid.Price,
				Quantity: bid.Amount,
				Bid:      true,
			})
		}
		for _, ask := range obData.Asks {
			if uint64(ask.Timestamp*1000) > ts {
				ts = uint64(ask.Timestamp * 1000)
			}
			asks = append(asks, gorderbook.OrderBookLevel{
				Price:    ask.Price,
				Quantity: ask.Amount,
				Bid:      false,
			})
		}

		snapshot := &models.OBL2Snapshot{
			Bids:      bids,
			Asks:      asks,
			Timestamp: utils.MilliToTimestamp(ts),
		}
		response.SnapshotL2 = snapshot
		response.SeqNum = ts
		response.Success = true
		context.Respond(response)
	})
	return nil
}

func (state *Executor) OnOrderStatusRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnPositionsRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnBalancesRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnNewOrderSingleRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnNewOrderBulkRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnOrderReplaceRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnOrderBulkReplaceRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnOrderCancelRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnOrderMassCancelRequest(context actor.Context) error {
	return nil
}
