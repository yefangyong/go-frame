package demo

import (
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
	ctx.JSON(200, baseFolder)
}

func (d *DemoApi) Demo(ctx *gin.Context) {
	envService := ctx.MustMake(contract.EnvKey).(contract.Env)
	logService := ctx.MustMake(contract.LogKey).(contract.Log)
	logService.Info(ctx, "this is test", map[string]interface{}{})
	//user := d.service.getUser()
	//UserDTO := UserModelsToUserDTOs(user)
	ctx.JSON(200, envService.All())
}
