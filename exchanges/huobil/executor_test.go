package huobil_test

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
			Exchange: &constants.HUOBIL,
			Symbol:   &types.StringValue{Value: "BTC-USDT"},
		},
		SecurityListRequest: true,
		MarketDataRequest:   false,
		OpenInterestRequest: true,
	})
}