package erc721

import (
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/protocols"
	gorderbook "gitlab.com/alphaticks/gorderbook/gorderbook.models"

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
	collection     []*models.Collection
	logger         *log.Logger
}

func NewExecutor() actor.Actor {
	return &Executor{
		queryRunnerETH: nil,
		collection:     nil,
		logger:         nil,
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

	client, err := ethclient.Dial(protocols.ETH_CLIENT_WS)
	if err != nil {
		return fmt.Errorf("error while dialing eth rpc client %v", err)
	}
	props := actor.PropsFromProducer(func() actor.Actor {
		return jobs.NewABIQuery(client)
	})
	state.queryRunnerETH = &QueryRunner{
		pid: context.Spawn(props),
	}
	return state.UpdateCollectionList(context)
}

func (state *Executor) OnHistoricalNftTransferDataRequest(context actor.Context) error {
	fmt.Println("STARTING HISTORICAL NFT TRANSFER REQUEST")
	req := context.Message().(*messages.HistoricalNftTransferDataRequest)
	msg := &messages.HistoricalNftTransferDataResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}

	if req.Collection == nil || req.Collection.Address == nil {
		msg.RejectionReason = messages.MissingInstrument
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
	copy(address[:], req.Collection.Address)
	fQuery := ethereum.FilterQuery{
		Addresses: []common.Address{address},
		FromBlock: big.NewInt(1).SetUint64(req.Start),
		ToBlock:   big.NewInt(1).SetUint64(req.Stop),
		Topics:    topics,
	}
	fmt.Println("request", fQuery)
	qr := state.queryRunnerETH
	future := context.RequestFuture(qr.pid, &jobs.PerformLogsQueryRequest{Query: fQuery}, 15*time.Second)
	context.AwaitFuture(future, func(res interface{}, err error) {
		fmt.Println("IN THE FUTURE REQUEST")
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

		events := make([]*models.NftTransfer, 0)
		for _, l := range resp.Logs {
			switch l.Topics[0] {
			case eabi.Events["Transfer"].ID:
				event := nft.ERC721Transfer{}
				if err := eth.UnpackLog(eabi, &event, "Transfer", l); err != nil {
					state.logger.Warn("error unpacking log", log.Error(err))
					msg.RejectionReason = messages.EthRPCError
					context.Respond(msg)
					return
				}
				t := &models.NftTransfer{
					Transfer: &gorderbook.NftTransfer{
						From:    event.From[:],
						To:      event.To[:],
						TokenId: event.TokenId.Bytes(),
					},
					Removed: l.Removed,
					Block:   l.BlockNumber,
				}
				events = append(events, t)
			}
		}
		msg.Transfers = events
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

func (state *Executor) UpdateCollectionList(context actor.Context) error {
	return nil
}
