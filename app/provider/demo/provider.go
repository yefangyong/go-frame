package demo

import (
	"fmt"

	"github.com/yefangyong/go-frame/framework"
)

type DemoServiceProvider struct {
}

func (d *DemoServiceProvider) Register(container framework.Container) framework.NewInstance {
	return NewDemoService
}

func (d *DemoServiceProvider) Boot(container framework.Container) error {
	fmt.Println("demo service boot")
	return nil
}

func (d *DemoServiceProvider) IsDefer() bool {
	return true
}

func (d *DemoServiceProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

func (d *DemoServiceProvider) Name() string {
	return DemoKey
}
