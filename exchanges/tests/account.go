package tests

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/types"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"reflect"
	"testing"
	"time"
)

type AccountTest struct {
	Account                 *models.Account
	Instrument              *models.Instrument
	OrderCancelRequest      bool
	OrderStatusRequest      bool
	NewOrderSingleRequest   bool
	NewOrderBulkRequest     bool
	GetPositionsLimit       bool
	OrderReplaceRequest     bool
	OrderBulkReplaceRequest bool
	GetPositionsMarket      bool
}

type AccountTestCtx struct {
	as       *actor.ActorSystem
	executor *actor.PID
}

func AccntTest(t *testing.T, tc AccountTest) {
	as, executor, cleaner := StartExecutor(t, tc.Instrument.Exchange, tc.Account)
	defer cleaner()

	ctx := AccountTestCtx{
		as:       as,
		executor: executor,
	}
	// Get security def
	res, err := as.Root.RequestFuture(executor, &messages.SecurityDefinitionRequest{
		Instrument: tc.Instrument,
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
		Account: tc.Account,
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
		Instrument:  tc.Instrument,
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

	if tc.OrderCancelRequest {
		OrderCancelRequest(t, ctx, tc)
	}
	if tc.OrderStatusRequest {
		OrderStatusRequest(t, ctx, tc)
	}
	if tc.GetPositionsLimit {
		GetPositionsLimit(t, ctx, tc)
	}
	if tc.GetPositionsMarket {
		GetPositionsMarket(t, ctx, tc)
	}
	if tc.NewOrderSingleRequest {
		NewOrderSingleRequest(t, ctx, tc)
	}
	if tc.NewOrderBulkRequest {
		NewOrderBulkRequest(t, ctx, tc)
	}
	if tc.OrderReplaceRequest {
		OrderReplaceRequest(t, ctx, tc)
	}
	if tc.OrderBulkReplaceRequest {
		OrderBulkReplaceRequest(t, ctx, tc)
	}
}

func OrderCancelRequest(t *testing.T, ctx AccountTestCtx, tc AccountTest) {
	orderID := fmt.Sprintf("%d", time.Now().UnixNano())
	// Test with one order
	res, err := ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   tc.Account,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    tc.Instrument,
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
	// Cancel from the account

	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderCancelRequest{
		RequestID:     0,
		Account:       tc.Account,
		Instrument:    tc.Instrument,
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

	time.Sleep(1 * time.Second)

	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Account:   tc.Account,
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
		for _, o := range orderList.Orders {
			fmt.Println(o.OrderStatus)
		}
		t.Fatal("was expecting no order")
	}
}

func OrderStatusRequest(t *testing.T, ctx AccountTestCtx, tc AccountTest) {
	// Test with no account
	res, err := ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderStatusRequest{
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   tc.Account,
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   tc.Account,
		Filter: &messages.OrderFilter{
			OrderID:       nil,
			ClientOrderID: nil,
			Instrument:    tc.Instrument,
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   tc.Account,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    tc.Instrument,
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

	// Test with instrument and order status
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   tc.Account,
		Filter: &messages.OrderFilter{
			OrderID:       nil,
			ClientOrderID: nil,
			Instrument:    tc.Instrument,
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
	if int(order.LeavesQuantity*1000) != 1 {
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

	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderCancelRequest{
		RequestID:  0,
		Account:    tc.Account,
		Instrument: tc.Instrument,
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   tc.Account,
		Filter: &messages.OrderFilter{
			OrderID:    &types.StringValue{Value: order.OrderID},
			Instrument: tc.Instrument,
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

func NewOrderSingleRequest(t *testing.T, ctx AccountTestCtx, tc AccountTest) {
	// Test Invalid account
	res, err := ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderSingleRequest{
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   tc.Account,
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   tc.Account,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    tc.Instrument,
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderMassCancelRequest{
		RequestID: 0,
		Account:   tc.Account,
		Filter: &messages.OrderFilter{
			Instrument: tc.Instrument,
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

func NewOrderBulkRequest(t *testing.T, ctx AccountTestCtx, tc AccountTest) {
	// Test Invalid account
	res, err := ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderBulkRequest{
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderBulkRequest{
		RequestID: 0,
		Account:   tc.Account,
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderBulkRequest{
		RequestID: 0,
		Account:   tc.Account,
		Orders: []*messages.NewOrder{{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    tc.Instrument,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      1.,
			Price:         &types.DoubleValue{Value: 30000.},
		}, {
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    tc.Instrument,
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderBulkRequest{
		RequestID: 0,
		Account:   tc.Account,
		Orders: []*messages.NewOrder{{
			ClientOrderID: order1ClID,
			Instrument:    tc.Instrument,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      2,
			Price:         &types.DoubleValue{Value: 30000.},
		}, {
			ClientOrderID: order2ClID,
			Instrument:    tc.Instrument,
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   tc.Account,
		Filter: &messages.OrderFilter{
			ClientOrderID: &types.StringValue{Value: order1ClID},
			Instrument:    tc.Instrument,
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   tc.Account,
		Filter: &messages.OrderFilter{
			ClientOrderID: &types.StringValue{Value: order2ClID},
			Instrument:    tc.Instrument,
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderMassCancelRequest{
		RequestID: 0,
		Account:   tc.Account,
		Filter: &messages.OrderFilter{
			Instrument: tc.Instrument,
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

func GetPositionsLimit(t *testing.T, ctx AccountTestCtx, tc AccountTest) {
	// Market buy
	orderID := fmt.Sprintf("%d", time.Now().UnixNano())
	res, err := ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderSingleRequest{
		Account: tc.Account,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    tc.Instrument,
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

	checkBalances(t, ctx.as, ctx.executor, tc.Account)

	fmt.Println("CLOSING")
	// Close position
	_, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderSingleRequest{
		Account: tc.Account,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    tc.Instrument,
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

	checkBalances(t, ctx.as, ctx.executor, tc.Account)
}

func OrderReplaceRequest(t *testing.T, ctx AccountTestCtx, tc AccountTest) {
	orderID := uuid.NewV1().String()
	// Post one order
	res, err := ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   tc.Account,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    tc.Instrument,
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderReplaceRequest{
		RequestID:  0,
		Account:    tc.Account,
		Instrument: tc.Instrument,
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   tc.Account,
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

func OrderBulkReplaceRequest(t *testing.T, ctx AccountTestCtx, tc AccountTest) {
	order1ClID := uuid.NewV1().String()
	order2ClID := uuid.NewV1().String()
	res, err := ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderBulkRequest{
		RequestID: 0,
		Account:   tc.Account,
		Orders: []*messages.NewOrder{{
			ClientOrderID: order1ClID,
			Instrument:    tc.Instrument,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      1,
			Price:         &types.DoubleValue{Value: 30000.},
		}, {
			ClientOrderID: order2ClID,
			Instrument:    tc.Instrument,
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderBulkReplaceRequest{
		RequestID:  0,
		Account:    tc.Account,
		Instrument: tc.Instrument,
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   tc.Account,
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
	res, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   tc.Account,
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

func GetPositionsMarket(t *testing.T, ctx AccountTestCtx, tc AccountTest) {
	// Market buy 1 contract
	orderID := fmt.Sprintf("%d", time.Now().UnixNano())
	res, err := ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderSingleRequest{
		Account: tc.Account,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    tc.Instrument,
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

	checkBalances(t, ctx.as, ctx.executor, tc.Account)

	fmt.Println("CHECK 1 GOOD")

	// Market sell 1 contract
	_, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderSingleRequest{
		Account: tc.Account,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    tc.Instrument,
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

	checkBalances(t, ctx.as, ctx.executor, tc.Account)

	fmt.Println("CLOSING")
	// Close position
	_, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderSingleRequest{
		Account: tc.Account,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    tc.Instrument,
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
	checkBalances(t, ctx.as, ctx.executor, tc.Account)

	// Close position
	_, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderSingleRequest{
		Account: tc.Account,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    tc.Instrument,
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
	checkBalances(t, ctx.as, ctx.executor, tc.Account)

	// Close position
	_, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderSingleRequest{
		Account: tc.Account,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    tc.Instrument,
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
	checkBalances(t, ctx.as, ctx.executor, tc.Account)

	// Close position
	_, err = ctx.as.Root.RequestFuture(ctx.executor, &messages.NewOrderSingleRequest{
		Account: tc.Account,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    tc.Instrument,
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
	checkBalances(t, ctx.as, ctx.executor, tc.Account)
}
