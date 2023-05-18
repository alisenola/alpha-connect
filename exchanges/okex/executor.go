package okex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

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
	"gitlab.com/alphaticks/xchanger/exchanges/okex"
	xutils "gitlab.com/alphaticks/xchanger/utils"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Executor struct {
	extypes.BaseExecutor
	client      *http.Client
	ws          map[string]*okex.Websocket
	securities  []*models.Security
	rateLimit   *exchanges.RateLimit
	queryRunner *actor.PID
	logger      *log.Logger
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

func (state *Executor) GetLogger() *log.Logger {
	return state.logger
}

func (state *Executor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(state).String()))

	state.client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1024,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	dialers := state.DialerPool.GetDialers()
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
		state.queryRunner = context.Spawn(props)
		state.rateLimit = exchanges.NewRateLimit(6, time.Second)
	}

	state.ws = make(map[string]*okex.Websocket)

	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	if state.ws != nil {
		for _, ws := range state.ws {
			if err := ws.Disconnect(); err != nil {
				state.logger.Info("error disconnecting socket", log.Error(err))
			}
		}
	}
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	if state.rateLimit.IsRateLimited() {
		return fmt.Errorf("rate limited")
	}

	var securities []*models.Security

	request, weight, err := okex.GetInstruments(okex.SPOT)
	if err != nil {
		return err
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

	var kResponse okex.SpotInstrumentsResponse
	err = json.Unmarshal(response, &kResponse)
	if err != nil {
		err = fmt.Errorf("error decoding query response: %v", err)
		return err
	}
	if kResponse.Code != "0" {
		return fmt.Errorf("error decoding query response: %s", kResponse.Msg)
	}

	for _, pair := range kResponse.Data {
		baseCurrency, ok := constants.GetAssetBySymbol(pair.BaseCurrency)
		if !ok {
			//state.logger.Info("unknown symbol " + pair.BaseCurrency + " for instrument " + pair.InstrumentID)
			continue
		}
		quoteCurrency, ok := constants.GetAssetBySymbol(pair.QuoteCurrency)
		if !ok {
			//state.logger.Info("unknown symbol " + pair.QuoteCurrency + " for instrument " + pair.InstrumentID)
			continue
		}

		security := models.Security{}
		security.Symbol = pair.InstrumentID
		security.Underlying = baseCurrency
		security.QuoteCurrency = quoteCurrency
		security.Status = models.InstrumentStatus_Trading
		security.Exchange = constants.OKEX
		security.SecurityType = enum.SecurityType_CRYPTO_SPOT
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinPriceIncrement = &wrapperspb.DoubleValue{Value: pair.TickSize}
		security.RoundLot = &wrapperspb.DoubleValue{Value: pair.LotSize}
		security.MinLimitQuantity = &wrapperspb.DoubleValue{Value: pair.MinOrderSize}
		securities = append(securities, &security)
	}
	state.SyncSecurities(securities, nil)

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: securities})

	return nil
}

func (state *Executor) OnHistoricalLiquidationsRequest(context actor.Context) error {
	msg := context.Message().(*messages.HistoricalLiquidationsRequest)
	response := &messages.HistoricalLiquidationsResponse{
		RequestID: msg.RequestID,
		Success:   false,
	}

	security, rej := state.InstrumentToSecurity(msg.Instrument)

	if rej != nil {
		response.RejectionReason = *rej
		context.Respond(response)
		return nil
	}
	req := okex.NewLiquidationsRequest(okex.SWAP)
	if msg.From != nil {
		frm := uint64(msg.From.Seconds*1000) + uint64(msg.From.Nanos/1000000)
		req.SetBefore(frm)
	}
	if msg.To != nil {
		req.SetAfter(uint64(msg.To.Seconds*1000) + uint64(msg.From.Nanos/1000000))
	}
	if err := req.SetState("filled"); err != nil {
		return fmt.Errorf("error setting liquidation state: %v", err)
	}
	if err := req.SetUnderlying(strings.Replace(security.Symbol, "-SWAP", "", -1)); err != nil {
		return fmt.Errorf("error setting underlying: %v", err)
	}
	request, weight, err := okex.GetLiquidations(req)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		response.Success = false
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

		var resx okex.LiquidationsResponse
		err = json.Unmarshal(queryResponse.Response, &resx)
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		if resx.Code != "0" {
			state.logger.Info("api error", log.Error(errors.New(resx.Msg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		var liquidations []*models.Liquidation
		for _, l := range resx.Data {
			for _, d := range l.Details {
				// Inverse will be in USD, Linear in BTC
				liquidations = append(liquidations, &models.Liquidation{
					Bid:       d.Side == "buy",
					Timestamp: utils.MilliToTimestamp(d.Ts),
					OrderID:   rand.Uint64(),
					Price:     d.BankruptcyPrice,
					Quantity:  d.Size,
				})
			}
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

func (state *Executor) OnBalancesRequest(context actor.Context) error {
	msg := context.Message().(*messages.BalancesRequest)
	sender := context.Sender()
	response := &messages.BalanceList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	if state.ws[msg.Account.ApiCredentials.APIKey] == nil {
		state.ws[msg.Account.ApiCredentials.APIKey] = okex.NewWebsocket()
	}

	if !state.ws[msg.Account.ApiCredentials.APIKey].Connected {
		if err := state.ws[msg.Account.ApiCredentials.APIKey].ConnectPrivate(nil); err != nil {
			state.logger.Warn("error fetching balances", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return nil
		}

		if err := state.ws[msg.Account.ApiCredentials.APIKey].Login(msg.Account.ApiCredentials, msg.Account.ApiCredentials.AccountID); err != nil {
			state.logger.Warn("error fetching balances", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return nil
		}
		if !state.ws[msg.Account.ApiCredentials.APIKey].ReadMessage() {
			state.logger.Warn("error fetching balances", log.Error(state.ws[msg.Account.ApiCredentials.APIKey].Err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return nil
		}
	}

	go func() {
		ws := state.ws[msg.Account.ApiCredentials.APIKey]
		// TODO Dialer

		if state.rateLimit.IsRateLimited() {
			response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
			context.Send(sender, response)
			return
		}

		state.rateLimit.Request(1)

		var balanceArgs []map[string]string
		balanceArg := okex.NewBalanceAndPositionRequest("balance_and_position")
		balanceArgs = append(balanceArgs, balanceArg)
		if err := ws.PrivateSubscribe(balanceArgs); err != nil {
			state.logger.Warn("error fetching balances", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		var balanceAndPositions []okex.WSBalanceAndPosition

		ready := false
		for !ready {
			if !ws.ReadMessage() {
				state.logger.Warn("error fetching balances", log.Error(ws.Err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
			switch msg := ws.Msg.Message.(type) {
			case []okex.WSBalanceAndPosition:
				balanceAndPositions = msg
			}
			ready = balanceAndPositions != nil
		}

		var unSubscribeArgs []map[string]string
		unSubscribeArg := okex.NewUnsubscribeRequest("balance_and_position")
		unSubscribeArgs = append(unSubscribeArgs, unSubscribeArg)
		if err := state.ws[msg.Account.ApiCredentials.APIKey].Unsubscribe(unSubscribeArgs); err != nil {
			state.logger.Warn("error fetching balances", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if !ws.ReadMessage() {
			state.logger.Warn("error fetching balances", log.Error(ws.Err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}

		for _, b := range balanceAndPositions[0].BalData {
			if b.CashBal == "0" {
				continue
			}
			asset := SymbolToAsset(b.Ccy)
			if asset == nil {
				state.logger.Error("got balance for unknown asset", log.String("asset", b.Ccy))
				continue
			}
			Qty, err := strconv.ParseFloat(b.CashBal, 64)
			if err != nil {
				fmt.Println("Error:", err)
			}
			response.Balances = append(response.Balances, &models.Balance{
				Account:  msg.Account.Name,
				Asset:    asset,
				Quantity: Qty,
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
	}

	if state.ws[msg.Account.ApiCredentials.APIKey] == nil {
		state.ws[msg.Account.ApiCredentials.APIKey] = okex.NewWebsocket()
	}

	if !state.ws[msg.Account.ApiCredentials.APIKey].Connected {
		if err := state.ws[msg.Account.ApiCredentials.APIKey].ConnectPrivate(nil); err != nil {
			state.logger.Warn("error fetching balances", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return nil
		}

		if err := state.ws[msg.Account.ApiCredentials.APIKey].Login(msg.Account.ApiCredentials, msg.Account.ApiCredentials.AccountID); err != nil {
			state.logger.Warn("error fetching balances", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return nil
		}
		if !state.ws[msg.Account.ApiCredentials.APIKey].ReadMessage() {
			state.logger.Warn("error fetching balances", log.Error(state.ws[msg.Account.ApiCredentials.APIKey].Err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return nil
		}
	}

	go func() {
		symbol := ""
		if msg.Instrument != nil {
			msg.Instrument.Symbol.Value = strings.ToUpper(msg.Instrument.Symbol.Value)
			s, rej := state.InstrumentToSymbol(msg.Instrument)
			if rej != nil {
				response.RejectionReason = *rej
				context.Send(sender, response)
				return
			}
			symbol = s
		}
		ws := state.ws[msg.Account.ApiCredentials.APIKey]

		if state.rateLimit.IsRateLimited() {
			response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
			context.Send(sender, response)
			return
		}

		state.rateLimit.Request(1)

		var positionArgs []map[string]string
		positionArg := okex.NewPositionRequest("positions", "ANY")
		positionArgs = append(positionArgs, positionArg)
		if err := ws.PrivateSubscribe(positionArgs); err != nil {
			state.logger.Warn("error fetching positions", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}

		var positions []okex.WSPosition
		ready := false
		for !ready {
			if !ws.ReadMessage() {
				state.logger.Warn("error fetching positions", log.Error(ws.Err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
			switch msg := ws.Msg.Message.(type) {
			case []okex.WSPosition:
				positions = msg
			}
			ready = positions != nil
		}

		var unSubscribeArgs []map[string]string
		unSubscribeArg := okex.NewUnsubscribeRequest("positions")
		unSubscribeArgs = append(unSubscribeArgs, unSubscribeArg)
		if err := state.ws[msg.Account.ApiCredentials.APIKey].Unsubscribe(unSubscribeArgs); err != nil {
			state.logger.Warn("error fetching positions", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if !ws.ReadMessage() {
			state.logger.Warn("error fetching positions", log.Error(ws.Err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}

		for _, p := range positions {
			if p.AvgPx == 0 {
				continue
			}
			if symbol != "" && p.InstId != symbol {
				continue
			}
			sec := state.SymbolToSecurity(p.InstId)
			if sec == nil {
				state.logger.Warn(fmt.Sprintf("unknown symbol %s", p.InstId))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Send(sender, response)
				return
			}
			var cost float64
			if sec.IsInverse {
				cost = ((1. / p.MarkPx) * sec.Multiplier.Value * p.AvgPx) - p.Upl
			} else {
				cost = (p.MarkPx * sec.Multiplier.Value * p.AvgPx) - p.Upl
			}
			pos := &models.Position{
				Account: msg.Account.Name,
				Instrument: &models.Instrument{
					Exchange:   constants.FBINANCE,
					Symbol:     &wrapperspb.StringValue{Value: p.InstId},
					SecurityID: &wrapperspb.UInt64Value{Value: sec.SecurityID},
				},
				Quantity:         p.AvgPx,
				Cost:             cost,
				Cross:            false,
				MarkPrice:        wrapperspb.Double(p.MarkPx),
				MaxNotionalValue: wrapperspb.Double(p.NotionalUsd),
			}
			response.Positions = append(response.Positions, pos)
		}
		response.Success = true
		context.Send(sender, response)
	}()

	return nil
}

func (state *Executor) OnHistoricalFundingRatesRequest(context actor.Context) error {
	msg := context.Message().(*messages.HistoricalFundingRatesRequest)
	response := &messages.HistoricalFundingRatesResponse{
		RequestID: msg.RequestID,
		Success:   false,
	}

	req := okex.NewFundingRateHistoryRequest("BTC-USDT-SWAP")
	if msg.From != nil {
		frm := uint64(msg.From.Seconds*1000) + uint64(msg.From.Nanos/1000000)
		req.SetBefore(frm)
	}
	if msg.To != nil {
		req.SetAfter(uint64(msg.To.Seconds*1000) + uint64(msg.From.Nanos/1000000))
	}
	request, weight, err := okex.GetFundingRateHistory(req)
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		response.Success = false
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

		var resx okex.FundingRateHistoryResponse
		err = json.Unmarshal(queryResponse.Response, &resx)
		if err != nil {
			state.logger.Info("http error", log.Error(err))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		if resx.Code != "0" {
			state.logger.Info("api error", log.Error(errors.New(resx.Msg)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Respond(response)
			return
		}

		var rates []*models.Stat
		for _, l := range resx.Data {
			// Inverse will be in USD, Linear in BTC
			rates = append(rates, &models.Stat{
				Timestamp: utils.MilliToTimestamp(l.FundingTime),
				Value:     l.FundingRate,
				StatType:  models.StatType_FundingRate,
			})
		}

		sort.Slice(rates, func(i, j int) bool {
			return utils.TimestampToMilli(rates[i].Timestamp) < utils.TimestampToMilli(rates[j].Timestamp)
		})
		response.Rates = rates
		response.Success = true
		context.Respond(response)
	})

	return nil
}

func (state *Executor) OnOrderStatusRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderStatusRequest)
	sender := context.Sender()
	response := &messages.OrderList{
		RequestID: msg.RequestID,
		Success:   false,
	}

	go func() {
		symbol := ""
		orderID := ""
		clOrderID := ""
		var orderStatus *models.OrderStatus
		if msg.Filter != nil {
			if msg.Filter.OrderStatus != nil {
				orderStatus = &msg.Filter.OrderStatus.Value
			}
			if msg.Filter.Side != nil {
				response.RejectionReason = messages.RejectionReason_UnsupportedFilter
				context.Send(sender, response)
				return
			}
			if msg.Filter.Instrument != nil {
				s, rej := state.InstrumentToSymbol(msg.Filter.Instrument)
				if rej != nil {
					response.RejectionReason = *rej
					context.Send(sender, response)
					return
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
		var data []okex.OrderData
		if orderID != "" || clOrderID != "" {
			params := okex.GetOrderDetailsRequest(symbol)
			if orderID != "" {
				params.SetOrdId(orderID)
			}
			if clOrderID != "" {
				params.SetClOrdId(clOrderID)
			}
			var err error
			var ordersDetail okex.OrderResponse
			request, weight, err = okex.GetOrderDetails(msg.Account.ApiCredentials, msg.Account.ApiCredentials.AccountID, "1", params)
			if err != nil {
				state.logger.Warn("error building request", log.Error(err))
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				context.Send(sender, response)
				return
			}
			if state.rateLimit.IsRateLimited() {
				response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
				context.Send(sender, response)
				return
			}

			state.rateLimit.Request(weight)
			if err := xutils.PerformJSONRequest(state.client, request, &ordersDetail); err != nil {
				state.logger.Warn("error fetching open orders", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Send(sender, response)
				return
			}
			if ordersDetail.Code != "0" {
				state.logger.Warn("error building request", log.Error(err))
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				context.Send(sender, response)
				return
			}
			for _, o := range ordersDetail.Data {
				data = append(data, o)
			}
		} else {
			var err error
			var orderHistoryData, openOrdersData okex.OrderResponse
			orderHistoryArgs := okex.GetOrderRequest("SPOT")
			request, weight, err = okex.GetOrdersHistory(msg.Account.ApiCredentials, msg.Account.ApiCredentials.AccountID, "1", orderHistoryArgs)
			if err != nil {
				state.logger.Warn("error building request", log.Error(err))
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				context.Send(sender, response)
				return
			}
			if err := xutils.PerformJSONRequest(state.client, request, &orderHistoryData); err != nil {
				state.logger.Warn("error building request", log.Error(err))
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				context.Send(sender, response)
				return
			}
			if orderHistoryData.Code != "0" {
				state.logger.Warn("error building request", log.Error(err))
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				context.Send(sender, response)
				return
			}

			openOrderArgs := okex.GetOrderRequest("SPOT")
			openOrderReq, _, err := okex.GetOpenOrders(msg.Account.ApiCredentials, msg.Account.ApiCredentials.AccountID, "1", openOrderArgs)
			if err != nil {
				state.logger.Warn("error building request", log.Error(err))
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				context.Send(sender, response)
				return
			}

			if err := xutils.PerformJSONRequest(state.client, openOrderReq, &openOrdersData); err != nil {
				state.logger.Warn("error building request", log.Error(err))
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				context.Send(sender, response)
				return
			}
			if openOrdersData.Code != "0" {
				state.logger.Warn("error building request", log.Error(err))
				response.RejectionReason = messages.RejectionReason_UnsupportedRequest
				context.Send(sender, response)
				return
			}
			sort.Slice(orderHistoryData.Data, func(i, j int) bool {
				return orderHistoryData.Data[i].UTime < orderHistoryData.Data[j].UTime
			})
			sort.Slice(openOrdersData.Data, func(i, j int) bool {
				return openOrdersData.Data[i].UTime < openOrdersData.Data[j].UTime
			})
			for _, o := range orderHistoryData.Data {
				data = append(data, o)
			}
			for _, o := range openOrdersData.Data {
				data = append(data, o)
			}
		}

		var morders []*models.Order
		for _, o := range data {
			sec := state.SymbolToSecurity(o.InstId)
			if sec == nil {
				response.RejectionReason = messages.RejectionReason_UnknownSymbol
				context.Send(sender, response)
				return
			}
			ord := WSOrderToModel(&o)
			if orderStatus != nil && ord.OrderStatus != *orderStatus {
				continue
			}
			ord.Instrument.SecurityID = &wrapperspb.UInt64Value{Value: sec.SecurityID}
			morders = append(morders, ord)
		}
		response.Success = true
		response.Orders = morders
		context.Send(sender, response)
	}()

	return nil
}

func (state *Executor) OnMarketStatisticsRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketStatisticsRequest)
	response := &messages.MarketStatisticsResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	if msg.Instrument == nil || msg.Instrument.Symbol == nil {
		response.RejectionReason = messages.RejectionReason_MissingInstrument
		context.Respond(response)
		return nil
	}
	security, rej := state.InstrumentToSecurity(msg.Instrument)

	if rej != nil {
		response.RejectionReason = *rej
		context.Respond(response)
		return nil
	}

	has := func(stat models.StatType) bool {
		for _, s := range response.Statistics {
			if s.StatType == stat {
				return true
			}
		}
		return false
	}
	for _, stat := range msg.Statistics {
		switch stat {
		case models.StatType_OpenInterest:
			if has(stat) {
				continue
			}
			request := okex.NewOpenInterestRequest(okex.SWAP)
			request.SetInstrumentID(security.Symbol)
			req, weight, err := okex.GetOpenInterest(request)
			if err != nil {
				return err
			}

			if state.rateLimit.IsRateLimited() {
				response.RejectionReason = messages.RejectionReason_IPRateLimitExceeded
				context.Respond(response)
				return nil
			}

			state.rateLimit.Request(weight)

			res, err := context.RequestFuture(state.queryRunner, &jobs.PerformHTTPQueryRequest{Request: req}, 10*time.Second).Result()
			if err != nil {
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
				return nil
			}
			queryResponse := res.(*jobs.PerformQueryResponse)
			if queryResponse.StatusCode != 200 {
				if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
					err := fmt.Errorf(
						"http client error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
					state.logger.Info("http client error", log.Error(err))
					response.RejectionReason = messages.RejectionReason_HTTPError
					context.Respond(response)
					return nil
				} else if queryResponse.StatusCode >= 500 {
					err := fmt.Errorf(
						"http server error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
					state.logger.Info("http client error", log.Error(err))
					response.RejectionReason = messages.RejectionReason_HTTPError
					context.Respond(response)
					return nil
				}
				return nil
			}

			var resp okex.OpenInterestResponse
			err = json.Unmarshal(queryResponse.Response, &resp)
			if err != nil {
				state.logger.Info("http error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return nil
			}

			if resp.Code != "0" {
				state.logger.Info("http error", log.Error(errors.New(resp.Msg)))
				response.RejectionReason = messages.RejectionReason_ExchangeAPIError
				context.Respond(response)
				return nil
			}
			for _, d := range resp.Data {
				response.Statistics = append(response.Statistics, &models.Stat{
					Timestamp: utils.MilliToTimestamp(d.Ts),
					StatType:  models.StatType_OpenInterest,
					Value:     d.OpenInterest * security.Multiplier.Value,
				})
			}
		}
	}

	response.Success = true
	context.Respond(response)

	return nil
}
