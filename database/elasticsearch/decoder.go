package elasticsearch

import "github.com/laconiz/metis/utils/json"

type Decoder struct {
}

func (decoder *Decoder) Decode(raw []byte, value interface{}) error {
	return json.Unmarshal(raw, value)
}
