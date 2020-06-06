package redis

import "github.com/laconiz/metis/database/redis/decoder"

type Set struct {
	key    string
	client *Client
}

/*
	SADD key member [member ...]

	将一个或多个 member 元素加入到集合 key 当中, 已经存在于集合的 member 元素将被忽略.
	假如 key 不存在, 则创建一个只包含 member 元素作成员的集合.
	当 key 不是集合类型时, 返回一个错误.

	可用版本: >= 1.0.0
	时间复杂度: O(N), N 为被添加的元素的数量.

	返回值:
		被添加到集合中的新元素的数量, 不包括被忽略的元素.
*/
func (set *Set) Add(members ...interface{}) error {
	_, err := set.client.Exec(SADD, append([]interface{}{set.key}, members...)...)
	return err
}

/*
	SREM key member [member ...]

	移除集合 key 中的一个或多个 member 元素, 不存在的 member 元素会被忽略.
	当 key 不是集合类型, 返回一个错误.

	可用版本: >= 1.0.0
	时间复杂度: O(N), N 为给定 member 元素的数量.

	返回值:
		被成功移除的元素的数量, 不包括被忽略的元素.
*/
func (set *Set) Remove(members ...interface{}) error {
	_, err := set.client.Exec(SREM, append([]interface{}{set.key}, members...)...)
	return err
}

/*
	SMEMBERS key

	返回集合 key 中的所有成员.
	不存在的 key 被视为空集合.

	可用版本: >= 1.0.0
	时间复杂度: O(N), N 为集合的基数.

	返回值:
		集合中的所有成员.
*/
func (set *Set) Keys(value interface{}) error {
	reply, err := set.client.Exec(SMEMBERS, set.key)
	return decoder.Decode(value, reply, err)
}
