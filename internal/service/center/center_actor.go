package center

import (
	"sync"

	"github.com/widrin/knights/internal/actor"
)

// CenterActor 中心服务，负责服务器管理、负载均衡等
type CenterActor struct {
	servers     map[string]*ServerInfo
	serverMutex sync.RWMutex
}

// ServerInfo 服务器信息
type ServerInfo struct {
	ServerID   string
	ServerType string // login, game, gateway
	Address    string
	Port       int
	PlayerCount int
	MaxPlayers int
	Status     ServerStatus
	PID        *actor.PID
}

// ServerStatus 服务器状态
type ServerStatus int

const (
	ServerStatusOffline ServerStatus = iota
	ServerStatusOnline
	ServerStatusMaintenance
	ServerStatusFull
)

// NewCenterActor 创建中心服务Actor
func NewCenterActor() actor.Actor {
	return &CenterActor{
		servers: make(map[string]*ServerInfo),
	}
}

// Receive 处理消息
func (c *CenterActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		c.onStarted(ctx)

	case *RegisterServerRequest:
		c.handleRegisterServer(ctx, msg)

	case *UnregisterServerRequest:
		c.handleUnregisterServer(ctx, msg)

	case *GetServerListRequest:
		c.handleGetServerList(ctx, msg)

	case *GetBestServerRequest:
		c.handleGetBestServer(ctx, msg)

	case *UpdateServerStatusRequest:
		c.handleUpdateServerStatus(ctx, msg)

	case *HeartbeatMessage:
		c.handleHeartbeat(ctx, msg)
	}
}

func (c *CenterActor) onStarted(ctx actor.Context) {
	// 初始化中心服务
}

func (c *CenterActor) handleRegisterServer(ctx actor.Context, msg *RegisterServerRequest) {
	c.serverMutex.Lock()
	defer c.serverMutex.Unlock()

	info := &ServerInfo{
		ServerID:   msg.ServerID,
		ServerType: msg.ServerType,
		Address:    msg.Address,
		Port:       msg.Port,
		MaxPlayers: msg.MaxPlayers,
		Status:     ServerStatusOnline,
		PID:        msg.PID,
	}

	c.servers[msg.ServerID] = info

	ctx.Respond(&RegisterServerResponse{
		Success: true,
	})
}

func (c *CenterActor) handleUnregisterServer(ctx actor.Context, msg *UnregisterServerRequest) {
	c.serverMutex.Lock()
	defer c.serverMutex.Unlock()

	delete(c.servers, msg.ServerID)

	ctx.Respond(&UnregisterServerResponse{
		Success: true,
	})
}

func (c *CenterActor) handleGetServerList(ctx actor.Context, msg *GetServerListRequest) {
	c.serverMutex.RLock()
	defer c.serverMutex.RUnlock()

	servers := make([]*ServerInfo, 0, len(c.servers))
	for _, server := range c.servers {
		if msg.ServerType == "" || server.ServerType == msg.ServerType {
			servers = append(servers, server)
		}
	}

	ctx.Respond(&GetServerListResponse{
		Servers: servers,
	})
}

func (c *CenterActor) handleGetBestServer(ctx actor.Context, msg *GetBestServerRequest) {
	c.serverMutex.RLock()
	defer c.serverMutex.RUnlock()

	var bestServer *ServerInfo
	minLoad := 100.0

	for _, server := range c.servers {
		if server.ServerType != msg.ServerType {
			continue
		}
		if server.Status != ServerStatusOnline {
			continue
		}

		// 计算负载率
		load := float64(server.PlayerCount) / float64(server.MaxPlayers)
		if load < minLoad {
			minLoad = load
			bestServer = server
		}
	}

	if bestServer != nil {
		ctx.Respond(&GetBestServerResponse{
			Success: true,
			Server:  bestServer,
		})
	} else {
		ctx.Respond(&GetBestServerResponse{
			Success: false,
			Error:   "no available server",
		})
	}
}

func (c *CenterActor) handleUpdateServerStatus(ctx actor.Context, msg *UpdateServerStatusRequest) {
	c.serverMutex.Lock()
	defer c.serverMutex.Unlock()

	if server, ok := c.servers[msg.ServerID]; ok {
		server.PlayerCount = msg.PlayerCount
		server.Status = msg.Status

		ctx.Respond(&UpdateServerStatusResponse{
			Success: true,
		})
	} else {
		ctx.Respond(&UpdateServerStatusResponse{
			Success: false,
			Error:   "server not found",
		})
	}
}

func (c *CenterActor) handleHeartbeat(ctx actor.Context, msg *HeartbeatMessage) {
	// 更新服务器心跳时间
	// TODO: 实现心跳超时检测
}

// Messages

type RegisterServerRequest struct {
	ServerID   string
	ServerType string
	Address    string
	Port       int
	MaxPlayers int
	PID        *actor.PID
}

type RegisterServerResponse struct {
	Success bool
	Error   string
}

type UnregisterServerRequest struct {
	ServerID string
}

type UnregisterServerResponse struct {
	Success bool
	Error   string
}

type GetServerListRequest struct {
	ServerType string // 空字符串表示获取所有类型
}

type GetServerListResponse struct {
	Servers []*ServerInfo
}

type GetBestServerRequest struct {
	ServerType string
}

type GetBestServerResponse struct {
	Success bool
	Server  *ServerInfo
	Error   string
}

type UpdateServerStatusRequest struct {
	ServerID    string
	PlayerCount int
	Status      ServerStatus
}

type UpdateServerStatusResponse struct {
	Success bool
	Error   string
}

type HeartbeatMessage struct {
	ServerID string
}
