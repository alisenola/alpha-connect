package tests

import (
	"gitlab.com/alphaticks/alpha-connect/config"
	"gitlab.com/alphaticks/alpha-connect/tests"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
)

type MDTest struct {
	IgnoreSizeResidue  bool
	IgnorePriceResidue bool
	Symbol             string
	SecurityType       string
	Exchange           *xchangerModels.Exchange
	MinPriceIncrement  float64
	RoundLot           float64
	HasMaturityDate    bool
	IsInverse          bool
	Status             models.InstrumentStatus
}

func CheckSecurityDefinition(t *testing.T, sec *models.Security, test MDTest) {
	// Test
	if sec == nil {
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
}

func MarketData(t *testing.T, test MDTest) {
	t.Parallel()

	C := &config.Config{
		Exchanges: []string{test.Exchange.Name},
	}
	as, executor, clean := tests.StartExecutor(t, C)
	defer clean()

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
	var sec *models.Security
	for _, s := range securityList.Securities {
		if s.Symbol == test.Symbol {
			sec = s
		}
	}

	if sec == nil {
		t.Fatalf("security not found")
	}

	CheckSecurityDefinition(t, sec, test)
	obChecker := as.Root.Spawn(actor.PropsFromProducer(NewMDCheckerProducer(sec, test)))
	defer as.Root.PoisonFuture(obChecker)

	time.Sleep(100 * time.Second)
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
