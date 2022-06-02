package dydx_test

import (
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	exTests "gitlab.com/alphaticks/alpha-connect/tests"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
)

func TestMarketData(t *testing.T) {
	exTests.LoadStatics(t)
	tests.MarketData(t, tests.MDTest{
		SecurityID:        13112609607273681222,
		Symbol:            "ETH-USD",
		SecurityType:      enum.SecurityType_CRYPTO_PERP,
		Exchange:          constants.DYDX,
		BaseCurrency:      constants.ETHEREUM,
		QuoteCurrency:     constants.USDC,
		MinPriceIncrement: 0.1,
		RoundLot:          0.001,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.InstrumentStatus_Trading,
	})
}
