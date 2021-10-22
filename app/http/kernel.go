package http

import "github.com/yefangyong/go-frame/framework/gin"

func NewHttpEngine() (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	Routes(r)
	return r, nil
}
