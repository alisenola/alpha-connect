package svm

import (
	goContext "context"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	chtypes "gitlab.com/alphaticks/alpha-connect/chains/types"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/chains/starknet"
	"reflect"
	"time"
)

type clientCkeck struct{}

type EventsSubscription struct {
	subscriber   *actor.PID
	query        *starknet.EventQuery
	seqNum       uint64
	lastPingTime time.Time
}

type Executor struct {
	chtypes.BaseExecutor
	logger        *log.Logger
	client        *starknet.Client
	subscriptions map[uint64]*EventsSubscription
	clientTicker  *time.Ticker
	rpc           string
}

func NewExecutor(config *chtypes.ExecutorConfig, rpc string) actor.Actor {
	e := Executor{
		rpc:          rpc,
		logger:       nil,
		client:       nil,
		clientTicker: nil,
	}
	e.ExecutorConfig = config
	return &e
}

func (state *Executor) Receive(context actor.Context) {
	switch context.Message().(type) {
	default:
		chtypes.ReceiveExecutor(state, context)
	}
}

func (state *Executor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("Type", reflect.TypeOf(*state).String()))

	cl, err := starknet.Dial(state.rpc)
	if err != nil {
		return fmt.Errorf("error dialing rpc url: %v", err)
	}
	state.client = cl

	state.subscriptions = make(map[uint64]*EventsSubscription)

	ticker := time.NewTicker(20 * time.Second)
	state.clientTicker = ticker
	go func(pid *actor.PID) {
		for {
			select {
			case <-ticker.C:
				context.Send(pid, &clientCkeck{})
			case <-time.After(40 * time.Second):
				if state.clientTicker != ticker {
					return
				}
			}
		}
	}(context.Self())
	return nil
}

func (state *Executor) OnBlockNumberRequest(context actor.Context) error {
	req := context.Message().(*messages.BlockNumberRequest)
	msg := &messages.BlockNumberResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	go func(pid *actor.PID) {
		ctx, cancel := goContext.WithTimeout(goContext.Background(), 20*time.Second)
		defer cancel()
		num, err := state.client.BlockNumber(ctx)
		if err != nil {
			state.logger.Error("error getting block number", log.Error(err))
			msg.RejectionReason = messages.RejectionReason_RPCError
			context.Send(pid, msg)
			return
		}
		msg.BlockNumber = num
		msg.Success = true
		context.Send(pid, msg)
	}(context.Sender())
	return nil
}

func (state *Executor) OnSVMEventsQueryRequest(context actor.Context) error {
	req := context.Message().(*messages.SVMEventsQueryRequest)
	msg := &messages.SVMEventsQueryResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	go func(pid *actor.PID) {
		ctx, cancel := goContext.WithTimeout(goContext.Background(), 30*time.Second)
		defer cancel()
		events, err := state.client.GetEvents(ctx, *req.Query)
		if err != nil {
			state.logger.Error("error getting svm events", log.Error(err))
			msg.RejectionReason = messages.RejectionReason_RPCError
			context.Send(pid, msg)
			return
		}
		msg.Events = events
		msg.Success = true
		context.Send(pid, msg)
	}(context.Sender())
	return nil
}

func (state *Executor) OnSVMBlockQueryRequest(context actor.Context) error {
	req := context.Message().(*messages.SVMBlockQueryRequest)
	msg := &messages.SVMBlockQueryResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	go func(pid *actor.PID) {
		ctx, cancel := goContext.WithTimeout(goContext.Background(), 20*time.Second)
		defer cancel()
		bl, err := state.client.BlockInfo(ctx, *req.Query)
		if err != nil {
			state.logger.Error("error getting svm block", log.Error(err))
			msg.RejectionReason = messages.RejectionReason_RPCError
			context.Send(pid, msg)
			return
		}
		msg.Block = bl
		msg.Success = true
		context.Send(pid, msg)
	}(context.Sender())
	return nil
}

func (state *Executor) OnTick(context actor.Context) error {
	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	if state.clientTicker != nil {
		state.clientTicker.Stop()
		state.clientTicker = nil
	}
	if state.client != nil {
		state.client.Close()
		state.client = nil
	}
	return nil
}

func (state *Executor) GetLogger() *log.Logger {
	return state.logger
}
