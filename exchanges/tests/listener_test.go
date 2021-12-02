package tests

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"math"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestMain(m *testing.M) {

	code := m.Run()
	os.Exit(code)
}

func StartExecutor(exchange *xchangerModels.Exchange) (*actor.ActorSystem, *actor.PID, func()) {
	assets := map[uint32]xchangerModels.Asset{
		constants.DOLLAR.ID:           constants.DOLLAR,
		constants.EURO.ID:             constants.EURO,
		constants.POUND.ID:            constants.POUND,
		constants.CANADIAN_DOLLAR.ID:  constants.CANADIAN_DOLLAR,
		constants.JAPENESE_YEN.ID:     constants.JAPENESE_YEN,
		constants.BITCOIN.ID:          constants.BITCOIN,
		constants.LITECOIN.ID:         constants.LITECOIN,
		constants.ETHEREUM.ID:         constants.ETHEREUM,
		constants.RIPPLE.ID:           constants.RIPPLE,
		constants.TETHER.ID:           constants.TETHER,
		constants.SOUTH_KOREAN_WON.ID: constants.SOUTH_KOREAN_WON,
		constants.USDC.ID:             constants.USDC,
		constants.DASH.ID:             constants.DASH,
		50: xchangerModels.Asset{
			Symbol: "BTT",
			Name:   "",
			ID:     50,
		},
	}

	_ = constants.LoadAssets(assets)
	exch := []*xchangerModels.Exchange{
		exchange,
	}
	as := actor.NewActorSystem()
	cfg := &exchanges.ExecutorConfig{
		Exchanges: exch,
		Strict:    true,
	}
	executor, _ := as.Root.SpawnNamed(actor.PropsFromProducer(exchanges.NewExecutorProducer(cfg)), "executor")
	return as, executor, func() { _ = as.Root.PoisonFuture(executor).Wait() }
}

type MDTest struct {
	securityID        uint64
	symbol            string
	securityType      string
	exchange          xchangerModels.Exchange
	baseCurrency      xchangerModels.Asset
	quoteCurrency     xchangerModels.Asset
	minPriceIncrement float64
	roundLot          float64
	hasMaturityDate   bool
	isInverse         bool
	status            models.InstrumentStatus
}

var MDTests = []MDTest{
	{
		securityID:        1416768858288349990,
		symbol:            "pi_xbtusd",
		securityType:      enum.SecurityType_CRYPTO_PERP,
		exchange:          constants.CRYPTOFACILITIES,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.DOLLAR,
		minPriceIncrement: 0.5,
		roundLot:          1,
		hasMaturityDate:   false,
		isInverse:         true,
		status:            models.Trading,
	},
	{
		securityID:        9281941173829172773,
		symbol:            "BTCUSDT",
		securityType:      enum.SecurityType_CRYPTO_SPOT,
		exchange:          constants.BINANCE,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.TETHER,
		minPriceIncrement: 0.01,
		roundLot:          0.000010,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},
	{
		securityID:        17873758715870285590,
		symbol:            "btcusd",
		securityType:      enum.SecurityType_CRYPTO_SPOT,
		exchange:          constants.BITFINEX,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.DOLLAR,
		minPriceIncrement: 0.1,
		roundLot:          1. / 100000000.,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},
	{
		securityID:        5391998915988476130,
		symbol:            "XBTUSD",
		securityType:      enum.SecurityType_CRYPTO_PERP,
		exchange:          constants.BITMEX,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.DOLLAR,
		minPriceIncrement: 0.5,
		roundLot:          100.,
		hasMaturityDate:   false,
		isInverse:         true,
		status:            models.Trading,
	},

	{
		securityID:        5279696656781449381,
		symbol:            "btcusd",
		securityType:      enum.SecurityType_CRYPTO_SPOT,
		exchange:          constants.BITSTAMP,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.DOLLAR,
		minPriceIncrement: 0.01,
		roundLot:          0.00000001,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},

	{
		securityID:        5485975358912730733,
		symbol:            "BTCUSDT",
		securityType:      enum.SecurityType_CRYPTO_PERP,
		exchange:          constants.FBINANCE,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.TETHER,
		minPriceIncrement: 0.01,
		roundLot:          0.001,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},

	{
		securityID:        4425198260936995601,
		symbol:            "BTC-PERP",
		securityType:      enum.SecurityType_CRYPTO_PERP,
		exchange:          constants.FTX,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.DOLLAR,
		minPriceIncrement: 1,
		roundLot:          0.0001,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},

	{
		securityID:        2028777944171259534,
		symbol:            "BTC/USDT",
		securityType:      enum.SecurityType_CRYPTO_SPOT,
		exchange:          constants.FTXUS,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.TETHER,
		minPriceIncrement: 1,
		roundLot:          0.0001,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},

	{
		securityID:        2195469462990134438,
		symbol:            "btcusdt",
		securityType:      enum.SecurityType_CRYPTO_SPOT,
		exchange:          constants.HUOBI,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.TETHER,
		minPriceIncrement: 0.01,
		roundLot:          1e-6,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},
	{
		securityID:        17496373742670049989,
		symbol:            "btcusd",
		securityType:      enum.SecurityType_CRYPTO_SPOT,
		exchange:          constants.GEMINI,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.DOLLAR,
		minPriceIncrement: 0.01,
		roundLot:          1e-8,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},

	{
		securityID:        12674447834540883135,
		symbol:            "BTCUSD",
		securityType:      enum.SecurityType_CRYPTO_SPOT,
		exchange:          constants.HITBTC,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.TETHER,
		minPriceIncrement: 0.01,
		roundLot:          1e-5,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},

	{
		securityID:        10955098577666557860,
		symbol:            "XBT/USD",
		securityType:      enum.SecurityType_CRYPTO_SPOT,
		exchange:          constants.KRAKEN,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.DOLLAR,
		minPriceIncrement: 0.1,
		roundLot:          1e-8,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},
	{
		securityID:        945944519923594006,
		symbol:            "BTC-USDT",
		securityType:      enum.SecurityType_CRYPTO_SPOT,
		exchange:          constants.OKEX,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.TETHER,
		minPriceIncrement: 0.1,
		roundLot:          1e-08,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},

	{
		securityID:        10652256150546133071,
		symbol:            "BTC-USDT-SWAP",
		securityType:      enum.SecurityType_CRYPTO_PERP,
		exchange:          constants.OKEXP,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.TETHER,
		minPriceIncrement: 0.1,
		roundLot:          1,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},
	{
		securityID:        2206542817128348325,
		symbol:            "BTC-PERPETUAL",
		securityType:      enum.SecurityType_CRYPTO_PERP,
		exchange:          constants.DERIBIT,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.DOLLAR,
		minPriceIncrement: 0.5,
		roundLot:          10.,
		hasMaturityDate:   false,
		isInverse:         true,
		status:            models.Trading,
	},
	{
		securityID:        10070367938184144403,
		symbol:            "BTC-USD",
		securityType:      enum.SecurityType_CRYPTO_PERP,
		exchange:          constants.HUOBIP,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.DOLLAR,
		minPriceIncrement: 0.1,
		roundLot:          1.,
		hasMaturityDate:   false,
		isInverse:         true,
		status:            models.Trading,
	},

	{
		securityID:        2402007053666382556,
		symbol:            "BTC210625",
		securityType:      enum.SecurityType_CRYPTO_FUT,
		exchange:          constants.HUOBIF,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.DOLLAR,
		minPriceIncrement: 0.01,
		roundLot:          1.,
		hasMaturityDate:   true,
		isInverse:         false,
		status:            models.Trading,
	},
	{
		securityID:        7374647908427501521,
		symbol:            "BTCUSD",
		securityType:      enum.SecurityType_CRYPTO_PERP,
		exchange:          constants.BYBITI,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.DOLLAR,
		minPriceIncrement: 0.5,
		roundLot:          1.,
		hasMaturityDate:   false,
		isInverse:         true,
		status:            models.Trading,
	},
	{
		securityID:        6789757764526280996,
		symbol:            "BTCUSDT",
		securityType:      enum.SecurityType_CRYPTO_PERP,
		exchange:          constants.BYBITL,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.TETHER,
		minPriceIncrement: 0.5,
		roundLot:          0.001,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},
	{
		securityID:        13641637530641868249,
		symbol:            "KRW-BTC",
		securityType:      enum.SecurityType_CRYPTO_SPOT,
		exchange:          constants.UPBIT,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.SOUTH_KOREAN_WON,
		minPriceIncrement: 0.,
		roundLot:          0.,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},
	{
		securityID:        8219034216918946889,
		symbol:            "BTC-USDT",
		securityType:      enum.SecurityType_CRYPTO_SPOT,
		exchange:          constants.BITHUMBG,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.TETHER,
		minPriceIncrement: 0.01,
		roundLot:          1e-06,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},
	{
		securityID:        13112609607273681222,
		symbol:            "ETH-USD",
		securityType:      enum.SecurityType_CRYPTO_PERP,
		exchange:          constants.DYDX,
		baseCurrency:      constants.ETHEREUM,
		quoteCurrency:     constants.USDC,
		minPriceIncrement: 0.1,
		roundLot:          0.001,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},
	{
		securityID:        11630614572540763252,
		symbol:            "BTC-USD",
		securityType:      enum.SecurityType_CRYPTO_SPOT,
		exchange:          constants.COINBASEPRO,
		baseCurrency:      constants.BITCOIN,
		quoteCurrency:     constants.DOLLAR,
		minPriceIncrement: 0.01,
		roundLot:          1e-08,
		hasMaturityDate:   false,
		isInverse:         false,
		status:            models.Trading,
	},
}

func TestAll(t *testing.T) {
	t.Parallel()
	for _, tc := range MDTests {
		tc := tc
		t.Run(tc.exchange.Name, func(t *testing.T) {
			t.Parallel()
			MarketData(t, tc)
		})
	}
}

func MarketData(t *testing.T, test MDTest) {
	//t.Parallel()
	as, executor, clean := StartExecutor(&test.exchange)
	defer clean()
	var obChecker *actor.PID
	defer func() {
		if obChecker != nil {
			_ = as.Root.PoisonFuture(obChecker).Wait()
		}
	}()
	securityID := []uint64{
		test.securityID,
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
	sec, ok := testedSecurities[test.securityID]
	if !ok {
		t.Fatalf("security not found")
	}
	if sec.Symbol != test.symbol {
		t.Fatalf("was expecting symbol %s, got %s", test.symbol, sec.Symbol)
	}
	if sec.SecurityType != test.securityType {
		t.Fatalf("was expecting %s type, got %s", test.securityType, sec.SecurityType)
	}
	if sec.Exchange.Name != test.exchange.Name {
		t.Fatalf("was expecting %s exchange, got %s", test.exchange.Name, sec.Exchange.Name)
	}
	if sec.Underlying.ID != test.baseCurrency.ID {
		t.Fatalf("was expecting %s base, got %s", test.baseCurrency.Symbol, sec.Underlying.Symbol)
	}
	if sec.QuoteCurrency.ID != test.quoteCurrency.ID {
		t.Fatalf("was expecting %s quote, got %s", test.quoteCurrency.Symbol, sec.QuoteCurrency.Symbol)
	}
	if sec.IsInverse != test.isInverse {
		t.Fatalf("was expecting different inverse")
	}
	if sec.Status != test.status {
		t.Fatal("was expecting enabled security")
	}
	if sec.MinPriceIncrement != nil && math.Abs(sec.MinPriceIncrement.Value-test.minPriceIncrement) > 0.000001 {
		t.Fatalf("was expecting %g min price increment, got %g", test.minPriceIncrement, sec.MinPriceIncrement.Value)
	}
	if sec.RoundLot != nil && math.Abs(sec.RoundLot.Value-test.roundLot) > 0.0000000001 {
		t.Fatalf("was expecting %g round lot increment, got %g", test.roundLot, sec.RoundLot.Value)
	}
	if (sec.MaturityDate != nil) != test.hasMaturityDate {
		t.Fatalf("was expecting different maturity date")
	}

	obChecker = as.Root.Spawn(actor.PropsFromProducer(NewOBCheckerProducer(sec)))
	time.Sleep(20 * time.Second)
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
