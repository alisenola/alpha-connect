package tests

import (
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
)

var ExPubTests = []ExPubTest{
	{
		Instrument: &models.Instrument{
			Exchange: &constants.BINANCE,
			Symbol:   &types.StringValue{Value: "BTCUSDT"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: false,
		MarketStatisticsRequest:       false,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.BITFINEX,
			Symbol:   &types.StringValue{Value: "btcusd"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: false,
		MarketStatisticsRequest:       false,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.BITSTAMP,
			Symbol:   &types.StringValue{Value: "btcusd"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: false,
		MarketStatisticsRequest:       false,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.FBINANCE,
			Symbol:   &types.StringValue{Value: "BTCUSDT"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: false,
		MarketStatisticsRequest:       true,
	},

	{
		Instrument: &models.Instrument{
			Exchange: &constants.FTXUS,
			Symbol:   &types.StringValue{Value: "BTC/USDT"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: false,
		MarketStatisticsRequest:       false,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.HUOBI,
			Symbol:   &types.StringValue{Value: "btcusdt"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: false,
		MarketStatisticsRequest:       false,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.GEMINI,
			Symbol:   &types.StringValue{Value: "btcusd"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: false,
		MarketStatisticsRequest:       false,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.HITBTC,
			Symbol:   &types.StringValue{Value: "BTCUSD"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: false,
		MarketStatisticsRequest:       false,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.KRAKEN,
			Symbol:   &types.StringValue{Value: "XBT/USD"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: false,
		MarketStatisticsRequest:       false,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.OKEX,
			Symbol:   &types.StringValue{Value: "BTC-USDT"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: false,
		MarketStatisticsRequest:       false,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.BITHUMBG,
			Symbol:   &types.StringValue{Value: "BTC-USDT"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: false,
		MarketStatisticsRequest:       false,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.UPBIT,
			Symbol:   &types.StringValue{Value: "KRW-BTC"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: false,
		MarketStatisticsRequest:       false,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.FBINANCE,
			Symbol:   &types.StringValue{Value: "BTCUSDT"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: true,
		MarketStatisticsRequest:       true,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.CRYPTOFACILITIES,
			Symbol:   &types.StringValue{Value: "pi_xbtusd"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: true,
		MarketStatisticsRequest:       true,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.BITMEX,
			Symbol:   &types.StringValue{Value: "XBTUSD"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: true,
		MarketStatisticsRequest:       true,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.FTX,
			Symbol:   &types.StringValue{Value: "BTC-PERP"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: true,
		MarketStatisticsRequest:       true,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.OKEXP,
			Symbol:   &types.StringValue{Value: "BTC-USDT-SWAP"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: true,
		MarketStatisticsRequest:       true,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.DERIBIT,
			Symbol:   &types.StringValue{Value: "BTC-PERPETUAL"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: true,
		MarketStatisticsRequest:       true,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.HUOBIP,
			Symbol:   &types.StringValue{Value: "BTC-USD"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: true,
		MarketStatisticsRequest:       true,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.BYBITI,
			Symbol:   &types.StringValue{Value: "BTCUSD"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: true,
		MarketStatisticsRequest:       true,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.BYBITL,
			Symbol:   &types.StringValue{Value: "BTCUSDT"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: true,
		MarketStatisticsRequest:       true,
	},
	{
		Instrument: &models.Instrument{
			Exchange: &constants.DYDX,
			Symbol:   &types.StringValue{Value: "BTC-USD"},
		},
		SecurityListRequest:           true,
		MarketDataRequest:             true,
		HistoricalLiquidationsRequest: true,
		MarketStatisticsRequest:       true,
	},
}

func TestAllExPub(t *testing.T) {
	t.Parallel()
	for _, tc := range ExPubTests {
		tc := tc
		t.Run(tc.Instrument.Exchange.Name, func(t *testing.T) {
			t.Parallel()
			ExPub(t, tc)
		})
	}
}
