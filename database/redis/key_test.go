package redis

import (
	"testing"
	"time"
)

func TestKey_Exist(t *testing.T) {

	key := client.Key()

	assert(key.Delete(KeyA) == nil)

	exist, err := key.Exist(KeyA)
	assert(err == nil, !exist)

	assert(key.Set(KeyA, ValueA) == nil)
	exist, err = key.Exist(KeyA)
	assert(err == nil, exist)

	assert(key.Delete(KeyA) == nil)
}

func TestKey_Set(t *testing.T) {

	key := client.Key()

	assert(key.Delete(KeyA) == nil)

	var value string
	ok, err := key.Get(KeyA, &value)
	assert(err == nil, !ok)

	assert(key.Set(KeyA, ValueA) == nil)
	ok, err = key.Get(KeyA, &value)
	assert(err == nil, ok, value == ValueA)

	assert(key.SetEX(KeyA, ValueB, 1) == nil)
	ok, err = key.Get(KeyA, &value)
	assert(err == nil, ok, value == ValueB)
	time.Sleep(time.Millisecond * 1100)
	ok, err = key.Get(KeyA, &value)
	assert(err == nil, !ok)

	ok, err = key.SetNX(KeyA, ValueA)
	assert(err == nil, ok)
	ok, err = key.Get(KeyA, &value)
	assert(err == nil, ok, value == ValueA)

	ok, err = key.SetNX(KeyA, ValueB)
	assert(err == nil, !ok)
	ok, err = key.Get(KeyA, &value)
	assert(err == nil, ok, value == ValueA)

	assert(key.Delete(KeyA) == nil)

	ok, err = key.SetNEX(KeyA, ValueB, 1)
	assert(err == nil, ok)
	ok, err = key.Get(KeyA, &value)
	assert(err == nil, ok, value == ValueB)
	ok, err = key.SetNEX(KeyA, ValueA, 1)
	assert(err == nil, !ok)
	ok, err = key.Get(KeyA, &value)
	assert(err == nil, ok, value == ValueB)
	time.Sleep(time.Millisecond * 1100)
	ok, err = key.Get(KeyA, &value)
	assert(err == nil, !ok)

	assert(key.Delete(KeyA) == nil)
}

func TestKey_Incr(t *testing.T) {

	key := client.Key()

	assert(key.Delete(KeyA) == nil)

	value, err := key.Incr(KeyA, 10)
	assert(err == nil, value == 10)
	value, err = key.Incr(KeyA, -5)
	assert(err == nil, value == 5)

	assert(key.Delete(KeyA) == nil)

	assert(key.Set(KeyA, ValueA) == nil)
	value, err = key.Incr(KeyA, 2)
	assert(err != nil)

	assert(key.Delete(KeyA) == nil)
}
