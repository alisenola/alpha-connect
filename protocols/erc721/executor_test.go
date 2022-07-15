package erc721_test

import (
	"context"
	"fmt"
	"github.com/melaurent/gotickfile/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.com/alphaticks/alpha-connect/config"
	exTests "gitlab.com/alphaticks/alpha-connect/tests"
	"gitlab.com/alphaticks/alpha-connect/utils"
	gorderbookModels "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	"gitlab.com/alphaticks/tickfunctors/protocols/erc721/evm"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/go-graphql-client"
	"gitlab.com/alphaticks/xchanger/constants"
	"math/big"
	"reflect"
	"testing"
	"time"
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
	// Load executor, chain and protocol asset
	exTests.LoadStatics(t)
	protocol := constants.ERC721
	testAsset := xchangerModels.Asset{
		ID:     275.,
		Symbol: "BAYC",
	}
	registryAddress := "127.0.0.1:8001"
	cfg, err := config.LoadConfig()
	if !assert.Nil(t, err, "LoadConfig err: %v", err) {
		t.Fatal()
	}
	cfg.RegistryAddress = registryAddress
	cfg.Protocols = []string{protocol.Name}
	as, executor, clean := exTests.StartExecutor(t, cfg)
	defer clean()

	res, err := as.Root.RequestFuture(executor, &messages.ProtocolAssetListRequest{}, 20*time.Second).Result()
	if !assert.Nil(t, err, "RequestFuture ProtocolAssetList err: %v", err) {
		t.Fatal()
	}
	assets, ok := res.(*messages.ProtocolAssetList)
	if !assert.True(t, ok, "incorrect for type assertion") {
		t.Fatal()
	}
	var coll *models.ProtocolAsset
	for _, asset := range assets.ProtocolAssets {
		if asset.Asset.Symbol == testAsset.Symbol {
			coll = asset
		}
	}
	if !assert.NotNil(t, coll, "missing collection") {
		t.Fatal()
	}
	r, err := as.Root.RequestFuture(executor, &messages.ProtocolAssetDefinitionRequest{
		RequestID:       uint64(time.Now().UnixNano()),
		ProtocolAssetID: utils.GetProtocolAssetID(&testAsset, constants.ERC721, constants.EthereumMainnet),
	}, 15*time.Second).Result()
	if !assert.Nil(t, err, "RequestFuture ProtocolAssetDefinition err: %v", err) {
		t.Fatal()
	}
	def, ok := r.(*messages.ProtocolAssetDefinitionResponse)
	if !assert.True(t, ok, "expected ProtocolAssetDefinitionResponse got %s", reflect.TypeOf(r).String()) {
		t.Fatal()
	}
	if !assert.Truef(t, def.Success, "request failed with %s", def.RejectionReason.String()) {
		t.Fatal()
	}
	var start uint64 = 12287507
	var stop uint64 = 12300000
	var step uint64 = 1000
	var updates []*models.ProtocolAssetUpdate
	for start != stop {
		//Execute the future request for the NFT historical data
		resp, err := as.Root.RequestFuture(
			executor,
			&messages.HistoricalProtocolAssetTransferRequest{
				RequestID:  uint64(time.Now().UnixNano()),
				AssetID:    &wrapperspb.UInt32Value{Value: testAsset.ID},
				ProtocolID: def.ProtocolAsset.Protocol.ID,
				ChainID:    def.ProtocolAsset.Chain.ID,
				Start:      start,
				Stop:       start + step,
			},
			30*time.Second,
		).Result()
		if !assert.Nil(t, err, "RequestFuture HistoricalProtocolAssetTransferRequest err: %v", err) {
			t.Fatal()
		}
		response, ok := resp.(*messages.HistoricalProtocolAssetTransferResponse)
		if !assert.True(t, ok, "expected HistoricalProtocolAssetTransferResponse, got %s", reflect.TypeOf(resp).String()) {
			t.Fatal()
		}
		if !assert.True(t, response.Success, "request failed with %s", response.RejectionReason.String()) {
			t.Fatal()
		}
		updates = append(updates, response.Update...)
		start += step + 1
		if start > stop {
			start = stop
		}
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
	if !assert.Nil(t, err, "Query err: %v", err) {
		t.Fatal()
	}
	for _, resT := range updates {
		for _, transfer := range resT.Transfers {
			if !assert.True(t, find(transfer, query.ERC721Contract.Transfers), "failed to find transfer") {
				t.Fatal()
			}
		}
	}
	//Run the transfers using the nft tracker
	tracker := evm.NewERC721Tracker()
	for _, updt := range updates {
		for _, transfer := range updt.Transfers {
			tokenID := big.NewInt(1)
			tokenID.SetBytes(transfer.Value)
			var from [20]byte
			var to [20]byte
			var tok [32]byte
			copy(from[:], transfer.From)
			copy(to[:], transfer.To)
			tokenID.FillBytes(tok[:])
			c := evm.ERC721Delta{
				From:    from,
				To:      to,
				TokenId: tok,
				Block:   updt.BlockNumber,
			}
			dlt := gotickfile.TickDeltas{
				Pointer: unsafe.Pointer(&c),
				Len:     1,
			}
			if !assert.Nil(t, tracker.ProcessDeltas(dlt), "TransferFrom err: %v", err) {
				t.Fatal()
			}
		}
	}
	//Capture snapshot of the tracker and extract owners and nft amount
	snap := tracker.GetTokens()
	owners := make(map[[20]byte]int32)
	for _, v := range snap {
		owners[v] += 1
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

	// test get all historical transfers
	updates = nil
	start = 12287507
	stop = 12287800
	step = 10
	for start != stop {
		//Execute the future request for the NFT historical data
		resp, err := as.Root.RequestFuture(
			executor,
			&messages.HistoricalProtocolAssetTransferRequest{
				RequestID:  uint64(time.Now().UnixNano()),
				ProtocolID: def.ProtocolAsset.Protocol.ID,
				ChainID:    def.ProtocolAsset.Chain.ID,
				Start:      start,
				Stop:       start + step,
			},
			30*time.Second,
		).Result()
		if !assert.Nil(t, err, "RequestFuture HistoricalProtocolAssetTransferRequest err: %v", err) {
			t.Fatal()
		}
		response, ok := resp.(*messages.HistoricalProtocolAssetTransferResponse)
		switch response.RejectionReason.String() {
		case "RPCTimeout":
			step /= 2
			continue
		}
		if !assert.True(t, ok, "expected HistoricalProtocolAssetTransferResponse, got %s", reflect.TypeOf(resp).String()) {
			t.Fatal()
		}
		if !assert.True(t, response.Success, "request failed with %s", response.RejectionReason.String()) {
			t.Fatal()
		}
		updates = append(updates, response.Update...)
		step = 10
		start += step + 1
		if start+step > stop {
			start = stop
		}
	}
}

func TestExecutorSVM(t *testing.T) {
	exTests.LoadStatics(t)
	cfg := config.Config{
		RegistryAddress: "registry.alphaticks.io:8001",
		Protocols:       []string{constants.ERC721.Name},
	}
	as, executor, clean := exTests.StartExecutor(t, &cfg)
	defer clean()

	testAsset := xchangerModels.Asset{
		ID:     275.,
		Symbol: "BAYC",
	}
	chain, ok := constants.GetChainByID(5)
	if !assert.True(t, ok, "missings svm") {
		t.Fatal()
	}
	res, err := as.Root.RequestFuture(executor, &messages.ProtocolAssetListRequest{}, 20*time.Second).Result()
	if !assert.Nil(t, err, "RequestFuture ProtocolAssetList err: %v", err) {
		t.Fatal()
	}
	assets, ok := res.(*messages.ProtocolAssetList)
	if !assert.True(t, ok, "incorrect type assertion") {
		t.Fatal()
	}
	var coll *models.ProtocolAsset
	for _, asset := range assets.ProtocolAssets {
		if asset.Asset.Symbol == testAsset.Symbol && asset.Chain.ID == chain.ID {
			fmt.Printf("asset %+v \n", asset)
			coll = asset
		}
	}
	if !assert.NotNil(t, coll, "missing collection") {
		t.Fatal()
	}
	r, err := as.Root.RequestFuture(executor, &messages.ProtocolAssetDefinitionRequest{
		RequestID:       uint64(time.Now().UnixNano()),
		ProtocolAssetID: utils.GetProtocolAssetID(&testAsset, constants.ERC721, constants.StarknetMainnet),
	}, 15*time.Second).Result()
	if !assert.Nil(t, err, "RequestFuture ProtocolAssetDefinition err: %v", err) {
		t.Fatal()
	}
	def, ok := r.(*messages.ProtocolAssetDefinitionResponse)
	if !assert.True(t, ok, "expected ProtocolAssetDefinitionResponse, got %s", reflect.TypeOf(r).String()) {
		t.Fatal()
	}
	if !assert.True(t, def.Success, "request failed with %v", def.RejectionReason.String()) {
		t.Fatal()
	}
	fmt.Println("Protocol Asset definition", def.ProtocolAsset)
	//Execute the future request for the NFT historical data
	resp, err := as.Root.RequestFuture(
		executor,
		&messages.HistoricalProtocolAssetTransferRequest{
			RequestID:  uint64(time.Now().UnixNano()),
			ProtocolID: def.ProtocolAsset.Protocol.ID,
			ChainID:    def.ProtocolAsset.Chain.ID,
			Start:      0,
			Stop:       3000,
		},
		30*time.Second,
	).Result()
	if !assert.Nil(t, err, "RequestFuture HistoricalProtocolAssetTransferRequest err: %v", err) {
		t.Fatal()
	}
	response, ok := resp.(*messages.HistoricalProtocolAssetTransferResponse)
	if !assert.True(t, ok, "expected HistoricalProtocolAssetTransferResponse, got %s", reflect.TypeOf(resp).String()) {
		t.Fatal()
	}
	if !assert.True(t, response.Success, "request failed with %s", response.RejectionReason.String()) {
		t.Fatal()
	}
	var events []*gorderbookModels.AssetTransfer
	for _, ev := range response.Update {
		events = append(events, ev.Transfers...)
	}
	if !assert.GreaterOrEqual(t, len(events), 20, "expected more than 20 events") {
		t.Fatal()
	}
	//TODO get all historical transfers
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
