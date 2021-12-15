package main

import (
	"github.com/yefangyong/go-frame/app/console"
	"github.com/yefangyong/go-frame/app/http"
	"github.com/yefangyong/go-frame/framework"
	"github.com/yefangyong/go-frame/framework/provider/app"
	"github.com/yefangyong/go-frame/framework/provider/config"
	"github.com/yefangyong/go-frame/framework/provider/distributed/local"
	"github.com/yefangyong/go-frame/framework/provider/env"
	"github.com/yefangyong/go-frame/framework/provider/kernel"
	"github.com/yefangyong/go-frame/framework/provider/log"
	"github.com/yefangyong/go-frame/framework/provider/log/formatter"
	"github.com/yefangyong/go-frame/framework/provider/orm"
)

func main() {
	container := framework.NewHadeContainer()
	container.Bind(&app.HadeAppProvider{})
	container.Bind(&env.HadeEnvProvider{})
	container.Bind(&local.DistributedProvider{})
	container.Bind(&config.HadeConfigProvider{})
	container.Bind(&orm.GormProvider{})
	container.Bind(&log.HadeLogServiceProvider{
		Driver:    "single",
		Formatter: formatter.JsonFormatter,
	})
	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.HadeKernelProvider{
			HttpEngine: engine,
		})
	}
	// 运行root命令
	console.RunCommand(container)
}
