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
	delta := 0
	for _, ev := range updates.Events {
		if i := ev.Initialize; i != nil {
			sqrt := big.NewInt(1).SetBytes(i.SqrtPriceX96)
			pool.Initialize(
				sqrt,
				i.Tick,
			)
		}
		if m := ev.Mint; m != nil {
			var own [20]byte
			copy(own[:], m.Owner)
			amount := big.NewInt(1).SetBytes(m.Amount)
			amount0 := big.NewInt(1).SetBytes(m.Amount0)
			amount1 := big.NewInt(1).SetBytes(m.Amount1)
			pool.Mint(
				own,
				m.TickLower,
				m.TickUpper,
				amount,
				amount0,
				amount1,
			)
		}
		if b := ev.Burn; b != nil {
			var own [20]byte
			copy(own[:], b.Owner)
			amount := big.NewInt(1).SetBytes(b.Amount)
			amount0 := big.NewInt(1).SetBytes(b.Amount0)
			amount1 := big.NewInt(1).SetBytes(b.Amount1)
			pool.Burn(
				own,
				b.TickLower,
				b.TickUpper,
				amount,
				amount0,
				amount1,
			)
		}
		if s := ev.Swap; s != nil {
			amount0 := big.NewInt(1).SetBytes(s.Amount0)
			amount1 := big.NewInt(1).SetBytes(s.Amount1)
			sqrt := big.NewInt(1).SetBytes(s.SqrtPriceX96)
			sqrtPrev := pool.GetSqrtPrice()
			if sqrt.Cmp(sqrtPrev) > 0 {
				sqrt.Add(sqrt, big.NewInt(1))
				delta = -1
				amount0.Neg(amount0)
			} else {
				sqrt.Sub(sqrt, big.NewInt(1))
				delta = 1
				amount1.Neg(amount1)
			}
			pool.Swap(
				amount0,
				amount1,
				sqrt,
				s.Tick,
			)
		}
		if c := ev.Collect; c != nil {
			var own [20]byte
			copy(own[:], c.Owner)
			amount0 := big.NewInt(1).SetBytes(c.AmountRequested0)
			amount1 := big.NewInt(1).SetBytes(c.AmountRequested1)
			pool.Collect(
				own,
				c.TickLower,
				c.TickUpper,
				amount0,
				amount1,
			)
		}
		if f := ev.Flash; f != nil {
			amount0 := big.NewInt(1).SetBytes(f.Amount0)
			amount1 := big.NewInt(1).SetBytes(f.Amount1)
			pool.Flash(
				amount0,
				amount1,
			)
		}
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
}
