package svm_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gitlab.com/alphaticks/alpha-connect/config"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/tests"
	"gitlab.com/alphaticks/xchanger/chains/svm"
	models2 "gitlab.com/alphaticks/xchanger/models"
	"math/big"
	"reflect"
	"testing"
	"time"
)

func TestExecutorBlockNumberRequest(t *testing.T) {
	chain := &models2.Chain{
		ID:   5,
		Name: "Starknet Mainnet",
		Type: "SVM",
	}
	cfg := config.Config{
		Chains: []string{"Starknet Mainnet"},
	}
	as, ex, poison := tests.StartExecutor(t, &cfg)
	defer poison()
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
	chain := &models2.Chain{
		ID:   5,
		Name: "Starknet Mainnet",
		Type: "SVM",
	}
	cfg := config.Config{
		Chains: []string{"Starknet Mainnet"},
	}
	as, ex, poison := tests.StartExecutor(t, &cfg)
	defer poison()
	v := validator.New()
	add := common.HexToHash("0x0276069eb59afc97d3f6dcd18ec236cbb3e71611324bfc6e86a6aee0851a65bb")
	q := svm.EventQuery{
		To:              big.NewInt(2790),
		ContractAddress: &add,
		PageSize:        1000,
		PageNumber:      0,
	}
	resp, err := as.Root.RequestFuture(ex, &messages.SVMEventsQueryRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Query:     q,
		Chain:     chain,
	}, 10*time.Second).Result()
	assert.Nil(t, err, "SVMEventsQueryRequest err: %v", err)
	events, ok := resp.(*messages.SVMEventsQueryResponse)
	assert.True(t, ok, "expected SVMEventsQueryRequest, got %s", reflect.TypeOf(resp).String())
	assert.True(t, events.Success, "failed SVMEventsQueryRequest, got: %s", events.RejectionReason.String())
	assert.GreaterOrEqual(t, len(events.Events), 6, "expected at least 6 events")

	for _, ev := range events.Events {
		err = v.Struct(ev)
		assert.Nil(t, err, "Validate struct err: %v", err)
	}
}

func TestExecutorSVMEventsQueryRequestTransfer(t *testing.T) {
	chain := &models2.Chain{
		ID:   5,
		Name: "Starknet Mainnet",
		Type: "SVM",
	}
	cfg := config.Config{
		Chains: []string{"Starknet Mainnet"},
	}
	as, ex, poison := tests.StartExecutor(t, &cfg)
	defer poison()
	v := validator.New()
	key := common.HexToHash("0x0099cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9")
	q := svm.EventQuery{
		To:       big.NewInt(2493),
		Keys:     &[]common.Hash{key},
		PageSize: 1000,
	}
	resp, err := as.Root.RequestFuture(ex, &messages.SVMEventsQueryRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Query:     q,
		Chain:     chain,
	}, 2*time.Minute).Result()
	assert.Nil(t, err, "SVMEventsQueryRequest err: %v", err)
	events, ok := resp.(*messages.SVMEventsQueryResponse)
	assert.True(t, ok, "expected SVMEventsQueryRequest, got %s", reflect.TypeOf(resp).String())
	assert.True(t, events.Success, "failed SVMEventsQueryRequest, got: %s", events.RejectionReason.String())
	assert.GreaterOrEqual(t, len(events.Events), 1057, "expected at least 1057 events")

	for _, ev := range events.Events {
		err = v.Struct(ev)
		assert.Nil(t, err, "Validate struct err: %v", err)
	}
	assert.Equal(t, len(events.Events), len(events.Times), "mismatched length")
}

func TestExecutorSVMBlockQueryRequest(t *testing.T) {
	chain := &models2.Chain{
		ID:   5,
		Name: "Starknet Mainnet",
		Type: "SVM",
	}
	cfg := config.Config{}
	as, ex, poison := tests.StartExecutor(t, &cfg)
	defer poison()
	v := validator.New()
	hash := common.HexToHash("0x69b96255bf7cc630ba99292ca1dd34130829fa3486fe44bfbb7f7aa13a4da29")

	q := &svm.BlockQuery{
		BlockHash: &hash,
		TxScope:   &svm.FULL_TXN_AND_RECEIPTS,
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
	assert.GreaterOrEqual(t, len(block.Block.Transactions), 265, "expected at least 265 transactions")

	err = v.Struct(block.Block)
	assert.Nil(t, err, "Validate struct err: %v", err)
	for _, tx := range block.Block.Transactions {
		err = v.Struct(tx)
		assert.Nil(t, err, "Validate struct err: %v", err)
	}

	blockN := big.NewInt(1828)
	q = &svm.BlockQuery{
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
	assert.Equal(t, len(block.Block.Transactions), 92, "expected at least 92 transactions")

	err = v.Struct(block.Block)
	assert.Nil(t, err, "Validate struct err: %v", err)
	for _, tx := range block.Block.Transactions {
		err = v.Struct(tx)
		assert.Nil(t, err, "Validate struct err: %v", err)
	}
}

func TestExecutorSVMTransactionByHashRequest(t *testing.T) {
	chain := &models2.Chain{
		ID:   5,
		Name: "Starknet Mainnet",
		Type: "SVM",
	}
	cfg := config.Config{
		Chains: []string{"Starknet Mainnet"},
	}
	as, ex, poison := tests.StartExecutor(t, &cfg)
	defer poison()
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
