package player

import (
	"fmt"
	"sync"

	"github.com/widrin/knights/internal/actor"
)

// PlayerManager manages all player actors
type PlayerManager struct {
	players sync.Map // map[string]*actor.PID
}

// NewPlayerManager creates a new player manager
func NewPlayerManager() actor.Actor {
	return &PlayerManager{}
}

// Receive handles incoming messages
func (pm *PlayerManager) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		pm.onStarted(ctx)
	case *CreatePlayerRequest:
		pm.handleCreatePlayer(ctx, msg)
	case *RemovePlayerRequest:
		pm.handleRemovePlayer(ctx, msg)
	case *GetPlayerRequest:
		pm.handleGetPlayer(ctx, msg)
	}
}

func (pm *PlayerManager) onStarted(ctx actor.Context) {
	// Initialize player manager
}

func (pm *PlayerManager) handleCreatePlayer(ctx actor.Context, msg *CreatePlayerRequest) {
	// Check if player already exists
	if _, exists := pm.players.Load(msg.PlayerID); exists {
		ctx.Respond(&CreatePlayerResponse{
			Success: false,
			Error:   fmt.Sprintf("player %s already exists", msg.PlayerID),
		})
		return
	}

	// Create player actor
	props := actor.NewProps(func() actor.Actor {
		return NewPlayerActor(msg.PlayerID)
	})
	pid := ctx.Spawn(props)

	// Store player PID
	pm.players.Store(msg.PlayerID, pid)

	ctx.Respond(&CreatePlayerResponse{
		Success:  true,
		PlayerID: msg.PlayerID,
		PID:      pid,
	})
}

func (pm *PlayerManager) handleRemovePlayer(ctx actor.Context, msg *RemovePlayerRequest) {
	if pid, exists := pm.players.LoadAndDelete(msg.PlayerID); exists {
		ctx.Stop(pid.(*actor.PID))
		ctx.Respond(&RemovePlayerResponse{Success: true})
	} else {
		ctx.Respond(&RemovePlayerResponse{
			Success: false,
			Error:   fmt.Sprintf("player %s not found", msg.PlayerID),
		})
	}
}

func (pm *PlayerManager) handleGetPlayer(ctx actor.Context, msg *GetPlayerRequest) {
	if pid, exists := pm.players.Load(msg.PlayerID); exists {
		ctx.Respond(&GetPlayerResponse{
			Success: true,
			PID:     pid.(*actor.PID),
		})
	} else {
		ctx.Respond(&GetPlayerResponse{
			Success: false,
			Error:   fmt.Sprintf("player %s not found", msg.PlayerID),
		})
	}
}
