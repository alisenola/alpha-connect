package main

import (
	"fmt"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/executor"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	chtypes "gitlab.com/alphaticks/alpha-connect/chains/types"
	"gitlab.com/alphaticks/alpha-connect/data"
	prtypes "gitlab.com/alphaticks/alpha-connect/protocols/types"
	"gitlab.com/alphaticks/alpha-connect/rpc"
	"gitlab.com/alphaticks/alpha-connect/utils"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	tickstore_grpc "gitlab.com/alphaticks/tickstore-grpc"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/models"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/credentials"
)

var done = make(chan os.Signal, 1)

var executorActor *actor.PID
var assetLoader *actor.PID

type GuardActor struct{}

func (state *GuardActor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		context.Watch(executorActor)
		context.Watch(assetLoader)

	case *actor.Terminated:
		done <- os.Signal(syscall.SIGTERM)
	}
}

func main() {

	as := actor.NewActorSystem()
	ctx := actor.NewRootContext(as, nil)

	actorAddress := os.Getenv("ACTOR_ADDRESS")
	actorAdvertisedAddress := os.Getenv("ACTOR_ADVERTISED_ADDRESS")
	if actorAddress != "" {
		address := strings.Split(actorAddress, ":")[0]
		port, err := strconv.ParseInt(strings.Split(actorAddress, ":")[1], 10, 64)
		if err != nil {
			panic(err)
		}
		conf := remote.Configure(address, int(port))
		if actorAdvertisedAddress != "" {
			conf.AdvertisedHost = actorAdvertisedAddress
		}
		//conf = conf.WithServerOptions()
		rem := remote.NewRemote(as, conf)
		rem.Start()
	}

	// Start actors
	exch := []*models.Exchange{
		constants.BINANCE,
		constants.BITFINEX,
		constants.BITSTAMP,
		constants.COINBASEPRO,
		constants.GEMINI,
		constants.KRAKEN,
		constants.CRYPTOFACILITIES,
		constants.OKCOIN,
		constants.FBINANCE,
		constants.HITBTC,
		constants.HUOBI,
		constants.FTX,
		constants.BITMEX,
		constants.BITSTAMP,
		constants.DERIBIT,
		constants.HUOBIP,
		constants.HUOBIF,
		constants.BYBITI,
		constants.BYBITL,
		constants.UPBIT,
		constants.BITHUMB,
		constants.BITHUMBG,
		constants.DYDX,
		constants.OKEXP,
	}
	// EXECUTOR //
	registryAddress := "registry.alphaticks.io:8001"
	if os.Getenv("REGISTRY_ADDRESS") != "" {
		registryAddress = os.Getenv("REGISTRY_ADDRESS")
	}
	conn, err := grpc.Dial(registryAddress, grpc.WithInsecure())
	if err != nil {
		err := fmt.Errorf("error connecting to public registry gRPC endpoint: %v", err)
		panic(err)
	}
	rgstr := registry.NewPublicRegistryClient(conn)
	assetLoader = ctx.Spawn(actor.PropsFromProducer(utils.NewStaticLoaderProducer(rgstr)))
	_, err = ctx.RequestFuture(assetLoader, &utils.Ready{}, 10*time.Second).Result()
	if err != nil {
		panic(err)
	}
	// TODO mongo env
	cfgExch := &extypes.ExecutorConfig{
		Exchanges:  exch,
		DialerPool: xchangerUtils.DefaultDialerPool,
		Strict:     true,
	}
	cfgPrt := &prtypes.ExecutorConfig{
		Registry:  nil,
		Protocols: nil,
	}
	cfgCh := &chtypes.ExecutorConfig{
		Registry: nil,
		Chains:   nil,
	}
	executorActor, _ = ctx.SpawnNamed(actor.PropsFromProducer(executor.NewExecutorProducer(cfgExch, cfgPrt, cfgCh)), "executor")

	// Spawn guard actor
	guardActor, err := ctx.SpawnNamed(
		actor.PropsFromProducer(func() actor.Actor { return &GuardActor{} }),
		"guard_actor")
	if err != nil {
		panic(err)
	}

	var dataServer *grpc.Server
	// Start live store gRPC server
	if address := os.Getenv("DATA_STORE_ADDRESS"); address != "" {
		lis, err := net.Listen("tcp", address)
		if err != nil {
			panic(err)
		}
		dataServer = grpc.NewServer()

		serverAddress := os.Getenv("DATA_SERVER_ADDRESS")
		if serverAddress == "" {
			panic("DATA_SERVER_ADDRESS undefined")
		}

		str, err := data.NewStorageClient("", serverAddress)
		if err != nil {
			panic(err)
		}
		dataER := rpc.NewDataER(ctx, str)

		go func() {
			err := dataServer.Serve(lis)
			if err != nil {
				fmt.Println("ERROR", err)
			}
			done <- os.Signal(syscall.SIGTERM)
		}()

		tickstore_grpc.RegisterStoreServer(dataServer, dataER)
	}

	// If interrupt or terminate signal is received, stop
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	if dataServer != nil {
		// Stop gRPC server
		c := make(chan bool, 1)
		go func() {
			dataServer.GracefulStop()
			c <- true
		}()

		select {
		case <-c:
		case <-time.After(time.Second * 10):
			dataServer.Stop()
		}
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
