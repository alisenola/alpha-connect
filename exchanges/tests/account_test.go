package tests

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	/*
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
	*/
}

var derivTests = []AccntTest{
	{
		account: &models.Account{
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
		},
		instrument: &models.Instrument{
			Exchange: &constants.DYDX,
			Symbol:   &types.StringValue{Value: "BTC-USD"},
		},
	},
}

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
			//OpenCancel(t, tc)
			MarketFill(t, tc)
			//CheckBalance(t, tc)
		})
	}
}

func TestCompute(t *testing.T) {
	size := 0.004
	entryPrice := 48404.285714

	newSize := 0.001
	newPrice := 48073.
	// fill of 0.001 at 48395

	cost := entryPrice*size + newPrice*newSize

	newEntryPrice := 48362.875000
	//241.690143 241.814375
	fmt.Println(cost/(size+newSize), newEntryPrice)
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

func OpenCancel(t *testing.T, tc AccntTest) {
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
		t.Fatalf("was expecting SecurityDefinitionResponse, got %s", reflect.TypeOf(res).String())
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

	clientOrderID := fmt.Sprintf("%d", time.Now().UnixNano())
	// Test with one order
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   tc.account,
		Order: &messages.NewOrder{
			ClientOrderID: clientOrderID,
			Instrument:    tc.instrument,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.GoodTillCancel,
			Quantity:      0.001,
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
		t.Fatalf("was expecting *messages.OrderList, got %s", reflect.TypeOf(res).String())
	}
	if !orderList.Success {
		t.Fatalf("was expecting sucessful request")
	}
	if len(orderList.Orders) != 1 {
		t.Fatal("was expecting one order")
	}
	fmt.Println(orderList.Orders[0])
	// Cancel from the account

	res, err = as.Root.RequestFuture(executor, &messages.OrderCancelRequest{
		RequestID:     0,
		Account:       tc.account,
		Instrument:    tc.instrument,
		ClientOrderID: &types.StringValue{Value: clientOrderID},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	mcResponse, ok := res.(*messages.OrderCancelResponse)
	if !ok {
		t.Fatalf("was expecting *messages.OrderCancelResponse, got %s", reflect.TypeOf(res).String())
	}
	if !mcResponse.Success {
		t.Fatalf("was expecting sucessful request, got %s", mcResponse.RejectionReason.String())
	}

	time.Sleep(2 * time.Second)

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
	orderList, ok = res.(*messages.OrderList)
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

func MarketFill(t *testing.T, tc AccntTest) {
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
		t.Fatalf("was expecting SecurityDefinitionResponse, got %s", reflect.TypeOf(res).String())
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

	clientOrderID := fmt.Sprintf("%d", time.Now().UnixNano())
	// Test with one order
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   tc.account,
		Order: &messages.NewOrder{
			ClientOrderID: clientOrderID,
			Instrument:    tc.instrument,
			OrderType:     models.Market,
			OrderSide:     models.Buy,
			TimeInForce:   models.GoodTillCancel,
			Quantity:      0.001,
			Price:         &types.DoubleValue{Value: 100000.},
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

	time.Sleep(5 * time.Second)
	checkPositions(t, as, executor, tc.account, tc.instrument)
	checkBalances(t, as, executor, tc.account)

	clientOrderID = fmt.Sprintf("%d", time.Now().UnixNano())
	// Test with one order
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   tc.account,
		Order: &messages.NewOrder{
			ClientOrderID: clientOrderID,
			Instrument:    tc.instrument,
			OrderType:     models.Market,
			OrderSide:     models.Buy,
			TimeInForce:   models.GoodTillCancel,
			Quantity:      0.001,
			Price:         &types.DoubleValue{Value: 100000.},
		},
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	response, ok = res.(*messages.NewOrderSingleResponse)
	if !ok {
		t.Fatalf("was expecting *messages.NewOrderSingleResponse, got %s", reflect.TypeOf(res).String())
	}
	if !response.Success {
		t.Fatalf("was expecting sucessful request: %s", response.RejectionReason.String())
	}

	time.Sleep(5 * time.Second)
	checkPositions(t, as, executor, tc.account, tc.instrument)
	checkBalances(t, as, executor, tc.account)

	clientOrderID = fmt.Sprintf("%d", time.Now().UnixNano())
	// Test with one order
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   tc.account,
		Order: &messages.NewOrder{
			ClientOrderID: clientOrderID,
			Instrument:    tc.instrument,
			OrderType:     models.Market,
			OrderSide:     models.Buy,
			TimeInForce:   models.GoodTillCancel,
			Quantity:      0.001,
			Price:         &types.DoubleValue{Value: 100000.},
		},
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	response, ok = res.(*messages.NewOrderSingleResponse)
	if !ok {
		t.Fatalf("was expecting *messages.NewOrderSingleResponse, got %s", reflect.TypeOf(res).String())
	}
	if !response.Success {
		t.Fatalf("was expecting sucessful request: %s", response.RejectionReason.String())
	}
	time.Sleep(5 * time.Second)
	checkPositions(t, as, executor, tc.account, tc.instrument)
	checkBalances(t, as, executor, tc.account)

	clientOrderID = fmt.Sprintf("%d", time.Now().UnixNano())
	// Test with one order
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   tc.account,
		Order: &messages.NewOrder{
			ClientOrderID: clientOrderID,
			Instrument:    tc.instrument,
			OrderType:     models.Market,
			OrderSide:     models.Sell,
			TimeInForce:   models.GoodTillCancel,
			Quantity:      0.003,
			Price:         &types.DoubleValue{Value: 100.},
		},
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	response, ok = res.(*messages.NewOrderSingleResponse)
	if !ok {
		t.Fatalf("was expecting *messages.NewOrderSingleResponse, got %s", reflect.TypeOf(res).String())
	}
	if !response.Success {
		t.Fatalf("was expecting sucessful request: %s", response.RejectionReason.String())
	}
	time.Sleep(5 * time.Second)

	checkPositions(t, as, executor, tc.account, tc.instrument)
	checkBalances(t, as, executor, tc.account)
	checkOrders(t, as, executor, tc.account)
}

func CheckBalance(t *testing.T, tc AccntTest) {
	as, executor, cleaner := StartExecutor(t, tc.instrument.Exchange, tc.account)
	defer cleaner()

	checkBalances(t, as, executor, tc.account)
	time.Sleep(10 * time.Second)
	checkBalances(t, as, executor, tc.account)
	time.Sleep(10 * time.Second)
	checkBalances(t, as, executor, tc.account)
	time.Sleep(10 * time.Second)
	checkBalances(t, as, executor, tc.account)
	time.Sleep(10 * time.Second)
	checkBalances(t, as, executor, tc.account)
}
