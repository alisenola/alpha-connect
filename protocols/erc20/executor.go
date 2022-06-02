package erc20

import (
	goContext "context"
	"fmt"
	"math/big"
	"reflect"
	"time"

	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	extype "gitlab.com/alphaticks/alpha-connect/protocols/types"
	"gitlab.com/alphaticks/xchanger/constants"
)

type QueryRunner struct {
	pid *actor.PID
}

type Executor struct {
	extype.BaseExecutor
	protocolAssets map[uint64]*models.ProtocolAsset
	logger         *log.Logger
	registry       registry.PublicRegistryClient
}

func NewExecutor(registry registry.PublicRegistryClient) actor.Actor {
	return &Executor{
		protocolAssets: nil,
		logger:         nil,
		registry:       registry,
	}
}

func (state *Executor) Receive(context actor.Context) {
	extype.ReceiveExecutor(state, context)
}

func (state *Executor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	return state.UpdateProtocolAssetList(context)
}

func (state *Executor) UpdateProtocolAssetList(context actor.Context) error {
	if state.registry == nil {
		return nil
	}
	assets := make([]*models.ProtocolAsset, 0)
	reg := state.registry

	ctx, cancel := goContext.WithTimeout(goContext.Background(), 10*time.Second)
	defer cancel()
	filter := registry.ProtocolAssetFilter{
		ProtocolId: []uint32{constants.ERC20.ID},
	}
	in := registry.ProtocolAssetsRequest{
		Filter: &filter,
	}
	res, err := reg.ProtocolAssets(ctx, &in)
	if err != nil {
		return fmt.Errorf("error updating protocol asset list: %v", err)
	}
	response := res.ProtocolAssets
	for _, protocolAsset := range response {
		if protocolAsset.ContractAddress == nil || len(protocolAsset.ContractAddress.Value) < 2 {
			state.logger.Warn("invalid protocol asset address")
			continue
		}
		_, ok := big.NewInt(1).SetString(protocolAsset.ContractAddress.Value[2:], 16)
		if !ok {
			state.logger.Warn("invalid protocol asset address", log.Error(err))
			continue
		}
		if protocolAsset.Decimals == nil {
			state.logger.Warn("unknown protocol asset decimals")
			continue
		}
		as, ok := constants.GetAssetByID(protocolAsset.AssetId)
		if !ok {
			state.logger.Warn(fmt.Sprintf("error getting asset with id %d", protocolAsset.AssetId))
			continue
		}
		ch, ok := constants.GetChainByID(protocolAsset.ChainId)
		if !ok {
			state.logger.Warn(fmt.Sprintf("error getting chain with id %d", protocolAsset.ChainId))
			continue
		}
		assets = append(
			assets,
			&models.ProtocolAsset{
				ProtocolAssetID: protocolAsset.ProtocolAssetId,
				Protocol:        constants.ERC20,
				Asset:           as,
				Chain:           ch,
				CreationBlock:   protocolAsset.CreationBlock,
				CreationDate:    protocolAsset.CreationDate,
				ContractAddress: protocolAsset.ContractAddress,
				Decimals:        protocolAsset.Decimals,
			})
	}
	state.protocolAssets = make(map[uint64]*models.ProtocolAsset)
	for _, a := range assets {
		state.protocolAssets[a.ProtocolAssetID] = a
	}

	context.Send(context.Parent(), &messages.ProtocolAssetList{
		ResponseID:     uint64(time.Now().UnixNano()),
		ProtocolAssets: assets,
		Success:        true,
	})

	return nil
}

func (state *Executor) OnProtocolAssetListRequest(context actor.Context) error {
	req := context.Message().(*messages.ProtocolAssetListRequest)
	passets := make([]*models.ProtocolAsset, len(state.protocolAssets))
	i := 0
	for _, v := range state.protocolAssets {
		passets[i] = v
		i += 1
	}
	context.Respond(&messages.ProtocolAssetList{
		RequestID:      req.RequestID,
		ResponseID:     uint64(time.Now().UnixNano()),
		Success:        true,
		ProtocolAssets: passets,
	})
	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) GetLogger() *log.Logger {
	return state.logger
}
