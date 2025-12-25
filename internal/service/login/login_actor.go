package login

import (
	"time"

	"github.com/widrin/knights/internal/actor"
)

// LoginActor 处理用户登录认证
type LoginActor struct {
	sessions map[string]*LoginSession
}

// LoginSession 登录会话信息
type LoginSession struct {
	UserID    string
	Token     string
	LoginTime time.Time
	IPAddress string
}

// NewLoginActor 创建登录服务Actor
func NewLoginActor() actor.Actor {
	return &LoginActor{
		sessions: make(map[string]*LoginSession),
	}
}

// Receive 处理消息
func (l *LoginActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		l.onStarted(ctx)

	case *LoginRequest:
		l.handleLogin(ctx, msg)

	case *LogoutRequest:
		l.handleLogout(ctx, msg)

	case *ValidateTokenRequest:
		l.handleValidateToken(ctx, msg)

	case *GetOnlineCountRequest:
		l.handleGetOnlineCount(ctx, msg)
	}
}

func (l *LoginActor) onStarted(ctx actor.Context) {
	// 初始化登录服务
}

func (l *LoginActor) handleLogin(ctx actor.Context, msg *LoginRequest) {
	// TODO: 验证用户名密码
	// TODO: 检查是否已经在线
	// TODO: 生成token

	token := l.generateToken(msg.Username)
	session := &LoginSession{
		UserID:    msg.Username,
		Token:     token,
		LoginTime: time.Now(),
		IPAddress: msg.IPAddress,
	}

	l.sessions[msg.Username] = session

	ctx.Respond(&LoginResponse{
		Success: true,
		Token:   token,
		UserID:  msg.Username,
	})
}

func (l *LoginActor) handleLogout(ctx actor.Context, msg *LogoutRequest) {
	delete(l.sessions, msg.UserID)

	ctx.Respond(&LogoutResponse{
		Success: true,
	})
}

func (l *LoginActor) handleValidateToken(ctx actor.Context, msg *ValidateTokenRequest) {
	if session, ok := l.sessions[msg.UserID]; ok {
		if session.Token == msg.Token {
			ctx.Respond(&ValidateTokenResponse{
				Valid:  true,
				UserID: msg.UserID,
			})
			return
		}
	}

	ctx.Respond(&ValidateTokenResponse{
		Valid: false,
		Error: "invalid token",
	})
}

func (l *LoginActor) handleGetOnlineCount(ctx actor.Context, msg *GetOnlineCountRequest) {
	ctx.Respond(&GetOnlineCountResponse{
		Count: len(l.sessions),
	})
}

func (l *LoginActor) generateToken(userID string) string {
	// TODO: 实现真实的token生成逻辑
	return "token_" + userID + "_" + time.Now().Format("20060102150405")
}

// Messages

type LoginRequest struct {
	Username  string
	Password  string
	IPAddress string
}

type LoginResponse struct {
	Success bool
	Token   string
	UserID  string
	Error   string
}

type LogoutRequest struct {
	UserID string
}

type LogoutResponse struct {
	Success bool
	Error   string
}

type ValidateTokenRequest struct {
	UserID string
	Token  string
}

type ValidateTokenResponse struct {
	Valid  bool
	UserID string
	Error  string
}

type GetOnlineCountRequest struct{}

type GetOnlineCountResponse struct {
	Count int
}
