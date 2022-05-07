package bitmex_test

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
			Exchange: constants.BITMEX,
			Symbol:   &wrapperspb.StringValue{Value: "XBTUSD"},
		},
		SecurityListRequest: true,
		MarketDataRequest:   true,
	})
}
