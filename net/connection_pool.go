package net

import (
	"knights/logger"
	"sync"
)

type ConnectionPool struct {
	mu            sync.RWMutex
	conns         map[uint64]interface{}
	maxSize       int
	autoReconnect bool
}

func NewConnectionPool(maxSize int) *ConnectionPool {
	return &ConnectionPool{
		conns:   make(map[uint64]interface{}),
		maxSize: maxSize,
	}
}

func (p *ConnectionPool) Add(cid uint64, conn interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.conns) >= p.maxSize {
		logger.Warn("Connection pool full")
		return
	}

	p.conns[cid] = conn
	logger.Debug("Connection added",
		logger.Uint64("cid", cid),
		logger.Int("total", len(p.conns)))
}

func (p *ConnectionPool) Remove(cid uint64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.conns, cid)
	logger.Debug("Connection removed",
		logger.Uint64("cid", cid),
		logger.Int("remaining", len(p.conns)))

	if p.autoReconnect {
		go p.reconnect(cid)
	}
}

func (p *ConnectionPool) reconnect(cid uint64) {
	// 实现自动重连逻辑
}
