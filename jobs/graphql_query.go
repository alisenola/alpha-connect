package jobs

import (
	goContext "context"
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/go-graphql-client"
)

// An api query actor execute a query and fits the result back into a given types

// The query actor panic when the request doesn't succeed (timeout etc..)
// The query actor doesn't panic when the request was successful but the server
// returned an error (client or server error)

type PerformGraphQueryRequest struct {
	Query     interface{}
	Variables map[string]interface{}
	Options   []graphql.Option
}

// We allow this query to fail without crashing
// because the failure is outside the system
// We are not responsible for the failure
type PerformGraphQueryResponse struct {
	Error error
}

type GraphQuery struct {
	client *graphql.Client
}

func NewGraphQuery(client *graphql.Client) actor.Actor {
	return &GraphQuery{client}
}

func (q *GraphQuery) Receive(context actor.Context) {
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

	case *PerformGraphQueryRequest:
		//Set API credentials
		err := q.PerformGraphQueryRequest(context)
		if err != nil {
			panic(err)
		}
	}
}

func (q *GraphQuery) Initialize(context actor.Context) error {
	return nil
}

func (q *GraphQuery) Clean(context actor.Context) error {
	return nil
}

func (q *GraphQuery) PerformGraphQueryRequest(context actor.Context) error {
	msg := context.Message().(*PerformGraphQueryRequest)
	go func(sender *actor.PID) {
		queryResponse := PerformGraphQueryResponse{}
		queryResponse.Error = q.client.Query(goContext.Background(), msg.Query, msg.Variables, msg.Options...)
		context.Send(sender, &queryResponse)
	}(context.Sender())

	return nil
}
