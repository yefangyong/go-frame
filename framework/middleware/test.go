package middleware

import (
	"fmt"
	"go-frame/framework"
)

func Test1() framework.ControllerHandle {
	return func(c *framework.Context) error {
		fmt.Println("this is a test1 start")
		err := c.Next()
		if err != nil {
			return err
		}
		fmt.Println("this is a test1 end")
		return nil
	}
}

func Test2() framework.ControllerHandle {
	return func(c *framework.Context) error {
		fmt.Println("this is a test2 start")
		err := c.Next()
		if err != nil {
			return err
		}
		fmt.Println("this is a test2 end")
		return nil
	}
}
