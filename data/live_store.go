package data

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	tickstore_go_client "gitlab.com/tachikoma.ai/tickstore-go-client"
	"gitlab.com/tachikoma.ai/tickstore-go-client/query"
	"gitlab.com/tachikoma.ai/tickstore/parsing"
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
	requestID uint64
	security  *SecurityInfo
	receiver  *utils.MDReceiver
	groupID   uint64
	tags      map[string]string
}

type LiveStore struct {
	executor   *actor.PID
	as         *actor.ActorSystem
	index      *utils.TagIndex
	securities map[uint64]*models.Security
	queries    []*LiveQuery
}

func NewLiveStore(as *actor.ActorSystem) (*LiveStore, error) {
	executor := as.NewLocalPID("executor")
	res, err := as.Root.RequestFuture(executor, &messages.SecurityListRequest{}, 20*time.Second).Result()
	if err != nil {
		return nil, fmt.Errorf("error fetching securities: %v", err)
	}
	sl := res.(*messages.SecurityList)
	// Build index
	index := utils.NewTagIndex(sl.Securities)
	securities := make(map[uint64]*models.Security)
	for _, sec := range sl.Securities {
		securities[sec.SecurityID] = sec
	}

	return &LiveStore{
		executor:   executor,
		as:         as,
		index:      index,
		securities: securities,
	}, nil
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

func (lt *LiveStore) NewTickWriter(measurement string, tags map[string]string, flushTime time.Duration) (tickstore_go_client.TickstoreWriter, error) {
	return nil, fmt.Errorf("no tick writer supported on live store")
}

func (lt *LiveStore) NewQuery(qs *query.Settings) (tickstore_go_client.TickstoreQuery, error) {
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
	inputSecurities, err := lt.index.Query(queryTags)
	if err != nil {
		return nil, fmt.Errorf("error querying input securities: %v", err)
	}

	// Need to subscribe to every securities, clone the functor for each group
	// for each security, compute which functor it has to go to, based on tags.
	// Place in map[securityID]Functor, use security ID as object ID
	feeds := make(map[uint64]*Feed)
	for _, s := range inputSecurities {
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

	q, err := NewLiveQuery(lt.as, sel, feeds)
	lt.queries = append(lt.queries, q)
	return q, nil
}
