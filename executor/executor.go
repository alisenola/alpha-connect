package executor

import (
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/chains"
	"gitlab.com/alphaticks/alpha-connect/config"
	"gitlab.com/alphaticks/alpha-connect/utils"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	"google.golang.org/grpc"
	"reflect"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	"gitlab.com/alphaticks/alpha-connect/models/commands"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/protocols"
)

type Executor struct {
	cfg       *config.Config
	exchanges *actor.PID
	protocols *actor.PID
	chains    *actor.PID
	logger    *log.Logger
}

func NewExecutorProducer(cfg *config.Config) actor.Producer {
	return func() actor.Actor {
		return NewExecutor(cfg)
	}
}

func NewExecutor(cfg *config.Config) actor.Actor {
	return &Executor{
		cfg: cfg,
	}
}

// TODO this implementation can easily lead to forgetting adding new messages
func (state *Executor) Receive(context actor.Context) {
	msg := context.Message()
	switch msg.(type) {
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
		if err := state.Initialize(context); err != nil {
			state.logger.Error("error restarting", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor restarting")
	case *messages.AccountDataRequest,
		*messages.MarketDataRequest,
		*messages.UnipoolV3DataRequest,
		*messages.MarketStatisticsRequest,
		*messages.HistoricalUnipoolV3DataRequest,
		*messages.HistoricalFundingRatesRequest,
		*messages.HistoricalLiquidationsRequest,
		*messages.HistoricalSalesRequest,
		*messages.SecurityDefinitionRequest,
		*messages.SecurityListRequest,
		*messages.SecurityList,
		*messages.MarketableProtocolAssetList,
		*messages.MarketableProtocolAssetListRequest,
		*messages.MarketableProtocolAssetDefinitionRequest,
		*messages.AccountMovementRequest,
		*messages.AccountInformationRequest,
		*messages.TradeCaptureReportRequest,
		*messages.PositionsRequest,
		*messages.BalancesRequest,
		*messages.OrderStatusRequest,
		*messages.NewOrderSingleRequest,
		*messages.NewOrderBulkRequest,
		*messages.OrderReplaceRequest,
		*messages.OrderBulkReplaceRequest,
		*messages.OrderCancelRequest,
		*messages.OrderMassCancelRequest,
		*commands.GetAccountRequest:
		if err := state.OnExchangesMessage(context); err != nil {
			state.logger.Error("error processing OnExchangesMessage", log.Error(err))
			panic(err)
		}
	case *messages.ProtocolAssetListRequest,
		*messages.ProtocolAssetList,
		*messages.HistoricalProtocolAssetTransferRequest,
		*messages.ProtocolAssetDataRequest,
		*messages.ProtocolAssetDefinitionRequest:
		if err := state.OnProtocolsMessage(context); err != nil {
			state.logger.Error("error processing OnProtocolsMessage", log.Error(err))
			panic(err)
		}

	case *messages.EVMLogsQueryRequest,
		*messages.EVMLogsSubscribeRequest,
		*messages.EVMContractCallRequest,
		*messages.SVMBlockQueryRequest,
		*messages.SVMEventsQueryRequest,
		*messages.SVMContractCallRequest,
		*messages.SVMTransactionByHashRequest,
		*messages.BlockNumberRequest,
		*messages.BlockInfoRequest:
		if err := state.OnChainsMessage(context); err != nil {
			state.logger.Error("error processing OnChainsMessage", log.Error(err))
			panic(err)
		}

	case *utils.Ready:
		context.Respond(&utils.Ready{})
	}
}

func (state *Executor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()),
	)

	registryAddress := "registry.alphaticks.io:8001"
	if state.cfg.RegistryAddress != "" {
		registryAddress = state.cfg.RegistryAddress
	}
	conn, err := grpc.Dial(registryAddress, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("error connecting to public registry gRPC endpoint: %v", err)
	}
	rgstr := registry.NewPublicRegistryClient(conn)
	assetLoader := context.Spawn(actor.PropsFromProducer(utils.NewStaticLoaderProducer(rgstr)))
	_, err = context.RequestFuture(assetLoader, &utils.Ready{}, 10*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error loading assets: %v", err)
	}

	exProducer := exchanges.NewExecutorProducer(state.cfg, rgstr)
	exProps := actor.PropsFromProducer(exProducer, actor.WithSupervisor(
		actor.NewExponentialBackoffStrategy(100*time.Second, time.Second),
	))
	exEx, err := context.SpawnNamed(exProps, "exchanges")
	if err != nil {
		return fmt.Errorf("error spawning exchanges executor: %v", err)
	}
	state.exchanges = exEx

	prProducer := protocols.NewExecutorProducer(state.cfg, rgstr)
	prProps := actor.PropsFromProducer(prProducer, actor.WithSupervisor(
		actor.NewExponentialBackoffStrategy(100*time.Second, time.Second),
	))
	prEx, err := context.SpawnNamed(prProps, "protocols")
	if err != nil {
		return fmt.Errorf("error spawning protocols executor: %v", err)
	}
	state.protocols = prEx

	chProducer := chains.NewExecutorProducer(state.cfg, rgstr)
	chProps := actor.PropsFromProducer(chProducer, actor.WithSupervisor(
		actor.NewExponentialBackoffStrategy(100*time.Second, time.Second),
	))
	chEx, err := context.SpawnNamed(chProps, "chains")
	if err != nil {
		return fmt.Errorf("error spawning chains executor: %v", err)
	}
	state.chains = chEx
	return nil
}

func (state *Executor) OnExchangesMessage(context actor.Context) error {
	if state.exchanges == nil {
		return fmt.Errorf("missing exchanges executor")
	}
	context.Forward(state.exchanges)
	return nil
}

func (state *Executor) OnProtocolsMessage(context actor.Context) error {
	if state.protocols == nil {
		return fmt.Errorf("missing protocols executor")
	}
	context.Forward(state.protocols)
	return nil
}

func (state *Executor) OnChainsMessage(context actor.Context) error {
	if state.chains == nil {
		return fmt.Errorf("missing protocols executor")
	}
	context.Forward(state.chains)
	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}
