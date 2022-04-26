package exchanges

import (
	goContext "context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/account"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/commands"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// The executor routes all the request to the underlying exchange executor & listeners
// He is the main part of the whole software..
type ExecutorConfig struct {
	Db                 *mongo.Database
	Exchanges          []*xchangerModels.Exchange
	Accounts           []*account.Account
	DialerPool         *xchangerUtils.DialerPool
	Registry           registry.PublicRegistryClient
	OpenseaCredentials *xchangerModels.APICredentials
	Strict             bool
}

type Executor struct {
	*ExecutorConfig
	accountManagers    map[string]*actor.PID
	executors          map[uint32]*actor.PID              // A map from exchange ID to executor
	securities         map[uint64]*models.Security        // A map from security ID to security
	marketAssets       map[uint64]*models.MarketableAsset // A map from MarketAsset ID to MarketAsset
	symbToSecs         map[uint32]map[string]*models.Security
	symbToMarketAssets map[uint32]map[uint64]*models.MarketableAsset
	instruments        map[uint64]*actor.PID // A map from security ID to market manager
	slSubscribers      map[uint64]*actor.PID // A map from request ID to security list subscribers
	malSubscribers     map[uint64]*actor.PID // A map from request ID to protocolAssets list subscribers
	logger             *log.Logger
	strict             bool
}

func NewExecutorProducer(cfg *ExecutorConfig) actor.Producer {
	return func() actor.Actor {
		return NewExecutor(cfg)
	}
}

func NewExecutor(cfg *ExecutorConfig) actor.Actor {
	return &Executor{
		ExecutorConfig: cfg,
	}
}

func (state *Executor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.logger.Error("error initializing", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor started")

	case *actor.Stopping:
		if err := state.Clean(context); err != nil {
			state.logger.Error("error stopping", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor stopping")

	case *actor.Stopped:
		state.logger.Info("actor stopped")

	case *actor.Restarting:
		if err := state.Clean(context); err != nil {
			state.logger.Error("error restarting", log.Error(err))
			// Attention, no panic in restarting or infinite loop
		}
		state.logger.Info("actor restarting")

	case *messages.AccountDataRequest:
		if err := state.OnAccountDataRequest(context); err != nil {
			state.logger.Error("error processing OnAccountDataRequest", log.Error(err))
			panic(err)
		}

	case *messages.MarketDataRequest:
		if err := state.OnMarketDataRequest(context); err != nil {
			state.logger.Error("error processing OnMarketDataRequest", log.Error(err))
			panic(err)
		}

	case *messages.UnipoolV3DataRequest:
		if err := state.OnUnipoolV3DataRequest(context); err != nil {
			state.logger.Error("error processing OnUnipoolV3DataRequest", log.Error(err))
			panic(err)
		}

	case *messages.HistoricalUnipoolV3DataRequest:
		if err := state.OnHistoricalUnipoolV3EventRequest(context); err != nil {
			state.logger.Error("error processing OnHistoricalUnipoolV3EventRequest", log.Error(err))
			panic(err)
		}

	case *messages.MarketStatisticsRequest:
		if err := state.OnMarketStatisticsRequest(context); err != nil {
			state.logger.Error("error processing OnMarketStatisticsRequest", log.Error(err))
			panic(err)
		}

	case *messages.HistoricalLiquidationsRequest:
		if err := state.OnHistoricalLiquidationsRequest(context); err != nil {
			state.logger.Error("error processing OnHistoricalLiquidationsRequest", log.Error(err))
			panic(err)
		}

	case *messages.HistoricalOpenInterestsRequest:
		if err := state.OnHistoricalOpenInterestsRequest(context); err != nil {
			state.logger.Error("error processing OnHistoricalOpenInterestsRequest", log.Error(err))
			panic(err)
		}

	case *messages.HistoricalFundingRatesRequest:
		if err := state.OnHistoricalFundingRatesRequest(context); err != nil {
			state.logger.Error("error processing OnHistoricalFundingRateRequest", log.Error(err))
			panic(err)
		}

	case *messages.HistoricalSalesRequest:
		if err := state.OnHistoricalSalesRequest(context); err != nil {
			state.logger.Error("error processing OnHistoricalSalesRequest", log.Error(err))
			panic(err)
		}

	case *messages.SecurityDefinitionRequest:
		if err := state.OnSecurityDefinitionRequest(context); err != nil {
			state.logger.Error("error processing OnSecurityDefinitionRequest", log.Error(err))
			panic(err)
		}

	case *messages.SecurityListRequest:
		if err := state.OnSecurityListRequest(context); err != nil {
			state.logger.Error("error processing OnSecurityListRequest", log.Error(err))
			panic(err)
		}

	case *messages.SecurityList:
		if err := state.OnSecurityList(context); err != nil {
			state.logger.Error("error processing OnSecurityList", log.Error(err))
			panic(err)
		}

	case *messages.MarketableAssetListRequest:
		if err := state.OnMarketableAssetListRequest(context); err != nil {
			state.logger.Error("error processing OnMarketableAssetListRequest", log.Error(err))
			panic(err)
		}

	case *messages.MarketableAssetList:
		if err := state.OnMarketableAssetList(context); err != nil {
			state.logger.Error("error processing OnMarketableAssetList", log.Error(err))
			panic(err)
		}
	case *messages.AccountMovementRequest:
		if err := state.OnAccountMovementRequest(context); err != nil {
			state.logger.Error("error processing OnAccountMovementRequest", log.Error(err))
			panic(err)
		}

	case *messages.TradeCaptureReportRequest:
		if err := state.OnTradeCaptureReportRequest(context); err != nil {
			state.logger.Error("error processing OnTradeCaptureReportRequest", log.Error(err))
			panic(err)
		}

	case *messages.PositionsRequest:
		if err := state.OnPositionsRequest(context); err != nil {
			state.logger.Error("error processing OnPositionRequest", log.Error(err))
			panic(err)
		}

	case *messages.BalancesRequest:
		if err := state.OnBalancesRequest(context); err != nil {
			state.logger.Error("error processing OnBalancesRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderStatusRequest:
		if err := state.OnOrderStatusRequest(context); err != nil {
			state.logger.Error("error processing OnOrderStatusRequest", log.Error(err))
			panic(err)
		}

	case *messages.NewOrderSingleRequest:
		if err := state.OnNewOrderSingleRequest(context); err != nil {
			state.logger.Error("error processing OnNewOrderSingle", log.Error(err))
			panic(err)
		}

	case *messages.NewOrderBulkRequest:
		if err := state.OnNewOrderBulkRequest(context); err != nil {
			state.logger.Error("error processing OnNewOrderBulk", log.Error(err))
			panic(err)
		}

	case *messages.OrderReplaceRequest:
		if err := state.OnOrderReplaceRequest(context); err != nil {
			state.logger.Error("error processing OrderReplaceRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderBulkReplaceRequest:
		if err := state.OnOrderBulkReplaceRequest(context); err != nil {
			state.logger.Error("error processing OrderBulkReplaceRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderCancelRequest:
		if err := state.OnOrderCancelRequest(context); err != nil {
			state.logger.Error("error processing OnOrderCancelRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderMassCancelRequest:
		if err := state.OnOrderMassCancelRequest(context); err != nil {
			state.logger.Error("error processing OnOrderMassCancelRequest", log.Error(err))
			panic(err)
		}

	case *commands.GetAccountRequest:
		if err := state.OnGetAccountRequest(context); err != nil {
			state.logger.Error("error processing OnListenAccountRequest", log.Error(err))
			panic(err)
		}

	case *actor.Terminated:
		if err := state.OnTerminated(context); err != nil {
			state.logger.Error("error processing OnTerminated", log.Error(err))
			panic(err)
		}
	}
}

func (state *Executor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	state.instruments = make(map[uint64]*actor.PID)
	state.slSubscribers = make(map[uint64]*actor.PID)
	state.accountManagers = make(map[string]*actor.PID)

	if state.DialerPool == nil {
		state.DialerPool = xchangerUtils.DefaultDialerPool
	}

	if state.Db != nil {
		unique := true
		mod := mongo.IndexModel{
			Keys: bson.M{
				"id": 1, // index in ascending order
			}, Options: &options.IndexOptions{Unique: &unique},
		}
		txs := state.Db.Collection("transactions")
		execs := state.Db.Collection("executions")
		if _, err := txs.Indexes().CreateOne(goContext.Background(), mod); err != nil {
			return fmt.Errorf("error creating index on transactions: %v", err)
		}
		if _, err := execs.Indexes().CreateOne(goContext.Background(), mod); err != nil {
			return fmt.Errorf("error creating index on executions: %v", err)
		}
	}

	// TODO add dialer pool test

	// Spawn all exchange executors
	state.executors = make(map[uint32]*actor.PID)
	for _, exch := range state.Exchanges {
		producer := NewExchangeExecutorProducer(exch, state.DialerPool, state.ExecutorConfig)
		if producer == nil {
			return fmt.Errorf("unknown exchange %s", exch.Name)
		}
		props := actor.PropsFromProducer(producer).WithSupervisor(actor.NewExponentialBackoffStrategy(100*time.Second, time.Second))

		state.executors[exch.ID], _ = context.SpawnNamed(props, exch.Name+"_executor")
	}

	// Request securities for each one of them
	var futures []*actor.Future
	request := &messages.SecurityListRequest{
		RequestID: 0,
		Subscribe: true,
	}
	for _, pid := range state.executors {
		fut := context.RequestFuture(pid, request, 20*time.Second)
		futures = append(futures, fut)
	}

	state.securities = make(map[uint64]*models.Security)
	state.symbToSecs = make(map[uint32]map[string]*models.Security)
	for _, fut := range futures {
		res, err := fut.Result()
		if err != nil {
			if state.strict {
				return fmt.Errorf("error fetching securities for one venue: %v", err)
			} else {
				state.logger.Error("error fetching securities for one venue", log.Error(err))
				continue
			}
		}
		response, ok := res.(*messages.SecurityList)
		if !ok {
			return fmt.Errorf("was expecting GetSecuritiesResponse, got %s", reflect.TypeOf(res).String())
		}
		if !response.Success {
			return errors.New(response.RejectionReason.String())
		}
		if len(response.Securities) > 0 {
			symbToSec := make(map[string]*models.Security)
			var exchID uint32
			for _, s := range response.Securities {
				if sec2, ok := state.securities[s.SecurityID]; ok {
					return fmt.Errorf("got two securities with the same ID: %s %s", sec2.Symbol, s.Symbol)
				}
				state.securities[s.SecurityID] = s
				symbToSec[s.Symbol] = s
				exchID = s.Exchange.ID
			}
			state.symbToSecs[exchID] = symbToSec
		}
	}

	//Request market assets for all of them
	var futs []*actor.Future
	req := messages.MarketableAssetListRequest{
		RequestID: 0,
		Subscribe: true,
	}
	for _, pid := range state.executors {
		f := context.RequestFuture(pid, &req, 10*time.Second)
		futs = append(futs, f)
	}

	state.marketAssets = make(map[uint64]*models.MarketableAsset)
	state.symbToMarketAssets = make(map[uint32]map[uint64]*models.MarketableAsset)
	for _, f := range futs {
		resp, err := f.Result()
		if err != nil {
			state.logger.Error("error fetching marketable assets for one venue", log.Error(err))
			continue
		}
		list, ok := resp.(*messages.MarketableAssetList)
		if !ok {
			return fmt.Errorf("was expecting MarketAssetList, got %s", reflect.TypeOf(resp).String())
		}
		if !list.Success {
			return errors.New(list.RejectionReason.String())
		}
		if len(list.MarketableAssets) > 0 {
			idToMarketAssets := make(map[uint64]*models.MarketableAsset)
			var marketID uint32
			for _, asset := range list.MarketableAssets {
				if asset2, ok := state.marketAssets[asset.ProtocolAsset.ProtocolAssetID]; ok {
					return fmt.Errorf("got two protocol assets with the same ID: %v and %v", asset.ProtocolAsset, asset2.ProtocolAsset)
				}
				state.marketAssets[asset.ProtocolAsset.ProtocolAssetID] = asset
				idToMarketAssets[asset.ProtocolAsset.ProtocolAssetID] = asset
				marketID = asset.Market.ID
			}
			state.symbToMarketAssets[marketID] = idToMarketAssets
		}
	}

	// Spawn all account listeners
	state.accountManagers = make(map[string]*actor.PID)
	for _, accnt := range state.Accounts {
		producer := NewAccountManagerProducer(accnt, state.Db, false)
		if producer == nil {
			return fmt.Errorf("unknown exchange %s", accnt.Exchange.Name)
		}
		props := actor.PropsFromProducer(producer).WithSupervisor(
			actor.NewExponentialBackoffStrategy(100*time.Second, time.Second))
		state.accountManagers[accnt.Name] = context.Spawn(props)
	}

	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) OnAccountDataRequest(context actor.Context) error {
	request := context.Message().(*messages.AccountDataRequest)
	if request.Account == nil {
		context.Respond(&messages.AccountDataResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: messages.UnknownAccount,
		})
		return nil
	}

	if pid, ok := state.accountManagers[request.Account.Name]; ok {
		context.Forward(pid)
	} else {
		accnt, err := NewAccount(request.Account)
		if err != nil {
			return fmt.Errorf("error creating account: %v", err)
		}
		producer := NewAccountManagerProducer(accnt, state.Db, false)
		if producer == nil {
			return fmt.Errorf("unknown exchange %s", accnt.Exchange.Name)
		}
		props := actor.PropsFromProducer(producer).WithSupervisor(
			actor.NewExponentialBackoffStrategy(100*time.Second, time.Second))
		pid := context.Spawn(props)
		state.accountManagers[request.Account.Name] = pid
		context.Forward(pid)
	}

	return nil
}

func (state *Executor) getSecurity(instr *models.Instrument) (*models.Security, *messages.RejectionReason) {
	if instr == nil {
		rej := messages.MissingInstrument
		return nil, &rej
	}
	if instr.SecurityID != nil {
		if sec, ok := state.securities[instr.SecurityID.Value]; ok {
			return sec, nil
		} else {
			rej := messages.UnknownSecurityID
			return nil, &rej
		}
	} else if instr.Symbol != nil {
		if instr.Exchange == nil {
			rej := messages.UnknownExchange
			return nil, &rej
		}
		symbolsToSecs, ok := state.symbToSecs[instr.Exchange.ID]
		if !ok {
			rej := messages.UnknownExchange
			return nil, &rej
		}
		sec, ok := symbolsToSecs[instr.Symbol.Value]
		if !ok {
			rej := messages.UnknownSymbol
			return nil, &rej
		}
		return sec, nil
	} else {
		rej := messages.MissingInstrument
		return nil, &rej
	}
}

func (state *Executor) getMarketableAssets(assetID uint32) ([]*models.MarketableAsset, *messages.RejectionReason) {
	var assets []*models.MarketableAsset
	for _, v := range state.marketAssets {
		if v.ProtocolAsset.Asset.ID == assetID {
			assets = append(assets, v)
		}
	}
	if len(assets) == 0 {
		rej := messages.UnknownAsset
		return nil, &rej
	}
	return assets, nil
}

func (state *Executor) OnMarketDataRequest(context actor.Context) error {
	request := context.Message().(*messages.MarketDataRequest)
	sec, rej := state.getSecurity(request.Instrument)
	if rej != nil {
		context.Respond(&messages.MarketDataResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: *rej,
		})
		return nil
	}
	if pid, ok := state.instruments[sec.SecurityID]; ok {
		context.Forward(pid)
	} else {
		props := actor.PropsFromProducer(NewDataManagerProducer(sec, state.DialerPool)).WithSupervisor(
			utils.NewExponentialBackoffStrategy(100*time.Second, time.Second, time.Second))
		pid := context.Spawn(props)
		state.instruments[sec.SecurityID] = pid
		context.Forward(pid)
	}

	return nil
}

func (state *Executor) OnUnipoolV3DataRequest(context actor.Context) error {
	request := context.Message().(*messages.UnipoolV3DataRequest)
	sec, rej := state.getSecurity(request.Instrument)
	if rej != nil {
		context.Respond(&messages.UnipoolV3DataResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: *rej,
		})
		return nil
	}
	if pid, ok := state.instruments[sec.SecurityID]; ok {
		context.Forward(pid)
	} else {
		props := actor.PropsFromProducer(NewDataManagerProducer(sec, state.DialerPool)).WithSupervisor(
			utils.NewExponentialBackoffStrategy(100*time.Second, time.Second, time.Second))
		pid := context.Spawn(props)
		state.instruments[sec.SecurityID] = pid
		context.Forward(pid)
	}

	return nil
}

func (state *Executor) OnHistoricalUnipoolV3EventRequest(context actor.Context) error {
	request := context.Message().(*messages.HistoricalUnipoolV3DataRequest)
	sec, rej := state.getSecurity(request.Instrument)
	if rej != nil {
		context.Respond(&messages.HistoricalUnipoolV3DataResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: *rej,
		})
		return nil
	}
	exchange, ok := state.executors[sec.Exchange.ID]
	if !ok {
		context.Respond(&messages.HistoricalUnipoolV3DataResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: messages.UnknownExchange,
		})
	}
	context.Forward(exchange)

	return nil
}

func (state *Executor) OnMarketStatisticsRequest(context actor.Context) error {
	request := context.Message().(*messages.MarketStatisticsRequest)
	sec, rej := state.getSecurity(request.Instrument)
	if rej != nil {
		context.Respond(&messages.MarketStatisticsResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: *rej,
		})
		return nil
	}
	exchange, ok := state.executors[sec.Exchange.ID]
	if !ok {
		context.Respond(&messages.MarketStatisticsResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: messages.UnknownExchange,
		})
		return nil
	}
	context.Forward(exchange)

	return nil
}

func (state *Executor) OnHistoricalLiquidationsRequest(context actor.Context) error {
	request := context.Message().(*messages.HistoricalLiquidationsRequest)
	sec, rej := state.getSecurity(request.Instrument)
	if rej != nil {
		context.Respond(&messages.HistoricalLiquidationsResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: *rej,
		})
		return nil
	}
	exchange, ok := state.executors[sec.Exchange.ID]
	if !ok {
		context.Respond(&messages.HistoricalLiquidationsResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: messages.UnknownExchange,
		})
		return nil
	}
	context.Forward(exchange)

	return nil
}

func (state *Executor) OnHistoricalOpenInterestsRequest(context actor.Context) error {
	request := context.Message().(*messages.HistoricalOpenInterestsRequest)
	sec, rej := state.getSecurity(request.Instrument)
	if rej != nil {
		context.Respond(&messages.HistoricalOpenInterestsResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: *rej,
		})
		return nil
	}
	exchange, ok := state.executors[sec.Exchange.ID]
	if !ok {
		context.Respond(&messages.HistoricalOpenInterestsResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: messages.UnknownExchange,
		})
		return nil
	}
	context.Forward(exchange)

	return nil
}

func (state *Executor) OnHistoricalFundingRatesRequest(context actor.Context) error {
	request := context.Message().(*messages.HistoricalFundingRatesRequest)
	sec, rej := state.getSecurity(request.Instrument)
	if rej != nil {
		context.Respond(&messages.HistoricalFundingRatesResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: *rej,
		})
		return nil
	}
	exchange, ok := state.executors[sec.Exchange.ID]
	if !ok {
		context.Respond(&messages.HistoricalFundingRatesResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: messages.UnknownExchange,
		})
		return nil
	}
	context.Forward(exchange)

	return nil
}

func (state *Executor) OnHistoricalSalesRequest(context actor.Context) error {
	request := context.Message().(*messages.HistoricalSalesRequest)
	assets, rej := state.getMarketableAssets(request.AssetID)
	if rej != nil {
		context.Respond(&messages.HistoricalSalesResponse{
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: *rej,
		})
		return nil
	}
	for _, asset := range assets {
		market, ok := state.executors[asset.Market.ID]
		if !ok {
			context.Respond(&messages.HistoricalSalesResponse{
				ResponseID:      uint64(time.Now().UnixNano()),
				Success:         false,
				RejectionReason: messages.UnknownExchange,
			})
			return nil
		}
		context.Forward(market)
	}
	return nil
}

func (state *Executor) OnSecurityDefinitionRequest(context actor.Context) error {
	request := context.Message().(*messages.SecurityDefinitionRequest)
	sec, rej := state.getSecurity(request.Instrument)
	if rej != nil {
		context.Respond(&messages.SecurityDefinitionResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: *rej,
		})
		return nil
	}

	context.Respond(&messages.SecurityDefinitionResponse{
		RequestID:  request.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Security:   sec,
		Success:    true,
	})

	return nil
}

func (state *Executor) OnSecurityListRequest(context actor.Context) error {
	request := context.Message().(*messages.SecurityListRequest)
	var securities []*models.Security
	for _, v := range state.securities {
		securities = append(securities, v)
	}
	response := &messages.SecurityList{
		RequestID:  request.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Securities: securities,
		Success:    true,
	}
	if request.Subscribe {
		context.Watch(request.Subscriber)
		state.slSubscribers[request.RequestID] = request.Subscriber
	}
	context.Respond(response)

	return nil
}

func (state *Executor) OnSecurityList(context actor.Context) error {
	securityList := context.Message().(*messages.SecurityList)
	// Do nothing
	if len(securityList.Securities) == 0 {
		return nil
	}
	exchangeID := securityList.Securities[0].Exchange.ID
	// It has to come from one exchange only, so delete all the known securities from that exchange
	for k, v := range state.securities {
		if v.Exchange.ID == exchangeID {
			delete(state.securities, k)
		}
	}
	delete(state.symbToSecs, exchangeID)

	// re-add them
	secs := make(map[string]*models.Security)
	for _, s := range securityList.Securities {
		state.securities[s.SecurityID] = s
		secs[s.Symbol] = s
	}
	state.symbToSecs[exchangeID] = secs

	// build the updated list
	var securities []*models.Security
	for _, v := range state.securities {
		securities = append(securities, v)
	}
	for k, v := range state.slSubscribers {
		securityList := &messages.SecurityList{
			RequestID:  k,
			ResponseID: uint64(time.Now().UnixNano()),
			Securities: securities,
			Success:    true,
		}
		context.Send(v, securityList)
	}

	return nil
}

func (state *Executor) OnMarketableAssetListRequest(context actor.Context) error {
	req := context.Message().(*messages.MarketableAssetListRequest)
	var marketAssets []*models.MarketableAsset
	for _, a := range state.marketAssets {
		marketAssets = append(marketAssets, a)
	}
	response := &messages.MarketableAssetList{
		RequestID:        req.RequestID,
		ResponseID:       uint64(time.Now().UnixNano()),
		MarketableAssets: marketAssets,
		Success:          true,
	}
	if req.Subscribe {
		context.Watch(req.Subscriber)
		state.malSubscribers[req.RequestID] = req.Subscriber
	}
	context.Respond(response)
	return nil
}

func (state *Executor) OnMarketableAssetList(context actor.Context) error {
	req := context.Message().(*messages.MarketableAssetList)
	if len(req.MarketableAssets) == 0 {
		return nil
	}
	market := req.MarketableAssets[0].Market
	// delete all assets from executor for this market
	for k, v := range state.marketAssets {
		if v.Market.ID == market.ID {
			delete(state.marketAssets, k)
		}
	}
	delete(state.symbToMarketAssets, market.ID)

	//replace all the assets
	m := make(map[uint64]*models.MarketableAsset)
	for _, a := range req.MarketableAssets {
		state.marketAssets[a.ProtocolAsset.ProtocolAssetID] = a
		m[a.ProtocolAsset.ProtocolAssetID] = a
	}
	state.symbToMarketAssets[market.ID] = m

	var maList []*models.MarketableAsset
	for _, v := range state.marketAssets {
		maList = append(maList, v)
	}
	for k, v := range state.malSubscribers {
		context.Send(v, &messages.MarketableAssetList{
			RequestID:        k,
			ResponseID:       uint64(time.Now().UnixNano()),
			Success:          true,
			MarketableAssets: maList,
		})
	}
	return nil
}

func (state *Executor) OnTradeCaptureReportRequest(context actor.Context) error {
	msg := context.Message().(*messages.TradeCaptureReportRequest)
	if msg.Account == nil {
		fmt.Println("NIL")
		context.Respond(&messages.TradeCaptureReport{
			RequestID:       msg.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		fmt.Println("NOT FOUND", state.accountManagers)
		context.Respond(&messages.TradeCaptureReport{
			RequestID:       msg.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnAccountMovementRequest(context actor.Context) error {
	msg := context.Message().(*messages.AccountMovementRequest)
	if msg.Account == nil {
		fmt.Println("NIL")
		context.Respond(&messages.AccountMovementResponse{
			RequestID:       msg.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		fmt.Println("NOT FOUND", state.accountManagers)
		context.Respond(&messages.AccountMovementResponse{
			RequestID:       msg.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnPositionsRequest(context actor.Context) error {
	msg := context.Message().(*messages.PositionsRequest)
	if msg.Account == nil {
		context.Respond(&messages.PositionList{
			RequestID:       msg.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.InvalidAccount,
			Positions:       nil,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		context.Respond(&messages.PositionList{
			RequestID:       msg.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.InvalidAccount,
			Positions:       nil,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnBalancesRequest(context actor.Context) error {
	msg := context.Message().(*messages.BalancesRequest)
	if msg.Account == nil {
		context.Respond(&messages.BalanceList{
			RequestID:       msg.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.InvalidAccount,
			Balances:        nil,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		context.Respond(&messages.BalanceList{
			RequestID:       msg.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.InvalidAccount,
			Balances:        nil,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnOrderStatusRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderStatusRequest)
	if msg.Account == nil {
		context.Respond(&messages.OrderList{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		context.Respond(&messages.OrderList{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnNewOrderSingleRequest(context actor.Context) error {
	msg := context.Message().(*messages.NewOrderSingleRequest)
	if msg.Account == nil {
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	if msg.Order == nil {
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidRequest,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnNewOrderBulkRequest(context actor.Context) error {
	msg := context.Message().(*messages.NewOrderBulkRequest)
	if msg.Account == nil {
		context.Respond(&messages.NewOrderBulkResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		context.Respond(&messages.NewOrderBulkResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	if len(msg.Orders) == 0 {
		context.Respond(&messages.NewOrderBulkResponse{
			RequestID: msg.RequestID,
			Success:   true,
			OrderIDs:  nil,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnOrderReplaceRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderReplaceRequest)
	if msg.Account == nil {
		context.Respond(&messages.OrderReplaceResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		context.Respond(&messages.OrderReplaceResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnOrderBulkReplaceRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderBulkReplaceRequest)
	if msg.Account == nil {
		context.Respond(&messages.OrderBulkReplaceResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		context.Respond(&messages.OrderBulkReplaceResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnOrderCancelRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderCancelRequest)
	if msg.Account == nil {
		context.Respond(&messages.OrderCancelResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		context.Respond(&messages.OrderCancelResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnOrderMassCancelRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderMassCancelRequest)
	if msg.Account == nil {
		context.Respond(&messages.OrderCancelResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		context.Respond(&messages.OrderCancelResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnGetAccountRequest(context actor.Context) error {
	request := context.Message().(*commands.GetAccountRequest)
	if request.Account == nil {
		context.Respond(&messages.AccountDataResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: messages.UnknownAccount,
		})
		return nil
	}

	if _, ok := state.accountManagers[request.Account.Name]; ok {
		context.Respond(&commands.GetAccountResponse{
			Account: Portfolio.GetAccount(request.Account.Name),
		})
	} else {
		accnt, err := NewAccount(request.Account)
		if err != nil {
			return fmt.Errorf("error creating account: %v", err)
		}
		producer := NewAccountManagerProducer(accnt, state.Db, false)
		if producer == nil {
			return fmt.Errorf("unknown exchange %s", accnt.Exchange.Name)
		}
		props := actor.PropsFromProducer(producer).WithSupervisor(
			actor.NewExponentialBackoffStrategy(100*time.Second, time.Second))
		pid := context.Spawn(props)
		state.accountManagers[request.Account.Name] = pid
		context.Respond(&commands.GetAccountResponse{
			Account: accnt,
		})
	}
	return nil
}

func (state *Executor) OnTerminated(context actor.Context) error {
	// Handle subscriber krash
	msg := context.Message().(*actor.Terminated)
	for k, v := range state.slSubscribers {
		if v.Id == msg.Who.Id {
			delete(state.slSubscribers, k)
		}
	}

	for k, v := range state.instruments {
		if v.Id == msg.Who.Id {
			delete(state.instruments, k)
		}
	}

	return nil
}
