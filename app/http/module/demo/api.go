package demo

import (
	"github.com/yefangyong/go-frame/app/provider/demo"
	"github.com/yefangyong/go-frame/framework/gin"
)

type DemoApi struct {
	service *Service
}

func Register(r *gin.Engine) error {
	demoApi := NewDemoApi()
	r.Bind(&demo.DemoServiceProvider{})
	r.GET("demo", demoApi.Demo)
	return nil
}

// 初始化demoApi
func NewDemoApi() *DemoApi {
	service := NewService()
	return &DemoApi{service: service}
}

func (d *DemoApi) Demo(ctx *gin.Context) {
	user := d.service.getUser()
	UserDTO := UserModelsToUserDTOs(user)
	ctx.JSON(200, UserDTO)
}
