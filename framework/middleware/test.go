package middleware

import (
	"fmt"

	"github.com/yefangyong/go-frame/framework/gin"
)

func Test1() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("this is a test1 start")
		c.Next()
		fmt.Println("this is a test1 end")
	}
}

func Test2() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("this is a test2 start")
		c.Next()
		fmt.Println("this is a test2 end")
	}
}
