package errors

import "fmt"

// ErrorCode represents an error code
type ErrorCode int

const (
	// Common errors
	ErrCodeSuccess ErrorCode = 0
	ErrCodeUnknown ErrorCode = 1000
	ErrCodeInvalidParam ErrorCode = 1001
	ErrCodeUnauthorized ErrorCode = 1002
	ErrCodeForbidden ErrorCode = 1003
	ErrCodeNotFound ErrorCode = 1004
	ErrCodeTimeout ErrorCode = 1005

	// Player errors
	ErrCodePlayerNotFound ErrorCode = 2001
	ErrCodePlayerAlreadyExists ErrorCode = 2002
	ErrCodePlayerOffline ErrorCode = 2003

	// Room errors
	ErrCodeRoomNotFound ErrorCode = 3001
	ErrCodeRoomFull ErrorCode = 3002
	ErrCodeRoomAlreadyExists ErrorCode = 3003
	ErrCodeNotInRoom ErrorCode = 3004

	// Match errors
	ErrCodeMatchFailed ErrorCode = 4001
	ErrCodeMatchTimeout ErrorCode = 4002
)

// GameError represents a game-specific error
type GameError struct {
	Code    ErrorCode
	Message string
}

// Error implements the error interface
func (e *GameError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// New creates a new GameError
func New(code ErrorCode, message string) *GameError {
	return &GameError{
		Code:    code,
		Message: message,
	}
}

// Newf creates a new GameError with formatted message
func Newf(code ErrorCode, format string, args ...interface{}) *GameError {
	return &GameError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}
