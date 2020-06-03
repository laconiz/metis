package redis

import (
	"github.com/laconiz/metis/database/redis/decoder"
)

type ZOrder struct {
	key    string
	client *Redis
}

/*
	ZINCRBY key increment member

	为有序集 key 的成员 member 的 score 值加上增量 increment .
	可以通过传递一个负数值 increment 让 score 减去相应的值.
	当 key 不存在或 member 不是 key 的成员时, ZINCRBY key increment member 等同于 ZADD key increment member .
	当 key 不是有序集类型时, 返回一个错误.
	score 值可以是整数值或双精度浮点数.

	可用版本: >= 1.2.0
	时间复杂度: O(log(N))

	返回值:
		member 成员的新 score 值, 以字符串形式表示.
*/

func (order *ZOrder) Incr(member interface{}, increment int64) (int64, error) {
	reply, err := order.client.Exec(ZINCRBY, order.key, increment, member)
	var value int64
	err = decoder.Decode(&value, reply, err)
	return value, err
}

/*
	ZREVRANGE key start stop [WITHSCORES]

	返回有序集 key 中, 指定区间内的成员.
	其中成员的位置按 score 值递增(从大到小)来排序.
	具有相同 score 值的成员按字典序的逆序(reverse lexicographical order)来排列.
	如果你需要成员按 score 值递减(从小到大)来排列, 请使用 ZRANGE 命令.

	下标参数 start 和 stop 都以 0 为底, 以 0 表示有序集第一个成员, 以 1 表示有序集第二个成员, 以此类推.
	也可以使用负数下标, 以 -1 表示最后一个成员, -2 表示倒数第二个成员, 以此类推.

	超出范围的下标并不会引起错误.
	当 start 的值比有序集的最大下标还要大, 或是 start > stop 时, ZRANGE 命令只是简单地返回一个空列表.
	假如 stop 参数的值比有序集的最大下标还要大, 那么 Redis 将 stop 当作最大下标来处理.
	可以通过使用 WITHSCORES 选项, 来让成员和它的 score 值一并返回, 返回列表以 value1,score1, ..., valueN,scoreN 的格式表示.

	可用版本: >= 1.2.0
	时间复杂度: O(log(N)+M), N 为有序集的基数, 而 M 为结果集的基数.

	返回值:
		指定区间内, 带有 score 值(可选)的有序集成员的列表.
*/
func (order *ZOrder) Range(value interface{}, start, stop int) error {
	reply, err := order.client.Exec(ZREVRANGE, order.key, start, stop, WITHSCORES)
	return decoder.Decode(value, reply, err)
}

/*
	ZSCORE key member

	返回有序集 key 中, 成员 member 的 score 值.
	如果 member 元素不是有序集 key 的成员, 或 key 不存在, 返回 nil .

	可用版本: >= 1.2.0
	时间复杂度: O(1)

	返回值:
		member 成员的 score 值, 以字符串形式表示.
*/
func (order *ZOrder) Score(member interface{}) (int64, bool, error) {

	reply, err := order.client.Exec(ZSCORE, order.key, member)
	if err == nil && reply == nil {
		return 0, false, nil
	}

	var score int64
	err = decoder.Int(&score, reply)
	return score, true, err
}

/*
	ZREVRANK key member

	返回有序集 key 中成员 member 的排名. 其中有序集成员按 score 值递减(从大到小)排序.
	排名以 0 开始, 也就是说, score 值最大的成员排名为 0.
	使用 ZRANK 命令可以获得成员按 score 值递增(从小到大)排列的排名.

	可用版本: >= 2.0.0
	时间复杂度: O(log(N))

	返回值:
		如果 member 是有序集 key 的成员, 返回 member 的排名.
		如果 member 不是有序集 key 的成员, 返回 nil.
*/
func (order *ZOrder) Rank(member interface{}) (int64, bool, error) {

	reply, err := order.client.Exec(ZREVRANK, order.key, member)
	if err == nil && reply == nil {
		return -1, false, nil
	}

	var rank int64
	err = decoder.Int(&rank, reply)
	return rank, true, err
}
