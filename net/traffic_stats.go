package net

import (
	"sync"

	"github.com/widrin/knights/logger"
)

type TrafficStats struct {
	mu         sync.RWMutex
	BytesRead  uint64
	BytesWrite uint64
}

func NewTrafficStats() *TrafficStats {
	return &TrafficStats{}
}

func (s *TrafficStats) RecordRead(n uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.BytesRead += n
	logger.Debug("Network read recorded bytes: %d, total: %d", n, s.BytesRead)
}

func (s *TrafficStats) RecordWrite(n uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.BytesWrite += n
	logger.Debug("Network write recorded bytes: %d, total: %d", n, s.BytesWrite)
}
