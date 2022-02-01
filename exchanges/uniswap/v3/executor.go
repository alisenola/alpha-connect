package v3

import (
	"fmt"
	"math/big"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"time"

	"gitlab.com/alphaticks/xchanger/eth"

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
	securities     []*models.Security
	queryRunnerETH *QueryRunner
	queryRunners   []*QueryRunner
	dialerPool     *xutils.DialerPool
	logger         *log.Logger
}

func NewExecutor(dialerPool *xutils.DialerPool) actor.Actor {
	return &Executor{
		queryRunnerETH: nil,
		queryRunners:   nil,
		logger:         nil,
		dialerPool:     dialerPool,
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

	client, err := ethclient.Dial(ETH_CLIENT_WS)
	if err != nil {
		return fmt.Errorf("error while dialing eth rpc client %v", err)
	}
	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewABIQuery(client)
	})
	state.queryRunnerETH = &QueryRunner{
		pid: context.Spawn(props),
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
				// state.logger.Info("unknown symbol " + pool.Token0.Symbol)
				continue
			}
			quoteCurrency, ok := constants.GetAssetBySymbol(string(pool.Token1.Symbol))
			if !ok {
				// state.logger.Info("unknown symbol " + pool.Token1.Symbol)
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
	msg := context.Message().(*messages.HistoricalUnipoolV3DataRequest)
	response := &messages.HistoricalUnipoolV3DataResponse{
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
		state.logger.Warn("error getting uniswap ABI", log.Error(err))
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
		uabi.Events["SetFeeProtocol"].ID,
		uabi.Events["CollectProtocol"].ID,
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

	qr := state.queryRunnerETH

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
		logs := queryResponse.Logs
		for _, l := range logs {
			switch l.Topics[0] {
			case uabi.Events["Initialize"].ID:
				event := uniswap.UniswapInitialize{}
				if err := eth.UnpackLog(uabi, &event, "Initialize", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
					response.RejectionReason = messages.EthRPCError
					context.Respond(response)
					return
				}
				update := &models.UPV3Update{
					Initialize: &gorderbook.UPV3Initialize{
						SqrtPriceX96: event.SqrtPriceX96.Bytes(),
						Tick:         int32(event.Tick.Int64()),
					},
					Removed: l.Removed,
					Block:   l.BlockNumber,
				}
				response.Events = append(response.Events, update)
			case uabi.Events["Mint"].ID:
				event := uniswap.UniswapMint{}
				if err := eth.UnpackLog(uabi, &event, "Mint", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
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
					Removed: l.Removed,
					Block:   l.BlockNumber,
				}
				response.Events = append(response.Events, update)
			case uabi.Events["Burn"].ID:
				event := uniswap.UniswapBurn{}
				if err := eth.UnpackLog(uabi, &event, "Burn", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
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
					Removed: l.Removed,
					Block:   l.BlockNumber,
				}
				response.Events = append(response.Events, update)
			case uabi.Events["Swap"].ID:
				event := uniswap.UniswapSwap{}
				if err := eth.UnpackLog(uabi, &event, "Swap", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
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
					Removed: l.Removed,
					Block:   l.BlockNumber,
				}
				response.Events = append(response.Events, update)
			case uabi.Events["Collect"].ID:
				event := uniswap.UniswapCollect{}
				if err := eth.UnpackLog(uabi, &event, "Collect", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
					response.RejectionReason = messages.EthRPCError
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
					Removed: l.Removed,
					Block:   l.BlockNumber,
				}
				response.Events = append(response.Events, update)
			case uabi.Events["Flash"].ID:
				event := uniswap.UniswapFlash{}
				if err := eth.UnpackLog(uabi, &event, "Flash", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
					response.RejectionReason = messages.EthRPCError
					context.Respond(response)
					return
				}
				update := &models.UPV3Update{
					Flash: &gorderbook.UPV3Flash{
						Amount0: event.Amount0.Bytes(),
						Amount1: event.Amount1.Bytes(),
					},
					Removed: l.Removed,
					Block:   l.BlockNumber,
				}
				response.Events = append(response.Events, update)
			case uabi.Events["SetFeeProtocol"].ID:
				event := uniswap.UniswapSetFeeProtocol{}
				if err := eth.UnpackLog(uabi, &event, "SetFeeProtocol", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
					response.RejectionReason = messages.EthRPCError
					context.Respond(response)
					return
				}
				update := &models.UPV3Update{
					SetFeeProtocol: &gorderbook.UPV3SetFeeProtocol{
						FeesProtocol: uint32(event.FeeProtocol0New) + uint32(event.FeeProtocol1New)<<8,
					},
					Removed: l.Removed,
					Block:   l.BlockNumber,
				}
				response.Events = append(response.Events, update)
			case uabi.Events["CollectProtocol"].ID:
				event := uniswap.UniswapCollectProtocol{}
				if err := eth.UnpackLog(uabi, &event, "CollectProtocol", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
					response.RejectionReason = messages.EthRPCError
					context.Respond(response)
					return
				}
				update := &models.UPV3Update{
					CollectProtocol: &gorderbook.UPV3CollectProtocol{
						AmountRequested0: event.Amount0.Bytes(),
						AmountRequested1: event.Amount1.Bytes(),
					},
					Block:   l.BlockNumber,
					Removed: l.Removed,
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
