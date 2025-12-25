package timer

import (
	"time"

	"github.com/widrin/knights/internal/actor"
)

// TimerActor manages scheduled tasks
type TimerActor struct {
	timers map[string]*time.Timer
}

// NewTimerActor creates a new timer actor
func NewTimerActor() actor.Actor {
	return &TimerActor{
		timers: make(map[string]*time.Timer),
	}
}

// Receive handles timer messages
func (t *TimerActor) Receive(ctx actor.Context) {
	// TODO: Implement timer management
}
