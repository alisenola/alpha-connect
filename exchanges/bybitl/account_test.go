package bybitl_test

import (
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/bybitl"
	m "gitlab.com/alphaticks/xchanger/models"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
)

var instrument = &models.Instrument{
	Exchange:   constants.BYBITL,
	Symbol:     wrapperspb.String("BTCUSDT"),
	SecurityID: wrapperspb.UInt64(6789757764526280996),
}

var bybitlAccount = &models.Account{
	Exchange: constants.BYBITL,
	ApiCredentials: &m.APICredentials{
		APIKey:    "21g5YOmxipsME0lHUr",
		APISecret: "x8D8Qc81hZ1kCGNqH1W6spl1gp1VrFbfkAuE",
	},
}

func TestNewAccountListener(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	bybitl.EnableTestNet()
	bybitl.EnableWebSocketTestNet()
	tests.AccntTest(t, tests.AccountTest{
		Account:                bybitlAccount,
		Instrument:             instrument,
		ExpiredOrder:           true,
		GetPositionsLimit:      false,
		GetPositionsMarket:     false,
		OrderReplaceRequest:    false,
		OrderMassCancelRequest: false,
	})
}
