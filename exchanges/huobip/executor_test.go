package huobip_test

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
			Exchange: constants.HUOBIP,
			Symbol:   &wrapperspb.StringValue{Value: "BTC-USD"},
		},
		SecurityListRequest: true,
		MarketDataRequest:   true,
		OpenInterestRequest: true,
	})
}
