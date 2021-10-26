package framework

import (
	"errors"
	"fmt"
	"sync"
)

// Container 是一个服务容器，提供绑定服务和获取服务的功能
type Container interface {
	// Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换
	Bind(provider ServiceProvider) error

	// IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool

	// Make 根据关键字凭证获取一个服务
	Make(key string) (interface{}, error)

	// MustMake 根据关键字凭证获取一个服务，如果这个关键字凭证未绑定服务提供者，那么会panic。
	// 所以在使用这个接口的时候请保证服务容器已经为这个关键字凭证绑定了服务提供者。
	MustMake(key string) interface{}

	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式的
	// 它是根据服务提供者注册的启动函数和传递的params参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, param []interface{}) (interface{}, error)
}

// HadeContainer 服务容器的具体实现
type HadeContainer struct {
	// providers 存储注册的服务提供者，key为字符串凭证
	providers map[string]ServiceProvider
	// instance 存储具体的实例，key为字符串凭证
	instances map[string]interface{}
	// lock 用于锁住对容器的变更操作
	lock sync.RWMutex
	Container
}

// 创建一个服务容器
func NewHadeContainer() *HadeContainer {
	return &HadeContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

// PrintProviders 输出服务容器中的注册的关键字
func (hade *HadeContainer) PrintProviders() []string {
	var ret []string
	for _, provider := range hade.providers {
		name := provider.Name()
		line := fmt.Sprintf(name)
		ret = append(ret, line)
	}
	return ret
}

func (hade *HadeContainer) Bind(provider ServiceProvider) error {
	hade.lock.Lock()
	key := provider.Name()
	hade.providers[key] = provider
	hade.lock.Unlock()
	if provider.IsDefer() == false {
		if err := provider.Boot(hade); err != nil {
			return err
		}

		params := provider.Params(hade)
		method := provider.Register(hade)
		instance, err := method(params...)
		if err != nil {
			return errors.New(err.Error())
		}
		hade.instances[key] = instance
	}
	return nil
}

// IsBind 是否绑定服务
func (hade *HadeContainer) IsBind(key string) bool {
	return hade.findServiceProvider(key) != nil
}

func (hade *HadeContainer) findServiceProvider(key string) ServiceProvider {
	//hade.lock.RLock()
	//defer hade.lock.RUnlock()
	if sp, ok := hade.providers[key]; ok {
		return sp
	}
	return nil
}

func (hade *HadeContainer) newInstance(provider ServiceProvider, params []interface{}) (interface{}, error) {
	if err := provider.Boot(hade); err != nil {
		return nil, err
	}
	if params == nil {
		params = provider.Params(hade)
	}

	method := provider.Register(hade)
	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return ins, err
}

func (hade *HadeContainer) Make(key string) (interface{}, error) {
	return hade.make(key, nil, false)
}

func (hade *HadeContainer) MustMake(key string) interface{} {
	ins, err := hade.make(key, nil, false)
	if err != nil {
		panic(err)
	}
	return ins
}

func (hade *HadeContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return hade.make(key, params, true)
}

// 真正实例化一个服务
func (hade *HadeContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	hade.lock.RLock()
	defer hade.lock.RUnlock()
	sp := hade.findServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}
	if forceNew {
		return hade.newInstance(sp, params)
	}

	// 不需要强制重新实例化，如果容器中已经实例化了，那么直接返回即可
	if ins, ok := hade.instances[key]; ok {
		return ins, nil
	}

	// 容器中还未实例化，则进行一次实例化
	inst, err := hade.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}
	hade.instances[key] = inst
	return inst, nil
}
