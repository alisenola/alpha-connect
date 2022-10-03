package fbinance

import (
	goContext "context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/config"
	"gitlab.com/alphaticks/alpha-connect/enum"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/fbinance"
	xutils "gitlab.com/alphaticks/xchanger/utils"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
	"net"
	"net/http"
	"reflect"
	"regexp"
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

var MakerFees = map[int]float64{
	0: 0.0002,
	1: 0.00016,
	2: 0.00014,
	3: 0.00012,
	4: 0.0001,
	5: 0.00008,
	6: 0.00006,
	7: 0.00004,
	8: 0.00002,
	9: 0,
}
var TakerFees = map[int]float64{
	0: 0.0004,
	1: 0.0004,
	2: 0.00035,
	3: 0.00032,
	4: 0.0003,
	5: 0.00027,
	6: 0.00025,
	7: 0.00022,
	8: 0.0002,
	9: 0.00017,
}

type AccountRateLimit struct {
	second *exchanges.RateLimit
	minute *exchanges.RateLimit
}

func NewAccountRateLimit(second, minute *exchanges.RateLimit) *AccountRateLimit {
	return &AccountRateLimit{
		second: second,
		minute: minute,
	}
}

func (rl *AccountRateLimit) Request() {
	rl.second.Request(1)
	rl.minute.Request(1)
}

func (rl *AccountRateLimit) IsRateLimited() bool {
	return rl.second.IsRateLimited() || rl.minute.IsRateLimited()
}

func (rl *AccountRateLimit) DurationBeforeNextRequest(weight int) time.Duration {
	dur1 := rl.second.DurationBeforeNextRequest(weight)
	dur2 := rl.minute.DurationBeforeNextRequest(weight)
	if dur1 > dur2 {
		return dur1
	} else {
		return dur2
	}
}

type QueryRunner struct {
	client          *http.Client
	globalRateLimit *exchanges.RateLimit
	endpoint        string
}

type Executor struct {
	extypes.BaseExecutor
	accountRateLimits   map[string]*AccountRateLimit
	newAccountRateLimit func() *AccountRateLimit
	queryRunners        []*QueryRunner
	logger              *log.Logger
	config              *config.Config
}

func NewExecutor(config *config.Config, dialerPool *xutils.DialerPool, registry registry.StaticClient, accountClients map[string]*http.Client) actor.Actor {
	ex := &Executor{
		queryRunners: nil,
		logger:       nil,
		config:       config,
	}
	ex.DialerPool = dialerPool
	ex.Registry = registry
	ex.AccountClients = accountClients
	return ex
}

func (state *Executor) Receive(context actor.Context) {
	extypes.ReceiveExecutor(state, context)
}

func (state *Executor) getQueryRunner(force bool) *QueryRunner {
	var qr *QueryRunner
	for _, q := range state.queryRunners {
		if !q.globalRateLimit.IsRateLimited() {
			qr = q
			break
		}
	}
	if qr == nil && force {
		min := time.Duration(math.MaxInt64)
		for _, q := range state.queryRunners {
			dur := q.globalRateLimit.DurationBeforeNextRequest(1)
			if dur < min {
				min = dur
				qr = q
			}
		}
	}

	return qr
}

func (state *Executor) durationBeforeNextRequest(weight int) time.Duration {
	var minDur time.Duration
	for _, q := range state.queryRunners {
		dur := q.globalRateLimit.DurationBeforeNextRequest(weight)
		if dur < minDur {
			minDur = dur
		}
	}

	return minDur
}

func (state *Executor) updateRateLimits(res fbinance.BaseResponse, ar *AccountRateLimit, qr *QueryRunner) {
	if ar != nil && res.Code == -1015 {
		// Order rate limit, find
		re := regexp.MustCompile(`Too many new orders; current limit is (\d+) orders per ([A-Z_]+)\.`)
		match := re.FindStringSubmatch(res.Message)
		if len(match) == 3 {
			switch match[2] {
			case "MINUTE":
				limit, err := strconv.ParseInt(match[1], 10, 64)
				if err != nil {
					state.logger.Warn("error parsing rate limit " + match[1])
				} else {
					state.logger.Info(fmt.Sprintf("updated %s rate limit to %d", match[2], limit))
					ar.minute.SetLimit(int(limit))
				}
			case "TEN_SECONDS":
				limit, err := strconv.ParseInt(match[1], 10, 64)
				if err != nil {
					state.logger.Warn("error parsing rate limit " + match[1])
				} else {
					state.logger.Info(fmt.Sprintf("updated %s rate limit to %d", match[2], limit))
					ar.second.SetLimit(int(limit))
				}
			default:
				state.logger.Warn(fmt.Sprintf("error matching message: %v", match))
			}
		} else {
			state.logger.Warn(fmt.Sprintf("error matching message: %v", match))
		}
	} else if qr != nil && res.Code == -1003 {
		re := regexp.MustCompile(`.* is (\d+) requests per minute.*`)
		match := re.FindStringSubmatch(res.Message)
		if len(match) == 2 {
			limit, err := strconv.ParseInt(match[1], 10, 64)
			if err != nil {
				state.logger.Warn("error parsing rate limit " + match[1])
			} else {
				state.logger.Info(fmt.Sprintf("updated global rate limit to %d", limit))
				qr.globalRateLimit.SetLimit(int(limit))
			}
		}
	}
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

	dialers := state.DialerPool.GetDialers()
	for _, dialer := range dialers {
		fmt.Println("SETTING UP", dialer.LocalAddr)
		d := dialer
		client := &http.Client{
			Transport: &http2.Transport{
				DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
					fmt.Println("DIALER !!", d.LocalAddr)
					return tls.DialWithDialer(d, network, addr, cfg)
				},
				/*
					MaxIdleConnsPerHost: 1024,
					TLSHandshakeTimeout: 10 * time.Second,
					DialContext:         dialer.DialContext,

				*/
			},
			Timeout: 10 * time.Second,
		}
		qr := &QueryRunner{
			client:          client,
			globalRateLimit: nil,
		}
		addr := strings.Split(d.LocalAddr.String(), ":")[0]
		for _, str := range state.config.FBinanceWhitelistedIPs {
			splits := strings.Split(str, "=>")
			fmt.Println("SPLITS", splits)
			if len(splits) == 2 && splits[0] == addr {
				fmt.Println("PRIVATE", addr, splits[1])
				qr.endpoint = splits[1]
			}
		}
		state.queryRunners = append(state.queryRunners, qr)
	}

	state.accountRateLimits = make(map[string]*AccountRateLimit)
	qr := state.queryRunners[0]
	//fbinance.SelectAddr(qr.client)
	request, weight, err := fbinance.GetExchangeInfo()
	if err != nil {
		return err
	}

	var data fbinance.ExchangeInfoResponse
	if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
		err := fmt.Errorf("error getting exchange info: %v", err)
		return err
	}
	if data.Code != 0 {
		err := fmt.Errorf("error getting exchange info: %v", errors.New(data.Message))
		return err
	}

	// Initialize rate limit
	var secondOrderInterval, minuteOrderInterval time.Duration
	var secondOrderLimit, minuteOrderLimit int
	for _, rateLimit := range data.RateLimits {
		if rateLimit.RateLimitType == "ORDERS" {
			if rateLimit.Interval == "MINUTE" {
				minuteOrderInterval = time.Duration(rateLimit.IntervalNum) * time.Minute
				minuteOrderLimit = rateLimit.Limit * 10 // Adaptive rate limit, wait for an error to update
			} else if rateLimit.Interval == "SECOND" {
				secondOrderInterval = time.Duration(rateLimit.IntervalNum) * time.Second
				secondOrderLimit = rateLimit.Limit * 10 // Adaptive rate limit, wait for an error to update
			}
		} else if rateLimit.RateLimitType == "REQUEST_WEIGHT" {
			for _, q := range state.queryRunners {
				q.globalRateLimit = exchanges.NewRateLimit(int(float64(rateLimit.Limit)*0.9), time.Minute)
				// Update rate limit with weight from the current exchange info fetch
				q.globalRateLimit.Request(weight * 10)
			}
		}
	}
	state.newAccountRateLimit = func() *AccountRateLimit {
		return NewAccountRateLimit(exchanges.NewRateLimit(secondOrderLimit, secondOrderInterval), exchanges.NewRateLimit(minuteOrderLimit, minuteOrderInterval))
	}
	if qr.globalRateLimit == nil {
		return fmt.Errorf("unable to set rate limit")
	}

	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	qr := state.getQueryRunner(false)
	if qr == nil {
		return fmt.Errorf("rate limited")
	}

	request, weight, err := fbinance.GetExchangeInfo()
	if err != nil {
		return err
	}

	qr.globalRateLimit.Request(weight)
	var data fbinance.ExchangeInfoResponse

	if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
		err := fmt.Errorf("error updating security list: %v", err)
		return err
	}
	if data.Code != 0 {
		err = fmt.Errorf(
			"fbinance api error: %d %s",
			data.Code,
			data.Message)
		return err
	}

	var securities []*models.Security
	for _, symbol := range data.Symbols {
		baseCurrency, ok := constants.GetAssetBySymbol(symbol.BaseAsset)
		if !ok {
			state.logger.Debug(fmt.Sprintf("unknown currency %s", symbol.BaseAsset))
			continue
		}
		quoteCurrency, ok := constants.GetAssetBySymbol(symbol.QuoteAsset)
		if !ok {
			state.logger.Debug(fmt.Sprintf("unknown currency %s", symbol.QuoteAsset))
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

	var historicalSecurities []*registry.Security
	if state.Registry != nil {
		rres, err := state.Registry.Securities(goContext.Background(), &registry.SecuritiesRequest{
			Filter: &registry.SecurityFilter{
				ExchangeId: []uint32{constants.FBINANCE.ID},
			},
		})
		if err != nil {
			return fmt.Errorf("error fetching historical securities: %v", err)
		}
		historicalSecurities = rres.Securities
	}

	state.SyncSecurities(securities, historicalSecurities)

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: securities})

	return nil
}

func (state *Executor) OnMarketStatisticsRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketStatisticsRequest)
	sender := context.Sender()
	response := &messages.MarketStatisticsResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	go func() {
		symbol, rej := state.InstrumentToSymbol(msg.Instrument)
		if rej != nil {
			response.RejectionReason = *rej
			context.Send(sender, response)
			return
		}
		for _, stat := range msg.Statistics {
			switch stat {
			case models.StatType_OpenInterest:
				request, weight, err := fbinance.GetOpenInterest(symbol)
				if err != nil {
					response.RejectionReason = messages.RejectionReason_UnsupportedRequest
					context.Send(sender, response)
					return
				}

				qr := state.getQueryRunner(false)
				if qr == nil {
					response.RejectionReason = messages.RejectionReason_RateLimitExceeded
					context.Send(sender, response)
					return
				}

				qr.globalRateLimit.Request(weight)

				var data fbinance.OpenInterestResponse
				err = xutils.PerformJSONRequest(qr.client, request, &data)
				if err != nil {
					state.logger.Warn("error fetching open interests", log.Error(err))
					response.RejectionReason = messages.RejectionReason_HTTPError
					context.Send(sender, response)
					return
				}
				if data.Code != 0 {
					state.logger.Warn("error fetching open interests", log.Error(errors.New(data.Message)))
					response.RejectionReason = messages.RejectionReason_ExchangeAPIError
					context.Send(sender, response)
					state.updateRateLimits(data.BaseResponse, nil, qr)
					return
				}
				response.Statistics = append(response.Statistics, &models.Stat{
					Timestamp: utils.MilliToTimestamp(uint64(data.Time)),
					StatType:  models.StatType_OpenInterest,
					Value:     data.OpenInterest,
				})
			}
		}
		response.Success = true
		context.Send(sender, response)
	}()

	return nil
}

func (state *Executor) OnMarketDataRequest(context actor.Context) error {
	var snapshot *models.OBL2Snapshot
	msg := context.Message().(*messages.MarketDataRequest)
	sender := context.Sender()
	response := &messages.MarketDataResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	if msg.Subscribe {
		response.RejectionReason = messages.RejectionReason_UnsupportedSubscription
		context.Send(sender, response)
		return nil
	}
	symbol, rej := state.InstrumentToSymbol(msg.Instrument)
	if rej != nil {
		response.RejectionReason = *rej
		context.Send(sender, response)
		return nil
	}
	go func() {
		// Get http request and the expected response
		request, weight, err := fbinance.GetOrderBook(symbol, 1000)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
			return
		}

		qr := state.getQueryRunner(false)

		if qr == nil {
			response.RejectionReason = messages.RejectionReason_RateLimitExceeded
			context.Send(sender, response)
			return
		}

		qr.globalRateLimit.Request(weight)
		var data fbinance.OrderBookResponse
		if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
			state.logger.Warn("error fetching order book", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.Code != 0 {
			state.logger.Warn("error fetching order book", log.Error(errors.New(data.Message)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			state.updateRateLimits(data.BaseResponse, nil, qr)
			return
		}

		bids, asks, err := data.ToBidAsk()
		if err != nil {
			err = fmt.Errorf("error converting orderbook: %v", err)
			state.logger.Info("http client error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}
		snapshot = &models.OBL2Snapshot{
			Bids:      bids,
			Asks:      asks,
			Timestamp: timestamppb.Now(),
		}
		response.Success = true
		response.SnapshotL2 = snapshot
		response.SeqNum = data.LastUpdateID
		context.Send(sender, response)
	}()

	return nil
}

func (state *Executor) OnAccountInformationRequest(context actor.Context) error {
	msg := context.Message().(*messages.AccountInformationRequest)
	sender := context.Sender()
	response := &messages.AccountInformationResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	go func() {
		request, weight, err := fbinance.GetAccountInfo(msg.Account.ApiCredentials)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
			return
		}

		qr := state.getQueryRunner(false)
		if qr == nil {
			response.RejectionReason = messages.RejectionReason_RateLimitExceeded
			context.Send(sender, response)
			return
		}

		qr.globalRateLimit.Request(weight)
		var data fbinance.AccountInfo
		if client, ok := state.AccountClients[msg.Account.Name]; ok {
			if err := xutils.PerformJSONRequest(client, request, &data); err != nil {
				state.logger.Warn("error fetching account information", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
		} else {
			if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
				state.logger.Warn("error fetching account information", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
		}

		if data.Code != 0 {
			state.logger.Warn("error fetching account information", log.Error(errors.New(data.Message)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			state.updateRateLimits(data.BaseResponse, nil, qr)
			return
		}
		if data.FeeTier < 0 || data.FeeTier > 9 {
			state.logger.Info(fmt.Sprintf("invalid fee tier: %d", data.FeeTier))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}

		makerFee := MakerFees[data.FeeTier]
		takerFee := TakerFees[data.FeeTier]

		fmt.Println("FEES", makerFee, takerFee)
		response.MakerFee = &wrapperspb.DoubleValue{Value: makerFee}
		response.TakerFee = &wrapperspb.DoubleValue{Value: takerFee}
		response.Success = true
		context.Send(sender, response)
	}()

	return nil
}

func (state *Executor) OnAccountMovementRequest(context actor.Context) error {
	fmt.Println("ON TRADE ACCOUNT MOVEMENT REQUEST !!!!")
	msg := context.Message().(*messages.AccountMovementRequest)
	sender := context.Sender()
	response := &messages.AccountMovementResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	go func() {
		params := fbinance.NewIncomeHistoryRequest()
		if msg.Filter != nil {
			if msg.Filter.Instrument != nil {
				symbol, rej := state.InstrumentToSymbol(msg.Filter.Instrument)
				if rej != nil {
					response.RejectionReason = *rej
					context.Send(sender, response)
					return
				}
				params.SetSymbol(symbol)
			}

			if msg.Filter.From != nil {
				ms := uint64(msg.Filter.From.Seconds*1000) + uint64(msg.Filter.From.Nanos/1000000)
				params.SetFrom(ms)
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
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
			return
		}

		qr := state.getQueryRunner(false)
		if qr == nil {
			response.RejectionReason = messages.RejectionReason_RateLimitExceeded
			context.Send(sender, response)
			return
		}

		qr.globalRateLimit.Request(weight)
		var data []fbinance.Income
		if client, ok := state.AccountClients[msg.Account.Name]; ok {
			if err := xutils.PerformJSONRequest(client, request, &data); err != nil {
				state.logger.Warn("error fetching account movement", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
		} else {
			if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
				state.logger.Warn("error fetching account movement", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
		}

		var movements []*messages.AccountMovement
		for _, t := range data {
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
				context.Send(sender, response)
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
				context.Send(sender, response)
				return
			}
			movements = append(movements, &mvt)
		}
		response.Success = true
		response.Movements = movements
		context.Send(sender, response)
	}()

	return nil
}

func (state *Executor) OnTradeCaptureReportRequest(context actor.Context) error {
	fmt.Println("ON TRADE CAPTURE REPORT REQUEST !!!!")
	msg := context.Message().(*messages.TradeCaptureReportRequest)
	sender := context.Sender()

	response := &messages.TradeCaptureReport{
		RequestID: msg.RequestID,
		Success:   false,
	}

	go func() {
		symbol := ""
		var from, to *uint64
		var fromID string
		if msg.Filter != nil {
			if msg.Filter.Side != nil || msg.Filter.OrderID != nil || msg.Filter.ClientOrderID != nil {
				response.RejectionReason = messages.RejectionReason_UnsupportedFilter
				context.Send(sender, response)
				return
			}

			if msg.Filter.Instrument != nil {
				s, rej := state.InstrumentToSymbol(msg.Filter.Instrument)
				if rej != nil {
					response.RejectionReason = *rej
					context.Send(sender, response)
					return
				}
				symbol = s
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
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
			return
		}

		qr := state.getQueryRunner(false)
		if qr == nil {
			response.RejectionReason = messages.RejectionReason_RateLimitExceeded
			context.Send(sender, response)
			return
		}

		qr.globalRateLimit.Request(weight)
		var data []fbinance.UserTrade
		if client, ok := state.AccountClients[msg.Account.Name]; ok {
			if err := xutils.PerformJSONRequest(client, request, &data); err != nil {
				state.logger.Warn("error fetching account trades", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
		} else {
			if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
				state.logger.Warn("error fetching account trades", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
		}
		var mtrades []*models.TradeCapture
		for _, t := range data {
			sec := state.SymbolToHistoricalSecurity(t.Symbol)
			if sec == nil {
				state.logger.Info("unknown symbol", log.String("symbol", t.Symbol))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Send(sender, response)
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
		context.Send(sender, response)
	}()

	return nil
}

func (state *Executor) OnOrderStatusRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderStatusRequest)
	sender := context.Sender()
	response := &messages.OrderList{
		RequestID: msg.RequestID,
		Success:   false,
	}

	go func() {
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
				context.Send(sender, response)
				return
			}
			if msg.Filter.Instrument != nil {
				s, rej := state.InstrumentToSymbol(msg.Filter.Instrument)
				if rej != nil {
					response.RejectionReason = *rej
					context.Send(sender, response)
					return
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
		if orderID != "" || clOrderID != "" {
			params := fbinance.NewQueryOrderRequest(symbol)
			if orderID != "" {
				orderIDInt, err := strconv.ParseInt(orderID, 10, 64)
				if err != nil {
					response.RejectionReason = messages.RejectionReason_UnsupportedFilter
					context.Send(sender, response)
					return
				}
				params.SetOrderID(orderIDInt)
			}
			if clOrderID != "" {
				params.SetOrigClientOrderID(clOrderID)
			}
			var err error
			request, weight, err = fbinance.QueryOrder(params, msg.Account.ApiCredentials)
			if err != nil {
				state.logger.Warn("error building request", log.Error(err))
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				context.Send(sender, response)
				return
			}
		} else {
			var err error
			request, weight, err = fbinance.QueryOpenOrders(symbol, msg.Account.ApiCredentials)
			if err != nil {
				state.logger.Warn("error building request", log.Error(err))
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				context.Send(sender, response)
				return
			}
		}
		qr := state.getQueryRunner(false)
		if qr == nil {
			response.RejectionReason = messages.RejectionReason_RateLimitExceeded
			context.Send(sender, response)
			return
		}

		qr.globalRateLimit.Request(weight)
		var data []fbinance.OrderData
		if client, ok := state.AccountClients[msg.Account.Name]; ok {
			if err := xutils.PerformJSONRequest(client, request, &data); err != nil {
				state.logger.Warn("error fetching open orders", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
		} else {
			if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
				state.logger.Warn("error fetching open orders", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
		}
		var morders []*models.Order
		for _, o := range data {
			sec := state.SymbolToSecurity(o.Symbol)
			if sec == nil {
				response.RejectionReason = messages.RejectionReason_UnknownSymbol
				context.Send(sender, response)
				return
			}
			ord := orderToModel(&o)
			if orderStatus != nil && ord.OrderStatus != *orderStatus {
				continue
			}
			ord.Instrument.SecurityID = &wrapperspb.UInt64Value{Value: sec.SecurityID}
			morders = append(morders, ord)
		}
		response.Success = true
		response.Orders = morders
		context.Send(sender, response)
	}()

	return nil
}

func (state *Executor) OnPositionsRequest(context actor.Context) error {
	msg := context.Message().(*messages.PositionsRequest)
	sender := context.Sender()
	response := &messages.PositionList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	go func() {
		symbol := ""
		if msg.Instrument != nil {
			s, rej := state.InstrumentToSymbol(msg.Instrument)
			if rej != nil {
				response.RejectionReason = *rej
				context.Send(sender, response)
				return
			}
			symbol = s
		}
		request, weight, err := fbinance.GetPositionRisk(msg.Account.ApiCredentials)
		if err != nil {
			state.logger.Warn("error building request", log.Error(err))
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
			return
		}

		qr := state.getQueryRunner(false)
		if qr == nil {
			response.RejectionReason = messages.RejectionReason_RateLimitExceeded
			context.Send(sender, response)
			return
		}

		qr.globalRateLimit.Request(weight)
		var data []fbinance.AccountPositionRisk
		if client, ok := state.AccountClients[msg.Account.Name]; ok {
			if err := xutils.PerformJSONRequest(client, request, &data); err != nil {
				state.logger.Warn("error fetching positions", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
		} else {
			if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
				state.logger.Warn("error fetching positions", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
		}
		for _, p := range data {
			if p.PositionAmount == 0 {
				continue
			}
			if symbol != "" && p.Symbol != symbol {
				continue
			}
			sec := state.SymbolToSecurity(p.Symbol)
			if sec == nil {
				state.logger.Warn(fmt.Sprintf("unknown symbol %s", p.Symbol))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Send(sender, response)
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
				Quantity:  p.PositionAmount,
				Cost:      cost,
				Cross:     false,
				MarkPrice: wrapperspb.Double(p.MarkPrice),
			}
			response.Positions = append(response.Positions, pos)
		}
		response.Success = true
		context.Send(sender, response)
	}()

	return nil
}

func (state *Executor) OnBalancesRequest(context actor.Context) error {
	msg := context.Message().(*messages.BalancesRequest)
	sender := context.Sender()
	response := &messages.BalanceList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	go func() {
		request, weight, err := fbinance.GetBalance(msg.Account.ApiCredentials)
		if err != nil {
			state.logger.Warn("error building request", log.Error(err))
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
			return
		}
		qr := state.getQueryRunner(false)
		if qr == nil {
			response.RejectionReason = messages.RejectionReason_RateLimitExceeded
			context.Send(sender, response)
			return
		}

		qr.globalRateLimit.Request(weight)

		var data []fbinance.AccountBalance
		if client, ok := state.AccountClients[msg.Account.Name]; ok {
			if err := xutils.PerformJSONRequest(client, request, &data); err != nil {
				state.logger.Warn("error fetching positions", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
		} else {
			if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
				state.logger.Warn("error fetching order book", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
		}

		for _, b := range data {
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
		}

		response.Success = true
		context.Send(sender, response)
	}()

	return nil
}

func (state *Executor) OnNewOrderSingleRequest(context actor.Context) error {
	req := context.Message().(*messages.NewOrderSingleRequest)
	if req.Expire != nil && req.Expire.AsTime().Before(time.Now()) {
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID:       req.RequestID,
			Success:         false,
			RejectionReason: messages.RejectionReason_RequestExpired,
		})
		return nil
	}
	sender := context.Sender()
	response := &messages.NewOrderSingleResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	ar, ok := state.accountRateLimits[req.Account.Name]
	if !ok {
		ar = state.newAccountRateLimit()
		state.accountRateLimits[req.Account.Name] = ar
	}

	if ar.IsRateLimited() {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		response.RateLimitDelay = durationpb.New(ar.DurationBeforeNextRequest(1))
		context.Send(sender, response)
		return nil
	}

	go func() {
		var tickPrecision, lotPrecision int
		sec, rej := state.InstrumentToSecurity(req.Order.Instrument)
		if rej != nil {
			response.RejectionReason = *rej
			context.Send(sender, response)
			return
		}

		tickPrecision = int(math.Ceil(math.Log10(1. / sec.MinPriceIncrement.Value)))
		lotPrecision = int(math.Ceil(math.Log10(1. / sec.RoundLot.Value)))

		params, rej := buildPostOrderRequest(sec.Symbol, req.Order, tickPrecision, lotPrecision)
		if rej != nil {
			response.RejectionReason = *rej
			context.Send(sender, response)
			return
		}

		qr := state.getQueryRunner(false)
		if qr == nil {
			response.RejectionReason = messages.RejectionReason_RateLimitExceeded
			response.RateLimitDelay = durationpb.New(state.durationBeforeNextRequest(1))
			context.Send(sender, response)
			return
		}

		var request *http.Request
		if qr.endpoint != "" {
			var err error
			request, _, err = fbinance.NewOrder(params, req.Account.ApiCredentials, qr.endpoint)
			if err != nil {
				state.logger.Warn("error building request", log.Error(err))
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				context.Send(sender, response)
				return
			}
		} else {
			var err error
			request, _, err = fbinance.NewOrder(params, req.Account.ApiCredentials)
			if err != nil {
				state.logger.Warn("error building request", log.Error(err))
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				context.Send(sender, response)
				return
			}
		}

		ar.Request()
		qr.globalRateLimit.Request(1)

		var data fbinance.OrderData
		start := time.Now()
		if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
			state.logger.Warn(fmt.Sprintf("error posting order for %s", req.Account.Name), log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		fmt.Println("PERFORM", time.Since(start))
		if data.Code != 0 {
			state.logger.Warn(fmt.Sprintf("error posting order for %s", req.Account.Name), log.Error(fmt.Errorf("%d: %s", data.Code, data.Message)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			state.updateRateLimits(data.BaseResponse, ar, qr)
			return
		}
		status := StatusToModel(data.Status)
		if status == nil {
			state.logger.Error(fmt.Sprintf("unknown status %s", data.Status))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}
		response.Success = true
		response.OrderStatus = *status
		response.CumQuantity = data.CumQuantity
		response.LeavesQuantity = data.OriginalQuantity - data.CumQuantity
		response.OrderID = fmt.Sprintf("%d", data.OrderID)
		fmt.Println("NEW SUCCESS", response.OrderID, response.OrderStatus.String())
		context.Send(sender, response)
	}()

	return nil
}

func (state *Executor) OnOrderCancelRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderCancelRequest)
	sender := context.Sender()
	response := &messages.OrderCancelResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	go func() {
		symbol := ""
		if req.Instrument != nil {
			if req.Instrument.Symbol != nil {
				symbol = req.Instrument.Symbol.Value
			} else if req.Instrument.SecurityID != nil {
				sec := state.IDToSecurity(req.Instrument.SecurityID.Value)
				if sec == nil {
					response.RejectionReason = messages.RejectionReason_UnknownSecurityID
					context.Send(sender, response)
					return
				}
				symbol = sec.Symbol
			}
		} else {
			response.RejectionReason = messages.RejectionReason_UnknownSecurityID
			context.Send(sender, response)
			return
		}
		params := fbinance.NewQueryOrderRequest(symbol)
		if req.OrderID != nil {
			orderIDInt, err := strconv.ParseInt(req.OrderID.Value, 10, 64)
			if err != nil {
				response.RejectionReason = messages.RejectionReason_UnknownOrder
				context.Send(sender, response)
				return
			}
			params.SetOrderID(orderIDInt)
		} else if req.ClientOrderID != nil {
			params.SetOrigClientOrderID(req.ClientOrderID.Value)
		} else {
			response.RejectionReason = messages.RejectionReason_UnknownOrder
			context.Send(sender, response)
			return
		}

		qr := state.getQueryRunner(true)
		var request *http.Request
		if qr.endpoint != "" {
			var err error
			request, _, err = fbinance.CancelOrder(params, req.Account.ApiCredentials, qr.endpoint)
			if err != nil {
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				context.Send(sender, response)
				return
			}
		} else {
			var err error
			request, _, err = fbinance.CancelOrder(params, req.Account.ApiCredentials)
			if err != nil {
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				context.Send(sender, response)
				return
			}
		}

		qr.globalRateLimit.Request(1)
		var data fbinance.OrderData
		if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
			state.logger.Warn("error cancelling order", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.Code != 0 {
			state.logger.Warn("error cancelling order", log.Error(errors.New(data.Message)))
			if data.Code == -1020 {
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			} else if data.Code < -1100 && data.Code > -1199 {
				response.RejectionReason = messages.RejectionReason_InvalidRequest
			} else if data.Code == -2011 && data.Message == "Unknown order sent." {
				response.RejectionReason = messages.RejectionReason_UnknownOrder
			}
			context.Send(sender, response)
			state.updateRateLimits(data.BaseResponse, nil, qr)
			return
		}

		response.Success = true
		context.Send(sender, response)
	}()

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

	go func() {
		symbol := ""
		if req.Filter != nil {
			if req.Filter.Instrument != nil {
				if req.Filter.Instrument.Symbol != nil {
					sec := state.SymbolToSecurity(req.Filter.Instrument.Symbol.Value)
					if sec == nil {
						response.RejectionReason = messages.RejectionReason_UnknownSymbol
						context.Send(sender, response)
						return
					}
					symbol = req.Filter.Instrument.Symbol.Value
				} else if req.Filter.Instrument.SecurityID != nil {
					sec := state.IDToSecurity(req.Filter.Instrument.SecurityID.Value)
					if sec == nil {
						response.RejectionReason = messages.RejectionReason_UnknownSecurityID
						context.Send(sender, response)
						return
					}
					symbol = sec.Symbol
				}
			}
			if req.Filter.Side != nil || req.Filter.OrderStatus != nil {
				response.RejectionReason = messages.RejectionReason_UnsupportedFilter
				context.Send(sender, response)
				return
			}
		}
		if symbol == "" {
			response.RejectionReason = messages.RejectionReason_UnknownSymbol
			context.Send(sender, response)
			return
		}

		request, weight, err := fbinance.CancelAllOrders(symbol, req.Account.ApiCredentials)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
			return
		}
		qr := state.getQueryRunner(false)
		if qr == nil {
			response.RejectionReason = messages.RejectionReason_RateLimitExceeded
			context.Send(sender, response)
			return
		}

		qr.globalRateLimit.Request(weight)
		var data fbinance.BaseResponse
		if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
			state.logger.Warn("error cancelling orders", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.Code != 200 {
			state.logger.Warn("error cancelling orders", log.Error(errors.New(data.Message)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}
		response.Success = true
		context.Send(sender, response)
	}()

	return nil
}
