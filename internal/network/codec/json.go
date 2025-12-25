package codec

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

// JSONCodec implements JSON-based encoding/decoding
type JSONCodec struct {
	maxMessageSize uint32
}

// NewJSONCodec creates a new JSON codec
func NewJSONCodec() Codec {
	return &JSONCodec{
		maxMessageSize: 1024 * 1024, // 1MB
	}
}

// Encode encodes a message to JSON
func (c *JSONCodec) Encode(conn net.Conn, message interface{}) error {
	// Marshal to JSON
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	// Write length prefix (4 bytes)
	length := uint32(len(data))
	if err := binary.Write(conn, binary.BigEndian, length); err != nil {
		return fmt.Errorf("write length error: %w", err)
	}

	// Write data
	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("write data error: %w", err)
	}

	return nil
}

// Decode decodes a JSON message
func (c *JSONCodec) Decode(conn net.Conn) (interface{}, error) {
	// Read length prefix
	var length uint32
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("read length error: %w", err)
	}

	// Check max size
	if length > c.maxMessageSize {
		return nil, fmt.Errorf("message too large: %d bytes", length)
	}

	// Read data
	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, fmt.Errorf("read data error: %w", err)
	}

	// Unmarshal JSON
	var message map[string]interface{}
	if err := json.Unmarshal(data, &message); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	return message, nil
}
