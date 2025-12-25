package match

import (
	"github.com/widrin/knights/internal/actor"
)

// MatchmakerActor handles player matchmaking
type MatchmakerActor struct {
	queue *MatchQueue
}

// NewMatchmakerActor creates a new matchmaker actor
func NewMatchmakerActor() actor.Actor {
	return &MatchmakerActor{
		queue: NewMatchQueue(),
	}
}

// Receive handles incoming messages
func (m *MatchmakerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		m.onStarted(ctx)
	case *JoinMatchRequest:
		m.handleJoinMatch(ctx, msg)
	case *CancelMatchRequest:
		m.handleCancelMatch(ctx, msg)
	case *TickMessage:
		m.handleTick(ctx)
	}
}

func (m *MatchmakerActor) onStarted(ctx actor.Context) {
	// Start matchmaking tick
}

func (m *MatchmakerActor) handleJoinMatch(ctx actor.Context, msg *JoinMatchRequest) {
	m.queue.AddPlayer(&MatchPlayer{
		PlayerID:  msg.PlayerID,
		PlayerPID: msg.PlayerPID,
		Rating:    msg.Rating,
	})

	ctx.Respond(&JoinMatchResponse{
		Success: true,
	})
}

func (m *MatchmakerActor) handleCancelMatch(ctx actor.Context, msg *CancelMatchRequest) {
	m.queue.RemovePlayer(msg.PlayerID)

	ctx.Respond(&CancelMatchResponse{
		Success: true,
	})
}

func (m *MatchmakerActor) handleTick(ctx actor.Context) {
	// Try to create matches
	matches := m.queue.CreateMatches()

	for _, match := range matches {
		// TODO: Create room and notify players
		_ = match
	}
}
