package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yefangyong/go-frame/framework"
	"github.com/yefangyong/go-frame/framework/contract"
)

// 获取 redis 基本的配置信息
func GetBaseConfig(container framework.Container) *contract.RedisConfig {
	logService := container.MustMake(contract.LogKey).(contract.Log)
	config := &contract.RedisConfig{Options: &redis.Options{}}
	opt := WithConfigPath("redis")
	err := opt(container, config)
	if err != nil {
		logService.Error(context.Background(), "parse cache config error", map[string]interface{}{})
		return nil
	}
	return config
}

// 表示自行配置 redis 的配置信息
func WithRedisConfig(f func(config *contract.RedisConfig)) contract.RedisOption {
	return func(container framework.Container, config *contract.RedisConfig) error {
		f(config)
		return nil
	}
}

// 加载配置地址
func WithConfigPath(configPath string) contract.RedisOption {
	return func(container framework.Container, config *contract.RedisConfig) error {
		configService := container.MustMake(contract.ConfigKey).(contract.Config)
		conf := configService.GetStringMapString(configPath)
		// 读取config配置
		/*
		   driver: redis # 连接驱动
		   host: localhost # ip地址
		   port: 3306 # 端口
		   db: 0 #db
		   username: jianfengye # 用户名
		   password: "123456789" # 密码
		   timeout: 10s # 连接超时
		   read_timeout: 2s # 读超时
		   write_timeout: 2s # 写超时
		   conn_min_idle: 10 # 连接池最小空闲连接数
		   conn_max_open: 20 # 连接池最大连接数
		   conn_max_lifetime: 1h # 连接数最大生命周期
		   conn_max_idletime: 1h # 连接数空闲时长
		*/
		if host, ok := conf["host"]; ok {
			if port, ok1 := conf["port"]; ok1 {
				config.Addr = host + ":" + port
			}
		}

		if db, ok := conf["db"]; ok {
			t, err := strconv.Atoi(db)
			if err != nil {
				return err
			}
			config.DB = t
		}

		if username, ok := conf["username"]; ok {
			config.Username = username
		}

		if password, ok := conf["password"]; ok {
			config.Password = password
		}

		if timeout, ok := conf["timeout"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.DialTimeout = t
		}

		if timeout, ok := conf["read_timeout"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.ReadTimeout = t
		}

		if timeout, ok := conf["write_timeout"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.WriteTimeout = t
		}

		if cnt, ok := conf["conn_min_idle"]; ok {
			t, err := strconv.Atoi(cnt)
			if err != nil {
				return err
			}
			config.MinIdleConns = t
		}

		if max, ok := conf["conn_max_open"]; ok {
			t, err := strconv.Atoi(max)
			if err != nil {
				return err
			}
			config.PoolSize = t
		}

		if timeout, ok := conf["conn_max_lifetime"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.MaxConnAge = t
		}

		if timeout, ok := conf["conn_max_idletime"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.IdleTimeout = t
		}

		return nil
	}
}
