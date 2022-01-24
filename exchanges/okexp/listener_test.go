package okexp_test

import (
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
)

func TestMarketData(t *testing.T) {
	tests.MarketData(t, tests.MDTest{
		SecurityID:        10652256150546133071,
		Symbol:            "BTC-USDT-SWAP",
		SecurityType:      enum.SecurityType_CRYPTO_PERP,
		Exchange:          constants.OKEXP,
		BaseCurrency:      constants.BITCOIN,
		QuoteCurrency:     constants.TETHER,
		MinPriceIncrement: 0.1,
		RoundLot:          1,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.Trading,
	})
}
