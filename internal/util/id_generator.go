package util

import (
	"fmt"
	"sync/atomic"
	"time"
)

// IDGenerator generates unique IDs
type IDGenerator struct {
	counter atomic.Uint64
}

// NewIDGenerator creates a new ID generator
func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

// Next generates the next ID
func (g *IDGenerator) Next() uint64 {
	return g.counter.Add(1)
}

// NextString generates the next ID as string
func (g *IDGenerator) NextString() string {
	return fmt.Sprintf("%d", g.Next())
}

// Snowflake generates a Twitter Snowflake-like ID
// Format: timestamp (41 bits) + machine ID (10 bits) + sequence (12 bits)
func Snowflake(machineID uint16) uint64 {
	now := time.Now().UnixMilli()
	// Simplified snowflake implementation
	// TODO: Add proper sequence handling and machine ID
	return uint64(now) << 22
}
