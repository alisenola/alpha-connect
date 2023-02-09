package modeling

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/utils"
	types "gitlab.com/alphaticks/tickstore-types"
	"gitlab.com/alphaticks/tickstore-types/tickobjects"
	"math"
	"reflect"
	"time"
)

type checkQueries struct{}

type securityInfo struct {
	requestID     uint64
	securityID    uint64
	tickPrecision float64
	lotPrecision  float64
	seqNum        uint64
	lastEventTime time.Time
	lastDeltaTime uint64
}

type Feed struct {
	tradeFunctor     tickobjects.TickFunctor
	orderBookFunctor tickobjects.TickFunctor
	feedID           uint64
}

type Modeler struct {
	executor      *actor.PID
	index         *utils.TagIndex
	store         types.TickstoreClient
	securityInfos map[uint64]*securityInfo
	subscriptions map[uint64]*securityInfo
	feeds         map[uint64]*Feed
	model         Model
	selectors     []string
	logger        *log.Logger
	ID            uint64
	frequency     uint64
	queries       []types.TickstoreQuery
	queryTicker   *time.Ticker
}

func NewModelerProducer(model Model, store types.TickstoreClient, selectors []string) actor.Producer {
	return func() actor.Actor {
		return NewModeler(model, store, selectors)
	}
}

func NewModeler(model Model, store types.TickstoreClient, selectors []string) actor.Actor {
	return &Modeler{
		model:     model,
		store:     store,
		selectors: selectors,
	}
}

func (state *Modeler) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.logger.Error("error initializing actor", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor started")

	case *actor.Stopping:
		state.logger.Info("actor stopping")
		state.Clean(context)

	case *actor.Stopped:
		state.logger.Info("actor stopped")

	case *actor.Restarting:
		state.logger.Info("actor restarting")
		state.Clean(context)

	case *checkQueries:
		if err := state.checkQueries(context); err != nil {
			state.logger.Error("error checking queries", log.Error(err))
			panic(err)
		}
	}
}

func (state *Modeler) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	state.executor = context.ActorSystem().NewLocalPID("executor")
	state.securityInfos = make(map[uint64]*securityInfo)
	state.subscriptions = make(map[uint64]*securityInfo)
	state.feeds = make(map[uint64]*Feed)
	state.frequency = state.model.Frequency()

	queryTicker := time.NewTicker(5 * time.Second)
	state.queryTicker = queryTicker
	go func(pid *actor.PID) {
		for {
			select {
			case _ = <-queryTicker.C:
				context.Send(pid, &checkQueries{})
			case <-time.After(10 * time.Second):
				// timer stopped, we leave
				return
			}
		}
	}(context.Self())

	return state.startQueries(context)
}

func (state *Modeler) startQueries(context actor.Context) error {
	var queries []types.TickstoreQuery
	for _, selector := range state.selectors {
		fmt.Println(selector)
		now := uint64(time.Now().UnixNano() / 1000000)
		qs := types.NewQuerySettings(
			types.WithSelector(selector),
			types.WithFrom(now),
			types.WithTo(math.MaxUint64),
			types.WithStreaming(true),
			types.WithTimeout(100*time.Millisecond),
			types.WithBatchSize(100),
		)
		qs.WithSelector(selector)

		q, err := state.store.NewQuery(qs)
		if err != nil {
			return fmt.Errorf("error querying store: %v", err)
		}

		queries = append(queries, q)
	}

	for _, q := range state.queries {
		_ = q.Close()
	}
	for i, q := range queries {
		go func(i int, q types.TickstoreQuery) {
			for q.Err() == nil {
				for q.Next() {
					tick, obj, groupID := q.Read()
					//fmt.Println(i, state.selectors[i], obj.(tickobjects.Float64Object).Float64())
					if err := state.model.Forward(i, tick, groupID, obj); err != nil {
						_ = q.Close()
					}
				}
			}
		}(i, q)
	}
	state.queries = queries

	return nil
}

func (state *Modeler) restartQuery(idx int) error {
	selector := state.selectors[idx]
	fmt.Println(selector)
	now := uint64(time.Now().UnixNano() / 1000000)
	qs := types.NewQuerySettings(
		types.WithSelector(selector),
		types.WithFrom(now),
		types.WithTo(math.MaxUint64),
		types.WithStreaming(true),
		types.WithTimeout(100*time.Millisecond),
		types.WithBatchSize(100),
	)
	qs.WithSelector(selector)

	q, err := state.store.NewQuery(qs)
	if err != nil {
		return fmt.Errorf("error querying store: %v", err)
	}
	_ = state.queries[idx].Close()
	go func(i int, q types.TickstoreQuery) {
		for q.Err() == nil {
			for q.Next() {
				tick, obj, groupID := q.Read()
				if err := state.model.Forward(i, tick, groupID, obj); err != nil {
					_ = q.Close()
				}
			}
		}
	}(idx, q)
	state.queries[idx] = q

	return nil
}

func (state *Modeler) Clean(context actor.Context) error {
	for _, q := range state.queries {
		_ = q.Close()
	}
	state.queries = nil
	if state.queryTicker != nil {
		state.queryTicker.Stop()
		state.queryTicker = nil
	}
	return nil
}

func (state *Modeler) checkQueries(context actor.Context) error {
	for i, q := range state.queries {
		if q.Err() != nil {
			state.logger.Error("error on query", log.Error(q.Err()), log.String("selector", state.selectors[i]))
			if err := state.restartQuery(i); err != nil {
				return fmt.Errorf("error restarting query: %v", err)
			}
		} else {
			_, obj, _ := q.Read()
			if f, ok := obj.(tickobjects.Float64Object); ok {
				if math.IsNaN(f.Float64()) || math.IsInf(f.Float64(), 0) {
					state.logger.Error("error on query", log.Error(fmt.Errorf("NAN of INF detected")), log.String("selector", state.selectors[i]))
					_ = q.Close()
					if err := state.restartQuery(i); err != nil {
						return fmt.Errorf("error restarting query: %v", err)
					}
				}
			}
		}
	}
	return nil
}
