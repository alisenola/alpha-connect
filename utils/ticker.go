package utils

import (
	"time"
)

func StartTicker(f func(), duration time.Duration) chan bool {
	done := make(chan bool, 1)
	go func() {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				f()
			case <-done:
				return
			}
		}
	}()
	return done
}
