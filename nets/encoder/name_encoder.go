package encoder

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/laconiz/metis/nets/packet"
)

const Sep = '&'

type NameEncoder struct {
}

func (encoder *NameEncoder) Marshal(msg interface{}) (*packet.Packet, error) {

	meta := packet.MetaByMsg(msg)
	if meta == nil {
		return nil, errors.New("meta cannot be found")
	}

	raw, err := meta.Encode(msg)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBufferString(meta.Name())
	buf.WriteByte(Sep)
	buf.Write(raw)

	return &packet.Packet{Meta: meta, Msg: msg, Raw: raw, Stream: buf.Bytes()}, nil
}

func (encoder *NameEncoder) Unmarshal(stream []byte) (*packet.Packet, error) {

	args := bytes.SplitN(stream, []byte{Sep}, 2)
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid stream: %s", string(stream))
	}

	meta := packet.MetaByName(string(args[0]))
	if meta == nil {
		return nil, errors.New("meta cannot be found")
	}

	msg, err := meta.Decode(args[1])
	if err != nil {
		return nil, err
	}

	return &packet.Packet{Meta: meta, Msg: msg, Raw: args[1], Stream: stream}, nil
}

type NameMaker struct {
	encoder NameEncoder
}

func (maker NameMaker) New() Encoder {
	return &maker.encoder
}
