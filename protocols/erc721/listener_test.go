package erc721_test

import (
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/protocols/tests"
	models2 "gitlab.com/alphaticks/xchanger/models"
)

func TestListener(t *testing.T) {
	protocol := models2.Protocol{Name: "ERC-721", ID: 0x01}
	as, ex, clean := tests.StartExecutor(t, &protocol)
	defer clean()

	res, err := as.Root.RequestFuture(ex, &messages.AssetListRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Subscribe: false,
	}, 20*time.Second).Result()
	if err != nil {
		t.Fatal()
	}
	response, ok := res.(*messages.AssetListResponse)
	if !ok {
		t.Fatal("incorrect type assertion")
	}
	assetTest := models.ProtocolAsset{
		Symbol: "BAYC",
	}
	var asset *models.ProtocolAsset
	for _, a := range response.Assets {
		if assetTest.Symbol == a.Symbol {
			asset = a
		}
	}
	if asset == nil {
		t.Fatal("asset not found")
	}
	props := actor.PropsFromProducer(tests.NewERC721CheckerProducer(asset))
	checker := as.Root.Spawn(props)
	as.Root.PoisonFuture(checker)
	time.Sleep(50 * time.Second)
}
