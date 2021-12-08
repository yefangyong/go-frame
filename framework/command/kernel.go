package command

import "github.com/yefangyong/go-frame/framework/cobra"

func AddKernelCommands(root *cobra.Command) {

	//绑定定时任务相关命令
	root.AddCommand(InitCronCommand())

	// 绑定 App 相关命令
	root.AddCommand(initAppCommand())

	// 绑定 build 相关命令
	root.AddCommand(initBuildCommand())

	// 绑定 Dev 调试相关命令
	root.AddCommand(initDevCommand())

	// 绑定自定义生成服务提供者的命令
	root.AddCommand(initProviderCommand())

	// 绑定中间件 middleware 相关命令
	root.AddCommand(initMiddlewareCommand())

	// 绑定 cmd 相关命令
	root.AddCommand(initCmdCommand())

	// 绑定 config 相关命令
	root.AddCommand(initConfigCommand())

	// 绑定 env 相关命令
	root.AddCommand(initEnvCommand())

	// 绑定 swagger 相关命令
	root.AddCommand(initSwaggerCommand())
}
