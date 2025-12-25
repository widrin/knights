package actor

import "fmt"

// PID (Process ID) uniquely identifies an actor
type PID struct {
	Address string // Node address for remote actors
	ID      string // Unique actor identifier
}

// NewPID creates a new PID
func NewPID(address, id string) *PID {
	return &PID{
		Address: address,
		ID:      id,
	}
}

// String returns the string representation of the PID
func (pid *PID) String() string {
	if pid.Address == "" {
		return pid.ID
	}
	return fmt.Sprintf("%s/%s", pid.Address, pid.ID)
}

// IsLocal returns true if the actor is on the local node
func (pid *PID) IsLocal() bool {
	return pid.Address == "" || pid.Address == "local"
}

// Equal checks if two PIDs are equal
func (pid *PID) Equal(other *PID) bool {
	if pid == nil || other == nil {
		return pid == other
	}
	return pid.Address == other.Address && pid.ID == other.ID
}
