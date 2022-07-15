package erc721_test

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gitlab.com/alphaticks/alpha-connect/config"
	exTests "gitlab.com/alphaticks/alpha-connect/tests"
	"gitlab.com/alphaticks/xchanger/constants"
	"reflect"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/protocols/tests"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
)

func TestListenerEVM(t *testing.T) {
	// Load executor, chain and protocol asset
	exTests.LoadStatics(t)
	protocol := constants.ERC721
	assetTest := models.ProtocolAsset{
		Asset: &xchangerModels.Asset{
			ID: 275,
		},
		Chain: &xchangerModels.Chain{
			ID:   1,
			Name: "Ethereum Mainnet",
			Type: "EVM",
		},
	}
	registryAddress := "127.0.0.1:8001"
	cfg, err := config.LoadConfig()
	if !assert.Nil(t, err, "LoadConfig err: %v", err) {
		t.Fatal()
	}
	cfg.RegistryAddress = registryAddress
	cfg.Protocols = []string{protocol.Name}
	as, executor, clean := exTests.StartExecutor(t, cfg)
	defer clean()

	res, err := as.Root.RequestFuture(executor, &messages.ProtocolAssetListRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Subscribe: false,
	}, 20*time.Second).Result()
	if !assert.Nil(t, err, "Request future ProtocolListRequest err: %v", err) {
		t.Fatal()
	}
	response, ok := res.(*messages.ProtocolAssetList)
	if !assert.True(t, ok, "expected *messages.ProtocolAssetList, got %s", reflect.TypeOf(res).String()) {
		t.Fatal()
	}
	var asset *models.ProtocolAsset
	for _, a := range response.ProtocolAssets {
		if assetTest.Asset.ID == a.Asset.ID && assetTest.Chain.ID == a.Chain.ID {
			asset = a
		}
	}
	if !assert.NotNil(t, asset, "asset not found") {
		t.Fatal()
	}
	s := asset.Asset

	// Listen on all protocol assets updates
	asset.Asset = nil
	props := actor.PropsFromProducer(tests.NewProtocolCheckerProducer(asset))
	checker := as.Root.Spawn(props)
	defer as.Root.PoisonFuture(checker)
	time.Sleep(2 * time.Minute)
	resp, err := as.Root.RequestFuture(checker, &tests.GetDataRequest{}, 10*time.Second).Result()
	if !assert.Nil(t, err, "Request future GetDataRequest err: %v", err) {
		t.Fatal()
	}
	d, ok := resp.(*tests.GetDataResponse)
	if !assert.True(t, ok, "expected *tests.GetDataResponse, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.Nil(t, d.Err, "GetDataResponse err: %v", err) {
		t.Fatal()
	}
	diff := false
	// check the contract addresses are different
	for _, u := range d.Updates {
		for _, tx := range u.Transfers {
			if asset.ContractAddress.Value != common.BytesToAddress(tx.Contract).String() {
				diff = true
			}
		}
	}
	if !assert.True(t, diff, "expected different contracts") {
		t.Fatal()
	}

	// Listen on a specific protocol asset updates
	asset.Asset = s
	props = actor.PropsFromProducer(tests.NewProtocolCheckerProducer(asset))
	c := as.Root.Spawn(props)
	defer as.Root.PoisonFuture(c)
	time.Sleep(2 * time.Minute)
	resp, err = as.Root.RequestFuture(c, &tests.GetDataRequest{}, 10*time.Second).Result()
	if !assert.Nil(t, err, "RequestFuture GetData err: %v", err) {
		t.Fatal()
	}
	d, ok = resp.(*tests.GetDataResponse)
	if !assert.True(t, ok, "expected *tests.GetDataResponse, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.Nil(t, d.Err, "listener error: %v", d.Err) {
		t.Fatal()
	}
	if !assert.GreaterOrEqual(t, len(d.Updates), 0, "expected length of updates > 0") {
		t.Fatal()
	}
	// check the contract address is what we expect for all transactions
	for _, u := range d.Updates {
		for _, tx := range u.Transfers {
			if !assert.Equal(t, asset.ContractAddress.Value, common.BytesToAddress(tx.Contract).String()) {
				t.Fatal()
			}
		}
	}
}

func TestListenerSVM(t *testing.T) {
	// TODO Listen on all protocol assets updates
	exTests.LoadStatics(t)
	cfg := config.Config{
		RegistryAddress: "registry.alphaticks.io:8001",
		Protocols:       []string{constants.ERC721.Name},
	}
	as, executor, clean := exTests.StartExecutor(t, &cfg)
	defer clean()
	res, err := as.Root.RequestFuture(executor, &messages.ProtocolAssetListRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Subscribe: false,
	}, 20*time.Second).Result()
	if err != nil {
		t.Fatal()
	}
	response, ok := res.(*messages.ProtocolAssetList)
	if !ok {
		t.Fatal("incorrect type assertion")
	}
	fmt.Println(response.ProtocolAssets)
	assetTest := models.ProtocolAsset{
		Asset: &xchangerModels.Asset{
			ID: 275,
		},
	}
	var asset *models.ProtocolAsset
	chain, ok := constants.GetChainByID(5)
	if !ok {
		t.Fatal("missing chain")
	}
	for _, a := range response.ProtocolAssets {
		if assetTest.Asset.ID == a.Asset.ID && a.Chain.ID == chain.ID {
			asset = a
		}
	}
	if asset == nil {
		t.Fatal("asset not found")
	}
	props := actor.PropsFromProducer(tests.NewProtocolCheckerProducer(asset))
	checker := as.Root.Spawn(props)
	defer as.Root.PoisonFuture(checker)
	time.Sleep(3 * time.Minute)
	resp, err := as.Root.RequestFuture(checker, &tests.GetDataRequest{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	d, ok := resp.(*tests.GetDataResponse)
	if !ok {
		t.Fatalf("expected *tests.GetDataResponse, got %s", reflect.TypeOf(resp).String())
	}
	if d.Err != nil {
		t.Fatal(d.Err)
	}
}
