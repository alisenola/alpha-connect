package data

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	pbtypes "github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	types "gitlab.com/alphaticks/tickstore-types"
	"gitlab.com/alphaticks/tickstore/parsing"
	"sync"
	"time"
)

type SecurityInfo struct {
	securityID        uint64
	minPriceIncrement *pbtypes.DoubleValue
	roundLot          *pbtypes.DoubleValue
	tickPrecision     uint64
	lotPrecision      uint64
}

type Feed struct {
	requestID uint64
	security  *SecurityInfo
	receiver  *utils.MDReceiver
	groupID   uint64
	tags      map[string]string
}

type LiveStore struct {
	sync.RWMutex
	as         *actor.ActorSystem
	executor   *actor.PID
	sub        *actor.PID
	index      *utils.TagIndex
	securities map[uint64]*models.Security
	queries    []*LiveQuery
}

func NewLiveStore(as *actor.ActorSystem, executor *actor.PID) (*LiveStore, error) {
	lt := &LiveStore{
		as:       as,
		executor: executor,
	}
	res, err := as.Root.RequestFuture(executor, &messages.SecurityListRequest{}, 10*time.Second).Result()
	if err != nil {
		return nil, err
	}
	sec := res.(*messages.SecurityList)
	if !sec.Success {
		return nil, fmt.Errorf("error fetching securities: %s", sec.RejectionReason.String())
	}
	lt.index = utils.NewTagIndex(sec.Securities)
	lt.securities = make(map[uint64]*models.Security)
	for _, sec := range sec.Securities {
		lt.securities[sec.SecurityID] = sec
	}

	props := actor.PropsFromFunc(func(c actor.Context) {
		switch res := c.Message().(type) {
		case *actor.Started:
			c.Request(executor, &messages.SecurityListRequest{
				Subscribe:  true,
				Subscriber: c.Self(),
			})
		case *messages.SecurityList:
			index := utils.NewTagIndex(res.Securities)
			securities := make(map[uint64]*models.Security)
			for _, sec := range res.Securities {
				securities[sec.SecurityID] = sec
			}
			lt.Lock()
			lt.index = index
			lt.securities = securities
			lt.Unlock()
		}
	})

	// Build index
	lt.sub = as.Root.Spawn(props)

	return lt, nil
}

func (lt *LiveStore) RegisterMeasurement(measurement string, typeID string) error {
	return nil
}

func (lt *LiveStore) DeleteMeasurement(measurement string, tags map[string]string, from, to uint64) error {
	return nil
}

func (lt *LiveStore) GetLastEventTime(measurement string, tags map[string]string) (uint64, error) {
	return uint64(time.Now().UnixNano()) / 1000000, nil
}

func (lt *LiveStore) NewTickWriter(measurement string, tags map[string]string, flushTime time.Duration) (types.TickstoreWriter, error) {
	return nil, fmt.Errorf("no tick writer supported on live store")
}

func (lt *LiveStore) NewQuery(qs *types.QuerySettings) (types.TickstoreQuery, error) {
	sel := parsing.Selector{}
	if err := parsing.SelectorParser.ParseString(qs.Selector, &sel); err != nil {
		return nil, fmt.Errorf("error parsing selector string: %v", err)
	}

	if sel.MetaSelector != nil {
		return nil, fmt.Errorf("meta selector not supported")
	}

	functor := &sel.TickSelector.Functor
	for functor.Measurement == nil {
		functor = functor.Application.Value
	}

	queryTags := make(map[string]string)
	for _, t := range sel.TickSelector.Tags {
		queryTags[t.Key] = t.Value
	}

	lt.RLock()
	defer lt.RUnlock()
	if lt.index == nil {
		return NewLiveQuery(lt.as, lt.executor, sel, nil)
	}
	inputSecurities, err := lt.index.Query(queryTags)
	if err != nil {
		return nil, fmt.Errorf("error querying input securities: %v", err)
	}
	// Need to subscribe to every securities, clone the functor for each group
	// for each security, compute which functor it has to go to, based on tags.
	// Place in map[securityID]Functor, use security ID as object ID
	feeds := make(map[uint64]*Feed)
	for _, s := range inputSecurities {
		if s.Status != models.Trading {
			continue
		}
		sec := lt.securities[s.SecurityID]
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

		if groupBy == nil {
			for k, _ := range tags {
				groupTags[k] = tags[k]
			}
		} else {
			for _, k := range groupBy.Tags {
				groupTags[k] = tags[k]
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

	q, err := NewLiveQuery(lt.as, lt.executor, sel, feeds)
	lt.queries = append(lt.queries, q)
	return q, nil
}

func (lt *LiveStore) Close() {
	if lt.sub != nil {
		lt.as.Root.Stop(lt.sub)
		lt.sub = nil
	}
}
