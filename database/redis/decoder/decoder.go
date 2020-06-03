package decoder

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/laconiz/metis/utils/json"
	"go/types"
	"reflect"
)

func Decode(recv interface{}, reply interface{}, err error) error {

	if err != nil {
		return err
	}

	switch recv.(type) {
	case types.Nil:
		return errors.New("nil receiver")
	case *string:
		return String(recv, reply)
	case *[]byte:
		return Bytes(recv, reply)
	case *bool:
		return Bool(recv, reply)
	case *int8, *int16, *int32, *int64, *int:
		return Int(recv, reply)
	case *uint8, *uint16, *uint32, *uint64, *uint:
		return Uint(recv, reply)
	}

	typo := reflect.TypeOf(recv)
	if typo.Kind() != reflect.Ptr {
		return fmt.Errorf("non-pointer receiver: %v", typo)
	}

	if bytes, err := redis.Bytes(reply, nil); err == nil {
		return json.Unmarshal(bytes, recv)
	}

	replies, err := redis.Values(reply, nil)
	if err != nil {
		return err
	}

	switch typo.Elem().Kind() {
	case reflect.Slice:
		return Slice(recv, replies)
	case reflect.Map:
		return Map(recv, replies)
	case reflect.Struct:
		return Struct(recv, replies)
	}

	return fmt.Errorf("unsupported type: %v", typo)
}

func String(recv interface{}, reply interface{}) error {

	value, err := redis.String(reply, nil)
	if err == redis.ErrNil {
		err = nil
	}

	reflect.ValueOf(recv).Elem().SetString(value)
	return err
}

func Bytes(recv interface{}, reply interface{}) error {

	value, err := redis.Bytes(reply, nil)
	if err == redis.ErrNil {
		err = nil
	}

	reflect.ValueOf(recv).Elem().SetBytes(value)
	return err
}

func Bool(recv interface{}, reply interface{}) error {

	value, err := redis.Bool(reply, nil)
	if err == redis.ErrNil {
		err = nil
	}

	reflect.ValueOf(recv).Elem().SetBool(value)
	return err
}

func Int(recv interface{}, reply interface{}) error {

	value, err := redis.Int64(reply, nil)
	if err == redis.ErrNil {
		err = nil
	}

	reflect.ValueOf(recv).Elem().SetInt(value)
	return err
}

func Uint(recv interface{}, reply interface{}) error {

	value, err := redis.Uint64(reply, nil)
	if err == redis.ErrNil {
		err = nil
	}

	reflect.ValueOf(recv).Elem().SetUint(value)
	return err
}
