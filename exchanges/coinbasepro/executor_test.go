package coinbasepro_test

import (
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	exTests "gitlab.com/alphaticks/alpha-connect/tests"
	"gitlab.com/alphaticks/xchanger/constants"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
)

func TestExecutorPublic(t *testing.T) {
	exTests.LoadStatics(t)
	tests.ExPub(t, tests.ExPubTest{
		Instrument: &models.Instrument{
			Exchange: constants.COINBASEPRO,
			Symbol:   &wrapperspb.StringValue{Value: "BTC-USD"},
		},
		SecurityListRequest: true,
		MarketDataRequest:   true,
	})
}
