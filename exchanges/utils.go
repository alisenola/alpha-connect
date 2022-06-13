package exchanges

import (
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/xchanger/models"
)

func GetExecutor(as *actor.ActorSystem, exchange *models.Exchange) *actor.PID {
	return actor.NewPID(as.Address(), "executor/exchanges/"+exchange.Name+"_executor")
}
