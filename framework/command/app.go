package command

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/sevlyar/go-daemon"

	"github.com/erikdubbelboer/gspt"

	"github.com/yefangyong/go-frame/framework/util"

	"github.com/yefangyong/go-frame/framework"

	"github.com/yefangyong/go-frame/framework/cobra"
	"github.com/yefangyong/go-frame/framework/contract"
)

// app启动地址
var appAddress = ""
var appDaemon = false

// appStartCommand 启动一个 Web 服务
var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动一个web服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 从Command中获取服务实例
		container := cmd.GetContainer()

		// 从服务容器中获取kernel的服务实例
		kernelService := container.MustMake(contract.KernelKey).(contract.Kernel)

		// 从kernel服务实例中获取引擎
		core := kernelService.HttpEngine()
		if appAddress == "" {
			envService := container.MustMake(contract.EnvKey).(contract.Env)
			if envService.Get("ADDRESS") != "" {
				appAddress = envService.Get("ADDRESS")
			} else {
				configService := container.MustMake(contract.ConfigKey).(contract.Config)
				if configService.IsExist("app.address") {
					appAddress = configService.GetString("app.address")
				} else {
					appAddress = ":8888"
				}
			}
		}

		// 创建一个Server服务
		server := &http.Server{
			Handler: core,
			Addr:    appAddress,
		}

		appService := container.MustMake(contract.AppKey).(contract.App)
		pidFolder := appService.RuntimeFolder()

		if !util.Exists(pidFolder) {
			if err := os.MkdirAll(pidFolder, os.ModePerm); err != nil {
				return err
			}
		}

		serverPidFile := filepath.Join(pidFolder, "app.pid")
		logFolder := appService.LogFolder()
		if !util.Exists(logFolder) {
			if err := os.MkdirAll(logFolder, os.ModePerm); err != nil {
				return err
			}
		}

		//// 应用日志文件
		serverLogFile := filepath.Join(logFolder, "app.log")
		currentFolder := util.GetExecDirectory()

		// daemon 模式
		if appDaemon {
			fmt.Println(serverPidFile)
			// 创建一个context
			cntxt := &daemon.Context{
				PidFileName: serverPidFile,
				PidFilePerm: 0644,
				LogFileName: serverLogFile,
				LogFilePerm: 0640,
				WorkDir:     currentFolder,
				Umask:       027,
				Args:        []string{"", "app", "start", "--daemon=true"},
			}
			// 启动子进程，d不为空表示当前是父进程，d为空表示当前是子进程
			d, err := cntxt.Reborn()
			if err != nil {
				return err
			}
			if d != nil {
				// 父进程直接打印启动成功信息，不做任何操作
				fmt.Println("app启动成功，pid:", d.Pid)
				fmt.Println("日志文件:", serverLogFile)
				return nil
			}
			defer cntxt.Release()
			// 子进程执行真正的app启动操作
			fmt.Println("daemon started")
			gspt.SetProcTitle("hade app")
			if err := startAppServer(server, container); err != nil {
				fmt.Println(err)
			}
			return nil
		}

		// 非 daemon 模式
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID]", content)
		err := ioutil.WriteFile(serverPidFile, []byte(content), 0644)
		if err != nil {
			return err
		}

		gspt.SetProcTitle("hade app")
		fmt.Println("app serve url:", appAddress)
		if err := startAppServer(server, container); err != nil {
			fmt.Println(err)
		}
		return nil
	},
}

// 初始化 app 并且启动服务
func initAppCommand() *cobra.Command {
	appStartCommand.Flags().BoolVarP(&appDaemon, "daemon", "d", false, "start app daemon")
	appStartCommand.Flags().StringVar(&appAddress, "address", ":8888", "设置app启动的地址，默认为:8888")
	appCommand.AddCommand(appStartCommand)
	appCommand.AddCommand(appStateCommand)
	appCommand.AddCommand(appStopCommand)
	appCommand.AddCommand(appRestartCommand)
	return appCommand
}

// AppCommand 是命令行参数第一级为app的命令，它没有实际功能，只是打印帮助文档
var appCommand = &cobra.Command{
	Use:   "app",
	Short: "业务应用控制命令",
	Long:  "业务应用控制命令，其包含业务启动，关闭，重启，查询等功能",
	RunE: func(c *cobra.Command, args []string) error {
		// 打印帮助文档
		c.Help()
		return nil
	},
}

var appRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "重启启动一个app服务",
	RunE: func(c *cobra.Command, args []string) error {
		err := appStopCommand.RunE(c, args)
		if err != nil {
			return err
		}
		fmt.Println("开启启动app服务")
		appDaemon = true
		err = appStartCommand.RunE(c, args)
		if err != nil {
			return err
		}
		fmt.Println("重启app服务成功")
		return nil
	},
}

// 启动AppServer, 这个函数会将当前goroutine阻塞
func startAppServer(server *http.Server, c framework.Container) error {
	go func() {
		server.ListenAndServe()
	}()

	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	closeWait := 5
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	if configService.IsExist("app.close_wait") {
		closeWait = configService.GetInt("app.close_wait")
	}
	// 调用Server.Shutdown graceful结束
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(closeWait)*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		return err
	}
	return nil
}

// 获取启动的app的pid
var appStateCommand = &cobra.Command{
	Use:   "state",
	Short: "获取启动的app的pid",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// 获取pid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")
		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				fmt.Println("app服务已经启动，pid:", pid)
				return nil
			}
		}
		fmt.Println("没有app服务存在")
		return nil
	},
}

// 停止一个已经启动的app服务
var appStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "停止一个已经启动的app服务",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// Get pid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}
		if content != nil && len(content) != 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			err = StopServe(pid, container)
			if err != nil {
				return err
			}
			if err := ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
				return err
			}
			fmt.Println("结束进程成功:" + strconv.Itoa(pid))
		}
		return nil
	},
}

// StopServe 发送信号，结束进程
func StopServe(pid int, container framework.Container) error {
	// 发送Sigterm命令
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return err
	}

	// 检查pid是否存在
	closeWait := 5
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	if configService.IsExist("app.close_wait") {
		closeWait = configService.GetInt("app.close_wait")
	}

	for i := 0; i < closeWait*2; i++ {
		if util.CheckProcessExist(pid) == false {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// 如果进程等待了2*closeWait之后还没有结束，返回错误，不进程后续的操作
	if util.CheckProcessExist(pid) == true {
		fmt.Println("结束进程失败："+strconv.Itoa(pid), "请查看原因")
		return errors.New("结束进程失败")
	}
	return nil
}
