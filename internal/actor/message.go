package actor

// System messages for actor lifecycle management

// Started is sent when an actor is started
type Started struct{}

// Stopping is sent when an actor is stopping
type Stopping struct{}

// Stopped is sent when an actor has stopped
type Stopped struct{}

// Restarting is sent when an actor is restarting
type Restarting struct{}

// Terminated is sent when a watched actor terminates
type Terminated struct {
	Who *PID
}

// ReceiveTimeout is sent when the receive timeout expires
type ReceiveTimeout struct{}

// ReceiveMiddleware allows intercepting and modifying message processing
type ReceiveMiddleware func(next ReceiveFunc) ReceiveFunc

// ReceiveFunc is a function that processes a message
type ReceiveFunc func(ctx Context)
