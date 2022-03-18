package account

import (
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"math"
	"testing"
)

var bitstampAccount = &models.Account{
	Name:     "1",
	Exchange: &constants.BITSTAMP,
}

var bitmexAccount = &models.Account{
	Name:     "1",
	Exchange: &constants.BITMEX,
}

var fbinanceAccount = &models.Account{
	Name:     "1",
	Exchange: &constants.FBINANCE,
}

var BTCUSD_PERP_SEC = &models.Security{
	SecurityID:        9999999,
	SecurityType:      enum.SecurityType_CRYPTO_PERP,
	Exchange:          &constants.BITMEX,
	Symbol:            "XBTUSD",
	MinPriceIncrement: &types.DoubleValue{Value: 0.05},
	RoundLot:          &types.DoubleValue{Value: 1},
	Underlying:        &constants.BITCOIN,
	QuoteCurrency:     &constants.DOLLAR,
	IsInverse:         true,
	MakerFee:          &types.DoubleValue{Value: -0.00025},
	TakerFee:          &types.DoubleValue{Value: 0.00075},
	Multiplier:        &types.DoubleValue{Value: -1.},
	MaturityDate:      nil,
}

var BTCUSDT_PERP_SEC = &models.Security{
	SecurityID:        7744455,
	SecurityType:      enum.SecurityType_CRYPTO_PERP,
	Exchange:          &constants.FBINANCE,
	Symbol:            "BTCUSDT",
	MinPriceIncrement: &types.DoubleValue{Value: 0.05},
	RoundLot:          &types.DoubleValue{Value: 1},
	Underlying:        &constants.BITCOIN,
	QuoteCurrency:     &constants.TETHER,
	IsInverse:         false,
	MakerFee:          &types.DoubleValue{Value: 0.0002},
	TakerFee:          &types.DoubleValue{Value: 0.0004},
	Multiplier:        &types.DoubleValue{Value: 1.},
	MaturityDate:      nil,
}

var ETHUSD_PERP_SEC = &models.Security{
	SecurityID:        8888888,
	SecurityType:      enum.SecurityType_CRYPTO_PERP,
	Exchange:          &constants.BITMEX,
	Symbol:            "ETHUSD",
	MinPriceIncrement: &types.DoubleValue{Value: 0.05},
	RoundLot:          &types.DoubleValue{Value: 1},
	Underlying:        &constants.ETHEREUM,
	QuoteCurrency:     &constants.DOLLAR,
	IsInverse:         false,
	MakerFee:          &types.DoubleValue{Value: -0.00025},
	TakerFee:          &types.DoubleValue{Value: 0.00075},
	Multiplier:        &types.DoubleValue{Value: 0.000001},
	MaturityDate:      nil,
}

var BTCUSD_SPOT_SEC = &models.Security{
	SecurityID:        7777777,
	SecurityType:      enum.SecurityType_CRYPTO_SPOT,
	Exchange:          &constants.BITSTAMP,
	Symbol:            "BTCUSD",
	MinPriceIncrement: &types.DoubleValue{Value: 0.05},
	RoundLot:          &types.DoubleValue{Value: 0.0001},
	Underlying:        &constants.BITCOIN,
	QuoteCurrency:     &constants.DOLLAR,
	IsInverse:         false,
	MakerFee:          &types.DoubleValue{Value: 0.0025},
	TakerFee:          &types.DoubleValue{Value: 0.0025},
	MaturityDate:      nil,
}

var ETHUSD_SPOT_SEC = &models.Security{
	SecurityID:        6666666,
	SecurityType:      enum.SecurityType_CRYPTO_SPOT,
	Exchange:          &constants.BITSTAMP,
	Symbol:            "ETHUSD",
	MinPriceIncrement: &types.DoubleValue{Value: 0.05},
	RoundLot:          &types.DoubleValue{Value: 0.0001},
	Underlying:        &constants.ETHEREUM,
	QuoteCurrency:     &constants.DOLLAR,
	IsInverse:         false,
	MakerFee:          &types.DoubleValue{Value: 0.0025},
	TakerFee:          &types.DoubleValue{Value: 0.0025},
	MaturityDate:      nil,
}

func TestAccount_ConfirmFill(t *testing.T) {
	accnt, err := NewAccount(bitmexAccount)
	if err != nil {
		t.Fatal(err)
	}
	err = accnt.Sync([]*models.Security{ETHUSD_PERP_SEC}, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add a buy order
	_, rej := accnt.NewOrder(&models.Order{
		OrderID:       "buy",
		ClientOrderID: "buy",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: ETHUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "ETHUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Buy,
		TimeInForce:    models.Session,
		LeavesQuantity: 10.,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: 10.},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err = accnt.ConfirmNewOrder("buy", "buy")
	if err != nil {
		t.Fatal(err)
	}

	// Add a sell order
	_, rej = accnt.NewOrder(&models.Order{
		OrderID:       "sell",
		ClientOrderID: "sell",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: ETHUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "ETHUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Sell,
		TimeInForce:    models.Session,
		LeavesQuantity: 10.,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: 10.},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err = accnt.ConfirmNewOrder("sell", "sell")
	if err != nil {
		t.Fatal(err)
	}

	fee1 := math.Floor(0.00025*200*2*0.000001*accnt.MarginPrecision) / accnt.MarginPrecision
	fee2 := math.Floor(0.00025*210*2*0.000001*accnt.MarginPrecision) / accnt.MarginPrecision
	expectedMarginChange := ((210 - 200) * 2 * 0.000001) + fee1 + fee2

	_, err = accnt.ConfirmFill("sell", "k1", 210., 1., false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = accnt.ConfirmFill("sell", "k1", 210., 1., false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = accnt.ConfirmFill("buy", "k1", 200., 2., false)
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(accnt.GetMargin(nil)-expectedMarginChange) > 0.000000001 {
		t.Fatalf("was expecting margin of %g, got %g", expectedMarginChange, accnt.GetMargin(nil))
	}

	_, err = accnt.ConfirmFill("buy", "k1", 200., 2., false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = accnt.ConfirmFill("sell", "k1", 210., 1., false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = accnt.ConfirmFill("sell", "k1", 210., 1., false)
	if err != nil {
		t.Fatal(err)
	}

	expectedMarginChange = expectedMarginChange + expectedMarginChange
	if math.Abs(accnt.GetMargin(nil)-expectedMarginChange) > 0.000000001 {
		t.Fatalf("was expecting margin of %g, got %g", expectedMarginChange, accnt.GetMargin(nil))
	}
}

func TestAccount_ConfirmFill_Inverse(t *testing.T) {
	accnt, err := NewAccount(bitmexAccount)
	if err != nil {
		t.Fatal(err)
	}
	err = accnt.Sync([]*models.Security{BTCUSD_PERP_SEC}, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add a buy order
	_, rej := accnt.NewOrder(&models.Order{
		OrderID:       "buy",
		ClientOrderID: "buy",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: BTCUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "XBTUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Buy,
		TimeInForce:    models.Session,
		LeavesQuantity: 10.,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: 10.},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err = accnt.ConfirmNewOrder("buy", "buy")
	if err != nil {
		t.Fatal(err)
	}

	// Add a sell order
	_, rej = accnt.NewOrder(&models.Order{
		OrderID:       "sell",
		ClientOrderID: "sell",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: BTCUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "XBTUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Sell,
		TimeInForce:    models.Session,
		LeavesQuantity: 10.,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: 10.},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err = accnt.ConfirmNewOrder("sell", "sell")
	if err != nil {
		t.Fatal(err)
	}

	_, err = accnt.ConfirmFill("sell", "k1", 210., 2., false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = accnt.ConfirmFill("buy", "k1", 200., 1., false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = accnt.ConfirmFill("buy", "k1", 200., 1., false)
	if err != nil {
		t.Fatal(err)
	}

	fee1 := math.Floor(0.00025*(1./200.)*2.*accnt.MarginPrecision) / accnt.MarginPrecision
	fee2 := math.Floor(0.00025*(1./210.)*2.*accnt.MarginPrecision) / accnt.MarginPrecision

	cost1 := (math.Round(1./200.*accnt.MarginPrecision) / accnt.MarginPrecision) * 2.
	cost2 := (math.Round(1./210.*accnt.MarginPrecision) / accnt.MarginPrecision) * 2.
	expectedMarginChange := (cost1 - cost2) + fee1 + fee2
	if math.Abs(accnt.GetMargin(nil)-expectedMarginChange) > 0.00000001 {
		t.Fatalf("was expecting margin of %g, got %g", expectedMarginChange, accnt.GetMargin(nil))
	}

	_, err = accnt.ConfirmFill("buy", "k1", 200., 1., false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = accnt.ConfirmFill("buy", "k1", 200., 1., false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = accnt.ConfirmFill("sell", "k1", 210., 2., false)
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(accnt.GetMargin(nil)-2*expectedMarginChange) > 0.00000001 {
		t.Fatalf("was expecting margin of %g, got %g", 2*expectedMarginChange, accnt.GetMargin(nil))
	}
}

func TestAccount_ConfirmFill_Replace(t *testing.T) {
	// Post a matching limit order and post a replace right after

	accnt, err := NewAccount(bitmexAccount)
	if err != nil {
		t.Fatal(err)
	}
	err = accnt.Sync([]*models.Security{BTCUSD_PERP_SEC}, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add a buy order
	_, rej := accnt.NewOrder(&models.Order{
		OrderID:       "buy",
		ClientOrderID: "buy",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: BTCUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "XBTUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Buy,
		TimeInForce:    models.Session,
		LeavesQuantity: 2.,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: 20000.},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err = accnt.ConfirmNewOrder("buy", "buy")
	if err != nil {
		t.Fatal(err)
	}

	// Add a sell order
	_, rej = accnt.NewOrder(&models.Order{
		OrderID:       "sell",
		ClientOrderID: "sell",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: BTCUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "XBTUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Sell,
		TimeInForce:    models.Session,
		LeavesQuantity: 10.,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: 10.},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err = accnt.ConfirmNewOrder("sell", "sell")
	if err != nil {
		t.Fatal(err)
	}

	_, err = accnt.ConfirmFill("sell", "k1", 210., 2., false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = accnt.ConfirmFill("buy", "k1", 200., 1., false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = accnt.ConfirmFill("buy", "k1", 200., 1., false)
	if err != nil {
		t.Fatal(err)
	}

	fee1 := math.Floor(0.00025*(1./200.)*2.*accnt.MarginPrecision) / accnt.MarginPrecision
	fee2 := math.Floor(0.00025*(1./210.)*2.*accnt.MarginPrecision) / accnt.MarginPrecision

	cost1 := (math.Round(1./200.*accnt.MarginPrecision) / accnt.MarginPrecision) * 2.
	cost2 := (math.Round(1./210.*accnt.MarginPrecision) / accnt.MarginPrecision) * 2.
	expectedMarginChange := (cost1 - cost2) + fee1 + fee2
	if math.Abs(accnt.GetMargin(nil)-expectedMarginChange) > 0.00000001 {
		t.Fatalf("was expecting margin of %g, got %g", expectedMarginChange, accnt.GetMargin(nil))
	}

	_, err = accnt.ConfirmFill("buy", "k1", 200., 1., false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = accnt.ConfirmFill("buy", "k1", 200., 1., false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = accnt.ConfirmFill("sell", "k1", 210., 2., false)
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(accnt.GetMargin(nil)-2*expectedMarginChange) > 0.00000001 {
		t.Fatalf("was expecting margin of %g, got %g", 2*expectedMarginChange, accnt.GetMargin(nil))
	}
}

func TestAccount_Compare(t *testing.T) {
	accnt1, err := NewAccount(bitmexAccount)
	if err != nil {
		t.Fatal(err)
	}
	err = accnt1.Sync([]*models.Security{BTCUSD_PERP_SEC}, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	accnt2, err := NewAccount(bitmexAccount)
	if err != nil {
		t.Fatal(err)
	}
	err = accnt2.Sync([]*models.Security{BTCUSD_PERP_SEC}, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !accnt1.Compare(accnt2) {
		t.Fatalf("different account")
	}

	// Add a buy order
	_, rej := accnt1.NewOrder(&models.Order{
		OrderID:       "buy",
		ClientOrderID: "buy",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: BTCUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "XBTUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Buy,
		TimeInForce:    models.Session,
		LeavesQuantity: 10.,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: 10.},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	// Add a buy order
	_, rej = accnt2.NewOrder(&models.Order{
		OrderID:       "buy",
		ClientOrderID: "buy",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: BTCUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "XBTUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Buy,
		TimeInForce:    models.Session,
		LeavesQuantity: 10.,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: 10.},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	if !accnt1.Compare(accnt2) {
		t.Fatalf("different account")
	}

	_, err = accnt1.ConfirmNewOrder("buy", "buy")
	if err != nil {
		t.Fatal(err)
	}
	if accnt1.Compare(accnt2) {
		t.Fatalf("same account")
	}
	_, err = accnt2.ConfirmNewOrder("buy", "buy")
	if err != nil {
		t.Fatal(err)
	}
	if !accnt1.Compare(accnt2) {
		t.Fatalf("different account")
	}

	// Add a sell order
	_, rej = accnt1.NewOrder(&models.Order{
		OrderID:       "sell",
		ClientOrderID: "sell",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: BTCUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "XBTUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Sell,
		TimeInForce:    models.Session,
		LeavesQuantity: 10.,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: 10.},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	// Add a sell order
	_, rej = accnt2.NewOrder(&models.Order{
		OrderID:       "sell",
		ClientOrderID: "sell",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: BTCUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "XBTUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Sell,
		TimeInForce:    models.Session,
		LeavesQuantity: 10.,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: 10.},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	if !accnt1.Compare(accnt2) {
		t.Fatalf("different account")
	}

	_, err = accnt1.ConfirmNewOrder("sell", "sell")
	if err != nil {
		t.Fatal(err)
	}
	if accnt1.Compare(accnt2) {
		t.Fatalf("same account")
	}
	_, err = accnt2.ConfirmNewOrder("sell", "sell")
	if err != nil {
		t.Fatal(err)
	}
	if !accnt1.Compare(accnt2) {
		t.Fatalf("different account")
	}

	_, err = accnt1.ConfirmFill("sell", "k1", 210., 2., false)
	if err != nil {
		t.Fatal(err)
	}
	if accnt1.Compare(accnt2) {
		t.Fatalf("same account")
	}
	_, err = accnt2.ConfirmFill("sell", "k1", 210., 2., false)
	if err != nil {
		t.Fatal(err)
	}

	if !accnt1.Compare(accnt2) {
		t.Fatalf("different account")
	}

	_, err = accnt1.ConfirmFill("buy", "k1", 200., 1., false)
	if err != nil {
		t.Fatal(err)
	}
	if accnt1.Compare(accnt2) {
		t.Fatalf("same account")
	}
	_, err = accnt2.ConfirmFill("buy", "k1", 200., 1., false)
	if err != nil {
		t.Fatal(err)
	}
	if !accnt1.Compare(accnt2) {
		t.Fatalf("different account")
	}

	_, err = accnt1.ConfirmFill("buy", "k1", 200., 1., false)
	if err != nil {
		t.Fatal(err)
	}
	if accnt1.Compare(accnt2) {
		t.Fatalf("same account")
	}
	_, err = accnt2.ConfirmFill("buy", "k1", 200., 1., false)
	if err != nil {
		t.Fatal(err)
	}
	if !accnt1.Compare(accnt2) {
		t.Fatalf("different account")
	}

	_, err = accnt1.ConfirmFill("buy", "k1", 200., 1., false)
	if err != nil {
		t.Fatal(err)
	}
	if accnt1.Compare(accnt2) {
		t.Fatalf("same account")
	}
	_, err = accnt2.ConfirmFill("buy", "k1", 200., 1., false)
	if err != nil {
		t.Fatal(err)
	}
	if !accnt1.Compare(accnt2) {
		t.Fatalf("different account")
	}

	_, err = accnt1.ConfirmFill("buy", "k1", 200., 1., false)
	if err != nil {
		t.Fatal(err)
	}
	if accnt1.Compare(accnt2) {
		t.Fatalf("same account")
	}
	_, err = accnt2.ConfirmFill("buy", "k1", 200., 1., false)
	if err != nil {
		t.Fatal(err)
	}
	if !accnt1.Compare(accnt2) {
		t.Fatalf("different account")
	}

	_, err = accnt1.ConfirmFill("sell", "k1", 210., 2., false)
	if err != nil {
		t.Fatal(err)
	}
	if accnt1.Compare(accnt2) {
		t.Fatalf("same account")
	}
	_, err = accnt2.ConfirmFill("sell", "k1", 210., 2., false)
	if err != nil {
		t.Fatal(err)
	}
	if !accnt1.Compare(accnt2) {
		t.Fatalf("different account")
	}
}
