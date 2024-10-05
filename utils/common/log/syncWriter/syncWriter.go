package syncWriter

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"math"
	"sync"
	"time"
	"tzgit.kaixinxiyou.com/utils/common/buffer"
	"tzgit.kaixinxiyou.com/utils/common/log"
)

const defaultQueueCap = math.MaxUint16 * 8

var QueueIsFullError = errors.New("queue is full")
var DropWriteMessageError = errors.New("message writing failure and drop it")

type SyncWriter struct {
	bufferPool buffer.Pool
	wg         sync.WaitGroup
	lock       sync.RWMutex
	channel    chan *buffer.Buffer
	self       log.ISyncWriter
	ticker     int32
	level      func(zapcore.Level) bool
}

func (w *SyncWriter) InitBufferWriter(self log.ISyncWriter, queueCap uint32, ticker int32, level log.Level) {
	if queueCap <= 0 {
		queueCap = defaultQueueCap
	}

	w.bufferPool = buffer.NewPool()
	w.channel = make(chan *buffer.Buffer, queueCap)
	w.wg.Add(1)
	w.self = self
	if ticker == 0 {
		w.ticker = 60 * 60
	} else {
		w.ticker = ticker
	}
	w.level = GetLevelFunc(level)
	go w.poller()

}

var LimitLevel = zap.DebugLevel
var (

	// 实现两个判断日志等级的interface
	debugLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
	infoLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	warnLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})
	errorLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
)

func GetLevelFunc(level log.Level) zap.LevelEnablerFunc {
	switch level {
	case log.DebugLevel:
		return debugLevel
	case log.InfoLevel:
		return infoLevel
	case log.WarnLevel:
		return warnLevel
	case log.ErrorLevel:
		return errorLevel
	default:
		return debugLevel
	}
}

func (w *SyncWriter) Write(p []byte) (int, error) {
	if w.lock.TryRLock() {
		defer w.lock.RUnlock()
		b := w.bufferPool.Get()
		count, err := b.Write(p)
		if err != nil {
			b.Free()
			return count, err
		}
		select {
		case w.channel <- b: // channel 内部传递的是 buffer 的指针，速度比传递对象快。
			break
		default:
			b.Free()
			return count, QueueIsFullError
		}
		return len(p), nil
	} else {
		return -1, DropWriteMessageError
	}
}
func (w *SyncWriter) Close() {
	w.lock.Lock()
	close(w.channel)
	w.wg.Wait()
	w.lock.Unlock()
}
func (w *SyncWriter) poller() {
	var (
		eb *buffer.Buffer
		//err error
	)
	ticker := time.NewTicker(time.Duration(w.ticker) * time.Second)
	defer w.wg.Done()
EXIT:
	for {
		select {
		case <-ticker.C:
			w.self.Ticker()
		case eb = <-w.channel:
			if eb == nil {
				break EXIT
			}
			w.self.SyncWrite(eb)
			eb.Free()
		}
	}
	w.self.OnClose()
	//for eb = range w.channel {
	//	w.self.SyncWrite(eb)
	//	//_, err = w.writer.Write(eb.Bytes())
	//	//if err != nil {
	//	//	log.Printf("log error %v", err.Error())
	//	//}
	//	//这里实现异步写入
	//
	//}

}
func (w *SyncWriter) SyncWrite(b *buffer.Buffer) {

}
func (w *SyncWriter) Ticker() {

}
func (w *SyncWriter) OnClose() {

}
func (w *SyncWriter) GetEncoder() zapcore.Encoder {
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey: "msg",
		LevelKey:   "level",
		TimeKey:    "ts",
		//CallerKey:      "file",
		CallerKey:     "caller",
		StacktraceKey: "trace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		//EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
	return encoder
}
func (w *SyncWriter) GetLevel() zap.LevelEnablerFunc {
	return w.level
}
