package encoder

import "github.com/laconiz/metis/nets/packet"

type Encoder interface {
	Marshal(interface{}) (*packet.Packet, error)
	Unmarshal([]byte) (*packet.Packet, error)
}

type Maker interface {
	New() Encoder
}
