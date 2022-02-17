package ftx

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
	"gitlab.com/alphaticks/xchanger/exchanges/ftx"
	xutils "gitlab.com/alphaticks/xchanger/utils"
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
	securities     map[uint64]*models.Security
	symbolToSec    map[string]*models.Security
	queryRunners   []*QueryRunner
	orderRateLimit *exchanges.RateLimit
	dialerPool     *xutils.DialerPool
	logger         *log.Logger
}

func NewExecutor(dialerPool *xutils.DialerPool) actor.Actor {
	return &Executor{
		client:         nil,
		queryRunners:   nil,
		dialerPool:     dialerPool,
		orderRateLimit: nil,
		logger:         nil,
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

	state.orderRateLimit = exchanges.NewRateLimit(30, time.Second)

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
		security.Exchange = &constants.FTX

		switch market.Type {
		case "spot":
			baseCurrency, ok := constants.GetAssetBySymbol(market.BaseCurrency)
			if !ok {
				//fmt.Printf("unknown currency symbol %s \n", market.BaseCurrency)
				continue
			}
			quoteCurrency, ok := constants.GetAssetBySymbol(market.QuoteCurrency)
			if !ok {
				//fmt.Printf("unknown currency symbol %s \n", market.BaseCurrency)
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
					if splits[1] == "PERP" {
						fmt.Printf("unknown currency symbol %s \n", market.Underlying)
					}
					continue
				}
				security.Underlying = underlying
				security.QuoteCurrency = &constants.DOLLAR
				if splits[1] == "PERP" {
					security.SecurityType = enum.SecurityType_CRYPTO_PERP
				} else {
					year := time.Now().Format("2006")
					date, err := time.Parse("20060102", year+splits[1])
					if err != nil {
						continue
					}
					security.SecurityType = enum.SecurityType_CRYPTO_FUT
					security.MaturityDate, err = types.TimestampProto(date)
					if err != nil {
						continue
					}
				}
				security.Multiplier = &types.DoubleValue{Value: 1.}
			} else {
				continue
			}

		default:
			continue
		}
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinPriceIncrement = &types.DoubleValue{Value: market.PriceIncrement}
		security.RoundLot = &types.DoubleValue{Value: market.SizeIncrement}
		security.Status = models.Trading

		securities = append(securities, &security)
	}
	state.securities = make(map[uint64]*models.Security)
	state.symbolToSec = make(map[string]*models.Security)

	for _, sec := range securities {
		state.securities[sec.SecurityID] = sec
		state.symbolToSec[sec.Symbol] = sec
	}
	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: securities})

	return nil
}

func (state *Executor) OnSecurityListRequest(context actor.Context) error {
	// Get http request and the expected response
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
		case models.OpenInterest, models.FundingRate:
			if has(stat) {
				continue
			}
			req, weight, err := ftx.GetFutureStats(symbol)
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

			var fstat ftx.FutureStatsResponse
			err = json.Unmarshal(queryResponse.Response, &fstat)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return nil
			}

			if !fstat.Success {
				state.logger.Info("http error", log.Error(errors.New(fstat.Error)))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return nil
			}
			ts := time.Now().UnixNano() / 1000000
			response.Statistics = append(response.Statistics, &models.Stat{
				Timestamp: utils.MilliToTimestamp(uint64(ts)),
				StatType:  models.OpenInterest,
				Value:     fstat.Result.OpenInterest,
			})
			if state.symbolToSec[symbol].SecurityType == enum.SecurityType_CRYPTO_PERP {
				response.Statistics = append(response.Statistics, &models.Stat{
					Timestamp: utils.MilliToTimestamp(uint64(fstat.Result.NextFundingTime.UnixNano() / 1000000)),
					StatType:  models.FundingRate,
					Value:     fstat.Result.NextFundingRate,
				})
			}
		}
	}

	response.Success = true
	context.Respond(response)

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
			if msg.Filter.Instrument.Symbol != nil {
				symbol = msg.Filter.Instrument.Symbol.Value
			} else if msg.Filter.Instrument.SecurityID != nil {
				sec, ok := state.securities[msg.Filter.Instrument.SecurityID.Value]
				if !ok {
					response.RejectionReason = messages.UnknownSecurityID
					context.Respond(response)
					return nil
				}
				symbol = sec.Symbol
			}
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
	case messages.Deposit:
		request, weight, err = ftx.GetDepositHistory(from, to, msg.Account.ApiCredentials)
	case messages.Withdrawal:
		request, weight, err = ftx.GetWithdrawalHistory(from, to, msg.Account.ApiCredentials)
	case messages.FundingFee:
		request, weight, err = ftx.GetFundingPayments(symbol, from, to, msg.Account.ApiCredentials)
	case messages.RealizedPnl,
		messages.WelcomeBonus,
		messages.Commission:
		response.RejectionReason = messages.UnsupportedRequest
		context.Respond(response)
		return nil
	}
	if err != nil {
		response.RejectionReason = messages.InvalidRequest
		context.Respond(response)
		return nil
	}

	qr := state.getQueryRunner()
	if qr == nil {
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
		var deposits []ftx.Deposit
		var withdrawals []ftx.Withdrawal
		var fundings []ftx.FundingPayment
		switch msg.Type {
		case messages.Deposit:
			res := ftx.DepositsResponse{}
			err = json.Unmarshal(queryResponse.Response, &res)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			if !res.Success {
				state.logger.Info("api error", log.Error(errors.New(res.Error)))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			deposits = res.Result
		case messages.Withdrawal:
			res := ftx.WithdrawalResponse{}
			err = json.Unmarshal(queryResponse.Response, &res)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			if !res.Success {
				state.logger.Info("api error", log.Error(errors.New(res.Error)))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			withdrawals = res.Result
		case messages.FundingFee:
			res := ftx.FundingPaymentsResponse{}
			err = json.Unmarshal(queryResponse.Response, &res)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			if !res.Success {
				state.logger.Info("api error", log.Error(errors.New(res.Error)))
				response.RejectionReason = messages.ExchangeAPIError
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
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			ts, _ := types.TimestampProto(t.Time)
			mvt := messages.AccountMovement{
				Asset:      asset,
				Change:     t.Size,
				MovementID: fmt.Sprintf("%s%d", msg.Account.ApiCredentials.AccountID, t.ID),
				Time:       ts,
				Type:       messages.Deposit,
			}
			movements = append(movements, &mvt)
		}
		for _, t := range withdrawals {
			asset, ok := constants.GetAssetBySymbol(t.Coin)
			if !ok {
				state.logger.Warn("unknown asset " + t.Coin)
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			ts, _ := types.TimestampProto(t.Time)
			mvt := messages.AccountMovement{
				Asset:      asset,
				Change:     -t.Size,
				MovementID: fmt.Sprintf("%s%d", msg.Account.ApiCredentials.AccountID, t.ID),
				Time:       ts,
				Type:       messages.Withdrawal,
			}
			movements = append(movements, &mvt)
		}
		for _, t := range fundings {
			ts, _ := types.TimestampProto(t.Time)
			mvt := messages.AccountMovement{
				Asset:      &constants.DOLLAR,
				Change:     -t.Payment,
				MovementID: fmt.Sprintf("%d", t.ID),
				Time:       ts,
				Type:       messages.FundingFee,
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
	response := &messages.TradeCaptureReport{
		RequestID: msg.RequestID,
		Success:   false,
	}

	symbol := ""
	var from, to *float64
	var orderID *uint64
	if msg.Filter != nil {
		if msg.Filter.Side != nil || msg.Filter.OrderID != nil || msg.Filter.ClientOrderID != nil {
			response.RejectionReason = messages.UnsupportedFilter
			context.Respond(response)
			return nil
		}

		if msg.Filter.Instrument != nil {
			if msg.Filter.Instrument.Symbol != nil {
				symbol = msg.Filter.Instrument.Symbol.Value
			} else if msg.Filter.Instrument.SecurityID != nil {
				sec, ok := state.securities[msg.Filter.Instrument.SecurityID.Value]
				if !ok {
					response.RejectionReason = messages.UnknownSecurityID
					context.Respond(response)
					return nil
				}
				symbol = sec.Symbol
			}
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
				response.RejectionReason = messages.InvalidRequest
				context.Respond(response)
				return nil
			}
			orderID = &v
		}
		if msg.Filter.FromID != nil {
			response.RejectionReason = messages.UnsupportedRequest
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
				fmt.Println(err)
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
			} else if queryResponse.StatusCode >= 500 {

				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.ExchangeAPIError
				fmt.Println(err)
				context.Respond(response)
			}
			return
		}
		var trades ftx.FillsResponse
		err = json.Unmarshal(queryResponse.Response, &trades)
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
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
			if _, ok := state.symbolToSec[t.Market]; ok {
				sec := state.symbolToSec[t.Market]
				instrument = &models.Instrument{
					Exchange:   &constants.FTX,
					Symbol:     &types.StringValue{Value: t.Market},
					SecurityID: &types.UInt64Value{Value: sec.SecurityID},
				}
			} else {
				instrument = &models.Instrument{
					Exchange:   &constants.FTX,
					Symbol:     &types.StringValue{Value: t.BaseCurrency + "-" + t.QuoteCurrency},
					SecurityID: nil,
				}
			}

			quantity := t.Size * quantityMul
			if t.Side == ftx.SELL {
				quantity *= -1
			}
			comAsset, ok := constants.GetAssetBySymbol(t.FeeCurrency)
			if !ok {
				state.logger.Info("api error", log.Error(fmt.Errorf("unknown commission asset %s", t.FeeCurrency)))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			ts, _ := types.TimestampProto(t.Time)
			trd := models.TradeCapture{
				Type:            models.Regular,
				Price:           t.Price,
				Quantity:        quantity,
				Commission:      t.Fee,
				CommissionAsset: comAsset,
				TradeID:         fmt.Sprintf("%d-%d", t.TradeID, t.OrderID),
				Instrument:      instrument,
				Trade_LinkID:    nil,
				OrderID:         &types.StringValue{Value: fmt.Sprintf("%d", t.OrderID)},
				TransactionTime: ts,
			}

			if t.Side == ftx.BUY {
				trd.Side = models.Buy
			} else {
				trd.Side = models.Sell
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
	response := &messages.OrderList{
		RequestID: msg.RequestID,
		Success:   false,
	}
	symbol := ""
	orderID := ""
	clOrderID := ""
	if msg.Filter != nil {
		if msg.Filter.OrderStatus != nil {
			response.RejectionReason = messages.UnsupportedFilter
			context.Respond(response)
			return nil
		}
		if msg.Filter.Side != nil {
			response.RejectionReason = messages.UnsupportedFilter
			context.Respond(response)
			return nil
		}
		if msg.Filter.Instrument != nil {
			if msg.Filter.Instrument.Symbol != nil {
				symbol = msg.Filter.Instrument.Symbol.Value
			} else if msg.Filter.Instrument.SecurityID != nil {
				sec, ok := state.securities[msg.Filter.Instrument.SecurityID.Value]
				if !ok {
					response.RejectionReason = messages.UnknownSecurityID
					context.Respond(response)
					return nil
				}
				symbol = sec.Symbol
			}
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
		var orders []ftx.Order
		if multi {
			var oor ftx.OpenOrdersResponse
			err = json.Unmarshal(queryResponse.Response, &oor)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			if !oor.Success {
				state.logger.Info("API error", log.String("error", oor.Error))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			orders = oor.Result
		} else {
			var oor ftx.OrderStatusResponse
			err = json.Unmarshal(queryResponse.Response, &oor)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			if !oor.Success {
				state.logger.Info("API error", log.String("error", oor.Error))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			orders = []ftx.Order{oor.Result}
		}

		var morders []*models.Order
		for _, o := range orders {
			sec, ok := state.symbolToSec[o.Market]
			if !ok {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			ord := models.Order{
				OrderID:       fmt.Sprintf("%d", o.ID),
				ClientOrderID: o.ClientID,
				Instrument: &models.Instrument{
					Exchange:   &constants.FTX,
					Symbol:     &types.StringValue{Value: o.Market},
					SecurityID: &types.UInt64Value{Value: sec.SecurityID},
				},
				LeavesQuantity: o.Size - o.FilledSize, // TODO check
				CumQuantity:    o.FilledSize,
			}

			if o.ReduceOnly {
				ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ReduceOnly)
			}
			if o.PostOnly {
				ord.ExecutionInstructions = append(ord.ExecutionInstructions, models.ParticipateDoNotInitiate)
			}

			switch o.Status {
			case ftx.NEW_ORDER:
				ord.OrderStatus = models.PendingNew
			case ftx.OPEN_ORDER:
				if o.FilledSize > 0 {
					ord.OrderStatus = models.PartiallyFilled
				} else {
					ord.OrderStatus = models.New
				}
			case ftx.CLOSED_ORDER:
				if o.FilledSize == o.Size {
					ord.OrderStatus = models.Filled
				} else {
					ord.OrderStatus = models.Canceled
				}
			default:
				fmt.Println("unknown ORDER STATUS", o.Status)
			}

			switch o.Type {
			case ftx.LIMIT_ORDER:
				ord.OrderType = models.Limit
			case ftx.MARKET_ORDER:
				ord.OrderType = models.Market
			default:
				fmt.Println("UNKNOWN ORDER TYPE", o.Type)
			}

			switch o.Side {
			case ftx.BUY:
				ord.Side = models.Buy
			case ftx.SELL:
				ord.Side = models.Sell
			default:
				fmt.Println("UNKNOWN ORDER SIDE", o.Side)
			}

			if o.Ioc {
				ord.TimeInForce = models.ImmediateOrCancel
			} else {
				ord.TimeInForce = models.GoodTillCancel
			}
			ord.Price = &types.DoubleValue{Value: o.Price}

			morders = append(morders, &ord)
		}
		response.Success = true
		response.Orders = morders
		context.Respond(response)
	})

	return nil
}

func (state *Executor) OnPositionsRequest(context actor.Context) error {
	msg := context.Message().(*messages.PositionsRequest)

	response := &messages.PositionList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
		Positions:  nil,
	}

	symbol := ""
	if msg.Instrument != nil {
		if msg.Instrument.Symbol != nil {
			symbol = msg.Instrument.Symbol.Value
		} else if msg.Instrument.SecurityID != nil {
			sec, ok := state.securities[msg.Instrument.SecurityID.Value]
			if !ok {
				response.RejectionReason = messages.UnknownSecurityID
				context.Respond(response)
				return nil
			}
			symbol = sec.Symbol
		}
	}

	request, weight, err := ftx.GetPositions(true, msg.Account.ApiCredentials)
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
	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.Other
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
				response.Success = false
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
			}
			return
		}
		var positions ftx.PositionsResponse
		err = json.Unmarshal(queryResponse.Response, &positions)
		if err != nil {
			state.logger.Info("unmarshaling error", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		if !positions.Success {
			state.logger.Info("api error", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
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
			sec, ok := state.symbolToSec[p.Future]
			if !ok {
				err := fmt.Errorf("unknown symbol %s", p.Future)
				state.logger.Info("unmarshaling error", log.Error(err))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			cost := *p.RecentAverageOpenPrice * sec.Multiplier.Value * p.NetSize
			size := p.NetSize
			pos := &models.Position{
				Account: msg.Account.Name,
				Instrument: &models.Instrument{
					Exchange:   &constants.BITMEX,
					Symbol:     &types.StringValue{Value: p.Future},
					SecurityID: &types.UInt64Value{Value: sec.SecurityID},
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
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.rateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
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
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http server error", log.Error(err))
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
			}
			return
		}
		var balances ftx.BalancesResponse
		err = json.Unmarshal(queryResponse.Response, &balances)
		if err != nil {
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		if !balances.Success {
			response.RejectionReason = messages.ExchangeAPIError
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

func buildPlaceOrderRequest(symbol string, order *messages.NewOrder, tickPrecision, lotPrecision int) (ftx.NewOrderRequest, *messages.RejectionReason) {
	request := ftx.NewOrderRequest{
		Market: symbol,
	}

	if order.OrderSide == models.Buy {
		request.Side = ftx.BUY
	} else {
		request.Side = ftx.SELL
	}
	switch order.OrderType {
	case models.Limit:
		request.Type = ftx.LIMIT_ORDER
	case models.Market:
		request.Type = ftx.MARKET_ORDER
	default:
		rej := messages.UnsupportedOrderType
		return request, &rej
	}

	request.Size = order.Quantity
	request.ClientID = order.ClientOrderID

	if order.OrderType != models.Market {
		switch order.TimeInForce {
		case models.ImmediateOrCancel:
			request.Ioc = true
		case models.GoodTillCancel:
			// Skip
		default:
			rej := messages.UnsupportedOrderTimeInForce
			return request, &rej
		}
	}

	if order.Price != nil {
		request.Price = &order.Price.Value
	}

	for _, exec := range order.ExecutionInstructions {
		switch exec {
		case models.ReduceOnly:
			request.ReduceOnly = true
		case models.ParticipateDoNotInitiate:
			request.PostOnly = true
		default:
			rej := messages.UnsupportedOrderCharacteristic
			return request, &rej
		}
	}

	return request, nil
}

func (state *Executor) OnNewOrderSingleRequest(context actor.Context) error {
	req := context.Message().(*messages.NewOrderSingleRequest)
	response := &messages.NewOrderSingleResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	symbol := ""
	var tickPrecision, lotPrecision int
	if req.Order.Instrument != nil {
		if req.Order.Instrument.Symbol != nil {
			symbol = req.Order.Instrument.Symbol.Value
			sec, ok := state.symbolToSec[symbol]
			if !ok {
				response.RejectionReason = messages.UnknownSymbol
				context.Respond(response)
				return nil
			}
			tickPrecision = int(math.Log10(math.Ceil(1. / sec.MinPriceIncrement.Value)))
			lotPrecision = int(math.Log10(math.Ceil(1. / sec.RoundLot.Value)))
		} else if req.Order.Instrument.SecurityID != nil {
			sec, ok := state.securities[req.Order.Instrument.SecurityID.Value]
			if !ok {
				response.RejectionReason = messages.UnknownSecurityID
				context.Respond(response)
				return nil
			}
			symbol = sec.Symbol
			tickPrecision = int(math.Log10(math.Ceil(1. / sec.MinPriceIncrement.Value)))
			lotPrecision = int(math.Log10(math.Ceil(1. / sec.RoundLot.Value)))
		}
	} else {
		response.RejectionReason = messages.UnknownSecurityID
		context.Respond(response)
		return nil
	}

	params, rej := buildPlaceOrderRequest(symbol, req.Order, tickPrecision, lotPrecision)
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
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	if state.orderRateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.orderRateLimit.Request(1)
	qr.rateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.HTTPError
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
			response.RejectionReason = messages.HTTPError
			context.Respond(response)
			return
		}
		var order ftx.NewOrderResponse
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		b, _ := json.Marshal(order)
		fmt.Println("ORDER RESPONSE", string(b))
		if !order.Success {
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
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
		response.RejectionReason = messages.UnknownOrder
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
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}
	if state.orderRateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.orderRateLimit.Request(1)
	qr.rateLimit.Request(weight)

	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.HTTPError
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
			response.RejectionReason = messages.HTTPError
			context.Respond(response)
			return
		}
		var order ftx.ModifyOrderResponse
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		if !order.Success {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
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
			response.RejectionReason = messages.InvalidRequest
			context.Respond(response)
			return nil
		}
	} else if req.ClientOrderID != nil {
		// TODO
	} else {
		response.RejectionReason = messages.InvalidRequest
		context.Respond(response)
		return nil
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.rateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.HTTPError
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
			response.RejectionReason = messages.HTTPError
			context.Respond(response)
			return
		}
		var order ftx.CancelOrderResponse
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		if !order.Success {
			if order.Error == "Order already queued for cancellation" || order.Error == "Order already closed" {
				fmt.Println("order already queued for cancel or already closed", req.OrderID)
				response.RejectionReason = messages.UnknownOrder
				context.Respond(response)
				return
			} else {
				state.logger.Info("error unmarshalling", log.String("error", order.Error))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
		}
		response.Success = true
		context.Respond(response)
	})
	return nil
}
