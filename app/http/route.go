package http

import "github.com/yefangyong/go-frame/framework/gin"

func Routes(r *gin.Engine) {
	r.Static("/dist/", "./dist/")

}
