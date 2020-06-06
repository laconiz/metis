package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/laconiz/metis/database/redis/decoder"
	"github.com/laconiz/metis/log"
	"time"
)

func New(option *Option) (*Client, error) {

	dial := func() (redis.Conn, error) {
		return redis.Dial(
			option.Network,
			option.Address,
			redis.DialPassword(option.Password),
			redis.DialDatabase(option.Database),
		)
	}

	client := &Client{
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

	if _, err := client.Exec(PING); err != nil {
		return nil, err
	}

	return client, nil
}

type Client struct {
	pool   *redis.Pool
	option *Option
}

func (client *Client) Exec(cmd string, args ...interface{}) (interface{}, error) {

	conn := client.pool.Get()
	defer conn.Close()

	params, err := decoder.Params(args)
	if err != nil {
		return nil, err
	}

	reply, err := conn.Do(cmd, params...)

	if client.option.Logger != nil {
		log := &ExecLog{Command: cmd, Request: params, Response: reply, Error: err}
		client.option.Logger.Debug(log)
	}

	return reply, err
}

func (client *Client) Key() *Key {
	return &Key{client: client}
}

func (client *Client) Hash(key string) *Hash {
	return &Hash{client: client, key: key}
}

func (client *Client) ZOrder(key string) *ZOrder {
	return &ZOrder{client: client, key: key}
}

func (client *Client) Set(key string) *Set {
	return &Set{client: client, key: key}
}

func (client *Client) Eval(script *Script) *Eval {
	return &Eval{client: client, script: script}
}

func (client *Client) Singleton(key string) *Singleton {
	return &Singleton{client: client, key: key}
}

func (client *Client) Atomic(key string) *Atomic {
	return (&Atomic{client: client, key: key}).
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
