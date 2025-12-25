package actor

// envelope wraps a message with metadata
type envelope struct {
	message interface{}
	sender  *PID
}

// Mailbox is a message queue for an actor
type Mailbox struct {
	queue  chan *envelope
	closed bool
}

// NewMailbox creates a new mailbox with the specified buffer size
func NewMailbox(size int) *Mailbox {
	return &Mailbox{
		queue: make(chan *envelope, size),
	}
}

// Send adds a message to the mailbox
func (m *Mailbox) Send(env *envelope) {
	if m.closed {
		return
	}
	select {
	case m.queue <- env:
	default:
		// Mailbox is full, message is dropped
		// TODO: Add overflow strategy
	}
}

// Receive retrieves a message from the mailbox (blocking)
func (m *Mailbox) Receive() *envelope {
	env, ok := <-m.queue
	if !ok {
		return nil
	}
	return env
}

// Close closes the mailbox
func (m *Mailbox) Close() {
	if !m.closed {
		m.closed = true
		close(m.queue)
	}
}

// Size returns the current number of messages in the mailbox
func (m *Mailbox) Size() int {
	return len(m.queue)
}
