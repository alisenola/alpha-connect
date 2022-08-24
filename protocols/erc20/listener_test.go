package erc20_test

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gitlab.com/alphaticks/alpha-connect/config"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/protocols/tests"
	exTests "gitlab.com/alphaticks/alpha-connect/tests"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"reflect"
	"testing"
	"time"
)

func TestListenerSVM(t *testing.T) {
	// Load statics
	exTests.LoadStatics(t)
	registryAddress := "127.0.0.1:8001"
	cfg, err := config.LoadConfig()
	if !assert.Nil(t, err, "LoadConfig err: %v", err) {
		t.Fatal()
	}
	cfg.RegistryAddress = registryAddress
	cfg.Protocols = []string{"ERC-20"}
	as, ex, clean := exTests.StartExecutor(t, cfg)
	defer clean()

	res, err := as.Root.RequestFuture(ex, &messages.ProtocolAssetListRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Subscribe: false,
	}, 20*time.Second).Result()
	if !assert.Nil(t, err, "RequestFuture ProtocolAssetList err: %v", err) {
		t.Fatal()
	}
	response, ok := res.(*messages.ProtocolAssetList)
	if !assert.True(t, ok, "expected *messages.ProtocolAssetList, got %s", reflect.TypeOf(res).String()) {
		t.Fatal()
	}
	assetTest := models.ProtocolAsset{
		Asset: &xchangerModels.Asset{
			ID: 57,
		},
	}
	var asset *models.ProtocolAsset
	chain, ok := constants.GetChainByID(5)
	if !assert.True(t, ok, "missing chain") {
		t.Fatal()
	}
	for _, a := range response.ProtocolAssets {
		if assetTest.Asset.ID == a.Asset.ID && a.Chain.ID == chain.ID {
			asset = a
		}
	}
	if !assert.NotNil(t, asset, "asset not found") {
		t.Fatal()
	}
	s := asset.Asset

	// listen on all protocol assets
	asset.Asset = nil
	props := actor.PropsFromProducer(tests.NewProtocolCheckerProducer(asset))
	checker := as.Root.Spawn(props)
	time.Sleep(2 * time.Minute)
	resp, err := as.Root.RequestFuture(checker, &tests.GetDataRequest{}, 10*time.Second).Result()
	as.Root.PoisonFuture(checker)
	if !assert.Nil(t, err, "RequestFuture GetData err: %v", err) {
		t.Fatal()
	}
	d, ok := resp.(*tests.GetDataResponse)
	if !assert.True(t, ok, "expected *tests.GetDataResponse, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.Nil(t, d.Err, "listener error: %v", d.Err) {
		t.Fatal()
	}

	// listen on a specific protocol asset
	asset.Asset = s
	props = actor.PropsFromProducer(tests.NewProtocolCheckerProducer(asset))
	checker = as.Root.Spawn(props)
	time.Sleep(2 * time.Minute)
	resp, err = as.Root.RequestFuture(checker, &tests.GetDataRequest{}, 10*time.Second).Result()
	as.Root.PoisonFuture(checker)
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
}

func TestListenerEVM(t *testing.T) {
	// Load statics
	exTests.LoadStatics(t)
	registryAddress := "127.0.0.1:8001"
	cfg, err := config.LoadConfig()
	if !assert.Nil(t, err, "LoadConfig err: %v", err) {
		t.Fatal()
	}
	cfg.RegistryAddress = registryAddress
	cfg.Protocols = []string{"ERC-20"}
	as, ex, clean := exTests.StartExecutor(t, cfg)
	defer clean()

	// get asset
	res, err := as.Root.RequestFuture(ex, &messages.ProtocolAssetListRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Subscribe: false,
	}, 20*time.Second).Result()
	if !assert.Nil(t, err, "RequestFuture ProtocolAssetList err: %v", err) {
		t.Fatal()
	}
	response, ok := res.(*messages.ProtocolAssetList)
	if !assert.True(t, ok, "expected *messages.ProtocolAssetList, got %s", reflect.TypeOf(res).String()) {
		t.Fatal()
	}
	assetTest := models.ProtocolAsset{
		Asset: &xchangerModels.Asset{
			ID: 35,
		},
	}
	var asset *models.ProtocolAsset
	chain, ok := constants.GetChainByID(1)
	if !assert.True(t, ok, "missing chain") {
		t.Fatal()
	}
	for _, a := range response.ProtocolAssets {
		if assetTest.Asset.ID == a.Asset.ID && a.Chain.ID == chain.ID {
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
	time.Sleep(2 * time.Minute)
	resp, err := as.Root.RequestFuture(checker, &tests.GetDataRequest{}, 10*time.Second).Result()
	if !assert.Nil(t, err, "RequestFuture GetData err: %v", err) {
		t.Fatal()
	}
	as.Root.PoisonFuture(checker)
	d, ok := resp.(*tests.GetDataResponse)
	if !assert.True(t, ok, "expected *tests.GetDataResponse, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.Nil(t, d.Err, "listener error: %v", d.Err) {
		t.Fatal()
	}
	if !assert.GreaterOrEqual(t, len(d.Updates), 0, "expected length of updates > 0") {
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

func TestListenerZKEVM(t *testing.T) {
	// Load statics
	exTests.LoadStatics(t)
	testAsset := models.ProtocolAsset{
		Asset: &xchangerModels.Asset{
			ID: 12,
		},
	}
	registryAddress := "127.0.0.1:8001"
	cfg, err := config.LoadConfig()
	if !assert.Nil(t, err, "LoadConfig err: %v", err) {
		t.Fatal()
	}
	cfg.RegistryAddress = registryAddress
	cfg.Protocols = []string{"ERC-20"}
	as, ex, clean := exTests.StartExecutor(t, cfg)
	defer clean()

	// get asset
	chain, ok := constants.GetChainByID(6)
	if !assert.True(t, ok, "missing zkevm") {
		t.Fatal()
	}
	res, err := as.Root.RequestFuture(ex, &messages.ProtocolAssetListRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Subscribe: false,
	}, 20*time.Second).Result()
	if !assert.Nil(t, err, "RequestFuture ProtocolAssetList err: %v", err) {
		t.Fatal()
	}
	response, ok := res.(*messages.ProtocolAssetList)
	if !assert.True(t, ok, "expected *messages.ProtocolAssetList, got %s", reflect.TypeOf(res).String()) {
		t.Fatal()
	}
	var asset *models.ProtocolAsset
	for _, a := range response.ProtocolAssets {
		if testAsset.Asset.ID == a.Asset.ID && a.Chain.ID == chain.ID {
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
	time.Sleep(2 * time.Minute)
	resp, err := as.Root.RequestFuture(checker, &tests.GetDataRequest{}, 10*time.Second).Result()
	if !assert.Nil(t, err, "RequestFuture GetData err: %v", err) {
		t.Fatal()
	}
	as.Root.PoisonFuture(checker)
	d, ok := resp.(*tests.GetDataResponse)
	if !assert.True(t, ok, "expected *tests.GetDataResponse, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.Nil(t, d.Err, "listener error: %v", d.Err) {
		t.Fatal()
	}
	if !assert.GreaterOrEqual(t, len(d.Updates), 0, "expected length of updates > 0") {
		t.Fatal()
	}
	diff := false
	// check the contract addresses are different
outer:
	for _, u := range d.Updates {
		for _, tx := range u.Transfers {
			if asset.ContractAddress.Value != common.BytesToAddress(tx.Contract).String() {
				diff = true
				break outer
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
	time.Sleep(1 * time.Minute)
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
