package erc20_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/alphaticks/alpha-connect/config"
	exTests "gitlab.com/alphaticks/alpha-connect/tests"
	"gitlab.com/alphaticks/alpha-connect/utils"
	gorderbookModels "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"

	"gitlab.com/alphaticks/xchanger/constants"
	"reflect"
	"testing"
	"time"

	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
)

func TestExecutorSVM(t *testing.T) {
	exTests.LoadStatics(t)
	protocol := constants.ERC20
	registryAddress := "127.0.0.1:8001"
	cfg := config.Config{
		RegistryAddress: registryAddress,
		Protocols:       []string{protocol.Name},
	}
	as, executor, clean := exTests.StartExecutor(t, &cfg)
	defer clean()

	testAsset := xchangerModels.Asset{
		ID:     57.,
		Symbol: "DAI",
	}
	chain, ok := constants.GetChainByID(5)
	assert.True(t, ok, "missings svm")
	res, err := as.Root.RequestFuture(executor, &messages.ProtocolAssetListRequest{}, 20*time.Second).Result()
	assert.Nil(t, err, "RequestFuture ProtocolAssetList err: %v", err)
	assets, ok := res.(*messages.ProtocolAssetList)
	assert.True(t, ok, "incorrect type assertion")
	var coll *models.ProtocolAsset
	for _, asset := range assets.ProtocolAssets {
		if asset.Asset.Symbol == testAsset.Symbol && asset.Chain.ID == chain.ID {
			fmt.Printf("asset %+v \n", asset)
			coll = asset
		}
	}
	assert.NotNil(t, coll, "missing collection")
	r, err := as.Root.RequestFuture(executor, &messages.ProtocolAssetDefinitionRequest{
		RequestID:       uint64(time.Now().UnixNano()),
		ProtocolAssetID: utils.GetProtocolAssetID(&testAsset, constants.ERC20, constants.StarknetMainnet),
	}, 15*time.Second).Result()
	assert.Nil(t, err, "RequestFuture ProtocolAssetDefinition err: %v", err)
	def, ok := r.(*messages.ProtocolAssetDefinitionResponse)
	assert.True(t, ok, "expected ProtocolAssetDefinitionResponse, got %s", reflect.TypeOf(r).String())
	assert.True(t, def.Success, "request failed with %v", def.RejectionReason.String())
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
	assert.Nil(t, err, "RequestFuture HistoricalProtocolAssetTransferRequest err: %v", err)
	response, ok := resp.(*messages.HistoricalProtocolAssetTransferResponse)
	assert.True(t, ok, "expected HistoricalProtocolAssetTransferResponse, got %s", reflect.TypeOf(resp).String())
	assert.True(t, response.Success, "request failed with %s", response.RejectionReason.String())
	var events []*gorderbookModels.AssetTransfer
	for _, ev := range response.Update {
		events = append(events, ev.Transfers...)
	}
	assert.GreaterOrEqual(t, len(events), 7000, "expected more than 20 events")
}
