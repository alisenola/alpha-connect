package exchanges_test

import (
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/config"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	"gitlab.com/alphaticks/alpha-connect/tests"
	tsc_config "gitlab.com/alphaticks/tickstore-go-client/config"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
	"time"
)

func TestAccountMonitor(t *testing.T) {
	var C = &config.Config{
		RegistryAddress: "registry.alphaticks.io:8021",
		Exchanges:       []string{"fbinance"},
		Accounts: []config.Account{{
			Exchange:  constants.FBINANCE.Name,
			Name:      "test",
			ApiKey:    "RezH1h4E6naaKhZ83dEVGkbiFWZssOOGFlHxfkXSNcOAvoFpd7uCg99zGprxN2Np",
			ApiSecret: "jvPpwMEEMlewlFwCYkcwHCBTU5e3VrCelMeMH3JI87MPRN6Mak55GshDENEhQ5Px",
			Listen:    true,
		}},
	}
	as, _, cleaner := tests.StartExecutor(t, C)
	defer cleaner()

	time.Sleep(1 * time.Second)
	// Spawn an actor
	monitor := exchanges.NewAccountMonitor(config.Account{
		Exchange: constants.FBINANCE.Name,
		Name:     "test",
	}, tsc_config.StoreClient{})

	as.Root.Spawn(actor.PropsFromFunc(func(context actor.Context) {
		monitor.Receive(context)
	}))

	time.Sleep(100 * time.Second)
}
