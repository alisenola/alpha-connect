package messages

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/alphaticks/xchanger/chains/svm"
	"gitlab.com/alphaticks/xchanger/models"
	"time"
)

type BlockNumberRequest struct {
	RequestID uint64
	Chain     *models.Chain
}

type BlockNumberResponse struct {
	RequestID       uint64
	ResponseID      uint64
	Success         bool
	RejectionReason RejectionReason
	BlockNumber     uint64
}

type BlockInfoRequest struct {
	RequestID   uint64
	Chain       *models.Chain
	BlockNumber uint64
}

type BlockInfoResponse struct {
	RequestID       uint64
	ResponseID      uint64
	Success         bool
	RejectionReason RejectionReason
	BlockTime       time.Time
}

type EVMContractCallRequest struct {
	RequestID   uint64
	Chain       *models.Chain
	Msg         ethereum.CallMsg
	BlockNumber uint64
}

type EVMContractCallResponse struct {
	RequestID       uint64
	ResponseID      uint64
	Out             []byte
	Success         bool
	RejectionReason RejectionReason
}

type EVMLogsQueryRequest struct {
	RequestID uint64
	Chain     *models.Chain
	Query     ethereum.FilterQuery
}

type EVMLogsQueryResponse struct {
	RequestID       uint64
	ResponseID      uint64
	Success         bool
	RejectionReason RejectionReason
	Logs            []types.Log
	Times           []uint64
}

type EVMLogsSubscribeRequest struct {
	RequestID  uint64
	Chain      *models.Chain
	Query      ethereum.FilterQuery
	Subscriber *actor.PID
}

type EVMLogsSubscribeResponse struct {
	RequestID       uint64
	ResponseID      uint64
	Success         bool
	RejectionReason RejectionReason
	SeqNum          uint64
}

type EVMLogsSubscribeRefresh struct {
	RequestID uint64
	SeqNum    uint64
	Update    *EVMLogs
}

type EVMLogs struct {
	BlockNumber uint64
	BlockTime   time.Time
	Logs        []types.Log
}

type SVMEventsQueryRequest struct {
	RequestID uint64
	Query     svm.EventQuery
	Chain     *models.Chain
}

type SVMEventsQueryResponse struct {
	RequestID       uint64
	ResponseID      uint64
	Events          []*svm.Event
	Times           []uint64
	Success         bool
	RejectionReason RejectionReason
}

type SVMContractCallRequest struct {
	RequestID      uint64
	Chain          *models.Chain
	CallParameters *svm.CallParameters
}

type SVMContractCallResponse struct {
	RequestID       uint64
	ResponseID      uint64
	Data            []*svm.Hash
	Success         bool
	RejectionReason RejectionReason
}

type SVMContractClassRequest struct {
	RequestID uint64
	Chain     *models.Chain
	Query     *svm.ContractClassQuery
}

type SVMContractClassResponse struct {
	RequestID       uint64
	ResponseID      uint64
	Class           *svm.Class
	Success         bool
	RejectionReason RejectionReason
}

type SVMBlockQueryRequest struct {
	RequestID uint64
	Query     *svm.BlockQuery
	Chain     *models.Chain
}

type SVMBlockQueryResponse struct {
	RequestID       uint64
	ResponseID      uint64
	Block           *svm.Block
	Success         bool
	RejectionReason RejectionReason
}

type SVMTransactionByHashRequest struct {
	RequestID uint64
	Hash      common.Hash
	Chain     *models.Chain
}

type SVMTransactionByHashResponse struct {
	RequestID       uint64
	ResponseID      uint64
	Transaction     *svm.Transaction
	Success         bool
	RejectionReason RejectionReason
}
