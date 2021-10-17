package demo

import (
	"fmt"

	"github.com/yefangyong/go-frame/framework"
)

type DemoService struct {
	Service

	c framework.Container
}

func NewDemoService(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)

	fmt.Println("new demo service")

	return &DemoService{c: c}, nil
}

func (d *DemoService) GetFoo() Foo {
	return Foo{
		Name: "this is a demo",
	}
}
