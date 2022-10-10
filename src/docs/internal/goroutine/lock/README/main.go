package main

import "sync"

type Stat struct {
	counters map[string]int64
	mutex    sync.RWMutex
}

func (s *Stat) getCounter(name string) int64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.counters[name]
}

func (s *Stat) setCounter(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.counters[name]++
}
