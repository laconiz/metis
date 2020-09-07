package log

import (
	"github.com/laconiz/metis/log/context"
)

type Logger interface {
	Level(level Level) Logger
	Field(key string, value interface{}) Logger
	Fields(fields context.Fields) Logger
	Data(data ...interface{}) Logger
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
}

const Module = "module"
