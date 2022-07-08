package erc721_test

import (
	"fmt"
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
			ID: 730,
		},
	}
	var asset *models.ProtocolAsset
	for _, a := range response.ProtocolAssets {
		if assetTest.Asset.ID == a.Asset.ID {
			asset = a
		}
	}
	if asset == nil {
		t.Fatal("asset not found")
	}
	props := actor.PropsFromProducer(tests.NewERC721CheckerProducer(asset))
	checker := as.Root.Spawn(props)
	defer as.Root.PoisonFuture(checker)
	time.Sleep(10 * time.Minute)
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
	for _, ev := range d.Updates {
		fmt.Println(ev)
	}
}

func TestListenerSVM(t *testing.T) {
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
	props := actor.PropsFromProducer(tests.NewERC721CheckerProducer(asset))
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
