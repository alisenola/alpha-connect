package okex_test

import (
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
)

func TestExecutorPublic(t *testing.T) {
	tests.ExPub(t, tests.ExPubTest{
		Instrument: &models.Instrument{
			Exchange: constants.OKEX,
			Symbol:   &wrapperspb.StringValue{Value: "BTC-USDT"},
		},
		SecurityListRequest: true,
		MarketDataRequest:   true,
	})
}
