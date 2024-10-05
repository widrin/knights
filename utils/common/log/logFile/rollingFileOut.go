package logFile

import (
	"bytes"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	slog "log"
	"os"
	"time"
	"tzgit.kaixinxiyou.com/utils/common/buffer"
	"tzgit.kaixinxiyou.com/utils/common/log"
	"tzgit.kaixinxiyou.com/utils/common/log/syncWriter"
)

const (
	output_dir = "/log/"
)

var workPath = ""

func init() {

	workPath, _ = os.Getwd()
}
func getLogPath() string {
	return workPath + output_dir
}

type rollingFileWrite struct {
	syncWriter.SyncWriter
	writer  io.Writer
	hook    *rotatelogs.RotateLogs
	buf     *bytes.Buffer
	maxSize int
}

func Create(fileName string, queueCap uint32, level log.Level) log.ISyncWriter {
	wr := &rollingFileWrite{}
	wr.writer = os.Stdout
	wr.buf = bytes.NewBuffer([]byte{})
	wr.maxSize = 1024 * 4
	var err error
	wr.hook, err = rotatelogs.New(
		// 没有使用go风格反人类的format格式
		getLogPath()+fileName+".%Y%m%d.log",
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(err)
	}

	wr.InitBufferWriter(wr, queueCap, 1, level)
	return wr
}

func (r *rollingFileWrite) SyncWrite(b *buffer.Buffer) {
	_, err := r.buf.Write(b.Bytes())
	if err != nil {
		slog.Printf("rollingFileWrite %v", err)
	}
	if r.buf.Len() > r.maxSize {
		r.sync()
	}
}

//关闭事件
func (r *rollingFileWrite) OnClose() {
	slog.Println("rollingFileWrite  close")
	if r.buf.Len() > 0 {
		r.sync()
	}
}

//每秒修改事件
func (r *rollingFileWrite) Ticker() {
	if r.buf.Len() > 0 {
		r.sync()
	}
}
func (r *rollingFileWrite) sync() error {
	defer r.buf.Reset()
	_, err := r.hook.Write(r.buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}
