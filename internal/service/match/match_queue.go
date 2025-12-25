package match

import (
	"sync"
	"time"

	"github.com/widrin/knights/internal/actor"
)

// MatchPlayer represents a player in the matchmaking queue
type MatchPlayer struct {
	PlayerID   string
	PlayerPID  *actor.PID
	Rating     int
	JoinedAt   time.Time
}

// Match represents a matched group of players
type Match struct {
	Players []*MatchPlayer
}

// MatchQueue manages the matchmaking queue
type MatchQueue struct {
	mu      sync.RWMutex
	players map[string]*MatchPlayer
}

// NewMatchQueue creates a new match queue
func NewMatchQueue() *MatchQueue {
	return &MatchQueue{
		players: make(map[string]*MatchPlayer),
	}
}

// AddPlayer adds a player to the queue
func (q *MatchQueue) AddPlayer(player *MatchPlayer) {
	q.mu.Lock()
	defer q.mu.Unlock()
	player.JoinedAt = time.Now()
	q.players[player.PlayerID] = player
}

// RemovePlayer removes a player from the queue
func (q *MatchQueue) RemovePlayer(playerID string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.players, playerID)
}

// CreateMatches attempts to create matches from queued players
func (q *MatchQueue) CreateMatches() []*Match {
	q.mu.Lock()
	defer q.mu.Unlock()

	matches := make([]*Match, 0)

	// Simple matchmaking: group players in pairs
	// TODO: Implement more sophisticated matchmaking algorithm
	var currentMatch []*MatchPlayer
	for _, player := range q.players {
		currentMatch = append(currentMatch, player)
		if len(currentMatch) >= 2 {
			matches = append(matches, &Match{
				Players: currentMatch,
			})
			// Remove matched players
			for _, p := range currentMatch {
				delete(q.players, p.PlayerID)
			}
			currentMatch = nil
		}
	}

	return matches
}

// Messages for matchmaker

type JoinMatchRequest struct {
	PlayerID  string
	PlayerPID *actor.PID
	Rating    int
}

type JoinMatchResponse struct {
	Success bool
	Error   string
}

type CancelMatchRequest struct {
	PlayerID string
}

type CancelMatchResponse struct {
	Success bool
	Error   string
}

type TickMessage struct{}
