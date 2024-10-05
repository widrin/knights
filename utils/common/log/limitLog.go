package log

import (
	"golang.org/x/time/rate"
	"time"
)

type LimitLog struct {
	limit *rate.Limiter
}

//interval 多少时间增加一个，maxNum：最多可以累积多少个
func Create(interval time.Duration, maxNum int) *LimitLog {
	l := new(LimitLog)
	l.limit = rate.NewLimiter(rate.Every(interval), maxNum)
	return l
}
func (l *LimitLog) Error(format string, a ...interface{}) {
	if l.limit.Allow() {
		Error(format, a...)
	}
}
