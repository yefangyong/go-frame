package http

import (
	"github.com/yefangyong/go-frame/app/http/module/demo"
	"github.com/yefangyong/go-frame/framework/gin"
	ginSwagger "github.com/yefangyong/go-frame/framework/middleware/gin-swagger"
	"github.com/yefangyong/go-frame/framework/middleware/gin-swagger/swaggerFiles"
	"github.com/yefangyong/go-frame/framework/middleware/static"
)

func Routes(r *gin.Engine) {
	//container := r.GetContainer()
	//configService := container.MustMake(contract.ConfigKey).(contract.Config)
	r.Static("/dist/", "./dist/")
	// /路径先去./dist目录下查找文件是否存在，找到使用文件服务提供服务
	r.Use(static.Serve("/", static.LocalFile("./dist", false)))

	//if configService.GetBool("app.swagger") == true {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//}

	// demo相关路由注册
	demo.Register(r)
}
