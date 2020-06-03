package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/laconiz/metis/database/redis/decoder"
	"github.com/laconiz/metis/log"
	"time"
)

func New(option *Option) (*Redis, error) {

	dial := func() (redis.Conn, error) {
		return redis.Dial(
			option.Network,
			option.Address,
			redis.DialPassword(option.Password),
			redis.DialDatabase(option.Database),
		)
	}

	redis := &Redis{
		pool: &redis.Pool{
			Dial:            dial,
			TestOnBorrow:    nil,
			MaxIdle:         option.MaxIdle,
			MaxActive:       option.MaxActive,
			IdleTimeout:     0,
			Wait:            true,
			MaxConnLifetime: 0,
		},
		option: option,
	}

	if _, err := redis.Exec(PING); err != nil {
		return nil, err
	}

	return redis, nil
}

type Redis struct {
	pool   *redis.Pool
	option *Option
}

func (redis *Redis) Exec(cmd string, args ...interface{}) (interface{}, error) {

	conn := redis.pool.Get()
	defer conn.Close()

	params, err := decoder.Params(args)
	if err != nil {
		return nil, err
	}

	reply, err := conn.Do(cmd, params...)

	if redis.option.Logger != nil {
		log := &ExecLog{Command: cmd, Request: params, Response: reply, Error: err}
		redis.option.Logger.Debug(log)
	}

	return reply, err
}

func (redis *Redis) Key() *Key {
	return &Key{client: redis}
}

func (redis *Redis) Hash(key string) *Hash {
	return &Hash{client: redis, key: key}
}

func (redis *Redis) ZOrder(key string) *ZOrder {
	return &ZOrder{client: redis, key: key}
}

func (redis *Redis) Set(key string) *Set {
	return &Set{client: redis, key: key}
}

func (redis *Redis) Eval(script *Script) *Eval {
	return &Eval{client: redis, script: script}
}

func (redis *Redis) Singleton(key string) *Singleton {
	return &Singleton{client: redis, key: key}
}

func (redis *Redis) Atomic(key string) *Atomic {
	return (&Atomic{redis: redis, key: key}).
		Expired(time.Second * 3).
		Timeout(time.Second * 6).
		Ticker(time.Millisecond * 50)
}

type Option struct {
	Network   string     // 网络类型
	Address   string     // 地址
	Password  string     // 密码
	Database  int        // 数据库
	MaxIdle   int        // 最大空闲连接数
	MaxActive int        // 最大活跃连接数
	Logger    log.Logger // 日志接口
}
