package bitmex_test

import (
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"testing"
)

var instrument1 = &models.Instrument{
	SecurityID: &types.UInt64Value{Value: 5391998915988476130},
	Exchange:   &constants.BITMEX,
	Symbol:     &types.StringValue{Value: "XBTUSD"},
}
var instrument2 = &models.Instrument{
	SecurityID: &types.UInt64Value{Value: 11093839049553737303},
	Exchange:   &constants.BITMEX,
	Symbol:     &types.StringValue{Value: "ETHUSD"},
}

var bitmexAccount = &models.Account{
	Name:     "299210",
	Exchange: &constants.BITMEX,
	ApiCredentials: &xchangerModels.APICredentials{
		APIKey:    "k5k6Mmaq3xe88Ph3fgIk9Vrt",
		APISecret: "0laIjZaKOMkJPtKy2ldJ18m4Dxjp66Vdim0k1-q4TXASZFZo",
	},
}

func TestAccountListener(t *testing.T) {
	tests.AccntTest(t, tests.AccountTest{
		Account:                 bitmexAccount,
		Instrument:              instrument1,
		OrderStatusRequest:      true,
		OrderCancelRequest:      true,
		NewOrderSingleRequest:   true,
		NewOrderBulkRequest:     true,
		OrderReplaceRequest:     true,
		OrderBulkReplaceRequest: true,
		GetPositionsLimit:       true,
		GetPositionsMarket:      true,
	})
	tests.AccntTest(t, tests.AccountTest{
		Account:    bitmexAccount,
		Instrument: instrument2,
	})
}
