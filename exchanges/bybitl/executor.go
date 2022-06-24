package bybitl

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
	"gitlab.com/alphaticks/xchanger/exchanges/bybitl"
	xutils "gitlab.com/alphaticks/xchanger/utils"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
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
	client         *http.Client
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
		l = exchanges.NewRateLimit(95, time.Minute)
		rl.rateLimits[symbol] = l
	}
	l.Request(weight)
}

func (rl *AccountRateLimit) IsRateLimited(symbol string) bool {
	l, ok := rl.rateLimits[symbol]
	if !ok {
		l = exchanges.NewRateLimit(95, time.Minute)
		rl.rateLimits[symbol] = l
	}
	return l.IsRateLimited()
}

func (rl *AccountRateLimit) DurationBeforeNextRequest(symbol string, weight int) time.Duration {
	l, ok := rl.rateLimits[symbol]
	if !ok {
		l = exchanges.NewRateLimit(95, time.Minute)
		rl.rateLimits[symbol] = l
	}
	return l.DurationBeforeNextRequest(weight)
}

type Executor struct {
	extypes.BaseExecutor
	queryRunners      []*QueryRunner
	accountRateLimits map[string]*AccountRateLimit
	logger            *log.Logger
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

func (state *Executor) durationBeforeNextRequest(post bool, weight int) time.Duration {
	var minDur time.Duration
	for _, q := range state.queryRunners {
		var dur time.Duration
		if post {
			dur1 := q.sPostRateLimit.DurationBeforeNextRequest(weight)
			dur2 := q.mPostRateLimit.DurationBeforeNextRequest(weight)
			if dur1 > dur2 {
				dur = dur1
			} else {
				dur = dur2
			}
		} else {
			dur1 := q.sGetRateLimit.DurationBeforeNextRequest(weight)
			dur2 := q.mGetRateLimit.DurationBeforeNextRequest(weight)
			if dur1 > dur2 {
				dur = dur1
			} else {
				dur = dur2
			}
		}
		if dur < minDur {
			minDur = dur
		}
	}

	return minDur
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
		state.queryRunners = append(state.queryRunners, &QueryRunner{
			client:         client,
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

	request, weight, err := bybitl.GetSymbols()
	if err != nil {
		return err
	}

	qr.Get(weight)

	var data bybitl.SymbolsResponse
	if err := xutils.PerformRequest(qr.client, request, &data); err != nil {
		err := fmt.Errorf("error updating security list: %v", err)
		return err
	}
	if data.RetCode != 0 {
		err := fmt.Errorf("error updating security list: %v", errors.New(data.RetMsg))
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

	var historicalSecurities []*registry.Security
	if state.Registry != nil {
		rres, err := state.Registry.Securities(goContext.Background(), &registry.SecuritiesRequest{
			Filter: &registry.SecurityFilter{
				ExchangeId: []uint32{constants.BYBITL.ID},
			},
		})
		if err != nil {
			return err
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
		context.Send(sender, response)
		return nil
	}
	securityID := msg.Instrument.SecurityID.Value
	security, ok := state.Securities[securityID]
	if !ok {
		response.Success = false
		response.RejectionReason = messages.RejectionReason_UnknownSecurityID
		context.Send(sender, response)
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
		context.Send(sender, response)
		return nil
	}

	qr.rateLimit.Request(weight)

	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

	context.ReenterAfter(future, func(res interface{}, err error) {
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
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Send(sender, response)
			} else if queryResponse.StatusCode >= 500 {

				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Send(sender, response)
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
			context.Send(sender, response)
			return
		}
		fmt.Println(string(bres))

		/*
		if bres.ReturnCode != 0 {
			state.logger.Info("http error", log.Error(errors.New(bres.ReturnMessage)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
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
		context.Send(sender, response)

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
	sender := context.Sender()
	response := &messages.MarketStatisticsResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	symbol, rej := state.getSymbol(msg.Instrument)
	if rej != nil {
		response.RejectionReason = *rej
		context.Send(sender, response)
		return nil
	}

	req, w, err := bybitl.GetTickers()
	if err != nil {
		return fmt.Errorf("error building request: %v", err)
	}
	qr := state.getQueryRunner(false)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Send(sender, response)
		return nil
	}
	qr.Get(w)
	go func() {
		var data bybitl.TickersResponse
		if err := xutils.PerformRequest(qr.client, req, &data); err != nil {
			state.logger.Warn("error fetching tickers", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.RetCode != 0 {
			state.logger.Warn("error fetching tickers", log.Error(errors.New(data.RetMsg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}
		for _, t := range data.Result {
			if symbol != "" && t.Symbol != symbol {
				continue
			}
			sec := state.SymbolToSecurity(symbol)
			if sec == nil {
				state.logger.Warn("error fetching tickers", log.Error(errors.New(data.RetMsg)))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Send(sender, response)
				return
			}
			response.Statistics = append(response.Statistics, &models.Stat{
				Timestamp:  timestamppb.Now(),
				SecurityID: sec.SecurityID,
				Value:      t.MarkPrice,
				StatType:   models.StatType_MarkPrice,
			})
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
		Balances:   nil,
	}
	if msg.Subscribe {
		response.RejectionReason = messages.RejectionReason_UnsupportedSubscription
		context.Send(sender, response)
		return nil
	}
	req, weight, err := bybitl.GetBalance("", msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner(false)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Send(sender, response)
		return nil
	}
	qr.Get(weight)
	go func() {
		var data bybitl.BalanceResponse
		if err := xutils.PerformRequest(qr.client, req, &data); err != nil {
			state.logger.Warn("error fetching balances", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.RetCode != 0 {
			state.logger.Warn("error fetching balances", log.Error(errors.New(data.RetMsg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}
		value := reflect.ValueOf(data.Balance)
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
		Positions:  nil,
	}
	if msg.Subscribe {
		response.RejectionReason = messages.RejectionReason_UnsupportedSubscription
		context.Send(sender, response)
		return nil
	}

	symbol, rej := state.getSymbol(msg.Instrument)
	if rej != nil {
		response.RejectionReason = *rej
		context.Send(sender, response)
		return nil
	}

	request, weight, err := bybitl.GetPositions("", msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner(false)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Send(sender, response)
		return nil
	}
	qr.Get(weight)
	// TODO check account rate limit

	go func() {
		var data bybitl.GetPositionsResponse
		if err := xutils.PerformRequest(qr.client, request, &data); err != nil {
			state.logger.Warn("error fetching positions", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.RetCode != 0 {
			state.logger.Warn("error fetching positions", log.Error(errors.New(data.RetMsg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}
		for _, pos := range data.Positions {
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
			sec := state.SymbolToSecurity(pos.Position.Symbol)
			if sec == nil {
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Send(sender, response)
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
		context.Send(sender, response)
	}()

	return nil
}

func (state *Executor) OnAccountMovementRequest(context actor.Context) error {
	fmt.Println("ON TRADE ACCOUNT MOVEMENT REQUEST !!!!")
	req := context.Message().(*messages.AccountMovementRequest)
	sender := context.Sender()
	response := &messages.AccountMovementResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	if req.Type == messages.AccountMovementType_FundingFee {
		symbol := ""
		var from, to *uint64
		if req.Filter != nil {
			var rej *messages.RejectionReason
			symbol, rej = state.getSymbol(req.Filter.Instrument)
			if rej != nil {
				response.RejectionReason = *rej
				context.Send(sender, response)
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
		}
		params := bybitl.NewGetTradeRecordsParams(symbol)
		if from != nil {
			params.SetStartTime(*from)
		}
		if to != nil {
			params.SetEndTime(*to)
		}
		params.SetExecType(bybitl.ExecFunding)
		request, weight, err := bybitl.GetTradeRecords(params, req.Account.ApiCredentials)
		if err != nil {
			return err
		}
		qr := state.getQueryRunner(false)
		if qr == nil {
			response.RejectionReason = messages.RejectionReason_RateLimitExceeded
			context.Send(sender, response)
			return nil
		}
		qr.Get(weight)

		/*
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
		*/

		go func() {
			var movements []*messages.AccountMovement
			done := false
			for !done {
				var data bybitl.TradingRecordResponse
				if err := xutils.PerformRequest(qr.client, request, &data); err != nil {
					state.logger.Warn("error fetching trade records", log.Error(err))
					response.RejectionReason = messages.RejectionReason_HTTPError
					context.Send(sender, response)
					return
				}
				if data.RetCode != 0 {
					state.logger.Warn("error fetching trade records", log.Error(errors.New(data.RetMsg)))
					response.RejectionReason = messages.RejectionReason_ExchangeAPIError
					context.Send(sender, response)
					return
				}

				sort.Slice(data.TradingRecords.Trades, func(i, j int) bool {
					return data.TradingRecords.Trades[i].TradeTimeMs < data.TradingRecords.Trades[j].TradeTimeMs
				})
				for _, t := range data.TradingRecords.Trades {
					fmt.Println(t)
					sec := state.SymbolToHistoricalSecurity(t.Symbol)
					if sec == nil {
						state.logger.Warn("error fetching trade records", log.Error(errors.New("unknown symbol")))
						response.RejectionReason = messages.RejectionReason_ExchangeAPIError
						context.Send(sender, response)
						return
					}
					mvt := &messages.AccountMovement{
						Asset:      constants.TETHER,
						Change:     t.ExecQty,
						MovementID: t.ExecId,
						Subtype:    t.Symbol,
						Time:       utils.MilliToTimestamp(uint64(t.TradeTimeMs)),
					}
					movements = append(movements, mvt)
				}
				done = len(data.TradingRecords.Trades) == 0
				if !done {
					params.SetPage(int(data.TradingRecords.CurrentPage + 1))
					request, weight, err = bybitl.GetTradeRecords(params, req.Account.ApiCredentials)
					if err != nil {
						panic(err)
					}
					qr.WaitGet(weight)
				}
			}
			response.Success = true
			//response.Movements = movements
			context.Send(sender, response)
		}()
	}

	return nil
}

func (state *Executor) OnTradeCaptureReportRequest(context actor.Context) error {
	req := context.Message().(*messages.TradeCaptureReportRequest)
	sender := context.Sender()
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
			context.Send(sender, response)
			return nil
		}

		var rej *messages.RejectionReason
		symbol, rej = state.getSymbol(req.Filter.Instrument)
		if rej != nil {
			response.RejectionReason = *rej
			context.Send(sender, response)
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
			context.Send(sender, response)
			return nil
		}
		if req.Filter.FromID != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
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
	params.SetExecType(bybitl.ExecTrade)
	request, weight, err := bybitl.GetTradeRecords(params, req.Account.ApiCredentials)
	if err != nil {
		return err
	}
	qr := state.getQueryRunner(false)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Send(sender, response)
		return nil
	}
	qr.Get(weight)

	go func() {
		var mtrades []*models.TradeCapture
		done := false
		for !done {
			var data bybitl.TradingRecordResponse
			if err := xutils.PerformRequest(qr.client, request, &data); err != nil {
				state.logger.Warn("error fetching trade records", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
			if data.RetCode != 0 {
				state.logger.Warn("error fetching trade records", log.Error(errors.New(data.RetMsg)))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Send(sender, response)
				return
			}

			sort.Slice(data.TradingRecords.Trades, func(i, j int) bool {
				return data.TradingRecords.Trades[i].TradeTimeMs < data.TradingRecords.Trades[j].TradeTimeMs
			})
			for _, t := range data.TradingRecords.Trades {
				quantityMul := 1.
				var instrument *models.Instrument
				sec := state.SymbolToHistoricalSecurity(t.Symbol)
				if sec == nil {
					state.logger.Warn("error fetching trade records", log.Error(errors.New("unknown symbol")))
					response.RejectionReason = messages.RejectionReason_ExchangeAPIError
					context.Send(sender, response)
					return
				}
				instrument = &models.Instrument{
					Exchange:   constants.BYBITL,
					Symbol:     &wrapperspb.StringValue{Value: t.Symbol},
					SecurityID: &wrapperspb.UInt64Value{Value: sec.SecurityId},
				}

				quantity := t.ExecQty * quantityMul
				if t.Side == "Sell" {
					quantity *= -1
				}
				comAsset := constants.TETHER
				ts := utils.MilliToTimestamp(uint64(t.TradeTimeMs))
				trd := models.TradeCapture{
					Type:            models.TradeType_Regular,
					Price:           t.ExecPrice,
					Quantity:        quantity,
					Commission:      t.ExecFee,
					CommissionAsset: comAsset,
					TradeID:         t.ExecId,
					Instrument:      instrument,
					Trade_LinkID:    nil,
					OrderID:         &wrapperspb.StringValue{Value: t.OrderId},
					TransactionTime: ts,
				}

				if t.Side == "Buy" {
					trd.Side = models.Side_Buy
				} else {
					trd.Side = models.Side_Sell
				}
				mtrades = append(mtrades, &trd)
			}
			done = len(data.TradingRecords.Trades) == 0
			if !done {
				params.SetPage(int(data.TradingRecords.CurrentPage + 1))
				request, weight, err = bybitl.GetTradeRecords(params, req.Account.ApiCredentials)
				if err != nil {
					panic(err)
				}
				qr.WaitGet(weight)
			}
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
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
		Orders:     nil,
	}
	if msg.Subscribe {
		response.RejectionReason = messages.RejectionReason_UnsupportedSubscription
		context.Send(sender, response)
		return nil
	}
	symbol := ""
	orderId := ""
	clOrderId := ""
	var orderStatus bybitl.OrderStatus
	if msg.Filter != nil {
		if msg.Filter.Side != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedFilter
			context.Send(sender, response)
			return nil
		}
		if msg.Filter.Instrument != nil {
			if msg.Filter.Instrument.Symbol != nil {
				symbol = msg.Filter.Instrument.Symbol.Value
			} else if msg.Filter.Instrument.SecurityID != nil {
				sec, ok := state.Securities[msg.Filter.Instrument.SecurityID.Value]
				if !ok {
					response.RejectionReason = messages.RejectionReason_UnknownSecurityID
					context.Send(sender, response)
					return nil
				}
				symbol = sec.Symbol
			}
		} else {
			response.RejectionReason = messages.RejectionReason_MissingInstrument
			context.Send(sender, response)
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
		context.Send(sender, response)
		return nil
	}
	qr.Get(1)
	// TODO check account rate limit
	go func() {
		var data bybitl.QueryActiveOrdersResponse
		if err := xutils.PerformRequest(qr.client, req, &data); err != nil {
			state.logger.Warn("error fetching orders", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.RetCode != 0 {
			state.logger.Warn("error fetching orders", log.Error(errors.New(data.RetMsg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}
		for _, ord := range data.Orders {
			if orderStatus != "" && ord.OrderStatus != orderStatus {
				continue
			}
			sec := state.SymbolToSecurity(ord.Symbol)
			if sec == nil {
				state.logger.Info("got order with unknown symbol: " + ord.Symbol)
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Send(sender, response)
				return
			}
			o := OrderToModel(&ord)
			o.Instrument.SecurityID = wrapperspb.UInt64(sec.SecurityID)
			o.Instrument.Symbol = wrapperspb.String(sec.Symbol)
			response.Orders = append(response.Orders, o)
		}
		response.Success = true
		context.Send(sender, response)
	}()

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

	symbol := ""
	var tickPrecision, lotPrecision int
	if req.Order.Instrument != nil {
		if req.Order.Instrument.Symbol != nil {
			symbol = req.Order.Instrument.Symbol.Value
			sec := state.SymbolToSecurity(symbol)
			if sec == nil {
				response.RejectionReason = messages.RejectionReason_UnknownSecurityID
				context.Send(sender, response)
				return nil
			}
			tickPrecision = int(math.Ceil(math.Log10(1. / sec.MinPriceIncrement.Value)))
			lotPrecision = int(math.Ceil(math.Log10(1. / sec.RoundLot.Value)))
		} else if req.Order.Instrument.SecurityID != nil {
			sec, ok := state.Securities[req.Order.Instrument.SecurityID.Value]
			if !ok {
				response.RejectionReason = messages.RejectionReason_UnknownSecurityID
				context.Send(sender, response)
				return nil
			}
			symbol = sec.Symbol
			tickPrecision = int(math.Ceil(math.Log10(1. / sec.MinPriceIncrement.Value)))
			lotPrecision = int(math.Ceil(math.Log10(1. / sec.RoundLot.Value)))
		}
	} else {
		response.RejectionReason = messages.RejectionReason_UnknownSecurityID
		context.Send(sender, response)
		return nil
	}

	ar, ok := state.accountRateLimits[req.Account.Name]
	if !ok {
		ar = NewAccountRateLimit()
		state.accountRateLimits[req.Account.Name] = ar
	}

	if ar.IsRateLimited(symbol) {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		response.RateLimitDelay = durationpb.New(ar.DurationBeforeNextRequest(symbol, 1))
		fmt.Println("DURATION", ar.DurationBeforeNextRequest(symbol, 1))
		context.Send(sender, response)
		return nil
	}
	ar.Request(symbol, 1)

	params, rej := buildPostOrderRequest(symbol, req.Order, tickPrecision, lotPrecision)
	if rej != nil {
		response.RejectionReason = *rej
		context.Send(sender, response)
		return nil
	}

	request, weight, err := bybitl.PostActiveOrder(params, req.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner(true)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		response.RateLimitDelay = durationpb.New(state.durationBeforeNextRequest(false, weight))
		context.Send(sender, response)
		return nil
	}

	qr.Post(weight)
	go func() {
		var data bybitl.PostOrderResponse
		if err := xutils.PerformRequest(qr.client, request, &data); err != nil {
			state.logger.Warn("error fetching positions", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.RetCode != 0 {
			state.logger.Warn("error fetching positions", log.Error(errors.New(data.RetMsg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}
		status := StatusToModel(data.Order.OrderStatus)
		if status == nil {
			state.logger.Error(fmt.Sprintf("unknown status %s", data.Order.OrderStatus))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}
		response.Success = true
		response.OrderStatus = *status
		response.CumQuantity = data.Order.CumExecQty
		response.LeavesQuantity = data.Order.Qty
		response.OrderID = data.Order.OrderId
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

	symbol := ""
	if req.Instrument != nil {
		if req.Instrument.Symbol != nil {
			symbol = req.Instrument.Symbol.Value
		} else if req.Instrument.SecurityID != nil {
			sec, ok := state.Securities[req.Instrument.SecurityID.Value]
			if !ok {
				response.RejectionReason = messages.RejectionReason_UnknownSecurityID
				context.Send(sender, response)
				return nil
			}
			symbol = sec.Symbol
		}
	} else {
		response.RejectionReason = messages.RejectionReason_UnknownSecurityID
		context.Send(sender, response)
		return nil
	}

	ar, ok := state.accountRateLimits[req.Account.Name]
	if !ok {
		ar = NewAccountRateLimit()
		state.accountRateLimits[req.Account.Name] = ar
	}
	ar.Request(symbol, 1)

	params := bybitl.NewCancelActiveOrderParams(symbol)
	if req.OrderID != nil {
		params.SetOrderId(req.OrderID.Value)
	} else if req.ClientOrderID != nil {
		params.SetOrderLinkId(req.ClientOrderID.Value)
	} else {
		response.RejectionReason = messages.RejectionReason_UnknownOrder
		context.Send(sender, response)
		return nil
	}

	request, weight, err := bybitl.CancelActiveOrder(params, req.Account.ApiCredentials)
	if err != nil {
		response.RejectionReason = messages.RejectionReason_InvalidRequest
		context.Send(sender, response)
		return nil
	}
	qr := state.getQueryRunner(true)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		response.RateLimitDelay = durationpb.New(state.durationBeforeNextRequest(true, weight))
		context.Send(sender, response)
		return nil
	}
	qr.Post(weight)

	go func() {
		// We ignore rate limits on cancel

		var data bybitl.CancelOrderResponse
		if err := xutils.PerformRequest(qr.client, request, &data); err != nil {
			state.logger.Warn("error cancelling order", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.RetCode != 0 {
			state.logger.Warn("error cancelling order", log.Error(errors.New(data.RetMsg)))
			switch data.RetCode {
			case 20001:
				response.RejectionReason = messages.RejectionReason_UnknownOrder
			default:
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			}
			context.Send(sender, response)
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
		ResponseID: uint64(time.Now().UnixNano()),
		RequestID:  req.RequestID,
		Success:    false,
	}

	qr := state.getQueryRunner(true)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Send(sender, response)
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
					context.Send(sender, response)
					return nil
				}
				symbol = sec.Symbol
			}
		} else {
			response.RejectionReason = messages.RejectionReason_MissingInstrument
			context.Send(sender, response)
			return nil
		}
	} else {
		response.RejectionReason = messages.RejectionReason_UnsupportedFilter
		context.Send(sender, response)
		return nil
	}

	params := bybitl.NewCancelAllActiveParams(symbol)
	request, _, err := bybitl.CancelAllActiveOrders(params, req.Account.ApiCredentials)
	if err != nil {
		return err
	}
	qr.Post(1)

	go func() {
		var data bybitl.CancelOrderResponse
		if err := xutils.PerformRequest(qr.client, request, &data); err != nil {
			state.logger.Warn("error canceling orders", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.RetCode != 0 {
			state.logger.Warn("error canceling orders", log.Error(errors.New(data.RetMsg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}
		response.Success = true
		context.Send(sender, response)
	}()
	return nil
}

func (state *Executor) OnOrderReplaceRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderReplaceRequest)
	sender := context.Sender()
	response := &messages.OrderReplaceResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	qr := state.getQueryRunner(true)
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Send(sender, response)
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
				context.Send(sender, response)
				return nil
			}
			symbol = sec.Symbol
		}
	} else {
		response.RejectionReason = messages.RejectionReason_MissingInstrument
		context.Send(sender, response)
		return nil
	}
	params := bybitl.NewAmendOrderParams(symbol)
	if req.Update.OrderID != nil {
		params.SetOrderId(req.Update.OrderID.Value)
	} else if req.Update.OrigClientOrderID != nil {
		params.SetOrderLinkId(req.Update.OrigClientOrderID.Value)
	} else {
		response.RejectionReason = messages.RejectionReason_UnknownOrder
		context.Send(sender, response)
		return nil
	}
	request, _, err := bybitl.AmendOrder(params, req.Account.ApiCredentials)
	if err != nil {
		return err
	}
	qr.Post(1)
	// TODO check account rate limits
	go func() {
		var data bybitl.AmendOrderResponse
		if err := xutils.PerformRequest(qr.client, request, &data); err != nil {
			state.logger.Warn("error fetching positions", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.RetCode != 0 {
			state.logger.Warn("error fetching positions", log.Error(errors.New(data.RetMsg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}
		response.Success = true
		response.OrderID = data.AmendedOrder.OrderId
		context.Send(sender, response)
	}()
	return nil
}
