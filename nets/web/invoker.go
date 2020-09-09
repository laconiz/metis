package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laconiz/metis/log"
	"github.com/laconiz/metis/nets/session"
	"github.com/laconiz/metis/utils/ioc"
	"net/http"
	"reflect"
	"time"
)

type Node struct {
	Path    string      // 路径
	Method  string      // 方法
	Handler interface{} // 接口
	Log     bool        // 日志
}

func NewInvoker(logger log.Logger) *Invoker {
	return &Invoker{squirt: ioc.NewSquirt().Params(&gin.Context{}), logger: logger}
}

type Invoker struct {
	squirt *ioc.Squirt
	logger log.Logger
}

func (invoker *Invoker) Arguments(args ...interface{}) *Invoker {
	invoker.squirt.Arguments(args...)
	return invoker
}

func (invoker *Invoker) Builders(funcs ...interface{}) *Invoker {
	invoker.squirt.Builders(funcs...)
	return invoker
}

func (invoker *Invoker) Register(router gin.IRouter, nodes []*Node) error {

	for _, node := range nodes {
		if err := invoker.RegisterNode(router, node); err != nil {
			return err
		}
	}

	return nil
}

func (invoker *Invoker) RegisterNode(router gin.IRouter, node *Node) error {

	arguments, err := invoker.squirt.Unknown(node.Handler)
	if err != nil {
		return err
	}

	if reflect.TypeOf(node.Handler).NumOut() != 1 {
		return errors.New("invalid handler num out")
	}

	switch len(arguments) {
	case 0:
	case 1:
		if arguments[0].Kind() != reflect.Ptr {
			return fmt.Errorf("invalid message %v", arguments[0])
		}
		invoker.squirt.Builder(arguments[0], invoker.bind(arguments[0]))
	default:
		return errors.New("invalid handler num in")
	}

	handler, err := invoker.squirt.Handle(node.Handler)
	if err != nil {
		return err
	}

	invoker.logger.Data(node.Path, node.Method).Info("registered")
	router.Handle(node.Method, node.Path, invoker.Handle(node, handler))
	return nil
}

func (invoker *Invoker) Handle(node *Node, handler ioc.Invoker) gin.HandlerFunc {

	logger := invoker.logger.Field("path", node.Path).Field("method", node.Method)

	return func(ctx *gin.Context) {

		now := time.Now()
		in, out, err := handler(ctx)
		entry := logger.Data("session", session.NewID()).Data("duration", time.Since(now).String())

		if err != nil {
			entry.Data("error", err).Info("invoke error")
			return
		}

		var requests []interface{}
		for _, value := range in {
			if value.CanInterface() {
				requests = append(requests, value.Interface())
			}
		}

		var responses []interface{}
		for _, value := range out {
			if value.CanInterface() {
				responses = append(responses, value.Interface())
			}
		}

		entry = entry.Data("request", requests)
		ctx.JSON(http.StatusOK, out[0].Interface())
		entry.Data("response", responses).Info("execute success")
	}
}

func (invoker *Invoker) bind(typo reflect.Type) func(ctx *gin.Context) (interface{}, error) {
	return func(ctx *gin.Context) (interface{}, error) {
		msg := reflect.New(typo.Elem()).Interface()
		err := ctx.Bind(msg)
		return msg, err
	}
}
