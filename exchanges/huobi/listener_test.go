package huobi_test

import (
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
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
		Status:            models.InstrumentStatus_Trading,
		IgnoreSizeResidue: true,
	})
}

func TestSlowMarketData(t *testing.T) {
	tests.MarketData(t, tests.MDTest{
		SecurityID:   4304365894942658747,
		Symbol:       "kaneth",
		SecurityType: enum.SecurityType_CRYPTO_SPOT,
		Exchange:     constants.HUOBI,
		BaseCurrency: &xchangerModels.Asset{
			Symbol: "KAN",
			Name:   "bitkan",
			ID:     220,
		},
		QuoteCurrency:     constants.ETHEREUM,
		MinPriceIncrement: 1e-10,
		RoundLot:          0.01,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.InstrumentStatus_Trading,
		IgnoreSizeResidue: true,
	})
}
