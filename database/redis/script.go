package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/laconiz/metis/database/redis/decoder"
)

type Eval struct {
	script *Script
	client *Client
}

func (eval *Eval) Load() error {

	conn := eval.client.pool.Get()
	defer conn.Close()

	return eval.script.Script.Load(conn)
}

func (eval *Eval) Do(args ...interface{}) (interface{}, error) {

	conn := eval.client.pool.Get()
	defer conn.Close()

	params, err := decoder.Params(args)
	if err != nil {
		return nil, err
	}

	reply, err := eval.script.Script.Do(conn, params...)

	if eval.client.option.Logger != nil {
		log := &EvalLog{Name: eval.script.Name, Request: args, Response: reply, Error: err}
		eval.client.option.Logger.Debug(log)
	}

	return reply, err
}

func (eval *Eval) Exec(value interface{}, args ...interface{}) error {
	reply, err := eval.Do(args...)
	return decoder.Decode(value, reply, err)
}

func NewScript(name string, keys int, src string) *Script {
	return &Script{Name: name, Script: redis.NewScript(keys, src)}
}

type Script struct {
	Name   string
	Script *redis.Script
}
