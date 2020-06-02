package connector

import (
	"github.com/laconiz/metis/nets/session"
	"time"
)

type Dialer func(addr string) (session.Conn, error)

type Option struct {
	Name      string
	Reconnect bool
	Delays    []time.Duration
	Session   session.Option
}

func (opt *Option) parse() {

	if opt.Name == "" {
		opt.Name = "connector"
	}

	if opt.Reconnect && len(opt.Delays) == 0 {
		opt.Delays = []time.Duration{
			time.Millisecond * 100,
			time.Second,
			time.Second * 5,
			time.Second * 15,
		}
	}

	opt.Session.Parse()
}
