package erc721

import (
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/big"
	"reflect"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/protocols/types"
	gorderbook_models "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	xutils "gitlab.com/alphaticks/xchanger/eth"
	nft "gitlab.com/alphaticks/xchanger/protocols/erc721"
)

type protoUpdates struct {
	Transfer  *gorderbook_models.AssetTransfer
	Timestamp *timestamppb.Timestamp
	Block     uint64
	Hash      common.Hash
	Index     uint
}

type Listener struct {
	types.BaseListener
	executor   *actor.PID
	address    [20]byte
	collection *models.ProtocolAsset
	abi        *abi.ABI
	requestID  uint64 // Stays the same after restart, so the log executor knows it is the same subscription
	seqNum     uint64
	logger     *log.Logger
}

func NewListenerProducer(collection *models.ProtocolAsset) actor.Producer {
	return func() actor.Actor {
		return NewListener(collection)
	}
}

func NewListener(collection *models.ProtocolAsset) actor.Actor {
	return &Listener{
		collection: collection,
		logger:     nil,
		requestID:  uint64(time.Now().UnixNano()),
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
	case *messages.EVMLogsSubscribeRefresh:
		if err := state.OnEVMLogsSubscribeRefresh(context); err != nil {
			state.logger.Error("error processing onEVMLogsSubscribeRefresh", log.Error(err))
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

	eabi, err := nft.ERC721MetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("error getting abi: %v", err)
	}
	state.abi = eabi
	state.seqNum = uint64(time.Now().UnixNano())
	state.executor = actor.NewPID(context.ActorSystem().Address(), "executor")
	if err := state.subscribeLogs(context); err != nil {
		return fmt.Errorf("error subscribing to logs: %v", err)
	}

	return nil
}

func (state *Listener) Clean(context actor.Context) error {
	return nil
}

func (state *Listener) subscribeLogs(context actor.Context) error {
	topicsv := [][]interface{}{{
		state.abi.Events["Transfer"].ID,
	}}
	topics, err := abi.MakeTopics(topicsv...)
	if err != nil {
		return fmt.Errorf("error making topics: %v", err)
	}
	query := ethereum.FilterQuery{
		Addresses: []common.Address{state.address},
		Topics:    topics,
	}

	res, err := context.RequestFuture(state.executor, &messages.EVMLogsSubscribeRequest{
		RequestID:  state.requestID,
		Chain:      state.collection.Chain,
		Query:      query,
		Subscriber: context.Self(),
	}, 10*time.Second).Result()

	if err != nil {
		return fmt.Errorf("error subscribing to EVM logs: %v", err)
	}
	subRes, ok := res.(*messages.EVMLogsSubscribeResponse)
	if !ok {
		return fmt.Errorf("was expecting EVMLogsSubscribeResponse, got %s", reflect.TypeOf(subRes).String())
	}
	state.seqNum = subRes.SeqNum

	return nil
}

func (state *Listener) OnProtocolAssetDataRequest(context actor.Context) error {
	req := context.Message().(*messages.ProtocolAssetDataRequest)
	context.Respond(&messages.ProtocolAssetDataResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		SeqNum:     state.seqNum,
	})
	return nil
}

func (state *Listener) OnEVMLogsSubscribeRefresh(context actor.Context) error {
	refresh := context.Message().(*messages.EVMLogsSubscribeRefresh)
	if refresh.SeqNum <= state.seqNum {
		return nil
	}
	if refresh.SeqNum != state.seqNum+1 {
		return fmt.Errorf("out of order sequence")
	}

	state.seqNum = refresh.SeqNum
	var update *models.ProtocolAssetUpdate
	if refresh.Update != nil {
		//fmt.Println("REFRESH", refresh.SeqNum, state.seqNum, len(refresh.Update.Logs))
		update := &models.ProtocolAssetUpdate{
			BlockNumber: refresh.Update.BlockNumber,
			BlockTime:   timestamppb.New(refresh.Update.BlockTime),
		}
		for _, l := range refresh.Update.Logs {
			switch l.Topics[0] {
			case state.abi.Events["Transfer"].ID:
				event := new(nft.ERC721Transfer)
				if err := xutils.UnpackLog(state.abi, event, "Transfer", l); err != nil {
					return fmt.Errorf("error unpacking log: %v", err)
				}
				//fmt.Println("TRANSFER", event.From, event.To, event.TokenId.Text(10))
				update.Transfers = append(update.Transfers, &gorderbook_models.AssetTransfer{
					From:    event.From[:],
					To:      event.To[:],
					TokenId: event.TokenId.Bytes(),
				})
			}
		}
	}
	context.Send(context.Parent(), &messages.ProtocolAssetDataIncrementalRefresh{
		SeqNum: state.seqNum,
		Update: update,
	})

	return nil
}
