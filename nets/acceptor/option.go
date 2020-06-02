package acceptor

import "github.com/laconiz/metis/nets/session"

type Option struct {
	Name    string
	Session session.Option
}

func (opt *Option) parse() {

	if opt.Name == "" {
		opt.Name = "acceptor"
	}

	opt.Session.Parse()
}
