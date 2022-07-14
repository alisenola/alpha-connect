package v3

import (
	"fmt"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	xmodels "gitlab.com/alphaticks/xchanger/models"
	xutils "gitlab.com/alphaticks/xchanger/utils"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math/big"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"time"

	"gitlab.com/alphaticks/xchanger/chains/evm"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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
)

type QueryRunner struct {
	pid *actor.PID
}

type Executor struct {
	extypes.BaseExecutor
	queryRunners []*QueryRunner
	logger       *log.Logger
	executor     *actor.PID
}

func NewExecutor(dialerPool *xutils.DialerPool, registry registry.PublicRegistryClient) actor.Actor {
	e := &Executor{}
	e.DialerPool = dialerPool
	e.Registry = registry
	return e
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
		log.String("type", reflect.TypeOf(state).String()))

	dialers := state.DialerPool.GetDialers()
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

	state.executor = actor.NewPID(context.ActorSystem().Address(), "executor")

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
			token0, ok := constants.GetAssetBySymbol(pool.Token0.Symbol)
			if !ok {
				// state.logger.Info("unknown symbol " + pool.Token0.Symbol)
				continue
			}
			token1, ok := constants.GetAssetBySymbol(pool.Token1.Symbol)
			if !ok {
				// state.logger.Info("unknown symbol " + pool.Token1.Symbol)
				continue
			}

			tickSpacing, err := pool.GetTickSpacing()
			if err != nil {
				continue
			}

			var baseCurrency, quoteCurrency *xmodels.Asset
			var inverse bool
			if token1.Symbol == "USDC" || token1.Symbol == "USDT" || token1.Symbol == "DAI" || token1.Symbol == "BUSD" {
				baseCurrency = token0
				quoteCurrency = token1
				inverse = false
			} else if token0.Symbol == "USDC" || token0.Symbol == "USDT" || token0.Symbol == "DAI" || token0.Symbol == "BUSD" {
				baseCurrency = token1
				quoteCurrency = token0
				inverse = true
			} else if token1.Symbol == "WBTC" || token1.Symbol == "WETH" {
				baseCurrency = token0
				quoteCurrency = token1
				inverse = false
			} else if token0.Symbol == "WBTC" || token0.Symbol == "WETH" {
				baseCurrency = token1
				quoteCurrency = token0
				inverse = true
			} else {
				baseCurrency = token0
				quoteCurrency = token1
				inverse = false
			}

			security := models.Security{}
			security.Symbol = pool.Id
			security.Underlying = baseCurrency
			security.QuoteCurrency = quoteCurrency
			security.IsInverse = inverse
			security.Status = models.InstrumentStatus_Trading
			security.Exchange = constants.UNISWAPV3
			security.SecurityType = enum.SecurityType_CRYPTO_AMM
			security.SecuritySubType = &wrapperspb.StringValue{Value: enum.SecuritySubType_UNIPOOLV3}
			security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
			security.MinPriceIncrement = &wrapperspb.DoubleValue{Value: float64(tickSpacing)}
			security.TakerFee = &wrapperspb.DoubleValue{Value: float64(pool.FeeTier)}
			security.CreationBlock = &wrapperspb.UInt64Value{Value: pool.CreatedAtBlockNumber.Uint64()}
			security.Protocol = constants.ERC20
			security.Chain = constants.EthereumMainnet
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

	state.SyncSecurities(securities, nil)

	context.Send(context.Parent(), &messages.SecurityList{
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: securities})

	return nil
}

func (state *Executor) OnHistoricalUnipoolV3DataRequest(context actor.Context) error {
	msg := context.Message().(*messages.HistoricalUnipoolV3DataRequest)
	response := &messages.HistoricalUnipoolV3DataResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	if msg.Instrument == nil || msg.Instrument.Symbol == nil {
		response.RejectionReason = messages.RejectionReason_MissingInstrument
		context.Respond(response)
		return nil
	}

	symbol := msg.Instrument.Symbol
	uabi, err := uniswap.UniswapMetaData.GetAbi()
	if err != nil {
		state.logger.Warn("error getting uniswap ABI", log.Error(err))
		response.RejectionReason = messages.RejectionReason_ABIError
		context.Respond(response)
		return nil
	}
	query := [][]interface{}{{
		uabi.Events["Initialize"].ID,
		uabi.Events["Mint"].ID,
		uabi.Events["Swap"].ID,
		uabi.Events["Burn"].ID,
		uabi.Events["Collect"].ID,
		uabi.Events["Flash"].ID,
		uabi.Events["SetFeeProtocol"].ID,
		uabi.Events["CollectProtocol"].ID,
	}}
	topics, err := abi.MakeTopics(query...)
	if err != nil {
		state.logger.Warn("error getting topics", log.Error(err))
		response.RejectionReason = messages.RejectionReason_ABIError
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

	future := context.RequestFuture(state.executor, &jobs.PerformLogsQueryRequest{Query: fQuery}, 50*time.Second)
	context.ReenterAfter(future, func(resp interface{}, err error) {
		if err != nil {
			state.logger.Warn("error at eth rpc server", log.Error(err))
			response.RejectionReason = messages.RejectionReason_RPCError
			context.Respond(response)
			return
		}
		queryResponse := resp.(*jobs.PerformLogsQueryResponse)
		if queryResponse.Error != nil {
			state.logger.Warn("error at eth rpc server", log.Error(queryResponse.Error))
			response.RejectionReason = messages.RejectionReason_RPCError
			context.Respond(response)
			return
		}
		logs := queryResponse.Logs
		for i, l := range logs {
			switch l.Topics[0] {
			case uabi.Events["Initialize"].ID:
				event := uniswap.UniswapInitialize{}
				if err := evm.UnpackLog(uabi, &event, "Initialize", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
					response.RejectionReason = messages.RejectionReason_RPCError
					context.Respond(response)
					return
				}
				update := &models.UPV3Update{
					Initialize: &gorderbook.UPV3Initialize{
						SqrtPriceX96: event.SqrtPriceX96.Bytes(),
						Tick:         int32(event.Tick.Int64()),
					},
					Block:     l.BlockNumber,
					Timestamp: utils.SecondToTimestamp(queryResponse.Times[i]),
				}
				response.Events = append(response.Events, update)
			case uabi.Events["Mint"].ID:
				event := uniswap.UniswapMint{}
				if err := evm.UnpackLog(uabi, &event, "Mint", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
					response.RejectionReason = messages.RejectionReason_RPCError
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
					Block:     l.BlockNumber,
					Timestamp: utils.SecondToTimestamp(queryResponse.Times[i]),
				}
				response.Events = append(response.Events, update)
			case uabi.Events["Burn"].ID:
				event := uniswap.UniswapBurn{}
				if err := evm.UnpackLog(uabi, &event, "Burn", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
					response.RejectionReason = messages.RejectionReason_RPCError
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
					Block:     l.BlockNumber,
					Timestamp: utils.SecondToTimestamp(queryResponse.Times[i]),
				}
				response.Events = append(response.Events, update)
			case uabi.Events["Swap"].ID:
				event := uniswap.UniswapSwap{}
				if err := evm.UnpackLog(uabi, &event, "Swap", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
					response.RejectionReason = messages.RejectionReason_RPCError
					context.Respond(response)
					return
				}
				update := &models.UPV3Update{
					Swap: &gorderbook.UPV3Swap{
						SqrtPriceX96: event.SqrtPriceX96.Bytes(),
						Tick:         int32(event.Tick.Int64()),
						Amount0:      event.Amount0.Bytes(),
						Amount1:      event.Amount1.Bytes(),
						Liquidity:    event.Liquidity.Bytes(),
					},
					Block:     l.BlockNumber,
					Timestamp: utils.SecondToTimestamp(queryResponse.Times[i]),
				}
				response.Events = append(response.Events, update)
			case uabi.Events["Collect"].ID:
				event := uniswap.UniswapCollect{}
				if err := evm.UnpackLog(uabi, &event, "Collect", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
					response.RejectionReason = messages.RejectionReason_RPCError
					context.Respond(response)
					return
				}
				update := &models.UPV3Update{
					Collect: &gorderbook.UPV3Collect{
						Owner:            event.Recipient[:],
						TickLower:        int32(event.TickLower.Int64()),
						TickUpper:        int32(event.TickUpper.Int64()),
						AmountRequested0: event.Amount0.Bytes(),
						AmountRequested1: event.Amount1.Bytes(),
					},
					Block:     l.BlockNumber,
					Timestamp: utils.SecondToTimestamp(queryResponse.Times[i]),
				}
				response.Events = append(response.Events, update)
			case uabi.Events["Flash"].ID:
				event := uniswap.UniswapFlash{}
				if err := evm.UnpackLog(uabi, &event, "Flash", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
					response.RejectionReason = messages.RejectionReason_RPCError
					context.Respond(response)
					return
				}
				update := &models.UPV3Update{
					Flash: &gorderbook.UPV3Flash{
						Amount0: event.Amount0.Bytes(),
						Amount1: event.Amount1.Bytes(),
					},
					Block:     l.BlockNumber,
					Timestamp: utils.SecondToTimestamp(queryResponse.Times[i]),
				}
				response.Events = append(response.Events, update)
			case uabi.Events["SetFeeProtocol"].ID:
				event := uniswap.UniswapSetFeeProtocol{}
				if err := evm.UnpackLog(uabi, &event, "SetFeeProtocol", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
					response.RejectionReason = messages.RejectionReason_RPCError
					context.Respond(response)
					return
				}
				update := &models.UPV3Update{
					SetFeeProtocol: &gorderbook.UPV3SetFeeProtocol{
						FeesProtocol: uint32(event.FeeProtocol0New) + uint32(event.FeeProtocol1New)<<8,
					},
					Block:     l.BlockNumber,
					Timestamp: utils.SecondToTimestamp(queryResponse.Times[i]),
				}
				response.Events = append(response.Events, update)
			case uabi.Events["CollectProtocol"].ID:
				event := uniswap.UniswapCollectProtocol{}
				if err := evm.UnpackLog(uabi, &event, "CollectProtocol", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
					response.RejectionReason = messages.RejectionReason_RPCError
					context.Respond(response)
					return
				}
				update := &models.UPV3Update{
					CollectProtocol: &gorderbook.UPV3CollectProtocol{
						AmountRequested0: event.Amount0.Bytes(),
						AmountRequested1: event.Amount1.Bytes(),
					},
					Block:     l.BlockNumber,
					Timestamp: utils.SecondToTimestamp(queryResponse.Times[i]),
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
