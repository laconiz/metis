package decoder

import (
	"reflect"
)

func Map(recv interface{}, replies []interface{}) error {

	typo := reflect.TypeOf(recv).Elem()
	mapValue := reflect.ValueOf(recv).Elem()
	mapValue.Set(reflect.MakeMap(typo))

	for i := 0; i < len(replies)-1; i += 2 {

		index, err := Value(typo.Key(), replies[i])
		if err != nil {
			return err
		}

		value, err := Value(typo.Elem(), replies[i+1])
		if err != nil {
			return err
		}

		mapValue.SetMapIndex(index, value)
	}

	return nil
}
