package ftxus_test

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
		SecurityID:        2028777944171259534,
		Symbol:            "BTC/USDT",
		SecurityType:      enum.SecurityType_CRYPTO_SPOT,
		Exchange:          constants.FTXUS,
		BaseCurrency:      constants.BITCOIN,
		QuoteCurrency:     constants.TETHER,
		MinPriceIncrement: 1,
		RoundLot:          0.0001,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.InstrumentStatus_Trading,
	})
}
