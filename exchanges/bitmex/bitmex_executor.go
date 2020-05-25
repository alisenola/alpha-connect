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

	// TODO rate limitting
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
		time.Sleep(state.rateLimit.DurationBeforeNextRequest(weight))
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
		Error:      "",
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
		Error:      "",
		Securities: securities})

	return nil
}

func (state *Executor) OnMarketDataRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnOrderStatusRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderStatusRequest)

	filters := make(map[string]interface{})
	params := bitmex.NewGetOrderParams()
	if msg.Filter != nil {
		if msg.Filter.Instrument != nil {
			if msg.Filter.Instrument.Symbol != nil {
				params.SetSymbol(msg.Filter.Instrument.Symbol.Value)
			} else if msg.Filter.Instrument.SecurityID != nil {
				sec, ok := state.securities[msg.Filter.Instrument.SecurityID.Value]
				if !ok {
					context.Respond(&messages.OrderList{
						RequestID:       msg.RequestID,
						Success:         false,
						RejectionReason: messages.UnknownSecurityID,
					})
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

	if len(filters) > 0 {
		params.SetFilters(filters)
	}

	request, weight, err := bitmex.GetOrder(msg.Account.Credentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		time.Sleep(state.rateLimit.DurationBeforeNextRequest(weight))
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			context.Respond(&messages.OrderList{
				RequestID:       msg.RequestID,
				Success:         false,
				RejectionReason: messages.Other,
			})
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				/*
					err := fmt.Errorf(
						"http client error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
				*/
				context.Respond(&messages.OrderList{
					RequestID:       msg.RequestID,
					Success:         false,
					RejectionReason: messages.ExchangeAPIError,
				})
			} else if queryResponse.StatusCode >= 500 {
				/*
					err := fmt.Errorf(
						"http server error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))

				*/
				context.Respond(&messages.OrderList{
					RequestID:       msg.RequestID,
					Success:         false,
					RejectionReason: messages.ExchangeAPIError,
				})
			}
			return
		}
		var orders []bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &orders)
		if err != nil {
			context.Respond(&messages.OrderList{
				RequestID:       msg.RequestID,
				Success:         false,
				RejectionReason: messages.ExchangeAPIError,
			})
			return
		}

		var morders []*models.Order
		for _, o := range orders {
			sec, ok := state.symbolToSec[o.Symbol]
			if !ok {
				context.Respond(&messages.OrderList{
					RequestID:       msg.RequestID,
					Success:         false,
					RejectionReason: messages.ExchangeAPIError,
				})
				return
			}
			ord := models.Order{
				OrderID: o.OrderID,
				Instrument: &models.Instrument{
					Exchange:   &constants.BITMEX,
					Symbol:     &types.StringValue{Value: o.Symbol},
					SecurityID: &types.UInt64Value{Value: sec.SecurityID},
				},
				LeavesQuantity: 0,
				CumQuantity:    0,
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

			morders = append(morders, &ord)
		}
		context.Respond(&messages.OrderList{
			RequestID: msg.RequestID,
			Success:   true,
			Orders:    morders,
		})
	})

	return nil
}

func (state *Executor) OnPositionsRequest(context actor.Context) error {
	msg := context.Message().(*messages.PositionsRequest)

	positionList := &messages.PositionList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
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
				positionList.Success = false
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
		time.Sleep(state.rateLimit.DurationBeforeNextRequest(weight))
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			positionList.Success = false
			positionList.RejectionReason = messages.Other
			context.Respond(positionList)
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				/*
					err := fmt.Errorf(
						"http client error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
				*/
				positionList.Success = false
				positionList.RejectionReason = messages.Other
				context.Respond(positionList)
			} else if queryResponse.StatusCode >= 500 {
				/*
					err := fmt.Errorf(
						"http server error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
				*/
				positionList.Success = false
				positionList.RejectionReason = messages.Other
				context.Respond(positionList)
			}
			return
		}
		var positions []bitmex.Position
		err = json.Unmarshal(queryResponse.Response, &positions)
		if err != nil {
			positionList.Success = false
			positionList.RejectionReason = messages.Other
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
				positionList.Success = false
				positionList.RejectionReason = messages.Other
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
		context.Respond(positionList)
	})

	return nil
}

func (state *Executor) OnBalancesRequest(context actor.Context) error {
	msg := context.Message().(*messages.BalancesRequest)

	balanceList := &messages.BalanceList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Balances:   nil,
	}

	request, weight, err := bitmex.GetMargin(msg.Account.Credentials, "XBt")
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		time.Sleep(state.rateLimit.DurationBeforeNextRequest(weight))
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			balanceList.Success = false
			balanceList.RejectionReason = messages.Other
			context.Respond(balanceList)
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				/*
					err := fmt.Errorf(
						"http client error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
				*/
				balanceList.Success = false
				balanceList.RejectionReason = messages.Other
				context.Respond(balanceList)
			} else if queryResponse.StatusCode >= 500 {
				/*
					err := fmt.Errorf(
						"http server error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
				*/
				balanceList.Success = false
				balanceList.RejectionReason = messages.Other
				context.Respond(balanceList)
			}
			return
		}
		var margin bitmex.Margin
		err = json.Unmarshal(queryResponse.Response, &margin)
		if err != nil {
			balanceList.Success = false
			balanceList.RejectionReason = messages.Other
			context.Respond(balanceList)
			return
		}
		balanceList.Balances = append(balanceList.Balances, &models.Balance{
			AccountID: msg.Account.AccountID,
			Asset:     &constants.BITCOIN,
			Quantity:  float64(margin.WalletBalance) * 0.00000001,
		})

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

	return request
}

func (state *Executor) OnNewOrderSingleRequest(context actor.Context) error {
	req := context.Message().(*messages.NewOrderSingleRequest)
	if req.Order.Instrument == nil || req.Order.Instrument.Symbol == nil {
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID:       req.RequestID,
			Success:         false,
			RejectionReason: messages.UnknownSymbol,
		})
		return nil
	}

	params := buildPostOrderRequest(req.Order)

	request, weight, err := bitmex.PostOrder(req.Account.Credentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		time.Sleep(state.rateLimit.DurationBeforeNextRequest(weight))
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			context.Respond(&messages.NewOrderSingleResponse{
				RequestID:       req.RequestID,
				Success:         false,
				RejectionReason: messages.Other,
			})
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				context.Respond(&messages.NewOrderSingleResponse{
					RequestID:       req.RequestID,
					Success:         false,
					RejectionReason: messages.Other,
				})
			} else if queryResponse.StatusCode >= 500 {
				context.Respond(&messages.NewOrderSingleResponse{
					RequestID:       req.RequestID,
					Success:         false,
					RejectionReason: messages.Other,
				})
			}
			return
		}
		var order bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &order)
		if err != nil {
			context.Respond(&messages.NewOrderSingleResponse{
				RequestID:       req.RequestID,
				Success:         false,
				RejectionReason: messages.Other,
			})
			return
		}
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID: req.RequestID,
			Success:   true,
			OrderID:   order.OrderID,
		})
	})
	return nil
}

func (state *Executor) OnNewOrderBulkRequest(context actor.Context) error {
	// Launch futures, and await on each of them
	req := context.Message().(*messages.NewOrderBulkRequest)
	if len(req.Orders) == 0 {
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
			context.Respond(&messages.NewOrderBulkResponse{
				RequestID:       req.RequestID,
				Success:         false,
				RejectionReason: messages.UnknownSymbol,
			})
			return nil
		} else if symbol != order.Instrument.Symbol.Value {
			context.Respond(&messages.NewOrderBulkResponse{
				RequestID:       req.RequestID,
				Success:         false,
				RejectionReason: messages.DifferentSymbols,
			})
			return nil
		}
		params = append(params, buildPostOrderRequest(order))
	}

	request, weight, err := bitmex.PostBulkOrder(req.Account.Credentials, params)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		time.Sleep(state.rateLimit.DurationBeforeNextRequest(weight))
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			context.Respond(&messages.NewOrderBulkResponse{
				RequestID:       req.RequestID,
				Success:         false,
				RejectionReason: messages.Other,
			})
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				context.Respond(&messages.NewOrderBulkResponse{
					RequestID:       req.RequestID,
					Success:         false,
					RejectionReason: messages.Other,
				})
			} else if queryResponse.StatusCode >= 500 {
				context.Respond(&messages.NewOrderBulkResponse{
					RequestID:       req.RequestID,
					Success:         false,
					RejectionReason: messages.Other,
				})
			}
			return
		}
		var orders []bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &orders)
		if err != nil {
			context.Respond(&messages.NewOrderBulkResponse{
				RequestID:       req.RequestID,
				Success:         false,
				RejectionReason: messages.Other,
			})
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
		context.Respond(&messages.NewOrderBulkResponse{
			RequestID: req.RequestID,
			Success:   true,
			OrderIDs:  orderIDs,
		})
	})
	return nil
}

func (state *Executor) OnOrderCancelRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderCancelRequest)

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
		time.Sleep(state.rateLimit.DurationBeforeNextRequest(weight))
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			context.Respond(&messages.OrderCancelResponse{
				RequestID:       req.RequestID,
				Success:         false,
				RejectionReason: messages.Other,
			})
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				context.Respond(&messages.OrderCancelResponse{
					RequestID:       req.RequestID,
					Success:         false,
					RejectionReason: messages.Other,
				})
			} else if queryResponse.StatusCode >= 500 {
				context.Respond(&messages.OrderCancelResponse{
					RequestID:       req.RequestID,
					Success:         false,
					RejectionReason: messages.Other,
				})
			}
			return
		}
		var orders []bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &orders)
		if err != nil {
			context.Respond(&messages.OrderCancelResponse{
				RequestID:       req.RequestID,
				Success:         false,
				RejectionReason: messages.Other,
			})
			return
		}
		context.Respond(&messages.OrderCancelResponse{
			RequestID:       req.RequestID,
			Success:         true,
			RejectionReason: messages.Other,
		})
	})
	return nil
}

func (state *Executor) OnOrderMassCancelRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderMassCancelRequest)

	filters := make(map[string]interface{})
	params := bitmex.NewCancelAllOrdersRequest()
	if req.Filter != nil {
		if req.Filter.Instrument != nil {
			if req.Filter.Instrument.Symbol != nil {
				if _, ok := state.symbolToSec[req.Filter.Instrument.Symbol.Value]; !ok {
					context.Respond(&messages.OrderMassCancelResponse{
						RequestID:       req.RequestID,
						Success:         false,
						RejectionReason: messages.UnknownSymbol,
					})
					return nil
				}
				params.SetSymbol(req.Filter.Instrument.Symbol.Value)
			} else if req.Filter.Instrument.SecurityID != nil {
				sec, ok := state.securities[req.Filter.Instrument.SecurityID.Value]
				if !ok {
					context.Respond(&messages.OrderMassCancelResponse{
						RequestID:       req.RequestID,
						Success:         false,
						RejectionReason: messages.UnknownSymbol,
					})
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
		time.Sleep(state.rateLimit.DurationBeforeNextRequest(weight))
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			context.Respond(&messages.OrderMassCancelResponse{
				RequestID:       req.RequestID,
				Success:         false,
				RejectionReason: messages.Other,
			})
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				context.Respond(&messages.OrderMassCancelResponse{
					RequestID:       req.RequestID,
					Success:         false,
					RejectionReason: messages.Other,
				})
			} else if queryResponse.StatusCode >= 500 {
				context.Respond(&messages.OrderMassCancelResponse{
					RequestID:       req.RequestID,
					Success:         false,
					RejectionReason: messages.Other,
				})
			}
			return
		}
		var orders []bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &orders)
		if err != nil {
			context.Respond(&messages.OrderMassCancelResponse{
				RequestID:       req.RequestID,
				Success:         false,
				RejectionReason: messages.Other,
			})
			return
		}
		context.Respond(&messages.OrderMassCancelResponse{
			RequestID:       req.RequestID,
			Success:         true,
			RejectionReason: messages.Other,
		})
	})
	return nil
}
