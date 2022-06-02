package huobil_test

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
		SecurityID:        7064821541481199344,
		Symbol:            "BTC-USDT",
		SecurityType:      enum.SecurityType_CRYPTO_PERP,
		Exchange:          constants.HUOBIL,
		BaseCurrency:      constants.BITCOIN,
		QuoteCurrency:     constants.TETHER,
		MinPriceIncrement: 0.1,
		RoundLot:          1,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.InstrumentStatus_Trading,
	})
}
