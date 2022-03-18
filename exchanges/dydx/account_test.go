package dydx_test

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"testing"
)

var account = &models.Account{
	Name:     "dydx",
	Exchange: &constants.DYDX,
	ApiCredentials: &xchangerModels.APICredentials{
		AccountID: "adc3dffd-9999-57ae-bd1a-34279aa9fa4b",
		APIKey:    "d678eb7d-e3c0-879d-f02d-9ea26ab26152",
		APISecret: "nazLBU9Xth156mVaktY7fZSXcLJ7Dv3wF5TpihMU:TIVcVRWBnZAc-AwDO8UC",
	},
	StarkCredentials: &xchangerModels.STARKCredentials{
		PositionId: 98309,
		PrivateKey: hexutil.MustDecode("0x0657d2ac46d36c1f06715ffb2ff6b8988403c9aad18e548e7a702b17109eeb76"),
	},
}

var instrument = &models.Instrument{
	Exchange: &constants.DYDX,
	Symbol:   &types.StringValue{Value: "BTC-USD"},
}

func TestAccountListener(t *testing.T) {
	tests.AccntTest(t, tests.AccountTest{
		Account:    account,
		Instrument: instrument,
	})
}
