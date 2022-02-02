package protocols

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/protocols/erc721"
	"gitlab.com/alphaticks/xchanger/constants"
	models2 "gitlab.com/alphaticks/xchanger/models"
)

func NewProtocolExecutorProducer(protocol *models2.Protocol, config *ExecutorConfig) actor.Producer {
	switch protocol.ID {
	case constants.ERC721.ID:
		return func() actor.Actor { return erc721.NewExecutor(config.Registry) }
	default:
		return nil
	}
}
