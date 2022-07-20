package chains

import (
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/chains/evm"
	"gitlab.com/alphaticks/alpha-connect/chains/svm"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	models2 "gitlab.com/alphaticks/xchanger/models"
)

func NewChainExecutorProducer(chain *models2.Chain, rpc string, registry registry.PublicRegistryClient) actor.Producer {
	switch chain.Type {
	case "EVM":
		return func() actor.Actor {
			return evm.NewExecutor(registry, rpc)
		}
	case "SVM":
		return func() actor.Actor {
			return svm.NewExecutor(registry, rpc)
		}
	default:
		return nil
	}
}
