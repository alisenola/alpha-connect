package v3_test

import (
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
)

func TestMarketData(t *testing.T) {
	tests.PoolData(t, tests.MDTest{
		IgnoreSizeResidue: true,
		SecurityID:        8096593996930596547,
		Symbol:            "BTC_USDT",
		SecurityType:      enum.SecurityType_CRYPTO_AMM,
		Exchange:          constants.UNISWAPV3,
		BaseCurrency:      constants.ETHEREUM,
		QuoteCurrency:     constants.USDC,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.Trading,
	})
}
