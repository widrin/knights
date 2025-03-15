package logrus

import (
	"os"
	"sync/atomic"

	"github.com/sirupsen/logrus"
	"github.com:widrin/knights/logger"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogrusLogger struct {
	logger *logrus.Logger
	level  *atomic.Int32
}

func NewLogrusLogger(cfg logger.Config) LogrusLogger {
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
		level:  atomic.NewInt32(int32(parseLevel(cfg.Level))),
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

// 实现Logger接口所有方法...
