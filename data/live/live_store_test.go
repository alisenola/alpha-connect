package live

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	exchanges2 "gitlab.com/alphaticks/alphac/exchanges"
	"gitlab.com/alphaticks/xchanger/constants"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	tickstore_grpc "gitlab.com/tachikoma.ai/tickstore-grpc"
	"testing"
	"time"
)

func TestDataProvider(t *testing.T) {
	exchanges := []*xchangerModels.Exchange{
		&constants.BITMEX,
	}
	executor, _ := actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(exchanges2.NewExecutorProducer(exchanges, nil, false)), "executor")
	liveStore, _ := actor.EmptyRootContext.SpawnNamed(actor.PropsFromProducer(NewLiveStoreProducer(0)), "live_store")

	res, err := actor.EmptyRootContext.RequestFuture(liveStore,
		&storeMessages.GetQueryReaderRequest{
			Query: &tickstore_grpc.QueryRequest{
				Streaming: true,
				BatchSize: 100,
				Timeout:   uint64(time.Millisecond * 300),
				From:      0,
				To:        0,
				Selector:  `SELECT OBPrice(orderbook, "0.0001") WHERE exchange="bitmex" base="BTC" quote="USD"`,
				Sampler:   nil,
			},
			RequestID: 0,
		}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}

	response := res.(*storeMessages.GetQueryReaderResponse)
	if response.Error != nil {
		t.Fatal(err)
	}

	res, err = actor.EmptyRootContext.RequestFuture(response.Reader,
		&storeMessages.ReadQueryEventBatchRequest{
			RequestID: 0,
			Timeout:   uint64(time.Millisecond * 300),
			BatchSize: 100,
		}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
	res, err = actor.EmptyRootContext.RequestFuture(response.Reader,
		&storeMessages.ReadQueryEventBatchRequest{
			RequestID: 0,
			Timeout:   uint64(time.Millisecond * 300),
			BatchSize: 100,
		}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)

	time.Sleep(10 * time.Second)

	actor.EmptyRootContext.Poison(executor)
	actor.EmptyRootContext.Poison(liveStore)
}
