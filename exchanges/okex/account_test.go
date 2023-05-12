package okex_test

import (
	"testing"

	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/okex"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var instrument = &models.Instrument{
	Exchange: constants.OKEX,
	Symbol:   wrapperspb.String("BTC-USDT"),
}

var okexAccount = &models.Account{
	Exchange: constants.OKEX,
	ApiCredentials: &xchangerModels.APICredentials{
		APIKey:    "0187a1c3-36e8-41bd-b84e-90b7347372c9",
		APISecret: "6F3CEECAD897A0E243832C49D62EB7A2",
		AccountID: "UStQhSVR7vK7h7g@",
	},
}

func TestAccountListener(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	okex.EnableTestnet()
	tests.AccntTest(t, tests.AccountTest{
		Account:            okexAccount,
		Instrument:         instrument,
		OrderStatusRequest: true,
		GetPositionsLimit:  true,
		GetPositionsMarket: true,
	})
}
