package rollFileLog

import (
	"bytes"
	"errors"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"sync"
	"time"
)

const DEFAULT_LOG_CHAN_NUM int = 50000
const DEFAULT_NEED_SYNC_LOG_SIZE int = 10000

// 日志输出接口
type IWriteSyncer interface {
	Write(topic string, bs []byte) (int, error)
	Sync() error
	RegisterTopic(topic, path string)
}

var wg sync.WaitGroup
var allRotateLogWriteSyncer []*rotateLogWriteSyncer

func CloseAllLogFile() {
	for i := 0; i < len(allRotateLogWriteSyncer); i++ {
		allRotateLogWriteSyncer[i].Stop()
	}
	wg.Wait()
}

// 创建分割日志的writer
func NewRollingFileLog() IWriteSyncer {
	var rlws = newRotateLogWriteSyncer()
	allRotateLogWriteSyncer = append(allRotateLogWriteSyncer, rlws)
	return rlws
}

type logData struct {
	topic string
	data  *bytes.Buffer
}

func (d *logData) reset() {
	d.topic = ""
	d.data.Reset()

}

var pool = sync.Pool{
	New: func() interface{} {
		ret := &logData{}
		ret.data = bytes.NewBuffer(make([]byte, 0, 1024))
		return ret
	},
}

func get() *logData {
	return pool.Get().(*logData)
}
func put(d *logData) {
	d.reset()
	pool.Put(d)
}

type rotateLogWriteSyncer struct {
	logger       map[string]*rotatelogs.RotateLogs //多个文件
	buf          map[string]*bytes.Buffer
	logChan      chan *logData
	closeChan    chan interface{}
	needSyncSize int
}

func newRotateLogWriteSyncer() *rotateLogWriteSyncer {
	ws := &rotateLogWriteSyncer{

		buf:          make(map[string]*bytes.Buffer),
		logger:       make(map[string]*rotatelogs.RotateLogs),
		logChan:      make(chan *logData, DEFAULT_LOG_CHAN_NUM),
		closeChan:    make(chan interface{}),
		needSyncSize: DEFAULT_NEED_SYNC_LOG_SIZE,
	}
	wg.Add(1)
	go ws.run()
	return ws
}
func (l *rotateLogWriteSyncer) RegisterTopic(topic, path string) {
	rl, err := rotatelogs.New(
		// 没有使用go风格反人类的format格式
		path+".%Y%m%d.log",
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(err)
	}
	l.logger[topic] = rl
	l.buf[topic] = bytes.NewBuffer(make([]byte, 0, 1024))
}

func (l *rotateLogWriteSyncer) run() {
	ticker := time.NewTicker(1 * time.Second)
	defer wg.Done() //减去一个计数
	for {
		select {
		case <-ticker.C:
			l.sync()
		case bs := <-l.logChan:
			v, ok := l.buf[bs.topic]
			if ok {
				_, err := v.Write(bs.data.Bytes())
				if err != nil {
					continue
				}
				if v.Len() > l.needSyncSize {
					l.sync()
				}
			}
			put(bs)

		case <-l.closeChan:
			l.sync()
			return
		}
	}
}

func (l *rotateLogWriteSyncer) Stop() {
	close(l.closeChan)
}

var lineSep = []byte("\r\n")

func (l *rotateLogWriteSyncer) Write(topic string, bs []byte) (int, error) {
	var err error = nil
	lb := len(bs)
	if lb > 0 {
		lb += 2
		b := get()
		b.topic = topic
		b.data.Write(bs)
		b.data.Write(lineSep)
		select {
		case l.logChan <- b:
		default:
			err = errors.New("rotateLogWriteSyncer.logChan full")
		}
	}
	return 0, err
}

func (l *rotateLogWriteSyncer) Sync() error {
	return nil
}

func (l *rotateLogWriteSyncer) sync() error {
	for topic, _buff := range l.buf { //todo 每秒都要循环map ，需要优化
		if _buff.Len() > 0 {
			_, err := l.logger[topic].Write(_buff.Bytes())
			if err != nil {
				return err
			}
			_buff.Reset()
		}
	}
	return nil
}
