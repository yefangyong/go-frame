package demo

import (
	"fmt"

	"github.com/yefangyong/go-frame/app/provider/demo"
	"github.com/yefangyong/go-frame/framework/contract"
	"github.com/yefangyong/go-frame/framework/gin"
)

type DemoApi struct {
	service *Service
}

func Register(r *gin.Engine) error {
	demoApi := NewDemoApi()
	r.Bind(&demo.DemoServiceProvider{})
	r.GET("demo", demoApi.Demo)
	r.GET("demo2", demoApi.Demo2)
	r.GET("demo3", demoApi.Demo3)
	return nil
}

// 初始化demoApi
func NewDemoApi() *DemoApi {
	service := NewService()
	return &DemoApi{service: service}
}

func (d *DemoApi) Demo2(ctx *gin.Context) {
	service := ctx.MustMake(demo.DemoKey).(demo.Service)
	students := service.GetAllStudent()
	data := StudentsToUserDTOs(students)
	ctx.JSON(200, data)
}

func (d *DemoApi) Demo3(ctx *gin.Context) {
	app := ctx.MustMake(contract.AppKey).(contract.App)
	baseFolder := app.BaseFolder()
	fmt.Println("this is test")
	ctx.JSON(200, baseFolder)
}

// Demo godoc
// @Summary 获取所有用户
// @Description 获取所有用户
// @Produce  json
// @Tags demo
// @Success 200 array []UserDTO
// @Router /demo [get]
func (api *DemoApi) Demo(c *gin.Context) {
	c.JSON(200, "this is demo for dev all")
}
