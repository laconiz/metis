package ioc

import (
	"fmt"
	"reflect"
)

type Builder struct {
	Type reflect.Type
	Func interface{}
}

type Builders []*Builder

func (builders Builders) Distinct() Builders {

	var ret Builders
	hits := map[reflect.Type]bool{}

	for _, builder := range builders {

		if _, ok := hits[builder.Type]; ok {
			continue
		}

		hits[builder.Type] = true
		ret = append(ret, builder)
	}

	return ret
}

func NewBuilder(function interface{}) (*Builder, error) {

	const errFormatter = "invalid builder: need func(...any) (any, error), got %v"

	typo := reflect.TypeOf(function)
	if typo == nil || typo.Kind() != reflect.Func {
		return nil, fmt.Errorf(errFormatter, typo)
	}

	errType := reflect.TypeOf((*error)(nil)).Elem()
	if typo.NumOut() != 2 || !typo.Out(1).Implements(errType) {
		return nil, fmt.Errorf(errFormatter, errType)
	}

	return &Builder{Type: typo, Func: function}, nil
}
