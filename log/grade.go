package log

const (
	GradeDebug Grade = "debug"
	GradeInfo  Grade = "info"
	GradeWarn  Grade = "warn"
	GradeError Grade = "error"
	GradeFatal Grade = "fatal"
)

type Grade string

func (grade Grade) Level() Level {
	switch grade {
	case GradeDebug:
		return DEBUG
	case GradeInfo:
		return INFO
	case GradeWarn:
		return WARN
	case GradeError:
		return ERROR
	case GradeFatal:
		return FATAL
	default:
		return INVALID
	}
}
