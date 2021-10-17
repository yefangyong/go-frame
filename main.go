package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/yefangyong/go-frame/provider/demo"

	"github.com/yefangyong/go-frame/framework/gin"
	"github.com/yefangyong/go-frame/framework/middleware"
)

func main() {
	core := gin.New()
	core.Bind(&demo.DemoServiceProvider{})
	core.Use(gin.Recovery())
	core.Use(middleware.RecordRequestLog())
	registerRouter(core)
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
