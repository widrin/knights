package log

var logger Logger

func init() {
	logger = NewZapLogger()
}

func Debug(format string, v ...any) {
	logger.Debug(format, v...)
}

func Info(format string, v ...any) {
	logger.Info(format, v...)
}

func Warn(format string, v ...any) {
	logger.Warn(format, v...)
}

func Error(format string, v ...any) {
	logger.Error(format, v...)
}

func Panic(format string, v ...any) {
	logger.Panic(format, v...)
}

func Sync() {
	logger.Sync()
}
