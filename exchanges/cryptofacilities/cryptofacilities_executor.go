package cryptofacilities

import (
	"encoding/json"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alphac/enum"
	_interface "gitlab.com/alphaticks/alphac/exchanges/interface"
	"gitlab.com/alphaticks/alphac/jobs"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/cryptofacilities"
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

	// TODO rate limitting
	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) OnSecurityListRequest(context actor.Context) error {
	msg := context.Message().(*messages.SecurityListRequest)
	request, _, err := cryptofacilities.GetInstruments()
	if err != nil {
		return err
	}

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

		var instrumentResponse cryptofacilities.InstrumentResponse
		err = json.Unmarshal(queryResponse.Response, &instrumentResponse)
		if err != nil {
			err = fmt.Errorf("error decoding query response: %v", err)
			context.Respond(&messages.SecurityList{
				RequestID:  msg.RequestID,
				ResponseID: uint64(time.Now().UnixNano()),
				Error:      err.Error(),
				Securities: nil})
			return
		}
		if instrumentResponse.Result != "success" {
			err = fmt.Errorf("error getting instruments: %s", instrumentResponse.Error)
			context.Respond(&messages.SecurityList{
				RequestID:  msg.RequestID,
				ResponseID: uint64(time.Now().UnixNano()),
				Error:      err.Error(),
				Securities: nil})
			return
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
			baseCurrency, ok := constants.SYMBOL_TO_ASSET[baseSymbol]
			if !ok {
				continue
			}

			quoteSymbol := underlying[3:]
			if sym, ok := cryptofacilities.CRYPTOFACILITIES_SYMBOL_TO_GLOBAL_SYMBOL[quoteSymbol]; ok {
				quoteSymbol = sym
			}
			quoteCurrency, ok := constants.SYMBOL_TO_ASSET[quoteSymbol]
			if !ok {
				continue
			}
			security := models.Security{}
			security.Symbol = instrument.Symbol
			security.Underlying = &baseCurrency
			security.QuoteCurrency = &quoteCurrency
			security.Enabled = instrument.Tradeable
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
			security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name)
			security.MinPriceIncrement = instrument.TickSize
			security.RoundLot = float64(instrument.ContractSize)
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
