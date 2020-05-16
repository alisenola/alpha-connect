package exchanges

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alphac/enum"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"math"
	"reflect"
	"testing"
	"time"
)

type GetStat struct {
	Error     error
	Trades    int
	OBUpdates int
}

type OBChecker struct {
	security  *models.Security
	orderbook *gorderbook.OrderBookL2
	seqNum    uint64
	trades    int
	OBUpdates int
	err       error
}

func NewOBCheckerProducer(security *models.Security) actor.Producer {
	return func() actor.Actor {
		return NewOBChecker(security)
	}
}

func NewOBChecker(security *models.Security) actor.Actor {
	return &OBChecker{
		security:  security,
		orderbook: nil,
		seqNum:    0,
		trades:    0,
		OBUpdates: 0,
		err:       nil,
	}
}

func (state *OBChecker) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.err = err
		}

	case *messages.MarketDataIncrementalRefresh:
		if state.err == nil {
			if err := state.OnMarketDataIncrementalRefresh(context); err != nil {
				state.err = err
			}
		}

	case *GetStat:
		context.Respond(&GetStat{
			Error:     state.err,
			Trades:    state.trades,
			OBUpdates: state.OBUpdates,
		})
	}
}

func (state *OBChecker) Initialize(context actor.Context) error {
	executor := actor.NewLocalPID("executor")
	res, err := actor.EmptyRootContext.RequestFuture(executor, &messages.MarketDataRequest{
		RequestID:  0,
		Subscribe:  true,
		Subscriber: context.Self(),
		Instrument: &messages.Instrument{
			SecurityID: state.security.SecurityID,
			Exchange:   state.security.Exchange,
			Symbol:     state.security.Symbol,
		},
		Aggregation: messages.L2,
	}, 10*time.Second).Result()
	if err != nil {
		return err
	}
	snapshot, ok := res.(*messages.MarketDataSnapshot)
	if !ok {
		return fmt.Errorf("was expecting market data snapshot, got %s", reflect.TypeOf(res).String())
	}

	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot))
	for _, b := range snapshot.SnapshotL2.Bids {
		rawPrice := b.Price * float64(tickPrecision)
		rawQty := b.Quantity * float64(lotPrecision)
		if (math.Round(rawPrice) - rawPrice) > 0.00001 {
			return fmt.Errorf("residue in price: %f %f", rawPrice, math.Round(rawPrice))
		}
		if (math.Round(rawQty) - rawQty) > 0.00001 {
			return fmt.Errorf("residue in qty: %f %f", rawQty, math.Round(rawQty))
		}
	}
	for _, a := range snapshot.SnapshotL2.Asks {
		rawPrice := a.Price * float64(tickPrecision)
		rawQty := a.Quantity * float64(lotPrecision)
		if (math.Round(rawPrice) - rawPrice) > 0.00001 {
			return fmt.Errorf("residue in price: %f %f", rawPrice, math.Round(rawPrice))
		}
		if (math.Round(rawQty) - rawQty) > 0.00001 {
			return fmt.Errorf("residue in qty: %f %f", rawQty, math.Round(rawQty))
		}
	}
	state.OBUpdates += 1
	state.orderbook = gorderbook.NewOrderBookL2(tickPrecision, lotPrecision, 10000)
	state.orderbook.Sync(snapshot.SnapshotL2.Bids, snapshot.SnapshotL2.Asks)
	state.seqNum = snapshot.SnapshotL2.SeqNum
	if state.orderbook.Crossed() {
		return fmt.Errorf("crossed OB on snapshot \n" + state.orderbook.String())
	}
	return nil
}

func (state *OBChecker) OnMarketDataIncrementalRefresh(context actor.Context) error {
	tickPrecision := uint64(math.Ceil(1. / state.security.MinPriceIncrement))
	lotPrecision := uint64(math.Ceil(1. / state.security.RoundLot))

	refresh := context.Message().(*messages.MarketDataIncrementalRefresh)
	if refresh.UpdateL2 != nil && refresh.UpdateL2.SeqNum > state.seqNum {
		if state.seqNum+1 != refresh.UpdateL2.SeqNum {
			return fmt.Errorf("out of order sequence %d %d", state.seqNum, refresh.UpdateL2.SeqNum)
		}
		for _, l := range refresh.UpdateL2.Levels {
			rawPrice := l.Price * float64(tickPrecision)
			rawQty := l.Quantity * float64(lotPrecision)
			if (math.Round(rawPrice) - rawPrice) > 0.00001 {
				return fmt.Errorf("residue in ob price: %f %f", rawPrice, math.Round(rawPrice))
			}
			if (math.Round(rawQty) - rawQty) > 0.00001 {
				return fmt.Errorf("residue in ob qty: %f %f", rawQty, math.Round(rawQty))
			}
			state.orderbook.UpdateOrderBookLevel(l)
		}
		if state.orderbook.Crossed() {
			return fmt.Errorf("crossed OB \n" + state.orderbook.String())
		}
		state.seqNum = refresh.UpdateL2.SeqNum
		state.OBUpdates += 1
	}

	for _, aggT := range refresh.Trades {
		for _, t := range aggT.Trades {
			rawPrice := t.Price * float64(tickPrecision)
			rawQty := t.Quantity * float64(lotPrecision)
			if (math.Round(rawPrice) - rawPrice) > 0.00001 {
				return fmt.Errorf("residue in trade price: %f %f", rawPrice, math.Round(rawPrice))
			}
			if (math.Round(rawQty) - rawQty) > 0.00001 {
				return fmt.Errorf("residue in trade qty: %f %f", rawQty, math.Round(rawQty))
			}
			state.trades += 1
		}
	}
	return nil
}

func (state *OBChecker) OnMarketDataSnapshot(context actor.Context) error {
	return nil
}

var obChecker *actor.PID
var executor *actor.PID

func clean() {
	_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
	_ = actor.EmptyRootContext.PoisonFuture(executor).Wait()
}

func TestBinance(t *testing.T) {
	defer clean()

	exchanges := []*xchangerModels.Exchange{&constants.BINANCE}
	securityID := []uint64{
		9281941173829172773, //BTCUSDT
	}
	testedSecurities := make(map[uint64]*models.Security)
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges)), "executor")

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

	// Test BTCUSDT
	Sec, ok := testedSecurities[9281941173829172773]
	if !ok {
		t.Fatalf("BTCUSDT not found")
	}
	if Sec.Symbol != "BTCUSDT" {
		t.Fatalf("was expecting symbol BTCEUR, got %s", Sec.Symbol)
	}
	if Sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", Sec.SecurityType)
	}
	if Sec.Exchange.Name != constants.BINANCE.Name {
		t.Fatalf("was expecting binance exchange, got %s", Sec.Exchange.Name)
	}
	if Sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", Sec.Underlying.ID)
	}
	if Sec.QuoteCurrency.ID != constants.TETHER.ID {
		t.Fatalf("was expecting EUR quote, got %d", Sec.QuoteCurrency.ID)
	}
	if Sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if math.Abs(Sec.MinPriceIncrement-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %f", Sec.MinPriceIncrement)
	}
	if math.Abs(Sec.RoundLot-0.000001) > 0.0000000001 {
		t.Fatalf("was expecting 0.000001 round lot increment, got %f", Sec.RoundLot)
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(Sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestBitfinex(t *testing.T) {
	defer clean()
	exchanges := []*xchangerModels.Exchange{&constants.BITFINEX}
	securityID := []uint64{
		17873758715870285590, //BTCUSDT
	}
	testedSecurities := make(map[uint64]*models.Security)
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges)), "executor")

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

	// Test BTCEUR
	BTCUSDSec, ok := testedSecurities[17873758715870285590]
	if !ok {
		t.Fatalf("BTCUSD not found")
	}
	if BTCUSDSec.Symbol != "btcusd" {
		t.Fatalf("was expecting symbol btcust, got %s", BTCUSDSec.Symbol)
	}
	if BTCUSDSec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRPERP type, got %s", BTCUSDSec.SecurityType)
	}
	if BTCUSDSec.Exchange.Name != constants.BITFINEX.Name {
		t.Fatalf("was expecting BITFINEX exchange, got %s", BTCUSDSec.Exchange.Name)
	}
	if BTCUSDSec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", BTCUSDSec.Underlying.ID)
	}
	if BTCUSDSec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", BTCUSDSec.QuoteCurrency.ID)
	}
	if BTCUSDSec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if math.Abs(BTCUSDSec.MinPriceIncrement-0.1) > 0.000001 {
		t.Fatalf("was expecting 0.1 min price increment, got %f", BTCUSDSec.MinPriceIncrement)
	}
	if math.Abs(BTCUSDSec.RoundLot-1./100000000.) > 0.0000000001 {
		t.Fatalf("was expecting 0.000001 round lot increment, got %f", BTCUSDSec.RoundLot)
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(BTCUSDSec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestBitmex(t *testing.T) {
	defer clean()
	exchanges := []*xchangerModels.Exchange{&constants.BITMEX}
	securityID := []uint64{
		5391998915988476130, //XBTUSD
	}
	testedSecurities := make(map[uint64]*models.Security)
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges)), "executor")

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
	Sec, ok := testedSecurities[5391998915988476130]
	if !ok {
		t.Fatalf("XBTUSD not found")
	}
	if Sec.Symbol != "XBTUSD" {
		t.Fatalf("was expecting symbol btcust, got %s", Sec.Symbol)
	}
	if Sec.SecurityType != enum.SecurityType_CRYPTO_PERP {
		t.Fatalf("was expecting CRSPOT type, got %s", Sec.SecurityType)
	}
	if Sec.Exchange.Name != constants.BITMEX.Name {
		t.Fatalf("was expecting BITFINEX exchange, got %s", Sec.Exchange.Name)
	}
	if Sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", Sec.Underlying.ID)
	}
	if Sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", Sec.QuoteCurrency.ID)
	}
	if Sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if math.Abs(Sec.MinPriceIncrement-0.5) > 0.000001 {
		t.Fatalf("was expecting 0.5 min price increment, got %g", Sec.MinPriceIncrement)
	}
	if math.Abs(Sec.RoundLot-1.) > 0.0000000001 {
		t.Fatalf("was expecting 1. round lot increment, got %g", Sec.RoundLot)
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(Sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestBitstamp(t *testing.T) {
	defer clean()
	exchanges := []*xchangerModels.Exchange{&constants.BITSTAMP}
	securityID := []uint64{
		5279696656781449381,
	}
	testedSecurities := make(map[uint64]*models.Security)
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges)), "executor")

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
	Sec, ok := testedSecurities[5279696656781449381]
	if !ok {
		t.Fatalf("BTCUSD not found")
	}
	if Sec.Symbol != "btcusd" {
		t.Fatalf("was expecting symbol BTC/USD, got %s", Sec.Symbol)
	}
	if Sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", Sec.SecurityType)
	}
	if Sec.Exchange.Name != constants.BITSTAMP.Name {
		t.Fatalf("was expecting BITFINEX exchange, got %s", Sec.Exchange.Name)
	}
	if Sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", Sec.Underlying.ID)
	}
	if Sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", Sec.QuoteCurrency.ID)
	}
	if Sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if math.Abs(Sec.MinPriceIncrement-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", Sec.MinPriceIncrement)
	}
	if math.Abs(Sec.RoundLot-0.00000001) > 0.0000000001 {
		t.Fatalf("was expecting 0.00000001 round lot increment, got %g", Sec.RoundLot)
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(Sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestBitz(t *testing.T) {
	defer clean()
	exchanges := []*xchangerModels.Exchange{&constants.BITZ}
	securityID := []uint64{
		243278890145991530,
	}
	testedSecurities := make(map[uint64]*models.Security)
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges)), "executor")

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
	Sec, ok := testedSecurities[243278890145991530]
	if !ok {
		t.Fatalf("BTCUSDT not found")
	}
	if Sec.Symbol != "btc_usdt" {
		t.Fatalf("was expecting symbol btc_usdt, got %s", Sec.Symbol)
	}
	if Sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", Sec.SecurityType)
	}
	if Sec.Exchange.Name != constants.BITZ.Name {
		t.Fatalf("was expecting bitz exchange, got %s", Sec.Exchange.Name)
	}
	if Sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", Sec.Underlying.ID)
	}
	if Sec.QuoteCurrency.ID != constants.TETHER.ID {
		t.Fatalf("was expecting USDT quote, got %d", Sec.QuoteCurrency.ID)
	}
	if Sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if math.Abs(Sec.MinPriceIncrement-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", Sec.MinPriceIncrement)
	}
	if math.Abs(Sec.RoundLot-0.0001) > 0.0000000001 {
		t.Fatalf("was expecting 0.0001 round lot increment, got %g", Sec.RoundLot)
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(Sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestCoinbasePro(t *testing.T) {
	defer clean()
	exchanges := []*xchangerModels.Exchange{&constants.COINBASEPRO}
	securityID := []uint64{
		11630614572540763252,
	}
	testedSecurities := make(map[uint64]*models.Security)
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges)), "executor")

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
	Sec, ok := testedSecurities[11630614572540763252]
	if !ok {
		t.Fatalf("BTCUSD not found")
	}
	if Sec.Symbol != "BTC-USD" {
		t.Fatalf("was expecting symbol BTC-USD, got %s", Sec.Symbol)
	}
	if Sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", Sec.SecurityType)
	}
	if Sec.Exchange.Name != constants.COINBASEPRO.Name {
		t.Fatalf("was expecting bitz exchange, got %s", Sec.Exchange.Name)
	}
	if Sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", Sec.Underlying.ID)
	}
	if Sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", Sec.QuoteCurrency.ID)
	}
	if Sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if math.Abs(Sec.MinPriceIncrement-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", Sec.MinPriceIncrement)
	}
	if math.Abs(Sec.RoundLot-0.00000001) > 0.0000000001 {
		t.Fatalf("was expecting 0.0001 round lot increment, got %g", Sec.RoundLot)
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(Sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestCryptofacilities(t *testing.T) {
	defer clean()
	exchanges := []*xchangerModels.Exchange{&constants.CRYPTOFACILITIES}
	securityID := []uint64{
		1416768858288349990,
	}
	testedSecurities := make(map[uint64]*models.Security)
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges)), "executor")

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
	Sec, ok := testedSecurities[1416768858288349990]
	if !ok {
		t.Fatalf("XBTUSD not found")
	}
	if Sec.Symbol != "pi_xbtusd" {
		t.Fatalf("was expecting symbol pi_xbtusd, got %s", Sec.Symbol)
	}
	if Sec.SecurityType != enum.SecurityType_CRYPTO_PERP {
		t.Fatalf("was expecting CRPERP type, got %s", Sec.SecurityType)
	}
	if Sec.Exchange.Name != constants.CRYPTOFACILITIES.Name {
		t.Fatalf("was expecting bitz exchange, got %s", Sec.Exchange.Name)
	}
	if Sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", Sec.Underlying.ID)
	}
	if Sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", Sec.QuoteCurrency.ID)
	}
	if !Sec.IsInverse {
		t.Fatalf("was expecting inverse, got non inverse")
	}
	if math.Abs(Sec.MinPriceIncrement-0.5) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", Sec.MinPriceIncrement)
	}
	if math.Abs(Sec.RoundLot-1.) > 0.0000000001 {
		t.Fatalf("was expecting 0.0001 round lot increment, got %g", Sec.RoundLot)
	}
	if Sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(Sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestFBinance(t *testing.T) {
	defer clean()
	exchanges := []*xchangerModels.Exchange{&constants.FBINANCE}
	securityID := []uint64{
		5485975358912730733,
	}
	testedSecurities := make(map[uint64]*models.Security)
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges)), "executor")

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
	Sec, ok := testedSecurities[5485975358912730733]
	if !ok {
		t.Fatalf("BTCUSDT not found")
	}
	if Sec.Symbol != "BTCUSDT" {
		t.Fatalf("was expecting symbol BTCUSDT, got %s", Sec.Symbol)
	}
	if Sec.SecurityType != enum.SecurityType_CRYPTO_PERP {
		t.Fatalf("was expecting CRPERP type, got %s", Sec.SecurityType)
	}
	if Sec.Exchange.Name != constants.FBINANCE.Name {
		t.Fatalf("was expecting bitz exchange, got %s", Sec.Exchange.Name)
	}
	if Sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", Sec.Underlying.ID)
	}
	if Sec.QuoteCurrency.ID != constants.TETHER.ID {
		t.Fatalf("was expecting USDT quote, got %d", Sec.QuoteCurrency.ID)
	}
	if Sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if math.Abs(Sec.MinPriceIncrement-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", Sec.MinPriceIncrement)
	}
	if math.Abs(Sec.RoundLot-0.001) > 0.0000000001 {
		t.Fatalf("was expecting 0.001 round lot increment, got %g", Sec.RoundLot)
	}
	if Sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(Sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestFTX(t *testing.T) {
	defer clean()
	exchanges := []*xchangerModels.Exchange{&constants.FTX}
	securityID := []uint64{
		4425198260936995601,
	}
	testedSecurities := make(map[uint64]*models.Security)
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges)), "executor")

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
	Sec, ok := testedSecurities[4425198260936995601]
	if !ok {
		t.Fatalf("security not found")
	}
	if Sec.Symbol != "BTC-PERP" {
		t.Fatalf("was expecting symbol BTC-PERP, got %s", Sec.Symbol)
	}
	if Sec.SecurityType != enum.SecurityType_CRYPTO_PERP {
		t.Fatalf("was expecting CRPERP type, got %s", Sec.SecurityType)
	}
	if Sec.Exchange.Name != constants.FTX.Name {
		t.Fatalf("was expecting ftx exchange, got %s", Sec.Exchange.Name)
	}
	if Sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", Sec.Underlying.ID)
	}
	if Sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", Sec.QuoteCurrency.ID)
	}
	if Sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if math.Abs(Sec.MinPriceIncrement-0.5) > 0.000001 {
		t.Fatalf("was expecting 0.5 min price increment, got %g", Sec.MinPriceIncrement)
	}
	if math.Abs(Sec.RoundLot-0.0001) > 0.0000000001 {
		t.Fatalf("was expecting 0.0001 round lot increment, got %g", Sec.RoundLot)
	}
	if Sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(Sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
	if err := actor.EmptyRootContext.PoisonFuture(obChecker).Wait(); err != nil {
		t.Fatal(err)
	}
	if err := actor.EmptyRootContext.PoisonFuture(executor).Wait(); err != nil {
		t.Fatal(err)
	}
}

func TestHuobi(t *testing.T) {
	defer clean()
	exchanges := []*xchangerModels.Exchange{&constants.HUOBI}
	securityID := []uint64{
		2195469462990134438,
	}
	testedSecurities := make(map[uint64]*models.Security)
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges)), "executor")

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
	Sec, ok := testedSecurities[2195469462990134438]
	if !ok {
		t.Fatalf("security not found")
	}
	if Sec.Symbol != "btcusdt" {
		t.Fatalf("was expecting symbol btcusdt, got %s", Sec.Symbol)
	}
	if Sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", Sec.SecurityType)
	}
	if Sec.Exchange.Name != constants.HUOBI.Name {
		t.Fatalf("was expecting huobi exchange, got %s", Sec.Exchange.Name)
	}
	if Sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", Sec.Underlying.ID)
	}
	if Sec.QuoteCurrency.ID != constants.TETHER.ID {
		t.Fatalf("was expecting USD quote, got %d", Sec.QuoteCurrency.ID)
	}
	if Sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if math.Abs(Sec.MinPriceIncrement-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", Sec.MinPriceIncrement)
	}
	if math.Abs(Sec.RoundLot-1e-6) > 0.0000000001 {
		t.Fatalf("was expecting 1e-6 round lot increment, got %g", Sec.RoundLot)
	}
	if Sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(Sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestGemini(t *testing.T) {
	defer clean()
	exchanges := []*xchangerModels.Exchange{&constants.GEMINI}
	securityID := []uint64{
		17496373742670049989,
	}
	testedSecurities := make(map[uint64]*models.Security)
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges)), "executor")

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
	Sec, ok := testedSecurities[17496373742670049989]
	if !ok {
		t.Fatalf("security not found")
	}
	if Sec.Symbol != "btcusd" {
		t.Fatalf("was expecting symbol btcusd, got %s", Sec.Symbol)
	}
	if Sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", Sec.SecurityType)
	}
	if Sec.Exchange.Name != constants.GEMINI.Name {
		t.Fatalf("was expecting gemini exchange, got %s", Sec.Exchange.Name)
	}
	if Sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", Sec.Underlying.ID)
	}
	if Sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", Sec.QuoteCurrency.ID)
	}
	if Sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if math.Abs(Sec.MinPriceIncrement-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", Sec.MinPriceIncrement)
	}
	if math.Abs(Sec.RoundLot-1e-08) > 0.0000000001 {
		t.Fatalf("was expecting 1e-6 round lot increment, got %g", Sec.RoundLot)
	}
	if Sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(Sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestHitbtc(t *testing.T) {
	defer clean()
	exchanges := []*xchangerModels.Exchange{&constants.HITBTC}
	securityID := []uint64{
		12674447834540883135,
	}
	testedSecurities := make(map[uint64]*models.Security)
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges)), "executor")

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
	Sec, ok := testedSecurities[12674447834540883135]
	if !ok {
		t.Fatalf("security not found")
	}
	if Sec.Symbol != "BTCUSD" {
		t.Fatalf("was expecting symbol btcusd, got %s", Sec.Symbol)
	}
	if Sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", Sec.SecurityType)
	}
	if Sec.Exchange.Name != constants.HITBTC.Name {
		t.Fatalf("was expecting hitbtc exchange, got %s", Sec.Exchange.Name)
	}
	if Sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", Sec.Underlying.ID)
	}
	if Sec.QuoteCurrency.ID != constants.TETHER.ID {
		t.Fatalf("was expecting USDT quote, got %d", Sec.QuoteCurrency.ID)
	}
	if Sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if math.Abs(Sec.MinPriceIncrement-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", Sec.MinPriceIncrement)
	}
	if math.Abs(Sec.RoundLot-1e-05) > 0.0000000001 {
		t.Fatalf("was expecting 1e-6 round lot increment, got %g", Sec.RoundLot)
	}
	if Sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(Sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestKraken(t *testing.T) {
	defer clean()
	exchanges := []*xchangerModels.Exchange{&constants.KRAKEN}
	securityID := []uint64{
		10955098577666557860,
	}
	testedSecurities := make(map[uint64]*models.Security)
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges)), "executor")

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
	Sec, ok := testedSecurities[10955098577666557860]
	if !ok {
		t.Fatalf("security not found")
	}
	if Sec.Symbol != "XBT/USD" {
		t.Fatalf("was expecting symbol btcusd, got %s", Sec.Symbol)
	}
	if Sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", Sec.SecurityType)
	}
	if Sec.Exchange.Name != constants.KRAKEN.Name {
		t.Fatalf("was expecting kraken exchange, got %s", Sec.Exchange.Name)
	}
	if Sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", Sec.Underlying.ID)
	}
	if Sec.QuoteCurrency.ID != constants.DOLLAR.ID {
		t.Fatalf("was expecting USD quote, got %d", Sec.QuoteCurrency.ID)
	}
	if Sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if math.Abs(Sec.MinPriceIncrement-0.1) > 0.000001 {
		t.Fatalf("was expecting 0.1 min price increment, got %g", Sec.MinPriceIncrement)
	}
	if math.Abs(Sec.RoundLot-1e-08) > 0.0000000001 {
		t.Fatalf("was expecting 1e-8 round lot increment, got %g", Sec.RoundLot)
	}
	if Sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(Sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}

func TestOKCoin(t *testing.T) {
	defer clean()
	exchanges := []*xchangerModels.Exchange{&constants.OKCOIN}
	securityID := []uint64{
		16235745264492357730,
	}
	testedSecurities := make(map[uint64]*models.Security)
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges)), "executor")

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
	Sec, ok := testedSecurities[16235745264492357730]
	if !ok {
		t.Fatalf("security not found")
	}
	if Sec.Symbol != "BTC-USDT" {
		t.Fatalf("was expecting symbol BTC-USDT, got %s", Sec.Symbol)
	}
	if Sec.SecurityType != enum.SecurityType_CRYPTO_SPOT {
		t.Fatalf("was expecting CRSPOT type, got %s", Sec.SecurityType)
	}
	if Sec.Exchange.Name != constants.OKCOIN.Name {
		t.Fatalf("was expecting kraken exchange, got %s", Sec.Exchange.Name)
	}
	if Sec.Underlying.ID != constants.BITCOIN.ID {
		t.Fatalf("was expecting bitcoin underlying, got %d", Sec.Underlying.ID)
	}
	if Sec.QuoteCurrency.ID != constants.TETHER.ID {
		t.Fatalf("was expecting USDT quote, got %d", Sec.QuoteCurrency.ID)
	}
	if Sec.IsInverse {
		t.Fatalf("was expecting non inverse, got inverse")
	}
	if math.Abs(Sec.MinPriceIncrement-0.1) > 0.000001 {
		t.Fatalf("was expecting 0.1 min price increment, got %g", Sec.MinPriceIncrement)
	}
	if math.Abs(Sec.RoundLot-0.0001) > 0.0000000001 {
		t.Fatalf("was expecting 0.0001 round lot increment, got %g", Sec.RoundLot)
	}
	if Sec.MaturityDate != nil {
		t.Fatalf("was expecting nil maturity date")
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(Sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	t.Logf("Trades: %d | OBUpdates: %d", stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}
