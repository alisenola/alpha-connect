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

	orderList := &messages.OrderList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Error:      "",
		Orders:     nil,
	}

	params := bitmex.NewGetOrderParams()
	if msg.Instrument != nil {
		if msg.Instrument.Symbol != nil {
			params.SetSymbol(msg.Instrument.Symbol.Value)
		} else if msg.Instrument.SecurityID != nil {
			sec, ok := state.securities[msg.Instrument.SecurityID.Value]
			if !ok {
				orderList.Error = "unknown security"
				context.Respond(orderList)
				return nil
			}
			params.SetSymbol(sec.Symbol)
		}
	}
	filters := make(map[string]interface{})
	if msg.OrderID != nil {
		filters["orderID"] = msg.OrderID.Value
	}
	if msg.ClientOrderID != nil {
		filters["clOrdID"] = msg.ClientOrderID.Value
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
			orderList.Error = err.Error()
			context.Respond(orderList)
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"http client error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				orderList.Error = err.Error()
				context.Respond(orderList)
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"http server error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				orderList.Error = err.Error()
				context.Respond(orderList)
			}
			return
		}
		var orders []bitmex.Order
		err = json.Unmarshal(queryResponse.Response, &orders)
		if err != nil {
			orderList.Error = err.Error()
			context.Respond(orderList)
			return
		}

		for _, o := range orders {
			sec, ok := state.symbolToSec[o.Symbol]
			if !ok {
				orderList.Error = "got order for unknown security"
				context.Respond(orderList)
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

			orderList.Orders = append(orderList.Orders, &ord)
		}
		context.Respond(orderList)
	})

	return nil
}

func (state *Executor) OnPositionsRequest(context actor.Context) error {
	msg := context.Message().(*messages.PositionsRequest)

	positionList := &messages.PositionList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Error:      "",
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
				positionList.Error = "unknown security"
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
			positionList.Error = err.Error()
			context.Respond(positionList)
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"http client error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				positionList.Error = err.Error()
				context.Respond(positionList)
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"http server error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				positionList.Error = err.Error()
				context.Respond(positionList)
			}
			return
		}

		var positions []bitmex.Position
		err = json.Unmarshal(queryResponse.Response, &positions)
		if err != nil {
			positionList.Error = err.Error()
			context.Respond(positionList)
			return
		}
		for _, p := range positions {
			sec, ok := state.symbolToSec[p.Symbol]
			if !ok {
				positionList.Error = "got order for unknown security"
				context.Respond(positionList)
				return
			}
			pos := &models.Position{
				AccountID: msg.Account.AccountID,
				Instrument: &models.Instrument{
					Exchange:   &constants.BITMEX,
					Symbol:     &types.StringValue{Value: p.Symbol},
					SecurityID: &types.UInt64Value{Value: sec.SecurityID},
				},
				Quantity: float64(p.CurrentQty),
			}
			positionList.Positions = append(positionList.Positions, pos)
		}
		context.Respond(positionList)
	})

	return nil
}

func (state *Executor) OnNewOrderSingleRequest(context actor.Context) error {
	msg := context.Message().(*messages.NewOrderSingleRequest)
	if msg.Instrument == nil || msg.Instrument.Symbol == nil {
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.UnknownSymbol,
		})
		return nil
	}
	params := bitmex.NewPostOrderRequest(msg.Instrument.Symbol.Value)

	params.SetOrderQty(msg.Quantity)
	params.SetClOrdID(msg.ClientOrderID)

	switch msg.OrderSide {
	case models.Buy:
		params.SetSide(bitmex.BUY_ORDER_SIDE)
	case models.Sell:
		params.SetSide(bitmex.SELL_ORDER_SIDE)
	default:
		params.SetSide(bitmex.BUY_ORDER_SIDE)
	}

	switch msg.OrderType {
	case models.Limit:
		params.SetOrderType(bitmex.LIMIT)
	case models.Market:
		params.SetOrderType(bitmex.MARKET)
	case models.Stop:
		params.SetOrderType(bitmex.STOP)
	case models.StopLimit:
		params.SetOrderType(bitmex.STOP_LIMIT)
	case models.LimitIfTouched:
		params.SetOrderType(bitmex.LIMIT_IF_TOUCHED)
	case models.MarketIfTouched:
		params.SetOrderType(bitmex.MARKET_IF_TOUCHED)
	default:
		params.SetOrderType(bitmex.LIMIT)
	}

	switch msg.TimeInForce {
	case models.Session:
		params.SetTimeInForce(bitmex.TIF_DAY)
	case models.GoodTillCancel:
		params.SetTimeInForce(bitmex.GOOD_TILL_CANCEL)
	case models.ImmediateOrCancel:
		params.SetTimeInForce(bitmex.IMMEDIATE_OR_CANCEL)
	case models.FillOrKill:
		params.SetTimeInForce(bitmex.FILL_OR_KILL)
	default:
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.UnsupportedOrderCharacteristic,
		})
		return nil
	}

	if msg.Price != nil {
		params.SetPrice(msg.Price.Value)
	}

	request, weight, err := bitmex.PostOrder(msg.Account.Credentials, params)
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
				RequestID:       msg.RequestID,
				Success:         false,
				RejectionReason: messages.Other,
			})
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				context.Respond(&messages.NewOrderSingleResponse{
					RequestID:       msg.RequestID,
					Success:         false,
					RejectionReason: messages.Other,
				})
			} else if queryResponse.StatusCode >= 500 {
				context.Respond(&messages.NewOrderSingleResponse{
					RequestID:       msg.RequestID,
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
				RequestID:       msg.RequestID,
				Success:         false,
				RejectionReason: messages.Other,
			})
			return
		}
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID: msg.RequestID,
			Success:   true,
			OrderID:   order.OrderID,
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
