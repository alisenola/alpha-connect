package tests

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"math"
	"reflect"
	"sort"
	"testing"
	"time"
)

func checkBalances(t *testing.T, as *actor.ActorSystem, executor *actor.PID, account *models.Account) {
	// Now check balance
	exchangeExecutor := as.NewLocalPID(fmt.Sprintf("executor/exchanges/%s_executor", account.Exchange.Name))
	res, err := as.Root.RequestFuture(executor, &messages.BalancesRequest{
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

	execBalances := balanceResponse.Balances

	res, err = as.Root.RequestFuture(exchangeExecutor, &messages.BalancesRequest{
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

	accntBalances := balanceResponse.Balances

	sort.Slice(accntBalances, func(i, j int) bool {
		return accntBalances[i].Asset.ID < accntBalances[j].Asset.ID
	})
	sort.Slice(execBalances, func(i, j int) bool {
		return execBalances[i].Asset.ID < execBalances[j].Asset.ID
	})
	for i, b1 := range execBalances {
		b2 := accntBalances[i]
		// TODO use exchange balance precision
		rawB1 := int(math.Round(b1.Quantity * 1e8))
		rawB2 := int(math.Round(b2.Quantity * 1e8))
		fmt.Println(rawB1, rawB2)
		if rawB1 != rawB2 {
			t.Fatalf("different balance: %f %f", b1.Quantity, b2.Quantity)
		}
	}
}

func checkPositions(t *testing.T, as *actor.ActorSystem, executor *actor.PID, account *models.Account, instrument *models.Instrument) {
	// Request the same from binance directly
	exchangeExecutor := as.NewLocalPID(fmt.Sprintf("executor/exchanges/%s_executor", account.Exchange.Name))

	res, err := as.Root.RequestFuture(executor, &messages.PositionsRequest{
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

	res, err = as.Root.RequestFuture(exchangeExecutor, &messages.PositionsRequest{
		RequestID:  0,
		Account:    account,
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
		fmt.Println(p1.Cost, p2.Cost)
		if math.Abs(p1.Cost-p2.Cost) > 0.000001 {
			t.Fatalf("different cost %f:%f", p1.Cost, p2.Cost)
		}
		if math.Abs(p1.Quantity-p2.Quantity) > 0.00001 {
			t.Fatalf("different quantity %f:%f", p1.Quantity, p2.Quantity)
		}
	}
}

func checkOrders(as *actor.ActorSystem, executor *actor.PID, account *models.Account, sec *models.Security, filter *messages.OrderFilter) error {
	// Request the same from binance directly
	exchangeExecutor := as.NewLocalPID(fmt.Sprintf("executor/exchanges/%s_executor", account.Exchange.Name))

	res, err := as.Root.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID: 0,
		Account:   account,
		Filter:    filter,
	}, 10*time.Second).Result()

	if err != nil {
		return err
	}
	response, ok := res.(*messages.OrderList)
	if !ok {
		return fmt.Errorf("was expecting *messages.OrderList, got %s", reflect.TypeOf(res).String())
	}
	if !response.Success {
		return fmt.Errorf("was expecting successful request: %s", response.RejectionReason.String())
	}

	orders1 := response.Orders
	fmt.Printf("got %d orders from main executor i.e. account \n", len(orders1))

	res, err = as.Root.RequestFuture(exchangeExecutor, &messages.OrderStatusRequest{
		RequestID: 0,
		Account:   account,
		Filter:    filter,
	}, 10*time.Second).Result()

	if err != nil {
		return err
	}
	response, ok = res.(*messages.OrderList)
	if !ok {
		return fmt.Errorf("was expecting *messages.OrderList, got %s", reflect.TypeOf(res).String())
	}
	if !response.Success {
		return fmt.Errorf("was expecting sucessful request: %s", response.RejectionReason.String())
	}
	orders2 := response.Orders
	fmt.Printf("got %d orders from exchange executor i.e. api \n", len(orders2))
	if len(orders1) != len(orders2) {
		return fmt.Errorf("differents lengths, got %d vs %d", len(orders1), len(orders2))
	}
	sort.Slice(orders1, func(i, j int) bool {
		return orders1[i].ClientOrderID < orders1[j].ClientOrderID
	})
	sort.Slice(orders2, func(i, j int) bool {
		return orders2[i].ClientOrderID < orders2[j].ClientOrderID
	})
	for i := range orders1 {
		if orders1[i].ClientOrderID != orders2[i].ClientOrderID {
			return fmt.Errorf("different ClientOrderIds %s vs %s", orders1[i].ClientOrderID, orders2[i].ClientOrderID)
		}
		p1 := int(orders1[i].Price.Value / sec.MinPriceIncrement.Value)
		p2 := int(orders2[i].Price.Value / sec.MinPriceIncrement.Value)
		if p1 != p2 {
			fmt.Println(p1, p2)
			return fmt.Errorf("different prices %f vs %f", orders1[i].Price.Value, orders2[i].Price.Value)
		}
		q1 := int(orders1[i].CumQuantity / sec.RoundLot.Value)
		q2 := int(orders2[i].CumQuantity / sec.RoundLot.Value)
		if q1 != q2 {
			return fmt.Errorf("different quantities %f vs %f", orders1[i].CumQuantity, orders2[i].CumQuantity)
		}
	}
	fmt.Println(orders1)
	fmt.Println(orders2)
	return nil
}
