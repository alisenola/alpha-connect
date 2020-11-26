package account

import (
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/modeling"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"math"
	"os"
	"testing"
)

var bitmex = constants.BITMEX
var model modeling.MarketModel

func TestMain(m *testing.M) {
	mdl := modeling.NewMapModel()
	mdl.SetPriceModel(uint64(constants.BITCOIN.ID)<<32|uint64(constants.TETHER.ID), modeling.NewConstantPriceModel(100.))
	mdl.SetPriceModel(uint64(constants.BITCOIN.ID)<<32|uint64(constants.DOLLAR.ID), modeling.NewConstantPriceModel(100.))
	mdl.SetPriceModel(uint64(constants.ETHEREUM.ID)<<32|uint64(constants.DOLLAR.ID), modeling.NewConstantPriceModel(10.))
	mdl.SetPriceModel(uint64(constants.DOLLAR.ID)<<32|uint64(constants.DOLLAR.ID), modeling.NewConstantPriceModel(1.))
	mdl.SetPriceModel(uint64(constants.TETHER.ID)<<32|uint64(constants.DOLLAR.ID), modeling.NewConstantPriceModel(1.))

	mdl.SetPriceModel(BTCUSDT_PERP_SEC.SecurityID, modeling.NewConstantPriceModel(100.))
	mdl.SetPriceModel(BTCUSD_PERP_SEC.SecurityID, modeling.NewConstantPriceModel(100.))
	mdl.SetPriceModel(ETHUSD_PERP_SEC.SecurityID, modeling.NewConstantPriceModel(10.))

	mdl.SetBuyTradeModel(BTCUSDT_PERP_SEC.SecurityID, modeling.NewConstantTradeModel(2))
	mdl.SetBuyTradeModel(BTCUSD_PERP_SEC.SecurityID, modeling.NewConstantTradeModel(20))
	mdl.SetBuyTradeModel(ETHUSD_PERP_SEC.SecurityID, modeling.NewConstantTradeModel(20))
	mdl.SetBuyTradeModel(BTCUSD_SPOT_SEC.SecurityID, modeling.NewConstantTradeModel(20))
	mdl.SetBuyTradeModel(ETHUSD_SPOT_SEC.SecurityID, modeling.NewConstantTradeModel(20))

	mdl.SetSellTradeModel(BTCUSDT_PERP_SEC.SecurityID, modeling.NewConstantTradeModel(2))
	mdl.SetSellTradeModel(BTCUSD_PERP_SEC.SecurityID, modeling.NewConstantTradeModel(20))
	mdl.SetSellTradeModel(ETHUSD_PERP_SEC.SecurityID, modeling.NewConstantTradeModel(20))
	mdl.SetSellTradeModel(BTCUSD_SPOT_SEC.SecurityID, modeling.NewConstantTradeModel(20))
	mdl.SetSellTradeModel(ETHUSD_SPOT_SEC.SecurityID, modeling.NewConstantTradeModel(20))

	model = mdl

	os.Exit(m.Run())
}

func TestAccount_GetAvailableMargin(t *testing.T) {
	account := NewAccount(account, &constants.BITCOIN, 1./0.00000001)
	if err := account.Sync([]*models.Security{BTCUSD_PERP_SEC, ETHUSD_PERP_SEC}, nil, nil, nil, 0.1, nil, nil); err != nil {
		t.Fatal(err)
	}
	expectedAv := 0.1
	avMargin := account.GetAvailableMargin(model, 1.)
	if math.Abs(avMargin-expectedAv) > 0.0000001 {
		t.Fatalf("was expecting %g, got %g", expectedAv, avMargin)
	}
	// Add a buy order. Using o.quantity allows us to check if the returned order's quantity is correct too
	_, rej := account.NewOrder(&models.Order{
		OrderID:       "buy1",
		ClientOrderID: "buy1",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: ETHUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "ETHUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Buy,
		TimeInForce:    models.Session,
		LeavesQuantity: 10,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: 90.},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err := account.ConfirmNewOrder("buy1", "buy1")
	if err != nil {
		t.Fatal(err)
	}
	account.ConfirmFill("buy1", "", 9., 10, false)
	// Balance + maker rebate + entry cost + PnL
	mul := ETHUSD_PERP_SEC.Multiplier.Value
	expectedAv = 0.1 + (0.00025 * 10 * 9 * mul) - (10 * 9 * mul) + (10.-9.)*mul*10.
	avMargin = account.GetAvailableMargin(model, 1.)
	if math.Abs(avMargin-expectedAv) > 0.0000001 {
		t.Fatalf("was expecting %g, got %g", expectedAv, avMargin)
	}
}

func TestAccount_GetAvailableMargin_Inverse(t *testing.T) {
	account := NewAccount(account, &constants.BITCOIN, 1./0.00000001)
	if err := account.Sync([]*models.Security{BTCUSD_PERP_SEC, ETHUSD_PERP_SEC}, nil, nil, nil, 0.1, nil, nil); err != nil {
		t.Fatal(err)
	}
	expectedAv := 0.1
	avMargin := account.GetAvailableMargin(model, 1.)
	if math.Abs(avMargin-expectedAv) > 0.0000001 {
		t.Fatalf("was expecting %g, got %g", expectedAv, avMargin)
	}
	// Add a buy order. Using o.quantity allows us to check if the returned order's quantity is correct too
	_, rej := account.NewOrder(&models.Order{
		OrderID:       "buy1",
		ClientOrderID: "buy1",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: BTCUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "XBTUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Buy,
		TimeInForce:    models.Session,
		LeavesQuantity: 10,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: 90.},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err := account.ConfirmNewOrder("buy1", "buy1")
	if err != nil {
		t.Fatal(err)
	}
	account.ConfirmFill("buy1", "", 90., 10, false)
	// Balance + maker rebate + entry cost + PnL
	expectedAv = 0.1 + (0.00025 * 10 * (1. / 90.)) - (10 * (1 / 90.)) + ((1./90.)-(1./100.))*10
	avMargin = account.GetAvailableMargin(model, 1.)
	if math.Abs(avMargin-expectedAv) > 0.0000001 {
		t.Fatalf("was expecting %g, got %g", expectedAv, avMargin)
	}
}

func TestAccount_PnL_Inverse(t *testing.T) {
	account := NewAccount(account, &constants.BITCOIN, 1./0.00000001)
	if err := account.Sync([]*models.Security{BTCUSD_PERP_SEC, ETHUSD_PERP_SEC}, nil, nil, nil, 0.1, nil, nil); err != nil {
		t.Fatal(err)
	}
	expectedAv := 0.1
	avMargin := account.GetAvailableMargin(model, 1.)
	if math.Abs(avMargin-expectedAv) > 0.0000001 {
		t.Fatalf("was expecting %g, got %g", expectedAv, avMargin)
	}
	// Add a buy order. Using o.quantity allows us to check if the returned order's quantity is correct too
	_, rej := account.NewOrder(&models.Order{
		OrderID:       "buy1",
		ClientOrderID: "buy1",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: BTCUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "XBTUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Buy,
		TimeInForce:    models.Session,
		LeavesQuantity: 100,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: 90.},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err := account.ConfirmNewOrder("buy1", "buy1")
	if err != nil {
		t.Fatal(err)
	}

	_, rej = account.NewOrder(&models.Order{
		OrderID:       "sell1",
		ClientOrderID: "sell1",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: BTCUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "XBTUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Sell,
		TimeInForce:    models.Session,
		LeavesQuantity: 100,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: 110.},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err = account.ConfirmNewOrder("sell1", "sell1")
	if err != nil {
		t.Fatal(err)
	}

	account.ConfirmFill("buy1", "", 90., 10, false)
	// Balance + maker rebate + entry cost + PnL
	expectedAv = 0.1 + (0.00025 * 10 * (1. / 90.)) - (10 * (1 / 90.)) + ((1./90.)-(1./100.))*10
	avMargin = account.GetAvailableMargin(model, 1.)
	if math.Abs(avMargin-expectedAv) > 0.0000001 {
		t.Fatalf("was expecting %g, got %g", expectedAv, avMargin)
	}

	account.ConfirmFill("sell1", "", 110., 20, false)

	account.ConfirmFill("buy1", "", 90., 10, false)
	account.ConfirmFill("sell1", "", 110., 10, false)
	account.ConfirmFill("buy1", "", 90., 10, false)

}

func TestPortfolio_Spot_ELR(t *testing.T) {
	account := NewAccount(account, &constants.DOLLAR, 1./0.00000001)
	dollarBalance := &models.Balance{
		AccountID: "1",
		Asset:     &constants.DOLLAR,
		Quantity:  100,
	}
	ethereumBalance := &models.Balance{
		AccountID: "1",
		Asset:     &constants.ETHEREUM,
		Quantity:  10,
	}
	if err := account.Sync([]*models.Security{BTCUSD_SPOT_SEC, ETHUSD_SPOT_SEC}, nil, nil, []*models.Balance{dollarBalance, ethereumBalance}, 0., nil, nil); err != nil {
		t.Fatal(err)
	}

	p := NewPortfolio(1000)
	p.AddAccount(account)

	expectedBaseChange := 2 - (2 * 0.0025)
	expectedQuoteChange := 2 * 50.
	expectedValueChange := expectedBaseChange*100. - expectedQuoteChange
	expectedElr := math.Log((p.Value(model) + expectedValueChange) / p.Value(model))

	// The trade needs to be profitable or the portfolio will return us a nil order and an ELR of 0
	elr, o := p.GetELROnLimitBid("1", BTCUSD_SPOT_SEC.SecurityID, model, 10, []float64{50}, []float64{1}, 100.)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	o.Quantity = math.Round(o.Quantity/BTCUSD_SPOT_SEC.RoundLot) * BTCUSD_SPOT_SEC.RoundLot
	// Add a buy order. Using o.quantity allows us to check if the returned order's quantity is correct too
	_, rej := account.NewOrder(&models.Order{
		OrderID:       "buy1",
		ClientOrderID: "buy1",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: BTCUSD_SPOT_SEC.SecurityID},
			Exchange:   &constants.BITSTAMP,
			Symbol:     &types.StringValue{Value: "BTCUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Buy,
		TimeInForce:    models.Session,
		LeavesQuantity: o.Quantity,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: o.Price},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err := account.ConfirmNewOrder("buy1", "buy1")
	if err != nil {
		t.Fatal(err)
	}

	elr = p.ExpectedLogReturn(model, 10)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	elr = p.ExpectedLogReturn(model, 11)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	account.CancelOrder("buy1")
	if _, err := account.ConfirmCancelOrder("buy1"); err != nil {
		t.Fatal(err)
	}

	expectedElr = 0.
	elr = p.ExpectedLogReturn(model, 11)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	expectedBaseChange = 10
	expectedQuoteChange = 10 * 20. * (1 - 0.0025)
	expectedValueChange = expectedQuoteChange - expectedBaseChange*10.
	expectedElr = math.Log((p.Value(model) + expectedValueChange) / p.Value(model))

	elr, o = p.GetELROnLimitAsk("1", ETHUSD_SPOT_SEC.SecurityID, model, 10, []float64{20}, []float64{1}, 10)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	o.Quantity = math.Round(o.Quantity/ETHUSD_SPOT_SEC.RoundLot) * ETHUSD_SPOT_SEC.RoundLot
	// Add a buy order. Using o.quantity allows us to check if the returned order's quantity is correct too
	_, rej = account.NewOrder(&models.Order{
		OrderID:       "buy2",
		ClientOrderID: "buy2",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: ETHUSD_SPOT_SEC.SecurityID},
			Exchange:   &constants.BITSTAMP,
			Symbol:     &types.StringValue{Value: "ETHUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Sell,
		TimeInForce:    models.Session,
		LeavesQuantity: o.Quantity,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: o.Price},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err = account.ConfirmNewOrder("buy2", "buy2")
	if err != nil {
		t.Fatal(err)
	}

	elr = p.ExpectedLogReturn(model, 10)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	elr = p.ExpectedLogReturn(model, 11)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	account.CancelOrder("buy2")
	if _, err := account.ConfirmCancelOrder("buy2"); err != nil {
		t.Fatal(err)
	}

	expectedElr = 0.
	elr = p.ExpectedLogReturn(model, 11)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
}

func TestPortfolio_Margin_ELR(t *testing.T) {

	account := NewAccount(account, &constants.BITCOIN, 1./0.00000001)
	if err := account.Sync([]*models.Security{BTCUSD_PERP_SEC, ETHUSD_PERP_SEC}, nil, nil, nil, 0.1, nil, nil); err != nil {
		t.Fatal(err)
	}

	p := NewPortfolio(1000)
	p.AddAccount(account)

	expectedMarginChange := ((1./90 - 1./100) * 9) - (0.00075 * (1. / 90) * 9)
	expectedElr := math.Log((p.Value(model) + (expectedMarginChange * 100)) / p.Value(model))
	elr, o := p.GetELROnMarketBuy("1", BTCUSD_PERP_SEC.SecurityID, model, 10, 90, 9, 0.1)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	expectedMarginChange = ((1./90 - 1./100) * 9) + (0.00025 * (1. / 90) * 9)
	expectedElr = math.Log((p.Value(model) + (expectedMarginChange * 100)) / p.Value(model))

	elr, o = p.GetELROnLimitBid("1", BTCUSD_PERP_SEC.SecurityID, model, 10, []float64{90}, []float64{1}, 0.1)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	// Add a buy order. Using o.quantity allows us to check if the returned order's quantity is correct too
	_, rej := account.NewOrder(&models.Order{
		OrderID:       "buy1",
		ClientOrderID: "buy1",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: BTCUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "XBTUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Buy,
		TimeInForce:    models.Session,
		LeavesQuantity: o.Quantity,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: o.Price},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err := account.ConfirmNewOrder("buy1", "buy1")
	if err != nil {
		t.Fatal(err)
	}
	account.UpdateBidOrderQueue(BTCUSD_PERP_SEC.SecurityID, "buy1", 1)
	// Try with same time to test value cache consistency
	elr = p.ExpectedLogReturn(model, 10)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	// Try with different time
	elr = p.ExpectedLogReturn(model, 11)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	elr = p.GetELROnCancelBid("1", BTCUSD_PERP_SEC.SecurityID, "buy1", model, 11)
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}

	account.CancelOrder("buy1")
	_, err = account.ConfirmCancelOrder("buy1")
	if err != nil {
		t.Fatal(err)
	}

	expectedElr = 0.
	elr = p.ExpectedLogReturn(model, 11)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	elr = p.ExpectedLogReturn(model, 10)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	// ETHUSD

	expectedMarginChange = ((10 - 9) * 0.000001 * 19) + (-0.00075 * 9. * 19 * 0.000001)
	expectedElr = math.Log((p.Value(model) + (expectedMarginChange * 100)) / p.Value(model))
	elr, o = p.GetELROnMarketBuy("1", ETHUSD_PERP_SEC.SecurityID, model, 10, 9, 19, 0.1)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	expectedMarginChange = ((10 - 9) * 0.000001 * 19) + (0.00025 * 9. * 19 * 0.000001)
	expectedElr = math.Log((p.Value(model) + (expectedMarginChange * 100)) / p.Value(model))

	elr, o = p.GetELROnLimitBid("1", ETHUSD_PERP_SEC.SecurityID, model, 10, []float64{9}, []float64{1}, 0.1)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	o.Quantity = math.Round(o.Quantity)
	// Add a buy order. Using o.quantity allows us to check if the returned order's quantity is correct too
	_, rej = account.NewOrder(&models.Order{
		OrderID:       "buy2",
		ClientOrderID: "buy2",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: ETHUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "ETHUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Buy,
		TimeInForce:    models.Session,
		LeavesQuantity: o.Quantity,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: o.Price},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err = account.ConfirmNewOrder("buy2", "buy2")
	if err != nil {
		t.Fatal(err)
	}
	account.UpdateBidOrderQueue(ETHUSD_PERP_SEC.SecurityID, "buy2", 1)

	elr = p.ExpectedLogReturn(model, 10)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	elr = p.ExpectedLogReturn(model, 11)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	elr = p.GetELROnCancelBid("1", ETHUSD_PERP_SEC.SecurityID, "buy2", model, 11)
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}
	account.CancelOrder("buy2")
	if _, err := account.ConfirmCancelOrder("buy2"); err != nil {
		t.Fatal(err)
	}

	expectedElr = 0.

	elr = p.ExpectedLogReturn(model, 11)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	elr = p.ExpectedLogReturn(model, 10)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	// BTC USD short
	// Short 11

	expectedMarginChange = ((1./100 - 1./110) * 1. * 11.) + (-0.00075 * (1. / 110) * 11.)
	expectedElr = math.Log((p.Value(model) + (expectedMarginChange * 100)) / p.Value(model))
	elr, o = p.GetELROnMarketSell("1", BTCUSD_PERP_SEC.SecurityID, model, 10, 110, 11, 0.1)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	expectedMarginChange = ((1./100 - 1./110) * 1. * 11.) + (0.00025 * (1. / 110) * 11.)
	expectedElr = math.Log((p.Value(model) + (expectedMarginChange * 100)) / p.Value(model))
	elr, o = p.GetELROnLimitAsk("1", BTCUSD_PERP_SEC.SecurityID, model, 10, []float64{110}, []float64{1}, 0.1)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	o.Quantity = math.Round(o.Quantity)

	// Add a sell order. Using o.quantity allows us to check if the returned order's quantity is correct too
	_, rej = account.NewOrder(&models.Order{
		OrderID:       "sell1",
		ClientOrderID: "sell1",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: BTCUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "XBTUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Sell,
		TimeInForce:    models.Session,
		LeavesQuantity: o.Quantity,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: o.Price},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err = account.ConfirmNewOrder("sell1", "sell1")
	if err != nil {
		t.Fatal(err)
	}
	account.UpdateAskOrderQueue(BTCUSD_PERP_SEC.SecurityID, "sell1", 1)

	elr = p.ExpectedLogReturn(model, 10)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	elr = p.ExpectedLogReturn(model, 11)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	elr = p.GetELROnCancelAsk("1", BTCUSD_PERP_SEC.SecurityID, "sell1", model, 11)
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}
	elr = p.GetELROnCancelAsk("1", BTCUSD_PERP_SEC.SecurityID, "sell1", model, 10)
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}

	account.CancelOrder("sell1")
	if _, err := account.ConfirmCancelOrder("sell1"); err != nil {
		t.Fatal(err)
	}
	elr = p.ExpectedLogReturn(model, 10)
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}
	elr = p.ExpectedLogReturn(model, 11)
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}

	// ETHUSD Short
	expectedMarginChange = ((11 - 10) * 0.000001 * 19) + (-0.00075 * 11. * 19 * 0.000001)
	expectedElr = math.Log((p.Value(model) + (expectedMarginChange * 100)) / p.Value(model))
	elr, o = p.GetELROnMarketSell("1", ETHUSD_PERP_SEC.SecurityID, model, 11, 11, 19, 0.1)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	expectedMarginChange = ((11 - 10) * 0.000001 * 19) + (0.00025 * 11. * 19 * 0.000001)
	expectedElr = math.Log((p.Value(model) + (expectedMarginChange * 100)) / p.Value(model))
	elr, o = p.GetELROnLimitAsk("1", ETHUSD_PERP_SEC.SecurityID, model, 11, []float64{11}, []float64{1}, 0.1)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	o.Quantity = math.Round(o.Quantity)

	// Add a sell order. Using o.quantity allows us to check if the returned order's quantity is correct too
	_, rej = account.NewOrder(&models.Order{
		OrderID:       "sell2",
		ClientOrderID: "sell2",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: ETHUSD_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "ETHUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Sell,
		TimeInForce:    models.Session,
		LeavesQuantity: o.Quantity,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: o.Price},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err = account.ConfirmNewOrder("sell2", "sell2")
	if err != nil {
		t.Fatal(err)
	}
	account.UpdateAskOrderQueue(ETHUSD_PERP_SEC.SecurityID, "sell2", 1)

	elr = p.ExpectedLogReturn(model, 10)

	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	elr = p.GetELROnCancelAsk("1", ETHUSD_PERP_SEC.SecurityID, "sell2", model, 10)
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}

	account.CancelOrder("sell2")
	if _, err := account.ConfirmCancelOrder("sell2"); err != nil {
		t.Fatal(err)
	}
	elr = p.ExpectedLogReturn(model, 10)
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}
}

/*
// Test case where you have multiple go routines accessing the portfolio together
// One go routine updates contracts, the other
func TestBitmexPortfolio_Parallel(t *testing.T) {
	gbm := NewGBMPriceModel(100., 100)
	priceModels := make(map[uint64]PriceModel)
	priceModels[bitmexInstruments[0].ID()] = gbm

	buyTradeModels := make(map[uint64]BuyTradeModel)
	buyTradeModels[bitmexInstruments[0].ID()] = NewConstantTradeModel(10)

	sellTradeModels := make(map[uint64]SellTradeModel)
	sellTradeModels[bitmexInstruments[0].ID()] = NewConstantTradeModel(10)

	ep, err := NewBitmexPortfolio(bitmexInstruments[:1], gbm, priceModels, buyTradeModels, sellTradeModels, 1000)
	if err != nil {
		t.Fatal(err)
	}
	p := NewPortfolio(map[uint64]ExchangePortfolio{bitmex: ep}, 1000)

	ep.SetWallet(1.)

	fmt.Println(p.ExpectedLogReturn(800))
	// One go routine place orders, and update wallet, etc

	// Others eval the price
}

*/

func TestPortfolio_Fbinance_Margin_ELR(t *testing.T) {

	account := NewAccount(account, &constants.TETHER, 1./0.00000001)
	if err := account.Sync([]*models.Security{BTCUSDT_PERP_SEC}, nil, nil, nil, 1000, nil, nil); err != nil {
		t.Fatal(err)
	}

	p := NewPortfolio(1000)
	p.AddAccount(account)

	expectedMarginChange := ((100 - 90) * 0.1) - (0.0004 * 90 * 0.1)
	expectedElr := math.Log((p.Value(model) + (expectedMarginChange)) / p.Value(model))
	elr, o := p.GetELROnMarketBuy("1", BTCUSDT_PERP_SEC.SecurityID, model, 10, 90, 0.1, 1000)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	// Match on bid of 2, queue of 1, match of 1
	expectedMarginChange = ((100 - 90) * 1) - (0.0002 * 90 * 1)
	expectedElr = math.Log((p.Value(model) + (expectedMarginChange)) / p.Value(model))

	elr, o = p.GetELROnLimitBid("1", BTCUSDT_PERP_SEC.SecurityID, model, 10, []float64{90}, []float64{1}, 1000)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	// Add a buy order. Using o.quantity allows us to check if the returned order's quantity is correct too
	o.Quantity = 1.
	_, rej := account.NewOrder(&models.Order{
		OrderID:       "buy1",
		ClientOrderID: "buy1",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: BTCUSDT_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "XBTUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Buy,
		TimeInForce:    models.Session,
		LeavesQuantity: o.Quantity,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: o.Price},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err := account.ConfirmNewOrder("buy1", "buy1")
	if err != nil {
		t.Fatal(err)
	}
	account.UpdateBidOrderQueue(BTCUSDT_PERP_SEC.SecurityID, "buy1", 1)
	// Try with same time to test value cache consistency
	elr = p.ExpectedLogReturn(model, 10)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	// Try with different time
	elr = p.ExpectedLogReturn(model, 11)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	elr = p.GetELROnCancelBid("1", BTCUSDT_PERP_SEC.SecurityID, "buy1", model, 11)
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}

	account.CancelOrder("buy1")
	_, err = account.ConfirmCancelOrder("buy1")
	if err != nil {
		t.Fatal(err)
	}

	expectedElr = 0.
	elr = p.ExpectedLogReturn(model, 11)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	elr = p.ExpectedLogReturn(model, 10)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	// SHORT

	// ETHUSD Short
	expectedMarginChange = ((110 - 100) * 1) - (0.0004 * 110 * 1)
	expectedElr = math.Log((p.Value(model) + (expectedMarginChange)) / p.Value(model))
	elr, o = p.GetELROnMarketSell("1", BTCUSDT_PERP_SEC.SecurityID, model, 11, 110, 1, 1000)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	expectedMarginChange = ((110 - 100) * 1) - (0.0002 * 110 * 1)
	expectedElr = math.Log((p.Value(model) + (expectedMarginChange)) / p.Value(model))
	elr, o = p.GetELROnLimitAsk("1", BTCUSDT_PERP_SEC.SecurityID, model, 11, []float64{110}, []float64{1}, 1000)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	o.Quantity = math.Round(o.Quantity)

	// Add a sell order. Using o.quantity allows us to check if the returned order's quantity is correct too
	_, rej = account.NewOrder(&models.Order{
		OrderID:       "sell2",
		ClientOrderID: "sell2",
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: BTCUSDT_PERP_SEC.SecurityID},
			Exchange:   &constants.BITMEX,
			Symbol:     &types.StringValue{Value: "ETHUSD"},
		},
		OrderStatus:    models.PendingNew,
		OrderType:      models.Limit,
		Side:           models.Sell,
		TimeInForce:    models.Session,
		LeavesQuantity: o.Quantity,
		CumQuantity:    0,
		Price:          &types.DoubleValue{Value: o.Price},
	})
	if rej != nil {
		t.Fatalf(rej.String())
	}
	_, err = account.ConfirmNewOrder("sell2", "sell2")
	if err != nil {
		t.Fatal(err)
	}
	account.UpdateAskOrderQueue(BTCUSDT_PERP_SEC.SecurityID, "sell2", 1)

	elr = p.ExpectedLogReturn(model, 10)

	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	elr = p.GetELROnCancelAsk("1", BTCUSDT_PERP_SEC.SecurityID, "sell2", model, 10)
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}

	account.CancelOrder("sell2")
	if _, err := account.ConfirmCancelOrder("sell2"); err != nil {
		t.Fatal(err)
	}
	elr = p.ExpectedLogReturn(model, 10)
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}
}
