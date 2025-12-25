package actor

import "time"

// SupervisorStrategy defines how a supervisor handles child actor failures
type SupervisorStrategy interface {
	// HandleFailure is called when a child actor fails
	HandleFailure(child *PID, reason interface{}) Directive
}

// Directive tells the supervisor what action to take
type Directive int

const (
	// ResumeDirective continues processing the next message
	ResumeDirective Directive = iota

	// RestartDirective stops the actor and creates a new instance
	RestartDirective

	// StopDirective permanently stops the actor
	StopDirective

	// EscalateDirective passes the failure to the parent supervisor
	EscalateDirective
)

// OneForOneStrategy restarts only the failed actor
type OneForOneStrategy struct {
	maxRetries int
	window     time.Duration
}

// NewOneForOneStrategy creates a new OneForOne supervisor strategy
func NewOneForOneStrategy(maxRetries int, window time.Duration) SupervisorStrategy {
	return &OneForOneStrategy{
		maxRetries: maxRetries,
		window:     window,
	}
}

func (s *OneForOneStrategy) HandleFailure(child *PID, reason interface{}) Directive {
	// TODO: Implement retry counting within time window
	return RestartDirective
}

// AllForOneStrategy restarts all sibling actors when one fails
type AllForOneStrategy struct {
	maxRetries int
	window     time.Duration
}

// NewAllForOneStrategy creates a new AllForOne supervisor strategy
func NewAllForOneStrategy(maxRetries int, window time.Duration) SupervisorStrategy {
	return &AllForOneStrategy{
		maxRetries: maxRetries,
		window:     window,
	}
}

func (s *AllForOneStrategy) HandleFailure(child *PID, reason interface{}) Directive {
	// TODO: Implement retry counting and restart all children
	return RestartDirective
}
