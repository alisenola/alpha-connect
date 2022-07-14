package utils

import (
	goContext "context"
	"fmt"
	"reflect"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/models"
)

type checkStatic struct{}
type Ready struct{}

type StaticLoader struct {
	registry registry.PublicRegistryClient
	logger   *log.Logger
	ticker   *time.Ticker
}

func NewStaticLoaderProducer(rgstry registry.PublicRegistryClient) actor.Producer {
	return func() actor.Actor {
		return NewStaticLoader(rgstry)
	}
}

func NewStaticLoader(rgstry registry.PublicRegistryClient) actor.Actor {
	return &StaticLoader{
		registry: rgstry,
	}
}

func (state *StaticLoader) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.logger.Error("error initializing actor", log.Error(err))
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

	case *checkStatic:
		if err := state.checkStatic(context); err != nil {
			state.logger.Error("error checkStatic", log.Error(err))
		}

	case *Ready:
		context.Respond(&Ready{})
	}
}

func (state *StaticLoader) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	ticker := time.NewTicker(20 * time.Second)
	state.ticker = ticker
	go func(pid *actor.PID) {
		for {
			select {
			case <-ticker.C:
				context.Send(pid, &checkStatic{})
			case <-time.After(2 * time.Minute):
				if state.ticker != ticker {
					return
				}
			}
		}
	}(context.Self())

	return state.checkStatic(context)
}

func (state *StaticLoader) Clean(context actor.Context) error {
	if state.ticker != nil {
		state.ticker.Stop()
		state.ticker = nil
	}

	return nil
}

func (state *StaticLoader) checkStatic(context actor.Context) error {
	res, err := state.registry.Assets(goContext.Background(), &registry.AssetsRequest{})
	if err != nil {
		return fmt.Errorf("error fetching assets: %v", err)
	}
	assets := make(map[uint32]*models.Asset)
	for _, a := range res.Assets {
		assets[a.AssetId] = &models.Asset{
			Symbol: a.Symbol,
			Name:   a.Name,
			ID:     a.AssetId,
		}
	}
	if err := constants.LoadAssets(assets); err != nil {
		return err
	}

	resp, err := state.registry.Protocols(goContext.Background(), &registry.ProtocolsRequest{})
	if err != nil {
		return fmt.Errorf("error fetching assets: %v", err)
	}
	protocols := make(map[uint32]*models.Protocol)
	for _, p := range resp.Protocols {
		protocols[p.ProtocolId] = &models.Protocol{
			ID:      p.ProtocolId,
			Name:    p.Name,
			ChainID: p.ChainId,
		}
	}
	if err := constants.LoadProtocols(protocols); err != nil {
		return err
	}

	resc, err := state.registry.Chains(goContext.Background(), &registry.ChainsRequest{})
	if err != nil {
		return fmt.Errorf("error fetching assets: %v", err)
	}
	chains := make(map[uint32]*models.Chain)
	for _, c := range resc.Chains {
		chains[c.ChainId] = &models.Chain{
			ID:   c.ChainId,
			Type: c.Type,
			Name: c.Name,
		}
	}
	if err := constants.LoadChains(chains); err != nil {
		return err
	}

	return nil
}
