package v3_test

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math/big"
	"reflect"
	"testing"
	"time"

	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	v3 "gitlab.com/alphaticks/alpha-connect/exchanges/uniswap/v3"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/go-graphql-client"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/constants"
	uniswap "gitlab.com/alphaticks/xchanger/exchanges/uniswap/V3"
)

func TestPoolData(t *testing.T) {
	as, executor, cleaner := tests.StartExecutor(t, constants.UNISWAPV3, nil)
	defer cleaner()

	securityID := []uint64{
		3923210566889873515,
	}

	res, err := as.Root.RequestFuture(executor, &messages.SecurityListRequest{}, 30*time.Second).Result()
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
	fmt.Println("GOT SEC", len(securityList.Securities))

	var sec *models.Security
	for _, s := range securityList.Securities {
		for _, secID := range securityID {
			if secID == s.SecurityID {
				sec = s
			}
		}
	}
	if sec == nil {
		t.Fatal("security not found")
	}

	res, err = as.Root.RequestFuture(executor, &messages.HistoricalUnipoolV3DataRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Instrument: &models.Instrument{
			SecurityID: &wrapperspb.UInt64Value{Value: sec.SecurityID},
			Exchange:   sec.Exchange,
			Symbol:     &wrapperspb.StringValue{Value: sec.Symbol},
		},
		Start: 12376729,
		End:   12376729 + 2000,
	}, 50*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	updates, ok := res.(*messages.HistoricalUnipoolV3DataResponse)
	if !ok {
		t.Fatalf("was expecting *messages.HistoricalUnipoolV3DataResponse, got %s", reflect.TypeOf(res).String())
	}
	if !updates.Success {
		t.Fatal(updates.RejectionReason.String())
	}

	fmt.Println(updates.Events)
	pool := gorderbook.NewUnipoolV3(int32(sec.TakerFee.Value))
	delta := -1
	for _, ev := range updates.Events {
		fmt.Println(ev.Timestamp)
		if err := v3.ProcessUpdate(pool, ev); err != nil {
			t.Fatal(err)
		}
	}

	graph := graphql.NewClient(uniswap.GRAPHQL_URL, nil)
	query, variables := uniswap.GetPoolSnapshotQuery(graphql.ID(sec.Symbol), graphql.Int(12376729+2000), graphql.ID(""))
	err = graph.Query(context.Background(), &query, variables)
	if err != nil {
		t.Fatal(err)
	}
	snapshot := pool.GetSnapshot()
	liquidity := big.NewInt(0).SetBytes(snapshot.Liquidity[:])
	if query.Pool.Liquidity.Cmp(liquidity) != 0 {
		t.Fatalf(
			"different liquidity: got %s for replay and %s for graph pool",
			liquidity.String(),
			query.Pool.Liquidity.String(),
		)
	}
	sqrtPriceX96 := big.NewInt(0).SetBytes(snapshot.SqrtPriceX96[:])
	if query.Pool.SqrtPrice.Cmp(sqrtPriceX96.Add(sqrtPriceX96, big.NewInt(int64(delta)))) != 0 {
		t.Fatalf(
			"different sqrtPriceX96: got %s for replay and %s for graph pool",
			sqrtPriceX96.String(),
			query.Pool.SqrtPrice.String(),
		)
	}

	for k := range snapshot.Ticks {
		fmt.Println("tick", k)
	}
}
