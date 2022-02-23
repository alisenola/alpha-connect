package tests

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"

	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/protocols"
	registry "gitlab.com/alphaticks/alpha-registry-grpc"
	models2 "gitlab.com/alphaticks/xchanger/models"
	"google.golang.org/grpc"
)

func StartExecutor(t *testing.T, protocol *models2.Protocol) (*actor.ActorSystem, *actor.PID, func()) {
	registryAddress := "registry.alphaticks.io:7001"
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

	return as, ex, func() { as.Root.PoisonFuture(ex).Wait() }
}

type ERC721Checker struct {
	asset   *models.ProtocolAsset
	logger  *log.Logger
	updates []*models.AssetUpdate
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
	case *messages.AssetDataIncrementalRefresh:
		if err := state.OnAssetDataIncrementalRefresh(context); err != nil {
			state.logger.Error("error processing AssetDataIncrementalRefresh", log.Error(err))
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
	res, err := context.RequestFuture(executor, &messages.AssetTransferRequest{
		RequestID:  0,
		Asset:      state.asset,
		Subscriber: context.Self(),
		Subscribe:  true,
	}, 30*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error fetching the asset transfer %v", err)
	}
	updt, ok := res.(*messages.AssetTransferResponse)
	if !ok {
		return fmt.Errorf("error for type assertion %v", err)
	}
	state.updates = updt.AssetUpdated
	return nil
}

func (state *ERC721Checker) OnAssetDataIncrementalRefresh(context actor.Context) error {
	res := context.Message().(*messages.AssetDataIncrementalRefresh)
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
