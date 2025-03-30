package logger

import (
	"errors"
	"sync"
)

var (
	gLogger Logger
	once    sync.Once
)

func InitLogger(l Logger) {
	once.Do(func() {
		gLogger = l
	})
}

func getInstance() Logger {
	if gLogger == nil {
		panic(errors.New("logger instance not initialized, call SetGlobalLogger first"))
	}
	return gLogger
}

// 全局代理方法
func Debug(msg string, fields ...interface{}) {
	getInstance().Debug(msg, fields...)
}

func Info(msg string, fields ...interface{}) {
	getInstance().Info(msg, fields...)
}

func Warn(msg string, fields ...interface{}) {
	getInstance().Warn(msg, fields...)
}

func Error(msg string, fields ...interface{}) {
	getInstance().Error(msg, fields...)
}
