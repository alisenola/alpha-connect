package bitmex_test

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
		Symbol:            "XBTUSD",
		SecurityType:      enum.SecurityType_CRYPTO_PERP,
		Exchange:          constants.BITMEX,
		MinPriceIncrement: 0.5,
		RoundLot:          100.,
		HasMaturityDate:   false,
		IsInverse:         true,
		Status:            models.InstrumentStatus_Trading,
	})
}
