package main

import (
	"sync"
	"time"
)

type Stats struct {
	mu sync.Mutex

	Success uint64
	Failed  uint64

	Latencies []time.Duration
	Ops       map[string]uint64
}

func NewStats() *Stats {
	return &Stats{
		Latencies: make([]time.Duration, 0, 100000),
		Ops:       make(map[string]uint64),
	}
}

func (s *Stats) Add(op string, latency time.Duration, success bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Latencies = append(s.Latencies, latency)

	if success {
		s.Success++
		s.Ops[op]++
	} else {
		s.Failed++
	}
}
