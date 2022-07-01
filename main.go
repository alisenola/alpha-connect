package main

import (
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/config"
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
	"gitlab.com/alphaticks/alpha-connect/data"
	"gitlab.com/alphaticks/alpha-connect/rpc"
	tickstore_grpc "gitlab.com/alphaticks/tickstore-grpc"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/credentials"
)

var done = make(chan os.Signal, 1)

var executorActor *actor.PID

type GuardActor struct{}

func (state *GuardActor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		context.Watch(executorActor)

	case *actor.Terminated:
		done <- os.Signal(syscall.SIGTERM)
	}
}

func main() {

	as := actor.NewActorSystem()
	ctx := actor.NewRootContext(as, nil)

	C, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	if C.ActorAddress != "" {
		address := strings.Split(C.ActorAddress, ":")[0]
		port, err := strconv.ParseInt(strings.Split(C.ActorAddress, ":")[1], 10, 64)
		if err != nil {
			panic(err)
		}
		conf := remote.Configure(address, int(port))
		if C.ActorAdvertisedAddress != "" {
			conf.AdvertisedHost = C.ActorAdvertisedAddress
		}
		//conf = conf.WithServerOptions()
		rem := remote.NewRemote(as, conf)
		rem.Start()
	}

	// EXECUTOR //
	executorActor, _ = ctx.SpawnNamed(actor.PropsFromProducer(executor.NewExecutorProducer(C)), "executor")

	// Spawn guard actor
	guardActor, err := ctx.SpawnNamed(
		actor.PropsFromProducer(func() actor.Actor { return &GuardActor{} }),
		"guard_actor")
	if err != nil {
		panic(err)
	}

	var dataServer *grpc.Server
	// Start live store gRPC server
	if C.DataStoreAddress != "" {
		lis, err := net.Listen("tcp", C.DataStoreAddress)
		if err != nil {
			panic(err)
		}
		dataServer = grpc.NewServer()

		if C.DataServerAddress == "" {
			panic("DataServerAddress undefined")
		}

		str, err := data.NewStorageClient("", C.DataServerAddress)
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
