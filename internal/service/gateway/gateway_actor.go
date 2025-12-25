package gateway

import (
	"sync"

	"github.com/widrin/knights/internal/actor"
	"github.com/widrin/knights/internal/network"
)

// GatewayActor 网关服务，负责消息转发和路由
type GatewayActor struct {
	sessions      map[string]*network.Session
	userToSession map[string]string // userID -> sessionID
	routes        map[string]*actor.PID
	mu            sync.RWMutex
}

// NewGatewayActor 创建网关服务Actor
func NewGatewayActor() actor.Actor {
	return &GatewayActor{
		sessions:      make(map[string]*network.Session),
		userToSession: make(map[string]string),
		routes:        make(map[string]*actor.PID),
	}
}

// Receive 处理消息
func (g *GatewayActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		g.onStarted(ctx)

	case *BindSessionRequest:
		g.handleBindSession(ctx, msg)

	case *UnbindSessionRequest:
		g.handleUnbindSession(ctx, msg)

	case *ForwardToClientRequest:
		g.handleForwardToClient(ctx, msg)

	case *ForwardToServerRequest:
		g.handleForwardToServer(ctx, msg)

	case *BroadcastRequest:
		g.handleBroadcast(ctx, msg)

	case *RegisterRouteRequest:
		g.handleRegisterRoute(ctx, msg)

	case *GetOnlineUsersRequest:
		g.handleGetOnlineUsers(ctx, msg)
	}
}

func (g *GatewayActor) onStarted(ctx actor.Context) {
	// 初始化网关服务
}

func (g *GatewayActor) handleBindSession(ctx actor.Context, msg *BindSessionRequest) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.sessions[msg.SessionID] = msg.Session
	g.userToSession[msg.UserID] = msg.SessionID

	ctx.Respond(&BindSessionResponse{
		Success: true,
	})
}

func (g *GatewayActor) handleUnbindSession(ctx actor.Context, msg *UnbindSessionRequest) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if sessionID, ok := g.userToSession[msg.UserID]; ok {
		delete(g.sessions, sessionID)
		delete(g.userToSession, msg.UserID)
	}

	ctx.Respond(&UnbindSessionResponse{
		Success: true,
	})
}

func (g *GatewayActor) handleForwardToClient(ctx actor.Context, msg *ForwardToClientRequest) {
	g.mu.RLock()
	sessionID, ok := g.userToSession[msg.UserID]
	g.mu.RUnlock()

	if !ok {
		ctx.Respond(&ForwardToClientResponse{
			Success: false,
			Error:   "user not online",
		})
		return
	}

	g.mu.RLock()
	session, ok := g.sessions[sessionID]
	g.mu.RUnlock()

	if !ok {
		ctx.Respond(&ForwardToClientResponse{
			Success: false,
			Error:   "session not found",
		})
		return
	}

	if err := session.Send(msg.Message); err != nil {
		ctx.Respond(&ForwardToClientResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.Respond(&ForwardToClientResponse{
		Success: true,
	})
}

func (g *GatewayActor) handleForwardToServer(ctx actor.Context, msg *ForwardToServerRequest) {
	g.mu.RLock()
	targetPID, ok := g.routes[msg.Route]
	g.mu.RUnlock()

	if !ok {
		ctx.Respond(&ForwardToServerResponse{
			Success: false,
			Error:   "route not found",
		})
		return
	}

	ctx.Send(targetPID, msg.Message)

	ctx.Respond(&ForwardToServerResponse{
		Success: true,
	})
}

func (g *GatewayActor) handleBroadcast(ctx actor.Context, msg *BroadcastRequest) {
	g.mu.RLock()
	sessions := make([]*network.Session, 0, len(g.sessions))
	for _, session := range g.sessions {
		sessions = append(sessions, session)
	}
	g.mu.RUnlock()

	successCount := 0
	for _, session := range sessions {
		if err := session.Send(msg.Message); err == nil {
			successCount++
		}
	}

	ctx.Respond(&BroadcastResponse{
		Success:      true,
		SendCount:    successCount,
		OnlineCount:  len(sessions),
	})
}

func (g *GatewayActor) handleRegisterRoute(ctx actor.Context, msg *RegisterRouteRequest) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.routes[msg.Route] = msg.TargetPID

	ctx.Respond(&RegisterRouteResponse{
		Success: true,
	})
}

func (g *GatewayActor) handleGetOnlineUsers(ctx actor.Context, msg *GetOnlineUsersRequest) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	users := make([]string, 0, len(g.userToSession))
	for userID := range g.userToSession {
		users = append(users, userID)
	}

	ctx.Respond(&GetOnlineUsersResponse{
		Users: users,
		Count: len(users),
	})
}

// Messages

type BindSessionRequest struct {
	UserID    string
	SessionID string
	Session   *network.Session
}

type BindSessionResponse struct {
	Success bool
	Error   string
}

type UnbindSessionRequest struct {
	UserID string
}

type UnbindSessionResponse struct {
	Success bool
	Error   string
}

type ForwardToClientRequest struct {
	UserID  string
	Message interface{}
}

type ForwardToClientResponse struct {
	Success bool
	Error   string
}

type ForwardToServerRequest struct {
	Route   string
	Message interface{}
}

type ForwardToServerResponse struct {
	Success bool
	Error   string
}

type BroadcastRequest struct {
	Message interface{}
}

type BroadcastResponse struct {
	Success     bool
	SendCount   int
	OnlineCount int
	Error       string
}

type RegisterRouteRequest struct {
	Route     string
	TargetPID *actor.PID
}

type RegisterRouteResponse struct {
	Success bool
	Error   string
}

type GetOnlineUsersRequest struct{}

type GetOnlineUsersResponse struct {
	Users []string
	Count int
}
