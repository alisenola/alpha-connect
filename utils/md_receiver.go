package utils

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"time"
)

// Wrap actor and sends messages to a channel

type MDReceiver struct {
	ch       chan interface{}
	as       *actor.ActorSystem
	executor *actor.PID
	pid      *actor.PID
	seq      uint64
}

func NewMDReceiver(as *actor.ActorSystem, executor *actor.PID, instrument *models.Instrument, aggregation models.OrderBookAggregation, requestID uint64, ch chan interface{}) *MDReceiver {
	r := &MDReceiver{
		ch: ch,
		as: as,
	}

	receiver := as.Root.Spawn(actor.PropsFromFunc(func(c actor.Context) {
		switch msg := c.Message().(type) {
		case *actor.Started:
			req := &messages.MarketDataRequest{
				RequestID:   requestID,
				Subscribe:   true,
				Subscriber:  c.Self(),
				Instrument:  instrument,
				Aggregation: aggregation,
			}
			res, err := c.RequestFuture(executor, req, time.Minute).Result()
			if err != nil {
				panic(fmt.Errorf("error getting market data"))
			}
			r.seq = res.(*messages.MarketDataResponse).SeqNum
			ch <- res

		case *messages.MarketDataIncrementalRefresh:
			if msg.SeqNum <= r.seq {
				return
			}
			if r.seq+1 != msg.SeqNum {
				panic(fmt.Errorf("out of order sequence"))
			}
			r.seq = msg.SeqNum
			ch <- msg
		}
	}))

	r.pid = receiver

	return r
}

func (r *MDReceiver) Close() {
	r.as.Root.Stop(r.pid)
}
