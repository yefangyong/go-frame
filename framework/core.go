package framework

import (
	"log"
	"net/http"
	"strings"
)

type Core struct {
	router map[string]*Tree
}

func NewCore() *Core {
	// 初始化路由
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	return &Core{
		router: router,
	}
}

// 匹配 GET 方法, 增加路由规则
func (c *Core) Get(url string, handler ControllerHandle) {
	if err := c.router["GET"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 匹配 Post 方法, 增加路由规则
func (c *Core) Post(url string, handler ControllerHandle) {
	if err := c.router["POST"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 匹配 Put 方法, 增加路由规则
func (c *Core) Put(url string, handler ControllerHandle) {
	if err := c.router["PUT"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 匹配 Delete 方法, 增加路由规则
func (c *Core) Delete(url string, handler ControllerHandle) {
	if err := c.router["DELETE"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 匹配路由，如果没有匹配到，则返回 nil
func (c *Core) FindRouteByRequest(request *http.Request) ControllerHandle {
	// url 和 method 全部转换为大写，保证大小写不敏感
	method := strings.ToUpper(request.Method)
	url := strings.ToUpper(request.URL.Path)
	if methodHandlers, ok := c.router[method]; ok {
		return methodHandlers.FindHandler(url)
	}
	return nil
}

// 分组路由
func (c *Core) Group(prefix string) *Group {
	return NewGroup(c, prefix)
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	// 封装自定义 context
	ctx := NewContext(request, response)

	// 寻找路由
	router := c.FindRouteByRequest(request)
	if router == nil {
		// 如果没有找到，这里打印日志
		ctx.Json(404, "not found")
		return
	}

	// 调用路由函数，如果返回 err 代表存在内部错误，返回 500 的状态码
	if err := router(ctx); err != nil {
		ctx.Json(500, "inner error")
		return
	}
}
