package protocols

import (
	goContext "context"
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
	registry "gitlab.com/alphaticks/alpha-registry-grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// The executor routes all the request to the underlying exchange executor & listeners
// He is the main part of the whole software..
type ExecutorConfig struct {
	Db        *mongo.Database
	Registry  registry.PublicRegistryClient
	Strict    bool
	Protocols []*models2.Protocol
}

type Executor struct {
	*ExecutorConfig
	executors     map[uint32]*actor.PID // A map from exchange ID to executor
	assets        map[[20]byte]*models.ProtocolAsset
	symToAsset    map[uint32]map[string]*models.ProtocolAsset
	alSubscribers map[uint64]*actor.PID // A map from request ID to asset list subscriber
	dataManagers  map[uint32]*actor.PID // A map from asset ID to data manager
	logger        *log.Logger
	strict        bool
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

	case *actor.Terminated:
		if err := state.OnTerminated(context); err != nil {
			state.logger.Error("error processing OnTerminated", log.Error(err))
			panic(err)
		}
	case *messages.AssetListRequest:
		if err := state.OnAssetListRequest(context); err != nil {
			state.logger.Error("error processing AssetListRequest", log.Error(err))
			panic(err)
		}
	case *messages.AssetList:
		if err := state.OnAssetList(context); err != nil {
			state.logger.Error("error processing AssetList", log.Error(err))
			panic(err)
		}
	case *messages.HistoricalAssetTransferRequest:
		if err := state.OnHistoricalAssetTransferRequest(context); err != nil {
			state.logger.Error("error processing HistoricalAssetTransferRequest", log.Error(err))
			panic(err)
		}
	case *messages.AssetTransferRequest:
		if err := state.OnAssetTransferRequest(context); err != nil {
			state.logger.Error("error processing AssetTransferRequest", log.Error(err))
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

	if state.Db != nil {
		unique := true
		mod := mongo.IndexModel{
			Keys: bson.M{
				"id": 1, // index in ascending order
			}, Options: &options.IndexOptions{Unique: &unique},
		}
		txs := state.Db.Collection("transactions")
		execs := state.Db.Collection("executions")
		if _, err := txs.Indexes().CreateOne(goContext.Background(), mod); err != nil {
			return fmt.Errorf("error creating index on transactions: %v", err)
		}
		if _, err := execs.Indexes().CreateOne(goContext.Background(), mod); err != nil {
			return fmt.Errorf("error creating index on executions: %v", err)
		}
	}
	state.symToAsset = make(map[uint32]map[string]*models.ProtocolAsset)
	state.assets = make(map[[20]byte]*models.ProtocolAsset)
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
	request := &messages.AssetListRequest{
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
		result, ok := res.(*messages.AssetListResponse)
		if !ok {
			return fmt.Errorf("was expecting AssetListResponse, got %s", reflect.TypeOf(res).String())
		}
		if !result.Success {
			return errors.New(result.RejectionReason.String())
		}

		symToAsset := make(map[string]*models.ProtocolAsset)
		var protoID uint32
		for _, asset := range result.Assets {
			var add [20]byte
			copy(add[:], asset.Address)
			if asset2, ok := state.assets[add]; ok {
				return fmt.Errorf("got two assets with the same contract address: %s and %s", asset2.Symbol, asset.Symbol)
			}
			state.assets[add] = asset
			symToAsset[asset.Symbol] = asset
			protoID = asset.Protocol.ID
		}
		state.symToAsset[protoID] = symToAsset
	}
	return nil
}

func (state *Executor) OnAssetListRequest(context actor.Context) error {
	req := context.Message().(*messages.AssetListRequest)
	assets := make([]*models.ProtocolAsset, 0)
	for _, asset := range state.assets {
		assets = append(assets, asset)
	}
	response := &messages.AssetListResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Assets:     assets,
	}
	if req.Subscribe {
		context.Watch(req.Subscriber)
		state.alSubscribers[req.RequestID] = req.Subscriber
	}
	context.Respond(response)
	return nil
}

func (state *Executor) OnAssetList(context actor.Context) error {
	msg := context.Message().(*messages.AssetList)
	proto := msg.Assets[0].Protocol.ID

	for k, v := range state.assets {
		if v.Protocol.ID == proto {
			delete(state.assets, k)
		}
	}
	for _, asset := range msg.Assets {
		var add [20]byte
		copy(add[:], asset.Address)
		state.assets[add] = asset
	}
	var assets []*models.ProtocolAsset
	for _, v := range state.assets {
		assets = append(assets, v)
	}
	for k, v := range state.alSubscribers {
		context.Send(v,
			&messages.AssetList{
				RequestID:  k,
				ResponseID: uint64(time.Now().UnixNano()),
				Assets:     assets,
				Success:    true,
			})
	}
	return nil
}

func (state *Executor) OnAssetTransferRequest(context actor.Context) error {
	req := context.Message().(*messages.AssetTransferRequest)
	a, rej := state.getAsset(req.Asset)
	if rej != nil {
		context.Respond(&messages.AssetListResponse{
			RequestID:       req.RequestID,
			RejectionReason: *rej,
			Success:         false,
		})
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

func (state *Executor) OnHistoricalAssetTransferRequest(context actor.Context) error {
	req := context.Message().(*messages.HistoricalAssetTransferRequest)
	asset, rej := state.getAsset(req.Asset)
	if rej != nil {
		context.Respond(&messages.HistoricalAssetTransferResponse{
			RequestID:       req.RequestID,
			RejectionReason: *rej,
			Success:         false,
		})
	}
	ex, ok := state.executors[asset.Protocol.ID]
	if !ok {
		context.Respond(&messages.HistoricalAssetTransferResponse{
			RequestID:       req.RequestID,
			RejectionReason: messages.UnknowProtocol,
			Success:         false,
		})
	}
	context.Forward(ex)
	return nil
}

func (state *Executor) getAsset(asset *models.ProtocolAsset) (*models.ProtocolAsset, *messages.RejectionReason) {
	if asset == nil {
		rej := messages.MissingAsset
		return nil, &rej
	}
	if len(asset.Address) == 20 {
		var add [20]byte
		copy(add[:], asset.Address)
		if a, ok := state.assets[add]; !ok {
			rej := messages.UnknownAsset
			return nil, &rej
		} else {
			return a, nil
		}
	} else if asset.Protocol != nil {
		p, ok := state.symToAsset[asset.Protocol.ID]
		if !ok {
			rej := messages.UnknowProtocol
			return nil, &rej
		}
		a, ok := p[asset.Symbol]
		if !ok {
			rej := messages.UnknownSymbol
			return nil, &rej
		}
		return a, nil
	} else {
		rej := messages.MissingAsset
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
