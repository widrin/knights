package actor

import (
	"sync/atomic"
)

// Router forwards messages to a group of actors
type Router interface {
	// Route selects an actor to receive the message
	Route(message interface{}) *PID
}

// RoundRobinRouter distributes messages in round-robin fashion
type RoundRobinRouter struct {
	routees []*PID
	index   atomic.Uint64
}

// NewRoundRobinRouter creates a new round-robin router
func NewRoundRobinRouter(routees []*PID) Router {
	return &RoundRobinRouter{
		routees: routees,
	}
}

func (r *RoundRobinRouter) Route(message interface{}) *PID {
	if len(r.routees) == 0 {
		return nil
	}
	idx := r.index.Add(1) % uint64(len(r.routees))
	return r.routees[idx]
}

// RandomRouter selects a random actor
type RandomRouter struct {
	routees []*PID
}

// NewRandomRouter creates a new random router
func NewRandomRouter(routees []*PID) Router {
	return &RandomRouter{
		routees: routees,
	}
}

func (r *RandomRouter) Route(message interface{}) *PID {
	if len(r.routees) == 0 {
		return nil
	}
	// TODO: Implement random selection
	return r.routees[0]
}

// BroadcastRouter sends messages to all actors
type BroadcastRouter struct {
	routees []*PID
}

// NewBroadcastRouter creates a new broadcast router
func NewBroadcastRouter(routees []*PID) Router {
	return &BroadcastRouter{
		routees: routees,
	}
}

func (r *BroadcastRouter) Route(message interface{}) *PID {
	// Broadcast doesn't return a single PID
	// This would need special handling in the actor system
	return nil
}
