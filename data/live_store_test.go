package data

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"gitlab.com/tachikoma.ai/tickobjects"
	"gitlab.com/tachikoma.ai/tickstore-go-client/query"
	"testing"
	"time"
)

func StartExecutor(exchange *xchangerModels.Exchange) (*actor.ActorSystem, *actor.PID, func()) {
	exch := []*xchangerModels.Exchange{
		exchange,
		&constants.COINBASEPRO,
	}
	as := actor.NewActorSystem()
	executor, _ := as.Root.SpawnNamed(actor.PropsFromProducer(exchanges.NewExecutorProducer(exch, nil, false, xchangerUtils.DefaultDialerPool)), "executor")
	return as, executor, func() { _ = as.Root.PoisonFuture(executor).Wait() }
}

func TestLiveQuery(t *testing.T) {
	as, _, clean := StartExecutor(&constants.BINANCE)
	defer clean()
	lt, err := NewLiveStore(as)
	if err != nil {
		t.Fatal(err)
	}
	qs := query.NewQuerySettings(query.WithSelector(`SELECT OBPrice(AggregateOB(orderbook), "0.2") WHERE exchange="^binance|coinbasepro$" base="^BTC$" quote="^USDT$" GROUPBY base`))
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
