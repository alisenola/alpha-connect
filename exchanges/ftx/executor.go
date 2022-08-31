package ftx

import (
	"encoding/json"
	"errors"
	"fmt"
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
	"gitlab.com/alphaticks/xchanger/exchanges/ftx"
	xutils "gitlab.com/alphaticks/xchanger/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"strconv"
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

// The global rate limit is per IP and the orderRateLimit is per
// account.

type QueryRunner struct {
	pid       *actor.PID
	rateLimit *exchanges.RateLimit
}

type Executor struct {
	extypes.BaseExecutor
	client         *http.Client
	queryRunners   []*QueryRunner
	orderRateLimit *exchanges.RateLimit
	logger         *log.Logger
}

func NewExecutor(dialerPool *xutils.DialerPool, registry registry.StaticClient) actor.Actor {
	e := &Executor{}
	e.DialerPool = dialerPool
	e.Registry = registry
	return e
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
		log.String("type", reflect.TypeOf(state).String()))

	state.orderRateLimit = exchanges.NewRateLimit(30, time.Second)

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
		state.queryRunners = append(state.queryRunners, &QueryRunner{
			pid:       context.Spawn(props),
			rateLimit: exchanges.NewRateLimit(500, time.Minute),
		})
	}

	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	request, weight, err := ftx.GetMarkets()
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		return fmt.Errorf("rate limited")
	}

	qr.rateLimit.Request(weight)

	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

	res, err := future.Result()
	if err != nil {
		return fmt.Errorf("http client error: %v", err)
	}
	resp := res.(*jobs.PerformQueryResponse)

	if resp.StatusCode != 200 {
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			err := fmt.Errorf(
				"http client error: %d %s",
				resp.StatusCode,
				string(resp.Response))
			return err
		} else if resp.StatusCode >= 500 {
			err := fmt.Errorf(
				"http server error: %d %s",
				resp.StatusCode,
				string(resp.Response))
			return err
		} else {
			err := fmt.Errorf("%d %s",
				resp.StatusCode,
				string(resp.Response))
			return err
		}
	}
	type Res struct {
		Success bool         `json:"success"`
		Markets []ftx.Market `json:"result"`
	}
	var queryRes Res
	err = json.Unmarshal(resp.Response, &queryRes)
	if err != nil {
		err = fmt.Errorf(
			"error unmarshaling response: %v",
			err)
		return err
	}
	if !queryRes.Success {
		err = fmt.Errorf("ftx unsuccessful request")
		return err
	}

	var securities []*models.Security
	for _, market := range queryRes.Markets {
		security := models.Security{}
		security.Symbol = market.Name
		security.Exchange = constants.FTX
		security.Price = wrapperspb.Double(market.Price)

		if strings.Contains(market.BaseCurrency, "BEAR") ||
			strings.Contains(market.BaseCurrency, "BULL") ||
			strings.Contains(market.BaseCurrency, "HALF") ||
			strings.Contains(market.BaseCurrency, "HEDGE") {
			continue
		}
		switch market.Type {
		case "spot":
			baseCurrency, ok := constants.GetAssetBySymbol(market.BaseCurrency)
			if !ok {
				state.logger.Debug(fmt.Sprintf("unknown currency %s", market.BaseCurrency))
				continue
			}
			quoteCurrency, ok := constants.GetAssetBySymbol(market.QuoteCurrency)
			if !ok {
				state.logger.Debug(fmt.Sprintf("unknown currency %s", market.QuoteCurrency))
				continue
			}
			security.Underlying = baseCurrency
			security.QuoteCurrency = quoteCurrency
			security.SecurityType = enum.SecurityType_CRYPTO_SPOT

		case "future":
			splits := strings.Split(market.Name, "-")
			if len(splits) == 2 {
				underlying, ok := constants.GetAssetBySymbol(market.Underlying)
				if !ok {
					/*
						if splits[1] == "PERP" {
							state.logger.Info(fmt.Sprintf("unknown currency %s", market.Underlying))
						}
					*/
					continue
				}
				security.Underlying = underlying
				security.QuoteCurrency = constants.DOLLAR
				if splits[1] == "PERP" {
					security.SecurityType = enum.SecurityType_CRYPTO_PERP
				} else {
					year := time.Now().Format("2006")
					date, err := time.Parse("20060102", year+splits[1])
					if err != nil {
						continue
					}
					security.SecurityType = enum.SecurityType_CRYPTO_FUT
					security.MaturityDate = timestamppb.New(date)
				}
				security.Multiplier = &wrapperspb.DoubleValue{Value: 1.}
			} else {
				continue
			}

		default:
			continue
		}
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinLimitQuantity = &wrapperspb.DoubleValue{Value: market.MinProvideSize}
		security.MinPriceIncrement = &wrapperspb.DoubleValue{Value: market.PriceIncrement}
		security.RoundLot = &wrapperspb.DoubleValue{Value: market.SizeIncrement}
		security.Status = models.InstrumentStatus_Trading

		securities = append(securities, &security)
	}

	state.SyncSecurities(securities, nil)

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: securities})

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
		response.RejectionReason = messages.RejectionReason_MissingInstrument
		context.Respond(response)
		return nil
	}

	symbol := msg.Instrument.Symbol.Value
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
		case models.StatType_OpenInterest, models.StatType_FundingRate:
			if has(stat) {
				continue
			}
			req, weight, err := ftx.GetFutureStats(symbol)
			if err != nil {
				return err
			}

			qr := state.getQueryRunner()

			if qr == nil {
				response.RejectionReason = messages.RejectionReason_RateLimitExceeded
				context.Respond(response)
				return nil
			}

			qr.rateLimit.Request(weight)

			res, err := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: req}, 10*time.Second).Result()
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

			var fstat ftx.FutureStatsResponse
			err = json.Unmarshal(queryResponse.Response, &fstat)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return nil
			}

			if !fstat.Success {
				state.logger.Info("http error", log.Error(errors.New(fstat.Error)))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return nil
			}
			ts := time.Now().UnixNano() / 1000000
			response.Statistics = append(response.Statistics, &models.Stat{
				Timestamp: utils.MilliToTimestamp(uint64(ts)),
				StatType:  models.StatType_OpenInterest,
				Value:     fstat.Result.OpenInterest,
			})
			if state.SymbolToSecurity(symbol).SecurityType == enum.SecurityType_CRYPTO_PERP {
				response.Statistics = append(response.Statistics, &models.Stat{
					Timestamp: utils.MilliToTimestamp(uint64(fstat.Result.NextFundingTime.UnixNano() / 1000000)),
					StatType:  models.StatType_FundingRate,
					Value:     fstat.Result.NextFundingRate,
				})
			}
		}
	}

	response.Success = true
	context.Respond(response)

	return nil
}

func (state *Executor) OnAccountInformationRequest(context actor.Context) error {
	msg := context.Message().(*messages.AccountInformationRequest)
	response := &messages.AccountInformationResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	request, weight, err := ftx.GetAccountInformation(msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.rateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

	context.ReenterAfter(future, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.RejectionReason_Other
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
				response.Success = false
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
			}
			return
		}
		var information ftx.AccountInformationResponse
		err = json.Unmarshal(queryResponse.Response, &information)
		if err != nil {
			state.logger.Info("unmarshaling error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if !information.Success {
			state.logger.Info("api error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		response.MakerFee = &wrapperspb.DoubleValue{Value: information.Result.MakerFee}
		response.TakerFee = &wrapperspb.DoubleValue{Value: information.Result.TakerFee}
		response.Success = true
		context.Respond(response)
	})

	return nil
}

func (state *Executor) OnAccountMovementRequest(context actor.Context) error {
	msg := context.Message().(*messages.AccountMovementRequest)
	response := &messages.AccountMovementResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	var from, to *uint64
	var symbol string
	if msg.Filter != nil {
		if msg.Filter.Instrument != nil {
			s, rej := state.InstrumentToSymbol(msg.Filter.Instrument)
			if rej != nil {
				response.RejectionReason = *rej
				context.Respond(response)
				return nil
			}
			symbol = s
		}

		if msg.Filter.From != nil {
			ms := uint64(msg.Filter.From.Seconds)
			from = &ms
		}
		if msg.Filter.To != nil {
			ms := uint64(msg.Filter.To.Seconds)
			to = &ms
		}
	}

	var request *http.Request
	var weight int
	var err error

	switch msg.Type {
	case messages.AccountMovementType_Deposit:
		request, weight, err = ftx.GetDepositHistory(from, to, msg.Account.ApiCredentials)
	case messages.AccountMovementType_Withdrawal:
		request, weight, err = ftx.GetWithdrawalHistory(from, to, msg.Account.ApiCredentials)
	case messages.AccountMovementType_FundingFee:
		request, weight, err = ftx.GetFundingPayments(symbol, from, to, msg.Account.ApiCredentials)
	case messages.AccountMovementType_RealizedPnl,
		messages.AccountMovementType_WelcomeBonus,
		messages.AccountMovementType_Commission:
		response.RejectionReason = messages.RejectionReason_UnsupportedRequest
		context.Respond(response)
		return nil
	}
	if err != nil {
		response.RejectionReason = messages.RejectionReason_InvalidRequest
		context.Respond(response)
		return nil
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.rateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

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
		var deposits []ftx.Deposit
		var withdrawals []ftx.Withdrawal
		var fundings []ftx.FundingPayment
		switch msg.Type {
		case messages.AccountMovementType_Deposit:
			res := ftx.DepositsResponse{}
			err = json.Unmarshal(queryResponse.Response, &res)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			if !res.Success {
				state.logger.Info("api error", log.Error(errors.New(res.Error)))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			deposits = res.Result
		case messages.AccountMovementType_Withdrawal:
			res := ftx.WithdrawalResponse{}
			err = json.Unmarshal(queryResponse.Response, &res)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			if !res.Success {
				state.logger.Info("api error", log.Error(errors.New(res.Error)))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			withdrawals = res.Result
		case messages.AccountMovementType_FundingFee:
			res := ftx.FundingPaymentsResponse{}
			err = json.Unmarshal(queryResponse.Response, &res)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			if !res.Success {
				state.logger.Info("api error", log.Error(errors.New(res.Error)))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			fundings = res.Result
		}

		var movements []*messages.AccountMovement
		for _, t := range deposits {
			asset, ok := constants.GetAssetBySymbol(t.Coin)
			if !ok {
				state.logger.Warn("unknown asset " + t.Coin)
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			ts := timestamppb.New(t.Time)
			mvt := messages.AccountMovement{
				Asset:      asset,
				Change:     t.Size,
				MovementID: fmt.Sprintf("%s%d", msg.Account.ApiCredentials.AccountID, t.ID),
				Time:       ts,
				Type:       messages.AccountMovementType_Deposit,
			}
			movements = append(movements, &mvt)
		}
		for _, t := range withdrawals {
			asset, ok := constants.GetAssetBySymbol(t.Coin)
			if !ok {
				state.logger.Warn("unknown asset " + t.Coin)
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			ts := timestamppb.New(t.Time)
			mvt := messages.AccountMovement{
				Asset:      asset,
				Change:     -t.Size,
				MovementID: fmt.Sprintf("%s%d", msg.Account.ApiCredentials.AccountID, t.ID),
				Time:       ts,
				Type:       messages.AccountMovementType_Withdrawal,
			}
			movements = append(movements, &mvt)
		}
		for _, t := range fundings {
			ts := timestamppb.New(t.Time)
			mvt := messages.AccountMovement{
				Asset:      constants.DOLLAR,
				Change:     -t.Payment,
				MovementID: fmt.Sprintf("%d", t.ID),
				Time:       ts,
				Type:       messages.AccountMovementType_FundingFee,
				Subtype:    t.Future,
			}
			movements = append(movements, &mvt)
		}

		response.Success = true
		response.Movements = movements
		context.Respond(response)
	})

	return nil
}

func (state *Executor) OnTradeCaptureReportRequest(context actor.Context) error {
	msg := context.Message().(*messages.TradeCaptureReportRequest)
	sender := context.Sender()
	response := &messages.TradeCaptureReport{
		RequestID: msg.RequestID,
		Success:   false,
	}

	symbol := ""
	var from, to *float64
	var orderID *uint64
	if msg.Filter != nil {
		if msg.Filter.Side != nil || msg.Filter.OrderID != nil || msg.Filter.ClientOrderID != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedFilter
			context.Send(sender, response)
			return nil
		}

		if msg.Filter.Instrument != nil {
			s, rej := state.InstrumentToSymbol(msg.Filter.Instrument)
			if rej != nil {
				response.RejectionReason = *rej
				context.Send(sender, response)
				return nil
			}
			symbol = s
		}

		if msg.Filter.From != nil {
			ts := utils.TimestampToMilli(msg.Filter.From)
			ms := float64(ts) / 1000
			from = &ms
		}
		if msg.Filter.To != nil {
			ts := utils.TimestampToMilli(msg.Filter.To)
			ms := float64(ts) / 1000
			to = &ms
		}
		if msg.Filter.OrderID != nil {
			v, err := strconv.ParseUint(msg.Filter.OrderID.Value, 10, 64)
			if err != nil {
				response.RejectionReason = messages.RejectionReason_InvalidRequest
				context.Respond(response)
				return nil
			}
			orderID = &v
		}
		if msg.Filter.FromID != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Respond(response)
			return nil
		}
	}
	request, weight, err := ftx.GetFills(symbol, from, to, "asc", orderID, msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.rateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

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
				fmt.Println(err)
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
			} else if queryResponse.StatusCode >= 500 {

				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				fmt.Println(err)
				context.Respond(response)
			}
			return
		}
		var trades ftx.FillsResponse
		err = json.Unmarshal(queryResponse.Response, &trades)
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		var mtrades []*models.TradeCapture
		sort.Slice(trades.Result, func(i, j int) bool {
			return trades.Result[i].Time.Before(trades.Result[j].Time)
		})
		for _, t := range trades.Result {
			quantityMul := 1.
			var instrument *models.Instrument
			market := t.Market
			if market == "LUNC-PERP" {
				market = "LUNA-PERP"
			}
			if market == "USTC-PERP" {
				market = "UST-PERP"
			}
			sec := state.SymbolToHistoricalSecurity(market)
			if sec == nil {
				state.logger.Warn(fmt.Sprintf("unknown security: %s", market))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}

			instrument = &models.Instrument{
				Exchange:   constants.FTX,
				Symbol:     &wrapperspb.StringValue{Value: market},
				SecurityID: &wrapperspb.UInt64Value{Value: sec.SecurityId},
			}

			quantity := t.Size * quantityMul
			if t.Side == ftx.SELL {
				quantity *= -1
			}
			comAsset, ok := constants.GetAssetBySymbol(t.FeeCurrency)
			if !ok {
				state.logger.Info("api error", log.Error(fmt.Errorf("unknown commission asset %s", t.FeeCurrency)))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			ts := timestamppb.New(t.Time)
			trd := models.TradeCapture{
				Type:            models.TradeType_Regular,
				Price:           t.Price,
				Quantity:        quantity,
				Commission:      t.Fee,
				CommissionAsset: comAsset,
				TradeID:         fmt.Sprintf("%d-%d", t.TradeID, t.OrderID),
				Instrument:      instrument,
				Trade_LinkID:    nil,
				OrderID:         &wrapperspb.StringValue{Value: fmt.Sprintf("%d", t.OrderID)},
				TransactionTime: ts,
			}

			if t.Side == ftx.BUY {
				trd.Side = models.Side_Buy
			} else {
				trd.Side = models.Side_Sell
			}
			mtrades = append(mtrades, &trd)
		}
		response.Success = true
		response.Trades = mtrades
		context.Respond(response)
	})

	return nil
}

func (state *Executor) OnOrderStatusRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderStatusRequest)
	sender := context.Sender()
	response := &messages.OrderList{
		RequestID: msg.RequestID,
		Success:   false,
	}
	symbol := ""
	orderID := ""
	clOrderID := ""
	var orderStatus *models.OrderStatus
	if msg.Filter != nil {
		if msg.Filter.OrderStatus != nil {
			orderStatus = &msg.Filter.OrderStatus.Value
		}
		if msg.Filter.Side != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedFilter
			context.Respond(response)
			return nil
		}
		if msg.Filter.Instrument != nil {
			s, rej := state.InstrumentToSymbol(msg.Filter.Instrument)
			if rej != nil {
				response.RejectionReason = *rej
				context.Send(sender, response)
				return nil
			}
			symbol = s
		}
		if msg.Filter.OrderID != nil {
			orderID = msg.Filter.OrderID.Value
		}
		if msg.Filter.ClientOrderID != nil {
			clOrderID = msg.Filter.ClientOrderID.Value
		}
	}

	var request *http.Request
	var weight int
	var multi bool
	if orderID != "" || clOrderID != "" {
		multi = false
		if orderID != "" {
			var err error
			request, weight, err = ftx.GetOrderStatus(orderID, msg.Account.ApiCredentials)
			if err != nil {
				return err
			}
		} else {
			var err error
			request, weight, err = ftx.GetOrderStatusByClientID(orderID, msg.Account.ApiCredentials)
			if err != nil {
				return err
			}
		}
	} else {
		multi = true
		var err error
		request, weight, err = ftx.GetOpenOrders(symbol, msg.Account.ApiCredentials)
		if err != nil {
			return err
		}
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.rateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

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
		var orders []ftx.Order
		if multi {
			var oor ftx.OpenOrdersResponse
			err = json.Unmarshal(queryResponse.Response, &oor)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			if !oor.Success {
				state.logger.Info("API error", log.String("error", oor.Error))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			orders = oor.Result
		} else {
			var oor ftx.OrderStatusResponse
			err = json.Unmarshal(queryResponse.Response, &oor)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			if !oor.Success {
				state.logger.Info("API error", log.String("error", oor.Error))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			orders = []ftx.Order{oor.Result}
		}

		var morders []*models.Order
		for _, o := range orders {
			sec := state.SymbolToSecurity(o.Market)
			if sec == nil {
				state.logger.Warn(fmt.Sprintf("unknown symbol: %s", o.Market))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			ord := OrderToModel(o)
			ord.Instrument.SecurityID = wrapperspb.UInt64(sec.SecurityID)
			if orderStatus != nil && ord.OrderStatus != *orderStatus {
				continue
			}
			morders = append(morders, ord)
		}
		response.Success = true
		response.Orders = morders
		context.Respond(response)
	})

	return nil
}

func (state *Executor) OnPositionsRequest(context actor.Context) error {
	msg := context.Message().(*messages.PositionsRequest)
	sender := context.Sender()
	response := &messages.PositionList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
		Positions:  nil,
	}

	symbol := ""
	if msg.Instrument != nil {
		s, rej := state.InstrumentToSymbol(msg.Instrument)
		if rej != nil {
			response.RejectionReason = *rej
			context.Send(sender, response)
			return nil
		}
		symbol = s
	}

	request, weight, err := ftx.GetPositions(true, msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.rateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

	context.ReenterAfter(future, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.RejectionReason_Other
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
				response.Success = false
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
			}
			return
		}
		var positions ftx.PositionsResponse
		err = json.Unmarshal(queryResponse.Response, &positions)
		if err != nil {
			state.logger.Info("unmarshaling error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if !positions.Success {
			state.logger.Info("api error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		for _, p := range positions.Result {
			if p.Size == 0 {
				continue
			}
			if symbol != "" && p.Future != symbol {
				continue
			}
			sec := state.SymbolToSecurity(p.Future)
			if sec == nil {
				state.logger.Warn(fmt.Sprintf("unkown symbol %s", p.Future))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			cost := *p.RecentAverageOpenPrice * sec.Multiplier.Value * p.NetSize
			size := p.NetSize
			pos := &models.Position{
				Account: msg.Account.Name,
				Instrument: &models.Instrument{
					Exchange:   constants.BITMEX,
					Symbol:     &wrapperspb.StringValue{Value: p.Future},
					SecurityID: &wrapperspb.UInt64Value{Value: sec.SecurityID},
				},
				Quantity: size,
				Cost:     cost,
				Cross:    false,
			}
			response.Positions = append(response.Positions, pos)
		}
		response.Success = true
		context.Respond(response)
	})

	return nil
}

func (state *Executor) OnBalancesRequest(context actor.Context) error {
	msg := context.Message().(*messages.BalancesRequest)

	response := &messages.BalanceList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
		Balances:   nil,
	}

	request, weight, err := ftx.GetBalances(msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.rateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

	context.ReenterAfter(future, func(res interface{}, err error) {
		if err != nil {
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
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http server error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			}
			return
		}
		var balances ftx.BalancesResponse
		err = json.Unmarshal(queryResponse.Response, &balances)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if !balances.Success {
			state.logger.Info("error unmarshalling", log.Error(errors.New(balances.Error)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		for _, b := range balances.Result {
			if b.Total == 0. {
				continue
			}
			asset, ok := constants.GetAssetBySymbol(b.Coin)
			if !ok {
				state.logger.Error("got balance for unknown asset", log.String("asset", b.Coin))
				continue
			}
			fmt.Println("BALANCE", asset, b.Total)
			response.Balances = append(response.Balances, &models.Balance{
				Account:  msg.Account.Name,
				Asset:    asset,
				Quantity: b.Total,
			})
		}

		response.Success = true
		context.Respond(response)
	})

	return nil
}

func (state *Executor) OnNewOrderSingleRequest(context actor.Context) error {
	req := context.Message().(*messages.NewOrderSingleRequest)
	sender := context.Sender()
	response := &messages.NewOrderSingleResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	var tickPrecision, lotPrecision int
	sec, rej := state.InstrumentToSecurity(req.Order.Instrument)
	if rej != nil {
		response.RejectionReason = *rej
		context.Send(sender, response)
		return nil
	}

	tickPrecision = int(math.Ceil(math.Log10(1. / sec.MinPriceIncrement.Value)))
	lotPrecision = int(math.Ceil(math.Log10(1. / sec.RoundLot.Value)))

	params, rej := buildPlaceOrderRequest(sec.Symbol, req.Order, tickPrecision, lotPrecision)
	if rej != nil {
		response.RejectionReason = *rej
		context.Respond(response)
		return nil
	}

	request, weight, err := ftx.PlaceOrder(params, req.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	if state.orderRateLimit.IsRateLimited() {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.orderRateLimit.Request(1)
	qr.rateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
	context.ReenterAfter(future, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode >= 500 {
			err := fmt.Errorf(
				"%d %s",
				queryResponse.StatusCode,
				string(queryResponse.Response))
			state.logger.Info("http server error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		var order ftx.NewOrderResponse
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		b, _ := json.Marshal(order)
		fmt.Println("ORDER RESPONSE", string(b))
		if !order.Success {
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		status := OrderToStatus(&order.Result)
		if status == nil {
			state.logger.Error(fmt.Sprintf("unknown status %s", order.Result.Status))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
		response.OrderStatus = *status
		response.CumQuantity = order.Result.FilledSize
		response.LeavesQuantity = order.Result.RemainingSize
		response.OrderID = fmt.Sprintf("%d", order.Result.ID)
		context.Respond(response)
	})
	return nil
}

func (state *Executor) OnOrderReplaceRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderReplaceRequest)
	response := &messages.OrderReplaceResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	params := ftx.ModifyOrderRequest{}
	if req.Update.OrderID == nil {
		response.RejectionReason = messages.RejectionReason_UnknownOrder
		context.Respond(response)
		return nil
	}

	if req.Update.Price != nil {
		params.Price = &req.Update.Price.Value
	}
	if req.Update.Quantity != nil {
		params.Size = &req.Update.Quantity.Value
	}

	request, weight, err := ftx.ModifyOrder(req.Update.OrderID.Value, params, req.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}
	if state.orderRateLimit.IsRateLimited() {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.orderRateLimit.Request(1)
	qr.rateLimit.Request(weight)

	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
	context.ReenterAfter(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode >= 500 {
			err := fmt.Errorf(
				"%d %s",
				queryResponse.StatusCode,
				string(queryResponse.Response))
			state.logger.Info("http server error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		var order ftx.ModifyOrderResponse
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if !order.Success {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
		response.OrderID = fmt.Sprintf("%d", order.Result.ID)
		context.Respond(response)
	})
	return nil
}

func (state *Executor) OnOrderCancelRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderCancelRequest)
	response := &messages.OrderCancelResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	var request *http.Request
	var weight int
	if req.OrderID != nil {
		var err error
		request, weight, err = ftx.CancelOrder(req.OrderID.Value, req.Account.ApiCredentials)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_InvalidRequest
			context.Respond(response)
			return nil
		}
	} else if req.ClientOrderID != nil {
		// TODO
	} else {
		response.RejectionReason = messages.RejectionReason_InvalidRequest
		context.Respond(response)
		return nil
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.rateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
	context.ReenterAfter(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode >= 500 {
			err := fmt.Errorf(
				"%d %s",
				queryResponse.StatusCode,
				string(queryResponse.Response))
			state.logger.Info("http server error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		var order ftx.CancelOrderResponse
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if !order.Success {
			if order.Error == "Order already queued for cancellation" || order.Error == "Order already closed" {
				fmt.Println("order already queued for cancel or already closed", req.OrderID)
				response.RejectionReason = messages.RejectionReason_UnknownOrder
				context.Respond(response)
				return
			} else {
				state.logger.Info("error unmarshalling", log.String("error", order.Error))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
		}
		response.Success = true
		context.Respond(response)
	})
	return nil
}

func (state *Executor) OnOrderMassCancelRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderMassCancelRequest)
	sender := context.Sender()
	response := &messages.OrderMassCancelResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	symbol := ""
	if req.Filter != nil {
		if req.Filter.Instrument != nil {
			s, rej := state.InstrumentToSymbol(req.Filter.Instrument)
			if rej != nil {
				response.RejectionReason = *rej
				context.Send(sender, response)
				return nil
			}
			symbol = s
		}
		if req.Filter.Side != nil || req.Filter.OrderStatus != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedFilter
			context.Respond(response)
			return nil
		}
	}
	if symbol == "" {
		response.RejectionReason = messages.RejectionReason_UnknownSymbol
		context.Respond(response)
		return nil
	}

	request, weight, err := ftx.CancelAllOrders(ftx.CancelAllOrdersRequest{
		Market: symbol,
	}, req.Account.ApiCredentials)
	if err != nil {
		return err
	}
	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.rateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
	context.ReenterAfter(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Info("http error", log.Error(err))
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
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http server error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			}
			return
		}
		var fres ftx.CancelAllOrdersResponse
		err = json.Unmarshal(queryResponse.Response, &fres)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		if !fres.Success {
			state.logger.Info("error unmarshalling", log.Error(errors.New(fres.Error)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
		context.Respond(response)
	})
	return nil
}
