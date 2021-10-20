package app

import (
	"github.com/yefangyong/go-frame/framework"
	"github.com/yefangyong/go-frame/framework/contract"
)

// HadeAppProvider 提供 App 的具体实现方法
type HadeAppProvider struct {
	BaseFolder string
}

func (h *HadeAppProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, h.BaseFolder}
}

func (h *HadeAppProvider) Register(container framework.Container) framework.NewInstance {
	return NewHadeApp
}

func (h *HadeAppProvider) Boot(container framework.Container) error {
	return nil
}

func (h *HadeAppProvider) IsDefer() bool {
	return false
}

func (h *HadeAppProvider) Name() string {
	return contract.AppKey
}
