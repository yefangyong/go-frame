package local

import (
	"github.com/yefangyong/go-frame/framework"
	"github.com/yefangyong/go-frame/framework/contract"
)

type DistributedProvider struct {
}

func (d *DistributedProvider) Register(container framework.Container) framework.NewInstance {
	return NewDistributedService
}

func (d *DistributedProvider) Boot(container framework.Container) error {
	return nil
}

func (d *DistributedProvider) IsDefer() bool {
	return false
}

func (d *DistributedProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

func (d *DistributedProvider) Name() string {
	return contract.DistributedKey
}
