package erc721

import (
	"fmt"
	sabi "gitlab.com/alphaticks/abigen-starknet/accounts/abi"
	snft "gitlab.com/alphaticks/xchanger/protocols/erc721/svm"
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
	xutils "gitlab.com/alphaticks/xchanger/chains/evm"
	nft "gitlab.com/alphaticks/xchanger/protocols/erc721/evm"
)

type checkTimeout struct{}
type updateRequest struct{}

type protoUpdates struct {
	Transfer  *gorderbook_models.AssetTransfer
	Timestamp *timestamppb.Timestamp
	Block     uint64
	Hash      common.Hash
	Index     uint
}

type Listener struct {
	types.BaseListener
	executor        *actor.PID
	address         [32]byte
	collection      *models.ProtocolAsset
	eabi            *abi.ABI
	sabi            *sabi.ABI
	seqNum          uint64
	lastBlock       uint64
	lastRefreshTime time.Time
	logger          *log.Logger
	updateTicker    *time.Ticker
	timeoutTicker   *time.Ticker
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
			panic(err)
		}

	case *updateRequest:
		if err := state.OnSVMUpdateRequest(context); err != nil {
			state.logger.Error("error processing OnUpdateRequest", log.Error(err))
			panic(err)
		}
	case *checkTimeout:
		if err := state.onCheckTimeout(context); err != nil {
			state.logger.Error("error processing onCheckTimeout", log.Error(err))
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
		log.String("protocol", "ERC-721"),
		log.String("chain", state.collection.Chain.Type),
	)
	addr := state.collection.ContractAddress
	if addr == nil || len(addr.Value) < 2 {
		return fmt.Errorf("invalid collection address")
	}
	state.logger.With(log.String("contract", addr.Value))

	state.executor = actor.NewPID(context.ActorSystem().Address(), "executor")
	switch state.collection.Chain.Type {
	case "EVM":
		addressBig, ok := big.NewInt(1).SetString(addr.Value[2:], 16)
		if !ok {
			return fmt.Errorf("invalid collection address: %s", addr.Value)
		}
		copy(state.address[:], addressBig.Bytes())
		if err := state.subscribeEVMLogs(context); err != nil {
			return err
		}
	case "SVM":
		add := common.HexToHash(addr.Value)
		copy(state.address[:], add.Bytes())
		if err := state.subscribeSVMEvents(context); err != nil {
			return err
		}
	}

	state.lastRefreshTime = time.Now()
	timeoutTicker := time.NewTicker(5 * time.Second)
	state.timeoutTicker = timeoutTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-timeoutTicker.C:
				context.Send(pid, &checkTimeout{})
			case <-time.After(15 * time.Second):
				if state.timeoutTicker != timeoutTicker {
					// Only stop if socket ticker has changed
					return
				}
			}
		}
	}(context.Self())

	return nil
}

func (state *Listener) subscribeEVMLogs(context actor.Context) error {
	eabi, err := nft.ERC721MetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("error getting eth abi: %v", err)
	}
	state.eabi = eabi

	topicsv := [][]interface{}{{
		state.eabi.Events["Transfer"].ID,
	}}
	topics, err := abi.MakeTopics(topicsv...)
	if err != nil {
		return fmt.Errorf("error making topics: %v", err)
	}
	var add [20]byte
	copy(add[:], state.address[:20])
	query := ethereum.FilterQuery{
		Addresses: []common.Address{add},
		Topics:    topics,
	}

	res, err := context.RequestFuture(state.executor, &messages.EVMLogsSubscribeRequest{
		RequestID:  state.collection.ProtocolAssetID,
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
	if !subRes.Success {
		return fmt.Errorf("error subscribing to EVM logs: %s", subRes.RejectionReason.String())
	}
	state.seqNum = subRes.SeqNum
	return nil
}

func (state *Listener) subscribeSVMEvents(context actor.Context) error {
	stabi, err := snft.ERC721MetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("error getting stark abi: %v", err)
	}
	state.sabi = stabi

	res, err := context.RequestFuture(state.executor, &messages.BlockNumberRequest{
		RequestID: state.collection.ProtocolAssetID,
		Chain:     state.collection.Chain,
	}, 10*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error fetching block number: %v", err)
	}
	b, ok := res.(*messages.BlockNumberResponse)
	if !ok {
		return fmt.Errorf("expected *messages.BlockNumberResponse, got %s", reflect.TypeOf(res).String())
	}
	if !b.Success {
		return fmt.Errorf("error fetching block number: %s", b.RejectionReason.String())
	}

	state.seqNum = b.ResponseID
	state.lastBlock = b.BlockNumber - 5

	ticker := time.NewTicker(20 * time.Second)
	state.updateTicker = ticker
	go func(pid *actor.PID) {
		for {
			select {
			case <-ticker.C:
				context.Send(pid, &updateRequest{})
			case <-time.After(40 * time.Second):
				if ticker != state.updateTicker {
					return
				}
			}
		}
	}(context.Self())
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
	state.lastRefreshTime = time.Now()
	var update *models.ProtocolAssetUpdate
	if refresh.Update != nil {
		//fmt.Println("REFRESH", refresh.SeqNum, state.seqNum, len(refresh.Update.Logs))
		update = &models.ProtocolAssetUpdate{
			BlockNumber: refresh.Update.BlockNumber,
			BlockTime:   timestamppb.New(refresh.Update.BlockTime),
		}
		for _, l := range refresh.Update.Logs {
			switch l.Topics[0] {
			case state.eabi.Events["Transfer"].ID:
				event := new(nft.ERC721Transfer)
				if err := xutils.UnpackLog(state.eabi, event, "Transfer", l); err != nil {
					return fmt.Errorf("error unpacking log: %v", err)
				}
				//fmt.Println("TRANSFER", event.From, event.To, event.TokenId.Text(10))
				update.Transfers = append(update.Transfers, &gorderbook_models.AssetTransfer{
					From:  event.From[:],
					To:    event.To[:],
					Value: event.TokenId.Bytes(),
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

func (state *Listener) OnSVMUpdateRequest(context actor.Context) error {
	res, err := context.RequestFuture(state.executor, &messages.BlockNumberRequest{
		RequestID: state.collection.ProtocolAssetID,
		Chain:     state.collection.Chain,
	}, 10*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error fetching block number: %v", err)
	}
	b, ok := res.(*messages.BlockNumberResponse)
	if !ok {
		return fmt.Errorf("expected *messages.BlockNumberResponse, got %s", reflect.TypeOf(res).String())
	}
	if !b.Success {
		return fmt.Errorf("error fetching block number: %s", b.RejectionReason.String())
	}

	var update *models.ProtocolAssetUpdate
	if b.BlockNumber >= state.lastBlock {
		resp, err := context.RequestFuture(state.executor, &messages.HistoricalProtocolAssetTransferRequest{
			RequestID:       uint64(time.Now().UnixNano()),
			ProtocolAssetID: state.collection.ProtocolAssetID,
			Start:           state.lastBlock,
			Stop:            state.lastBlock,
		}, 1*time.Minute).Result()
		if err != nil {
			return fmt.Errorf("error fetching svm events: %v", err)
		}
		evs, ok := resp.(*messages.HistoricalProtocolAssetTransferResponse)
		if !ok {
			return fmt.Errorf("expected *messages.HistoricalProtocolAssetTransferResponse, got %s", reflect.TypeOf(resp).String())
		}
		if !evs.Success {
			return fmt.Errorf("error fetching svm events, got %s", evs.RejectionReason)
		}
		if len(evs.Update) > 1 {
			return fmt.Errorf("fetched more than one block")
		}
		if len(evs.Update) == 1 {
			update = &models.ProtocolAssetUpdate{
				Transfers:   evs.Update[0].Transfers,
				BlockNumber: evs.Update[0].BlockNumber,
				BlockTime:   evs.Update[0].BlockTime,
			}
		}
		state.lastBlock += 1
	}

	context.Send(context.Parent(), &messages.ProtocolAssetDataIncrementalRefresh{
		Update: update,
		SeqNum: state.seqNum + 1,
	})
	state.seqNum += 1
	state.lastRefreshTime = time.Now()
	return nil
}

func (state *Listener) Clean(context actor.Context) error {
	if state.timeoutTicker != nil {
		state.timeoutTicker.Stop()
		state.timeoutTicker = nil
	}
	if state.updateTicker != nil {
		state.updateTicker.Stop()
		state.updateTicker = nil
	}
	return nil
}

func (state *Listener) onCheckTimeout(context actor.Context) error {
	if time.Since(state.lastRefreshTime) > 30*time.Second {
		return fmt.Errorf("timed-out")
	}
	return nil
}
