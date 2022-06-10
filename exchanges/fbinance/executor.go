package fbinance

import (
	goContext "context"
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
	"gitlab.com/alphaticks/xchanger/exchanges/fbinance"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"strconv"
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
	pid             *actor.PID
	globalRateLimit *exchanges.RateLimit
}

type Executor struct {
	extypes.BaseExecutor
	client                *http.Client
	securities            map[uint64]*models.Security
	historicalSecurities  map[uint64]*registry.Security
	symbolToSec           map[string]*models.Security
	symbolToHistoricalSec map[string]*registry.Security
	queryRunners          []*QueryRunner
	secondOrderRateLimit  *exchanges.RateLimit
	minuteOrderRateLimit  *exchanges.RateLimit
	logger                *log.Logger
}

func NewExecutor(config *extypes.ExecutorConfig) actor.Actor {
	ex := &Executor{
		queryRunners: nil,
		logger:       nil,
	}
	ex.ExecutorConfig = config
	return ex
}

func (state *Executor) Receive(context actor.Context) {
	extypes.ReceiveExecutor(state, context)
}

func (state *Executor) getQueryRunner() *QueryRunner {
	sort.Slice(state.queryRunners, func(i, j int) bool {
		return rand.Uint64()%2 == 0
	})

	var qr *QueryRunner
	for _, q := range state.queryRunners {
		if !q.globalRateLimit.IsRateLimited() {
			qr = q
			break
		}
	}

	return qr
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

	dialers := state.ExecutorConfig.DialerPool.GetDialers()
	for _, dialer := range dialers {
		fmt.Println("SETTING UP", dialer.LocalAddr)
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
			pid:             context.Spawn(props),
			globalRateLimit: nil,
		})
	}

	request, weight, err := fbinance.GetExchangeInfo()
	if err != nil {
		return err
	}

	future := context.RequestFuture(state.queryRunners[0].pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
	res, err := future.Result()
	if err != nil {
		return err
	}
	queryResponse := res.(*jobs.PerformQueryResponse)
	if queryResponse.StatusCode != 200 {
		return fmt.Errorf("error getting exchange info: status code %d", queryResponse.StatusCode)
	}

	var exchangeInfo fbinance.ExchangeInfoResponse
	err = json.Unmarshal(queryResponse.Response, &exchangeInfo)
	if err != nil {
		return fmt.Errorf("error decoding query response: %v", err)
	}
	if exchangeInfo.Code != 0 {
		return fmt.Errorf("error getting exchange info: %s", exchangeInfo.Message)
	}

	// Initialize rate limit
	for _, rateLimit := range exchangeInfo.RateLimits {
		if rateLimit.RateLimitType == "ORDERS" {
			if rateLimit.Interval == "MINUTE" {
				state.minuteOrderRateLimit = exchanges.NewRateLimit(rateLimit.Limit, time.Minute)
			} else if rateLimit.Interval == "SECOND" {
				state.secondOrderRateLimit = exchanges.NewRateLimit(rateLimit.Limit, time.Duration(rateLimit.IntervalNum)*time.Second)
			}
		} else if rateLimit.RateLimitType == "REQUEST_WEIGHT" {
			for _, qr := range state.queryRunners {
				qr.globalRateLimit = exchanges.NewRateLimit(rateLimit.Limit, time.Minute)
			}
		}
	}
	if state.queryRunners[0].globalRateLimit == nil || state.minuteOrderRateLimit == nil || state.secondOrderRateLimit == nil {
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
	request, weight, err := fbinance.GetExchangeInfo()
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		return fmt.Errorf("rate limited")
	}

	qr.globalRateLimit.Request(weight)
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
	var exchangeInfo fbinance.ExchangeInfoResponse
	err = json.Unmarshal(resp.Response, &exchangeInfo)
	if err != nil {
		err = fmt.Errorf(
			"error unmarshaling response: %v",
			err)
		return err
	}
	if exchangeInfo.Code != 0 {
		err = fmt.Errorf(
			"fbinance api error: %d %s",
			exchangeInfo.Code,
			exchangeInfo.Message)
		return err
	}

	var securities []*models.Security
	for _, symbol := range exchangeInfo.Symbols {
		baseCurrency, ok := constants.GetAssetBySymbol(symbol.BaseAsset)
		if !ok {
			state.logger.Info(fmt.Sprintf("unknown currency %s", symbol.BaseAsset))
			continue
		}
		quoteCurrency, ok := constants.GetAssetBySymbol(symbol.QuoteAsset)
		if !ok {
			continue
		}
		security := models.Security{}
		security.Symbol = symbol.Symbol
		security.Underlying = baseCurrency
		security.QuoteCurrency = quoteCurrency
		switch symbol.Status {
		case "PRE_TRADING":
			security.Status = models.InstrumentStatus_PreTrading
		case "TRADING":
			security.Status = models.InstrumentStatus_Trading
		case "POST_TRADING":
			security.Status = models.InstrumentStatus_PostTrading
		case "END_OF_DAY":
			security.Status = models.InstrumentStatus_EndOfDay
		case "HALT":
			security.Status = models.InstrumentStatus_Halt
		case "AUCTION_MATCH":
			security.Status = models.InstrumentStatus_AuctionMatch
		case "BREAK":
			security.Status = models.InstrumentStatus_Break
		default:
			security.Status = models.InstrumentStatus_Disabled
		}
		security.Exchange = constants.FBINANCE
		switch symbol.ContractType {
		case "PERPETUAL":
			security.SecurityType = enum.SecurityType_CRYPTO_PERP
		default:
			continue
		}
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		for _, f := range symbol.Filters {
			if f.FilterType == "PRICE_FILTER" {
				security.MinPriceIncrement = &wrapperspb.DoubleValue{Value: symbol.Filters[0].TickSize}
			}
		}
		if security.MinPriceIncrement == nil {
			fmt.Println("NO MIN PRICE INCREMENT", symbol.Symbol)
			continue
		}
		security.RoundLot = &wrapperspb.DoubleValue{Value: 1. / math.Pow10(symbol.QuantityPrecision)}
		security.IsInverse = false
		security.Multiplier = &wrapperspb.DoubleValue{Value: 1.}
		// Default fee
		security.MakerFee = &wrapperspb.DoubleValue{Value: 0.0002}
		security.TakerFee = &wrapperspb.DoubleValue{Value: 0.0004}
		for _, filter := range symbol.Filters {
			switch filter.FilterType {
			case fbinance.LOT_SIZE:
				security.MinLimitQuantity = &wrapperspb.DoubleValue{Value: filter.MinQty}
				security.MaxLimitQuantity = &wrapperspb.DoubleValue{Value: filter.MaxQty}
			case fbinance.MARKET_LOT_SIZE:
				security.MinMarketQuantity = &wrapperspb.DoubleValue{Value: filter.MinQty}
				security.MaxMarketQuantity = &wrapperspb.DoubleValue{Value: filter.MaxQty}
			}
		}
		securities = append(securities, &security)
	}

	state.securities = make(map[uint64]*models.Security)
	state.symbolToSec = make(map[string]*models.Security)
	for _, s := range securities {
		state.securities[s.SecurityID] = s
		state.symbolToSec[s.Symbol] = s
	}

	state.historicalSecurities = make(map[uint64]*registry.Security)
	state.symbolToHistoricalSec = make(map[string]*registry.Security)
	if state.Registry != nil {
		rres, err := state.Registry.Securities(goContext.Background(), &registry.SecuritiesRequest{
			Filter: &registry.SecurityFilter{
				ExchangeId: []uint32{constants.FBINANCE.ID},
			},
		})
		if err != nil {
			return fmt.Errorf("error fetching historical securities: %v", err)
		}
		for _, sec := range rres.Securities {
			state.historicalSecurities[sec.SecurityId] = sec
			state.symbolToHistoricalSec[sec.Symbol] = sec
		}
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
		response.RejectionReason = messages.RejectionReason_MissingInstrument
		context.Respond(response)
		return nil
	}
	symbol := msg.Instrument.Symbol.Value
	for _, stat := range msg.Statistics {
		switch stat {
		case models.StatType_OpenInterest:
			req, weight, err := fbinance.GetOpenInterest(symbol)
			if err != nil {
				return err
			}

			qr := state.getQueryRunner()

			if qr == nil {
				response.RejectionReason = messages.RejectionReason_RateLimitExceeded
				context.Respond(response)
				return nil
			}

			qr.globalRateLimit.Request(weight)

			res, err := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: req}, 10*time.Second).Result()
			if err != nil {
				state.logger.Warn("http client error", log.Error(err))
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

			var oires fbinance.OpenInterestResponse
			err = json.Unmarshal(queryResponse.Response, &oires)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return nil
			}

			if oires.Code != 0 {
				state.logger.Info("http error", log.Error(errors.New(oires.Message)))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return nil
			}
			response.Statistics = append(response.Statistics, &models.Stat{
				Timestamp: utils.MilliToTimestamp(uint64(oires.Time)),
				StatType:  models.StatType_OpenInterest,
				Value:     oires.OpenInterest,
			})
		}
	}

	response.Success = true
	context.Respond(response)

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
		response.RejectionReason = messages.RejectionReason_UnsupportedSubscription
		context.Respond(response)
		return nil
	}
	if msg.Instrument == nil || msg.Instrument.Symbol == nil {
		response.RejectionReason = messages.RejectionReason_MissingInstrument
		context.Respond(response)
		return nil
	}
	symbol := msg.Instrument.Symbol.Value
	// Get http request and the expected response
	request, weight, err := fbinance.GetOrderBook(symbol, 1000)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()

	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.globalRateLimit.Request(weight)

	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

	context.ReenterAfter(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Info("http client error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
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
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
				return
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"http server error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
				return
			}
			return
		}
		var obData fbinance.OrderBookResponse
		err = json.Unmarshal(queryResponse.Response, &obData)
		if err != nil {
			err = fmt.Errorf("error decoding query response: %v", err)
			state.logger.Info("http client error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if obData.Code != 0 {
			err = fmt.Errorf("error getting orderbook: %d %s", obData.Code, obData.Message)
			state.logger.Info("http client error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		bids, asks, err := obData.ToBidAsk()
		if err != nil {
			err = fmt.Errorf("error converting orderbook: %v", err)
			state.logger.Info("http client error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		snapshot = &models.OBL2Snapshot{
			Bids:      bids,
			Asks:      asks,
			Timestamp: timestamppb.Now(),
		}
		response.Success = true
		response.SnapshotL2 = snapshot
		response.SeqNum = obData.LastUpdateID
		context.Respond(response)
	})

	return nil
}

func (state *Executor) OnAccountMovementRequest(context actor.Context) error {
	fmt.Println("ON TRADE ACCOUNT MOVEMENT REQUEST !!!!")
	msg := context.Message().(*messages.AccountMovementRequest)
	response := &messages.AccountMovementResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	params := fbinance.NewIncomeHistoryRequest()
	if msg.Filter != nil {
		if msg.Filter.Instrument != nil {
			if msg.Filter.Instrument.Symbol != nil {
				params.SetSymbol(msg.Filter.Instrument.Symbol.Value)
			} else if msg.Filter.Instrument.SecurityID != nil {
				sec, ok := state.securities[msg.Filter.Instrument.SecurityID.Value]
				if !ok {
					response.RejectionReason = messages.RejectionReason_UnknownSecurityID
					context.Respond(response)
					return nil
				}
				params.SetSymbol(sec.Symbol)
			}
		}

		if msg.Filter.From != nil {
			ms := uint64(msg.Filter.From.Seconds*1000) + uint64(msg.Filter.From.Nanos/1000000)
			params.SetFrom(ms)
			fmt.Println("SET FROM", ms)
		}
		if msg.Filter.To != nil {
			ms := uint64(msg.Filter.To.Seconds*1000) + uint64(msg.Filter.To.Nanos/1000000)
			params.SetTo(ms)
		}
	}

	params.SetLimit(1000)

	switch msg.Type {
	case messages.AccountMovementType_Commission:
		params.SetIncomeType(fbinance.COMMISSION)
	case messages.AccountMovementType_Deposit:
		params.SetIncomeType(fbinance.TRANSFER)
	case messages.AccountMovementType_Withdrawal:
		params.SetIncomeType(fbinance.TRANSFER)
	case messages.AccountMovementType_FundingFee:
		params.SetIncomeType(fbinance.FUNDING_FEE)
	case messages.AccountMovementType_RealizedPnl:
		params.SetIncomeType(fbinance.REALIZED_PNL)
	case messages.AccountMovementType_WelcomeBonus:
		params.SetIncomeType(fbinance.WELCOME_BONUS)
	}

	request, weight, err := fbinance.GetIncomeHistory(params, msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.globalRateLimit.Request(weight)
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
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			} else if queryResponse.StatusCode >= 500 {

				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			}
			return
		}
		var incomes []fbinance.Income
		err = json.Unmarshal(queryResponse.Response, &incomes)
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		var movements []*messages.AccountMovement
		for _, t := range incomes {
			if msg.Type == messages.AccountMovementType_Deposit && t.Income < 0 {
				continue
			}
			if msg.Type == messages.AccountMovementType_Withdrawal && t.Income > 0 {
				continue
			}
			asset, ok := constants.GetAssetBySymbol(t.Asset)
			if !ok {
				state.logger.Warn("unknown asset " + t.Asset)
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			mvt := messages.AccountMovement{
				Asset:      asset,
				Change:     t.Income,
				MovementID: fmt.Sprintf("%s%s", string(t.IncomeType), t.TransferID),
				Time:       utils.MilliToTimestamp(t.Time),
			}
			switch t.IncomeType {
			case fbinance.FUNDING_FEE:
				mvt.Type = messages.AccountMovementType_FundingFee
				mvt.Subtype = t.Symbol
			case fbinance.WELCOME_BONUS:
				mvt.Type = messages.AccountMovementType_WelcomeBonus
			case fbinance.COMMISSION:
				mvt.Type = messages.AccountMovementType_Commission
			case fbinance.TRANSFER:
				if mvt.Change > 0 {
					mvt.Type = messages.AccountMovementType_Deposit
				} else {
					mvt.Type = messages.AccountMovementType_Withdrawal
				}
			case fbinance.REALIZED_PNL:
				mvt.Type = messages.AccountMovementType_RealizedPnl
			default:
				state.logger.Warn("unknown income type " + string(t.IncomeType))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
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
	fmt.Println("ON TRADE CAPTURE REPORT REQUEST !!!!")
	msg := context.Message().(*messages.TradeCaptureReportRequest)
	response := &messages.TradeCaptureReport{
		RequestID: msg.RequestID,
		Success:   false,
	}

	symbol := ""
	var from, to *uint64
	var fromID string
	if msg.Filter != nil {
		if msg.Filter.Side != nil || msg.Filter.OrderID != nil || msg.Filter.ClientOrderID != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedFilter
			context.Respond(response)
			return nil
		}

		if msg.Filter.Instrument != nil {
			if msg.Filter.Instrument.Symbol != nil {
				symbol = msg.Filter.Instrument.Symbol.Value
			} else if msg.Filter.Instrument.SecurityID != nil {
				sec, ok := state.securities[msg.Filter.Instrument.SecurityID.Value]
				if !ok {
					response.RejectionReason = messages.RejectionReason_UnknownSecurityID
					context.Respond(response)
					return nil
				}
				symbol = sec.Symbol
			}
		}

		if msg.Filter.From != nil {
			ms := uint64(msg.Filter.From.Seconds*1000) + uint64(msg.Filter.From.Nanos/1000000)
			from = &ms
		}
		if msg.Filter.To != nil {
			ms := uint64(msg.Filter.To.Seconds*1000) + uint64(msg.Filter.To.Nanos/1000000)
			to = &ms
		}
		if msg.Filter.FromID != nil {
			fromID = msg.Filter.FromID.Value
		}
	}
	params := fbinance.NewUserTradesRequest(symbol)

	// If from is not set, but to is set,
	// If from is set, but to is not set, ok

	if fromID != "" {
		fromIDInt, _ := strconv.ParseInt(fromID, 10, 64)
		params.SetFromID(int(fromIDInt))
	} else {
		if from == nil || *from == 0 {
			params.SetFromID(0)
		} else {
			params.SetFrom(*from)
			if to != nil {
				if *to-*from > (7 * 24 * 60 * 60 * 1000) {
					*to = *from + (7 * 24 * 60 * 60 * 1000)
				}
				params.SetTo(*to)
			}
		}
	}

	request, weight, err := fbinance.GetUserTrades(params, msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.globalRateLimit.Request(weight)
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
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			} else if queryResponse.StatusCode >= 500 {

				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			}
			return
		}
		var trades []fbinance.UserTrade
		err = json.Unmarshal(queryResponse.Response, &trades)
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		var mtrades []*models.TradeCapture
		for _, t := range trades {
			sec, ok := state.symbolToHistoricalSec[t.Symbol]
			if !ok {
				state.logger.Info("unknown symbol", log.String("symbol", t.Symbol))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			quantity := t.Quantity
			if t.Side == fbinance.SELL_ODER {
				quantity *= -1
			}
			trd := models.TradeCapture{
				Type:       models.TradeType_Regular,
				Price:      t.Price,
				Quantity:   quantity,
				Commission: t.Commission,
				TradeID:    fmt.Sprintf("%d-%d", t.TradeID, t.OrderID),
				Instrument: &models.Instrument{
					Exchange:   constants.FBINANCE,
					Symbol:     &wrapperspb.StringValue{Value: t.Symbol},
					SecurityID: &wrapperspb.UInt64Value{Value: sec.SecurityId},
				},
				Trade_LinkID:    nil,
				OrderID:         &wrapperspb.StringValue{Value: fmt.Sprintf("%d", t.OrderID)},
				TransactionTime: utils.MilliToTimestamp(t.Timestamp),
			}

			if t.Side == fbinance.BUY_ORDER {
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
	response := &messages.OrderList{
		RequestID: msg.RequestID,
		Success:   false,
	}
	symbol := ""
	orderID := ""
	clOrderID := ""
	if msg.Filter != nil {
		if msg.Filter.OrderStatus != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedFilter
			context.Respond(response)
			return nil
		}
		if msg.Filter.Side != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedFilter
			context.Respond(response)
			return nil
		}
		if msg.Filter.Instrument != nil {
			if msg.Filter.Instrument.Symbol != nil {
				symbol = msg.Filter.Instrument.Symbol.Value
			} else if msg.Filter.Instrument.SecurityID != nil {
				sec, ok := state.securities[msg.Filter.Instrument.SecurityID.Value]
				if !ok {
					response.RejectionReason = messages.RejectionReason_UnknownSecurityID
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
	if orderID != "" || clOrderID != "" {
		params := fbinance.NewQueryOrderRequest(symbol)
		if orderID != "" {
			orderIDInt, err := strconv.ParseInt(orderID, 10, 64)
			if err != nil {
				response.RejectionReason = messages.RejectionReason_UnsupportedFilter
				context.Respond(response)
				return nil
			}
			params.SetOrderID(orderIDInt)
		}
		if clOrderID != "" {
			params.SetOrigClientOrderID(clOrderID)
		}
		var err error
		request, weight, err = fbinance.QueryOrder(params, msg.Account.ApiCredentials)
		if err != nil {
			return err
		}
	} else {
		var err error
		request, weight, err = fbinance.QueryOpenOrders(symbol, msg.Account.ApiCredentials)
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

	qr.globalRateLimit.Request(weight)
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
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			} else if queryResponse.StatusCode >= 500 {

				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			}
			return
		}
		var orders []fbinance.OrderData
		err = json.Unmarshal(queryResponse.Response, &orders)
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		var morders []*models.Order
		for _, o := range orders {
			sec, ok := state.symbolToSec[o.Symbol]
			if !ok {
				response.RejectionReason = messages.RejectionReason_UnknownSymbol
				context.Respond(response)
				return
			}
			ord := orderToModel(&o)
			ord.Instrument.SecurityID = &wrapperspb.UInt64Value{Value: sec.SecurityID}
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
				response.RejectionReason = messages.RejectionReason_UnknownSecurityID
				context.Respond(response)
				return nil
			}
			symbol = sec.Symbol
		}
	}

	request, weight, err := fbinance.GetPositionRisk(msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.globalRateLimit.Request(weight)
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
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			} else if queryResponse.StatusCode >= 500 {

				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http error", log.Error(err))
				response.Success = false
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			}
			return
		}
		var positions []fbinance.AccountPositionRisk
		err = json.Unmarshal(queryResponse.Response, &positions)
		if err != nil {
			state.logger.Info("unmarshalling error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		for _, p := range positions {
			if p.PositionAmount == 0 {
				continue
			}
			if symbol != "" && p.Symbol != symbol {
				continue
			}
			sec, ok := state.symbolToSec[p.Symbol]
			if !ok {
				err := fmt.Errorf("unknown symbol %s", p.Symbol)
				state.logger.Info("unknown symbol error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			var cost float64
			if sec.IsInverse {
				cost = ((1. / p.MarkPrice) * sec.Multiplier.Value * p.PositionAmount) - p.UnrealizedProfit
			} else {
				cost = (p.MarkPrice * sec.Multiplier.Value * p.PositionAmount) - p.UnrealizedProfit
			}
			pos := &models.Position{
				Account: msg.Account.Name,
				Instrument: &models.Instrument{
					Exchange:   constants.FBINANCE,
					Symbol:     &wrapperspb.StringValue{Value: p.Symbol},
					SecurityID: &wrapperspb.UInt64Value{Value: sec.SecurityID},
				},
				Quantity: p.PositionAmount,
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

	request, weight, err := fbinance.GetBalance(msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.globalRateLimit.Request(weight)
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
		var balances []fbinance.AccountBalance
		err = json.Unmarshal(queryResponse.Response, &balances)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		for _, b := range balances {
			if b.Balance == 0. {
				continue
			}
			asset, ok := constants.GetAssetBySymbol(b.Asset)
			if !ok {
				state.logger.Error("got balance for unknown asset", log.String("asset", b.Asset))
				continue
			}
			response.Balances = append(response.Balances, &models.Balance{
				Account:  msg.Account.Name,
				Asset:    asset,
				Quantity: b.Balance,
			})
			fmt.Println(response.Balances)
		}

		response.Success = true
		context.Respond(response)
	})

	return nil
}

func (state *Executor) OnNewOrderSingleRequest(context actor.Context) error {
	req := context.Message().(*messages.NewOrderSingleRequest)
	response := &messages.NewOrderSingleResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	if state.secondOrderRateLimit.IsRateLimited() {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	if state.minuteOrderRateLimit.IsRateLimited() {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	symbol := ""
	var tickPrecision, lotPrecision int
	if req.Order.Instrument != nil {
		if req.Order.Instrument.Symbol != nil {
			symbol = req.Order.Instrument.Symbol.Value
			sec, ok := state.symbolToSec[symbol]
			if !ok {
				response.RejectionReason = messages.RejectionReason_UnknownSecurityID
				context.Respond(response)
				return nil
			}
			tickPrecision = int(math.Ceil(math.Log10(1. / sec.MinPriceIncrement.Value)))
			lotPrecision = int(math.Ceil(math.Log10(1. / sec.RoundLot.Value)))
		} else if req.Order.Instrument.SecurityID != nil {
			sec, ok := state.securities[req.Order.Instrument.SecurityID.Value]
			if !ok {
				response.RejectionReason = messages.RejectionReason_UnknownSecurityID
				context.Respond(response)
				return nil
			}
			symbol = sec.Symbol
			tickPrecision = int(math.Ceil(math.Log10(1. / sec.MinPriceIncrement.Value)))
			lotPrecision = int(math.Ceil(math.Log10(1. / sec.RoundLot.Value)))
		}
	} else {
		response.RejectionReason = messages.RejectionReason_UnknownSecurityID
		context.Respond(response)
		return nil
	}

	params, rej := buildPostOrderRequest(symbol, req.Order, tickPrecision, lotPrecision)
	if rej != nil {
		response.RejectionReason = *rej
		context.Respond(response)
		return nil
	}

	request, weight, err := fbinance.NewOrder(params, req.Account.ApiCredentials)
	if err != nil {
		return err
	}

	state.secondOrderRateLimit.Request(1)
	state.minuteOrderRateLimit.Request(1)
	qr.globalRateLimit.Request(weight)
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
		var order fbinance.OrderData
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
		response.OrderID = fmt.Sprintf("%d", order.OrderID)
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
	symbol := ""
	if req.Instrument != nil {
		if req.Instrument.Symbol != nil {
			symbol = req.Instrument.Symbol.Value
		} else if req.Instrument.SecurityID != nil {
			sec, ok := state.securities[req.Instrument.SecurityID.Value]
			if !ok {
				response.RejectionReason = messages.RejectionReason_UnknownSecurityID
				context.Respond(response)
				return nil
			}
			symbol = sec.Symbol
		}
	} else {
		response.RejectionReason = messages.RejectionReason_UnknownSecurityID
		context.Respond(response)
		return nil
	}
	params := fbinance.NewQueryOrderRequest(symbol)
	if req.OrderID != nil {
		orderIDInt, err := strconv.ParseInt(req.OrderID.Value, 10, 64)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_UnknownOrder
			context.Respond(response)
			return nil
		}
		params.SetOrderID(orderIDInt)
	} else if req.ClientOrderID != nil {
		params.SetOrigClientOrderID(req.ClientOrderID.Value)
	} else {
		response.RejectionReason = messages.RejectionReason_UnknownOrder
		context.Respond(response)
		return nil
	}

	request, weight, err := fbinance.CancelOrder(params, req.Account.ApiCredentials)
	if err != nil {
		return err
	}

	fmt.Println(request.URL)
	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.globalRateLimit.Request(weight)
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
				response.RejectionReason = messages.RejectionReason_HTTPError

				var res fbinance.BaseResponse
				if err := json.Unmarshal(queryResponse.Response, &res); err != nil {
					err := fmt.Errorf(
						"%d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
					state.logger.Info("http client error", log.Error(err))
					context.Respond(response)
					return
				}
				if res.Code == -1020 {
					response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				} else if res.Code < -1100 && res.Code > -1199 {
					response.RejectionReason = messages.RejectionReason_InvalidRequest
				} else if res.Code == -2011 && res.Message == "Unknown order sent." {
					response.RejectionReason = messages.RejectionReason_UnknownOrder
				}

				err := fmt.Errorf(
					"%d %s",
					res.Code,
					res.Message)
				state.logger.Info("http client error", log.Error(err))
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
		var order fbinance.OrderData
		fmt.Println(string(queryResponse.Response))
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
		context.Respond(response)
	})
	return nil
}

func (state *Executor) OnOrderMassCancelRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderMassCancelRequest)
	response := &messages.OrderMassCancelResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	symbol := ""
	if req.Filter != nil {
		if req.Filter.Instrument != nil {
			if req.Filter.Instrument.Symbol != nil {
				if _, ok := state.symbolToSec[req.Filter.Instrument.Symbol.Value]; !ok {
					response.RejectionReason = messages.RejectionReason_UnknownSymbol
					context.Respond(response)
					return nil
				}
				symbol = req.Filter.Instrument.Symbol.Value
			} else if req.Filter.Instrument.SecurityID != nil {
				sec, ok := state.securities[req.Filter.Instrument.SecurityID.Value]
				if !ok {
					fmt.Println(state.securities, req.Filter.Instrument.SecurityID.Value)
					response.RejectionReason = messages.RejectionReason_UnknownSecurityID
					context.Respond(response)
					return nil
				}
				symbol = sec.Symbol
			}
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

	request, weight, err := fbinance.CancelAllOrders(symbol, req.Account.ApiCredentials)
	if err != nil {
		return err
	}
	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.globalRateLimit.Request(weight)
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
		var fres fbinance.BaseResponse
		err = json.Unmarshal(queryResponse.Response, &fres)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if fres.Code != 200 {
			state.logger.Info("error unmarshalling", log.Error(errors.New(fres.Message)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
		context.Respond(response)
	})
	return nil
}
