package rpc

import (
	"context"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/config"
	"gitlab.com/alphaticks/alpha-connect/data"
	"gitlab.com/alphaticks/alpha-connect/executor"
	tickstore_go_client "gitlab.com/alphaticks/tickstore-go-client"
	tsconfig "gitlab.com/alphaticks/tickstore-go-client/config"
	tickstore_grpc "gitlab.com/alphaticks/tickstore-grpc"
	types "gitlab.com/alphaticks/tickstore-types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"net"
	"testing"
)

func TestDataEr(t *testing.T) {
	C, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Start the server
	lis, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	dataServer := grpc.NewServer()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Append license-id and license-key
		md := metadata.New(map[string]string{"license-id": "39KpGZmHKj4=", "license-key": "rxS6ux0OoXybnfY9o2EdoA=="})
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}))
	opts = append(opts, grpc.WithStreamInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (stream grpc.ClientStream, err error) {
		// Append license-id and license-key
		md := metadata.New(map[string]string{"license-id": "39KpGZmHKj4=", "license-key": "rxS6ux0OoXybnfY9o2EdoA=="})
		ctx = metadata.NewOutgoingContext(ctx, md)
		return streamer(ctx, desc, cc, method, opts...)
	}))
	str, err := data.NewStorageClient("/tmp/data", "store.alphaticks.io", opts...)
	if err != nil {
		panic(err)
	}

	as := actor.NewActorSystem()
	ctx := actor.NewRootContext(as, nil)
	executorActor, _ := ctx.SpawnNamed(actor.PropsFromProducer(executor.NewExecutorProducer(C)), "executor")

	lstr, err := data.NewLiveStore(as, executorActor)

	dataER := NewDataER(str, lstr)

	go func() {
		err := dataServer.Serve(lis)
		if err != nil {
			fmt.Println("ERROR", err)
		}
	}()

	tickstore_grpc.RegisterStoreServer(dataServer, dataER)

	// Start the client
	cl, err := tickstore_go_client.NewRemoteClient(tsconfig.StoreClient{Address: "127.0.0.1:8080"}, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}

	qs := types.NewQuerySettings(
		types.WithSelector("SELECT AggOHLCV(ohlcv, \"60000\") WHERE ID=\"9281941173829172773\""),
		types.WithFrom(uint64(1667862714185)),
		types.WithTo(uint64(1667866314185)),
		types.WithStreaming(false),
		types.WithTimeout(100))
	q, err := cl.NewQuery(qs)
	if err != nil {
		t.Fatal(err)
	}
	for q.Next() {
		fmt.Println(q.Read())
	}
	fmt.Println(q.Err())
}
