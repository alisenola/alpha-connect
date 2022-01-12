package v3

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/enum"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/jobs"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/go-graphql-client"
	gorderbook "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	"gitlab.com/alphaticks/xchanger/constants"
	uniswap "gitlab.com/alphaticks/xchanger/exchanges/uniswap/V3"
	xutils "gitlab.com/alphaticks/xchanger/utils"
)

type QueryRunner struct {
	pid *actor.PID
}

type Executor struct {
	extypes.BaseExecutor
	securities   []*models.Security
	queryRunners []*QueryRunner
	dialerPool   *xutils.DialerPool
	logger       *log.Logger
}

func NewExecutor(dialerPool *xutils.DialerPool) actor.Actor {
	return &Executor{
		queryRunners: nil,
		logger:       nil,
		dialerPool:   dialerPool,
	}
}

func (state *Executor) getQueryRunner() *QueryRunner {
	sort.Slice(state.queryRunners, func(i, j int) bool {
		return rand.Uint64()%2 == 0
	})

	return state.queryRunners[0]
}

func (state *Executor) Receive(context actor.Context) {
	fmt.Println("received")
	extypes.ReceiveExecutor(state, context)
}

func (state *Executor) GetLogger() *log.Logger {
	return state.logger
}

func (state *Executor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	dialers := state.dialerPool.GetDialers()
	for _, dialer := range dialers {
		httpClient := &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 1024,
				TLSHandshakeTimeout: 10 * time.Second,
				DialContext:         dialer.DialContext,
			},
			Timeout: 10 * time.Second,
		}
		uniClient := graphql.NewClient(uniswap.GRAPHQL_URL, httpClient)
		props := actor.PropsFromProducer(func() actor.Actor {
			return jobs.NewGraphQuery(uniClient)
		})
		state.queryRunners = append(state.queryRunners, &QueryRunner{
			pid: context.Spawn(props),
		})
	}
	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {

	var securities []*models.Security

	query, variables := uniswap.GetPoolDefinitionsQuery(graphql.ID(""))
	qr := state.getQueryRunner()
	if qr == nil {
		return fmt.Errorf("rate limited")
	}

	future := context.RequestFuture(qr.pid, &jobs.PerformGraphQueryRequest{Query: &query, Variables: variables}, 10*time.Second)
	res, err := future.Result()
	if err != nil {
		return fmt.Errorf("error updating security list: %v", err)
	}
	gqr := res.(*jobs.PerformGraphQueryResponse)
	if gqr.Error != nil {
		return fmt.Errorf("error updating security list: %v", gqr.Error)
	}
	done := false
	for !done {
		for _, pool := range query.Pools {
			baseCurrency, ok := constants.GetAssetBySymbol(pool.Token0.Symbol)
			if !ok {
				state.logger.Info("unknown symbol " + pool.Token0.Symbol)
				continue
			}
			quoteCurrency, ok := constants.GetAssetBySymbol(string(pool.Token1.Symbol))
			if !ok {
				state.logger.Info("unknown symbol " + pool.Token1.Symbol)
				continue
			}
			tickSpacing, err := pool.GetTickSpacing()
			if err != nil {
				continue
			}
			security := models.Security{}
			security.Symbol = pool.Id
			security.Underlying = baseCurrency
			security.QuoteCurrency = quoteCurrency
			security.Status = models.Trading
			security.Exchange = &constants.UNISWAPV3
			security.IsInverse = false
			security.SecurityType = enum.SecurityType_CRYPTO_AMM
			security.SecuritySubType = &types.StringValue{Value: enum.SecuritySubType_UNIPOOLV3}
			security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
			security.MinPriceIncrement = &types.DoubleValue{Value: float64(tickSpacing)} // TODO in bps ?
			security.RoundLot = nil                                                      // TODO Token precision ?
			security.TakerFee = nil                                                      // TODO pool fees
			securities = append(securities, &security)
		}
		if len(query.Pools) != 1000 {
			done = true
			continue
		}
		nextID := query.Pools[len(query.Pools)-1].Id
		query, variables = uniswap.GetPoolDefinitionsQuery(graphql.ID(nextID))
		qr = state.getQueryRunner()
		if qr == nil {
			return fmt.Errorf("rate limited")
		}

		future = context.RequestFuture(qr.pid, &jobs.PerformGraphQueryRequest{Query: &query, Variables: variables}, 10*time.Second)
		res, err = future.Result()
		if err != nil {
			return fmt.Errorf("error updating security list: %v", err)
		}
		gqr = res.(*jobs.PerformGraphQueryResponse)
		if gqr.Error != nil {
			return fmt.Errorf("error updating security list: %v", gqr.Error)
		}
	}

	state.securities = securities

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: state.securities})

	return nil
}

func (state *Executor) OnSecurityListRequest(context actor.Context) error {
	msg := context.Message().(*messages.SecurityListRequest)
	context.Respond(&messages.SecurityList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: state.securities})

	return nil
}

func (state *Executor) OnUnipoolV3DataRequest(context actor.Context) error {
	msg := context.Message().(*messages.UnipoolV3DataRequest)
	response := &messages.UnipoolV3DataResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	if msg.Subscribe {
		response.RejectionReason = messages.UnsupportedSubscription
		context.Respond(response)
		return nil
	}
	if msg.Instrument == nil || msg.Instrument.Symbol == nil {
		response.RejectionReason = messages.MissingInstrument
		context.Respond(response)
		return nil
	}
	symbol := msg.Instrument.Symbol.Value
	// Symbol is pool id
	query, variables := uniswap.GetPoolSnapshotQuery(graphql.ID(symbol), graphql.Int(0), graphql.ID(""))

	qr := state.getQueryRunner()
	if qr == nil {
		return fmt.Errorf("rate limited")
	}

	future := context.RequestFuture(qr.pid, &jobs.PerformGraphQueryRequest{Query: &query, Variables: variables}, 10*time.Second)
	res, err := future.Result()
	if err != nil {
		state.logger.Warn("error getting pool snapshot", log.Error(err))
		response.RejectionReason = messages.ExchangeAPIError
		context.Respond(response)
		return nil
	}
	qresp := res.(*jobs.PerformGraphQueryResponse)
	if qresp.Error != nil {
		state.logger.Warn("error getting pool snapshot", log.Error(qresp.Error))
		response.RejectionReason = messages.ExchangeAPIError
		context.Respond(response)
		return nil
	}
	done := false
	var t []*gorderbook.UPV3Tick
	for !done {
		// Store all ticks in gorderbook.UPV3Tick structures
		for _, tick := range query.Pool.Ticks {
			t = append(t, &gorderbook.UPV3Tick{
				LiquidityNet:          tick.LiquidityNet.Bytes(),
				LiquidityGross:        tick.LiquidityGross.Bytes(),
				FeeGrowthOutside0X128: tick.FeeGrowthOutside0X128.Bytes(),
				FeeGrowthOutside1X128: tick.FeeGrowthOutside1X128.Bytes(),
			})
		}

		if len(query.Pool.Ticks) != 1000 {
			done = true
			continue
		}
		nextID := query.Pool.Ticks[len(query.Pool.Ticks)-1]
		query, variables = uniswap.GetPoolSnapshotQuery(graphql.ID(symbol), graphql.Int(0), graphql.ID(nextID))

		qr = state.getQueryRunner()
		if qr == nil {
			return fmt.Errorf("rate limited")
		}

		future = context.RequestFuture(qr.pid, &jobs.PerformGraphQueryRequest{Query: &query, Variables: variables}, 10*time.Second)
		res, err = future.Result()
		if err != nil {
			return fmt.Errorf("error getting pool snapshot %v", err)
		}
		qresp = res.(*jobs.PerformGraphQueryResponse)
		if qresp.Error != nil {
			return fmt.Errorf("error getting pool snapshot %v", err)
		}
	}

	qp, vp := uniswap.GetPositionsQuery(graphql.ID(symbol), graphql.Int(0), graphql.ID(""))
	qrun := state.getQueryRunner()
	if qrun == nil {
		return fmt.Errorf("rate limited")
	}

	f := context.RequestFuture(qrun.pid, &jobs.PerformGraphQueryRequest{Query: qp, Variables: vp}, 10*time.Second)
	resp, err := f.Result()
	if err != nil {
		return fmt.Errorf("error getting positions %v", err)
	}
	qrespP := resp.(*jobs.PerformGraphQueryResponse)
	if qrespP.Error != nil {
		return fmt.Errorf("error getting pool snapshot %v", err)
	}
	done = false
	var pos []*gorderbook.UPV3Position
	for !done {
		// Store all the positions in gorderbook.UPV3Position structures
		for _, p := range qp.Positions {
			idByte, ownerByte, err := p.StringToBytes()
			if err != nil {
				continue
			}
			pos = append(pos, &gorderbook.UPV3Position{
				ID:        idByte,
				Owner:     ownerByte,
				TickLower: p.TickLower.TickIdx,
				TickUpper: p.TickUpper.TickIdx,
			})
		}
		if len(qp.Positions) != 1000 {
			done = true
			continue
		}
		nextID := qp.Positions[len(qp.Positions)-1]
		qp, vp = uniswap.GetPositionsQuery(graphql.ID(symbol), graphql.Int(0), graphql.ID(nextID))
		qrun = state.getQueryRunner()
		if qrun == nil {
			return fmt.Errorf("rate limited")
		}

		f = context.RequestFuture(qrun.pid, &jobs.PerformGraphQueryRequest{Query: qp, Variables: vp}, 10*time.Second)
		resp, err = f.Result()
		if err != nil {
			return fmt.Errorf("error getting positions %v", err)
		}
		qrespP = resp.(*jobs.PerformGraphQueryResponse)
		if qrespP.Error != nil {
			return fmt.Errorf("error getting pool snapshot %v", err)
		}
	}

	response.Snapshot = &models.UPV3Snapshot{
		Ticks:                 t,
		Positions:             pos,
		Liquidity:             query.Pool.Liquidity.Bytes(),
		SqrtPrice:             query.Pool.SqrtPrice.Bytes(),
		FeeGrowthGlobal_0X128: query.Pool.FeeGrowthGlobal0X128.Bytes(),
		FeeGrowthGlobal_1X128: query.Pool.FeeGrowthGlobal1X128.Bytes(),
		Tick:                  query.Pool.Tick,
	}
	response.Success = true
	response.SeqNum = 0 // TODO
	fmt.Printf("got this %+v", response)
	context.Respond(response)

	return nil
}
