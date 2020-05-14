package ftx

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
	"gitlab.com/alphaticks/xchanger/exchanges/ftx"
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

// The global rate limit is per IP and the orderRateLimit is per
// account.

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

	state.rateLimit = exchanges.NewRateLimit(30, time.Second)
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
	// Get http request and the expected response
	msg := context.Message().(*messages.SecurityListRequest)
	request, weight, err := ftx.GetMarkets()
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		time.Sleep(state.rateLimit.DurationBeforeNextRequest(weight))
		return nil
	}

	state.rateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res interface{}, err error) {
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
		type Res struct {
			Success bool         `json:"success"`
			Markets []ftx.Market `json:"result"`
		}
		var queryRes Res
		err = json.Unmarshal(queryResponse.Response, &queryRes)
		if err != nil {
			err = fmt.Errorf(
				"error unmarshaling response: %v",
				err)
			context.Respond(&messages.SecurityList{
				RequestID:  msg.RequestID,
				ResponseID: uint64(time.Now().UnixNano()),
				Error:      err.Error(),
				Securities: nil})
			return
		}
		if !queryRes.Success {
			err = fmt.Errorf("ftx unsuccessful request")
			context.Respond(&messages.SecurityList{
				RequestID:  msg.RequestID,
				ResponseID: uint64(time.Now().UnixNano()),
				Error:      err.Error(),
				Securities: nil})
			return
		}

		var securities []*models.Security
		for _, market := range queryRes.Markets {
			security := models.Security{}
			security.Symbol = market.Name
			security.Exchange = &constants.FTX

			switch market.Type {
			case "spot":
				baseCurrency, ok := constants.SYMBOL_TO_ASSET[market.BaseCurrency]
				if !ok {
					//fmt.Printf("unknown currency symbol %s \n", market.BaseCurrency)
					continue
				}
				quoteCurrency, ok := constants.SYMBOL_TO_ASSET[market.QuoteCurrency]
				if !ok {
					//fmt.Printf("unknown currency symbol %s \n", market.BaseCurrency)
					continue
				}
				security.Underlying = &baseCurrency
				security.QuoteCurrency = &quoteCurrency
				security.SecurityType = enum.SecurityType_CRYPTO_SPOT

			case "future":
				splits := strings.Split(market.Name, "-")
				if len(splits) == 2 {
					underlying, ok := constants.SYMBOL_TO_ASSET[market.Underlying]
					if !ok {
						fmt.Printf("unknown currency symbol %s \n", market.Underlying)
						continue
					}
					security.Underlying = &underlying
					security.QuoteCurrency = &constants.DOLLAR
					if splits[1] == "PERP" {
						security.SecurityType = enum.SecurityType_CRYPTO_PERP
					} else {
						year := time.Now().Format("2006")
						date, err := time.Parse("20060102", year+splits[1])
						if err != nil {
							continue
						}
						security.SecurityType = enum.SecurityType_CRYPTO_FUT
						security.MaturityDate, err = types.TimestampProto(date)
						if err != nil {
							continue
						}
					}
				} else {
					continue
				}

			default:
				continue
			}
			security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name)
			security.MinPriceIncrement = market.PriceIncrement
			security.RoundLot = market.SizeIncrement

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
