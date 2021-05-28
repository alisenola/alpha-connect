package bybitl

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
	"gitlab.com/alphaticks/xchanger/exchanges/bybitl"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"
	"unicode"
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
	request, weight, err := bybitl.GetSymbols()
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
		ReturnCode    int             `json:"ret_code"`
		ReturnMessage string          `json:"ret_msg"`
		ExitCode      string          `json:"ext_code"`
		ExitInfo      string          `json:"ext_info"`
		Result        []bybitl.Symbol `json:"result"`
	}
	var data Res
	err = json.Unmarshal(response, &data)
	if err != nil {
		err = fmt.Errorf(
			"error unmarshaling response: %v",
			err)
		return err
	}
	if data.ReturnCode != 0 {
		err = fmt.Errorf(
			"got wrong retrun code: %s",
			data.ReturnMessage)
		return err
	}

	var securities []*models.Security
	for _, symbol := range data.Result {
		// Only linear
		if symbol.QuoteCurrency == "USD" {
			continue
		}
		baseCurrency, ok := constants.GetAssetBySymbol(symbol.BaseCurrency)
		if !ok {
			//state.logger.Info(fmt.Sprintf("unknown currency %s", baseStr))
			continue
		}
		quoteCurrency, ok := constants.GetAssetBySymbol(symbol.QuoteCurrency)
		if !ok {
			//state.logger.Info(fmt.Sprintf("unknown currency %s", baseStr))
			continue
		}
		security := models.Security{}
		security.Symbol = symbol.Name
		security.Underlying = baseCurrency
		security.QuoteCurrency = quoteCurrency
		security.Status = models.Trading
		security.Exchange = &constants.BYBITL
		day := symbol.Alias[len(symbol.Alias)-2:]
		if unicode.IsNumber(rune(day[0])) && unicode.IsNumber(rune(day[1])) {
			security.SecurityType = enum.SecurityType_CRYPTO_FUT

			month := symbol.Alias[len(symbol.Alias)-4 : len(symbol.Alias)-2]
			monthInt, err := strconv.ParseInt(month, 10, 64)
			if err != nil {
				state.logger.Info(fmt.Sprintf("error parsing month %s", month))
				continue
			}
			dayInt, err := strconv.ParseInt(day, 10, 64)
			if err != nil {
				state.logger.Info(fmt.Sprintf("error parsing day: %d", dayInt))
				continue
			}
			year := time.Now().Year()
			if int(monthInt) < int(time.Now().Month()) {
				year += 1
			}
			date := time.Date(year, time.Month(monthInt), int(dayInt), 0, 0, 0, 0, time.UTC)
			ts, err := types.TimestampProto(date)
			if err != nil {
				state.logger.Info(fmt.Sprintf("error converting date"))
				continue
			}
			security.MaturityDate = ts
		} else {
			security.SecurityType = enum.SecurityType_CRYPTO_PERP
		}
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinPriceIncrement = &types.DoubleValue{Value: symbol.PriceFilter.TickSize}
		security.RoundLot = &types.DoubleValue{Value: symbol.LotSizeFilter.QuantityStep}
		security.IsInverse = true
		security.MakerFee = &types.DoubleValue{Value: symbol.MakerFee}
		security.TakerFee = &types.DoubleValue{Value: symbol.TakerFee}

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
