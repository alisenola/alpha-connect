package opensea_test

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/xchanger/constants"
	"io/ioutil"
	"math/big"
	"reflect"
	"sort"
	"testing"
	"time"

	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	models2 "gitlab.com/alphaticks/xchanger/models"
)

func TestExecutor(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	as, executor, cancel := tests.StartExecutor(t, &models2.Exchange{ID: 0x23, Name: "opensea"}, nil)
	defer cancel()
	testAsset := models2.Asset{
		ID:     276.,
		Symbol: "MEEBTS",
	}
	//testAsset := models2.Asset{
	//	ID:     275.,
	//	Symbol: "BAYC",
	//}
	res, err := as.Root.RequestFuture(executor, &messages.MarketableProtocolAssetListRequest{}, 60*time.Second).Result()
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
	now := time.Now().Unix()
	step := int64(86400 / 3)
	from := int64(1620054232)
	done := false
	var sales *messages.HistoricalSalesResponse
	highestTs := uint64(0)
	bundleCount := 0
OUTER:
	for !done {
		r, err := as.Root.RequestFuture(executor, &messages.HistoricalSalesRequest{
			RequestID:                 uint64(time.Now().UnixNano()),
			MarketableProtocolAssetID: coll.MarketableProtocolAssetID,
			From:                      utils.SecondToTimestamp(uint64(from)),
			To:                        utils.SecondToTimestamp(uint64(from + step)),
		}, 5*time.Minute).Result()
		if err != nil {
			t.Fatal(err)
		}
		sales, ok = r.(*messages.HistoricalSalesResponse)
		if !ok {
			t.Fatalf("was expecting HistoricalSalesResponse got %s", reflect.TypeOf(r).String())
		}
		if !sales.Success {
			t.Fatal(sales.RejectionReason)
		}
		sort.Slice(sales.Sale, func(i, j int) bool {
			return utils.TimestampToMilli(sales.Sale[i].Timestamp) < utils.TimestampToMilli(sales.Sale[j].Timestamp)
		})
		done = from >= now
		from += step + 1
		for _, sale := range sales.Sale {
			if len(sale.Transfer) > 1 {
				fmt.Println("bundle")
				if bundleCount > 5 {
					break OUTER
				}
				bundleCount++
			}
		}
		sort.Slice(sales.Sale, func(i, j int) bool {
			return utils.TimestampToMilli(sales.Sale[i].Timestamp) < utils.TimestampToMilli(sales.Sale[j].Timestamp)
		})
		for _, s := range sales.Sale {
			if utils.TimestampToMilli(s.Timestamp) < highestTs {
				t.Fatal("mismatched timestamp")
			}
			highestTs = utils.TimestampToMilli(s.Timestamp)
		}
	}
	type SaleFile struct {
		From        common.Address
		To          common.Address
		TokenID     *big.Int
		Price       *big.Int
		SaleId      uint64
		BlockNumber uint64
	}
	var s []SaleFile
	for _, sale := range sales.Sale {
		for _, tr := range sale.Transfer {
			s = append(s, SaleFile{
				From:        common.BytesToAddress(tr.From),
				To:          common.BytesToAddress(tr.To),
				TokenID:     big.NewInt(1).SetBytes(tr.TokenId),
				Price:       big.NewInt(1).SetBytes(sale.Price),
				SaleId:      sale.Id,
				BlockNumber: sale.Block,
			})
		}
	}
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("./sales_test.json", b, 0644)
	if err != nil {
		t.Fatal(err)
	}
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
