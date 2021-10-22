package kernel

import (
	"net/http"

	"github.com/yefangyong/go-frame/framework/gin"
)

type HadeKernelService struct {
	engine *gin.Engine
}

// 初始化web引擎服务实例
func NewHadeKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.Engine)
	return &HadeKernelService{engine: httpEngine}, nil
}

//返回引擎
func (s *HadeKernelService) HttpEngine() http.Handler {
	return s.engine
}
