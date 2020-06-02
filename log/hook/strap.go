package hook

import (
	"github.com/laconiz/metis/log"
)

func NewStrap(hook *Hook) *Strap {
	return &Strap{level: hook.level, hooks: []*Hook{hook}}
}

type Strap struct {
	level log.Level
	hooks []*Hook
}

func (strap *Strap) Enable(level log.Level) bool {
	return strap.level.Enable(level)
}

func (strap *Strap) Hook(hook *Hook) *Strap {

	if !strap.level.Enable(hook.level) {
		strap.level = hook.level
	}

	strap.hooks = append(strap.hooks, hook)
	return strap
}

func (strap *Strap) Invoke(log *log.Log) {
	for _, hook := range strap.hooks {
		if hook.level.Enable(log.Level) {
			hook.Hook(log)
		}
	}
}
