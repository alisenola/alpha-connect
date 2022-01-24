package coinbasepro_test

import (
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
)

func TestMarketData(t *testing.T) {
	tests.MarketData(t, tests.MDTest{
		SecurityID:        11630614572540763252,
		Symbol:            "BTC-USD",
		SecurityType:      enum.SecurityType_CRYPTO_SPOT,
		Exchange:          constants.COINBASEPRO,
		BaseCurrency:      constants.BITCOIN,
		QuoteCurrency:     constants.DOLLAR,
		MinPriceIncrement: 0.01,
		RoundLot:          1e-08,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.Trading,
	})
}
