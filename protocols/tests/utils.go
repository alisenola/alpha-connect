package tests

import (
	"os"
	"testing"

	"github.com/AsynkronIT/protoactor-go/actor"
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
	as := actor.NewActorSystem()
	ex, err := as.Root.SpawnNamed(actor.PropsFromProducer(protocols.NewExecutorProducer(&cfg)), "executor")
	if err != nil {
		t.Fatal(err)
	}

	return as, ex, func() { as.Root.PoisonFuture(ex).Wait() }
}
