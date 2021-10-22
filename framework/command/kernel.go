package command

import "github.com/yefangyong/go-frame/framework/cobra"

func AddKernelCommands(root *cobra.Command) {

	//绑定定时任务相关命令
	root.AddCommand(InitCronCommand())

	// 绑定App相关命令
	root.AddCommand(initAppCommand())
}
