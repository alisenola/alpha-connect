package erc721

import (
	goContext "context"
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"time"

	"gitlab.com/alphaticks/alpha-connect/utils"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
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
	protocolAssets map[uint64]*models.ProtocolAsset
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

	client, err := ethclient.Dial(xutils.ETH_CLIENT_WS_2)
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
	if state.registry == nil {
		return nil
	}
	assets := make([]*models.ProtocolAsset, 0)

	reg := state.registry

	ctx, cancel := goContext.WithTimeout(goContext.Background(), 10*time.Second)
	defer cancel()
	filter := registry.ProtocolAssetFilter{
		ProtocolId: []uint32{constants.ERC721.ID},
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
				Protocol: &models2.Protocol{
					ID:   constants.ERC721.ID,
					Name: "ERC-721",
				},
				Asset: &models2.Asset{
					Name:   as.Name,
					Symbol: as.Symbol,
					ID:     as.ID,
				},
				Chain: &models2.Chain{
					ID:   ch.ID,
					Name: ch.Name,
					Type: ch.Type,
				},
				CreationDate:    protocolAsset.CreationDate,
				CreationBlock:   protocolAsset.CreationBlock,
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

func (state *Executor) OnHistoricalProtocolAssetTransferRequest(context actor.Context) error {
	req := context.Message().(*messages.HistoricalProtocolAssetTransferRequest)
	msg := &messages.HistoricalProtocolAssetTransferResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	pa, ok := state.protocolAssets[req.ProtocolAssetID]
	if !ok {
		msg.RejectionReason = messages.RejectionReason_UnknownProtocolAsset
		context.Respond(msg)
		return nil
	}

	eabi, err := nft.ERC721MetaData.GetAbi()
	if err != nil {
		state.logger.Warn("error getting erc721 ABI", log.Error(err))
		msg.RejectionReason = messages.RejectionReason_ABIError
		context.Respond(msg)
		return nil
	}

	topics := [][]common.Hash{{
		eabi.Events["Transfer"].ID,
	}}
	var address [20]byte
	addressBig, ok := big.NewInt(1).SetString(pa.ContractAddress.Value[2:], 16)
	if !ok {
		state.logger.Warn("invalid protocol asset address", log.Error(err))
		msg.RejectionReason = messages.RejectionReason_UnknownProtocolAsset
		context.Respond(msg)
		return nil
	}

	copy(address[:], addressBig.Bytes())
	fQuery := ethereum.FilterQuery{
		Addresses: []common.Address{address},
		FromBlock: big.NewInt(1).SetUint64(req.Start),
		ToBlock:   big.NewInt(1).SetUint64(req.Stop),
		Topics:    topics,
	}
	qr := state.queryRunnerETH
	future := context.RequestFuture(qr.pid, &jobs.PerformLogsQueryRequest{Query: fQuery}, 15*time.Second)
	context.ReenterAfter(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Warn("error at eth rpc server", log.Error(err))
			switch err.Error() {
			case "future: timeout":
				msg.RejectionReason = messages.RejectionReason_EthRPCTimeout
			default:
				msg.RejectionReason = messages.RejectionReason_EthRPCError
			}
			context.Respond(msg)
			return
		}

		resp := res.(*jobs.PerformLogsQueryResponse)
		if resp.Error != nil {
			state.logger.Warn("error at eth rpc server", log.Error(resp.Error))
			msg.RejectionReason = messages.RejectionReason_EthRPCError
			context.Respond(msg)
			return
		}

		events := make([]*models.ProtocolAssetUpdate, 0)
		sort.Slice(resp.Logs, func(i, j int) bool {
			if resp.Logs[i].BlockNumber == resp.Logs[j].BlockNumber {
				return resp.Logs[i].Index < resp.Logs[j].Index
			}
			return resp.Logs[i].BlockNumber < resp.Logs[j].BlockNumber
		})
		var update *models.ProtocolAssetUpdate
		for i, l := range resp.Logs {
			switch l.Topics[0] {
			case eabi.Events["Transfer"].ID:
				event := nft.ERC721Transfer{}
				if err := eth.UnpackLog(eabi, &event, "Transfer", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
					msg.RejectionReason = messages.RejectionReason_EthRPCError
					context.Respond(msg)
					return
				}
				if update == nil || update.Block != l.BlockNumber {
					if update != nil {
						events = append(events, update)
					}
					update = &models.ProtocolAssetUpdate{
						Transfers: nil,
						Block:     l.BlockNumber,
						Timestamp: utils.SecondToTimestamp(resp.Times[i]),
					}
				}
				update.Transfers = append(update.Transfers, &gorderbook.AssetTransfer{
					From:    event.From[:],
					To:      event.To[:],
					TokenId: event.TokenId.Bytes(),
				})
			}
		}
		if update != nil {
			events = append(events, update)
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
