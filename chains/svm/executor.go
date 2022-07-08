package svm

import (
	goContext "context"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	chtypes "gitlab.com/alphaticks/alpha-connect/chains/types"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	"gitlab.com/alphaticks/xchanger/chains/svm"
	"math/big"
	"reflect"
	"time"
)

type Executor struct {
	chtypes.BaseExecutor
	logger *log.Logger
	client *svm.Client
	rpc    string
}

func NewExecutor(registry registry.PublicRegistryClient, rpc string) actor.Actor {
	e := Executor{
		rpc:    rpc,
		logger: nil,
		client: nil,
	}
	e.Registry = registry
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

	cl, err := svm.Dial(state.rpc)
	if err != nil {
		return fmt.Errorf("error dialing rpc url: %v", err)
	}
	state.client = cl
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
		done := false
		q := req.Query
		var events []*svm.Event
		for !done {
			evs, err := state.client.GetEvents(ctx, q)
			if err != nil {
				state.logger.Error("error getting svm events", log.Error(err))
				msg.RejectionReason = messages.RejectionReason_RPCError
				context.Send(pid, msg)
				return
			}
			events = append(events, evs.Events...)
			q.PageNumber += 1
			done = evs.IsLastPage
		}
		lastBlock := 0
		var times []uint64
		var bl *svm.Block
		for _, ev := range events {
			if ev.BlockNumber != lastBlock {
				ctx, cancel := goContext.WithTimeout(goContext.Background(), 5*time.Second)
				qB := svm.BlockQuery{
					BlockNumber: big.NewInt(int64(ev.BlockNumber)),
				}
				var err error
				bl, err = state.client.BlockInfo(ctx, qB)
				if err != nil {
					state.logger.Error("error getting svm block times", log.Error(err))
					msg.RejectionReason = messages.RejectionReason_RPCError
					context.Send(pid, msg)
					cancel()
					return
				}
				cancel()
			}
			lastBlock = ev.BlockNumber
			times = append(times, bl.AcceptedTime)
		}
		msg.Events = events
		msg.Times = times
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

func (state *Executor) OnSVMTransactionByHashRequest(context actor.Context) error {
	req := context.Message().(*messages.SVMTransactionByHashRequest)
	msg := &messages.SVMTransactionByHashResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	go func(pid *actor.PID) {
		ctx, cancel := goContext.WithTimeout(goContext.Background(), 20*time.Second)
		defer cancel()
		tx, err := state.client.TransactionReceipt(ctx, req.Hash)
		if err != nil {
			state.logger.Error("error getting svm transaction by hash", log.Error(err))
			msg.RejectionReason = messages.RejectionReason_RPCError
			context.Send(pid, msg)
			return
		}
		msg.Transaction = tx
		msg.Success = true
		context.Send(pid, msg)
	}(context.Sender())
	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	if state.client != nil {
		state.client.Close()
		state.client = nil
	}
	return nil
}

func (state *Executor) GetLogger() *log.Logger {
	return state.logger
}
