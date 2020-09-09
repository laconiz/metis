package context

import (
	"fmt"
	"github.com/laconiz/metis/utils/json"
	"reflect"
)

func NewData() *Data {
	return &Data{
		counter: map[string]int32{},
		values:  map[string]interface{}{},
	}
}

type Data struct {
	raw     []byte
	counter map[string]int32
	values  map[string]interface{}
}

func (data *Data) index(key string, index int32) string {
	return fmt.Sprintf("%s#%d", key, index)
}

func (data *Data) value(value interface{}) {

	index := func(key string, count int32) string {
		const format = "%s#%d"
		return fmt.Sprintf(format, key, count)
	}

	typo := reflect.TypeOf(value)
	key := "nil"
	if typo != nil {
		key = typo.Name()
	}

	count := data.counter[key]
	switch count {
	case 0:
		data.values[key] = value
	case 1:
		first := int32(0)
		data.values[index(key, first)] = data.values[key]
		delete(data.values, key)
		data.values[index(key, count)] = value
	default:
		data.values[index(key, count)] = value
	}

	data.counter[key]++
	data.raw = nil
}

func (data *Data) Value(values ...interface{}) *Data {

	copy := NewData()

	for key, value := range data.counter {
		copy.counter[key] = value
	}

	for key, value := range data.values {
		copy.values[key] = value
	}

	for _, value := range values {

		type Error string

		if err, ok := value.(error); !ok {
			copy.value(value)
		} else {
			copy.value(Error(err.Error()))
		}
	}

	return copy
}

func (data *Data) Raw() []byte {

	if len(data.values) == 0 {
		return nil
	}

	if data.raw != nil {
		return data.raw
	}

	raw, err := json.Marshal(data.values)
	if err == nil {
		data.raw = raw
	} else {
		data.raw = []byte(fmt.Sprint(data.values))
	}

	return data.raw
}
