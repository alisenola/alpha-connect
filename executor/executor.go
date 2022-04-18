package executor

import (
	"fmt"
	"reflect"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	"gitlab.com/alphaticks/alpha-connect/models/commands"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/protocols"
)

type Executor struct {
	cfgEx      *exchanges.ExecutorConfig
	cfgPr      *protocols.ExecutorConfig
	executorEx *actor.PID
	executorPr *actor.PID
	logger     *log.Logger
}

func NewExecutorProducer(cfgEx *exchanges.ExecutorConfig, cfgPr *protocols.ExecutorConfig) actor.Producer {
	return func() actor.Actor {
		return NewExecutor(cfgEx, cfgPr)
	}
}

func NewExecutor(cfgEx *exchanges.ExecutorConfig, cfgPr *protocols.ExecutorConfig) actor.Actor {
	return &Executor{
		cfgEx: cfgEx,
		cfgPr: cfgPr,
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
		*messages.SecurityDefinitionRequest,
		*messages.SecurityListRequest,
		*messages.SecurityList,
		*messages.AccountMovementRequest,
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
		//state.logger.Info("message forwarded to exchange executor")
	case *messages.ProtocolAssetListRequest,
		*messages.ProtocolAssetList,
		*messages.HistoricalProtocolAssetTransferRequest,
		*messages.ProtocolAssetTransferRequest,
		*messages.ProtocolAssetDefinitionRequest:
		if err := state.OnProtocolsMessage(context); err != nil {
			state.logger.Error("error processing OnProtocolsMessage", log.Error(err))
			panic(err)
		}
		//state.logger.Info("message forwarded to to protocols executor")
	}
}

func (state *Executor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()),
	)

	exProducer := exchanges.NewExecutorProducer(state.cfgEx)
	exProps := actor.PropsFromProducer(exProducer).WithSupervisor(
		actor.NewExponentialBackoffStrategy(100*time.Second, time.Second),
	)
	exEx, err := context.SpawnNamed(exProps, "exchanges")
	if err != nil {
		return fmt.Errorf("error spawning exchanges executor: %v", err)
	}
	state.executorEx = exEx

	prProducer := protocols.NewExecutorProducer(state.cfgPr)
	prProps := actor.PropsFromProducer(prProducer).WithSupervisor(
		actor.NewExponentialBackoffStrategy(100*time.Second, time.Second),
	)
	prEx, err := context.SpawnNamed(prProps, "protocols")
	if err != nil {
		return fmt.Errorf("error spawning protocols executor: %v", err)
	}
	state.executorPr = prEx
	return nil
}

func (state *Executor) OnExchangesMessage(context actor.Context) error {
	if state.executorEx == nil {
		return fmt.Errorf("missing exchanges executor")
	}
	context.Forward(state.executorEx)
	return nil
}

func (state *Executor) OnProtocolsMessage(context actor.Context) error {
	if state.executorPr == nil {
		return fmt.Errorf("missing protocols executor")
	}
	context.Forward(state.executorPr)
	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}
