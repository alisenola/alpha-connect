package tests

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/binance"
	"gitlab.com/alphaticks/xchanger/exchanges/fbinance"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"reflect"
	"testing"
	"time"
)

type AccntTest struct {
	account    *models.Account
	instrument *models.Instrument
}

var spotTests = []AccntTest{
	{
		account: &models.Account{
			Name:     "binance",
			Exchange: &constants.BINANCE,
			ApiCredentials: &xchangerModels.APICredentials{
				APIKey:    "DNySYXVSG7xrM7S8dGSTvLRmRGJIzoU80Uj78IpFOYxiI9veS54VCu8bxQNLloz2",
				APISecret: "OrIJH6qYynrVFgEi62cknUpf02MoA82l45ySfP3ZTKRPainFJzG377BoJJmTuwXv",
			},
		},
		instrument: &models.Instrument{
			Exchange: &constants.BINANCE,
			Symbol:   &types.StringValue{Value: "BTCUSDT"},
		},
	},
}

var derivTests = []AccntTest{}

func TestAllAccount(t *testing.T) {
	binance.EnableTestNet()
	fbinance.EnableTestNet()
	t.Parallel()
	for _, tc := range spotTests {
		tc := tc
		t.Run(tc.instrument.Exchange.Name, func(t *testing.T) {
			t.Parallel()
			SpotAccount(t, tc)
		})
	}
	for _, tc := range derivTests {
		tc := tc
		t.Run(tc.instrument.Exchange.Name, func(t *testing.T) {
			t.Parallel()
			DerivAccount(t, tc)
		})
	}
}

func SpotAccount(t *testing.T, tc AccntTest) {
	as, executor, cleaner := StartExecutor(t, tc.instrument.Exchange, tc.account)
	defer cleaner()

	// Get security def
	res, err := as.Root.RequestFuture(executor, &messages.SecurityDefinitionRequest{
		Instrument: tc.instrument,
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	sd, ok := res.(*messages.SecurityDefinitionResponse)
	if !ok {
		t.Fatalf("was expecting balance SecurityDefinitionRequest, got %s", reflect.TypeOf(res).String())
	}
	if !sd.Success {
		t.Fatal(sd.RejectionReason.String())
	}

	sec := sd.Security

	// Get balances
	res, err = as.Root.RequestFuture(executor, &messages.BalancesRequest{
		Asset:   nil,
		Account: tc.account,
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	bl, ok := res.(*messages.BalanceList)
	if !ok {
		t.Fatalf("was expecting balance list, got %s", reflect.TypeOf(res).String())
	}
	if !bl.Success {
		t.Fatal(bl.RejectionReason.String())
	}

	var available float64
	for _, bl := range bl.Balances {
		if bl.Asset.ID == sec.QuoteCurrency.ID {
			available = bl.Quantity
		}
	}
	if available == 0. {
		t.Fatal("quote balance of 0, cannot test")
	}

	// Get market data
	res, err = as.Root.RequestFuture(executor, &messages.MarketDataRequest{
		RequestID:   0,
		Instrument:  tc.instrument,
		Aggregation: models.L2,
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	v, ok := res.(*messages.MarketDataResponse)
	if !ok {
		t.Fatalf("was expecting *messages.SecurityList, got %s", reflect.TypeOf(res).String())
	}
	if !v.Success {
		t.Fatalf("was expecting success, go %s", v.RejectionReason.String())
	}

	orderID := fmt.Sprintf("%d", time.Now().UnixNano())
	// Test with one order
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   tc.account,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    tc.instrument,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.GoodTillCancel,
			Quantity:      0.00040,
			Price:         &types.DoubleValue{Value: 30000.},
		},
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	response, ok := res.(*messages.NewOrderSingleResponse)
	if !ok {
		t.Fatalf("was expecting *messages.NewOrderSingleResponse, got %s", reflect.TypeOf(res).String())
	}
	if !response.Success {
		t.Fatalf("was expecting sucessful request: %s", response.RejectionReason.String())
	}

	time.Sleep(2 * time.Second)
	// Cancel from the account

	res, err = as.Root.RequestFuture(executor, &messages.OrderCancelRequest{
		RequestID:  0,
		Account:    tc.account,
		Instrument: tc.instrument,
		OrderID:    &types.StringValue{Value: orderID},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	mcResponse, ok := res.(*messages.OrderCancelResponse)
	if !ok {
		t.Fatalf("was expecting *messages.OrderCancelResponse, got %s", reflect.TypeOf(res).String())
	}
	if !mcResponse.Success {
		t.Fatalf("was expecting sucessful request")
	}

	time.Sleep(1 * time.Second)

	res, err = as.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Account:   tc.account,
		Filter: &messages.OrderFilter{
			OrderStatus: &messages.OrderStatusValue{Value: models.New},
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	orderList, ok := res.(*messages.OrderList)
	if !ok {
		t.Fatalf("was expecting *messages.OrderCancelResponse, got %s", reflect.TypeOf(res).String())
	}
	if !orderList.Success {
		t.Fatalf("was expecting sucessful request")
	}
	if len(orderList.Orders) > 0 {
		t.Fatal("was expecting no order")
	}
}

func DerivAccount(t *testing.T, tc AccntTest) {
}
