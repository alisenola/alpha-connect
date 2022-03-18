package ftx_test

import (
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"testing"
)

var instrument = &models.Instrument{
	Exchange: &constants.FTX,
	Symbol:   &types.StringValue{Value: "BTC-PERP"},
}

var ftxAccount = &models.Account{
	Name:     "299211",
	Exchange: &constants.FTX,
	ApiCredentials: &xchangerModels.APICredentials{
		APIKey:    "ej2YRRJMMwQD2qjOaRQnB18K7EWTpy3fRTP1ZZoX",
		APISecret: "MjJSlt1ix9OLpQLzkaQPlD_-y-N7c_2-6ZQJiwlY",
	},
}

func TestAccountListener(t *testing.T) {
	tests.AccntTest(t, tests.AccountTest{
		Account:                 ftxAccount,
		Instrument:              instrument,
		OrderStatusRequest:      true,
		OrderCancelRequest:      true,
		NewOrderSingleRequest:   true,
		OrderReplaceRequest:     true,
		OrderBulkReplaceRequest: true,
		GetPositionsLimit:       true,
		GetPositionsMarket:      true,
	})
}
