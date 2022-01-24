package deribit_test

import (
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
)

func TestMarketData(t *testing.T) {
	tests.MarketData(t, tests.MDTest{
		SecurityID:        2206542817128348325,
		Symbol:            "BTC-PERPETUAL",
		SecurityType:      enum.SecurityType_CRYPTO_PERP,
		Exchange:          constants.DERIBIT,
		BaseCurrency:      constants.BITCOIN,
		QuoteCurrency:     constants.DOLLAR,
		MinPriceIncrement: 0.5,
		RoundLot:          10.,
		HasMaturityDate:   false,
		IsInverse:         true,
		Status:            models.Trading,
	})
}
