package bittrex

/*
import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/jobs"
	"gitlab.com/alphaticks/alpha-connect/models/messages/executor"
	"gitlab.com/alphaticks/xchanger/asset"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/bittrex"
	"net/http"
	"reflect"
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
		client:      nil,
		rateLimit:   nil,
		queryRunner: nil,
		logger:      nil,
	}
}

func (state *Executor) Receive(context actor.Context) {
	_interface.ExecutorReceive(state, context)
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
		Timeout:   10 * time.Second,
	}
	state.rateLimit = exchanges.NewRateLimit(3, time.Second)

	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewAPIQuery(state.client)
	})
	state.queryRunner = context.Spawn(props)

	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) GetInstrumentsRequest(context actor.Context) error {
	// Get http request and the expected response
	msg := context.Message().(*executor.GetInstrumentsRequest)
	request, weight, err := bittrex.GetMarkets()
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		time.Sleep(state.rateLimit.DurationBeforeNextRequest(weight))
		return nil
	}
	// TODO Rate limit

	future := context.RequestFuture(state.queryRunner, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res types{}, err error) {
		if err != nil {
			context.Respond(&executor.GetInstrumentsResponse{
				RequestID:   msg.RequestID,
				Error:       err,
				Instruments: nil})
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"http client error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				context.Respond(&executor.GetInstrumentsResponse{
					RequestID:   msg.RequestID,
					Error:       err,
					Instruments: nil})
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"http server error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				context.Respond(&executor.GetInstrumentsResponse{
					RequestID:   msg.RequestID,
					Error:       err,
					Instruments: nil})
			}
			return
		}
		type Result struct {
			Message string           `json:"message"`
			Success bool             `json:"success"`
			Result  []bittrex.Market `json:"result"`
		}
		var response Result
		err = json.Unmarshal(queryResponse.Response, &response)
		if err != nil {
			err = fmt.Errorf(
				"error unmarshaling response: %v",
				err)
			context.Respond(&executor.GetInstrumentsResponse{
				RequestID:   msg.RequestID,
				Error:       err,
				Instruments: nil})
			return
		}
		if !response.Success {
			context.Respond(&executor.GetInstrumentsResponse{
				RequestID:   msg.RequestID,
				Error:       errors.New(response.Message),
				Instruments: nil})
			return
		}

		var instruments []*exchanges.Instrument
		for _, market := range response.Result {
			if !market.IsActive {
				continue
			}
			baseCurrency, ok := constants.GetAssetBySymbol(market.MarketCurrency)
			if !ok {
				//state.logger.Info(fmt.Sprintf("unknown currency %s", market.MarketCurrency))
				continue
			}
			quoteCurrency, ok := constants.GetAssetBySymbol(market.BaseCurrency)
			if !ok {
				//state.logger.Info(fmt.Sprintf("unknown currency %s", market.BaseCurrency))
				continue
			}
			pair := asset.NewPair(&baseCurrency, &quoteCurrency)
			instrument := exchanges.Instrument{}
			instrument.Exchange = constants.BITTREX
			instrument.Pair = pair
			instrument.Type = exchanges.SPOT

			pres, ok := bittrex.TickPrecisions[pair.String()]
			if !ok {
				continue
			}
			instrument.TickPrecision = pres
			instrument.LotPrecision = 1e8

			instruments = append(instruments, &instrument)
		}
		context.Respond(&executor.GetInstrumentsResponse{
			RequestID:   msg.RequestID,
			Error:       nil,
			Instruments: instruments})
	})

	return nil
}

func (state *Executor) GetOrderBookL2Request(context actor.Context) error {
	return nil
}

func (state *Executor) GetOrderBookL3Request(context actor.Context) error {
	return nil
}

func (state *Executor) GetOpenOrdersRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OpenOrdersRequest(context actor.Context) error {
	return nil
}

func (state *Executor) CloseOrdersRequest(context actor.Context) error {
	return nil
}

func (state *Executor) CloseAllOrdersRequest(context actor.Context) error {
	return nil
}

*/
