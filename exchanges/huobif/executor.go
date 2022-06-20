package huobif

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
	"gitlab.com/alphaticks/xchanger/exchanges/huobif"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		return jobs.NewHTTPQuery(state.client)
	})
	state.queryRunner = context.Spawn(props)

	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {
	request, weight, err := huobif.GetContractInfo()
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

	type Data struct {
		Status string                `json:"status"`
		Data   []huobif.ContractInfo `json:"data"`
	}
	var data Data
	err = json.Unmarshal(response, &data)
	if err != nil {
		err = fmt.Errorf(
			"error unmarshaling response: %v",
			err)
		return err
	}
	if data.Status != "ok" {
		err = fmt.Errorf(
			"got wrong status: %s",
			data.Status)
		return err
	}

	var securities []*models.Security
	for _, symbol := range data.Data {
		baseCurrency, ok := constants.GetAssetBySymbol(symbol.Symbol)
		if !ok {
			//state.logger.Info(fmt.Sprintf("unknown currency %s", baseStr))
			continue
		}
		/*
			0: Delisting,1: Listing,2: Pending Listing,3: Suspension,4: Suspending of Listing,5: In Settlement,6: Delivering,7: Settlement Completed,8: Delivered,9: Suspending of Trade
		*/
		quoteCurrency := constants.DOLLAR
		security := models.Security{}
		security.Symbol = symbol.ContractCode
		security.Underlying = baseCurrency
		security.QuoteCurrency = quoteCurrency
		switch symbol.ContractStatus {
		case 0:
			security.Status = models.InstrumentStatus_Disabled
		case 1:
			security.Status = models.InstrumentStatus_Trading
		case 2:
			security.Status = models.InstrumentStatus_PreTrading
		case 3:
			security.Status = models.InstrumentStatus_Break
		default:
			security.Status = models.InstrumentStatus_Disabled
		}
		security.Exchange = constants.HUOBIF
		security.SecurityType = enum.SecurityType_CRYPTO_FUT
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinPriceIncrement = &wrapperspb.DoubleValue{Value: symbol.PriceTick}
		security.RoundLot = &wrapperspb.DoubleValue{Value: 1}
		security.IsInverse = true
		security.Multiplier = &wrapperspb.DoubleValue{Value: symbol.ContractSize}

		date, err := time.Parse("20060102", symbol.DeliveryDate)
		if err != nil {
			return fmt.Errorf("error parsing delivery date %s: %v", symbol.DeliveryDate, err)
		}
		// Note: deliver at 8 AM GMT
		date = date.Add(8 * time.Hour)
		security.MaturityDate = timestamppb.New(date)
		if err != nil {
			return fmt.Errorf("error parsing delivery date %s: %v", symbol.DeliveryDate, err)
		}
		securities = append(securities, &security)
	}
	state.SyncSecurities(securities, nil)

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: securities})

	return nil
}
