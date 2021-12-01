package tests

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/exchanges"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
)

func StartExecutor(exchange *xchangerModels.Exchange) (*actor.ActorSystem, *actor.PID, func()) {
	exch := []*xchangerModels.Exchange{
		exchange,
	}
	as := actor.NewActorSystem()
	cfg := &exchanges.ExecutorConfig{
		Exchanges: exch,
		Strict:    true,
	}
	executor, _ := as.Root.SpawnNamed(actor.PropsFromProducer(exchanges.NewExecutorProducer(cfg)), "executor")
	return as, executor, func() { _ = as.Root.PoisonFuture(executor).Wait() }
}
