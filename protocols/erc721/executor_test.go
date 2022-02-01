package erc721

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/go-graphql-client"
	"gitlab.com/alphaticks/gorderbook"
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
	as := actor.NewActorSystem()
	ex, err := as.Root.SpawnNamed(actor.PropsFromProducer(
		func() actor.Actor {
			return NewExecutor()
		},
	), "executor_erc")
	if err != nil {
		t.Fatal(err)
	}
	contract, ok := big.NewInt(1).SetString("bc4ca0eda7647a8ab7c2061c2e118a18a936f13d", 16)
	if !ok {
		t.Fatal("error in conversion from string to hexa")
	}
	resp, err := as.Root.RequestFuture(
		ex,
		&messages.HistoricalNftTransferDataRequest{
			RequestID: uint64(time.Now().UnixNano()),
			Collection: &models.Collection{
				Address:     contract.Bytes(),
				Name:        "BoredApeYachtClub",
				Symbol:      "BAYC",
				TotalSupply: big.NewInt(10000).Bytes(),
			},
			Start: 12200000,
			Stop:  12300000,
		},
		30*time.Second,
	).Result()
	if err != nil {
		t.Fatal()
	}
	response, _ := resp.(*messages.HistoricalNftTransferDataResponse)
	if !response.Success {
		t.Fatal("error in the transfers request")
	}

	client := graphql.NewClient("https://api.thegraph.com/subgraphs/name/ryry79261/mainnet-erc721-erc1155", nil)
	query := ERC721Contract{}
	variables := map[string]interface{}{
		"id":        graphql.ID("0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d"),
		"number":    graphql.Int(12300000),
		"timestamp": graphql.Int(1619228194),
	}
	err = client.Query(context.Background(), &query, variables)
	if err != nil {
		t.Fatal(err)
	}
	for _, resT := range response.Transfers {
		if ok := find(resT, query.ERC721Contract.Transfers); !ok {
			t.Fatal()
		}
	}

	tracker := gorderbook.NewNftTracker()
	for _, t := range response.Transfers {
		var from [20]byte
		var to [20]byte
		tokenID := big.NewInt(1)
		copy(from[:], t.Transfer.From)
		copy(to[:], t.Transfer.To)
		tokenID.SetBytes(t.Transfer.TokenId)
		tracker.TransferFromNft(from, to, tokenID)
	}
	snap := tracker.GetNftTrackerSnapshot()
	for k, v := range snap.Coins {
		fmt.Println("Token", big.NewInt(1).SetBytes(k[:]), "has owner", common.Bytes2Hex(v.Owner[:]))
	}
}

func find(t *models.NftTransfer, arr []*Transfer) bool {
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
