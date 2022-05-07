package binance_test

import (
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/binance"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
)

var binanceInstrument = &models.Instrument{
	Exchange: constants.BINANCE,
	Symbol:   &wrapperspb.StringValue{Value: "BTCUSDT"},
}

var binanceAccount = &models.Account{
	Exchange: constants.BINANCE,
	ApiCredentials: &xchangerModels.APICredentials{
		APIKey:    "Ht0uwwcEDGChidsBgua1utw8I3Qwsin8fS5CEmZIPvHFfM5rapbKAzR1SInIe41f",
		APISecret: "3rUAN1uUeMpG9C0HOOmOIQTEijs1rz3VV0PaqExusBJLB0FzRrKE1CjzkVK4ZiSM",
	},
}

func TestAccountListener(t *testing.T) {
	binance.EnableTestNet()
	tests.AccntTest(t, tests.AccountTest{
		Account:            binanceAccount,
		Instrument:         binanceInstrument,
		OrderStatusRequest: true,
	})
}
