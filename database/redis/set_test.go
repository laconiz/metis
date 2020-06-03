package redis

import (
	"testing"
)

func TestSet(t *testing.T) {

	assert(client.Key().Delete(KeyA) == nil)

	set := client.Set(KeyA)

	client.Key().Set(KeyA, ValueA)

	assert(set.Add(MemberA) != nil)
	assert(set.Remove(MemberA) != nil)

	var keys []string
	assert(set.Keys(&keys) != nil)

	assert(client.Key().Delete(KeyA) == nil)

	assert(set.Add(MemberA, MemberB) == nil)
	assert(set.Keys(&keys) == nil, len(keys) == 2)
	assert(set.Remove(MemberA) == nil)
	assert(set.Keys(&keys) == nil, len(keys) == 1, keys[0] == MemberB)

	assert(client.Key().Delete(KeyA) == nil)

	assert(set.Keys(&keys) == nil, len(keys) == 0)
}
