package utils

import (
	"runtime"
	"sync"
	"time"
)

type item struct {
	value      interface{}
	lastAccess int64
}

type TTLMap struct {
	m       map[string]*item
	l       sync.Mutex
	ttl     int64
	stopped bool
}

func NewTTLMap(TTL int) (m *TTLMap) {
	m = &TTLMap{
		m:       make(map[string]*item),
		l:       sync.Mutex{},
		ttl:     int64(TTL),
		stopped: true}
	return
}

func (m *TTLMap) Start() {
	m.stopped = false
	go func() {
		for now := range time.NewTicker(time.Second).C {
			m.l.Lock()
			if m.stopped {
				return
			}
			for k, v := range m.m {
				if now.Unix()-v.lastAccess > int64(m.ttl) {
					delete(m.m, k)
				}
			}
			m.l.Unlock()
		}
	}()
	runtime.SetFinalizer(m, m.Stop)
}

func (m *TTLMap) Stop() {
	m.l.Lock()
	m.stopped = true
	runtime.SetFinalizer(m, nil)
	m.l.Unlock()
}

func (m *TTLMap) Len() int {
	return len(m.m)
}

func (m *TTLMap) Put(k string, v interface{}) {
	m.l.Lock()
	it, ok := m.m[k]
	if !ok {
		it = &item{value: v}
		m.m[k] = it
	}
	it.value = v
	it.lastAccess = time.Now().Unix()
	m.l.Unlock()
}

func (m *TTLMap) Get(k string) (v interface{}, ok bool) {
	m.l.Lock()
	if it, ok := m.m[k]; ok {
		v = it.value
	}
	m.l.Unlock()
	return
}
