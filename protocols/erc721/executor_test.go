package erc721_test

import (
	"context"
	"fmt"
	ctypes "gitlab.com/alphaticks/alpha-connect/chains/types"
	xtypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/protocols/types"
	exTests "gitlab.com/alphaticks/alpha-connect/tests"
	"gitlab.com/alphaticks/alpha-connect/utils"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	gorderbookModels "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"

	"gitlab.com/alphaticks/xchanger/constants"
	"google.golang.org/grpc"
	"math/big"
	"reflect"
	"testing"
	"time"

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

func TestExecutorEVM(t *testing.T) {
	exTests.LoadStatics(t)
	protocol := constants.ERC721
	registryAddress := "registry.alphaticks.io:8001"
	conn, err := grpc.Dial(registryAddress, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	reg := registry.NewPublicRegistryClient(conn)

	prCfg := &types.ExecutorConfig{
		Registry:  reg,
		Protocols: []*xchangerModels.Protocol{protocol},
	}
	as, executor, clean := exTests.StartExecutor(t, &xtypes.ExecutorConfig{}, prCfg, &ctypes.ExecutorConfig{}, nil)
	defer clean()

	testAsset := xchangerModels.Asset{
		ID:     275.,
		Symbol: "BAYC",
	}
	res, err := as.Root.RequestFuture(executor, &messages.ProtocolAssetListRequest{}, 20*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	assets, ok := res.(*messages.ProtocolAssetList)
	if !ok {
		t.Fatal("incorrect for type assertion")
	}
	var coll *models.ProtocolAsset
	for _, asset := range assets.ProtocolAssets {
		fmt.Printf("asset %+v \n", asset)
		if asset.Asset.Symbol == testAsset.Symbol {
			coll = asset
		}
	}
	if coll == nil {
		t.Fatal("missing collection")
	}
	r, err := as.Root.RequestFuture(executor, &messages.ProtocolAssetDefinitionRequest{
		RequestID:       uint64(time.Now().UnixNano()),
		ProtocolAssetID: utils.GetProtocolAssetID(&testAsset, constants.ERC721, constants.EthereumMainnet),
	}, 15*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	def, ok := r.(*messages.ProtocolAssetDefinitionResponse)
	if !ok {
		t.Fatalf("was expecting ProtocolAssetDefinitionResponse got %s", reflect.TypeOf(r).String())
	}
	if !def.Success {
		t.Fatal(def.RejectionReason)
	}
	fmt.Println("Protocol Asset definition", def.ProtocolAsset)
	//Execute the future request for the NFT historical data
	resp, err := as.Root.RequestFuture(
		executor,
		&messages.HistoricalProtocolAssetTransferRequest{
			RequestID:       uint64(time.Now().UnixNano()),
			ProtocolAssetID: coll.ProtocolAssetID,
			Start:           12200000,
			Stop:            12300000,
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
		for _, transfer := range resT.Transfers {
			if ok := find(transfer, query.ERC721Contract.Transfers); !ok {
				t.Fatal()
			}
		}
	}
	//Run the transfers using the nft tracker
	tracker := gorderbook.NewERC721Tracker()
	for _, updt := range response.Update {
		for _, transfer := range updt.Transfers {
			var from [20]byte
			var to [20]byte
			tokenID := big.NewInt(1)
			copy(from[:], transfer.From)
			copy(to[:], transfer.To)
			tokenID.SetBytes(transfer.TokenId)
			if err := tracker.TransferFrom(from, to, tokenID); err != nil {
				t.Fatal(err)
			}
		}
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

func TestExecutorSVM(t *testing.T) {
	exTests.LoadStatics(t)
	protocol := constants.ERC721
	registryAddress := "127.0.0.1:8001"
	conn, err := grpc.Dial(registryAddress, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	reg := registry.NewPublicRegistryClient(conn)

	prCfg := &types.ExecutorConfig{
		Registry:  reg,
		Protocols: []*xchangerModels.Protocol{protocol},
	}
	as, executor, clean := exTests.StartExecutor(t, &xtypes.ExecutorConfig{}, prCfg, &ctypes.ExecutorConfig{}, nil)
	defer clean()

	testAsset := xchangerModels.Asset{
		ID:     275.,
		Symbol: "BAYC",
	}
	chain, ok := constants.GetChainByID(5)
	if !ok {
		t.Fatal("missing svm")
	}
	res, err := as.Root.RequestFuture(executor, &messages.ProtocolAssetListRequest{}, 20*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	assets, ok := res.(*messages.ProtocolAssetList)
	if !ok {
		t.Fatal("incorrect for type assertion")
	}
	var coll *models.ProtocolAsset
	for _, asset := range assets.ProtocolAssets {
		if asset.Asset.Symbol == testAsset.Symbol && asset.Chain.ID == chain.ID {
			fmt.Printf("asset %+v \n", asset)
			coll = asset
		}
	}
	if coll == nil {
		t.Fatal("missing collection")
	}
	r, err := as.Root.RequestFuture(executor, &messages.ProtocolAssetDefinitionRequest{
		RequestID:       uint64(time.Now().UnixNano()),
		ProtocolAssetID: utils.GetProtocolAssetID(&testAsset, constants.ERC721, constants.StarnetMainnet),
	}, 15*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	def, ok := r.(*messages.ProtocolAssetDefinitionResponse)
	if !ok {
		t.Fatalf("was expecting ProtocolAssetDefinitionResponse got %s", reflect.TypeOf(r).String())
	}
	if !def.Success {
		t.Fatal(def.RejectionReason)
	}
	fmt.Println("Protocol Asset definition", def.ProtocolAsset)
	//Execute the future request for the NFT historical data
	resp, err := as.Root.RequestFuture(
		executor,
		&messages.HistoricalProtocolAssetTransferRequest{
			RequestID:       uint64(time.Now().UnixNano()),
			ProtocolAssetID: coll.ProtocolAssetID,
			Start:           0,
			Stop:            3000,
		},
		30*time.Second,
	).Result()
	if err != nil {
		t.Fatal()
	}
	response, ok := resp.(*messages.HistoricalProtocolAssetTransferResponse)
	if !ok {
		t.Fatalf("expected *messages.HistoricalProtocolAssetTransferResponse, got %s", reflect.TypeOf(resp).String())
	}
	if !response.Success {
		t.Fatal("error in the transfers request", response.RejectionReason)
	}
}

func find(t *gorderbookModels.AssetTransfer, arr []*Transfer) bool {
	from := big.NewInt(1).SetBytes(t.From)
	to := big.NewInt(1).SetBytes(t.To)
	for _, tOther := range arr {
		fromOther, _ := big.NewInt(1).SetString(tOther.From.Id[2:], 16)
		toOther, _ := big.NewInt(1).SetString(tOther.To.Id[2:], 16)
		if from.Cmp(fromOther) == 0 && to.Cmp(toOther) == 0 {
			return true
		}
	}
	return false
}
