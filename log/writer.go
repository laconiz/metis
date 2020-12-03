package log

import (
	"fmt"
	"io"
	"os"
)

type Writer interface {
	Enable(Level) bool
	Write(log *Log, raw []byte)
}

func NewWriter(level Level, writer io.Writer) Writer {
	return &stdWriter{level: level, writer: writer}
}

type stdWriter struct {
	level  Level
	writer io.Writer
}

// 检测日志等级
func (writer *stdWriter) Enable(level Level) bool {
	return writer.level.Enable(level)
}

// 写入日志
func (writer *stdWriter) Write(_ *Log, raw []byte) {
	if _, err := writer.writer.Write(raw); err != nil {
		const format = "write log[%s] error: %v"
		str := fmt.Sprintf(format, string(raw), err)
		os.Stderr.WriteString(str)
	}
}
