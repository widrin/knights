package network

import (
	"fmt"

	"github.com/widrin/knights/internal/actor"
)

// DefaultHandler is the default message handler
type DefaultHandler struct {
	actorSystem *actor.ActorSystem
	routes      map[string]*actor.PID
}

// NewDefaultHandler creates a new default handler
func NewDefaultHandler(actorSystem *actor.ActorSystem) MessageHandler {
	return &DefaultHandler{
		actorSystem: actorSystem,
		routes:      make(map[string]*actor.PID),
	}
}

// RegisterRoute registers a message route to an actor
func (h *DefaultHandler) RegisterRoute(messageType string, pid *actor.PID) {
	h.routes[messageType] = pid
}

// HandleMessage handles incoming messages
func (h *DefaultHandler) HandleMessage(session *Session, message interface{}) error {
	// TODO: Determine message type and route to appropriate actor
	fmt.Printf("Received message from session %s: %+v\n", session.ID(), message)

	// Example routing logic
	// messageType := getMessageType(message)
	// if pid, ok := h.routes[messageType]; ok {
	//     h.actorSystem.Send(pid, message)
	// }

	return nil
}
