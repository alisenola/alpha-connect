package ftx

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
	"gitlab.com/alphaticks/xchanger/exchanges/ftx"
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

// The global rate limit is per IP and the orderRateLimit is per
// account.

type Executor struct {
	_interface.ExchangeExecutorBase
	client      *http.Client
	securities  []*models.Security
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

	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	request, weight, err := ftx.GetMarkets()
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
	type Res struct {
		Success bool         `json:"success"`
		Markets []ftx.Market `json:"result"`
	}
	var queryRes Res
	err = json.Unmarshal(response, &queryRes)
	if err != nil {
		err = fmt.Errorf(
			"error unmarshaling response: %v",
			err)
		return err
	}
	if !queryRes.Success {
		err = fmt.Errorf("ftx unsuccessful request")
		return err
	}

	var securities []*models.Security
	for _, market := range queryRes.Markets {
		security := models.Security{}
		security.Symbol = market.Name
		security.Exchange = &constants.FTX

		switch market.Type {
		case "spot":
			baseCurrency, ok := constants.GetAssetBySymbol(market.BaseCurrency)
			if !ok {
				//fmt.Printf("unknown currency symbol %s \n", market.BaseCurrency)
				continue
			}
			quoteCurrency, ok := constants.GetAssetBySymbol(market.QuoteCurrency)
			if !ok {
				//fmt.Printf("unknown currency symbol %s \n", market.BaseCurrency)
				continue
			}
			security.Underlying = baseCurrency
			security.QuoteCurrency = quoteCurrency
			security.SecurityType = enum.SecurityType_CRYPTO_SPOT

		case "future":
			splits := strings.Split(market.Name, "-")
			if len(splits) == 2 {
				underlying, ok := constants.GetAssetBySymbol(market.Underlying)
				if !ok {
					//fmt.Printf("unknown currency symbol %s \n", market.Underlying)
					continue
				}
				security.Underlying = underlying
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
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinPriceIncrement = &types.DoubleValue{Value: market.PriceIncrement}
		security.RoundLot = &types.DoubleValue{Value: market.SizeIncrement}
		security.Status = models.Trading

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
	// Get http request and the expected response
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

func (state *Executor) OnMarketStatisticsRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketStatisticsResponse)
	context.Respond(&messages.MarketStatisticsResponse{
		RequestID:       msg.RequestID,
		Success:         false,
		RejectionReason: messages.UnsupportedRequest,
	})
	return nil
}

func (state *Executor) OnMarketDataRequest(context actor.Context) error {
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
