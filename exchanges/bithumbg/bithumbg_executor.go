package bithumbg

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
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/bithumbg"
	"io/ioutil"
	"math"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Executor struct {
	client      *http.Client
	securities  map[uint64]*models.Security
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
	_interface.ExchangeExecutorReceive(state, context)
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

	state.client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1024,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
	// TODO
	state.rateLimit = exchanges.NewRateLimit(10, time.Second)
	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewAPIQuery(state.client)
	})
	state.queryRunner = context.Spawn(props)
	state.securities = make(map[uint64]*models.Security)
	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	request, weight, err := bithumbg.GetSpotConfig()
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

	var res bithumbg.Response
	if err := json.Unmarshal(response, &res); err != nil {
		err = fmt.Errorf(
			"error unmarshaling response: %v",
			err)
		return err
	}

	var config bithumbg.Config
	if err := json.Unmarshal(res.Data, &config); err != nil {
		err = fmt.Errorf(
			"error unmarshaling config: %v",
			err)
		return err
	}

	var securities []*models.Security
	for _, symbol := range config.Pairs {
		splits := strings.Split(symbol.Symbol, "-")
		baseStr := strings.ToUpper(splits[0])
		quoteStr := strings.ToUpper(splits[1])

		baseCurrency, ok := constants.GetAssetBySymbol(baseStr)
		if !ok {
			//state.logger.Info(fmt.Sprintf("unknown currency %s", baseStr))
			continue
		}
		quoteCurrency, ok := constants.GetAssetBySymbol(quoteStr)
		if !ok {
			//state.logger.Info(fmt.Sprintf("unknown currency %s", quoteStr))
			continue
		}
		security := models.Security{}
		security.Symbol = symbol.Symbol
		security.Underlying = baseCurrency
		security.QuoteCurrency = quoteCurrency
		security.Status = models.Trading
		security.Exchange = &constants.BITHUMBG
		security.SecurityType = enum.SecurityType_CRYPTO_SPOT
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name)
		security.IsInverse = false
		security.MinPriceIncrement = &types.DoubleValue{Value: 1. / math.Pow(10, float64(symbol.Accuracy[0]))}
		security.RoundLot = &types.DoubleValue{Value: 1. / math.Pow(10, float64(symbol.Accuracy[1]))}
		securities = append(securities, &security)
	}
	state.securities = make(map[uint64]*models.Security)
	for _, sec := range securities {
		state.securities[sec.SecurityID] = sec
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
	msg := context.Message().(*messages.MarketDataRequest)
	response := &messages.MarketDataResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	if state.rateLimit.IsRateLimited() {
		response.RejectionReason = messages.RateLimitExceeded
		context.Respond(response)
		return nil
	}
	if msg.Subscribe {
		response.RejectionReason = messages.UnsupportedSubscription
		context.Respond(response)
		return nil
	}
	if msg.Instrument == nil {
		response.RejectionReason = messages.MissingInstrument
		context.Respond(response)
		return nil
	}
	var symbol = ""
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
	if symbol == "" {
		response.RejectionReason = messages.UnknownSymbol
		context.Respond(response)
		return nil
	}

	if msg.Aggregation == models.L2 {
		var snapshot *models.OBL2Snapshot
		// Get http request and the expected response
		request, weight, err := bithumbg.GetSpotOrderBook(symbol)
		if err != nil {
			return err
		}

		state.rateLimit.Request(weight)
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
				} else if queryResponse.StatusCode >= 500 {
					err := fmt.Errorf(
						"http server error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
					state.logger.Info("http client error", log.Error(err))
					response.RejectionReason = messages.HTTPError
					context.Respond(response)
				}
				return
			}
			var bresponse bithumbg.Response
			if err := json.Unmarshal(queryResponse.Response, &bresponse); err != nil {
				state.logger.Info("error decoding query response", log.Error(err))
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
				return
			}
			if bresponse.Code != "0" {
				state.logger.Info("error getting order book data", log.Error(errors.New(bresponse.Msg)))
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
				return
			}
			var obData bithumbg.OrderBook
			if err := json.Unmarshal(bresponse.Data, &obData); err != nil {
				state.logger.Info("error decoding query response", log.Error(err))
				response.RejectionReason = messages.HTTPError
				context.Respond(response)
				return
			}

			bidst, askst := obData.ToBidAsk()
			var bids, asks []gorderbook.OrderBookLevel
			for _, b := range bidst {
				if b.Quantity == 0. {
					continue
				}
				bids = append(bids, b)
			}
			for _, a := range askst {
				if a.Quantity == 0. {
					continue
				}
				asks = append(asks, a)
			}
			snapshot = &models.OBL2Snapshot{
				Bids:      bids,
				Asks:      asks,
				Timestamp: utils.MilliToTimestamp(0),
			}
			response.SnapshotL2 = snapshot
			response.SeqNum = uint64(obData.Version)
			response.Success = true
			context.Respond(response)
		})

		return nil
	} else {

		return nil
	}
}

func (state *Executor) OnOrderStatusRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnPositionsRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnBalancesRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnNewOrderSingleRequest(context actor.Context) error {
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
	return nil
}

func (state *Executor) OnOrderMassCancelRequest(context actor.Context) error {
	return nil
}
