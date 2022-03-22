package rpc

import (
	"context"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/data"
	_ "gitlab.com/alphaticks/tickfunctors/market"
	"gitlab.com/alphaticks/tickstore-grpc"
	"gitlab.com/alphaticks/tickstore/query"
	"gitlab.com/alphaticks/tickstore/readers"
	"io"
	"log"
	"time"
	"unsafe"
)

type DataER struct {
	ctx    *actor.RootContext
	store  *data.StorageClient
	Logger *log.Logger
}

func NewDataER(ctx *actor.RootContext, store *data.StorageClient) *DataER {
	return &DataER{
		ctx:    ctx,
		store:  store,
		Logger: nil,
	}
}

func (s *DataER) GetMeasurements(context context.Context, request *tickstore_grpc.GetMeasurementsRequest) (*tickstore_grpc.GetMeasurementsResponse, error) {
	res := &tickstore_grpc.GetMeasurementsResponse{
		Measurements: []*tickstore_grpc.Measurement{{
			Name:   "trade",
			TypeID: "RawTrade",
		}, {
			Name:   "ohlcv",
			TypeID: "OHLCV",
		},
			{
				Name:   "obliquidity",
				TypeID: "OBLiquidity",
			},
		}}
	return res, nil
}

func (s *DataER) RegisterMeasurement(context context.Context, request *tickstore_grpc.RegisterMeasurementRequest) (*tickstore_grpc.RegisterMeasurementResponse, error) {
	return &tickstore_grpc.RegisterMeasurementResponse{}, nil
}

func (s *DataER) DeleteMeasurement(context context.Context, request *tickstore_grpc.DeleteMeasurementRequest) (*tickstore_grpc.DeleteMeasurementResponse, error) {
	return &tickstore_grpc.DeleteMeasurementResponse{}, nil
}

func (s *DataER) GetLastEventTime(context context.Context, request *tickstore_grpc.GetLastEventTimeRequest) (*tickstore_grpc.GetLastEventTimeResponse, error) {
	return &tickstore_grpc.GetLastEventTimeResponse{}, nil
}

// TODO logging
func (s *DataER) Write(stream tickstore_grpc.Store_WriteServer) error {
	return fmt.Errorf("data store doesn't support writes")
}

func (s *DataER) Query(qr *tickstore_grpc.StoreQueryRequest, stream tickstore_grpc.Store_QueryServer) error {

	// TODO more checks, prevent too big batch size, etc..
	if qr.BatchSize == 0 || qr.Timeout == 0 {
		return fmt.Errorf("batch size or time out unspecified")
	}

	done := false
	go func() {
		<-stream.Context().Done()
		done = true
	}()

	var freq int64 = 0
	if sampler := qr.GetTickSampler(); sampler != nil {
		freq = int64(sampler.Interval)
	}

	str, _, err := s.store.GetStore(freq)
	if err != nil {
		return fmt.Errorf("error getting store: %v", err)
	}

	qs := query.NewQuerySettings(
		query.WithStreaming(qr.Streaming),
		query.WithFrom(qr.From),
		query.WithTo(qr.To))

	sampler := qr.GetSampler()
	if sampler != nil {
		switch sampler := sampler.(type) {
		case *tickstore_grpc.StoreQueryRequest_TickSampler:
			qs.WithSampler(readers.NewTickSampler(sampler.TickSampler.Interval))
		}
	}
	if err := qs.WithSelectorString(qr.Selector); err != nil {
		return fmt.Errorf("error compiling query: %v", err)
	}

	q, err := str.Query(qs)
	if err != nil {
		return fmt.Errorf("error querying store: %v", err)
	}
	defer func() {
		fmt.Println("CLOSING QUERY")
		_ = q.Close()
	}()

	objects := make(map[uint64]bool)
	for {
		batch := &tickstore_grpc.StoreQueryBatch{
			Ticks: nil,
		}

		events := make([]*tickstore_grpc.StoreQueryTick, 0, qr.BatchSize)

		endTime := time.Now().Add(time.Duration(qr.Timeout))

		// We are asked to read a batch with a timeout.
		// that means we need to return something after the time has passed
		q.SetNextDeadline(endTime)
		for uint64(len(events)) < qr.BatchSize && time.Now().Before(endTime) && q.Next() {
			var event *tickstore_grpc.StoreQueryTick
			tick, deltas, groupID := q.ReadDeltas()
			// The first event has to be a snapshot. So send snapshot if not seen objectID yet
			if _, ok := objects[groupID]; !ok {
				objects[groupID] = true
				event = &tickstore_grpc.StoreQueryTick{
					Tick:    tick,
					GroupId: groupID,
					Event: &tickstore_grpc.StoreQueryTick_Snapshot{
						Snapshot: &tickstore_grpc.StoreQuerySnapshot{
							Tags:     q.Tags(),
							Snapshot: q.TickObject().ToSnapshot(),
						},
					},
				}
			} else {
				// TODO BUG HERERE !!
				typ := q.DeltaType()
				buff := make([]byte, int(typ.Size())*deltas.Len)
				for i := 0; i < len(buff); i++ {
					buff[i] = *(*byte)(unsafe.Pointer(uintptr(deltas.Pointer) + uintptr(i)))
				}
				event = &tickstore_grpc.StoreQueryTick{
					Tick:    tick,
					GroupId: groupID,
					Event: &tickstore_grpc.StoreQueryTick_Deltas{
						Deltas: buff,
					},
				}
			}

			events = append(events, event)

			// Reset the deadline
			q.SetNextDeadline(endTime)
		}

		if len(events) > 0 {
			batch.Ticks = events
		} else {
			batch = nil
		}

		// we can have an error on msg but with still some events in the batch
		// or it can be an empty batch in case of streaming
		if batch != nil && !done {
			if err := stream.Send(batch); err != nil {
				return fmt.Errorf("error sending batch: %v", err)
			}
		}

		if q.Err() != nil {
			if q.Err() == io.EOF {
				return nil
			} else {
				return q.Err()
			}
		}

		// if context is done, stop
		if done {
			break
		}
	}

	return nil
}
