package fbinance_test

import (
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/fbinance"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
)

var instrument = &models.Instrument{
	SecurityID: &wrapperspb.UInt64Value{Value: 5485975358912730733},
	Exchange:   constants.FBINANCE,
	Symbol:     &wrapperspb.StringValue{Value: "BTCUSDT"},
}

var fbinanceAccount = &models.Account{
	Exchange: constants.FBINANCE,
	ApiCredentials: &xchangerModels.APICredentials{
		APIKey:    "584ed906986e68940c1cdeba22be16cbe0a459678dcb58e4761c2899ce8a82af",
		APISecret: "05fc33219ab552d416d8506b24276d9930af18cbbe93c0a45df23d4e11eaf3b8",
	},
}

func TestAccountListener(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	fbinance.EnableTestNet()
	tests.AccntTest(t, tests.AccountTest{
		Account:             fbinanceAccount,
		Instrument:          instrument,
		OrderStatusRequest:  true,
		GetPositionsLimit:   false,
		GetPositionsMarket:  false,
		OrderReplaceRequest: false,
	})
}
