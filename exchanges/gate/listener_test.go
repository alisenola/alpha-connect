package gate_test

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
		IgnoreSizeResidue: true,
		Symbol:            "BTC_USDT",
		SecurityType:      enum.SecurityType_CRYPTO_SPOT,
		Exchange:          constants.GATE,
		MinPriceIncrement: 0.01,
		RoundLot:          1e-04,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.InstrumentStatus_Trading,
	})
}
