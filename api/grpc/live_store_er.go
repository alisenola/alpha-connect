package grpc

import (
	"context"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pkg/errors"
	_ "gitlab.com/tachikoma.ai/tickobjects/market"
	"gitlab.com/tachikoma.ai/tickstore-grpc"
	"gitlab.com/tachikoma.ai/tickstore/actors/messages"
	"gitlab.com/tachikoma.ai/tickstore/db"
	"io"
	"log"
	"time"
)

type LiveStoreER struct {
	StoreActor *actor.PID
	Logger     *log.Logger
}

func (s *LiveStoreER) GetMeasurements(context context.Context, request *tickstore_grpc.GetMeasurementsRequest) (*tickstore_grpc.GetMeasurementsResponse, error) {
	res := &tickstore_grpc.GetMeasurementsResponse{
		Measurements: nil,
	}
	return res, nil
}

func (s *LiveStoreER) RegisterMeasurement(context context.Context, request *tickstore_grpc.RegisterMeasurementRequest) (*tickstore_grpc.RegisterMeasurementResponse, error) {
	if err := db.RegisterMeasurement(request.Measurement.Name, request.Measurement.TypeID); err != nil {
		return nil, err
	}
	return &tickstore_grpc.RegisterMeasurementResponse{}, nil
}

func (s *LiveStoreER) DeleteMeasurement(context context.Context, request *tickstore_grpc.DeleteMeasurementRequest) (*tickstore_grpc.DeleteMeasurementResponse, error) {
	return &tickstore_grpc.DeleteMeasurementResponse{}, nil
}

func (s *LiveStoreER) GetLastEventTime(context context.Context, request *tickstore_grpc.GetLastEventTimeRequest) (*tickstore_grpc.GetLastEventTimeResponse, error) {
	return &tickstore_grpc.GetLastEventTimeResponse{
		Tick: uint64(time.Now().UnixNano() / 1000000),
	}, nil
}

// TODO logging
func (s *LiveStoreER) Write(stream tickstore_grpc.Remoting_WriteServer) error {
	return fmt.Errorf("live store doesn't support writes")
}

func (s *LiveStoreER) Query(query *tickstore_grpc.QueryRequest, stream tickstore_grpc.Remoting_QueryServer) error {

	// TODO more checks, prevent too big batch size, etc..
	if query.BatchSize == 0 || query.Timeout == 0 {
		return fmt.Errorf("batch size or time out unspecified")
	}

	done := false
	go func() {
		<-stream.Context().Done()
		done = true
	}()

	res, err := actor.EmptyRootContext.RequestFuture(s.StoreActor, &messages.GetQueryReaderRequest{
		Query:     query,
		RequestID: 0,
	}, 30*time.Second).Result()
	if err != nil {
		return fmt.Errorf("unable to get reader: %v", err)
	}
	msg := res.(*messages.GetQueryReaderResponse)
	if msg.Error != nil {
		return msg.Error
	}
	queryReader := msg.Reader
	defer func() {
		// stop the query
		actor.EmptyRootContext.Send(
			queryReader,
			&messages.CloseQueryReader{})
	}()

	for {
		res, err := actor.EmptyRootContext.RequestFuture(
			queryReader,
			&messages.ReadQueryEventBatchRequest{
				RequestID: 0,
				Timeout:   query.Timeout,
				BatchSize: query.BatchSize,
			},
			500*time.Second).Result()
		if err != nil {
			return fmt.Errorf("error reading batch from actor: %v", err)
		}
		msg := res.(*messages.ReadQueryEventBatchResponse)
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
