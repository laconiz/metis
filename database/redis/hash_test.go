package redis

import (
	"testing"
)

func TestHash_Set(t *testing.T) {

	assert(client.Key().Delete(KeyA) == nil)

	hash := client.Hash(KeyA)

	assert(hash.Delete(MemberA) == nil)

	assert(hash.Set(MemberA, ValueA) == nil)
	var value string
	ok, err := hash.Get(MemberA, &value)
	assert(err == nil, ok, value == ValueA)

	assert(hash.Delete(MemberA) == nil)
	ok, err = hash.Get(MemberA, &value)
	assert(err == nil, !ok, value != ValueA)

	assert(client.Key().Delete(KeyA) == nil)
}

func TestHash_Get(t *testing.T) {

	assert(client.Key().Delete(KeyA) == nil)

	hash := client.Hash(KeyA)

	var value string
	ok, err := hash.Get(MemberA, &value)
	assert(err == nil, !ok)

	assert(client.Key().Delete(KeyA) == nil)

	hash.Set(MemberA, ValueA)
	ok, err = hash.Get(MemberA, &value)
	assert(err == nil, ok, value == ValueA)

	assert(client.Key().Delete(KeyA) == nil)
}

func TestHash_Gets(t *testing.T) {

	assert(client.Key().Delete(KeyA) == nil)

	assert(client.Key().Set(MemberA, ValueA) == nil)

	hash := client.Hash(KeyA)

	var value []string
	assert(hash.Gets(&value, MemberA) == nil)

	assert(client.Key().Delete(KeyA) == nil)

	assert(hash.Set(MemberA, ValueA) == nil)
	assert(hash.Set(MemberB, ValueB) == nil)

	assert(hash.Gets(&value, MemberA, MemberB) == nil, len(value) == 2)
	assert(hash.Gets(&value, MemberA) == nil, len(value) == 1)

	assert(client.Key().Delete(KeyA) == nil)
}

func TestHash_GetAll(t *testing.T) {

	assert(client.Key().Delete(KeyA) == nil)

	assert(client.Key().Set(MemberA, ValueA) == nil)

	hash := client.Hash(KeyA)

	var value map[string]string

	assert(hash.Set(MemberA, ValueA) == nil)
	assert(hash.GetAll(&value) == nil)
	assert(len(value) == 1, value[MemberA] == ValueA)

	assert(hash.Set(MemberB, ValueB) == nil)
	assert(hash.GetAll(&value) == nil)
	assert(len(value) == 2, value[MemberA] == ValueA, value[MemberB] == ValueB)

	assert(client.Key().Delete(KeyA) == nil)
}

func TestHash_Increase(t *testing.T) {

	assert(client.Key().Delete(KeyA) == nil)

	hash := client.Hash(KeyA)

	assert(client.Key().Set(KeyA, ValueA) == nil)

	value, err := hash.Increase(MemberA, 2)
	assert(err != nil)

	assert(client.Key().Delete(KeyA) == nil)

	value, err = hash.Increase(MemberA, 3)
	assert(err == nil && value == 3)

	assert(client.Key().Delete(KeyA) == nil)
}

func TestHash_Consume(t *testing.T) {

	assert(client.Key().Delete(KeyA) == nil)

	hash := client.Hash(KeyA)

	assert(client.Key().Set(KeyA, ValueA) == nil)

	value, success, err := hash.Consume(MemberA, -2)
	assert(err != nil)

	assert(client.Key().Delete(KeyA) == nil)

	value, success, err = hash.Consume(MemberA, 10)
	assert(err == nil, success, value == 10)

	value, success, err = hash.Consume(MemberA, -6)
	assert(err == nil, success, value == 4)

	value, success, err = hash.Consume(MemberA, -5)
	assert(err == nil, !success, value == 4)

	assert(client.Key().Delete(KeyA) == nil)
}
