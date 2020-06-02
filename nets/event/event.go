package event

import (
	"github.com/laconiz/metis/nets/codec"
	"github.com/laconiz/metis/nets/packet"
	"sync"
)

type Session interface {
	ID() uint64
	Addr() string
	Data() *sync.Map
	Send(interface{}) error
	SendRaw([]byte) error
	Close() error
}

type Event struct {
	*packet.Packet
	Ses Session
}

type Connected struct{}

type Disconnected struct{}

type ConnectFailed struct{}

func NewConnected(ses Session) *Event {
	return &Event{Ses: ses, Packet: pktConnected}
}

func NewDisconnected(ses Session) *Event {
	return &Event{Ses: ses, Packet: pktDisconnected}
}

func NewConnectFailed() *Event {
	return &Event{Packet: pktConnectFailed}
}

var (
	pktConnected     *packet.Packet
	pktDisconnected  *packet.Packet
	pktConnectFailed *packet.Packet
)

func init() {

	meta, err := packet.Register(Connected{}, codec.Json())
	if err != nil {
		panic(err)
	}
	pktConnected = &packet.Packet{Meta: meta, Msg: &Connected{}}

	meta, err = packet.Register(Disconnected{}, codec.Json())
	if err != nil {
		panic(err)
	}
	pktDisconnected = &packet.Packet{Meta: meta, Msg: &Disconnected{}}

	meta, err = packet.Register(ConnectFailed{}, codec.Json())
	if err != nil {
		panic(err)
	}
	pktConnectFailed = &packet.Packet{Meta: meta, Msg: &ConnectFailed{}}
}
