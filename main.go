package main

import (
	"github.com/yefangyong/go-frame/app/console"
	"github.com/yefangyong/go-frame/app/http"
	"github.com/yefangyong/go-frame/framework"
	"github.com/yefangyong/go-frame/framework/provider/app"
	"github.com/yefangyong/go-frame/framework/provider/distributed/local"
	"github.com/yefangyong/go-frame/framework/provider/env"
	"github.com/yefangyong/go-frame/framework/provider/kernel"
)

func main() {
	container := framework.NewHadeContainer()
	container.Bind(&app.HadeAppProvider{})
	container.Bind(&env.HadeEnvProvider{})
	container.Bind(&local.DistributedProvider{})
	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.HadeKernelProvider{
			HttpEngine: engine,
		})
	}
	// 运行root命令
	console.RunCommand(container)
}
