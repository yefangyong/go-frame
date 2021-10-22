package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/yefangyong/go-frame/framework/util"

	"github.com/erikdubbelboer/gspt"

	"github.com/yefangyong/go-frame/framework/cobra"
	"github.com/yefangyong/go-frame/framework/contract"
)

// 初始化定时任务命令
func InitCronCommand() *cobra.Command {
	// 启动
	cronCommand.AddCommand(cronStartCommand)
	// 查看定时任务列表
	cronCommand.AddCommand(cronListCommand)
	// 重启
	cronCommand.AddCommand(cronRestartCommand)
	// 停止
	cronCommand.AddCommand(cronStopCommand)
	// 状态
	cronCommand.AddCommand(cronStateCommand)
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

var cronListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出所有的定时任务",
	RunE: func(command *cobra.Command, args []string) error {
		cronSpecs := command.Root().CronSpec
		ps := [][]string{}
		for _, cronSpec := range cronSpecs {
			line := []string{
				cronSpec.Type, cronSpec.Spec, cronSpec.Cmd.Use, cronSpec.Cmd.Short, cronSpec.ServiceName,
			}
			ps = append(ps, line)
		}
		util.PrettyPrint(ps)
		return nil
	},
}

var cronStateCommand = &cobra.Command{
	Use:   "state",
	Short: "cron常驻进程状态",
	RunE: func(command *cobra.Command, args []string) error {
		container := command.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		serverPidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return nil
		}

		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				fmt.Println("cron server started pid:", pid)
				return nil
			}

		}
		fmt.Println("no cron server start")
		return nil
	},
}

var cronRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "重启cron的常驻进程",
	RunE: func(command *cobra.Command, args []string) error {
		container := command.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		serverPidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return nil
		}

		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}

			if util.CheckProcessExist(pid) {
				if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
					return err
				}
			}
			for i := 0; i < 10; i++ {
				if util.CheckProcessExist(pid) == false {
					break
				}
				time.Sleep(1 * time.Second)
			}
			fmt.Println("kill process:" + strconv.Itoa(pid))
		}
		cronStartCommand.RunE(command, args)
		return nil
	},
}

var cronStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "停止cron的常驻进程",
	RunE: func(command *cobra.Command, args []string) error {
		container := command.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		serverPidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return nil
		}

		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}

			if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
				return err
			}

			if err := ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
				return err
			}
			fmt.Println("stop pid:", pid)
		}
		return nil
	},
}
