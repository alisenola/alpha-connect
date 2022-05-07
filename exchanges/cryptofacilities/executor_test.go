package cryptofacilities_test

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
			Exchange: constants.CRYPTOFACILITIES,
			Symbol:   &wrapperspb.StringValue{Value: "pi_xbtusd"},
		},
		SecurityListRequest: true,
		MarketDataRequest:   true,
	})
}
