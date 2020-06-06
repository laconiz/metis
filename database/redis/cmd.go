package redis

const (
	OK   = "OK"
	PING = "PING"
	PONG = "PONG"

	// key
	DEL    = "DEL"
	EXISTS = "EXISTS"
	INCRBY = "INCRBY"

	// string
	SET = "SET"
	NX  = "NX"
	EX  = "EX"
	GET = "GET"

	// hash
	HDEL    = "HDEL"
	HSET    = "HSET"
	HGET    = "HGET"
	HMGET   = "HMGET"
	HGETALL = "HGETALL"
	HINCRBY = "HINCRBY"
	HEXISTS = "HEXISTS"

	// set
	SADD      = "SADD"
	SREM      = "SREM"
	SMEMBERS  = "SMEMBERS"
	SISMEMBER = "SISMEMBER"

	// sorted set
	ZINCRBY    = "ZINCRBY"
	ZREVRANGE  = "ZREVRANGE"
	ZSCORE     = "ZSCORE"
	ZREVRANK   = "ZREVRANK"
	WITHSCORES = "WITHSCORES"
)
