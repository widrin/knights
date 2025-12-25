package room

import (
	"github.com/widrin/knights/internal/actor"
)

// RoomActor represents a game room/battle instance
type RoomActor struct {
	roomID  string
	state   *RoomState
	players map[string]*actor.PID
}

// NewRoomActor creates a new room actor
func NewRoomActor(roomID string, maxPlayers int) actor.Actor {
	return &RoomActor{
		roomID:  roomID,
		state:   NewRoomState(roomID, maxPlayers),
		players: make(map[string]*actor.PID),
	}
}

// Receive handles incoming messages
func (r *RoomActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		r.onStarted(ctx)
	case *actor.Stopping:
		r.onStopping(ctx)
	case *JoinRoomRequest:
		r.handleJoinRoom(ctx, msg)
	case *LeaveRoomRequest:
		r.handleLeaveRoom(ctx, msg)
	case *BroadcastMessage:
		r.handleBroadcast(ctx, msg)
	}
}

func (r *RoomActor) onStarted(ctx actor.Context) {
	// Initialize room
}

func (r *RoomActor) onStopping(ctx actor.Context) {
	// Cleanup room
}

func (r *RoomActor) handleJoinRoom(ctx actor.Context, msg *JoinRoomRequest) {
	if len(r.players) >= r.state.MaxPlayers {
		ctx.Respond(&JoinRoomResponse{
			Success: false,
			Error:   "room is full",
		})
		return
	}

	r.players[msg.PlayerID] = msg.PlayerPID
	r.state.PlayerCount = len(r.players)

	ctx.Respond(&JoinRoomResponse{
		Success: true,
		RoomID:  r.roomID,
	})

	// Notify other players
	r.broadcastToPlayers(ctx, &PlayerJoinedEvent{
		PlayerID: msg.PlayerID,
	})
}

func (r *RoomActor) handleLeaveRoom(ctx actor.Context, msg *LeaveRoomRequest) {
	delete(r.players, msg.PlayerID)
	r.state.PlayerCount = len(r.players)

	ctx.Respond(&LeaveRoomResponse{Success: true})

	// Notify other players
	r.broadcastToPlayers(ctx, &PlayerLeftEvent{
		PlayerID: msg.PlayerID,
	})

	// Close room if empty
	if len(r.players) == 0 {
		ctx.StopSelf()
	}
}

func (r *RoomActor) handleBroadcast(ctx actor.Context, msg *BroadcastMessage) {
	r.broadcastToPlayers(ctx, msg.Message)
}

func (r *RoomActor) broadcastToPlayers(ctx actor.Context, message interface{}) {
	for _, playerPID := range r.players {
		ctx.Send(playerPID, message)
	}
}
