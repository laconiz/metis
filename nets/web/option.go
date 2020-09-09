package web

import "github.com/laconiz/metis/log"

// ---------------------------------------------------------------------------------------------------------------------

const module = "httpis"

// ---------------------------------------------------------------------------------------------------------------------

type AcceptorOption struct {
	Name     string        // 名称
	Addr     string        // 侦听地址
	Nodes    []*Node       // 接口
	Params   []interface{} // 注入参数
	Creators []interface{} // 参数生成器
	Level    log.Level     // 日志等级
}

func (option *AcceptorOption) parse() {

	if option.Name == "" {
		option.Name = "acceptor"
	}

	if option.Addr == "" {
		option.Addr = "0.0.0.0:8080"
	}

	if !option.Level.Valid() {
		option.Level = log.INFO
	}
}
