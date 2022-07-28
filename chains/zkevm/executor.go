package zkevm

import (
	"container/list"
	goContext "context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"sync"

	"math/big"
	"reflect"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"github.com/ethereum/go-ethereum/ethclient"
	chtypes "gitlab.com/alphaticks/alpha-connect/chains/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/chains/zks"
)

type flushLogs struct{}

type logsSubscription struct {
	sync.RWMutex
	logs         *list.List
	subscription ethereum.Subscription
	subscriber   *actor.PID
	seqNum       uint64
	lastPingTime time.Time
	ch           chan types.Log
}

type Executor struct {
	chtypes.BaseExecutor
	protocolAssets map[uint64]*models.ProtocolAsset
	logger         *log.Logger
	client         *ethclient.Client
	rpc            *rpc.Client
	rateLimit      *exchanges.RateLimit
	subscriptions  map[uint64]*logsSubscription
	flushTicker    *time.Ticker
	endpoint       string
}

func NewExecutor(registry registry.PublicRegistryClient, rpc string) actor.Actor {
	e := &Executor{
		protocolAssets: nil,
		logger:         nil,
		endpoint:       rpc,
	}
	e.Registry = registry
	return e
}

func (state *Executor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *flushLogs:
		if err := state.onFlushLogs(context); err != nil {
			panic(err)
		}
	case *actor.Terminated:
		if err := state.OnTerminated(context); err != nil {
			panic(err)
		}
	default:
		chtypes.ReceiveExecutor(state, context)
	}
}

func (state *Executor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	r, err := rpc.DialContext(goContext.Background(), state.endpoint)
	if err != nil {
		return fmt.Errorf("error while dialing zksync rpc client %v", err)
	}
	state.client = ethclient.NewClient(r)
	state.rpc = r
	// state.rateLimit =
	state.subscriptions = make(map[uint64]*logsSubscription)

	flushTicker := time.NewTicker(5 * time.Second)
	state.flushTicker = flushTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-flushTicker.C:
				context.Send(pid, &flushLogs{})
			case <-time.After(15 * time.Second):
				if state.flushTicker != flushTicker {
					// Only stop if socket ticker has changed
					return
				}
			}
		}
	}(context.Self())

	return nil
}

func (state *Executor) OnBlockNumberRequest(context actor.Context) error {
	req := context.Message().(*messages.BlockNumberRequest)
	go func(sender *actor.PID) {
		res := &messages.BlockNumberResponse{
			RequestID:  req.RequestID,
			ResponseID: uint64(time.Now().UnixNano()),
		}

		current, err := state.client.BlockNumber(goContext.Background())
		if err != nil {
			state.logger.Warn("error getting block number", log.Error(err))
			res.RejectionReason = messages.RejectionReason_RPCError
			context.Send(sender, res)
			return
		}
		res.BlockNumber = current
		res.Success = true
		context.Send(sender, res)
	}(context.Sender())

	return nil
}

func (state *Executor) OnBlockInfoRequest(context actor.Context) error {
	req := context.Message().(*messages.BlockInfoRequest)
	go func(sender *actor.PID) {
		res := &messages.BlockInfoResponse{
			RequestID:  req.RequestID,
			ResponseID: uint64(time.Now().UnixNano()),
		}

		var header *zks.Header
		ctx, cancel := goContext.WithTimeout(goContext.Background(), 10*time.Second)
		defer cancel()
		num := hexutil.EncodeBig(big.NewInt(int64(req.BlockNumber)))
		err := state.rpc.CallContext(ctx, &header, "eth_getBlockByNumber", num, false)
		if err != nil {
			state.logger.Warn("error getting header", log.Error(err))
			res.RejectionReason = messages.RejectionReason_RPCError
			context.Send(sender, res)
			return
		}
		res.BlockTime = time.Unix(int64(header.Time), 0)
		res.Success = true
		context.Send(sender, res)
	}(context.Sender())

	return nil
}

func (state *Executor) OnEVMContractCallRequest(context actor.Context) error {
	req := context.Message().(*messages.EVMContractCallRequest)
	go func(sender *actor.PID) {
		res := &messages.EVMContractCallResponse{
			RequestID:  req.RequestID,
			ResponseID: uint64(time.Now().UnixNano()),
		}

		out, err := state.client.CallContract(goContext.Background(), req.Msg, big.NewInt(int64(req.BlockNumber)))
		if err != nil {
			state.logger.Warn("error with contract call", log.Error(err))
			res.RejectionReason = messages.RejectionReason_RPCError
			switch err.Error() {
			case "invalid jump destination":
				res.RejectionReason = messages.RejectionReason_InvalidJumpDestination
			}
			context.Send(sender, res)
			return
		}
		res.Out = out
		res.Success = true
		context.Send(sender, res)
	}(context.Sender())
	return nil
}

func (state *Executor) OnEVMLogsQueryRequest(context actor.Context) error {
	req := context.Message().(*messages.EVMLogsQueryRequest)
	/*
		query := &ethereum.FilterQuery{}
		if req.BlockHash != nil {
			if len(req.BlockHash) != 32 {
				context.Respond(&messages.EVMLogsQueryResponse{
					// TODO
				})
				return nil
			}
			var tmp common.Hash
			copy(tmp[:], req.BlockHash)
			query.BlockHash = &tmp
		}
		if req.FromBlock != nil {
			query.FromBlock = big.NewInt(int64(req.FromBlock.Value))
		}
		if req.ToBlock != nil {
			query.ToBlock = big.NewInt(int64(req.ToBlock.Value))
		}
		for _, addr := range req.Addresses {
			if len(addr) != 20 {
				context.Respond(&messages.EVMLogsQueryResponse{
					// TODO
				})
				return nil
			}
			var tmp common.Address
			copy(tmp[:], addr[:])
			query.Addresses = append(query.Addresses, tmp)
		}
		for _, topics := range req.Topics {
			var tops []common.Hash
			for _, topic := range topics.Topic {
				if len(topic) != 32 {
					context.Respond(&messages.EVMLogsQueryResponse{
						// TODO
					})
					return nil
				}
				var tmp common.Hash
				copy(tmp[:], topic)
				tops = append(tops, tmp)
			}
			query.Topics = append(query.Topics, tops)
		}
	*/

	go func(sender *actor.PID) {
		res := &messages.EVMLogsQueryResponse{
			RequestID:  req.RequestID,
			ResponseID: uint64(time.Now().UnixNano()),
		}

		logs, err := state.client.FilterLogs(goContext.Background(), req.Query)
		if err != nil {
			state.logger.Warn("error filtering logs", log.Error(err))
			res.RejectionReason = messages.RejectionReason_RPCError
			context.Send(sender, res)
			return
		}

		res.Logs = logs
		var lastBlock uint64 = 0
		var lastTime uint64 = 0
		for _, l := range logs {
			if lastBlock != l.BlockNumber {
				var header *zks.Header
				ctx, cancel := goContext.WithTimeout(goContext.Background(), 15*time.Second)
				num := hexutil.EncodeBig(big.NewInt(int64(l.BlockNumber)))
				err := state.rpc.CallContext(ctx, &header, "eth_getBlockByNumber", num, false)
				if err != nil {
					state.logger.Warn("error getting header", log.Error(err))
					res.RejectionReason = messages.RejectionReason_RPCError
					context.Send(sender, res)
					cancel()
					return
				}
				lastTime = header.Time
				lastBlock = l.BlockNumber
				cancel()
			}
			res.Times = append(res.Times, lastTime)
		}
		res.Success = true
		context.Send(sender, res)
	}(context.Sender())

	// When you get a filter log query, you need to
	return nil
}

func (state *Executor) OnEVMLogsSubscribeRequest(context actor.Context) error {
	req := context.Message().(*messages.EVMLogsSubscribeRequest)
	res := &messages.EVMLogsSubscribeResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
	}

	ch := make(chan types.Log)
	sub, err := state.client.SubscribeFilterLogs(goContext.Background(), req.Query, ch)
	if err != nil {
		res.RejectionReason = messages.RejectionReason_RPCError
		context.Respond(res)
		return nil
	}
	subs, ok := state.subscriptions[req.RequestID]
	if !ok {
		subs = &logsSubscription{
			logs:         list.New(),
			subscription: sub,
			seqNum:       uint64(time.Now().UnixNano()),
			subscriber:   req.Subscriber,
			ch:           ch,
		}
		go func() {
			for {
				l, ok := <-ch
				if !ok {
					return
				}
				subs.Lock()
				if l.Removed {
					el := subs.logs.Front()
					removed := false
					for ; el != nil; el = el.Next() {
						update := el.Value.(*types.Log)
						if update.TxHash == l.TxHash {
							subs.logs.Remove(el)
							removed = true
							break
						}
					}
					if !removed {
						// TODO restart subscription
						fmt.Println("COULD NOT FIND REMOVED TX")
					}
				} else {
					subs.logs.PushBack(&l)
				}
				subs.Unlock()
			}
		}()
		state.subscriptions[req.RequestID] = subs
		fmt.Println("NEW SUBSCRIPTION")
	} else {
		fmt.Println("REUSING SUB")
	}

	res.Success = true
	res.SeqNum = subs.seqNum
	context.Respond(res)
	context.Watch(req.Subscriber)

	return nil
}

func (state *Executor) onFlushLogs(context actor.Context) error {
	ctx, canc := goContext.WithTimeout(goContext.Background(), 10*time.Second)
	defer canc()
	current, err := state.client.BlockNumber(ctx)
	if err != nil {
		state.logger.Warn("error fetching block number", log.Error(err))
		return nil
	}
	for k, sub := range state.subscriptions {
		sub.RLock()
		var logs *messages.EVMLogs
		for el := sub.logs.Front(); el != nil; el = sub.logs.Front() {
			updt := el.Value.(*types.Log)
			if updt.BlockNumber > current-3 {
				// we are done
				break
			}
			if logs == nil || logs.BlockNumber != updt.BlockNumber {
				// No update or new block, publish
				if logs != nil {
					context.Send(sub.subscriber, &messages.EVMLogsSubscribeRefresh{
						RequestID: k,
						SeqNum:    sub.seqNum + 1,
						Update:    logs,
					})
					sub.seqNum += 1
				}
				ctx, cancel := goContext.WithTimeout(goContext.Background(), 10*time.Second)
				var head *zks.Header
				num := hexutil.EncodeBig(big.NewInt(int64(updt.BlockNumber)))
				err := state.rpc.CallContext(ctx, &head, "eth_getBlockByNumber", num, false)
				if err != nil {
					state.logger.Warn("error getting block header", log.Error(err))
					sub.RUnlock()
					cancel()
					return nil
				}
				logs = &messages.EVMLogs{
					BlockNumber: updt.BlockNumber,
					BlockTime:   time.Unix(int64(head.Time), 0),
				}
				cancel()
			}
			logs.Logs = append(logs.Logs, *updt)
			sub.logs.Remove(el)
		}
		// Publish it
		if logs != nil {
			context.Send(sub.subscriber, &messages.EVMLogsSubscribeRefresh{
				RequestID: k,
				SeqNum:    sub.seqNum + 1,
				Update:    logs,
			})
			sub.seqNum += 1
			sub.lastPingTime = time.Now()
		} else if time.Since(sub.lastPingTime) > 10*time.Second {
			context.Send(sub.subscriber, &messages.EVMLogsSubscribeRefresh{
				RequestID: k,
				SeqNum:    sub.seqNum + 1,
			})
			sub.seqNum += 1
			sub.lastPingTime = time.Now()
		}
		sub.RUnlock()
	}
	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	for _, sub := range state.subscriptions {
		close(sub.ch)
		context.Unwatch(sub.subscriber)
	}
	state.subscriptions = nil
	if state.flushTicker == nil {
		state.flushTicker.Stop()
		state.flushTicker = nil
	}
	return nil
}

func (state *Executor) OnTerminated(context actor.Context) error {
	msg := context.Message().(*actor.Terminated)
	for k, sub := range state.subscriptions {
		if sub.subscriber.String() == msg.Who.String() {
			close(sub.ch)
			delete(state.subscriptions, k)
		}
	}
	return nil
}

func (state *Executor) GetLogger() *log.Logger {
	return state.logger
}
