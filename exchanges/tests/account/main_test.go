package account

import (
	"fmt"
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
	//BitmexAccount,
	//FBinanceTestnetAccount,
	FTXAccount,
}

var As *actor.ActorSystem

func TestMain(m *testing.M) {
	assets := map[uint32]xchangerModels.Asset{
		constants.DOLLAR.ID:           constants.DOLLAR,
		constants.EURO.ID:             constants.EURO,
		constants.POUND.ID:            constants.POUND,
		constants.CANADIAN_DOLLAR.ID:  constants.CANADIAN_DOLLAR,
		constants.JAPENESE_YEN.ID:     constants.JAPENESE_YEN,
		constants.BITCOIN.ID:          constants.BITCOIN,
		constants.LITECOIN.ID:         constants.LITECOIN,
		constants.ETHEREUM.ID:         constants.ETHEREUM,
		constants.RIPPLE.ID:           constants.RIPPLE,
		constants.TETHER.ID:           constants.TETHER,
		constants.SOUTH_KOREAN_WON.ID: constants.SOUTH_KOREAN_WON,
		constants.USDC.ID:             constants.USDC,
		constants.DASH.ID:             constants.DASH,
	}

	fmt.Println("LOAD ASSETS")
	if err := constants.LoadAssets(assets); err != nil {
		panic(err)
	}
	exch := []*xchangerModels.Exchange{
		//&constants.FBINANCE,
		//&constants.BITMEX,
		&constants.FTX,
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

	executor, _ = As.Root.SpawnNamed(actor.PropsFromProducer(exchanges.NewExecutorProducer(exch, accnts, xchangerUtils.DefaultDialerPool)), "executor")
	bitmex.EnableTestNet()
	fbinance.EnableTestNet()

	code := m.Run()

	bitmex.DisableTestNet()
	fbinance.DisableTestNet()

	As.Root.PoisonFuture(executor)
	os.Exit(code)
}
