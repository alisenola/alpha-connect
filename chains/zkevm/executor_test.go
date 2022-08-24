package zkevm_test

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gitlab.com/alphaticks/alpha-connect/chains/tests"
	"gitlab.com/alphaticks/alpha-connect/config"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	exTests "gitlab.com/alphaticks/alpha-connect/tests"
	"gitlab.com/alphaticks/xchanger/chains/evm"
	"gitlab.com/alphaticks/xchanger/constants"
	evm20 "gitlab.com/alphaticks/xchanger/protocols/erc20/evm"
	"math/big"
	"reflect"
	"testing"
	"time"
)

func TestExecutorBlockNumberRequest(t *testing.T) {
	// Load executor, chain and protocol asset
	exTests.LoadStatics(t)
	chain, ok := constants.GetChainByID(6)
	if !assert.True(t, ok, "Missing ZKSync chain") {
		t.Fatal()
	}
	cfg, err := config.LoadConfig()
	if !assert.Nil(t, err, "LoadConfig err: %v", err) {
		t.Fatal()
	}
	as, ex, clean := exTests.StartExecutor(t, cfg)
	defer clean()

	// get current block number
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
	// check greater than 2000
	if !assert.Greater(t, num.BlockNumber, uint64(2000), "expected block to be higher than 2000") {
		t.Fatal()
	}
}

func TestExecutorBlockInfoRequest(t *testing.T) {
	// Load executor, chain and protocol asset
	exTests.LoadStatics(t)
	chain, ok := constants.GetChainByID(6)
	if !assert.True(t, ok, "Missing ZKSync chain") {
		t.Fatal()
	}
	cfg, err := config.LoadConfig()
	if !assert.Nil(t, err, "LoadConfig err: %v", err) {
		t.Fatal()
	}
	as, ex, clean := exTests.StartExecutor(t, cfg)
	defer clean()

	// get block info at block 123456
	resp, err := as.Root.RequestFuture(ex, &messages.BlockInfoRequest{
		RequestID:   uint64(time.Now().UnixNano()),
		Chain:       chain,
		BlockNumber: uint64(123456),
	}, 10*time.Second).Result()
	if !assert.Nil(t, err, "SVMEventsQueryRequest err: %v", err) {
		t.Fatal()
	}
	info, ok := resp.(*messages.BlockInfoResponse)
	if !assert.True(t, ok, "expected BlockInfoResponse, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.True(t, info.Success, "failed BlockInfoResponse, got: %s", info.RejectionReason.String()) {
		t.Fatal()
	}
	loc, err := time.LoadLocation("Local")
	if !assert.Nil(t, err, "LoadLocation err: %v", err) {
		t.Fatal()
	}
	// check block time is 4/7/2022 16:28:58 CEST
	if !assert.EqualValues(t, info.BlockTime, time.Date(2022, 7, 4, 4, 28, 58, 0, loc), "expected time greater than 0") {
		t.Fatal()
	}
}

func TestExecutorContractCallRequest(t *testing.T) {
	// Load executor, chain and protocol asset
	exTests.LoadStatics(t)
	chain, ok := constants.GetChainByID(6)
	if !assert.True(t, ok, "Missing ZKSync chain") {
		t.Fatal()
	}
	cfg, err := config.LoadConfig()
	if !assert.Nil(t, err, "LoadConfig err: %v", err) {
		t.Fatal()
	}
	as, ex, clean := exTests.StartExecutor(t, cfg)
	defer clean()

	// get current block number
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

	abi, err := evm20.ERC20MetaData.GetAbi()
	if !assert.Nil(t, err, "GetAbi err: %v", err) {
		t.Fatal()
	}
	add := common.HexToAddress("0x46384918127fBd1679C757DF7b495C3F61481467")
	res, err := as.Root.RequestFuture(ex, &messages.EVMContractCallRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Chain:     chain,
		Msg: ethereum.CallMsg{
			To:   &add,
			Data: abi.Methods["decimals"].ID,
		},
		BlockNumber: num.BlockNumber,
	}, 15*time.Second).Result()
	if !assert.Nil(t, err, "EVMContractCallRequest err: %v", err) {
		t.Fatal()
	}
	call, ok := res.(*messages.EVMContractCallResponse)
	if !assert.True(t, ok, "expected EVMContractCallResponse, got %s", reflect.TypeOf(res).String()) {
		t.Fatal()
	}
	if !assert.True(t, call.Success, "EVMContractCall failed with %s", call.RejectionReason.String()) {
		t.Fatal()
	}
	i := big.NewInt(1).SetBytes(call.Out)
	if !assert.EqualValues(t, i, big.NewInt(6), "expected 18 decimals") {
		t.Fatal()
	}
}

func TestExecutorEVMLogsQueryRequest(t *testing.T) {
	// Load executor, chain and protocol asset
	exTests.LoadStatics(t)
	chain, ok := constants.GetChainByID(6)
	if !assert.True(t, ok, "Missing ZKSync chain") {
		t.Fatal()
	}
	cfg, err := config.LoadConfig()
	if !assert.Nil(t, err, "LoadConfig err: %v", err) {
		t.Fatal()
	}
	as, ex, clean := exTests.StartExecutor(t, cfg)
	defer clean()

	// get current block number
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

	// make ethereum logs query
	abi, err := evm20.ERC20MetaData.GetAbi()
	if !assert.Nil(t, err, "GetAbi err: %v", err) {
		t.Fatal()
	}
	add := common.HexToAddress("0x46384918127fBd1679C757DF7b495C3F61481467")
	q := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(num.BlockNumber - 99)),
		ToBlock:   big.NewInt(int64(num.BlockNumber)),
		Addresses: []common.Address{add},
		Topics: [][]common.Hash{
			{
				abi.Events["Transfer"].ID,
			},
		},
	}
	resp, err = as.Root.RequestFuture(ex, &messages.EVMLogsQueryRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Chain:     chain,
		Query:     q,
	}, 15*time.Second).Result()
	if !assert.Nil(t, err, "EVMLogsQueryRequest err: %v", err) {
		t.Fatal()
	}
	logs, ok := resp.(*messages.EVMLogsQueryResponse)
	if !assert.True(t, ok, "expected EVMLogsQueryResponse, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.True(t, logs.Success, "failed EVMLogsQuery, got: %s", logs.RejectionReason.String()) {
		t.Fatal()
	}
	if !assert.Greater(t, len(logs.Logs), 0, "expected at least one log") {
		t.Fatal()
	}
	if !assert.Equal(t, len(logs.Logs), len(logs.Times), "expected same length of logs and times") {
		t.Fatal()
	}
}

func TestExecutorEVMLogsSubscribeRequest(t *testing.T) {
	// Load executor, chain and protocol asset
	exTests.LoadStatics(t)
	chain, ok := constants.GetChainByID(6)
	if !assert.True(t, ok, "Missing ZKSync chain") {
		t.Fatal()
	}
	cfg, err := config.LoadConfig()
	if !assert.Nil(t, err, "LoadConfig err: %v", err) {
		t.Fatal()
	}
	as, ex, clean := exTests.StartExecutor(t, cfg)
	defer clean()

	// make ethereum logs query
	abi, err := evm20.ERC20MetaData.GetAbi()
	if !assert.Nil(t, err, "GetAbi err: %v", err) {
		t.Fatal()
	}
	q := ethereum.FilterQuery{
		Addresses: []common.Address{
			common.HexToAddress("0x46384918127fBd1679C757DF7b495C3F61481467"),
		},
		Topics: [][]common.Hash{
			{
				abi.Events["Transfer"].ID,
			},
		},
	}
	// start chain checker
	props := actor.PropsFromProducer(tests.NewChainCheckerProducer(chain, q))
	pid, err := as.Root.SpawnNamed(props, "chain_checker")
	defer as.Root.PoisonFuture(pid)
	if !assert.Nil(t, err, "SpawnNamed err: %v", err) {
		t.Fatal()
	}
	time.Sleep(5 * time.Minute)

	// get data from checker
	resp, err := as.Root.RequestFuture(pid, &tests.GetDataRequest{}, 15*time.Second).Result()
	if !assert.Nil(t, err, "GetDataRequest err: %v", err) {
		t.Fatal()
	}
	data, ok := resp.(*tests.GetDataResponse)
	if !assert.True(t, ok, "expected GetDataResponse, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.Nil(t, data.Err, "GetData err: %v", data.Err) {
		t.Fatal()
	}

	// verify logs by querying via EVMLogsQueryRequest
	for _, up := range data.Updates {
		q.FromBlock = big.NewInt(int64(up.BlockNumber))
		q.ToBlock = big.NewInt(int64(up.BlockNumber))
		r, err := as.Root.RequestFuture(ex, &messages.EVMLogsQueryRequest{
			RequestID: uint64(time.Now().UnixNano()),
			Chain:     chain,
			Query:     q,
		}, 15*time.Second).Result()
		if !assert.Nil(t, err, "EVMLogsQueryRequest err: %v", err) {
			t.Fatal()
		}
		logs, ok := r.(*messages.EVMLogsQueryResponse)
		if !assert.True(t, ok, "expected EVMLogsQueryResponse, got %s", reflect.TypeOf(resp).String()) {
			t.Fatal()
		}
		if !assert.True(t, logs.Success, "failed EVMLogsQuery, got: %s", logs.RejectionReason.String()) {
			t.Fatal()
		}
		var futTxs []evm20.ERC20Transfer
		var txs []evm20.ERC20Transfer
		for i, l := range up.Logs {
			transfer := evm20.ERC20Transfer{}
			err := evm.UnpackLog(abi, &transfer, "Transfer", l)
			if !assert.Nil(t, err, "UnpackLog err: %v", err) {
				t.Fatal()
			}
			txs = append(txs, transfer)
			transfer = evm20.ERC20Transfer{}
			err = evm.UnpackLog(abi, &transfer, "Transfer", logs.Logs[i])
			if !assert.Nil(t, err, "UnpackLog err: %v", err) {
				t.Fatal()
			}
			futTxs = append(futTxs, transfer)
		}
		if !assert.Equal(t, len(txs), len(futTxs), "expected same amount of transactions") {
			t.Fatal()
		}
		for _, tx := range futTxs {
			if !assert.True(t, found(tx, txs), "did not find transaction: %+v", tx) {
				t.Fatal()
			}
		}
	}
}

func found(i evm20.ERC20Transfer, arr []evm20.ERC20Transfer) bool {
	for _, v := range arr {
		if v.To == i.To && v.From == i.From && v.Value.Cmp(i.Value) == 0 {
			return true
		}
	}
	return false
}
