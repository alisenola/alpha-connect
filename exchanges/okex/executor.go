package okex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/enum"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/jobs"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/okex"
	xutils "gitlab.com/alphaticks/xchanger/utils"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Executor struct {
	extypes.BaseExecutor
	client      *http.Client
	securities  []*models.Security
	rateLimit   *exchanges.RateLimit
	queryRunner *actor.PID
	logger      *log.Logger
}

func NewExecutor(dialerPool *xutils.DialerPool, registry registry.StaticClient) actor.Actor {
	e := &Executor{}
	e.DialerPool = dialerPool
	e.Registry = registry
	return e
}

func (state *Executor) Receive(context actor.Context) {
	extypes.ReceiveExecutor(state, context)
}

func (state *Executor) GetLogger() *log.Logger {
	return state.logger
}

func (state *Executor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(state).String()))

	state.client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1024,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	dialers := state.DialerPool.GetDialers()
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
			return jobs.NewHTTPQuery(client)
		})
		state.queryRunner = context.Spawn(props)
		state.rateLimit = exchanges.NewRateLimit(6, time.Second)
	}

	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	if state.rateLimit.IsRateLimited() {
		return fmt.Errorf("rate limited")
	}

	var securities []*models.Security

	request, weight, err := okex.GetInstruments(okex.SPOT)
	if err != nil {
		return err
	}
	state.rateLimit.Request(weight)

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

	var kResponse okex.SpotInstrumentsResponse
	err = json.Unmarshal(response, &kResponse)
	if err != nil {
		err = fmt.Errorf("error decoding query response: %v", err)
		return err
	}
	if kResponse.Code != "0" {
		return fmt.Errorf("error decoding query response: %s", kResponse.Msg)
	}

	for _, pair := range kResponse.Data {
		baseCurrency, ok := constants.GetAssetBySymbol(pair.BaseCurrency)
		if !ok {
			//state.logger.Info("unknown symbol " + pair.BaseCurrency + " for instrument " + pair.InstrumentID)
			continue
		}
		quoteCurrency, ok := constants.GetAssetBySymbol(pair.QuoteCurrency)
		if !ok {
			//state.logger.Info("unknown symbol " + pair.QuoteCurrency + " for instrument " + pair.InstrumentID)
			continue
		}

		security := models.Security{}
		security.Symbol = pair.InstrumentID
		security.Underlying = baseCurrency
		security.QuoteCurrency = quoteCurrency
		security.Status = models.InstrumentStatus_Trading
		security.Exchange = constants.OKEX
		security.SecurityType = enum.SecurityType_CRYPTO_SPOT
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinPriceIncrement = &wrapperspb.DoubleValue{Value: pair.TickSize}
		security.RoundLot = &wrapperspb.DoubleValue{Value: pair.LotSize}
		securities = append(securities, &security)
	}
	state.SyncSecurities(securities, nil)

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: securities})

	return nil
}

func (state *Executor) OnHistoricalLiquidationsRequest(context actor.Context) error {
	msg := context.Message().(*messages.HistoricalLiquidationsRequest)
	response := &messages.HistoricalLiquidationsResponse{
		RequestID: msg.RequestID,
		Success:   false,
	}

	security, rej := state.InstrumentToSecurity(msg.Instrument)

	if rej != nil {
		response.RejectionReason = *rej
		context.Respond(response)
		return nil
	}
	req := okex.NewLiquidationsRequest(okex.SWAP)
	if msg.From != nil {
		frm := uint64(msg.From.Seconds*1000) + uint64(msg.From.Nanos/1000000)
		req.SetBefore(frm)
	}
	if msg.To != nil {
		req.SetAfter(uint64(msg.To.Seconds*1000) + uint64(msg.From.Nanos/1000000))
	}
	if err := req.SetState("filled"); err != nil {
		return fmt.Errorf("error setting liquidation state: %v", err)
	}
	if err := req.SetUnderlying(strings.Replace(security.Symbol, "-SWAP", "", -1)); err != nil {
		return fmt.Errorf("error setting underlying: %v", err)
	}
	request, weight, err := okex.GetLiquidations(req)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		response.Success = false
		response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.rateLimit.Request(weight)

	future := context.RequestFuture(state.queryRunner, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

	context.ReenterAfter(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Info("request error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
			} else if queryResponse.StatusCode >= 500 {

				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
			}
			return
		}

		var resx okex.LiquidationsResponse
		err = json.Unmarshal(queryResponse.Response, &resx)
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		if resx.Code != "0" {
			state.logger.Info("api error", log.Error(errors.New(resx.Msg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		var liquidations []*models.Liquidation
		for _, l := range resx.Data {
			for _, d := range l.Details {
				// Inverse will be in USD, Linear in BTC
				liquidations = append(liquidations, &models.Liquidation{
					Bid:       d.Side == "buy",
					Timestamp: utils.MilliToTimestamp(d.Ts),
					OrderID:   rand.Uint64(),
					Price:     d.BankruptcyPrice,
					Quantity:  d.Size,
				})
			}
		}

		sort.Slice(liquidations, func(i, j int) bool {
			return utils.TimestampToMilli(liquidations[i].Timestamp) < utils.TimestampToMilli(liquidations[j].Timestamp)
		})
		response.Liquidations = liquidations
		response.Success = true
		context.Respond(response)
	})

	return nil
}

func (state *Executor) OnHistoricalFundingRatesRequest(context actor.Context) error {
	msg := context.Message().(*messages.HistoricalFundingRatesRequest)
	response := &messages.HistoricalFundingRatesResponse{
		RequestID: msg.RequestID,
		Success:   false,
	}

	req := okex.NewFundingRateHistoryRequest("BTC-USDT-SWAP")
	if msg.From != nil {
		frm := uint64(msg.From.Seconds*1000) + uint64(msg.From.Nanos/1000000)
		req.SetBefore(frm)
	}
	if msg.To != nil {
		req.SetAfter(uint64(msg.To.Seconds*1000) + uint64(msg.From.Nanos/1000000))
	}
	request, weight, err := okex.GetFundingRateHistory(req)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		response.Success = false
		response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.rateLimit.Request(weight)

	future := context.RequestFuture(state.queryRunner, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

	context.ReenterAfter(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Info("request error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
			} else if queryResponse.StatusCode >= 500 {

				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
			}
			return
		}

		var resx okex.FundingRateHistoryResponse
		err = json.Unmarshal(queryResponse.Response, &resx)
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		if resx.Code != "0" {
			state.logger.Info("api error", log.Error(errors.New(resx.Msg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		var rates []*models.Stat
		for _, l := range resx.Data {
			// Inverse will be in USD, Linear in BTC
			rates = append(rates, &models.Stat{
				Timestamp: utils.MilliToTimestamp(l.FundingTime),
				Value:     l.FundingRate,
				StatType:  models.StatType_FundingRate,
			})
		}

		sort.Slice(rates, func(i, j int) bool {
			return utils.TimestampToMilli(rates[i].Timestamp) < utils.TimestampToMilli(rates[j].Timestamp)
		})
		response.Rates = rates
		response.Success = true
		context.Respond(response)
	})

	return nil
}

func (state *Executor) OnMarketStatisticsRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketStatisticsRequest)
	fmt.Println("qqqqqqqqqqqqqqqq", msg)
	response := &messages.MarketStatisticsResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	if msg.Instrument == nil || msg.Instrument.Symbol == nil {
		response.RejectionReason = messages.RejectionReason_MissingInstrument
		context.Respond(response)
		return nil
	}
	security, rej := state.InstrumentToSecurity(msg.Instrument)

	if rej != nil {
		response.RejectionReason = *rej
		context.Respond(response)
		return nil
	}
	fmt.Println("wwwwwwwwwwwwwwwwww", security, rej)

	has := func(stat models.StatType) bool {
		for _, s := range response.Statistics {
			if s.StatType == stat {
				return true
			}
		}
		return false
	}
	for _, stat := range msg.Statistics {
		switch stat {
		case models.StatType_OpenInterest:
			if has(stat) {
				continue
			}
			request := okex.NewOpenInterestRequest(okex.SWAP)
			request.SetInstrumentID(security.Symbol)
			req, weight, err := okex.GetOpenInterest(request)
			if err != nil {
				return err
			}

			if state.rateLimit.IsRateLimited() {
				response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
				context.Respond(response)
				return nil
			}

			state.rateLimit.Request(weight)

			res, err := context.RequestFuture(state.queryRunner, &jobs.PerformHTTPQueryRequest{Request: req}, 10*time.Second).Result()
			if err != nil {
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
				return nil
			}
			queryResponse := res.(*jobs.PerformQueryResponse)
			if queryResponse.StatusCode != 200 {
				if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
					err := fmt.Errorf(
						"http client error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
					state.logger.Info("http client error", log.Error(err))
					response.RejectionReason = messages.RejectionReason_HTTPError
					context.Respond(response)
					return nil
				} else if queryResponse.StatusCode >= 500 {
					err := fmt.Errorf(
						"http server error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
					state.logger.Info("http client error", log.Error(err))
					response.RejectionReason = messages.RejectionReason_HTTPError
					context.Respond(response)
					return nil
				}
				return nil
			}

			var resp okex.OpenInterestResponse
			err = json.Unmarshal(queryResponse.Response, &resp)
			fmt.Println("eeeeeeeeeeeeeeeee", resp)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return nil
			}

			if resp.Code != "0" {
				state.logger.Info("http error", log.Error(errors.New(resp.Msg)))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return nil
			}
			for _, d := range resp.Data {
				response.Statistics = append(response.Statistics, &models.Stat{
					Timestamp: utils.MilliToTimestamp(d.Ts),
					StatType:  models.StatType_OpenInterest,
					Value:     d.OpenInterest * security.Multiplier.Value,
				})
			}
		}
	}

	response.Success = true
	context.Respond(response)

	return nil
}
