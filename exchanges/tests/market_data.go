package tests

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
)

type MDTest struct {
	IgnoreSizeResidue  bool
	IgnorePriceResidue bool
	SecurityID         uint64
	Symbol             string
	SecurityType       string
	Exchange           xchangerModels.Exchange
	BaseCurrency       xchangerModels.Asset
	QuoteCurrency      xchangerModels.Asset
	MinPriceIncrement  float64
	RoundLot           float64
	HasMaturityDate    bool
	IsInverse          bool
	Status             models.InstrumentStatus
}

func MarketData(t *testing.T, test MDTest) {
	//t.Parallel()
	as, executor, clean := StartExecutor(t, &test.Exchange, nil)
	defer clean()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		test.SecurityID,
	}
	testedSecurities := make(map[uint64]*models.Security)

	res, err := as.Root.RequestFuture(executor, &messages.SecurityListRequest{}, 10*time.Second).Result()
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

	// Test
	sec, ok := testedSecurities[test.SecurityID]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != test.Symbol {
		t.Fatalf("was expecting symbol %s, got %s", test.Symbol, sec.Symbol)
	}
	if sec.SecurityType != test.SecurityType {
		t.Fatalf("was expecting %s type, got %s", test.SecurityType, sec.SecurityType)
	}
	if sec.Exchange.Name != test.Exchange.Name {
		t.Fatalf("was expecting %s exchange, got %s", test.Exchange.Name, sec.Exchange.Name)
	}
	if sec.Underlying.ID != test.BaseCurrency.ID {
		t.Fatalf("was expecting %s base, got %s", test.BaseCurrency.Symbol, sec.Underlying.Symbol)
	}
	if sec.QuoteCurrency.ID != test.QuoteCurrency.ID {
		t.Fatalf("was expecting %s quote, got %s", test.QuoteCurrency.Symbol, sec.QuoteCurrency.Symbol)
	}
	if sec.IsInverse != test.IsInverse {
		t.Fatalf("was expecting different inverse")
	}
	if sec.Status != test.Status {
		t.Fatal("was expecting enabled security")
	}
	if sec.MinPriceIncrement != nil && math.Abs(sec.MinPriceIncrement.Value-test.MinPriceIncrement) > 0.000001 {
		t.Fatalf("was expecting %g min price increment, got %g", test.MinPriceIncrement, sec.MinPriceIncrement.Value)
	}
	if sec.RoundLot != nil && math.Abs(sec.RoundLot.Value-test.RoundLot) > 0.0000000001 {
		t.Fatalf("was expecting %g round lot increment, got %g", test.RoundLot, sec.RoundLot.Value)
	}
	if (sec.MaturityDate != nil) != test.HasMaturityDate {
		t.Fatalf("was expecting different maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec, test)))
	time.Sleep(80 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}
