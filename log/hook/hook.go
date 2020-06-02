package hook

import (
	"fmt"
	"github.com/laconiz/metis/log"
	"github.com/laconiz/metis/log/formatter"
	"os"
)

func NewHook(formatter formatter.Formatter, writer *Writer) *Hook {
	return &Hook{level: writer.level, formatter: formatter, writers: []*Writer{writer}}
}

type Hook struct {
	level     log.Level
	formatter formatter.Formatter
	writers   []*Writer
}

func (hook *Hook) Hook(log *log.Log) {

	raw, err := hook.formatter.Format(log)
	if err != nil {
		const format = "format log[%+v] error: %v"
		str := fmt.Sprintf(format, log, err)
		os.Stderr.WriteString(str)
		return
	}

	for _, writer := range hook.writers {
		writer.Write(log.Level, raw)
	}
}

func (hook *Hook) Writer(writer *Writer) *Hook {

	if writer != nil && writer.level.Valid() {

		if !hook.level.Enable(writer.level) {
			hook.level = writer.level
		}

		hook.writers = append(hook.writers, writer)
	}

	return hook
}

func (hook *Hook) Strap() *Strap {
	return NewStrap(hook)
}
