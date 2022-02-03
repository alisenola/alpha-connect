package v3_test

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
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
	as, executor, cleaner := tests.StartExecutor(t, &constants.UNISWAPV3, nil)
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
			SecurityID: &types.UInt64Value{Value: sec.SecurityID},
			Exchange:   sec.Exchange,
			Symbol:     &types.StringValue{Value: sec.Symbol},
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
		v3.ProcessUpdate(pool, ev)
	}

	graph := graphql.NewClient(uniswap.GRAPHQL_URL, nil)
	query, variables := uniswap.GetPoolSnapshotQuery(graphql.ID(sec.Symbol), graphql.Int(12376729+2000), graphql.ID(""))
	err = graph.Query(context.Background(), &query, variables)
	if err != nil {
		t.Fatal(err)
	}
	snapshot := pool.GetSnapshot()
	if query.Pool.FeeGrowthGlobal0X128.Cmp(snapshot.FeeGrowthGlobal0X128) != 0 {
		t.Fatalf(
			"different fee growth global 0: got %s for replay and %s for graph pool, with compare %d",
			snapshot.FeeGrowthGlobal0X128.String(),
			query.Pool.FeeGrowthGlobal0X128.String(),
			query.Pool.FeeGrowthGlobal0X128.Cmp(snapshot.FeeGrowthGlobal0X128),
		)
	}
	if query.Pool.FeeGrowthGlobal1X128.Cmp(snapshot.FeeGrowthGlobal1X128) != 0 {
		t.Fatalf(
			"different fee growth global 1: got %s for replay and %s for graph pool",
			snapshot.FeeGrowthGlobal1X128.String(),
			query.Pool.FeeGrowthGlobal1X128.String(),
		)
	}
	if query.Pool.Liquidity.Cmp(snapshot.Liquidity) != 0 {
		t.Fatalf(
			"different liquidity: got %s for replay and %s for graph pool",
			snapshot.Liquidity.String(),
			query.Pool.Liquidity.String(),
		)
	}
	if query.Pool.SqrtPrice.Cmp(snapshot.SqrtPriceX96.Add(snapshot.SqrtPriceX96, big.NewInt(int64(delta)))) != 0 {
		t.Fatalf(
			"different sqrtPriceX96: got %s for replay and %s for graph pool",
			snapshot.SqrtPriceX96.String(),
			query.Pool.SqrtPrice.String(),
		)
	}

	for k := range snapshot.Ticks {
		fmt.Println("tick", k)
		tick := pool.GetTickValue(k)
		fmt.Println("FeeGrowthOutside0X128", tick.FeeGrowthOutside0X128)
		fmt.Println("FeeGrowthOutside1X128", tick.FeeGrowthOutside1X128)
		fmt.Println("Liquidity Net", tick.LiquidityNet)
		fmt.Println("Liquidity Gross", tick.LiquidityGross)
	}
}
