package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/erikdubbelboer/gspt"

	"github.com/yefangyong/go-frame/framework/cobra"
	"github.com/yefangyong/go-frame/framework/contract"
)

// 初始化定时任务命令
func InitCronCommand() *cobra.Command {
	cronCommand.AddCommand(cronStartCommand)
	return cronCommand
}

var cronCommand = &cobra.Command{
	Use:   "cron",
	Short: "定时任务的相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

var cronStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动cron常驻进程",
	RunE: func(c *cobra.Command, args []string) error {
		// 获取容器
		container := c.GetContainer()

		// 获取容器中的App服务
		appService := container.MustMake(contract.AppKey).(contract.App)

		// 设置cron的日志地址和进程ID地址
		pidFolder := appService.RuntimeFolder()
		serverPidFile := filepath.Join(pidFolder, "cron.pid")
		//logFolder := appService.LogFolder()
		//serverLogFile := filepath.Join(logFolder, "cron.log")
		//currentFolder := appService.BaseFolder()

		// no deamon mode
		fmt.Println("start cron job")
		pid := strconv.Itoa(os.Getpid())
		fmt.Println("【pid】", pid)
		err := ioutil.WriteFile(serverPidFile, []byte(pid), 0644)
		if err != nil {
			return err
		}
		gspt.SetProcTitle("hade cron")
		c.Root().Cron.Run()
		return nil
	},
}
