package bitmex

import (
	"encoding/json"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/enum"
	"gitlab.com/alphaticks/alphac/exchanges/interface"
	"gitlab.com/alphaticks/alphac/jobs"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/bitmex"
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
type Executor struct {
	client      *http.Client
	rateLimit   *exchanges.RateLimit
	queryRunner *actor.PID
	logger      *log.Logger
}

func NewExecutor() actor.Actor {
	return &Executor{nil, nil, nil, nil}
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

	// TODO rate limitting
	// Launch an APIQuery actor with the given request and target
	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewAPIQuery(state.client)
	})
	state.queryRunner = context.Spawn(props)
	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

/*
// Fetches all the instruments on the website
func (b *BitmexAPI) Initialize() error {


	return nil
}
*/

func (state *Executor) OnSecurityListRequest(context actor.Context) error {
	msg := context.Message().(*messages.SecurityListRequest)
	request, _, err := bitmex.GetActiveInstruments()
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

		var activeInstruments []bitmex.Instrument
		err = json.Unmarshal(queryResponse.Response, &activeInstruments)
		if err != nil {
			err = fmt.Errorf("error decoding query response: %v", err)
			context.Respond(&messages.SecurityList{
				RequestID:  msg.RequestID,
				ResponseID: uint64(time.Now().UnixNano()),
				Error:      err.Error(),
				Securities: nil})
			return
		}
		var securities []*models.Security
		for _, activeInstrument := range activeInstruments {
			if activeInstrument.State != "Open" {
				continue
			}
			switch activeInstrument.Typ {
			/*
				case "FFCCSX":
					instr.Type = exchanges.FUTURE
			*/
			case "FFWCSX":
				symbolStr := strings.ToUpper(activeInstrument.Underlying)
				if sym, ok := bitmex.BITMEX_SYMBOL_TO_GLOBAL_SYMBOL[symbolStr]; ok {
					symbolStr = sym
				}
				baseCurrency, ok := constants.SYMBOL_TO_ASSET[symbolStr]
				if !ok {
					continue
				}
				symbolStr = strings.ToUpper(activeInstrument.QuoteCurrency)
				if sym, ok := bitmex.BITMEX_SYMBOL_TO_GLOBAL_SYMBOL[symbolStr]; ok {
					symbolStr = sym
				}
				quoteCurrency, ok := constants.SYMBOL_TO_ASSET[symbolStr]
				if !ok {
					continue
				}

				security := models.Security{}
				security.Symbol = activeInstrument.Symbol
				security.Underlying = &baseCurrency
				security.QuoteCurrency = &quoteCurrency
				security.Enabled = activeInstrument.State == "Open"
				security.Exchange = &constants.BITMEX
				security.SecurityType = enum.SecurityType_CRYPTO_PERP
				security.MinPriceIncrement = activeInstrument.TickSize
				security.RoundLot = float64(activeInstrument.LotSize)
				security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name)

				securities = append(securities, &security)

			default:
				// case "OCECCS":
				//instr.Type = exchanges.CALL_OPTION
				// case "OPECCS":
				//instr.Type = exchanges.PUT_OPTION
				// Non-supported instrument, passing..
				continue
			}

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
