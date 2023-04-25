package krakenf_test

import (
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/krakenf"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
)

var instrument = &models.Instrument{
	SecurityID: wrapperspb.UInt64(1199421337905542595),
	Exchange:   constants.KRAKENF,
	Symbol:     wrapperspb.String("pf_ethusd"),
}

var krakenfAccount = &models.Account{
	Exchange: constants.KRAKENF,
	ApiCredentials: &xchangerModels.APICredentials{
		APIKey:    "vYm2R19p4XgYL09UxKXkH0nA5nn3IR2q3KVh/qAu4/6k8WfrtORy75Vb",
		APISecret: "dLuqPIvpaqmMzUpbN2YhgxFpJE2uDBVCZB2rRCKvqFoVMC0bXao3QaOritkTuuISvzBa+n6tnBHf0bUd62dYxqI7",
	},
}

func TestAccountListener(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	krakenf.EnableTestNet()
	tests.AccntTest(t, tests.AccountTest{
		Account:             krakenfAccount,
		Instrument:          instrument,
		OrderStatusRequest:  true,
		GetPositionsLimit:   true,
		GetPositionsMarket:  true,
		OrderReplaceRequest: true,
	})
}
