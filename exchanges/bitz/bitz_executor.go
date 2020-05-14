package bitz

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
	"gitlab.com/alphaticks/xchanger/exchanges/bitz"
	"math"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Executor struct {
	client      *http.Client
	rateLimit   *exchanges.RateLimit
	queryRunner *actor.PID
	logger      *log.Logger
}

func NewExecutor() actor.Actor {
	return &Executor{
		client:    nil,
		rateLimit: nil,
		logger:    nil,
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
	//TODO state.httpRateLimit = exchanges.NewRateLimit()
	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewAPIQuery(state.client)
	})
	state.queryRunner = context.Spawn(props)

	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) OnSecurityListRequest(context actor.Context) error {
	msg := context.Message().(*messages.SecurityListRequest)
	request, _, err := bitz.GetSymbolList()
	if err != nil {
		return err
	}

	/* TODO
	if state.httpRateLimit.IsRateLimited() {
		time.Sleep(state.httpRateLimit.DurationBeforeNextRequest(weight))
	}
	*/

	// TODO state.httpRateLimit.Request(weight)

	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
	context.AwaitFuture(future, func(res interface{}, err error) {
		context.Stop(state.queryRunner)
		if err != nil {
			context.Respond(&messages.SecurityList{
				RequestID:  msg.RequestID,
				ResponseID: uint64(time.Now().UnixNano()),
				Error:      err.Error(),
				Securities: nil})
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"http client error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				context.Respond(&messages.SecurityList{
					RequestID:  msg.RequestID,
					ResponseID: uint64(time.Now().UnixNano()),
					Error:      err.Error(),
					Securities: nil})
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"http server error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				context.Respond(&messages.SecurityList{
					RequestID:  msg.RequestID,
					ResponseID: uint64(time.Now().UnixNano()),
					Error:      err.Error(),
					Securities: nil})
			}
			return
		}

		var response bitz.SymbolList
		err = json.Unmarshal(queryResponse.Response, &response)
		if err != nil {
			err = fmt.Errorf("error decoding query response: %v", err)
			context.Respond(&messages.SecurityList{
				RequestID:  msg.RequestID,
				ResponseID: uint64(time.Now().UnixNano()),
				Error:      err.Error(),
				Securities: nil})
			return
		}

		if response.Status != 200 {
			err = fmt.Errorf("got incorrect status code: %d", response.Status)
			context.Respond(&messages.SecurityList{
				RequestID:  msg.RequestID,
				ResponseID: uint64(time.Now().UnixNano()),
				Error:      err.Error(),
				Securities: nil})
			return
		}

		var securities []*models.Security
		for _, pair := range response.Data {
			baseCurrency, ok := constants.SYMBOL_TO_ASSET[strings.ToUpper(pair.CoinFrom)]
			if !ok {
				//state.logger.Info("unknown symbol " + pair.BaseCurrency + " for instrument " + pair.InstrumentID)
				continue
			}
			quoteCurrency, ok := constants.SYMBOL_TO_ASSET[strings.ToUpper(pair.CoinTo)]
			if !ok {
				//state.logger.Info("unknown symbol " + pair.QuoteCurrency + " for instrument " + pair.InstrumentID)
				continue
			}

			security := models.Security{}
			security.Symbol = pair.Name
			security.Underlying = &baseCurrency
			security.QuoteCurrency = &quoteCurrency
			security.Enabled = true
			security.Exchange = &constants.BITZ
			security.SecurityType = enum.SecurityType_CRYPTO_SPOT
			security.RoundLot = 1. / math.Pow(10, float64(pair.NumberFloat))
			security.MinPriceIncrement = 1. / math.Pow(10, float64(pair.PriceFloat))
			security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name)
			securities = append(securities, &security)
		}
		context.Respond(&messages.SecurityList{
			RequestID:  msg.RequestID,
			ResponseID: uint64(time.Now().UnixNano()),
			Error:      "",
			Securities: securities})
	})

	return nil
}

func (state *Executor) OnMarketDataRequest(context actor.Context) error {
	return nil
}
