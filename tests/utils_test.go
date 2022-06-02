package tests

import (
	"context"
	"encoding/json"
	"fmt"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"testing"
)

func TestCaptureStatics(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	registryAddress := "registry.alphaticks.io:8001"
	if os.Getenv("REGISTRY_ADDRESS") != "" {
		registryAddress = os.Getenv("REGISTRY_ADDRESS")
	}
	conn, err := grpc.Dial(registryAddress, grpc.WithInsecure())
	if err != nil {
		err := fmt.Errorf("error connecting to public registry gRPC endpoint: %v", err)
		panic(err)
	}
	rgstr := registry.NewPublicRegistryClient(conn)

	res, err := rgstr.Assets(context.Background(), &registry.AssetsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	assets := make(map[uint32]*xchangerModels.Asset)
	for _, a := range res.Assets {
		assets[a.AssetId] = &xchangerModels.Asset{
			Symbol: a.Symbol,
			Name:   a.Name,
			ID:     a.AssetId,
		}
	}
	b, _ := json.Marshal(assets)
	err = ioutil.WriteFile("assets.json", b, 0644)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := rgstr.Protocols(context.Background(), &registry.ProtocolsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	protocols := make(map[uint32]*xchangerModels.Protocol)
	for _, p := range resp.Protocols {
		protocols[p.ProtocolId] = &xchangerModels.Protocol{
			ID:   p.ProtocolId,
			Name: p.Name,
		}
	}
	b, _ = json.Marshal(protocols)
	err = ioutil.WriteFile("protocols.json", b, 0644)
	if err != nil {
		t.Fatal(err)
	}

	resc, err := rgstr.Chains(context.Background(), &registry.ChainsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	chains := make(map[uint32]*xchangerModels.Chain)
	for _, c := range resc.Chains {
		chains[c.ChainId] = &xchangerModels.Chain{
			ID:   c.ChainId,
			Type: c.Type,
			Name: c.Name,
		}
	}
	b, _ = json.Marshal(chains)
	err = ioutil.WriteFile("chains.json", b, 0644)
	if err != nil {
		t.Fatal(err)
	}
}
