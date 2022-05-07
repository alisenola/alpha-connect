package memory

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"log"
	"os"
	"reflect"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"
)

func TestMemoryLeak(t *testing.T) {

	exch := constants.COINBASEPRO
	as, ex, cancel := tests.StartExecutor(t, exch, nil)
	defer cancel()

	f, err := os.Create("profiles/mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close() // error handling omitted for example

	securityID := []uint64{
		11630614572540763252,
	}
	testedSecurities := make(map[uint64]*models.Security)

	res, err := as.Root.RequestFuture(ex, &messages.SecurityListRequest{}, 10*time.Second).Result()
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
	fmt.Println(sec)

	test := tests.MDTest{
		IgnoreSizeResidue: true,
		SecurityID:        11630614572540763252,
		Symbol:            "BTC_USDT",
		SecurityType:      enum.SecurityType_CRYPTO_SPOT,
		Exchange:          constants.GATE,
		BaseCurrency:      constants.BITCOIN,
		QuoteCurrency:     constants.TETHER,
		MinPriceIncrement: 0.01,
		RoundLot:          1e-04,
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.InstrumentStatus_Trading,
	}
	//TODO Check if MDChecker replaced OBChecker
	for i := 0; i < 10; i++ {
		obChecker := as.Root.Spawn(actor.PropsFromProducer(tests.NewMDCheckerProducer(sec, test)))
		time.Sleep(5 * time.Second)
		err = as.Root.PoisonFuture(obChecker).Wait()
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(25 * time.Millisecond)
	}
	err = as.Root.PoisonFuture(ex).Wait()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	runtime.GC() // get up-to-date statistics

	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

}
