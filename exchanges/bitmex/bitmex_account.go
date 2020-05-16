package bitmex

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges/bitmex"
	"gitlab.com/alphaticks/xchanger/exchanges/coinbasepro"
	"reflect"
)

type Account struct {
	credentials     exchanges.APICredentials
	ws              *bitmex.Websocket
	wsChan          chan *coinbasepro.WebsocketMessage
	executorManager *actor.PID
	logger          *log.Logger
}

func NewAccountProducer(credentials exchanges.APICredentials) actor.Producer {
	return func() actor.Actor {
		return NewAccount(credentials)
	}
}

func NewAccount(credentials exchanges.APICredentials) actor.Actor {
	return &Account{
		credentials:     credentials,
		ws:              nil,
		wsChan:          nil,
		executorManager: nil,
		logger:          nil,
	}
}

func (state *Account) Receive(context actor.Context) {
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

	case *readSocket:
		if err := state.readSocket(context); err != nil {
			state.logger.Error("error processing readSocket", log.Error(err))
			panic(err)
		}
	}
}

func (state *Account) Initialize(context actor.Context) error {
	// When initialize is done, the account must be aware of all the settings / assets / portofilio
	// so as to be able to answer to FIX messages

	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	if err := state.subscribeAccount(context); err != nil {
		return fmt.Errorf("error subscribing to order book: %v", err)
	}
	// Then fetch positions
	// Then fetch orders
	// Then reconcile with execution feed

	context.Send(context.Self(), &readSocket{})

	return nil
}

// TODO
func (state *Account) Clean(context actor.Context) error {
	if state.ws != nil {
		if err := state.ws.Disconnect(); err != nil {
			state.logger.Info("error disconnecting socket", log.Error(err))
		}
	}

	return nil
}

func (state *Account) readSocket(context actor.Context) error {
	select {
	case msg := <-state.wsChan:
		switch msg.Message.(type) {

		case error:
			return fmt.Errorf("socket error: %v", msg)

		case bitmex.WSExecutionData:
			execData := msg.Message.(bitmex.WSExecutionData)
			fmt.Println(execData)
		}
	}

	return nil
}

func (state *Account) subscribeAccount(context actor.Context) error {
	if state.ws != nil {
		_ = state.ws.Disconnect()
	}

	ws := bitmex.NewWebsocket()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to bitmex websocket: %v", err)
	}

	if err := ws.Auth(&state.credentials); err != nil {
		return fmt.Errorf("error sending auth request: %v", err)
	}

	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	receivedMessage, ok := ws.Msg.Message.(bitmex.WSResponse)
	if !ok {
		errorMessage, ok := ws.Msg.Message.(bitmex.WSErrorResponse)
		if ok {
			return fmt.Errorf("error auth: %s", errorMessage.Error)
		}
		return fmt.Errorf("error casting message to WSResponse")
	}

	if !receivedMessage.Success {
		return fmt.Errorf("auth unsuccessful")
	}

	if err := ws.Subscribe(bitmex.WSExecutionStreamName); err != nil {
		return fmt.Errorf("error sending subscription request: %v", err)
	}
	if !ws.ReadMessage() {
		return fmt.Errorf("error reading message: %v", ws.Err)
	}
	subResponse, ok := ws.Msg.Message.(bitmex.WSSubscribeResponse)
	if !ok {
		errorMessage, ok := ws.Msg.Message.(bitmex.WSErrorResponse)
		if ok {
			return fmt.Errorf("error auth: %s", errorMessage.Error)
		}
		return fmt.Errorf("error casting message to WSSubscribeResponse")
	}
	fmt.Println(subResponse.Success)

	return nil
}
