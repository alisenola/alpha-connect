package bybitl

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
	"gitlab.com/alphaticks/xchanger/exchanges/bybitl"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"time"
	"unicode"
)

func (state *Executor) getSymbol(instrument *models.Instrument) (string, *messages.RejectionReason) {
	symbol := ""
	if instrument != nil {
		if instrument.Symbol != nil {
			symbol = instrument.Symbol.Value
		} else if instrument.SecurityID != nil {
			sec, ok := state.Securities[instrument.SecurityID.Value]
			if !ok {
				rej := messages.RejectionReason_UnknownSecurityID
				return symbol, &rej
			}
			symbol = sec.Symbol
		}
	}
	return symbol, nil
}

type QueryRunner struct {
	pid            *actor.PID
	sPostRateLimit *exchanges.RateLimit
	mPostRateLimit *exchanges.RateLimit
	sGetRateLimit  *exchanges.RateLimit
	mGetRateLimit  *exchanges.RateLimit
}

func (qr *QueryRunner) Get(weight int) {
	qr.sGetRateLimit.Request(weight)
	qr.mGetRateLimit.Request(weight)
}

func (qr *QueryRunner) WaitGet(weight int) {
	qr.sGetRateLimit.WaitRequest(weight)
	qr.mGetRateLimit.WaitRequest(weight)
}

func (qr *QueryRunner) Post(weight int) {
	qr.sPostRateLimit.Request(weight)
	qr.mPostRateLimit.Request(weight)
}

type AccountRateLimit struct {
	rateLimits map[string]*exchanges.RateLimit
}

func NewAccountRateLimit() *AccountRateLimit {
	return &AccountRateLimit{
		rateLimits: make(map[string]*exchanges.RateLimit),
	}
}

func (rl *AccountRateLimit) Request(symbol string, weight int) {
	l, ok := rl.rateLimits[symbol]
	if !ok {
		l = exchanges.NewRateLimit(100, time.Minute)
		rl.rateLimits[symbol] = l
	}
	l.Request(weight)
}

func (rl *AccountRateLimit) IsRateLimited(symbol string) bool {
	l, ok := rl.rateLimits[symbol]
	if !ok {
		l = exchanges.NewRateLimit(100, time.Minute)
		rl.rateLimits[symbol] = l
	}
	return l.IsRateLimited()
}

type Executor struct {
	extypes.BaseExecutor
	historicalSecurities  map[uint64]*registry.Security
	symbolToHistoricalSec map[string]*registry.Security
	queryRunners          []*QueryRunner
	accountRateLimits     map[string]*AccountRateLimit
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

func (state *Executor) getQueryRunner(post bool) *QueryRunner {
	sort.Slice(state.queryRunners, func(i, j int) bool {
		return rand.Uint64()%2 == 0
	})

	var qr *QueryRunner
	for _, q := range state.queryRunners {
		if post {
			if !q.sPostRateLimit.IsRateLimited() && !q.mPostRateLimit.IsRateLimited() {
				qr = q
				break
			}
		} else {
			if !q.sGetRateLimit.IsRateLimited() && !q.mGetRateLimit.IsRateLimited() {
				qr = q
				break
			}
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
		log.String("type", reflect.TypeOf(*state).String()))

	dialers := state.ExecutorConfig.DialerPool.GetDialers()
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
			pid:            context.Spawn(props),
			sGetRateLimit:  exchanges.NewRateLimit(50*100, 2*time.Minute),
			mGetRateLimit:  exchanges.NewRateLimit(70*4, 5*time.Second),
			sPostRateLimit: exchanges.NewRateLimit(20*100, 2*time.Minute),
			mPostRateLimit: exchanges.NewRateLimit(50*4, 5*time.Second),
		})
	}

	state.accountRateLimits = make(map[string]*AccountRateLimit)

	if err := state.UpdateSecurityList(context); err != nil {
		state.logger.Warn("error updating security list: %v", log.Error(err))
	}

	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	qr := state.getQueryRunner(false)
	if qr == nil {
		return fmt.Errorf("rate limited")
	}

	request, _, err := bybitl.GetSymbols()
	if err != nil {
		return err
	}

	qr.Get(1)

	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
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

	var data bybitl.SymbolsResponse
	err = json.Unmarshal(queryResponse.Response, &data)
	if err != nil {
		err = fmt.Errorf(
			"error unmarshaling response: %v",
			err)
		return err
	}
	if data.RetCode != 0 {
		err = fmt.Errorf(
			"got wrong return code: %s",
			data.RetMsg)
		return err
	}

	var securities []*models.Security
	for _, symbol := range data.Result {
		// Only linear
		if symbol.QuoteCurrency == "USD" {
			continue
		}
		baseCurrency, ok := constants.GetAssetBySymbol(symbol.BaseCurrency)
		if !ok {
			//state.logger.Info(fmt.Sprintf("unknown currency %s", symbol.BaseCurrency))
			continue
		}
		quoteCurrency, ok := constants.GetAssetBySymbol(symbol.QuoteCurrency)
		if !ok {
			//state.logger.Info(fmt.Sprintf("unknown currency %s", baseStr))
			continue
		}
		security := models.Security{}
		security.Symbol = symbol.Name
		security.Underlying = baseCurrency
		security.QuoteCurrency = quoteCurrency
		security.Status = models.InstrumentStatus_Trading
		security.Exchange = constants.BYBITL
		day := symbol.Alias[len(symbol.Alias)-2:]
		if unicode.IsNumber(rune(day[0])) && unicode.IsNumber(rune(day[1])) {
			security.SecurityType = enum.SecurityType_CRYPTO_FUT

			month := symbol.Alias[len(symbol.Alias)-4 : len(symbol.Alias)-2]
			monthInt, err := strconv.ParseInt(month, 10, 64)
			if err != nil {
				state.logger.Info(fmt.Sprintf("error parsing month %s", month))
				continue
			}
			dayInt, err := strconv.ParseInt(day, 10, 64)
			if err != nil {
				state.logger.Info(fmt.Sprintf("error parsing day: %d", dayInt))
				continue
			}
			year := time.Now().Year()
			if int(monthInt) < int(time.Now().Month()) {
				year += 1
			}
			date := time.Date(year, time.Month(monthInt), int(dayInt), 0, 0, 0, 0, time.UTC)
			ts := timestamppb.New(date)
			security.MaturityDate = ts
		} else {
			security.SecurityType = enum.SecurityType_CRYPTO_PERP
		}
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinPriceIncrement = &wrapperspb.DoubleValue{Value: symbol.PriceFilter.TickSize}
		security.RoundLot = &wrapperspb.DoubleValue{Value: symbol.LotSizeFilter.QuantityStep}
		security.IsInverse = false
		security.MakerFee = &wrapperspb.DoubleValue{Value: symbol.MakerFee}
		security.TakerFee = &wrapperspb.DoubleValue{Value: symbol.TakerFee}
		security.Multiplier = &wrapperspb.DoubleValue{Value: 1}
		security.MinLimitQuantity = &wrapperspb.DoubleValue{Value: symbol.LotSizeFilter.MinTradingQuantity}
		security.MinMarketQuantity = &wrapperspb.DoubleValue{Value: symbol.LotSizeFilter.MinTradingQuantity}
		security.MaxLimitQuantity = &wrapperspb.DoubleValue{Value: symbol.LotSizeFilter.MaxTradingQuantity}
		security.MaxMarketQuantity = &wrapperspb.DoubleValue{Value: symbol.LotSizeFilter.MaxTradingQuantity}

		securities = append(securities, &security)
	}

	state.Securities = make(map[uint64]*models.Security)
	state.SymbolToSec = make(map[string]*models.Security)
	for _, s := range securities {
		state.Securities[s.SecurityID] = s
		state.SymbolToSec[s.Symbol] = s
	}

	state.historicalSecurities = make(map[uint64]*registry.Security)
	state.symbolToHistoricalSec = make(map[string]*registry.Security)
	if state.Registry != nil {
		rres, err := state.Registry.Securities(goContext.Background(), &registry.SecuritiesRequest{
			Filter: &registry.SecurityFilter{
				ExchangeId: []uint32{constants.BYBITL.ID},
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
	securities := make([]*models.Security, len(state.Securities))
	i := 0
	for _, v := range state.Securities {
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

/*
func (state *Executor) OnHistoricalLiquidationsRequest(context actor.Context) error {
	msg := context.Message().(*messages.HistoricalLiquidationsRequest)
	response := &messages.HistoricalLiquidationsResponse{
		RequestID: msg.RequestID,
		Success:   false,
	}

	// Find security
	if msg.Instrument == nil || msg.Instrument.SecurityID == nil {
		response.Success = false
		response.RejectionReason = messages.RejectionReason_UnknownSecurityID
		context.Respond(response)
		return nil
	}
	securityID := msg.Instrument.SecurityID.Value
	security, ok := state.Securities[securityID]
	if !ok {
		response.Success = false
		response.RejectionReason = messages.RejectionReason_UnknownSecurityID
		context.Respond(response)
		return nil
	}

	p := bybitl.NewLiquidatedOrdersParams(security.Symbol)
	if msg.From != nil {
		p.SetStartTime(uint64(msg.From.Seconds*1000) + uint64(msg.From.Nanos/1000000))
	}
	if msg.To != nil {
		p.SetEndTime(uint64(msg.To.Seconds*1000) + uint64(msg.From.Nanos/1000000))
	}

	request, weight, err := bybitl.GetLiquidatedOrders(p)
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

		type Res struct {
			ReturnCode    int                  `json:"ret_code"`
			ReturnMessage string               `json:"ret_msg"`
			ExitCode      string               `json:"ext_code"`
			ExitInfo      string               `json:"ext_info"`
			Result        []bybiti.Liquidation `json:"result"`
		}
		var bres json.RawMessage
		err = json.Unmarshal(queryResponse.Response, &bres)
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		fmt.Println(string(bres))

		/*
		if bres.ReturnCode != 0 {
			state.logger.Info("http error", log.Error(errors.New(bres.ReturnMessage)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		var liquidations []*models.Liquidation
		for _, l := range bres.Result {
			liquidations = append(liquidations, &models.Liquidation{
				Bid:       l.Side == "Sell",
				Timestamp: utils.MilliToTimestamp(l.Time),
				OrderID:   uint64(l.ID),
				Price:     l.Price,
				Quantity:  l.Quantity,
			})
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
*/

func (state *Executor) OnMarketStatisticsRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketStatisticsRequest)
	if len(msg.Statistics) != 1 || msg.Statistics[0] != models.StatType_MarkPrice {
		context.Respond(&messages.MarketStatisticsResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.RejectionReason_UnsupportedRequest,
		})
	}
	response := &messages.MarketStatisticsResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	symbol, rej := state.getSymbol(msg.Instrument)
	if rej != nil {
		response.RejectionReason = *rej
		context.Respond(response)
		return nil
	}

	req, w, err := bybitl.GetTickers()
	if err != nil {
		return fmt.Errorf("error building request: %v", err)
	}
	qr := state.getQueryRunner(false)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}
	qr.Get(w)
	f := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: req}, 15*time.Second)
	context.ReenterAfter(f, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		qResponse := res.(*jobs.PerformQueryResponse)
		if qResponse.StatusCode != 200 {
			if qResponse.StatusCode >= 400 && qResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"%d %s",
					qResponse.StatusCode,
					string(qResponse.Response))
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			} else if qResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"%d %s",
					qResponse.StatusCode,
					string(qResponse.Response))
				state.logger.Info("http server error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			}
			return
		}
		var tickers bybitl.TickersResponse
		err = json.Unmarshal(qResponse.Response, &tickers)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		if tickers.RetCode != 0 {
			state.logger.Info("error getting orders", log.Error(errors.New(tickers.RetMsg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		for _, t := range tickers.Result {
			if symbol != "" {
				if t.Symbol == symbol {
					sec := state.SymbolToSec[symbol]
					response.Statistics = append(response.Statistics, &models.Stat{
						Timestamp:  timestamppb.Now(),
						SecurityID: sec.SecurityID,
						Value:      t.MarkPrice,
						StatType:   models.StatType_MarkPrice,
					})
				}
			} else {
				sec, ok := state.SymbolToSec[t.Symbol]
				if ok {
					response.Statistics = append(response.Statistics, &models.Stat{
						Timestamp:  timestamppb.Now(),
						SecurityID: sec.SecurityID,
						Value:      t.MarkPrice,
						StatType:   models.StatType_MarkPrice,
					})
				}

			}
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
	if msg.Subscribe {
		response.RejectionReason = messages.RejectionReason_UnsupportedSubscription
		context.Respond(response)
		return nil
	}
	req, _, err := bybitl.GetBalance("", msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner(false)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}
	qr.Get(1)
	f := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: req}, 15*time.Second)
	context.ReenterAfter(f, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		qResponse := res.(*jobs.PerformQueryResponse)
		if qResponse.StatusCode != 200 {
			if qResponse.StatusCode >= 400 && qResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"%d %s",
					qResponse.StatusCode,
					string(qResponse.Response))
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			} else if qResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"%d %s",
					qResponse.StatusCode,
					string(qResponse.Response))
				state.logger.Info("http server error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			}
			return
		}
		var balances bybitl.BalanceResponse
		err = json.Unmarshal(qResponse.Response, &balances)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		if balances.RetCode != 0 {
			state.logger.Info("error getting orders", log.Error(errors.New(balances.RetMsg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		value := reflect.ValueOf(balances.Balance)
		for i := 0; i < value.NumField(); i++ {
			coin := value.Field(i).Interface().(bybitl.Coin)
			if coin.WalletBalance == 0 {
				continue
			}
			asset, ok := constants.GetAssetBySymbol(value.Type().Field(i).Name)
			if !ok {
				state.logger.Error("got balance for unknown asset", log.String("asset", value.Type().Field(i).Name))
				continue
			}
			response.Balances = append(response.Balances, &models.Balance{
				Account:  msg.Account.Name,
				Asset:    asset,
				Quantity: coin.WalletBalance,
			})
		}
		response.Success = true
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
	if msg.Subscribe {
		response.RejectionReason = messages.RejectionReason_UnsupportedSubscription
		context.Respond(response)
		return nil
	}

	symbol, rej := state.getSymbol(msg.Instrument)
	if rej != nil {
		response.RejectionReason = *rej
		context.Respond(response)
		return nil
	}

	req, _, err := bybitl.GetPositions("", msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner(false)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}
	qr.Get(1)
	// TODO check account rate limit

	f := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: req}, 15*time.Second)
	context.ReenterAfter(f, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		qResponse := res.(*jobs.PerformQueryResponse)
		if qResponse.StatusCode != 200 {
			if qResponse.StatusCode >= 400 && qResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"%d %s",
					qResponse.StatusCode,
					string(qResponse.Response))
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			} else if qResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"%d %s",
					qResponse.StatusCode,
					string(qResponse.Response))
				state.logger.Info("http server error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			}
			return
		}
		var positions bybitl.GetPositionsResponse
		err = json.Unmarshal(qResponse.Response, &positions)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if positions.RetCode != 0 {
			state.logger.Info("error getting orders", log.Error(errors.New(positions.RetMsg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		for _, pos := range positions.Positions {
			if pos.Position.Size == 0 {
				continue
			}
			if !pos.IsValid {
				state.logger.Warn("got invalid position", log.String("position", fmt.Sprintf("%+v", pos.Position)))
				continue
			}
			if symbol != "" && symbol != pos.Position.Symbol {
				continue
			}
			sec, ok := state.SymbolToSec[pos.Position.Symbol]
			if !ok {
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			size := pos.Position.Size
			if pos.Position.Side == bybitl.Sell {
				size *= -1
			}
			cost := size * pos.Position.EntryPrice

			response.Positions = append(response.Positions, &models.Position{
				Account: msg.Account.Name,
				Instrument: &models.Instrument{
					Exchange:   constants.BYBITL,
					Symbol:     &wrapperspb.StringValue{Value: pos.Position.Symbol},
					SecurityID: &wrapperspb.UInt64Value{Value: sec.SecurityID},
				},
				Quantity: size,
				Cost:     cost,
				Cross:    false,
			})
		}
		response.Success = true
		context.Respond(response)
	})
	return nil
}

func (state *Executor) OnTradeCaptureReportRequest(context actor.Context) error {
	req := context.Message().(*messages.TradeCaptureReportRequest)
	response := &messages.TradeCaptureReport{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	symbol := ""
	var from, to *uint64
	if req.Filter != nil {
		if req.Filter.Side != nil || req.Filter.OrderID != nil || req.Filter.ClientOrderID != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedFilter
			context.Respond(response)
			return nil
		}

		var rej *messages.RejectionReason
		symbol, rej = state.getSymbol(req.Filter.Instrument)
		if rej != nil {
			response.RejectionReason = *rej
			context.Respond(response)
			return nil
		}

		if req.Filter.From != nil {
			ts := utils.TimestampToMilli(req.Filter.From)
			from = &ts
		}
		if req.Filter.To != nil {
			ts := utils.TimestampToMilli(req.Filter.To)
			to = &ts
		}
		if req.Filter.OrderID != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Respond(response)
			return nil
		}
		if req.Filter.FromID != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Respond(response)
			return nil
		}
	}
	params := bybitl.NewGetTradeRecordsParams(symbol)
	if from != nil {
		params.SetStartTime(*from)
	}
	if to != nil {
		params.SetEndTime(*to)
	}

	r, weight, err := bybitl.GetTradeRecords(params, req.Account.ApiCredentials)
	if err != nil {
		return err
	}
	qr := state.getQueryRunner(false)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}
	qr.Get(weight)

	var mtrades []*models.TradeCapture
	sender := context.Sender()
	var processFuture func(res interface{}, err error)
	processFuture = func(res interface{}, err error) {
		if err != nil {
			state.logger.Info("request error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
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
		var trades bybitl.TradingRecordResponse
		err = json.Unmarshal(queryResponse.Response, &trades)
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		if trades.RetCode != 0 {
			state.logger.Info("error getting trades", log.Error(errors.New(trades.RetMsg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		sort.Slice(trades.TradingRecords.Trades, func(i, j int) bool {
			return trades.TradingRecords.Trades[i].TradeTimeMs < trades.TradingRecords.Trades[j].TradeTimeMs
		})
		for _, t := range trades.TradingRecords.Trades {
			quantityMul := 1.
			var instrument *models.Instrument
			if _, ok := state.symbolToHistoricalSec[t.Symbol]; ok {
				sec := state.symbolToHistoricalSec[t.Symbol]
				instrument = &models.Instrument{
					Exchange:   constants.BYBITL,
					Symbol:     &wrapperspb.StringValue{Value: t.Symbol},
					SecurityID: &wrapperspb.UInt64Value{Value: sec.SecurityId},
				}
			} else {
				instrument = &models.Instrument{
					Exchange:   constants.BYBITL,
					Symbol:     &wrapperspb.StringValue{Value: t.Symbol},
					SecurityID: nil,
				}
			}

			quantity := t.ExecQty * quantityMul
			if t.Side == "sell" {
				quantity *= -1
			}
			comAsset := constants.TETHER
			ts := utils.MilliToTimestamp(uint64(t.TradeTimeMs))
			trd := models.TradeCapture{
				Type:            models.TradeType_Regular,
				Price:           t.Price,
				Quantity:        quantity,
				Commission:      t.ExecFee,
				CommissionAsset: comAsset,
				TradeID:         t.ExecId,
				Instrument:      instrument,
				Trade_LinkID:    nil,
				OrderID:         &wrapperspb.StringValue{Value: t.OrderId},
				TransactionTime: ts,
			}

			if t.Side == "buy" {
				trd.Side = models.Side_Buy
			} else {
				trd.Side = models.Side_Sell
			}
			mtrades = append(mtrades, &trd)
		}
		done := len(trades.TradingRecords.Trades) == 0
		if !done {
			fmt.Println("WE KEEP GOING", int(trades.TradingRecords.CurrentPage+1))
			params.SetPage(int(trades.TradingRecords.CurrentPage + 1))
			r, weight, err := bybitl.GetTradeRecords(params, req.Account.ApiCredentials)
			qr.WaitGet(weight)
			if err != nil {
				response.RejectionReason = messages.RejectionReason_UnsupportedOrderCharacteristic
				context.Send(sender, response)
				return
			}
			fut := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: r}, 15*time.Second)
			context.ReenterAfter(fut, processFuture)
			return
		} else {
			response.Success = true
			response.Trades = mtrades
			context.Send(sender, response)
			return
		}
	}
	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: r}, 10*time.Minute)
	context.ReenterAfter(future, processFuture)
	return nil
}

func (state *Executor) OnOrderStatusRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderStatusRequest)
	response := &messages.OrderList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
		Orders:     nil,
	}
	if msg.Subscribe {
		response.RejectionReason = messages.RejectionReason_UnsupportedSubscription
		context.Respond(response)
		return nil
	}
	symbol := ""
	orderId := ""
	clOrderId := ""
	var orderStatus bybitl.OrderStatus
	if msg.Filter != nil {
		if msg.Filter.Side != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedFilter
			context.Respond(response)
			return nil
		}
		if msg.Filter.Instrument != nil {
			if msg.Filter.Instrument.Symbol != nil {
				symbol = msg.Filter.Instrument.Symbol.Value
			} else if msg.Filter.Instrument.SecurityID != nil {
				sec, ok := state.Securities[msg.Filter.Instrument.SecurityID.Value]
				if !ok {
					response.RejectionReason = messages.RejectionReason_UnknownSecurityID
					context.Respond(response)
					return nil
				}
				symbol = sec.Symbol
			}
		} else {
			response.RejectionReason = messages.RejectionReason_MissingInstrument
			context.Respond(response)
			return nil
		}
		if msg.Filter.OrderStatus != nil {
			orderStatus = StatusToBybitl(msg.Filter.OrderStatus.Value)
		}
		if msg.Filter.OrderID != nil {
			orderId = msg.Filter.OrderID.Value
		}
		if msg.Filter.ClientOrderID != nil {
			clOrderId = msg.Filter.ClientOrderID.Value
		}
	}

	params := bybitl.NewQueryActiveOrderParams(symbol)
	if clOrderId != "" {
		params.SetOrderLinkID(clOrderId)
	}
	if orderId != "" {
		params.SetOrderID(orderId)
	}
	req, _, err := bybitl.QueryActiveOrdersRT(params, msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner(false)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}
	qr.Get(1)
	// TODO check account rate limit
	f := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: req}, 15*time.Second)
	context.ReenterAfter(f, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		qResponse := res.(*jobs.PerformQueryResponse)
		if qResponse.StatusCode != 200 {
			if qResponse.StatusCode >= 400 && qResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"%d %s",
					qResponse.StatusCode,
					string(qResponse.Response))
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			} else if qResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"%d %s",
					qResponse.StatusCode,
					string(qResponse.Response))
				state.logger.Info("http server error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			}
			return
		}
		var orders bybitl.QueryActiveOrdersResponse
		err = json.Unmarshal(qResponse.Response, &orders)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if orders.RetCode != 0 {
			state.logger.Info("error getting orders", log.Error(errors.New(orders.RetMsg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		for _, ord := range orders.Orders {
			if orderStatus != "" && ord.OrderStatus != orderStatus {
				continue
			}
			sec, ok := state.SymbolToSec[ord.Symbol]
			if !ok {
				state.logger.Info("got order with unknown symbol: " + ord.Symbol)
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			o := OrderToModel(&ord)
			o.Instrument.SecurityID = wrapperspb.UInt64(sec.SecurityID)
			o.Instrument.Symbol = wrapperspb.String(sec.Symbol)
			response.Orders = append(response.Orders, o)
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
	qr := state.getQueryRunner(true)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	symbol := ""
	var tickPrecision, lotPrecision int
	if req.Order.Instrument != nil {
		if req.Order.Instrument.Symbol != nil {
			symbol = req.Order.Instrument.Symbol.Value
			sec, ok := state.SymbolToSec[symbol]
			if !ok {
				response.RejectionReason = messages.RejectionReason_UnknownSecurityID
				context.Respond(response)
				return nil
			}
			tickPrecision = int(math.Ceil(math.Log10(1. / sec.MinPriceIncrement.Value)))
			lotPrecision = int(math.Ceil(math.Log10(1. / sec.RoundLot.Value)))
		} else if req.Order.Instrument.SecurityID != nil {
			sec, ok := state.Securities[req.Order.Instrument.SecurityID.Value]
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

	ar, ok := state.accountRateLimits[req.Account.Name]
	if !ok {
		ar = NewAccountRateLimit()
		state.accountRateLimits[req.Account.Name] = ar
	}

	if ar.IsRateLimited(symbol) {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	params, rej := buildPostOrderRequest(symbol, req.Order, tickPrecision, lotPrecision)
	if rej != nil {
		response.RejectionReason = *rej
		context.Respond(response)
		return nil
	}

	request, _, err := bybitl.PostActiveOrder(params, req.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr.Post(1)
	// TODO check account rate limit
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
		var order bybitl.PostOrderResponse
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if order.RetCode != 0 {
			state.logger.Info("error posting order", log.Error(errors.New(fmt.Sprintf("%d: ", order.RetCode)+order.RetMsg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		status := StatusToModel(order.Order.OrderStatus)
		if status == nil {
			state.logger.Error(fmt.Sprintf("unknown status %s", order.Order.OrderStatus))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
		response.OrderStatus = *status
		response.CumQuantity = order.Order.CumExecQty
		response.LeavesQuantity = order.Order.Qty
		response.OrderID = order.Order.OrderId
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

	qr := state.getQueryRunner(true)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	symbol := ""
	if req.Instrument != nil {
		if req.Instrument.Symbol != nil {
			symbol = req.Instrument.Symbol.Value
		} else if req.Instrument.SecurityID != nil {
			sec, ok := state.Securities[req.Instrument.SecurityID.Value]
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

	ar, ok := state.accountRateLimits[req.Account.Name]
	if !ok {
		ar = NewAccountRateLimit()
		state.accountRateLimits[req.Account.Name] = ar
	}

	if ar.IsRateLimited(symbol) {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	params := bybitl.NewCancelActiveOrderParams(symbol)
	if req.OrderID != nil {
		params.SetOrderId(req.OrderID.Value)
	} else if req.ClientOrderID != nil {
		params.SetOrderLinkId(req.ClientOrderID.Value)
	} else {
		response.RejectionReason = messages.RejectionReason_UnknownOrder
		context.Respond(response)
		return nil
	}

	request, weight, err := bybitl.CancelActiveOrder(params, req.Account.ApiCredentials)
	if err != nil {
		return err
	}

	fmt.Println(request.URL)

	qr.Post(weight)
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
		var order bybitl.CancelOrderResponse
		fmt.Println(string(queryResponse.Response))
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if order.RetCode != 0 {
			state.logger.Info("error cancelling order", log.Error(errors.New(order.RetMsg)))
			switch order.RetCode {
			case 20001:
				response.RejectionReason = messages.RejectionReason_UnknownOrder
			default:
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			}
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
		ResponseID: uint64(time.Now().UnixNano()),
		RequestID:  req.RequestID,
		Success:    false,
	}

	qr := state.getQueryRunner(true)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}
	symbol := ""
	if req.Filter != nil {
		if req.Filter.Instrument != nil {
			if req.Filter.Instrument.Symbol != nil {
				symbol = req.Filter.Instrument.Symbol.Value
			} else if req.Filter.Instrument.SecurityID != nil {
				sec, ok := state.Securities[req.Filter.Instrument.SecurityID.Value]
				if !ok {
					response.RejectionReason = messages.RejectionReason_UnknownSecurityID
					context.Respond(response)
					return nil
				}
				symbol = sec.Symbol
			}
		} else {
			response.RejectionReason = messages.RejectionReason_MissingInstrument
			context.Respond(response)
			return nil
		}
	} else {
		response.RejectionReason = messages.RejectionReason_UnsupportedFilter
		context.Respond(response)
		return nil
	}

	params := bybitl.NewCancelAllActiveParams(symbol)
	request, _, err := bybitl.CancelAllActiveOrders(params, req.Account.ApiCredentials)
	if err != nil {
		return err
	}
	qr.Post(1)

	f := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
	context.ReenterAfter(f, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				err := fmt.Errorf("%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf("%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http server error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
			}
			return
		}
		var cancelResponse bybitl.CancelOrdersResponse
		err = json.Unmarshal(queryResponse.Response, &cancelResponse)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if cancelResponse.RetCode != 0 {
			state.logger.Info("error mass cancelling orders", log.Error(errors.New(cancelResponse.RetMsg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
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

	qr := state.getQueryRunner(true)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}

	symbol := ""
	if req.Instrument != nil {
		if req.Instrument.Symbol != nil {
			symbol = req.Instrument.Symbol.Value
		} else if req.Instrument.SecurityID != nil {
			sec, ok := state.Securities[req.Instrument.SecurityID.Value]
			if !ok {
				response.RejectionReason = messages.RejectionReason_UnknownSecurityID
				context.Respond(response)
				return nil
			}
			symbol = sec.Symbol
		}
	} else {
		response.RejectionReason = messages.RejectionReason_MissingInstrument
		context.Respond(response)
		return nil
	}
	params := bybitl.NewAmendOrderParams(symbol)
	if req.Update.OrderID != nil {
		params.SetOrderId(req.Update.OrderID.Value)
	} else if req.Update.OrigClientOrderID != nil {
		params.SetOrderLinkId(req.Update.OrigClientOrderID.Value)
	} else {
		response.RejectionReason = messages.RejectionReason_UnknownOrder
		context.Respond(response)
		return nil
	}
	request, _, err := bybitl.AmendOrder(params, req.Account.ApiCredentials)
	if err != nil {
		return err
	}
	qr.Post(1)
	// TODO check account rate limits
	f := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
	context.ReenterAfter(f, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				err := fmt.Errorf("%d %s", queryResponse.StatusCode, string(queryResponse.Response))
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
				return
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf("%d %s", queryResponse.StatusCode, string(queryResponse.Response))
				state.logger.Info("http server error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
				return
			}
		}
		var amendResponse bybitl.AmendOrderResponse
		err = json.Unmarshal(queryResponse.Response, &amendResponse)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if amendResponse.RetCode != 0 {
			state.logger.Info("error amending order", log.Error(errors.New(amendResponse.RetMsg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
		response.OrderID = amendResponse.AmendedOrder.OrderId
		context.Respond(response)
	})
	return nil
}
