package binance

import (
	"encoding/json"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/interface"
	"gitlab.com/alphaticks/alpha-connect/jobs"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/binance"
	xutils "gitlab.com/alphaticks/xchanger/utils"
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

type QueryRunner struct {
	pid                  *actor.PID
	secondOrderRateLimit *exchanges.RateLimit
	dayOrderRateLimit    *exchanges.RateLimit
	globalRateLimit      *exchanges.RateLimit
}

type Executor struct {
	client       *http.Client
	securities   []*models.Security
	queryRunners []*QueryRunner
	dialerPool   *xutils.DialerPool
	logger       *log.Logger
}

func NewExecutor(dialerPool *xutils.DialerPool) actor.Actor {
	return &Executor{
		client:       nil,
		queryRunners: nil,
		dialerPool:   dialerPool,
	}
}

func (state *Executor) Receive(context actor.Context) {
	_interface.ExchangeExecutorReceive(state, context)
}

func (state *Executor) GetLogger() *log.Logger {
	return state.logger
}

func (state *Executor) Initialize(context actor.Context) error {
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

	dialers := state.dialerPool.GetDialers()
	for _, dialer := range dialers {
		client := &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 1024,
				TLSHandshakeTimeout: 10 * time.Second,
				DialContext:         dialer.DialContext,
			},
			Timeout: 10 * time.Second,
		}
		props := actor.PropsFromProducer(func() actor.Actor {
			return jobs.NewAPIQuery(client)
		})
		state.queryRunners = append(state.queryRunners, &QueryRunner{
			pid:                  context.Spawn(props),
			secondOrderRateLimit: nil,
			dayOrderRateLimit:    nil,
			globalRateLimit:      nil,
		})
	}

	request, weight, err := binance.GetExchangeInfo()
	// Launch an APIQuery actor with the given request and target

	future := context.RequestFuture(state.queryRunners[0].pid, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
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
				for _, qr := range state.queryRunners {
					qr.secondOrderRateLimit = exchanges.NewRateLimit(rateLimit.Limit, time.Second)
				}
			} else if rateLimit.Interval == "DAY" {
				for _, qr := range state.queryRunners {
					qr.dayOrderRateLimit = exchanges.NewRateLimit(rateLimit.Limit, time.Hour*24)
				}
			}
		} else if rateLimit.RateLimitType == "REQUEST_WEIGHT" {
			for _, qr := range state.queryRunners {
				qr.globalRateLimit = exchanges.NewRateLimit(rateLimit.Limit, time.Minute)
			}
		}
	}
	if state.queryRunners[0].secondOrderRateLimit == nil || state.queryRunners[0].dayOrderRateLimit == nil {
		return fmt.Errorf("unable to set second or day rate limit")
	}

	// Update rate limit with weight from the current exchange info fetch
	state.queryRunners[0].globalRateLimit.Request(weight)

	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	// Get http request and the expected response
	request, weight, err := binance.GetExchangeInfo()
	if err != nil {
		return err
	}

	var qr *QueryRunner
	for _, q := range state.queryRunners {
		if !q.globalRateLimit.IsRateLimited() {
			qr = q
			break
		}
	}

	if qr == nil {
		return fmt.Errorf("rate limited")
	}

	qr.globalRateLimit.Request(weight)

	future := context.RequestFuture(qr.pid, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

	res, err := future.Result()
	if err != nil {
		return fmt.Errorf("http client error: %v", err)
	}
	queryResponse := res.(*jobs.PerformQueryResponse)
	if queryResponse.StatusCode != 200 {
		if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
			return fmt.Errorf(
				"http client error: %d %s",
				queryResponse.StatusCode,
				string(queryResponse.Response))
		} else if queryResponse.StatusCode >= 500 {
			return fmt.Errorf(
				"http server error: %d %s",
				queryResponse.StatusCode,
				string(queryResponse.Response))
		}
	}

	var exchangeInfo binance.ExchangeInfo
	err = json.Unmarshal(queryResponse.Response, &exchangeInfo)
	if err != nil {
		return fmt.Errorf(
			"error unmarshaling response: %v",
			err)
	}
	if exchangeInfo.Code != 0 {
		return fmt.Errorf(
			"binance api error: %d %s",
			exchangeInfo.Code,
			exchangeInfo.Msg)
	}

	var securities []*models.Security
	for _, symbol := range exchangeInfo.Symbols {
		base := symbol.BaseAsset
		if sym, ok := binance.BINANCE_TO_GLOBAL[base]; ok {
			base = sym
		}
		baseCurrency, ok := constants.GetAssetBySymbol(symbol.BaseAsset)
		if !ok {
			if symbol.Status == "TRADING" {
				//fmt.Println("UNKNOWN BASE CURRENCY", symbol.BaseAsset)
			}
			continue
		}
		quote := symbol.QuoteAsset
		if sym, ok := binance.BINANCE_TO_GLOBAL[quote]; ok {
			quote = sym
		}
		quoteCurrency, ok := constants.GetAssetBySymbol(quote)
		if !ok {
			//fmt.Println("UNKNOWN QUOTE CURRENCY", symbol.QuoteAsset)
			continue
		}
		security := models.Security{}
		security.Symbol = symbol.Symbol
		security.Underlying = baseCurrency
		security.QuoteCurrency = quoteCurrency
		switch symbol.Status {
		case "PRE_TRADING":
			security.Status = models.PreTrading
		case "TRADING":
			security.Status = models.Trading
		case "POST_TRADING":
			security.Status = models.PostTrading
		case "END_OF_DAY":
			security.Status = models.EndOfDay
		case "HALT":
			security.Status = models.Halt
		case "AUCTION_MATCH":
			security.Status = models.AuctionMatch
		case "BREAK":
			security.Status = models.Break
		default:
			security.Status = models.Disabled
		}
		security.Exchange = &constants.BINANCE
		security.SecurityType = enum.SecurityType_CRYPTO_SPOT
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		for _, filter := range symbol.Filters {
			if filter.FilterType == "PRICE_FILTER" {
				security.MinPriceIncrement = &types.DoubleValue{Value: filter.TickSize}
			} else if filter.FilterType == "LOT_SIZE" {
				security.RoundLot = &types.DoubleValue{Value: filter.StepSize}
			}
		}
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

func (state *Executor) OnHistoricalLiquidationsRequest(context actor.Context) error {
	msg := context.Message().(*messages.HistoricalLiquidationsRequest)
	context.Respond(&messages.HistoricalLiquidationsResponse{
		RequestID:       msg.RequestID,
		Liquidations:    nil,
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *Executor) OnMarketStatisticsRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketStatisticsResponse)
	context.Respond(&messages.MarketStatisticsResponse{
		RequestID:       msg.RequestID,
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *Executor) OnMarketDataRequest(context actor.Context) error {
	var snapshot *models.OBL2Snapshot
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
	if msg.Instrument == nil || msg.Instrument.Symbol == nil {
		response.RejectionReason = messages.MissingInstrument
		context.Respond(response)
		return nil
	}
	symbol := msg.Instrument.Symbol.Value
	// Get http request and the expected response
	request, weight, err := binance.GetOrderBook(symbol, 1000)
	if err != nil {
		return err
	}

	var qr *QueryRunner
	for _, q := range state.queryRunners {
		if !q.globalRateLimit.IsRateLimited() {
			qr = q
			break
		}
	}

	if qr == nil {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}
	qr.globalRateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

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
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"http server error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
			}
			return
		}
		var obData binance.OrderBookData
		err = json.Unmarshal(queryResponse.Response, &obData)
		if err != nil {
			state.logger.Info("error decoding query response", log.Error(err))
			response.RejectionReason = messages.HTTPError
			context.Respond(response)
			return
		}
		if obData.Code != 0 {
			err = fmt.Errorf("error getting orderbook: %d %s", obData.Code, obData.Msg)
			state.logger.Info("error decoding query response", log.Error(err))
			response.RejectionReason = messages.HTTPError
			context.Respond(response)
			return
		}

		bids, asks, err := obData.ToBidAsk()
		if err != nil {
			state.logger.Info("error decoding query response", log.Error(err))
			response.RejectionReason = messages.HTTPError
			context.Respond(response)
			return
		}
		snapshot = &models.OBL2Snapshot{
			Bids:      bids,
			Asks:      asks,
			Timestamp: &types.Timestamp{Seconds: 0, Nanos: 0},
		}
		response.SnapshotL2 = snapshot
		response.SeqNum = obData.LastUpdateID
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
