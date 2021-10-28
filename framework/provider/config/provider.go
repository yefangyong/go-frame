package config

import (
	"path/filepath"

	"github.com/yefangyong/go-frame/framework"
	"github.com/yefangyong/go-frame/framework/contract"
)

type HadeConfigProvider struct {
}

func (h *HadeConfigProvider) Register(container framework.Container) framework.NewInstance {
	return NewHadeConfig
}

func (h *HadeConfigProvider) Boot(container framework.Container) error {
	return nil
}

func (h *HadeConfigProvider) IsDefer() bool {
	return false
}

func (h *HadeConfigProvider) Params(container framework.Container) []interface{} {
	appService := container.MustMake(contract.AppKey).(contract.App)
	envService := container.MustMake(contract.EnvKey).(contract.Env)
	env := envService.AppEnv()
	configFolder := appService.ConfigFolder()
	envFolder := filepath.Join(configFolder, env)
	return []interface{}{container, envFolder, envService.All()}
}

func (h *HadeConfigProvider) Name() string {
	return contract.ConfigKey
}
