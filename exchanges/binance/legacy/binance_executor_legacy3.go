package legacy

/*
import (
	"encoding/json"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"github.com/quickfixgo/enum"
	fix50sl "github.com/quickfixgo/fix50/securitylist"
	"github.com/shopspring/decimal"
	"gitlab.com/alphaticks/alpha-connect/utils"
	//fix50nso "github.com/quickfixgo/quickfix/fix50/newordersingle"
	fix50slr "github.com/quickfixgo/fix50/securitylistrequest"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/jobs"
	models "gitlab.com/alphaticks/alpha-connect/messages/exchanges"
	"gitlab.com/alphaticks/alpha-connect/messages/executor"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/binance"
	"net/http"
	"reflect"
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

// The global rate limit is per IP and the orderRateLimit is per
// account.

type Executor struct {
	client               *http.Client
	secondOrderRateLimit *exchanges.RateLimit
	dayOrderRateLimit    *exchanges.RateLimit
	globalRateLimit      *exchanges.RateLimit
	queryRunner          *actor.PID
	exchangeInfo         *binance.ExchangeInfo
	logger               *log.Logger
}

func NewExecutor() actor.Actor {
	return &Executor{
		client:               nil,
		secondOrderRateLimit: nil,
		dayOrderRateLimit:    nil,
		globalRateLimit:      nil,
		queryRunner:          nil,
		logger:               nil,
	}
}

func (state *Executor) Receive(context actor.Context) {
	_interface.ExecutorReceive(state, context)
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

	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewAPIQuery(state.client)
	})
	state.queryRunner = context.Spawn(props)

	request, weight, err := binance.GetExchangeInfo()
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)
	res, err := future.Result()
	if err != nil {
		return err
	}
	queryResponse := res.(*jobs.PerformQueryResponse)
	if queryResponse.StatusCode != 200 {
		return fmt.Errorf("error getting exchange info: status code %d", queryResponse.StatusCode)
	}

	var exchangeInfo binance.ExchangeInfo
	err = json.Unmarshal(queryResponse.Response, &exchangeInfo)
	if err != nil {
		return fmt.Errorf("error decoding query response: %v", err)
	}
	if exchangeInfo.Code != 0 {
		return fmt.Errorf("error getting exchange info: %s", exchangeInfo.Msg)
	}

	state.exchangeInfo = &exchangeInfo
	// Initialize rate limit
	for _, rateLimit := range exchangeInfo.RateLimits {
		if rateLimit.RateLimitType == "ORDERS" {
			if rateLimit.Interval == "SECOND" {
				state.secondOrderRateLimit = exchanges.NewRateLimit(rateLimit.Limit, time.Second)
			} else if rateLimit.Interval == "DAY" {
				state.dayOrderRateLimit = exchanges.NewRateLimit(rateLimit.Limit, time.Hour*24)
			}
		} else if rateLimit.RateLimitType == "REQUEST_WEIGHT" {
			state.globalRateLimit = exchanges.NewRateLimit(rateLimit.Limit, time.Minute)
		}
	}
	if state.secondOrderRateLimit == nil || state.dayOrderRateLimit == nil {
		return fmt.Errorf("unable to set second or day rate limit")
	}

	// Update rate limit with weight for the current exchange info fetch
	state.globalRateLimit.Request(weight)
	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) updateExchangeInfos(context actor.Context) error {
	// Get http request and the expected response
	request, weight, err := binance.GetExchangeInfo()
	if err != nil {
		return err
	}

	if state.globalRateLimit.IsRateLimited() {
		time.Sleep(state.globalRateLimit.DurationBeforeNextRequest(weight))
		return nil
	}

	// Launch an APIQuery actor with the given request and target
	state.globalRateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res types{}, err error) {

		if err != nil {
			// TODO LOG
			return
		}

		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			// TODO LOG
			return
		}

		var exchangeInfo binance.ExchangeInfo
		err = json.Unmarshal(queryResponse.Response, &exchangeInfo)
		if err != nil {
			// TODO LOG
			return
		}

		if exchangeInfo.Code != 0 {
			// TODO LOG
			return
		}

		state.exchangeInfo = &exchangeInfo
	})

	return nil
}

func (state *Executor) GetOrderBookL2Request(context actor.Context) error {
	var snapshot *models.OBL2Snapshot
	msg := context.Message().(*executor.GetOrderBookL2Request)

	symbol := msg.Instrument.Format(binance.SymbolFormat)
	// Get http request and the expected response
	request, weight, err := binance.GetOrderBook(symbol, 1000)
	if err != nil {
		return err
	}

	if state.globalRateLimit.IsRateLimited() {
		time.Sleep(state.globalRateLimit.DurationBeforeNextRequest(weight))
		return nil
	}

	state.globalRateLimit.Request(weight)
	future := context.RequestFuture(state.queryRunner, &jobs.PerformQueryRequest{Request: request}, 10*time.Second)

	context.AwaitFuture(future, func(res types{}, err error) {
		if err != nil {
			context.Respond(&executor.GetOrderBookL2Response{
				RequestID: msg.RequestID,
				Error:     err,
				Snapshot:  nil})
			return
		}
		queryResponse := res.(*jobs.PerformQueryResponse)
		if queryResponse.StatusCode != 200 {
			if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
				err := fmt.Errorf(
					"http client error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				context.Respond(&executor.GetOrderBookL2Response{
					RequestID: msg.RequestID,
					Error:     err,
					Snapshot:  nil})
			} else if queryResponse.StatusCode >= 500 {
				err := fmt.Errorf(
					"http server error: %d %s",
					queryResponse.StatusCode,
					string(queryResponse.Response))
				context.Respond(&executor.GetOrderBookL2Response{
					RequestID: msg.RequestID,
					Error:     err,
					Snapshot:  nil})
			}
			return
		}
		var obData binance.OrderBookData
		err = json.Unmarshal(queryResponse.Response, &obData)
		if err != nil {
			err = fmt.Errorf("error decoding query response: %v", err)
			context.Respond(&executor.GetOrderBookL2Response{
				RequestID: msg.RequestID,
				Error:     err,
				Snapshot:  nil})
			return
		}
		if obData.Code != 0 {
			err = fmt.Errorf("error getting orderbook: %d %s", obData.Code, obData.Msg)
			context.Respond(&executor.GetOrderBookL2Response{
				RequestID: msg.RequestID,
				Error:     err,
				Snapshot:  nil})
			return
		}

		bids, asks, err := obData.ToRawBidAsk(msg.Instrument.TickPrecision, msg.Instrument.LotPrecision)
		if err != nil {
			err = fmt.Errorf("error converting orderbook: %v", err)
			context.Respond(&executor.GetOrderBookL2Response{
				RequestID: msg.RequestID,
				Error:     err,
				Snapshot:  nil})
			return
		}
		snapshot = &models.OBL2Snapshot{
			Instrument: msg.Instrument,
			Bids:       bids,
			Asks:       asks,
			Timestamp:  &types.Timestamp{Seconds: 0, Nanos: 0},
			ID:         obData.LastUpdateID,
		}
		context.Respond(&executor.GetOrderBookL2Response{
			RequestID: msg.RequestID,
			Error:     nil,
			Snapshot:  snapshot})
	})

	return nil
}

func (state *Executor) OnFIX50NewOrderSingle(context actor.Context) error {
	//msg := context.Message().(*fix50nso.NewOrderSingle)
	return nil
}

func (state *Executor) OnFIX50SecurityListRequest(context actor.Context) error {
	req := context.Message().(*fix50slr.SecurityListRequest)
	reqTyp, ferr := req.GetSecurityListRequestType()
	if ferr != nil {
		return fmt.Errorf("error getting SecurityListRequestType field: %v", ferr)
	}
	if reqTyp != enum.SecurityListRequestType_ALL_SECURITIES {
		return fmt.Errorf("SecurityListRequestType not supported")
	}
	reqID, ferr := req.GetSecurityReqID()
	if ferr != nil {
		return fmt.Errorf("error getting SecurityListRequestID field: %v", ferr)
	}

	response := fix50sl.New()
	response.SetSecurityReqID(reqID)
	response.SetSecurityResponseID(reqID)

	if state.exchangeInfo == nil {
		response.SetSecurityRequestResult(enum.SecurityRequestResult_INSTRUMENT_DATA_TEMPORARILY_UNAVAILABLE)
		context.Respond(&response)
		return nil
	}

	securities := fix50sl.NewNoRelatedSymRepeatingGroup()
	for _, symbol := range state.exchangeInfo.Symbols {
		// TODO logging
		baseCurrency, ok := constants.SYMBOL_TO_ASSET[symbol.BaseAsset]
		if !ok {
			continue
		}
		quoteCurrency, ok := constants.SYMBOL_TO_ASSET[symbol.QuoteAsset]
		if !ok {
			continue
		}
		security := securities.Add()
		if symbol.Status == "TRADING" {
			security.SetSecurityStatus(enum.SecurityStatus_ACTIVE)
		} else {
			security.SetSecurityStatus(enum.SecurityStatus_INACTIVE)
		}

		var tickPrecision int
		var lotPrecision int
		for _, filter := range symbol.Filters {
			if filter.FilterType == "PRICE_FILTER" {
				dec := decimal.NewFromFloat(filter.TickSize)
				tickPrecision = int(1. / filter.TickSize)
				security.SetMinPriceIncrement(dec, dec.Exponent())
			} else if filter.FilterType == "LOT_SIZE" {
				dec := decimal.NewFromFloat(filter.StepSize)
				lotPrecision = int(1. / filter.StepSize)
				security.SetRoundLot(dec, dec.Exponent())
			}
		}

		secID := fmt.Sprintf("%d", utils.SecurityID("SPOT", baseCurrency.Symbol, quoteCurrency.Symbol, constants.BINANCE, tickPrecision, lotPrecision))
		security.SetSymbol(symbol.Symbol)
		security.SetSecurityExchange(constants.BINANCE)
		security.SetCurrency(quoteCurrency.Symbol)
		security.SetProduct(enum.Product_CURRENCY)
		security.SetSecurityID(secID)
		security.SetSecurityType(enum.SecurityType)

	}

	return nil
}

func (state *Executor) OnFIX50SecurityDefinitionRequest(context actor.Context) error {
	return nil
}
*/
