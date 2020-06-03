package decoder

import (
	"fmt"
	"github.com/laconiz/metis/utils/json"
)

func Params(args []interface{}) ([]interface{}, error) {

	result := make([]interface{}, len(args))

	for idx, arg := range args {

		switch arg.(type) {
		case string:
			result[idx] = arg
		case int8, int16, int32, int64, int,
			uint8, uint16, uint32, uint64, uint,
			float32, float64, complex64, complex128, bool:
			result[idx] = arg
		case nil:
			return nil, fmt.Errorf("got nil arguments: %v", args)
		default:
			raw, err := json.Marshal(arg)
			if err != nil {
				return nil, err
			}
			result[idx] = string(raw)
		}
	}

	return result, nil
}

func Reply(reply interface{}) interface{} {

	switch reply.(type) {
	case []byte:
		return string(reply.([]byte))
	case []interface{}:
		var formats []interface{}
		for _, value := range reply.([]interface{}) {
			formats = append(formats, Reply(value))
		}
		return formats
	}

	return reply
}
