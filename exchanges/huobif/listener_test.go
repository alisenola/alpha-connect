package huobif_test

import (
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	exTests "gitlab.com/alphaticks/alpha-connect/tests"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
)

func TestMarketData(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	exTests.LoadStatics(t)
	tests.MarketData(t, tests.MDTest{
		Symbol:            "BTC210625",
		SecurityType:      enum.SecurityType_CRYPTO_FUT,
		Exchange:          constants.HUOBIF,
		MinPriceIncrement: 0.01,
		RoundLot:          1.,
		HasMaturityDate:   true,
		IsInverse:         false,
		Status:            models.InstrumentStatus_Trading,
	})
}
