package main

import (
	"context"
	"fmt"
	"go-frame/framework"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	core := framework.NewCore()
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
