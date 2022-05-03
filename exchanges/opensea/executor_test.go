package opensea_test

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/xchanger/constants"
	"io/ioutil"
	"math/big"
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
	res, err := as.Root.RequestFuture(executor, &messages.MarketableProtocolAssetListRequest{}, 20*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	assets, ok := res.(*messages.MarketableProtocolAssetList)
	if !ok {
		t.Fatal("incorrect for type assertion")
	}
	var coll *models.MarketableProtocolAsset
	for _, asset := range assets.MarketableProtocolAssets {
		fmt.Printf("asset %+v \n", asset)
		if asset.ProtocolAsset.Asset.Symbol == testAsset.Symbol {
			coll = asset
		}
	}
	if coll == nil {
		t.Fatal("missing collection")
	}
	r, err := as.Root.RequestFuture(executor, &messages.HistoricalSalesRequest{
		RequestID:                 uint64(time.Now().UnixNano()),
		MarketableProtocolAssetID: coll.MarketableProtocolAssetID,
		From:                      &types.Timestamp{Seconds: 0},
	}, 5*time.Minute).Result()
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
	type SaleFile struct {
		From        common.Address
		To          common.Address
		TokenID     *big.Int
		Price       *big.Int
		BlockNumber uint64
	}
	var s []SaleFile
	for _, sale := range sales.Sale {
		s = append(s, SaleFile{
			From:        common.BytesToAddress(sale.Transfer.From),
			To:          common.BytesToAddress(sale.Transfer.To),
			TokenID:     big.NewInt(1).SetBytes(sale.Transfer.TokenId),
			Price:       big.NewInt(1).SetBytes(sale.Price),
			BlockNumber: sale.Block,
		})
	}
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("./sales_test.json", b, 0644)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("TimeStamp", sales.Sale[0].Timestamp)
}

func TestMarketableProtocolAssetID(t *testing.T) {
	market := constants.OPENSEA
	protocolAssetID := uint64(12345678)
	id := utils.MarketableProtocolAssetID(protocolAssetID, market.ID)
	for i := 0; i < 100; i++ {
		tmp := utils.MarketableProtocolAssetID(protocolAssetID, market.ID)
		if tmp != id {
			t.Fatal("failed")
		}
	}
}
