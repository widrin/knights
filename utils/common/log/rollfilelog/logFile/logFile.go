package logFile

import (
	"bytes"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"tzgit.kaixinxiyou.com/utils/common/buffer"
	"tzgit.kaixinxiyou.com/utils/common/log"
	rollFileLog "tzgit.kaixinxiyou.com/utils/common/log/rollfilelog"
	"tzgit.kaixinxiyou.com/utils/common/log/syncWriter"
)

// 日志文件输出
type logFile struct {
	syncWriter.SyncWriter
	clsTopic string
	syncer   rollFileLog.IWriteSyncer
}

func Create(clsTopic string, level log.Level, syncer rollFileLog.IWriteSyncer) log.ISyncWriter {
	r := &logFile{}
	//注册
	r.clsTopic = clsTopic
	r.syncer = syncer
	r.InitBufferWriter(r, 0, 0, level)
	return r
}
func (c *logFile) SyncWrite(b *buffer.Buffer) {
	c.syncer.Write(c.clsTopic, bytes.TrimRight(b.Bytes(), "\n"))
}
func (c *logFile) GetEncoder() zapcore.Encoder {
	kafkaEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	return kafkaEncoder
}

// 关闭事件
func (c *logFile) OnClose() {

}
