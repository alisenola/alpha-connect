package account

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"reflect"
	"testing"
	"time"
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

/*
var binanceAccount = &models.Account{
	AccountID: "299211",
	Exchange:  &constants.binance,
	ApiCredentials: &xchangerModels.APICredentials{
		APIKey:    "MpYkeK3pGP80gGiIrqWtLNwjJmyK2DTREYzNx8Cyc3AWTkl2T0iWnQEtdCIlvAoE",
		APISecret: "CJcJZEkktzhGzEdQhclfHcfJz5k01OY6n42MeF9B3oQWGqba3RrXEnG4bZktXQNu",
	},
}
*/

func TestBinanceAccountListener_OnOrderCancelRequest(t *testing.T) {
	as, executor := start(t, binanceAccount)

	orderID := fmt.Sprintf("%d", time.Now().UnixNano())
	// Test with one order
	res, err := as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   binanceAccount,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    binanceInstrument,
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
		Account:    binanceAccount,
		Instrument: instrument,
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
		Account:   binanceAccount,
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

func TestBinanceAccountListener_OnOrderStatusRequest(t *testing.T) {
	as, executor := start(t, binanceAccount)

	// Test with no account
	res, err := as.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   nil,
	}, 30*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}

	orderList, ok := res.(*messages.OrderList)
	if !ok {
		t.Fatalf("was expecting *messages.OrderList, got %s", reflect.TypeOf(res).String())
	}
	if orderList.Success {
		t.Fatalf("wasn't expecting success")
	}
	if orderList.RejectionReason != messages.InvalidAccount {
		t.Fatalf("was expecting %s, got %s", messages.InvalidAccount.String(), orderList.RejectionReason.String())
	}
	if len(orderList.Orders) > 0 {
		t.Fatalf("was expecting no order")
	}

	// Test with account
	res, err = as.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   binanceAccount,
	}, 30*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}

	orderList, ok = res.(*messages.OrderList)
	if !ok {
		t.Fatalf("was expecting *messages.OrderList, got %s", reflect.TypeOf(res).String())
	}
	if !orderList.Success {
		t.Fatalf("was expecting success: %s", orderList.RejectionReason.String())
	}

	// Test with instrument and order status
	res, err = as.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   binanceAccount,
		Filter: &messages.OrderFilter{
			OrderID:       nil,
			ClientOrderID: nil,
			Instrument:    binanceInstrument,
			OrderStatus:   &messages.OrderStatusValue{Value: models.New},
		},
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}

	orderList, ok = res.(*messages.OrderList)
	if !ok {
		t.Fatalf("was expecting *messages.OrderList, got %s", reflect.TypeOf(res).String())
	}
	if !orderList.Success {
		t.Fatalf("was expecting success: %s", orderList.RejectionReason.String())
	}
	if len(orderList.Orders) > 0 {
		t.Fatalf("was expecting no open order, got %d", len(orderList.Orders))
	}

	orderID := fmt.Sprintf("%d", time.Now().UnixNano())
	// Test with one order
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   binanceAccount,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    binanceInstrument,
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

	// Test with instrument and order status
	res, err = as.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   binanceAccount,
		Filter: &messages.OrderFilter{
			OrderID:       nil,
			ClientOrderID: nil,
			Instrument:    binanceInstrument,
			OrderStatus:   &messages.OrderStatusValue{Value: models.New},
		},
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}

	orderList, ok = res.(*messages.OrderList)
	if !ok {
		t.Fatalf("was expecting *messages.OrderList, got %s", reflect.TypeOf(res).String())
	}
	if !orderList.Success {
		t.Fatalf("was expecting success: %s", orderList.RejectionReason.String())
	}
	if len(orderList.Orders) != 1 {
		t.Fatalf("was expecting 1 open order, got %d", len(orderList.Orders))
	}
	order := orderList.Orders[0]
	if order.OrderStatus != models.New {
		t.Fatalf("order status not new")
	}
	if int(order.LeavesQuantity*10000) != 4 {
		t.Fatalf("was expecting leaves quantity of 0.0004, got %f", order.LeavesQuantity)
	}
	if int(order.CumQuantity) != 0 {
		t.Fatalf("was expecting cum quantity of 0")
	}
	if order.OrderType != models.Limit {
		t.Fatalf("was expecting limit order type")
	}
	if order.TimeInForce != models.GoodTillCancel {
		t.Fatalf("was expecting GoodTillCancel time in force")
	}
	if order.Side != models.Buy {
		t.Fatalf("was expecting buy side order")
	}

	res, err = as.Root.RequestFuture(executor, &messages.OrderCancelRequest{
		RequestID:  0,
		Account:    binanceAccount,
		Instrument: binanceInstrument,
		OrderID:    &types.StringValue{Value: order.OrderID},
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	mcResponse, ok := res.(*messages.OrderCancelResponse)
	if !ok {
		t.Fatalf("was expecting *messages.OrderCancelResponse, got %s", reflect.TypeOf(res).String())
	}
	if !mcResponse.Success {
		t.Fatalf("was expecting successful request: %s", response.RejectionReason.String())
	}

	time.Sleep(3 * time.Second)

	// Query order and check if got canceled
	res, err = as.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   binanceAccount,
		Filter: &messages.OrderFilter{
			OrderID:    &types.StringValue{Value: order.OrderID},
			Instrument: binanceInstrument,
		},
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}

	orderList, ok = res.(*messages.OrderList)
	if !ok {
		t.Fatalf("was expecting *messages.OrderList, got %s", reflect.TypeOf(res).String())
	}
	if !orderList.Success {
		t.Fatalf("was expecting success: %s", orderList.RejectionReason.String())
	}
	if len(orderList.Orders) != 1 {
		t.Fatalf("was expecting 1 order, got %d", len(orderList.Orders))
	}
	order = orderList.Orders[0]
	if order.OrderStatus != models.Canceled {
		t.Fatalf("order status not Canceled, but %s", order.OrderStatus)
	}
	if int(order.LeavesQuantity) != 0 {
		t.Fatalf("was expecting leaves quantity of 0")
	}
	if int(order.CumQuantity) != 0 {
		t.Fatalf("was expecting cum quantity of 0")
	}
}

func TestBinanceAccountListener_OnNewOrderSingleRequest(t *testing.T) {
	as, executor := start(t, binanceAccount)

	// Test Invalid account
	res, err := as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   nil,
		Order:     nil,
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	response, ok := res.(*messages.NewOrderSingleResponse)
	if !ok {
		t.Fatalf("was expecting *messages.NewOrderSingleResponse, got %s", reflect.TypeOf(res).String())
	}
	if response.Success {
		t.Fatalf("was expecting unsucessful request")
	}
	if response.RejectionReason != messages.InvalidAccount {
		t.Fatalf("was expecting %s got %s", messages.InvalidAccount.String(), response.RejectionReason.String())
	}

	// Test no order
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   binanceAccount,
		Order:     nil,
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	response, ok = res.(*messages.NewOrderSingleResponse)
	if !ok {
		t.Fatalf("was expecting *messages.NewOrderSingleResponse, got %s", reflect.TypeOf(res).String())
	}
	if response.Success {
		t.Fatalf("was expecting unsucessful request")
	}
	if response.RejectionReason != messages.InvalidRequest {
		t.Fatalf("was expecting %s got %s", messages.InvalidRequest.String(), response.RejectionReason.String())
	}

	// Test with one order
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   binanceAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    binanceInstrument,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.GoodTillCancel,
			Quantity:      0.0004,
			Price:         &types.DoubleValue{Value: 35000.},
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

	// Delete orders
	res, err = as.Root.RequestFuture(executor, &messages.OrderMassCancelRequest{
		RequestID: 0,
		Account:   binanceAccount,
		Filter: &messages.OrderFilter{
			Instrument: binanceInstrument,
		},
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	mcResponse, ok := res.(*messages.OrderMassCancelResponse)
	if !ok {
		t.Fatalf("was expecting *messages.NewOrderBulkResponse, got %s", reflect.TypeOf(res).String())
	}
	if !mcResponse.Success {
		t.Fatalf("was expecting successful request: %s", response.RejectionReason.String())
	}
}

func TestBinanceAccountListener_OnBalancesRequest(t *testing.T) {
	as, executor := start(t, binanceAccount)

	res, err := as.Root.RequestFuture(executor, &messages.BalancesRequest{
		RequestID: 0,
		Account:   binanceAccount,
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	balanceResponse, ok := res.(*messages.BalanceList)
	if !ok {
		t.Fatalf("was expecting *messages.BalanceList, got %s", reflect.TypeOf(res).String())
	}
	if !balanceResponse.Success {
		t.Fatalf("was expecting sucessful request: %s", balanceResponse.RejectionReason.String())
	}
	if len(balanceResponse.Balances) != 2 {
		t.Fatalf("was expecting 2 balance, got %d", len(balanceResponse.Balances))
	}
	fmt.Println(balanceResponse.Balances)
}

func TestBinanceAccountListener_OnGetPositionsLimit(t *testing.T) {
	as, executor := start(t, binanceAccount)

	// Market buy
	orderID := fmt.Sprintf("%d", time.Now().UnixNano())
	res, err := as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: binanceAccount,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    binanceInstrument,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			Price:         &types.DoubleValue{Value: 80000},
			Quantity:      0.0002,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	newOrderResponse := res.(*messages.NewOrderSingleResponse)
	if !newOrderResponse.Success {
		t.Fatalf("error creating new order: %s", newOrderResponse.RejectionReason.String())
	}

	time.Sleep(3 * time.Second)

	checkBalances(t, as, executor, binanceAccount)

	fmt.Println("CLOSING")
	// Close position
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: binanceAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    binanceInstrument,
			OrderType:     models.Limit,
			OrderSide:     models.Sell,
			Price:         &types.DoubleValue{Value: 42000},
			Quantity:      0.0002,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	checkBalances(t, as, executor, binanceAccount)
}

func TestBinanceAccountListener_OnOrderReplaceRequest(t *testing.T) {
	as, executor := start(t, binanceAccount)

	orderID := uuid.NewV1().String()
	// Post one order
	res, err := as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   binanceAccount,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    instrument,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      1.,
			Price:         &types.DoubleValue{Value: 35000.},
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

	time.Sleep(1 * time.Second)

	// Test replace quantity
	res, err = as.Root.RequestFuture(executor, &messages.OrderReplaceRequest{
		RequestID:  0,
		Account:    binanceAccount,
		Instrument: instrument,
		Update: &messages.OrderUpdate{
			OrderID:           nil,
			OrigClientOrderID: &types.StringValue{Value: orderID},
			Quantity:          &types.DoubleValue{Value: 2},
			Price:             nil,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	replaceResponse, ok := res.(*messages.OrderReplaceResponse)
	if !ok {
		t.Fatalf("was expecting *messages.OrderReplaceResponse. got %s", reflect.TypeOf(res).String())
	}
	if !replaceResponse.Success {
		t.Fatalf("was expecting successful request: %s", replaceResponse.RejectionReason.String())
	}

	time.Sleep(1 * time.Second)
	// Fetch orders
	res, err = as.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   binanceAccount,
		Filter: &messages.OrderFilter{
			ClientOrderID: &types.StringValue{Value: orderID},
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
		t.Fatalf("was expecting success: %s", orderList.RejectionReason.String())
	}
	if len(orderList.Orders) != 1 {
		t.Fatalf("was expecting one order, got %d", len(orderList.Orders))
	}
	if orderList.Orders[0].LeavesQuantity != 2 {
		t.Fatalf("was expecting quantity of 2")
	}
}

func TestBinanceAccountListener_OnGetPositionsMarket(t *testing.T) {
	as, executor := start(t, binanceAccount)
	// Market buy 1 contract
	orderID := fmt.Sprintf("%d", time.Now().UnixNano())
	res, err := as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: binanceAccount,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    instrument,
			OrderType:     models.Market,
			OrderSide:     models.Buy,
			Quantity:      0.005,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	newOrderResponse := res.(*messages.NewOrderSingleResponse)
	if !newOrderResponse.Success {
		t.Fatalf("error creating new order: %s", newOrderResponse.RejectionReason.String())
	}

	time.Sleep(3 * time.Second)

	checkBalances(t, as, executor, binanceAccount)

	fmt.Println("CHECK 1 GOOD")

	// Market sell 1 contract
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: binanceAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument,
			OrderType:     models.Market,
			OrderSide:     models.Sell,
			TimeInForce:   models.Session,
			Quantity:      0.001,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	checkBalances(t, as, executor, binanceAccount)

	fmt.Println("CLOSING")
	// Close position
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: binanceAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument,
			OrderType:     models.Market,
			OrderSide:     models.Sell,
			TimeInForce:   models.Session,
			Quantity:      0.001,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)
	checkBalances(t, as, executor, binanceAccount)

	// Close position
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: binanceAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument,
			OrderType:     models.Market,
			OrderSide:     models.Sell,
			TimeInForce:   models.Session,
			Quantity:      0.001,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)
	checkBalances(t, as, executor, binanceAccount)

	// Close position
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: binanceAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument,
			OrderType:     models.Market,
			OrderSide:     models.Sell,
			TimeInForce:   models.Session,
			Quantity:      0.001,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)
	checkBalances(t, as, executor, binanceAccount)

	// Close position
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: binanceAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument,
			OrderType:     models.Market,
			OrderSide:     models.Sell,
			TimeInForce:   models.Session,
			Quantity:      0.001,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)
	checkBalances(t, as, executor, binanceAccount)

}
