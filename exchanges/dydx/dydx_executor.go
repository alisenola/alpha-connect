package dydx

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
	"gitlab.com/alphaticks/xchanger/exchanges/dydx"
	"io/ioutil"
	"net/http"
	"reflect"
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

	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	// TODO
	// Fetch currencies first
	if state.rateLimit.IsRateLimited() {
		return fmt.Errorf("rate limited")
	}

	req, weight, err := dydx.GetMarkets()
	state.rateLimit.Request(weight)

	resp, err := state.client.Do(req)
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

	var res dydx.Markets
	err = json.Unmarshal(response, &res)
	if err != nil {
		err = fmt.Errorf("error decoding query response: %v", err)
		return err
	}

	state.securities = nil
	var securities []*models.Security
	for _, m := range res.Markets {
		security := models.Security{}
		//ONLINE, OFFLINE, POST_ONLY or CANCEL_ONLY.
		switch m.Status {
		case "ONLINE":
			security.Status = models.Trading
		case "OFFLINE":
			security.Status = models.Halt
		case "POST_ONLY":
			security.Status = models.PreTrading
		case "CANCEL_ONLY":
			security.Status = models.PostTrading
		}
		security.Symbol = m.Market
		baseStr := m.BaseAsset
		if sym, ok := dydx.SYMBOLS[baseStr]; ok {
			baseStr = sym
		}
		quoteStr := m.QuoteAsset
		if sym, ok := dydx.SYMBOLS[quoteStr]; ok {
			quoteStr = sym
		}
		baseCurrency, ok := constants.GetAssetBySymbol(baseStr)
		if !ok {
			continue
		}
		security.Underlying = baseCurrency
		quoteCurrency, ok := constants.GetAssetBySymbol(quoteStr)
		if !ok {
			continue
		}
		security.QuoteCurrency = quoteCurrency
		security.Exchange = &constants.DYDX
		security.SecurityType = enum.SecurityType_CRYPTO_PERP
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name)
		security.MinPriceIncrement = &types.DoubleValue{Value: m.TickSize}
		security.RoundLot = &types.DoubleValue{Value: m.StepSize}
		security.IsInverse = false

		fmt.Println(security)
		securities = append(securities, &security)
	}

	state.securities = make(map[uint64]*models.Security)
	for _, s := range securities {
		state.securities[s.SecurityID] = s
	}

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: securities,
	})

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
		Securities: securities,
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
