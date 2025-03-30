package logrus

import (
	"os"
	"sync/atomic"

	"github.com/widrin/knights/logger"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogrusLogger struct {
	logger *logrus.Logger
	level  *atomic.Int32
}

func NewLogrusLogger(cfg logger.Config) *LogrusLogger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z07:00",
	})

	writer := &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	// 同时输出到文件和控制台
	logger.SetOutput(os.Stdout)
	logger.AddHook(&writerHook{writer})

	return &LogrusLogger{
		logger: logger,
		level:  new(atomic.Int32),
	}
}

type writerHook struct {
	writer *lumberjack.Logger
}

func (h *writerHook) Fire(entry *logrus.Entry) error {
	line, _ := entry.String()
	h.writer.Write([]byte(line))
	return nil
}

func (h *writerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (l *LogrusLogger) Debug(msg string, fields ...interface{}) {
	l.logger.Debugf(msg, fields...)
}

func (l *LogrusLogger) Info(msg string, fields ...interface{}) {
	l.logger.Infof(msg, fields...)
}

func (l *LogrusLogger) Warn(msg string, fields ...interface{}) {
	l.logger.Warnf(msg, fields...)
}

func (l *LogrusLogger) Error(msg string, fields ...interface{}) {
	l.logger.Errorf(msg, fields...)
}

func (l *LogrusLogger) With(fields ...interface{}) *LogrusLogger {
	return &LogrusLogger{
		logger: l.logger.WithFields(toLogrusFields(fields...)).Logger,
		level:  l.level,
	}
}

func (l *LogrusLogger) SetLevel(level int32) {
	l.level.Store(level)
	l.logger.SetLevel(logrus.Level(l.level.Load()))
}

func toLogrusFields(fields ...interface{}) logrus.Fields {
	// 转换字段到logrus支持的格式...

	// 这里假设函数返回一个空的 logrus.Fields 类型
	return logrus.Fields{}
}
