package account

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/account"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/bitmex"
	"gitlab.com/alphaticks/xchanger/exchanges/fbinance"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"math"
	"reflect"
	"testing"
	"time"
)

var assets = map[uint32]xchangerModels.Asset{
	constants.DOLLAR.ID:           constants.DOLLAR,
	constants.EURO.ID:             constants.EURO,
	constants.POUND.ID:            constants.POUND,
	constants.CANADIAN_DOLLAR.ID:  constants.CANADIAN_DOLLAR,
	constants.JAPENESE_YEN.ID:     constants.JAPENESE_YEN,
	constants.BITCOIN.ID:          constants.BITCOIN,
	constants.LITECOIN.ID:         constants.LITECOIN,
	constants.ETHEREUM.ID:         constants.ETHEREUM,
	constants.RIPPLE.ID:           constants.RIPPLE,
	constants.TETHER.ID:           constants.TETHER,
	constants.SOUTH_KOREAN_WON.ID: constants.SOUTH_KOREAN_WON,
	constants.USDC.ID:             constants.USDC,
	constants.DASH.ID:             constants.DASH,
}

func start(t *testing.T, accnt *models.Account) (*actor.ActorSystem, *actor.PID) {
	fmt.Println("LOAD ASSETS")
	if err := constants.LoadAssets(assets); err != nil {
		panic(err)
	}
	exch := []*xchangerModels.Exchange{
		//&constants.FBINANCE,
		//&constants.BITMEX,
		&constants.BINANCE,
	}
	accnts := []*account.Account{}
	aa, err := exchanges.NewAccount(accnt)
	if err != nil {
		t.Fatal(err)
	}
	accnts = append(accnts, aa)

	as := actor.NewActorSystem()

	cfg := exchanges.ExecutorConfig{
		Exchanges:  exch,
		Accounts:   accnts,
		DialerPool: xchangerUtils.DefaultDialerPool,
		Strict:     true,
	}
	executor, _ := as.Root.SpawnNamed(actor.PropsFromProducer(exchanges.NewExecutorProducer(&cfg)), "executor")
	bitmex.EnableTestNet()
	fbinance.EnableTestNet()
	return as, executor
}

func checkBalances(t *testing.T, as *actor.ActorSystem, executor *actor.PID, account *models.Account) {
	// Now check balance
	binanceExecutor := as.NewLocalPID("executor/binance_executor")
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

	bal1 := make(map[string]float64)
	for _, b := range balanceResponse.Balances {
		bal1[b.Asset.Name] = b.Quantity
	}
	fmt.Println("ACCOUNT BALANCE", bal1)

	res, err = as.Root.RequestFuture(binanceExecutor, &messages.BalancesRequest{
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

	bal2 := make(map[string]float64)
	for _, b := range balanceResponse.Balances {
		bal2[b.Asset.Name] = b.Quantity
	}
	fmt.Println("EXECUTOR BALANCE", bal2)

	for k, b1 := range bal1 {
		if b2, ok := bal2[k]; ok {
			if math.Abs(b1-b2) > 0.00001 {
				t.Fatalf("different balance for %s %f:%f", k, b1, b2)
			}
		} else {
			t.Fatalf("different balance for %s %f:%f", k, b1, 0.)
		}
	}
}

func checkPositions(t *testing.T, as *actor.ActorSystem, executor *actor.PID, account *models.Account, instrument *models.Instrument) {
	// Request the same from binance directly
	binanceExecutor := as.NewLocalPID("executor/binance_executor")

	res, err := as.Root.RequestFuture(binanceExecutor, &messages.PositionsRequest{
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

	res, err = as.Root.RequestFuture(executor, &messages.PositionsRequest{
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
