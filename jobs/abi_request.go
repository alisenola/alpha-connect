package jobs

import (
	goContext "context"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// An api query actor execute a query and fits the result back into a given types

// The query actor panic when the request doesn't succeed (timeout etc..)
// The query actor doesn't panic when the request was successful but the server
// returned an error (client or server error)

type PerformLogsQueryRequest struct {
	Query ethereum.FilterQuery
}

// We allow this query to fail without crashing
// because the failure is outside the system
// We are not responsible for the failure
type PerformLogsQueryResponse struct {
	Error error
	Logs  []types.Log
}

type ABIQuery struct {
	client *ethclient.Client
}

func NewABIQuery(client *ethclient.Client) actor.Actor {
	return &ABIQuery{client}
}

func (q *ABIQuery) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		err := q.Initialize(context)
		if err != nil {
			panic(err)
		}

	case *actor.Stopping:
		if err := q.Clean(context); err != nil {
			panic(err)
		}

	case *actor.Restarting:
		if err := q.Clean(context); err != nil {
			// Attention, no panic in restarting or infinite loop
		}

	case *PerformLogsQueryRequest:
		//Set API credentials
		err := q.PerformLogsQueryRequest(context)
		if err != nil {
			panic(err)
		}
	}
}

func (q *ABIQuery) Initialize(context actor.Context) error {
	return nil
}

func (q *ABIQuery) Clean(context actor.Context) error {
	return nil
}

func (q *ABIQuery) PerformLogsQueryRequest(context actor.Context) error {
	msg := context.Message().(*PerformLogsQueryRequest)
	go func(sender *actor.PID) {
		queryResponse := PerformLogsQueryResponse{}
		logs, err := q.client.FilterLogs(goContext.Background(), msg.Query)
		if err != nil {
			queryResponse.Error = err
		} else {
			queryResponse.Logs = logs
		}
		context.Send(sender, &queryResponse)
	}(context.Sender())

	return nil
}
