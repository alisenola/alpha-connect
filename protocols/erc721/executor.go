package erc721

import (
	goContext "context"
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	sabi "gitlab.com/alphaticks/abigen-starknet/accounts/abi"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	extype "gitlab.com/alphaticks/alpha-connect/protocols/types"
	"gitlab.com/alphaticks/alpha-connect/utils"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	gorderbook "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	"gitlab.com/alphaticks/xchanger/chains/evm"
	"gitlab.com/alphaticks/xchanger/chains/svm"
	"gitlab.com/alphaticks/xchanger/constants"
	models2 "gitlab.com/alphaticks/xchanger/models"
	nft "gitlab.com/alphaticks/xchanger/protocols/erc721/evm"
	snft "gitlab.com/alphaticks/xchanger/protocols/erc721/svm"
)

type QueryRunner struct {
	pid *actor.PID
}

type Executor struct {
	extype.BaseExecutor
	addresses      map[string]bool
	key            *big.Int
	executor       *actor.PID
	protocolAssets map[uint64]*models.ProtocolAsset
	eabi           *abi.ABI
	sabi           *sabi.ABI
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

	state.executor = actor.NewPID(context.ActorSystem().Address(), "executor")

	eabi, err := nft.ERC721MetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("error getting ethereum abi: %v", err)
	}
	state.eabi = eabi

	stabi, err := snft.ERC721MetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("error getting starknet abi: %v", err)
	}
	state.sabi = stabi

	state.addresses = make(map[string]bool)
	state.key = svm.GetSelectorFromName("ownerOf")

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
	var pa *models.ProtocolAsset
	chain, ok := constants.GetChainByID(req.ChainID)
	if !ok {
		msg.RejectionReason = messages.RejectionReason_UnknownChain
		context.Respond(msg)
		return nil
	}
	if req.AssetID != nil {
		var ok bool
		asset, ok := constants.GetAssetByID(req.AssetID.Value)
		if !ok {
			msg.RejectionReason = messages.RejectionReason_UnknownAsset
			context.Respond(msg)
			return nil
		}
		id := utils.GetProtocolAssetID(asset, constants.ERC721, chain)
		pa, ok = state.protocolAssets[id]
		if !ok {
			msg.RejectionReason = messages.RejectionReason_UnknownProtocolAsset
			context.Respond(msg)
			return nil
		}
	}

	var future *actor.Future
	var query interface{}
	var delay int64
	switch chain.Type {
	case "EVM":
		topics := [][]common.Hash{{
			state.eabi.Events["Transfer"].ID,
		}}
		fQuery := ethereum.FilterQuery{
			FromBlock: big.NewInt(1).SetUint64(req.Start),
			ToBlock:   big.NewInt(1).SetUint64(req.Stop),
			Topics:    topics,
		}
		r := &messages.EVMLogsQueryRequest{
			RequestID: uint64(time.Now().UnixNano()),
			Chain:     chain,
			Query:     fQuery,
		}
		if pa != nil {
			a := common.HexToAddress(pa.ContractAddress.Value)
			r.Query.Addresses = []common.Address{a}
		}
		delay = 20
		query = r
	case "SVM":
		q := svm.EventQuery{
			Keys:       &[]common.Hash{state.sabi.Events["Transfer"].ID},
			From:       big.NewInt(1).SetUint64(req.Start),
			To:         big.NewInt(1).SetUint64(req.Stop),
			PageSize:   1000,
			PageNumber: 0,
		}
		r := &messages.SVMEventsQueryRequest{
			RequestID: uint64(time.Now().UnixNano()),
			Chain:     chain,
			Query:     q,
		}
		if pa != nil {
			add := common.HexToHash(pa.ContractAddress.Value)
			r.Query.ContractAddress = &add
		}
		delay = 120
		query = r
	}
	future = context.RequestFuture(state.executor, query, time.Duration(delay)*time.Second)
	context.ReenterAfter(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Warn("error at rpc server", log.Error(err))
			switch err.Error() {
			case "future: timeout":
				msg.RejectionReason = messages.RejectionReason_RPCTimeout
			default:
				msg.RejectionReason = messages.RejectionReason_RPCError
			}
			context.Respond(msg)
			return
		}

		var datas []interface{}
		var times []uint64
		switch resp := res.(type) {
		case *messages.EVMLogsQueryResponse:
			if !resp.Success {
				state.logger.Warn("error at eth rpc server", log.String("rejection reason", resp.RejectionReason.String()))
				msg.RejectionReason = messages.RejectionReason_RPCError
				context.Respond(msg)
				return
			}
			if len(resp.Times) != len(resp.Logs) {
				state.logger.Warn("mismatched logs and times array length", log.Int("times length", len(resp.Times)), log.Int("events length", len(resp.Logs)))
				msg.RejectionReason = messages.RejectionReason_Other
				context.Respond(msg)
				return
			}
			var c []types.Log
			for _, l := range resp.Logs {
				if len(l.Topics) != 4 {
					continue
				}
				c = append(c, l)
			}
			sort.Slice(c, func(i, j int) bool {
				if c[i].BlockNumber == c[j].BlockNumber {
					return c[i].Index < c[j].Index
				}
				return c[i].BlockNumber < c[j].BlockNumber
			})
			for _, l := range c {
				datas = append(datas, l)
			}
			times = resp.Times
		case *messages.SVMEventsQueryResponse:
			if !resp.Success {
				state.logger.Warn("error at stark rpc server", log.String("rejection reason", resp.RejectionReason.String()))
				msg.RejectionReason = messages.RejectionReason_RPCError
				context.Respond(msg)
				return
			}
			if len(resp.Times) != len(resp.Events) {
				state.logger.Warn("mismatched events and times array length", log.Int("times length", len(resp.Times)), log.Int("events length", len(resp.Events)))
				msg.RejectionReason = messages.RejectionReason_Other
				context.Respond(msg)
				return
			}
			//filter out erc721 events
			var e []*svm.Event
			for _, l := range resp.Events {
				if v, ok := state.addresses[l.FromAddress.String()]; !ok {
					add := common.BytesToHash(l.FromAddress.Bytes())
					q := &svm.ContractClassQuery{
						ContractAddress: &add,
					}
					resp, err := context.RequestFuture(state.executor, &messages.SVMContractClassRequest{
						RequestID: uint64(time.Now().UnixNano()),
						Chain:     chain,
						Query:     q,
					}, 10*time.Second).Result()
					if err != nil {
						state.logger.Error("error at stark rpc server", log.Error(err))
						msg.RejectionReason = messages.RejectionReason_RPCError
						context.Respond(msg)
						return
					}
					c, ok := resp.(*messages.SVMContractClassResponse)
					if !ok {
						msg.RejectionReason = messages.RejectionReason_Other
						context.Respond(msg)
						return
					}
					if !c.Success {
						msg.RejectionReason = c.RejectionReason
						context.Respond(msg)
						return
					}
					state.addresses[l.FromAddress.String()] = false
					for _, ext := range c.Class.EntryPointsByType.External {
						i, _ := big.NewInt(0).SetString(ext.Selector[2:], 16)
						if state.key.Cmp(i) == 0 {
							e = append(e, l)
							state.addresses[l.FromAddress.String()] = true
							break
						}
					}
				} else if v {
					e = append(e, l)
				}
			}
			sort.SliceStable(e, func(i, j int) bool {
				return e[i].BlockNumber < e[j].BlockNumber
			})
			for _, l := range e {
				datas = append(datas, l)
			}
			times = resp.Times
		}
		events := make([]*models.ProtocolAssetUpdate, 0)
		var update *models.ProtocolAssetUpdate
		for i, d := range datas {
			from := make([]byte, 32)
			to := make([]byte, 32)
			tok := make([]byte, 32)
			contract := make([]byte, 32)
			var num uint64
			switch l := d.(type) {
			case types.Log:
				switch l.Topics[0] {
				case state.eabi.Events["Transfer"].ID:
					event := nft.ERC721Transfer{}
					if err := evm.UnpackLog(state.eabi, &event, "Transfer", l); err != nil {
						state.logger.Warn("error unpacking eth log", log.Error(err))
						msg.RejectionReason = messages.RejectionReason_RPCError
						context.Respond(msg)
						return
					}
					from = event.From[:]
					to = event.To[:]
					tok = event.TokenId.Bytes()
					contract = l.Address.Bytes()
					num = l.BlockNumber
				}
			case *svm.Event:
				key := [32]byte(l.Keys[0])
				id := [32]byte(state.sabi.Events["Transfer"].ID)
				switch key {
				case id:
					event := snft.ERC721Transfer{}
					if err := svm.UnpackLog(state.sabi, &event, "Transfer", l); err != nil {
						state.logger.Warn("error unpacking stark event", log.Error(err))
						msg.RejectionReason = messages.RejectionReason_RPCError
						context.Respond(msg)
						return
					}
					event.From.FillBytes(from[:])
					event.To.FillBytes(to[:])
					event.TokenId.FillBytes(tok[:])
					contract = l.FromAddress[:]
					num = uint64(l.BlockNumber)
				}
			}
			if update == nil || update.BlockNumber != num {
				if update != nil {
					events = append(events, update)
				}
				update = &models.ProtocolAssetUpdate{
					Transfers:   nil,
					BlockNumber: num,
					BlockTime:   utils.SecondToTimestamp(times[i]),
				}
			}
			update.Transfers = append(update.Transfers, &gorderbook.AssetTransfer{
				From:     from,
				To:       to,
				Value:    tok,
				Contract: contract,
			})
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
