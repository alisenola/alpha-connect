package tests

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"reflect"
	"testing"
	"time"

	"gitlab.com/alphaticks/xchanger/constants"

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
	if sec.Underlying.ID != test.BaseCurrency.ID {
		t.Fatalf("was expecting %d:%s base, got %d:%s",
			test.BaseCurrency.ID,
			test.BaseCurrency.Symbol,
			sec.Underlying.ID,
			sec.Underlying.Symbol)
	}
	if sec.QuoteCurrency.ID != test.QuoteCurrency.ID {
		t.Fatalf("was expecting %d:%s quote, got %d:%s",
			test.QuoteCurrency.ID,
			test.QuoteCurrency.Symbol,
			sec.QuoteCurrency.ID,
			sec.QuoteCurrency.Symbol)
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
	b, err := ioutil.ReadFile("../../assets.json")
	if err != nil {
		t.Fatalf("failed to load assets: %v", err)
	}
	assets := make(map[uint32]xchangerModels.Asset)
	err = json.Unmarshal(b, &assets)
	if err != nil {
		t.Fatalf("error parsing assets: %v", err)
	}
	if err := constants.LoadAssets(assets); err != nil {
		t.Fatalf("error loading assets: %v", err)
	}
	as, executor, clean := StartExecutor(t, &test.Exchange, nil)
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

	if sec.SecurityID != test.SecurityID {
		t.Fatalf("different security ID %d: %d", sec.SecurityID, test.SecurityID)
	}

	CheckSecurityDefinition(t, sec, test)
	obChecker := as.Root.Spawn(actor.PropsFromProducer(NewMDCheckerProducer(sec, test)))
	defer as.Root.PoisonFuture(obChecker)

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
