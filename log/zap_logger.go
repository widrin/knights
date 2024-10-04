package log

import (
	"fmt"

	"go.uber.org/zap"
)

type zapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger() *zapLogger {
	logger, err := zap.NewDevelopment()

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &zapLogger{
		logger: logger.Sugar(),
	}
}

func (l *zapLogger) Debug(format string, v ...any) {
	l.logger.Debugf(format+"\n", v...)
}

func (l *zapLogger) Info(format string, v ...any) {
	l.logger.Infof(format+"\n", v...)
}

func (l *zapLogger) Warn(format string, v ...any) {
	l.logger.Warnf(format+"\n", v...)
}

func (l *zapLogger) Error(format string, v ...any) {
	l.logger.Errorf(format+"\n", v...)
}

func (l *zapLogger) Panic(format string, v ...any) {
	l.logger.Panicf(format+"\n", v...)
}

func (l *zapLogger) Sync() {
	l.logger.Sync()
}
