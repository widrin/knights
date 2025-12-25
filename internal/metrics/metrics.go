package metrics

import "sync/atomic"

// Metrics collects system metrics
type Metrics struct {
	ActivePlayers   atomic.Int64
	ActiveRooms     atomic.Int64
	MessagesSent    atomic.Uint64
	MessagesReceived atomic.Uint64
}

var globalMetrics = &Metrics{}

// Global returns the global metrics instance
func Global() *Metrics {
	return globalMetrics
}

// IncrementActivePlayers increments active player count
func (m *Metrics) IncrementActivePlayers() {
	m.ActivePlayers.Add(1)
}

// DecrementActivePlayers decrements active player count
func (m *Metrics) DecrementActivePlayers() {
	m.ActivePlayers.Add(-1)
}

// IncrementActiveRooms increments active room count
func (m *Metrics) IncrementActiveRooms() {
	m.ActiveRooms.Add(1)
}

// DecrementActiveRooms decrements active room count
func (m *Metrics) DecrementActiveRooms() {
	m.ActiveRooms.Add(-1)
}

// IncrementMessagesSent increments messages sent counter
func (m *Metrics) IncrementMessagesSent() {
	m.MessagesSent.Add(1)
}

// IncrementMessagesReceived increments messages received counter
func (m *Metrics) IncrementMessagesReceived() {
	m.MessagesReceived.Add(1)
}

// GetActivePlayers returns active player count
func (m *Metrics) GetActivePlayers() int64 {
	return m.ActivePlayers.Load()
}

// GetActiveRooms returns active room count
func (m *Metrics) GetActiveRooms() int64 {
	return m.ActiveRooms.Load()
}
