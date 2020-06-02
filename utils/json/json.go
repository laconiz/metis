package json

import (
	"github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

type RawMessage = jsoniter.RawMessage

var json jsoniter.API

func Marshal(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func Unmarshal(raw []byte, value interface{}) error {
	return json.Unmarshal(raw, value)
}

func MarshalText(v interface{}) (string, error) {

	raw, err := Marshal(v)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}

func UnmarshalText(str string, value interface{}) error {
	return Unmarshal([]byte(str), value)
}

func init() {

	extra.RegisterFuzzyDecoders()

	json = jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		UseNumber:              true,
	}.Froze()
}
