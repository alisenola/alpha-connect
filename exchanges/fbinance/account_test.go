package fbinance_test

import (
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/fbinance"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"testing"
)

var instrument = &models.Instrument{
	SecurityID: &types.UInt64Value{Value: 5485975358912730733},
	Exchange:   &constants.FBINANCE,
	Symbol:     &types.StringValue{Value: "BTCUSDT"},
}

var fbinanceAccount = &models.Account{
	Exchange: &constants.FBINANCE,
	ApiCredentials: &xchangerModels.APICredentials{
		APIKey:    "74f122652da74f6e1bcc34b8c23fc91e0239b502e68440632ae9a3cb7cefa18e",
		APISecret: "c3e0d76ee014b597b93616478dc789e6bb6616ad59ddbe384d2554ace4a60f86",
	},
}

func TestAccountListener(t *testing.T) {
	fbinance.EnableTestNet()
	tests.AccntTest(t, tests.AccountTest{
		Account:               fbinanceAccount,
		Instrument:            instrument,
		OrderStatusRequest:    true,
		OrderCancelRequest:    true,
		GetPositionsLimit:     true,
		GetPositionsMarket:    true,
		NewOrderSingleRequest: true,
		OrderReplaceRequest:   true,
	})
}
