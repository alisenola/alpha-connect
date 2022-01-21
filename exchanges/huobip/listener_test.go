package huobip_test

import (
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
)

func TestMarketData(t *testing.T) {
	tests.MarketData(t, tests.MDTest{
		SecurityID:        13641637530641868249,
		Symbol:            "KRW-BTC",
		SecurityType:      enum.SecurityType_CRYPTO_SPOT,
		Exchange:          constants.UPBIT,
		BaseCurrency:      constants.BITCOIN,
		QuoteCurrency:     constants.SOUTH_KOREAN_WON,
		MinPriceIncrement: 0.,
		RoundLot:          0.,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.Trading,
	})
}
