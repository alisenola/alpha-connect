package utils

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"math/rand"
	"time"
)

func NewExponentialBackoffStrategy(backoffWindow time.Duration, initialBackoff time.Duration, noise time.Duration) actor.SupervisorStrategy {
	return &exponentialBackoffStrategy{
		backoffWindow:  backoffWindow,
		initialBackoff: initialBackoff,
		noise:          noise,
	}
}

type exponentialBackoffStrategy struct {
	backoffWindow  time.Duration
	initialBackoff time.Duration
	noise          time.Duration
}

func (strategy *exponentialBackoffStrategy) HandleFailure(as *actor.ActorSystem, supervisor actor.Supervisor, child *actor.PID, rs *actor.RestartStatistics, reason interface{}, message interface{}) {
	strategy.setFailureCount(rs)

	backoff := rs.FailureCount() * int(strategy.initialBackoff.Nanoseconds())
	noise := rand.Intn(int(strategy.noise.Nanoseconds()))
	dur := time.Duration(backoff + noise)
	time.AfterFunc(dur, func() {
		supervisor.RestartChildren(child)
	})
}

func (strategy *exponentialBackoffStrategy) setFailureCount(rs *actor.RestartStatistics) {
	if rs.NumberOfFailures(strategy.backoffWindow) == 0 {
		rs.Reset()
	}

	rs.Fail()
}
