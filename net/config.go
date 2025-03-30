package net

import (
	"errors"
)

type ProtocolType = string

const (
	TCP       ProtocolType = "tcp"
	UDP       ProtocolType = "udp"
	WebSocket ProtocolType = "ws"
	KCP       ProtocolType = "kcp"
)

type NetworkConfig struct {
	Protocol   ProtocolType
	Host       string
	Port       int
	MaxConn    int
	BufferSize int
	Endianness string
}

func (c *NetworkConfig) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return errors.New("invalid port number")
	}
	if c.MaxConn < 0 {
		return errors.New("max connections cannot be negative")
	}
	return nil
}

type NetworkDriver interface {
	Initialize(config *NetworkConfig) error
	Start() error
	Shutdown() error
	GetConnections() int
}
