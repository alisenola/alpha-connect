package data

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alphac/enum"
	"gitlab.com/alphaticks/alphac/exchanges"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"math"
	"os"
	"reflect"
	"testing"
	"time"
)

type GetStat struct {
	Error     error
	Trades    int
	AggTrades int
	OBUpdates int
}

type OBChecker struct {
	security      *models.Security
	orderbook     *gorderbook.OrderBookL2
	tickPrecision uint64
	lotPrecision  uint64
	seqNum        uint64
	synced        bool
	trades        int
	aggTrades     int
	aggTradeIDs   map[uint64]bool
	OBUpdates     int
	err           error
}

func NewOBCheckerProducer(security *models.Security) actor.Producer {
	return func() actor.Actor {
		return NewOBChecker(security)
	}
}

func NewOBChecker(security *models.Security) actor.Actor {
	return &OBChecker{
		security:    security,
		orderbook:   nil,
		seqNum:      0,
		synced:      false,
		trades:      0,
		aggTrades:   0,
		aggTradeIDs: make(map[uint64]bool),
		OBUpdates:   0,
		err:         nil,
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
			AggTrades: state.aggTrades,
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
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: state.security.SecurityID},
			Exchange:   state.security.Exchange,
			Symbol:     &types.StringValue{Value: state.security.Symbol},
		},
		Aggregation: models.L2,
	}, 20*time.Second).Result()
	if err != nil {
		return err
	}
	response, ok := res.(*messages.MarketDataResponse)
	if !ok {
		return fmt.Errorf("was expecting market data snapshot, got %s", reflect.TypeOf(res).String())
	}

	var tickPrecision uint64
	if response.SnapshotL2.TickPrecision != nil {
		tickPrecision = response.SnapshotL2.TickPrecision.Value
	} else if state.security.MinPriceIncrement != nil {
		tickPrecision = uint64(math.Ceil(1. / state.security.MinPriceIncrement.Value))
	} else {
		return fmt.Errorf("unable to get tick precision")
	}

	var lotPrecision uint64
	if response.SnapshotL2.LotPrecision != nil {
		lotPrecision = response.SnapshotL2.LotPrecision.Value
	} else if state.security.RoundLot != nil {
		lotPrecision = uint64(math.Ceil(1. / state.security.RoundLot.Value))
	} else {
		return fmt.Errorf("unable to get lo precision")
	}

	state.tickPrecision = tickPrecision
	state.lotPrecision = lotPrecision

	for _, b := range response.SnapshotL2.Bids {
		rawPrice := b.Price * float64(tickPrecision)
		rawQty := b.Quantity * float64(lotPrecision)
		if (math.Round(rawPrice) - rawPrice) > 0.00001 {
			return fmt.Errorf("residue in price: %f %f", rawPrice, math.Round(rawPrice))
		}
		if (math.Round(rawQty) - rawQty) > 0.01 {
			return fmt.Errorf("residue in qty: %f %f", rawQty, math.Round(rawQty))
		}
	}
	for _, a := range response.SnapshotL2.Asks {
		rawPrice := a.Price * float64(tickPrecision)
		rawQty := a.Quantity * float64(lotPrecision)
		if (math.Round(rawPrice) - rawPrice) > 0.00001 {
			return fmt.Errorf("residue in price: %f %f", rawPrice, math.Round(rawPrice))
		}
		if (math.Round(rawQty) - rawQty) > 0.01 {
			return fmt.Errorf("residue in qty: %f %f", rawQty, math.Round(rawQty))
		}
	}
	state.OBUpdates += 1
	state.orderbook = gorderbook.NewOrderBookL2(uint64(tickPrecision), uint64(lotPrecision), 10000)
	state.orderbook.Sync(response.SnapshotL2.Bids, response.SnapshotL2.Asks)
	state.seqNum = response.SeqNum
	if state.orderbook.Crossed() {
		return fmt.Errorf("crossed OB on snapshot \n" + state.orderbook.String())
	}
	return nil
}

func (state *OBChecker) OnMarketDataIncrementalRefresh(context actor.Context) error {
	refresh := context.Message().(*messages.MarketDataIncrementalRefresh)

	if !state.synced && refresh.SeqNum <= state.seqNum {
		//fmt.Println("SKIPPING", refresh.SeqNum, state.securityInfo.seqNum)
		return nil
	}
	state.synced = true
	if state.seqNum+1 != refresh.SeqNum {
		//fmt.Println("OUT OF SYNC", state.securityInfo.seqNum, refresh.SeqNum)
		return fmt.Errorf("out of order sequence %d %d", state.seqNum, refresh.SeqNum)
	}

	if refresh.UpdateL2 != nil {
		for _, l := range refresh.UpdateL2.Levels {
			rawPrice := l.Price * float64(state.tickPrecision)
			rawQty := l.Quantity * float64(state.lotPrecision)
			if (math.Round(rawPrice) - rawPrice) > 0.00001 {
				return fmt.Errorf("residue in ob price: %f %f", rawPrice, math.Round(rawPrice))
			}
			if (math.Round(rawQty) - rawQty) > 0.0001 {
				return fmt.Errorf("residue in ob qty: %f %f", rawQty, math.Round(rawQty))
			}
			state.orderbook.UpdateOrderBookLevel(l)
		}
		if state.orderbook.Crossed() {
			fmt.Println("CROSSED")
			for _, l := range refresh.UpdateL2.Levels {
				fmt.Println(l)
			}
			return fmt.Errorf("crossed OB \n" + state.orderbook.String())
		}
		state.OBUpdates += 1
	}

	for _, aggT := range refresh.Trades {
		if _, ok := state.aggTradeIDs[aggT.AggregateID]; ok {
			return fmt.Errorf("duplicate aggregate ID")
		}
		state.aggTradeIDs[aggT.AggregateID] = true
		state.aggTrades += 1
		for _, t := range aggT.Trades {
			rawPrice := t.Price * float64(state.tickPrecision)
			rawQty := t.Quantity * float64(state.lotPrecision)
			if (math.Round(rawPrice) - rawPrice) > 0.00001 {
				return fmt.Errorf("residue in trade price: %f %f", rawPrice, math.Round(rawPrice))
			}
			if (math.Round(rawQty) - rawQty) > 0.0001 {
				return fmt.Errorf("residue in trade qty: %f %f", rawQty, math.Round(rawQty))
			}
			state.trades += 1
		}
	}
	state.seqNum = refresh.SeqNum
	return nil
}

func (state *OBChecker) OnMarketDataSnapshot(context actor.Context) error {
	return nil
}

var executor *actor.PID

func TestMain(m *testing.M) {
	exch := []*xchangerModels.Exchange{
		&constants.BINANCE,
		&constants.BITFINEX,
		&constants.BITSTAMP,
		&constants.COINBASEPRO,
		&constants.GEMINI,
		&constants.KRAKEN,
		&constants.CRYPTOFACILITIES,
		&constants.OKCOIN,
		&constants.FBINANCE,
		&constants.HITBTC,
		&constants.BITZ,
		&constants.HUOBI,
		&constants.FTX,
		&constants.BITMEX,
		&constants.BITSTAMP,
		&constants.DERIBIT,
		&constants.HUOBIP,
		&constants.HUOBIF,
		&constants.BYBITI,
		&constants.BYBITL,
	}
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(exchanges.NewExecutorProducer(exch, nil, false, xchangerUtils.DefaultDialerPool)), "executor")
	code := m.Run()
	_ = actor.EmptyRootContext.PoisonFuture(executor).Wait()
	os.Exit(code)
}

func TestBinance(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()

	securityID := []uint64{
		9281941173829172773, //BTCUSDT
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
	if !sec.Enabled {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %f", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-0.000001) > 0.0000000001 {
		t.Fatalf("was expecting 0.000001 round lot increment, got %f", sec.RoundLot.Value)
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
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
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()

	securityID := []uint64{
		17873758715870285590, //BTCUSDT
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
	if !sec.Enabled {
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

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
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
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		5391998915988476130, //XBTUSD
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
	if !sec.Enabled {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.5) > 0.000001 {
		t.Fatalf("was expecting 0.5 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-1.) > 0.0000000001 {
		t.Fatalf("was expecting 1. round lot increment, got %g", sec.RoundLot.Value)
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
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
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		5279696656781449381,
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
	if !sec.Enabled {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-0.00000001) > 0.0000000001 {
		t.Fatalf("was expecting 0.00000001 round lot increment, got %g", sec.RoundLot.Value)
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
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
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		243278890145991530,
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
	if !sec.Enabled {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-0.0001) > 0.0000000001 {
		t.Fatalf("was expecting 0.0001 round lot increment, got %g", sec.RoundLot.Value)
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
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
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
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
	if !sec.Enabled {
		t.Fatal("was expecting enabled security")
	}
	if math.Abs(sec.MinPriceIncrement.Value-0.01) > 0.000001 {
		t.Fatalf("was expecting 0.01 min price increment, got %g", sec.MinPriceIncrement.Value)
	}
	if math.Abs(sec.RoundLot.Value-0.00000001) > 0.0000000001 {
		t.Fatalf("was expecting 0.0001 round lot increment, got %g", sec.RoundLot.Value)
	}

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
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
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		1416768858288349990,
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
	if !sec.Enabled {
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

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
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
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		5485975358912730733,
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
	if !sec.Enabled {
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

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
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
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		4425198260936995601,
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
	if !sec.Enabled {
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

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
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

func TestHuobi(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		2195469462990134438,
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
	if !sec.Enabled {
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

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
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
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		17496373742670049989,
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
	if !sec.Enabled {
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

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
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
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		12674447834540883135,
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
	if !sec.Enabled {
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

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
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
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		10955098577666557860,
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
	if !sec.Enabled {
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

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
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
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		16235745264492357730,
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
	if !sec.Enabled {
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

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
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

func TestDeribit(t *testing.T) {
	t.Parallel()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		2206542817128348325,
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
	if !sec.Enabled {
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

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
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
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		10070367938184144403,
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
	if !sec.Enabled {
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

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
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
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		5362427922299651408,
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
	if !sec.Enabled {
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

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
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
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		7374647908427501521,
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
	if !sec.Enabled {
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

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
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
			_ = actor.EmptyRootContext.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		6789757764526280996,
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
	if !sec.Enabled {
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

	obChecker = actor.EmptyRootContext.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
	res, err = actor.EmptyRootContext.RequestFuture(obChecker, &GetStat{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	stats := res.(*GetStat)
	t.Logf("Agg Trades: %d | Trades: %d | OBUpdates: %d", stats.AggTrades, stats.Trades, stats.OBUpdates)
	if stats.Error != nil {
		t.Fatal(stats.Error)
	}
}
