package okex_test

import (
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
	tests.ExPub(t, tests.ExPubTest{
		Instrument: &models.Instrument{
			Exchange: constants.OKEX,
			Symbol:   &wrapperspb.StringValue{Value: "BTC-USDT"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: true,
		HistoricalFundingRateRequest:  true,
	})
}
