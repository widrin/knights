package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"tzgit.kaixinxiyou.com/utils/common/buffer"
)

type ISyncWriter interface {
	SyncWrite(b *buffer.Buffer)
	Close()
	Write(p []byte) (int, error)
	Ticker()
	OnClose()
	GetEncoder() zapcore.Encoder
	GetLevel() zap.LevelEnablerFunc
}

// A Level is a logging priority. Higher levels are more important.
type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel

	_minLevel = DebugLevel
	_maxLevel = FatalLevel
)

var Log = make(map[string]*zap.SugaredLogger)
var defaultLog *zap.SugaredLogger

var logList []ISyncWriter

// 日志组
var coreList = make(map[string][]zapcore.Core)

func RegisterLog(hook ISyncWriter, _group ...string) {
	logList = append(logList, hook)
	group := ""
	if len(_group) > 0 {
		group = _group[0]
	}
	coreList[group] = append(coreList[group], zapcore.NewCore(hook.GetEncoder(), zapcore.AddSync(hook), hook.GetLevel()))
	logger := zap.New(zapcore.NewTee(coreList[group]...), zap.AddCaller(), zap.AddCallerSkip(2), zap.AddStacktrace(zap.WarnLevel))
	Log[group] = logger.Sugar()
	if group == "" {
		defaultLog = Log[group]
	}
}
func RegisterLogWithField(hook ISyncWriter, fields ...interface{}) {
	logList = append(logList, hook)
	group := ""
	coreList[group] = append(coreList[group], zapcore.NewCore(hook.GetEncoder(), zapcore.AddSync(hook), hook.GetLevel()))
	logger := zap.New(zapcore.NewTee(coreList[group]...), zap.AddCaller(), zap.AddCallerSkip(2), zap.AddStacktrace(zap.WarnLevel))
	if len(fields) > 0 {
		f := make([]zap.Field, 0)
		for j := 0; j < len(fields); j += 2 {
			f = append(f, zap.Any(fields[j].(string), fields[j+1]))
		}
		logger = logger.With(f...)
	}

	Log[group] = logger.Sugar()
	defaultLog = Log[group]
}
func Zap_debug(format string, args ...interface{}) {
	if defaultLog != nil {
		defaultLog.Debugf(format, args...)
	}
}

// 指定log 名字输出
func Zap_DebugName(name string, format string, args ...interface{}) {
	if v, ok := Log[name]; ok {
		v.Infof(format, args...)
	}
}
func Zap_Release(format string, args ...interface{}) {
	if defaultLog != nil {
		defaultLog.Infof(format, args...)
	}
}
func Zap_Error(format string, args ...interface{}) {
	if defaultLog != nil {
		defaultLog.Errorf(format, args...)
	}
}
func Zap_Warn(format string, args ...interface{}) {
	if defaultLog != nil {
		defaultLog.Warnf(format, args...)
	}
}
func Zap_Fatal(format string, args ...interface{}) {
	if defaultLog != nil {
		defaultLog.Fatalf(format, args...)
	}

}
