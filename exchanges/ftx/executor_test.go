package ftx_test

import (
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
)

func TestExecutorPublic(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	tests.ExPub(t, tests.ExPubTest{
		Instrument: &models.Instrument{
			Exchange: constants.FTX,
			Symbol:   &wrapperspb.StringValue{Value: "ZIL-PERP"},
		},
		SecurityListRequest: true,
		MarketDataRequest:   true,
		OpenInterestRequest: true,
	})
}
