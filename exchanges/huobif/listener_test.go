package huobif_test

import (
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
)

func TestMarketData(t *testing.T) {
	tests.MarketData(t, tests.MDTest{
		SecurityID:        2402007053666382556,
		Symbol:            "BTC210625",
		SecurityType:      enum.SecurityType_CRYPTO_FUT,
		Exchange:          constants.HUOBIF,
		BaseCurrency:      constants.BITCOIN,
		QuoteCurrency:     constants.DOLLAR,
		MinPriceIncrement: 0.01,
		RoundLot:          1.,
		HasMaturityDate:   true,
		IsInverse:         false,
		Status:            models.Trading,
	})
}
