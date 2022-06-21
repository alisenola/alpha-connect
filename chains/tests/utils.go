package tests

import (
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/chains"
	"gitlab.com/alphaticks/alpha-connect/chains/types"
)

func StartExecutor() (*actor.ActorSystem, *actor.PID, func(), error) {
	config := &actor.Config{DeveloperSupervisionLogging: true, DeadLetterRequestLogging: true}
	as := actor.NewActorSystemWithConfig(config)
	cfg := types.ExecutorConfig{}
	ex, err := as.Root.SpawnNamed(actor.PropsFromProducer(chains.NewExecutorProducer(&cfg)), "executor")
	if err != nil {
		return nil, nil, nil, err
	}

	return as, ex, func() { as.Root.PoisonFuture(ex).Wait() }, nil
}
