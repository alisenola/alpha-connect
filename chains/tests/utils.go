package tests

import (
	"fmt"
	"github.com/ethereum/go-ethereum"
	"gitlab.com/alphaticks/xchanger/models"
	"reflect"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"

	"gitlab.com/alphaticks/alpha-connect/models/messages"
)

type GetDataRequest struct{}
type GetDataResponse struct {
	Updates []*messages.EVMLogs
	Err     error
}

type ChainChecker struct {
	logger  *log.Logger
	chain   *models.Chain
	query   ethereum.FilterQuery
	updates []*messages.EVMLogs
	err     error
	seqNum  uint64
}

func NewChainCheckerProducer(chain *models.Chain, q ethereum.FilterQuery) actor.Producer {
	return func() actor.Actor {
		return NewChainChecker(chain, q)
	}
}

func NewChainChecker(chain *models.Chain, q ethereum.FilterQuery) actor.Actor {
	return &ChainChecker{
		query:   q,
		chain:   chain,
		updates: nil,
		logger:  nil,
		seqNum:  0,
	}
}

func (state *ChainChecker) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.logger.Error("error starting the actor", log.Error(err))
			state.err = err
		}
	case *actor.Stopping:
		state.logger.Info("actor stopping")
	case *actor.Stopped:
		state.logger.Info("actor stopped")

	case *messages.EVMLogsSubscribeRefresh:
		if err := state.OnEVMLogsSubscribeRefresh(context); err != nil {
			state.logger.Error("error processing OnEVMLogsSubscribeRefresh", log.Error(err))
			state.err = err
		}
	case *GetDataRequest:
		context.Respond(&GetDataResponse{
			Updates: state.updates,
			Err:     state.err,
		})
	}
}

func (state *ChainChecker) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(state).String()),
	)
	executor := context.ActorSystem().NewLocalPID("executor")
	res, err := context.RequestFuture(executor, &messages.EVMLogsSubscribeRequest{
		RequestID:  uint64(time.Now().UnixNano()),
		Chain:      state.chain,
		Query:      state.query,
		Subscriber: context.Self(),
	}, 30*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error subscribing to evm logs for chain %s: %v", state.chain.Name, err)
	}
	updt, ok := res.(*messages.EVMLogsSubscribeResponse)
	if !ok {
		return fmt.Errorf("expected EVMLogsSubscribeResponse, got %s", reflect.TypeOf(res).String())
	}
	if !updt.Success {
		return fmt.Errorf("error on EVMLogsSubscribeResponse, got %s", updt.RejectionReason.String())
	}
	state.seqNum = updt.SeqNum
	return nil
}

func (state *ChainChecker) OnEVMLogsSubscribeRefresh(context actor.Context) error {
	if state.err != nil {
		return nil
	}
	res := context.Message().(*messages.EVMLogsSubscribeRefresh)
	if res.SeqNum <= state.seqNum {
		return nil
	}
	if res.SeqNum != state.seqNum+1 {
		return fmt.Errorf("seqNum not in sequence, expected %d, got %d", state.seqNum+1, res.SeqNum)
	}
	if res.Update != nil {
		state.updates = append(state.updates, res.Update)
	}
	state.seqNum = res.SeqNum
	return nil
}
