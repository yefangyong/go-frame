package main

import (
	"go-frame/framework"
	"go-frame/framework/middleware"
)

func registerRouter(core *framework.Core) {
	core.Use(middleware.RecordRequestLog()) // 全局中间件
	// 静态路由+HTTP方法匹配
	core.Get("/user/login", middleware.Test2(), UserLoginController) // 路由中间件

	// 批量通用前缀
	subjectApi := core.Group("/subject")
	{
		subjectApi.Use(middleware.Test2()) // 路由组中间件
		// 动态路由
		subjectApi.Delete("/:id", SubjectDelController)
		subjectApi.Put("/:id", SubjectUpdateController)
		subjectApi.Get("/:id", SubjectGetController)
		subjectApi.Get("/list/all", SubjectListController)

		subjectInnerApi := subjectApi.Group("/info")
		{
			subjectInnerApi.Get("/name", SubjectNameController)
		}
	}
}
