package framework

import (
	"fmt"
	"net/http"
)

type Core struct {
	Router map[string]ControllerHandle
}

func NewCore() *Core {
	return &Core{
		Router: map[string]ControllerHandle{},
	}
}

func (c *Core) Get(url string, handler ControllerHandle) {
	c.Router[url] = handler
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request, response)

	// 一个简单的路由选择器
	router := c.Router["foo"]
	if router == nil {
		return
	}
	err := router(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
