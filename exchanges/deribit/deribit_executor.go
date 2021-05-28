package deribit

import (
	"encoding/json"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/enum"
	_interface "gitlab.com/alphaticks/alpha-connect/exchanges/interface"
	"gitlab.com/alphaticks/alpha-connect/jobs"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/deribit"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

type Executor struct {
	client      *http.Client
	securities  map[uint64]*models.Security
	symbolToSec map[string]*models.Security
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
	state.client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1024,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
	state.rateLimit = exchanges.NewRateLimit(8000, 10*time.Minute)
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))
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
	// Fetch currencies first
	if state.rateLimit.IsRateLimited() {
		return fmt.Errorf("rate limited")
	}

	req, weight, err := deribit.GetCurrencies()
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

	type Res struct {
		Result []deribit.Currency `json:"result"`
		Error  *deribit.RPCError  `json:"error"`
	}
	var res Res
	err = json.Unmarshal(response, &res)
	if err != nil {
		err = fmt.Errorf("error decoding query response: %v", err)
		return err
	}

	if res.Error != nil {
		return fmt.Errorf("error in query response: %s", res.Error.Message)
	}

	var securities []*models.Security
	for _, c := range res.Result {
		req, weight, err = deribit.GetInstruments(deribit.GetInstrumentsParams{
			Currency: c.Currency,
			Kind:     nil,
			Expired:  nil,
		})
		if err != nil {
			return err
		}
		if state.rateLimit.IsRateLimited() {
			return fmt.Errorf("rate limited")
		}
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

		type Res struct {
			Result []deribit.Instrument `json:"result"`
			Error  *deribit.RPCError    `json:"error"`
		}
		var res Res
		err = json.Unmarshal(response, &res)
		if err != nil {
			err = fmt.Errorf("error decoding query response: %v", err)
			return err
		}

		if res.Error != nil {
			return fmt.Errorf("error in query response: %s", res.Error.Message)
		}

		for _, i := range res.Result {
			if i.Kind == "option" {
				continue
			}
			security := models.Security{}
			if i.IsActive {
				security.Status = models.Trading
			} else {
				security.Status = models.Disabled
			}
			security.Symbol = i.InstrumentName
			baseCurrency, ok := constants.GetAssetBySymbol(i.BaseCurrency)
			if !ok {
				continue
			}
			security.Underlying = baseCurrency
			quoteCurrency, ok := constants.GetAssetBySymbol(i.QuoteCurrency)
			if !ok {
				continue
			}
			security.QuoteCurrency = quoteCurrency
			security.Exchange = &constants.DERIBIT
			if i.SettlementPeriod == "perpetual" {
				security.SecurityType = enum.SecurityType_CRYPTO_PERP
			} else {
				security.SecurityType = enum.SecurityType_CRYPTO_FUT
				security.MaturityDate = utils.MilliToTimestamp(i.ExpirationTimestamp)
			}
			security.MakerFee = &types.DoubleValue{Value: i.MakerCommission}
			security.TakerFee = &types.DoubleValue{Value: i.TakerCommission}

			security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
			security.MinPriceIncrement = &types.DoubleValue{Value: i.TickSize}
			security.RoundLot = &types.DoubleValue{Value: i.ContractSize}
			security.IsInverse = true

			securities = append(securities, &security)
		}
	}

	state.securities = make(map[uint64]*models.Security)
	state.symbolToSec = make(map[string]*models.Security)
	for _, s := range securities {
		state.securities[s.SecurityID] = s
		state.symbolToSec[s.Symbol] = s
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

func (state *Executor) OnHistoricalLiquidationsRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnMarketDataRequest(context actor.Context) error {
	return nil
}

func (state *Executor) OnOrderStatusRequest(context actor.Context) error {
	// TODO forward to private executor
	return nil
}

func (state *Executor) OnPositionsRequest(context actor.Context) error {
	// TODO forward to private executor
	return nil
}

func (state *Executor) OnBalancesRequest(context actor.Context) error {
	// TODO forward to private executor
	return nil
}

func (state *Executor) OnNewOrderSingleRequest(context actor.Context) error {
	// TODO forward to private executor
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
	// TODO forward to private executor
	return nil
}

func (state *Executor) OnOrderMassCancelRequest(context actor.Context) error {
	// TODO forward to private executor
	return nil
}
