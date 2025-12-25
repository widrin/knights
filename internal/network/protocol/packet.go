package protocol

// PacketType represents different packet types
type PacketType uint16

const (
	PacketTypeHandshake PacketType = iota
	PacketTypeHeartbeat
	PacketTypeData
	PacketTypeKick
)

// Packet represents a network packet
type Packet struct {
	Type    PacketType
	ID      uint32
	RouteID uint16
	Data    []byte
}

// Message represents a decoded message
type Message struct {
	Route string
	Data  interface{}
}
