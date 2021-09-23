package main

import "go-frame/framework"

func registerRouter(core *framework.Core) {
	core.Get("foo", FooControllerHandler)
}
