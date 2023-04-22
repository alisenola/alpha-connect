package krakenf

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
	"gitlab.com/alphaticks/xchanger/exchanges/krakenf"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
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
		client:      nil,
		rateLimit:   nil,
		queryRunner: nil,
		logger:      nil,
	}
}

func (state *Executor) Receive(context actor.Context) {
	extypes.ReceiveExecutor(state, context)
}

func (state *Executor) GetLogger() *log.Logger {
	return state.logger
}

func (state *Executor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(state).String()))

	state.client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1024,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewHTTPQuery(state.client)
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
	request, _, err := krakenf.GetInstruments()
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

	var instrumentResponse krakenf.InstrumentResponse
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
		underlying := strings.ToUpper(strings.Split(instrument.Symbol, "_")[1])
		baseSymbol := underlying[:3]
		baseCurrency := SymbolToAsset(baseSymbol)
		if baseCurrency == nil {
			continue
		}

		quoteSymbol := underlying[3:]
		quoteCurrency := SymbolToAsset(quoteSymbol)
		if quoteCurrency == nil {
			continue
		}
		security := models.Security{}
		security.Symbol = instrument.Symbol
		security.Underlying = baseCurrency
		security.QuoteCurrency = quoteCurrency
		if instrument.Tradeable {
			security.Status = models.InstrumentStatus_Trading
		} else {
			security.Status = models.InstrumentStatus_Disabled
		}
		security.Exchange = constants.KRAKENF
		splits := strings.Split(instrument.Symbol, "_")
		switch splits[0] {
		case "pv":
			// Perpetual vanilla
			security.SecurityType = enum.SecurityType_CRYPTO_PERP
			security.IsInverse = false
			// TODO check multiplier
			security.Multiplier = wrapperspb.Double(1)
		case "pi":
			// Perpetual inverse
			security.SecurityType = enum.SecurityType_CRYPTO_PERP
			security.IsInverse = true
			// TODO check multiplier
			security.Multiplier = wrapperspb.Double(1)
		case "pf":
			// Perpetual flexible
			security.SecurityType = enum.SecurityType_CRYPTO_PERP
			security.IsInverse = false
			// TODO check multiplier
			security.Multiplier = wrapperspb.Double(1)
		case "fv":
			// Future vanilla
			security.SecurityType = enum.SecurityType_CRYPTO_FUT
			security.IsInverse = false
			security.Multiplier = wrapperspb.Double(1)
			year := time.Now().Format("2006")
			date, err := time.Parse("20060102", year+splits[2][2:])
			if err != nil {
				continue
			}
			security.MaturityDate = timestamppb.New(date)
		case "fi":
			// Future inverse
			security.SecurityType = enum.SecurityType_CRYPTO_FUT
			security.IsInverse = true
			security.Multiplier = wrapperspb.Double(1)
			year := time.Now().Format("2006")
			date, err := time.Parse("20060102", year+splits[2][2:])
			if err != nil {
				continue
			}
			security.MaturityDate = timestamppb.New(date)
		}
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinPriceIncrement = &wrapperspb.DoubleValue{Value: instrument.TickSize}
		security.RoundLot = &wrapperspb.DoubleValue{Value: float64(instrument.ContractSize)}
		securities = append(securities, &security)
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
