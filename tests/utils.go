package tests

import (
	"encoding/json"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/config"
	"gitlab.com/alphaticks/alpha-connect/executor"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"testing"
	"time"
)

func LoadStatics() error {
	assets := make(map[uint32]*xchangerModels.Asset)
	err := json.Unmarshal(AssetsJSON, &assets)
	if err != nil {
		return fmt.Errorf("error parsing assets: %v", err)
	}
	if err := constants.LoadAssets(assets); err != nil {
		return fmt.Errorf("error loading assets: %v", err)
	}

	prtcls := make(map[uint32]*xchangerModels.Protocol)
	err = json.Unmarshal(ProtocolsJSON, &prtcls)
	if err != nil {
		return fmt.Errorf("error parsing protocols: %v", err)
	}
	if err := constants.LoadProtocols(prtcls); err != nil {
		return fmt.Errorf("error loading protocols: %v", err)
	}

	chns := make(map[uint32]*xchangerModels.Chain)
	err = json.Unmarshal(ChainsJSON, &chns)
	if err != nil {
		return fmt.Errorf("error parsing chains: %v", err)
	}
	if err := constants.LoadChains(chns); err != nil {
		return fmt.Errorf("error loading chains: %v", err)
	}
	return nil
}

func StartExecutor(t *testing.T, cfg *config.Config) (*actor.ActorSystem, *actor.PID, func()) {
	if err := LoadStatics(); err != nil {
		t.Fatalf("error loading statics: %v", err)
	}
	aconfig := actor.Config{DeadLetterThrottleInterval: time.Second, DeveloperSupervisionLogging: true, DeadLetterRequestLogging: true}
	as := actor.NewActorSystemWithConfig(&aconfig)
	exec, _ := as.Root.SpawnNamed(actor.PropsFromProducer(executor.NewExecutorProducer(cfg)), "executor")
	return as, exec, func() { _ = as.Root.PoisonFuture(exec).Wait() }
}
