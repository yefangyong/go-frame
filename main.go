package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	httpHade "github.com/yefangyong/go-frame/app/http"

	demo2 "github.com/yefangyong/go-frame/app/provider/demo"

	"github.com/yefangyong/go-frame/framework/gin"
	"github.com/yefangyong/go-frame/framework/middleware"
	"github.com/yefangyong/go-frame/framework/provider/app"
)

func main() {
	core := gin.New()
	// 绑定具体的服务
	core.Bind(&app.HadeAppProvider{})
	core.Bind(&demo2.DemoServiceProvider{})
	httpHade.Routes(core)
	core.Use(gin.Recovery())
	core.Use(middleware.RecordRequestLog())
	server := &http.Server{
		Handler: core,
		Addr:    ":8082",
	}
	go func() {
		server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, // kill -SIGINT XXXX 或 Ctrl+c
		os.Interrupt,
		syscall.SIGINT, // register that too, it should be ok
		// os.Kill等同于syscall.Kill
		os.Kill,
		syscall.SIGKILL, // register that too, it should be ok
		// kill -SIGTERM XXXX
		syscall.SIGTERM)
	select {
	case <-quit:
		err := server.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("shutdown error:%v\n", err)
		}
		fmt.Printf("shutdown ok")
	}

}
