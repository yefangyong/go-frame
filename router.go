package main

import (
	"github.com/yefangyong/go-frame/framework/gin"
	"github.com/yefangyong/go-frame/framework/middleware"
)

func registerRouter(core *gin.Engine) {
	core.Use(middleware.RecordRequestLog()) // 全局中间件
	// 静态路由+HTTP方法匹配
	core.GET("/user/login", middleware.Test2(), UserLoginController) // 路由中间件

	// 批量通用前缀
	subjectApi := core.Group("/subject")
	{
		subjectApi.Use(middleware.Test2()) // 路由组中间件
		// 动态路由
		subjectApi.DELETE("/:id", SubjectDelController)
		subjectApi.PUT("/:id", SubjectUpdateController)
		subjectApi.GET("/:id", SubjectGetController)
		subjectApi.GET("/list/all", SubjectListController)

		subjectInnerApi := subjectApi.Group("/info")
		{
			subjectInnerApi.GET("/name", SubjectNameController)
		}
	}
}
