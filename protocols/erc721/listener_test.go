package erc721_test

import (
	"fmt"
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
	as, ex, _, clean := tests.StartExecutor(t, &protocol)
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
	assetTest := models.ProtocolAsset{
		Asset: &models2.Asset{
			Symbol: "BAYC",
		},
	}
	var asset *models.ProtocolAsset
	for _, a := range response.ProtocolAssets {
		if assetTest.Asset.Symbol == a.Asset.Symbol {
			asset = a
		}
	}
	if asset == nil {
		t.Fatal("asset not found")
	}
	fmt.Println("asset", asset)
	props := actor.PropsFromProducer(tests.NewERC721CheckerProducer(asset))
	checker := as.Root.Spawn(props)
	defer as.Root.Poison(checker)
	time.Sleep(3000 * time.Second)
}
