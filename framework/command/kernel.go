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
}
