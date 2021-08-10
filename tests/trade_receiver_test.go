package tests

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"testing"
	"time"
)

var FBinanceAccount = &models.Account{
	AccountID: "299211",
	Exchange:  &constants.FBINANCE,
	Credentials: &xchangerModels.APICredentials{
		APIKey:    "MpYkeK3pGP80gGiIrqWtLNwjJmyK2DTREYzNx8Cyc3AWTkl2T0iWnQEtdCIlvAoE",
		APISecret: "CJcJZEkktzhGzEdQhclfHcfJz5k01OY6n42MeF9B3oQWGqba3RrXEnG4bZktXQNu",
	},
}

func TestTradeReceiver(t *testing.T) {
	exchgs := []*xchangerModels.Exchange{
		&constants.FBINANCE,
	}
	as := actor.NewActorSystem()
	assetLoader := as.Root.Spawn(actor.PropsFromProducer(utils.NewAssetLoaderProducer("../assets.json")))
	_, err := as.Root.RequestFuture(assetLoader, &utils.Ready{}, 10*time.Second).Result()
	if err != nil {
		panic(err)
	}
	executor, _ := as.Root.SpawnNamed(actor.PropsFromProducer(exchanges.NewExecutorProducer(exchgs, xchangerUtils.DefaultDialerPool)), "executor")
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
			}, 0, uint64(time.Now().UnixNano()), ch)
			receivers = append(receivers, receiver)
		}
	}

	for {
		msg := <-ch
		fmt.Println(msg)
	}
}
