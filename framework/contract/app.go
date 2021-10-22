package contract

const AppKey = "hade:app"

type App interface {

	//AppId

	APPID() string
	// app 版本
	Version() string

	// BaseFolder 基础路径
	BaseFolder() string

	//ConfigFolder 定义配置路径
	ConfigFolder() string

	//LogFolder 定义了日志路径
	LogFolder() string

	// ProviderFolder 定义业务服务提供者的路径
	ProviderFolder() string

	// MiddlewareFolder 定义业务中间件的路径
	MiddlewareFolder() string

	// CommandFolder 定义业务定义的命令
	CommandFolder() string
	// RuntimeFolder 定义业务的运行中间态信息
	RuntimeFolder() string
	// TestFolder 存放测试所需要的信息
	TestFolder() string
}
