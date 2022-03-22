package utils_test

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/account"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"google.golang.org/grpc"
	"testing"
	"time"
)

var FBinanceAccount = &models.Account{
	Name:     "299211",
	Exchange: &constants.FBINANCE,
	ApiCredentials: &xchangerModels.APICredentials{
		APIKey:    "MpYkeK3pGP80gGiIrqWtLNwjJmyK2DTREYzNx8Cyc3AWTkl2T0iWnQEtdCIlvAoE",
		APISecret: "CJcJZEkktzhGzEdQhclfHcfJz5k01OY6n42MeF9B3oQWGqba3RrXEnG4bZktXQNu",
	},
}

func TestTradeReceiver(t *testing.T) {
	exchgs := []*xchangerModels.Exchange{
		&constants.FBINANCE,
	}
	accnt, err := account.NewAccount(FBinanceAccount)
	if err != nil {
		t.Fatal(err)
	}
	as := actor.NewActorSystem()
	registryAddress := "registry.alphaticks.io:7001"
	conn, err := grpc.Dial(registryAddress, grpc.WithInsecure())
	if err != nil {
		err := fmt.Errorf("error connecting to public registry gRPC endpoint: %v", err)
		panic(err)
	}
	rgstr := registry.NewPublicRegistryClient(conn)
	assetLoader := as.Root.Spawn(actor.PropsFromProducer(utils.NewAssetLoaderProducer(rgstr)))
	_, err = as.Root.RequestFuture(assetLoader, &utils.Ready{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	cfg := &exchanges.ExecutorConfig{
		Db:         nil,
		Exchanges:  exchgs,
		Accounts:   []*account.Account{accnt},
		DialerPool: xchangerUtils.DefaultDialerPool,
		Strict:     false,
	}
	executor, _ := as.Root.SpawnNamed(actor.PropsFromProducer(exchanges.NewExecutorProducer(cfg)), "executor")
	res, err := as.Root.RequestFuture(executor, &messages.SecurityListRequest{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	securities := res.(*messages.SecurityList)
	ch := make(chan *models.TradeCapture)
	var receivers []*utils.TradeReceiver
	for _, sec := range securities.Securities {
		if sec.Symbol == "BZRXUSDT" {
			receiver := utils.NewTradeReceiver(as, FBinanceAccount, &models.Instrument{
				SecurityID: &types.UInt64Value{Value: sec.SecurityID},
				Exchange:   sec.Exchange,
				Symbol:     &types.StringValue{Value: sec.Symbol},
			}, 0, false, uint64(time.Now().UnixNano()), ch)
			receivers = append(receivers, receiver)
		}
	}
	fmt.Println(receivers)
	for {
		msg := <-ch
		fmt.Println(msg)
	}
}
