package memory

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"log"
	"os"
	"reflect"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"
)

func TestMemoryLeak(t *testing.T) {

	exch := []*xchangerModels.Exchange{
		&constants.COINBASEPRO,
	}
	executor, _ := actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(exchanges.NewExecutorProducer(exch, nil, false, xchangerUtils.DefaultDialerPool)), "executor")

	f, err := os.Create("profiles/mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close() // error handling omitted for example

	securityID := []uint64{
		11630614572540763252,
	}
	testedSecurities := make(map[uint64]*models.Security)

	res, err := actor.EmptyRootContext.RequestFuture(executor, &messages.SecurityListRequest{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	securityList, ok := res.(*messages.SecurityList)
	if !ok {
		t.Fatalf("was expecting *messages.SecurityList, got %s", reflect.TypeOf(res).String())
	}
	if !securityList.Success {
		t.Fatal(securityList.RejectionReason.String())
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

	// Test BTCEUR
	sec, ok := testedSecurities[11630614572540763252]
	if !ok {
		t.Fatalf("BTCUSD not found")
	}

	for i := 0; i < 10; i++ {
		obChecker := actor.EmptyRootContext.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
		time.Sleep(5 * time.Second)
		err = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(25 * time.Millisecond)
	}
	err = actor.EmptyRootContext.PoisonFuture(executor).Wait()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	runtime.GC() // get up-to-date statistics

	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

}
