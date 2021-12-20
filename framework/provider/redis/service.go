package redis

import (
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/yefangyong/go-frame/framework"
	"github.com/yefangyong/go-frame/framework/contract"
)

// 设置 redis 服务结构体
type HadeRedisService struct {
	container framework.Container      //服务容器
	clients   map[string]*redis.Client // key 为UniqKey,value为redis.Client
	lock      *sync.RWMutex            // 读写锁
}

// 初始化 redis 服务
func NewHadeRedisService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	clients := make(map[string]*redis.Client)
	lock := &sync.RWMutex{}
	return &HadeRedisService{
		container: container,
		clients:   clients,
		lock:      lock,
	}, nil
}

func (app *HadeRedisService) GetClient(options ...contract.RedisOption) (*redis.Client, error) {
	// 读取默认配置
	config := GetBaseConfig(app.container)

	// 修改默认配置
	for _, opt := range options {
		err := opt(app.container, config)
		if err != nil {
			return nil, err
		}
	}

	key := config.UniqKey()

	// 如果已经实例化，则直接返回
	app.lock.Lock()
	if db, ok := app.clients[key]; ok {
		app.lock.Unlock()
		return db, nil
	}
	app.lock.Unlock()

	// 实例化redis
	app.lock.Lock()
	defer app.lock.Unlock()

	client := redis.NewClient(config.Options)
	app.clients[key] = client
	return client, nil
}
