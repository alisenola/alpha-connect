package erc721

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
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/protocols/types"
	gorderbook_models "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	utils "gitlab.com/alphaticks/xchanger/eth"
	"gitlab.com/alphaticks/xchanger/protocols"
	nft "gitlab.com/alphaticks/xchanger/protocols/erc721"
)

type checkSockets struct{}

//Check socket for uniswap also

type InstrumentData struct {
	events          []*models.AssetUpdate
	seqNum          uint64
	lastBlockUpdate uint64
	lastHB          time.Time
}

type Listener struct {
	types.Listener
	client       *ethclient.Client
	instrument   *InstrumentData
	iterator     *utils.LogIterator
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
	case *messages.AssetTransferRequest:
		if err := state.OnAssetTransferRequest(context); err != nil {
			state.logger.Error("error processing AssetTransferRequest", log.Error(err))
		}
	}
}

func (state *Listener) Initialize(context actor.Context) error {
	address := common.Bytes2Hex(state.collection.Address)
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()),
		log.String("protocol", "ERC-721"),
		log.String("contract", address),
	)

	state.instrument = &InstrumentData{
		events:          make([]*models.AssetUpdate, 0),
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

func (state *Listener) OnAssetTransferRequest(context actor.Context) error {
	req := context.Message().(*messages.AssetTransferRequest)
	context.Respond(&messages.AssetTransferResponse{
		RequestID:    req.RequestID,
		ResponseID:   uint64(time.Now().UnixNano()),
		AssetUpdated: state.instrument.events,
		Success:      true,
		SeqNum:       state.instrument.seqNum,
	})
	return nil
}

func (state *Listener) subscribeLogs(context actor.Context) error {
	if state.iterator != nil {
		state.iterator.Close()
	}

	address := common.BytesToAddress(state.collection.Address)
	eabi, err := nft.ERC721MetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("error getting abi: %v", err)
	}
	it := utils.NewLogIterator(eabi)

	query := [][]interface{}{{
		eabi.Events["Transfer"].ID,
	}}
	topics, err := abi.MakeTopics(query...)
	if err != nil {
		return fmt.Errorf("error making topics: %v", err)
	}
	fQuery := ethereum.FilterQuery{
		Addresses: []common.Address{address},
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

	var updt *models.AssetUpdate
	switch req.Topics[0] {
	case eabi.Events["Transfer"].ID:
		event := nft.ERC721Transfer{}
		if err := utils.UnpackLog(eabi, &event, "Transfer", *req); err != nil {
			return fmt.Errorf("error unpacking log: %v", err)
		}
		updt = &models.AssetUpdate{
			Transfer: &gorderbook_models.AssetTransfer{
				From:    event.From[:],
				To:      event.To[:],
				TokenId: event.TokenId.Bytes(),
			},
			Block:   event.Raw.BlockNumber,
			Removed: event.Raw.Removed,
		}
	}

	context.Send(
		context.Parent(),
		&messages.AssetDataIncrementalRefresh{
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
		context.Send(context.Parent(), &messages.AssetDataIncrementalRefresh{
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
