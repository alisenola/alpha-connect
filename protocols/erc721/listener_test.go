package erc721_test

import (
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/chains"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	"gitlab.com/alphaticks/alpha-connect/protocols"
	exTests "gitlab.com/alphaticks/alpha-connect/tests"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	"gitlab.com/alphaticks/xchanger/constants"
	"google.golang.org/grpc"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/protocols/tests"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
)

func TestListener(t *testing.T) {
	exTests.LoadStatics(t)
	protocol := constants.ERC721
	registryAddress := "registry.alphaticks.io:8001"
	conn, err := grpc.Dial(registryAddress, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	reg := registry.NewPublicRegistryClient(conn)

	prCfg := &protocols.ExecutorConfig{
		Registry:  reg,
		Protocols: []*xchangerModels.Protocol{protocol},
	}
	as, ex, clean := exTests.StartExecutor(t, &exchanges.ExecutorConfig{}, prCfg, &chains.ExecutorConfig{}, nil)
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
	time.Sleep(800 * time.Second)
	as.Root.Poison(checker)
	time.Sleep(10 * time.Second)
	// TODO check if
}
