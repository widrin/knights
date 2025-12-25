package player

import (
	"github.com/widrin/knights/internal/actor"
)

// PlayerActor represents a connected player
type PlayerActor struct {
	playerID   string
	sessionPID *actor.PID
	state      *PlayerState
}

// NewPlayerActor creates a new player actor
func NewPlayerActor(playerID string) actor.Actor {
	return &PlayerActor{
		playerID: playerID,
		state:    NewPlayerState(playerID),
	}
}

// Receive handles incoming messages
func (p *PlayerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		p.onStarted(ctx)
	case *actor.Stopping:
		p.onStopping(ctx)
	case *LoginRequest:
		p.handleLogin(ctx, msg)
	case *LogoutRequest:
		p.handleLogout(ctx, msg)
	case *MoveRequest:
		p.handleMove(ctx, msg)
	default:
		// Unknown message
	}
}

func (p *PlayerActor) onStarted(ctx actor.Context) {
	// Initialize player actor
}

func (p *PlayerActor) onStopping(ctx actor.Context) {
	// Cleanup player actor
}

func (p *PlayerActor) handleLogin(ctx actor.Context, msg *LoginRequest) {
	// TODO: Implement login logic
	ctx.Respond(&LoginResponse{
		Success: true,
		PlayerID: p.playerID,
	})
}

func (p *PlayerActor) handleLogout(ctx actor.Context, msg *LogoutRequest) {
	// TODO: Implement logout logic
	ctx.Respond(&LogoutResponse{
		Success: true,
	})
}

func (p *PlayerActor) handleMove(ctx actor.Context, msg *MoveRequest) {
	// TODO: Implement movement logic
	p.state.Position = msg.Position
	ctx.Respond(&MoveResponse{
		Success: true,
		Position: p.state.Position,
	})
}
