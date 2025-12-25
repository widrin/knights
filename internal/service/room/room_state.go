package room

import (
	"time"

	"github.com/widrin/knights/internal/actor"
)

// RoomState holds the room's state
type RoomState struct {
	RoomID      string
	Name        string
	MaxPlayers  int
	PlayerCount int
	Status      RoomStatus
	CreatedAt   time.Time
}

// RoomStatus represents the current status of a room
type RoomStatus int

const (
	RoomStatusWaiting RoomStatus = iota
	RoomStatusPlaying
	RoomStatusFinished
)

// NewRoomState creates a new room state
func NewRoomState(roomID string, maxPlayers int) *RoomState {
	return &RoomState{
		RoomID:      roomID,
		MaxPlayers:  maxPlayers,
		PlayerCount: 0,
		Status:      RoomStatusWaiting,
		CreatedAt:   time.Now(),
	}
}

// Messages for room actor

type JoinRoomRequest struct {
	PlayerID  string
	PlayerPID *actor.PID
}

type JoinRoomResponse struct {
	Success bool
	RoomID  string
	Error   string
}

type LeaveRoomRequest struct {
	PlayerID string
}

type LeaveRoomResponse struct {
	Success bool
	Error   string
}

type BroadcastMessage struct {
	Message interface{}
}

type PlayerJoinedEvent struct {
	PlayerID string
}

type PlayerLeftEvent struct {
	PlayerID string
}

// Messages for room manager

type CreateRoomRequest struct {
	RoomID     string
	MaxPlayers int
}

type CreateRoomResponse struct {
	Success bool
	RoomID  string
	PID     *actor.PID
	Error   string
}

type RemoveRoomRequest struct {
	RoomID string
}

type RemoveRoomResponse struct {
	Success bool
	Error   string
}

type GetRoomRequest struct {
	RoomID string
}

type GetRoomResponse struct {
	Success bool
	PID     *actor.PID
	Error   string
}
