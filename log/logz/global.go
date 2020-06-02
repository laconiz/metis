package logz

import (
	"github.com/laconiz/metis/log"
	"github.com/laconiz/metis/log/context"
	"github.com/laconiz/metis/log/formatter"
	"github.com/laconiz/metis/log/hook"
	"os"
)

var global = log.New(hook.NewHook(formatter.Text(), hook.NewWriter(log.DEBUG, os.Stdout)).Strap())

func Global() log.Logger {
	return global
}

func Level(level log.Level) log.Logger {
	return global.Level(level)
}

func Field(key string, value interface{}) log.Logger {
	return global.Field(key, value)
}

func Fields(fields context.Fields) log.Logger {
	return global.Fields(fields)
}

func Data(data ...interface{}) log.Logger {
	return global.Data(data...)
}

func Debug(args ...interface{}) {
	global.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	global.Debugf(format, args...)
}

func Info(args ...interface{}) {
	global.Info(args...)
}

func Infof(format string, args ...interface{}) {
	global.Infof(format, args...)
}

func Warn(args ...interface{}) {
	global.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	global.Warnf(format, args...)
}

func Error(args ...interface{}) {
	global.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	global.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	global.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	global.Fatalf(format, args...)
}
