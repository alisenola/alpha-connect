package jobs

import (
	goContext "context"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/alphaticks/xchanger/exchanges/uniswap/v3"
)

// An api query actor execute a query and fits the result back into a given types

// The query actor panic when the request doesn't succeed (timeout etc..)
// The query actor doesn't panic when the request was successful but the server
// returned an error (client or server error)

type PerformABIRequest struct {
	Address common.Address
	Start   uint64
	End     uint64
}

// We allow this query to fail without crashing
// because the failure is outside the system
// We are not responsible for the failure
type PerformABIResponse struct {
	Error    error
	Response *ABIResponse
}

type ABIResponse struct {
	Initialize *uniswap.UniswapInitializeIterator
	Mints      *uniswap.UniswapMintIterator
	Burns      *uniswap.UniswapBurnIterator
	Swaps      *uniswap.UniswapSwapIterator
	Collects   *uniswap.UniswapCollectIterator
	Flashes    *uniswap.UniswapFlashIterator
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

	case *PerformABIRequest:
		//Set API credentials
		err := q.PerformABIRequest(context)
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

func (q *ABIQuery) PerformABIRequest(context actor.Context) error {
	msg := context.Message().(*PerformABIRequest)
	go func(sender *actor.PID) {
		queryResponse := PerformABIResponse{
			Response: &ABIResponse{},
		}
		instance, err := uniswap.NewUniswap(msg.Address, q.client)
		if err != nil {
			queryResponse.Error = err
		} else {
			iit, err := instance.FilterInitialize(&bind.FilterOpts{
				Start:   msg.Start,
				End:     &msg.End,
				Context: goContext.Background(),
			})
			if err == nil {
				queryResponse.Response.Initialize = iit
			}
			mit, err := instance.FilterMint(&bind.FilterOpts{
				Start:   msg.Start,
				End:     &msg.End,
				Context: goContext.Background(),
			}, nil, nil, nil)
			if err == nil {
				queryResponse.Response.Mints = mit
			}
			bit, err := instance.FilterBurn(&bind.FilterOpts{
				Start:   msg.Start,
				End:     &msg.End,
				Context: goContext.Background(),
			}, nil, nil, nil)
			if err == nil {
				queryResponse.Response.Burns = bit
			}
			cit, err := instance.FilterCollect(&bind.FilterOpts{
				Start:   msg.Start,
				End:     &msg.End,
				Context: goContext.Background(),
			}, nil, nil, nil)
			if err == nil {
				queryResponse.Response.Collects = cit
			}
			fit, err := instance.FilterFlash(&bind.FilterOpts{
				Start:   msg.Start,
				End:     &msg.End,
				Context: goContext.Background(),
			}, nil, nil)
			if err == nil {
				queryResponse.Response.Flashes = fit
			}

		}
		context.Send(sender, &queryResponse)
	}(context.Sender())

	return nil
}
