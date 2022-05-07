package data

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	types "gitlab.com/alphaticks/tickstore-types"
	"gitlab.com/alphaticks/tickstore-types/tickobjects"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"testing"
	"time"
)

func StartExecutor(exchange *xchangerModels.Exchange) (*actor.ActorSystem, *actor.PID, func()) {
	exch := []*xchangerModels.Exchange{
		exchange,
		constants.COINBASEPRO,
	}
	as := actor.NewActorSystem()
	cfg := exchanges.ExecutorConfig{
		Exchanges:  exch,
		Strict:     false,
		DialerPool: xchangerUtils.DefaultDialerPool,
	}
	executor, _ := as.Root.SpawnNamed(actor.PropsFromProducer(exchanges.NewExecutorProducer(&cfg)), "executor")
	return as, executor, func() { _ = as.Root.PoisonFuture(executor).Wait() }
}

func TestLiveQuery(t *testing.T) {
	as, executor, clean := StartExecutor(constants.BINANCE)
	defer clean()
	lt, err := NewLiveStore(as, executor)
	if err != nil {
		t.Fatal(err)
	}
	qs := types.NewQuerySettings(types.WithSelector(`SELECT OBPrice(AggregateOB(orderbook), "0.2") WHERE exchange="^binance|coinbasepro$" base="^BTC$" quote="^USDT$" GROUPBY base`))
	q, err := lt.NewQuery(qs)
	if err != nil {
		t.Fatal(err)
	}
	i := 0
	for i < 100 {
		for q.Next() {
			tick, to, gid := q.Read()
			fmt.Println(tick, to.(tickobjects.Float64Object).Float64(), gid)
			i += 1
		}
		if q.Err() != nil {
			t.Fatal(q.Err())
		}
	}

	if err := q.Close(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Second)
}
