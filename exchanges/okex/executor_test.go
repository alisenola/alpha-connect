package okex_test

import (
	"gitlab.com/alphaticks/xchanger/exchanges/okex"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"testing"

	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestExecutorPublic(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	okex.EnableTestnet()
	tests.ExPub(t, tests.ExPubTest{
		Instrument: &models.Instrument{
			Exchange: constants.OKEX,
			Symbol:   &wrapperspb.StringValue{Value: "BTC-USDT"},
		},
		Account: &models.Account{
			Exchange: constants.OKEX,
			ApiCredentials: &xchangerModels.APICredentials{
				APIKey:    "0187a1c3-36e8-41bd-b84e-90b7347372c9",
				APISecret: "6F3CEECAD897A0E243832C49D62EB7A2",
				AccountID: "UStQhSVR7vK7h7g@",
			},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: true,
		HistoricalFundingRateRequest:  true,
		BalancesRequest:               true,
	})
}
