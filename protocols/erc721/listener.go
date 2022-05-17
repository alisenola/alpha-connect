package erc721

import (
	goContext "context"
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/protocols/types"
	"gitlab.com/alphaticks/alpha-connect/utils"
	gorderbook_models "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	xutils "gitlab.com/alphaticks/xchanger/eth"
	"gitlab.com/alphaticks/xchanger/protocols"
	nft "gitlab.com/alphaticks/xchanger/protocols/erc721"
)

type checkSockets struct{}

//Check socket for uniswap also

type InstrumentData struct {
	seqNum          uint64
	lastBlockUpdate uint64
	lastHB          time.Time
}

type Listener struct {
	types.Listener
	client       *ethclient.Client
	instrument   *InstrumentData
	iterator     *xutils.LogIterator
	address      [20]byte
	collection   *models.ProtocolAsset
	logger       *log.Logger
	socketTicker *time.Ticker
	lastPingTime time.Time
}

func NewListenerProducer(collection *models.ProtocolAsset) actor.Producer {
	return func() actor.Actor {
		return NewListener(collection)
	}
}

func NewListener(collection *models.ProtocolAsset) actor.Actor {
	return &Listener{
		instrument:   nil,
		iterator:     nil,
		collection:   collection,
		logger:       nil,
		socketTicker: nil,
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
		state.logger.Info("actor stopped")
	case *actor.Stopped:
		state.logger.Info("actor stopped")
	case *actor.Restarting:
		if err := state.Clean(context); err != nil {
			state.logger.Error("error restarting", log.Error(err))
			// No panic or we get an infinite loop
		}
		state.logger.Info("actor restarting")
	case *checkSockets:
		if err := state.onCheckSockets(context); err != nil {
			state.logger.Error("error checking sockets", log.Error(err))
			panic(err)
		}
	case *ethTypes.Log:
		if err := state.onLog(context); err != nil {
			state.logger.Error("error processing log", log.Error(err))
			panic(err)
		}
	case *messages.ProtocolAssetTransferRequest:
		if err := state.OnProtocolAssetTransferRequest(context); err != nil {
			state.logger.Error("error processing ProtocolAssetTransferRequest", log.Error(err))
		}
	}
}

func (state *Listener) Initialize(context actor.Context) error {
	addr := state.collection.ContractAddress
	if addr == nil || len(addr.Value) < 2 {
		return fmt.Errorf("invalid collection address")
	}
	addressBig, ok := big.NewInt(1).SetString(addr.Value[2:], 16)
	if !ok {
		return fmt.Errorf("invalid collection address: %s", addr.Value)
	}
	copy(state.address[:], addressBig.Bytes())

	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()),
		log.String("protocol", "ERC-721"),
		log.String("contract", addr.Value),
	)

	state.instrument = &InstrumentData{
		seqNum:          uint64(time.Now().UnixNano()),
		lastBlockUpdate: 0,
	}

	client, err := ethclient.Dial(protocols.ETH_CLIENT_WS)
	if err != nil {
		return fmt.Errorf("error dialing eth rpc client: %v", err)
	}
	state.client = client

	if err := state.subscribeLogs(context); err != nil {
		return fmt.Errorf("error subscribing to logs: %v", err)
	}

	socketTicker := time.NewTicker(5 * time.Second)
	state.socketTicker = socketTicker
	go func(pid *actor.PID) {
		select {
		case <-socketTicker.C:
			context.Send(pid, &checkSockets{})
		case <-time.After(10 * time.Second):
			return
		}
	}(context.Self())
	return nil
}

func (state *Listener) OnProtocolAssetTransferRequest(context actor.Context) error {
	req := context.Message().(*messages.ProtocolAssetTransferRequest)
	context.Respond(&messages.ProtocolAssetTransferResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		SeqNum:     state.instrument.seqNum,
	})
	return nil
}

func (state *Listener) subscribeLogs(context actor.Context) error {
	if state.iterator != nil {
		state.iterator.Close()
	}

	eabi, err := nft.ERC721MetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("error getting abi: %v", err)
	}
	it := xutils.NewLogIterator(eabi)

	query := [][]interface{}{{
		eabi.Events["Transfer"].ID,
	}}
	topics, err := abi.MakeTopics(query...)
	if err != nil {
		return fmt.Errorf("error making topics: %v", err)
	}
	fQuery := ethereum.FilterQuery{
		Addresses: []common.Address{state.address},
		Topics:    topics,
	}

	ctx, _ := goContext.WithTimeout(goContext.Background(), 10*time.Second)
	err = it.WatchLogs(state.client, ctx, fQuery)
	if err != nil {
		return fmt.Errorf("error watching logs: %v", err)
	}
	state.iterator = it
	go func(pid *actor.PID) {
		for state.iterator.Next() {
			context.Send(pid, state.iterator.Log)
		}
	}(context.Self())

	return nil
}

func (state *Listener) onLog(context actor.Context) error {
	req := context.Message().(*ethTypes.Log)
	eabi, err := nft.ERC721MetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("error getting abi: %v", err)
	}

	var updt *models.ProtocolAssetUpdate
	switch req.Topics[0] {
	case eabi.Events["Transfer"].ID:
		event := nft.ERC721Transfer{}
		if err := xutils.UnpackLog(eabi, &event, "Transfer", *req); err != nil {
			return fmt.Errorf("error unpacking log: %v", err)
		}
		block, err := state.client.BlockByNumber(goContext.Background(), big.NewInt(int64(req.BlockNumber)))
		if err != nil {
			return fmt.Errorf("error getting block number: %v", err)
		}
		updt = &models.ProtocolAssetUpdate{
			Transfer: &gorderbook_models.AssetTransfer{
				From:    event.From[:],
				To:      event.To[:],
				TokenId: event.TokenId.Bytes(),
			},
			Removed:   event.Raw.Removed,
			Block:     event.Raw.BlockNumber,
			Timestamp: utils.SecondToTimestamp(block.Time()),
		}
	}
	context.Send(
		context.Parent(),
		&messages.ProtocolAssetDataIncrementalRefresh{
			ResponseID: uint64(time.Now().UnixNano()),
			SeqNum:     state.instrument.seqNum + 1,
			Update:     updt,
		})
	state.instrument.lastBlockUpdate = updt.Block
	state.instrument.seqNum += 1
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

	if time.Since(state.instrument.lastHB) > 2*time.Second {
		context.Send(context.Parent(), &messages.ProtocolAssetDataIncrementalRefresh{
			SeqNum: state.instrument.seqNum + 1,
		})
		state.instrument.seqNum += 1
		state.instrument.lastHB = time.Now()
	}
	return nil
}

func (state *Listener) Clean(context actor.Context) error {
	return nil
}
