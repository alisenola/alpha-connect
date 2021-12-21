package dydx

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
	"gitlab.com/alphaticks/xchanger/exchanges/dydx"
	"gitlab.com/alphaticks/xchanger/exchanges/fbinance"
	xutils "gitlab.com/alphaticks/xchanger/utils"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"time"
)

type QueryRunner struct {
	pid             *actor.PID
	globalRateLimit *exchanges.RateLimit
	getRateLimit    *exchanges.RateLimit
}

type AccountRateLimit struct {
	cancelOrderRateLimit *exchanges.RateLimit
	placeOrderRateLimit  *exchanges.RateLimit
}

type Executor struct {
	extypes.BaseExecutor
	securities        map[uint64]*models.Security
	symbolToSec       map[string]*models.Security
	accountRateLimits map[string]*AccountRateLimit
	dialerPool        *xutils.DialerPool
	queryRunners      []*QueryRunner
	logger            *log.Logger
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
		log.String("type", reflect.TypeOf(*state).String()))

	dialers := state.dialerPool.GetDialers()
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
			return jobs.NewAPIQuery(client)
		})
		state.queryRunners = append(state.queryRunners, &QueryRunner{
			pid:             context.Spawn(props),
			globalRateLimit: exchanges.NewRateLimit(10, time.Minute),
			getRateLimit:    exchanges.NewRateLimit(100, 10*time.Second),
		})
	}
	// TODO

	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	req, weight, err := dydx.GetMarkets()

	qr := state.getQueryRunner()
	if qr == nil {
		return fmt.Errorf("rate limited")
	}

	qr.getRateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformQueryRequest{Request: req}, 10*time.Second)

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

	var markets dydx.MarketsResponse
	err = json.Unmarshal(resp.Response, &markets)
	if err != nil {
		err = fmt.Errorf("error decoding query response: %v", err)
		return err
	}

	state.securities = nil
	var securities []*models.Security
	for _, m := range markets.Markets {
		security := models.Security{}
		//ONLINE, OFFLINE, POST_ONLY or CANCEL_ONLY.
		switch m.Status {
		case "ONLINE":
			security.Status = models.Trading
		case "OFFLINE":
			security.Status = models.Halt
		case "POST_ONLY":
			security.Status = models.PreTrading
		case "CANCEL_ONLY":
			security.Status = models.PostTrading
		}
		security.Symbol = m.Market
		baseStr := m.BaseAsset
		if sym, ok := dydx.SYMBOLS[baseStr]; ok {
			baseStr = sym
		}
		quoteStr := m.QuoteAsset
		if sym, ok := dydx.SYMBOLS[quoteStr]; ok {
			quoteStr = sym
		}
		baseCurrency, ok := constants.GetAssetBySymbol(baseStr)
		if !ok {
			continue
		}
		security.Underlying = baseCurrency
		quoteCurrency, ok := constants.GetAssetBySymbol(quoteStr)
		if !ok {
			continue
		}
		security.QuoteCurrency = quoteCurrency
		security.Exchange = &constants.DYDX
		security.SecurityType = enum.SecurityType_CRYPTO_PERP
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinPriceIncrement = &types.DoubleValue{Value: m.TickSize}
		security.RoundLot = &types.DoubleValue{Value: m.StepSize}
		security.IsInverse = false

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
		Securities: securities,
	})

	return nil
}

func (state *Executor) OnOrderStatusRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderStatusRequest)
	response := &messages.OrderList{
		RequestID: msg.RequestID,
		Success:   false,
	}
	var err error
	var request *http.Request
	var weight int
	var single bool
	if msg.Filter != nil {
		if msg.Filter.OrderID != nil || msg.Filter.ClientOrderID != nil {
			single = true
			if msg.Filter.OrderID != nil {
				request, weight, err = dydx.GetOrder(msg.Filter.OrderID.Value, msg.Account.ApiCredentials)
				if err != nil {
					return err
				}
			} else {
				request, weight, err = dydx.GetOrderClient(msg.Filter.ClientOrderID.Value, msg.Account.ApiCredentials)
				if err != nil {
					return err
				}
			}
		} else {
			single = false
			params := dydx.NewGetOrdersParams()
			if msg.Filter.OrderStatus != nil {
				status := statusToDydx(msg.Filter.OrderStatus.Value)
				if status == "" {
					response.RejectionReason = messages.UnsupportedFilter
					context.Respond(response)
					return nil
				}
				params.SetStatus(status)
			}
			if msg.Filter.Side != nil {
				side := sideToDydx(msg.Filter.Side.Value)
				if side == "" {
					response.RejectionReason = messages.UnsupportedFilter
					context.Respond(response)
					return nil
				}
				params.SetSide(side)
			}
			if msg.Filter.Instrument != nil {
				if msg.Filter.Instrument.Symbol != nil {
					params.SetMarket(msg.Filter.Instrument.Symbol.Value)
				} else if msg.Filter.Instrument.SecurityID != nil {
					sec, ok := state.securities[msg.Filter.Instrument.SecurityID.Value]
					if !ok {
						response.RejectionReason = messages.UnknownSecurityID
						context.Respond(response)
						return nil
					}
					params.SetMarket(sec.Symbol)
				}
			}
			request, weight, err = dydx.GetOrders(params, msg.Account.ApiCredentials)
			if err != nil {
				return err
			}
		}
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.getRateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

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
		var orders []dydx.Order
		if single {
			var order dydx.Order
			err = json.Unmarshal(queryResponse.Response, &order)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			orders = append(orders, order)
		} else {
			err = json.Unmarshal(queryResponse.Response, &orders)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
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
			ord := orderToModel(o)
			ord.Instrument.SecurityID = &types.UInt64Value{Value: sec.SecurityID}
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
				response.RejectionReason = messages.UnknownSecurityID
				context.Respond(response)
				return nil
			}
			symbol = sec.Symbol
		}
	}

	params := dydx.NewGetPositionsParams()
	if symbol != "" {
		params.SetMarket(symbol)
	}
	request, weight, err := dydx.GetPositions(params, msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.globalRateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

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
		var positions dydx.PositionsResponse
		err = json.Unmarshal(queryResponse.Response, &positions)
		if err != nil {
			state.logger.Info("unmarshalling error", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		if len(positions.Errors) > 0 {
			errStr := ""
			for _, e := range positions.Errors {
				errStr += e.Msg + " - "
			}
			state.logger.Info("api error", log.Error(errors.New(errStr)))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		for _, p := range positions.Positions {
			if p.Size == 0 {
				continue
			}
			sec, ok := state.symbolToSec[p.Market]
			if !ok {
				err := fmt.Errorf("unknown symbol %s", p.Market)
				state.logger.Info("unmarshalling error", log.Error(err))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			cost := p.Size * p.EntryPrice
			pos := &models.Position{
				Account: msg.Account.Name,
				Instrument: &models.Instrument{
					Exchange:   &constants.BITMEX,
					Symbol:     &types.StringValue{Value: p.Market},
					SecurityID: &types.UInt64Value{Value: sec.SecurityID},
				},
				Quantity: p.Size,
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

	request, weight, err := dydx.GetAccount(msg.Account.ApiCredentials)
	if err != nil {
		return err
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.globalRateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

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
		var account dydx.AccountResponse
		err = json.Unmarshal(queryResponse.Response, &account)
		if err != nil {
			response.RejectionReason = messages.HTTPError
			context.Respond(response)
			return
		}
		if len(account.Errors) > 0 {
			errStr := ""
			for _, e := range account.Errors {
				errStr += e.Msg + " - "
			}
			state.logger.Info("api error", log.Error(errors.New(errStr)))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Balances = append(response.Balances, &models.Balance{
			Account:  msg.Account.Name,
			Asset:    &constants.USDC,
			Quantity: account.Account.QuoteBalance,
		})

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

	arl, ok := state.accountRateLimits[req.Account.Name]
	if !ok {
		arl = &AccountRateLimit{
			cancelOrderRateLimit: exchanges.NewRateLimit(10, time.Second),
			placeOrderRateLimit:  exchanges.NewRateLimit(10, time.Second),
		}
		state.accountRateLimits[req.Account.Name] = arl
	}
	if arl.placeOrderRateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	symbol := ""
	if req.Order.Instrument != nil {
		if req.Order.Instrument.Symbol != nil {
			symbol = req.Order.Instrument.Symbol.Value
		} else if req.Order.Instrument.SecurityID != nil {
			sec, ok := state.securities[req.Order.Instrument.SecurityID.Value]
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

	params := buildCreateOrderParams(symbol, req.Order)
	request, weight, err := dydx.CreateOrder(params, req.Account.StarkCredentials, req.Account.ApiCredentials)
	if err != nil {
		return err
	}

	arl.placeOrderRateLimit.Request(weight)
	qr.globalRateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
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
		var order fbinance.OrderData
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			response.RejectionReason = messages.HTTPError
			context.Respond(response)
			return
		}
		response.Success = true
		response.OrderID = fmt.Sprintf("%d", order.OrderID)
		context.Respond(response)
	})
	return nil
}
