package redis

import (
	"testing"
)

func TestZOrder_Incr(t *testing.T) {

	assert(client.Key().Delete(KeyA) == nil)

	order := client.ZOrder(KeyA)

	Int64, err := order.Incr(MemberA, 4)
	assert(err == nil && Int64 == 4)
	Int64, err = order.Incr(MemberA, -1)
	assert(err == nil && Int64 == 3)

	assert(client.Key().Delete(KeyA) == nil)

	client.Key().Set(KeyA, ValueA)

	_, err = order.Incr(MemberA, 2)
	assert(err != nil)
	t.Log(err)

	assert(client.Key().Delete(KeyA) == nil)
}

func TestZOrder_Range(t *testing.T) {

	assert(client.Key().Delete(KeyA) == nil)

	order := client.ZOrder(KeyA)

	client.Key().Set(KeyA, ValueA)

	scores := map[string]int64{}
	err := order.Range(&scores, 0, 1)
	assert(err != nil)

	assert(client.Key().Delete(KeyA) == nil)

	_, err = order.Incr(MemberA, 2)
	assert(err == nil)
	_, err = order.Incr(MemberB, 4)
	assert(err == nil)

	scores = map[string]int64{}
	err = order.Range(&scores, 0, 1)
	assert(err == nil, len(scores) == 2, scores[MemberA] == 2, scores[MemberB] == 4)

	assert(client.Key().Delete(KeyA) == nil)
}

func TestZOrder_Score(t *testing.T) {

	assert(client.Key().Delete(KeyA) == nil)

	order := client.ZOrder(KeyA)

	client.Key().Set(KeyA, ValueA)

	score, exist, err := order.Score(MemberA)
	assert(err != nil)

	assert(client.Key().Delete(KeyA) == nil)

	_, err = order.Incr(MemberA, 2)
	assert(err == nil)

	score, exist, err = order.Score(MemberB)
	assert(err == nil, !exist)

	score, exist, err = order.Score(MemberA)
	assert(err == nil, exist, score == 2)

	assert(client.Key().Delete(KeyA) == nil)
}

func TestZOrder_Rank(t *testing.T) {

	assert(client.Key().Delete(KeyA) == nil)

	order := client.ZOrder(KeyA)

	client.Key().Set(KeyA, ValueA)

	score, exist, err := order.Rank(MemberA)
	assert(err != nil)

	assert(client.Key().Delete(KeyA) == nil)

	_, err = order.Incr(MemberA, 2)
	assert(err == nil)

	score, exist, err = order.Rank(MemberB)
	assert(err == nil, !exist)

	_, err = order.Incr(MemberB, 4)

	score, exist, err = order.Rank(MemberA)
	assert(err == nil, exist, score == 1)

	assert(client.Key().Delete(KeyA) == nil)
}
