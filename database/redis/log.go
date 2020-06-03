package redis

import (
	"fmt"
	"github.com/laconiz/metis/database/redis/decoder"
)

type ExecLog struct {
	Command  string
	Request  []interface{}
	Response interface{}
	Error    error
}

func (log *ExecLog) String() string {

	head := append([]interface{}{log.Command}, log.Request...)
	tail := decoder.Reply(log.Response)
	str := fmt.Sprint(head) + " => " + fmt.Sprint(tail)

	if log.Error != nil && log.Error != log.Response {
		str += " => " + log.Error.Error()
	}

	return str
}

type EvalLog struct {
	Name     string
	Request  []interface{}
	Response interface{}
	Error    error
}

func (log *EvalLog) String() string {

	head := append([]interface{}{log.Name}, log.Request...)
	tail := decoder.Reply(log.Response)
	str := fmt.Sprint(head) + " => " + fmt.Sprint(tail)

	if log.Error != nil && log.Error != log.Response {
		str += " => " + log.Error.Error()
	}

	return str
}
