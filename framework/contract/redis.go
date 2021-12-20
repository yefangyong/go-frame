package contract

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/yefangyong/go-frame/framework"
)

const RedisKey = "hade:redis"

// RedisOption 代表初始化的时候的选项
type RedisOption func(container framework.Container, config *RedisConfig) error

type RedisService interface {
	GetClient(options ...RedisOption) (*redis.Client, error)
}

type RedisConfig struct {
	*redis.Options
}

// UniqKey 用来唯一标识 redisConfig 的值
func (config *RedisConfig) UniqKey() string {
	return fmt.Sprintf("%v_%v_%v_%v", config.Addr, config.DB, config.Username, config.Network)
}
