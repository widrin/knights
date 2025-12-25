package actor

// Actor defines the interface that all actors must implement
type Actor interface {
	// Receive processes incoming messages
	Receive(ctx Context)
}

// Producer is a function that creates a new actor instance
type Producer func() Actor

// Props contains configuration for creating an actor
type Props struct {
	producer       Producer
	mailboxSize    int
	dispatcher     Dispatcher
	supervisor     SupervisorStrategy
	middlewares    []ReceiveMiddleware
}

// NewProps creates a new Props with the given producer
func NewProps(producer Producer) *Props {
	return &Props{
		producer:    producer,
		mailboxSize: 100, // default mailbox size
	}
}

// WithMailboxSize sets the mailbox buffer size
func (p *Props) WithMailboxSize(size int) *Props {
	p.mailboxSize = size
	return p
}

// WithDispatcher sets a custom dispatcher
func (p *Props) WithDispatcher(dispatcher Dispatcher) *Props {
	p.dispatcher = dispatcher
	return p
}

// WithSupervisor sets the supervisor strategy
func (p *Props) WithSupervisor(strategy SupervisorStrategy) *Props {
	p.supervisor = strategy
	return p
}

// WithMiddleware adds a receive middleware
func (p *Props) WithMiddleware(middleware ReceiveMiddleware) *Props {
	p.middlewares = append(p.middlewares, middleware)
	return p
}
