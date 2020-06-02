package main

import (
	"github.com/laconiz/metis/log"
	"github.com/laconiz/metis/log/logz"
	"github.com/laconiz/metis/nets/connector"
	"github.com/laconiz/metis/nets/event"
	"github.com/laconiz/metis/nets/example"
	"github.com/laconiz/metis/nets/session"
	"github.com/laconiz/metis/nets/socket"
	"math/rand"
	"time"
)

var times time.Duration
var duration time.Duration

func NewConnector(logger log.Logger) *connector.Connector {

	var bytes []byte
	for i := 0; i < 40960; i++ {
		bytes = append(bytes, byte(i))
	}

	random := rand.New(rand.NewSource(time.Now().Unix()))

	invoker := event.NewStdInvoker()
	invoker.Register(example.ACK{}, func(event *event.Event) {
		ack := event.Msg.(*example.ACK)
		times++
		duration += time.Since(ack.Time)
	})

	option := connector.Option{Reconnect: true, Session: session.Option{Queue: 16, Invoker: invoker}}

	dialer := socket.Connector()

	connector := connector.New(example.Addr, dialer, option, logger)
	connector.Run()

	go func() {
		for {
			connector.Send(&example.REQ{Time: time.Now(), Bytes: bytes[:random.Intn(10)+1]})
		}
	}()

	return connector
}

func main() {

	logger := logz.Level(log.INFO).Field("module", "main")

	for i := 0; i < 2; i++ {
		NewConnector(logger)
	}

	for {
		<-time.After(time.Second)
		if times > 0 {
			logger.Infof("%v / %d = %v", duration, times, duration/times)
			times = 0
			duration = 0
		}
	}
}
