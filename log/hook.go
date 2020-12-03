package log

import (
	"fmt"
	"os"
)

func NewHook(formatter Formatter) *Hook {
	return &Hook{formatter: formatter}
}

type Hook struct {
	formatter Formatter
	writers   []Writer
}

// 调用接口
func (hook *Hook) Hook(log *Log) {
	// 检测日志等级T
	if !hook.Enable(log.Level) {
		return
	}
	// 序列化日志
	raw, err := hook.formatter.Format(log)
	if err != nil {
		const format = "format log[%+v] error: %v"
		str := fmt.Sprintf(format, log, err)
		os.Stderr.WriteString(str)
		return
	}
	// 写入日志
	for _, writer := range hook.writers {
		if writer.Enable(log.Level) {
			writer.Write(log, raw)
		}
	}
}

// 设置日志等级
func (hook *Hook) Enable(level Level) bool {
	for _, writer := range hook.writers {
		if writer.Enable(level) {
			return true
		}
	}
	return false
}

// 设置写入器
func (hook *Hook) Writer(writer Writer) *Hook {
	hook.writers = append(hook.writers, writer)
	return hook
}
