package framework

type IGroup interface {
	// 实现HTTPMethod方法
	Get(string, ControllerHandle)
	Delete(string, ControllerHandle)
	Put(string, ControllerHandle)
	Post(string, ControllerHandle)

	// 实现嵌套 group
	Group(string) IGroup
}

// Group struct 实现了IGroup
type Group struct {
	core   *Core  // 指向 core 结构
	parent *Group // 指向上一个Group,如果有的话
	prefix string // 这个group的通用前缀
}

// 初始化Group
func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		parent: nil,
		prefix: prefix,
	}
}

// 实现 Get 方法
func (group *Group) Get(url string, handler ControllerHandle) {
	url = group.getAbsolutePrefix() + url
	group.core.Get(url, handler)
}

// 实现 Put 方法
func (group *Group) Put(url string, handler ControllerHandle) {
	url = group.getAbsolutePrefix() + url
	group.core.Put(url, handler)
}

// 实现 Post 方法
func (group *Group) Post(url string, handler ControllerHandle) {
	url = group.getAbsolutePrefix() + url
	group.core.Post(url, handler)
}

// 实现 Delete 方法
func (group *Group) Delete(url string, handler ControllerHandle) {
	url = group.getAbsolutePrefix() + url
	group.core.Delete(url, handler)
}

// 递归获取当前group的绝对路径
func (group *Group) getAbsolutePrefix() string {
	if group.parent == nil {
		return group.prefix
	}
	return group.parent.getAbsolutePrefix() + group.prefix
}

// 实现 Group 方法
func (group *Group) Group(url string) IGroup {
	cgroup := NewGroup(group.core, url)
	cgroup.parent = group
	return cgroup
}
