package protocols

import (
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/protocols/erc20"
	"gitlab.com/alphaticks/alpha-connect/protocols/erc721"
	"gitlab.com/alphaticks/xchanger/constants"
	models2 "gitlab.com/alphaticks/xchanger/models"
)

func NewProtocolExecutorProducer(protocol *models2.Protocol, config *ExecutorConfig) actor.Producer {
	switch protocol.ID {
	case constants.ERC20.ID:
		return func() actor.Actor { return erc20.NewExecutor(config.Registry) }
	case constants.ERC721.ID:
		return func() actor.Actor { return erc721.NewExecutor(config.Registry) }
	default:
		return nil
	}
}

func NewProtocolAssetListenerProducer(asset *models.ProtocolAsset) actor.Producer {
	switch asset.Protocol.ID {
	case constants.ERC721.ID:
		return func() actor.Actor { return erc721.NewListener(asset) }
	}
	return nil
}
