package v3_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/enum"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/go-graphql-client"
	"gitlab.com/alphaticks/gorderbook"
	mod "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	"gitlab.com/alphaticks/xchanger/constants"
	uniswap "gitlab.com/alphaticks/xchanger/exchanges/uniswap/V3"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
)

func TestMarketData(t *testing.T) {
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

	client := graphql.NewClient(uniswap.GRAPHQL_URL, nil)
	query, variables := uniswap.GetPoolSnapshotQuery(graphql.ID(sec.Symbol), graphql.Int(12976959), graphql.ID(""))
	err = client.Query(context.Background(), &query, variables)
	if err != nil {
		t.Fatal()
	}
	feeTier := query.Pool.FeeTier
	ticks := make([]*mod.UPV3Tick, 0)
	for len(query.Pool.Ticks) == 1000 {
		fmt.Println("last tick is", query.Pool.Ticks[999])
		for _, t := range query.Pool.Ticks {
			ticks = append(
				ticks,
				&mod.UPV3Tick{
					FeeGrowthOutside0X128: t.FeeGrowthOutside0X128.Bytes(),
					FeeGrowthOutside1X128: t.FeeGrowthOutside0X128.Bytes(),
					LiquidityNet:          t.LiquidityNet.Bytes(),
					LiquidityGross:        t.LiquidityGross.Bytes(),
				},
			)
		}
		query, variables = uniswap.GetPoolSnapshotQuery(graphql.ID(sec.Symbol), graphql.Int(12976959), graphql.ID(query.Pool.Ticks[999].Id))
		err = client.Query(context.Background(), &query, variables)
		if err != nil {
			t.Fatal()
		}
	}
	for _, t := range query.Pool.Ticks {
		ticks = append(
			ticks,
			&mod.UPV3Tick{
				FeeGrowthOutside0X128: t.FeeGrowthOutside0X128.Bytes(),
				FeeGrowthOutside1X128: t.FeeGrowthOutside0X128.Bytes(),
				LiquidityNet:          t.LiquidityNet.Bytes(),
				LiquidityGross:        t.LiquidityGross.Bytes(),
			},
		)
	}
	fmt.Println("GOT", len(ticks), "ticks")
	fmt.Println("Fee tier", feeTier)

	pool := gorderbook.NewUnipoolV3(feeTier)
	pool.Initialize(query.Pool.SqrtPrice, feeTier)
	pool.Sync(ticks, query.Pool.Liquidity.Bytes(), make([]byte, 0), make([]byte, 0), query.Pool.FeeGrowthGlobal0X128.Bytes(), query.Pool.FeeGrowthGlobal1X128.Bytes(), make([]byte, 0), make([]byte, 0), feeTier, make([]*mod.UPV3Position, 0))

	test := tests.MDTest{
		IgnoreSizeResidue: true,
		SecurityID:        3923210566889873515,
		Symbol:            "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640",
		SecurityType:      enum.SecurityType_CRYPTO_AMM,
		Exchange:          constants.UNISWAPV3,
		QuoteCurrency: xchangerModels.Asset{
			Symbol: "WETH",
			Name:   "wrapped-ether",
			ID:     0,
		},
		BaseCurrency: xchangerModels.Asset{
			Symbol: "USDC",
			Name:   "usdc",
			ID:     1,
		},
		HasMaturityDate:   false,
		IsInverse:         false,
		Status:            models.Trading,
		MinPriceIncrement: 10,
	}

	checker := as.Root.Spawn(actor.PropsFromProducer(tests.NewPoolV3CheckerProducer(sec, test)))
	defer as.Root.Poison(checker)

	fmt.Println("Pause")
	time.Sleep(200 * time.Second)
	resp, err := as.Root.RequestFuture(checker, &tests.GetPool{}, 20*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	p := resp.(*tests.GetPool).Pool.GetSnapshot()
	fmt.Println("Pool Fee0", p.FeeGrowthGlobal0X128.String())
	fmt.Println("Pool Fee1", p.FeeGrowthGlobal1X128.String())
	fmt.Println("Pool Liquidity", p.Liquidity.String())
	fmt.Println("Pool SqrtPrice", p.SqrtPriceX96.String())

}
