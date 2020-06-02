package example

import (
	"github.com/laconiz/metis/nets/codec"
	"github.com/laconiz/metis/nets/packet"
	"time"
)

const Addr = "127.0.0.1:6000"

type REQ struct {
	Time  time.Time
	Bytes []byte
}

type ACK struct {
	Time time.Time
}

func init() {
	packet.Register(REQ{}, codec.Json())
	packet.Register(ACK{}, codec.Json())
}
