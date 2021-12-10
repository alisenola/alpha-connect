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

var ftxInstrument = &models.Instrument{
	Exchange: &constants.FTX,
	Symbol:   &types.StringValue{Value: "BTC-PERP"},
}

var ftxAccount = &models.Account{
	Name:     "299211",
	Exchange: &constants.FTX,
	ApiCredentials: &xchangerModels.APICredentials{
		APIKey:    "ej2YRRJMMwQD2qjOaRQnB18K7EWTpy3fRTP1ZZoX",
		APISecret: "MjJSlt1ix9OLpQLzkaQPlD_-y-N7c_2-6ZQJiwlY",
	},
}

/*
var FBinanceAccount = &models.Account{
	AccountID: "299211",
	Exchange:  &constants.FBINANCE,
	ApiCredentials: &xchangerModels.APICredentials{
		APIKey:    "MpYkeK3pGP80gGiIrqWtLNwjJmyK2DTREYzNx8Cyc3AWTkl2T0iWnQEtdCIlvAoE",
		APISecret: "CJcJZEkktzhGzEdQhclfHcfJz5k01OY6n42MeF9B3oQWGqba3RrXEnG4bZktXQNu",
	},
}
*/

func TestftxAccountListener_OnOrderCancelRequest(t *testing.T) {
	as, executor := start(t, ftxAccount)

	orderID := fmt.Sprintf("%d", time.Now().UnixNano())
	// Test with one order
	res, err := as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   ftxAccount,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    ftxInstrument,
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

	// Now cancel from the account

	res, err = as.Root.RequestFuture(executor, &messages.OrderCancelRequest{
		RequestID:     0,
		Account:       ftxAccount,
		Instrument:    ftxInstrument,
		ClientOrderID: &types.StringValue{Value: orderID},
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

	time.Sleep(2 * time.Second)
}

func TestftxAccountListener_OnOrderStatusRequest(t *testing.T) {
	as, executor := start(t, ftxAccount)

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
		Account:   ftxAccount,
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

	// Test with ftxInstrument and order status
	res, err = as.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   ftxAccount,
		Filter: &messages.OrderFilter{
			OrderID:       nil,
			ClientOrderID: nil,
			Instrument:    ftxInstrument,
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
		Account:   ftxAccount,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    ftxInstrument,
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

	// Test with ftxInstrument and order status
	res, err = as.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   ftxAccount,
		Filter: &messages.OrderFilter{
			OrderID:       nil,
			ClientOrderID: nil,
			Instrument:    ftxInstrument,
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
	if order.LeavesQuantity != 0.001 {
		t.Fatalf("was expecting leaves quantity of 1")
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
		Account:    ftxAccount,
		Instrument: ftxInstrument,
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
	/*
		// Now delete
		res, err = as.Root.RequestFuture(executor, &messages.OrderMassCancelRequest{
			RequestID: 0,
			Account:   ftxAccount,
			Filter: &messages.OrderFilter{
				Instrument: ftxInstrument,
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
	*/

	time.Sleep(3 * time.Second)

	// Query order and check if got canceled
	res, err = as.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   ftxAccount,
		Filter: &messages.OrderFilter{
			OrderID:    &types.StringValue{Value: order.OrderID},
			Instrument: ftxInstrument,
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

func TestftxAccountListener_OnNewOrderSingleRequest(t *testing.T) {
	as, executor := start(t, ftxAccount)

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
		Account:   ftxAccount,
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
		Account:   ftxAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    ftxInstrument,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.GoodTillCancel,
			Quantity:      0.001,
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
		Account:   ftxAccount,
		Filter: &messages.OrderFilter{
			Instrument: ftxInstrument,
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

func TestftxAccountListener_OnBalancesRequest(t *testing.T) {
	as, executor := start(t, ftxAccount)

	res, err := as.Root.RequestFuture(executor, &messages.BalancesRequest{
		RequestID: 0,
		Account:   ftxAccount,
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
		t.Fatalf("was expecting one balance, got %d", len(balanceResponse.Balances))
	}
	fmt.Println(balanceResponse.Balances)
}

func TestftxAccountListener_OnGetPositionsLimit(t *testing.T) {
	as, executor := start(t, ftxAccount)

	// Market buy 2 contract
	orderID := fmt.Sprintf("%d", time.Now().UnixNano())
	res, err := as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: ftxAccount,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    ftxInstrument,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			Price:         &types.DoubleValue{Value: 50000},
			Quantity:      0.002,
			TimeInForce:   models.GoodTillCancel,
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

	checkPositions(t, as, executor, ftxAccount, ftxInstrument)

	// Market sell 1 contract
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: ftxAccount,
		Order: &messages.NewOrder{
			ClientOrderID:         uuid.NewV1().String(),
			Instrument:            ftxInstrument,
			OrderType:             models.Limit,
			OrderSide:             models.Sell,
			Price:                 &types.DoubleValue{Value: 42000},
			Quantity:              0.001,
			TimeInForce:           models.GoodTillCancel,
			ExecutionInstructions: []models.ExecutionInstruction{models.ReduceOnly},
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	checkPositions(t, as, executor, ftxAccount, ftxInstrument)

	fmt.Println("CLOSING")
	// Close position
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: ftxAccount,
		Order: &messages.NewOrder{
			ClientOrderID:         uuid.NewV1().String(),
			Instrument:            ftxInstrument,
			OrderType:             models.Limit,
			OrderSide:             models.Sell,
			Price:                 &types.DoubleValue{Value: 42000},
			Quantity:              0.001,
			TimeInForce:           models.GoodTillCancel,
			ExecutionInstructions: []models.ExecutionInstruction{models.ReduceOnly},
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	// Only check balances at the end, due to how FTX realize the PnL every
	// second
	checkPositions(t, as, executor, ftxAccount, ftxInstrument)
	checkBalances(t, as, executor, ftxAccount)
}

func TestftxAccountListener_OnOrderReplaceRequest(t *testing.T) {
	as, executor := start(t, ftxAccount)

	orderID := uuid.NewV1().String()
	// Post one order
	res, err := as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   ftxAccount,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    ftxInstrument,
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

	time.Sleep(1 * time.Second)

	// Test replace quantity
	res, err = as.Root.RequestFuture(executor, &messages.OrderReplaceRequest{
		RequestID:  0,
		Account:    ftxAccount,
		Instrument: ftxInstrument,
		Update: &messages.OrderUpdate{
			OrderID:           nil,
			OrigClientOrderID: &types.StringValue{Value: orderID},
			Quantity:          &types.DoubleValue{Value: 0.002},
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
		Account:   ftxAccount,
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
	if orderList.Orders[0].LeavesQuantity != 0.002 {
		t.Fatalf("was expecting quantity of 0.002")
	}

	// Delete orders
	res, err = as.Root.RequestFuture(executor, &messages.OrderMassCancelRequest{
		RequestID: 0,
		Account:   ftxAccount,
		Filter: &messages.OrderFilter{
			Instrument: ftxInstrument,
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

func TestftxAccountListener_OnGetPositionsMarket(t *testing.T) {
	as, executor := start(t, ftxAccount)

	// Market buy 1 contract
	orderID := fmt.Sprintf("%d", time.Now().UnixNano())
	res, err := as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: ftxAccount,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    ftxInstrument,
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

	checkPositions(t, as, executor, ftxAccount, ftxInstrument)
	fmt.Println("CHECK 1 GOOD")

	// Market sell 1 contract
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: ftxAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    ftxInstrument,
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

	checkPositions(t, as, executor, ftxAccount, ftxInstrument)

	fmt.Println("CLOSING")
	// Close position
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: ftxAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    ftxInstrument,
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

	checkPositions(t, as, executor, ftxAccount, ftxInstrument)

	// Close position
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: ftxAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    ftxInstrument,
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

	checkPositions(t, as, executor, ftxAccount, ftxInstrument)

	// Close position
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: ftxAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    ftxInstrument,
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

	checkPositions(t, as, executor, ftxAccount, ftxInstrument)

	// Close position
	res, err = as.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: ftxAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    ftxInstrument,
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

	checkPositions(t, as, executor, ftxAccount, ftxInstrument)
	checkBalances(t, as, executor, ftxAccount)
}
