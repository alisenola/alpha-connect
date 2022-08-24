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
		ChainRPCs: []config.ChainRPC{
			{Chain: 5, Endpoint: "http://127.0.0.1:9545"},
		},
	}
	as, ex, poison := tests.StartExecutor(t, &cfg)
	defer poison()
	resp, err := as.Root.RequestFuture(ex, &messages.BlockNumberRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Chain:     chain,
	}, 10*time.Second).Result()
	if !assert.Nil(t, err, "BlockNumberRequest err: %v", err) {
		t.Fatal()
	}
	num, ok := resp.(*messages.BlockNumberResponse)
	if !assert.True(t, ok, "expected BlockNumberRequest, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.True(t, num.Success, "failed BlockNumberRequest, got %s", num.RejectionReason.String()) {
		t.Fatal()
	}
	if !assert.Greater(t, num.BlockNumber, uint64(2000), "expected block to be higher than 2000") {
		t.Fatal()
	}
}

func TestExecutorSVMEventsQueryRequest(t *testing.T) {
	chain := &models2.Chain{
		ID:   5,
		Name: "Starknet Mainnet",
		Type: "SVM",
	}
	cfg := config.Config{
		ChainRPCs: []config.ChainRPC{
			{Chain: 5, Endpoint: "http://127.0.0.1:9545"},
		},
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
	if !assert.Nil(t, err, "SVMEventsQueryRequest err: %v", err) {
		t.Fatal()
	}
	events, ok := resp.(*messages.SVMEventsQueryResponse)
	if !assert.True(t, ok, "expected SVMEventsQueryRequest, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.True(t, events.Success, "failed SVMEventsQueryRequest, got: %s", events.RejectionReason.String()) {
		t.Fatal()
	}
	if !assert.GreaterOrEqual(t, len(events.Events), 6, "expected at least 6 events") {
		t.Fatal()
	}

	for _, ev := range events.Events {
		err = v.Struct(ev)
		assert.Nil(t, err, "Validate struct err: %v", err)
	}
}

func TestExecutorSVMContractCallRequest(t *testing.T) {
	chain := &models2.Chain{
		ID:   5,
		Name: "Starknet Mainnet",
		Type: "SVM",
	}
	cfg := config.Config{
		ChainRPCs: []config.ChainRPC{
			{Chain: 5, Endpoint: "http://127.0.0.1:9545"},
		},
	}
	as, ex, poison := tests.StartExecutor(t, &cfg)
	defer poison()
	add := common.HexToHash("0x00da114221cb83fa859dbdb4c44beeaa0bb37c7537ad5ae66fe5e0efd20e6eb3")
	entrypoint := common.BigToHash(svm.GetSelectorFromName("decimals"))
	q := svm.CallParameters{
		CallData: &svm.CallData{
			ContractAddress:    &add,
			EntryPointSelector: &entrypoint,
		},
		BlockTag: &svm.LATEST_BLOCK,
	}
	res, err := as.Root.RequestFuture(ex, &messages.SVMContractCallRequest{
		RequestID:      uint64(time.Now().UnixNano()),
		Chain:          chain,
		CallParameters: &q,
	}, 15*time.Second).Result()
	if !assert.Nil(t, err, "SVMContractCall err: %v", err) {
		t.Fatal()
	}
	call, ok := res.(*messages.SVMContractCallResponse)
	if !assert.True(t, ok, "expected SVMContractCallResponse, got %s", reflect.TypeOf(res).String()) {
		t.Fatal()
	}
	if !assert.True(t, call.Success, "SVMContractCall failed with %s", call.RejectionReason.String()) {
		t.Fatal()
	}
	if !assert.Equal(t, len(call.Data), 1, "expected one array") {
		t.Fatal()
	}
	i := big.NewInt(1).SetBytes(call.Data[0].Bytes())
	if !assert.EqualValues(t, i, big.NewInt(18), "expected 18 decimals") {
		t.Fatal()
	}
}

func TestExecutorSVMContractClassRequest(t *testing.T) {
	chain := &models2.Chain{
		ID:   5,
		Name: "Starknet Mainnet",
		Type: "SVM",
	}
	cfg := config.Config{
		ChainRPCs: []config.ChainRPC{
			{Chain: 5, Endpoint: "http://127.0.0.1:9545"},
		},
	}
	as, ex, poison := tests.StartExecutor(t, &cfg)
	defer poison()
	add := common.HexToHash("0x00da114221cb83fa859dbdb4c44beeaa0bb37c7537ad5ae66fe5e0efd20e6eb3")
	// get class erc20
	q := svm.ContractClassQuery{
		ContractAddress: &add,
	}
	resp, err := as.Root.RequestFuture(ex, &messages.SVMContractClassRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Chain:     chain,
		Query:     &q,
	}, 15*time.Second).Result()
	if !assert.Nil(t, err, "SVMContractClassRequest err: %v", err) {
		t.Fatal()
	}
	class, ok := resp.(*messages.SVMContractClassResponse)
	if !assert.True(t, ok, "expected SVMContractClassResponse, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.True(t, class.Success, "failed SVMContractClassResponse, got: %s", class.RejectionReason.String()) {
		t.Fatal()
	}
	found := false
	selector := svm.GetSelectorFromName("decimals")
	for _, c := range class.Class.EntryPointsByType.External {
		i, ok := big.NewInt(0).SetString(c.Selector[2:], 16)
		if !assert.True(t, ok, "error with big.Int SetString for string %s", c.Selector[2:]) {
			t.Fatal()
		}
		if selector.Cmp(i) == 0 {
			found = true
		}
	}
	if !assert.True(t, found, "expected contract to have decimals") {
		t.Fatal()
	}
}

func TestExecutorSVMEventsQueryRequestTransfer(t *testing.T) {
	chain := &models2.Chain{
		ID:   5,
		Name: "Starknet Mainnet",
		Type: "SVM",
	}
	cfg := config.Config{
		ChainRPCs: []config.ChainRPC{
			{Chain: 5, Endpoint: "http://127.0.0.1:9545"},
		},
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
	if !assert.Nil(t, err, "SVMEventsQueryRequest err: %v", err) {
		t.Fatal()
	}
	events, ok := resp.(*messages.SVMEventsQueryResponse)
	if !assert.True(t, ok, "expected SVMEventsQueryRequest, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.True(t, events.Success, "failed SVMEventsQueryRequest, got: %s", events.RejectionReason.String()) {
		t.Fatal()
	}
	if !assert.GreaterOrEqual(t, len(events.Events), 1057, "expected at least 1057 events") {
		t.Fatal()
	}

	for _, ev := range events.Events {
		err = v.Struct(ev)
		if !assert.Nil(t, err, "Validate struct err: %v", err) {
			t.Fatal()
		}
	}
	if !assert.Equal(t, len(events.Events), len(events.Times), "mismatched length") {
		t.Fatal()
	}
}

func TestExecutorSVMBlockQueryRequest(t *testing.T) {
	chain := &models2.Chain{
		ID:   5,
		Name: "Starknet Mainnet",
		Type: "SVM",
	}
	cfg := config.Config{
		ChainRPCs: []config.ChainRPC{
			{Chain: 5, Endpoint: "http://127.0.0.1:9545"},
		},
	}
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
	if !assert.Nil(t, err, "SVMEventsQueryRequest err: %v", err) {
		t.Fatal()
	}
	block, ok := resp.(*messages.SVMBlockQueryResponse)
	if !assert.True(t, ok, "expected SVMBlockQueryResponse, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.True(t, block.Success, "failed SVMBlockQueryResponse, got: %s", block.RejectionReason.String()) {
		t.Fatal()
	}
	if !assert.GreaterOrEqual(t, len(block.Block.Transactions), 265, "expected at least 265 transactions") {
		t.Fatal()
	}

	err = v.Struct(block.Block)
	if !assert.Nil(t, err, "Validate struct err: %v", err) {
		t.Fatal()
	}
	for _, tx := range block.Block.Transactions {
		err = v.Struct(tx)
		if !assert.Nil(t, err, "Validate struct err: %v", err) {
			t.Fatal()
		}
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
	if !assert.Nil(t, err, "SVMEventsQueryRequest err: %v", err) {
		t.Fatal()
	}
	block, ok = resp.(*messages.SVMBlockQueryResponse)
	if !assert.True(t, ok, "expected SVMBlockQueryResponse, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.True(t, block.Success, "failed SVMBlockQueryResponse, got: %s", block.RejectionReason.String()) {
		t.Fatal()
	}
	if !assert.Equal(t, len(block.Block.Transactions), 92, "expected at least 92 transactions") {
		t.Fatal()
	}

	err = v.Struct(block.Block)
	if !assert.Nil(t, err, "Validate struct err: %v", err) {
		t.Fatal()
	}
	for _, tx := range block.Block.Transactions {
		err = v.Struct(tx)
		if !assert.Nil(t, err, "Validate struct err: %v", err) {
			t.Fatal()
		}
	}
}

func TestExecutorSVMTransactionByHashRequest(t *testing.T) {
	chain := &models2.Chain{
		ID:   5,
		Name: "Starknet Mainnet",
		Type: "SVM",
	}
	cfg := config.Config{
		ChainRPCs: []config.ChainRPC{
			{Chain: 5, Endpoint: "http://127.0.0.1:9545"},
		},
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
	if !assert.Nil(t, err, "SVMTransactionByHashRequest err: %v", err) {
		t.Fatal()
	}
	tx, ok := resp.(*messages.SVMTransactionByHashResponse)
	if !assert.True(t, ok, "expected SVMTransactionByHashResponse, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.True(t, tx.Success, "failed SVMTransactionByHashResponse, got: %s", tx.RejectionReason.String()) {
		t.Fatal()
	}
	if !assert.Equal(t, tx.Transaction.TxnHash.String(), "0x065f171da62b350b4dbb5a56161cddb1ce0bd12130cd767cc03f2f65e4d5a23f", "mismatched hash") {
		t.Fatal()
	}

	err = v.Struct(tx.Transaction)
	if !assert.Nil(t, err, "Validate struct err: %v", err) {
		t.Fatal()
	}
}
