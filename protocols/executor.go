package protocols

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	models2 "gitlab.com/alphaticks/xchanger/models"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
)

// The executor routes all the request to the underlying exchange executor & listeners
// He is the main part of the whole software..
type ExecutorConfig struct {
	Registry  registry.PublicRegistryClient
	Protocols []*models2.Protocol
}

type Executor struct {
	*ExecutorConfig
	executors          map[uint32]*actor.PID // A map from exchange ID to executor
	protocolAssets     map[uint64]*models.ProtocolAsset
	symToProtocolAsset map[uint32]map[string]*models.ProtocolAsset
	alSubscribers      map[uint64]*actor.PID // A map from request ID to asset list subscriber
	dataManagers       map[uint32]*actor.PID // A map from asset ID to data manager
	logger             *log.Logger
	strict             bool
}

func NewExecutorProducer(cfg *ExecutorConfig) actor.Producer {
	return func() actor.Actor {
		return NewExecutor(cfg)
	}
}

func NewExecutor(cfg *ExecutorConfig) actor.Actor {
	return &Executor{
		ExecutorConfig: cfg,
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
	case *messages.ProtocolAssetTransferRequest:
		if err := state.OnProtocolAssetTransferRequest(context); err != nil {
			state.logger.Error("error processing ProtocolAssetTransferRequest", log.Error(err))
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

	state.symToProtocolAsset = make(map[uint32]map[string]*models.ProtocolAsset)
	state.protocolAssets = make(map[uint64]*models.ProtocolAsset)
	state.executors = make(map[uint32]*actor.PID)
	state.dataManagers = make(map[uint32]*actor.PID)

	// Spawn all exchange executors
	for _, protocol := range state.ExecutorConfig.Protocols {
		producer := NewProtocolExecutorProducer(protocol, state.ExecutorConfig)
		if producer == nil {
			return fmt.Errorf("unknown protocol %s", protocol.Name)
		}
		props := actor.PropsFromProducer(producer).WithSupervisor(actor.NewExponentialBackoffStrategy(100*time.Second, time.Second))

		state.executors[protocol.ID], _ = context.SpawnNamed(props, protocol.Name+"_executor")
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
				return fmt.Errorf("error fetching assets for one venue: %v", err)
			} else {
				state.logger.Error("error fetching assets for one venue: %v", log.Error(err))
			}
		}
		result, ok := res.(*messages.ProtocolAssetListResponse)
		if !ok {
			return fmt.Errorf("was expecting ProtocolAssetListResponse, got %s", reflect.TypeOf(res).String())
		}
		if !result.Success {
			return errors.New(result.RejectionReason.String())
		}

		symToProtocolAsset := make(map[string]*models.ProtocolAsset)
		var protoID uint32
		for _, asset := range result.ProtocolAssets {
			id := uint64(asset.Asset.ID)<<32 + uint64(asset.Protocol.ID)
			if asset2, ok := state.protocolAssets[id]; ok {
				return fmt.Errorf("got two assets with the same ID: %s and %s", asset2.Asset.Symbol, asset.Asset.Symbol)
			}
			state.protocolAssets[id] = asset
			symToProtocolAsset[asset.Asset.Symbol] = asset
			protoID = asset.Protocol.ID
		}
		state.symToProtocolAsset[protoID] = symToProtocolAsset
	}
	return nil
}

func (state *Executor) OnProtocolAssetListRequest(context actor.Context) error {
	req := context.Message().(*messages.ProtocolAssetListRequest)
	assets := make([]*models.ProtocolAsset, 0)
	for _, asset := range state.protocolAssets {
		assets = append(assets, asset)
	}
	response := &messages.ProtocolAssetListResponse{
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
		id := uint64(asset.Asset.ID)<<32 + uint64(asset.Protocol.ID)
		state.protocolAssets[id] = asset
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

func (state *Executor) OnProtocolAssetTransferRequest(context actor.Context) error {
	req := context.Message().(*messages.ProtocolAssetTransferRequest)
	a, rej := state.getProtocolAsset(req.ProtocolAsset)
	if rej != nil {
		context.Respond(&messages.ProtocolAssetListResponse{
			RequestID:       req.RequestID,
			RejectionReason: *rej,
			Success:         false,
		})
		return nil
	}
	if pid, ok := state.dataManagers[a.Protocol.ID]; ok {
		context.Forward(pid)
	} else {
		props := actor.PropsFromProducer(NewDataManagerProducer(a)).WithSupervisor(
			utils.NewExponentialBackoffStrategy(100*time.Second, time.Second, time.Second))
		pid := context.Spawn(props)
		state.dataManagers[a.Protocol.ID] = pid
		context.Forward(pid)
	}
	return nil
}

func (state *Executor) OnHistoricalProtocolAssetTransferRequest(context actor.Context) error {
	req := context.Message().(*messages.HistoricalProtocolAssetTransferRequest)
	a, rej := state.getProtocolAsset(req.ProtocolAsset)
	if rej != nil {
		context.Respond(&messages.HistoricalProtocolAssetTransferResponse{
			RequestID:       req.RequestID,
			RejectionReason: *rej,
			Success:         false,
		})
		return nil
	}
	ex, ok := state.executors[a.Protocol.ID]
	if !ok {
		context.Respond(&messages.HistoricalProtocolAssetTransferResponse{
			RequestID:       req.RequestID,
			RejectionReason: messages.UnknowProtocol,
			Success:         false,
		})
	}
	context.Forward(ex)
	return nil
}

func (state *Executor) OnProtocolAssetDefinition(context actor.Context) error {
	req := context.Message().(*messages.ProtocolAssetDefinitionRequest)
	id := req.ProtocolAssetID
	assetID := id >> 32
	protocolID := id - assetID<<32
	asset := &models.ProtocolAsset{
		Protocol: &models2.Protocol{
			ID: uint32(protocolID),
		},
		Asset: &models2.Asset{
			ID: uint32(assetID),
		},
	}
	a, rej := state.getProtocolAsset(asset)
	if rej != nil {
		context.Respond(&messages.ProtocolAssetDefinitionResponse{
			RequestID:       req.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: *rej,
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
		rej := messages.MissingProtocolAsset
		return nil, &rej
	}
	if asset.Protocol != nil && asset.Asset != nil {
		id := uint64(asset.Asset.ID)<<32 + uint64(asset.Protocol.ID)
		if a, ok := state.protocolAssets[id]; !ok {
			rej := messages.UnknownProtocolAsset
			return nil, &rej
		} else {
			return a, nil
		}
	} else if asset.Protocol != nil {
		p, ok := state.symToProtocolAsset[asset.Protocol.ID]
		if !ok {
			rej := messages.UnknowProtocol
			return nil, &rej
		}
		a, ok := p[asset.Asset.Symbol]
		if !ok {
			rej := messages.UnknownSymbol
			return nil, &rej
		}
		return a, nil
	} else {
		rej := messages.MissingProtocolAsset
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
