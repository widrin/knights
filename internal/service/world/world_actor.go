package world

import (
	"github.com/widrin/knights/internal/actor"
)

// WorldActor manages the global game world state
type WorldActor struct {
	state *WorldState
}

// NewWorldActor creates a new world actor
func NewWorldActor() actor.Actor {
	return &WorldActor{
		state: NewWorldState(),
	}
}

// Receive handles incoming messages
func (w *WorldActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		w.onStarted(ctx)
	case *actor.Stopping:
		w.onStopping(ctx)
	case *GetWorldStateRequest:
		w.handleGetWorldState(ctx, msg)
	case *UpdateWorldStateRequest:
		w.handleUpdateWorldState(ctx, msg)
	}
}

func (w *WorldActor) onStarted(ctx actor.Context) {
	// Initialize world
}

func (w *WorldActor) onStopping(ctx actor.Context) {
	// Cleanup world
}

func (w *WorldActor) handleGetWorldState(ctx actor.Context, msg *GetWorldStateRequest) {
	ctx.Respond(&GetWorldStateResponse{
		Success: true,
		State:   w.state,
	})
}

func (w *WorldActor) handleUpdateWorldState(ctx actor.Context, msg *UpdateWorldStateRequest) {
	// TODO: Update world state based on message
	ctx.Respond(&UpdateWorldStateResponse{
		Success: true,
	})
}

// WorldState represents the global game world state
type WorldState struct {
	ServerTime  int64
	OnlineCount int
	// Add more global state fields
}

// NewWorldState creates a new world state
func NewWorldState() *WorldState {
	return &WorldState{
		ServerTime:  0,
		OnlineCount: 0,
	}
}

// Messages

type GetWorldStateRequest struct{}

type GetWorldStateResponse struct {
	Success bool
	State   *WorldState
	Error   string
}

type UpdateWorldStateRequest struct {
	// Add update fields
}

type UpdateWorldStateResponse struct {
	Success bool
	Error   string
}
