package bithumbg_test

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
		Symbol:            "BTC-USDT",
		SecurityType:      enum.SecurityType_CRYPTO_SPOT,
		Exchange:          constants.BITHUMBG,
		MinPriceIncrement: 0.01,
		RoundLot:          1e-06,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.InstrumentStatus_Trading,
	})
}
