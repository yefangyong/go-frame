package contract

import "time"

const DistributedKey = "hade:distributed"

type Distributed interface {
	// Select 分布式选择器，所有节点对同一个服务进行抢占，只返回其中一个节点
	// serviceName 服务名称
	// appId 当前的appId
	// holdTime 占用的时间
	// 返回值：
	// selectAppID 节点ID
	// err 错误值
	Select(serviceName string, appId string, holdTime time.Duration) (selectAppID string, err error)
}
