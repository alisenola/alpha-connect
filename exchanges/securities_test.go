package exchanges

import (
	"encoding/json"
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

func TestSecurities(t *testing.T) {
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
	}
	executor, _ = actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewExecutorProducer(exchanges, nil, false)), "executor")
	res, err := actor.EmptyRootContext.RequestFuture(executor, &messages.SecurityListRequest{}, 10*time.Second).Result()
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
