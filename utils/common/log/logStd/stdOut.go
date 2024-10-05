package logStd

import (
	"io"
	slog "log"
	"os"
	"tzgit.kaixinxiyou.com/utils/common/buffer"
	"tzgit.kaixinxiyou.com/utils/common/log"
	"tzgit.kaixinxiyou.com/utils/common/log/syncWriter"
)

type StdWriter struct {
	syncWriter.SyncWriter
	writer io.Writer
}

func Create(queueCap uint32, level log.Level) *StdWriter {
	wr := &StdWriter{}
	wr.writer = os.Stdout
	wr.InitBufferWriter(wr, queueCap, 0, level)
	return wr
}

func (w *StdWriter) SyncWrite(b *buffer.Buffer) {
	_, err := w.writer.Write(b.Bytes())
	if err != nil {
		slog.Printf("StdWriter %v", err)
	}
}
