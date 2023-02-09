package dydx

import (
	"encoding/json"
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
	"gitlab.com/alphaticks/xchanger/exchanges/dydx"
	xutils "gitlab.com/alphaticks/xchanger/utils"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
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
	accountRateLimits map[string]*AccountRateLimit
	queryRunners      []*QueryRunner
	logger            *log.Logger
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
	// TODO filter
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
			globalRateLimit: exchanges.NewRateLimit(10, time.Minute),
			getRateLimit:    exchanges.NewRateLimit(100, 10*time.Second),
		})
	}
	state.accountRateLimits = make(map[string]*AccountRateLimit)
	// TODO

	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	req, weight, err := dydx.GetMarkets()
	if err != nil {
		return fmt.Errorf("error getting markets: %v", err)
	}

	qr := state.getQueryRunner()
	if qr == nil {
		return fmt.Errorf("rate limited")
	}

	qr.getRateLimit.Request(weight)
	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: req}, 10*time.Second)

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
	if len(markets.Errors) > 0 {
		err = fmt.Errorf("api error: %v", markets.Errors)
		return err
	}

	var securities []*models.Security
	for _, m := range markets.Markets {
		security := models.Security{}
		//ONLINE, OFFLINE, POST_ONLY or CANCEL_ONLY.
		switch m.Status {
		case "ONLINE":
			security.Status = models.InstrumentStatus_Trading
		case "OFFLINE":
			security.Status = models.InstrumentStatus_Halt
		case "POST_ONLY":
			security.Status = models.InstrumentStatus_PreTrading
		case "CANCEL_ONLY":
			security.Status = models.InstrumentStatus_PostTrading
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
		security.Exchange = constants.DYDX
		security.SecurityType = enum.SecurityType_CRYPTO_PERP
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinPriceIncrement = &wrapperspb.DoubleValue{Value: m.TickSize}
		security.Multiplier = &wrapperspb.DoubleValue{Value: 1}
		security.RoundLot = &wrapperspb.DoubleValue{Value: m.StepSize}
		security.IsInverse = false

		securities = append(securities, &security)
	}

	state.SyncSecurities(securities, nil)

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: securities})

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
					response.RejectionReason = messages.RejectionReason_UnsupportedFilter
					context.Respond(response)
					return nil
				}
				params.SetStatus(status)
			}
			if msg.Filter.Side != nil {
				side := sideToDydx(msg.Filter.Side.Value)
				if side == "" {
					response.RejectionReason = messages.RejectionReason_UnsupportedFilter
					context.Respond(response)
					return nil
				}
				params.SetSide(side)
			}
			if msg.Filter.Instrument != nil {
				symbol, rej := state.InstrumentToSymbol(msg.Filter.Instrument)
				if rej != nil {
					response.RejectionReason = *rej
					context.Respond(response)
					return nil
				}
				params.SetMarket(symbol)
			}
			request, weight, err = dydx.GetOrders(params, msg.Account.ApiCredentials)
			if err != nil {
				return err
			}
		}
	} else {
		single = false
		request, weight, err = dydx.GetOrders(dydx.NewGetOrdersParams(), msg.Account.ApiCredentials)
		if err != nil {
			return err
		}
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.getRateLimit.Request(weight)
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
		var orders []dydx.Order
		if single {
			var res dydx.OrderResponse
			err = json.Unmarshal(queryResponse.Response, &res)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			if len(res.Errors) > 0 {
				err = fmt.Errorf("%v", res.Errors)
				state.logger.Info("api error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			orders = append(orders, res.Order)
		} else {
			var res dydx.OrdersResponse
			err = json.Unmarshal(queryResponse.Response, &res)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			if len(res.Errors) > 0 {
				err = fmt.Errorf("%v", res.Errors)
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			orders = res.Orders
		}

		var morders []*models.Order
		for _, o := range orders {
			sec := state.SymbolToSecurity(o.Market)
			if sec == nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			ord := orderToModel(o)
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
		s, rej := state.InstrumentToSymbol(msg.Instrument)
		if rej != nil {
			response.RejectionReason = *rej
			context.Respond(response)
			return nil
		}
		symbol = s
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
		response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.getRateLimit.Request(weight)
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
		var positions dydx.PositionsResponse
		err = json.Unmarshal(queryResponse.Response, &positions)
		if err != nil {
			state.logger.Info("unmarshalling error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if len(positions.Errors) > 0 {
			err = fmt.Errorf("%v", positions.Errors)
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		fmt.Println("POSITIONS", positions)
		for _, p := range positions.Positions {
			if p.Size == 0 {
				continue
			}
			sec := state.SymbolToSecurity(p.Market)
			if sec == nil {
				state.logger.Warn(fmt.Sprintf("unknown symbol: %s", p.Market))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			// Calculate the cost using unrealized pnl
			//cost := p.Size
			cost := math.Round(1e6*(p.SumOpen*p.EntryPrice-p.SumClose*p.ExitPrice)) / 1e6
			pos := &models.Position{
				Account: msg.Account.Name,
				Instrument: &models.Instrument{
					Exchange:   constants.DYDX,
					Symbol:     &wrapperspb.StringValue{Value: p.Market},
					SecurityID: &wrapperspb.UInt64Value{Value: sec.SecurityID},
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
		response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Respond(response)
		return nil
	}

	qr.getRateLimit.Request(weight)
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
		var account dydx.AccountResponse
		err = json.Unmarshal(queryResponse.Response, &account)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		if len(account.Errors) > 0 {
			err = fmt.Errorf("%v", account.Errors)
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		qb := account.Account.QuoteBalance
		for _, p := range account.Account.OpenPositions {
			cost := math.Round(1e6*(p.SumOpen*p.EntryPrice-p.SumClose*p.ExitPrice)) / 1e6
			qb += cost
			fmt.Printf("FC: %f Unrealized %f QB: %f \n", account.Account.FreeCollateral, p.UnrealizedPnl, account.Account.QuoteBalance)
		}
		fmt.Println("MARGIN BALANCE", qb)
		response.Balances = append(response.Balances, &models.Balance{
			Account:  msg.Account.Name,
			Asset:    constants.USDC,
			Quantity: qb,
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
		response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
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
		response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Respond(response)
		return nil
	}

	sec, rej := state.InstrumentToSecurity(req.Order.Instrument)
	if rej != nil {
		response.RejectionReason = *rej
		context.Respond(response)
		return nil
	}

	params := buildCreateOrderParams(sec.Symbol, req.Order)
	request, weight, err := dydx.CreateOrder(params, req.Account.StarkCredentials, req.Account.ApiCredentials)
	if err != nil {
		return err
	}

	arl.placeOrderRateLimit.Request(weight)
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
		var order dydx.CreateOrderResponse
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		if len(order.Errors) > 0 {
			err = fmt.Errorf("%v", order.Errors)
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
		response.OrderID = order.Order.ID
		context.Respond(response)
	})
	return nil
}

func (state *Executor) OnNewOrderBulkRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnOrderReplaceRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnOrderBulkReplaceRequest(context actor.Context) error {
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
	var err error
	// TODO cancel by client order ID ?
	if req.OrderID != nil {
		request, weight, err = dydx.CancelOrder(req.OrderID.Value, req.Account.ApiCredentials)
		if err != nil {
			return err
		}
	} else {
		response.RejectionReason = messages.RejectionReason_UnknownOrder
		context.Respond(response)
		return nil
	}

	qr := state.getQueryRunner()
	if qr == nil {
		response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Respond(response)
		return nil
	}

	// TODO account cancel rate limit
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
		var order dydx.CancelOrderResponse
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if len(order.Errors) > 0 {
			err = fmt.Errorf("%v", order.Errors)
			state.logger.Info("http error", log.Error(err))
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
	return nil
}
