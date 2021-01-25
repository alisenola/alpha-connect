package rpc

/*
import (
	"context"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"log"
	"reflect"
	"time"
)

type ExecutorEndpointReader struct {
	Ctx      *actor.RootContext
	Executor *actor.PID
	Logger   *log.Logger
}

func (s *ExecutorEndpointReader) MarketData(context context.Context, request *messages.MarketDataRequest) (*messages.MarketDataIncrementalRefresh, error) {
	res, err := s.Ctx.RequestFuture(s.Executor, request, 10*time.Second).Result()
	if err != nil {
		return nil, err
	}
	md, ok := res.(*messages.MarketDataResponse)
	if !ok {
		return nil, fmt.Errorf("was expecting *messages.MarketDataResponse got %s", reflect.TypeOf(md).String())
	}

	receiver := s.Ctx.Spawn(actor.ReceiverFunc())
	s.Ctx.Send

	return nil, nil
}

*/
