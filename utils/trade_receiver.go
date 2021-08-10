package utils

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"time"
)

// Wrap actor and sends messages to a channel

type TradeReceiver struct {
	ch   chan *models.TradeCapture
	as   *actor.ActorSystem
	pid  *actor.PID
	seq  uint64
	from uint64
}

func NewTradeReceiver(as *actor.ActorSystem, account *models.Account, instrument *models.Instrument, from uint64, requestID uint64, ch chan *models.TradeCapture) *TradeReceiver {
	r := &TradeReceiver{
		ch:   ch,
		as:   as,
		from: from,
	}

	// First, fetch all the trades
	executor := as.NewLocalPID("executor")
	receiver := as.Root.Spawn(actor.PropsFromFunc(func(c actor.Context) {
		switch msg := c.Message().(type) {
		case *actor.Started:
			req := &messages.AccountDataRequest{
				RequestID:  requestID,
				Subscribe:  true,
				Subscriber: c.Self(),
				Account:    account,
			}
			res, err := c.RequestFuture(executor, req, time.Minute).Result()
			if err != nil {
				panic(fmt.Errorf("error getting market data"))
			}
			r.seq = res.(*messages.AccountDataResponse).SeqNum
			done := false
			for !done {
				res, err := as.Root.RequestFuture(executor, &messages.TradeCaptureReportRequest{
					RequestID: 0,
					Filter: &messages.TradeCaptureReportFilter{
						From:       MilliToTimestamp(r.from),
						Instrument: instrument,
					},
					Account: account,
				}, 20*time.Second).Result()
				if err != nil {
					panic(fmt.Errorf("error getting trade capture report"))
				}
				trades := res.(*messages.TradeCaptureReport)
				if !trades.Success {
					panic(fmt.Errorf("error getting trade capture report: %s", trades.RejectionReason.String()))
				}
				for _, trd := range trades.Trades {
					ch <- trd
					r.from = TimestampToMilli(trd.TransactionTime) + 1
				}
				if len(trades.Trades) == 0 {
					fmt.Println("FINISH", r.from)
					done = true
				}
			}

		case *messages.AccountDataIncrementalRefresh:
			if msg.Report.SeqNum <= r.seq {
				return
			}
			if r.seq+1 != msg.Report.SeqNum {
				panic(fmt.Errorf("out of order sequence"))
			}
			r.seq = msg.Report.SeqNum
			fmt.Println(msg)
		}
	}))

	r.pid = receiver

	return r
}

func (r *TradeReceiver) Close() {
	r.as.Root.Stop(r.pid)
	// TODO chans ?
}
