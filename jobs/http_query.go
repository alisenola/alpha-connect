package jobs

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"io/ioutil"
	"net/http"
)

// An api query actor execute a query and fits the result back into a given types

// The query actor panic when the request doesn't succeed (timeout etc..)
// The query actor doesn't panic when the request was successful but the server
// returned an error (client or server error)

type PerformHTTPQueryRequest struct {
	Request *http.Request
}

// We allow this query to fail without crashing
// because the failure is outside the system
// We are not responsible for the failure
type PerformQueryResponse struct {
	StatusCode int64
	Response   []byte
}

type HTTPQuery struct {
	client *http.Client
}

func NewHTTPQuery(client *http.Client) actor.Actor {
	return &HTTPQuery{client}
}

func (q *HTTPQuery) Receive(context actor.Context) {
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
		q.Clean(context)

	case *PerformHTTPQueryRequest:
		//Set API credentials
		err := q.PerformHTTPQueryRequest(context)
		if err != nil {
			panic(err)
		}
	}
}

func (q *HTTPQuery) Initialize(context actor.Context) error {
	return nil
}

func (q *HTTPQuery) Clean(context actor.Context) error {
	return nil
}

func (q *HTTPQuery) PerformHTTPQueryRequest(context actor.Context) error {
	msg := context.Message().(*PerformHTTPQueryRequest)
	go func(sender *actor.PID) {
		queryResponse := PerformQueryResponse{}
		resp, err := q.client.Do(msg.Request)
		if err != nil {
			queryResponse.Response = nil
			queryResponse.StatusCode = 500
		} else {
			defer resp.Body.Close()
			queryResponse.StatusCode = int64(resp.StatusCode)
			queryResponse.Response, _ = ioutil.ReadAll(resp.Body)
		}
		context.Send(sender, &queryResponse)
	}(context.Sender())

	return nil
}
