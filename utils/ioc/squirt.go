package ioc

import (
	"fmt"
	"github.com/codegangsta/inject"
	"reflect"
)

type Invoker func(...interface{}) ([]reflect.Value, []reflect.Value, error)

func NewSquirt() *Squirt {

	squirt := &Squirt{
		arguments: inject.New(),
		params:    inject.New(),
		builders:  map[reflect.Type]*Builder{},
	}
	squirt.params.SetParent(squirt.arguments)

	return squirt
}

type Squirt struct {
	arguments inject.Injector
	params    inject.Injector
	builders  map[reflect.Type]*Builder
	err       error
}

func (squirt *Squirt) execute(operate func() error) *Squirt {
	if squirt.err == nil {
		squirt.err = operate()
	}
	return squirt
}

// 是否实现指定参数类型
func (squirt *Squirt) implement(typo reflect.Type) bool {
	return squirt.params.Get(typo).IsValid() || squirt.builders[typo] != nil
}

// 写入默认参数
func (squirt *Squirt) Arguments(arguments ...interface{}) *Squirt {

	return squirt.execute(func() error {

		for _, argument := range arguments {

			typo := reflect.TypeOf(argument)
			if typo == nil {
				return fmt.Errorf("invalid argument %v", typo)
			}

			if squirt.implement(typo) {
				return fmt.Errorf("duplicate argument %v", typo)
			}

			squirt.arguments.Map(argument)
		}

		return nil
	})
}

// 写入参数生成函数
func (squirt *Squirt) Builders(functions ...interface{}) *Squirt {

	return squirt.execute(func() error {

		for _, function := range functions {

			builder, err := NewBuilder(function)
			if err != nil {
				return err
			}

			if squirt.implement(builder.Type) {
				return fmt.Errorf("duplicate builder %v", builder.Type)
			}

			squirt.builders[builder.Type] = builder
		}

		return nil
	})
}

// 写入临时参数
func (squirt *Squirt) Params(params ...interface{}) *Squirt {

	return squirt.execute(func() error {

		squirt.params = inject.New()
		squirt.params.SetParent(squirt.arguments)

		for _, param := range params {
			squirt.params.Map(param)
		}

		return nil
	})
}

// 写入参数生成函数
func (squirt *Squirt) Builder(typo reflect.Type, function interface{}) *Squirt {

	return squirt.execute(func() error {

		if squirt.implement(typo) {
			return fmt.Errorf("duplicate builder %v", typo)
		}

		builder, err := NewBuilder(function)
		if err != nil {
			return err
		}

		builder.Type = typo
		squirt.builders[typo] = builder

		return nil
	})
}

// 检查接口类型
func (squirt *Squirt) handler(handler interface{}) (reflect.Type, error) {

	typo := reflect.TypeOf(handler)
	if typo == nil || typo.Kind() != reflect.Func {
		return nil, fmt.Errorf("invalid handler %v", typo)
	}

	return typo, nil
}

func (squirt *Squirt) Handle(handler interface{}) (Invoker, error) {

	if squirt.err != nil {
		return nil, squirt.err
	}

	typo, err := squirt.handler(handler)
	if err != nil {
		return nil, err
	}

	builders := Builders{}
	for i := 0; i < typo.NumIn(); i++ {

		if squirt.params.Get(typo.In(i)).IsValid() {
			continue
		}

		builder, ok := squirt.builders[typo.In(i)]
		if !ok {
			return nil, fmt.Errorf("invalid argument %v", typo.In(i))
		}

		chain, err := squirt.chain(builder, 0)
		if err != nil {
			return nil, err
		}

		builders = append(builders, chain...)
	}

	return squirt.invoker(builders.Distinct(), handler), nil
}

// 组建参数生成链
func (squirt *Squirt) chain(builder *Builder, depth int) (Builders, error) {

	if depth >= 20 {
		return nil, fmt.Errorf("too deep argument %v", builder.Type)
	}

	builders := Builders{}
	typo := reflect.TypeOf(builder.Func)

	for i := 0; i < typo.NumIn(); i++ {

		if squirt.params.Get(typo.In(i)).IsValid() {
			continue
		}

		sub, ok := squirt.builders[typo.In(i)]
		if !ok {
			return nil, fmt.Errorf("invalid argument %v", typo.In(i))
		}

		chain, err := squirt.chain(sub, depth+1)
		if err != nil {
			return nil, err
		}

		builders = append(builders, chain...)
	}

	return append(builders, builder), nil
}

func (squirt *Squirt) invoker(builders Builders, handler interface{}) Invoker {

	typo := reflect.TypeOf(handler)
	value := reflect.ValueOf(handler)

	return func(arguments ...interface{}) ([]reflect.Value, []reflect.Value, error) {

		injector := inject.New()
		injector.SetParent(squirt.arguments)
		for _, argument := range arguments {
			injector.Map(argument)
		}

		for _, builder := range builders {

			values, err := injector.Invoke(builder.Func)
			if err != nil {
				return nil, nil, err
			}

			if !values[1].IsNil() {
				return nil, nil, values[1].Interface().(error)
			}

			if values[0].CanInterface() {
				injector.Map(values[0].Interface())
			}
		}

		in := make([]reflect.Value, typo.NumIn())

		for i := 0; i < typo.NumIn(); i++ {

			if value := injector.Get(typo.In(i)); value.IsValid() {
				in[i] = value
				continue
			}

			return in, nil, fmt.Errorf("value not found by type %v", typo.In(i))
		}

		out := value.Call(in)

		return in, out, nil
	}
}

// 获取未知参数
func (squirt *Squirt) Unknown(handler interface{}) ([]reflect.Type, error) {

	typo, err := squirt.handler(handler)
	if err != nil {
		return nil, err
	}

	var types []reflect.Type
	for i := 0; i < typo.NumIn(); i++ {
		if !squirt.implement(typo.In(i)) {
			types = append(types, typo.In(i))
		}
	}

	return types, nil
}
