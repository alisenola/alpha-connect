package kraken

import (
	"encoding/json"
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
	"gitlab.com/alphaticks/xchanger/exchanges/kraken"
	"io/ioutil"
	"math"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Executor struct {
	client      *http.Client
	securities  []*models.Security
	rateLimit   *exchanges.RateLimit
	queryRunner *actor.PID
	logger      *log.Logger
}

func NewExecutor() actor.Actor {
	return &Executor{
		client:      nil,
		securities:  nil,
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
	state.rateLimit = exchanges.NewRateLimit(1, time.Second)
	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewAPIQuery(state.client)
	})
	state.queryRunner = context.Spawn(props)

	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	request, weight, err := kraken.GetAssetPairs()
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
	var kResponse struct {
		Error  []string                    `json:"error"`
		Result map[string]kraken.AssetPair `json:"result"`
	}
	err = json.Unmarshal(response, &kResponse)
	if err != nil {
		err = fmt.Errorf("error decoding query response: %v", err)
		return err
	}

	var securities []*models.Security
	for _, pair := range kResponse.Result {
		if pair.WSName == "" {
			// Dark pool pair
			continue
		}
		baseName := strings.Split(pair.WSName, "/")[0]
		if sym, ok := kraken.KRAKEN_SYMBOL_TO_GLOBAL_SYMBOL[baseName]; ok {
			baseName = sym
		}
		baseCurrency, ok := constants.GetAssetBySymbol(baseName)
		if !ok {
			//state.logger.Info("unknown symbol " + baseName)
			continue
		}
		quoteName := strings.Split(pair.WSName, "/")[1]
		if sym, ok := kraken.KRAKEN_SYMBOL_TO_GLOBAL_SYMBOL[quoteName]; ok {
			quoteName = sym
		}
		quoteCurrency, ok := constants.GetAssetBySymbol(quoteName)
		if !ok {
			//state.logger.Info("unknown symbol " + quoteName)
			continue
		}
		security := models.Security{}
		security.Symbol = pair.WSName
		security.Underlying = baseCurrency
		security.QuoteCurrency = quoteCurrency
		security.Status = models.Trading
		security.Exchange = &constants.KRAKEN
		security.SecurityType = enum.SecurityType_CRYPTO_SPOT
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinPriceIncrement = &types.DoubleValue{Value: 1. / math.Pow10(pair.PairDecimals)}
		security.RoundLot = &types.DoubleValue{Value: 1. / math.Pow10(pair.LotDecimals)}

		securities = append(securities, &security)
	}
	state.securities = securities

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: state.securities})

	return nil
}

func (state *Executor) OnSecurityListRequest(context actor.Context) error {
	msg := context.Message().(*messages.SecurityListRequest)

	context.Respond(&messages.SecurityList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: state.securities})

	return nil
}

func (state *Executor) OnHistoricalLiquidationsRequest(context actor.Context) error {
	msg := context.Message().(*messages.HistoricalLiquidationsRequest)
	context.Respond(&messages.HistoricalLiquidationsResponse{
		RequestID:       msg.RequestID,
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
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
	// Get http request and the expected response
	symbol := msg.Instrument.Symbol.Value
	request, weight, err := kraken.GetOrderBook(symbol)
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

		var apiResponse struct {
			Error  []string                      `json:"error"`
			Result map[string]kraken.OrderBookL2 `json:"result"`
		}
		err = json.Unmarshal(queryResponse.Response, &apiResponse)
		if err != nil {
			err = fmt.Errorf("error decoding query response: %v", err)
			state.logger.Info("http client error", log.Error(err))
			response.RejectionReason = messages.ExchangeAPIError
			context.Respond(response)
			return
		}

		bids, asks := apiResponse.Result[symbol].ToBidAsk()
		var maxTs uint64 = 0
		for _, bid := range apiResponse.Result[symbol].Bids {
			if bid.Time > maxTs {
				maxTs = bid.Time
			}
		}
		for _, ask := range apiResponse.Result[symbol].Asks {
			if ask.Time > maxTs {
				maxTs = ask.Time
			}
		}
		snapshot = &models.OBL2Snapshot{
			Bids:      bids,
			Asks:      asks,
			Timestamp: utils.MicroToTimestamp(maxTs),
		}
		response.Success = true
		response.SnapshotL2 = snapshot
		response.SeqNum = maxTs
		context.Respond(response)
	})
	return nil
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
