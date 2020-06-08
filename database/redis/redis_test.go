package redis

import (
	"fmt"
	"github.com/laconiz/metis/log/logz"
)

func assert(assertions ...bool) {
	for i, a := range assertions {
		if !a {
			panic(fmt.Errorf("assert index %d failed", i))
		}
	}
}

const (
	KeyA    = "TestKeyA"
	MemberA = "TestMemberA"
	MemberB = "TestMemberB"
	ValueA  = "TestValueA"
	ValueB  = "TestValueB"
)

var client *Client

type Struct struct {
	A int64
	B string
	C []float32
}

func (s Struct) Equal(v Struct) bool {
	return s.A == v.A && s.B == v.B
}

var ComplexValue = Struct{A: 100, B: "complex value", C: []float32{1.1, 2.2, 3.3}}
var ComplexPointer = &Struct{A: 200, B: "complex pointer", C: []float32{4.4, 5.5}}

func init() {

	conf := Option{
		Network:   "tcp",
		Address:   "192.168.3.45:6379",
		Database:  15,
		MaxIdle:   5,
		MaxActive: 50,
		Logger:    logz.Global(),
	}

	var err error
	if client, err = New(conf); err != nil {
		panic(err)
	}
}
