package v3

import (
	"container/list"
	goContext "context"
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/xchanger/chains/evm"
	"gitlab.com/alphaticks/xchanger/constants"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math/big"
	"reflect"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	gorderbook "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	uniswap "gitlab.com/alphaticks/xchanger/exchanges/uniswap/V3"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
)

type checkSockets struct{}
type flush struct{}

type InstrumentData struct {
	seqNum uint64
	lastHB time.Time
}

type Listener struct {
	extypes.Listener
	client         *ethclient.Client
	iterator       *evm.LogIterator
	securityID     uint64
	security       *models.Security
	executor       *actor.PID
	instrumentData *InstrumentData
	logger         *log.Logger
	socketTicker   *time.Ticker
	flushTicker    *time.Ticker
	lastPingTime   time.Time
	updates        *list.List
}

func NewListenerProducer(securityID uint64, dialerPool *xchangerUtils.DialerPool) actor.Producer {
	return func() actor.Actor {
		return NewListener(securityID, dialerPool)
	}
}

func NewListener(securityID uint64, dialerPool *xchangerUtils.DialerPool) actor.Actor {
	return &Listener{
		securityID: securityID,
	}
}

func (state *Listener) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.logger.Error("error initializing", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor started")

	case *actor.Stopping:
		if err := state.Clean(context); err != nil {
			state.logger.Error("error stopping", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor stopping")

	case *actor.Stopped:
		state.logger.Info("actor stopped")

	case *actor.Restarting:
		if err := state.Clean(context); err != nil {
			state.logger.Error("error restarting", log.Error(err))
			// Attention, no panic in restarting or infinite loop
		}
		state.logger.Info("actor restarting")

	case *messages.UnipoolV3DataRequest:
		if err := state.OnUnipoolV3DataRequest(context); err != nil {
			state.logger.Error("error processing OnUnipoolV3DataRequest", log.Error(err))
			panic(err)
		}

	case *types.Log:
		if err := state.onLog(context); err != nil {
			state.logger.Error("error processing log", log.Error(err))
			panic(err)
		}

	case *flush:
		if err := state.onFlush(context); err != nil {
			state.logger.Error("error flushing updates", log.Error(err))
			panic(err)
		}

	case *checkSockets:
		if err := state.onCheckSockets(context); err != nil {
			state.logger.Error("error checking sockets", log.Error(err))
			panic(err)
		}
	}
}

func (state *Listener) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()),
		log.String("security-id", fmt.Sprintf("%d", state.securityID)))
	state.executor = actor.NewPID(context.ActorSystem().Address(), "executor/exchanges/"+constants.BITZ.Name+"_executor")

	res, err := context.RequestFuture(state.executor, &messages.SecurityDefinitionRequest{
		RequestID:  0,
		Instrument: &models.Instrument{SecurityID: wrapperspb.UInt64(state.securityID)},
	}, 5*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error fetching security definition: %v", err)
	}
	def := res.(*messages.SecurityDefinitionResponse)
	if !def.Success {
		return fmt.Errorf("error fetching security definition: %s", def.RejectionReason.String())
	}
	state.security = def.Security
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()),
		log.String("security-id", fmt.Sprintf("%d", state.securityID)),
		log.String("exchange", state.security.Exchange.Name),
		log.String("symbol", state.security.Symbol))
	state.instrumentData = &InstrumentData{
		seqNum: uint64(time.Now().UnixNano()),
	}

	client, err := ethclient.Dial(evm.ETH_CLIENT_WS)
	if err != nil {
		return fmt.Errorf("error while dialing eth rpc client %v", err)
	}
	state.client = client

	if err := state.subscribeLogs(context); err != nil {
		return fmt.Errorf("error subscribing to logs %v", err)
	}

	socketTicker := time.NewTicker(5 * time.Second)
	state.socketTicker = socketTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-socketTicker.C:
				context.Send(pid, &checkSockets{})
			case <-time.After(10 * time.Second):
				if state.socketTicker != socketTicker {
					return
				}
			}
		}
	}(context.Self())

	flushTicker := time.NewTicker(10 * time.Second)
	state.flushTicker = flushTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-flushTicker.C:
				context.Send(pid, &flush{})
			case <-time.After(20 * time.Second):
				if state.flushTicker != flushTicker {
					return
				}
			}
		}
	}(context.Self())

	state.updates = list.New()

	return nil
}

func (state *Listener) OnUnipoolV3DataRequest(context actor.Context) error {
	msg := context.Message().(*messages.UnipoolV3DataRequest)
	resp := &messages.UnipoolV3DataResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		SeqNum:     state.instrumentData.seqNum,
	}

	context.Respond(resp)
	return nil
}

func (state *Listener) subscribeLogs(context actor.Context) error {
	if state.iterator != nil {
		state.iterator.Close()
	}

	symbol := state.security.Symbol

	uabi, err := uniswap.UniswapMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("error getting contract abi %v", err)
	}
	it := evm.NewLogIterator(uabi)

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
		return fmt.Errorf("error getting topics %v", err)
	}
	fQuery := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(symbol)},
		Topics:    topics,
	}

	ctx, cancel := goContext.WithTimeout(goContext.Background(), 10*time.Second)
	err = it.WatchLogs(state.client, ctx, fQuery)
	if err != nil {
		cancel()
		return fmt.Errorf("error watching logs: %v", err)
	}
	state.iterator = it

	go func(pid *actor.PID) {
		defer cancel()
		for state.iterator.Next() {
			context.Send(pid, state.iterator.Log)
		}
	}(context.Self())

	return nil
}

func (state *Listener) Clean(context actor.Context) error {
	if state.socketTicker != nil {
		state.socketTicker.Stop()
		state.socketTicker = nil
	}

	return nil
}

func (state *Listener) onLog(context actor.Context) error {
	msg := context.Message().(*types.Log)
	if msg.Removed {
		el := state.updates.Front()
		removed := false
		for ; el != nil; el = el.Next() {
			update := el.Value.(*models.UPV3Update)
			if update.Block == msg.BlockNumber {
				state.updates.Remove(el)
				removed = true
				break
			}
		}
		if !removed {
			return fmt.Errorf("removed log not found for block %d", msg.BlockNumber)
		} else {
			return nil
		}
	}
	uabi, err := uniswap.UniswapMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("error getting contract abi %v", err)
	}
	var update *models.UPV3Update
	header, err := state.client.HeaderByNumber(goContext.Background(), big.NewInt(int64(msg.BlockNumber)))
	if err != nil {
		return fmt.Errorf("error getting block number: %v", err)
	}
	switch msg.Topics[0] {
	case uabi.Events["Initialize"].ID:
		event := new(uniswap.UniswapInitialize)
		if err := evm.UnpackLog(uabi, event, "Initialize", *msg); err != nil {
			return fmt.Errorf("error unpacking the log %v", err)
		}
		update = &models.UPV3Update{
			Initialize: &gorderbook.UPV3Initialize{
				SqrtPriceX96: event.SqrtPriceX96.Bytes(),
				Tick:         int32(event.Tick.Int64()),
			},
			Block:     msg.BlockNumber,
			Timestamp: utils.SecondToTimestamp(header.Time),
		}
	case uabi.Events["Mint"].ID:
		event := new(uniswap.UniswapMint)
		if err := evm.UnpackLog(uabi, event, "Mint", *msg); err != nil {
			return fmt.Errorf("error unpacking the log %v", err)
		}
		update = &models.UPV3Update{
			Mint: &gorderbook.UPV3Mint{
				Owner:     event.Owner[:],
				TickLower: int32(event.TickLower.Int64()),
				TickUpper: int32(event.TickUpper.Int64()),
				Amount:    event.Amount.Bytes(),
				Amount0:   event.Amount0.Bytes(),
				Amount1:   event.Amount1.Bytes(),
			},
			Block:     msg.BlockNumber,
			Timestamp: utils.SecondToTimestamp(header.Time),
		}
	case uabi.Events["Burn"].ID:
		event := new(uniswap.UniswapBurn)
		if err := evm.UnpackLog(uabi, event, "Burn", *msg); err != nil {
			return fmt.Errorf("error unpacking the log %v", err)
		}
		update = &models.UPV3Update{
			Burn: &gorderbook.UPV3Burn{
				Owner:     event.Owner[:],
				TickLower: int32(event.TickLower.Int64()),
				TickUpper: int32(event.TickUpper.Int64()),
				Amount:    event.Amount.Bytes(),
				Amount0:   event.Amount0.Bytes(),
				Amount1:   event.Amount1.Bytes(),
			},
			Block:     msg.BlockNumber,
			Timestamp: utils.SecondToTimestamp(header.Time),
		}
	case uabi.Events["Swap"].ID:
		event := new(uniswap.UniswapSwap)
		if err := evm.UnpackLog(uabi, event, "Swap", *msg); err != nil {
			return fmt.Errorf("error unpacking the log %v", err)
		}
		update = &models.UPV3Update{
			Swap: &gorderbook.UPV3Swap{
				SqrtPriceX96: event.SqrtPriceX96.Bytes(),
				Tick:         int32(event.Tick.Int64()),
				Amount0:      event.Amount0.Bytes(),
				Amount1:      event.Amount1.Bytes(),
			},
			Block:     msg.BlockNumber,
			Timestamp: utils.SecondToTimestamp(header.Time),
		}
	case uabi.Events["Collect"].ID:
		event := new(uniswap.UniswapCollect)
		if err := evm.UnpackLog(uabi, event, "Collect", *msg); err != nil {
			return fmt.Errorf("error unpacking the log %v", err)
		}
		update = &models.UPV3Update{
			Collect: &gorderbook.UPV3Collect{
				Owner:            event.Recipient[:],
				TickLower:        int32(event.TickLower.Int64()),
				TickUpper:        int32(event.TickUpper.Int64()),
				AmountRequested0: event.Amount0.Bytes(),
				AmountRequested1: event.Amount1.Bytes(),
			},
			Block:     msg.BlockNumber,
			Timestamp: utils.SecondToTimestamp(header.Time),
		}
	case uabi.Events["Flash"].ID:
		event := new(uniswap.UniswapFlash)
		if err := evm.UnpackLog(uabi, event, "Flash", *msg); err != nil {
			return fmt.Errorf("error unpacking the log %v", err)
		}
		update = &models.UPV3Update{
			Flash: &gorderbook.UPV3Flash{
				Amount0: event.Amount0.Bytes(),
				Amount1: event.Amount1.Bytes(),
			},
			Block:     msg.BlockNumber,
			Timestamp: utils.SecondToTimestamp(header.Time),
		}
	case uabi.Events["SetFeeProtocol"].ID:
		event := new(uniswap.UniswapSetFeeProtocol)
		if err := evm.UnpackLog(uabi, event, "SetFeeProtocol", *msg); err != nil {
			return fmt.Errorf("error unpacking the log %v", err)
		}
		update = &models.UPV3Update{
			SetFeeProtocol: &gorderbook.UPV3SetFeeProtocol{
				FeesProtocol: uint32(event.FeeProtocol0New) + uint32(event.FeeProtocol1New)<<8,
			},
			Block:     msg.BlockNumber,
			Timestamp: utils.SecondToTimestamp(header.Time),
		}
	case uabi.Events["CollectProtocol"].ID:
		event := new(uniswap.UniswapCollectProtocol)
		if err := evm.UnpackLog(uabi, event, "CollectProtocol", *msg); err != nil {
			return fmt.Errorf("error unpacking the log %v", err)
		}
		update = &models.UPV3Update{
			CollectProtocol: &gorderbook.UPV3CollectProtocol{
				AmountRequested0: event.Amount0.Bytes(),
				AmountRequested1: event.Amount1.Bytes(),
			},
			Block:     msg.BlockNumber,
			Timestamp: utils.SecondToTimestamp(header.Time),
		}
	default:
		return fmt.Errorf("received unknown event: %v", msg.Topics[0])
	}

	state.updates.PushBack(update)
	return nil
}

func (state *Listener) onFlush(context actor.Context) error {
	current, err := state.client.BlockNumber(goContext.Background())
	if err != nil {
		return fmt.Errorf("error fetching current block number")
	}
	// Post all the updates with block < current - 4
	for el := state.updates.Front(); el != nil; el = state.updates.Front() {
		update := el.Value.(*models.UPV3Update)
		if update.Block <= current-4 {
			context.Send(context.Parent(), &messages.UnipoolV3DataIncrementalRefresh{
				SeqNum: state.instrumentData.seqNum + 1,
				Update: update,
			})
			state.instrumentData.seqNum += 1
			state.updates.Remove(el)
		} else {
			break
		}
	}
	return nil
}

func (state *Listener) onCheckSockets(context actor.Context) error {
	if time.Since(state.lastPingTime) > 5*time.Second {
		ctx, cancel := goContext.WithTimeout(goContext.Background(), 5*time.Second)
		defer cancel()
		_, err := state.client.BlockNumber(ctx)
		if err != nil {
			state.logger.Info("eth client err", log.Error(err))
			if err := state.subscribeLogs(context); err != nil {
				return fmt.Errorf("error subscribing to logs: %v", err)
			}
		}
	}

	if time.Since(state.instrumentData.lastHB) > 2*time.Second {
		context.Send(context.Parent(), &messages.UnipoolV3DataIncrementalRefresh{
			SeqNum: state.instrumentData.seqNum + 1,
		})
		state.instrumentData.seqNum += 1
		state.instrumentData.lastHB = time.Now()
	}
	return nil
}
