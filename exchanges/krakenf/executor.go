package krakenf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

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
	xutils "gitlab.com/alphaticks/xchanger/utils"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
	req := context.Message().(*messages.HistoricalLiquidationsRequest)
	sender := context.Sender()
	response := &messages.HistoricalLiquidationsResponse{
		RequestID:       req.RequestID,
		Success:         false,
		RejectionReason: messages.RejectionReason_UnsupportedRequest,
	}
	go func() {
		symbol := ""
		if req.Instrument != nil {
			if req.Instrument.Symbol != nil {
				symbol = req.Instrument.Symbol.Value
			} else if req.Instrument.SecurityID != nil {
				sec := state.IDToSecurity(req.Instrument.SecurityID.Value)
				if sec == nil {
					response.RejectionReason = messages.RejectionReason_UnknownSecurityID
					context.Send(sender, response)
					return
				}
				symbol = sec.Symbol
			}
		} else {
			response.RejectionReason = messages.RejectionReason_UnknownSecurityID
			context.Send(sender, response)
			return
		}
		params := krakenf.NewFundingRateRequest(symbol)

		var request *http.Request
		var err error
		request, _, err = krakenf.GetHistoricalFundingRates(params)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
			return
		}

		state.rateLimit.Request(1)
		var data krakenf.HistoricalFundingRatesResponse
		if err := xutils.PerformJSONRequest(state.client, request, &data); err != nil {
			state.logger.Warn("error getting historicalFundingRates", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.Error != "" {
			state.logger.Warn(fmt.Sprintf("error getting historicalFundingRates"), log.Error(fmt.Errorf("%s", data.Error)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}

		response.Success = true
		context.Send(sender, response)
	}()
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

func (state *Executor) OnNewOrderSingleRequest(context actor.Context) error {
	msg := context.Message().(*messages.NewOrderSingleRequest)
	context.Respond(&messages.NewOrderSingleResponse{
		RequestID:       msg.RequestID,
		Success:         false,
		RejectionReason: messages.RejectionReason_UnsupportedRequest,
	})
	sender := context.Sender()
	response := &messages.NewOrderSingleResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	ar := state.rateLimit
	if ar == nil {
		var newAr *exchanges.RateLimit
		ar = newAr
	}

	if ar.IsRateLimited() {
		response.RejectionReason = messages.RejectionReason_AccountRateLimitExceeded
		response.RateLimitDelay = durationpb.New(ar.DurationBeforeNextRequest(1))
		context.Send(sender, response)
		return nil
	}
	if ar.GetCapacity() > 0.9 && msg.Order.IsForceMaker() {
		response.RejectionReason = messages.RejectionReason_TakerOnly
		context.Send(sender, response)
		return nil
	}

	go func() {
		var tickPrecision, lotPrecision int
		sec, rej := state.InstrumentToSecurity(msg.Order.Instrument)
		if rej != nil {
			response.RejectionReason = *rej
			context.Send(sender, response)
			return
		}

		tickPrecision = int(math.Ceil(math.Log10(1. / sec.MinPriceIncrement.Value)))
		lotPrecision = int(math.Ceil(math.Log10(1. / sec.RoundLot.Value)))

		params, rej := buildPostOrderRequest(sec.Symbol, msg.Order, tickPrecision, lotPrecision)
		if rej != nil {
			response.RejectionReason = *rej
			context.Send(sender, response)
			return
		}

		var request *http.Request
		var err error
		request, _, err = krakenf.SendOrder(msg.Account.ApiCredentials, params)
		if err != nil {
			state.logger.Warn("error building request", log.Error(err))
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
			return
		}

		ar.Request(1)

		var data krakenf.SendOrderResponse
		start := time.Now()
		if err := xutils.PerformJSONRequest(state.client, request, &data); err != nil {
			state.logger.Warn(fmt.Sprintf("error posting order for %s", msg.Account.Name), log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.Error != "" {
			state.logger.Warn(fmt.Sprintf("error posting order for %s", msg.Account.Name), log.Error(fmt.Errorf("%s", data.Error)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}
		status := StatusToModel(data.SendStatus.Status)
		if status == nil {
			state.logger.Error(fmt.Sprintf("unknown status %s", data.SendStatus.Status))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}
		response.NetworkRtt = durationpb.New(time.Since(start))
		response.Success = true
		response.OrderStatus = *status
		response.OrderID = fmt.Sprintf("%s", data.SendStatus.OrderId)
		fmt.Println("NEW SUCCESS", response.OrderID, response.OrderStatus.String())
		context.Send(sender, response)
	}()

	return nil
}

func (state *Executor) OnOrderCancelRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderCancelRequest)
	sender := context.Sender()
	response := &messages.OrderCancelResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	go func() {
		symbol := ""
		if req.Instrument != nil {
			if req.Instrument.Symbol != nil {
				symbol = req.Instrument.Symbol.Value
			} else if req.Instrument.SecurityID != nil {
				sec := state.IDToSecurity(req.Instrument.SecurityID.Value)
				if sec == nil {
					response.RejectionReason = messages.RejectionReason_UnknownSecurityID
					context.Send(sender, response)
					return
				}
				symbol = sec.Symbol
			}
		} else {
			response.RejectionReason = messages.RejectionReason_UnknownSecurityID
			context.Send(sender, response)
			return
		}
		params := krakenf.NewCancelOrderRequest(symbol)
		if req.OrderID != nil {
			orderIDInt, err := strconv.ParseInt(req.OrderID.Value, 10, 64)
			if err != nil {
				response.RejectionReason = messages.RejectionReason_UnknownOrder
				context.Send(sender, response)
				return
			}
			params.SetOrderId(fmt.Sprint(orderIDInt))
		} else if req.ClientOrderID != nil {
			params.SetCliOrdID(req.ClientOrderID.Value)
		} else {
			response.RejectionReason = messages.RejectionReason_UnknownOrder
			context.Send(sender, response)
			return
		}

		var request *http.Request
		var err error
		request, _, err = krakenf.CancelOrder(req.Account.ApiCredentials, params)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
			return
		}

		state.rateLimit.Request(1)
		var data krakenf.CancelOrderResponse
		start := time.Now()
		if err := xutils.PerformJSONRequest(state.client, request, &data); err != nil {
			state.logger.Warn("error cancelling order", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.Error != "" {
			state.logger.Warn(fmt.Sprintf("error posting order for %s", req.Account.Name), log.Error(fmt.Errorf("%s", data.Error)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}

		response.NetworkRtt = durationpb.New(time.Since(start))
		response.Success = true
		context.Send(sender, response)
	}()

	return nil
}

func (state *Executor) OnOrderMassCancelRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderMassCancelRequest)
	sender := context.Sender()
	response := &messages.OrderMassCancelResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	go func() {
		symbol := ""
		if req.Filter != nil {
			if req.Filter.Instrument != nil {
				if req.Filter.Instrument.Symbol != nil {
					sec := state.SymbolToSecurity(req.Filter.Instrument.Symbol.Value)
					if sec == nil {
						response.RejectionReason = messages.RejectionReason_UnknownSymbol
						context.Send(sender, response)
						return
					}
					symbol = req.Filter.Instrument.Symbol.Value
				} else if req.Filter.Instrument.SecurityID != nil {
					sec := state.IDToSecurity(req.Filter.Instrument.SecurityID.Value)
					if sec == nil {
						response.RejectionReason = messages.RejectionReason_UnknownSecurityID
						context.Send(sender, response)
						return
					}
					symbol = sec.Symbol
				}
			}
			if req.Filter.Side != nil || req.Filter.OrderStatus != nil {
				response.RejectionReason = messages.RejectionReason_UnsupportedFilter
				context.Send(sender, response)
				return
			}
		}
		if symbol == "" {
			response.RejectionReason = messages.RejectionReason_UnknownSymbol
			context.Send(sender, response)
			return
		}

		params := krakenf.NewCancelOrderRequest(symbol)

		request, weight, err := krakenf.CancelAllOrders(req.Account.ApiCredentials, params)
		if err != nil {
			response.RejectionReason = messages.RejectionReason_UnsupportedRequest
			context.Send(sender, response)
			return
		}

		state.rateLimit.Request(weight)
		var data krakenf.CancelAllOrdersResponse
		if err := xutils.PerformJSONRequest(state.client, request, &data); err != nil {
			state.logger.Warn("error cancelling orders", log.Error(err))
			response.RejectionReason = messages.RejectionReason_HTTPError
			context.Send(sender, response)
			return
		}
		if data.Error != "" {
			state.logger.Warn(fmt.Sprintf("error posting order for %s", req.Account.Name), log.Error(fmt.Errorf("%s", data.Error)))
			response.RejectionReason = messages.RejectionReason_ExchangeAPIError
			context.Send(sender, response)
			return
		}
		response.Success = true
		context.Send(sender, response)
	}()
	return nil
}
