package binance_test

import (
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/binance"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"testing"
)

func TestAccountListener(t *testing.T) {
	binance.EnableTestNet()
	tests.MarketFill(t, tests.AccntTest{
		Account: &models.Account{
			Name:     "binance",
			Exchange: &constants.BINANCE,
			ApiCredentials: &xchangerModels.APICredentials{
				APIKey:    "DNySYXVSG7xrM7S8dGSTvLRmRGJIzoU80Uj78IpFOYxiI9veS54VCu8bxQNLloz2",
				APISecret: "OrIJH6qYynrVFgEi62cknUpf02MoA82l45ySfP3ZTKRPainFJzG377BoJJmTuwXv",
			},
		},
		Instrument: &models.Instrument{
			Exchange: &constants.BINANCE,
			Symbol:   &types.StringValue{Value: "BTCUSDT"},
		},
	})
}
