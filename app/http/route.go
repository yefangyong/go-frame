package http

import (
	"github.com/yefangyong/go-frame/app/http/module/demo"
	"github.com/yefangyong/go-frame/framework/gin"
	"github.com/yefangyong/go-frame/framework/middleware/static"
)

func Routes(r *gin.Engine) {
	r.Static("/dist/", "./dist/")
	// /路径先去./dist目录下查找文件是否存在，找到使用文件服务提供服务
	r.Use(static.Serve("/", static.LocalFile("./dist", false)))
	// demo相关路由注册
	demo.Register(r)
}
