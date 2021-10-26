package env

import (
	"github.com/yefangyong/go-frame/framework"
	"github.com/yefangyong/go-frame/framework/contract"
)

type HadeEnvProvider struct {
	Folder string
}

func (e *HadeEnvProvider) Register(container framework.Container) framework.NewInstance {
	return NewHadeEnv
}

func (e *HadeEnvProvider) Boot(container framework.Container) error {
	appService := container.MustMake(contract.AppKey).(contract.App)
	e.Folder = appService.BaseFolder()
	return nil
}

func (e *HadeEnvProvider) IsDefer() bool {
	return false
}

func (e *HadeEnvProvider) Params(container framework.Container) []interface{} {
	return []interface{}{e.Folder}
}

func (e *HadeEnvProvider) Name() string {
	return contract.EnvKey
}
