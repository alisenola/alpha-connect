// Actor that is connected to tickstore with all the public data
package live

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	_ "gitlab.com/tachikoma.ai/tickobjects/market"
	"gitlab.com/tachikoma.ai/tickstore/actors/messages/remote/store"
	"gitlab.com/tachikoma.ai/tickstore/parsing"
	"reflect"
	"time"
)

type SecurityInfo struct {
	securityID        uint64
	minPriceIncrement *types.DoubleValue
	roundLot          *types.DoubleValue
	tickPrecision     uint64
	lotPrecision      uint64
}

type Feed struct {
	requestID     uint64
	seqNum        uint64
	lastEventTime time.Time
	lastDeltaTime uint64
	security      *SecurityInfo
	groupID       uint64
	tags          map[string]string
}

type LiveStore struct {
	ID         uint64
	executor   *actor.PID
	index      *utils.TagIndex
	securities map[uint64]*models.Security
	logger     *log.Logger
}

func NewLiveStoreProducer(ID uint64) actor.Producer {
	return func() actor.Actor {
		return NewLiveStore(ID)
	}
}

func NewLiveStore(ID uint64) actor.Actor {
	return &LiveStore{
		ID: ID,
	}
}

func (state *LiveStore) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.logger.Error("error initializing", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor started")

	case *actor.Stopping:
		if err := state.Clean(context); err != nil {
			state.logger.Error("error stopping", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor stopping")

	case *actor.Stopped:
		state.logger.Info("actor stopped")

	case *actor.Restarting:
		if err := state.Clean(context); err != nil {
			state.logger.Error("error restarting", log.Error(err))
			// Attention, no panic in restarting or infinite loop
		}
		state.logger.Info("actor restarting")

	case *messages.SecurityList:
		if err := state.OnSecurityList(context); err != nil {
			state.logger.Error("error processing OnSecurityList", log.Error(err))
			panic(err)
		}

	case *store.GetQueryReaderRequest:
		if err := state.GetQueryReaderRequest(context); err != nil {
			state.logger.Error("error processing GetQueryReaderRequest", log.Error(err))
		}
	}
}

func (state *LiveStore) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	state.executor = context.ActorSystem().NewLocalPID("executor")
	// Populate index
	res, err := context.RequestFuture(state.executor, &messages.SecurityListRequest{
		RequestID:  state.ID,
		Subscribe:  true,
		Subscriber: context.Self(),
	}, 10*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error fetching securities: %v", err)
	}
	secRes, ok := res.(*messages.SecurityList)
	if !ok {
		return fmt.Errorf("was expecting SecurityList, got %s", reflect.TypeOf(res))
	}
	state.index = utils.NewTagIndex(secRes.Securities)
	state.securities = make(map[uint64]*models.Security)

	for _, sec := range secRes.Securities {
		state.securities[sec.SecurityID] = sec
	}

	return nil
}

func (state *LiveStore) OnSecurityList(context actor.Context) error {
	msg := context.Message().(*messages.SecurityList)
	state.index = utils.NewTagIndex(msg.Securities)
	state.securities = make(map[uint64]*models.Security)
	for _, sec := range msg.Securities {
		state.securities[sec.SecurityID] = sec
	}
	return nil
}

func (state *LiveStore) GetQueryReaderRequest(context actor.Context) error {
	// Need to hash the query, get firstEventTime and lastEventTime,
	// If missing, fetch from client
	// Launch a query reader actor which will ingest the data from alphaC listeners
	msg := context.Message().(*store.GetQueryReaderRequest)
	selector := msg.Query.Selector
	feeds := make(map[uint64]*Feed)
	sel := parsing.Selector{}
	if err := parsing.SelectorParser.ParseString(selector, &sel); err != nil {
		return fmt.Errorf("error parsing query: %v", err)
	}

	// Now you have everything you need
	if sel.MetaSelector != nil {
		return fmt.Errorf("meta selector not supported")
	}

	functor := &sel.TickSelector.Functor
	for functor.Measurement == nil {
		functor = functor.Application.Value
	}

	queryTags := make(map[string]string)
	for _, t := range sel.TickSelector.Tags {
		queryTags[t.Key] = t.Value
	}
	inputSecurities, err := state.index.Query(queryTags)
	if err != nil {
		return fmt.Errorf("error querying input securities: %v", err)
	}

	// Need to subscribe to every securities, clone the functor for each group
	// for each security, compute which functor it has to go to, based on tags.
	// Place in map[securityID]Functor, use security ID as object ID
	for _, s := range inputSecurities {
		sec := state.securities[s.SecurityID]
		tags := map[string]string{
			"ID":       fmt.Sprintf("%d", sec.SecurityID),
			"type":     sec.SecurityType,
			"base":     sec.Underlying.Symbol,
			"quote":    sec.QuoteCurrency.Symbol,
			"exchange": sec.Exchange.Name,
			"symbol":   sec.Symbol,
		}
		// Extract group
		groupTags := make(map[string]string)
		groupBy := sel.TickSelector.GroupBy
		if len(groupBy) == 0 {
			groupBy = []string{"ID", "type", "base", "quote", "exchange", "symbol"}
		}
		for _, t := range groupBy {
			if tag, ok := tags[t]; ok {
				groupTags[t] = tag
			} else {
				return fmt.Errorf("group by unknown tag %s", t)
			}
		}
		groupID := utils.HashTags(groupTags)
		feed := &Feed{
			requestID: uint64(time.Now().UnixNano()),
			tags:      groupTags,
			groupID:   groupID,
		}
		feed.security = &SecurityInfo{
			securityID:        sec.SecurityID,
			minPriceIncrement: sec.MinPriceIncrement,
			roundLot:          sec.RoundLot,
		}
		feeds[sec.SecurityID] = feed
	}

	prop := actor.PropsFromProducer(NewQueryReaderProducer(feeds, sel.TickSelector.Functor))
	pid := context.Spawn(prop)

	context.Respond(&store.GetQueryReaderResponse{
		RequestID: msg.RequestID,
		Error:     "",
		Reader:    pid,
	})

	return nil
}

func (state *LiveStore) Clean(context actor.Context) error {
	return nil
}
