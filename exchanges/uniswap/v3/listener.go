package v3

import (
	goContext "context"
	"fmt"
	"reflect"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	gorderbook "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	"gitlab.com/alphaticks/xchanger/eth"
	uniswap "gitlab.com/alphaticks/xchanger/exchanges/uniswap/V3"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
)

type checkSockets struct{}
type postAggTrade struct{}

type InstrumentData struct {
	events          []*models.UPV3Update
	seqNum          uint64
	lastBlockUpdate uint64
}

type Listener struct {
	extypes.Listener
	iterator       *eth.LogIterator
	security       *models.Security
	dialerPool     *xchangerUtils.DialerPool
	instrumentData *InstrumentData
	logger         *log.Logger
	socketTicker   *time.Ticker
}

func NewListenerProducer(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Producer {
	return func() actor.Actor {
		return NewListener(security, dialerPool)
	}
}

func NewListener(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Actor {
	return &Listener{
		iterator:       nil,
		security:       security,
		dialerPool:     dialerPool,
		instrumentData: nil,
		logger:         nil,
		socketTicker:   nil,
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
	}
}

func (state *Listener) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()),
		log.String("exchange", state.security.Exchange.Name),
		log.String("symbol", state.security.Symbol))

	state.instrumentData = &InstrumentData{
		events:          []*models.UPV3Update{},
		seqNum:          uint64(time.Now().UnixNano()),
		lastBlockUpdate: 0,
	}

	if err := state.subscribeLogs(context); err != nil {
		return fmt.Errorf("error subscribing to logs %v", err)
	}

	socketTicker := time.NewTicker(5 * time.Second)
	state.socketTicker = socketTicker
	go func(pid *actor.PID) {
		for {
			select {
			case _ = <-socketTicker.C:
				context.Send(pid, &checkSockets{})
			case <-time.After(10 * time.Second):
				// timer stopped, we leave
				return
			}
		}
	}(context.Self())

	return nil
}

func (state *Listener) OnUnipoolV3DataRequest(context actor.Context) error {
	msg := context.Message().(*messages.UnipoolV3DataRequest)
	resp := &messages.UnipoolV3DataResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Update:     state.instrumentData.events,
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

	client, err := ethclient.Dial(ETH_CLIENT_URL)
	if err != nil {
		return fmt.Errorf("error while dialing eth rpc client %v", err)
	}
	uabi, err := uniswap.UniswapMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("error getting contract abi %v", err)
	}
	query := [][]interface{}{{
		uabi.Events["Initialize"].ID,
		uabi.Events["Mint"].ID,
		uabi.Events["Swap"].ID,
		uabi.Events["Burn"].ID,
		uabi.Events["Collect"].ID,
		uabi.Events["Flash"].ID,
	}}
	topics, err := abi.MakeTopics(query...)
	if err != nil {
		return fmt.Errorf("error getting topics %v", err)
	}
	fQuery := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(symbol)},
		Topics:    topics,
	}
	it := eth.NewLogIterator(uabi)
	ctx, _ := goContext.WithTimeout(goContext.Background(), 10*time.Second)
	it.WatchLogs(client, ctx, fQuery)

	go func(it *eth.LogIterator, pid *actor.PID) {
		for it.Next() {
			context.Send(pid, it.Log)
		}
	}(it, context.Self())

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
	uabi, err := uniswap.UniswapMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("error getting contract abi %v", err)
	}
	updt := &models.UPV3Update{}
	switch msg.Topics[0] {
	case uabi.Events["Initialize"].ID:
		event := new(uniswap.UniswapInitialize)
		if err := eth.UnpackLog(uabi, event, "Initialize", *msg); err != nil {
			return fmt.Errorf("error unpacking the log %v", err)
		}
		updt = &models.UPV3Update{
			Initialize: &gorderbook.UPV3Initialize{
				SqrtPriceX96: event.SqrtPriceX96.Bytes(),
				Tick:         int32(event.Tick.Int64()),
			},
			Block: msg.BlockNumber,
		}
		state.instrumentData.events = append(state.instrumentData.events, updt)
	case uabi.Events["Mint"].ID:
		event := new(uniswap.UniswapMint)
		if err := eth.UnpackLog(uabi, event, "Mint", *msg); err != nil {
			return fmt.Errorf("error unpacking the log %v", err)
		}
		updt = &models.UPV3Update{
			Mint: &gorderbook.UPV3Mint{
				Owner:     event.Owner[:],
				TickLower: int32(event.TickLower.Int64()),
				TickUpper: int32(event.TickUpper.Int64()),
				Amount:    event.Amount.Bytes(),
				Amount0:   event.Amount0.Bytes(),
				Amount1:   event.Amount1.Bytes(),
			},
			Block: msg.BlockNumber,
		}
		state.instrumentData.events = append(state.instrumentData.events, updt)
	case uabi.Events["Burn"].ID:
		event := new(uniswap.UniswapBurn)
		if err := eth.UnpackLog(uabi, event, "Burn", *msg); err != nil {
			return fmt.Errorf("error unpacking the log %v", err)
		}
		updt = &models.UPV3Update{
			Burn: &gorderbook.UPV3Burn{
				Owner:     event.Owner[:],
				TickLower: int32(event.TickLower.Int64()),
				TickUpper: int32(event.TickUpper.Int64()),
				Amount:    event.Amount.Bytes(),
				Amount0:   event.Amount0.Bytes(),
				Amount1:   event.Amount1.Bytes(),
			},
			Block: msg.BlockNumber,
		}
		state.instrumentData.events = append(state.instrumentData.events, updt)
	case uabi.Events["Swap"].ID:
		event := new(uniswap.UniswapSwap)
		if err := eth.UnpackLog(uabi, event, "Swap", *msg); err != nil {
			return fmt.Errorf("error unpacking the log %v", err)
		}
		updt = &models.UPV3Update{
			Swap: &gorderbook.UPV3Swap{
				SqrtPriceX96: event.SqrtPriceX96.Bytes(),
				Tick:         int32(event.Tick.Int64()),
				Amount0:      event.Amount0.Bytes(),
				Amount1:      event.Amount1.Bytes(),
			},
			Block: msg.BlockNumber,
		}
		state.instrumentData.events = append(state.instrumentData.events, updt)
	case uabi.Events["Collect"].ID:
		event := new(uniswap.UniswapCollect)
		if err := eth.UnpackLog(uabi, event, "Collect", *msg); err != nil {
			return fmt.Errorf("error unpacking the log %v", err)
		}
		updt = &models.UPV3Update{
			Collect: &gorderbook.UPV3Collect{
				Owner:            event.Recipient[:],
				TickLower:        int32(event.TickLower.Int64()),
				TickUpper:        int32(event.TickUpper.Int64()),
				AmountRequested0: event.Amount0.Bytes(),
				AmountRequested1: event.Amount1.Bytes(),
			},
			Block: msg.BlockNumber,
		}
		state.instrumentData.events = append(state.instrumentData.events, updt)
	case uabi.Events["Flash"].ID:
		event := new(uniswap.UniswapFlash)
		if err := eth.UnpackLog(uabi, event, "Flash", *msg); err != nil {
			return fmt.Errorf("error unpacking the log %v", err)
		}
		updt = &models.UPV3Update{
			Flash: &gorderbook.UPV3Flash{
				Amount0: event.Amount0.Bytes(),
				Amount1: event.Amount1.Bytes(),
			},
			Block: msg.BlockNumber,
		}
		state.instrumentData.events = append(state.instrumentData.events, updt)
	}

	state.instrumentData.lastBlockUpdate = msg.BlockNumber
	context.Send(context.Parent(), &messages.UnipoolV3DataIncrementalRefresh{
		SeqNum: state.instrumentData.seqNum + 1,
		Update: updt,
	})
	state.instrumentData.seqNum += 1

	return nil
}
