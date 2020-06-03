package decoder

import "reflect"

func Value(typo reflect.Type, reply interface{}) (reflect.Value, error) {

	ptr := false
	if typo.Kind() == reflect.Ptr {
		ptr = true
		typo = typo.Elem()
	}

	value := reflect.New(typo)
	err := Decode(value.Interface(), reply, nil)
	if !ptr {
		value = value.Elem()
	}

	return value, err
}
