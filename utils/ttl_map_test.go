package utils

import (
	"testing"
	"time"
)

func TestTTLMap(t *testing.T) {
	m := NewTTLMap(1)
	m.Start()
	m.Put("kek", 1)
	time.Sleep(5 * time.Second)
	m.Stop()
	time.Sleep(5 * time.Second)
}
