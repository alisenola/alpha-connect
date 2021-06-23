package cryptofacilities

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
	"gitlab.com/alphaticks/xchanger/exchanges/cryptofacilities"
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
	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewAPIQuery(state.client)
	})
	state.queryRunner = context.Spawn(props)

	state.rateLimit = exchanges.NewRateLimit(500, 10*time.Second)
	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	if state.rateLimit.IsRateLimited() {
		return fmt.Errorf("rate limited")
	}
	request, _, err := cryptofacilities.GetInstruments()
	if err != nil {
		return err
	}
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

	var instrumentResponse cryptofacilities.InstrumentResponse
	err = json.Unmarshal(response, &instrumentResponse)
	if err != nil {
		err = fmt.Errorf("error decoding query response: %v", err)
		return err
	}
	if instrumentResponse.Result != "success" {
		err = fmt.Errorf("error getting instruments: %s", instrumentResponse.Error)
		return err
	}

	var securities []*models.Security
	for _, instrument := range instrumentResponse.Instruments {
		if !instrument.Tradeable {
			continue
		}
		underlying := strings.ToUpper(strings.Split(instrument.Underlying, "_")[1])
		baseSymbol := underlying[:3]
		if sym, ok := cryptofacilities.CRYPTOFACILITIES_SYMBOL_TO_GLOBAL_SYMBOL[baseSymbol]; ok {
			baseSymbol = sym
		}
		baseCurrency, ok := constants.GetAssetBySymbol(baseSymbol)
		if !ok {
			continue
		}

		quoteSymbol := underlying[3:]
		if sym, ok := cryptofacilities.CRYPTOFACILITIES_SYMBOL_TO_GLOBAL_SYMBOL[quoteSymbol]; ok {
			quoteSymbol = sym
		}
		quoteCurrency, ok := constants.GetAssetBySymbol(quoteSymbol)
		if !ok {
			continue
		}
		security := models.Security{}
		security.Symbol = instrument.Symbol
		security.Underlying = baseCurrency
		security.QuoteCurrency = quoteCurrency
		if instrument.Tradeable {
			security.Status = models.Trading
		} else {
			security.Status = models.Disabled
		}
		security.Exchange = &constants.CRYPTOFACILITIES
		splits := strings.Split(instrument.Symbol, "_")
		switch splits[0] {
		case "pv":
			// Perpetual vanilla
			security.SecurityType = enum.SecurityType_CRYPTO_PERP
			security.IsInverse = false
		case "pi":
			// Perpetual inverse
			security.SecurityType = enum.SecurityType_CRYPTO_PERP
			security.IsInverse = true
		case "fv":
			// Future vanilla
			security.SecurityType = enum.SecurityType_CRYPTO_FUT
			security.IsInverse = false
			year := time.Now().Format("2006")
			date, err := time.Parse("20060102", year+splits[2][2:])
			if err != nil {
				continue
			}
			security.MaturityDate, err = types.TimestampProto(date)
			if err != nil {
				continue
			}
		case "fi":
			// Future inverse
			security.SecurityType = enum.SecurityType_CRYPTO_FUT
			security.IsInverse = true
			year := time.Now().Format("2006")
			date, err := time.Parse("20060102", year+splits[2][2:])
			if err != nil {
				continue
			}
			security.MaturityDate, err = types.TimestampProto(date)
			if err != nil {
				continue
			}
		}
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinPriceIncrement = &types.DoubleValue{Value: instrument.TickSize}
		security.RoundLot = &types.DoubleValue{Value: float64(instrument.ContractSize)}
		securities = append(securities, &security)
	}

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: state.securities})

	state.securities = securities
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
