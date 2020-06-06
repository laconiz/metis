package redis

import (
	"github.com/laconiz/metis/database/redis/decoder"
)

type Key struct {
	client *Client
}

/*
	EXISTS key

	检查给定 key 是否存在.

	可用版本: >= 1.0.0
	时间复杂度: O(1)

	返回值:
		若 key 存在, 返回 1 , 否则返回 0 .
*/
func (k *Key) Exist(key interface{}) (bool, error) {
	reply, err := k.client.Exec(EXISTS, key)
	var exist bool
	err = decoder.Decode(&exist, reply, err)
	return exist, err
}

/*
	SET key value [EX seconds] [PX milliseconds] [NX|XX]

	将字符串值 value 关联到 key .
	如果 key 已经持有其他值, SET 就覆写旧值, 无视类型.
	对于某个原本带有生存时间(TTL)的键来说, 当 SET 命令成功在这个键上执行时, 这个键原有的 TTL 将被清除.

	可选参数:
		EX second: 设置键的过期时间为 second 秒.
		PX millisecond: 设置键的过期时间为 millisecond 毫秒.
		NX: 只在键不存在时, 才对键进行设置操作.
		XX: 只在键已经存在时, 才对键进行设置操作.

	可用版本: >= 2.6.12
	时间复杂度: O(1)

	返回值：
		SET 在设置操作成功完成时, 才返回 OK .
		如果设置了 NX 或者 XX , 但因为条件没达到而造成设置操作未执行, 那么命令返回空批量回复(NULL Bulk Reply).
*/
func (k *Key) Set(key, value interface{}) error {
	_, err := k.client.Exec(SET, key, value)
	return err
}

func (k *Key) SetNX(key, value interface{}) (bool, error) {
	reply, err := k.client.Exec(SET, key, value, NX)
	var ok string
	err = decoder.Decode(&ok, reply, err)
	return ok == OK, err
}

func (k *Key) SetEX(key, value interface{}, second int64) error {
	_, err := k.client.Exec(SET, key, value, EX, second)
	return err
}

func (k *Key) SetNEX(key, value interface{}, second int64) (bool, error) {
	reply, err := k.client.Exec(SET, key, value, EX, second, NX)
	var ok string
	err = decoder.Decode(&ok, reply, err)
	return ok == OK, err
}

/*
	GET key

	返回 key 所关联的字符串值.
	如果 key 不存在那么返回特殊值 nil .
	假如 key 储存的值不是字符串类型, 返回一个错误, 因为 GET 只能用于处理字符串值.

	可用版本: >= 1.0.0
	时间复杂度: O(1)
	返回值:
		当 key 不存在时, 返回 nil , 否则，返回 key 的值.
		如果 key 不是字符串类型, 那么返回一个错误.
*/
func (k *Key) Get(key, value interface{}) (bool, error) {
	reply, err := k.client.Exec(GET, key)
	err = decoder.Decode(value, reply, err)
	return reply != nil && err == nil, err
}

/*
	DEL key [key ...]

	删除给定的一个或多个 key .
	不存在的 key 会被忽略.

	可用版本: >= 1.0.0
	时间复杂度:
		O(N), N 为被删除的 key 的数量.
		删除单个字符串类型的 key , 时间复杂度为O(1).
		删除单个列表、集合、有序集合或哈希表类型的 key , 时间复杂度为O(M) , M 为以上数据结构内的元素数量.

	返回值:
		被删除 key 的数量.
*/
func (k *Key) Delete(keys ...interface{}) error {
	_, err := k.client.Exec(DEL, keys...)
	return err
}

/*
	INCRBY key increment

	将 key 所储存的值加上增量 increment .
	如果 key 不存在, 那么 key 的值会先被初始化为 0 , 然后再执行 INCRBY 命令.
	如果值包含错误的类型, 或字符串类型的值不能表示为数字, 那么返回一个错误.
	本操作的值限制在 64 位(bit)有符号数字表示之内.
	关于递增(increment) / 递减(decrement)操作的更多信息, 参见 INCR 命令.

	可用版本: >= 1.0.0
	时间复杂度: O(1)

	返回值:
		加上 increment 之后 key 的值.
*/
func (k *Key) Incr(key interface{}, increment int64) (int64, error) {
	reply, err := k.client.Exec(INCRBY, key, increment)
	var value int64
	err = decoder.Decode(&value, reply, err)
	return value, err
}
