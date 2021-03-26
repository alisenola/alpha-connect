package fbinance

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/interface"
	"gitlab.com/alphaticks/alpha-connect/jobs"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/fbinance"
	"io/ioutil"
	"math"
	"net/http"
	"reflect"
	"strconv"
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

type Executor struct {
	client               *http.Client
	securities           map[uint64]*models.Security
	symbolToSec          map[string]*models.Security
	minuteOrderRateLimit *exchanges.RateLimit
	globalRateLimit      *exchanges.RateLimit
	queryRunner          *actor.PID
	logger               *log.Logger
}

func NewExecutor() actor.Actor {
	return &Executor{
		client:               nil,
		minuteOrderRateLimit: nil,
		globalRateLimit:      nil,
		queryRunner:          nil,
		logger:               nil,
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

	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewAPIQuery(state.client)
	})
	state.queryRunner = context.Spawn(props)

	request, weight, err := fbinance.GetExchangeInfo()

	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
	res, err := future.Result()
	if err != nil {
		return err
	}
	queryResponse := res.(*jobs.PerformQueryResponse)
	if queryResponse.StatusCode != 200 {
		return fmt.Errorf("error getting exchange info: status code %d", queryResponse.StatusCode)
	}

	var exchangeInfo fbinance.ExchangeInfo
	err = json.Unmarshal(queryResponse.Response, &exchangeInfo)
	if err != nil {
		return fmt.Errorf("error decoding query response: %v", err)
	}
	if exchangeInfo.Code != 0 {
		return fmt.Errorf("error getting exchange info: %s", exchangeInfo.Msg)
	}

	// Initialize rate limit
	for _, rateLimit := range exchangeInfo.RateLimits {
		if rateLimit.RateLimitType == "ORDERS" {
			if rateLimit.Interval == "MINUTE" {
				state.minuteOrderRateLimit = exchanges.NewRateLimit(rateLimit.Limit, time.Minute)
			}
		} else if rateLimit.RateLimitType == "REQUEST_WEIGHT" {
			state.globalRateLimit = exchanges.NewRateLimit(rateLimit.Limit, time.Minute)
		}
	}
	if state.minuteOrderRateLimit == nil || state.globalRateLimit == nil {
		return fmt.Errorf("unable to set minute or global rate limit")
	}

	// Update rate limit with weight from the current exchange info fetch
	state.globalRateLimit.Request(weight)
	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	request, weight, err := fbinance.GetExchangeInfo()
	if err != nil {
		return err
	}

	if state.globalRateLimit.IsRateLimited() {
		return fmt.Errorf("rate limited")
	}

	state.globalRateLimit.Request(weight)
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
	var exchangeInfo fbinance.ExchangeInfo
	err = json.Unmarshal(response, &exchangeInfo)
	if err != nil {
		err = fmt.Errorf(
			"error unmarshaling response: %v",
			err)
		return err
	}
	if exchangeInfo.Code != 0 {
		err = fmt.Errorf(
			"fbinance api error: %d %s",
			exchangeInfo.Code,
			exchangeInfo.Msg)
		return err
	}

	var securities []*models.Security
	for _, symbol := range exchangeInfo.Symbols {
		baseCurrency, ok := constants.GetAssetBySymbol(symbol.BaseAsset)
		if !ok {
			continue
		}
		quoteCurrency, ok := constants.GetAssetBySymbol(symbol.QuoteAsset)
		if !ok {
			continue
		}
		security := models.Security{}
		security.Symbol = symbol.Symbol
		security.Underlying = baseCurrency
		security.QuoteCurrency = quoteCurrency
		switch symbol.Status {
		case "PRE_TRADING":
			security.Status = models.PreTrading
		case "TRADING":
			security.Status = models.Trading
		case "POST_TRADING":
			security.Status = models.PostTrading
		case "END_OF_DAY":
			security.Status = models.EndOfDay
		case "HALT":
			security.Status = models.Halt
		case "AUCTION_MATCH":
			security.Status = models.AuctionMatch
		case "BREAK":
			security.Status = models.Break
		default:
			security.Status = models.Disabled
		}
		security.Exchange = &constants.FBINANCE
		switch symbol.ContractType {
		case "PERPETUAL":
			security.SecurityType = enum.SecurityType_CRYPTO_PERP
		default:
			continue
		}
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name)
		security.MinPriceIncrement = &types.DoubleValue{Value: 1. / math.Pow10(symbol.PricePrecision)}
		security.RoundLot = &types.DoubleValue{Value: 1. / math.Pow10(symbol.QuantityPrecision)}
		security.IsInverse = false
		security.Multiplier = &types.DoubleValue{Value: 1.}
		// Default fee
		security.MakerFee = &types.DoubleValue{Value: 0.0002}
		security.TakerFee = &types.DoubleValue{Value: 0.0004}
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

func (state *Executor) OnMarketDataRequest(context actor.Context) error {
	var snapshot *models.OBL2Snapshot
	msg := context.Message().(*messages.MarketDataRequest)
	response := &messages.MarketDataResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	if msg.Subscribe {
		response.RejectionReason = messages.UnsupportedSubscription
		context.Respond(response)
		return nil
	}
	if msg.Instrument == nil || msg.Instrument.Symbol == nil {
		response.RejectionReason = messages.MissingInstrument
		context.Respond(response)
		return nil
	}
	symbol := msg.Instrument.Symbol.Value
	// Get http request and the expected response
	request, weight, err := fbinance.GetOrderBook(symbol, 1000)
	if err != nil {
		return err
	}

	if state.globalRateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.globalRateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Info("http client error", log.Error(err))
			response.RejectionReason = messages.HTTPError
			context.Respond(response)
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"http client error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
				return
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"http server error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
				return
			}
			return
		}
		var obData fbinance.OrderBookData
		err = json.Unmarshal(queryResponse.Response, &obData)
		if err != nil {
			err = fmt.Errorf("error decoding query response: %v", err)
			state.logger.Info("http client error", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		if obData.Code != 0 {
			err = fmt.Errorf("error getting orderbook: %d %s", obData.Code, obData.Msg)
			state.logger.Info("http client error", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}

		bids, asks, err := obData.ToBidAsk()
		if err != nil {
			err = fmt.Errorf("error converting orderbook: %v", err)
			state.logger.Info("http client error", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		snapshot = &models.OBL2Snapshot{
			Bids:      bids,
			Asks:      asks,
			Timestamp: &types.Timestamp{Seconds: 0, Nanos: 0},
		}
		response.Success = true
		response.SnapshotL2 = snapshot
		response.SeqNum = obData.LastUpdateID
		context.Respond(response)
	})

	return nil
}

func (state *Executor) OnOrderStatusRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderStatusRequest)
	response := &messages.OrderList{
		RequestID: msg.RequestID,
		Success:   false,
	}
	symbol := ""
	orderID := ""
	clOrderID := ""
	if msg.Filter != nil {
		if msg.Filter.OrderStatus != nil {
			response.RejectionReason = messages.UnsupportedFilter
			context.Respond(response)
			return nil
		}
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
	if orderID != "" || clOrderID != "" {
		params := fbinance.NewQueryOrderRequest(symbol)
		if orderID != "" {
			orderIDInt, err := strconv.ParseInt(orderID, 10, 64)
			if err != nil {
				response.RejectionReason = messages.UnsupportedFilter
				context.Respond(response)
				return nil
			}
			params.SetOrderID(orderIDInt)
		}
		if clOrderID != "" {
			params.SetOrigClientOrderID(clOrderID)
		}
		var err error
		request, weight, err = fbinance.QueryOrder(msg.Account.Credentials, params)
		if err != nil {
			return err
		}
	} else {
		var err error
		request, weight, err = fbinance.QueryOpenOrders(msg.Account.Credentials, symbol)
		if err != nil {
			return err
		}
	}

	if state.globalRateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.globalRateLimit.Request(weight)
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
		var orders []fbinance.OrderData
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
				OrderID:       fmt.Sprintf("%d", o.OrderID),
				ClientOrderID: o.ClientOrderID,
				Instrument: &models.Instrument{
					Exchange:   &constants.BITMEX,
					Symbol:     &types.StringValue{Value: o.Symbol},
					SecurityID: &types.UInt64Value{Value: sec.SecurityID},
				},
				LeavesQuantity: o.OriginalQuantity - o.ExecutedQuantity, // TODO check
				CumQuantity:    o.CumQuantity,
			}

			switch o.Status {
			case fbinance.OS_NEW:
				ord.OrderStatus = models.New
			case fbinance.OS_CANCELED:
				ord.OrderStatus = models.Canceled
			default:
				fmt.Println("UNKNWOEN ORDER STATUS", o.Status)
			}

			/*
				const LIMIT = OrderType("LIMIT")
				const MARKET = OrderType("MARKET")
				const STOP_LOSS = OrderType("STOP_LOSS")
				const STOP_MARKET = OrderType("STOP_MARKET")
				const STOP_LOSS_LIMIT = OrderType("STOP_LOSS_LIMIT")
				const TAKE_PROFIT = OrderType("TAKE_PROFIT")
				const TAKE_PROFIT_LIMIT = OrderType("TAKE_PROFIT_LIMIT")
				const TAKE_PROFIT_MARKET = OrderType("TAKE_PROFIT_MARKET")
				const LIMIT_MAKER = OrderType("LIMIT_MAKER")
				const TRAILING_STOP_MARKET = OrderType("TRAILING_STOP_MARKET")

			*/

			switch o.Type {
			case fbinance.LIMIT:
				ord.OrderType = models.Limit
			case fbinance.MARKET:
				ord.OrderType = models.Market
			case fbinance.STOP_LOSS:
				ord.OrderType = models.Stop
			case fbinance.STOP_LOSS_LIMIT:
				ord.OrderType = models.StopLimit
			default:
				fmt.Println("UNKNOWN ORDER TYPE", o.Type)
			}

			switch o.Side {
			case fbinance.BUY_SIDE:
				ord.Side = models.Buy
			case fbinance.SELL_SIDE:
				ord.Side = models.Sell
			default:
				fmt.Println("UNKNOWN ORDER SIDE", o.Side)
			}

			switch o.TimeInForce {
			case fbinance.GOOD_TILL_CANCEL:
				ord.TimeInForce = models.GoodTillCancel
			case fbinance.FILL_OR_KILL:
				ord.TimeInForce = models.FillOrKill
			case fbinance.IMMEDIATE_OR_CANCEL:
				ord.TimeInForce = models.ImmediateOrCancel
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

	symbol := ""
	if msg.Instrument != nil {
		if msg.Instrument.Symbol != nil {
			symbol = msg.Instrument.Symbol.Value
		} else if msg.Instrument.SecurityID != nil {
			sec, ok := state.securities[msg.Instrument.SecurityID.Value]
			if !ok {
				positionList.RejectionReason = messages.UnknownSecurityID
				context.Respond(positionList)
				return nil
			}
			symbol = sec.Symbol
		}
	}

	request, weight, err := fbinance.GetPositionRisk(msg.Account.Credentials)
	if err != nil {
		return err
	}

	if state.globalRateLimit.IsRateLimited() {
		positionList.RejectionReason = messages.RateLimitExceeded
		context.Respond(positionList)
		return nil
	}
	state.globalRateLimit.Request(weight)

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
		var positions []fbinance.AccountPositionRisk
		err = json.Unmarshal(queryResponse.Response, &positions)
		if err != nil {
			state.logger.Info("unmarshaling error", log.Error(err))
			positionList.RejectionReason = messages.ExchangeAPIError
			context.Respond(positionList)
			return
		}
		for _, p := range positions {
			if symbol != "" && p.Symbol != symbol {
				continue
			}
			sec, ok := state.symbolToSec[p.Symbol]
			if !ok {
				err := fmt.Errorf("unknown symbol %s", p.Symbol)
				state.logger.Info("unmarshaling error", log.Error(err))
				positionList.RejectionReason = messages.ExchangeAPIError
				context.Respond(positionList)
				return
			}
			var cost float64
			if sec.IsInverse {
				cost = ((1. / p.MarkPrice) * sec.Multiplier.Value * p.PositionAmount) - p.UnrealizedProfit
			} else {
				cost = (p.MarkPrice * sec.Multiplier.Value * p.PositionAmount) - p.UnrealizedProfit
			}
			pos := &models.Position{
				AccountID: msg.Account.AccountID,
				Instrument: &models.Instrument{
					Exchange:   &constants.BITMEX,
					Symbol:     &types.StringValue{Value: p.Symbol},
					SecurityID: &types.UInt64Value{Value: sec.SecurityID},
				},
				Quantity: p.PositionAmount,
				Cost:     cost,
				Cross:    false,
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

	request, weight, err := fbinance.GetBalance(msg.Account.Credentials)
	if err != nil {
		return err
	}

	if state.globalRateLimit.IsRateLimited() {
		balanceList.RejectionReason = messages.RateLimitExceeded
		context.Respond(balanceList)
		return nil
	}

	state.globalRateLimit.Request(weight)
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
		var balances []fbinance.AccountBalance
		err = json.Unmarshal(queryResponse.Response, &balances)
		if err != nil {
			balanceList.RejectionReason = messages.HTTPError
			context.Respond(balanceList)
			return
		}
		for _, b := range balances {
			asset, ok := constants.GetAssetBySymbol(b.Asset)
			if !ok {
				state.logger.Error("got balance for unknown asset", log.String("asset", b.Asset))
			}
			balanceList.Balances = append(balanceList.Balances, &models.Balance{
				AccountID: msg.Account.AccountID,
				Asset:     asset,
				Quantity:  b.Balance,
			})
		}

		balanceList.Success = true
		context.Respond(balanceList)
	})

	return nil
}

func buildPostOrderRequest(order *messages.NewOrder) (fbinance.NewOrderRequest, *messages.RejectionReason) {
	var side fbinance.OrderSide
	var typ fbinance.OrderType
	if order.OrderSide == models.Buy {
		side = fbinance.BUY_SIDE
	} else {
		side = fbinance.SELL_SIDE
	}
	switch order.OrderType {
	case models.Limit:
		typ = fbinance.LIMIT
	case models.Market:
		typ = fbinance.MARKET
	case models.Stop:
		typ = fbinance.STOP_LOSS
	case models.StopLimit:
		typ = fbinance.STOP_LOSS_LIMIT
	default:
		rej := messages.UnsupportedOrderType
		return nil, &rej
	}

	request := fbinance.NewNewOrderRequest(order.Instrument.Symbol.Value, side, typ)

	request.SetQuantity(order.Quantity)
	request.SetNewClientOrderID(order.ClientOrderID)

	if order.OrderType != models.Market {
		switch order.TimeInForce {
		case models.Session:
			request.SetTimeInForce(fbinance.GOOD_TILL_CANCEL)
		case models.GoodTillCancel:
			request.SetTimeInForce(fbinance.GOOD_TILL_CANCEL)
		case models.ImmediateOrCancel:
			request.SetTimeInForce(fbinance.IMMEDIATE_OR_CANCEL)
		case models.FillOrKill:
			request.SetTimeInForce(fbinance.FILL_OR_KILL)
		default:
			rej := messages.UnsupportedOrderTimeInForce
			return nil, &rej
		}
	}

	if order.Price != nil {
		request.SetPrice(order.Price.Value)
	}

	// TODO handle multiple exec inst
	if len(order.ExecutionInstructions) > 0 {
		rej := messages.UnsupportedOrderCharacteristic
		return nil, &rej
	}

	return request, nil
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

	params, rej := buildPostOrderRequest(req.Order)
	if rej != nil {
		response.RejectionReason = *rej
		context.Respond(response)
		return nil
	}

	request, weight, err := fbinance.NewOrder(req.Account.Credentials, params)
	if err != nil {
		return err
	}

	if state.globalRateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.globalRateLimit.Request(weight)
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
	params := fbinance.NewQueryOrderRequest(symbol)
	if req.OrderID != nil {
		orderIDInt, err := strconv.ParseInt(req.OrderID.Value, 10, 64)
		if err != nil {
			response.RejectionReason = messages.UnknownOrder
			context.Respond(response)
			return nil
		}
		params.SetOrderID(orderIDInt)
	} else if req.ClientOrderID != nil {
		params.SetOrigClientOrderID(req.ClientOrderID.Value)
	} else {
		response.RejectionReason = messages.UnknownOrder
		context.Respond(response)
		return nil
	}

	request, weight, err := fbinance.CancelOrder(req.Account.Credentials, params)
	if err != nil {
		return err
	}

	if state.globalRateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.globalRateLimit.Request(weight)
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
		var fres fbinance.Response
		err = json.Unmarshal(queryResponse.Response, &fres)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		if fres.Code != 200 {
			state.logger.Info("error unmarshalling", log.Error(errors.New(fres.Message)))
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
	symbol := ""
	if req.Filter != nil {
		if req.Filter.Instrument != nil {
			if req.Filter.Instrument.Symbol != nil {
				if _, ok := state.symbolToSec[req.Filter.Instrument.Symbol.Value]; !ok {
					response.RejectionReason = messages.UnknownSymbol
					context.Respond(response)
					return nil
				}
				symbol = req.Filter.Instrument.Symbol.Value
			} else if req.Filter.Instrument.SecurityID != nil {
				sec, ok := state.securities[req.Filter.Instrument.SecurityID.Value]
				if !ok {
					response.RejectionReason = messages.UnknownSymbol
					context.Respond(response)
					return nil
				}
				symbol = sec.Symbol
			}
		}
		if req.Filter.Side != nil || req.Filter.OrderStatus != nil {
			response.RejectionReason = messages.UnsupportedFilter
			context.Respond(response)
			return nil
		}
	}
	if symbol == "" {
		response.RejectionReason = messages.UnknownSymbol
		context.Respond(response)
		return nil
	}

	request, weight, err := fbinance.CancelAllOrders(req.Account.Credentials, symbol)
	if err != nil {
		return err
	}

	if state.globalRateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}

	state.globalRateLimit.Request(weight)
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
		var fres fbinance.Response
		err = json.Unmarshal(queryResponse.Response, &fres)
		if err != nil {
			state.logger.Info("error unmarshalling", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		if fres.Code != 200 {
			state.logger.Info("error unmarshalling", log.Error(errors.New(fres.Message)))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}
		response.Success = true
		context.Respond(response)
	})
	return nil
}
