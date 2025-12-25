package network

import (
	"sync"
)

// SessionManager manages all active sessions
type SessionManager struct {
	sessions sync.Map // map[string]*Session
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

// AddSession adds a session
func (sm *SessionManager) AddSession(session *Session) {
	sm.sessions.Store(session.ID(), session)
}

// RemoveSession removes a session
func (sm *SessionManager) RemoveSession(sessionID string) {
	sm.sessions.Delete(sessionID)
}

// GetSession retrieves a session by ID
func (sm *SessionManager) GetSession(sessionID string) (*Session, bool) {
	if session, ok := sm.sessions.Load(sessionID); ok {
		return session.(*Session), true
	}
	return nil, false
}

// CloseAll closes all sessions
func (sm *SessionManager) CloseAll() {
	sm.sessions.Range(func(key, value interface{}) bool {
		session := value.(*Session)
		session.Close()
		return true
	})
}

// Count returns the number of active sessions
func (sm *SessionManager) Count() int {
	count := 0
	sm.sessions.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// Broadcast sends a message to all sessions
func (sm *SessionManager) Broadcast(message interface{}) {
	sm.sessions.Range(func(key, value interface{}) bool {
		session := value.(*Session)
		session.Send(message)
		return true
	})
}
