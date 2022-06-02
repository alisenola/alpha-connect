package bybitl_test

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
		SecurityID:        6789757764526280996,
		Symbol:            "BTCUSDT",
		SecurityType:      enum.SecurityType_CRYPTO_PERP,
		Exchange:          constants.BYBITL,
		BaseCurrency:      constants.BITCOIN,
		QuoteCurrency:     constants.TETHER,
		MinPriceIncrement: 0.5,
		RoundLot:          0.001,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.InstrumentStatus_Trading,
	})
}
