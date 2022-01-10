package v3_test

import (
	"testing"

	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
)

func TestMarketData(t *testing.T) {
	tests.PoolData(t, tests.MDTest{
		IgnoreSizeResidue: true,
		SecurityID:        6115381056087219401,
		Symbol:            "0x000ea4a83acefdd62b1b43e9ccc281f442651520",
		SecurityType:      enum.SecurityType_CRYPTO_AMM,
		Exchange:          constants.UNISWAPV3,
		QuoteCurrency: xchangerModels.Asset{
			Symbol: "WETH",
			Name:   "wrapped-ether",
			ID:     0,
		},
		BaseCurrency: xchangerModels.Asset{
			Symbol: "BUSD",
			Name:   "binance-usd",
			ID:     2,
		},
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.Trading,
		MinPriceIncrement: 60,
	})
}
