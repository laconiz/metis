package codec

import (
	"github.com/laconiz/metis/utils/json"
)

type patcher interface {
	Patch()
}

type jsonCodec struct {
}

func (codec *jsonCodec) Encode(msg interface{}) ([]byte, error) {
	if patcher, ok := msg.(patcher); ok {
		patcher.Patch()
	}
	return json.Marshal(msg)
}

func (codec *jsonCodec) Decode(raw []byte, msg interface{}) error {
	return json.Unmarshal(raw, msg)
}

var global = &jsonCodec{}

func Json() Codec {
	return global
}
