package main

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	fix50mdr "github.com/quickfixgo/fix50/marketdatarequest"
	fix50sl "github.com/quickfixgo/fix50/securitylist"
	fix50slr "github.com/quickfixgo/fix50/securitylistrequest"
	"gitlab.com/alphaticks/alphac/exchanges"
)

type FIXRouter struct {
	exchanges                  []string
	executors                  map[string]*actor.PID
	instrumentListeners        map[string]*actor.PID
	runningSecurityListRequest map[string][]*fix50sl.SecurityList
	logger                     *log.Logger
}

func NewFIXRouter(exchanges []string) actor.Actor {
	return &FIXRouter{
		exchanges: exchanges,
		logger:    nil,
	}
}

func (state *FIXRouter) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		err := state.Initialize(context)
		if err != nil {
			panic(err)
		}

	case *fix50mdr.MarketDataRequest:
		if err := state.OnFIX50MarketDataRequest(context); err != nil {
			state.logger.Error("error processing FIX50MarketDataRequest", log.Error(err))
			panic(err)
		}

	case *fix50slr.SecurityListRequest:
		if err := state.OnFIX50SecurityListRequest(context); err != nil {
			state.logger.Error("error processing FIX50SecurityListRequest", log.Error(err))
			panic(err)
		}

	case *fix50sl.SecurityList:
		if err := state.OnFIX50SecurityList(context); err != nil {
			state.logger.Error("error processing FIX50SecurityList", log.Error(err))
			panic(err)
		}
	}
}

func (state *FIXRouter) Initialize(context actor.Context) error {
	// Instantiate exchanges
	state.executors = make(map[string]*actor.PID)
	state.instrumentListeners = make(map[string]*actor.PID)
	state.runningSecurityListRequest = make(map[string][]*fix50sl.SecurityList)

	for _, exchg := range state.exchanges {
		producer := exchanges.NewExchangeExecutorProducer(exchg)
		if producer == nil {
			return fmt.Errorf("unknown exchange %s", exchg)
		}
		props := actor.PropsFromProducer(producer)
		pid := context.Spawn(props)
		state.executors[exchg] = pid
	}
	return nil
}

func (state *FIXRouter) OnFIX50MarketDataRequest(context actor.Context) error {
	//msg := context.Message().(*fix50mdr.MarketDataRequest)
	// Spawn instrument listener if not exist already
	return nil
}

func (state *FIXRouter) OnFIX50SecurityListRequest(context actor.Context) error {
	req := context.Message().(*fix50slr.SecurityListRequest)
	reqID, err := req.GetSecurityReqID()
	if err != nil {
		return err
	}
	for _, executor := range state.executors {
		context.Request(executor, req)
	}
	state.runningSecurityListRequest[reqID] = nil
	return nil
}

func (state *FIXRouter) OnFIX50SecurityList(context actor.Context) error {
	res := context.Message().(*fix50sl.SecurityList)
	reqID, _ := res.GetSecurityReqID()
	if _, ok := state.runningSecurityListRequest[reqID]; ok {
		state.runningSecurityListRequest[reqID] = append(state.runningSecurityListRequest[reqID], res)
		if len(state.runningSecurityListRequest[reqID]) == len(state.exchanges) {
			totNoRelatedSym := 0
			for _, sl := range state.runningSecurityListRequest[reqID] {
				syms, _ := sl.GetNoRelatedSym()
				totNoRelatedSym += syms.Len()
			}
			for i, sl := range state.runningSecurityListRequest[reqID] {
				sl.SetTotNoRelatedSym(totNoRelatedSym)
				sl.SetLastFragment(i == len(state.runningSecurityListRequest[reqID]))
				context.Send(context.Parent(), sl)
			}
			delete(state.runningSecurityListRequest, reqID)
		}
	}

	return nil
}
