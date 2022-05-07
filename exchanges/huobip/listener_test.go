package huobip_test

import (
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	xmodels "gitlab.com/alphaticks/xchanger/models"
	"testing"
)

func TestMarketData(t *testing.T) {
	tests.MarketData(t, tests.MDTest{
		SecurityID:   16932771734232838537,
		Symbol:       "OMG-USD",
		SecurityType: enum.SecurityType_CRYPTO_PERP,
		Exchange:     constants.HUOBIP,
		BaseCurrency: &xmodels.Asset{
			ID:       32,
			Name:     "",
			Symbol:   "OMG",
			Fungible: false,
		},
		QuoteCurrency:     constants.DOLLAR,
		MinPriceIncrement: 0.0001,
		RoundLot:          1,
		HasMaturityDate:   false,
		IsInverse:         true,
		Status:            models.InstrumentStatus_Trading,
	})
}
