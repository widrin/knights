package network

import "net"

// Codec defines the interface for encoding/decoding messages
type Codec interface {
	// Encode encodes a message and writes it to the connection
	Encode(conn net.Conn, message interface{}) error

	// Decode reads from the connection and decodes a message
	Decode(conn net.Conn) (interface{}, error)
}
