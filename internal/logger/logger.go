package logger

import (
	"fmt"
	"log"
	"os"
)

// Logger provides logging interface
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
}

// defaultLogger is a simple logger implementation
type defaultLogger struct {
	logger *log.Logger
}

var std Logger

func init() {
	std = &defaultLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *defaultLogger) Debug(msg string, fields ...interface{}) {
	l.logger.Printf("[DEBUG] "+msg, fields...)
}

func (l *defaultLogger) Info(msg string, fields ...interface{}) {
	l.logger.Printf("[INFO] "+msg, fields...)
}

func (l *defaultLogger) Warn(msg string, fields ...interface{}) {
	l.logger.Printf("[WARN] "+msg, fields...)
}

func (l *defaultLogger) Error(msg string, fields ...interface{}) {
	l.logger.Printf("[ERROR] "+msg, fields...)
}

func (l *defaultLogger) Fatal(msg string, fields ...interface{}) {
	l.logger.Fatalf("[FATAL] "+msg, fields...)
}

// Global logger functions
func Debug(msg string, fields ...interface{}) {
	std.Debug(msg, fields...)
}

func Info(msg string, fields ...interface{}) {
	std.Info(msg, fields...)
}

func Warn(msg string, fields ...interface{}) {
	std.Warn(msg, fields...)
}

func Error(msg string, fields ...interface{}) {
	std.Error(msg, fields...)
}

func Fatal(msg string, fields ...interface{}) {
	std.Fatal(msg, fields...)
}

func Errorf(format string, args ...interface{}) {
	std.Error(fmt.Sprintf(format, args...))
}
