package v3

import (
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"sort"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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

	ETH_CLIENT_URL := "wss://eth-mainnet.alchemyapi.io/v2/hdodrT8DMC-Ow9rd6qIOcjpZgr5_Ixdg"
	client, err := ethclient.Dial(ETH_CLIENT_URL)
	if err != nil {
		return fmt.Errorf("error while dialing eth rpc client %v", err)
	}
	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewABIQuery(client)
	})
	state.queryRunners = append(state.queryRunners, &QueryRunner{
		pid: context.Spawn(props),
	})
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
			security.TakerFee = &types.DoubleValue{Value: float64(pool.FeeTier)}         // TODO pool fees
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

func (state *Executor) OnHistoricalUnipoolV3EventRequest(context actor.Context) error {
	fmt.Println("EXECUTOR HISTORICAL DATA REQUEST")
	msg := context.Message().(*messages.HistoricalUnipoolV3EventRequest)
	response := &messages.HistoricalUnipoolV3EventResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	if msg.Instrument == nil || msg.Instrument.Symbol == nil {
		response.RejectionReason = messages.MissingInstrument
		context.Respond(response)
		return nil
	}

	symbol := msg.Instrument.Symbol
	uabi, err := uniswap.UniswapMetaData.GetAbi()
	if err != nil {
		state.logger.Warn("error getting abi", log.Error(err))
		response.RejectionReason = messages.ABIError
		context.Respond(response)
		return nil
	}
	query := append([][]interface{}{{
		uabi.Events["Initialize"].ID,
		uabi.Events["Mint"].ID,
		uabi.Events["Swap"].ID,
		uabi.Events["Burn"].ID,
		uabi.Events["Collect"].ID,
		uabi.Events["Flash"].ID,
	}})
	topics, err := abi.MakeTopics(query...)
	if err != nil {
		state.logger.Warn("error getting topics", log.Error(err))
		response.RejectionReason = messages.ABIError
		context.Respond(response)
		return nil
	}
	start := msg.Start
	end := msg.End
	fQuery := ethereum.FilterQuery{
		BlockHash: nil,
		FromBlock: big.NewInt(int64(start)),
		ToBlock:   big.NewInt(int64(end)),
		Addresses: []common.Address{common.HexToAddress(symbol.Value)},
		Topics:    topics,
	}

	qr := state.getQueryRunner()
	if qr == nil {
		return fmt.Errorf("rate limited")
	}

	future := context.RequestFuture(qr.pid, &jobs.PerformLogsQueryRequest{Query: fQuery}, 20*time.Second)
	context.AwaitFuture(future, func(resp interface{}, err error) {
		if err != nil {
			state.logger.Warn("error at eth rpc server", log.Error(err))
			response.RejectionReason = messages.EthRPCError
			context.Respond(response)
			return
		}
		queryResponse := resp.(*jobs.PerformLogsQueryResponse)
		if queryResponse.Error != nil {
			state.logger.Warn("error at eth rpc server", log.Error(err))
			response.RejectionReason = messages.EthRPCError
			context.Respond(response)
			return
		}
		//TODO Order the logs based on block and txIndex
		logs := queryResponse.Logs
		for _, l := range logs {
			switch l.Topics[0] {
			case uabi.Events["Initialize"].ID:
				event := uniswap.UniswapInitialize{}
				if err := uabi.UnpackIntoInterface(&event, "Initialize", l.Data); err != nil {
					state.logger.Warn("error unpacking the log", log.Error(err))
					response.RejectionReason = messages.EthRPCError
					context.Respond(response)
					return
				}
				update := &models.UPV3Update{
					Initialize: &gorderbook.UPV3Initialize{
						SqrtPriceX96: event.SqrtPriceX96.Bytes(),
						Tick:         int32(event.Tick.Int64()),
					},
				}
				response.Events = append(response.Events, update)
			case uabi.Events["Mint"].ID:
				event := uniswap.UniswapMint{}
				if err := uabi.UnpackIntoInterface(&event, "Mint", l.Data); err != nil {
					state.logger.Warn("error unpacking the log", log.Error(err))
					response.RejectionReason = messages.EthRPCError
					context.Respond(response)
					return
				}
				update := &models.UPV3Update{
					Mint: &gorderbook.UPV3Mint{
						Owner:     event.Owner[:],
						TickLower: int32(event.TickLower.Int64()),
						TickUpper: int32(event.TickUpper.Int64()),
						Amount:    event.Amount.Bytes(),
						Amount0:   event.Amount0.Bytes(),
						Amount1:   event.Amount1.Bytes(),
					},
				}
				response.Events = append(response.Events, update)
			case uabi.Events["Burns"].ID:
				event := uniswap.UniswapBurn{}
				if err := uabi.UnpackIntoInterface(&event, "Burn", l.Data); err != nil {
					state.logger.Warn("error unpacking the log", log.Error(err))
					response.RejectionReason = messages.EthRPCError
					context.Respond(response)
					return
				}
				update := &models.UPV3Update{
					Burn: &gorderbook.UPV3Burn{
						Owner:     event.Owner[:],
						TickLower: int32(event.TickLower.Int64()),
						TickUpper: int32(event.TickUpper.Int64()),
						Amount:    event.Amount.Bytes(),
						Amount0:   event.Amount0.Bytes(),
						Amount1:   event.Amount1.Bytes(),
					},
				}
				response.Events = append(response.Events, update)
			case uabi.Events["Swap"].ID:
				event := uniswap.UniswapSwap{}
				if err := uabi.UnpackIntoInterface(&event, "Swap", l.Data); err != nil {
					state.logger.Warn("error unpacking the log", log.Error(err))
					response.RejectionReason = messages.EthRPCError
					context.Respond(response)
					return
				}
				update := &models.UPV3Update{
					Swap: &gorderbook.UPV3Swap{
						SqrtPriceX96: event.SqrtPriceX96.Bytes(),
						Tick:         int32(event.Tick.Int64()),
						Amount0:      event.Amount0.Bytes(),
						Amount1:      event.Amount1.Bytes(),
					},
				}
				response.Events = append(response.Events, update)
			case uabi.Events["Collect"].ID:
				event := uniswap.UniswapCollect{}
				if err := uabi.UnpackIntoInterface(&event, "Collect", l.Data); err != nil {
					state.logger.Warn("error unpacking the log", log.Error(err))
					response.RejectionReason = messages.EthRPCError
					context.Respond(response)
					return
				}
				update := &models.UPV3Update{
					Collect: &gorderbook.UPV3Collect{
						Owner:            event.Owner.Bytes(),
						TickLower:        int32(event.TickLower.Int64()),
						TickUpper:        int32(event.TickUpper.Int64()),
						AmountRequested0: event.Amount0.Bytes(),
						AmountRequested1: event.Amount1.Bytes(),
					},
				}
				response.Events = append(response.Events, update)
			case uabi.Events["Flash"].ID:
				event := uniswap.UniswapFlash{}
				if err := uabi.UnpackIntoInterface(&event, "Flash", l.Data); err != nil {
					state.logger.Warn("error unpacking the log", log.Error(err))
					response.RejectionReason = messages.EthRPCError
					context.Respond(response)
					return
				}
				update := &models.UPV3Update{
					Flash: &gorderbook.UPV3Flash{
						Amount0: event.Amount0.Bytes(),
						Amount1: event.Amount1.Bytes(),
					},
				}
				response.Events = append(response.Events, update)
			}
		}
		response.Success = true
		response.SeqNum = uint64(time.Now().UnixNano())
		context.Respond(response)

	})
	return nil
}

func (state *Executor) OnUnipoolV3DataRequest(context actor.Context) error {
	fmt.Println("EXECUTOR UNIPOOLV3 DATA REQUEST")
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
		state.logger.Warn("error getting first pool snapshot", log.Error(err))
		response.RejectionReason = messages.ExchangeAPIError
		context.Respond(response)
		return nil
	}
	qresp := res.(*jobs.PerformGraphQueryResponse)
	if qresp.Error != nil {
		state.logger.Warn("error in query response", log.Error(qresp.Error))
		response.RejectionReason = messages.ExchangeAPIError
		context.Respond(response)
		return nil
	}

	var t []*gorderbook.UPV3Tick
	// Store all ticks in gorderbook.UPV3Tick structures
	for _, tick := range query.Pool.Ticks {
		t = append(t, &gorderbook.UPV3Tick{
			ID:                    tick.TickIdx,
			LiquidityNet:          tick.LiquidityNet.Bytes(),
			LiquidityGross:        tick.LiquidityGross.Bytes(),
			FeeGrowthOutside0X128: tick.FeeGrowthOutside0X128.Bytes(),
			FeeGrowthOutside1X128: tick.FeeGrowthOutside1X128.Bytes(),
		})
	}

	qp, vp := uniswap.GetPositionsQuery(graphql.ID(symbol), graphql.Int(0), graphql.ID(""))
	qrun := state.getQueryRunner()
	if qrun == nil {
		return fmt.Errorf("rate limited")
	}

	f := context.RequestFuture(qrun.pid, &jobs.PerformGraphQueryRequest{Query: &qp, Variables: vp}, 10*time.Second)
	resp, err := f.Result()
	if err != nil {
		state.logger.Warn("error getting positions", log.Error(err))
		response.RejectionReason = messages.ExchangeAPIError
		context.Respond(response)
		return nil
	}
	qrespP := resp.(*jobs.PerformGraphQueryResponse)
	if qrespP.Error != nil {
		state.logger.Warn("error in query response", log.Error(qrespP.Error))
		response.RejectionReason = messages.ExchangeAPIError
		context.Respond(response)
		return nil
	}

	var pos []*gorderbook.UPV3Position
	for _, p := range qp.Positions {
		idByte, ownerByte, err := p.StringToBytes()
		if err != nil {
			continue
		}
		//TODO update the proto in order to have the correct values
		pos = append(pos, &gorderbook.UPV3Position{
			ID:                        idByte,
			Owner:                     ownerByte,
			TickLower:                 p.TickLower.TickIdx,
			TickUpper:                 p.TickUpper.TickIdx,
			Liquidity:                 p.Liquidity.Bytes(),
			FeeGrowthInside_0LastX128: p.FeeGrowthInside0LastX128.Bytes(),
			FeeGrowthInside_1LastX128: p.FeeGrowthInside1LastX128.Bytes(),
			TokensOwed0:               p.TokensOwed0.Bytes(),
			TokensOwed1:               p.TokensOwed1.Bytes(),
		})
	}

	var timestamps = []int32{}
	if len(query.Pool.Mints) > 0 {
		timestamps = append(timestamps, query.Pool.Mints[0].Timestamp)
	}
	if len(query.Pool.Burns) > 0 {
		timestamps = append(timestamps, query.Pool.Burns[0].Timestamp)
	}
	if len(query.Pool.Swaps) > 0 {
		timestamps = append(timestamps, query.Pool.Swaps[0].Timestamp)
	}
	if len(query.Pool.Collects) > 0 {
		timestamps = append(timestamps, query.Pool.Collects[0].Timestamp)
	}
	var minTime = int32(^uint32(0) >> 1)
	for _, time := range timestamps {
		if time < minTime {
			minTime = time
		}
	}

	response.Snapshot = &models.UPV3Snapshot{
		Ticks:                   t,
		Positions:               pos,
		Liquidity:               query.Pool.Liquidity.Bytes(),
		SqrtPrice:               query.Pool.SqrtPrice.Bytes(),
		FeeGrowthGlobal_0X128:   query.Pool.FeeGrowthGlobal0X128.Bytes(),
		FeeGrowthGlobal_1X128:   query.Pool.FeeGrowthGlobal1X128.Bytes(),
		ProtocolFees_0:          query.Pool.ProtocolFees0.Bytes(),
		ProtocolFees_1:          query.Pool.ProtocolFees1.Bytes(),
		TotalValueLockedToken_0: query.Pool.TotalValueLockedToken0.Bytes(),
		TotalValueLockedToken_1: query.Pool.TotalValueLockedToken1.Bytes(),
		Tick:                    query.Pool.Tick,
		FeeTier:                 query.Pool.FeeTier,
		Timestamp:               &types.Timestamp{Seconds: int64(minTime)},
	}
	fmt.Println("TIMESTAMP", response.Snapshot.Timestamp)

	response.Success = true
	response.SeqNum = uint64(time.Now().UnixNano())
	context.Respond(response)

	return nil
}
