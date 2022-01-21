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
		SecurityID:        3923210566889873515,
		Symbol:            "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640",
		SecurityType:      enum.SecurityType_CRYPTO_AMM,
		Exchange:          constants.UNISWAPV3,
		QuoteCurrency: xchangerModels.Asset{
			Symbol: "WETH",
			Name:   "wrapped-ether",
			ID:     0,
		},
		BaseCurrency: xchangerModels.Asset{
			Symbol: "USDC",
			Name:   "usdc",
			ID:     1,
		},
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.Trading,
		MinPriceIncrement: 10,
	})
}
