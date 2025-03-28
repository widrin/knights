package net

import (
	"knights/logger"
	"sync"
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
	logger.Debug("Network read recorded",
		logger.Uint64("bytes", n),
		logger.Uint64("total", s.BytesRead))
}

func (s *TrafficStats) RecordWrite(n uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.BytesWrite += n
	logger.Debug("Network write recorded",
		logger.Uint64("bytes", n),
		logger.Uint64("total", s.BytesWrite))
}
