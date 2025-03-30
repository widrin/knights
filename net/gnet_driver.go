package net

import (
	"fmt"
	"sync"
	"time"

	"github.com/widrin/knights/logger"

	"github.com/panjf2000/gnet/v2"
)

type GNetDriver struct {
	gnet.BuiltinEventEngine

	config      *NetworkConfig
	server      gnet.Engine
	connections sync.Map
	stats       *TrafficStats
	endian      *EndianHandler
	pool        *ConnectionPool
}

func (g *GNetDriver) Initialize(config *NetworkConfig) error {
	if err := config.Validate(); err != nil {
		return err
	}

	g.config = config
	g.stats = NewTrafficStats()
	g.endian = NewEndianHandler(config.Endianness)
	g.pool = NewConnectionPool(config.MaxConn)

	logger.Info("Initializing network driver protocol: %d, host: %s, port: %d",
		config.Protocol, config.Host, config.Port)
	return nil
}

func (g *GNetDriver) Start() error {
	options := []gnet.Option{
		gnet.WithMulticore(true),
		gnet.WithTCPKeepAlive(time.Minute),
		gnet.WithSocketSendBuffer(g.config.BufferSize),
		gnet.WithSocketRecvBuffer(g.config.BufferSize),
	}
	protoAddr := fmt.Sprintf("%s://%s:%d", g.config.Protocol, g.config.Host, g.config.Port)

	logger.Info("Starting network server addr: %s", protoAddr)

	return gnet.Run(g, protoAddr, options...)
}

// 实现Shutdown和GetConnections方法...
