package svm_test

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gitlab.com/alphaticks/alpha-connect/chains/tests"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/chains/starknet"
	models2 "gitlab.com/alphaticks/xchanger/models"
	"math/big"
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	as    *actor.ActorSystem
	ex    *actor.PID
	chain *models2.Chain
)

func TestMain(m *testing.M) {
	chain = &models2.Chain{
		ID:   5,
		Name: "Starknet Mainnet",
		Type: "SVM",
	}
	var err error
	var poison func()
	as, ex, poison, err = tests.StartExecutor()
	defer poison()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestExecutorBlockNumberRequest(t *testing.T) {
	resp, err := as.Root.RequestFuture(ex, &messages.BlockNumberRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Chain:     chain,
	}, 10*time.Second).Result()
	assert.Nil(t, err, "BlockNumberRequest err: %v", err)
	num, ok := resp.(*messages.BlockNumberResponse)
	assert.True(t, ok, "expected BlockNumberRequest, got %s", reflect.TypeOf(resp).String())
	assert.True(t, num.Success, "failed BlockNumberRequest, got %s", num.RejectionReason.String())
	assert.Greater(t, num.BlockNumber, uint64(2000), "expected block to be higher than 2000")
}

func TestExecutorSVMEventsQueryRequest(t *testing.T) {
	v := validator.New()
	add := common.HexToHash("0x0276069eb59afc97d3f6dcd18ec236cbb3e71611324bfc6e86a6aee0851a65bb")
	q := &starknet.EventQuery{
		To:              big.NewInt(2790),
		ContractAddress: &add,
		PageSize:        1000,
		PageNumber:      0,
	}
	done := false
	for !done {
		resp, err := as.Root.RequestFuture(ex, &messages.SVMEventsQueryRequest{
			RequestID: uint64(time.Now().UnixNano()),
			Query:     q,
			Chain:     chain,
		}, 10*time.Second).Result()
		assert.Nil(t, err, "SVMEventsQueryRequest err: %v", err)
		events, ok := resp.(*messages.SVMEventsQueryResponse)
		assert.True(t, ok, "expected SVMEventsQueryRequest, got %s", reflect.TypeOf(resp).String())
		assert.True(t, events.Success, "failed SVMEventsQueryRequest, got: %s", events.RejectionReason.String())
		assert.Equal(t, len(events.Events.Events), 6, "expected more than 5 events")

		err = v.Struct(events.Events)
		assert.Nil(t, err, "Validate struct err: %v", err)
		for _, ev := range events.Events.Events {
			err = v.Struct(ev)
			assert.Nil(t, err, "Validate struct err: %v", err)
		}
		q.PageNumber = uint64(events.Events.PageNumber + 1)
		done = events.Events.IsLastPage
	}
}

func TestExecutorSVMBlockQueryRequest(t *testing.T) {
	v := validator.New()
	hash := common.HexToHash("0x69b96255bf7cc630ba99292ca1dd34130829fa3486fe44bfbb7f7aa13a4da29")

	q := &starknet.BlockQuery{
		BlockHash: &hash,
		TxScope:   &starknet.FULL_TXN_AND_RECEIPTS,
	}
	resp, err := as.Root.RequestFuture(ex, &messages.SVMBlockQueryRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Query:     q,
		Chain:     chain,
	}, 10*time.Second).Result()
	assert.Nil(t, err, "SVMEventsQueryRequest err: %v", err)
	block, ok := resp.(*messages.SVMBlockQueryResponse)
	assert.True(t, ok, "expected SVMBlockQueryResponse, got %s", reflect.TypeOf(resp).String())
	assert.True(t, block.Success, "failed SVMBlockQueryResponse, got: %s", block.RejectionReason.String())
	assert.Equal(t, len(block.Block.Transactions), 265, "expected more than 5 events")

	err = v.Struct(block.Block)
	assert.Nil(t, err, "Validate struct err: %v", err)
	for _, tx := range block.Block.Transactions {
		err = v.Struct(tx)
		assert.Nil(t, err, "Validate struct err: %v", err)
	}

	blockN := big.NewInt(1828)
	q = &starknet.BlockQuery{
		BlockNumber: blockN,
	}
	resp, err = as.Root.RequestFuture(ex, &messages.SVMBlockQueryRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Query:     q,
		Chain:     chain,
	}, 10*time.Second).Result()
	assert.Nil(t, err, "SVMEventsQueryRequest err: %v", err)
	block, ok = resp.(*messages.SVMBlockQueryResponse)
	assert.True(t, ok, "expected SVMBlockQueryResponse, got %s", reflect.TypeOf(resp).String())
	assert.True(t, block.Success, "failed SVMBlockQueryResponse, got: %s", block.RejectionReason.String())
	assert.Equal(t, len(block.Block.Transactions), 92, "expected more than 5 events")

	err = v.Struct(block.Block)
	assert.Nil(t, err, "Validate struct err: %v", err)
	for _, tx := range block.Block.Transactions {
		err = v.Struct(tx)
		assert.Nil(t, err, "Validate struct err: %v", err)
	}
}

func TestExecutorSVMTransactionByHashRequest(t *testing.T) {
	v := validator.New()
	hash := common.HexToHash("0x65f171da62b350b4dbb5a56161cddb1ce0bd12130cd767cc03f2f65e4d5a23f")

	resp, err := as.Root.RequestFuture(ex, &messages.SVMTransactionByHashRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Hash:      hash,
		Chain:     chain,
	}, 10*time.Second).Result()
	assert.Nil(t, err, "SVMTransactionByHashRequest err: %v", err)
	tx, ok := resp.(*messages.SVMTransactionByHashResponse)
	assert.True(t, ok, "expected SVMTransactionByHashResponse, got %s", reflect.TypeOf(resp).String())
	assert.True(t, tx.Success, "failed SVMTransactionByHashResponse, got: %s", tx.RejectionReason.String())
	assert.Equal(t, tx.Transaction.TxnHash.String(), "0x65f171da62b350b4dbb5a56161cddb1ce0bd12130cd767cc03f2f65e4d5a23f", "mismatched hash")

	err = v.Struct(tx.Transaction)
	assert.Nil(t, err, "Validate struct err: %v", err)
}
