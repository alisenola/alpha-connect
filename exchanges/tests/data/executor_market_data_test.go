package data

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"math"
	"os"
	"reflect"
	"testing"
	"time"
)

var executor *actor.PID
var as *actor.ActorSystem

func TestMain(m *testing.M) {
	exch := []*xchangerModels.Exchange{
		&constants.CRYPTOFACILITIES,
	}
	as = actor.NewActorSystem()
	executor, _ = as.Root.SpawnNamed(actor.PropsFromProducer(exchanges.NewExecutorProducer(exch, nil, false, xchangerUtils.DefaultDialerPool)), "executor")
	code := m.Run()
	_ = as.Root.PoisonFuture(executor).Wait()
	os.Exit(code)
}

func TestBinance(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()

	securityID := []uint64{
		9281941173829172773, //BTCUSDT
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

	// Test BTCUSDT
	sec, ok := testedSecurities[9281941173829172773]
	if !ok {
		t.Fatalf("BTCUSDT not found")
	}
	if sec.Symbol != "BTCUSDT" {
		t.Fatalf("was expecting symbol BTCEUR, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.BINANCE.Name {
		t.Fatalf("was expecting binance exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.TETHER.ID {
		t.Fatalf("was expecting EUR quote, got %d", sec.QuoteCurrency.ID)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %f", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-0.000001) > 0.0000000001 {
		t.Fatalf("was expecting 0.000001 round lot increment, got %f", sec.RoundLot.Value)
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestBitfinex(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()

	securityID := []uint64{
		17873758715870285590, //BTCUSDT
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

	// Test BTCEUR
	sec, ok := testedSecurities[17873758715870285590]
	if !ok {
		t.Fatalf("BTCUSD not found")
	}
	if sec.Symbol != "btcusd" {
		t.Fatalf("was expecting symbol btcust, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRPERP type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.BITFINEX.Name {
		t.Fatalf("was expecting BITFINEX exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", sec.QuoteCurrency.ID)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	/*
		if math.Abs(sec.MinPriceIncrement.Value-0.1) > 0.000001 {
			t.Fatalf("was expecting 0.1 min price increment, got %f", sec.MinPriceIncrement.Value)
		}
	*/
	if math.Abs(sec.RoundLot.Value-1./100000000.) > 0.0000000001 {
		t.Fatalf("was expecting 0.000001 round lot increment, got %f", sec.RoundLot.Value)
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestBitmex(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		5391998915988476130, //XBTUSD
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
	sec, ok := testedSecurities[5391998915988476130]
	if !ok {
		t.Fatalf("XBTUSD not found")
	}
	if sec.Symbol != "XBTUSD" {
		t.Fatalf("was expecting symbol btcust, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_PERP {
		t.Fatalf("was expecting CRSPOT type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.BITMEX.Name {
		t.Fatalf("was expecting BITFINEX exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", sec.QuoteCurrency.ID)
	}
	if !sec.IsInverse {
		t.Fatalf("was expecting inverse, got noninverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.5) > 0.000001 {
		t.Fatalf("was expecting 0.5 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-1.) > 0.0000000001 {
		t.Fatalf("was expecting 1. round lot increment, got %g", sec.RoundLot.Value)
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Agg Trades: %d | Trades: %d | OBUpdates: %d", stats.AggTrades, stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestBitstamp(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		5279696656781449381,
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
	sec, ok := testedSecurities[5279696656781449381]
	if !ok {
		t.Fatalf("BTCUSD not found")
	}
	if sec.Symbol != "btcusd" {
		t.Fatalf("was expecting symbol BTC/USD, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.BITSTAMP.Name {
		t.Fatalf("was expecting BITFINEX exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", sec.QuoteCurrency.ID)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-0.00000001) > 0.0000000001 {
		t.Fatalf("was expecting 0.00000001 round lot increment, got %g", sec.RoundLot.Value)
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestBitz(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		243278890145991530,
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
	sec, ok := testedSecurities[243278890145991530]
	if !ok {
		t.Fatalf("BTCUSDT not found")
	}
	if sec.Symbol != "btc_usdt" {
		t.Fatalf("was expecting symbol btc_usdt, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.BITZ.Name {
		t.Fatalf("was expecting bitz exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.TETHER.ID {
		t.Fatalf("was expecting USDT quote, got %d", sec.QuoteCurrency.ID)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-0.0001) > 0.0000000001 {
		t.Fatalf("was expecting 0.0001 round lot increment, got %g", sec.RoundLot.Value)
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestCoinbasePro(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		11630614572540763252,
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
	sec, ok := testedSecurities[11630614572540763252]
	if !ok {
		t.Fatalf("BTCUSD not found")
	}
	if sec.Symbol != "BTC-USD" {
		t.Fatalf("was expecting symbol BTC-USD, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.COINBASEPRO.Name {
		t.Fatalf("was expecting bitz exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", sec.QuoteCurrency.ID)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-0.00000001) > 0.0000000001 {
		t.Fatalf("was expecting 0.0001 round lot increment, got %g", sec.RoundLot.Value)
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestCryptofacilities(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		1416768858288349990,
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
	sec, ok := testedSecurities[1416768858288349990]
	if !ok {
		t.Fatalf("XBTUSD not found")
	}
	if sec.Symbol != "pi_xbtusd" {
		t.Fatalf("was expecting symbol pi_xbtusd, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_PERP {
		t.Fatalf("was expecting CRPERP type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.CRYPTOFACILITIES.Name {
		t.Fatalf("was expecting bitz exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", sec.QuoteCurrency.ID)
	}
	if !sec.IsInverse {
		t.Fatalf("was expecting inverse, got non inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.5) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-1.) > 0.0000000001 {
		t.Fatalf("was expecting 0.0001 round lot increment, got %g", sec.RoundLot.Value)
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestFBinance(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		5485975358912730733,
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
	sec, ok := testedSecurities[5485975358912730733]
	if !ok {
		t.Fatalf("BTCUSDT not found")
	}
	if sec.Symbol != "BTCUSDT" {
		t.Fatalf("was expecting symbol BTCUSDT, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_PERP {
		t.Fatalf("was expecting CRPERP type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.FBINANCE.Name {
		t.Fatalf("was expecting bitz exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.TETHER.ID {
		t.Fatalf("was expecting USDT quote, got %d", sec.QuoteCurrency.ID)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-0.001) > 0.0000000001 {
		t.Fatalf("was expecting 0.001 round lot increment, got %g", sec.RoundLot.Value)
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestFTX(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		4425198260936995601,
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
	sec, ok := testedSecurities[4425198260936995601]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != "BTC-PERP" {
		t.Fatalf("was expecting symbol BTC-PERP, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_PERP {
		t.Fatalf("was expecting CRPERP type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.FTX.Name {
		t.Fatalf("was expecting ftx exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", sec.QuoteCurrency.ID)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.5) > 0.000001 {
		t.Fatalf("was expecting 0.5 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-0.0001) > 0.0000000001 {
		t.Fatalf("was expecting 0.0001 round lot increment, got %g", sec.RoundLot.Value)
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestHuobi(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		2195469462990134438,
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
	sec, ok := testedSecurities[2195469462990134438]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != "btcusdt" {
		t.Fatalf("was expecting symbol btcusdt, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.HUOBI.Name {
		t.Fatalf("was expecting huobi exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.TETHER.ID {
		t.Fatalf("was expecting USD quote, got %d", sec.QuoteCurrency.ID)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-1e-6) > 0.0000000001 {
		t.Fatalf("was expecting 1e-6 round lot increment, got %g", sec.RoundLot.Value)
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestGemini(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		17496373742670049989,
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
	sec, ok := testedSecurities[17496373742670049989]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != "btcusd" {
		t.Fatalf("was expecting symbol btcusd, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.GEMINI.Name {
		t.Fatalf("was expecting gemini exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", sec.QuoteCurrency.ID)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-1e-08) > 0.0000000001 {
		t.Fatalf("was expecting 1e-6 round lot increment, got %g", sec.RoundLot.Value)
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestHitbtc(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		12674447834540883135,
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
	sec, ok := testedSecurities[12674447834540883135]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != "BTCUSD" {
		t.Fatalf("was expecting symbol btcusd, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.HITBTC.Name {
		t.Fatalf("was expecting hitbtc exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.TETHER.ID {
		t.Fatalf("was expecting USDT quote, got %d", sec.QuoteCurrency.ID)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-1e-05) > 0.0000000001 {
		t.Fatalf("was expecting 1e-6 round lot increment, got %g", sec.RoundLot.Value)
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestKraken(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		10955098577666557860,
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
	sec, ok := testedSecurities[10955098577666557860]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != "XBT/USD" {
		t.Fatalf("was expecting symbol btcusd, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.KRAKEN.Name {
		t.Fatalf("was expecting kraken exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", sec.QuoteCurrency.ID)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.1) > 0.000001 {
		t.Fatalf("was expecting 0.1 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-1e-08) > 0.0000000001 {
		t.Fatalf("was expecting 1e-8 round lot increment, got %g", sec.RoundLot.Value)
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestOKCoin(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		16235745264492357730,
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
	sec, ok := testedSecurities[16235745264492357730]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != "BTC-USDT" {
		t.Fatalf("was expecting symbol BTC-USDT, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.OKCOIN.Name {
		t.Fatalf("was expecting kraken exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.TETHER.ID {
		t.Fatalf("was expecting USDT quote, got %d", sec.QuoteCurrency.ID)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.1) > 0.000001 {
		t.Fatalf("was expecting 0.1 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-0.0001) > 0.0000000001 {
		t.Fatalf("was expecting 0.0001 round lot increment, got %g", sec.RoundLot.Value)
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestOKex(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		945944519923594006,
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
	sec, ok := testedSecurities[945944519923594006]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != "BTC-USDT" {
		t.Fatalf("was expecting symbol BTC-USDT, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.OKEX.Name {
		t.Fatalf("was expecting OKEX exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.TETHER.ID {
		t.Fatalf("was expecting USDT quote, got %d", sec.QuoteCurrency.ID)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.1) > 0.000001 {
		t.Fatalf("was expecting 0.1 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-1e-08) > 0.0000000001 {
		t.Fatalf("was expecting 1e-08 round lot increment, got %g", sec.RoundLot.Value)
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestDeribit(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		2206542817128348325,
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
		fmt.Println(s)
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
	sec, ok := testedSecurities[2206542817128348325]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != "BTC-PERPETUAL" {
		t.Fatalf("was expecting symbol BTC-PERPETUAL, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_PERP {
		t.Fatalf("was expecting CRPERP type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.DERIBIT.Name {
		t.Fatalf("was expecting DERIBIT exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", sec.QuoteCurrency.ID)
	}
	if !sec.IsInverse {
		t.Fatalf("was expecting inverse, got non inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.5) > 0.000001 {
		t.Fatalf("was expecting 0.1 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-10.) > 0.0000000001 {
		t.Fatalf("was expecting 0.0001 round lot increment, got %g", sec.RoundLot.Value)
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Agg Trades: %d | Trades: %d | OBUpdates: %d", stats.AggTrades, stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestHuobip(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		10070367938184144403,
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
		fmt.Println(s)
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
	sec, ok := testedSecurities[10070367938184144403]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != "BTC-USD" {
		t.Fatalf("was expecting symbol BTC-USD, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_PERP {
		t.Fatalf("was expecting CRPERP type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.HUOBIP.Name {
		t.Fatalf("was expecting HUOBIP exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", sec.QuoteCurrency.ID)
	}
	if !sec.IsInverse {
		t.Fatalf("was expecting inverse, got non inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.1) > 0.000001 {
		t.Fatalf("was expecting 0.1 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-1.) > 0.0000000001 {
		t.Fatalf("was expecting 1 round lot increment, got %g", sec.RoundLot.Value)
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Agg Trades: %d | Trades: %d | OBUpdates: %d", stats.AggTrades, stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestHuobif(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		5362427922299651408,
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
		fmt.Println(s)
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
	sec, ok := testedSecurities[5362427922299651408]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != "BTC201225" {
		t.Fatalf("was expecting symbol BTC200814, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_FUT {
		t.Fatalf("was expecting CRFUT type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.HUOBIF.Name {
		t.Fatalf("was expecting HUOBIF exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", sec.QuoteCurrency.ID)
	}
	if !sec.IsInverse {
		t.Fatalf("was expecting inverse, got non inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.1 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-1.) > 0.0000000001 {
		t.Fatalf("was expecting 1 round lot increment, got %g", sec.RoundLot.Value)
	}
	if sec.MaturityDate == nil {
		t.Fatalf("was expecting maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Agg Trades: %d | Trades: %d | OBUpdates: %d", stats.AggTrades, stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestBybiti(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		7374647908427501521,
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
		fmt.Println(s)
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
	sec, ok := testedSecurities[7374647908427501521]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != "BTCUSD" {
		t.Fatalf("was expecting symbol BTCUSD, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_PERP {
		t.Fatalf("was expecting CRPERP type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.BYBITI.Name {
		t.Fatalf("was expecting BYBITI exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", sec.QuoteCurrency.ID)
	}
	if !sec.IsInverse {
		t.Fatalf("was expecting inverse, got non inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.5) > 0.000001 {
		t.Fatalf("was expecting 0.5 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-1.) > 0.0000000001 {
		t.Fatalf("was expecting 1 round lot increment, got %g", sec.RoundLot.Value)
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was not expecting maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Agg Trades: %d | Trades: %d | OBUpdates: %d", stats.AggTrades, stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestBybitl(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		6789757764526280996,
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
	sec, ok := testedSecurities[6789757764526280996]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != "BTCUSDT" {
		t.Fatalf("was expecting symbol BTCUSDT, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_PERP {
		t.Fatalf("was expecting CRPERP type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.BYBITL.Name {
		t.Fatalf("was expecting BYBITL exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.TETHER.ID {
		t.Fatalf("was expecting USDT quote, got %d", sec.QuoteCurrency.ID)
	}
	if !sec.IsInverse {
		t.Fatalf("was expecting inverse, got non inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.5) > 0.000001 {
		t.Fatalf("was expecting 0.5 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-0.001) > 0.0000000001 {
		t.Fatalf("was expecting 1 round lot increment, got %g", sec.RoundLot.Value)
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was not expecting maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Agg Trades: %d | Trades: %d | OBUpdates: %d", stats.AggTrades, stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestUpbit(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		13641637530641868249,
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
	sec, ok := testedSecurities[13641637530641868249]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != "KRW-BTC" {
		t.Fatalf("was expecting symbol BTCUSDT, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRPERP type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.UPBIT.Name {
		t.Fatalf("was expecting UPBIT exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.SOUTH_KOREAN_WON.ID {
		t.Fatalf("was expecting KRW quote, got %d", sec.QuoteCurrency.ID)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if sec.MinPriceIncrement != nil {
		t.Fatalf("was not expecting min price increment")
	}
	if sec.RoundLot != nil {
		t.Fatalf("was not expecting round lot")
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was not expecting maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Agg Trades: %d | Trades: %d | OBUpdates: %d", stats.AggTrades, stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestBithumb(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		3580725199953192415,
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
	sec, ok := testedSecurities[3580725199953192415]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != "BTC_KRW" {
		t.Fatalf("was expecting symbol BTCUSDT, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.BITHUMB.Name {
		t.Fatalf("was expecting BITHUMB exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting BTC underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.SOUTH_KOREAN_WON.ID {
		t.Fatalf("was expecting SOUTH_KOREAN_WON quote, got %s", sec.QuoteCurrency.Name)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if sec.MinPriceIncrement != nil {
		t.Fatalf("was expecting nil min price increment")
	}
	if sec.RoundLot != nil {
		t.Fatalf("was expecting nil round lot")
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was not expecting maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Agg Trades: %d | Trades: %d | OBUpdates: %d", stats.AggTrades, stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestBithumbg(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		8219034216918946889,
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
	sec, ok := testedSecurities[8219034216918946889]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != "BTC-USDT" {
		t.Fatalf("was expecting symbol BTCUSDT, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRPERP type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.BITHUMBG.Name {
		t.Fatalf("was expecting BITHUMBG exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting BTC underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.TETHER.ID {
		t.Fatalf("was expecting USDT quote, got %d", sec.QuoteCurrency.ID)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if sec.MinPriceIncrement == nil || sec.MinPriceIncrement.Value != 0.01 {
		t.Fatalf("was expecting min price increment of 0.01")
	}
	if sec.RoundLot == nil || sec.RoundLot.Value != 0.000001 {
		t.Fatalf("was expecting round lot of 0.000001")
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was not expecting maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Agg Trades: %d | Trades: %d | OBUpdates: %d", stats.AggTrades, stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestDydx(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		5390294902314234278,
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
	sec, ok := testedSecurities[5390294902314234278]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != "BTC-USD" {
		t.Fatalf("was expecting symbol BTCUSD, got %s", sec.Symbol)
	}
	if sec.SecurityType != enum.SecurityType_CRYPTO_PERP {
		t.Fatalf("was expecting CRPERP type, got %s", sec.SecurityType)
	}
	if sec.Exchange.Name != constants.DYDX.Name {
		t.Fatalf("was expecting BITHUMBG exchange, got %s", sec.Exchange.Name)
	}
	if sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting BTC underlying, got %d", sec.Underlying.ID)
	}
	if sec.QuoteCurrency.ID != constants.USDC.ID {
		t.Fatalf("was expecting USDC quote, got %d", sec.QuoteCurrency.ID)
	}
	if sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if sec.Status != models.Trading {
		t.Fatal("was expecting enabled security")
	}
	if sec.MinPriceIncrement == nil || sec.MinPriceIncrement.Value != 1 {
		t.Fatalf("was expecting min price increment of 0.01")
	}
	if sec.RoundLot == nil || sec.RoundLot.Value != 0.0001 {
		t.Fatalf("was expecting round lot of 0.0001")
	}
	if sec.MaturityDate != nil {
		t.Fatalf("was not expecting maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(tests.NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = as.Root.RequestFuture(obChecker, &tests.GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*tests.GetStat)
	t.Logf("Agg Trades: %d | Trades: %d | OBUpdates: %d", stats.AggTrades, stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}
