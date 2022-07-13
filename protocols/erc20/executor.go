package erc20

import (
	goContext "context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	sabi "gitlab.com/alphaticks/abigen-starknet/accounts/abi"
	"gitlab.com/alphaticks/alpha-connect/utils"
	gorderbook "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	"gitlab.com/alphaticks/xchanger/chains/svm"
	token "gitlab.com/alphaticks/xchanger/protocols/erc20/svm"
	"math/big"
	"reflect"
	"sort"
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
	executor       *actor.PID
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

	stabi, err := token.ERC20MetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("error getting starknet abi: %v", err)
	}
	state.sabi = stabi
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
	var future *actor.Future
	var query interface{}
	var delay int64
	switch pa.Chain.Type {
	case "EVM":
		return nil
	case "SVM":
		add := common.HexToHash(pa.ContractAddress.Value)
		q := svm.EventQuery{
			ContractAddress: &add,
			Keys:            &[]common.Hash{state.sabi.Events["Transfer"].ID},
			From:            big.NewInt(1).SetUint64(req.Start),
			To:              big.NewInt(1).SetUint64(req.Stop),
			PageSize:        1000,
			PageNumber:      0,
		}
		query = &messages.SVMEventsQueryRequest{
			RequestID: uint64(time.Now().UnixNano()),
			Chain:     pa.Chain,
			Query:     q,
		}
		delay = 120
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
			sort.Slice(resp.Logs, func(i, j int) bool {
				if resp.Logs[i].BlockNumber == resp.Logs[j].BlockNumber {
					return resp.Logs[i].Index < resp.Logs[j].Index
				}
				return resp.Logs[i].BlockNumber < resp.Logs[j].BlockNumber
			})
			for _, l := range resp.Logs {
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
			sort.SliceStable(resp.Events, func(i, j int) bool {
				return resp.Events[i].BlockNumber < resp.Events[j].BlockNumber
			})
			for _, l := range resp.Events {
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
			var num uint64
			switch l := d.(type) {
			case *svm.Event:
				key := [32]byte(l.Keys[0])
				id := [32]byte(state.sabi.Events["Transfer"].ID)
				switch key {
				case id:
					event := token.ERC20Transfer{}
					if err := svm.UnpackLog(state.sabi, &event, "Transfer", l); err != nil {
						state.logger.Warn("error unpacking stark event", log.Error(err))
						msg.RejectionReason = messages.RejectionReason_RPCError
						context.Respond(msg)
						return
					}
					event.Sender.FillBytes(from[:])
					event.Recipient.FillBytes(to[:])
					event.Value.FillBytes(tok[:])
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
				From:  from,
				To:    to,
				Value: tok,
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
