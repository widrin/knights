package zap

import (
	"github.com:widrin/knights/logger"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
	core   zapcore.Core
	level  *atomic.Int32
}

func NewZapLogger(cfg logger.Config) ZapLogger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
	}

	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writer,
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.Level(cfg.level.Load())
		}),
	)

	return &ZapLogger{
		logger: zap.New(core).Sugar(),
		core:   core,
		level:  atomic.NewInt32(int32(parseLevel(cfg.Level))),
	}
}

func (z *ZapLogger) Debug(msg string, fields ...interface{}) {
	z.logger.Debugw(msg, fields...)
}

func (z *ZapLogger) With(fields ...interface{}) Logger {
	return &ZapLogger{
		logger: z.logger.With(fields...),
		core:   z.core,
		level:  z.level,
	}
}

func (z *ZapLogger) SetLevel(level string) {
	z.level.Store(int32(parseLevel(level)))
}
