package bybitl_test

import (
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/bybitl"
	m "gitlab.com/alphaticks/xchanger/models"
	"testing"
)

var instrument = &models.Instrument{
	Exchange: &constants.BYBITL,
	Symbol:   &types.StringValue{Value: "BTCUSDT"},
}

var bybitlAccount = &models.Account{
	Exchange: &constants.BYBITL,
	ApiCredentials: &m.APICredentials{
		APIKey:    "21g5YOmxipsME0lHUr",
		APISecret: "x8D8Qc81hZ1kCGNqH1W6spl1gp1VrFbfkAuE",
	},
}

func TestNewAccountListener(t *testing.T) {
	bybitl.EnableTestNet()
	bybitl.EnableWebSocketTestNet()
	tests.AccntTest(t, tests.AccountTest{
		Account:                bybitlAccount,
		Instrument:             instrument,
		OrderStatusRequest:     false,
		GetPositionsLimit:      true,
		GetPositionsMarket:     false,
		OrderReplaceRequest:    false,
		OrderMassCancelRequest: false,
	})
}
