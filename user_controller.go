package main

import (
	"go-frame/framework"
	"time"
)

func UserLoginController(c *framework.Context) error {
	time.Sleep(10 * time.Second)
	c.Json(200, "ok, UserLoginController")
	return nil
}
