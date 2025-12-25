package actor

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ActorSystem manages the lifecycle of all actors
type ActorSystem struct {
	name       string
	address    string
	actors     sync.Map // map[string]*actorProcess
	idCounter  atomic.Uint64
	dispatcher Dispatcher
	root       *PID
}

// NewActorSystem creates a new actor system
func NewActorSystem(name string) *ActorSystem {
	system := &ActorSystem{
		name:       name,
		address:    "local",
		dispatcher: NewDefaultDispatcher(100), // 100 worker goroutines
	}
	return system
}

// Spawn creates a new root-level actor
func (s *ActorSystem) Spawn(props *Props) *PID {
	id := s.generateID()
	return s.spawnActor(props, id, nil)
}

// SpawnNamed creates a new named root-level actor
func (s *ActorSystem) SpawnNamed(props *Props, name string) (*PID, error) {
	if _, exists := s.actors.Load(name); exists {
		return nil, fmt.Errorf("actor with name %s already exists", name)
	}
	return s.spawnActor(props, name, nil), nil
}

// spawnActor creates a new actor with the given ID and parent
func (s *ActorSystem) spawnActor(props *Props, id string, parent *PID) *PID {
	pid := NewPID(s.address, id)

	// Create mailbox
	mailbox := NewMailbox(props.mailboxSize)

	// Create actor context
	ctx := &actorContext{
		self:     pid,
		parent:   parent,
		children: make(map[string]*PID),
		system:   s,
	}

	// Create actor instance
	actor := props.producer()
	ctx.actor = actor

	// Create actor process
	process := &actorProcess{
		pid:     pid,
		actor:   actor,
		context: ctx,
		mailbox: mailbox,
		system:  s,
	}

	// Register actor
	s.actors.Store(id, process)

	// Start processing messages
	go process.run()

	return pid
}

// Send sends a message to an actor
func (s *ActorSystem) Send(target *PID, message interface{}) {
	s.SendWithSender(target, nil, message)
}

// SendWithSender sends a message to an actor with sender information
func (s *ActorSystem) SendWithSender(target *PID, sender *PID, message interface{}) {
	if !target.IsLocal() {
		// TODO: Handle remote messaging
		return
	}

	if process, ok := s.actors.Load(target.ID); ok {
		proc := process.(*actorProcess)
		proc.mailbox.Send(&envelope{
			message: message,
			sender:  sender,
		})
	}
}

// Stop stops an actor
func (s *ActorSystem) Stop(pid *PID) {
	if process, ok := s.actors.LoadAndDelete(pid.ID); ok {
		proc := process.(*actorProcess)
		proc.stop()
	}
}

// Shutdown gracefully shuts down the actor system
func (s *ActorSystem) Shutdown() {
	// Stop all actors
	s.actors.Range(func(key, value interface{}) bool {
		proc := value.(*actorProcess)
		proc.stop()
		return true
	})
}

// generateID generates a unique actor ID
func (s *ActorSystem) generateID() string {
	id := s.idCounter.Add(1)
	return fmt.Sprintf("actor-%d", id)
}

// actorProcess represents a running actor
type actorProcess struct {
	pid     *PID
	actor   Actor
	context *actorContext
	mailbox *Mailbox
	system  *ActorSystem
	stopped atomic.Bool
}

// run processes messages from the mailbox
func (p *actorProcess) run() {
	for {
		if p.stopped.Load() {
			return
		}

		env := p.mailbox.Receive()
		if env == nil {
			continue
		}

		p.context.message = env.message
		p.context.sender = env.sender

		// Process message
		p.actor.Receive(p.context)
	}
}

// stop stops the actor process
func (p *actorProcess) stop() {
	p.stopped.Store(true)
	p.mailbox.Close()
}
