package erc721

import (
	"container/list"
	goContext "context"
	"fmt"
	"math/big"
	"reflect"
	"sort"
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
type flush struct{}
type protoUpdates struct {
	*models.ProtocolAssetUpdate
	Hash  common.Hash
	Index uint
}

//Check socket for uniswap also

type InstrumentData struct {
	seqNum          uint64
	lastBlockUpdate uint64
	lastHB          time.Time
}

type Listener struct {
	types.BaseListener
	client       *ethclient.Client
	instrument   *InstrumentData
	iterator     *xutils.LogIterator
	address      [20]byte
	collection   *models.ProtocolAsset
	logger       *log.Logger
	socketTicker *time.Ticker
	flushTicker  *time.Ticker
	lastPingTime time.Time
	updates      *list.List
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
		flushTicker:  nil,
		updates:      nil,
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
	case *flush:
		if err := state.onFlush(context); err != nil {
			state.logger.Error("error flushing updates", log.Error(err))
			panic(err)
		}
	case *ethTypes.Log:
		if err := state.onLog(context); err != nil {
			state.logger.Error("error processing log", log.Error(err))
			panic(err)
		}
	case *messages.ProtocolAssetDataRequest:
		if err := state.OnProtocolAssetDataRequest(context); err != nil {
			state.logger.Error("error processing OnProtocolAssetDataRequest", log.Error(err))
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

	client, err := ethclient.Dial(protocols.ETH_CLIENT_WS_2)
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
		for {
			select {
			case <-socketTicker.C:
				context.Send(pid, &checkSockets{})
			case <-time.After(10 * time.Second):
				if socketTicker != state.socketTicker {
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
				if flushTicker != state.flushTicker {
					return
				}
			}
		}
	}(context.Self())

	state.updates = list.New()

	return nil
}

func (state *Listener) OnProtocolAssetDataRequest(context actor.Context) error {
	req := context.Message().(*messages.ProtocolAssetDataRequest)
	context.Respond(&messages.ProtocolAssetDataResponse{
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
	if req.Removed {
		el := state.updates.Front()
		removed := false
		for ; el != nil; el = el.Next() {
			update := el.Value.(*protoUpdates)
			if update.Hash == req.TxHash {
				state.updates.Remove(el)
				removed = true
				break
			}
		}
		if !removed {
			return fmt.Errorf("removed log not found for hash %v", req.TxHash)
		}
	}
	eabi, err := nft.ERC721MetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("error getting abi: %v", err)
	}

	var updt *models.ProtocolAssetUpdate
	switch req.Topics[0] {
	case eabi.Events["Transfer"].ID:
		event := new(nft.ERC721Transfer)
		if err := xutils.UnpackLog(eabi, event, "Transfer", *req); err != nil {
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
			Block:     req.BlockNumber,
			Timestamp: utils.SecondToTimestamp(block.Time()),
		}
	}
	uptHash := &protoUpdates{
		ProtocolAssetUpdate: updt,
		Hash:                req.TxHash,
		Index:               req.Index,
	}
	state.updates.PushBack(uptHash)
	return nil
}

func (state *Listener) onCheckSockets(context actor.Context) error {
	if time.Since(state.lastPingTime) > 5*time.Second {
		ctx, cancel := goContext.WithTimeout(goContext.Background(), 10*time.Second)
		defer cancel()
		_, err := state.client.BlockNumber(ctx)
		if err != nil {
			state.logger.Info("eth client err", log.Error(err))
			return fmt.Errorf("error checking sockets: %v", err)
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

func (state *Listener) onFlush(context actor.Context) error {
	current, err := state.client.BlockNumber(goContext.Background())
	if err != nil {
		return fmt.Errorf("error fetching block number: %v", err)
	}
	var updts []*protoUpdates
	for el := state.updates.Front(); el != nil; el = state.updates.Front() {
		update := el.Value.(*protoUpdates)
		nextBlock := uint64(0)
		if el.Next() != nil {
			nextBlock = el.Next().Value.(*protoUpdates).Block
		}
		if update.Block <= current-4 {
			updts = append(updts, update)
			if update.Block != nextBlock {
				sort.Slice(updts, func(i, j int) bool {
					return updts[i].Index < updts[j].Index
				})
				var u []*models.ProtocolAssetUpdate
				for _, updt := range updts {
					u = append(u, updt.ProtocolAssetUpdate)
				}
				context.Send(context.Parent(),
					&messages.ProtocolAssetDataIncrementalRefresh{
						ResponseID: uint64(time.Now().UnixNano()),
						SeqNum:     state.instrument.seqNum + 1,
						Update:     u,
					})
				state.instrument.lastBlockUpdate = update.Block
				state.instrument.seqNum += 1
				updts = nil
			}
			state.updates.Remove(el)
		} else {
			break
		}
	}
	return nil
}

func (state *Listener) Clean(context actor.Context) error {
	if state.socketTicker != nil {
		state.socketTicker.Stop()
		state.socketTicker = nil
	}
	if state.flushTicker != nil {
		state.flushTicker.Stop()
		state.flushTicker = nil
	}
	return nil
}
