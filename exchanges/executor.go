package exchanges

import (
	"errors"
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/config"
	"gitlab.com/alphaticks/alpha-connect/exchanges/types"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	_ "gitlab.com/alphaticks/tickfunctors/market/portfolio"
	tickstore_go_client "gitlab.com/alphaticks/tickstore-go-client"
	tickstore_types "gitlab.com/alphaticks/tickstore-types"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/hitbtc"
	xmodels "gitlab.com/alphaticks/xchanger/models"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/commands"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// The executor routes all the request to the underlying exchange executor & listeners
// He is the main part of the whole software..

type Executor struct {
	*config.Config
	store                    tickstore_types.TickstoreClient
	db                       *gorm.DB
	registry                 registry.StaticClient
	dialerPool               *xchangerUtils.DialerPool
	wsPools                  map[uint32]*xchangerUtils.WebsocketPool
	accountClients           map[string]map[string]*http.Client
	accountManagers          map[string]*actor.PID
	executors                map[uint32]*actor.PID                      // A map from exchange ID to executor
	securities               map[uint64]*models.Security                // A map from security ID to security
	marketableProtocolAssets map[uint64]*models.MarketableProtocolAsset // A map from MarketAsset ID to MarketAsset
	symbToSecs               map[uint32]map[string]*models.Security
	instruments              map[uint64]*actor.PID // A map from security ID to market manager
	slSubscribers            map[uint64]*actor.PID // A map from request ID to security list subscribers
	malSubscribers           map[uint64]*actor.PID // A map from request ID to protocolAssets list subscribers
	logger                   *log.Logger
	strict                   bool
}

func NewExecutorProducer(cfg *config.Config, registry registry.StaticClient) actor.Producer {
	return func() actor.Actor {
		return NewExecutor(cfg, registry)
	}
}

func NewExecutor(cfg *config.Config, registry registry.StaticClient) actor.Actor {
	return &Executor{
		Config:   cfg,
		registry: registry,
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

	case *messages.MarketableProtocolAssetDefinitionRequest:
		if err := state.OnMarketableProtocolsAssetDefinitionRequest(context); err != nil {
			state.logger.Error("error processing OnMarketableProtocolsAssetDefinitionRequest", log.Error(err))
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

	case *messages.MarketableProtocolAssetListRequest:
		if err := state.OnMarketableProtocolAssetListRequest(context); err != nil {
			state.logger.Error("error processing OnMarketableProtocolAssetListRequest", log.Error(err))
			panic(err)
		}

	case *messages.MarketableProtocolAssetList:
		if err := state.OnMarketableProtocolAssetList(context); err != nil {
			state.logger.Error("error processing OnMarketableProtocolAssetList", log.Error(err))
			panic(err)
		}

	case *messages.AccountInformationRequest:
		if err := state.OnAccountInformationRequest(context); err != nil {
			state.logger.Error("error processing OnAccountInformationRequest", log.Error(err))
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
	state.wsPools = make(map[uint32]*xchangerUtils.WebsocketPool)

	var dialerPool *xchangerUtils.DialerPool
	interfaceName := state.DialerPoolInterface
	ips := state.DialerPoolIPs
	if interfaceName != "" {
		addresses, err := localAddresses(interfaceName)
		if err != nil {
			panic(fmt.Errorf("error getting interface addressses: %v", err))
		}
		var dialers []*net.Dialer
		for _, a := range addresses {
			state.logger.Info("using interface address " + a.String())
			dialers = append(dialers, &net.Dialer{
				LocalAddr: a,
			})
		}
		dialerPool, err = xchangerUtils.NewDialerPool(dialers)
		if err != nil {
			return fmt.Errorf("error creating dialer pool: %v", err)
		}
	} else if len(ips) > 0 {
		var dialers []*net.Dialer
		for _, ip := range ips {
			tcpAddr := &net.TCPAddr{
				IP: net.ParseIP(ip),
			}
			dialers = append(dialers, &net.Dialer{
				LocalAddr: tcpAddr,
			})
		}
		var err error
		dialerPool, err = xchangerUtils.NewDialerPool(dialers)
		if err != nil {
			return fmt.Errorf("error creating dialer pool: %v", err)
		}
	} else {
		dialerPool = xchangerUtils.DefaultDialerPool
	}
	state.dialerPool = dialerPool

	if state.Store != nil {
		client, err := tickstore_go_client.NewRemoteClient(*state.Store)
		if err != nil {
			return fmt.Errorf("error creating monitor store client: %v", err)
		}
		state.store = client
		// register measurements
		if err := state.store.RegisterMeasurement("portfolio", "PortfolioTracker"); err != nil {
			return fmt.Errorf("error registering portfolio measurement: %v", err)
		}
	}

	if state.DB != nil {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			state.DB.PostgresHost,
			state.DB.PostgresUser,
			state.DB.PostgresPassword,
			state.DB.PostgresDB,
			state.DB.PostgresPort)
		sql, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		if state.DB.Migrate {
			fmt.Println("MIGRATING")
			if err := sql.AutoMigrate(&types.Account{}); err != nil {
				return fmt.Errorf("error migrating account type: %v", err)
			}
			if err := sql.AutoMigrate(&types.Transaction{}); err != nil {
				return fmt.Errorf("error migrating transaction type: %v", err)
			}
			if err := sql.AutoMigrate(&types.Fill{}); err != nil {
				return fmt.Errorf("error migrating transaction type: %v", err)
			}
			if err := sql.AutoMigrate(&types.Movement{}); err != nil {
				return fmt.Errorf("error migrating transaction type: %v", err)
			}
		}
		state.db = sql
	}

	// TODO add dialer pool test

	state.accountClients = make(map[string]map[string]*http.Client)
	// Spawn all exchange executors
	state.executors = make(map[uint32]*actor.PID)
	for _, exchStr := range state.Exchanges {
		exch, ok := constants.GetExchangeByName(exchStr)
		if !ok {
			state.logger.Warn(fmt.Sprintf("unknown exchange %s", exchStr))
			continue
		}
		state.accountClients[exch.Name] = make(map[string]*http.Client)
		for _, accntCfg := range state.Accounts {
			accntExch, ok := constants.GetExchangeByName(accntCfg.Exchange)
			if !ok {
				return fmt.Errorf("unknown exchange %s", accntCfg.Exchange)
			}
			if accntExch.ID == exch.ID && accntCfg.SOCKS5 != "" {
				dialSocksProxy, err := proxy.SOCKS5("tcp", accntCfg.SOCKS5, nil, &net.Dialer{
					Timeout: 10 * time.Second,
				})
				if err != nil {
					return fmt.Errorf("error creating SOCKS5 proxy: %v", err)
				}
				state.accountClients[exch.Name][accntCfg.Name] = &http.Client{
					Transport: &http.Transport{
						Proxy:                 http.ProxyFromEnvironment,
						DialContext:           dialSocksProxy.(proxy.ContextDialer).DialContext,
						IdleConnTimeout:       5 * time.Second,
						TLSHandshakeTimeout:   5 * time.Second,
						ExpectContinueTimeout: 1 * time.Second,
						MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
					},
				}
			}
		}
		producer := NewExchangeExecutorProducer(exch, state.Config, state.dialerPool, state.registry, state.accountClients[exch.Name])
		if producer == nil {
			state.logger.Warn(fmt.Sprintf("unknown exchange %s", exchStr))
			continue
		}
		props := actor.PropsFromProducer(
			producer,
			actor.WithMailbox(actor.UnboundedPriority()),
			actor.WithSupervisor(actor.NewExponentialBackoffStrategy(100*time.Second, time.Second)))

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

	//Request marketable protocol assets for all of them
	var futs []*actor.Future
	req := messages.MarketableProtocolAssetListRequest{
		RequestID: 0,
		Subscribe: true,
	}
	for _, pid := range state.executors {
		f := context.RequestFuture(pid, &req, 15*time.Second)
		futs = append(futs, f)
	}

	state.marketableProtocolAssets = make(map[uint64]*models.MarketableProtocolAsset)
	for _, f := range futs {
		resp, err := f.Result()
		if err != nil {
			state.logger.Error("error fetching marketable protocol assets for one venue", log.Error(err))
			continue
		}
		list, ok := resp.(*messages.MarketableProtocolAssetList)
		if !ok {
			return fmt.Errorf("was expecting MarketAssetList, got %s", reflect.TypeOf(resp).String())
		}
		if !list.Success {
			return errors.New(list.RejectionReason.String())
		}
		if len(list.MarketableProtocolAssets) > 0 {
			idToMarketAssets := make(map[uint64]*models.MarketableProtocolAsset)
			for _, asset := range list.MarketableProtocolAssets {
				if asset2, ok := state.marketableProtocolAssets[asset.MarketableProtocolAssetID]; ok {
					return fmt.Errorf("got two protocol assets with the same ID: %v and %v", asset.ProtocolAsset, asset2.ProtocolAsset)
				}
				state.marketableProtocolAssets[asset.MarketableProtocolAssetID] = asset
				idToMarketAssets[asset.MarketableProtocolAssetID] = asset
			}
		}
	}

	// Spawn all account listeners
	state.accountManagers = make(map[string]*actor.PID)
	for _, accntCfg := range state.Accounts {
		exch, ok := constants.GetExchangeByName(accntCfg.Exchange)
		if !ok {
			return fmt.Errorf("unknown exchange %s", accntCfg.Exchange)
		}
		account, err := NewAccount(&models.Account{
			Portfolio: accntCfg.Portfolio,
			Name:      accntCfg.Name,
			Exchange:  exch,
			ApiCredentials: &xmodels.APICredentials{
				APIKey:    accntCfg.ApiKey,
				APISecret: accntCfg.ApiSecret,
				AccountID: accntCfg.ID,
			},
		})
		if err != nil {
			return fmt.Errorf("error creating new account: %v", err)
		}
		client := state.accountClients[exch.Name][accntCfg.Name]
		producer := NewAccountManagerProducer(accntCfg, account, state.store, state.db, state.registry, client)
		if producer == nil {
			return fmt.Errorf("unknown exchange %s", accntCfg.Exchange)
		}
		props := actor.PropsFromProducer(producer, actor.WithSupervisor(
			actor.NewExponentialBackoffStrategy(100*time.Second, time.Second)))
		state.accountManagers[accntCfg.Name], err = context.SpawnNamed(props, accntCfg.Name+"_account")
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
			RejectionReason: messages.RejectionReason_UnknownAccount,
		})
		return nil
	}

	if pid, ok := state.accountManagers[request.Account.Name]; ok {
		context.Forward(pid)
	} else {
		context.Respond(&messages.AccountDataResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: messages.RejectionReason_UnknownAccount,
		})
	}

	return nil
}

func (state *Executor) getSecurity(instr *models.Instrument) (*models.Security, *messages.RejectionReason) {
	if instr == nil {
		rej := messages.RejectionReason_MissingInstrument
		return nil, &rej
	}
	if instr.SecurityID != nil {
		if sec, ok := state.securities[instr.SecurityID.Value]; ok {
			return sec, nil
		} else {
			rej := messages.RejectionReason_UnknownSecurityID
			return nil, &rej
		}
	} else if instr.Symbol != nil {
		if instr.Exchange == nil {
			rej := messages.RejectionReason_UnknownExchange
			return nil, &rej
		}
		symbolsToSecs, ok := state.symbToSecs[instr.Exchange.ID]
		if !ok {
			rej := messages.RejectionReason_UnknownExchange
			return nil, &rej
		}
		sec, ok := symbolsToSecs[instr.Symbol.Value]
		if !ok {
			rej := messages.RejectionReason_UnknownSymbol
			return nil, &rej
		}
		return sec, nil
	} else {
		rej := messages.RejectionReason_MissingInstrument
		return nil, &rej
	}
}

func (state *Executor) getExchange(instr *models.Instrument) (*xmodels.Exchange, *messages.RejectionReason) {
	if instr == nil {
		rej := messages.RejectionReason_MissingInstrument
		return nil, &rej
	}
	if instr.SecurityID != nil {
		if sec, ok := state.securities[instr.SecurityID.Value]; ok {
			return sec.Exchange, nil
		} else {
			rej := messages.RejectionReason_UnknownSecurityID
			return nil, &rej
		}
	} else if instr.Symbol != nil {
		if instr.Exchange == nil {
			rej := messages.RejectionReason_UnknownExchange
			return nil, &rej
		}
		symbolsToSecs, ok := state.symbToSecs[instr.Exchange.ID]
		if !ok {
			rej := messages.RejectionReason_UnknownExchange
			return nil, &rej
		}
		sec, ok := symbolsToSecs[instr.Symbol.Value]
		if !ok {
			rej := messages.RejectionReason_UnknownSymbol
			return nil, &rej
		}
		return sec.Exchange, nil
	} else if instr.Exchange != nil {
		return instr.Exchange, nil
	} else {
		rej := messages.RejectionReason_MissingInstrument
		return nil, &rej
	}
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
		var wsp *xchangerUtils.WebsocketPool
		var err error
		if sec.Exchange.ID == constants.HITBTC.ID {
			if _, ok := state.wsPools[sec.Exchange.ID]; !ok {
				wsp, err = hitbtc.NewWebsocketPool(state.dialerPool)
				// TODO logging
				if err != nil {
					context.Respond(&messages.MarketDataResponse{
						RequestID:       request.RequestID,
						Success:         false,
						RejectionReason: messages.RejectionReason_Other,
					})
					return nil
				}
				state.wsPools[sec.Exchange.ID] = wsp
			}
			wsp = state.wsPools[sec.Exchange.ID]
		}

		props := actor.PropsFromProducer(NewDataManagerProducer(sec, state.dialerPool, wsp), actor.WithSupervisor(
			utils.NewExponentialBackoffStrategy(100*time.Second, time.Second, time.Second)))
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
		props := actor.PropsFromProducer(NewDataManagerProducer(sec, state.dialerPool, nil), actor.WithSupervisor(
			utils.NewExponentialBackoffStrategy(100*time.Second, time.Second, time.Second)))
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
			RejectionReason: messages.RejectionReason_UnknownExchange,
		})
	}
	context.Forward(exchange)

	return nil
}

func (state *Executor) OnMarketStatisticsRequest(context actor.Context) error {
	request := context.Message().(*messages.MarketStatisticsRequest)
	ex, rej := state.getExchange(request.Instrument)
	if rej != nil {
		context.Respond(&messages.MarketStatisticsResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: *rej,
		})
		return nil
	}
	exchange, ok := state.executors[ex.ID]
	if !ok {
		context.Respond(&messages.MarketStatisticsResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: messages.RejectionReason_UnknownExchange,
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
			RejectionReason: messages.RejectionReason_UnknownExchange,
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
			RejectionReason: messages.RejectionReason_UnknownExchange,
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
			RejectionReason: messages.RejectionReason_UnknownExchange,
		})
		return nil
	}
	context.Forward(exchange)

	return nil
}

func (state *Executor) OnHistoricalSalesRequest(context actor.Context) error {
	request := context.Message().(*messages.HistoricalSalesRequest)
	passet, ok := state.marketableProtocolAssets[request.MarketableProtocolAssetID]
	if !ok {
		context.Respond(&messages.HistoricalSalesResponse{
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.RejectionReason_UnknownProtocolAsset,
		})
		return nil
	}
	market, ok := state.executors[passet.Market.ID]
	if !ok {
		context.Respond(&messages.HistoricalSalesResponse{
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.RejectionReason_UnknownExchange,
		})
		return nil
	}
	context.Forward(market)
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

func (state *Executor) OnMarketableProtocolsAssetDefinitionRequest(context actor.Context) error {
	req := context.Message().(*messages.MarketableProtocolAssetDefinitionRequest)
	ma, ok := state.marketableProtocolAssets[req.MarketableProtocolAssetID]
	if !ok {
		context.Respond(&messages.MarketableProtocolAssetDefinitionResponse{
			RequestID:       req.RequestID,
			Success:         false,
			RejectionReason: messages.RejectionReason_UnknownProtocolAsset,
		})
	}
	context.Respond(&messages.MarketableProtocolAssetDefinitionResponse{
		RequestID:               req.RequestID,
		ResponseID:              uint64(time.Now().UnixNano()),
		Success:                 true,
		MarketableProtocolAsset: ma,
	})
	return nil
}

func (state *Executor) OnMarketableProtocolAssetListRequest(context actor.Context) error {
	req := context.Message().(*messages.MarketableProtocolAssetListRequest)
	var marketAssets []*models.MarketableProtocolAsset
	for _, a := range state.marketableProtocolAssets {
		marketAssets = append(marketAssets, a)
	}
	response := &messages.MarketableProtocolAssetList{
		RequestID:                req.RequestID,
		ResponseID:               uint64(time.Now().UnixNano()),
		MarketableProtocolAssets: marketAssets,
		Success:                  true,
	}
	if req.Subscribe {
		context.Watch(req.Subscriber)
		state.malSubscribers[req.RequestID] = req.Subscriber
	}
	context.Respond(response)
	return nil
}

func (state *Executor) OnMarketableProtocolAssetList(context actor.Context) error {
	req := context.Message().(*messages.MarketableProtocolAssetList)
	if len(req.MarketableProtocolAssets) == 0 {
		return nil
	}
	market := req.MarketableProtocolAssets[0].Market
	// delete all assets from executor for this market
	for k, v := range state.marketableProtocolAssets {
		if v.Market.ID == market.ID {
			delete(state.marketableProtocolAssets, k)
		}
	}

	//replace all the assets
	for _, a := range req.MarketableProtocolAssets {
		state.marketableProtocolAssets[a.MarketableProtocolAssetID] = a
	}

	var maList []*models.MarketableProtocolAsset
	for _, v := range state.marketableProtocolAssets {
		maList = append(maList, v)
	}
	for k, v := range state.malSubscribers {
		context.Send(v, &messages.MarketableProtocolAssetList{
			RequestID:                k,
			ResponseID:               uint64(time.Now().UnixNano()),
			Success:                  true,
			MarketableProtocolAssets: maList,
		})
	}
	return nil
}

func (state *Executor) OnTradeCaptureReportRequest(context actor.Context) error {
	msg := context.Message().(*messages.TradeCaptureReportRequest)
	if msg.Account == nil {
		context.Respond(&messages.TradeCaptureReport{
			RequestID:       msg.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.RejectionReason_InvalidAccount,
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
			RejectionReason: messages.RejectionReason_InvalidAccount,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnAccountInformationRequest(context actor.Context) error {
	msg := context.Message().(*messages.AccountInformationRequest)
	if msg.Account == nil {
		fmt.Println("NIL")
		context.Respond(&messages.AccountInformationResponse{
			RequestID:       msg.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.RejectionReason_InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		fmt.Println("NOT FOUND", state.accountManagers)
		context.Respond(&messages.AccountInformationResponse{
			RequestID:       msg.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.RejectionReason_InvalidAccount,
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
			RejectionReason: messages.RejectionReason_InvalidAccount,
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
			RejectionReason: messages.RejectionReason_InvalidAccount,
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
			RejectionReason: messages.RejectionReason_InvalidAccount,
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
			RejectionReason: messages.RejectionReason_InvalidAccount,
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
			RejectionReason: messages.RejectionReason_InvalidAccount,
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
			RejectionReason: messages.RejectionReason_InvalidAccount,
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
			RejectionReason: messages.RejectionReason_InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		context.Respond(&messages.OrderList{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.RejectionReason_InvalidAccount,
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
			RejectionReason: messages.RejectionReason_InvalidAccount,
		})
		return nil
	}
	if msg.Order == nil {
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.RejectionReason_InvalidRequest,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.RejectionReason_InvalidAccount,
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
			RejectionReason: messages.RejectionReason_InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		context.Respond(&messages.NewOrderBulkResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.RejectionReason_InvalidAccount,
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
			RejectionReason: messages.RejectionReason_InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		context.Respond(&messages.OrderReplaceResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.RejectionReason_InvalidAccount,
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
			RejectionReason: messages.RejectionReason_InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		context.Respond(&messages.OrderBulkReplaceResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.RejectionReason_InvalidAccount,
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
			RejectionReason: messages.RejectionReason_InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		context.Respond(&messages.OrderCancelResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.RejectionReason_InvalidAccount,
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
			RejectionReason: messages.RejectionReason_InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.Name]
	if !ok {
		context.Respond(&messages.OrderCancelResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.RejectionReason_InvalidAccount,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnGetAccountRequest(context actor.Context) error {
	req := context.Message().(*commands.GetAccountRequest)
	if req.Account == nil {
		context.Respond(&commands.GetAccountResponse{
			Err: fmt.Errorf("unknown account"),
		})
		return nil
	}

	account := GetAccount(req.Account.Name)
	if account == nil {
		exch, ok := constants.GetExchangeByName(req.Account.Exchange)
		if !ok {
			context.Respond(&commands.GetAccountResponse{
				Err: fmt.Errorf("unknown exchange: %s", req.Account.Exchange),
			})
			return nil
		}
		if req.Account.SOCKS5 != "" {
			dialSocksProxy, err := proxy.SOCKS5("tcp", req.Account.SOCKS5, nil, &net.Dialer{
				Timeout: 10 * time.Second,
			})
			if err != nil {
				return fmt.Errorf("error creating SOCKS5 proxy: %v", err)
			}
			state.accountClients[exch.Name][req.Account.Name] = &http.Client{
				Transport: &http.Transport{
					Proxy:                 http.ProxyFromEnvironment,
					DialContext:           dialSocksProxy.(proxy.ContextDialer).DialContext,
					IdleConnTimeout:       60 * time.Second,
					TLSHandshakeTimeout:   10 * time.Second,
					ExpectContinueTimeout: 1 * time.Second,
					//MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
				},
			}
		}
		var err error
		account, err = NewAccount(&models.Account{
			Name:      req.Account.Name,
			Portfolio: req.Account.Portfolio,
			Exchange:  exch,
			ApiCredentials: &xmodels.APICredentials{
				APIKey:    req.Account.ApiKey,
				APISecret: req.Account.ApiSecret,
				AccountID: req.Account.ID,
			},
		})
		if err != nil {
			return fmt.Errorf("error creating new account: %v", err)
		}
		producer := NewAccountManagerProducer(*req.Account, account, state.store, state.db, state.registry, state.accountClients[exch.Name][req.Account.Name])
		if producer == nil {
			return fmt.Errorf("unknown exchange %s", req.Account.Exchange)
		}
		props := actor.PropsFromProducer(producer, actor.WithSupervisor(
			actor.NewExponentialBackoffStrategy(100*time.Second, time.Second)))
		state.accountManagers[req.Account.Name] = context.Spawn(props)
	}
	context.Respond(&commands.GetAccountResponse{
		Account: account,
	})
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
