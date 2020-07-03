package bitmex

import (
	"encoding/json"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alphac/enum"
	"gitlab.com/alphaticks/alphac/exchanges/interface"
	"gitlab.com/alphaticks/alphac/jobs"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/bitmex"
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
	client      *http.Client
	securities  map[uint64]*models.Security
	symbolToSec map[string]*models.Security
	rateLimit   *exchanges.RateLimit
	queryRunner *actor.PID
	logger      *log.Logger
}

func NewExecutor() actor.Actor {
	return &Executor{
		client:      nil,
		securities:  nil,
		symbolToSec: nil,
		rateLimit:   nil,
		queryRunner: nil,
		logger:      nil,
	}
}

func (state *Executor) Receive(context actor.Context) {
	_interface.ExchangeExecutorReceive(state, context)
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

	state.rateLimit = exchanges.NewRateLimit(1, time.Second)
	// Launch an APIQuery actor with the given request and target
	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewAPIQuery(state.client)
	})
	state.queryRunner = context.Spawn(props)
	state.rateLimit = exchanges.NewRateLimit(60, time.Minute)
	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
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
		/*
			case "FFCCSX":
				instr.Type = exchanges.FUTURE
		*/
		case "FFWCSX":
			symbolStr := strings.ToUpper(activeInstrument.Underlying)
			if sym, ok := bitmex.BITMEX_SYMBOL_TO_GLOBAL_SYMBOL[symbolStr]; ok {
				symbolStr = sym
			}
			baseCurrency, ok := constants.SYMBOL_TO_ASSET[symbolStr]
			if !ok {
				continue
			}
			symbolStr = strings.ToUpper(activeInstrument.QuoteCurrency)
			if sym, ok := bitmex.BITMEX_SYMBOL_TO_GLOBAL_SYMBOL[symbolStr]; ok {
				symbolStr = sym
			}
			quoteCurrency, ok := constants.SYMBOL_TO_ASSET[symbolStr]
			if !ok {
				continue
			}

			security := models.Security{}
			security.Symbol = activeInstrument.Symbol
			security.Underlying = &baseCurrency
			security.QuoteCurrency = &quoteCurrency
			security.Enabled = activeInstrument.State == "Open"
			security.Exchange = &constants.BITMEX
			security.SecurityType = enum.SecurityType_CRYPTO_PERP
			security.MinPriceIncrement = activeInstrument.TickSize
			security.RoundLot = float64(activeInstrument.LotSize)
			security.Multiplier = &types.DoubleValue{Value: float64(activeInstrument.Multiplier) * 0.00000001}
			security.MakerFee = &types.DoubleValue{Value: activeInstrument.MakerFee}
			security.TakerFee = &types.DoubleValue{Value: activeInstrument.TakerFee}
			security.IsInverse = activeInstrument.IsInverse
			security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name)
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

func (state *Executor) OnMarketDataRequest(context actor.Context) error {
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
				sec, ok := state.securities[msg.Filter.Instrument.SecurityID.Value]
				if !ok {
					response.RejectionReason = messages.UnknownSecurityID
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

	request, weight, err := bitmex.GetOrder(msg.Account.Credentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

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
		var orders []bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &orders)
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}

		var morders []*models.Order
		for _, o := range orders {
			sec, ok := state.symbolToSec[o.Symbol]
			if !ok {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.ExchangeAPIError
				context.Respond(response)
				return
			}
			ord := models.Order{
				OrderID: o.OrderID,
				Instrument: &models.Instrument{
					Exchange:   &constants.BITMEX,
					Symbol:     &types.StringValue{Value: o.Symbol},
					SecurityID: &types.UInt64Value{Value: sec.SecurityID},
				},
				LeavesQuantity: float64(o.LeavesQty),
				CumQuantity:    float64(o.CumQty),
			}

			if o.ClOrdID != nil {
				ord.ClientOrderID = *o.ClOrdID
			}

			switch o.OrdStatus {
			case "New":
				ord.OrderStatus = models.New
			case "Canceled":
				ord.OrderStatus = models.Canceled
			case "Filled":
				ord.OrderStatus = models.Filled
			default:
				fmt.Println("UNKNWOEN ORDER STATUS", o.OrdStatus)
			}

			switch bitmex.OrderType(o.OrdType) {
			case bitmex.LIMIT:
				ord.OrderType = models.Limit
			case bitmex.MARKET:
				ord.OrderType = models.Market
			case bitmex.STOP:
				ord.OrderType = models.Stop
			case bitmex.STOP_LIMIT:
				ord.OrderType = models.StopLimit
			case bitmex.LIMIT_IF_TOUCHED:
				ord.OrderType = models.LimitIfTouched
			case bitmex.MARKET_IF_TOUCHED:
				ord.OrderType = models.MarketIfTouched
			default:
				fmt.Println("UNKNOWN ORDER TYPE", o.OrdType)
			}

			switch bitmex.OrderSide(o.Side) {
			case bitmex.BUY_ORDER_SIDE:
				ord.Side = models.Buy
			case bitmex.SELL_ORDER_SIDE:
				ord.Side = models.Sell
			default:
				fmt.Println("UNKNOWN ORDER SIDE", o.Side)
			}

			switch bitmex.TimeInForce(o.TimeInForce) {
			case bitmex.TIF_DAY:
				ord.TimeInForce = models.Session
			case bitmex.GOOD_TILL_CANCEL:
				ord.TimeInForce = models.GoodTillCancel
			case bitmex.IMMEDIATE_OR_CANCEL:
				ord.TimeInForce = models.ImmediateOrCancel
			case bitmex.FILL_OR_KILL:
				ord.TimeInForce = models.FillOrKill
			default:
				fmt.Println("UNKNOWN TOF", o.TimeInForce)
			}

			ord.Price = &types.DoubleValue{Value: o.Price}

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
			sec, ok := state.securities[msg.Instrument.SecurityID.Value]
			if !ok {
				positionList.RejectionReason = messages.UnknownSecurityID
				context.Respond(positionList)
				return nil
			}
			filters["symbol"] = sec.Symbol
		}
	}
	if len(filters) > 0 {
		params.SetFilters(filters)
	}

	request, weight, err := bitmex.GetPosition(msg.Account.Credentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		positionList.RejectionReason = messages.RateLimitExceeded
		context.Respond(positionList)
		return nil
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			positionList.RejectionReason = messages.Other
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
				positionList.RejectionReason = messages.ExchangeAPIError
				context.Respond(positionList)
			} else if queryResponse.StatusCode >= 500 {

				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http error", log.Error(err))
				positionList.Success = false
				positionList.RejectionReason = messages.ExchangeAPIError
				context.Respond(positionList)
			}
			return
		}
		var positions []bitmex.Position
		err = json.Unmarshal(queryResponse.Response, &positions)
		if err != nil {
			state.logger.Info("unmarshaling error", log.Error(err))
			positionList.RejectionReason = messages.ExchangeAPIError
			context.Respond(positionList)
			return
		}
		fmt.Println(string(queryResponse.Response))
		for _, p := range positions {
			if p.CurrentQty == 0 {
				continue
			}
			sec, ok := state.symbolToSec[p.Symbol]
			if !ok {
				positionList.RejectionReason = messages.ExchangeAPIError
				context.Respond(positionList)
				return
			}
			quantity := float64(p.CurrentQty)
			cost := float64(p.UnrealisedCost) * 0.00000001
			pos := &models.Position{
				AccountID: msg.Account.AccountID,
				Instrument: &models.Instrument{
					Exchange:   &constants.BITMEX,
					Symbol:     &types.StringValue{Value: p.Symbol},
					SecurityID: &types.UInt64Value{Value: sec.SecurityID},
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

	request, weight, err := bitmex.GetMargin(msg.Account.Credentials, "XBt")
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		balanceList.RejectionReason = messages.RateLimitExceeded
		context.Respond(balanceList)
		return nil
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			balanceList.RejectionReason = messages.HTTPError
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
				balanceList.RejectionReason = messages.HTTPError
				context.Respond(balanceList)
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"%d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http server error", log.Error(err))
				balanceList.RejectionReason = messages.HTTPError
				context.Respond(balanceList)
			}
			return
		}
		var margin bitmex.Margin
		err = json.Unmarshal(queryResponse.Response, &margin)
		if err != nil {
			balanceList.RejectionReason = messages.HTTPError
			context.Respond(balanceList)
			return
		}
		balanceList.Balances = append(balanceList.Balances, &models.Balance{
			AccountID: msg.Account.AccountID,
			Asset:     &constants.BITCOIN,
			Quantity:  float64(margin.WalletBalance) * 0.00000001,
		})
		balanceList.Success = true
		context.Respond(balanceList)
	})

	return nil
}

func buildPostOrderRequest(order *messages.NewOrder) bitmex.PostOrderRequest {
	request := bitmex.NewPostOrderRequest(order.Instrument.Symbol.Value)

	request.SetOrderQty(order.Quantity)
	request.SetClOrdID(order.ClientOrderID)

	switch order.OrderSide {
	case models.Buy:
		request.SetSide(bitmex.BUY_ORDER_SIDE)
	case models.Sell:
		request.SetSide(bitmex.SELL_ORDER_SIDE)
	default:
		request.SetSide(bitmex.BUY_ORDER_SIDE)
	}

	switch order.OrderType {
	case models.Limit:
		request.SetOrderType(bitmex.LIMIT)
	case models.Market:
		request.SetOrderType(bitmex.MARKET)
	case models.Stop:
		request.SetOrderType(bitmex.STOP)
	case models.StopLimit:
		request.SetOrderType(bitmex.STOP_LIMIT)
	case models.LimitIfTouched:
		request.SetOrderType(bitmex.LIMIT_IF_TOUCHED)
	case models.MarketIfTouched:
		request.SetOrderType(bitmex.MARKET_IF_TOUCHED)
	default:
		request.SetOrderType(bitmex.LIMIT)
	}

	switch order.TimeInForce {
	case models.Session:
		request.SetTimeInForce(bitmex.TIF_DAY)
	case models.GoodTillCancel:
		request.SetTimeInForce(bitmex.GOOD_TILL_CANCEL)
	case models.ImmediateOrCancel:
		request.SetTimeInForce(bitmex.IMMEDIATE_OR_CANCEL)
	case models.FillOrKill:
		request.SetTimeInForce(bitmex.FILL_OR_KILL)
	default:
		request.SetTimeInForce(bitmex.TIF_DAY)
	}

	if order.Price != nil {
		request.SetPrice(order.Price.Value)
	}

	// TODO handle multiple exec inst
	if len(order.ExecutionInstructions) > 0 {
		switch order.ExecutionInstructions[0] {
		case messages.ParticipateDoNotInitiate:
			request.SetExecInst(bitmex.EI_PARTICIPATE_DO_NOT_INITIATE)
		}
	}

	return request
}

func (state *Executor) OnNewOrderSingleRequest(context actor.Context) error {
	req := context.Message().(*messages.NewOrderSingleRequest)
	response := &messages.NewOrderSingleResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	if req.Order.Instrument == nil || req.Order.Instrument.Symbol == nil {
		response.RejectionReason = messages.UnknownSymbol
		context.Respond(response)
		return nil
	}

	params := buildPostOrderRequest(req.Order)

	request, weight, err := bitmex.PostOrder(req.Account.Credentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
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
		var order bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			response.RejectionReason = messages.HTTPError
			context.Respond(response)
			return
		}
		response.Success = true
		response.OrderID = order.OrderID
		context.Respond(response)
	})
	return nil
}

func (state *Executor) OnNewOrderBulkRequest(context actor.Context) error {
	// Launch futures, and await on each of them
	req := context.Message().(*messages.NewOrderBulkRequest)
	response := &messages.NewOrderBulkResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	if len(req.Orders) == 0 {
		response.Success = true
		context.Respond(response)
		return nil
	}

	var params []bitmex.PostOrderRequest
	// Validation
	symbol := ""
	if req.Orders[0].Instrument != nil && req.Orders[0].Instrument.Symbol != nil {
		symbol = req.Orders[0].Instrument.Symbol.Value
	}

	for _, order := range req.Orders {
		if order.Instrument == nil || order.Instrument.Symbol == nil {
			response.RejectionReason = messages.UnknownSymbol
			context.Respond(response)
			return nil
		} else if symbol != order.Instrument.Symbol.Value {
			fmt.Println("DIFFERENT SUYMBOL")
			response.RejectionReason = messages.DifferentSymbols
			fmt.Println("RES", response)
			return nil
		}
		params = append(params, buildPostOrderRequest(order))
	}
	fmt.Println("HELLOOO")
	request, weight, err := bitmex.PostBulkOrder(req.Account.Credentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			response.RejectionReason = messages.ExchangeAPIError
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
		var orders []bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &orders)
		if err != nil {
			state.logger.Info("unmarshalling error", log.Error(err))
			response.Success = false
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		var orderIDs []string
		for _, reqOrder := range req.Orders {
			for _, order := range orders {
				if order.ClOrdID != nil && reqOrder.ClientOrderID == *order.ClOrdID {
					orderIDs = append(orderIDs, order.OrderID)
					break
				}
			}
		}
		response.OrderIDs = orderIDs
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

	request, weight, err := bitmex.AmendOrder(req.Account.Credentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.HTTPError
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
		var order bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
		context.Respond(response)
	})
	return nil
}

func (state *Executor) OnOrderBulkReplaceRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderBulkReplaceRequest)
	response := &messages.OrderBulkReplaceResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	var params []bitmex.AmendOrderRequest
	for _, u := range req.Updates {
		param := bitmex.NewAmendOrderRequest()
		if u.OrderID != nil {
			param.SetOrderID(u.OrderID.Value)
		} else if u.OrigClientOrderID != nil {
			param.SetOrigClOrdID(u.OrigClientOrderID.Value)
		}

		if u.Price != nil {
			param.SetPrice(u.Price.Value)
		}
		if u.Quantity != nil {
			param.SetOrderQty(u.Quantity.Value)
		}
		params = append(params, param)
	}

	request, weight, err := bitmex.AmendBulkOrder(req.Account.Credentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
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
		var orders []bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &orders)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
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

	request, weight, err := bitmex.CancelOrder(req.Account.Credentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
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
		var orders []bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &orders)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
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
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	filters := make(map[string]interface{})
	params := bitmex.NewCancelAllOrdersRequest()
	if req.Filter != nil {
		if req.Filter.Instrument != nil {
			if req.Filter.Instrument.Symbol != nil {
				if _, ok := state.symbolToSec[req.Filter.Instrument.Symbol.Value]; !ok {
					response.RejectionReason = messages.UnknownSymbol
					context.Respond(response)
					return nil
				}
				params.SetSymbol(req.Filter.Instrument.Symbol.Value)
			} else if req.Filter.Instrument.SecurityID != nil {
				sec, ok := state.securities[req.Filter.Instrument.SecurityID.Value]
				if !ok {
					response.RejectionReason = messages.UnknownSymbol
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

	request, weight, err := bitmex.CancelAllOrders(req.Account.Credentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
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
		var orders []bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &orders)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
		context.Respond(response)
	})
	return nil
}
