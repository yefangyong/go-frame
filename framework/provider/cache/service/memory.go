package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/yefangyong/go-frame/framework/contract"

	"github.com/yefangyong/go-frame/framework"
)

type MemoryCache struct {
	datas     map[string]*MemoryData
	container framework.Container
	lock      *sync.RWMutex
}

func (m *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	var val string
	if err := m.GetObj(ctx, key, &val); err != nil {
		return "", err
	}
	return val, nil
}

func (m *MemoryCache) GetObj(ctx context.Context, key string, model interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if md, ok := m.datas[key]; ok {
		if md.ttl != NoneDuration {
			if time.Now().Sub(md.createTime) > md.ttl {
				delete(m.datas, key)
				return ErrKeyNotFound
			}
		}

		bs, _ := json.Marshal(md.val)
		err := json.Unmarshal(bs, model)
		if err != nil {
			return err
		}
		return nil
	}
	return ErrKeyNotFound
}

func (m *MemoryCache) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	errs := make([]string, 0, len(keys))
	rets := make(map[string]string)
	for _, key := range keys {
		val, err := m.Get(ctx, key)
		if err != nil {
			errs = append(errs, err.Error())
		}
		rets[key] = val
	}
	if len(errs) == 0 {
		return rets, nil
	}
	return rets, errors.New(strings.Join(errs, "||"))
}

func (m *MemoryCache) Set(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	return m.SetObj(ctx, key, val, timeout)
}

func (m *MemoryCache) SetObj(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	md := &MemoryData{
		val:        val,
		createTime: time.Now(),
		ttl:        timeout,
	}
	m.datas[key] = md
	return nil
}

func (m *MemoryCache) SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error {
	var errs []string
	for k, v := range data {
		err := m.Set(ctx, k, v, timeout)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "||"))
	}
	return nil
}

func (m *MemoryCache) SetForever(ctx context.Context, key string, val string) error {
	return m.Set(ctx, key, val, NoneDuration)
}

func (m *MemoryCache) SetForeverObj(ctx context.Context, key string, val interface{}) error {
	return m.Set(ctx, key, val, NoneDuration)
}

func (m *MemoryCache) SetTTL(ctx context.Context, key string, timeout time.Duration) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if md, ok := m.datas[key]; ok {
		md.ttl = timeout
		return nil
	}
	return ErrKeyNotFound
}

func (m *MemoryCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if md, ok := m.datas[key]; ok {
		return md.ttl, nil
	}
	return 0, ErrKeyNotFound
}

// cache-aside方式
func (m *MemoryCache) Remember(ctx context.Context, key string, timeout time.Duration, remember contract.RememberFunc, model interface{}) error {
	err := m.GetObj(ctx, key, model)
	if err == nil {
		return nil
	}

	if !errors.Is(err, ErrKeyNotFound) {
		return err
	}

	// key not found
	objNew, err := remember(ctx, m.container)
	if err != nil {
		return err
	}

	if err := m.SetObj(ctx, key, objNew, timeout); err != nil {
		return err
	}

	if err := m.GetObj(ctx, key, &model); err != nil {
		return err
	}

	return nil
}

func (m *MemoryCache) Calc(ctx context.Context, key string, step int64) (int64, error) {
	var val int64
	err := m.GetObj(ctx, key, val)
	val = val + step
	if err == nil {
		m.datas[key].val = val
		return val, nil
	}

	if !errors.Is(err, ErrKeyNotFound) {
		return 0, err
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	m.datas[key] = &MemoryData{
		val:        val,
		createTime: time.Now(),
		ttl:        NoneDuration,
	}
	return val, nil
}

func (m *MemoryCache) Increment(ctx context.Context, key string) (int64, error) {
	return m.Calc(ctx, key, 1)
}

func (m *MemoryCache) Decrement(ctx context.Context, key string) (int64, error) {
	return m.Calc(ctx, key, -1)
}

func (m *MemoryCache) Del(ctx context.Context, key string) error {
	m.lock.Lock()
	m.lock.Unlock()
	delete(m.datas, key)
	return nil
}

func (m *MemoryCache) DelMany(ctx context.Context, keys []string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, key := range keys {
		delete(m.datas, key)
	}
	return nil
}

type MemoryData struct {
	val        interface{}
	createTime time.Time
	ttl        time.Duration
}

// 初始化内存缓存
func NewMemoryCache(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	lock := &sync.RWMutex{}
	datas := make(map[string]*MemoryData)
	return &MemoryCache{
		datas:     datas,
		lock:      lock,
		container: container,
	}, nil
}
