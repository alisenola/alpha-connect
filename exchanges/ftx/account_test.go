package ftx_test

import (
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
)

var instrument = &models.Instrument{
	Exchange: constants.FTX,
	Symbol:   &wrapperspb.StringValue{Value: "BTC-PERP"},
}

var ftxAccount = &models.Account{
	Name:     "299211",
	Exchange: constants.FTX,
	ApiCredentials: &xchangerModels.APICredentials{
		AccountID: "integration-test",
		APIKey:    "9VSHd0TKq5i5pRif3rJN1oB3fi2pUaNqweRSYMFZ",
		APISecret: "R3ImbLQWsY3oJ5dXMogDP-0VmLT0BoOaTPwttmTR",
	},
}

func TestAccountListener(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	tests.AccntTest(t, tests.AccountTest{
		Account:    ftxAccount,
		Instrument: instrument,
		//OrderStatusRequest:      true,
		//OrderReplaceRequest:     true,
		//OrderBulkReplaceRequest: true,
		GetPositionsLimit: true,
		SkipCheckBalance:  true,
		//GetPositionsMarket:      true,
	})
}
