package hook

import (
	"fmt"
	"github.com/laconiz/metis/log"
	"io"
	"os"
)

func NewWriter(level log.Level, writer io.Writer) *Writer {
	return &Writer{level: level, writer: writer}
}

type Writer struct {
	level  log.Level
	writer io.Writer
}

func (writer *Writer) Write(level log.Level, raw []byte) {

	if !writer.level.Enable(level) {
		return
	}

	if _, err := writer.writer.Write(raw); err != nil {
		const format = "write log[%s] error: %v"
		str := fmt.Sprintf(format, string(raw), err)
		os.Stderr.WriteString(str)
	}
}
