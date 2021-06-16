package account

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/account"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/bitmex"
	"gitlab.com/alphaticks/xchanger/exchanges/fbinance"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"os"
	"testing"
)

var executor *actor.PID

var accounts = []*models.Account{
	BitmexAccount,
	FBinanceTestnetAccount,
}

var As *actor.ActorSystem

func TestMain(m *testing.M) {
	exch := []*xchangerModels.Exchange{
		&constants.FBINANCE,
		&constants.BITMEX,
	}
	accnts := []*account.Account{}
	for _, a := range accounts {
		aa, err := exchanges.NewAccount(a)
		if err != nil {
			panic(err)
		}
		accnts = append(accnts, aa)
	}

	As = actor.NewActorSystem()

	executor, _ = As.Root.SpawnNamed(actor.PropsFromProducer(exchanges.NewExecutorProducer(exch, accnts, false, xchangerUtils.DefaultDialerPool)), "executor")
	bitmex.EnableTestNet()
	fbinance.EnableTestNet()

	code := m.Run()

	bitmex.DisableTestNet()
	fbinance.DisableTestNet()

	As.Root.PoisonFuture(executor)
	os.Exit(code)
}
