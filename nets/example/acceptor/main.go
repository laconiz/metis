package main

import (
	"github.com/laconiz/metis/log"
	"github.com/laconiz/metis/log/logz"
	"github.com/laconiz/metis/nets/acceptor"
	"github.com/laconiz/metis/nets/event"
	"github.com/laconiz/metis/nets/example"
	"github.com/laconiz/metis/nets/session"
	"github.com/laconiz/metis/nets/socket"
	"time"
)

var times uint64

func main() {

	logger := logz.Level(log.INFO).Field("module", "main")

	invoker := event.NewStdInvoker()
	invoker.Register(example.REQ{}, func(event *event.Event) {
		req := event.Msg.(*example.REQ)
		event.Ses.Send(&example.ACK{Time: req.Time})
		times++
	})

	option := acceptor.Option{Session: session.Option{Invoker: invoker}}

	dialer := socket.Acceptor(102400)

	acceptor := acceptor.New(example.Addr, dialer, option, logger)
	acceptor.Run()

	for {
		<-time.After(time.Second)
		if times > 0 {
			logz.Infof("%d: %d", acceptor.Count(), times)
			times = 0
		}
	}
}
