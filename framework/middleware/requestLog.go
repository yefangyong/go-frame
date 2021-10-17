package middleware

import (
	"log"
	"time"

	"github.com/yefangyong/go-frame/framework/gin"
)

func RecordRequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		cost := end.Sub(start)
		log.Printf("api request url:%s,time:%v\n", c.Request.URL.Path, cost.Microseconds())
	}
}
