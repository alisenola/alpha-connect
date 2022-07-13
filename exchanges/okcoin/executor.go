package okcoin

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
	"gitlab.com/alphaticks/xchanger/exchanges/okcoin"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"net/http"
	"reflect"
	"time"
)

type QueryRunner struct {
	pid       *actor.PID
	rateLimit *exchanges.RateLimit
}

type Executor struct {
	extypes.BaseExecutor
	securities  []*models.Security
	queryRunner *QueryRunner
	logger      *log.Logger
}

func NewExecutor() actor.Actor {
	return &Executor{
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

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1024,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewHTTPQuery(client)
	})
	state.queryRunner = &QueryRunner{
		pid:       context.Spawn(props),
		rateLimit: exchanges.NewRateLimit(6, time.Second),
	}
	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	if state.queryRunner.rateLimit.IsRateLimited() {
		return fmt.Errorf("rate limited")
	}
	request, weight, err := okcoin.GetSpotInstruments()
	if err != nil {
		return err
	}
	qr := state.queryRunner

	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

	res, err := future.Result()
	if err != nil {
		return fmt.Errorf("http client error: %v", err)
	}
	resp := res.(*jobs.PerformQueryResponse)

	qr.rateLimit.Request(weight)

	if resp.StatusCode != 200 {
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			err := fmt.Errorf(
				"http client error: %d %s",
				resp.StatusCode,
				string(resp.Response))
			return err
		} else if resp.StatusCode >= 500 {
			err := fmt.Errorf(
				"http server error: %d %s",
				resp.StatusCode,
				string(resp.Response))
			return err
		} else {
			err := fmt.Errorf("%d %s",
				resp.StatusCode,
				string(resp.Response))
			return err
		}
	}

	var kResponse []okcoin.SpotInstrument
	err = json.Unmarshal(resp.Response, &kResponse)
	if err != nil {
		err = fmt.Errorf("error decoding query response: %v", err)
		return err
	}

	var securities []*models.Security
	for _, pair := range kResponse {
		baseCurrency, ok := constants.GetAssetBySymbol(pair.BaseCurrency)
		if !ok {
			//state.logger.Info("unknown symbol " + pair.BaseCurrency + " for instrument " + pair.InstrumentID)
			continue
		}
		quoteCurrency, ok := constants.GetAssetBySymbol(pair.QuoteCurrency)
		if !ok {
			//state.logger.Info("unknown symbol " + pair.QuoteCurrency + " for instrument " + pair.InstrumentID)
			continue
		}

		security := models.Security{}
		security.Symbol = pair.InstrumentID
		security.Underlying = baseCurrency
		security.QuoteCurrency = quoteCurrency
		security.Status = models.InstrumentStatus_Trading
		security.Exchange = constants.OKCOIN
		security.SecurityType = enum.SecurityType_CRYPTO_SPOT
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinPriceIncrement = &wrapperspb.DoubleValue{Value: pair.TickSize}
		security.RoundLot = &wrapperspb.DoubleValue{Value: pair.SizeIncrement}
		securities = append(securities, &security)
	}
	state.SyncSecurities(securities, nil)

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: securities})

	return nil
}
