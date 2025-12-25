package network

import (
	"fmt"
	"net"
	"sync"

	"github.com/widrin/knights/internal/actor"
)

// Server represents a network server
type Server struct {
	address        string
	listener       net.Listener
	sessionManager *SessionManager
	codec          Codec
	handler        MessageHandler
	actorSystem    *actor.ActorSystem
	running        bool
	mu             sync.RWMutex
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Address        string
	Codec          Codec
	Handler        MessageHandler
	ActorSystem    *actor.ActorSystem
}

// NewServer creates a new network server
func NewServer(config *ServerConfig) *Server {
	return &Server{
		address:     config.Address,
		codec:       config.Codec,
		handler:     config.Handler,
		actorSystem: config.ActorSystem,
		sessionManager: NewSessionManager(),
	}
}

// Start starts the server
func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("server already running")
	}

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.listener = listener
	s.running = true

	fmt.Printf("Server listening on %s\n", s.address)

	go s.acceptLoop()

	return nil
}

// Stop stops the server
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	s.running = false
	if s.listener != nil {
		s.listener.Close()
	}

	// Close all sessions
	s.sessionManager.CloseAll()

	return nil
}

// acceptLoop accepts incoming connections
func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if !s.isRunning() {
				return
			}
			fmt.Printf("Accept error: %v\n", err)
			continue
		}

		session := NewSession(conn, s.codec, s.handler)
		s.sessionManager.AddSession(session)

		go session.Start()
	}
}

func (s *Server) isRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// MessageHandler processes decoded messages
type MessageHandler interface {
	HandleMessage(session *Session, message interface{}) error
}
