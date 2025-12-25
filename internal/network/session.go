package network

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
)

// Session represents a client connection
type Session struct {
	id       string
	conn     net.Conn
	codec    Codec
	handler  MessageHandler
	sendChan chan interface{}
	closed   atomic.Bool
	userdata sync.Map
}

// NewSession creates a new session
func NewSession(conn net.Conn, codec Codec, handler MessageHandler) *Session {
	return &Session{
		id:       generateSessionID(),
		conn:     conn,
		codec:    codec,
		handler:  handler,
		sendChan: make(chan interface{}, 100),
	}
}

// Start starts the session
func (s *Session) Start() {
	go s.readLoop()
	go s.writeLoop()
}

// Close closes the session
func (s *Session) Close() {
	if s.closed.CompareAndSwap(false, true) {
		close(s.sendChan)
		s.conn.Close()
	}
}

// Send sends a message to the client
func (s *Session) Send(message interface{}) error {
	if s.closed.Load() {
		return fmt.Errorf("session closed")
	}

	select {
	case s.sendChan <- message:
		return nil
	default:
		return fmt.Errorf("send channel full")
	}
}

// ID returns the session ID
func (s *Session) ID() string {
	return s.id
}

// RemoteAddr returns the remote address
func (s *Session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

// SetUserData stores user data
func (s *Session) SetUserData(key, value interface{}) {
	s.userdata.Store(key, value)
}

// GetUserData retrieves user data
func (s *Session) GetUserData(key interface{}) (interface{}, bool) {
	return s.userdata.Load(key)
}

// readLoop reads messages from the connection
func (s *Session) readLoop() {
	defer s.Close()

	for {
		message, err := s.codec.Decode(s.conn)
		if err != nil {
			fmt.Printf("Session %s decode error: %v\n", s.id, err)
			return
		}

		if err := s.handler.HandleMessage(s, message); err != nil {
			fmt.Printf("Session %s handle error: %v\n", s.id, err)
		}
	}
}

// writeLoop writes messages to the connection
func (s *Session) writeLoop() {
	defer s.Close()

	for message := range s.sendChan {
		if err := s.codec.Encode(s.conn, message); err != nil {
			fmt.Printf("Session %s encode error: %v\n", s.id, err)
			return
		}
	}
}

var sessionIDCounter atomic.Uint64

func generateSessionID() string {
	id := sessionIDCounter.Add(1)
	return fmt.Sprintf("session-%d", id)
}
