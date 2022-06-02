package cryptofacilities_test

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
		SecurityID:        1416768858288349990,
		Symbol:            "pi_xbtusd",
		SecurityType:      enum.SecurityType_CRYPTO_PERP,
		Exchange:          constants.CRYPTOFACILITIES,
		BaseCurrency:      constants.BITCOIN,
		QuoteCurrency:     constants.DOLLAR,
		MinPriceIncrement: 0.5,
		RoundLot:          1,
		HasMaturityDate:   false,
		IsInverse:         true,
		Status:            models.InstrumentStatus_Trading,
	})
}
