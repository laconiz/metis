package formatter

import (
	"bytes"
	"fmt"
	"github.com/laconiz/metis/log"
)

func Text() *TextFormatter {
	return (&TextFormatter{}).TimeLayout(DefaultTimeLayout)
}

type TextFormatter struct {
	timeLayout string
}

func (formatter *TextFormatter) TimeLayout(layout string) *TextFormatter {
	return &TextFormatter{timeLayout: layout}
}

func (formatter *TextFormatter) Format(log *log.Log) ([]byte, error) {

	var buf bytes.Buffer

	buf.WriteString(formatter.Level(log.Level))

	buf.WriteByte(' ')
	buf.WriteString(log.Time.Format(formatter.timeLayout))

	if context := log.Context.Raw(); len(context) > 0 {
		buf.WriteString(" context:")
		buf.Write(context)
	}

	if data := log.Data.Raw(); len(data) > 0 {
		buf.WriteString(" data:")
		buf.Write(data)
	}

	buf.WriteString(" $ ")
	buf.WriteString(log.Message)
	if len(log.Message) == 0 || log.Message[len(log.Message)-1] != '\n' {
		buf.WriteByte('\n')
	}

	return buf.Bytes(), nil
}

func (formatter *TextFormatter) Level(level log.Level) string {
	switch level {
	case log.DEBUG:
		return "[DEBUG]"
	case log.INFO:
		return "[INFO] "
	case log.WARN:
		return "[WARN] "
	case log.ERROR:
		return "[ERROR]"
	case log.FATAL:
		return "[FATAL]"
	default:
		const format = "[UNKNOWN<%d>]"
		return fmt.Sprintf(format, level)
	}
}
