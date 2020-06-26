package legacy

import (
	"fmt"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/models"
	"math"
	"testing"
)

var spotInstruments = []*exchanges.Instrument{
	&exchanges.Instrument{
		Exchange: constants.BINANCE,
		Pair: &models.Pair{
			Base:  &constants.ETHEREUM,
			Quote: &constants.TETHER,
		},
		TickPrecision: 100,
		LotPrecision:  100,
		Type:          exchanges.SPOT,
	},
}

func TestSpotPortfolio_SetAsset(t *testing.T) {
	binance := constants.EXCHANGE_NAME_TO_EXCHANGE_ID[constants.BINANCE]

	priceModels := make(map[uint32]PriceModel)
	priceModels[constants.ETHEREUM.ID] = NewConstantPriceModel(100.)
	priceModels[constants.TETHER.ID] = NewConstantPriceModel(1.)

	buyTradeModels := make(map[uint64]BuyTradeModel)
	buyTradeModels[spotInstruments[0].ID()] = NewConstantTradeModel(10)

	sellTradeModels := make(map[uint64]SellTradeModel)
	sellTradeModels[spotInstruments[0].ID()] = NewConstantTradeModel(10)

	p, err := NewSpotPortfolio(binance, spotInstruments, priceModels, buyTradeModels, sellTradeModels, 1000)
	if err != nil {
		t.Fatal(err)
	}
	p.SetAsset(constants.ETHEREUM.ID, 1.)
	val := p.GetAsset(constants.ETHEREUM.ID)
	if math.Abs(val-1.0) > gorderbook.EPSILON_QTY {
		t.Fatal("incorrect asset qty")
	}
}

func TestPortfolio_SetAskOrder(t *testing.T) {
	binance := constants.EXCHANGE_NAME_TO_EXCHANGE_ID[constants.BINANCE]

	priceModels := make(map[uint32]PriceModel)
	priceModels[constants.ETHEREUM.ID] = NewConstantPriceModel(100.)
	priceModels[constants.TETHER.ID] = NewConstantPriceModel(1.)

	buyTradeModels := make(map[uint64]BuyTradeModel)
	buyTradeModels[spotInstruments[0].ID()] = NewConstantTradeModel(10)

	sellTradeModels := make(map[uint64]SellTradeModel)
	sellTradeModels[spotInstruments[0].ID()] = NewConstantTradeModel(10)

	ep, err := NewSpotPortfolio(binance, spotInstruments, priceModels, buyTradeModels, sellTradeModels, 1000)
	if err != nil {
		t.Fatal(err)
	}
	priceModels[constants.ETHEREUM.ID] = NewConstantPriceModel(100)
	priceModels[constants.TETHER.ID] = NewConstantPriceModel(1.)
	p := NewPortfolio(map[uint64]ExchangePortfolio{binance: ep}, 1000)
	ep.SetAsset(constants.TETHER.ID, 100)

	order := Order{
		Price:    100,
		Quantity: 10,
		Queue:    1,
	}
	ep.SetAskOrder(spotInstruments[0].ID(), order)

	fmt.Println(p.ExpectedLogReturn(0))
}

func TestSpotPortfolio_SampleValueChange(t *testing.T) {
	// Change bid order
	// Change fee
	// Change sample match

	binance := constants.EXCHANGE_NAME_TO_EXCHANGE_ID[constants.BINANCE]
	priceModels := make(map[uint32]PriceModel)
	priceModels[constants.ETHEREUM.ID] = NewConstantPriceModel(100.)
	priceModels[constants.TETHER.ID] = NewConstantPriceModel(1.)

	buyTradeModels := make(map[uint64]BuyTradeModel)
	buyTradeModels[spotInstruments[0].ID()] = NewConstantTradeModel(10)

	sellTradeModels := make(map[uint64]SellTradeModel)
	sellTradeModels[spotInstruments[0].ID()] = NewConstantTradeModel(10)

	ep, err := NewSpotPortfolio(binance, spotInstruments, priceModels, buyTradeModels, sellTradeModels, 1000)
	if err != nil {
		t.Fatal(err)
	}
	ep.SetMakerFee(0.02)
	ep.SetAsset(constants.TETHER.ID, 10000.)

	order := Order{
		Price:    90,
		Quantity: 0.5,
		Queue:    1,
	}
	values := make([]float64, 1000, 1000)

	ep.SetBidOrder(spotInstruments[0].ID(), order)
	ep.AddSampleValues(10, values)

	// We expect to match 0.5 ETH at a fee of 0.02, so we expect to get (0.5 - 0.5 * 0.02) * price
	expectedBaseValueChange := (0.5 - 0.5*0.02) * 100.
	expectedQuoteValueChange := -0.5 * 90
	expectedValueChange := expectedBaseValueChange + expectedQuoteValueChange

	valueChange := ep.instruments[spotInstruments[0].ID()].SampleValueChange[0]

	if math.Abs(expectedValueChange-valueChange) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedValueChange, valueChange)
	}
}

func TestSpotPortfolio_Precision(t *testing.T) {
	// Change bid order
	// Change fee
	// Change sample match

	binance := constants.EXCHANGE_NAME_TO_EXCHANGE_ID[constants.BINANCE]
	priceModels := make(map[uint32]PriceModel)
	priceModels[constants.ETHEREUM.ID] = NewConstantPriceModel(100.)
	priceModels[constants.TETHER.ID] = NewConstantPriceModel(1.)

	buyTradeModels := make(map[uint64]BuyTradeModel)
	buyTradeModels[spotInstruments[0].ID()] = NewConstantTradeModel(10)

	sellTradeModels := make(map[uint64]SellTradeModel)
	sellTradeModels[spotInstruments[0].ID()] = NewConstantTradeModel(10)

	ep, err := NewSpotPortfolio(binance, spotInstruments, priceModels, buyTradeModels, sellTradeModels, 1000)
	if err != nil {
		t.Fatal(err)
	}
	ep.SetMakerFee(0.02)
	ep.SetAsset(constants.TETHER.ID, 10000.)

	order1 := Order{
		Price:    100,
		Quantity: 0.5,
		Queue:    1,
	}
	order2 := Order{
		Price:    98,
		Quantity: 0.5,
		Queue:    1,
	}

	ep.SetBidOrder(spotInstruments[0].ID(), order1)

	expectedVal := ep.instruments[spotInstruments[0].ID()].SampleValueChange[0]
	for i := 0; i < 100000; i++ {
		ep.SetMakerFee(0.01)
		ep.SetAsset(constants.TETHER.ID, 10005)
		ep.SetBidOrder(spotInstruments[0].ID(), order2)
		// back to normal
		ep.SetMakerFee(0.02)
		ep.SetBidOrder(spotInstruments[0].ID(), order1)
		ep.SetAsset(constants.TETHER.ID, 10000)
	}
	val := ep.instruments[spotInstruments[0].ID()].SampleValueChange[0]

	if math.Abs(val-expectedVal) > 0.000001 {
		t.Fatalf("precision error of %f", val-expectedVal)
	}
}

func TestSpotPortfolio_ELR(t *testing.T) {
	binance := constants.EXCHANGE_NAME_TO_EXCHANGE_ID[constants.BINANCE]

	priceModels := make(map[uint32]PriceModel)
	priceModels[constants.ETHEREUM.ID] = NewConstantPriceModel(100.)
	priceModels[constants.TETHER.ID] = NewConstantPriceModel(1.)

	buyTradeModels := make(map[uint64]BuyTradeModel)
	buyTradeModels[spotInstruments[0].ID()] = NewConstantTradeModel(10)

	sellTradeModels := make(map[uint64]SellTradeModel)
	sellTradeModels[spotInstruments[0].ID()] = NewConstantTradeModel(10)

	ep, err := NewSpotPortfolio(binance, spotInstruments, priceModels, buyTradeModels, sellTradeModels, 1000)
	if err != nil {
		t.Fatal(err)
	}
	ep.SetAsset(constants.TETHER.ID, 1000.)
	ep.SetMakerFee(0.002)

	p := NewPortfolio(map[uint64]ExchangePortfolio{binance: ep}, 1000)

	order := Order{
		Price:    90,
		Quantity: 10,
		Queue:    1,
	}
	expectedBaseChange := 9 - (9 * 0.002)
	expectedQuoteChange := 9 * 90.
	expectedValueChange := expectedBaseChange*100. - expectedQuoteChange

	expectedElr := math.Log((p.Value() + expectedValueChange) / p.Value())

	elr, o := p.ELROnLimitBid(10, binance, spotInstruments[0].ID(), []float64{90}, []float64{1}, 1000.)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	order.Quantity = o.Quantity

	ep.SetBidOrder(spotInstruments[0].ID(), order)

	elr = p.ExpectedLogReturn(10)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	elr = p.ELROnCancelBid(10, binance, spotInstruments[0].ID())
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}
	ep.CancelBidOrder(spotInstruments[0].ID())

	elr = p.ExpectedLogReturn(10)
	expectedElr = 0.
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	elr = p.ExpectedLogReturn(20)
	expectedElr = 0.
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	ep.SetAsset(constants.ETHEREUM.ID, 9.)
	// Sell 9, but only 8.8 available to short due to margin requirements
	order = Order{
		Price:    110,
		Quantity: 10,
		Queue:    1,
	}

	expectedBaseChange = 9
	expectedQuoteChange = 9 * 110. * (1 - 0.002)
	expectedValueChange = expectedQuoteChange - expectedBaseChange*100.

	expectedElr = math.Log((p.Value() + expectedValueChange) / p.Value())

	elr, o = p.ELROnLimitAsk(10, binance, spotInstruments[0].ID(), []float64{110}, []float64{1}, 9.)
	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}
	order.Quantity = o.Quantity

	ep.SetAskOrder(spotInstruments[0].ID(), order)
	elr = p.ExpectedLogReturn(50)

	if math.Abs(elr-expectedElr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", expectedElr, elr)
	}

	elr = p.ELROnCancelAsk(50, binance, spotInstruments[0].ID())
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}

	ep.CancelAskOrder(spotInstruments[0].ID())
	elr = p.ExpectedLogReturn(60)
	if math.Abs(elr) > 0.000001 {
		t.Fatalf("was expecting %f got %f", 0., elr)
	}
}
