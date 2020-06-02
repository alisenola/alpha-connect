package kraken

import (
	"encoding/json"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/enum"
	"gitlab.com/alphaticks/alphac/exchanges/interface"
	"gitlab.com/alphaticks/alphac/jobs"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
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
	state.rateLimit = exchanges.NewRateLimit(8000, 10*time.Minute)
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
		baseCurrency, ok := constants.SYMBOL_TO_ASSET[baseName]
		if !ok {
			//state.logger.Info("unknown symbol " + baseName)
			continue
		}
		quoteName := strings.Split(pair.WSName, "/")[1]
		if sym, ok := kraken.KRAKEN_SYMBOL_TO_GLOBAL_SYMBOL[quoteName]; ok {
			quoteName = sym
		}
		quoteCurrency, ok := constants.SYMBOL_TO_ASSET[quoteName]
		if !ok {
			//state.logger.Info("unknown symbol " + quoteName)
			continue
		}
		security := models.Security{}
		security.Symbol = pair.WSName
		security.Underlying = &baseCurrency
		security.QuoteCurrency = &quoteCurrency
		security.Enabled = true
		security.Exchange = &constants.KRAKEN
		security.SecurityType = enum.SecurityType_CRYPTO_SPOT
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name)
		security.MinPriceIncrement = 1. / math.Pow10(pair.PairDecimals)
		security.RoundLot = 1. / math.Pow10(pair.LotDecimals)

		securities = append(securities, &security)
	}
	state.securities = securities

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Error:      "",
		Securities: state.securities})

	return nil
}

func (state *Executor) OnSecurityListRequest(context actor.Context) error {
	msg := context.Message().(*messages.SecurityListRequest)

	context.Respond(&messages.SecurityList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Error:      "",
		Securities: state.securities})

	return nil
}

func (state *Executor) OnMarketDataRequest(context actor.Context) error {
	var snapshot *models.OBL2Snapshot
	msg := context.Message().(*messages.MarketDataRequest)
	if msg.Subscribe {
		context.Respond(&messages.MarketDataRequestReject{
			RequestID: msg.RequestID,
			Reason:    "market data subscription not supported on executor"})
	}
	if msg.Instrument == nil || msg.Instrument.Symbol == nil {
		context.Respond(&messages.MarketDataRequestReject{
			RequestID: msg.RequestID,
			Reason:    "symbol needed"})
	}
	// Get http request and the expected response
	symbol := msg.Instrument.Symbol.Value
	request, weight, err := kraken.GetOrderBook(symbol)
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
			context.Respond(&messages.MarketDataRequestReject{
				RequestID: msg.RequestID,
				Reason:    err.Error()})
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)

		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"http client error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				context.Respond(&messages.MarketDataRequestReject{
					RequestID: msg.RequestID,
					Reason:    err.Error()})
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"http server error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				context.Respond(&messages.MarketDataRequestReject{
					RequestID: msg.RequestID,
					Reason:    err.Error()})
			}
			return
		}

		var response struct {
			Error  []string                      `json:"error"`
			Result map[string]kraken.OrderBookL2 `json:"result"`
		}
		err = json.Unmarshal(queryResponse.Response, &response)
		if err != nil {
			err = fmt.Errorf("error decoding query response: %v", err)
			context.Respond(&messages.MarketDataRequestReject{
				RequestID: msg.RequestID,
				Reason:    err.Error()})
			return
		}

		bids, asks := response.Result[symbol].ToBidAsk()
		var maxTs uint64 = 0
		for _, bid := range response.Result[symbol].Bids {
			if bid.Time > maxTs {
				maxTs = bid.Time
			}
		}
		for _, ask := range response.Result[symbol].Asks {
			if ask.Time > maxTs {
				maxTs = ask.Time
			}
		}
		snapshot = &models.OBL2Snapshot{
			Bids:      bids,
			Asks:      asks,
			Timestamp: utils.MicroToTimestamp(maxTs),
		}
		context.Respond(&messages.MarketDataSnapshot{
			RequestID:  msg.RequestID,
			ResponseID: uint64(time.Now().UnixNano()),
			SnapshotL2: snapshot,
			SeqNum:     maxTs,
		})
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

func (state *Executor) OnOrderCancelRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnOrderMassCancelRequest(context actor.Context) error {
	return nil
}
