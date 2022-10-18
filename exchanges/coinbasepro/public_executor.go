package coinbasepro

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
	"gitlab.com/alphaticks/xchanger/exchanges/coinbasepro"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io/ioutil"
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

// The role of a CoinbasePro Executor is to
// process api request
type PublicExecutor struct {
	extypes.BaseExecutor
	client       *http.Client
	rateLimit    *exchanges.RateLimit
	queryRunners []*actor.PID
	qrIdx        int
	logger       *log.Logger
}

func NewPublicExecutor() actor.Actor {
	return &PublicExecutor{
		client:       nil,
		rateLimit:    nil,
		queryRunners: nil,
		logger:       nil,
	}
}

func (state *PublicExecutor) Receive(context actor.Context) {
	extypes.ReceiveExecutor(state, context)
}

func (state *PublicExecutor) GetLogger() *log.Logger {
	return state.logger
}

func (state *PublicExecutor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(state).String()))

	state.rateLimit = exchanges.NewRateLimit(3, time.Second)
	state.client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1024,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	for i := 0; i < 4; i++ {
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
		state.queryRunners = append(state.queryRunners, context.Spawn(props))
	}

	return state.UpdateSecurityList(context)
}

func (state *PublicExecutor) Clean(context actor.Context) error {
	return nil
}

func (state *PublicExecutor) UpdateSecurityList(context actor.Context) error {
	request, weight, err := coinbasepro.GetProducts()
	if err != nil {
		return err
	}

	if state.rateLimit.IsRateLimited() {
		return fmt.Errorf("rate limit exceeded")
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
	var products []coinbasepro.Product
	err = json.Unmarshal(response, &products)
	if err != nil {
		err = fmt.Errorf(
			"error unmarshaling response: %v",
			err)
		return err
	}

	var securities []*models.Security
	for _, product := range products {
		if product.Status != "online" {
			continue
		}
		baseCurrency, ok := constants.GetAssetBySymbol(product.BaseCurrency)
		if !ok {
			continue
		}
		quoteCurrency, ok := constants.GetAssetBySymbol(product.QuoteCurrency)
		if !ok {
			continue
		}
		security := models.Security{}
		security.Symbol = product.ID
		security.Underlying = baseCurrency
		security.QuoteCurrency = quoteCurrency
		security.Status = models.InstrumentStatus_Trading
		security.Exchange = constants.COINBASEPRO
		security.SecurityType = enum.SecurityType_CRYPTO_SPOT
		security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
		security.MinPriceIncrement = &wrapperspb.DoubleValue{Value: product.QuoteIncrement}
		security.RoundLot = &wrapperspb.DoubleValue{Value: 1. / 100000000.}
		securities = append(securities, &security)
	}

	fmt.Println("SYNC", len(securities))
	state.SyncSecurities(securities, nil)

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: securities})

	return nil
}

func (state *PublicExecutor) OnHistoricalLiquidationsRequest(context actor.Context) error {
	msg := context.Message().(*messages.HistoricalLiquidationsRequest)
	context.Respond(&messages.HistoricalLiquidationsResponse{
		RequestID:       msg.RequestID,
		Success:         false,
		RejectionReason: messages.RejectionReason_UnsupportedRequest,
	})
	return nil
}

func (state *PublicExecutor) OnMarketStatisticsRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketStatisticsResponse)
	context.Respond(&messages.MarketStatisticsResponse{
		RequestID:       msg.RequestID,
		Success:         false,
		RejectionReason: messages.RejectionReason_UnsupportedRequest,
	})
	return nil
}

func (state *PublicExecutor) OnMarketDataRequest(context actor.Context) error {
	msg := context.Message().(*messages.MarketDataRequest)
	response := &messages.MarketDataResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	if state.rateLimit.IsRateLimited() {
		response.RejectionReason = messages.RejectionReason_RateLimitExceeded
		context.Respond(response)
		return nil
	}
	if msg.Subscribe {
		response.RejectionReason = messages.RejectionReason_UnsupportedSubscription
		context.Respond(response)
		return nil
	}
	if msg.Instrument == nil {
		response.RejectionReason = messages.RejectionReason_MissingInstrument
		context.Respond(response)
		return nil
	}
	var symbol = ""
	if msg.Instrument.Symbol != nil {
		symbol = msg.Instrument.Symbol.Value
	} else if msg.Instrument.SecurityID != nil {
		sec := state.IDToSecurity(msg.Instrument.SecurityID.Value)
		symbol = sec.Symbol
	}
	if symbol == "" {
		response.RejectionReason = messages.RejectionReason_UnknownSymbol
		context.Respond(response)
		return nil
	}

	if msg.Aggregation == models.OrderBookAggregation_L2 {
		var snapshot *models.OBL2Snapshot
		// Get http request and the expected response
		request, weight, err := coinbasepro.GetProductOrderBook(symbol, coinbasepro.L2ORDERBOOK)
		if err != nil {
			return err
		}

		state.rateLimit.Request(weight)
		queryRunner := state.queryRunners[state.qrIdx]
		state.qrIdx = (state.qrIdx + 1) % len(state.queryRunners)
		future := context.RequestFuture(queryRunner, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

		context.ReenterAfter(future, func(res interface{}, err error) {
			if err != nil {
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
				return
			}
			queryResponse := res.(*jobs.PerformQueryResponse)
			if queryResponse.StatusCode != 200 {
				if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
					err := fmt.Errorf(
						"http client error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
					state.logger.Info("http client error", log.Error(err))
					response.RejectionReason = messages.RejectionReason_HTTPError
					context.Respond(response)
				} else if queryResponse.StatusCode >= 500 {
					err := fmt.Errorf(
						"http server error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
					state.logger.Info("http client error", log.Error(err))
					response.RejectionReason = messages.RejectionReason_HTTPError
					context.Respond(response)
				}
				return
			}
			var obData coinbasepro.OrderBookL2
			err = json.Unmarshal(queryResponse.Response, &obData)
			if err != nil {
				state.logger.Info("error decoding query response", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
				return
			}

			bids, asks := obData.ToBidAsk()
			snapshot = &models.OBL2Snapshot{
				Bids:      bids,
				Asks:      asks,
				Timestamp: utils.MilliToTimestamp(0),
			}
			response.SnapshotL2 = snapshot
			response.SeqNum = obData.Sequence
			response.Success = true
			context.Respond(response)
		})

		return nil
	} else {
		var snapshot *models.OBL3Snapshot

		// Get http request and the expected response
		request, weight, err := coinbasepro.GetProductOrderBook(symbol, coinbasepro.L3ORDERBOOK)
		if err != nil {
			return err
		}

		state.rateLimit.Request(weight)
		queryRunner := state.queryRunners[state.qrIdx]
		state.qrIdx = (state.qrIdx + 1) % len(state.queryRunners)
		future := context.RequestFuture(queryRunner, &jobs.PerformHTTPQueryRequest{Request: request}, 10*time.Second)

		context.ReenterAfter(future, func(res interface{}, err error) {
			if err != nil {
				state.logger.Info("http client error", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
				return
			}
			queryResponse := res.(*jobs.PerformQueryResponse)
			if queryResponse.StatusCode != 200 {
				if queryResponse.StatusCode >= 400 && queryResponse.StatusCode < 500 {
					err := fmt.Errorf(
						"http client error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
					state.logger.Info("http client error", log.Error(err))
					response.RejectionReason = messages.RejectionReason_HTTPError
					context.Respond(response)
				} else if queryResponse.StatusCode >= 500 {
					err := fmt.Errorf(
						"http server error: %d %s",
						queryResponse.StatusCode,
						string(queryResponse.Response))
					state.logger.Info("http client error", log.Error(err))
					response.RejectionReason = messages.RejectionReason_HTTPError
					context.Respond(response)
				}
				return
			}

			var obData coinbasepro.OrderBookL3
			err = json.Unmarshal(queryResponse.Response, &obData)
			if err != nil {
				state.logger.Info("error decoding query response", log.Error(err))
				response.RejectionReason = messages.RejectionReason_HTTPError
				context.Respond(response)
				return
			}

			bids, asks := obData.ToBidAsk()
			snapshot = &models.OBL3Snapshot{
				Bids:      bids,
				Asks:      asks,
				Timestamp: nil,
			}
			response.SnapshotL3 = snapshot
			response.SeqNum = obData.Sequence
			response.Success = true
			context.Respond(response)
		})

		return nil
	}
}
