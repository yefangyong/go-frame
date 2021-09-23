package main

import (
	"context"
	"fmt"
	"go-frame/framework"
	"log"
	"time"
)

func FooControllerHandler(ctx *framework.Context) error {
	fmt.Println(1212)
	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)
	durationCtx, cancel := context.WithTimeout(ctx.BaseContext(), 5*time.Second)
	defer cancel()

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		// do real action
		time.Sleep(time.Second * 2)
		_ = ctx.Json(200, "ok")

		finish <- struct{}{}
	}()

	select {
	case p := <-panicChan:
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		log.Println(p)
		ctx.Json(500, "panic")
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		ctx.Json(500, "time out")
		ctx.SetHasTimeout()
	}
	return nil

}
