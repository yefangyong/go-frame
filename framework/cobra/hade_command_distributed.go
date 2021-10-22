package cobra

import (
	"log"
	"time"

	"github.com/yefangyong/go-frame/framework/contract"

	"github.com/robfig/cron/v3"
)

func (c *Command) AddDistributedCronCommand(serviceName string, spec string, cmd *Command, holdTime time.Duration) {
	root := c.Root()
	if root.Cron == nil {
		root.Cron = cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))
		root.CronSpec = []CronSpec{}
	}
	// 增加说明信息
	root.CronSpec = append(root.CronSpec, CronSpec{
		Type:        "normal-cron",
		Cmd:         cmd,
		Spec:        spec,
		ServiceName: serviceName,
	})

	container := root.GetContainer()
	appService := container.MustMake(contract.AppKey).(contract.App)
	distributedService := container.MustMake(contract.DistributedKey).(contract.Distributed)
	appID := appService.APPID()

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

		// 节点进行选举，返回选举结果
		selectAppID, err := distributedService.Select(serviceName, appID, holdTime)
		if err != nil {
			return
		}

		// 如果自己没有被选择到，直接返回
		if selectAppID != appID {
			return
		}

		err = cronCmd.ExecuteContext(ctx)
		if err != nil {
			log.Println(err)
		}
	})

}
