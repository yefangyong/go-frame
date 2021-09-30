package middleware

import (
	"go-frame/framework"
	"log"
	"time"
)

func RecordRequestLog() framework.ControllerHandle {
	return func(c *framework.Context) error {
		start := time.Now()
		err := c.Next()
		if err != nil {
			return err
		}
		end := time.Now()
		cost := end.Sub(start)
		log.Printf("api request url:%s,time:%v\n", c.GetRequest().URL, cost.Microseconds())
		return nil
	}
}
