package log

import (
	"fmt"
	"github.com/laconiz/metis/log/context"
	"os"
	"time"
)

// 创建日志接口
func NewEntry(hooks ...*Hook) Logger {
	return &Entry{
		level:   DEBUG,
		data:    context.New(),
		context: context.New(),
		hooks:   hooks,
	}
}

// 日志接口
type Entry struct {
	level   Level            // 等级
	data    *context.Context // 数据
	context *context.Context //
	hooks   []*Hook          //
}

// 设置日志等级
func (entry *Entry) Level(level Level) Logger {
	copy := *entry
	copy.level = level
	return &copy
}

func (entry *Entry) Field(key string, value interface{}) Logger {
	return entry.Fields(context.Fields{key: value})
}

func (entry *Entry) Fields(fields context.Fields) Logger {
	copy := *entry
	copy.context = entry.context.Fields(fields)
	return &copy
}

func (entry *Entry) Data(key string, value interface{}) Logger {
	copy := *entry
	copy.data = entry.data.Fields(context.Fields{key: value})
	return &copy
}

func (entry *Entry) Debug(args ...interface{}) {
	entry.Log(DEBUG, args...)
}

func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.Logf(DEBUG, format, args...)
}

func (entry *Entry) Info(args ...interface{}) {
	entry.Log(INFO, args...)
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.Logf(INFO, format, args...)
}

func (entry *Entry) Warn(args ...interface{}) {
	entry.Log(WARN, args...)
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.Logf(WARN, format, args...)
}

func (entry *Entry) Error(args ...interface{}) {
	entry.Log(ERROR, args...)
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.Logf(ERROR, format, args...)
}

func (entry *Entry) Fatal(args ...interface{}) {
	entry.Log(FATAL, args...)
	os.Exit(1)
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.Logf(FATAL, format, args...)
	os.Exit(1)
}

func (entry *Entry) Log(level Level, args ...interface{}) {
	if !entry.level.Enable(level) {
		return
	}
	entry.log(level, fmt.Sprint(args...))
}

func (entry *Entry) Logf(level Level, format string, args ...interface{}) {
	if !entry.level.Enable(level) {
		return
	}
	entry.log(level, fmt.Sprintf(format, args...))
}

// 调用接口
func (entry *Entry) log(level Level, message string) {
	log := &Log{Level: level, Time: time.Now(), Data: entry.data, Context: entry.context, Message: message}
	for _, hook := range entry.hooks {
		hook.Hook(log)
	}
}
