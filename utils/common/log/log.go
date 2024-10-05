package log

import (
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"time"
)

var limit = rate.NewLimiter(rate.Every(time.Second), 1000)

func Debug(format string, a ...interface{}) {
	Zap_debug(format, a...)
}

func Release(format string, a ...interface{}) {
	Zap_Release(format, a...)
}

func Warn(format string, a ...interface{}) {
	if limit.Allow() {
		Zap_Warn(format, a...)
	}
}
func Error(format string, a ...interface{}) {
	if limit.Allow() {
		Zap_Error(format, a...)
	}
	//format = LightRed(format)
	//gLogger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func Fatal(format string, a ...interface{}) {
	//gLogger.doPrintf(fatalLevel, printFatalLevel, format, a...)
	log.Println(fmt.Sprintf(format, a...))
	Zap_Fatal(format, a...)

}

func Recover(r interface{}) {
	//buf := make([]byte, 4096)
	//l := runtime.Stack(buf, false)
	Error("catch err %v: ", r)
}
func Close() {
	log.Println("close>>>>>>>>>>>>>>")
	for _, v := range logList {
		v.Close()
	}
}
