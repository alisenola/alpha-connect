package ftxus

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
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/ftxus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
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

func NewExecutor(config *extypes.ExecutorConfig) actor.Actor {
	e := &Executor{}
	e.ExecutorConfig = config
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
	request, weight, err := ftxus.GetMarkets()
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
		Success bool           `json:"success"`
		Markets []ftxus.Market `json:"result"`
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
		security.Exchange = constants.FTXUS
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
					//fmt.Printf("unknown currency symbol %s \n", market.Underlying)
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
			request, weight, err = ftxus.GetOrderStatus(orderID, msg.Account.ApiCredentials)
			if err != nil {
				return err
			}
		} else {
			var err error
			request, weight, err = ftxus.GetOrderStatusByClientID(orderID, msg.Account.ApiCredentials)
			if err != nil {
				return err
			}
		}
	} else {
		multi = true
		var err error
		request, weight, err = ftxus.GetOpenOrders(symbol, msg.Account.ApiCredentials)
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
		var orders []ftxus.Order
		if multi {
			var oor ftxus.OpenOrdersResponse
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
			var oor ftxus.OrderStatusResponse
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
			orders = []ftxus.Order{oor.Result}
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
			ord := models.Order{
				OrderID:       fmt.Sprintf("%d", o.ID),
				ClientOrderID: o.ClientID,
				Instrument: &models.Instrument{
					Exchange:   constants.FBINANCE,
					Symbol:     &wrapperspb.StringValue{Value: o.Market},
					SecurityID: &wrapperspb.UInt64Value{Value: sec.SecurityID},
				},
				LeavesQuantity: o.Size - o.FilledSize, // TODO check
				CumQuantity:    o.FilledSize,
			}

			switch o.Status {
			case ftxus.NEW_ORDER:
				ord.OrderStatus = models.OrderStatus_PendingNew
			case ftxus.OPEN_ORDER:
				if o.FilledSize > 0 {
					ord.OrderStatus = models.OrderStatus_PartiallyFilled
				} else {
					ord.OrderStatus = models.OrderStatus_New
				}
			case ftxus.CLOSED_ORDER:
				if o.FilledSize == o.Size {
					ord.OrderStatus = models.OrderStatus_Filled
				} else {
					ord.OrderStatus = models.OrderStatus_Canceled
				}
			default:
				fmt.Println("unknown ORDER STATUS", o.Status)
			}

			switch o.Type {
			case ftxus.LIMIT_ORDER:
				ord.OrderType = models.OrderType_Limit
			case ftxus.MARKET_ORDER:
				ord.OrderType = models.OrderType_Market
			default:
				fmt.Println("UNKNOWN ORDER TYPE", o.Type)
			}

			switch o.Side {
			case ftxus.BUY_ORDER:
				ord.Side = models.Side_Buy
			case ftxus.SELL_ORDER:
				ord.Side = models.Side_Sell
			default:
				fmt.Println("UNKNOWN ORDER SIDE", o.Side)
			}

			if o.Ioc {
				ord.TimeInForce = models.TimeInForce_ImmediateOrCancel
			} else {
				ord.TimeInForce = models.TimeInForce_GoodTillCancel
			}
			ord.Price = &wrapperspb.DoubleValue{Value: o.Price}

			morders = append(morders, &ord)
		}
		response.Success = true
		response.Orders = morders
		fmt.Println(morders)
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

	request, weight, err := ftxus.GetPositions(true, msg.Account.ApiCredentials)
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
		var positions ftxus.PositionsResponse
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
			fmt.Println(p)
			if p.Size == 0 {
				continue
			}
			if symbol != "" && p.Future != symbol {
				continue
			}

			sec := state.SymbolToSecurity(p.Future)
			if sec == nil {
				state.logger.Warn(fmt.Sprintf("unknown symbol %s", p.Future))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			cost := *p.RecentAverageOpenPrice * sec.Multiplier.Value * p.Size
			size := p.Size
			if p.Side == "sell" {
				size *= -1
			}
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

	request, weight, err := ftxus.GetBalances(msg.Account.ApiCredentials)
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
		var balances ftxus.BalancesResponse
		err = json.Unmarshal(queryResponse.Response, &balances)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if !balances.Success {
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		for _, b := range balances.Result {
			fmt.Println(b)
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
			fmt.Println(response.Balances)
		}

		response.Success = true
		context.Respond(response)
	})

	return nil
}

func buildPlaceOrderRequest(symbol string, order *messages.NewOrder, tickPrecision, lotPrecision int) (ftxus.NewOrderRequest, *messages.RejectionReason) {
	request := ftxus.NewOrderRequest{
		Market: symbol,
	}

	if order.OrderSide == models.Side_Buy {
		request.Side = ftxus.BUY_ORDER
	} else {
		request.Side = ftxus.SELL_ORDER
	}
	switch order.OrderType {
	case models.OrderType_Limit:
		request.Type = ftxus.LIMIT_ORDER
	case models.OrderType_Market:
		request.Type = ftxus.MARKET_ORDER
	default:
		rej := messages.RejectionReason_UnsupportedOrderType
		return request, &rej
	}

	request.Size = order.Quantity
	request.ClientID = order.ClientOrderID

	if order.OrderType != models.OrderType_Market {
		switch order.TimeInForce {
		case models.TimeInForce_ImmediateOrCancel:
			request.Ioc = true
		case models.TimeInForce_GoodTillCancel:
			// Skip
		default:
			rej := messages.RejectionReason_UnsupportedOrderTimeInForce
			return request, &rej
		}
	}

	if order.Price != nil {
		request.Price = &order.Price.Value
	}

	for _, exec := range order.ExecutionInstructions {
		switch exec {
		case models.ExecutionInstruction_ReduceOnly:
			request.ReduceOnly = true
		case models.ExecutionInstruction_ParticipateDoNotInitiate:
			request.PostOnly = true
		default:
			rej := messages.RejectionReason_UnsupportedOrderCharacteristic
			return request, &rej
		}
	}

	return request, nil
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
		sec, rej := state.InstrumentToSecurity(req.Order.Instrument)
		if rej != nil {
			response.RejectionReason = *rej
			context.Send(sender, response)
			return nil
		}
		tickPrecision = int(math.Ceil(math.Log10(1. / sec.MinPriceIncrement.Value)))
		lotPrecision = int(math.Ceil(math.Log10(1. / sec.RoundLot.Value)))
	} else {
		response.RejectionReason = messages.RejectionReason_UnknownSecurityID
		context.Respond(response)
		return nil
	}

	params, rej := buildPlaceOrderRequest(symbol, req.Order, tickPrecision, lotPrecision)
	if rej != nil {
		response.RejectionReason = *rej
		context.Respond(response)
		return nil
	}

	request, weight, err := ftxus.PlaceOrder(params, req.Account.ApiCredentials)
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
		var order ftxus.NewOrderResponse
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if !order.Success {
			fmt.Println(order.Error)
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		fmt.Println("ORDER", order, order.Result.ID)
		response.Success = true
		response.OrderID = fmt.Sprintf("%d", order.Result.ID)
		context.Respond(response)
	})
	return nil
}

func (state *Executor) OnNewOrderBulkRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnOrderReplaceRequest(context actor.Context) error {
	fmt.Println("ORDER REPLACE REQUEST")
	req := context.Message().(*messages.OrderReplaceRequest)
	response := &messages.OrderReplaceResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	params := ftxus.ModifyOrderRequest{}
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

	request, weight, err := ftxus.ModifyOrder(req.Update.OrderID.Value, params, req.Account.ApiCredentials)
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
		var order ftxus.ModifyOrderResponse
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
	fmt.Println("ON ORDER CANCEL REQUEST")

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
		request, weight, err = ftxus.CancelOrder(req.OrderID.Value, req.Account.ApiCredentials)
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

	fmt.Println(request.URL)
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
		var order ftxus.CancelOrderResponse
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		if !order.Success {
			state.logger.Info("error unmarshalling", log.String("error", order.Error))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
		context.Respond(response)
	})
	return nil
}
