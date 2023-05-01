package okex_test

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
	exTests.LoadStatics()
	tests.MarketData(t, tests.MDTest{
		Symbol:            "BTC-USDT",
		SecurityType:      enum.SecurityType_CRYPTO_SPOT,
		Exchange:          constants.OKEX,
		MinPriceIncrement: 0.1,
		RoundLot:          1e-08,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.InstrumentStatus_Trading,
	})
}
