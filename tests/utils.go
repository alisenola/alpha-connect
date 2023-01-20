package tests

import (
	"encoding/json"
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/config"
	"gitlab.com/alphaticks/alpha-connect/executor"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"testing"
	"time"
)

func LoadStatics(t *testing.T) {
	assets := make(map[uint32]*xchangerModels.Asset)
	err := json.Unmarshal(AssetsJSON, &assets)
	if err != nil {
		t.Fatalf("error parsing assets: %v", err)
	}
	if err := constants.LoadAssets(assets); err != nil {
		t.Fatalf("error loading assets: %v", err)
	}

	prtcls := make(map[uint32]*xchangerModels.Protocol)
	err = json.Unmarshal(ProtocolsJSON, &prtcls)
	if err != nil {
		t.Fatalf("error parsing protocols: %v", err)
	}
	if err := constants.LoadProtocols(prtcls); err != nil {
		t.Fatalf("error loading protocols: %v", err)
	}

	chns := make(map[uint32]*xchangerModels.Chain)
	err = json.Unmarshal(ChainsJSON, &chns)
	if err != nil {
		t.Fatalf("error parsing chains: %v", err)
	}
	if err := constants.LoadChains(chns); err != nil {
		t.Fatalf("error loading chains: %v", err)
	}
}

func StartExecutor(t *testing.T, cfg *config.Config) (*actor.ActorSystem, *actor.PID, func()) {
	LoadStatics(t)
	config := actor.Config{DeadLetterThrottleInterval: time.Second, DeveloperSupervisionLogging: true, DeadLetterRequestLogging: true}
	as := actor.NewActorSystemWithConfig(&config)
	exec, _ := as.Root.SpawnNamed(actor.PropsFromProducer(executor.NewExecutorProducer(cfg)), "executor")
	return as, exec, func() { _ = as.Root.PoisonFuture(exec).Wait() }
}
