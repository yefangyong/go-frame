package demo

import (
	"github.com/yefangyong/go-frame/framework/contract"
	"github.com/yefangyong/go-frame/framework/gin"
	"github.com/yefangyong/go-frame/framework/provider/redis"
)

func (api *DemoApi) DemoRedis(c *gin.Context) {
	logger := c.MustMake(contract.LogKey).(contract.Log)
	logger.Info(c, "request start", map[string]interface{}{})
	redisService := c.MustMake(contract.RedisKey).(contract.RedisService)
	redisClient, err := redisService.GetClient(redis.WithConfigPath("cache.redis"), redis.WithRedisConfig(func(config *contract.RedisConfig) {
		config.MaxRetries = 3
	}))
	if err != nil {
		logger.Error(c, err.Error(), map[string]interface{}{})
		c.AbortWithError(500, err)
		return
	}
	key := "test"
	err = redisClient.Set(c, key, "1234", 0).Err()
	if err != nil {
		logger.Error(c, err.Error(), map[string]interface{}{})
		c.AbortWithError(500, err)
		return
	}
	logger.Info(c, "设置缓存成功", map[string]interface{}{})
	err = redisClient.Get(c, key).Err()
	if err != nil {
		logger.Error(c, err.Error(), map[string]interface{}{})
		c.AbortWithError(500, err)
		return
	}
	logger.Info(c, "获取缓存成功", map[string]interface{}{})
	c.JSON(200, "2323")
}
