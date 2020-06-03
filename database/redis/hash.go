package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/laconiz/metis/database/redis/decoder"
)

type Hash struct {
	key    string
	client *Redis
}

/*
	HDEL key field [field ...]

	删除哈希表 key 中的一个或多个指定域, 不存在的域将被忽略.

	可用版本: >= 2.0.0
	时间复杂度: O(N), N 为要删除的域的数量.

	返回值:
		被成功移除的域的数量, 不包括被忽略的域.
*/
func (hash *Hash) Delete(fields ...interface{}) error {
	_, err := hash.client.Exec(HDEL, append([]interface{}{hash.key}, fields...)...)
	return err
}

/*
	HSET key field value

	将哈希表 key 中的域 field 的值设为 value .
	如果 key 不存在, 一个新的哈希表被创建并进行 HSET 操作.
	如果域 field 已经存在于哈希表中, 旧值将被覆盖.

	可用版本: >= 2.0.0
	时间复杂度: O(1)

	返回值:
		如果 field 是哈希表中的一个新建域, 并且值设置成功, 返回 1 .
		如果哈希表中域 field 已经存在且旧值已被新值覆盖, 返回 0 .
*/
func (hash *Hash) Set(field, value interface{}) error {
	_, err := hash.client.Exec(HSET, hash.key, field, value)
	return err
}

/*
	HGET key field

	返回哈希表 key 中给定域 field 的值。

	可用版本: >= 2.0.0
	时间复杂度: O(1)

	返回值:
		给定域的值.
		当给定域不存在或是给定 key 不存在时, 返回 nil .
*/
func (hash *Hash) Get(field, value interface{}) error {
	reply, err := hash.client.Exec(HGET, hash.key, field)
	return decoder.Decode(value, reply, err)
}

/*
	HMGET key field [field ...]

	返回哈希表 key 中, 一个或多个给定域的值.
	如果给定的域不存在于哈希表, 那么返回一个 nil 值.
	因为不存在的 key 被当作一个空哈希表来处理, 所以对一个不存在的 key 进行 HMGET 操作将返回一个只带有 nil 值的表.

	可用版本: >= 2.0.0
	时间复杂度: O(N), N 为给定域的数量.

	返回值:
		一个包含多个给定域的关联值的表, 表值的排列顺序和给定域参数的请求顺序一样.
*/
func (hash *Hash) Gets(value interface{}, fields ...interface{}) error {
	reply, err := hash.client.Exec(HMGET, append([]interface{}{hash.key}, fields...)...)
	return decoder.Decode(value, reply, err)
}

/*
	HGETALL key

	返回哈希表 key 中, 所有的域和值.
	在返回值里, 紧跟每个域名(field name)之后是域的值(value), 所以返回值的长度是哈希表大小的两倍.

	可用版本: >= 2.0.0
	时间复杂度: O(N), N 为哈希表的大小.

	返回值:
		以列表形式返回哈希表的域和域的值.
		若 key 不存在, 返回空列表.
*/
func (hash *Hash) GetAll(value interface{}) error {
	reply, err := hash.client.Exec(HGETALL, hash.key)
	return decoder.Decode(value, reply, err)
}

/*
	HINCRBY key field increment

	为哈希表 key 中的域 field 的值加上增量 increment .
	增量也可以为负数, 相当于对给定域进行减法操作.
	如果 key 不存在, 一个新的哈希表被创建并执行 HINCRBY 命令.
	如果域 field 不存在, 那么在执行命令前, 域的值被初始化为 0 .
	对一个储存字符串值的域 field 执行 HINCRBY 命令将造成一个错误.
	本操作的值被限制在 64 位(bit)有符号数字表示之内.

	可用版本: >= 2.0.0
	时间复杂度: O(1)
	返回值:
		执行 HINCRBY 命令之后, 哈希表 key 中域 field 的值.
*/
func (hash *Hash) Increase(field interface{}, increment int64) (int64, error) {
	reply, err := hash.client.Exec(HINCRBY, hash.key, field, increment)
	var value int64
	err = decoder.Decode(&value, reply, err)
	return value, err
}

func (hash *Hash) Consume(field interface{}, increment int64) (int64, bool, error) {

	if increment >= 0 {
		value, err := hash.Increase(field, increment)
		return value, err == nil, err
	}

	type consumeResult struct {
		Value   int64
		Success bool
	}

	result := &consumeResult{}
	eval := hash.client.Eval(scriptHashConsume)
	reply, err := eval.Exec(hash.key, field, increment)
	err = decoder.Decode(result, reply, err)
	return result.Value, result.Success, err
}

var scriptHashConsume = &Script{Name: "HashConsume", Script: redis.NewScript(1, luaConsume)}

var luaConsume = `

	local new = redis.call('HINCRBY', KEYS[1], ARGV[1], ARGV[2])
	if new >= 0 then
		return {new, 1}
	end

	new = redis.call('HINCRBY', KEYS[1], ARGV[1], -ARGV[2])
	return {new, 0}
`
