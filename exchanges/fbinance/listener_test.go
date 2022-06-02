package fbinance_test

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
		SecurityID:        5485975358912730733,
		Symbol:            "BTCUSDT",
		SecurityType:      enum.SecurityType_CRYPTO_PERP,
		Exchange:          constants.FBINANCE,
		BaseCurrency:      constants.BITCOIN,
		QuoteCurrency:     constants.TETHER,
		MinPriceIncrement: 0.1,
		RoundLot:          0.001,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.InstrumentStatus_Trading,
	})
}
