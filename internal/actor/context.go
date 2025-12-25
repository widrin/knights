package actor

import "time"

// Context provides methods for actors to interact with the system
type Context interface {
	// Message returns the current message being processed
	Message() interface{}

	// Sender returns the PID of the message sender
	Sender() *PID

	// Self returns the PID of the current actor
	Self() *PID

	// Send sends a message to the target PID
	Send(target *PID, message interface{})

	// Request sends a message and expects a response
	Request(target *PID, message interface{})

	// Respond sends a response to the current sender
	Respond(message interface{})

	// Spawn creates a new child actor
	Spawn(props *Props) *PID

	// SpawnNamed creates a new child actor with a specific name
	SpawnNamed(props *Props, name string) (*PID, error)

	// Stop stops the specified actor
	Stop(pid *PID)

	// StopSelf stops the current actor
	StopSelf()

	// Watch monitors another actor for termination
	Watch(pid *PID)

	// Unwatch stops monitoring an actor
	Unwatch(pid *PID)

	// SetReceiveTimeout sets a timeout for receiving messages
	SetReceiveTimeout(timeout time.Duration)

	// CancelReceiveTimeout cancels the receive timeout
	CancelReceiveTimeout()

	// Children returns all child actor PIDs
	Children() []*PID

	// Parent returns the parent actor PID
	Parent() *PID
}

// actorContext is the default implementation of Context
type actorContext struct {
	message  interface{}
	sender   *PID
	self     *PID
	parent   *PID
	children map[string]*PID
	system   *ActorSystem
	actor    Actor
}

func (c *actorContext) Message() interface{} {
	return c.message
}

func (c *actorContext) Sender() *PID {
	return c.sender
}

func (c *actorContext) Self() *PID {
	return c.self
}

func (c *actorContext) Send(target *PID, message interface{}) {
	// TODO: Implement message sending
}

func (c *actorContext) Request(target *PID, message interface{}) {
	// TODO: Implement request-response pattern
}

func (c *actorContext) Respond(message interface{}) {
	if c.sender != nil {
		c.Send(c.sender, message)
	}
}

func (c *actorContext) Spawn(props *Props) *PID {
	// TODO: Implement actor spawning
	return nil
}

func (c *actorContext) SpawnNamed(props *Props, name string) (*PID, error) {
	// TODO: Implement named actor spawning
	return nil, nil
}

func (c *actorContext) Stop(pid *PID) {
	// TODO: Implement actor stopping
}

func (c *actorContext) StopSelf() {
	c.Stop(c.self)
}

func (c *actorContext) Watch(pid *PID) {
	// TODO: Implement actor monitoring
}

func (c *actorContext) Unwatch(pid *PID) {
	// TODO: Implement unwatch
}

func (c *actorContext) SetReceiveTimeout(timeout time.Duration) {
	// TODO: Implement receive timeout
}

func (c *actorContext) CancelReceiveTimeout() {
	// TODO: Implement cancel timeout
}

func (c *actorContext) Children() []*PID {
	result := make([]*PID, 0, len(c.children))
	for _, pid := range c.children {
		result = append(result, pid)
	}
	return result
}

func (c *actorContext) Parent() *PID {
	return c.parent
}
