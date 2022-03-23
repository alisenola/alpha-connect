package erc721

import (
	goContext "context"
	"fmt"
	"math/big"
	"reflect"
	"time"

	"gitlab.com/alphaticks/alpha-connect/utils"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	gorderbook "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	"gitlab.com/alphaticks/xchanger/constants"
	models2 "gitlab.com/alphaticks/xchanger/models"
	xutils "gitlab.com/alphaticks/xchanger/protocols"

	"gitlab.com/alphaticks/alpha-connect/jobs"
	extype "gitlab.com/alphaticks/alpha-connect/protocols/types"
	"gitlab.com/alphaticks/xchanger/eth"
	nft "gitlab.com/alphaticks/xchanger/protocols/erc721"
)

type QueryRunner struct {
	pid *actor.PID
}

type Executor struct {
	extype.BaseExecutor
	queryRunnerETH *QueryRunner
	protocolAssets []*models.ProtocolAsset
	logger         *log.Logger
	registry       registry.PublicRegistryClient
}

func NewExecutor(registry registry.PublicRegistryClient) actor.Actor {
	return &Executor{
		queryRunnerETH: nil,
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

	client, err := ethclient.Dial(xutils.ETH_CLIENT_WS)
	if err != nil {
		return fmt.Errorf("error while dialing eth rpc client %v", err)
	}
	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewETHQuery(client)
	})
	state.queryRunnerETH = &QueryRunner{
		pid: context.Spawn(props),
	}
	return state.UpdateProtocolAssetList(context)
}

func (state *Executor) UpdateProtocolAssetList(context actor.Context) error {
	assets := make([]*models.ProtocolAsset, 0)
	reg := state.registry

	ctx, cancel := goContext.WithTimeout(goContext.Background(), 10*time.Second)
	defer cancel()
	filter := registry.ProtocolAssetFilter{}
	in := registry.ProtocolAssetsRequest{
		Filter: &filter,
	}
	res, err := reg.ProtocolAssets(ctx, &in)
	if err != nil {
		return fmt.Errorf("error updating protocol asset list: %v", err)
	}
	response := res.ProtocolAssets
	for _, protocolAsset := range response {
		address, ok := big.NewInt(1).SetString(protocolAsset.Meta["address"][2:], 16)
		if !ok {
			state.logger.Warn("incorrect address parsing for asset")
			continue
		}
		as, ok := constants.GetAssetByID(protocolAsset.AssetId)
		if !ok {
			state.logger.Warn("error getting asset with id", log.String("asset-id", fmt.Sprint(protocolAsset.AssetId)))
			continue
		}
		assets = append(
			assets,
			&models.ProtocolAsset{
				Address: address.Bytes(),
				Protocol: &models2.Protocol{
					ID:   constants.ERC721.ID,
					Name: "ERC-721",
				},
				Asset: &models2.Asset{
					Name:   as.Name,
					Symbol: as.Symbol,
					ID:     as.ID,
				},
				Meta: protocolAsset.Meta,
			})
	}
	state.protocolAssets = assets

	context.Send(context.Parent(), &messages.ProtocolAssetList{
		ResponseID:     uint64(time.Now().UnixNano()),
		ProtocolAssets: assets,
		Success:        true,
	})

	return nil
}

func (state *Executor) OnProtocolAssetListRequest(context actor.Context) error {
	req := context.Message().(*messages.ProtocolAssetListRequest)
	context.Respond(&messages.ProtocolAssetListResponse{
		RequestID:      req.RequestID,
		ResponseID:     uint64(time.Now().UnixNano()),
		Success:        true,
		ProtocolAssets: state.protocolAssets,
	})
	return nil
}

func (state *Executor) OnHistoricalProtocolAssetTransferRequest(context actor.Context) error {
	req := context.Message().(*messages.HistoricalProtocolAssetTransferRequest)
	msg := &messages.HistoricalProtocolAssetTransferResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	if req.ProtocolAsset == nil || req.ProtocolAsset.Address == nil {
		msg.RejectionReason = messages.MissingProtocolAsset
		context.Respond(msg)
		return nil
	}

	eabi, err := nft.ERC721MetaData.GetAbi()
	if err != nil {
		state.logger.Warn("error getting erc721 ABI", log.Error(err))
		msg.RejectionReason = messages.ABIError
		context.Respond(msg)
		return nil
	}

	topics := [][]common.Hash{{
		eabi.Events["Transfer"].ID,
	}}
	var address [20]byte
	copy(address[:], req.ProtocolAsset.Address)
	fQuery := ethereum.FilterQuery{
		Addresses: []common.Address{address},
		FromBlock: big.NewInt(1).SetUint64(req.Start),
		ToBlock:   big.NewInt(1).SetUint64(req.Stop),
		Topics:    topics,
	}
	qr := state.queryRunnerETH
	future := context.RequestFuture(qr.pid, &jobs.PerformLogsQueryRequest{Query: fQuery}, 15*time.Second)
	context.AwaitFuture(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Warn("error at eth rpc server", log.Error(err))
			msg.RejectionReason = messages.EthRPCError
			context.Respond(msg)
			return
		}

		resp := res.(*jobs.PerformLogsQueryResponse)
		if resp.Error != nil {
			state.logger.Warn("error at eth rpc server", log.Error(err))
			msg.RejectionReason = messages.EthRPCError
			context.Respond(msg)
			return
		}

		events := make([]*models.ProtocolAssetUpdate, 0)
		for i, l := range resp.Logs {
			switch l.Topics[0] {
			case eabi.Events["Transfer"].ID:
				event := nft.ERC721Transfer{}
				if err := eth.UnpackLog(eabi, &event, "Transfer", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
					msg.RejectionReason = messages.EthRPCError
					context.Respond(msg)
					return
				}
				t := &models.ProtocolAssetUpdate{
					Transfer: &gorderbook.AssetTransfer{
						From:    event.From[:],
						To:      event.To[:],
						TokenId: event.TokenId.Bytes(),
					},
					Removed:   l.Removed,
					Block:     l.BlockNumber,
					Timestamp: utils.SecondToTimestamp(resp.Times[i]),
				}
				events = append(events, t)
			}
		}
		msg.Update = events
		msg.Success = true
		msg.SeqNum = uint64(time.Now().UnixNano())
		context.Respond(msg)

	})

	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) GetLogger() *log.Logger {
	return state.logger
}
