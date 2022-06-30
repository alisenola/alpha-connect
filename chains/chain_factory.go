package chains

import (
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/chains/evm"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	models2 "gitlab.com/alphaticks/xchanger/models"
)

func NewChainExecutorProducer(chain *models2.Chain, registry registry.PublicRegistryClient) actor.Producer {
	switch chain.Type {
	case "EVM":
		switch chain.ID {
		case 1:
			return func() actor.Actor {
				return evm.NewExecutor(registry, "wss://eth-mainnet.alchemyapi.io/v2/4kjvftiD6NzHc6kkD1ih3-5wilV--3mz")
			}
		case 10:
			return func() actor.Actor {
				return evm.NewExecutor(registry, "wss://opt-mainnet.g.alchemy.com/v2/MAdidiXxtFnW5b4q9pTmBLcTW73SHoMN")
			}
		case 147:
			return func() actor.Actor {
				return evm.NewExecutor(registry, "wss://polygon-mainnet.g.alchemy.com/v2/PYNN12EJrMrmlxWjy9KZyrYK6GHrErCM")
			}
		default:
			return nil
		}
	default:
		return nil
	}
}
