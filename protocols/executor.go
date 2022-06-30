package protocols

import (
	"errors"
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/config"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	"gitlab.com/alphaticks/xchanger/constants"
	"reflect"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
)

// The executor routes all the request to the underlying exchange executor & listeners
// He is the main part of the whole software..

type Executor struct {
	*config.Config
	registry       registry.PublicRegistryClient
	executors      map[uint32]*actor.PID // A map from exchange ID to executor
	protocolAssets map[uint64]*models.ProtocolAsset
	alSubscribers  map[uint64]*actor.PID // A map from request ID to asset list subscriber
	dataManagers   map[uint64]*actor.PID // A map from asset ID to data manager
	logger         *log.Logger
	strict         bool
}

func NewExecutorProducer(cfg *config.Config, registry registry.PublicRegistryClient) actor.Producer {
	return func() actor.Actor {
		return NewExecutor(cfg, registry)
	}
}

func NewExecutor(cfg *config.Config, registry registry.PublicRegistryClient) actor.Actor {
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
	case *messages.ProtocolAssetListRequest:
		if err := state.OnProtocolAssetListRequest(context); err != nil {
			state.logger.Error("error processing ProtocolAssetListRequest", log.Error(err))
			panic(err)
		}
	case *messages.ProtocolAssetList:
		if err := state.OnProtocolAssetList(context); err != nil {
			state.logger.Error("error processing ProtocolAssetList", log.Error(err))
			panic(err)
		}
	case *messages.HistoricalProtocolAssetTransferRequest:
		if err := state.OnHistoricalProtocolAssetTransferRequest(context); err != nil {
			state.logger.Error("error processing HistoricalProtocolAssetTransferRequest", log.Error(err))
			panic(err)
		}
	case *messages.ProtocolAssetDataRequest:
		if err := state.OnProtocolAssetDataRequest(context); err != nil {
			state.logger.Error("error processing ProtocolAssetDataRequest", log.Error(err))
			panic(err)
		}
	case *messages.ProtocolAssetDefinitionRequest:
		if err := state.OnProtocolAssetDefinition(context); err != nil {
			state.logger.Error("error processing OnProtocolAssetDefinition", log.Error(err))
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

	state.protocolAssets = make(map[uint64]*models.ProtocolAsset)
	state.executors = make(map[uint32]*actor.PID)
	state.dataManagers = make(map[uint64]*actor.PID)

	// Spawn all exchange executors
	for _, protocolStr := range state.Config.Protocols {
		prtcl, ok := constants.GetProtocolByName(protocolStr)
		if !ok {
			return fmt.Errorf("unknown protocol %s", protocolStr)
		}
		producer := NewProtocolExecutorProducer(prtcl, state.registry)
		if producer == nil {
			return fmt.Errorf("unknown protocol %s", prtcl.Name)
		}
		props := actor.PropsFromProducer(producer, actor.WithSupervisor(actor.NewExponentialBackoffStrategy(100*time.Second, time.Second)))

		state.executors[prtcl.ID], _ = context.SpawnNamed(props, prtcl.Name+"_executor")
	}

	// Request securities for each one of them
	var futures []*actor.Future
	request := &messages.ProtocolAssetListRequest{
		RequestID: 0,
		Subscribe: true,
	}
	for _, pid := range state.executors {
		fut := context.RequestFuture(pid, request, 20*time.Second)
		futures = append(futures, fut)
	}

	for _, fut := range futures {
		res, err := fut.Result()
		if err != nil {
			if state.strict {
				return fmt.Errorf("error fetching assets for one protocol: %v", err)
			} else {
				state.logger.Error("error fetching assets for one protocol: %v", log.Error(err))
				continue
			}
		}
		result, ok := res.(*messages.ProtocolAssetList)
		if !ok {
			return fmt.Errorf("was expecting ProtocolAssetList, got %s", reflect.TypeOf(res).String())
		}
		if !result.Success {
			return errors.New(result.RejectionReason.String())
		}

		for _, asset := range result.ProtocolAssets {
			if asset2, ok := state.protocolAssets[asset.ProtocolAssetID]; ok {
				return fmt.Errorf("got two protocol assets with the same ID: %s and %s", asset2.Asset.Symbol, asset.Asset.Symbol)
			}
			state.protocolAssets[asset.ProtocolAssetID] = asset
		}
	}
	return nil
}

func (state *Executor) OnProtocolAssetListRequest(context actor.Context) error {
	req := context.Message().(*messages.ProtocolAssetListRequest)
	assets := make([]*models.ProtocolAsset, 0)
	for _, asset := range state.protocolAssets {
		assets = append(assets, asset)
	}
	response := &messages.ProtocolAssetList{
		RequestID:      req.RequestID,
		ResponseID:     uint64(time.Now().UnixNano()),
		Success:        true,
		ProtocolAssets: assets,
	}
	if req.Subscribe {
		context.Watch(req.Subscriber)
		state.alSubscribers[req.RequestID] = req.Subscriber
	}
	context.Respond(response)
	return nil
}

func (state *Executor) OnProtocolAssetList(context actor.Context) error {
	msg := context.Message().(*messages.ProtocolAssetList)
	var protoId uint32
	if len(msg.ProtocolAssets) > 0 {
		protoId = msg.ProtocolAssets[0].Protocol.ID
	}

	for k, v := range state.protocolAssets {
		if v.Protocol.ID == protoId {
			delete(state.protocolAssets, k)
		}
	}
	for _, asset := range msg.ProtocolAssets {
		state.protocolAssets[asset.ProtocolAssetID] = asset
	}
	var assets []*models.ProtocolAsset
	for _, v := range state.protocolAssets {
		assets = append(assets, v)
	}
	for k, v := range state.alSubscribers {
		context.Send(v,
			&messages.ProtocolAssetList{
				RequestID:      k,
				ResponseID:     uint64(time.Now().UnixNano()),
				ProtocolAssets: assets,
				Success:        true,
			})
	}
	return nil
}

func (state *Executor) OnProtocolAssetDataRequest(context actor.Context) error {
	req := context.Message().(*messages.ProtocolAssetDataRequest)
	a, ok := state.protocolAssets[req.ProtocolAssetID]
	if !ok {
		context.Respond(&messages.ProtocolAssetList{
			RequestID:       req.RequestID,
			RejectionReason: messages.RejectionReason_UnknownProtocolAsset,
			Success:         false,
		})
		return nil
	}
	if pid, ok := state.dataManagers[a.ProtocolAssetID]; ok {
		context.Forward(pid)
	} else {
		props := actor.PropsFromProducer(NewDataManagerProducer(a), actor.WithSupervisor(
			utils.NewExponentialBackoffStrategy(100*time.Second, time.Second, time.Second)))
		pid := context.Spawn(props)
		state.dataManagers[a.ProtocolAssetID] = pid
		context.Forward(pid)
	}
	return nil
}

func (state *Executor) OnHistoricalProtocolAssetTransferRequest(context actor.Context) error {
	req := context.Message().(*messages.HistoricalProtocolAssetTransferRequest)
	a, ok := state.protocolAssets[req.ProtocolAssetID]
	if !ok {
		context.Respond(&messages.ProtocolAssetList{
			RequestID:       req.RequestID,
			RejectionReason: messages.RejectionReason_UnknownProtocolAsset,
			Success:         false,
		})
		return nil
	}
	ex, ok := state.executors[a.Protocol.ID]
	if !ok {
		context.Respond(&messages.HistoricalProtocolAssetTransferResponse{
			RequestID:       req.RequestID,
			RejectionReason: messages.RejectionReason_UnknownProtocol,
			Success:         false,
		})
	}
	context.Forward(ex)
	return nil
}

func (state *Executor) OnProtocolAssetDefinition(context actor.Context) error {
	req := context.Message().(*messages.ProtocolAssetDefinitionRequest)
	a, ok := state.protocolAssets[req.ProtocolAssetID]
	if !ok {
		context.Respond(&messages.ProtocolAssetDefinitionResponse{
			RequestID:       req.RequestID,
			RejectionReason: messages.RejectionReason_UnknownProtocolAsset,
			Success:         false,
		})
		return nil
	}
	context.Respond(&messages.ProtocolAssetDefinitionResponse{
		RequestID:     req.RequestID,
		ResponseID:    uint64(time.Now().UnixNano()),
		Success:       true,
		ProtocolAsset: a,
	})
	return nil
}

func (state *Executor) getProtocolAsset(asset *models.ProtocolAsset) (*models.ProtocolAsset, *messages.RejectionReason) {
	if asset == nil {
		rej := messages.RejectionReason_MissingProtocolAsset
		return nil, &rej
	}
	if asset.Protocol != nil && asset.Asset != nil && asset.Chain != nil {
		ID := utils.GetProtocolAssetID(asset.Asset, asset.Protocol, asset.Chain)
		if a, ok := state.protocolAssets[ID]; !ok {
			rej := messages.RejectionReason_UnknownProtocolAsset
			return nil, &rej
		} else {
			return a, nil
		}
	} else {
		rej := messages.RejectionReason_MissingProtocolAsset
		return nil, &rej
	}
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) OnTerminated(context actor.Context) error {
	// Handle subscriber krash
	req := context.Message().(*actor.Terminated)
	for k, v := range state.alSubscribers {
		if v.Id == req.Who.Id {
			delete(state.alSubscribers, k)
		}
	}

	for k, v := range state.dataManagers {
		if v.Id == req.Who.Id {
			delete(state.dataManagers, k)
		}
	}

	return nil
}
