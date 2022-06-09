package executor_test

import (
	"fmt"
	chtypes "gitlab.com/alphaticks/alpha-connect/chains/types"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	prtypes "gitlab.com/alphaticks/alpha-connect/protocols/types"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"os"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/executor"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"google.golang.org/grpc"
)

func TestMainExecutor(t *testing.T) {
	registryAddress := "registry.alphaticks.io:7001"
	if os.Getenv("REGISTRY_ADDRESS") != "" {
		registryAddress = os.Getenv("REGISTRY_ADDRESS")
	}
	conn, err := grpc.Dial(registryAddress, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	reg := registry.NewPublicRegistryClient(conn)

	cfgEx := extypes.ExecutorConfig{
		Registry: reg,
		Exchanges: []*xchangerModels.Exchange{
			{
				Name: "uniswapv3",
				ID:   0x20,
			},
		},
	}
	cfgPr := prtypes.ExecutorConfig{
		Registry: reg,
		Protocols: []*xchangerModels.Protocol{
			{
				Name: "ERC721",
				ID:   0x01,
			},
		},
	}
	chPr := chtypes.ExecutorConfig{}
	prod := executor.NewExecutorProducer(&cfgEx, &cfgPr, &chPr)
	as := actor.NewActorSystem()
	ex, err := as.Root.SpawnNamed(actor.PropsFromProducer(prod), "executor")
	if err != nil {
		t.Fatal(err)
	}
	res, err := as.Root.RequestFuture(ex, &messages.SecurityListRequest{
		RequestID: 0,
	}, 20*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	response, ok := res.(*messages.SecurityList)
	if !ok {
		t.Fatal("incorrect type assertion")
	}
	var s *models.Security
	for _, sec := range response.Securities {
		if sec.Exchange.Name == "uniswapv3" && sec.Symbol == "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640" {
			s = sec
		}
	}
	if s == nil {
		t.Fatal("missing security")
	}
	res, err = as.Root.RequestFuture(ex, &messages.HistoricalUnipoolV3DataRequest{
		RequestID: 0,
		Instrument: &models.Instrument{
			SecurityID: &wrapperspb.UInt64Value{Value: s.SecurityID},
			Exchange:   s.Exchange,
			Symbol:     &wrapperspb.StringValue{Value: s.Symbol},
		},
		Start: 14268513 - 100,
		End:   14268513,
	}, 50*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	uniResponse, ok := res.(*messages.HistoricalUnipoolV3DataResponse)
	if !ok {
		t.Fatal("incorrect type assertion")
	}
	if !uniResponse.Success {
		t.Fatal(uniResponse.RejectionReason.String())
	}
	for _, trade := range uniResponse.Events {
		fmt.Println(trade)
	}

	pro, err := as.Root.RequestFuture(ex, &messages.ProtocolAssetListRequest{
		RequestID: uint64(time.Now().UnixNano()),
	}, 15*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	msg, ok := pro.(*messages.ProtocolAssetList)
	if !ok {
		t.Fatal("incorrect type assertiob")
	}
	if !msg.Success {
		t.Fatal(msg.RejectionReason.String())
	}
	var a *models.ProtocolAsset
	for _, asset := range msg.ProtocolAssets {
		if asset.Asset.Symbol == "BAYC" {
			a = asset
		}
	}
	if a == nil {
		t.Fatal("Missing asset")
	}
	r, err := as.Root.RequestFuture(ex, &messages.HistoricalProtocolAssetTransferRequest{
		RequestID:       uint64(time.Now().UnixNano()),
		ProtocolAssetID: a.ProtocolAssetID,
		Start:           14268513 - 500,
		Stop:            14268513,
	}, 40*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	events, ok := r.(*messages.HistoricalProtocolAssetTransferResponse)
	if !ok {
		t.Fatal("incorrect type assertiob")
	}
	if !events.Success {
		t.Fatal(events.RejectionReason.String())
	}
	for _, e := range events.Update {
		fmt.Println(e)
	}
}
