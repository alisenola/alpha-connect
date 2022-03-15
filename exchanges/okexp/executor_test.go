package okexp_test

import (
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
)

func TestExecutorPublic(t *testing.T) {
	tests.ExPub(t, tests.ExPubTest{
		Instrument: &models.Instrument{
			Exchange: &constants.OKEXP,
			Symbol:   &types.StringValue{Value: "BTC-USD-SWAP"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: true,
		OpenInterestRequest:           true,
	})
}
