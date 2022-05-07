package tests

import (
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"

	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/protocols"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	models2 "gitlab.com/alphaticks/xchanger/models"
	"google.golang.org/grpc"
)

func StartExecutor(t *testing.T, protocol *models2.Protocol) (*actor.ActorSystem, *actor.PID, *actor.PID, func()) {
	registryAddress := "127.0.0.1:7001"
	if os.Getenv("REGISTRY_ADDRESS") != "" {
		registryAddress = os.Getenv("REGISTRY_ADDRESS")
	}
	conn, err := grpc.Dial(registryAddress, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	reg := registry.NewPublicRegistryClient(conn)

	cfg := protocols.ExecutorConfig{
		Registry:  reg,
		Protocols: []*models2.Protocol{protocol},
	}
	config := actor.Config{}
	config = config.WithDeveloperSupervisionLogging(true).WithDeadLetterRequestLogging(true)
	as := actor.NewActorSystemWithConfig(config)
	ex, err := as.Root.SpawnNamed(actor.PropsFromProducer(protocols.NewExecutorProducer(&cfg)), "executor")
	if err != nil {
		t.Fatal(err)
	}
	loader := as.Root.Spawn(actor.PropsFromProducer(utils.NewStaticLoaderProducer(reg)))

	return as, ex, loader, func() { as.Root.PoisonFuture(ex).Wait() }
}

type ERC721Checker struct {
	asset   *models.ProtocolAsset
	logger  *log.Logger
	updates []*models.ProtocolAssetUpdate
	seqNum  uint64
}

func NewERC721CheckerProducer(asset *models.ProtocolAsset) actor.Producer {
	return func() actor.Actor {
		return NewERC721Checker(asset)
	}
}

func NewERC721Checker(asset *models.ProtocolAsset) actor.Actor {
	return &ERC721Checker{
		asset:   asset,
		updates: nil,
		logger:  nil,
		seqNum:  0,
	}
}

func (state *ERC721Checker) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.logger.Error("error starting the actor", log.Error(err))
			panic(err)
		}
	case *messages.ProtocolAssetDataIncrementalRefresh:
		if err := state.OnProtocolAssetDataIncrementalRefresh(context); err != nil {
			state.logger.Error("error processing ProtocolAssetDataIncrementalRefresh", log.Error(err))
			panic(err)
		}
	}
}

func (state *ERC721Checker) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(state).String()),
	)
	executor := context.ActorSystem().NewLocalPID("executor")
	res, err := context.RequestFuture(executor, &messages.ProtocolAssetTransferRequest{
		RequestID:       0,
		ProtocolAssetID: state.asset.ProtocolAssetID,
		Subscriber:      context.Self(),
		Subscribe:       true,
	}, 30*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error fetching the asset transfer %v", err)
	}
	updt, ok := res.(*messages.ProtocolAssetTransferResponse)
	if !ok {
		return fmt.Errorf("error for type assertion %v", err)
	}
	state.updates = updt.Update
	return nil
}

func (state *ERC721Checker) OnProtocolAssetDataIncrementalRefresh(context actor.Context) error {
	res := context.Message().(*messages.ProtocolAssetDataIncrementalRefresh)
	if res.Update == nil {
		return nil
	}
	if state.seqNum > res.SeqNum {
		return nil
	}
	state.seqNum = res.SeqNum
	state.updates = append(state.updates, res.Update)
	fmt.Println(state.updates)
	return nil
}
