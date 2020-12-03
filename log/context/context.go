package context

import (
	"fmt"
	"github.com/laconiz/metis/utils/json"
)

type Fields map[string]interface{}

func New() *Context {
	return &Context{fields: Fields{}}
}

type Context struct {
	raw    []byte
	fields Fields
}

func (context *Context) Fields(fields Fields) *Context {
	// 拷贝旧数据
	copy := New()
	for key, value := range context.fields {
		copy.fields[key] = value
	}
	// 写入新数据
	for key, value := range fields {
		if err, ok := value.(error); ok {
			copy.fields[key] = err.Error()
		} else {
			copy.fields[key] = value
		}
	}
	return copy
}

func (context *Context) Raw() []byte {
	// 空接口
	if len(context.fields) == 0 {
		return nil
	}
	// 已序列化数据
	if context.raw != nil {
		return context.raw
	}
	// 序列化数据
	raw, err := json.Marshal(context.fields)
	if err == nil {
		context.raw = raw
	} else {
		context.raw = []byte(fmt.Sprint(context.fields, err))
	}
	return context.raw
}
