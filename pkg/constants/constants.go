package constants

const (
	// Server constants
	DefaultServerPort = 8080
	DefaultMaxPlayers = 10000

	// Game constants
	DefaultRoomMaxPlayers = 4
	DefaultMatchmakingInterval = 1000 // milliseconds
	DefaultHeartbeatInterval = 30 // seconds

	// Network constants
	DefaultMaxMessageSize = 1024 * 1024 // 1MB
	DefaultReadBufferSize = 4096
	DefaultWriteBufferSize = 4096
)
