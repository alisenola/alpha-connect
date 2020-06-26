package portfolio

import (
	"fmt"
	"gitlab.com/alphaticks/xchanger/constants"
	"math"
	"testing"
)

var bitmex = constants.BITMEX

func TestBitmexPortfolio_UpdateContract(t *testing.T) {
	priceModels := make(map[uint64]PriceModel)
	priceModels[bitmexInstruments[0].ID()] = NewConstantPriceModel(100.)
	priceModels[bitmexInstruments[1].ID()] = NewConstantPriceModel(10.)

	buyTradeModels := make(map[uint64]BuyTradeModel)
	buyTradeModels[bitmexInstruments[0].ID()] = NewConstantTradeModel(10)
	buyTradeModels[bitmexInstruments[1].ID()] = NewConstantTradeModel(10)

	sellTradeModels := make(map[uint64]SellTradeModel)
	sellTradeModels[bitmexInstruments[0].ID()] = NewConstantTradeModel(10)
	sellTradeModels[bitmexInstruments[1].ID()] = NewConstantTradeModel(10)

	p, err := NewBitmexPortfolio(bitmexInstruments, NewConstantPriceModel(1000.), priceModels, buyTradeModels, sellTradeModels, 1000)
	if err != nil {
		t.Fatal(err)
	}
	p.BuyContract(bitmexInstruments[0].ID(), 1, 99)
	val := p.GetContract(bitmexInstruments[0].ID())
	if val != 1 {
		t.Fatal("incorrect asset qty")
	}
}

func TestBitmexPortfolio_ELR(t *testing.T) {

	priceModels := make(map[uint64]PriceModel)
	priceModels[bitmexInstruments[0].ID()] = NewConstantPriceModel(100.)
	priceModels[bitmexInstruments[1].ID()] = NewConstantPriceModel(10.)

	buyTradeModels := make(map[uint64]BuyTradeModel)
	buyTradeModels[bitmexInstruments[0].ID()] = NewConstantTradeModel(10)
	buyTradeModels[bitmexInstruments[1].ID()] = NewConstantTradeModel(10)

	sellTradeModels := make(map[uint64]SellTradeModel)
	sellTradeModels[bitmexInstruments[0].ID()] = NewConstantTradeModel(10)
	sellTradeModels[bitmexInstruments[1].ID()] = NewConstantTradeModel(10)

	ep, err := NewBitmexPortfolio(bitmexInstruments, NewConstantPriceModel(100.), priceModels, buyTradeModels, sellTradeModels, 1000)
	if err != nil {
		t.Fatal(err)
	}
	ep.SetWallet(0.08)
	availableMargin := ep.AvailableMargin(1.)

	p := NewPortfolio(map[uint64]ExchangePortfolio{bitmex: ep}, 1000)

	order := Order{
		Price:    90,
		Quantity: 10,
		Queue:    1,
	}
	expectedMarginChange := ((1./90 - 1./100) * 1. * 7.2) + (0.00025 * (1. / 90) * 7.2)
	expectedElr := math.Log((p.Value() + (expectedMarginChange * 100)) / p.Value())

	elr, o := p.ELROnLimitBid(10, bitmex, bitmexInstruments[0].ID(), []float64{90}, []float64{1}, availableMargin)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	order.Quantity = o.Quantity

	ep.SetBidOrder(bitmexInstruments[0].ID(), order)

	elr = p.ExpectedLogReturn(10)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	elr = p.ELROnCancelBid(10, bitmex, bitmexInstruments[0].ID())
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}
	ep.CancelBidOrder(bitmexInstruments[0].ID())

	elr = p.ExpectedLogReturn(20)
	expectedElr = 0.
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	order = Order{
		Price:    9,
		Quantity: 10,
		Queue:    1,
	}

	expectedMarginChange = ((10 - 9) * 0.000001 * 9) + (0.00025 * 9. * 9 * 0.000001)
	expectedElr = math.Log((p.Value() + (expectedMarginChange * 100)) / p.Value())

	elr, _ = p.ELROnLimitBid(10, bitmex, bitmexInstruments[1].ID(), []float64{9}, []float64{1}, availableMargin)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	ep.SetBidOrder(bitmexInstruments[1].ID(), order)

	elr = p.ExpectedLogReturn(30)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	elr = p.ELROnCancelBid(10, bitmex, bitmexInstruments[1].ID())
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}
	ep.CancelBidOrder(bitmexInstruments[1].ID())

	elr = p.ExpectedLogReturn(40)
	expectedElr = 0.
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	// Short 20, but only 8.8 available to short due to margin requirements
	order = Order{
		Price:    110,
		Quantity: 20,
		Queue:    1,
	}

	expectedMarginChange = ((1./100 - 1./110) * 1. * 8.8) + (0.00025 * (1. / 110) * 8.8)
	expectedElr = math.Log((p.Value() + (expectedMarginChange * 100)) / p.Value())

	elr, o = p.ELROnLimitAsk(10, bitmex, bitmexInstruments[0].ID(), []float64{110}, []float64{1}, availableMargin)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	order.Quantity = o.Quantity

	ep.SetAskOrder(bitmexInstruments[0].ID(), order)
	elr = p.ExpectedLogReturn(50)

	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	elr = p.ELROnCancelAsk(50, bitmex, bitmexInstruments[0].ID())
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}

	ep.CancelAskOrder(bitmexInstruments[0].ID())
	elr = p.ExpectedLogReturn(60)
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}

	order = Order{
		Price:    11,
		Quantity: 10,
		Queue:    1,
	}

	expectedMarginChange = ((11 - 10) * 0.000001 * 9) + (0.00025 * 11. * 9 * 0.000001)
	expectedElr = math.Log((p.Value() + (expectedMarginChange * 100)) / p.Value())

	elr, _ = p.ELROnLimitAsk(10, bitmex, bitmexInstruments[1].ID(), []float64{11}, []float64{1}, availableMargin)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	ep.SetAskOrder(bitmexInstruments[1].ID(), order)
	elr = p.ExpectedLogReturn(70)

	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	elr = p.ELROnCancelAsk(70, bitmex, bitmexInstruments[1].ID())
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}

	ep.CancelAskOrder(bitmexInstruments[1].ID())
	elr = p.ExpectedLogReturn(80)
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}
}

func TestBitmexPortfolio_BuySell(t *testing.T) {
	bitmex := constants.EXCHANGE_NAME_TO_EXCHANGE_ID[constants.BITMEX]

	priceModels := make(map[uint64]PriceModel)
	priceModels[bitmexInstruments[0].ID()] = NewConstantPriceModel(100.)
	priceModels[bitmexInstruments[1].ID()] = NewConstantPriceModel(10.)

	buyTradeModels := make(map[uint64]BuyTradeModel)
	buyTradeModels[bitmexInstruments[0].ID()] = NewConstantTradeModel(10)
	buyTradeModels[bitmexInstruments[1].ID()] = NewConstantTradeModel(10)

	sellTradeModels := make(map[uint64]SellTradeModel)
	sellTradeModels[bitmexInstruments[0].ID()] = NewConstantTradeModel(10)
	sellTradeModels[bitmexInstruments[1].ID()] = NewConstantTradeModel(10)

	ep, err := NewBitmexPortfolio(bitmexInstruments, NewConstantPriceModel(100.), priceModels, buyTradeModels, sellTradeModels, 1000)
	if err != nil {
		t.Fatal(err)
	}
	p := NewPortfolio(map[uint64]ExchangePortfolio{bitmex: ep}, 1000)

	ep.SetWallet(0.1)

	// Buy 100 at 100
	ep.BuyContract(bitmexInstruments[0].ID(), 10, 100)
	if math.Abs(ep.AvailableMargin(1.)-0.000025) > 0.000001 {
		t.Fatalf("was expecting available margin of %f, got %f", 0.000025, ep.AvailableMargin(1.))
	}

	// Look at elr on limit ask with an available margin of 0.
	elr, o := p.ELROnLimitAsk(10, bitmex, bitmexInstruments[0].ID(), []float64{110}, []float64{1}, 0.)
	if math.Abs(o.Quantity-9.) > 0.000001 {
		t.Fatalf("was expecting an order quantity of 9, got %f", o.Quantity)
	}
	expectedMarginChange := ((1./100 - 1./110) * 1. * 9) + (0.00025 * (1. / 110) * 9)
	expectedElr := math.Log((p.Value() + (expectedMarginChange * 100)) / p.Value())
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	// Sell 100 at 110
	ep.SellContract(bitmexInstruments[0].ID(), 10, 110)
	sampleMargin := ep.contracts[bitmexInstruments[0].ID()].sampleMargin[0]
	if math.Abs(sampleMargin) > 0.0000001 {
		t.Fatalf("was expecting sample margin of 0., got %f", sampleMargin)
	}

	// Add profit
	expectedWallet := 0.1 + (((1. / 100.) - (1. / 110.)) * 10)
	// Add maker rebate on buy
	expectedWallet += 0.00025 * (1. / 100.) * 10
	// Add maker rebate on sell
	expectedWallet += 0.00025 * (1. / 110.) * 10

	if math.Abs(ep.wallet-expectedWallet) > 0.000001 {
		t.Fatalf("was expecting wallet of %f, got %f", expectedWallet, ep.wallet)
	}

	if math.Abs(ep.AvailableMargin(1.)-expectedWallet) > 0.000001 {
		t.Fatalf("was expecting available margin of %f, got %f", expectedWallet, ep.AvailableMargin(1.))
	}

	ep.SetWallet(0.1)

	// Buy 10 at 100
	ep.BuyContract(bitmexInstruments[0].ID(), 10, 100)

	if math.Abs(ep.AvailableMargin(1.)-0.000025) > 0.000001 {
		t.Fatalf("was expecting available margin of %f, got %f", 0.000025, ep.AvailableMargin(1.))
	}

	// Sell 20 at 110
	ep.SellContract(bitmexInstruments[0].ID(), 20, 110)
	// Buy 10 at 90
	ep.BuyContract(bitmexInstruments[0].ID(), 10, 90)

	expectedWallet = 0.1 + (((1. / 100.) - (1. / 110.)) * 10) - (((1. / 110.) - (1. / 90.)) * 10)
	// Add maker rebate on buy
	expectedWallet += 0.00025 * (1. / 100.) * 10
	// Add maker rebate on sell
	expectedWallet += 0.00025 * (1. / 110.) * 20
	// Add maker rebate on buy
	expectedWallet += 0.00025 * (1. / 90.) * 10

	if math.Abs(ep.wallet-expectedWallet) > 0.000001 {
		t.Fatalf("was expecting wallet of %f, got %f", expectedWallet, ep.wallet)
	}
	if math.Abs(ep.AvailableMargin(1.)-expectedWallet) > 0.000001 {
		t.Fatalf("was expecting available margin of %f, got %f", expectedWallet, ep.AvailableMargin(1.))
	}

	ep.SetWallet(1.)

	// Buy 100 at 10
	ep.BuyContract(bitmexInstruments[1].ID(), 100, 10)

	// Look at elr on limit ask with an available margin of 0.
	elr, o = p.ELROnLimitAsk(10, bitmex, bitmexInstruments[1].ID(), []float64{11}, []float64{1}, 0.)
	if math.Abs(o.Quantity-9.) > 0.000001 {
		t.Fatalf("was expecting an order quantity of 9, got %f", o.Quantity)
	}
	expectedMarginChange = ((11 - 10) * 1. * 0.000001 * 9) + (0.00025 * 11 * 0.000001 * 9)
	expectedElr = math.Log((p.Value() + (expectedMarginChange * 100)) / p.Value())
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	// Sell 100 at 11
	ep.SellContract(bitmexInstruments[1].ID(), 100, 11)

	expectedWallet = 1. + ((11. - 10.) * 100 * 0.000001)
	// Add maker rebate on buy
	expectedWallet += 0.00025 * 100. * 10 * 0.000001
	// Add maker rebate on sell
	expectedWallet += 0.00025 * 100. * 11 * 0.000001

	if math.Abs(ep.wallet-expectedWallet) > 0.0000001 {
		t.Fatalf("was expecting wallet of %f, got %f", expectedWallet, ep.wallet)
	}
	if math.Abs(ep.AvailableMargin(1.)-expectedWallet) > 0.000001 {
		t.Fatalf("was expecting available margin of %f, got %f", expectedWallet, ep.AvailableMargin(1.))
	}

	ep.SetWallet(1.)

	// Buy 100 at 10
	ep.BuyContract(bitmexInstruments[1].ID(), 100, 10)
	// Sell 200 at 11
	ep.SellContract(bitmexInstruments[1].ID(), 200, 11)
	// Buy 100 at 9
	ep.BuyContract(bitmexInstruments[1].ID(), 100, 9)

	expectedWallet = 1. + ((11. - 10.) * 100 * 0.000001) - ((9. - 11.) * 100 * 0.000001)
	// Add maker rebate on buy
	expectedWallet += 0.00025 * 100. * 10 * 0.000001
	// Add maker rebate on sell
	expectedWallet += 0.00025 * 200. * 11 * 0.000001
	// Add maker rebate on buy
	expectedWallet += 0.00025 * 100. * 9 * 0.000001

	if math.Abs(ep.wallet-expectedWallet) > 0.000001 {
		t.Fatalf("was expecting wallet of %f, got %f", expectedWallet, ep.wallet)
	}
	if math.Abs(ep.AvailableMargin(1.)-expectedWallet) > 0.000001 {
		t.Fatalf("was expecting available margin of %f, got %f", expectedWallet, ep.AvailableMargin(1.))
	}
}

func TestBitmexPortfolio_Precision(t *testing.T) {
	// Change bid order
	// Change fee
	// Change sample match

	priceModels := make(map[uint64]PriceModel)
	priceModels[bitmexInstruments[0].ID()] = NewConstantPriceModel(100.)
	priceModels[bitmexInstruments[1].ID()] = NewConstantPriceModel(10.)

	buyTradeModels := make(map[uint64]BuyTradeModel)
	buyTradeModels[bitmexInstruments[0].ID()] = NewConstantTradeModel(10)
	buyTradeModels[bitmexInstruments[1].ID()] = NewConstantTradeModel(10)

	sellTradeModels := make(map[uint64]SellTradeModel)
	sellTradeModels[bitmexInstruments[0].ID()] = NewConstantTradeModel(10)
	sellTradeModels[bitmexInstruments[1].ID()] = NewConstantTradeModel(10)

	ep, err := NewBitmexPortfolio(bitmexInstruments, NewConstantPriceModel(100.), priceModels, buyTradeModels, sellTradeModels, 1000)
	if err != nil {
		t.Fatal(err)
	}

	order1 := Order{
		Price:    100,
		Quantity: 1,
		Queue:    1,
	}
	order2 := Order{
		Price:    98,
		Quantity: 1,
		Queue:    1,
	}

	ep.SetMakerFee(bitmexInstruments[0].ID(), -0.00025)
	ep.SetBidOrder(bitmexInstruments[0].ID(), order1)

	expectedVal := ep.contracts[bitmexInstruments[0].ID()].sampleSize[0]

	for i := 0; i < 100000; i++ {
		ep.SetMakerFee(bitmexInstruments[0].ID(), 0.01)
		ep.BuyContract(bitmexInstruments[0].ID(), 10, 100.)
		ep.SetBidOrder(bitmexInstruments[0].ID(), order2)
		// back to normal
		ep.SetMakerFee(bitmexInstruments[0].ID(), -0.00025)
		ep.BuyContract(bitmexInstruments[0].ID(), -10, 100.)
		ep.SetBidOrder(bitmexInstruments[0].ID(), order1)
	}
	val := ep.contracts[bitmexInstruments[0].ID()].sampleSize[0]

	if math.Abs(val-expectedVal) > 0.000001 {
		t.Fatalf("precision error of %f", val-expectedVal)
	}

	ep.CancelBidOrder(bitmexInstruments[0].ID())

	// Do the same for ETHUSD

	order1 = Order{
		Price:    10,
		Quantity: 1,
		Queue:    1,
	}
	order2 = Order{
		Price:    9.8,
		Quantity: 1,
		Queue:    1,
	}

	ep.SetMakerFee(bitmexInstruments[1].ID(), -0.00025)
	ep.SetBidOrder(bitmexInstruments[1].ID(), order1)

	expectedVal = ep.contracts[bitmexInstruments[1].ID()].sampleSize[0]

	for i := 0; i < 100000; i++ {
		ep.SetMakerFee(bitmexInstruments[1].ID(), 0.01)
		ep.BuyContract(bitmexInstruments[1].ID(), 10, 100.)
		ep.SetBidOrder(bitmexInstruments[1].ID(), order2)
		// back to normal
		ep.SetMakerFee(bitmexInstruments[1].ID(), -0.00025)
		ep.BuyContract(bitmexInstruments[1].ID(), -10, 100.)
		ep.SetBidOrder(bitmexInstruments[1].ID(), order1)
	}
	val = ep.contracts[bitmexInstruments[1].ID()].sampleSize[0]

	if math.Abs(val-expectedVal) > 0.000001 {
		t.Fatalf("precision error of %f", val-expectedVal)
	}
}

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
