package okexp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/enum"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/jobs"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/okex"
	xutils "gitlab.com/alphaticks/xchanger/utils"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"
)

type QueryRunner struct {
	pid       *actor.PID
	rateLimit *exchanges.RateLimit
}

type Executor struct {
	extypes.BaseExecutor
	client       *http.Client
	securities   map[uint64]*models.Security
	queryRunners []*QueryRunner
	dialerPool   *xutils.DialerPool
	logger       *log.Logger
}

func NewExecutor(dialerPool *xutils.DialerPool) actor.Actor {
	return &Executor{
		client:     nil,
		dialerPool: dialerPool,
		logger:     nil,
	}
}

func (state *Executor) getQueryRunner() *QueryRunner {
	sort.Slice(state.queryRunners, func(i, j int) bool {
		return rand.Uint64()%2 == 0
	})

	var qr *QueryRunner
	for _, q := range state.queryRunners {
		if !q.rateLimit.IsRateLimited() {
			qr = q
			break
		}
	}

	return qr
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
			return jobs.NewHTTPQuery(client)
		})
		state.queryRunners = append(state.queryRunners, &QueryRunner{
			pid:       context.Spawn(props),
			rateLimit: exchanges.NewRateLimit(6, time.Second),
		})
	}

	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	qr := state.getQueryRunner()
	if qr == nil {
		return fmt.Errorf("rate limited")
	}

	var securities []*models.Security

	request, weight, err := okex.GetInstruments(okex.SWAP)
	if err != nil {
		return err
	}
	qr.rateLimit.Request(weight)

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

	var perpResponse okex.PerpetualInstrumentsResponse
	err = json.Unmarshal(response, &perpResponse)
	if err != nil {
		err = fmt.Errorf("error decoding query response: %v", err)
		return err
	}
	if perpResponse.Code != "0" {
		return fmt.Errorf("error decoding query response: %s", perpResponse.Msg)
	}

	for _, pair := range perpResponse.Data {
		splits := strings.Split(pair.Underlying, "-")
		baseCurrency, ok := constants.GetAssetBySymbol(splits[0])
		if !ok {
			//state.logger.Info("unknown symbol " + pair.BaseCurrency + " for instrument " + pair.InstrumentID)
			continue
		}
		quoteCurrency, ok := constants.GetAssetBySymbol(splits[1])
		if !ok {
			//state.logger.Info("unknown symbol " + pair.QuoteCurrency + " for instrument " + pair.InstrumentID)
			continue
		}

		security := models.Security{}
		security.Symbol = pair.InstrumentID
		security.Underlying = baseCurrency
		security.QuoteCurrency = quoteCurrency
		security.Status = models.Trading
		security.Exchange = &constants.OKEXP
		security.IsInverse = pair.ContractType == okex.INVERSE
		security.SecurityType = enum.SecurityType_CRYPTO_PERP
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinPriceIncrement = &types.DoubleValue{Value: pair.TickSize}
		security.RoundLot = &types.DoubleValue{Value: pair.LotSize}
		security.Multiplier = &types.DoubleValue{Value: pair.ContractValue}
		securities = append(securities, &security)
	}

	state.securities = make(map[uint64]*models.Security)
	for _, s := range securities {
		state.securities[s.SecurityID] = s
	}

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: securities})

	return nil
}

func (state *Executor) OnSecurityListRequest(context actor.Context) error {
	msg := context.Message().(*messages.SecurityListRequest)
	securities := make([]*models.Security, len(state.securities))
	i := 0
	for _, v := range state.securities {
		securities[i] = v
		i += 1
	}
	context.Respond(&messages.SecurityList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: securities})

	return nil
}

func (state *Executor) OnHistoricalLiquidationsRequest(context actor.Context) error {
	fmt.Println("HISTORICAL LIQUIDATIONS")
	msg := context.Message().(*messages.HistoricalLiquidationsRequest)
	response := &messages.HistoricalLiquidationsResponse{
		RequestID: msg.RequestID,
		Success:   false,
	}

	var security *models.Security
	if msg.Instrument != nil {
		if msg.Instrument.Symbol != nil {
			for _, sec := range state.securities {
				if sec.Symbol == msg.Instrument.Symbol.Value {
					security = sec
				}
			}
		} else if msg.Instrument.SecurityID != nil {
			securityID := msg.Instrument.SecurityID.Value
			sec, ok := state.securities[securityID]
			if ok {
				security = sec
			}
		}
	}

	if security == nil {
		response.Success = false
		response.RejectionReason = messages.UnknownSecurityID
		context.Respond(response)
		return nil
	}
	req := okex.NewLiquidationsRequest(okex.SWAP)
	if msg.From != nil {
		frm := uint64(msg.From.Seconds*1000) + uint64(msg.From.Nanos/1000000)
		req.SetBefore(frm)
		fmt.Println("FROM", frm)
	}
	if msg.To != nil {
		fmt.Println("TO", uint64(msg.To.Seconds*1000)+uint64(msg.From.Nanos/1000000))
		req.SetAfter(uint64(msg.To.Seconds*1000) + uint64(msg.From.Nanos/1000000))
	}
	if err := req.SetState("filled"); err != nil {
		return fmt.Errorf("error setting liquidation state: %v", err)
	}
	fmt.Println("underlying", strings.Replace(security.Symbol, "-SWAP", "", -1))
	if err := req.SetUnderlying(strings.Replace(security.Symbol, "-SWAP", "", -1)); err != nil {
		return fmt.Errorf("error setting underlying: %v", err)
	}
	request, weight, err := okex.GetLiquidations(req)
	if err != nil {
		return err
	}

	var qr *QueryRunner
	for _, q := range state.queryRunners {
		if !q.rateLimit.IsRateLimited() {
			qr = q
			break
		}
	}

	if qr == nil {
		response.Success = false
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.rateLimit.Request(weight)

	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Info("request error", log.Error(err))
			response.RejectionReason = messages.HTTPError
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
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
			} else if queryResponse.StatusCode >= 500 {

				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
			}
			return
		}

		var resx okex.LiquidationsResponse
		err = json.Unmarshal(queryResponse.Response, &resx)
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}

		if resx.Code != "0" {
			state.logger.Info("api error", log.Error(errors.New(resx.Msg)))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}

		var liquidations []*models.Liquidation
		for _, l := range resx.Data {
			fmt.Println("LEN", len(l.Details))
			for _, d := range l.Details {
				// Inverse will be in USD, Linear in BTC
				liquidations = append(liquidations, &models.Liquidation{
					Bid:       d.Side == "sell",
					Timestamp: utils.MilliToTimestamp(d.Ts),
					OrderID:   rand.Uint64(),
					Price:     d.BankruptcyPrice,
					Quantity:  d.Size * security.Multiplier.Value,
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
	fmt.Println("HISTORICAL FUNDING RATES")
	msg := context.Message().(*messages.HistoricalFundingRatesRequest)
	response := &messages.HistoricalFundingRatesResponse{
		RequestID: msg.RequestID,
		Success:   false,
	}

	var security *models.Security
	if msg.Instrument != nil {
		if msg.Instrument.Symbol != nil {
			for _, sec := range state.securities {
				if sec.Symbol == msg.Instrument.Symbol.Value {
					security = sec
				}
			}
		} else if msg.Instrument.SecurityID != nil {
			securityID := msg.Instrument.SecurityID.Value
			sec, ok := state.securities[securityID]
			if ok {
				security = sec
			}
		}
	}

	if security == nil {
		response.Success = false
		response.RejectionReason = messages.UnknownSecurityID
		context.Respond(response)
		return nil
	}
	req := okex.NewFundingRateHistoryRequest(security.Symbol)
	if msg.From != nil {
		frm := uint64(msg.From.Seconds*1000) + uint64(msg.From.Nanos/1000000)
		req.SetBefore(frm)
		fmt.Println("FROM", frm)
	}
	if msg.To != nil {
		fmt.Println("TO", uint64(msg.To.Seconds*1000)+uint64(msg.From.Nanos/1000000))
		req.SetAfter(uint64(msg.To.Seconds*1000) + uint64(msg.From.Nanos/1000000))
	}
	request, weight, err := okex.GetFundingRateHistory(req)
	if err != nil {
		return err
	}

	var qr *QueryRunner
	for _, q := range state.queryRunners {
		if !q.rateLimit.IsRateLimited() {
			qr = q
			break
		}
	}

	if qr == nil {
		response.Success = false
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.rateLimit.Request(weight)

	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Info("request error", log.Error(err))
			response.RejectionReason = messages.HTTPError
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
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
			} else if queryResponse.StatusCode >= 500 {

				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
			}
			return
		}

		var resx okex.FundingRateHistoryResponse
		err = json.Unmarshal(queryResponse.Response, &resx)
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}

		if resx.Code != "0" {
			state.logger.Info("api error", log.Error(errors.New(resx.Msg)))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}

		var rates []*models.Stat
		for _, l := range resx.Data {
			// Inverse will be in USD, Linear in BTC
			rates = append(rates, &models.Stat{
				Timestamp: utils.MilliToTimestamp(l.FundingTime),
				Value:     l.FundingRate,
				StatType:  models.FundingRate,
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
	response := &messages.MarketStatisticsResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	if msg.Instrument == nil || msg.Instrument.Symbol == nil {
		response.RejectionReason = messages.MissingInstrument
		context.Respond(response)
		return nil
	}

	var security *models.Security
	if msg.Instrument != nil {
		if msg.Instrument.Symbol != nil {
			for _, sec := range state.securities {
				if sec.Symbol == msg.Instrument.Symbol.Value {
					security = sec
				}
			}
		} else if msg.Instrument.SecurityID != nil {
			securityID := msg.Instrument.SecurityID.Value
			sec, ok := state.securities[securityID]
			if ok {
				security = sec
			}
		}
	}

	if security == nil {
		response.Success = false
		response.RejectionReason = messages.UnknownSecurityID
		context.Respond(response)
		return nil
	}

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
		case models.OpenInterest:
			if has(stat) {
				continue
			}
			request := okex.NewOpenInterestRequest(okex.SWAP)
			request.SetInstrumentID(security.Symbol)
			req, weight, err := okex.GetOpenInterest(request)
			if err != nil {
				return err
			}

			qr := state.getQueryRunner()

			if qr == nil {
				response.RejectionReason = messages.RateLimitExceeded
				context.Respond(response)
				return nil
			}

			qr.rateLimit.Request(weight)

			res, err := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: req}, 10*time.Second).Result()
			if err != nil {
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.HTTPError
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
					response.RejectionReason = messages.HTTPError
					context.Respond(response)
					return nil
				} else if queryResponse.StatusCode >= 500 {
					err := fmt.Errorf(
						"http server error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
					state.logger.Info("http client error", log.Error(err))
					response.RejectionReason = messages.HTTPError
					context.Respond(response)
					return nil
				}
				return nil
			}

			var resp okex.OpenInterestResponse
			err = json.Unmarshal(queryResponse.Response, &resp)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return nil
			}

			if resp.Code != "0" {
				state.logger.Info("http error", log.Error(errors.New(resp.Msg)))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return nil
			}
			for _, d := range resp.Data {
				response.Statistics = append(response.Statistics, &models.Stat{
					Timestamp: utils.MilliToTimestamp(d.Ts),
					StatType:  models.OpenInterest,
					Value:     d.OpenInterest * security.Multiplier.Value,
				})
			}
		}
	}

	response.Success = true
	context.Respond(response)

	return nil
}
