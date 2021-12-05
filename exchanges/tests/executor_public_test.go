package tests

import (
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"reflect"
	"testing"
	"time"
)

type ExPubTest struct {
	instrument                    *models.Instrument
	securityListRequest           bool
	historicalLiquidationsRequest bool
	marketStatisticsRequest       bool
	marketDataRequest             bool
}

var ExPubTests = []ExPubTest{
	{
		instrument: &models.Instrument{
			Exchange: &constants.BINANCE,
			Symbol:   &types.StringValue{Value: "BTCUSDT"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: false,
		marketStatisticsRequest:       false,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.BITFINEX,
			Symbol:   &types.StringValue{Value: "btcusd"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: false,
		marketStatisticsRequest:       false,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.BITSTAMP,
			Symbol:   &types.StringValue{Value: "btcusd"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: false,
		marketStatisticsRequest:       false,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.FBINANCE,
			Symbol:   &types.StringValue{Value: "BTCUSDT"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: false,
		marketStatisticsRequest:       true,
	},

	{
		instrument: &models.Instrument{
			Exchange: &constants.FTXUS,
			Symbol:   &types.StringValue{Value: "BTC/USDT"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: false,
		marketStatisticsRequest:       false,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.HUOBI,
			Symbol:   &types.StringValue{Value: "btcusdt"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: false,
		marketStatisticsRequest:       false,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.GEMINI,
			Symbol:   &types.StringValue{Value: "btcusd"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: false,
		marketStatisticsRequest:       false,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.HITBTC,
			Symbol:   &types.StringValue{Value: "BTCUSD"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: false,
		marketStatisticsRequest:       false,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.KRAKEN,
			Symbol:   &types.StringValue{Value: "XBT/USD"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: false,
		marketStatisticsRequest:       false,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.OKEX,
			Symbol:   &types.StringValue{Value: "BTC-USDT"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: false,
		marketStatisticsRequest:       false,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.BITHUMBG,
			Symbol:   &types.StringValue{Value: "BTC-USDT"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: false,
		marketStatisticsRequest:       false,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.UPBIT,
			Symbol:   &types.StringValue{Value: "KRW-BTC"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: false,
		marketStatisticsRequest:       false,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.FBINANCE,
			Symbol:   &types.StringValue{Value: "BTCUSDT"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: true,
		marketStatisticsRequest:       true,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.CRYPTOFACILITIES,
			Symbol:   &types.StringValue{Value: "pi_xbtusd"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: true,
		marketStatisticsRequest:       true,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.BITMEX,
			Symbol:   &types.StringValue{Value: "XBTUSD"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: true,
		marketStatisticsRequest:       true,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.FTX,
			Symbol:   &types.StringValue{Value: "BTC-PERP"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: true,
		marketStatisticsRequest:       true,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.OKEXP,
			Symbol:   &types.StringValue{Value: "BTC-USDT-SWAP"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: true,
		marketStatisticsRequest:       true,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.DERIBIT,
			Symbol:   &types.StringValue{Value: "BTC-PERPETUAL"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: true,
		marketStatisticsRequest:       true,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.HUOBIP,
			Symbol:   &types.StringValue{Value: "BTC-USD"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: true,
		marketStatisticsRequest:       true,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.BYBITI,
			Symbol:   &types.StringValue{Value: "BTCUSD"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: true,
		marketStatisticsRequest:       true,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.BYBITL,
			Symbol:   &types.StringValue{Value: "BTCUSDT"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: true,
		marketStatisticsRequest:       true,
	},
	{
		instrument: &models.Instrument{
			Exchange: &constants.DYDX,
			Symbol:   &types.StringValue{Value: "BTC-USD"},
		},
		securityListRequest:           true,
		marketDataRequest:             true,
		historicalLiquidationsRequest: true,
		marketStatisticsRequest:       true,
	},
}

func TestAllExPub(t *testing.T) {
	t.Parallel()
	for _, tc := range ExPubTests {
		tc := tc
		t.Run(tc.instrument.Exchange.Name, func(t *testing.T) {
			t.Parallel()
			ExPub(t, tc)
		})
	}
}

func ExPub(t *testing.T, tc ExPubTest) {
	as, executor, cleaner := StartExecutor(t, tc.instrument.Exchange, nil)
	defer cleaner()

	if tc.securityListRequest {
		t.Run("SecurityListRequest", func(t *testing.T) {
			res, err := as.Root.RequestFuture(executor, &messages.SecurityListRequest{
				RequestID: 0,
			}, 10*time.Second).Result()
			if err != nil {
				t.Fatal(err)
			}
			v, ok := res.(*messages.SecurityList)
			if !ok {
				t.Fatalf("was expecting *messages.SecurityList, got %s", reflect.TypeOf(res).String())
			}
			if !v.Success {
				t.Fatalf("was expecting success, go %s", v.RejectionReason.String())
			}
		})
	}

	if tc.marketDataRequest {
		t.Run("MarketDataRequest", func(t *testing.T) {
			res, err := as.Root.RequestFuture(executor, &messages.MarketDataRequest{
				RequestID:   0,
				Instrument:  tc.instrument,
				Aggregation: models.L2,
			}, 10*time.Second).Result()
			if err != nil {
				t.Fatal(err)
			}
			v, ok := res.(*messages.MarketDataResponse)
			if !ok {
				t.Fatalf("was expecting *messages.SecurityList, got %s", reflect.TypeOf(res).String())
			}
			if !v.Success {
				t.Fatalf("was expecting success, go %s", v.RejectionReason.String())
			}
		})
	}

	if tc.marketStatisticsRequest {
		t.Run("MarketStatisticsRequest", func(t *testing.T) {
			res, err := as.Root.RequestFuture(executor, &messages.MarketStatisticsRequest{
				RequestID:  0,
				Instrument: tc.instrument,
				Statistics: []models.StatType{models.OpenInterest},
			}, 10*time.Second).Result()
			if err != nil {
				t.Fatal(err)
			}
			v, ok := res.(*messages.MarketStatisticsResponse)
			if !ok {
				t.Fatalf("was expecting *messages.MarketStatisticsResponse, got %s", reflect.TypeOf(res).String())
			}
			if !v.Success {
				t.Fatalf("was expecting success, go %s", v.RejectionReason.String())
			}
		})
	}
}
