package bitfinex_test

import (
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
)

func TestMarketData(t *testing.T) {
	tests.MarketData(t, tests.MDTest{
		SecurityID:        17873758715870285590,
		Symbol:            "btcusd",
		SecurityType:      enum.SecurityType_CRYPTO_SPOT,
		Exchange:          constants.BITFINEX,
		BaseCurrency:      constants.BITCOIN,
		QuoteCurrency:     constants.DOLLAR,
		MinPriceIncrement: 0.1,
		RoundLot:          1. / 100000000.,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.Trading,
	})
}
