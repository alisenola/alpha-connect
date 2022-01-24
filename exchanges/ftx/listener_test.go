package ftx_test

import (
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
)

func TestMarketData(t *testing.T) {
	tests.MarketData(t, tests.MDTest{
		SecurityID:        4425198260936995601,
		Symbol:            "BTC-PERP",
		SecurityType:      enum.SecurityType_CRYPTO_PERP,
		Exchange:          constants.FTX,
		BaseCurrency:      constants.BITCOIN,
		QuoteCurrency:     constants.DOLLAR,
		MinPriceIncrement: 1,
		RoundLot:          0.0001,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.Trading,
	})
}
