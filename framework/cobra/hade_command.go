package cobra

import (
	"log"

	"github.com/robfig/cron/v3"
	"github.com/yefangyong/go-frame/framework"
)

// CronSpec 保存cron命令的信息，用于展示
type CronSpec struct {
	Type        string
	Cmd         *Command
	Spec        string
	ServiceName string
}

func (c *Command) SetParentNull() {
	c.parent = nil
}

func (c *Command) AddCronCommand(spec string, cmd *Command) {
	root := c.Root()
	if root.Cron == nil {
		root.Cron = cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))
		root.CronSpec = []CronSpec{}
	}
	// 增加说明信息
	root.CronSpec = append(root.CronSpec, CronSpec{
		Type: "normal-cron",
		Cmd:  cmd,
		Spec: spec,
	})

	// 制作一个rootCommand
	var cronCmd Command
	ctx := root.Context()
	cronCmd = *cmd
	cronCmd.args = []string{}
	cronCmd.SetParentNull()
	cronCmd.SetContainer(root.GetContainer())

	// 增加调用函数
	root.Cron.AddFunc(spec, func() {
		// 每个goroutine都是平等的，其中一个panic，其他的都会退出
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
		err := cronCmd.ExecuteContext(ctx)
		if err != nil {
			log.Println(err)
		}
	})

}

// 设置容器
func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

// 获取容器
func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}
