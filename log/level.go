package log

import "fmt"

const (
	DEBUG Level = 1 << iota
	INFO
	WARN
	ERROR
	FATAL
	INVALID
)

type Level int8

func (level Level) Valid() bool {
	switch level {
	case DEBUG, INFO, WARN, ERROR, FATAL:
		return true
	default:
		return false
	}
}

func (level Level) Enable(other Level) bool {
	return level <= other
}

func (level Level) Grade() Grade {
	switch level {
	case DEBUG:
		return GradeDebug
	case INFO:
		return GradeInfo
	case WARN:
		return GradeWarn
	case ERROR:
		return GradeError
	case FATAL:
		return GradeFatal
	default:
		return Grade(fmt.Sprintf("unknown[%d]", level))
	}
}
