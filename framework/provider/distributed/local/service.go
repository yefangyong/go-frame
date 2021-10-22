package local

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/yefangyong/go-frame/framework/contract"

	"github.com/yefangyong/go-frame/framework"
)

type DistributedService struct {
	container framework.Container
}

func NewDistributedService(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("params error")
	}
	container := params[0].(framework.Container)
	return &DistributedService{container: container}, nil
}

func (d *DistributedService) Select(serviceName string, appId string, hold time.Duration) (selectID string, error error) {
	appService := d.container.MustMake(contract.AppKey).(contract.App)
	runtimeFolder := appService.RuntimeFolder()
	lockFile := filepath.Join(runtimeFolder, "distributed_"+serviceName)

	// 打开文件锁
	lock, err := os.OpenFile(lockFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	//尝试独占文件锁
	err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	// 抢不到文件锁
	if err != nil {
		// 读取被选择的appId
		selectIDByte, err := ioutil.ReadAll(lock)
		if err != nil {
			return "", err
		}
		return string(selectIDByte), nil
	}

	//开启一个携程，定时释放锁
	go func() {
		defer func() {
			// 释放文件锁
			_ = syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
			// 释放文件
			_ = lock.Close()
			// 删除文件
			_ = os.Remove(lockFile)
		}()
		//创建选举结果有效的计时器
		timer := time.NewTimer(hold)
		<-timer.C
	}()

	// 已经抢到了文件锁，把AppId写入文件中
	if _, err := lock.WriteString(appId); err != nil {
		return "", err
	}

	return appId, nil
}
