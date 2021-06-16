package account

import (
	"github.com/gogo/protobuf/types"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"math"
	"reflect"
	"testing"
	"time"
)

var instrument1 = &models.Instrument{
	SecurityID: &types.UInt64Value{Value: 5391998915988476130},
	Exchange:   &constants.BITMEX,
	Symbol:     &types.StringValue{Value: "XBTUSD"},
}
var instrument2 = &models.Instrument{
	SecurityID: &types.UInt64Value{Value: 11093839049553737303},
	Exchange:   &constants.BITMEX,
	Symbol:     &types.StringValue{Value: "ETHUSD"},
}

var BitmexAccount = &models.Account{
	AccountID: "299210",
	Exchange:  &constants.BITMEX,
	Credentials: &xchangerModels.APICredentials{
		APIKey:    "k5k6Mmaq3xe88Ph3fgIk9Vrt",
		APISecret: "0laIjZaKOMkJPtKy2ldJ18m4Dxjp66Vdim0k1-q4TXASZFZo",
	},
}

func TestBitmexAccountListener_OnOrderStatusRequest(t *testing.T) {
	// Test with no account
	res, err := As.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   nil,
	}, 10*time.Second).Result()

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
	res, err = As.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   BitmexAccount,
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

	// Test with instrument and order status
	res, err = As.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   BitmexAccount,
		Filter: &messages.OrderFilter{
			OrderID:       nil,
			ClientOrderID: nil,
			Instrument:    instrument1,
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

	// Test with one order
	res, err = As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   BitmexAccount,
		Order: &messages.NewOrder{
			ClientOrderID:         uuid.NewV1().String(),
			Instrument:            instrument1,
			OrderType:             models.Limit,
			ExecutionInstructions: []messages.ExecutionInstruction{messages.ParticipateDoNotInitiate},
			OrderSide:             models.Buy,
			TimeInForce:           models.Session,
			Quantity:              10.,
			Price:                 &types.DoubleValue{Value: 30000.},
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
	res, err = As.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   BitmexAccount,
		Filter: &messages.OrderFilter{
			OrderID:       nil,
			ClientOrderID: nil,
			Instrument:    instrument1,
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
	if int(order.LeavesQuantity) != 1 {
		t.Fatalf("was expecting leaves quantity of 1")
	}
	if int(order.CumQuantity) != 0 {
		t.Fatalf("was expecting cum quantity of 0")
	}
	if order.OrderType != models.Limit {
		t.Fatalf("was expecting limit order type")
	}
	if order.TimeInForce != models.Session {
		t.Fatalf("was expecting session time in force")
	}
	if order.Side != models.Buy {
		t.Fatalf("was expecting buy side order")
	}

	// Now delete
	res, err = As.Root.RequestFuture(executor, &messages.OrderMassCancelRequest{
		RequestID: 0,
		Account:   BitmexAccount,
		Filter: &messages.OrderFilter{
			Instrument: instrument1,
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

	time.Sleep(1 * time.Second)

	// Query order and check if got canceled
	res, err = As.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   BitmexAccount,
		Filter: &messages.OrderFilter{
			OrderID:    &types.StringValue{Value: order.OrderID},
			Instrument: instrument1,
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
		t.Fatalf("order status not Canceled")
	}
	if int(order.LeavesQuantity) != 0 {
		t.Fatalf("was expecting leaves quantity of 0")
	}
	if int(order.CumQuantity) != 0 {
		t.Fatalf("was expecting cum quantity of 0")
	}
}

func TestBitmexAccountListener_OnNewOrderSingleRequest(t *testing.T) {
	// Test Invalid account
	res, err := As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
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
	res, err = As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   BitmexAccount,
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
	res, err = As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   BitmexAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument1,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      1.,
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

	// Delete orders
	res, err = As.Root.RequestFuture(executor, &messages.OrderMassCancelRequest{
		RequestID: 0,
		Account:   BitmexAccount,
		Filter: &messages.OrderFilter{
			Instrument: instrument1,
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

func TestBitmexAccountListener_OnNewOrderBulkRequest(t *testing.T) {
	// Test Invalid account
	res, err := As.Root.RequestFuture(executor, &messages.NewOrderBulkRequest{
		RequestID: 0,
		Account:   nil,
		Orders:    nil,
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	response, ok := res.(*messages.NewOrderBulkResponse)
	if !ok {
		t.Fatalf("was expecting *messages.NewOrderBulkResponse, got %s", reflect.TypeOf(res).String())
	}
	if response.Success {
		t.Fatalf("was expecting unsucessful request")
	}
	if response.RejectionReason != messages.InvalidAccount {
		t.Fatalf("was expecting %s got %s", messages.InvalidAccount.String(), response.RejectionReason.String())
	}

	// Test no orders
	res, err = As.Root.RequestFuture(executor, &messages.NewOrderBulkRequest{
		RequestID: 0,
		Account:   BitmexAccount,
		Orders:    nil,
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	response, ok = res.(*messages.NewOrderBulkResponse)
	if !ok {
		t.Fatalf("was expecting *messages.NewOrderBulkResponse, got %s", reflect.TypeOf(res).String())
	}
	if !response.Success {
		t.Fatalf("was expecting successful request: %s", response.RejectionReason.String())
	}

	// Test with two orders diff symbols
	res, err = As.Root.RequestFuture(executor, &messages.NewOrderBulkRequest{
		RequestID: 0,
		Account:   BitmexAccount,
		Orders: []*messages.NewOrder{{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument1,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      1.,
			Price:         &types.DoubleValue{Value: 30000.},
		}, {
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument2,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      1.,
			Price:         &types.DoubleValue{Value: 30000.},
		}},
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	response, ok = res.(*messages.NewOrderBulkResponse)
	if !ok {
		t.Fatalf("was expecting *messages.NewOrderBulkResponse, got %s", reflect.TypeOf(res).String())
	}
	if response.Success {
		t.Fatalf("was expecting unsuccessful")
	}
	if response.RejectionReason != messages.DifferentSymbols {
		t.Fatalf("was expecting %s, got %s", messages.DifferentSymbols.String(), response.RejectionReason.String())
	}

	order1ClID := uuid.NewV1().String()
	order2ClID := uuid.NewV1().String()
	// Test with two orders same symbol diff price
	res, err = As.Root.RequestFuture(executor, &messages.NewOrderBulkRequest{
		RequestID: 0,
		Account:   BitmexAccount,
		Orders: []*messages.NewOrder{{
			ClientOrderID: order1ClID,
			Instrument:    instrument1,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      2,
			Price:         &types.DoubleValue{Value: 30000.},
		}, {
			ClientOrderID: order2ClID,
			Instrument:    instrument1,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      1,
			Price:         &types.DoubleValue{Value: 30010.},
		}},
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	response, ok = res.(*messages.NewOrderBulkResponse)
	if !ok {
		t.Fatalf("was expecting *messages.NewOrderBulkResponse, got %s", reflect.TypeOf(res).String())
	}
	if !response.Success {
		t.Fatalf("was expecting successful request: %s", response.RejectionReason.String())
	}
	time.Sleep(1 * time.Second)

	// Query order
	res, err = As.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   BitmexAccount,
		Filter: &messages.OrderFilter{
			ClientOrderID: &types.StringValue{Value: order1ClID},
			Instrument:    instrument1,
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
		t.Fatalf("was expecting 1 order, got %d", len(orderList.Orders))
	}
	order1 := orderList.Orders[0]
	if order1.OrderStatus != models.New {
		t.Fatalf("order status not new %s", order1.OrderStatus.String())
	}
	if int(order1.LeavesQuantity) != 2 {
		t.Fatalf("was expecting leaves quantity of 2")
	}
	if int(order1.CumQuantity) != 0 {
		t.Fatalf("was expecting cum quantity of 0")
	}
	if order1.OrderType != models.Limit {
		t.Fatalf("was expecting limit order type")
	}
	if order1.TimeInForce != models.Session {
		t.Fatalf("was expecting session time in force")
	}
	if order1.Side != models.Buy {
		t.Fatalf("was expecting buy side order")
	}

	// Query order
	res, err = As.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   BitmexAccount,
		Filter: &messages.OrderFilter{
			ClientOrderID: &types.StringValue{Value: order2ClID},
			Instrument:    instrument1,
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
	order2 := orderList.Orders[0]
	if order2.OrderStatus != models.New {
		t.Fatalf("order status not new")
	}
	if int(order2.LeavesQuantity) != 1 {
		t.Fatalf("was expecting leaves quantity of 2")
	}
	if int(order2.CumQuantity) != 0 {
		t.Fatalf("was expecting cum quantity of 0")
	}
	if order2.OrderType != models.Limit {
		t.Fatalf("was expecting limit order type")
	}
	if order2.TimeInForce != models.Session {
		t.Fatalf("was expecting session time in force")
	}
	if order2.Side != models.Buy {
		t.Fatalf("was expecting buy side order")
	}

	// Delete orders
	res, err = As.Root.RequestFuture(executor, &messages.OrderMassCancelRequest{
		RequestID: 0,
		Account:   BitmexAccount,
		Filter: &messages.OrderFilter{
			Instrument: instrument1,
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

func TestBitmexAccountListener_OnOrderReplaceRequest(t *testing.T) {
	orderID := uuid.NewV1().String()
	// Post one order
	res, err := As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   BitmexAccount,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    instrument1,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      1.,
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
	res, err = As.Root.RequestFuture(executor, &messages.OrderReplaceRequest{
		RequestID:  0,
		Account:    BitmexAccount,
		Instrument: instrument1,
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
	res, err = As.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   BitmexAccount,
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
		t.Fatalf("was expecting one order, ogt %d", len(orderList.Orders))
	}
	if orderList.Orders[0].LeavesQuantity != 2 {
		t.Fatalf("was expecting quantity of 2")
	}
}

func TestBitmexAccountListener_OnOrderBulkReplaceRequest(t *testing.T) {
	order1ClID := uuid.NewV1().String()
	order2ClID := uuid.NewV1().String()
	res, err := As.Root.RequestFuture(executor, &messages.NewOrderBulkRequest{
		RequestID: 0,
		Account:   BitmexAccount,
		Orders: []*messages.NewOrder{{
			ClientOrderID: order1ClID,
			Instrument:    instrument1,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      1,
			Price:         &types.DoubleValue{Value: 30000.},
		}, {
			ClientOrderID: order2ClID,
			Instrument:    instrument1,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      1,
			Price:         &types.DoubleValue{Value: 30010.},
		}},
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	response, ok := res.(*messages.NewOrderBulkResponse)
	if !ok {
		t.Fatalf("was expecting *messages.NewOrderBulkResponse, got %s", reflect.TypeOf(res).String())
	}
	if !response.Success {
		t.Fatalf("was expecting successful request: %s", response.RejectionReason.String())
	}
	time.Sleep(1 * time.Second)

	// Test replace quantity
	res, err = As.Root.RequestFuture(executor, &messages.OrderBulkReplaceRequest{
		RequestID:  0,
		Account:    BitmexAccount,
		Instrument: instrument1,
		Updates: []*messages.OrderUpdate{
			{
				OrderID:           nil,
				OrigClientOrderID: &types.StringValue{Value: order1ClID},
				Quantity:          &types.DoubleValue{Value: 2},
				Price:             nil,
			}, {
				OrderID:           nil,
				OrigClientOrderID: &types.StringValue{Value: order2ClID},
				Quantity:          &types.DoubleValue{Value: 2},
				Price:             nil,
			},
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	replaceResponse, ok := res.(*messages.OrderBulkReplaceResponse)
	if !ok {
		t.Fatalf("was expecting *messages.OrderBulkReplaceResponse. got %s", reflect.TypeOf(res).String())
	}
	if !replaceResponse.Success {
		t.Fatalf("was expecting successful request: %s", replaceResponse.RejectionReason.String())
	}

	time.Sleep(1 * time.Second)

	// Fetch orders
	res, err = As.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   BitmexAccount,
		Filter: &messages.OrderFilter{
			ClientOrderID: &types.StringValue{Value: order1ClID},
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
		t.Fatalf("was expecting one order, ogt %d", len(orderList.Orders))
	}
	if orderList.Orders[0].LeavesQuantity != 2 {
		t.Fatalf("was expecting quantity of 2")
	}

	// Fetch orders
	res, err = As.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   BitmexAccount,
		Filter: &messages.OrderFilter{
			ClientOrderID: &types.StringValue{Value: order2ClID},
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
		t.Fatalf("was expecting one order, ogt %d", len(orderList.Orders))
	}
	if orderList.Orders[0].LeavesQuantity != 2 {
		t.Fatalf("was expecting quantity of 2")
	}
}

func TestBitmexAccountListener_OnBalancesRequest(t *testing.T) {
	res, err := As.Root.RequestFuture(executor, &messages.BalancesRequest{
		RequestID: 0,
		Account:   BitmexAccount,
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
	if len(balanceResponse.Balances) != 1 {
		t.Fatalf("was expecting one balance, got %d", len(balanceResponse.Balances))
	}
}

func checkBitmexBalances(t *testing.T) {
	// Now check balance
	bitmexExecutor := As.NewLocalPID("executor/bitmex_executor")

	res, err := As.Root.RequestFuture(executor, &messages.BalancesRequest{
		RequestID: 0,
		Account:   BitmexAccount,
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
	if len(balanceResponse.Balances) != 1 {
		t.Fatalf("was expecting one balance, got %d", len(balanceResponse.Balances))
	}
	bal1 := balanceResponse.Balances[0]

	res, err = As.Root.RequestFuture(bitmexExecutor, &messages.BalancesRequest{
		RequestID: 0,
		Account:   BitmexAccount,
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	balanceResponse, ok = res.(*messages.BalanceList)
	if !ok {
		t.Fatalf("was expecting *messages.BalanceList, got %s", reflect.TypeOf(res).String())
	}
	if !balanceResponse.Success {
		t.Fatalf("was expecting sucessful request: %s", balanceResponse.RejectionReason.String())
	}
	if len(balanceResponse.Balances) != 1 {
		t.Fatalf("was expecting one balance, got %d", len(balanceResponse.Balances))
	}
	bal2 := balanceResponse.Balances[0]

	if math.Abs(bal1.Quantity-bal2.Quantity) > 0.00001 {
		t.Fatalf("different balance %f:%f", bal1.Quantity, bal2.Quantity)
	}
}

func checkBitmexPositions(t *testing.T, instrument *models.Instrument) {
	// Request the same from bitmex directly
	bitmexExecutor := As.NewLocalPID("executor/bitmex_executor")

	res, err := As.Root.RequestFuture(bitmexExecutor, &messages.PositionsRequest{
		RequestID:  0,
		Account:    BitmexAccount,
		Instrument: instrument,
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	response, ok := res.(*messages.PositionList)
	if !ok {
		t.Fatalf("was expecting *messages.PositionList, got %s", reflect.TypeOf(res).String())
	}
	if !response.Success {
		t.Fatalf("was expecting sucessful request: %s", response.RejectionReason.String())
	}

	pos1 := response.Positions

	res, err = As.Root.RequestFuture(executor, &messages.PositionsRequest{
		RequestID:  0,
		Account:    BitmexAccount,
		Instrument: instrument,
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	response, ok = res.(*messages.PositionList)
	if !ok {
		t.Fatalf("was expecting *messages.NewOrderBulkResponse, got %s", reflect.TypeOf(res).String())
	}
	if !response.Success {
		t.Fatalf("was expecting sucessful request: %s", response.RejectionReason.String())
	}
	pos2 := response.Positions

	if len(pos1) != len(pos2) {
		t.Fatalf("got different number of positions: %d %d", len(pos1), len(pos2))
	}

	for i := range pos1 {
		p1 := pos1[i]
		p2 := pos2[i]
		// Compare the two
		if math.Abs(p1.Cost-p2.Cost) > 0.000001 {
			t.Fatalf("different cost %f:%f", p1.Cost, p2.Cost)
		}
		if math.Abs(p1.Quantity-p2.Quantity) > 0.00001 {
			t.Fatalf("different quantity %f:%f", p1.Quantity, p2.Quantity)
		}
	}
}

func TestBitmexAccountListener_OnGetPositions(t *testing.T) {

	// Market buy 1 contract
	res, err := As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: BitmexAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument2,
			OrderType:     models.Market,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      2.,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	newOrderResponse := res.(*messages.NewOrderSingleResponse)
	if !newOrderResponse.Success {
		t.Fatalf("error creating new order: %s", newOrderResponse.RejectionReason.String())
	}

	time.Sleep(1 * time.Second)

	checkBitmexPositions(t, instrument2)
	checkBitmexBalances(t)

	// Market sell 1 contract
	res, err = As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: BitmexAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument2,
			OrderType:     models.Market,
			OrderSide:     models.Sell,
			TimeInForce:   models.Session,
			Quantity:      1.,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	checkBitmexPositions(t, instrument2)
	checkBitmexBalances(t)

	// Close position
	res, err = As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: BitmexAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument2,
			OrderType:     models.Market,
			OrderSide:     models.Sell,
			TimeInForce:   models.Session,
			Quantity:      1.,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	checkBitmexPositions(t, instrument2)
	checkBitmexBalances(t)
}

func TestBitmexAccountListener_OnGetPositions_Inverse(t *testing.T) {

	// Market buy 2 contracts
	res, err := As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: BitmexAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument1,
			OrderType:     models.Market,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      2.,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	newOrderResponse := res.(*messages.NewOrderSingleResponse)
	if !newOrderResponse.Success {
		t.Fatalf("error creating new order: %s", newOrderResponse.RejectionReason.String())
	}

	time.Sleep(1 * time.Second)

	checkBitmexPositions(t, instrument1)
	checkBitmexBalances(t)

	// Market sell 1 contract
	res, err = As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: BitmexAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument1,
			OrderType:     models.Market,
			OrderSide:     models.Sell,
			TimeInForce:   models.Session,
			Quantity:      1.,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	checkBitmexPositions(t, instrument1)
	checkBitmexBalances(t)

	// Close position
	res, err = As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: BitmexAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument1,
			OrderType:     models.Market,
			OrderSide:     models.Sell,
			TimeInForce:   models.Session,
			Quantity:      1.,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	// Request the same from bitmex directly

	checkBitmexPositions(t, instrument1)
	checkBitmexBalances(t)
}

func TestBitmexAccountListener_ConfirmFillReplace(t *testing.T) {

	order1ClID := uuid.NewV1().String()
	// Test with two orders same symbol diff price
	res, err := As.Root.RequestFuture(executor, &messages.NewOrderBulkRequest{
		RequestID: 0,
		Account:   BitmexAccount,
		Orders: []*messages.NewOrder{{
			ClientOrderID: order1ClID,
			Instrument:    instrument1,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      1,
			Price:         &types.DoubleValue{Value: 100000.},
		}},
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	response, ok := res.(*messages.NewOrderBulkResponse)
	if !ok {
		t.Fatalf("was expecting *messages.NewOrderBulkResponse, got %s", reflect.TypeOf(res).String())
	}
	if !response.Success {
		t.Fatalf("was expecting successful request: %s", response.RejectionReason.String())
	}

	time.Sleep(1 * time.Second)

	// Send a replace right after
	res, err = As.Root.RequestFuture(executor, &messages.OrderReplaceRequest{
		RequestID:  0,
		Account:    BitmexAccount,
		Instrument: instrument1,
		Update: &messages.OrderUpdate{
			OrderID:           nil,
			OrigClientOrderID: &types.StringValue{Value: order1ClID},
			Quantity:          &types.DoubleValue{Value: 5},
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
	if replaceResponse.Success {
		t.Fatalf("was expecting unsuccessful request")
	}

	time.Sleep(1 * time.Second)

	// Query order
	res, err = As.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   BitmexAccount,
		Filter: &messages.OrderFilter{
			ClientOrderID: &types.StringValue{Value: order1ClID},
			Instrument:    instrument1,
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
		t.Fatalf("was expecting 1 order, got %d", len(orderList.Orders))
	}
	order1 := orderList.Orders[0]
	if order1.OrderStatus != models.Filled {
		t.Fatalf("order status not filled %s", order1.OrderStatus.String())
	}
	if int(order1.LeavesQuantity) != 0 {
		t.Fatalf("was expecting leaves quantity of 0")
	}
	if int(order1.CumQuantity) != 1 {
		t.Fatalf("was expecting cum quantity of 1")
	}
	if order1.OrderType != models.Limit {
		t.Fatalf("was expecting limit order type")
	}
	if order1.TimeInForce != models.Session {
		t.Fatalf("was expecting session time in force")
	}
	if order1.Side != models.Buy {
		t.Fatalf("was expecting buy side order")
	}
}
