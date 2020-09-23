package exchanges

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/types"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/alphaticks/alphac/account"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/bitmex"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"reflect"
	"testing"
	"time"
)

type GetAccntCheckerStats struct {
	Error   error
	Reports []*messages.ExecutionReport
}

type AccountChecker struct {
	account    *models.Account
	instrument *models.Instrument
	seqNum     uint64
	err        error
	reports    []*messages.ExecutionReport
}

func NewAccountCheckerProducer(account *models.Account, instrument *models.Instrument) actor.Producer {
	return func() actor.Actor {
		return NewAccountChecker(account, instrument)
	}
}

func NewAccountChecker(account *models.Account, instrument *models.Instrument) actor.Actor {
	return &AccountChecker{
		account:    account,
		instrument: instrument,
		seqNum:     0,
		err:        nil,
		reports:    nil,
	}
}

func (state *AccountChecker) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.err = err
		}

	case *messages.ExecutionReport:
		if err := state.OnExecutionReport(context); err != nil {
			state.err = err
		}

	case *GetAccntCheckerStats:
		context.Respond(&GetAccntCheckerStats{
			Error:   state.err,
			Reports: state.reports,
		})
	}
}

func (state *AccountChecker) Initialize(context actor.Context) error {
	executor := actor.NewLocalPID("executor")
	// Get orders & subscribe
	res, err := context.RequestFuture(executor, &messages.OrderStatusRequest{
		RequestID:  0,
		Subscribe:  true,
		Subscriber: context.Self(),
		Account:    state.account,
	}, 10*time.Second).Result()
	if err != nil {
		return err
	}
	_, ok := res.(*messages.OrderList)
	if !ok {
		return fmt.Errorf("was expecting *messages.OrderList, got %s", reflect.TypeOf(res))
	}
	// Place limit buy for 1 contract
	res, err = context.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: state.account,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    state.instrument,
			OrderType:     models.Limit,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      1.,
			Price:         &types.DoubleValue{Value: 100.},
		},
	}, 10*time.Second).Result()
	if err != nil {
		return err
	}
	newOrderResponse := res.(*messages.NewOrderSingleResponse)
	if !newOrderResponse.Success {
		return fmt.Errorf("error creating new order: %s", newOrderResponse.RejectionReason.String())
	}

	res, err = context.RequestFuture(executor, &messages.OrderCancelRequest{
		OrderID:    &types.StringValue{Value: newOrderResponse.OrderID},
		Instrument: state.instrument,
		Account:    state.account,
	}, 5*time.Second).Result()
	if err != nil {
		return err
	}
	cancelResponse := res.(*messages.OrderCancelResponse)
	if !cancelResponse.Success {
		return fmt.Errorf("cancel unsucessful %s", cancelResponse.RejectionReason.String())
	}

	// Market buy 1 contract
	res, err = context.RequestFuture(executor, &messages.NewOrderSingleRequest{
		Account: state.account,
		Order: &messages.NewOrder{
			ClientOrderID: uuid.NewV1().String(),
			Instrument:    state.instrument,
			OrderType:     models.Market,
			OrderSide:     models.Buy,
			TimeInForce:   models.Session,
			Quantity:      1.,
		},
	}, 10*time.Second).Result()
	if err != nil {
		return err
	}
	newOrderResponse = res.(*messages.NewOrderSingleResponse)
	if !newOrderResponse.Success {
		return fmt.Errorf("error creating new order: %s", newOrderResponse.RejectionReason.String())
	}

	return nil
}

func (state *AccountChecker) OnExecutionReport(context actor.Context) error {
	report := context.Message().(*messages.ExecutionReport)
	state.reports = append(state.reports, report)
	return nil
}

func cleanAccount() {
	_ = actor.EmptyRootContext.PoisonFuture(executor).Wait()
	_ = actor.EmptyRootContext.PoisonFuture(accntChecker).Wait()
}

var accntChecker *actor.PID

func TestBitmexAccount(t *testing.T) {
	defer cleanAccount()

	bitmex.EnableTestNet()

	accnt := &models.Account{
		AccountID: "299210",
		Exchange:  &constants.BITMEX,
		Credentials: &xchangerModels.APICredentials{
			APIKey:    "k5k6Mmaq3xe88Ph3fgIk9Vrt",
			APISecret: "0laIjZaKOMkJPtKy2ldJ18m4Dxjp66Vdim0k1-q4TXASZFZo",
		},
	}

	instrument := &models.Instrument{
		SecurityID: &types.UInt64Value{Value: 5391998915988476130},
		Exchange:   &constants.BITMEX,
		Symbol:     &types.StringValue{Value: "XBTUSD"},
	}

	exchanges := []*xchangerModels.Exchange{&constants.BITMEX}
	acnt, err := NewAccount(accnt)
	if err != nil {
		t.Fatal(err)
	}
	accounts := []*account.Account{acnt}
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges, accounts, false)), "executor")

	accntChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewAccountCheckerProducer(accnt, instrument)))
	time.Sleep(5 * time.Second)
	res, err := actor.EmptyRootContext.RequestFuture(accntChecker, &GetAccntCheckerStats{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetAccntCheckerStats)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
	/*
		if len(stats.Reports) != 4 {
			t.Fatalf("was expecting 4 reports, got %d", len(stats.Reports))
		}

	*/
	if stats.Reports[0].ExecutionType != messages.PendingNew {
		t.Fatalf("was expecting PendingNew, got %s", stats.Reports[0].ExecutionType.String())
	}
	if stats.Reports[1].ExecutionType != messages.New {
		t.Fatalf("was expecting New, got %s", stats.Reports[1].ExecutionType.String())
	}
	if stats.Reports[2].ExecutionType != messages.PendingCancel {
		t.Fatalf("was expecting PendingCancel, got %s", stats.Reports[2].ExecutionType.String())
	}
	if stats.Reports[3].ExecutionType != messages.Canceled {
		t.Fatalf("was expecting Canceled, got %s", stats.Reports[3].ExecutionType.String())
	}
	if stats.Reports[4].ExecutionType != messages.PendingNew {
		t.Fatalf("was expecting PendingNew, got %s", stats.Reports[4].ExecutionType.String())
	}
	if stats.Reports[5].ExecutionType != messages.New {
		t.Fatalf("was expecting New, got %s", stats.Reports[5].ExecutionType.String())
	}
	if stats.Reports[6].ExecutionType != messages.Trade {
		t.Fatalf("was expecting Trade, got %s", stats.Reports[6].ExecutionType.String())
	}
	if stats.Reports[6].OrderStatus != models.Filled {
		t.Fatalf("was expecting Filled, got %s", stats.Reports[6].ExecutionType.String())
	}
	fmt.Println(stats.Reports[6].FillPrice, stats.Reports[6].FillQuantity)
	fmt.Println(stats.Reports)
}
