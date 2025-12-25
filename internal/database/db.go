package database

import "context"

// Database provides database operations interface
type Database interface {
	Connect(ctx context.Context) error
	Close() error
	Query(ctx context.Context, query string, args ...interface{}) (interface{}, error)
	Execute(ctx context.Context, query string, args ...interface{}) error
}
