package http

import (
	"github.com/yefangyong/go-frame/app/http/module/demo"
	"github.com/yefangyong/go-frame/framework/gin"
)

func Routes(r *gin.Engine) {
	r.Static("/dist/", "./dist/")
	// demo相关路由注册
	demo.Register(r)
}
