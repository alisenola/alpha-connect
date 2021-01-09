package main

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/data/live"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	"gitlab.com/alphaticks/alpha-connect/rpc"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/models"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	tickstore_grpc "gitlab.com/tachikoma.ai/tickstore-grpc"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/credentials"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var done = make(chan os.Signal, 1)

var liveStoreActor *actor.PID
var executorActor *actor.PID
var assetLoader *actor.PID

type GuardActor struct{}

func (state *GuardActor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		context.Watch(liveStoreActor)
		context.Watch(executorActor)
		context.Watch(assetLoader)

	case *actor.Terminated:
		done <- os.Signal(syscall.SIGTERM)
	}
}

func main() {
	as := actor.NewActorSystem()
	ctx := actor.NewRootContext(as, nil)

	// Start actors
	exch := []*models.Exchange{
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
		&constants.BITMEX,
		&constants.BITSTAMP,
		&constants.DERIBIT,
		&constants.HUOBIP,
		&constants.HUOBIF,
		&constants.BYBITI,
		&constants.BYBITL,
	}
	// EXECUTOR //
	assetLoader = ctx.Spawn(actor.PropsFromProducer(utils.NewAssetLoaderProducer("gs://patrick-configs/assets.json")))
	_, err := ctx.RequestFuture(assetLoader, &utils.Ready{}, 10*time.Second).Result()
	if err != nil {
		panic(err)
	}
	executorActor, _ = ctx.SpawnNamed(actor.PropsFromProducer(exchanges.NewExecutorProducer(exch, nil, false, xchangerUtils.DefaultDialerPool)), "executor")
	liveStoreActor, _ = ctx.SpawnNamed(actor.PropsFromProducer(live.NewLiveStoreProducer(0)), "live_store")

	// Spawn guard actor
	guardActor, err := ctx.SpawnNamed(
		actor.PropsFromProducer(func() actor.Actor { return &GuardActor{} }),
		"guard_actor")
	if err != nil {
		panic(err)
	}

	// TODO only if in config
	// Start live store gRPC server
	address := "127.0.0.1:7965"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()

	liveER := rpc.NewLiveER(ctx, as.NewLocalPID("live_store"))

	go func() {
		err := server.Serve(lis)
		if err != nil {
			fmt.Println("ERROR", err)
		}
		done <- os.Signal(syscall.SIGTERM)
	}()

	tickstore_grpc.RegisterStoreServer(server, liveER)

	// If interrupt or terminate signal is received, stop
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	// Stop gRPC server
	c := make(chan bool, 1)
	go func() {
		server.GracefulStop()
		c <- true
	}()

	select {
	case <-c:
	case <-time.After(time.Second * 10):
		server.Stop()
	}

	// Stop guard actor first
	_ = ctx.PoisonFuture(guardActor).Wait()

	tries := 0
	for {
		if err := ctx.PoisonFuture(executorActor).Wait(); err == nil {
			break
		}
		tries += 1
		if tries > 80 {
			panic("error shutting down executor actor")
		}
	}
}
