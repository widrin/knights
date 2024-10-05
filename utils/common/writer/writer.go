package writer

import (
	"bytes"
	"fmt"
	"strings"
)

type Writer struct {
	Content      bytes.Buffer
	indent       string
	LineBreak    string
	NotPrevCount bool //为true 说明没有prevCount
}

func Create() *Writer {
	return &Writer{}
}
func (writer *Writer) AddLine(msg string) {
	if strings.HasSuffix(msg, "}") && !strings.HasSuffix(msg, "; }") {
		writer.Out()
	}
	if strings.HasSuffix(msg, "});") {
		writer.Out()
	}

	if !writer.NotPrevCount && len(writer.indent) > 0 {
		writer.Content.WriteString(writer.indent)
	}

	writer.Content.WriteString(msg)
	if writer.LineBreak == "" {
		writer.Content.WriteString("\r\n")
	} else {
		writer.Content.WriteString(writer.LineBreak)
	}

	if strings.HasSuffix(msg, "{") {
		writer.In()
	}
}
func (writer *Writer) AddLineFmt(format string, a ...interface{}) {
	writer.AddLine(fmt.Sprintf(format, a...))
}
func (writer *Writer) P(str ...string) {
	writer.Content.WriteString(writer.indent)
	for _, v := range str {
		writer.Content.WriteString(v)
	}
	if writer.LineBreak == "" {
		writer.Content.WriteString("\r\n")
	} else {
		writer.Content.WriteString(writer.LineBreak)
	}

}

// In Indents the output one tab stop.
func (writer *Writer) In() { writer.indent += "\t" }

// Out unindents the output one tab stop.
func (writer *Writer) Out() {
	if len(writer.indent) > 0 {
		writer.indent = writer.indent[1:]
	}
}

func (writer *Writer) String() string {
	return writer.Content.String()
}
