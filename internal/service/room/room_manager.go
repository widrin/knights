package room

import (
	"fmt"
	"sync"

	"github.com/widrin/knights/internal/actor"
)

// RoomManager manages all game rooms
type RoomManager struct {
	rooms sync.Map // map[string]*actor.PID
}

// NewRoomManager creates a new room manager
func NewRoomManager() actor.Actor {
	return &RoomManager{}
}

// Receive handles incoming messages
func (rm *RoomManager) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		rm.onStarted(ctx)
	case *CreateRoomRequest:
		rm.handleCreateRoom(ctx, msg)
	case *RemoveRoomRequest:
		rm.handleRemoveRoom(ctx, msg)
	case *GetRoomRequest:
		rm.handleGetRoom(ctx, msg)
	}
}

func (rm *RoomManager) onStarted(ctx actor.Context) {
	// Initialize room manager
}

func (rm *RoomManager) handleCreateRoom(ctx actor.Context, msg *CreateRoomRequest) {
	if _, exists := rm.rooms.Load(msg.RoomID); exists {
		ctx.Respond(&CreateRoomResponse{
			Success: false,
			Error:   fmt.Sprintf("room %s already exists", msg.RoomID),
		})
		return
	}

	props := actor.NewProps(func() actor.Actor {
		return NewRoomActor(msg.RoomID, msg.MaxPlayers)
	})
	pid := ctx.Spawn(props)

	rm.rooms.Store(msg.RoomID, pid)

	ctx.Respond(&CreateRoomResponse{
		Success: true,
		RoomID:  msg.RoomID,
		PID:     pid,
	})
}

func (rm *RoomManager) handleRemoveRoom(ctx actor.Context, msg *RemoveRoomRequest) {
	if pid, exists := rm.rooms.LoadAndDelete(msg.RoomID); exists {
		ctx.Stop(pid.(*actor.PID))
		ctx.Respond(&RemoveRoomResponse{Success: true})
	} else {
		ctx.Respond(&RemoveRoomResponse{
			Success: false,
			Error:   fmt.Sprintf("room %s not found", msg.RoomID),
		})
	}
}

func (rm *RoomManager) handleGetRoom(ctx actor.Context, msg *GetRoomRequest) {
	if pid, exists := rm.rooms.Load(msg.RoomID); exists {
		ctx.Respond(&GetRoomResponse{
			Success: true,
			PID:     pid.(*actor.PID),
		})
	} else {
		ctx.Respond(&GetRoomResponse{
			Success: false,
			Error:   fmt.Sprintf("room %s not found", msg.RoomID),
		})
	}
}
