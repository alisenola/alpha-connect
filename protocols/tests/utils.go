package tests

import (
	"fmt"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

type ProtocolChecker struct {
	asset   *models.ProtocolAsset
	logger  *log.Logger
	updates []*models.ProtocolAssetUpdate
	err     error
	seqNum  uint64
}

func NewProtocolCheckerProducer(asset *models.ProtocolAsset) actor.Producer {
	return func() actor.Actor {
		return NewProtocolChecker(asset)
	}
}

func NewProtocolChecker(asset *models.ProtocolAsset) actor.Actor {
	return &ProtocolChecker{
		asset:   asset,
		updates: nil,
		logger:  nil,
		seqNum:  0,
	}
}

func (state *ProtocolChecker) Receive(context actor.Context) {
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

func (state *ProtocolChecker) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(state).String()),
	)
	executor := context.ActorSystem().NewLocalPID("executor")
	req := &messages.ProtocolAssetDataRequest{
		RequestID:  0,
		ChainID:    state.asset.Chain.ID,
		ProtocolID: state.asset.Protocol.ID,
		Subscriber: context.Self(),
		Subscribe:  true,
	}
	if state.asset.Asset != nil {
		req.AssetID = &wrapperspb.UInt32Value{Value: state.asset.Asset.ID}
	}
	res, err := context.RequestFuture(executor, req, 30*time.Second).Result()
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

func (state *ProtocolChecker) OnProtocolAssetDataIncrementalRefresh(context actor.Context) error {
	if state.err != nil {
		return nil
	}
	res := context.Message().(*messages.ProtocolAssetDataIncrementalRefresh)
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
