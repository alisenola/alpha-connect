package chains

import (
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/chains/evm"
	models2 "gitlab.com/alphaticks/xchanger/models"
)

func NewChainExecutorProducer(chain *models2.Chain, config *ExecutorConfig) actor.Producer {
	switch chain.Type {
	case "EVM":
		return func() actor.Actor { return evm.NewExecutor(config.Registry) }
	default:
		return nil
	}
}
