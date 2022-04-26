package opensea_test

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	"reflect"
	"testing"
	"time"

	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	models2 "gitlab.com/alphaticks/xchanger/models"
)

func TestExecutor(t *testing.T) {
	as, executor, cancel := tests.StartExecutor(t, &models2.Exchange{ID: 0x23, Name: "opensea"}, nil)
	defer cancel()
	testAsset := models2.Asset{
		ID:     275.,
		Symbol: "BAYC",
	}
	res, err := as.Root.RequestFuture(executor, &messages.MarketableAssetListRequest{}, 20*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	assets, ok := res.(*messages.MarketableAssetList)
	if !ok {
		t.Fatal("incorrect for type assertion")
	}
	var coll *models.MarketableAsset
	for _, asset := range assets.MarketableAssets {
		fmt.Printf("asset %+v \n", asset)
		if asset.ProtocolAsset.Asset.Symbol == testAsset.Symbol {
			coll = asset
		}
	}
	if coll == nil {
		t.Fatal("missing collection")
	}
	r, err := as.Root.RequestFuture(executor, &messages.HistoricalSalesRequest{
		RequestID: uint64(time.Now().UnixNano()),
		AssetID:   coll.ProtocolAsset.Asset.ID,
		From:      &types.Timestamp{Seconds: 1650891168},
	}, 15*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	sales, ok := r.(*messages.HistoricalSalesResponse)
	if !ok {
		t.Fatalf("was expecting HistoricalSalesResponse got %s", reflect.TypeOf(r).String())
	}
	if !sales.Success {
		t.Fatal(sales.RejectionReason)
	}
	fmt.Println("Sales", sales.Sale)
}
