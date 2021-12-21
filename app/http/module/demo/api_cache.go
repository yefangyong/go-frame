package demo

import (
	"time"

	"github.com/yefangyong/go-frame/framework/contract"
	"github.com/yefangyong/go-frame/framework/gin"
	"github.com/yefangyong/go-frame/framework/provider/redis"
)

func (api *DemoApi) DemoCache(c *gin.Context) {
	logger := c.MustMake(contract.LogKey).(contract.Log)
	logger.Info(c, "request start", map[string]interface{}{})
	// 初始化cache服务
	cacheService := c.MustMake(contract.CacheKey).(contract.CacheService)
	// 设置key为foo
	err := cacheService.Set(c, "foo", "bar", 1*time.Hour)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	// 获取key为foo
	val, err := cacheService.Get(c, "foo")
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	logger.Info(c, "cache get", map[string]interface{}{
		"val": val,
	})
	// 删除key为foo
	//if err := cacheService.Del(c, "foo"); err != nil {
	//	c.AbortWithError(500, err)
	//	return
	//}
	c.JSON(200, "ok")
}

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
