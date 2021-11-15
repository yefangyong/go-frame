package command

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"syscall"

	"github.com/yefangyong/go-frame/framework/cobra"

	"github.com/yefangyong/go-frame/framework"
	"github.com/yefangyong/go-frame/framework/contract"
)

type devConfig struct {
	Port    string // 调试模式最终监听的端口
	Backend struct {
		RefreshTime   int    // 调试模式后端更新时间，如果文件变更，等待3s才进行一次更新，能让频繁保存变更更为顺畅, 默认1s
		Port          string // 后端监听端口
		MonitorFolder string // 监听文件夹，默认为AppFolder
	}
	Frontend struct { // 前端调试模式配置
		Port string // 前端启动端口，默认8071
	}
}

func initDevConfig(container framework.Container) *devConfig {
	devConfig := &devConfig{
		Port: "8087",
		Backend: struct {
			RefreshTime   int
			Port          string
			MonitorFolder string
		}{RefreshTime: 1, Port: "8072", MonitorFolder: ""},
		Frontend: struct{ Port string }{Port: "8071"},
	}

	configService := container.MustMake(contract.ConfigKey).(contract.Config)

	// 对每个配置项进行检查
	if configService.IsExist("app.dev.port") {
		devConfig.Port = configService.GetString("app.dev.port")
	}

	if configService.IsExist("app.dev.backend.port") {
		devConfig.Port = configService.GetString("app.dev.backend.port")
	}

	if configService.IsExist("app.dev.backend.refresh_time") {
		devConfig.Port = configService.GetString("app.dev.backend.refresh_time")
	}

	// monitorFolder 默认使用目录服务的 AppFolder()
	monitorFolder := configService.GetString("app.dev.backend.monitor_folder")
	if monitorFolder == "" {
		appService := container.MustMake(contract.AppKey).(contract.App)
		devConfig.Backend.MonitorFolder = appService.AppFolder()
	}

	if configService.IsExist("app.dev.fronted.port") {
		devConfig.Port = configService.GetString("app.dev.fronted.port")
	}

	return devConfig
}

type Proxy struct {
	devConfig   *devConfig // 配置文件
	backendPid  int
	frontendPid int
}

// 初始化一个Proxy
func NewProxy(container framework.Container) *Proxy {
	devConfig := initDevConfig(container)
	return &Proxy{
		devConfig: devConfig,
	}
}

// 重新启动一个 proxy 网关
func (p *Proxy) newProxyReverseProxy(fronted, backend *url.URL) *httputil.ReverseProxy {
	if p.frontendPid == 0 && p.backendPid == 0 {
		fmt.Println("前端和后端服务都不存在")
		return nil
	}

	//后端服务存在
	if p.frontendPid == 0 && p.backendPid != 0 {
		return httputil.NewSingleHostReverseProxy(backend)
	}

	//前端服务存在
	if p.frontendPid != 0 && p.backendPid == 0 {
		return httputil.NewSingleHostReverseProxy(fronted)
	}

	// 两个都有进程
	// 先创建一个后端服务的directory
	directory := func(req *http.Request) {
		if req.URL.Path == "/" || req.URL.Path == "/app.js" {
			req.URL.Scheme = fronted.Scheme
			req.URL.Host = fronted.Host
		} else {
			req.URL.Scheme = backend.Scheme
			req.URL.Host = backend.Host
		}
	}

	NotFoundErr := errors.New("response is 404, need to redirect")
	return &httputil.ReverseProxy{
		Director: directory, // 先转发到后端
		ModifyResponse: func(response *http.Response) error {
			if response.StatusCode == 404 {
				return NotFoundErr
			}
			return nil
		},
		ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
			// 判断 Error 是否为NotFoundError, 是的话则进行前端服务的转发，重新修改writer
			if errors.Is(err, NotFoundErr) {
				httputil.NewSingleHostReverseProxy(fronted).ServeHTTP(writer, request)
			}
		},
	}
}

// 启动proxy服务，并且根据参数启动前端服务或者后端服务
func (p *Proxy) startProxy(startFronted, startBackend bool) error {
	var backendURL, frontendURL *url.URL
	var err error

	// 启动后端
	if startBackend {
		if err := p.restartBackend(); err != nil {
			return err
		}
	}

	// 启动前端
	if startFronted {
		if err := p.restartFrontend(); err != nil {
			return err
		}
	}

	if frontendURL, err = url.Parse(fmt.Sprintf("%s%s", "http://127.0.0.1:", p.devConfig.Frontend.Port)); err != nil {
		return err
	}

	if backendURL, err = url.Parse(fmt.Sprintf("%s%s", "http://127.0.0.1:", p.devConfig.Backend.Port)); err != nil {
		return err
	}

	// 设置反向代理
	proxyReverse := p.newProxyReverseProxy(frontendURL, backendURL)
	proxyServer := &http.Server{
		Addr:    "127.0.0.1:" + p.devConfig.Port,
		Handler: proxyReverse,
	}
	fmt.Println("代理服务启动：", "http://"+proxyServer.Addr)

	// 启动proxy服务
	err = proxyServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
	return nil

}

func (p *Proxy) restartBackend() error {
	// 先杀死旧的进程
	if p.backendPid != 0 {
		syscall.Kill(p.backendPid, syscall.SIGKILL)
		p.backendPid = 0
	}

	// 设置随机端口，真实后端的端口
	port := p.devConfig.Backend.Port
	hadeAddress := fmt.Sprintf(":" + port)

	// 使用命令行启动后端进程
	cmd := exec.Command("./hade", "app", "start", "--address="+hadeAddress)
	cmd.Stdout = os.NewFile(0, os.DevNull)
	cmd.Stderr = os.Stderr
	fmt.Println("启动后端服务：", "http://127.0.0.1:"+port)
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	p.backendPid = cmd.Process.Pid
	fmt.Println("后端服务Pid:", p.backendPid)
	return nil
}

// 重启启动前端服务
func (p *Proxy) restartFrontend() error {
	// 启动前端调试模式
	// 先杀死旧的进程
	if p.frontendPid != 0 {
		syscall.Kill(p.frontendPid, syscall.SIGKILL)
		p.frontendPid = 0
	}

	// 否则开启 npm run serve
	port := p.devConfig.Frontend.Port
	path, err := exec.LookPath("npm")
	if err != nil {
		return err
	}
	cmd := exec.Command(path, "run", "dev")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s%s", "PORT=", port))
	cmd.Stdout = os.NewFile(0, os.DevNull)
	cmd.Stderr = os.Stderr

	// 因为npm run serve 是控制台挂起模式，所以这里使用go routine启动
	err = cmd.Start()
	fmt.Println("启动前端服务：", "http://127.0.0.1:"+port)
	if err != nil {
		fmt.Println(err)
	}
	p.frontendPid = cmd.Process.Pid
	fmt.Println("前端服务Pid:", p.frontendPid)
	return nil
}

var devCommand = &cobra.Command{
	Use:   "dev",
	Short: "调试模式",
	RunE: func(c *cobra.Command, args []string) error {
		c.Help()
		return nil
	},
}

// 启动后端调试模式
var devBackendCommand = &cobra.Command{
	Use:   "backend",
	Short: "启动后端调试模式",
	RunE: func(c *cobra.Command, args []string) error {
		proxy := NewProxy(c.GetContainer())
		if err := proxy.startProxy(false, true); err != nil {
			return err
		}
		return nil
	},
}

// devFrontendCommand 启动前端调试模式
var devFrontendCommand = &cobra.Command{
	Use:   "frontend",
	Short: "前端调试模式",
	RunE: func(c *cobra.Command, args []string) error {
		// 启动前端服务
		proxy := NewProxy(c.GetContainer())
		return proxy.startProxy(true, false)
	},
}

var devAllCommand = &cobra.Command{
	Use:   "all",
	Short: "同时启动前端和后端进行调试",
	RunE: func(c *cobra.Command, args []string) error {
		proxy := NewProxy(c.GetContainer())
		if err := proxy.startProxy(true, true); err != nil {
			return err
		}
		return nil
	},
}

// 初始化Dev命令
func initDevCommand() *cobra.Command {
	devCommand.AddCommand(devBackendCommand)
	devCommand.AddCommand(devFrontendCommand)
	devCommand.AddCommand(devAllCommand)
	return devCommand
}