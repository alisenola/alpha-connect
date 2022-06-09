package huobi_test

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
		Symbol:            "btcusdt",
		SecurityType:      enum.SecurityType_CRYPTO_SPOT,
		Exchange:          constants.HUOBI,
		MinPriceIncrement: 0.01,
		RoundLot:          1e-6,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.InstrumentStatus_Trading,
		IgnoreSizeResidue: true,
	})
}

func TestSlowMarketData(t *testing.T) {
	exTests.LoadStatics(t)
	tests.MarketData(t, tests.MDTest{
		Symbol:            "kaneth",
		SecurityType:      enum.SecurityType_CRYPTO_SPOT,
		Exchange:          constants.HUOBI,
		MinPriceIncrement: 1e-10,
		RoundLot:          0.01,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.InstrumentStatus_Trading,
		IgnoreSizeResidue: true,
	})
}
