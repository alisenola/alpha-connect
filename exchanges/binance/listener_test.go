package binance_test

import (
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	exTests "gitlab.com/alphaticks/alpha-connect/tests"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
)

func TestMarketData(t *testing.T) {
	exTests.LoadStatics(t)
	tests.MarketData(t, tests.MDTest{
		SecurityID:        9281941173829172773,
		Symbol:            "BTCUSDT",
		SecurityType:      enum.SecurityType_CRYPTO_SPOT,
		Exchange:          constants.BINANCE,
		BaseCurrency:      constants.BITCOIN,
		QuoteCurrency:     constants.TETHER,
		MinPriceIncrement: 0.01,
		RoundLot:          0.000010,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.InstrumentStatus_Trading,
	})
}
