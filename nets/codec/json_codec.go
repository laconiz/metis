package codec

import "github.com/laconiz/metis/utils/json"

type jsonCodec struct {
}

func (codec *jsonCodec) Encode(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}

func (codec *jsonCodec) Decode(raw []byte, msg interface{}) error {
	return json.Unmarshal(raw, msg)
}

var global = &jsonCodec{}

func Json() Codec {
	return global
}
