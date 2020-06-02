package context

import (
	"fmt"
	"github.com/laconiz/metis/utils/json"
)

type Fields map[string]interface{}

func NewContext() *Context {
	return &Context{fields: Fields{}}
}

type Context struct {
	raw    []byte
	fields Fields
}

func (context *Context) Fields(fields Fields) *Context {

	copy := NewContext()

	for key, value := range context.fields {
		copy.fields[key] = value
	}

	for key, value := range fields {
		copy.fields[key] = value
	}

	return copy
}

func (context *Context) Raw() []byte {

	if len(context.fields) == 0 {
		return nil
	}

	if context.raw != nil {
		return context.raw
	}

	raw, err := json.Marshal(context.fields)
	if err == nil {
		context.raw = raw
	} else {
		context.raw = []byte(fmt.Sprint(context.fields))
	}

	return context.raw
}
