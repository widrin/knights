package net

import (
	"context"
	"knights/logger"
	"sync"
	"time"

	"github.com/panjf2000/gnet/v2"
)

type GNetDriver struct {
	config      *NetworkConfig
	server      gnet.Server
	connections sync.Map
	stats       *TrafficStats
	endian      EndianHandler
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

	logger.Info("Initializing network driver",
		logger.String("protocol", config.Protocol.String()),
		logger.Int("port", config.Port))
	return nil
}

func (g *GNetDriver) Start() error {
	options := []gnet.Option{
		gnet.WithMulticore(true),
		gnet.WithTCPKeepAlive(time.Minute),
		gnet.WithSocketSendBuffer(g.config.BufferSize),
		gnet.WithSocketRecvBuffer(g.config.BufferSize),
	}

	server, err := gnet.NewServer(&gnetHandler{driver: g}, options...)
	if err != nil {
		return err
	}

	g.server = server
	logger.Info("Starting network server",
		logger.String("host", g.config.Host),
		logger.Int("port", g.config.Port))
	return server.Run(context.Background())
}

// 实现Shutdown和GetConnections方法...
