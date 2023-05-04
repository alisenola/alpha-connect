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
		APIKey:    "0b7927eb-5dd7-4dd0-bc34-6dd576a007c5",
		APISecret: "35227C3FEF945545DA45F419F100BE61",
		AccountID: "npe9vnf*efp@drk!CRG",
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
		//GetPositionsLimit:   true,
		//GetPositionsMarket:  true,
		//OrderReplaceRequest: true,
	})
}
