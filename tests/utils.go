package tests

import (
	"encoding/json"
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/account"
	chtypes "gitlab.com/alphaticks/alpha-connect/chains/types"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/executor"
	"gitlab.com/alphaticks/alpha-connect/models"
	prtypes "gitlab.com/alphaticks/alpha-connect/protocols/types"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func LoadStatics(t *testing.T) {
	pwd, _ := os.Getwd()
	for path.Base(pwd) != "alpha-connect" {
		pwd = path.Dir(pwd)
	}
	b, err := ioutil.ReadFile(pwd + "/tests/static/assets.json")
	if err != nil {
		t.Fatalf("failed to load assets: %v", err)
	}
	assets := make(map[uint32]*xchangerModels.Asset)
	err = json.Unmarshal(b, &assets)
	if err != nil {
		t.Fatalf("error parsing assets: %v", err)
	}
	if err := constants.LoadAssets(assets); err != nil {
		t.Fatalf("error loading assets: %v", err)
	}

	b, err = ioutil.ReadFile(pwd + "/tests/static/protocols.json")
	if err != nil {
		t.Fatalf("failed to load protocols: %v", err)
	}
	prtcls := make(map[uint32]*xchangerModels.Protocol)
	err = json.Unmarshal(b, &prtcls)
	if err != nil {
		t.Fatalf("error parsing protocols: %v", err)
	}
	if err := constants.LoadProtocols(prtcls); err != nil {
		t.Fatalf("error loading protocols: %v", err)
	}

	b, err = ioutil.ReadFile(pwd + "/tests/static/chains.json")
	if err != nil {
		t.Fatalf("failed to load chains: %v", err)
	}
	chns := make(map[uint32]*xchangerModels.Chain)
	err = json.Unmarshal(b, &chns)
	if err != nil {
		t.Fatalf("error parsing chains: %v", err)
	}
	if err := constants.LoadChains(chns); err != nil {
		t.Fatalf("error loading chains: %v", err)
	}
}

func StartExecutor(t *testing.T, cfgEx *extypes.ExecutorConfig, cfgPr *prtypes.ExecutorConfig, cfgCh *chtypes.ExecutorConfig, acc *models.Account) (*actor.ActorSystem, *actor.PID, func()) {
	LoadStatics(t)
	if cfgEx != nil && acc != nil {
		accnt, err := exchanges.NewAccount(acc)
		if err != nil {
			t.Fatal(err)
		}
		cfgEx.Accounts = []*account.Account{accnt}
	}
	config := actor.Config{DeadLetterThrottleInterval: time.Second, DeveloperSupervisionLogging: true, DeadLetterRequestLogging: true}
	as := actor.NewActorSystemWithConfig(&config)
	exec, _ := as.Root.SpawnNamed(actor.PropsFromProducer(executor.NewExecutorProducer(cfgEx, cfgPr, cfgCh)), "executor")
	return as, exec, func() { _ = as.Root.PoisonFuture(exec).Wait() }
}
