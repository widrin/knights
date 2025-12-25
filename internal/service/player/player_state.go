package player

import (
	"time"

	"github.com/widrin/knights/internal/actor"
)

// PlayerState holds the player's game state
type PlayerState struct {
	PlayerID     string
	Name         string
	Level        int
	Position     Position
	HP           int
	MaxHP        int
	MP           int
	MaxMP        int
	Gold         int
	LastLoginAt  time.Time
	CreatedAt    time.Time
}

// Position represents a 3D position in the game world
type Position struct {
	X float64
	Y float64
	Z float64
}

// NewPlayerState creates a new player state
func NewPlayerState(playerID string) *PlayerState {
	now := time.Now()
	return &PlayerState{
		PlayerID:    playerID,
		Level:       1,
		Position:    Position{X: 0, Y: 0, Z: 0},
		HP:          100,
		MaxHP:       100,
		MP:          50,
		MaxMP:       50,
		Gold:        0,
		LastLoginAt: now,
		CreatedAt:   now,
	}
}

// Messages for player actor

type LoginRequest struct {
	PlayerID string
	Token    string
}

type LoginResponse struct {
	Success  bool
	PlayerID string
	Error    string
}

type LogoutRequest struct{}

type LogoutResponse struct {
	Success bool
	Error   string
}

type MoveRequest struct {
	Position Position
}

type MoveResponse struct {
	Success  bool
	Position Position
	Error    string
}

// Messages for player manager

type CreatePlayerRequest struct {
	PlayerID string
}

type CreatePlayerResponse struct {
	Success  bool
	PlayerID string
	PID      *actor.PID
	Error    string
}

type RemovePlayerRequest struct {
	PlayerID string
}

type RemovePlayerResponse struct {
	Success bool
	Error   string
}

type GetPlayerRequest struct {
	PlayerID string
}

type GetPlayerResponse struct {
	Success bool
	PID     *actor.PID
	Error   string
}
