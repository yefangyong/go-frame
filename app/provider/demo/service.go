package demo

import (
	"github.com/yefangyong/go-frame/framework"
)

type DemoService struct {
	container framework.Container
}

func NewDemoService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	return &DemoService{container: container}, nil
}

func (d *DemoService) GetAllStudent() []Student {
	return []Student{
		{ID: 1, Name: "yfy323"}, {ID: 2, Name: "jsz12"},
	}
}
