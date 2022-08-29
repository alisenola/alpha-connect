package deribit

import (
	"encoding/json"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/enum"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/jobs"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/deribit"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

type Executor struct {
	extypes.BaseExecutor
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
	extypes.ReceiveExecutor(state, context)
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
		log.String("type", reflect.TypeOf(state).String()))
	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewHTTPQuery(state.client)
	})
	state.queryRunner = context.Spawn(props)

	return state.UpdateSecurityList(context)
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	// Fetch currencies first
	if state.rateLimit.IsRateLimited() {
		return fmt.Errorf("rate limited")
	}

	req, weight, err := deribit.GetCurrencies()
	if err != nil {
		return err
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
			security := models.Security{}
			if i.Kind == "future" {
				if i.SettlementPeriod == "perpetual" {
					security.SecurityType = enum.SecurityType_CRYPTO_PERP
				} else {
					security.SecurityType = enum.SecurityType_CRYPTO_FUT
					security.MaturityDate = utils.MilliToTimestamp(i.ExpirationTimestamp)
				}
			} else if i.Kind == "option" {
				security.SecurityType = enum.SecurityType_CRYPTO_OPT
				if i.OptionType == "call" {
					security.SecuritySubType = wrapperspb.String(enum.SecuritySubType_CALL)
				} else {
					security.SecuritySubType = wrapperspb.String(enum.SecuritySubType_PUT)
				}
				security.StrikePrice = wrapperspb.Double(i.Strike)
			} else {
				// Skip combos
				continue
			}

			if i.IsActive {
				security.Status = models.InstrumentStatus_Trading
			} else {
				security.Status = models.InstrumentStatus_Disabled
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
			security.Exchange = constants.DERIBIT

			security.MakerFee = wrapperspb.Double(i.MakerCommission)
			security.TakerFee = wrapperspb.Double(i.TakerCommission)

			security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
			security.MinPriceIncrement = wrapperspb.Double(i.TickSize)
			security.RoundLot = wrapperspb.Double(i.ContractSize)
			security.IsInverse = true

			securities = append(securities, &security)
		}
	}

	state.SyncSecurities(securities, nil)

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: securities})

	return nil
}

func (state *Executor) OnHistoricalLiquidationsRequest(context actor.Context) error {
	msg := context.Message().(*messages.HistoricalLiquidationsRequest)
	context.Respond(&messages.HistoricalLiquidationsResponse{
		RequestID:       msg.RequestID,
		Success:         false,
		RejectionReason: messages.RejectionReason_UnsupportedRequest,
	})
	return nil
}

func (state *Executor) OnMarketStatisticsRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketStatisticsRequest)
	context.Respond(&messages.MarketStatisticsResponse{
		RequestID:       msg.RequestID,
		Success:         false,
		RejectionReason: messages.RejectionReason_UnsupportedRequest,
	})
	return nil
}
