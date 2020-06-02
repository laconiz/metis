package encoder

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/laconiz/metis/nets/packet"
)

type IDEncoder struct {
}

func (encoder *IDEncoder) Marshal(msg interface{}) (*packet.Packet, error) {

	meta := packet.MetaByMsg(msg)
	if meta == nil {
		return nil, errors.New("meta cannot be found")
	}

	raw, err := meta.Encode(msg)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, meta.ID())
	buf.Write(raw)

	return &packet.Packet{Meta: meta, Msg: msg, Raw: raw, Stream: buf.Bytes()}, nil
}

func (encoder *IDEncoder) Unmarshal(stream []byte) (*packet.Packet, error) {

	buf := bytes.NewBuffer(stream)

	var id packet.MetaID
	if err := binary.Read(buf, binary.LittleEndian, &id); err != nil {
		return nil, err
	}

	meta := packet.MetaByID(id)
	if meta == nil {
		return nil, errors.New("meta cannot be found")
	}

	raw := buf.Bytes()

	msg, err := meta.Decode(raw)
	if err != nil {
		return nil, err
	}

	return &packet.Packet{Meta: meta, Msg: msg, Raw: raw, Stream: stream}, nil
}

type IDMaker struct {
	encoder IDEncoder
}

func (maker IDMaker) New() Encoder {
	return &maker.encoder
}
