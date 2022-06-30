package tests

import (
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/config"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/tests"
	"gitlab.com/alphaticks/xchanger/constants"
	"reflect"
	"testing"
	"time"
)

type ExPubTest struct {
	Instrument                    *models.Instrument
	SecurityListRequest           bool
	HistoricalLiquidationsRequest bool
	OpenInterestRequest           bool
	MarketDataRequest             bool
	MarketStatisticsRequest       bool
}

func ExPub(t *testing.T, tc ExPubTest) {
	C := &config.Config{
		Exchanges: []string{tc.Instrument.Exchange.Name},
	}
	as, executor, cleaner := tests.StartExecutor(t, C)
	defer cleaner()

	if tc.SecurityListRequest {
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
			for _, s := range v.Securities {
				fmt.Println(s)
			}
			if !v.Success {
				t.Fatalf("was expecting success, go %s", v.RejectionReason.String())
			}
		})
	}

	if tc.MarketDataRequest {
		t.Run("MarketDataRequest", func(t *testing.T) {
			res, err := as.Root.RequestFuture(executor, &messages.MarketDataRequest{
				RequestID:   0,
				Instrument:  tc.Instrument,
				Aggregation: models.OrderBookAggregation_L2,
			}, 20*time.Second).Result()
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

	if tc.MarketStatisticsRequest {
		t.Run("MarketStatisticsRequest", func(t *testing.T) {
			res, err := as.Root.RequestFuture(executor, &messages.MarketStatisticsRequest{
				RequestID:  0,
				Instrument: &models.Instrument{Exchange: constants.BYBITL},
				Statistics: []models.StatType{models.StatType_MarkPrice},
			}, 20*time.Second).Result()
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

	if tc.OpenInterestRequest {
		t.Run("OpenInterestRequest", func(t *testing.T) {
			res, err := as.Root.RequestFuture(executor, &messages.MarketStatisticsRequest{
				RequestID:  0,
				Instrument: tc.Instrument,
				Statistics: []models.StatType{models.StatType_OpenInterest},
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
			hasStat := false
			for _, stat := range v.Statistics {
				if stat.StatType == models.StatType_OpenInterest {
					fmt.Println(stat.Value)
					hasStat = true
				}
			}
			if !hasStat {
				t.Fatal("open interest not found")
			}
		})
	}

	if tc.HistoricalLiquidationsRequest {
		t.Run("HistoricalLiquidationsRequest", func(t *testing.T) {
			res, err := as.Root.RequestFuture(executor, &messages.HistoricalLiquidationsRequest{
				RequestID:  0,
				Instrument: tc.Instrument,
				From:       nil,
				To:         nil,
			}, 10*time.Second).Result()
			if err != nil {
				t.Fatal(err)
			}
			v, ok := res.(*messages.HistoricalLiquidationsResponse)
			if !ok {
				t.Fatalf("was expecting *messages.HistoricalLiquidationsResponse, got %s", reflect.TypeOf(res).String())
			}
			if !v.Success {
				t.Fatalf("was expecting success, go %s", v.RejectionReason.String())
			}
		})
	}
}
