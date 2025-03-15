package logger

// 通用日志接口
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	With(fields ...interface{}) Logger
	SetLevel(level string)
}

// 日志配置
type Config struct {
	Driver     string // zap/logrus
	Level      string
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// 创建日志实例
func NewLogger(cfg Config) Logger {
	switch cfg.Driver {
	case "zap":
		return newZapLogger(cfg)
	case "logrus":
		return newLogrusLogger(cfg)
	default:
		panic("unsupported logger driver")
	}
}
