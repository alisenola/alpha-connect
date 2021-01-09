package rpc

import (
	"context"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pkg/errors"
	_ "gitlab.com/tachikoma.ai/tickobjects/market"
	"gitlab.com/tachikoma.ai/tickstore-grpc"
	"gitlab.com/tachikoma.ai/tickstore/actors/messages/remote/query"
	"gitlab.com/tachikoma.ai/tickstore/actors/messages/remote/store"
	"io"
	"log"
	"time"
)

type LiveER struct {
	ctx    *actor.RootContext
	store  *actor.PID
	Logger *log.Logger
}

func NewLiveER(ctx *actor.RootContext, store *actor.PID) *LiveER {
	return &LiveER{
		ctx:    ctx,
		store:  store,
		Logger: nil,
	}
}

func (s *LiveER) GetMeasurements(context context.Context, request *tickstore_grpc.GetMeasurementsRequest) (*tickstore_grpc.GetMeasurementsResponse, error) {
	res := &tickstore_grpc.GetMeasurementsResponse{
		Measurements: []*tickstore_grpc.Measurement{{
			Name:   "trade",
			TypeID: "RawTrade",
		}},
	}
	return res, nil
}

func (s *LiveER) RegisterMeasurement(context context.Context, request *tickstore_grpc.RegisterMeasurementRequest) (*tickstore_grpc.RegisterMeasurementResponse, error) {
	return &tickstore_grpc.RegisterMeasurementResponse{}, nil
}

func (s *LiveER) DeleteMeasurement(context context.Context, request *tickstore_grpc.DeleteMeasurementRequest) (*tickstore_grpc.DeleteMeasurementResponse, error) {
	return &tickstore_grpc.DeleteMeasurementResponse{}, nil
}

func (s *LiveER) GetLastEventTime(context context.Context, request *tickstore_grpc.GetLastEventTimeRequest) (*tickstore_grpc.GetLastEventTimeResponse, error) {
	return &tickstore_grpc.GetLastEventTimeResponse{
		Tick: uint64(time.Now().UnixNano() / 1000000),
	}, nil
}

// TODO logging
func (s *LiveER) Write(stream tickstore_grpc.Store_WriteServer) error {
	return fmt.Errorf("live store doesn't support writes")
}

func (s *LiveER) Query(qr *tickstore_grpc.StoreQueryRequest, stream tickstore_grpc.Store_QueryServer) error {

	// TODO more checks, prevent too big batch size, etc..
	if qr.BatchSize == 0 || qr.Timeout == 0 {
		return fmt.Errorf("batch size or time out unspecified")
	}

	done := false
	go func() {
		<-stream.Context().Done()
		done = true
	}()

	res, err := s.ctx.RequestFuture(s.store, &store.GetQueryReaderRequest{
		Query:     qr,
		RequestID: 0,
	}, 30*time.Second).Result()
	if err != nil {
		return fmt.Errorf("unable to get reader: %v", err)
	}
	msg := res.(*store.GetQueryReaderResponse)
	if msg.Error != "" {
		return errors.New(msg.Error)
	}
	queryReader := msg.Reader
	defer func() {
		// stop the query
		s.ctx.Send(
			queryReader,
			&query.CloseQueryReader{})
	}()

	for {
		res, err := s.ctx.RequestFuture(
			queryReader,
			&query.ReadQueryEventBatchRequest{
				RequestID: 0,
				Timeout:   qr.Timeout,
				BatchSize: qr.BatchSize,
			},
			500*time.Second).Result()
		if err != nil {
			return fmt.Errorf("error reading batch from actor: %v", err)
		}
		msg := res.(*query.ReadQueryEventBatchResponse)
		// we can have an error on msg but with still some events in the batch
		// or it can be an empty batch in case of streaming
		if msg.Batch != nil && !done {
			if err := stream.Send(msg.Batch); err != nil {
				return fmt.Errorf("error sending batch: %v", err)
			}
		}

		if msg.Error != "" {
			if msg.Error == io.EOF.Error() {
				return nil
			} else {
				return errors.New(msg.Error)
			}
		}

		// if context is done, stop
		if done {
			break
		}
	}

	return nil
}
