package bithumbg_test

import (
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
)

func TestExecutorPublic(t *testing.T) {
	tests.ExPub(t, tests.ExPubTest{
		Instrument: &models.Instrument{
			Exchange: constants.BITHUMBG,
			Symbol:   &wrapperspb.StringValue{Value: "BTC-USDT"},
		},
		SecurityListRequest: true,
		MarketDataRequest:   true,
	})
}
