package main

import (
	"go-frame/framework"
	"net/http"
)

func main() {
	core := framework.NewCore()
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    "localhost:8082",
	}
	server.ListenAndServe()
}
