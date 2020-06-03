package decoder

import (
	"reflect"
)

func Slice(recv interface{}, replies []interface{}) error {

	typo := reflect.TypeOf(recv).Elem()
	slice := reflect.ValueOf(recv).Elem()
	slice.Set(reflect.MakeSlice(typo, 0, len(replies)))

	for _, reply := range replies {

		value, err := Value(typo.Elem(), reply)
		if err != nil {
			return err
		}

		slice.Set(reflect.Append(slice, value))
	}

	return nil
}
