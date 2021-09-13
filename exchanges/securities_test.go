package exchanges

import (
	"encoding/json"
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

func TestSecurities(t *testing.T) {
	as := actor.NewActorSystem()
	exchanges := []*xchangerModels.Exchange{
		&constants.BITMEX,
		&constants.BINANCE,
		&constants.BITFINEX,
		&constants.BITSTAMP,
		&constants.COINBASEPRO,
		&constants.GEMINI,
		&constants.KRAKEN,
		&constants.CRYPTOFACILITIES,
		&constants.OKCOIN,
		&constants.FBINANCE,
		&constants.HITBTC,
		&constants.BITZ,
		&constants.HUOBI,
		&constants.FTX,
		&constants.UPBIT,
	}
	assetLoader := as.Root.Spawn(actor.PropsFromProducer(utils.NewAssetLoaderProducer("gs://patrick-configs/assets.json")))
	_, err := as.Root.RequestFuture(assetLoader, &utils.Ready{}, 10*time.Second).Result()
	if err != nil {
		panic(err)
	}

	executor, _ := as.Root.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(nil, exchanges, nil, xchangerUtils.DefaultDialerPool)), "executor")
	res, err := as.Root.RequestFuture(executor, &messages.SecurityListRequest{}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	securityList, ok := res.(*messages.SecurityList)
	if !ok {
		t.Fatalf("was expecting *messages.SecurityList, got %s", reflect.TypeOf(res).String())
	}
	if !securityList.Success {
		t.Fatal(securityList.RejectionReason.String())
	}
	b, err := json.Marshal(securityList.Securities)
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("securities.json", b, 0644)
}
