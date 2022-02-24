package erc721_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/protocols/tests"
	"gitlab.com/alphaticks/go-graphql-client"
	"gitlab.com/alphaticks/gorderbook"
	models2 "gitlab.com/alphaticks/xchanger/models"
)

type ERC721Data struct {
	Name      string
	Tokens    []*Token
	Transfers []*Transfer `graphql:"transfers(where:{timestamp_lte:$timestamp})"`
}

type Transfer struct {
	From      Owner
	To        Owner
	Timestamp string
}

type Token struct {
	Id    string
	Owner Owner
}

type Owner struct {
	Id string
}

type ERC721Contract struct {
	ERC721Contract ERC721Data `graphql:"erc721Contract(block:{number:$number} id:$id)"`
}

func TestExecutor(t *testing.T) {
	as, executor, cancel := tests.StartExecutor(t, &models2.Protocol{Name: "ERC-721", ID: 0x01})
	defer cancel()

	testAsset := []models.ProtocolAsset{{
		Symbol: "BAYC",
	}}
	res, err := as.Root.RequestFuture(executor, &messages.ProtocolAssetListRequest{}, 20*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	assets, ok := res.(*messages.ProtocolAssetListResponse)
	if !ok {
		t.Fatal("incorrect typying")
	}
	var coll models.ProtocolAsset
	for _, asset := range assets.ProtocolAssets {
		for _, c := range testAsset {
			if asset.Symbol == c.Symbol {
				coll = *asset
			}
		}
	}
	//Execute the future request for the NFT historical data
	resp, err := as.Root.RequestFuture(
		executor,
		&messages.HistoricalProtocolAssetTransferRequest{
			RequestID: uint64(time.Now().UnixNano()),
			ProtocolAsset: &models.ProtocolAsset{
				Address:     coll.Address,
				Name:        coll.Name,
				Symbol:      coll.Symbol,
				TotalSupply: big.NewInt(0).Bytes(),
			},
			Start: 12200000,
			Stop:  12300000,
		},
		30*time.Second,
	).Result()
	if err != nil {
		t.Fatal()
	}
	response, _ := resp.(*messages.HistoricalProtocolAssetTransferResponse)
	if !response.Success {
		t.Fatal("error in the transfers request", response.RejectionReason)
	}
	//Execute the graphql query for the same period querying transfers for the BAYC contract
	client := graphql.NewClient("https://api.thegraph.com/subgraphs/name/ryry79261/mainnet-erc721-erc1155", nil)
	query := ERC721Contract{}
	variables := map[string]interface{}{
		"id":        graphql.ID("0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d"),
		"number":    graphql.Int(12300000),
		"timestamp": graphql.Int(1619228194), //Exact cut off time for block 12300000
	}
	err = client.Query(context.Background(), &query, variables)
	if err != nil {
		t.Fatal(err)
	}
	for _, resT := range response.Update {
		if ok := find(resT, query.ERC721Contract.Transfers); !ok {
			t.Fatal()
		}
	}
	//Run the transfers using the nft tracker
	tracker := gorderbook.NewNftTracker()
	for _, t := range response.Update {
		var from [20]byte
		var to [20]byte
		tokenID := big.NewInt(1)
		copy(from[:], t.Transfer.From)
		copy(to[:], t.Transfer.To)
		tokenID.SetBytes(t.Transfer.TokenId)
		tracker.TransferFrom(from, to, tokenID, "")
	}
	//Capture snapshot of the tracker and extract owners and nft amount
	snap := tracker.GetSnapshot()
	owners := make(map[[20]byte]int32)
	for k, v := range snap.Coins {
		owners[v.Owner] += 1
		fmt.Println("Token", big.NewInt(1).SetBytes(k[:]), "has owner", common.Bytes2Hex(v.Owner[:]))
	}
	//Check owners using the transfers, from graphql
	ownersCheck := make(map[[20]byte]int32)
	for _, t := range query.ERC721Contract.Transfers {
		toOther, _ := big.NewInt(1).SetString(t.To.Id[2:], 16)
		fromOther, _ := big.NewInt(1).SetString(t.From.Id[2:], 16)
		var addTo [20]byte
		var addFrom [20]byte
		copy(addTo[:], toOther.Bytes())
		copy(addFrom[:], fromOther.Bytes())
		ownersCheck[addTo] += 1
		ownersCheck[addFrom] -= 1
	}
	//Compare the final amount of owners and amount of token owned
	for k, v := range owners {
		if ownersCheck[k] != v {
			fmt.Println("error with", common.Bytes2Hex(k[:]))
		}
	}
}

func find(t *models.ProtocolAssetUpdate, arr []*Transfer) bool {
	from := big.NewInt(1).SetBytes(t.Transfer.From)
	to := big.NewInt(1).SetBytes(t.Transfer.To)
	for _, tOther := range arr {
		fromOther, _ := big.NewInt(1).SetString(tOther.From.Id[2:], 16)
		toOther, _ := big.NewInt(1).SetString(tOther.To.Id[2:], 16)
		if from.Cmp(fromOther) == 0 && to.Cmp(toOther) == 0 {
			return true
		}
	}
	return false
}
