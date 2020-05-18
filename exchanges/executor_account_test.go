package exchanges

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/bitmex"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"reflect"
	"testing"
	"time"
)

type AccountChecker struct {
	account  *models.Account
	security *models.Security
	seqNum   uint64
	err      error
}

func NewAccountCheckerProducer(account *models.Account, security *models.Security) actor.Producer {
	return func() actor.Actor {
		return NewAccountChecker(account, security)
	}
}

func NewAccountChecker(account *models.Account, security *models.Security) actor.Actor {
	return &AccountChecker{
		account:  account,
		security: security,
		seqNum:   0,
		err:      nil,
	}
}

func (state *AccountChecker) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.err = err
		}

	case *messages.ExecutionReport:
		fmt.Println("GOT EXEC REPORT", time.Now())
		fmt.Println(context.Message())

	case *GetStat:
		context.Respond(&GetStat{
			Error:     state.err,
			Trades:    0,
			OBUpdates: 0,
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
		Instrument: nil,
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
	fmt.Println("PLACING LIMIT BUY")
	context.Send(executor, &messages.NewOrderSingle{
		ClientOrderID: "ord1",
		Instrument: &models.Instrument{
			Exchange:   state.security.Exchange,
			SecurityID: &types.UInt64Value{Value: state.security.SecurityID},
			Symbol:     &types.StringValue{Value: state.security.Symbol},
		},
		Account:     state.account,
		OrderType:   models.Limit,
		OrderSide:   models.Buy,
		TimeInForce: models.Session,
		Quantity:    1.,
		Price:       &types.DoubleValue{Value: 100.},
	})

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

	account := &models.Account{
		AccountID: "299210",
		Exchange:  &constants.BITMEX,
		Credentials: &xchangerModels.APICredentials{
			APIKey:    "k5k6Mmaq3xe88Ph3fgIk9Vrt",
			APISecret: "0laIjZaKOMkJPtKy2ldJ18m4Dxjp66Vdim0k1-q4TXASZFZo",
		},
	}

	exchanges := []*xchangerModels.Exchange{&constants.BITMEX}
	accounts := []*models.Account{account}
	securityID := []uint64{
		5391998915988476130, //XBTUSD
	}
	testedSecurities := make(map[uint64]*models.Security)
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges, accounts)), "executor")

	res, err := actor.EmptyRootContext.RequestFuture(executor, &messages.SecurityListRequest{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	securityList, ok := res.(*messages.SecurityList)
	if !ok {
		t.Fatalf("was expecting *messages.SecurityList, got %s", reflect.TypeOf(res).String())
	}
	if securityList.Error != "" {
		t.Fatal(securityList.Error)
	}

	for _, s := range securityList.Securities {
		tested := false
		for _, secID := range securityID {
			if secID == s.SecurityID {
				tested = true
				break
			}
		}
		if tested {
			testedSecurities[s.SecurityID] = s
		}
	}

	// Test
	sec, ok := testedSecurities[5391998915988476130]
	if !ok {
		t.Fatalf("XBTUSD not found")
	}

	accntChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewAccountCheckerProducer(account, sec)))
	time.Sleep(5 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(accntChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
	fmt.Println(res)
}
