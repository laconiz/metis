package decoder

import (
	"fmt"
	"reflect"
)

func Struct(recv interface{}, replies []interface{}) error {

	typo := reflect.TypeOf(recv).Elem()

	if len(replies) != typo.NumField() {
		format := "replies's length does not equal struct's field num: %d != %d"
		return fmt.Errorf(format, len(replies), typo.NumField())
	}

	structValue := reflect.ValueOf(recv).Elem()

	for i := 0; i < typo.NumField(); i++ {

		value, err := Value(typo.Field(i).Type, replies[i])
		if err != nil {
			return err
		}

		structValue.Field(i).Set(value)
	}

	return nil
}
