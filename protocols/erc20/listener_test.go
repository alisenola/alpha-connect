package erc20_test

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	ctypes "gitlab.com/alphaticks/alpha-connect/chains/types"
	xtypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/protocols/tests"
	"gitlab.com/alphaticks/alpha-connect/protocols/types"
	exTests "gitlab.com/alphaticks/alpha-connect/tests"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"google.golang.org/grpc"
	"reflect"
	"testing"
	"time"
)

func TestListenerSVM(t *testing.T) {
	exTests.LoadStatics(t)
	protocol := constants.ERC20
	registryAddress := "127.0.0.1:8001"
	conn, err := grpc.Dial(registryAddress, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	reg := registry.NewPublicRegistryClient(conn)

	prCfg := &types.ExecutorConfig{
		Registry:  reg,
		Protocols: []*xchangerModels.Protocol{protocol},
	}
	as, ex, clean := exTests.StartExecutor(t, &xtypes.ExecutorConfig{}, prCfg, &ctypes.ExecutorConfig{}, nil)
	defer clean()

	res, err := as.Root.RequestFuture(ex, &messages.ProtocolAssetListRequest{
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
			ID: 57,
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
}
