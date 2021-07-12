package account

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
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

var instrument = &models.Instrument{
	SecurityID: &types.UInt64Value{Value: 5485975358912730733},
	Exchange:   &constants.FBINANCE,
	Symbol:     &types.StringValue{Value: "BTCUSDT"},
}

var FBinanceTestnetAccount = &models.Account{
	AccountID: "299211",
	Exchange:  &constants.FBINANCE,
	Credentials: &xchangerModels.APICredentials{
		APIKey:    "74f122652da74f6e1bcc34b8c23fc91e0239b502e68440632ae9a3cb7cefa18e",
		APISecret: "c3e0d76ee014b597b93616478dc789e6bb6616ad59ddbe384d2554ace4a60f86",
	},
}

/*
var FBinanceAccount = &models.Account{
	AccountID: "299211",
	Exchange:  &constants.FBINANCE,
	Credentials: &xchangerModels.APICredentials{
		APIKey:    "MpYkeK3pGP80gGiIrqWtLNwjJmyK2DTREYzNx8Cyc3AWTkl2T0iWnQEtdCIlvAoE",
		APISecret: "CJcJZEkktzhGzEdQhclfHcfJz5k01OY6n42MeF9B3oQWGqba3RrXEnG4bZktXQNu",
	},
}
*/

func TestFBinanceAccountListener_OnOrderCancelRequest(t *testing.T) {
	orderID := fmt.Sprintf("%d", time.Now().UnixNano())
	// Test with one order
	res, err := As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   FBinanceTestnetAccount,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    instrument,
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
	// cancel directly from the executor
	fbinanceExecutor := actor.NewPID(As.Address(), "executor/"+constants.FBINANCE.Name+"_executor")

	res, err = As.Root.RequestFuture(fbinanceExecutor, &messages.OrderCancelRequest{
		RequestID:     0,
		Account:       FBinanceTestnetAccount,
		Instrument:    instrument,
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
		t.Fatalf("was expecting successful request: %s", mcResponse.RejectionReason.String())
	}

	res, err = As.Root.RequestFuture(fbinanceExecutor, &messages.OrderCancelRequest{
		RequestID:     0,
		Account:       FBinanceTestnetAccount,
		Instrument:    instrument,
		ClientOrderID: &types.StringValue{Value: orderID},
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	mcResponse, ok = res.(*messages.OrderCancelResponse)
	if !ok {
		t.Fatalf("was expecting *messages.OrderCancelResponse, got %s", reflect.TypeOf(res).String())
	}
	if mcResponse.Success {
		t.Fatalf("was expecting unsuccessful request")
	}
	if mcResponse.RejectionReason != messages.UnknownOrder {
		t.Fatalf("was expecting unknown order")
	}

	// Now cancel from the account

	res, err = As.Root.RequestFuture(executor, &messages.OrderCancelRequest{
		RequestID:  0,
		Account:    FBinanceTestnetAccount,
		Instrument: instrument,
		OrderID:    &types.StringValue{Value: response.OrderID},
	}, 10*time.Second).Result()

	if err != nil {
		t.Fatal(err)
	}
	mcResponse, ok = res.(*messages.OrderCancelResponse)
	if !ok {
		t.Fatalf("was expecting *messages.OrderCancelResponse, got %s", reflect.TypeOf(res).String())
	}
	if mcResponse.Success {
		t.Fatalf("was expecting unsuccessful request")
	}
	if mcResponse.RejectionReason != messages.UnknownOrder {
		t.Fatalf("was expecting unknown order")
	}
}

func TestFBinanceAccountListener_OnOrderStatusRequest(t *testing.T) {
	// Test with no account
	res, err := As.Root.RequestFuture(executor, &messages.OrderStatusRequest{
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
	res, err = As.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   FBinanceTestnetAccount,
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
	res, err = As.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   FBinanceTestnetAccount,
		Filter: &messages.OrderFilter{
			OrderID:       nil,
			ClientOrderID: nil,
			Instrument:    instrument,
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
	res, err = As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   FBinanceTestnetAccount,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    instrument,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.GoodTillCancel,
			Quantity:      1.,
			Price:         &types.DoubleValue{Value: 10000.},
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
		Account:   FBinanceTestnetAccount,
		Filter: &messages.OrderFilter{
			OrderID:       nil,
			ClientOrderID: nil,
			Instrument:    instrument,
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
	if order.TimeInForce != models.GoodTillCancel {
		t.Fatalf("was expecting GoodTillCancel time in force")
	}
	if order.Side != models.Buy {
		t.Fatalf("was expecting buy side order")
	}

	res, err = As.Root.RequestFuture(executor, &messages.OrderCancelRequest{
		RequestID:  0,
		Account:    FBinanceTestnetAccount,
		Instrument: instrument,
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
		res, err = As.Root.RequestFuture(executor, &messages.OrderMassCancelRequest{
			RequestID: 0,
			Account:   FBinanceTestnetAccount,
			Filter: &messages.OrderFilter{
				Instrument: instrument,
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
	res, err = As.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   FBinanceTestnetAccount,
		Filter: &messages.OrderFilter{
			OrderID:    &types.StringValue{Value: order.OrderID},
			Instrument: instrument,
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

func TestFBinanceAccountListener_OnNewOrderSingleRequest(t *testing.T) {
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
		Account:   FBinanceTestnetAccount,
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
		Account:   FBinanceTestnetAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument,
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
	res, err = As.Root.RequestFuture(executor, &messages.OrderMassCancelRequest{
		RequestID: 0,
		Account:   FBinanceTestnetAccount,
		Filter: &messages.OrderFilter{
			Instrument: instrument,
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

func TestListen(t *testing.T) {
	time.Sleep(1 * time.Minute)
}

func TestFBinanceAccountListener_OnBalancesRequest(t *testing.T) {
	res, err := As.Root.RequestFuture(executor, &messages.BalancesRequest{
		RequestID: 0,
		Account:   FBinanceTestnetAccount,
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
	fmt.Println(balanceResponse.Balances)
}

func checkFBinanceBalances(t *testing.T, account *models.Account) {
	// Now check balance
	fbinanceExecutor := As.NewLocalPID("executor/fbinance_executor")
	res, err := As.Root.RequestFuture(executor, &messages.BalancesRequest{
		RequestID: 0,
		Account:   account,
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

	res, err = As.Root.RequestFuture(fbinanceExecutor, &messages.BalancesRequest{
		RequestID: 0,
		Account:   account,
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

	fmt.Println(bal1.Quantity, bal2.Quantity)
	if math.Abs(bal1.Quantity-bal2.Quantity) > 0.00001 {
		t.Fatalf("different balance %f:%f", bal1.Quantity, bal2.Quantity)
	}
}

func checkFBinancePositions(t *testing.T, account *models.Account, instrument *models.Instrument) {
	// Request the same from fbinance directly
	fbinanceExecutor := As.NewLocalPID("executor/fbinance_executor")

	res, err := As.Root.RequestFuture(fbinanceExecutor, &messages.PositionsRequest{
		RequestID:  0,
		Account:    account,
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
		Account:    FBinanceTestnetAccount,
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

func TestFBinanceAccountListener_OnGetPositionsLimit(t *testing.T) {

	// Market buy 1 contract
	orderID := fmt.Sprintf("%d", time.Now().UnixNano())
	res, err := As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: FBinanceTestnetAccount,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    instrument,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			Price:         &types.DoubleValue{Value: 39000},
			Quantity:      0.002,
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

	checkFBinancePositions(t, FBinanceTestnetAccount, instrument)
	checkFBinanceBalances(t, FBinanceTestnetAccount)

	// Market sell 1 contract
	res, err = As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: FBinanceTestnetAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument,
			OrderType:     models.Limit,
			OrderSide:     models.Sell,
			TimeInForce:   models.Session,
			Price:         &types.DoubleValue{Value: 30000},
			Quantity:      0.001,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	checkFBinancePositions(t, FBinanceTestnetAccount, instrument)
	checkFBinanceBalances(t, FBinanceTestnetAccount)

	fmt.Println("CLOSING")
	// Close position
	res, err = As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: FBinanceTestnetAccount,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    instrument,
			OrderType:     models.Limit,
			OrderSide:     models.Sell,
			TimeInForce:   models.Session,
			Price:         &types.DoubleValue{Value: 30000},
			Quantity:      0.001,
		},
	}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	checkFBinancePositions(t, FBinanceTestnetAccount, instrument)
	checkFBinanceBalances(t, FBinanceTestnetAccount)
}

func TestFBinanceAccountListener_OnOrderReplaceRequest(t *testing.T) {
	orderID := uuid.NewV1().String()
	// Post one order
	res, err := As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		RequestID: 0,
		Account:   FBinanceTestnetAccount,
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
	res, err = As.Root.RequestFuture(executor, &messages.OrderReplaceRequest{
		RequestID:  0,
		Account:    FBinanceTestnetAccount,
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
	res, err = As.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Subscribe: false,
		Account:   FBinanceTestnetAccount,
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

func TestFBinanceAccountListener_OnGetPositionsMarket(t *testing.T) {

	// Market buy 1 contract
	orderID := fmt.Sprintf("%d", time.Now().UnixNano())
	res, err := As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: FBinanceTestnetAccount,
		Order: &messages.NewOrder{
			ClientOrderID: orderID,
			Instrument:    instrument,
			OrderType:     models.Market,
			OrderSide:     models.Buy,
			Quantity:      0.002,
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

	checkFBinancePositions(t, FBinanceTestnetAccount, instrument)
	checkFBinanceBalances(t, FBinanceTestnetAccount)
	fmt.Println("CHECK 1 GOOD")

	// Market sell 1 contract
	res, err = As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: FBinanceTestnetAccount,
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

	checkFBinancePositions(t, FBinanceTestnetAccount, instrument)
	checkFBinanceBalances(t, FBinanceTestnetAccount)

	fmt.Println("CLOSING")
	// Close position
	res, err = As.Root.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: FBinanceTestnetAccount,
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

	checkFBinancePositions(t, FBinanceTestnetAccount, instrument)
	checkFBinanceBalances(t, FBinanceTestnetAccount)
}
