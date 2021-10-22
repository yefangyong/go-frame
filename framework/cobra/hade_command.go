package cobra

import "github.com/yefangyong/go-frame/framework"

// 设置容器
func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

// 获取容器
func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}
