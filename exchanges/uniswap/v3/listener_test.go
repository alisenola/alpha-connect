package v3_test

import (
	"context"
	"fmt"
	"math/big"
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
	pool := gorderbook.NewUnipoolV3(feeTier)
	pool.Initialize(query.Pool.SqrtPrice, feeTier)
	snap := gorderbook.PoolSnapShot{
		Tick:                 query.Pool.Tick,
		FeeTier:              query.Pool.FeeTier,
		Liquidity:            query.Pool.Liquidity,
		SqrtPriceX96:         query.Pool.SqrtPrice,
		ProtocolFees0:        big.NewInt(0),
		ProtocolFees1:        big.NewInt(0),
		FeeGrowthGlobal0X128: query.Pool.FeeGrowthGlobal0X128,
		FeeGrowthGlobal1X128: query.Pool.FeeGrowthGlobal1X128,
		Tvl0:                 big.NewInt(0),
		Tvl1:                 big.NewInt(0),
		FeeProtocol:          0,
		Ticks:                make(map[int32]*gorderbook.UnipoolV3Tick),
		Positions:            make(map[[32]byte]*gorderbook.UnipoolV3Position),
	}
	ticks := make(map[int32]*gorderbook.UnipoolV3Tick, 0)
	for len(query.Pool.Ticks) == 1000 {
		fmt.Println("last tick is", query.Pool.Ticks[999])
		for _, t := range query.Pool.Ticks {
			snap.SetTick(t.TickIdx, t.FeeGrowthOutside0X128, t.FeeGrowthOutside1X128, t.LiquidityNet, t.LiquidityGross)
		}
		query, variables = uniswap.GetPoolSnapshotQuery(graphql.ID(sec.Symbol), graphql.Int(12976959), graphql.ID(query.Pool.Ticks[999].Id))
		err = client.Query(context.Background(), &query, variables)
		if err != nil {
			t.Fatal()
		}
	}
	for _, t := range query.Pool.Ticks {
		snap.SetTick(t.TickIdx, t.FeeGrowthOutside0X128, t.FeeGrowthOutside1X128, t.LiquidityNet, t.LiquidityGross)
	}
	fmt.Println("GOT", len(ticks), "ticks")
	fmt.Println("Fee tier", feeTier)

	pool.Sync(&snap)

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
