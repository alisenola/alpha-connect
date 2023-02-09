package bitmex

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
	"gitlab.com/alphaticks/xchanger/exchanges/bitmex"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io/ioutil"
	"net/http"
	"reflect"
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
type Executor struct {
	extypes.BaseExecutor
	client      *http.Client
	rateLimit   *exchanges.RateLimit
	queryRunner *actor.PID
	logger      *log.Logger
}

func NewExecutor() actor.Actor {
	return &Executor{
		client:      nil,
		rateLimit:   nil,
		queryRunner: nil,
		logger:      nil,
	}
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
			IdleConnTimeout:     0,
		},
		Timeout: 10 * time.Second,
	}
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(state).String()))

	state.rateLimit = exchanges.NewRateLimit(1, time.Second)
	// Launch an APIQuery actor with the given request and target
	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewHTTPQuery(state.client)
	})
	state.queryRunner = context.Spawn(props)
	state.rateLimit = exchanges.NewRateLimit(60, time.Minute)
	return state.UpdateSecurityList(context)
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	request, weight, err := bitmex.GetActiveInstruments()
	if err != nil {
		return err
	}
	if state.rateLimit.IsRateLimited() {
		return fmt.Errorf("rate limited")
	}
	state.rateLimit.Request(weight)

	resp, err := state.client.Do(request)
	if err != nil {
		return err
	}
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			err := fmt.Errorf(
				"http client error: %d %s",
				resp.StatusCode,
				string(response))
			return err
		} else if resp.StatusCode >= 500 {
			err := fmt.Errorf(
				"http server error: %d %s",
				resp.StatusCode,
				string(response))
			return err
		} else {
			err := fmt.Errorf("%d %s",
				resp.StatusCode,
				string(response))
			return err
		}
	}

	var activeInstruments []bitmex.Instrument
	err = json.Unmarshal(response, &activeInstruments)
	if err != nil {
		err = fmt.Errorf("error decoding query response: %v", err)
		return err
	}
	var securities []*models.Security
	for _, activeInstrument := range activeInstruments {
		if activeInstrument.State != "Open" {
			continue
		}
		switch activeInstrument.Typ {
		case "FFWCSX", "FFCCSX":
			symbolStr := strings.ToUpper(activeInstrument.Underlying)
			if sym, ok := bitmex.BITMEX_SYMBOL_TO_GLOBAL_SYMBOL[symbolStr]; ok {
				symbolStr = sym
			}
			baseCurrency, ok := constants.GetAssetBySymbol(symbolStr)
			if !ok {
				continue
			}
			symbolStr = strings.ToUpper(activeInstrument.QuoteCurrency)
			if sym, ok := bitmex.BITMEX_SYMBOL_TO_GLOBAL_SYMBOL[symbolStr]; ok {
				symbolStr = sym
			}
			quoteCurrency, ok := constants.GetAssetBySymbol(symbolStr)
			if !ok {
				continue
			}

			security := models.Security{}
			security.Symbol = activeInstrument.Symbol
			security.Underlying = baseCurrency
			security.QuoteCurrency = quoteCurrency
			switch activeInstrument.State {
			case "Open":
				security.Status = models.InstrumentStatus_Trading
			default:
				security.Status = models.InstrumentStatus_Disabled
			}
			security.Exchange = constants.BITMEX
			if activeInstrument.Typ == "FFCCSX" {
				security.SecurityType = enum.SecurityType_CRYPTO_FUT
				security.MaturityDate = timestamppb.New(activeInstrument.Settle)
			} else {
				security.SecurityType = enum.SecurityType_CRYPTO_PERP
			}
			security.MinPriceIncrement = &wrapperspb.DoubleValue{Value: activeInstrument.TickSize}
			security.RoundLot = &wrapperspb.DoubleValue{Value: float64(activeInstrument.LotSize)}
			security.Multiplier = &wrapperspb.DoubleValue{Value: float64(activeInstrument.Multiplier) * 0.00000001}
			security.MakerFee = &wrapperspb.DoubleValue{Value: activeInstrument.MakerFee}
			security.TakerFee = &wrapperspb.DoubleValue{Value: activeInstrument.TakerFee}
			security.IsInverse = activeInstrument.IsInverse
			security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
			securities = append(securities, &security)

		default:
			// case "OCECCS":
			//instr.Type = exchanges.CALL_OPTION
			// case "OPECCS":
			//instr.Type = exchanges.PUT_OPTION
			// Non-supported instrument, passing..
			continue
		}

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
	filters := make(map[string]interface{})
	params := bitmex.NewGetOrderParams()
	if msg.Filter != nil {
		if msg.Filter.Instrument != nil {
			if msg.Filter.Instrument.Symbol != nil {
				params.SetSymbol(msg.Filter.Instrument.Symbol.Value)
			} else if msg.Filter.Instrument.SecurityID != nil {
				sec := state.IDToSecurity(msg.Filter.Instrument.SecurityID.Value)
				if sec == nil {
					response.RejectionReason = messages.RejectionReason_UnknownSecurityID
					context.Respond(response)
					return nil
				}
				params.SetSymbol(sec.Symbol)
			}
		}
		if msg.Filter.OrderID != nil {
			filters["orderID"] = msg.Filter.OrderID.Value
		}
		if msg.Filter.ClientOrderID != nil {
			filters["clOrdID"] = msg.Filter.ClientOrderID.Value
		}
	}
	filters["open"] = true

	if len(filters) > 0 {
		params.SetFilters(filters)
	}

	request, weight, err := bitmex.GetOrder(msg.Account.ApiCredentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

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
		var orders []bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &orders)
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		var morders []*models.Order
		for _, o := range orders {
			sec := state.SymbolToSecurity(o.Symbol)
			if sec == nil {
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return
			}
			ord := models.Order{
				OrderID: o.OrderID,
				Instrument: &models.Instrument{
					Exchange:   constants.BITMEX,
					Symbol:     &wrapperspb.StringValue{Value: o.Symbol},
					SecurityID: &wrapperspb.UInt64Value{Value: sec.SecurityID},
				},
				LeavesQuantity: float64(o.LeavesQty),
				CumQuantity:    float64(o.CumQty),
			}

			if o.ClOrdID != nil {
				ord.ClientOrderID = *o.ClOrdID
			}

			switch o.OrdStatus {
			case "New":
				ord.OrderStatus = models.OrderStatus_New
			case "Canceled":
				ord.OrderStatus = models.OrderStatus_Canceled
			case "Filled":
				ord.OrderStatus = models.OrderStatus_Filled
			default:
				fmt.Println("UNKNWOEN ORDER STATUS", o.OrdStatus)
			}

			switch bitmex.OrderType(o.OrdType) {
			case bitmex.LIMIT:
				ord.OrderType = models.OrderType_Limit
			case bitmex.MARKET:
				ord.OrderType = models.OrderType_Market
			case bitmex.STOP:
				ord.OrderType = models.OrderType_Stop
			case bitmex.STOP_LIMIT:
				ord.OrderType = models.OrderType_StopLimit
			case bitmex.LIMIT_IF_TOUCHED:
				ord.OrderType = models.OrderType_LimitIfTouched
			case bitmex.MARKET_IF_TOUCHED:
				ord.OrderType = models.OrderType_MarketIfTouched
			default:
				fmt.Println("UNKNOWN ORDER TYPE", o.OrdType)
			}

			switch bitmex.OrderSide(o.Side) {
			case bitmex.BUY_ORDER_SIDE:
				ord.Side = models.Side_Buy
			case bitmex.SELL_ORDER_SIDE:
				ord.Side = models.Side_Sell
			default:
				fmt.Println("UNKNOWN ORDER SIDE", o.Side)
			}

			switch bitmex.TimeInForce(o.TimeInForce) {
			case bitmex.TIF_DAY:
				ord.TimeInForce = models.TimeInForce_Session
			case bitmex.GOOD_TILL_CANCEL:
				ord.TimeInForce = models.TimeInForce_GoodTillCancel
			case bitmex.IMMEDIATE_OR_CANCEL:
				ord.TimeInForce = models.TimeInForce_ImmediateOrCancel
			case bitmex.FILL_OR_KILL:
				ord.TimeInForce = models.TimeInForce_FillOrKill
			default:
				fmt.Println("UNKNOWN TOF", o.TimeInForce)
			}

			ord.Price = &wrapperspb.DoubleValue{Value: o.Price}

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

	positionList := &messages.PositionList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
		Positions:  nil,
	}

	params := bitmex.NewGetPositionParams()
	filters := make(map[string]interface{})

	if msg.Instrument != nil {
		if msg.Instrument.Symbol != nil {
			filters["symbol"] = msg.Instrument.Symbol.Value
		} else if msg.Instrument.SecurityID != nil {
			sec := state.IDToSecurity(msg.Instrument.SecurityID.Value)
			if sec == nil {
				positionList.RejectionReason = messages.RejectionReason_UnknownSecurityID
				context.Respond(positionList)
				return nil
			}
			filters["symbol"] = sec.Symbol
		}
	}
	if len(filters) > 0 {
		params.SetFilters(filters)
	}

	request, weight, err := bitmex.GetPosition(msg.Account.ApiCredentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		positionList.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Respond(positionList)
		return nil
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

	context.ReenterAfter(future, func(res interface{}, err error) {
		if err != nil {
			positionList.RejectionReason = messages.RejectionReason_Other
			context.Respond(positionList)
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
				positionList.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(positionList)
			} else if queryResponse.StatusCode >= 500 {

				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http error", log.Error(err))
				positionList.Success = false
				positionList.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(positionList)
			}
			return
		}
		var positions []bitmex.Position
		err = json.Unmarshal(queryResponse.Response, &positions)
		if err != nil {
			state.logger.Info("unmarshaling error", log.Error(err))
			positionList.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(positionList)
			return
		}
		fmt.Println(string(queryResponse.Response))
		for _, p := range positions {
			if p.CurrentQty == 0 {
				continue
			}
			sec := state.SymbolToSecurity(p.Symbol)
			if sec == nil {
				positionList.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(positionList)
				return
			}
			quantity := float64(p.CurrentQty)
			cost := float64(p.UnrealisedCost) * 0.00000001
			pos := &models.Position{
				Account: msg.Account.Name,
				Instrument: &models.Instrument{
					Exchange:   constants.BITMEX,
					Symbol:     &wrapperspb.StringValue{Value: p.Symbol},
					SecurityID: &wrapperspb.UInt64Value{Value: sec.SecurityID},
				},
				Quantity: quantity,
				Cost:     cost,
				Cross:    p.CrossMargin,
			}
			positionList.Positions = append(positionList.Positions, pos)
		}
		positionList.Success = true
		context.Respond(positionList)
	})

	return nil
}

func (state *Executor) OnBalancesRequest(context actor.Context) error {
	msg := context.Message().(*messages.BalancesRequest)

	balanceList := &messages.BalanceList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
		Balances:   nil,
	}

	request, weight, err := bitmex.GetMargin(msg.Account.ApiCredentials, "XBt")
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		balanceList.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Respond(balanceList)
		return nil
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

	context.ReenterAfter(future, func(res interface{}, err error) {
		if err != nil {
			balanceList.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(balanceList)
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
				balanceList.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(balanceList)
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http server error", log.Error(err))
				balanceList.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(balanceList)
			}
			return
		}
		var margin bitmex.Margin
		err = json.Unmarshal(queryResponse.Response, &margin)
		if err != nil {
			balanceList.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(balanceList)
			return
		}
		balanceList.Balances = append(balanceList.Balances, &models.Balance{
			Account:  msg.Account.Name,
			Asset:    constants.BITCOIN,
			Quantity: float64(margin.WalletBalance) * 0.00000001,
		})
		balanceList.Success = true
		context.Respond(balanceList)
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
	if req.Order.Instrument == nil || req.Order.Instrument.Symbol == nil {
		response.RejectionReason = messages.RejectionReason_UnknownSymbol
		context.Respond(response)
		return nil
	}

	params, rej := buildPostOrderRequest(req.Order)
	if rej != nil {
		response.RejectionReason = *rej
		context.Respond(response)
		return nil
	}

	request, weight, err := bitmex.PostOrder(req.Account.ApiCredentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.rateLimit.Request(weight)
	requestStart := time.Now()
	future := context.RequestFuture(state.queryRunner, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
	context.ReenterAfter(future, func(res interface{}, err error) {
		requestDur := time.Since(requestStart)
		state.logger.Info(fmt.Sprintf("post order request done in %v", requestDur))
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
		var order bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		response.Success = true
		response.OrderID = order.OrderID
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
	params := bitmex.NewAmendOrderRequest()
	if req.Update.OrderID != nil {
		params.SetOrderID(req.Update.OrderID.Value)
	} else if req.Update.OrigClientOrderID != nil {
		params.SetOrigClOrdID(req.Update.OrigClientOrderID.Value)
	}

	if req.Update.Price != nil {
		params.SetPrice(req.Update.Price.Value)
	}
	if req.Update.Quantity != nil {
		params.SetOrderQty(req.Update.Quantity.Value)
	}

	request, weight, err := bitmex.AmendOrder(req.Account.ApiCredentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Respond(response)
		return nil
	}

	future := context.RequestFuture(state.queryRunner, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
	context.ReenterAfter(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Respond(response)
			return
		}
		// Only update rate limit if successful request
		state.rateLimit.Request(weight)
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
		var order bitmex.Order
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

func (state *Executor) OnOrderCancelRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderCancelRequest)
	response := &messages.OrderCancelResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	params := bitmex.NewCancelOrderRequest()
	if req.OrderID != nil {
		params.SetOrderID(req.OrderID.Value)
	} else if req.ClientOrderID != nil {
		params.SetClOrdID(req.ClientOrderID.Value)
	}

	request, weight, err := bitmex.CancelOrder(req.Account.ApiCredentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
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
		var orders []bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &orders)
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
	filters := make(map[string]interface{})
	params := bitmex.NewCancelAllOrdersRequest()
	if req.Filter != nil {
		if req.Filter.Instrument != nil {
			if req.Filter.Instrument.Symbol != nil {
				if sec := state.SymbolToSecurity(req.Filter.Instrument.Symbol.Value); sec == nil {
					response.RejectionReason = messages.RejectionReason_UnknownSymbol
					context.Respond(response)
					return nil
				}
				params.SetSymbol(req.Filter.Instrument.Symbol.Value)
			} else if req.Filter.Instrument.SecurityID != nil {
				sec := state.IDToSecurity(req.Filter.Instrument.SecurityID.Value)
				if sec == nil {
					response.RejectionReason = messages.RejectionReason_UnknownSecurityID
					context.Respond(response)
					return nil
				}
				params.SetSymbol(sec.Symbol)
			}
		}
		if req.Filter.Side != nil {
			filters["side"] = req.Filter.Side.Value.String()
		}
	}

	request, weight, err := bitmex.CancelAllOrders(req.Account.ApiCredentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)
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
		var orders []bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &orders)
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
