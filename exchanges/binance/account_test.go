package binance_test

import (
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"testing"
)

var binanceInstrument = &models.Instrument{
	Exchange: &constants.BINANCE,
	Symbol:   &types.StringValue{Value: "BTCUSDT"},
}

var binanceAccount = &models.Account{
	Exchange: &constants.BINANCE,
	ApiCredentials: &xchangerModels.APICredentials{
		APIKey:    "XgHlK1xiQcq5IBP6KsmD7ewZONJHjhlq9JqopTuoB7lanSw3TYdLAcFn5fudyevO",
		APISecret: "GAiZ4UOztJzoB3Qonv4nE2X8KgEyl0jxSBtLBkpwgXeZYtSqGogwlh89YerLSqlu",
	},
}

func TestAccountListener(t *testing.T) {
	tests.AccntTest(t, tests.AccountTest{
		Account:    binanceAccount,
		Instrument: binanceInstrument,
	})
}
