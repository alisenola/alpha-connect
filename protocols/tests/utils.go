package tests

import (
	"fmt"
	"reflect"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"

	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
)

type GetDataRequest struct{}
type GetDataResponse struct {
	Updates []*models.ProtocolAssetUpdate
	Err     error
}

type ERC721Checker struct {
	asset   *models.ProtocolAsset
	logger  *log.Logger
	updates []*models.ProtocolAssetUpdate
	err     error
	seqNum  uint64
}

func NewERC721CheckerProducer(asset *models.ProtocolAsset) actor.Producer {
	return func() actor.Actor {
		return NewERC721Checker(asset)
	}
}

func NewERC721Checker(asset *models.ProtocolAsset) actor.Actor {
	return &ERC721Checker{
		asset:   asset,
		updates: nil,
		logger:  nil,
		seqNum:  0,
	}
}

func (state *ERC721Checker) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.logger.Error("error starting the actor", log.Error(err))
			state.err = err
		}
	case *messages.ProtocolAssetDataIncrementalRefresh:
		if err := state.OnProtocolAssetDataIncrementalRefresh(context); err != nil {
			state.logger.Error("error processing ProtocolAssetDataIncrementalRefresh", log.Error(err))
			state.err = err
		}
	case *GetDataRequest:
		context.Respond(&GetDataResponse{
			Updates: state.updates,
			Err:     state.err,
		})
	}
}

func (state *ERC721Checker) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(state).String()),
	)
	executor := context.ActorSystem().NewLocalPID("executor")
	res, err := context.RequestFuture(executor, &messages.ProtocolAssetDataRequest{
		RequestID:       0,
		ProtocolAssetID: state.asset.ProtocolAssetID,
		Subscriber:      context.Self(),
		Subscribe:       true,
	}, 30*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error fetching the asset transfer %v", err)
	}
	updt, ok := res.(*messages.ProtocolAssetDataResponse)
	if !ok {
		return fmt.Errorf("error for type assertion %v", err)
	}
	if !updt.Success {
		return fmt.Errorf("error on ProtocolAssetTransferResponse, got %s", updt.RejectionReason.String())
	}
	state.seqNum = updt.SeqNum
	return nil
}

func (state *ERC721Checker) OnProtocolAssetDataIncrementalRefresh(context actor.Context) error {
	if state.err != nil {
		return nil
	}
	res := context.Message().(*messages.ProtocolAssetDataIncrementalRefresh)
	if res.Update == nil {
		return nil
	}
	if res.SeqNum <= state.seqNum {
		return nil
	}
	if res.SeqNum != state.seqNum+1 {
		return fmt.Errorf("seqNum not in sequence, expected %d, got %d", state.seqNum+1, res.SeqNum)
	}
	state.seqNum = res.SeqNum
	return nil
}
