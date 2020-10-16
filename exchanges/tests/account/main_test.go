package account

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alphac/account"
	"gitlab.com/alphaticks/alphac/exchanges"
	"gitlab.com/alphaticks/alphac/models"
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
	FBinanceAccount,
}

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
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(exchanges.NewExecutorProducer(exch, accnts, false, xchangerUtils.DefaultDialerPool)), "executor")
	bitmex.EnableTestNet()
	fbinance.EnableTestNet()

	code := m.Run()

	bitmex.DisableTestNet()
	fbinance.DisableTestNet()

	actor.EmptyRootContext.PoisonFuture(executor)
	os.Exit(code)
}
