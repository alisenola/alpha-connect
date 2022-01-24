package huobi_test

import (
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
)

func TestMarketData(t *testing.T) {
	tests.MarketData(t, tests.MDTest{
		SecurityID:        2195469462990134438,
		Symbol:            "btcusdt",
		SecurityType:      enum.SecurityType_CRYPTO_SPOT,
		Exchange:          constants.HUOBI,
		BaseCurrency:      constants.BITCOIN,
		QuoteCurrency:     constants.TETHER,
		MinPriceIncrement: 0.01,
		RoundLot:          1e-6,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.Trading,
	})
}
