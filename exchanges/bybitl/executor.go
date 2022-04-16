package bybitl

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
	"gitlab.com/alphaticks/xchanger/exchanges/bybitl"
	xutils "gitlab.com/alphaticks/xchanger/utils"
	"math"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"time"
	"unicode"
)

type QueryRunner struct {
	pid       *actor.PID
	rateLimit *exchanges.RateLimit
}

type Executor struct {
	extypes.BaseExecutor
	securities   map[uint64]*models.Security
	symbolToSec  map[string]*models.Security
	queryRunners []*QueryRunner
	dialerPool   *xutils.DialerPool
	logger       *log.Logger
}

func NewExecutor(dialerPool *xutils.DialerPool) actor.Actor {
	return &Executor{
		queryRunners: nil,
		logger:       nil,
		dialerPool:   dialerPool,
	}
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
		if !q.rateLimit.IsRateLimited() {
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
		log.String("type", reflect.TypeOf(*state).String()))

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
			rateLimit: exchanges.NewRateLimit(150, time.Second),
		})
	}

	if err := state.UpdateSecurityList(context); err != nil {
		state.logger.Warn("error updating security list: %v", log.Error(err))
	}

	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	request, weight, err := bybitl.GetSymbols()
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
		return fmt.Errorf("rate limited")
	}

	qr.rateLimit.Request(weight)

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
			//state.logger.Info(fmt.Sprintf("unknown currency %s", baseStr))
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
		security.Status = models.Trading
		security.Exchange = &constants.BYBITL
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
			ts, err := types.TimestampProto(date)
			if err != nil {
				state.logger.Info("error converting date")
				continue
			}
			security.MaturityDate = ts
		} else {
			security.SecurityType = enum.SecurityType_CRYPTO_PERP
		}
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinPriceIncrement = &types.DoubleValue{Value: symbol.PriceFilter.TickSize}
		security.RoundLot = &types.DoubleValue{Value: symbol.LotSizeFilter.QuantityStep}
		security.IsInverse = false
		security.MakerFee = &types.DoubleValue{Value: symbol.MakerFee}
		security.TakerFee = &types.DoubleValue{Value: symbol.TakerFee}
		security.Multiplier = &types.DoubleValue{Value: 1}
		securities = append(securities, &security)
	}

	state.securities = make(map[uint64]*models.Security)
	state.symbolToSec = make(map[string]*models.Security)
	for _, s := range securities {
		state.securities[s.SecurityID] = s
		state.symbolToSec[s.Symbol] = s
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
		response.RejectionReason = messages.UnknownSecurityID
		context.Respond(response)
		return nil
	}
	securityID := msg.Instrument.SecurityID.Value
	security, ok := state.securities[securityID]
	if !ok {
		response.Success = false
		response.RejectionReason = messages.UnknownSecurityID
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
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		fmt.Println(string(bres))

		/*
		if bres.ReturnCode != 0 {
			state.logger.Info("http error", log.Error(errors.New(bres.ReturnMessage)))
			response.RejectionReason = messages.ExchangeAPIError
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
	context.Respond(&messages.MarketStatisticsResponse{
		RequestID:       msg.RequestID,
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
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
		response.RejectionReason = messages.UnsupportedSubscription
		context.Respond(response)
		return nil
	}
	req, rate, err := bybitl.GetBalance("", msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}
	qr.rateLimit.Request(rate)
	f := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: req}, 15*time.Second)
	context.AwaitFuture(f, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.HTTPError
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
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
			} else if qResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"%d %s",
					qResponse.StatusCode,
					string(qResponse.Response))
				state.logger.Info("http server error", log.Error(err))
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
			}
			return
		}
		var balances bybitl.BalanceResponse
		err = json.Unmarshal(qResponse.Response, &balances)
		if err != nil {
			response.RejectionReason = messages.HTTPError
			context.Respond(response)
			return
		}
		if balances.RetCode != 0 {
			state.logger.Info("error getting orders", log.Error(errors.New(balances.RetMsg)))
			response.RejectionReason = messages.ExchangeAPIError
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
		response.RejectionReason = messages.UnsupportedSubscription
		context.Respond(response)
		return nil
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

	req, rate, err := bybitl.GetPositions("", msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}
	qr.rateLimit.Request(rate)
	f := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: req}, 15*time.Second)
	context.AwaitFuture(f, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.HTTPError
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
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
			} else if qResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"%d %s",
					qResponse.StatusCode,
					string(qResponse.Response))
				state.logger.Info("http server error", log.Error(err))
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
			}
			return
		}
		var positions bybitl.GetPositionsResponse
		err = json.Unmarshal(qResponse.Response, &positions)
		if err != nil {
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		if positions.RetCode != 0 {
			state.logger.Info("error getting orders", log.Error(errors.New(positions.RetMsg)))
			response.RejectionReason = messages.ExchangeAPIError
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
			sec, ok := state.symbolToSec[pos.Position.Symbol]
			if !ok {
				response.RejectionReason = messages.ExchangeAPIError
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
					Exchange:   &constants.BYBITL,
					Symbol:     &types.StringValue{Value: pos.Position.Symbol},
					SecurityID: &types.UInt64Value{Value: sec.SecurityID},
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

func (state *Executor) OnOrderStatusRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderStatusRequest)
	response := &messages.OrderList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
		Orders:     nil,
	}
	if msg.Subscribe {
		response.RejectionReason = messages.UnsupportedSubscription
		context.Respond(response)
		return nil
	}
	symbol := ""
	orderId := ""
	clOrderId := ""
	var orderStatus bybitl.OrderStatus
	if msg.Filter != nil {
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
		} else {
			response.RejectionReason = messages.MissingInstrument
			context.Respond(response)
			return nil
		}
		if msg.Filter.OrderStatus != nil {
			orderStatus = statusToBybitl(msg.Filter.OrderStatus.Value)
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
	req, rate, err := bybitl.QueryActiveOrdersRT(params, msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}
	qr.rateLimit.Request(rate)
	f := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: req}, 15*time.Second)
	context.AwaitFuture(f, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.HTTPError
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
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
			} else if qResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"%d %s",
					qResponse.StatusCode,
					string(qResponse.Response))
				state.logger.Info("http server error", log.Error(err))
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
			}
			return
		}
		var orders bybitl.QueryActiveOrdersResponse
		err = json.Unmarshal(qResponse.Response, &orders)
		if err != nil {
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		if orders.RetCode != 0 {
			state.logger.Info("error getting orders", log.Error(errors.New(orders.RetMsg)))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		for _, ord := range orders.Orders {
			if orderStatus != "" && ord.OrderStatus != orderStatus {
				continue
			}
			sec, ok := state.symbolToSec[ord.Symbol]
			if !ok {
				response.RejectionReason = messages.UnknownSymbol
				context.Respond(response)
				return
			}
			o := orderToModel(&ord)
			o.Instrument.SecurityID = &types.UInt64Value{Value: sec.SecurityID}
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
	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RateLimitExceeded
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
				response.RejectionReason = messages.UnknownSecurityID
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

	params, rej := buildPostOrderRequest(symbol, req.Order, tickPrecision, lotPrecision)
	if rej != nil {
		response.RejectionReason = *rej
		context.Respond(response)
		return nil
	}

	request, weight, err := bybitl.PostActiveOrder(params, req.Account.ApiCredentials)
	if err != nil {
		return err
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
		var order bybitl.PostOrderResponse
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		if order.RetCode != 0 {
			state.logger.Info("error posting order", log.Error(errors.New(fmt.Sprintf("%d: ", order.RetCode)+order.RetMsg)))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
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
	symbol := ""
	if req.Instrument != nil {
		if req.Instrument.Symbol != nil {
			symbol = req.Instrument.Symbol.Value
		} else if req.Instrument.SecurityID != nil {
			sec, ok := state.securities[req.Instrument.SecurityID.Value]
			if !ok {
				response.RejectionReason = messages.UnknownSecurityID
				context.Respond(response)
				return nil
			}
			symbol = sec.Symbol
		}
	} else {
		response.RejectionReason = messages.UnknownSecurityID
		context.Respond(response)
		return nil
	}
	params := bybitl.NewCancelActiveOrderParams(symbol)
	if req.OrderID != nil {
		params.SetOrderId(req.OrderID.Value)
	} else if req.ClientOrderID != nil {
		params.SetOrderLinkId(req.ClientOrderID.Value)
	} else {
		response.RejectionReason = messages.UnknownOrder
		context.Respond(response)
		return nil
	}

	request, weight, err := bybitl.CancelActiveOrder(params, req.Account.ApiCredentials)
	if err != nil {
		return err
	}

	fmt.Println(request.URL)
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
		var order bybitl.CancelOrderResponse
		fmt.Println(string(queryResponse.Response))
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		if order.RetCode != 0 {
			state.logger.Info("error cancelling order", log.Error(errors.New(order.RetMsg)))
			response.RejectionReason = messages.ExchangeAPIError
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
	symbol := ""
	if req.Filter != nil {
		if req.Filter.Instrument != nil {
			if req.Filter.Instrument.Symbol != nil {
				symbol = req.Filter.Instrument.Symbol.Value
			} else if req.Filter.Instrument.SecurityID != nil {
				sec, ok := state.securities[req.Filter.Instrument.SecurityID.Value]
				if !ok {
					response.RejectionReason = messages.UnknownSecurityID
					context.Respond(response)
					return nil
				}
				symbol = sec.Symbol
			}
		} else {
			response.RejectionReason = messages.MissingInstrument
			context.Respond(response)
			return nil
		}
	} else {
		response.RejectionReason = messages.UnsupportedFilter
		context.Respond(response)
		return nil
	}

	params := bybitl.NewCancelAllActiveParams(symbol)
	request, weight, err := bybitl.CancelAllActiveOrders(params, req.Account.ApiCredentials)
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

	f := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
	context.AwaitFuture(f, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.HTTPError
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
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf("%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http server error", log.Error(err))
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
			}
			return
		}
		var cancelResponse bybitl.CancelOrdersResponse
		err = json.Unmarshal(queryResponse.Response, &cancelResponse)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		if cancelResponse.RetCode != 0 {
			state.logger.Info("error mass cancelling orders", log.Error(errors.New(cancelResponse.RetMsg)))
			response.RejectionReason = messages.ExchangeAPIError
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

	symbol := ""
	if req.Instrument != nil {
		if req.Instrument.Symbol != nil {
			symbol = req.Instrument.Symbol.Value
		} else if req.Instrument.SecurityID != nil {
			sec, ok := state.securities[req.Instrument.SecurityID.Value]
			if !ok {
				response.RejectionReason = messages.UnknownSecurityID
				context.Respond(response)
				return nil
			}
			symbol = sec.Symbol
		}
	} else {
		response.RejectionReason = messages.MissingInstrument
		context.Respond(response)
		return nil
	}
	params := bybitl.NewAmendOrderParams(symbol)
	if req.Update.OrderID != nil {
		params.SetOrderId(req.Update.OrderID.Value)
	} else if req.Update.OrigClientOrderID != nil {
		params.SetOrderLinkId(req.Update.OrigClientOrderID.Value)
	} else {
		response.RejectionReason = messages.UnknownOrder
		context.Respond(response)
		return nil
	}
	request, weight, err := bybitl.AmendOrder(params, req.Account.ApiCredentials)
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

	f := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
	context.AwaitFuture(f, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.HTTPError
			context.Respond(response)
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				err := fmt.Errorf("%d %s", queryResponse.StatusCode, string(queryResponse.Response))
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
				return
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf("%d %s", queryResponse.StatusCode, string(queryResponse.Response))
				state.logger.Info("http server error", log.Error(err))
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
				return
			}
		}
		var amendResponse bybitl.AmendOrderResponse
		err = json.Unmarshal(queryResponse.Response, &amendResponse)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		if amendResponse.RetCode != 0 {
			state.logger.Info("error amending order", log.Error(errors.New(amendResponse.RetMsg)))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
		response.OrderID = amendResponse.AmendedOrder.OrderId
		context.Respond(response)
	})
	return nil
}
