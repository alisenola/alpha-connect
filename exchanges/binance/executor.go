package binance

import (
	goContext "context"
	"errors"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/enum"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/binance"
	xutils "gitlab.com/alphaticks/xchanger/utils"
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
	client          *http.Client
	globalRateLimit *exchanges.RateLimit
}

type AccountRateLimit struct {
	second *exchanges.RateLimit
	day    *exchanges.RateLimit
}

func NewAccountRateLimit(second, day *exchanges.RateLimit) *AccountRateLimit {
	return &AccountRateLimit{
		second: second,
		day:    day,
	}
}

func (rl *AccountRateLimit) Request() {
	rl.second.Request(1)
	rl.day.Request(1)
}

func (rl *AccountRateLimit) IsRateLimited() bool {
	return rl.second.IsRateLimited() || rl.day.IsRateLimited()
}

type Executor struct {
	extypes.BaseExecutor
	accountRateLimits   map[string]*AccountRateLimit
	newAccountRateLimit func() *AccountRateLimit
	queryRunners        []*QueryRunner
	logger              *log.Logger
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
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(state).String()))

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
		state.queryRunners = append(state.queryRunners, &QueryRunner{
			client:          client,
			globalRateLimit: nil,
		})
	}

	state.accountRateLimits = make(map[string]*AccountRateLimit)

	for _, qr := range state.queryRunners {
		request, weight, err := binance.GetExchangeInfo()
		if err != nil {
			return err
		}
		var data binance.ExchangeInfo
		if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
			err := fmt.Errorf("error updating security list: %v", err)
			return err
		}
		if data.Code != 0 {
			err := fmt.Errorf("error updating security list: %v", errors.New(data.Message))
			return err
		}

		// Initialize rate limit
		var secondOrderLimit, dayOrderLimit int
		// Initialize rate limit
		for _, rateLimit := range data.RateLimits {
			if rateLimit.RateLimitType == "ORDERS" {
				if rateLimit.Interval == "SECOND" {
					secondOrderLimit = rateLimit.Limit
				} else if rateLimit.Interval == "DAY" {
					dayOrderLimit = rateLimit.Limit
				}
			} else if rateLimit.RateLimitType == "REQUEST_WEIGHT" {
				qr.globalRateLimit = exchanges.NewRateLimit(rateLimit.Limit, time.Minute)
			}
		}
		state.newAccountRateLimit = func() *AccountRateLimit {
			return NewAccountRateLimit(exchanges.NewRateLimit(secondOrderLimit, time.Second), exchanges.NewRateLimit(dayOrderLimit, time.Hour*24))
		}
		qr.globalRateLimit.Request(weight)
	}

	return state.UpdateSecurityList(context)
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	// Get http request and the expected response
	request, weight, err := binance.GetExchangeInfo()
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		return fmt.Errorf("rate limited")
	}

	qr.globalRateLimit.Request(weight)

	var data binance.ExchangeInfo
	if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
		err := fmt.Errorf("error updating security list: %v", err)
		return err
	}
	if data.Code != 0 {
		return fmt.Errorf(
			"binance api error: %d %s",
			data.Code,
			data.Message)
	}

	var securities []*models.Security
	for _, symbol := range data.Symbols {
		baseCurrency, ok := constants.GetAssetBySymbol(symbol.BaseAsset)
		if !ok {
			if symbol.Status == "TRADING" {
				state.logger.Debug(fmt.Sprintf("unknown currency %s", symbol.BaseAsset))
			}
			continue
		}
		quote := symbol.QuoteAsset
		if sym, ok := binance.BINANCE_TO_GLOBAL[quote]; ok {
			quote = sym
		}
		quoteCurrency, ok := constants.GetAssetBySymbol(quote)
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
		security.Exchange = constants.BINANCE
		security.SecurityType = enum.SecurityType_CRYPTO_SPOT
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		for _, filter := range symbol.Filters {
			if filter.FilterType == "PRICE_FILTER" {
				security.MinPriceIncrement = &wrapperspb.DoubleValue{Value: filter.TickSize}
			} else if filter.FilterType == "LOT_SIZE" {
				security.RoundLot = &wrapperspb.DoubleValue{Value: filter.StepSize}
			}
		}
		securities = append(securities, &security)
	}

	var historicalSecurities []*registry.Security
	if state.Registry != nil {
		rres, err := state.Registry.Securities(goContext.Background(), &registry.SecuritiesRequest{
			Filter: &registry.SecurityFilter{
				ExchangeId: []uint32{constants.BINANCE.ID},
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
	if msg.Instrument == nil || msg.Instrument.Symbol == nil {
		response.RejectionReason = messages.RejectionReason_MissingInstrument
		context.Send(sender, response)
		return nil
	}

	go func() {
		symbol := msg.Instrument.Symbol.Value
		// Get http request and the expected response
		request, weight, err := binance.GetOrderBook(symbol, 1000)
		if err != nil {
			state.logger.Warn("error building request", log.Error(err))
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
			return
		}

		qr := state.getQueryRunner()
		if qr == nil {
			response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
			context.Send(sender, response)
			return
		}
		qr.globalRateLimit.Request(weight)

		var data binance.OrderBookData
		if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
			state.logger.Warn("error fetching open interests", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.Code != 0 {
			state.logger.Warn("error fetching open interests", log.Error(errors.New(data.Message)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}

		bids, asks, err := data.ToBidAsk()
		if err != nil {
			state.logger.Info("error decoding query response", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		snapshot = &models.OBL2Snapshot{
			Bids:      bids,
			Asks:      asks,
			Timestamp: timestamppb.Now(),
		}
		response.SnapshotL2 = snapshot
		response.SeqNum = data.LastUpdateID
		response.Success = true
		context.Send(sender, response)
	}()

	return nil
}

func (state *Executor) OnAccountMovementRequest(context actor.Context) error {
	fmt.Println("ON TRADE ACCOUNT MOVEMENT REQUEST !!!!")
	msg := context.Message().(*messages.AccountMovementRequest)
	switch msg.Type {
	case messages.AccountMovementType_Deposit:
		return state.onDepositHistoryRequest(context)
	case messages.AccountMovementType_Withdrawal:
		return state.onWithdrawalHistoryRequest(context)
	case messages.AccountMovementType_Exchange:
		return state.onConvertHistoryRequest(context)
	}
	return nil
}

func (state *Executor) onDepositHistoryRequest(context actor.Context) error {
	msg := context.Message().(*messages.AccountMovementRequest)
	sender := context.Sender()
	response := &messages.AccountMovementResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	ut1Params := binance.NewUserUniversalTransferHistoryRequest(binance.FUNDING_MAIN)
	ut2Params := binance.NewUserUniversalTransferHistoryRequest(binance.UMFUTURE_MAIN)
	adParams := binance.NewAccountDepositHistoryRequest()
	var startTime, endTime *time.Time

	if msg.Filter != nil {
		if msg.Filter.From != nil {
			v := msg.Filter.From.AsTime()
			startTime = &v
		}
		if msg.Filter.To != nil {
			v := msg.Filter.To.AsTime()
			if startTime != nil && v.Sub(*startTime) > 80*24*time.Hour {
				v = startTime.Add(80 * 24 * time.Hour)
			}
			endTime = &v
		}
	}
	if startTime != nil {
		ut1Params.SetStartTime(uint64(startTime.UnixMilli()))
		ut2Params.SetStartTime(uint64(startTime.UnixMilli()))
		adParams.SetStartTime(uint64(startTime.UnixMilli()))
		fmt.Println("from", startTime.String())
	}
	if endTime != nil {
		ut1Params.SetEndTime(uint64(endTime.UnixMilli()))
		ut2Params.SetEndTime(uint64(endTime.UnixMilli()))
		adParams.SetEndTime(uint64(endTime.UnixMilli()))
		fmt.Println("to", endTime.String())
	}
	adParams.SetLimit(1000)
	accountDepositRequest, weight, err := binance.GetAccountDepositHistory(adParams, msg.Account.ApiCredentials)
	if err != nil {
		response.RejectionReason = messages.RejectionReason_UnsupportedRequest
		context.Send(sender, response)
		return nil
	}
	universalTransferRequest1, weight, err := binance.GetUserUniversalTransferHistory(ut1Params, msg.Account.ApiCredentials)
	if err != nil {
		response.RejectionReason = messages.RejectionReason_UnsupportedRequest
		context.Send(sender, response)
		return nil
	}
	universalTransferRequest2, weight, err := binance.GetUserUniversalTransferHistory(ut2Params, msg.Account.ApiCredentials)
	if err != nil {
		response.RejectionReason = messages.RejectionReason_UnsupportedRequest
		context.Send(sender, response)
		return nil
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Send(sender, response)
		return nil
	}

	qr.globalRateLimit.Request(weight)
	var depositData []binance.DepositData
	var universalTransferResponse1 binance.UniversalTransferResponse
	var universalTransferResponse2 binance.UniversalTransferResponse
	if client, ok := state.AccountClients[msg.Account.Name]; ok {
		if err := xutils.PerformJSONRequest(client, accountDepositRequest, &depositData); err != nil {
			state.logger.Warn("error fetching account movement", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return nil
		}
		if err := xutils.PerformJSONRequest(client, universalTransferRequest1, &universalTransferResponse1); err != nil {
			state.logger.Warn("error fetching account movement", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return nil
		}
		if err := xutils.PerformJSONRequest(client, universalTransferRequest2, &universalTransferResponse2); err != nil {
			state.logger.Warn("error fetching account movement", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return nil
		}
	} else {
		if err := xutils.PerformJSONRequest(qr.client, accountDepositRequest, &depositData); err != nil {
			state.logger.Warn("error fetching account movement", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return nil
		}
		if err := xutils.PerformJSONRequest(qr.client, universalTransferRequest1, &universalTransferResponse1); err != nil {
			state.logger.Warn("error fetching account movement", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return nil
		}
		if err := xutils.PerformJSONRequest(qr.client, universalTransferRequest2, &universalTransferResponse2); err != nil {
			state.logger.Warn("error fetching account movement", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return nil
		}
	}

	var movements []*messages.AccountMovement
	for _, t := range depositData {
		fmt.Println("DEPOSIT HISTORY", t)
		asset, ok := constants.GetAssetBySymbol(t.Coin)
		if !ok {
			state.logger.Warn("unknown asset " + t.Coin)
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return nil
		}
		mvt := messages.AccountMovement{
			Asset:      asset,
			Change:     t.Amount,
			MovementID: fmt.Sprintf("%s", t.TxId),
			Time:       utils.MilliToTimestamp(t.InsertTime),
			Type:       messages.AccountMovementType_Deposit,
		}
		movements = append(movements, &mvt)
	}
	fmt.Println("UNIVERSAL TRANSFER TOTAL", universalTransferResponse1.Total)
	for _, t := range universalTransferResponse1.Rows {
		fmt.Println("TRANSFER HISTORY", t)
		asset, ok := constants.GetAssetBySymbol(t.Asset)
		if !ok {
			state.logger.Warn("unknown asset " + t.Asset)
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return nil
		}
		mvt := messages.AccountMovement{
			Asset:      asset,
			Change:     t.Amount,
			MovementID: fmt.Sprintf("%d", t.TranId),
			Time:       utils.MilliToTimestamp(t.Timestamp),
			Type:       messages.AccountMovementType_Deposit,
		}
		movements = append(movements, &mvt)
	}
	for _, t := range universalTransferResponse2.Rows {
		fmt.Println("TRANSFER HISTORY", t)
		asset, ok := constants.GetAssetBySymbol(t.Asset)
		if !ok {
			state.logger.Warn("unknown asset " + t.Asset)
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return nil
		}
		mvt := messages.AccountMovement{
			Asset:      asset,
			Change:     t.Amount,
			MovementID: fmt.Sprintf("%d", t.TranId),
			Time:       utils.MilliToTimestamp(t.Timestamp),
			Type:       messages.AccountMovementType_Deposit,
		}
		movements = append(movements, &mvt)
	}
	response.Success = true
	sort.Slice(movements, func(i, j int) bool {
		return movements[i].Time.AsTime().Before(movements[j].Time.AsTime())
	})
	response.Movements = movements
	context.Send(sender, response)
	return nil
}

func (state *Executor) onWithdrawalHistoryRequest(context actor.Context) error {
	msg := context.Message().(*messages.AccountMovementRequest)
	sender := context.Sender()
	response := &messages.AccountMovementResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	utParams := binance.NewUserUniversalTransferHistoryRequest(binance.MAIN_UMFUTURE)
	var startTime, endTime *time.Time
	if msg.Filter != nil {
		if msg.Filter.From != nil {
			v := msg.Filter.From.AsTime()
			startTime = &v
		}
		if msg.Filter.To != nil {
			v := msg.Filter.To.AsTime()
			if startTime != nil && v.Sub(*startTime) > 80*24*time.Hour {
				v = startTime.Add(80 * 24 * time.Hour)
			}
			endTime = &v
		}
	}
	if startTime != nil {
		utParams.SetStartTime(uint64(startTime.UnixMilli()))
		fmt.Println("from", startTime.String())
	}
	if endTime != nil {
		utParams.SetEndTime(uint64(endTime.UnixMilli()))
		fmt.Println("to", endTime.String())
	}
	universalTransferRequest, weight, err := binance.GetUserUniversalTransferHistory(utParams, msg.Account.ApiCredentials)
	if err != nil {
		response.RejectionReason = messages.RejectionReason_UnsupportedRequest
		context.Send(sender, response)
		return nil
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Send(sender, response)
		return nil
	}

	qr.globalRateLimit.Request(weight)
	var universalTransferResponse binance.UniversalTransferResponse
	if client, ok := state.AccountClients[msg.Account.Name]; ok {
		if err := xutils.PerformJSONRequest(client, universalTransferRequest, &universalTransferResponse); err != nil {
			state.logger.Warn("error fetching account movement", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return nil
		}
	} else {
		if err := xutils.PerformJSONRequest(qr.client, universalTransferRequest, &universalTransferResponse); err != nil {
			state.logger.Warn("error fetching account movement", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return nil
		}
	}

	var movements []*messages.AccountMovement
	fmt.Println("UNIVERSAL TRANSFER TOTAL", universalTransferResponse.Total)
	for _, t := range universalTransferResponse.Rows {
		fmt.Println("TRANSFER HISTORY", t, utils.MilliToTimestamp(t.Timestamp).AsTime().String())
		asset, ok := constants.GetAssetBySymbol(t.Asset)
		if !ok {
			state.logger.Warn("unknown asset " + t.Asset)
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return nil
		}
		mvt := messages.AccountMovement{
			Asset:      asset,
			Change:     -t.Amount,
			MovementID: fmt.Sprintf("%d", t.TranId),
			Time:       utils.MilliToTimestamp(t.Timestamp),
			Type:       messages.AccountMovementType_Withdrawal,
		}
		movements = append(movements, &mvt)
	}
	response.Success = true
	sort.Slice(movements, func(i, j int) bool {
		return movements[i].Time.AsTime().Before(movements[j].Time.AsTime())
	})
	response.Movements = movements
	context.Send(sender, response)
	return nil
}

func (state *Executor) onConvertHistoryRequest(context actor.Context) error {
	fmt.Println("ON ACCOUNT Convert REQUEST !!!!")
	msg := context.Message().(*messages.AccountMovementRequest)
	sender := context.Sender()

	response := &messages.AccountMovementResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	var startTime, endTime *time.Time
	cParams := binance.NewConvertTradeHistoryRequest()
	if msg.Filter != nil {
		if msg.Filter.From != nil {
			v := msg.Filter.From.AsTime()
			startTime = &v
		}
		if msg.Filter.To != nil {
			v := msg.Filter.To.AsTime()
			if startTime != nil && v.Sub(*startTime) > 80*24*time.Hour {
				v = startTime.Add(29 * 24 * time.Hour)
			}
			endTime = &v
		}
		if startTime != nil {
			cParams.SetStartTime(uint64(startTime.UnixMilli()))
			fmt.Println("from", startTime.String())
		}
		if endTime != nil {
			cParams.SetEndTime(uint64(endTime.UnixMilli()))
			fmt.Println("to", endTime.String())
		}
	}

	request, weight, err := binance.GetConvertTradeHistory(cParams, msg.Account.ApiCredentials)
	if err != nil {
		response.RejectionReason = messages.RejectionReason_UnsupportedRequest
		context.Send(sender, response)
		return nil
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Send(sender, response)
		return nil
	}
	qr.globalRateLimit.Request(weight)
	var cdata binance.ConvertTradeHistoryResponse
	if client, ok := state.AccountClients[msg.Account.Name]; ok {
		if err := xutils.PerformJSONRequest(client, request, &cdata); err != nil {
			state.logger.Warn("error fetching account trades", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
		}
	} else {
		if err := xutils.PerformJSONRequest(qr.client, request, &cdata); err != nil {
			state.logger.Warn("error fetching account trades", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
		}
	}
	var movements []*messages.AccountMovement
	fmt.Println(cdata)
	for _, t := range cdata.List {
		to, ok := constants.GetAssetBySymbol(t.ToAsset)
		if !ok {
			state.logger.Warn("unknown asset " + t.ToAsset)
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return nil
		}
		from, ok := constants.GetAssetBySymbol(t.FromAsset)
		if !ok {
			state.logger.Warn("unknown asset " + t.FromAsset)
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return nil
		}
		movements = append(movements, &messages.AccountMovement{
			Asset:      from,
			Change:     -t.FromAmount,
			MovementID: fmt.Sprintf("from-%s-%d", t.QuoteId, t.OrderId),
			Time:       utils.MilliToTimestamp(t.CreateTime),
			Type:       messages.AccountMovementType_Exchange,
		})
		movements = append(movements, &messages.AccountMovement{
			Asset:      to,
			Change:     t.ToAmount,
			MovementID: fmt.Sprintf("to-%s-%d", t.QuoteId, t.OrderId),
			Time:       utils.MilliToTimestamp(t.CreateTime),
			Type:       messages.AccountMovementType_Exchange,
		})
	}
	response.Success = true
	response.Movements = movements
	context.Send(sender, response)
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
		tParams := binance.NewAccountTradeListRequest(symbol)
		// If from is not set, but to is set,
		// If from is set, but to is not set, ok

		if fromID != "" {
			fromIDInt, _ := strconv.ParseUint(fromID, 10, 64)
			tParams.SetFromID(fromIDInt)
		} else {
			if from == nil || *from == 0 {
				tParams.SetFromID(0)
			} else {
				tParams.SetStartTime(*from)
				if to != nil {
					if *to-*from > (7 * 24 * 60 * 60 * 1000) {
						*to = *from + (7 * 24 * 60 * 60 * 1000)
					}
					tParams.SetEndTime(*to)
				}
			}
		}

		trequest, weight, err := binance.GetAccountTradeList(tParams, msg.Account.ApiCredentials)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
			return
		}
		qr := state.getQueryRunner()
		if qr == nil {
			response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
			context.Send(sender, response)
			return
		}

		qr.globalRateLimit.Request(weight)
		var tdata []binance.TradeData
		if client, ok := state.AccountClients[msg.Account.Name]; ok {
			if err := xutils.PerformJSONRequest(client, trequest, &tdata); err != nil {
				state.logger.Warn("error fetching account trades", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
		} else {
			if err := xutils.PerformJSONRequest(qr.client, trequest, &tdata); err != nil {
				state.logger.Warn("error fetching account trades", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
		}
		var mtrades []*models.TradeCapture
		for _, t := range tdata {
			sec := state.SymbolToHistoricalSecurity(t.Symbol)
			if sec == nil {
				state.logger.Info("unknown symbol", log.String("symbol", t.Symbol))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Send(sender, response)
				return
			}
			quantity := t.Quantity
			if !t.IsBuyer {
				quantity *= -1
			}
			asset, ok := constants.GetAssetBySymbol(t.CommissionAsset)
			if !ok {
				state.logger.Warn("unknown asset " + t.CommissionAsset)
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Send(sender, response)
				return
			}
			trd := models.TradeCapture{
				Type:            models.TradeType_Regular,
				Price:           t.Price,
				Quantity:        quantity,
				Commission:      t.Commission,
				CommissionAsset: asset,
				TradeID:         fmt.Sprintf("%d-%d", t.ID, t.OrderID),
				Instrument: &models.Instrument{
					Exchange:   constants.BINANCE,
					Symbol:     &wrapperspb.StringValue{Value: t.Symbol},
					SecurityID: &wrapperspb.UInt64Value{Value: sec.SecurityId},
				},
				Trade_LinkID:    nil,
				OrderID:         &wrapperspb.StringValue{Value: fmt.Sprintf("%d", t.OrderID)},
				TransactionTime: utils.MilliToTimestamp(t.Time),
			}

			if t.IsBuyer {
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
		if msg.Filter != nil {
			if msg.Filter.OrderStatus != nil {
				response.RejectionReason = messages.RejectionReason_UnsupportedFilter
				context.Send(sender, response)
				return
			}
			if msg.Filter.Side != nil {
				response.RejectionReason = messages.RejectionReason_UnsupportedFilter
				context.Send(sender, response)
				return
			}
			if msg.Filter.Instrument != nil {
				if msg.Filter.Instrument.Symbol != nil {
					symbol = msg.Filter.Instrument.Symbol.Value
				} else if msg.Filter.Instrument.SecurityID != nil {
					state.SecuritiesLock.RLock()
					sec, ok := state.Securities[msg.Filter.Instrument.SecurityID.Value]
					state.SecuritiesLock.RUnlock()
					if !ok {
						response.RejectionReason = messages.RejectionReason_UnknownSecurityID
						context.Send(sender, response)
						return
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
			params := binance.NewQueryOrderRequest(symbol)
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
			request, weight, err = binance.QueryOrder(params, msg.Account.ApiCredentials)
			if err != nil {
				state.logger.Warn("error building request", log.Error(err))
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				context.Send(sender, response)
				return
			}
		} else {
			multi = true
			var err error
			request, weight, err = binance.GetOpenOrders(symbol, msg.Account.ApiCredentials)
			if err != nil {
				state.logger.Warn("error building request", log.Error(err))
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				context.Send(sender, response)
				return
			}
		}

		qr := state.getQueryRunner()
		if qr == nil {
			response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
			context.Send(sender, response)
			return
		}

		qr.globalRateLimit.Request(weight)
		var data []binance.OrderData
		if multi {
			if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
				state.logger.Warn("error fetching open orders", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
		} else {
			var sorder binance.QueryOrderResponse
			if err := xutils.PerformJSONRequest(qr.client, request, &sorder); err != nil {
				state.logger.Warn("error fetching open orders", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
			data = []binance.OrderData{sorder.OrderData}
		}

		fmt.Println("ORDERS", data)
		var morders []*models.Order
		for _, o := range data {
			sec := state.SymbolToSecurity(o.Symbol)
			if sec == nil {
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Send(sender, response)
				return
			}
			ord := OrderToModel(&o)
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
		Success:    true,
		Positions:  nil,
	}
	context.Send(sender, response)
	return nil
}

func (state *Executor) OnBalancesRequest(context actor.Context) error {
	msg := context.Message().(*messages.BalancesRequest)
	sender := context.Sender()
	response := &messages.BalanceList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
		Balances:   nil,
	}

	go func() {
		request, weight, err := binance.GetAccountInformation(msg.Account.ApiCredentials)
		if err != nil {
			state.logger.Warn("error building request", log.Error(err))
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
			return
		}

		qr := state.getQueryRunner()
		if qr == nil {
			response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
			context.Send(sender, response)
			return
		}

		qr.globalRateLimit.Request(weight)
		var data binance.AccountInformationResponse
		if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
			state.logger.Warn("error fetching positions", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		for _, b := range data.Balances {
			tot := b.Locked + b.Free
			if tot == 0. {
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
				Quantity: tot,
			})
		}

		response.Success = true
		context.Send(sender, response)
	}()

	return nil
}

func (state *Executor) OnNewOrderSingleRequest(context actor.Context) error {
	fmt.Println("EXEC NEW SINGLE")
	req := context.Message().(*messages.NewOrderSingleRequest)
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
		response.RejectionReason = messages.RejectionReason_AccountRateLimitExceeded
		context.Send(sender, response)
		return nil
	}

	go func() {
		qr := state.getQueryRunner()
		if qr == nil {
			response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
			context.Send(sender, response)
			return
		}

		symbol := ""
		var tickPrecision, lotPrecision int
		if req.Order.Instrument != nil {
			if req.Order.Instrument.Symbol != nil {
				symbol = req.Order.Instrument.Symbol.Value
				sec := state.SymbolToSecurity(symbol)
				if sec == nil {
					response.RejectionReason = messages.RejectionReason_UnknownSecurityID
					context.Send(sender, response)
					return
				}
				tickPrecision = int(math.Ceil(math.Log10(1. / sec.MinPriceIncrement.Value)))
				lotPrecision = int(math.Ceil(math.Log10(1. / sec.RoundLot.Value)))
			} else if req.Order.Instrument.SecurityID != nil {
				sec := state.IDToSecurity(req.Order.Instrument.SecurityID.Value)
				if sec == nil {
					response.RejectionReason = messages.RejectionReason_UnknownSecurityID
					context.Send(sender, response)
				}
				symbol = sec.Symbol
				tickPrecision = int(math.Ceil(math.Log10(1. / sec.MinPriceIncrement.Value)))
				lotPrecision = int(math.Ceil(math.Log10(1. / sec.RoundLot.Value)))
			}
		} else {
			response.RejectionReason = messages.RejectionReason_UnknownSecurityID
			context.Send(sender, response)
			return
		}

		params, rej := buildPostOrderRequest(symbol, req.Order, tickPrecision, lotPrecision)
		if rej != nil {
			response.RejectionReason = *rej
			context.Send(sender, response)
			return
		}

		request, weight, err := binance.NewOrder(params, req.Account.ApiCredentials)
		if err != nil {
			state.logger.Warn("error building request", log.Error(err))
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
			return
		}

		fmt.Println(request.URL)

		qr.globalRateLimit.Request(weight)
		ar.Request()

		var data binance.NewOrderResponse
		if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
			state.logger.Warn("error fetching order book", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.Code != 0 {
			state.logger.Warn("error posting order", log.Error(errors.New(data.Message)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}
		response.Success = true
		response.OrderID = fmt.Sprintf("%d", data.OrderId)
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
		if req.OrderID == nil {
			response.RejectionReason = messages.RejectionReason_UnknownOrder
			context.Send(sender, response)
			return
		}
		orderID, err := strconv.ParseInt(req.OrderID.Value, 10, 64)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_UnknownOrder
			context.Send(sender, response)
			return
		}

		qr := state.getQueryRunner()
		if qr == nil {
			response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
			context.Send(sender, response)
			return
		}

		request, weight, err := binance.CancelOrder(symbol, orderID, req.Account.ApiCredentials)
		if err != nil {
			state.logger.Warn("error building request", log.Error(err))
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
			return
		}

		fmt.Println(request.URL)
		qr.globalRateLimit.Request(weight)

		var data binance.OrderData
		if err := xutils.PerformJSONRequest(qr.client, request, &data); err != nil {
			state.logger.Warn("error posting order", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.Code != 0 {
			state.logger.Warn("error posting order", log.Error(errors.New(data.Message)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}
		response.Success = true
		context.Send(sender, response)
	}()

	return nil
}
