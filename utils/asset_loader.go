package utils

import (
	goContext "context"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/models"
	"reflect"
	"time"
)

type checkAsset struct{}
type Ready struct{}

type AssetLoader struct {
	registry registry.PublicRegistryClient
	logger   *log.Logger
	ticker   *time.Ticker
}

func NewAssetLoaderProducer(rgstry registry.PublicRegistryClient) actor.Producer {
	return func() actor.Actor {
		return NewAssetLoader(rgstry)
	}
}

func NewAssetLoader(rgstry registry.PublicRegistryClient) actor.Actor {
	return &AssetLoader{
		registry: rgstry,
	}
}

func (state *AssetLoader) Receive(context actor.Context) {
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

	case *checkAsset:
		if err := state.checkAsset(context); err != nil {
			state.logger.Error("error checkAsset", log.Error(err))
		}

	case *Ready:
		if err := state.onReady(context); err != nil {
			state.logger.Error("error onReady", log.Error(err))
			panic(err)
		}
	}
}

func (state *AssetLoader) Initialize(context actor.Context) error {
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
			case _ = <-ticker.C:
				context.Send(pid, &checkAsset{})
			case <-time.After(2 * time.Minute):
				// timer stopped, we leave
				return
			}
		}
	}(context.Self())

	return state.checkAsset(context)
}

func (state *AssetLoader) Clean(context actor.Context) error {
	if state.ticker != nil {
		state.ticker.Stop()
		state.ticker = nil
	}

	return nil
}

func (state *AssetLoader) onReady(context actor.Context) error {
	context.Respond(&Ready{})
	return nil
}

func (state *AssetLoader) checkAsset(context actor.Context) error {
	res, err := state.registry.Assets(goContext.Background(), &registry.AssetsRequest{
		Filter: &registry.AssetFilter{
			Fungible: true,
		},
	})
	if err != nil {
		return fmt.Errorf("error fetching assets: %v", err)
	}
	assets := make(map[uint32]models.Asset)
	for _, a := range res.Assets {
		assets[a.AssetId] = models.Asset{
			Symbol: a.Symbol,
			Name:   a.Name,
			ID:     a.AssetId,
		}
	}

	if err := constants.LoadAssets(assets); err != nil {
		return err
	}

	return nil
}
