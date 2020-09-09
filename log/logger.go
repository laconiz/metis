package log

import (
	"github.com/laconiz/metis/log/context"
)

type Logger interface {
	Level(Level) Logger
	Field(string, interface{}) Logger
	Fields(context.Fields) Logger
	Data(string, interface{}) Logger
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
}

const Module = "module"
