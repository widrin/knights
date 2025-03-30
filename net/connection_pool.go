package net

import (
	"sync"

	"github.com/widrin/knights/logger"
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
	logger.Debug("Connection added cid: %d, total: %d", cid, len(p.conns))
}

func (p *ConnectionPool) Remove(cid uint64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.conns, cid)
	logger.Debug("Connection removed cid: %d, remaining: %d", cid, len(p.conns))

	if p.autoReconnect {
		go p.reconnect(cid)
	}
}

func (p *ConnectionPool) reconnect(cid uint64) {
	// 实现自动重连逻辑
}
