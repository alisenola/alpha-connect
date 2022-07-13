package erc20_test

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gitlab.com/alphaticks/alpha-connect/config"
	exTests "gitlab.com/alphaticks/alpha-connect/tests"
	"gitlab.com/alphaticks/alpha-connect/utils"
	gorderbookModels "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"gitlab.com/alphaticks/xchanger/constants"
	"reflect"
	"testing"
	"time"

	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
)

func TestExecutorSVM(t *testing.T) {
	exTests.LoadStatics(t)
	protocol := constants.ERC20
	registryAddress := "127.0.0.1:8001"
	cfg := config.Config{
		RegistryAddress: registryAddress,
		Protocols:       []string{protocol.Name},
	}
	as, executor, clean := exTests.StartExecutor(t, &cfg)
	defer clean()

	testAsset := xchangerModels.Asset{
		ID:     57.,
		Symbol: "DAI",
	}
	chain, ok := constants.GetChainByID(5)
	if !assert.True(t, ok, "missings svm") {
		t.Fatal()
	}
	res, err := as.Root.RequestFuture(executor, &messages.ProtocolAssetListRequest{}, 20*time.Second).Result()
	if !assert.Nil(t, err, "RequestFuture ProtocolAssetList err: %v", err) {
		t.Fatal()
	}
	assets, ok := res.(*messages.ProtocolAssetList)
	if !assert.True(t, ok, "incorrect type assertion") {
		t.Fatal()
	}
	var coll *models.ProtocolAsset
	for _, asset := range assets.ProtocolAssets {
		if asset.Asset.Symbol == testAsset.Symbol && asset.Chain.ID == chain.ID {
			fmt.Printf("asset %+v \n", asset)
			coll = asset
		}
	}
	if !assert.NotNil(t, coll, "missing collection") {
		t.Fatal()
	}
	r, err := as.Root.RequestFuture(executor, &messages.ProtocolAssetDefinitionRequest{
		RequestID:       uint64(time.Now().UnixNano()),
		ProtocolAssetID: utils.GetProtocolAssetID(&testAsset, constants.ERC20, constants.StarknetMainnet),
	}, 15*time.Second).Result()
	if !assert.Nil(t, err, "RequestFuture ProtocolAssetDefinition err: %v", err) {
		t.Fatal()
	}
	def, ok := r.(*messages.ProtocolAssetDefinitionResponse)
	if !assert.True(t, ok, "expected ProtocolAssetDefinitionResponse, got %s", reflect.TypeOf(r).String()) {
		t.Fatal()
	}
	if !assert.True(t, def.Success, "request failed with %v", def.RejectionReason.String()) {
		t.Fatal()
	}
	fmt.Println("Protocol Asset definition", def.ProtocolAsset)
	//Execute the future request for the ERC20 historical data
	resp, err := as.Root.RequestFuture(
		executor,
		&messages.HistoricalProtocolAssetTransferRequest{
			RequestID:  uint64(time.Now().UnixNano()),
			ProtocolID: protocol.ID,
			ChainID:    chain.ID,
			Start:      0,
			Stop:       3000,
		},
		30*time.Second,
	).Result()
	if !assert.Nil(t, err, "RequestFuture HistoricalProtocolAssetTransferRequest err: %v", err) {
		t.Fatal()
	}
	response, ok := resp.(*messages.HistoricalProtocolAssetTransferResponse)
	if !assert.True(t, ok, "expected HistoricalProtocolAssetTransferResponse, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.True(t, response.Success, "request failed with %s", response.RejectionReason.String()) {
		t.Fatal()
	}
	var events []*gorderbookModels.AssetTransfer
	for _, ev := range response.Update {
		events = append(events, ev.Transfers...)
	}
	if !assert.GreaterOrEqual(t, len(events), 7000, "expected more than 20 events") {
		t.Fatal()
	}

	//TODO test for HistoricalProtocolAssetTransferRequest with assetID
}

func TestExecutorEVM(t *testing.T) {
	// Load executor, chain and protocol asset
	exTests.LoadStatics(t)
	protocol := constants.ERC20
	testAsset := xchangerModels.Asset{
		ID: 35.,
	}
	registryAddress := "127.0.0.1:8001"
	cfg := config.Config{
		RegistryAddress: registryAddress,
		Protocols:       []string{protocol.Name},
	}
	as, executor, clean := exTests.StartExecutor(t, &cfg)
	defer clean()
	chain, ok := constants.GetChainByID(1)
	if !assert.True(t, ok, "missings evm") {
		t.Fatal()
	}
	res, err := as.Root.RequestFuture(executor, &messages.ProtocolAssetListRequest{}, 20*time.Second).Result()
	if !assert.Nil(t, err, "RequestFuture ProtocolAssetList err: %v", err) {
		t.Fatal()
	}
	assets, ok := res.(*messages.ProtocolAssetList)
	if !assert.True(t, ok, "incorrect type assertion") {
		t.Fatal()
	}
	if !assert.GreaterOrEqual(t, len(assets.ProtocolAssets), 0, "expected at least one asset") {
		t.Fatal()
	}
	var asset *models.ProtocolAsset
	for _, pa := range assets.ProtocolAssets {
		if pa.Protocol.ID == protocol.ID && chain.ID == chain.ID && pa.Asset.ID == testAsset.ID {
			asset = pa
		}
	}
	if !assert.NotNil(t, asset, "Missing protocol asset") {
		t.Fatal()
	}

	var step uint64 = 10
	var start uint64 = 15110232
	var stop uint64 = 15110300
	done := false
	var events []*gorderbookModels.AssetTransfer

	// test get all historical transfers
	for !done {
		//Execute the future request for the ERC20 historical data
		resp, err := as.Root.RequestFuture(
			executor,
			&messages.HistoricalProtocolAssetTransferRequest{
				RequestID:  uint64(time.Now().UnixNano()),
				ProtocolID: protocol.ID,
				ChainID:    chain.ID,
				Start:      start,
				Stop:       step + start,
			},
			30*time.Second,
		).Result()
		if !assert.Nil(t, err, "RequestFuture HistoricalProtocolAssetTransferRequest err: %v", err) {
			t.Fatal()
		}
		response, ok := resp.(*messages.HistoricalProtocolAssetTransferResponse)
		if !assert.True(t, ok, "expected HistoricalProtocolAssetTransferResponse, got %s", reflect.TypeOf(resp).String()) {
			t.Fatal()
		}
		if !assert.True(t, response.Success, "request failed with %s", response.RejectionReason.String()) {
			t.Fatal()
		}
		for _, ev := range response.Update {
			events = append(events, ev.Transfers...)
		}
		done = step+start > stop
		start += step + 1
	}
	if !assert.GreaterOrEqual(t, len(events), 1000, "expected more than 1000 events") {
		t.Fatal()
	}

	// test get specific asset historical transfer
	start = 15110232
	stop = 15110300
	done = false
	events = nil
	for !done {
		//Execute the future request for the ERC20 historical data
		resp, err := as.Root.RequestFuture(
			executor,
			&messages.HistoricalProtocolAssetTransferRequest{
				RequestID:  uint64(time.Now().UnixNano()),
				ProtocolID: protocol.ID,
				ChainID:    chain.ID,
				AssetID: &wrapperspb.UInt32Value{
					Value: asset.Asset.ID,
				},
				Start: start,
				Stop:  step + start,
			},
			30*time.Second,
		).Result()
		if !assert.Nil(t, err, "RequestFuture HistoricalProtocolAssetTransferRequest err: %v", err) {
			t.Fatal()
		}
		response, ok := resp.(*messages.HistoricalProtocolAssetTransferResponse)
		if !assert.True(t, ok, "expected HistoricalProtocolAssetTransferResponse, got %s", reflect.TypeOf(resp).String()) {
			t.Fatal()
		}
		if !assert.True(t, response.Success, "request failed with %s", response.RejectionReason.String()) {
			t.Fatal()
		}
		for _, ev := range response.Update {
			events = append(events, ev.Transfers...)
		}
		done = step+start > stop
		start += step + 1
	}
	if !assert.GreaterOrEqual(t, len(events), 700, "expected more than 700 events") {
		t.Fatal()
	}
	for _, tx := range events {
		if !assert.Equal(t, asset.ContractAddress.Value, common.BytesToAddress(tx.Contract).String(), "contract address error") {
			t.Fatal()
		}
	}
}
