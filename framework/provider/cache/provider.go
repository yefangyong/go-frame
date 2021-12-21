package cache

import (
	"strings"

	"github.com/yefangyong/go-frame/framework"
	"github.com/yefangyong/go-frame/framework/contract"
	"github.com/yefangyong/go-frame/framework/provider/cache/service"
)

type HadeCacheProvider struct {
	Driver string // Driver
}

// 根据不同的驱动，使用不同的缓存方式
func (h *HadeCacheProvider) Register(container framework.Container) framework.NewInstance {
	if h.Driver == "" {
		tcs, err := container.Make(contract.ConfigKey)
		if err != nil {
			return service.NewMemoryCache
		}
		configService := tcs.(contract.Config)
		h.Driver = strings.ToLower(configService.GetString("cache.driver"))
	}
	switch h.Driver {
	case "redis":
		return service.NewRedisCache
	case "memory":
		return service.NewMemoryCache
	default:
		return service.NewMemoryCache
	}
}

func (h *HadeCacheProvider) Boot(container framework.Container) error {
	return nil
}

func (h *HadeCacheProvider) IsDefer() bool {
	return true
}

func (h *HadeCacheProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

func (h *HadeCacheProvider) Name() string {
	return contract.CacheKey
}
