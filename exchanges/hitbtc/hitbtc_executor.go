package hitbtc

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
	"gitlab.com/alphaticks/xchanger/exchanges/hitbtc"
	"io/ioutil"
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
	state.rateLimit = exchanges.NewRateLimit(100, time.Second)
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
	request, weight, err := hitbtc.GetSymbols()
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

	var symbols []hitbtc.Symbol
	err = json.Unmarshal(response, &symbols)
	if err != nil {
		err = fmt.Errorf(
			"error unmarshaling response: %v",
			err)
		return err
	}

	var securities []*models.Security
	for _, symbol := range symbols {
		baseStr := strings.ToUpper(symbol.BaseCurrency)
		if sym, ok := hitbtc.HITBTC_SYMBOL_TO_GLOBAL_SYMBOL[baseStr]; ok {
			baseStr = sym
		}
		quoteStr := strings.ToUpper(symbol.QuoteCurrency)
		if sym, ok := hitbtc.HITBTC_SYMBOL_TO_GLOBAL_SYMBOL[quoteStr]; ok {
			quoteStr = sym
		}

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
		security.Symbol = symbol.ID
		security.Underlying = baseCurrency
		security.QuoteCurrency = quoteCurrency
		security.Status = models.Trading
		security.Exchange = &constants.HITBTC
		security.SecurityType = enum.SecurityType_CRYPTO_SPOT
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name)
		security.MinPriceIncrement = &types.DoubleValue{Value: symbol.TickSize}
		security.RoundLot = &types.DoubleValue{Value: symbol.QuantityIncrement}

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
