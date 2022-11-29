package exchanges

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/config"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	tickstore_go_client "gitlab.com/alphaticks/tickstore-go-client"
	tgc_config "gitlab.com/alphaticks/tickstore-go-client/config"
	"reflect"
	"time"
)

// Subscribe to an account and write it in the monitor store

type AccountMonitor struct {
	config.Account
	tgc_config.StoreClient
	monitorStore *tickstore_go_client.RemoteClient
	executor     *actor.PID
	logger       *log.Logger
}

func NewAccountMonitorProducer(config config.Account, client tgc_config.StoreClient) actor.Producer {
	return func() actor.Actor {
		return NewAccountMonitor(config, client)
	}
}

func NewAccountMonitor(config config.Account, client tgc_config.StoreClient) actor.Actor {
	return &AccountMonitor{
		Account:     config,
		StoreClient: client,
	}
}

func (state *AccountMonitor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.logger.Error("error initializing", log.Error(err))
			panic(err)
		}

		state.logger.Info("actor started")
	}
}

func (state *AccountMonitor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	fmt.Println("ACCOUNT MONITOR STARTING")

	state.executor = GetExchangesExecutor(context.ActorSystem())

	res, err := context.RequestFuture(state.executor, &messages.AccountDataRequest{
		RequestID:  0,
		Subscribe:  true,
		Subscriber: context.Self(),
		Account: &models.Account{
			Name: state.Name,
		},
	}, 20*time.Second).Result()

	if err != nil {
		return err
	}
	fmt.Println(res)

	return nil
}
